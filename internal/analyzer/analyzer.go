package analyzer

import "log.analizer/internal/domain"

func Analyze(entries []domain.LogEntry) (domain.AnalysisResult, error) {
	if len(entries) == 0 {
		return domain.AnalysisResult{}, domain.ErrEmptyInput
	}

	result := domain.AnalysisResult{}
	errorCount := make(map[string]int)

	result.Total = len(entries)

	for _, entry := range entries {
		switch entry.Level {
		case domain.LevelInfo:
			result.InfoCount++
		case domain.LevelWarn:
			result.WarnCount++
		case domain.LevelError:
			result.ErrorCount++
			errorCount[entry.Message]++

			if errorCount[entry.Message] > errorCount[result.MostCommonError] {
				result.MostCommonError = entry.Message
			}
		}
	}

	return result, nil
}
