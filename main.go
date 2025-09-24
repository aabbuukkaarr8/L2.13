package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	fields    []int
	delimiter string
	separated bool
}

func main() {
	cfg := parseFlags()

	lines, err := readLines()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}

	for _, line := range lines {
		result := cutLine(line, cfg)
		if result != "" {
			fmt.Println(result)
		}
	}
}

func parseFlags() Config {
	var f, d string
	var s bool
	flag.StringVar(&f, "f", "", "fields to extract (e.g., 1,3-5)")
	flag.StringVar(&d, "d", "\t", "delimiter (default tab)")
	flag.BoolVar(&s, "s", false, "only lines with delimiter")
	flag.Parse()

	if f == "" {
		fmt.Fprintln(os.Stderr, "fields (-f) must be specified")
		os.Exit(1)
	}

	fields := parseFields(f)
	return Config{
		fields:    fields,
		delimiter: d,
		separated: s,
	}
}

func readLines() ([]string, error) {
	if flag.NArg() == 0 {
		var lines []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return lines, scanner.Err()
	}

	var lines []string
	for _, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	return lines, nil
}

func cutLine(line string, cfg Config) string {
	if cfg.separated && !strings.Contains(line, cfg.delimiter) {
		return ""
	}

	cols := strings.Split(line, cfg.delimiter)
	var result []string
	for _, idx := range cfg.fields {
		if idx <= 0 || idx > len(cols) {
			continue
		}
		result = append(result, cols[idx-1])
	}
	return strings.Join(result, cfg.delimiter)
}

func parseFields(f string) []int {
	var fields []int
	parts := strings.Split(f, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(rangeParts[0])
			end, err2 := strconv.Atoi(rangeParts[1])
			if err1 != nil || err2 != nil || start > end {
				continue
			}
			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			idx, err := strconv.Atoi(part)
			if err != nil {
				continue
			}
			fields = append(fields, idx)
		}
	}
	return fields
}
