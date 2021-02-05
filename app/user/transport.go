package user

import (
	"net/http"

	app_pb "github.com/DoNewsCode/skeleton/app/proto"
	usersvc "github.com/DoNewsCode/skeleton/app/user/gen"
	"github.com/DoNewsCode/std/pkg/contract"
	"github.com/DoNewsCode/std/pkg/key"
	"github.com/DoNewsCode/std/pkg/kitmw"
	"github.com/DoNewsCode/std/pkg/srvgrpc"
	"github.com/DoNewsCode/std/pkg/srvhttp"
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
	keyer := key.NewKeyManager("module", module, "service", service)
	endpoints := usersvc.NewEndpoints(server)
	endpoints.WrapAllLabeledExcept(kitmw.MakeLabeledMetricsMiddleware(metrics, keyer))
	endpoints.WrapAllLabeledExcept(kitmw.MakeLabeledTraceServerMiddleware(tracer, keyer))
	endpoints.WrapAllExcept(kitmw.NewValidationMiddleware())
	endpoints.WrapAllExcept(kitmw.MakeErrorMarshallerMiddleware(kitmw.ErrorOption{
		AlwaysHTTP200: true,
		AlwaysGRPCOk:  true,
		ShouldRecover: env.IsProduction(),
	}))
	endpoints.WrapAllLabeledExcept(kitmw.MakeLabeledLoggingMiddleware(logger, keyer, env.IsLocal()))
	httpHandler := usersvc.MakeHTTPHandler(endpoints,
		httptransport.ServerBefore(
			srvhttp.IpToContext(),
			kittracing.HTTPToContext(tracer, "http", logger),
		),
		httptransport.ServerErrorEncoder(kitmw.ErrorEncoder),
	)
	grpcHandler := usersvc.MakeGRPCServer(endpoints,
		grpctransport.ServerBefore(
			kittracing.GRPCToContext(tracer, "grpc", logger),
			srvgrpc.IpToContext(),
		),
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	)
	return Transport{Handler: httpHandler, UserServer: grpcHandler}
}
