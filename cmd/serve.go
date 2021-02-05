package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DoNewsCode/std/pkg/core"
	"github.com/DoNewsCode/std/pkg/srvhttp"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewServeCommand(c *core.C) *cobra.Command {
	getHttpHandler := func(ln net.Listener, providers ...func(*mux.Router)) http.Handler {
		c.Info(fmt.Sprintf("http service is listening at %s", ln.Addr()))

		var handler http.Handler
		var router = mux.NewRouter()
		for _, p := range providers {
			p(router)
		}
		handler = srvhttp.MakeUnsafeCorsMiddleware()(router)
		return handler
	}

	getGRPCServer := func(ln net.Listener, providers ...func(s *grpc.Server)) *grpc.Server {
		c.Info(fmt.Sprintf("gRPC service is listening at %s", ln.Addr()))

		s := grpc.NewServer(grpc.ConnectionTimeout(time.Second))
		for _, p := range providers {
			p(s)
		}
		return s
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Long:  `Start the gRPC server and HTTP server`,
		Run: func(cmd *cobra.Command, args []string) {

			var g run.Group

			// Start HTTP Server
			{
				httpAddr := c.String("http.addr")
				ln, err := net.Listen("tcp", httpAddr)
				c.CheckErr(err)

				h := getHttpHandler(ln, c.HttpProviders...)
				srv := &http.Server{
					Handler:      h,
					IdleTimeout:  2 * time.Second,
					ReadTimeout:  2 * time.Second,
					WriteTimeout: 2 * time.Second,
				}
				g.Add(func() error {
					return srv.Serve(ln)
				}, func(err error) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					err = srv.Shutdown(ctx)
					c.CheckErr(err)
					_ = ln.Close()
				})
			}

			// Start gRPC server
			{
				grpcAddr := c.String("grpc.addr")
				ln, err := net.Listen("tcp", grpcAddr)
				c.CheckErr(err)

				s := getGRPCServer(ln, c.GrpcProviders...)
				g.Add(func() error {
					return s.Serve(ln)
				}, func(err error) {
					s.GracefulStop()
					_ = ln.Close()
				})
			}

			// Add Crontab
			{
				tab := cron.New()
				for _, f := range c.CronProviders {
					f(tab)
				}
				g.Add(func() error {
					tab.Run()
					return nil
				}, func(err error) {
					<-tab.Stop().Done()
				})
			}

			// Graceful shutdown
			{
				sig := make(chan os.Signal, 1)
				signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
				g.Add(func() error {
					terminateError := fmt.Errorf("%s", <-sig)
					return terminateError
				}, func(err error) {
					close(sig)
				})
			}

			// Additional run groups
			for _, s := range c.RunProviders {
				s(&g)
			}

			if err := g.Run(); err != nil {
				c.Err(err)
				os.Exit(1)
			}

			c.Info("graceful shutdown complete; see you next time")
		},
	}
	return serveCmd
}
