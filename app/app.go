package app

import (
	"fmt"
	"github.com/KevLehmann/packky-tracker-api/controllers"
	"github.com/KevLehmann/packky-tracker-api/middleware"
	"github.com/KevLehmann/packky-tracker-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"

	"os"
	"strconv"

	// Postgres dialect sideeffects
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

var globalEnv controllers.Env
var err error

func init() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	fmt.Println(fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, pass))
	globalEnv.DB, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, pass))
	if err != nil {
		panic(err)
	}

	migrate, _ := strconv.ParseBool(os.Getenv("AUTO_MIGRATE"))
	seed, _ := strconv.ParseBool(os.Getenv("AUTO_SEED"))
	logDb, _ := strconv.ParseBool(os.Getenv("DB_DEBUG"))
	isProduction := os.Getenv("GIN_MODE") == "release"

	if logDb {
		log.Println("DB Logging is enabled")
		globalEnv.DB.LogMode(true)
	}
	if migrate {
		log.Println("DB Migration is enabled")
		models.AutoMigrate(globalEnv.DB)
	}
	if seed {
		log.Println("[WARNING] DB Seeding is enabled. This may occasionate duplicate data")
		models.Seed(globalEnv.DB, isProduction)
	}
}

// InitRoutes initializes route engine on API
func InitRoutes(r *gin.Engine) {
	// Global Middleware
	r.Use(middleware.ValidateAuthorizationHeader)

	// API routes
	r.POST("/api/intent", globalEnv.IntentController)
}
