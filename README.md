# Cloudant AirportDB Hub Finder
This is a simple demo application that allows the user to find the nearest transport hubs to a given location. It prompts the user for latitude, longitude and distance in kilometers, and then uses the "https://mikerhodes.cloudant.com/airportdb" database to find the nearest hubs within the specified distance.

# Dependencies
The only external dependency is the [cloudant-go-sdk](https://github.com/IBM/cloudant-go-sdk), which is used to interact with the Cloudant database.

# Usage
## Build and running the application
You can build the application using the following command:
```bash
go build -o hubfinder ./cmd/hubfinder
```
Then you can run the application:
```bash
./hubfinder
```
## Running the application directly with `go run`
Alternatively, you can run the application directly without building it first:
```bash
go run ./cmd/hubfinder
```

# Testing
The application includes unit tests for the distance calculations and the hub finding logic. You can run the tests using the following command:
```bash
go test ./...
```

To ensure the system was tested against a diverse range of inputs and edge cases, a portion of the test data was synthesized using AI. These cases were manually reviewed, verified for accuracy against expected outcomes, and adjusted to guarantee validity.
