package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/image_processing_service/internal/config"
	"github.com/Falokut/image_processing_service/internal/default_processing"
	"github.com/Falokut/image_processing_service/internal/service"
	image_service "github.com/Falokut/image_processing_service/pkg/image_processing_service/v1/protos"
	jaegerTracer "github.com/Falokut/image_processing_service/pkg/jaeger"
	"github.com/Falokut/image_processing_service/pkg/metrics"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func main() {

	logging.NewEntry(logging.FileAndConsoleOutput)
	logger := logging.GetLogger()

	appCfg := config.GetConfig()
	log_level, err := logrus.ParseLevel(appCfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(log_level)
	var metric metrics.Metrics
	if appCfg.EnableMetrics {
		tracer, closer, err := jaegerTracer.InitJaeger(appCfg.JaegerConfig)
		if err != nil {
			logger.Fatal("cannot create tracer", err)
		}
		logger.Info("Jaeger connected")
		defer closer.Close()

		opentracing.SetGlobalTracer(tracer)

		logger.Info("Metrics initializing")
		metric, err = metrics.CreateMetrics(appCfg.PrometheusConfig.Name)
		if err != nil {
			logger.Fatal(err)
		}

		go func() {
			logger.Info("Metrics server running")
			if err := metrics.RunMetricServer(appCfg.PrometheusConfig.ServerConfig); err != nil {
				logger.Fatal(err)
			}
		}()
	} else {
		metric = &metrics.EmptyMetrics{}
	}

	imagesProcessing := default_processing.NewImageProcessingService(logger.Logger)
	logger.Info("Service initializing")
	service := service.NewImagesProcessingService(logger.Logger, imagesProcessing)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)
	s.Run(getListenServerConfig(appCfg), metric, nil, nil)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.Shutdown()
}

const (
	kb = 8 << 10
	mb = kb << 10
)

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Host:        cfg.Listen.Host,
		Port:        cfg.Listen.Port,
		Mode:        cfg.Listen.Mode,
		ServiceDesc: &image_service.ImageProcessingServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(image_service.ImageProcessingServiceV1Server)
			if !ok {
				return errors.New("error while creating images processing service")
			}
			return image_service.RegisterImageProcessingServiceV1HandlerServer(ctx, mux, serv)
		},
		MaxRequestSize:  cfg.Listen.MaxRequestSize * mb,
		MaxResponceSize: cfg.Listen.MaxResponseSize * mb,
	}
}
