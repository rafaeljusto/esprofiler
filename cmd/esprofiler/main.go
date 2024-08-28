package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rafaeljusto/esprofiler/internal/config"
	"github.com/rafaeljusto/esprofiler/internal/web"
	"go.uber.org/multierr"
)

func main() {
	defer handleExit()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	config, errs := config.ParseFromEnvs()
	if errs != nil {
		for _, err := range multierr.Errors(errs) {
			logger.Error("failed to parse configuration",
				slog.String("error", err.Error()),
			)
		}
		exit(exitCodeInvalidInput)
	}

	listener, err := net.Listen("tcp", ":"+strconv.FormatInt(config.Port, 10))
	if err != nil {
		logger.Error("failed to listen",
			slog.String("error", err.Error()),
		)
		exit(exitCodeSetupFailure)
	}

	logger.Info("starting web server",
		slog.String("address", listener.Addr().String()),
	)

	router := http.NewServeMux()
	web.RegisterHandlers(router, logger)
	server := http.Server{
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				logger.Error("failed to serve",
					slog.String("error", err.Error()),
				)
			}
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed",
			slog.String("error", err.Error()),
		)
	}
	logger.Info("server stopped")
}

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeInvalidInput
	exitCodeSetupFailure
)

type exitData struct {
	code exitCode
}

// exit allows to abort the program while still executing all defer statements.
func exit(code exitCode) {
	panic(exitData{code: code})
}

// handleExit exit code handler.
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(exitData); ok {
			os.Exit(int(exit.code))
		}
		panic(e)
	}
}
