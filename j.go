package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func Diff(t []*Repo, starred []*Repo) []*Repo {
	var c []*Repo
	for _, i := range t {
		if !contains(starred, i) {
			c = append(c, i)
		}
	}
	return c
}

func contains(a []*Repo, b *Repo) bool {
	for _, i := range a {
		if i.FullName == b.FullName {
			return true
		}
	}
	return false
}

func LoadRepos(path string) ([]*Repo, error) {
	var repos []*Repo
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var s *Repo
		if err := json.Unmarshal(scanner.Bytes(), &s); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)
		}
		repos = append(repos, s)
	}
	return repos, nil
}
