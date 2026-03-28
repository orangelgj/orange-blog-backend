package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gblog/config"
	_ "gblog/docs"
	"gblog/router"
	"gblog/utils"
)

// @title           橘子的个人博客系统 API
// @version         1.0
// @description     这是一个个人博客的后端接口
// @host      localhost:8081
// @BasePath  /api/v1
func main() {
	utils.InitLogger()
	utils.Logger.Info("Starting application...")

	utils.Logger.Info("Loading configuration...")
	config.InitConfig()
	utils.Logger.Info("Configuration loaded successfully")

	utils.Logger.Info("Setting up router...")
	r := router.SetupRouter()
	utils.Logger.Info("Router setup completed")

	port := ":8081"
	utils.Logger.WithField("port", port).Info("Starting HTTP server...")
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	utils.Logger.Info("Received shutdown signal, gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Logger.WithError(err).Error("Server shutdown error")
	}

	utils.Logger.Info("Server exited successfully")
}
