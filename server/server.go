package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/jobs"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store"
	"go.bobheadxi.dev/zapx/zhttp"
)

// Server contains Seer's core application logic
type Server struct {
	l   *zap.Logger
	mux *chi.Mux

	jobEngine *jobs.Engine
}

// New instantiates a new server
func New(
	l *zap.Logger,
	riotAPI riot.API,
	backend store.Store,
	jobsEngine jobs.Engine,
) (*Server, error) {
	srv := &Server{
		l:   l,
		mux: &chi.Mux{},
	}

	srv.mux.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		zhttp.NewMiddleware(l.Named("requests"), nil).Logger,
	)

	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return nil, err
	}
	srv.mux.Route("/team", teamAPI{l.Named("teams"), riotAPI, backend, jobsEngine, sid}.Group)

	return srv, nil
}

// Start spins up the API server
func (s *Server) Start(addr string) error {
	if addr == "" {
		addr = ":8080"
	}

	// TODO background jobs here?

	serve := &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
	return serve.ListenAndServe()
}

func read(r *http.Request, out interface{}) error {
	body, err := r.GetBody()
	if err != nil {
		return err
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}
