package middleware

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	util "github.com/HoneySinghDev/go-templ-htmx-template/pkg/utils"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RequestBodyLogSkipper defines a function to skip logging certain request bodies.
// Returning true skips logging the payload of the request.
type RequestBodyLogSkipper func(req *http.Request) bool

// DefaultRequestBodyLogSkipper returns true for all requests with Content-Type
// application/x-www-form-urlencoded or multipart/form-data as those might contain
// binary or URL-encoded file uploads unfit for logging purposes.
func DefaultRequestBodyLogSkipper(req *http.Request) bool {
	contentType := req.Header.Get(echo.HeaderContentType)
	switch {
	case strings.HasPrefix(contentType, echo.MIMEApplicationForm),
		strings.HasPrefix(contentType, echo.MIMEMultipartForm):
		return true
	default:
		return false
	}
}

// ResponseBodyLogSkipper defines a function to skip logging certain response bodies.
// Returning true skips logging the payload of the response.
type ResponseBodyLogSkipper func(req *http.Request, res *echo.Response) bool

// DefaultResponseBodyLogSkipper returns false for all responses with Content-Type
// application/json, preventing logging for all other types of payloads as those
// might contain binary or URL-encoded data unfit for logging purposes.
func DefaultResponseBodyLogSkipper(_ *http.Request, res *echo.Response) bool {
	contentType := res.Header().Get(echo.HeaderContentType)
	switch {
	case strings.HasPrefix(contentType, echo.MIMEApplicationJSON):
		return false
	default:
		return true
	}
}

// BodyLogReplacer defines a function to replace certain parts of a body before logging it,
// mainly used to strip sensitive information from a request or response payload.
// The []byte returned should contain a sanitized payload ready for logging.
type BodyLogReplacer func(body []byte) []byte

// DefaultBodyLogReplacer returns the body received without any modifications.
func DefaultBodyLogReplacer(body []byte) []byte {
	return body
}

// HeaderLogReplacer defines a function to replace certain parts of a header before logging it,
// mainly used to strip sensitive information from a request or response header.
// The http.Header returned should be a sanitized copy of the original header as not to modify
// the request or response while logging.
type HeaderLogReplacer func(header http.Header) http.Header

// DefaultHeaderLogReplacer replaces all Authorization, X-CSRF-Token and Proxy-Authorization
// header entries with a redacted string, indicating their presence without revealing actual,
// potentially sensitive values in the logs.
func DefaultHeaderLogReplacer(header http.Header) http.Header {
	sanitizedHeader := http.Header{}

	for k, vv := range header {
		shouldRedact := strings.EqualFold(k, echo.HeaderAuthorization) ||
			strings.EqualFold(k, echo.HeaderXCSRFToken) ||
			strings.EqualFold(k, "Proxy-Authorization")

		for _, v := range vv {
			if shouldRedact {
				sanitizedHeader.Add(k, "*****REDACTED*****")
			} else {
				sanitizedHeader.Add(k, v)
			}
		}
	}

	return sanitizedHeader
}

// QueryLogReplacer defines a function to replace certain parts of a URL query before logging it,
// mainly used to strip sensitive information from a request query.
// The url.Values returned should be a sanitized copy of the original query as not to modify the
// request while logging.
type QueryLogReplacer func(query url.Values) url.Values

// DefaultQueryLogReplacer returns the query received without any modifications.
func DefaultQueryLogReplacer(query url.Values) url.Values {
	return query
}

type LoggerConfig struct {
	Skipper                   middleware.Skipper
	Level                     zerolog.Level
	LogRequestBody            bool
	LogRequestHeader          bool
	LogRequestQuery           bool
	LogCaller                 bool
	RequestBodyLogSkipper     RequestBodyLogSkipper
	RequestBodyLogReplacer    BodyLogReplacer
	RequestHeaderLogReplacer  HeaderLogReplacer
	RequestQueryLogReplacer   QueryLogReplacer
	LogResponseBody           bool
	LogResponseHeader         bool
	ResponseBodyLogSkipper    ResponseBodyLogSkipper
	ResponseBodyLogReplacer   BodyLogReplacer
	ResponseHeaderLogReplacer HeaderLogReplacer
}

