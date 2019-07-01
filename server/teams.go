package server

import (
	"net/http"

	"go.bobheadxi.dev/seer/jobs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"go.bobheadxi.dev/res"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store"
)

type teamAPI struct {
	l *zap.Logger

	riotAPI   riot.API
	backend   store.Store
	jobEngine jobs.Engine
}

func (t *teamAPI) Group(r chi.Router) {
	r.Post("/", t.postTeam)
	r.Get("/{teamID}", t.getTeam)
	r.Post("/update/{teamID}", t.postUpdateTeam)
}

type newTeamRequest struct {
	Region  riot.Region `json:"region"`
	Members []string    `json:"members"`
}

func (t *teamAPI) postTeam(w http.ResponseWriter, r *http.Request) {
	var n newTeamRequest
	if err := read(r, &n); err != nil {
		res.R(w, r, res.ErrBadRequest(err.Error()))
		return
	}
	if len(n.Members) < 5 {
		res.R(w, r, res.ErrBadRequest("insufficent members in this team (minimum: 5)"))
		return
	}
	if len(n.Members) > 10 {
		res.R(w, r, res.ErrBadRequest("too many members in team (maximum: 10)"))
		return
	}

	log := t.l.With(zap.String("request.id", middleware.GetReqID(r.Context())))
	log.Info("preparing to create new team", zap.Any("request.new_team", n))

	// look for team members
	team := &store.Team{
		Region:  n.Region,
		Members: make([]*riot.Summoner, len(n.Members)),
	}
	riotAPI := t.riotAPI.WithRegion(n.Region)
	for i, name := range n.Members {
		s, err := riotAPI.Summoner(r.Context(), name)
		if err != nil { // TODO: better responses
			log.Error("failed to find summoner", zap.Error(err))
			res.R(w, r, res.ErrInternalServer("failed to find summoner", err,
				"summoner", name))
			return
		}
		team.Members[i] = s
	}

	teamID := team.GenerateTeamID()
	log = log.With(zap.String("team.id", teamID))
	log.Info("discovered team", zap.Any("team.members", team.Members))

	// commit to datastore
	if err := t.backend.Create(r.Context(), teamID, team); err != nil {
		log.Error("failed to store team data", zap.Error(err))
		res.R(w, r, res.ErrInternalServer("failed to store team data", err))
		return
	}

	log.Info("team stored")
	res.R(w, r, res.MsgOK("team created",
		"team.id", teamID))
}

func (t *teamAPI) getTeam(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "teamID")
	log := t.l.With(
		zap.String("request.id", middleware.GetReqID(r.Context())),
		zap.String("team.id", teamID))

	team, matches, err := t.backend.Get(r.Context(), teamID)
	if err != nil { // TODO better responses
		log.Error("failed to find team", zap.Error(err))
		res.R(w, r, res.ErrInternalServer("failed to find team", err))
		return
	}

	log.Info("team found",
		zap.Int("members", len(team.Members)),
		zap.Int("matches", len(matches)))
	res.R(w, r, res.MsgOK("team found",
		"team", team,
		"matches", matches))
}

func (t *teamAPI) postUpdateTeam(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "teamID")
	requestID := middleware.GetReqID(r.Context())
	log := t.l.With(zap.String("team.id", teamID), zap.String("request.id", requestID))

	// TODO: check for last updated?

	jobID, err := t.jobEngine.Queue(jobs.NewMatchesSyncJob(teamID, requestID))
	if err != nil {
		log.Error("failed to queue team update", zap.Error(err))
		res.R(w, r, res.ErrInternalServer("failed to queue team update", err))
		return
	}
	log.Info("update queued", zap.String("job.id", jobID))

	res.R(w, r, res.MsgOK("team update queued",
		"job_id", jobID))
}
