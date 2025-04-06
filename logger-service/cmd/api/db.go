package main

import (
  "fmt"
  "log"
  "sync"
  "time"
)

const (
  // timeout for database operations
  dbTimeout = time.Second * 15
)

// Models is the wrapper for database collections
type Models struct {
  LogEntry LogEntryModel
}

// New is the function used to create an instance of the Models struct
func New(mongoClient interface{}) Models {
  return Models{
    LogEntry: LogEntryModel{
      Logs: make([]LogEntry, 0),
    },
  }
}

// LogEntry is the structure for log entries
type LogEntry struct {
  ID        string    `json:"id,omitempty"`
  Name      string    `json:"name"`
  Data      string    `json:"data"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

// LogEntryModel wraps the in-memory log storage
type LogEntryModel struct {
  Logs  []LogEntry
  mu    sync.Mutex
  count int
}

// Insert adds a new log entry to the in-memory storage
func (l *LogEntryModel) Insert(entry LogEntry) error {
  l.mu.Lock()
  defer l.mu.Unlock()
  
  l.count++
  entry.ID = fmt.Sprintf("%d", l.count)
  entry.CreatedAt = time.Now()
  entry.UpdatedAt = time.Now()
  
  l.Logs = append(l.Logs, entry)
  log.Printf("Added log entry #%s: %s\n", entry.ID, entry.Name)
  
  return nil
}

// GetAll returns all log entries sorted by creation time in descending order
func (l *LogEntryModel) GetAll() ([]LogEntry, error) {
  l.mu.Lock()
  defer l.mu.Unlock()
  
  // Create a copy of the logs slice to avoid race conditions
  logs := make([]LogEntry, len(l.Logs))
  copy(logs, l.Logs)
  
  // Sort in descending order by creation time (newest first)
  // We could implement sorting here, but for simplicity we'll assume
  // the logs are already in order with newest appended at the end
  
  // Reverse the slice to get newest first
  for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
    logs[i], logs[j] = logs[j], logs[i]
  }
  
  return logs, nil
}

// GetOne returns a single log entry by ID
func (l *LogEntryModel) GetOne(id string) (*LogEntry, error) {
  l.mu.Lock()
  defer l.mu.Unlock()
  
  for _, entry := range l.Logs {
    if entry.ID == id {
      return &entry, nil
    }
  }
  
  return nil, fmt.Errorf("log entry with ID %s not found", id)
}

// DropCollection clears all logs
func (l *LogEntryModel) DropCollection() error {
  l.mu.Lock()
  defer l.mu.Unlock()
  
  l.Logs = make([]LogEntry, 0)
  l.count = 0
  
  return nil
}
