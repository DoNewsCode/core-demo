package user

import (
	"github.com/DoNewsCode/core-kit/option"
	"net/http"

	"github.com/DoNewsCode/core-kit/mw"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/key"
	app_pb "github.com/DoNewsCode/skeleton/app/proto"
	usersvc "github.com/DoNewsCode/skeleton/app/user/gen"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kittracing "github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/opentracing/opentracing-go"
)

const module = "app"
const service = "user"

type Transport struct {
	http.Handler
	app_pb.UserServer
}

func NewTransport(
	server app_pb.UserServer,
	logger log.Logger,
	env contract.Env,
	tracer opentracing.Tracer,
	metrics metrics.Histogram,
) Transport {
	keyer := key.New("module", module, "service", service)
	endpoints := usersvc.NewEndpoints(server)
	endpoints.WrapAllLabeledExcept(mw.LabeledMetrics(metrics, keyer))
	endpoints.WrapAllLabeledExcept(mw.LabeledTraceServer(tracer, keyer))
	endpoints.WrapAllExcept(mw.Validate())
	endpoints.WrapAllExcept(mw.Error(mw.ErrorOption{
		AlwaysHTTP200: true,
		ShouldRecover: env.IsProduction(),
	}))
	endpoints.WrapAllLabeledExcept(mw.LabeledLog(logger, keyer, env.IsLocal()))
	httpHandler := usersvc.MakeHTTPHandler(endpoints,
		httptransport.ServerBefore(
			option.IPToHTTPContext(),
			kittracing.HTTPToContext(tracer, "http", logger),
		),
		httptransport.ServerErrorEncoder(option.ErrorEncoder),
	)
	grpcHandler := usersvc.MakeGRPCServer(endpoints,
		grpctransport.ServerBefore(
			kittracing.GRPCToContext(tracer, "grpc", logger),
			option.IPToGRPCContext(),
		),
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	)
	return Transport{Handler: httpHandler, UserServer: grpcHandler}
}
