package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v48/github"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
)

func gitHubClient(ctx context.Context, token string) *github.Client {
	// fmt.Print("GitHub Token: ")
	// byteToken, _ := terminal.ReadPassword(int(syscall.Stdin))
	// println()
	// token := string(byteToken)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func GithubUser(ctx context.Context, client *github.Client, u string) {
	user, resp, err := client.Users.Get(ctx, u)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	// Rate.Limit should most likely be 5000 when authorized.
	log.Printf("Rate: %#v\n", resp.Rate)

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	fmt.Printf("\n%v\n", github.Stringify(user))
}

func getLastPage(
	ctx context.Context,
	client *github.Client,
	username string,
	itemPerPage int,
) (int, error) {
	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: itemPerPage,
		},
	}

	_, resp, err := client.Activity.ListStarred(ctx, username, opts)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to query the list of starred repositories: %w",
			err)
	}
	return resp.LastPage, nil
}

func getStarred(
	ctx context.Context,
	client *github.Client,
	username string,
	itemPerPage, pageStart, pageCount int,
) ([]*Repo, error) {
	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			Page:    pageStart,
			PerPage: itemPerPage,
		},
	}

	var starred []*Repo
	for {
		stars, resp, err := client.Activity.ListStarred(ctx, username, opts)
		if err != nil {
			var rle *github.RateLimitError
			if ok := errors.As(err, &rle); ok {
				slog.Info(
					"rate limit exceeded",
					slog.Time("sleeping until", rle.Rate.Reset.Time))
				// TODO: check timezone
				time.Sleep(time.Until(rle.Rate.Reset.Time))
				continue
			}

			var ae *github.AcceptedError
			if ok := errors.As(err, &ae); ok {
				slog.Info("accepted error, scheduled on GitHub side, retrying in 1s")
				// TODO: exponential backoff
				time.Sleep(1 * time.Second)
				continue
			}
			return nil, fmt.Errorf(
				"failed to query the list of starred repositories: %w",
				err)
		}

		for _, star := range stars {
			r := convertGhRepo(star)
			b, jerr := json.Marshal(r)
			if jerr != nil {
				return nil, fmt.Errorf("failed to marshal repo: %w", jerr)
			}
			// Send to a channel to streamline the processing of the data
			// here we just log it for simplicity
			fmt.Printf("%s\n", b)
			starred = append(
				starred,
				r) // not really needed but for the sake of the example
		}

		if resp.NextPage == 0 || resp.NextPage > pageCount {
			break
		}

		opts.Page = resp.NextPage

	}

	return starred, nil
}
