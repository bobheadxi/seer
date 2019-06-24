package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"go.bobheadxi.dev/res"
	"go.bobheadxi.dev/seer/jobs"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store"
	"go.bobheadxi.dev/zapx/zhttp"
)

// Server contains Seer's core application logic
type Server struct {
	l   *zap.Logger
	mux *chi.Mux
	srv *http.Server

	jobsEngine jobs.Engine
}

// New instantiates a new server
func New(
	l *zap.Logger,
	riotAPI riot.API,
	backend store.Store,
	jobsEngine jobs.Engine,
) (*Server, error) {
	srv := &Server{
		l:          l,
		srv:        &http.Server{},
		jobsEngine: jobsEngine,
	}

	mux := chi.NewMux()
	mux.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		zhttp.NewMiddleware(l.Named("requests"), nil).Logger,
	)

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) { res.R(w, r, res.MsgOK("server online")) })

	// TODO: hash of names instead of shortid?
	teams := &teamAPI{l.Named("teams"), riotAPI, backend, jobsEngine}
	mux.Route("/team", teams.Group)

	srv.srv.Handler = mux
	return srv, nil
}

// Start spins up the API server
func (s *Server) Start(addr string, stop chan bool) error {
	if addr == "" {
		addr = ":8080"
	}

	s.jobsEngine.Start()

	s.srv.Addr = addr
	go func() {
		<-stop
		s.Stop(context.Background())
	}()
	return s.srv.ListenAndServe()
}

// Stop shuts down this server and its associated resources
func (s *Server) Stop(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.l.Error(err.Error())
	}
	s.jobsEngine.Close()
}
