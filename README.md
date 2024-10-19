# lambdamux

[![Test](https://github.com/D-Andreev/lambdamux/actions/workflows/test.yml/badge.svg)](https://github.com/D-Andreev/lambdamux/actions/workflows/test.yml)

A simple and lightweight high performance HTTP router specifically designed for AWS Lambda functions handling API Gateway requests. 

## Features
- Fast and efficient routing with static and dynamic route support
- Seamless handling of path parameters in routes
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
 
The following table shows the benchmark results for `lambdamux` compared to other popular routers, including those using `aws-lambda-go-api-proxy`:

| Benchmark | Operations | Time per Operation | Bytes per Operation | Allocations per Operation | Using aws-lambda-go-api-proxy | % Slower than LambdaMux |
|-----------|------------|---------------------|---------------------|---------------------------|-------------------------------|--------------------------|
| LambdaMux | 535,137 | 2,229 ns/op | 1,852 B/op | 29 allocs/op | No | 0% |
| LmdRouter | 508,597 | 2,329 ns/op | 1,615 B/op | 25 allocs/op | No | 4.49% |
| AWSLambdaGoAPIProxyWithGin | 372,218 | 3,169 ns/op | 3,430 B/op | 38 allocs/op | Yes | 42.17% |
| AWSLambdaGoAPIProxyWithChi | 348,360 | 3,394 ns/op | 3,786 B/op | 40 allocs/op | Yes | 52.27% |
| AWSLambdaGoAPIProxyWithFiber | 259,388 | 4,572 ns/op | 5,770 B/op | 52 allocs/op | Yes | 105.11% |

`lambdamux` performs competitively, being slightly faster than LmdRouter and using marginally more memory. It's much faster than using `aws-lambda-go-api-proxy` with Gin, Fiber, and Chi.

Benchmarks can be run with `make benchmark` and the full benchmark code can be found [here](https://github.com/D-Andreev/lambdamux/blob/main/lambdamux_benchmark_test.go).
