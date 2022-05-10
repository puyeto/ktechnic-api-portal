package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	"github.com/jinzhu/gorm"
	"github.com/ktechnics/ktechnics-api/api/app"
	"github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

// Server ...
type Server struct {
	DB *gorm.DB
}

// Initialize ...
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string, logger *logrus.Logger) {

	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			logger.Fatalf("This is the error:", err)
		} else {
			logger.Printf("We are connected to the %v database\n", Dbdriver)
		}
	}

	// server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) //database migration

	router := routing.New()

	router.Use(
		// app.Init(logger),
		content.TypeNegotiator(content.JSON),
		// app.Transactional(testdata.DB),
	)

	// server.initializeRoutes()
}

func (server *Server) buildRouter(logger *logrus.Logger) *routing.Router {
	router := routing.New()

	router.To("GET,HEAD", "/ping", func(c *routing.Context) error {
		c.Abort() // skip all other middlewares/handlers
		return c.Write("OK ")
	})

	router.Use(
		app.InitLog(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	rg := router.Group("/api")

	// rg.Post("/login", server.Login())
	server.InitializeRoutes(rg)

	return router
}

func (server *Server) Run(httpAddr, httpsAddr string, logger *logrus.Logger) {
	http.Handle("/", server.buildRouter(logger))
	go func() {
		// serve HTTP, which will redirect automatically to HTTPS
		log.Fatal(http.ListenAndServe(httpAddr, nil))
	}()

	// serve HTTPS!
	log.Fatal(http.ListenAndServeTLS(httpsAddr, "server.crt", "server.key", nil))
}

// cacheDir makes a consistent cache directory inside /tmp. Returns "" on error.
func cacheDir() (dir string) {
	if u, _ := user.Current(); u != nil {
		dir = filepath.Join(os.TempDir(), "cache-"+u.Username)
		if err := os.MkdirAll(dir, 0700); err == nil {
			return dir
		}
	}
	return ""
}
