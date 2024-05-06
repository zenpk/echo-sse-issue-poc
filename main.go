package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.GET("/sse-get", func(c echo.Context) error {
		log.Printf("SSE GET")

		w := c.Response()
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		ticker := time.NewTicker(10 * time.Millisecond)
		timeout, _ := context.WithTimeout(context.Background(), 1*time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-timeout.Done():
				log.Printf("SSE client disconnected, ip: %v", c.RealIP())
				return nil
			case <-ticker.C:
				event := Event{
					Data: []byte("time: " + time.Now().Format(time.RFC3339Nano)),
				}
				if err := event.MarshalTo(w); err != nil {
					return err
				}
				w.Flush()
			}
		}
	})
	e.POST("/sse-post", func(c echo.Context) error {
		log.Printf("SSE POST")

		w := c.Response()
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		ticker := time.NewTicker(10 * time.Millisecond)
		timeout, _ := context.WithTimeout(context.Background(), 1*time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-timeout.Done():
				log.Printf("SSE client disconnected, ip: %v", c.RealIP())
				return nil
			case <-ticker.C:
				event := Event{
					Data: []byte("time: " + time.Now().Format(time.RFC3339Nano)),
				}
				if err := event.MarshalTo(w); err != nil {
					return err
				}
				w.Flush()
			}
		}
	})
	e.StartServer(&http.Server{
		Addr: ":8080",
	})
}
