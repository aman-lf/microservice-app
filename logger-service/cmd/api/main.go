package main

import (
  "errors"
  "fmt"
  "log"
  "net"
  "net/http"
  "net/rpc"

  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
)

const (
  webPort = "8005"
  rpcPort = "5001"
)

type Config struct {
  Models Models
  router *gin.Engine
}

func main() {
  // Set up application config with in-memory database
  app := Config{
    Models: New(nil), // nil client means we're using in-memory storage
  }

  // Register the RPC server
  if err := rpc.Register(&RPCServer{app: &app}); err != nil {
    log.Panic(err)
  }

  // Start the RPC server in a goroutine
  go app.rpcListen()

  // Set up Gin router with middleware
  router := gin.New()
  router.Use(gin.Recovery())
  router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://*", "http://*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposeHeaders:    []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300,
  }))
  app.router = router

  // Set up routes
  app.setupRoutes()

  // Start the HTTP server
  log.Printf("Starting logger service on port %s\n", webPort)
  srv := &http.Server{
    Addr:    fmt.Sprintf("0.0.0.0:%s", webPort),
    Handler: app.router,
  }

  if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
    log.Fatalf("Failed to listen and serve: %v", err)
  }
}

func (app *Config) rpcListen() {
  log.Println("Starting RPC server on port", rpcPort)

  listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
  if err != nil {
    log.Fatalf("Failed to listen on port %s: %v", rpcPort, err)
  }
  defer listen.Close()

  for {
    rpcConn, err := listen.Accept()
    if err != nil {
      continue
    }
    go rpc.ServeConn(rpcConn)
  }
}

// writeJSON is a helper to write JSON responses
func (app *Config) writeJSON(c *gin.Context, status int, data any, headers ...http.Header) error {
  // Add headers if they exist
  if len(headers) > 0 {
    for key, value := range headers[0] {
      for _, v := range value {
        c.Header(key, v)
      }
    }
  }

  c.JSON(status, data)
  return nil
}

// readJSON tries to read the body of a request and convert it into JSON
func (app *Config) readJSON(c *gin.Context, data any) error {
  if err := c.ShouldBindJSON(data); err != nil {
    return err
  }

  return nil
}

// errorJSON takes an error, and optionally a status code, and generates and sends a JSON error response
func (app *Config) errorJSON(c *gin.Context, err error, status ...int) {
  statusCode := http.StatusBadRequest

  if len(status) > 0 {
    statusCode = status[0]
  }

  var payload struct {
    Error   bool   `json:"error"`
    Message string `json:"message"`
  }

  payload.Error = true
  payload.Message = err.Error()

  c.JSON(statusCode, payload)
}
