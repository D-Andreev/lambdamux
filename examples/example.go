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
