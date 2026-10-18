package domain

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID     int
	Input  []LogEntry
	Status TaskStatus
	Result *AnalysisResult
}
