package main

import (
	"encoding/json"
	"github.com/Alan15r/GasPriceTestTask/ethereum"
	"github.com/valyala/fasthttp"
	"log"
	"sync"
)

func main() {
	server := &fasthttp.Server{
		Handler: statistics,
	}

	log.Println("run API")
	log.Fatal(server.ListenAndServe(":9090"))
}

func statistics(ctx *fasthttp.RequestCtx) {
	var data ethereum.Data
	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		log.Println("ERROR:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Response.SetBody([]byte(err.Error()))
	} else {
		answer := ethereum.Answer{}
		var wg sync.WaitGroup
		wg.Add(4)
		go data.Ethereum.AveragePricePerDay(&wg, &answer)
		go data.Ethereum.SpentInMonth(&wg, &answer)
		go data.Ethereum.TotalСosts(&wg, &answer)
		go data.Ethereum.FrequencyDistribution(&wg, &answer)
		wg.Wait()
		setResponseBody(ctx, &answer)
	}
}

func setResponseBody(ctx *fasthttp.RequestCtx, answer *ethereum.Answer) {
	body, err := json.Marshal(answer)
	if err != nil {
		log.Println("ERROR:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Response.SetBody([]byte(err.Error()))
	} else {
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.SetBody(body)
	}
}
