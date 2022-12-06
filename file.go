package main

//
// NOT USED, kept for documentation purpose
//

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v48/github"
	"golang.org/x/exp/slog"
)

// SaveStarred saves the starred repositories to a file.
// It will append the new stars to the file.
// The file will be created if it does not exist.
func SaveStarred(ctx context.Context, client *github.Client, username string, perPage int) error {
	filename := fmt.Sprintf("%s-stars.json", username)
	slog.Info("starred repos stored", slog.String("filename", filename))

	allStars, err := resume(filename)
	if err != nil {
		return fmt.Errorf("failed to resume: %w", err)
	}

	GithubUser(ctx, client, username)

	slog.Info("fetching user", slog.String("username", username))

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			Page:    resumePage(allStars, perPage),
			PerPage: perPage,
		},
		// Sort:        "full_name", // created, updated, pushed, full_name (default).
	}

	// var allStars []*github.StarredRepository
	cpt := len(allStars)

	for {
		stars, resp, err := client.Activity.ListStarred(ctx, username, opts)
		if err != nil {
			var rle *github.RateLimitError
			if ok := errors.As(err, &rle); ok {
				slog.Info("rate limit exceeded", slog.Time("sleeping until", rle.Rate.Reset.Time))
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

		b := new(bytes.Buffer)
		jerr := json.NewEncoder(b).Encode(stars)
		if jerr != nil {
			return jerr
		}

		if _, err := f.Write(b.Bytes()); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		ferr := f.Sync()
		if ferr != nil {
			return ferr
		}

		cpt = cpt + len(stars)

		slog.Info("starred repo", slog.Int("counter", cpt), slog.Int("quota remaining", resp.Rate.Remaining))

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
	return nil
}

func resumePage(stars []*github.StarredRepository, perPage int) int {
	if len(stars) == 0 {
		return 1
	}
	return len(stars)/perPage + 1
}

func resume(path string) ([]*github.StarredRepository, error) {
	var stars []*github.StarredRepository

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create file %q: %w", path, err)
		}
		_ = f.Close()
		return stars, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var s []*github.StarredRepository
		if err := json.Unmarshal(scanner.Bytes(), &s); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)
		}
		stars = append(stars, s...)
	}

	// the error `invalid character '[' after top-level value` occurs
	// because we write one json array per line in a single line
	// so this does not work:
	//
	// content, err := os.ReadFile(path)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read file %q: %w", path, err)
	// }
	// if err := json.Unmarshal(content, &stars); err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal file %q: %w", path, err)
	// }

	return stars, nil
}
