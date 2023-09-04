package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v48/github"
	"golang.org/x/exp/slog"
)

func SaveTopic(
	ctx context.Context,
	client *github.Client,
	topic string,
	maxStars int,
	perPage int,
) error {
	filename := fmt.Sprintf("%s-topic.json", topic)
	if err := upsertFile(filename); err != nil {
		return fmt.Errorf("file: %w", err)
	}
	slog.Info(
		"topic repos stored",
		slog.String("filename", filename),
		slog.String("topic", topic))
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: perPage,
		},
	}
	query := fmt.Sprintf("topic:%s stars:<%d", topic, maxStars)

	for {
		result, resp, err := client.Search.Repositories(ctx, query, opts)
		if err != nil {
			var rle *github.RateLimitError
			if ok := errors.As(err, &rle); ok {
				slog.Info(
					"rate limit exceeded",
					slog.Time("sleeping until", rle.Rate.Reset.Time))
				time.Sleep(time.Until(rle.Rate.Reset.Time))
				continue
			}

			var ae *github.AcceptedError
			if ok := errors.As(err, &ae); ok {
				slog.Info("accepted error, scheduled on GitHub side, retrying in 1s")
				time.Sleep(1 * time.Second)
				continue
			}

			return fmt.Errorf("failed to list starred repo: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			slog.Error(
				"request",
				slog.String("status", resp.Status),
				slog.Any("body", resp.Body))
			return fmt.Errorf("status: %s", resp.Status)
		}

		for _, repo := range result.Repositories {
			if err = json.NewEncoder(f).Encode(cleanRepoFromTopic(repo)); err != nil {
				return err
			}
		}

		if err = f.Sync(); err != nil {
			return err
		}

		slog.Info(
			"topic repo",
			slog.Int("page", opts.Page),
			slog.Int("quota remaining", resp.Rate.Remaining))

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
	return nil
}

func upsertFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file %q: %w", path, err)
		}
		_ = f.Close()
		return nil
	}
	defer file.Close()

	return nil
}
