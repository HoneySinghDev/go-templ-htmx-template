package config

import (
	"fmt"
	"strings"
)

// DBConnectionString generates a connection string to be passed to sql.Open or equivalents, assuming Postgres syntax.
func (s *Server) DBConnectionString() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", s.Database.PSQLHOST, s.Database.PSQLPORT,
		s.Database.PSQLUSER, s.Database.PSQLPASS, s.Database.PSQLDB))

	if _, ok := s.Database.AdditionalParams["sslmode"]; !ok {
		b.WriteString(" sslmode=disable")
	}

	for key, value := range s.Database.AdditionalParams {
		b.WriteString(fmt.Sprintf(" %s=%s", key, value))
	}

	// Include pool configurations directly in the connection string
	if s.Database.DBMaxOpenConns > 0 {
		b.WriteString(fmt.Sprintf(" pool_max_conns=%d", s.Database.DBMaxOpenConns))
	}
	if s.Database.MaxIdleConns > 0 {
		b.WriteString(fmt.Sprintf(" pool_min_conns=%d", s.Database.MaxIdleConns))
	}
	if s.Database.ConnectionMaxLifetime.GoDuration() > 0 {
		b.WriteString(fmt.Sprintf(" pool_max_conn_lifetime=%s", s.Database.ConnectionMaxLifetime.GoDuration()))
	}

	return b.String()
}
