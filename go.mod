module github.com/DoNewsCode/skeleton

go 1.14

replace github.com/DoNewsCode/core => /Users/donew/src/std

require (
	github.com/DoNewsCode/core v0.7.3
	github.com/DoNewsCode/core-gin v0.1.0
	github.com/DoNewsCode/core-kit v0.1.0
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-kit/kit v0.11.0
	github.com/go-redis/redis/v8 v8.8.2
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.4.3
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.1.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	google.golang.org/genproto a8c4777a87af
	google.golang.org/grpc v1.37.0
	gorm.io/gorm v1.21.9
)
