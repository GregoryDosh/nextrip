package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var filesHandler = fasthttp.FSHandler(staticFilePath, 0)

func routeHandler(ctx *fasthttp.RequestCtx) {
	// 127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
	log.Infof("%s %s [%s] %s \"%s\"", ctx.RemoteIP(), ctx.UserAgent(), time.Now().Format("02/01/2006 15:04:05 -0700"), ctx.Method(), ctx.Path())
	switch string(ctx.Path()) {
	case "/departures":
		busHandler(ctx)
	case "/graphql":
		graphqlHandler(ctx)
	default:
		if staticFilePath != "" {
			filesHandler(ctx)
		}
	}
}
