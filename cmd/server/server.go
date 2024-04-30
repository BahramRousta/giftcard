package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"time"
)

type Server struct {
	*echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			e.Logger.Error("Failed to open log file for body dump:", err)
			return
		}
		defer logFile.Close()

		currentTime := time.Now()

		if _, err := fmt.Fprintf(
			logFile,
			"[%s] - Method: %s, "+
				"URI: %s\n"+
				"Request Body: %s\n"+
				"Response Body: %s\n",
			currentTime.Format(time.RFC3339),
			c.Request().Method,
			c.Request().RequestURI,
			string(reqBody),
			string(resBody),
		); err != nil {
			e.Logger.Error("Failed to write request and response payloads to log file:", err)
			return
		}
	}))

	return &Server{e}
}
