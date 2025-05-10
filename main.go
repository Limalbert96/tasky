package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	controller "github.com/jeffthorne/tasky/controllers"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func main() {
	godotenv.Overload()

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize New Relic with debug logging
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("tasky"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize New Relic: %v", err))
	}

	// Wait for New Relic to connect
	if err := app.WaitForConnection(5 * time.Second); err != nil {
		panic(fmt.Sprintf("Failed to connect to New Relic: %v", err))
	}

	router := gin.Default()
	// Add New Relic middleware
	router.Use(nrgin.Middleware(app))
	router.LoadHTMLGlob("assets/*.html")
	router.Static("/assets", "./assets")

	// Handle both GET and HEAD requests for root path
	router.GET("/", index)
	router.HEAD("/", index)

	// Handle API routes
	router.GET("/todos/:userid", controller.GetTodos)
	router.HEAD("/todos/:userid", controller.GetTodos)
	router.GET("/todo/:id", controller.GetTodo)
	router.HEAD("/todo/:id", controller.GetTodo)
	router.POST("/todo/:userid", controller.AddTodo)
	router.DELETE("/todo/:userid/:id", controller.DeleteTodo)
	router.DELETE("/todos/:userid", controller.ClearAll)
	router.PUT("/todo", controller.UpdateTodo)

	router.POST("/signup", controller.SignUp)
	router.POST("/login", controller.Login)
	router.GET("/todo", controller.Todo)

	router.Run(":8080")

}
