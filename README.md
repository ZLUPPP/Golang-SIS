# Log Analyzer SIS

A console application written in Go for analyzing log entries.

## Features

The application can:

- create log analysis tasks;
- assign a unique ID to each task;
- store multiple tasks in memory;
- retrieve tasks by ID;
- count `INFO`, `WARN`, and `ERROR` log entries;
- find the most common error message;
- store the analysis result;
- track the current task status;
- reject empty input and invalid operations.

## Task Lifecycle

Each task goes through the following states:

`pending -> processing -> completed`

- `pending` — the task has been created and is waiting to be processed;
- `processing` — the task is currently being processed;
- `completed` — processing has finished and the result has been stored.

Only tasks with the `pending` status can be processed.

Attempting to process a completed task again returns an error.

## Project Structure

```text
cmd/sis/              Application entry point
internal/domain/      Domain models and application errors
internal/analyzer/    Log analysis logic
internal/application/ Application service
internal/storage/     In-memory task repository
```

## Log File Format

The application reads log entries from a text file. Each log entry must be placed on a separate line.

The default format is:

```text
%Y-%m-%d %LEVEL %MESSAGE
```

Example:

```text
2026-10-18 INFO application started
2026-10-18 WARN response is slow
2026-10-18 ERROR connection failed
```

### Supported Format Tokens

| Token | Description |
|---|---|
| `%Y` | Four-digit year |
| `%m` | Two-digit month |
| `%d` | Two-digit day |
| `%H` | Two-digit hour |
| `%M` | Two-digit minute |
| `%S` | Two-digit second |
| `%LEVEL` | Log level: `INFO`, `WARN`, or `ERROR` |
| `%MESSAGE` | Log message |

The date, level, and message tokens are required.

If the time is included, `%H`, `%M`, and `%S` must all be present.

## Command-Line Usage

Run the application with the default format:

```powershell
go run ./cmd/sis -p "test.log"
```

Also works with pipes

```powershell
go build -o la ./cmd/sis
cat test.log | ./la
```

The `-p` flag can be omitted when input is provided through stdin:

```bash
./la < test.log
```

**Use a custom format:**

```powershell
go run ./cmd/sis -p "test.log" -f "%Y-%m-%d [%LEVEL] %MESSAGE"
```

or

```powershell
cat test.log | la -f "%Y-%m-%d [%LEVEL] %MESSAGE"
la -p - -f "%Y-%m-%d [%LEVEL] %MESSAGE" < test.log
```

Display help:

```powershell
go run ./cmd/sis -h
```

## Building

Build the executable:

```powershell
go build -o sis.exe ./cmd/sis
```

Run the compiled application:

```powershell
.\sis.exe -p "test.log"
```

The compiled executable does not require `go run`. All application packages are included in the executable.

## Project Startup

```powershell
go run ./cmd/sis
```

## Project Verification

```powershell
gofmt -l .
go vet ./...
go test ./...
go build ./...
```