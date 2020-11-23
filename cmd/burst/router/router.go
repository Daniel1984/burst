package router

import (
	"github.com/bitburst/burstconsumer/cmd/burst/burstcfg"
	"github.com/bitburst/burstconsumer/cmd/burst/handlers/callback"
	"github.com/julienschmidt/httprouter"
)

func Get(cfg *burstcfg.BurstCfg) *httprouter.Router {
	mux := httprouter.New()
	mux.POST("/callback", callback.Do(cfg))
	return mux
}
