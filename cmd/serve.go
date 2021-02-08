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
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewServeCommand(c *core.C) *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Long:  `Start the gRPC server, HTTP server, and cron job runner.`,
		Run: func(cmd *cobra.Command, args []string) {

			var g run.Group

			// Start HTTP server
			{
				httpAddr := c.String("http.addr")
				ln, err := net.Listen("tcp", httpAddr)
				c.CheckErr(err)

				srv := &http.Server{
					Handler:      collectHttpHandler(c.HttpProviders),
					IdleTimeout:  60 * time.Second,
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 60 * time.Second,
				}
				g.Add(func() error {
					c.Info(fmt.Sprintf("http service is listening at %s", ln.Addr()))
					return srv.Serve(ln)
				}, func(err error) {
					_ = srv.Shutdown(context.Background())
					_ = ln.Close()
				})
			}

			// Start gRPC server
			{
				grpcAddr := c.String("grpc.addr")
				ln, err := net.Listen("tcp", grpcAddr)
				c.CheckErr(err)

				s := collectGrpcServer(c.GrpcProviders)
				g.Add(func() error {
					c.Info(fmt.Sprintf("gRPC service is listening at %s", ln.Addr()))
					return s.Serve(ln)
				}, func(err error) {
					s.GracefulStop()
					_ = ln.Close()
				})
			}

			// Start cron runner
			{
				crontab := collectCron(c.CronProviders)
				g.Add(func() error {
					c.Info(fmt.Sprintf("cron runner started"))
					crontab.Run()
					return nil
				}, func(err error) {
					<-crontab.Stop().Done()
				})
			}

			// Config hot reload
			{
				ctx, cancel := context.WithCancel(context.Background())
				g.Add(func() error {
					return c.ConfigAccessor.(interface{ Watch(context.Context) error }).Watch(ctx)
				}, func(err error) {
					cancel()
				})
			}

			// Graceful shutdown
			{
				sig := make(chan os.Signal, 1)
				signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
				g.Add(func() error {
					c.Err(fmt.Errorf("signal received: %s", <-sig))
					return nil
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
				os.Exit(255)
			}

			c.Info("graceful shutdown complete; see you next time :)")
		},
	}
	return serveCmd
}

func collectHttpHandler(providers []func(*mux.Router)) http.Handler {
	var router = mux.NewRouter()
	for _, p := range providers {
		p(router)
	}
	return router
}

func collectGrpcServer(providers []func(s *grpc.Server)) *grpc.Server {
	s := grpc.NewServer()
	for _, p := range providers {
		p(s)
	}
	return s
}

func collectCron(providers []func(s *cron.Cron)) *cron.Cron {
	c := cron.New()
	for _, p := range providers {
		p(c)
	}
	return c
}
