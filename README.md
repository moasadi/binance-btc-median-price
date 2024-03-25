# Binance Trade App

This application connects to the Binance WebSocket API, calculates the median trade price, and sends it to the frontend via WebSocket.

### Usage
Once the application is running, you can access the frontend by opening a web browser and navigating to http://localhost:8080.


### Prerequisites

- Go 1.20 or later
- A `.env` file in the root directory with the following environment variable:
  - `BINANCE_WS_URL`: The URL of the Binance WebSocket API

