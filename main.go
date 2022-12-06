package main

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/ardanlabs/conf/v3"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

const (
	devVersion = "dev"
	// exitCodeOK    = 0
	exitCodeError = 1
)

var (
	// NOTE: use ldflags to set the version
	gitSha    = devVersion
	buildTime = devVersion // time.Now().UTC().Format(time.RFC3339)

	ErrFieldRequired = errors.New("field is required")
)

func main() {
	batchID := uuid.New().String()

	log := makeLogger(slog.String("batchID", batchID))
	if err := run(log); err != nil {
		slog.Error("run", err)

		os.Exit(exitCodeError)
	}
}

func run(log *slog.Logger) error {
	cfg := struct {
		conf.Version
		Username    string `conf:"default:veggiemonk,short:u,help:github username (default: veggiemonk)"`
		SecretURL   string `conf:"short:s,mask,help:secret manager url"`
		GithubToken string `conf:"short:t,mask,help:github token"`
		ItemPerPage int    `conf:"default:10,short:i,help:github starred repo per page (default: 10)"`
		Info        bool   `conf:"default:false,mask,help:show build info"`
	}{
		Version: conf.Version{
			Build: gitSha,
			Desc:  "Download all starred repositories from GitHub",
		},
	}

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		return fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Info {
		printBuildInfo()
		return nil
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}

	fmt.Println(out)

	if cfg.GithubToken == "" && cfg.SecretURL == "" {
		msg := "github token or secret manager url is not set."
		msg += "one of them is required: %w"
		return fmt.Errorf(msg, ErrFieldRequired)
	}

	if cfg.GithubToken == "" && cfg.SecretURL != "" {
		token, serr := AccessSecretVersion(cfg.SecretURL)
		if serr != nil {
			return fmt.Errorf("failed accessing secret: %w", serr)
		}
		cfg.GithubToken = string(token)
	}

	taskCount, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_COUNT"))
	if err != nil {
		taskCount = 1
	}
	taskIndex, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_INDEX"))
	if err != nil {
		taskIndex = 0
	}
	taskAttempt, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_ATTEMPT"))
	if err != nil {
		taskAttempt = 0
	}

	cre := os.Getenv("CLOUD_RUN_EXECUTION")

	log.Info(
		"starting",
		slog.String("cloud_run_execution", cre),
		slog.String("username", cfg.Username),
		slog.Int("taskCount", taskCount),
		slog.Int("taskIndex", taskIndex),
		slog.Int("taskAttempt", taskAttempt),
		slog.Int("itemPerPage", cfg.ItemPerPage),
	)

	ctx := context.Background()
	client := gitHubClient(ctx, cfg.GithubToken)

	pageCount, err := getLastPage(ctx, client, cfg.Username, cfg.ItemPerPage)
	if err != nil {
		return fmt.Errorf("failed to get last page for %s: %w", cfg.Username, err)
	}

	// see https://go.dev/play/p/fTzUGh2B3hq
	pageStart := taskIndex * pageCount / taskCount
	pageEnd := ((taskIndex + 1) * pageCount) / taskCount

	log.Info(
		"page range",
		slog.Int("pageStart", pageStart),
		slog.Int("pageEnd", pageEnd),
		slog.Int("pageCount", pageCount),
	)

	_, err = getStarred(ctx, client, cfg.Username, cfg.ItemPerPage, pageStart, pageEnd)
	if err != nil {
		slog.Error("failed to get starred repo", err, slog.String("username", cfg.Username))
		return nil
	}

	return nil
}

// AccessSecretVersion accesses the payload for the given secret version if one
// exists. The version can be a version number as a string (e.g. "5") or an
// alias (e.g. "latest").
// taken from https://github.com/GoogleCloudPlatform/golang-samples/blob/main/secretmanager/access_secret_version.go
func AccessSecretVersion(name string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("secret manager error: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("secret manager error: %w", err)
	}

	// WARNING: Do not print the secret in a production environment - this snippet
	// is showing how to access the secret material.
	// fmt.Println("Plaintext:", string(result.Payload.Data))

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return nil, fmt.Errorf("data corruption detected.")
	}

	return result.Payload.Data, nil
}
