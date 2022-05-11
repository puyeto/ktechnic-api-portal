package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ktechnics/ktechnics-api/api/app"
	"github.com/ktechnics/ktechnics-api/api/controllers"
	"github.com/sirupsen/logrus"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file not found")
	}
}

// Run ...
func Run(logger *logrus.Logger) {

	var err error
	err = godotenv.Load()
	if err != nil {
		logger.Fatalf("Error getting env, %v", err)
	} else {
		logger.Println("We are getting the env values")
	}

	app.MongoDB = app.InitializeMongoDB(os.Getenv("MONGO_PROD_DNS"), os.Getenv("MONGO_DB_NAME"), logger)

	if os.Getenv("GO_ENV") == "production" {
		server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_PROD_USER"), os.Getenv("DB_PROD_PASSWORD"), os.Getenv("DB_PROD_PORT"), os.Getenv("DB_PROD_HOST"), os.Getenv("DB_PROD_NAME"), logger)
	} else {
		server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), logger)
	}
	// seed.Load(server.DB)

	server.Run(":9002", ":9003", logger)

}
