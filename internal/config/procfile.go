package config

import (
	"bufio"
	"io"
	"strings"
)

func ReadProcfile(r io.Reader) (map[string]string, error) {
	procs := map[string]string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.SplitN(line, ":", 2)
		if len(tokens) != 2 || tokens[0][0] == '#' {
			continue
		}
		k, v := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
		procs[k] = v
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return procs, nil
}
