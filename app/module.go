//go:generate trs proto/user.proto --doc=../docs/swagger --lib=../third_party
package app

import (
	"net/http"

	"github.com/DoNewsCode/core/otgorm"
	"github.com/DoNewsCode/skeleton/app/book"
	pb "github.com/DoNewsCode/skeleton/app/proto"
	"github.com/DoNewsCode/skeleton/app/user"
	"github.com/DoNewsCode/skeleton/internal/repositories"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type AppModule struct {
	UserTransport user.Transport
	BookTransport book.Transport
}

func (a AppModule) ProvideHTTP(router *mux.Router) {
	router.PathPrefix("/app/user/").Handler(http.StripPrefix("/app/user", a.UserTransport))
	router.PathPrefix("/app/book/").Handler(http.StripPrefix("/app/book", a.BookTransport))
}

func (a AppModule) ProvideGRPC(server *grpc.Server) {
	pb.RegisterUserServer(server, a.UserTransport)
}

func (a AppModule) ProvideMigration() []*otgorm.Migration {
	return repositories.ProvideMigrator()
}

func (a AppModule) ProvideSeed() []*otgorm.Seed {
	return repositories.ProvideSeed()
}

func (a AppModule) SeedRedis() func(client redis.UniversalClient) error {
	return book.Seed
}
