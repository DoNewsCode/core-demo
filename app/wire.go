//+build wireinject

package app

import (
	"github.com/DoNewsCode/skeleton/app/book"
	pb "github.com/DoNewsCode/skeleton/app/proto"
	"github.com/DoNewsCode/skeleton/app/user"
	"github.com/DoNewsCode/skeleton/internal/repositories"
	"github.com/DoNewsCode/std/pkg/contract"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

func InjectAppModule(db *gorm.DB, histogram metrics.Histogram, tracer opentracing.Tracer, logger log.Logger, env contract.Env, redis redis.UniversalClient) (AppModule, error) {
	panic(wire.Build(
		user.NewTransport,
		book.NewTransport,
		repositories.NewUserDao,
		wire.Struct(new(AppModule), "*"),
		wire.Struct(new(book.Handler), "*"),
		wire.Struct(new(user.Service), "*"),
		wire.Bind(new(pb.UserServer), new(user.Service)),
		wire.Bind(new(user.UserRepository), new(*repositories.UserDao)),
	))
}
