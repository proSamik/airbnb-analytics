# Airbnb Room Analytics API

A Go service that provides analytics for Airbnb room bookings, including occupancy rates and price analytics (Mock data)

## Prerequisites

Before running the application, ensure you have the following installed:
* Go 1.23 or higher
* PostgreSQL 12 or higher

## Project Setup

### PostgreSQL Setup

1. **Install PostgreSQL**
   ```bash
   # For Ubuntu/Debian
   sudo apt update
   sudo apt install postgresql postgresql-contrib

   # For MacOS using Homebrew
   brew install postgresql
   ```

2. **Start PostgreSQL Service**
   ```bash
   # For Ubuntu/Debian
   sudo systemctl start postgresql
   sudo systemctl enable postgresql

   # For MacOS
   brew services start postgresql@14
   ```

3. **Configure PostgreSQL**
   ```bash
   # Access PostgreSQL prompt
   sudo -u postgres psql

   # Create user (if needed)
   CREATE USER your_username WITH PASSWORD 'your_password';

   # Grant privileges (optional)
   ALTER USER your_username WITH SUPERUSER;

   # Exit PostgreSQL prompt
   \q
   ```

### Go Setup

1. **Install Go**
   ```bash
   # Download Go from official website
   # For Linux: https://golang.org/dl/
   wget https://go.dev/dl/go1.20.linux-amd64.tar.gz

   # Extract and install
   sudo rm -rf /usr/local/go
   sudo tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz

   # Add to PATH in ~/.bashrc or ~/.zshrc
   export PATH=$PATH:/usr/local/go/bin
   ```

2. **Verify Go Installation**
   ```bash
   go version
   ```

### Project Dependencies

1. **Install Required Go Packages**
   ```bash
   # PostgreSQL driver
   go get github.com/lib/pq

   # Environment variables
   go get github.com/joho/godotenv

   # Router
   go get github.com/gorilla/mux
   ```

### Project Configuration

1. **Clone Repository**
   ```bash
   git clone https://github.com/proSamik/airbnb-analytics
   cd airbnb-analytics
   ```

2. **Create Environment File**
   ```bash
   # Create .env file in project root
   touch .env

   # Add following configurations
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=airbnb_analytics
   ```

3. **Initialize Database**
   ```bash
   # Run database setup script
   go run scripts/db_setup.go
   ```

4. **Verify Setup**
   ```bash
   # Connect to database
   psql -U your_username -d airbnb_analytics

   # Check tables
   \dt

   # Check sample data
   SELECT * FROM room_bookings LIMIT 5;
   ```

### Running the Application

1. **Start the Server**
   ```bash
   go run cmd/api/main.go
   ```

2. **Test the API**
   ```bash
   # Get all room IDs
   curl http://localhost:8080/rooms

   # Get analytics for a specific room
   curl http://localhost:8080/{roomId}
   ```

### Common Issues and Solutions

1. **PostgreSQL Connection Issues**
  - Verify PostgreSQL is running:
    ```bash
    sudo systemctl status postgresql
    ```
  - Check connection settings in .env file
  - Ensure database user has proper permissions

2. **Database Setup Issues**
  - Make sure PostgreSQL user has permission to create databases
  - Check if database already exists:
    ```bash
    psql -U postgres -l
    ```

3. **Go Module Issues**
  - Ensure you're in the project directory
  - Run `go mod tidy` to clean up dependencies
  - Verify go.mod file exists and is correct

## API Usage

### Get All Room IDs
```bash
GET /rooms
```
Example request:
```bash
curl http://localhost:8080/rooms
```

### Get Room Analytics
```bash
GET /{roomId}
```
Example request:
```bash
curl http://localhost:8080/{roomId}
```

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
* PostgreSQL must be running and accessible
* Environment variables must be properly configured
* The API supports CORS for cross-origin requests
* Analytics are calculated for:
  - Occupancy: Next 5 months
  - Rates: Next 30 days

## Error Handling
The API returns appropriate HTTP status codes and error messages:
* 400: Bad Request (invalid room ID)
* 404: Room not found
* 500: Internal server error

Error responses are in JSON format:
```json
{
    "error": "error message"
}
```
---
