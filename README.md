
# lambdamux

`lambdamux` is a simple and lightweight high performance HTTP router specifically designed for AWS Lambda functions handling API Gateway requests. 

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

## Installation

```
go get github.com/D-Andreev/lambdahttp
```

## Usage

```go
package main

import (
	"context"

	"github.com/D-Andreev/lambdamux"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	router := lambdamux.NewRouter()
	router.POST("/users", func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Handle create user route
		return events.APIGatewayProxyResponse{
			StatusCode: 201,
			Body:       "User created successfully",
		}, nil
	})
	router.GET("/users/:id", func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Handle get user route
		userID := req.PathParameters["id"]
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "User details for ID: " + userID,
		}, nil
	})

	lambda.Start(router.Handle)
}

```
##  Benchmarks

 
The following table shows the benchmark results for `lambdamux` compared to using http routers with `aws-lambda-go-api-proxy`

| Benchmark | Operations | Time per Operation | Bytes per Operation | Allocations per Operation |
|-----------|------------|---------------------|---------------------|---------------------------|
| LambdaHTTPRouter | 1,061,984 | 1,114 ns/op | 1,310 B/op | 19 allocs/op |
| AWSLambdaGoAPIProxyWithGin | 397,654 | 3,028 ns/op | 3,562 B/op | 34 allocs/op |

`lambdamux` is  ~`3x` faster and uses ~`60%` less memory.
Benchmarks can be found [here](https://github.com/D-Andreev/lambdamux/blob/main/router_benchmark_test.go).