package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chayutK/skill-management-incubator/backend/repository"
	"github.com/gin-gonic/gin"

	skill "github.com/chayutK/skill-management-incubator/backend/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	repository.Sync()
	db := repository.DB
	defer db.Close()

	r := gin.Default()

	skill.DB = db
	r.GET("/", skill.HelloWorldHandler)
	r.GET("/api/v1/skills", skill.GetAllHandler)
	r.GET("/api/v1/skills/:key", skill.GetByKeyHandler)

	r.POST("/api/v1/skills", skill.CreateHandler)

	r.PUT("/api/v1/skills/:key", skill.UpdateHandler)

	r.PATCH("/api/v1/skills/:key/actions/name", skill.UpdateNameHandler)
	r.PATCH("/api/v1/skills/:key/actions/description", skill.UpdateDescriptionHandler)
	r.PATCH("/api/v1/skills/:key/actions/logo", skill.UpdateLogoHandler)
	r.PATCH("/api/v1/skills/:key/actions/tags", skill.UpdateTagsHandler)

	r.DELETE("/api/v1/skills/:key", skill.DeleteHandler)

	srv := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	cancelChannel := make(chan struct{})
	go func() {
		<-ctx.Done()
		fmt.Println("Server is shutting down......")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			//  ทำไม error วะ
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println("Error while shutting down the server.", err.Error())
			}
		}

		close(cancelChannel)
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Println("Error while starting server.", err.Error())
	}

	<-cancelChannel
	log.Println("Server is closed.")
}
