package main

import (
  "log"
  "time"
)

// RPCServer is the RPC server for the logger service
type RPCServer struct {
  app *Config
}

// RPCPayload is the type for data we receive from RPC
type RPCPayload struct {
  Name string
  Data string
}

// LogInfo logs information and returns true if successful
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
  // Insert log entry into in-memory storage
  err := r.app.Models.LogEntry.Insert(LogEntry{
    Name:      payload.Name,
    Data:      payload.Data,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  })
  
  if err != nil {
    log.Println("Error inserting log via RPC:", err)
    return err
  }

  *resp = "Processed payload via RPC: " + payload.Name
  return nil
}