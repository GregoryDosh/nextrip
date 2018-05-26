package main

import (
	"github.com/valyala/fasthttp"
)

var filesHandler = fasthttp.FSHandler(".", 0)

func routeHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/departures":
		busHandler(ctx)
	case "/graphql":
		graphqlHandler(ctx)
	default:
		filesHandler(ctx)
	}
}
