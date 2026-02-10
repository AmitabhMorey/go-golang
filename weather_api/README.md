# Weather API Wrapper Service

A Go-based weather API wrapper that fetches weather data from Visual Crossing API with Redis caching support.

## Features

- Fetches real-time weather data from Visual Crossing API
- Redis-based caching to reduce API calls and improve performance
- Configurable cache expiration (default: 12 hours)
- Environment variable configuration
- Clean JSON API responses
- Health check endpoint

## Prerequisites

- Go 1.21 or higher
- Redis server (optional, but recommended for caching)
- Visual Crossing API key (free tier available at https://www.visualcrossing.com/)

## Installation

1. Clone the repository and navigate to the project directory:
```bash
cd weather_api
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file from the example:
```bash
cp .env.example .env
```

4. Edit `.env` and add your Visual Crossing API key:
```
WEATHER_API_KEY=your_actual_api_key_here
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_EXPIRATION=43200
PORT=8080
```

## Running the Application

### Start Redis (if using caching)
```bash
redis-server
```

### Run the application
```bash
go run main.go
```

Or build and run:
```bash
go build -o weather-api
./weather-api
```

## API Endpoints

### Get Weather Data
```
GET /weather?city={city_name}
```

Example:
```bash
curl "http://localhost:8080/weather?city=London"
```

Response:
```json
{
  "location": "London, England, United Kingdom",
  "temperature": 15.5,
  "conditions": "Partly cloudy",
  "humidity": 72.0,
  "windSpeed": 12.5,
  "timestamp": "2026-02-11T10:30:00Z"
}
```

### Health Check
```
GET /health
```

Example:
```bash
curl "http://localhost:8080/health"
```

Response:
```json
{
  "status": "ok"
}
```

## Configuration

Environment variables:

- `WEATHER_API_KEY` (required): Your Visual Crossing API key
- `REDIS_ADDR` (optional): Redis server address (default: localhost:6379)
- `REDIS_PASSWORD` (optional): Redis password if required
- `REDIS_DB` (optional): Redis database number (default: 0)
- `CACHE_EXPIRATION` (optional): Cache expiration time in seconds (default: 43200 = 12 hours)
- `PORT` (optional): Server port (default: 8080)

## How It Works

1. Client requests weather data for a city
2. Server checks Redis cache for existing data
3. If cache hit: Returns cached data immediately
4. If cache miss: Fetches fresh data from Visual Crossing API
5. Stores the response in Redis with expiration time
6. Returns weather data to client

## Caching Strategy

- Cache key: City name (as provided by user)
- Cache expiration: Configurable (default 12 hours)
- Automatic cache cleanup via Redis TTL
- Graceful degradation if Redis is unavailable

## Error Handling

- Returns appropriate HTTP status codes
- Validates required parameters
- Handles API failures gracefully
- Logs errors for debugging

## Project Structure

```
weather_api/
├── main.go           # Main application code
├── go.mod            # Go module dependencies
├── .env.example      # Example environment configuration
├── .gitignore        # Git ignore rules
└── README.md         # This file
```

## Testing

Test the API with different cities:
```bash
curl "http://localhost:8080/weather?city=New%20York"
curl "http://localhost:8080/weather?city=Tokyo"
curl "http://localhost:8080/weather?city=Paris"
```

## License

MIT License
