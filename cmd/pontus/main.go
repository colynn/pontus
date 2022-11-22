package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/colynn/pontus/config"
	"github.com/colynn/pontus/internal/db"
	"github.com/colynn/pontus/internal/db/gorm"
	"github.com/colynn/pontus/internal/initialize"
	"github.com/colynn/pontus/internal/router"

	istools "github.com/isbrick/tools"

	log "unknwon.dev/clog/v2"
)

func main() {
	env := os.Getenv("ENV")
	if !istools.IsSliceContainsStr([]string{"dev", "prod", "local"}, env) {
		fmt.Printf("Usage: ENV=[local|dev|prod] ./pontus\n")
		return
	}
	config.Init(env)

	c := config.GetConfig()

	initialize.InitLogging(c)

	log.Info("config env is %s", env)
	// db  models init
	db.Init()
	err := gorm.AutoMigrate(db.Eloquent)
	if err != nil {
		log.Fatal("数据库初始化失败 err: %v", err)
	}

	// router init
	r := router.InitRouter(c)

	defer db.Eloquent.Close()

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%v", c.GetString("application.host"), c.GetInt("application.port")),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// Graceful restart or stop
	// https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s\n", err)
		}
	}()

	if c.GetBool("application.initialize") {
		// role init
		err = initialize.InitRole()
		if err != nil {
			log.Error("init role error: %s", err.Error())
		}

		err = initialize.InitMenu()
		if err != nil {
			log.Error("init menu error: %s", err.Error())
		}

		err = initialize.InitPermission()
		if err != nil {
			log.Error("init role permission error: %s", err.Error())
		}
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Info("timeout of 5 seconds.")

	log.Info("Server exiting")
	log.Stop()
}
