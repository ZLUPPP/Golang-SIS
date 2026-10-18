package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"log.analizer/internal/domain"
)

type tokenDefinition struct {
	token   string
	pattern string
}

var tokenDefinitions = []tokenDefinition{
	{token: "%MESSAGE", pattern: `(.+)`},
	{token: "%LEVEL", pattern: `(INFO|WARN|ERROR)`},
	{token: "%Y", pattern: `(\d{4})`},
	{token: "%m", pattern: `(\d{2})`},
	{token: "%d", pattern: `(\d{2})`},
	{token: "%H", pattern: `(\d{2})`},
	{token: "%M", pattern: `(\d{2})`},
	{token: "%S", pattern: `(\d{2})`},
}

type LineParser struct {
	pattern *regexp.Regexp
	fields  []string
	hasTime bool
}

func New(format string) (*LineParser, error) {
	var pattern strings.Builder
	var fields []string

	seen := make(map[string]bool)
	remaining := format

	pattern.WriteString("^")

	for remaining != "" {
		matched := false

		for _, definition := range tokenDefinitions {
			if !strings.HasPrefix(remaining, definition.token) {
				continue
			}

			if seen[definition.token] {
				return nil, fmt.Errorf("token '%s' already defined", definition.token)
			}

			pattern.WriteString(definition.pattern)
			fields = append(fields, definition.token)
			seen[definition.token] = true

			remaining = remaining[len(definition.token):]
			matched = true
			break
		}

		if matched {
			continue
		}

		r, size := utf8.DecodeRuneInString(remaining)
		pattern.WriteString(regexp.QuoteMeta(string(r)))
		remaining = remaining[size:]
	}

	pattern.WriteString("$")

	required := []string{"%Y", "%m", "%d", "%LEVEL", "%MESSAGE"}
	for _, token := range required {
		if !seen[token] {
			return nil, fmt.Errorf(
				"%w: missing token %s",
				domain.ErrLogFormatInvalid,
				token,
			)
		}
	}

	hasAnyTime := seen["%H"] || seen["%M"] || seen["%S"]
	hasFullTime := seen["%H"] && seen["%M"] && seen["%S"]

	if hasAnyTime && !hasFullTime {
		return nil, fmt.Errorf(
			"%w: time requires %%H, %%M and %%S",
			domain.ErrLogFormatInvalid,
		)
	}

	compiled, err := regexp.Compile(pattern.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLogFormatInvalid, err)
	}

	return &LineParser{
		pattern: compiled,
		fields:  fields,
		hasTime: hasFullTime,
	}, nil

}

func (p *LineParser) ParseLine(line string) (domain.LogEntry, error) {
	matches := p.pattern.FindStringSubmatch(strings.TrimSpace(line))
	if matches == nil {
		return domain.LogEntry{}, fmt.Errorf(
			"%w: %q",
			domain.ErrLogFormatInvalid,
			line,
		)
	}

	values := make(map[string]string)

	for i, field := range p.fields {
		values[field] = matches[i+1]
	}

	dateValue := fmt.Sprintf(
		"%s-%s-%s",
		values["%Y"],
		values["%m"],
		values["%d"],
	)
	dateLayout := "2006-01-02"

	if p.hasTime {
		dateValue += fmt.Sprintf(
			" %s:%s:%s",
			values["%H"],
			values["%M"],
			values["%S"],
		)
		dateLayout += " 15:04:05"
	}

	timestamp, err := time.Parse(dateLayout, dateValue)
	if err != nil {
		return domain.LogEntry{}, fmt.Errorf(
			"%w: invalid date: %v",
			domain.ErrLogFormatInvalid,
			err,
		)
	}

	message := strings.TrimSpace(values["%MESSAGE"])
	if message == "" {
		return domain.LogEntry{}, fmt.Errorf(
			"%w: message is empty",
			domain.ErrLogFormatInvalid,
		)
	}

	return domain.LogEntry{
		Timestamp: timestamp,
		Level:     domain.LogLevel(values["%LEVEL"]),
		Message:   message,
	}, nil
}
