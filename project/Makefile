BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
INVENTORY_BINARY=inventoryApp
MENU_BINARY=menuApp
ORDER_BINARY=orderApp

## up: starts all containers in the background without forcing build
up:
  @echo "Starting Docker images..."
  docker compose up -d
  @echo "Docker images started!"

## up_build: stops docker compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_inventory build_menu build_order
  @echo "Stopping docker images (if running...)"
  docker compose down
  @echo "Building (when required) and starting docker images..."
  docker compose up --build -d
  @echo "Docker images built and started!"

## down: stop docker compose
down:
  @echo "Stopping docker compose..."
  docker compose down
  @echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
  @echo "Building broker binary..."
  cd ../broker-service && go build -o ${BROKER_BINARY} ./cmd/api
  @echo "Broker binary built!"

## build_auth: builds the auth binary as a linux executable
build_auth:
  @echo "Building auth binary..."
  cd ../authentication-service && go build -o ${AUTH_BINARY} ./cmd/api
  @echo "Auth binary built!"

## build_logger: builds the logger binary as a linux executable
build_logger:
  @echo "Building logger binary..."
  cd ../logger-service && go build -o ${LOGGER_BINARY} ./cmd/api
  @echo "Logger binary built!"

## build_inventory: builds the inventory binary as a linux executable
build_inventory:
  @echo "Building inventory binary..."
  cd ../inventory-service && go build -o ${INVENTORY_BINARY} ./cmd/api
  @echo "Inventory binary built!"

## build_menu: builds the menu binary as a linux executable
build_menu:
  @echo "Building menu binary..."
  cd ../menu-service && go build -o ${MENU_BINARY} ./cmd/api
  @echo "Menu binary built!"

## build_order: builds the order binary as a linux executable
build_order:
  @echo "Building order binary..."
  cd ../order-service && go build -o ${ORDER_BINARY} ./cmd/api
  @echo "Order binary built!"