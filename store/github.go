package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.bobheadxi.dev/seer/config"
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
	repo config.GitHubStoreRepo

	teams *teamsToIDCache

	// versionID *int
}

// NewGitHubStore instantiates a new Store backed by GitHub issues
func NewGitHubStore(ctx context.Context, l *zap.Logger, auth oauth2.TokenSource, repo config.GitHubStoreRepo) (Store, error) {
	c := github.NewClient(oauth2.NewClient(ctx, auth))

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
		repo:  repo,
		// versionID: &versionID,
	}, nil
}

func (g *gitHubStore) Create(ctx context.Context, teamID string, team *Team) error {
	log := g.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))

	// check for existence of team
	if teamIssue, _ := g.teams.getID(ctx, teamID); teamIssue != 0 {
		return fmt.Errorf("generated team ID '%s' already exists - there is probably already a team with the exact same members",
			teamID)
	}

	// generate meta-issue for team
	var labels []string
	for _, m := range team.Members {
		summID := "p=" + m.PlayerID
		labels = append(labels, summID)
		log.Debug("creating label", zap.String("label", summID))
		if _, resp, err := g.c.Issues.CreateLabel(ctx, g.repo.Owner, g.repo.Repo, &github.Label{
			Name:        github.String(summID),
			Color:       github.String("f29513"),
			Description: github.String(m.Name),
		}); err != nil {
			if resp != nil && resp.StatusCode == 422 {
				log.Info("ignoring label creation error", zap.Error(err),
					zap.String("status_message", resp.Status))
				continue
			}
			return err
		}
	}
	labels = append(labels, "team")
	log.Info("labels generated", zap.Strings("labels", labels))

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

	log.Info("team successfully created")
	return nil
}

func (g *gitHubStore) Get(ctx context.Context, teamID string) (*Team, Matches, error) {
	log := g.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))

	teamIssue, err := g.teams.getID(ctx, teamID)
	if err != nil {
		return nil, nil, err
	}

	// get team from issue
	issue, _, err := g.c.Issues.Get(ctx, g.repo.Owner, g.repo.Repo, teamIssue)
	if err != nil {
		return nil, nil, err
	}
	if issue == nil {
		return nil, nil, fmt.Errorf("team %s (%d) not found", teamID, teamIssue)
	}
	var team Team
	if err := json.Unmarshal([]byte(issue.GetBody()), &team); err != nil {
		log.Error("could not read team contents", zap.Error(err),
			zap.Int("issue.number", teamIssue),
			zap.Any("issue.received", issue))
		return nil, nil, err
	}

	// get matches from comments
	// TODO: might need to page eventually
	comments, _, err := g.c.Issues.ListComments(ctx, g.repo.Owner, g.repo.Repo, teamIssue, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, nil, err
	}
	log.Info("found comments", zap.Int("comments", len(comments)))
	matches := make(Matches, 0)
	for _, c := range comments {
		if c.GetUser().GetLogin() != g.repo.Owner {
			log.Debug("skipping comment from unknown user", zap.String("unknown_user", c.GetUser().GetLogin()))
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
	sort.Sort(matches)

	log.Info("team retrieved", zap.Int("matches", len(matches)))
	return &team, matches, nil
}

func (g *gitHubStore) Add(ctx context.Context, teamID string, matches Matches) error {
	log := g.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))

	teamIssue, err := g.teams.getID(ctx, teamID)
	if err != nil {
		return err
	}

	var added int
	sort.Sort(matches)
	for _, m := range matches {
		log.Debug("storing match", zap.Int64("game_id", m.Details.GameID))
		b, err := json.Marshal(m.Details)
		if err != nil {
			log.Error("failed to marshal match details", zap.Error(err))
			continue
		}
		if _, _, err := g.c.Issues.CreateComment(ctx, g.repo.Owner, g.repo.Repo, teamIssue, &github.IssueComment{
			Body: github.String(string(b)),
		}); err != nil {
			log.Error("failed to create comment for match", zap.Error(err))
			return fmt.Errorf("failed to save match '%d': %v", m.Details.GameID, err)
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
	repo config.GitHubStoreRepo
}

func (t *teamsToIDCache) setID(teamID string, ghIssueID int) {
	t.cache.Set(teamID, ghIssueID, cache.DefaultExpiration)
}

func (t *teamsToIDCache) getID(ctx context.Context, teamID string) (int, error) {
	v, found := t.cache.Get(teamID)
	if !found {
		issues, _, err := t.c.Search.Issues(ctx,
			newIssueQuery(t.repo).team(teamID).build(),
			&github.SearchOptions{},
		)
		if err != nil {
			return 0, err
		}
		if issues.GetTotal() == 0 {
			return 0, errors.New("no such team found")
		}
		for _, i := range issues.Issues {
			if strings.Contains(i.GetTitle(), teamID) && !i.IsPullRequest() {
				t.setID(teamID, int(i.GetNumber()))
				v = int(i.GetNumber())
			}
		}
	}
	return v.(int), nil
}

type githubQuery string

func newIssueQuery(repo config.GitHubStoreRepo) *githubQuery {
	q := githubQuery(fmt.Sprintf("user:%s repo:%s is:open", repo.Owner, repo.Repo))
	return &q
}

func (q *githubQuery) team(t string) *githubQuery {
	*q += githubQuery(fmt.Sprintf(" team.id=%s in:title ", t))
	return q
}

func (q *githubQuery) build() string { return string(*q) }
