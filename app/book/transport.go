package book

import (
	"net/http"

	"github.com/DoNewsCode/core-gin/mw"
	"github.com/DoNewsCode/core/key"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/opentracing/opentracing-go"
)

const (
	module  = "app"
	service = "book"
)

type Transport struct {
	http.Handler
}

func NewTransport(b Handler, logger log.Logger, hist metrics.Histogram, tracer opentracing.Tracer) Transport {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	keyer := key.New("module", module, "service", service)
	r.Use(mw.Context())
	r.Use(mw.Log(logger, keyer))
	r.Use(mw.Metrics(hist, keyer, false))
	r.Use(mw.Trace(tracer, keyer))
	r.Use(gin.Recovery())
	r.GET("/", b.Find)
	return Transport{Handler: r}
}
