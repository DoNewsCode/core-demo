package book

import (
	"net/http"

	"github.com/DoNewsCode/core/ginmw"
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
	r.Use(ginmw.WithContext())
	r.Use(ginmw.WithLogger(logger, keyer))
	r.Use(ginmw.WithMetrics(hist, keyer, false))
	r.Use(ginmw.WithTrace(tracer, keyer))
	r.Use(gin.Recovery())
	r.GET("/", b.Find)
	return Transport{Handler: r}
}
