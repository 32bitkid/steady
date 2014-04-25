package main

import "bufio"
import "log"
import "strings"

func ScanGit(options ParserOptions, otherArgs string) Commits {
	commits := make(Commits, 0)
	args := strings.Split(otherArgs, " ")
	logArgs := append(options.GitArgs, args...)
	cmd := ExecGit(logArgs...)

	data, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(data)
	s.Split(options.SplitFunc)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	for s.Scan() {
		commits = append(commits, options.CommitParser(s.Text()))
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	cmd.Wait()

	return commits
}
