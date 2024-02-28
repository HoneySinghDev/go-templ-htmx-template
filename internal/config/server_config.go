package config

import (
	"context"

	"github.com/HoneySinghDev/go-templ-htmx-template/pkl/pklgen"
	"github.com/HoneySinghDev/go-templ-htmx-template/pkl/pklgen/loggerlevel"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Database   *pklgen.DatabaseConfig
	Echo       *pklgen.EchoConfig
	Supabase   *pklgen.SupabaseConfig
	Logger     *pklgen.LoggerConfig
	Auth       *pklgen.AuthServerConfig
	Management *pklgen.ManagementServerConfig
}

func DefaultServiceConfig() Server {
	cfg, err := pklgen.LoadFromPath(context.Background(), "pkl/local/appConfig.pkl")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	return Server{
		Database:   cfg.DB,
		Echo:       cfg.Echo,
		Supabase:   cfg.SB,
		Logger:     cfg.Logger,
		Auth:       cfg.Auth,
		Management: cfg.Management,
	}
}

func (s *Server) LogLevelFromString(l loggerlevel.LoggerLevel) zerolog.Level {
	lv, err := zerolog.ParseLevel(l.String())
	if err != nil || lv == zerolog.NoLevel {
		log.Error().Err(err).Msgf("Failed to parse log level, defaulting to %s", zerolog.DebugLevel)
		return zerolog.DebugLevel
	}

	return lv
}
