taskCount, _ := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_COUNT"))
taskIndex, _ := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_INDEX"))

pageCount, err := getLastPage(ctx, client, cfg.Username, cfg.ItemPerPage)
if err != nil {
	return fmt.Errorf("failed to get last page for %s: %w", cfg.Username, err)
}

// see https://go.dev/play/p/fTzUGh2B3hq
pageStart := taskIndex * pageCount / taskCount
pageEnd := ((taskIndex + 1) * pageCount) / taskCount