// Initialize default configuration within the function scope to avoid global variable.
func getDefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Skipper:                  middleware.DefaultSkipper,
		Level:                    zerolog.DebugLevel,
		LogRequestBody:           false,
		LogRequestHeader:         false,
		LogRequestQuery:          false,
		RequestBodyLogSkipper:    DefaultRequestBodyLogSkipper,
		RequestBodyLogReplacer:   DefaultBodyLogReplacer,
		RequestHeaderLogReplacer: DefaultHeaderLogReplacer,
		RequestQueryLogReplacer:  DefaultQueryLogReplacer,
		LogResponseBody:          false,
		LogResponseHeader:        false,
		ResponseBodyLogSkipper:   DefaultResponseBodyLogSkipper,
		ResponseBodyLogReplacer:  DefaultBodyLogReplacer,
	}
}

// Logger with default logger output and configuration.
func Logger() echo.MiddlewareFunc {
	return WithConfig(getDefaultLoggerConfig(), nil)
}

// WithConfig returns a new MiddlewareFunc which creates a logger with the desired configuration.
// If output is set to nil, the default output is used. If more output params are provided, the first is being used.
func WithConfig(config LoggerConfig, output ...io.Writer) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = getDefaultLoggerConfig().Skipper
	}
	if config.RequestBodyLogSkipper == nil {
		config.RequestBodyLogSkipper = DefaultRequestBodyLogSkipper
	}
	if config.RequestBodyLogReplacer == nil {
		config.RequestBodyLogReplacer = DefaultBodyLogReplacer
	}
	if config.RequestHeaderLogReplacer == nil {
		config.RequestHeaderLogReplacer = DefaultHeaderLogReplacer
	}
	if config.RequestQueryLogReplacer == nil {
		config.RequestQueryLogReplacer = DefaultQueryLogReplacer
	}
	if config.ResponseBodyLogSkipper == nil {
		config.ResponseBodyLogSkipper = DefaultResponseBodyLogSkipper
	}
	if config.ResponseBodyLogReplacer == nil {
		config.ResponseBodyLogReplacer = DefaultBodyLogReplacer
	}
	if config.ResponseHeaderLogReplacer == nil {
		config.ResponseHeaderLogReplacer = DefaultHeaderLogReplacer
	}

	return processRequestAndResponse(config, output...)
}

//nolint:funlen,gocognit,cyclop // This function is long and complex by design.
func processRequestAndResponse(config LoggerConfig, output ...io.Writer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			l := log.With().
				Str("id", id).
				Str("host", req.Host).
				Str("method", req.Method).
				Str("url", req.URL.String()).
				Str("bytes_in", req.Header.Get(echo.HeaderContentLength)).
				Logger()

			if len(output) > 0 {
				l = l.Output(output[0])
			}

			if config.LogCaller {
				l = l.With().Caller().Logger()
			}

			le := l.WithLevel(config.Level)
			req = req.WithContext(context.WithValue(req.Context(), util.CTXKeyRequestID, id))

			var err error
			if config.LogRequestBody && !config.RequestBodyLogSkipper(req) {
				reqBody, err := io.ReadAll(req.Body)
				if err != nil {
					le.Err(err).Msg("Failed to read body while logging request")
					return err
				}
				reqBody = config.RequestBodyLogReplacer(reqBody)
				le.Bytes("req_body", reqBody)
				req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			if config.LogRequestHeader {
				header := zerolog.Dict()
				for k, v := range config.RequestHeaderLogReplacer(req.Header) {
					header.Strs(k, v)
				}
				le.Dict("req_header", header)
			}

			if config.LogRequestQuery {
				query := zerolog.Dict()
				for k, v := range req.URL.Query() {
					query.Strs(k, v)
				}
				le.Dict("req_query", query)
			}

			le.Msg("Request received")

			c.SetRequest(req)

			var resBody bytes.Buffer
			if config.LogResponseBody {
				mw := io.MultiWriter(res.Writer, &resBody)
				res.Writer = &bodyDumpResponseWriter{Writer: mw, ResponseWriter: res.Writer}
			}

			start := time.Now()
			err = next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			ll := util.LogFromEchoContext(c)
			lle := ll.WithLevel(config.Level).
				Dict("res", zerolog.Dict().
					Int("status", res.Status).
					Int64("bytes_out", res.Size).
					TimeDiff("duration_ms", stop, start).
					Err(err),
				)

			if config.LogResponseBody && !config.ResponseBodyLogSkipper(req, res) {
				lle.Bytes("res_body", config.ResponseBodyLogReplacer(resBody.Bytes()))
			}

			if config.LogResponseHeader {
				header := zerolog.Dict()
				for k, v := range config.RequestHeaderLogReplacer(req.Header) {
					header.Strs(k, v)
				}
				lle.Dict("res_header", header)
			}

			lle.Msg("Response sent")

			return nil
		}
	}
}

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
