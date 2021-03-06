package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var filesHandler = fasthttp.FSHandler(staticFilePath, 0)

func routeHandler(ctx *fasthttp.RequestCtx) {
	// 127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
	userIP := ctx.RemoteIP().String()
	if tempIP := string(ctx.Request.Header.Peek("X-Forwarded-For")); tempIP != "" {
		userIP = tempIP
	}
	ctx.Response.Header.Add("Access-Control-Allow-Origin", cors)
	log.Infof("%s %s [%s] %s \"%s\"", userIP, ctx.UserAgent(), time.Now().Format("02/01/2006 15:04:05 -0700"), ctx.Method(), ctx.Path())
	switch string(ctx.Path()) {
	case "/departures":
		ctx.SetContentType("application/json")
		busHandler(ctx)
	case "/graphql":
		ctx.SetContentType("application/json")
		graphqlHandler(ctx)
	case "/graphiql":
		graphiqlHandler(ctx)
	default:
		if staticFilePath != "" {
			filesHandler(ctx)
		}
	}
}
