<p align="center"> 
	<img src="https://github.com/user-attachments/assets/09405fba-3c4a-4dc7-ad4c-e4977360a39a" alt="LambdaMux logo">
</p>
<p align="center">
  <a href="https://github.com/D-Andreev/lambdamux/actions/workflows/test.yml">
    <img src="https://github.com/D-Andreev/lambdamux/actions/workflows/test.yml/badge.svg" alt="Test">
  </a>
  <a href="https://godoc.org/github.com/D-Andreev/lambdamux">
    <img src="https://godoc.org/github.com/D-Andreev/lambdamux?status.svg" alt="GoDoc">
  </a>
</p>

A simple and lightweight high performance HTTP router specifically designed for AWS Lambda functions handling API Gateway requests. 

## Features
- Fast and efficient routing for static routes
- Seamless handling of path parameters in routes (e.g. `/users/:id`)
- Simple and intuitive API for easy integration 

## Motivation
When deploying REST APIs on AWS Lambda, a common approach is to use one of the popular HTTP routers with a proxy like [aws-lambda-go-api-proxy](https://github.com/awslabs/aws-lambda-go-api-proxy). While this leverages existing routers, it introduces overhead by converting the `APIGatewayProxyRequest` event to `http.Request`, so that it can be processed by the router.

`lambdamux` is designed to work directly with API Gateway events, offering several advantages:
  
1.  **Efficient Storage**: Uses a radix tree to compactly store routes, reducing memory usage.
2.  **Fast Matching**: Achieves O(m) time complexity for route matching, where m is the url length.
3.  **No Conversion Overhead**: Processes API Gateway events directly, eliminating request conversion time.

These features make `lambdamux` ideal for serverless environments, optimizing both memory usage and execution time - crucial factors in Lambda performance and cost.

**Note**: While `lambdamux` offers performance improvements, it's important to recognize that HTTP routers are typically not the main bottleneck in API performance. If you're looking to drastically improve your application's performance, you should primarily focus on optimizing database operations, external network calls, and other potentially time-consuming operations within your application logic. The router's performance becomes more significant in high-throughput scenarios or when dealing with very large numbers of routes.

## Installation

```
go get github.com/D-Andreev/lambdamux
```

## Usage

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/D-Andreev/lambdamux"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	router := lambdamux.NewLambdaMux()

	// GET request
	router.GET("/users", listUsers)

	// GET request with path parameter
	router.GET("/users/:id", getUser)

	// POST request
	router.POST("/users", createUser)

	// PUT request with path parameter
	router.PUT("/users/:id", updateUser)

	// DELETE request with path parameter
	router.DELETE("/users/:id", deleteUser)

	lambda.Start(router.Handle)
}

func listUsers(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	users := []string{"Alice", "Bob", "Charlie"}
	body, _ := json.Marshal(users)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func getUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("User details for ID: %s", userID),
	}, nil
}

func createUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Here you would typically parse the request body and create a user
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       "User created successfully",
	}, nil
}

func updateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("User %s updated successfully", userID),
	}, nil
}

func deleteUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("User %s deleted successfully", userID),
	}, nil
}

```

## Running the Examples

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop) installed and running
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) installed
- [Go](https://golang.org/doc/install) installed

### Local Testing

To run the example locally:

1. Ensure Docker is running on your machine.

2. Navigate to the `examples` directory:
   ```
   cd examples
   ```

3. Build the Lambda function:
   ```
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap example.go
   ```

4. Run the example using the AWS SAM CLI:
   ```
   sam local start-api
   ```

   This command will start a local API Gateway emulator using Docker.

5. In a new terminal, you can now test your API with curl:
   ```
   curl http://localhost:3000/users
   curl http://localhost:3000/users/123
   curl -X POST http://localhost:3000/users -d '{"name": "John Doe"}'
   curl -X PUT http://localhost:3000/users/123 -d '{"name": "Jane Doe"}'
   curl -X DELETE http://localhost:3000/users/123
   ```

##  Benchmarks
 
Benchmarks can be run with `make benchmark` and the full benchmark code can be found [here](https://github.com/D-Andreev/lambdamux/blob/main/lambdamux_benchmark_test.go).
The router used in the benchmarks consists of 50 routes in total, some static and some dynamic.
| Benchmark                                                                | Operations | Time per Operation | Bytes per Operation | Allocations per Operation | Using aws-lambda-go-api-proxy | % Slower than LambdaMux |
|--------------------------------------------------------------------------|------------|---------------------|---------------------|---------------------------|-------------------------------|--------------------------|
| LambdaMux                                                                | 378,866    | 3,130 ns/op         | 2,444 B/op          | 40 allocs/op              | No                            | 0%                       |
| [LmdRouter](https://github.com/aquasecurity/lmdrouter)                   | 322,635    | 3,707 ns/op         | 2,060 B/op          | 33 allocs/op              | No                            | 18.43%                   |
| [Gin](https://github.com/gin-gonic/gin)                                  | 294,595    | 4,069 ns/op         | 3,975 B/op          | 47 allocs/op              | Yes                           | 29.99%                   |
| [Chi](https://github.com/go-chi/chi)                                     | 276,445    | 4,360 ns/op         | 4,312 B/op          | 49 allocs/op              | Yes                           | 39.30%                   |
| [Standard Library](https://pkg.go.dev/net/http#ServeMux)                          | 266,296    | 4,552 ns/op         | 3,989 B/op          | 48 allocs/op              | Yes                            | 45.43%                   |
| [Fiber](https://github.com/gofiber/fiber)                                | 211,684    | 5,653 ns/op         | 6,324 B/op          | 61 allocs/op              | Yes                           | 80.61%                   |
