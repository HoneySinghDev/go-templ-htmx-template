package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HoneySinghDev/go-templ-htmx-template/internal/router"
	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/sb"
	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"

	"github.com/rs/zerolog/log"

	"github.com/HoneySinghDev/go-templ-htmx-template/internal/config"
	"github.com/rs/zerolog"
)

func App() {
	c := config.DefaultServiceConfig()

	// Setup logger
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(c.LogLevelFromString(c.Logger.Level))

	if c.Logger.PreetyPrintConsole {
		log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "15:04:05"
		}))
	}

	s := server.NewServer(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := s.InitDB(ctx); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	cancel()

	//Init SupaBase Client
	s.SB = sb.InitSB(c.Supabase.ApiUrl, c.Supabase.SecretKey)

	router.Init(s)

	go func() {
		if err := s.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info().Msg("Server closed")
			} else {
				log.Fatal().Err(err).Msg("Failed to start server")
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Failed to gracefully shut down server")
	}
}
