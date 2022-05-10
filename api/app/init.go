package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/access"
	"github.com/go-ozzo/ozzo-routing/v2/fault"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/sirupsen/logrus"
)

// InitLog Init returns a middleware that prepares the request context and processing environment.
// The middleware will populate RequestContext, handle possible panics and errors from the processing
// handlers, and add an access log entry.
func InitLog(logger *logrus.Logger) routing.Handler {
	return func(rc *routing.Context) error {
		now := time.Now()

		rc.Response = &access.LogResponseWriter{rc.Response, http.StatusOK, 0}

		ac := newRequestScope(now, logger, rc.Request)
		rc.Set("Context", ac)

		fault.Recovery(ac.Errorf, convertError)(rc)
		logAccess(rc, ac.Infof, ac.Now())

		return nil
	}
}

// InitLogger Initialize Logger
func InitLogger(logger *logrus.Logger) {
	// The API for setting attributes is a little different than the package level
	// exported Logger. See Godoc.
	logger.Out = os.Stdout
	logger.Formatter = &logrus.JSONFormatter{}

	if os.Getenv("GO_ENV") == "production" {
		logger.SetLevel(logrus.WarnLevel)
		// You could set this to any `io.Writer` such as a file
		t := time.Now()
		name := "data_" + t.Format("2006-01-02")
		file, err := os.OpenFile("./logs/"+name+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.Out = file
		} else {
			logger.Infof("Failed to log to file, using default stderr %v", err)
		}
	}
}

// GetRequestScope returns the RequestScope of the current request.
func GetRequestScope(c *routing.Context) RequestScope {
	return c.Get("Context").(RequestScope)
}

// logAccess logs a message describing the current request.
func logAccess(c *routing.Context, logFunc access.LogFunc, start time.Time) {
	rw := c.Response.(*access.LogResponseWriter)
	elapsed := float64(time.Now().Sub(start).Nanoseconds()) / 1e6
	requestLine := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)
	logFunc(`[%.3fms] %s %d %d`, elapsed, requestLine, rw.Status, rw.BytesWritten)
}

// convertError converts an error into an APIError so that it can be properly sent to the response.
// You may need to customize this method by adding conversion logic for more error types.
func convertError(c *routing.Context, err error) error {
	if err == sql.ErrNoRows {
		return errors.NotFound("the requested resource was not found")
	}
	if err == errors.New("sql: no rows in result set") {
		return errors.NoContentFound("No result found")
	}
	if err.Error() == "sql: no rows in result set" {
		return errors.NoContentFound("No result found")
	}
	switch err.(type) {
	case *errors.APIError:
		return err
	case validation.Errors:
		return errors.InvalidData(err.(validation.Errors))
	case routing.HTTPError:
		switch err.(routing.HTTPError).StatusCode() {
		case http.StatusUnauthorized:
			return errors.Unauthorized(err.Error())
		case http.StatusNotFound:
			return errors.NotFound("the requested resource")
		}
	}
	return errors.InternalServerError(err.Error())
}
