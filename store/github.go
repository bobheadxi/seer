package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.bobheadxi.dev/seer/riot"

	"github.com/go-chi/chi/middleware"
	"github.com/google/go-github/v26/github"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	cacheTTLMins = 60
)

// TODO: is this sustainable? lots of caching opportunities here I think
type gitHubStore struct {
	l    *zap.Logger
	c    *github.Client
	repo GitHubStoreRepo

	teams *teamsToIDCache

	// versionID *int
}

// GitHubStoreRepo configures where data goes
type GitHubStoreRepo struct {
	Owner string
	Repo  string
}

// NewGitHubStore instantiates a new Store backed by GitHub issues
func NewGitHubStore(ctx context.Context, l *zap.Logger, auth oauth2.TokenSource, repo GitHubStoreRepo) (Store, error) {
	c := github.NewClient(oauth2.NewClient(ctx, auth))

	md, _, err := c.APIMeta(ctx)
	if err != nil {
		return nil, err
	}
	l.Info("connection established with GitHub", zap.Any("githubapi.metadata", md))

	lims, _, err := c.RateLimits(ctx)
	if err != nil {
		return nil, err
	}
	l.Info("rate limit check ok", zap.Any("githubapi.ratelimits", lims))

	/*
		versions, _, err := c.Issues.ListMilestones(ctx, repo.Owner, repo.Repo, &github.MilestoneListOptions{})
		if err != nil {
			return nil, err
		}
		var versionID int
		var found bool
		for _, v := range versions {
			if v.GetTitle() == schemaVersion {
				found = true
				versionID = int(v.GetID())
			}
		}
		if !found {
			milestone, _, err := c.Issues.CreateMilestone(ctx, repo.Owner, repo.Repo, &github.Milestone{
				Title:       github.String("schema:" + schemaVersion),
				Description: github.String("metadata milestone for indicating schema version of an issue"),
			})
			if err != nil {
				return nil, err
			}
			versionID = int(milestone.GetID())
		}
	*/

	return &gitHubStore{
		l:     l,
		c:     c,
		teams: &teamsToIDCache{cache.New(cacheTTLMins*time.Minute, cacheTTLMins*time.Minute), c, repo},
		// versionID: &versionID,
	}, nil
}

func (g *gitHubStore) Create(ctx context.Context, teamID string, team *Team) error {
	log := g.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))

	// generate meta-issue for team
	var labels []string
	for _, m := range team.Members {
		summID := "summoner.id=" + m.AccountID
		labels = append(labels, summID)
		if _, _, err := g.c.Issues.CreateLabel(ctx, g.repo.Owner, g.repo.Repo, &github.Label{
			Name: github.String(summID),
		}); err != nil {
			return err
		}
	}
	labels = append(labels, "team")

	// generate body
	body, err := json.Marshal(team)
	if err != nil {
		return err
	}

	// generate id for team
	teamIDLabel := "team.id=" + teamID
	if _, _, err := g.c.Issues.Create(ctx, g.repo.Owner, g.repo.Repo, &github.IssueRequest{
		Title:  github.String(teamIDLabel),
		Labels: &labels,
		Body:   github.String(string(body)),
		// Milestone: g.versionID,
	}); err != nil {
		return err
	}

	log.Info("team successfully created", zap.Strings("labels", labels))
	return nil
}

func (g *gitHubStore) Get(ctx context.Context, teamID string) (*Team, []MatchData, error) {
	log := g.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))

	teamIssue, err := g.teams.getID(ctx, teamID)
	if err != nil {
		return nil, nil, err
	}

	// get team from issue
	issue, _, err := g.c.Issues.Get(ctx, g.repo.Owner, g.repo.Repo, teamIssue)
	if err == nil {
		return nil, nil, err
	}
	var team Team
	if err := json.Unmarshal([]byte(issue.GetBody()), &team); err != nil {
		return nil, nil, err
	}

	// get matches from comments
	comments, _, err := g.c.Issues.ListComments(ctx, g.repo.Owner, g.repo.Repo, teamIssue, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, nil, err
	}
	var matches []MatchData
	for _, c := range comments {
		if c.GetUser().GetName() != g.repo.Owner {
			continue
		}
		var details riot.MatchDetails
		if err := json.Unmarshal([]byte(c.GetBody()), &details); err != nil {
			log.Error("failed to unmarshal comment",
				zap.Error(err),
				zap.Int64("comment.id", c.GetID()),
				zap.String("comment.url", c.GetURL()))
			continue
		}
		matches = append(matches, MatchData{
			Details: &details,
		})
	}

	log.Info("team retrieved")
	return &team, matches, nil
}

func (g *gitHubStore) Add(ctx context.Context, teamID string, matches []MatchData) error {
	log := g.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))

	teamIssue, err := g.teams.getID(ctx, teamID)
	if err != nil {
		return err
	}

	var added int
	for _, m := range matches {
		b, err := json.Marshal(m.Details)
		if err != nil {
			log.Error("failed to marshal match details", zap.Error(err))
			continue
		}
		if _, _, err := g.c.Issues.CreateComment(ctx, g.repo.Owner, g.repo.Repo, teamIssue, &github.IssueComment{
			Body: github.String(string(b)),
		}); err != nil {
			return err
		}
		added++
	}

	log.Info("matches added", zap.Int("matches.count", added))
	return nil
}

// TODO: need?
func (g *gitHubStore) LastUpdated(ctx context.Context, teamID string) (*time.Time, error) {
	return nil, nil
}

type teamsToIDCache struct {
	cache *cache.Cache

	c    *github.Client
	repo GitHubStoreRepo
}

func (t *teamsToIDCache) setID(teamID string, ghIssueID int) {
	t.cache.Set(teamID, ghIssueID, cache.DefaultExpiration)
}

func (t *teamsToIDCache) getID(ctx context.Context, teamID string) (int, error) {
	v, found := t.cache.Get(teamID)
	if !found {
		issues, _, err := t.c.Search.Issues(ctx, newQuery(t.repo).team(teamID).build(), &github.SearchOptions{})
		if err == nil {
			return 0, err
		}
		if issues.GetTotal() == 0 {
			return 0, errors.New("no such team found")
		}
		for _, i := range issues.Issues {
			if strings.Contains(i.GetTitle(), teamID) {
				t.setID(teamID, int(i.GetID()))
				v = int(i.GetID())
			}
		}
	}
	return v.(int), nil
}

type githubQuery string

func newQuery(repo GitHubStoreRepo) *githubQuery {
	q := githubQuery(fmt.Sprintf("user:%s repo:%s", repo.Owner, repo.Repo))
	return &q
}

func (q *githubQuery) team(t string) *githubQuery {
	*q += githubQuery(fmt.Sprintf(" team.id=%s in:title ", t))
	return q
}

func (q *githubQuery) build() string { return string(*q) }
