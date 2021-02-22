package user

import (
	"net/http"

	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/key"
	"github.com/DoNewsCode/core/kitmw"
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
	endpoints.WrapAllLabeledExcept(kitmw.MakeLabeledMetricsMiddleware(metrics, keyer))
	endpoints.WrapAllLabeledExcept(kitmw.MakeLabeledTraceServerMiddleware(tracer, keyer))
	endpoints.WrapAllExcept(kitmw.MakeValidationMiddleware())
	endpoints.WrapAllExcept(kitmw.MakeErrorConversionMiddleware(kitmw.ErrorOption{
		AlwaysHTTP200: true,
		ShouldRecover: env.IsProduction(),
	}))
	endpoints.WrapAllLabeledExcept(kitmw.MakeLabeledLoggingMiddleware(logger, keyer, env.IsLocal()))
	httpHandler := usersvc.MakeHTTPHandler(endpoints,
		httptransport.ServerBefore(
			kitmw.IpToHTTPContext(),
			kittracing.HTTPToContext(tracer, "http", logger),
		),
		httptransport.ServerErrorEncoder(kitmw.ErrorEncoder),
	)
	grpcHandler := usersvc.MakeGRPCServer(endpoints,
		grpctransport.ServerBefore(
			kittracing.GRPCToContext(tracer, "grpc", logger),
			kitmw.IpToGRPCContext(),
		),
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	)
	return Transport{Handler: httpHandler, UserServer: grpcHandler}
}
