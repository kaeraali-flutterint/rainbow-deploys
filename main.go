package main

import (
	"fmt"
  "log/slog"
	"net"
	"net/http"
	"os"
	"time"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var logger *slog.Logger

func getColor() string {
	return os.Getenv("COLOR")
}

func tcpHandler(c net.Conn) {
	for {
		c.Write([]byte(fmt.Sprintf("The color is #%s", getColor())))
		c.Write([]byte(fmt.Sprintln()))
		time.Sleep(5 * time.Second)

	}
}

func serveTCP() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error but not today
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error but not today
		}
		go tcpHandler(conn)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	color := getColor()

  logger.Info("Serving color", "color", color)

	fmt.Fprintf(w, "<body bgcolor=\"#%s\"><h1>#%s</h1></body>", color, color)
}

func main() {
  logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	color := getColor()

  logger.Info("Booted", "color", color)

	go serveTCP()

	http.HandleFunc("/", httpHandler)
  http.Handle("/metrics", promhttp.Handler())
  logger.Info("listening with http on :8080 and tcp on :8081")
	http.ListenAndServe(":8080", nil)
}
