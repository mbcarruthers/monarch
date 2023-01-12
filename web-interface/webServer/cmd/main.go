package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	__port = 3000
)

var (
	port = fmt.Sprintf(":%d", __port)
)

func main() {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./client", true)))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error":   "false",
			"message": "Fingers Crossed",
		})
		return
	})

	server := &http.Server{
		Addr:    port,
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listening err: %s \n",
				err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*750)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("webServer shutdown err: %s \n",
			err.Error())
	}

	select {
	case <-ctx.Done():
		break;
	}
	log.Println("Web Server has shutdown properly")
}
