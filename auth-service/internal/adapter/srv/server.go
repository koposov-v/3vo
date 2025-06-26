package srv

import (
	"authjwt/internal/init/config"
	"context"
	"net/http"
)

type server struct {
	server *http.Server
	config *config.Config
}

func NewServer(cfg *config.Config, handler http.Handler) Server {
	httpCfg := cfg.HTTP

	return &server{
		config: cfg,
		server: &http.Server{
			Addr:    httpCfg.Host + ":" + httpCfg.Port,
			Handler: handler,
		},
	}
}

func (s *server) Start() error {
	return s.server.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.HTTP.ShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
