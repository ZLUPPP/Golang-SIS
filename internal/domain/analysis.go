package domain

type AnalysisResult struct {
	Total           int
	InfoCount       int
	WarnCount       int
	ErrorCount      int
	MostCommonError string
}
