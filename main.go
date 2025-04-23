package main

import (
	"context"
	"datn_backend/config"
	_ "datn_backend/docs"
	"datn_backend/middleware"
	"datn_backend/middleware/logger"
	"datn_backend/router"
	sso_utils "datn_backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"os/signal"
	"strconv"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if hasEnv := godotenv.Load(); hasEnv != nil {
		log.Printf(hasEnv.Error())
	}

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Tracer())
	engine.Use(middleware.Cors())
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	config.InitializeDatabase()
	config.SetRedisStore(engine)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	if os.Getenv("TELEGRAM_BOT_TOKEN") != "" {
		sso_utils.CreateBot(ctx, os.Getenv("TELEGRAM_BOT_TOKEN"))
	}

	engine.Use(middleware.LogResponse(logger.InitLogger()))
	router.RegisterRoutes(engine)
	config.InitializeFirebase()
	//handle websocket monitoring login device
	//engine.GET("/ws", func(c *gin.Context) {
	//	websocket.HandleWebSocket(c.Writer, c.Request)
	//})
	//engine.GET("/ws/monitor", func(c *gin.Context) {
	//	websocket.HandleMonitorWebSocket(c.Writer, c.Request)
	//})
	log.Println(fmt.Sprintf("✅ WebSocket server listening on :%s", os.Getenv("PORT")))
	errRun := engine.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if errRun != nil {
		return
	}
}
