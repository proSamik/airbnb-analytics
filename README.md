# Airbnb Room Analytics API

A Go service that provides analytics for Airbnb room bookings, including occupancy rates and price analytics.

## Prerequisites

Before running the application, ensure you have the following installed:

* Go 1.20 or higher
* Node.js and npm
* json-server (Install globally: `npm install -g json-server`)

## Setup & Running

### 1. Clone the Repository

```bash
git clone <repository-url>
cd <repository-name>
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Start Mock API Server

Open a new terminal and run:
```bash
cd mock_data
json-server --watch db.json --port 3001
```

### 4. Start API Server

In a separate terminal, from the project root:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## API Usage

### Get Room Analytics

```bash
GET /{roomId}
```

Example request:
```bash
curl http://localhost:8080/{roomId}
```

Replace `{roomId}` with a valid room ID from `mock_data/README.md`.

### Response Format

```json
{
    "room_id": "string",
    "monthly_occupancy": [
        {
            "month": "YYYY-MM",
            "occupancy_percentage": 85.5
        }
    ],
    "rate_analytics": {
        "average_rate": 150.00,
        "highest_rate": 200.00,
        "lowest_rate": 100.00
    }
}
```

## Important Notes

* Mock API server must be running on port 3001 before starting the main server
* Main server runs on port 8080
* Check  [mock_data/README.md](mock_data/README.md) for valid room IDs
* The API supports CORS for cross-origin requests

## Error Handling

The API returns appropriate HTTP status codes and error messages:

* 404: Room not found
* 500: Internal server error
* Error responses are in JSON format:
  ```json
  {
      "error": "error message"
  }
  ```
---

## Assumptions

1. **Data Source**
- The room occupancy and daily rate data is assumed to be provided by Airbnb's API
- A mock API server has been implemented to simulate this data source

2. **Data Format**
- Room rates are stored as floating-point values
- Dates in the data are in "YYYY-MM-DD" format

3. **Data Continuity**
- The dataset contains continuous dates without any gaps
- Each room has complete booking information for all consecutive dates in the range

---

## Development Challenges & Solutions

### 1. Data Source Design
**Challenge:**
- Needed to determine appropriate data structure and source for room analytics
- Had to decide the starting point for the data timeline

**Solution:**
- Implemented a mock JSON API to simulate Airbnb's data structure
- Used current date as the reference point for all calculations
- Structured data with daily room rates and booking status

### 2. Rate Analytics Implementation
**Challenge:**
- Initially calculated rates for the entire dataset instead of the next 30 days
- Needed to ensure accurate rate calculations within the specified timeframe

**Solution:**
- Implemented date filtering to only consider next 30 days from current date
- Added validation to ensure rate calculations only use relevant data points
- Used floating-point values for precise rate calculations

### 3. Occupancy Rate Calculations
**Challenge:**
- Complex logic required for calculating monthly occupancy rates
- Needed to handle edge cases at month boundaries
- Required consistent calculation method across different months

**Solution:**
- Implemented a rolling 5-month window starting from current month
- Created a structured approach to aggregate daily bookings into monthly statistics
- Added date validation to ensure accuracy of occupancy calculations
- Implemented sorting to ensure consistent month-wise data presentation

---