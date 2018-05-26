package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func busHandler(ctx *fasthttp.RequestCtx) {
	queryArgs := ctx.QueryArgs()
	stop, err := strconv.Atoi(string(queryArgs.Peek("stop_id")))
	if err != nil {
		fmt.Fprint(ctx, "bad stop specified")
		log.Error(err)
		return
	}
	s, err := env.GetDepartures(stop)
	if err != nil {
		log.Error(err)
		return
	}

	b, err := json.Marshal(s)
	if err != nil {
		log.Error(err)
	}

	if _, err := fmt.Fprint(ctx, string(b)); err != nil {
		log.Error(err)
	}
}
