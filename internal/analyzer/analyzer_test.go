package analyzer

import (
	"errors"
	"testing"

	"log.analizer/internal/domain"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name    string
		entries []domain.LogEntry
		want    domain.AnalysisResult
	}{
		{
			name: "mixed levels",
			entries: []domain.LogEntry{
				{Level: domain.LevelInfo, Message: "started"},
				{Level: domain.LevelWarn, Message: "slow"},
				{Level: domain.LevelError, Message: "failed"},
				{Level: domain.LevelError, Message: "failed"},
			},
			want: domain.AnalysisResult{
				Total:           4,
				InfoCount:       1,
				WarnCount:       1,
				ErrorCount:      2,
				MostCommonError: "failed",
			},
		},
		{
			name: "only info",
			entries: []domain.LogEntry{
				{Level: domain.LevelInfo, Message: "first"},
				{Level: domain.LevelInfo, Message: "second"},
			},
			want: domain.AnalysisResult{
				Total:     2,
				InfoCount: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Analyze(tt.entries)
			if err != nil {
				t.Fatalf("Analyze() error: %v", err)
			}

			if got != tt.want {
				t.Errorf("Analyze() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestAnalyze_EmptyInput(t *testing.T) {
	_, err := Analyze([]domain.LogEntry{})

	if !errors.Is(err, domain.ErrEmptyInput) {
		t.Errorf("expected ErrEmptyInput, got %v", err)
	}
}
