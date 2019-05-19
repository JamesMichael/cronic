package cronic

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type crontab struct {
	path        string
	environment []string
	commands    []command
	error       error
}

func newCrontabFromFile(path string) (*crontab, error) {
	file, err := os.Open(path)
	if err != nil {
		return &crontab{
			path:        path,
			commands:    []command{},
			environment: []string{},
			error:       err,
		}, nil
	}
	defer file.Close()

	var environment []string
	var commands []command

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		if isSkippableLine(line) {
			continue
		}

		if isCommandLine(line) {
			commands = append(commands, splitCommandLine(string(line)))
			continue
		}

		if isEnvironmentLine(line) {
			environment = append(environment, string(line))
			continue
		}

		return nil, fmt.Errorf("Unable to parse file '%s', line: '%s'", path, line)
	}

	return &crontab{
		path:        path,
		environment: environment,
		commands:    commands,
		error:       nil,
	}, nil
}

func (c crontab) Commands() []command {
	return c.commands
}

func (c crontab) Environment() []string {
	return c.environment
}

func (c crontab) Error() error {
	return c.error
}

func (c crontab) WriteCommand(cmd string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "cronic")
	if err != nil {
		return "", err
	}

	const scriptTemplate = `#!/usr/bin/env bash

# environment
{{range .Environment}}
export {{ . }}
{{end}}

# command, pipe output to syslog
{
    {{ .Command }}
} 2>&1 | logger -t cronic

# attempt to clean up after ourselves
rm "$0" || true
`

	scriptContents := struct {
		Environment []string
		Command     string
	}{
		c.Environment(),
		cmd,
	}

	t := template.Must(template.New("script").Parse(scriptTemplate))
	var script bytes.Buffer
	if err := t.Execute(&script, scriptContents); err != nil {
		return "", err
	}

	if _, err = tmpfile.WriteString(script.String()); err != nil {
		os.Remove(tmpfile.Name())
		return "", err
	}

	tmpfile.Close()

	if err = os.Chmod(tmpfile.Name(), 0755); err != nil {
		os.Remove(tmpfile.Name())
		return "", err
	}

	return tmpfile.Name(), nil
}

func splitCommandLine(line string) command {
	r := regexp.MustCompile(`\s+`)

	var parts []string
	if strings.HasPrefix(line, "@") {
		parts = r.Split(line, 3)
	} else {
		words := r.Split(line, 7)
		parts = []string{
			words[0] + " " + words[1] + " " + words[2] + " " + words[3] + " " + words[4],
			words[5],
			words[6],
		}
	}

	return command{
		time:    parts[0],
		user:    parts[1],
		command: parts[2],
	}
}

func isBlankLine(line []byte) bool {
	matched, err := regexp.Match(`^\s*$`, line)
	if err != nil {
		panic(err)
	}

	return matched
}

func isCommandLine(line []byte) bool {
	matched, err := regexp.Match(`^\s*(?:\*|@|\d)`, line)
	if err != nil {
		panic(err)
	}

	return matched
}

func isCommentLine(line []byte) bool {
	matched, err := regexp.Match(`^\s*#`, line)
	if err != nil {
		panic(err)
	}

	return matched
}

func isEnvironmentLine(line []byte) bool {
	matched, err := regexp.Match(`^\s*[^=]+=`, line)
	if err != nil {
		panic(err)
	}

	return matched
}

func isSkippableLine(line []byte) bool {
	return isBlankLine(line) || isCommentLine(line)
}
