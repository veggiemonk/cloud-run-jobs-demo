package main

import (
	"time"

	"github.com/google/go-github/v48/github"
)

type Repo struct {
	StarredAt       time.Time       `json:"starred_at"`
	FullName        string          `json:"full_name"`
	Description     string          `json:"description"`
	Homepage        string          `json:"homepage"`
	CreatedAt       time.Time       `json:"created_at"`
	PushedAt        time.Time       `json:"pushed_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Language        string          `json:"language"`
	Fork            bool            `json:"fork"`
	ForksCount      int             `json:"forks_count"`
	OpenIssuesCount int             `json:"open_issues_count"`
	StargazersCount int             `json:"stargazers_count"`
	Size            int             `json:"size"`
	AllowForking    bool            `json:"allow_forking"`
	Topics          []string        `json:"topics"`
	Archived        bool            `json:"archived"`
	Disabled        bool            `json:"disabled"`
	Private         bool            `json:"private"`
	HasIssues       bool            `json:"has_issues"`
	HasWiki         bool            `json:"has_wiki"`
	HasPages        bool            `json:"has_pages"`
	HasProjects     bool            `json:"has_projects"`
	HasDownloads    bool            `json:"has_downloads"`
	IsTemplate      bool            `json:"is_template"`
	License         *github.License `json:"license"`
}

func cleanRepoFromTopic(repo *github.Repository) *Repo {
	if repo == nil {
		return nil
	}
	return &Repo{
		FullName:        repo.GetFullName(),
		Description:     repo.GetDescription(),
		Homepage:        repo.GetHomepage(),
		CreatedAt:       repo.GetCreatedAt().Time,
		PushedAt:        repo.GetPushedAt().Time,
		UpdatedAt:       repo.GetUpdatedAt().Time,
		Language:        repo.GetLanguage(),
		Fork:            repo.GetFork(),
		ForksCount:      repo.GetForksCount(),
		OpenIssuesCount: ptrInt(repo.OpenIssuesCount),
		StargazersCount: ptrInt(repo.StargazersCount),
		Size:            ptrInt(repo.Size),
		AllowForking:    ptrBool(repo.AllowForking),
		Topics:          repo.Topics,
		Archived:        ptrBool(repo.Archived),
		Disabled:        ptrBool(repo.Disabled),
		Private:         ptrBool(repo.Private),
		HasIssues:       ptrBool(repo.HasIssues),
		HasWiki:         ptrBool(repo.HasWiki),
		HasPages:        ptrBool(repo.HasPages),
		HasProjects:     ptrBool(repo.HasProjects),
		HasDownloads:    ptrBool(repo.HasDownloads),
		IsTemplate:      ptrBool(repo.IsTemplate),
		License:         repo.License,
	}
}

func convertGhRepo(star *github.StarredRepository) *Repo {
	r := star.GetRepository()
	if r == nil {
		return nil
	}
	return &Repo{
		StarredAt:       star.StarredAt.Time, // UTC ??  should be modified depending on use case
		FullName:        r.GetFullName(),
		Description:     r.GetDescription(),
		Homepage:        r.GetHomepage(),
		CreatedAt:       r.GetCreatedAt().Time,
		PushedAt:        r.GetPushedAt().Time,
		UpdatedAt:       r.GetUpdatedAt().Time,
		Language:        r.GetLanguage(),
		Fork:            r.GetFork(),
		ForksCount:      r.GetForksCount(),
		OpenIssuesCount: ptrInt(r.OpenIssuesCount),
		StargazersCount: ptrInt(r.StargazersCount),
		Size:            ptrInt(r.Size),
		AllowForking:    ptrBool(r.AllowForking),
		Topics:          r.Topics,
		Archived:        ptrBool(r.Archived),
		Disabled:        ptrBool(r.Disabled),
		Private:         ptrBool(r.Private),
		HasIssues:       ptrBool(r.HasIssues),
		HasWiki:         ptrBool(r.HasWiki),
		HasPages:        ptrBool(r.HasPages),
		HasProjects:     ptrBool(r.HasProjects),
		HasDownloads:    ptrBool(r.HasDownloads),
		IsTemplate:      ptrBool(r.IsTemplate),
		License:         r.License,
	}
}

func ptrInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func ptrBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
