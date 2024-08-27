package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/rafaeljusto/esprofiler/internal/parser"
)

// staticContent is our static web server content.
//
//go:embed static/*
var staticContent embed.FS
var staticContentHandler = http.FileServer(http.FS(staticContent))

// templateContent is our template web server content.
//
//go:embed templates/*
var templateContent embed.FS
var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"milliseconds": func(took int64) string {
		return (time.Duration(took) * time.Millisecond).String()
	},
	"widthPercentage": func(timeInNanos, took int64) string {
		return fmt.Sprintf("%.3f%%", float64(timeInNanos*100)/float64(took*1e+6))
	},
	"dict": func(values ...any) map[string]any {
		if len(values)%2 != 0 {
			panic("invalid dict call")
		}
		dict := make(map[string]any, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				panic("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict
	},
}).ParseFS(templateContent, "templates/*.html"))

var esClient = http.Client{
	Timeout: 30 * time.Second,
}

// RegisterHandlers registers the handlers for the web server.
func RegisterHandlers(router *http.ServeMux, logger *slog.Logger) {
	router.HandleFunc("/", loggerWrapper(logger, staticContentRerouteHandler))
	router.HandleFunc("/analyze", loggerWrapper(logger, analyzeHandler(logger)))
}

func staticContentRerouteHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "static/" + r.URL.Path
	staticContentHandler.ServeHTTP(w, r)
}

func analyzeHandler(logger *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpLogger := logger.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		server := r.FormValue("server")
		if server == "" {
			http.Error(w, "server is required", http.StatusBadRequest)
			return
		}
		method := r.FormValue("method")
		if method == "" {
			http.Error(w, "method is required", http.StatusBadRequest)
			return
		}
		path := r.FormValue("path")
		if path == "" {
			http.Error(w, "path is required", http.StatusBadRequest)
			return
		}
		query := r.FormValue("query")
		if query == "" {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		}

		server = strings.TrimSuffix(server, "/")
		method = strings.ToUpper(method)
		path = strings.TrimPrefix(path, "/")

		esRequest, err := http.NewRequest(method, server+"/"+path, bytes.NewBufferString(query))
		if err != nil {
			httpLogger.Error("failed to create elasticsearch request",
				slog.String("error", err.Error()),
			)
			http.Error(w, "failed to create elasticsearch request", http.StatusInternalServerError)
			return
		}
		esRequest.Header.Set("Content-Type", "application/json")

		esResponse, err := esClient.Do(esRequest)
		if err != nil {
			httpLogger.Error("failed to send elasticsearch request",
				slog.String("error", err.Error()),
			)
			http.Error(w, "failed to send elasticsearch request", http.StatusInternalServerError)
			return
		}
		defer esResponse.Body.Close()

		if esResponse.StatusCode != http.StatusOK {
			msg := fmt.Sprintf("elasticsearch request failed with status code %d", esResponse.StatusCode)
			http.Error(w, msg, esResponse.StatusCode)
			return
		}

		esResponseParsed, err := parser.Parse(esResponse.Body)
		if err != nil {
			httpLogger.Error("failed to parse elasticsearch response",
				slog.String("error", err.Error()),
			)
			http.Error(w, "failed to parse elasticsearch response", http.StatusInternalServerError)
			return
		}

		httpLogger.Debug("elasticsearch response",
			slog.Int64("took", esResponseParsed.Took),
			slog.Int("shards", len(esResponseParsed.Profile.Shards)),
		)

		if err := templates.ExecuteTemplate(w, "analyze.html", esResponseParsed); err != nil {
			httpLogger.Error("failed to execute template",
				slog.String("error", err.Error()),
			)
			http.Error(w, "failed to execute template", http.StatusInternalServerError)
			return
		}
	}
}

func loggerWrapper(logger *slog.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("request received",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		handler(w, r)
	}
}
