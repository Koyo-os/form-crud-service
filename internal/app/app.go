package app

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Koyo-os/form-crud-service/internal/config"
	"github.com/Koyo-os/form-crud-service/internal/metrics"
	"github.com/Koyo-os/form-crud-service/internal/server"
	"github.com/Koyo-os/form-crud-service/pkg/api/pb"
	"github.com/Koyo-os/form-crud-service/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	server *server.Server
	logger *logger.Logger
	cfg    *config.Config
	grcpServer *grpc.Server
	healthCheckServer *http.Server
}

func NewApp(server *server.Server) *App {
	return &App{
		server: server,
		logger: logger.Get(),
		cfg:    config.NewConfig(),
	}
}

func (a *App) Start(ctx context.Context) error {
	now := time.Now()

	var wg sync.WaitGroup

	errCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()

		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		httpServer := &http.Server{
			Addr: a.cfg.HealthCheckAddr,
			Handler: mux,
		}

		a.healthCheckServer = httpServer

		a.logger.Info("metrics server starting",
			zap.String("addr", a.cfg.HealthCheckAddr))

		if err := httpServer.ListenAndServe(); err != nil {
			a.logger.Error("error start metrics server",
				zap.String("addr", a.cfg.HealthCheckAddr),
				zap.Error(err))

			errCh <- err

			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		server := grpc.NewServer()

		pb.RegisterFormServiceServer(server, a.server)

		lis, err := net.Listen("tcp", a.cfg.Addr)
		if err != nil {
			a.logger.Error("failed listen",
				zap.String("addr", a.cfg.Addr),
				zap.Error(err))

			errCh <- err

			return
		}

		if err = server.Serve(lis); err != nil {
			a.logger.Error("failed serve",
				zap.String("addr", a.cfg.Addr),
				zap.Error(err))

			errCh <- err

			return
		}

		a.grcpServer = server
	}()

	wg.Wait()

	for err := range errCh {
		return err
	}

	metrics.StartTime.WithLabelValues("form").Observe(time.Since(now).Seconds())

	a.logger.Info("form service started successfully",
		zap.String("healthcheck", a.cfg.HealthCheckAddr),
		zap.String("grpc_server", a.cfg.Addr))

	return nil
}

func (a *App) Close() error {
	a.grcpServer.Stop()

	return a.healthCheckServer.Close()
}