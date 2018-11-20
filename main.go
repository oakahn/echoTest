package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/oakahn/echoTest/fatca"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Echo instance
	e := echo.New()
	hostname, _ := os.Hostname()
	filename := hostname + ".log"
	out := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1, // megabytes
		MaxBackups: 100,
		MaxAge:     1,    //days
		Compress:   true, // disabled by default
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\n",
		Output: out,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		c.Logger().SetOutput(out)
		c.Logger().Print(string(reqBody))
	}))

	// Routes
	e.POST("/", fatca.Call)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
