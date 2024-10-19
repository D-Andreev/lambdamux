// Package lambdamux provides a lightweight HTTP router for AWS Lambda functions.
package lambdamux

import (
	"context"
	"net/http"

	"github.com/D-Andreev/lambdamux/internal/radix"
	"github.com/aws/aws-lambda-go/events"
)

// HandlerFunc defines the function signature for request handlers
type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// LambdaMux is a request multiplexer for AWS Lambda functions
type LambdaMux struct {
	tree *radix.Node
}

// NewLambdaMux creates and returns a new LambdaMux instance
func NewLambdaMux() *LambdaMux {
	return &LambdaMux{
		tree: radix.NewNode("", false),
	}
}

func (r *LambdaMux) addRoute(method, path string, handler HandlerFunc) {
	fullPath := method + " " + path
	r.tree.InsertWithHandler(fullPath, handler)
}

// GET registers a new GET route with the given path and handler
func (r *LambdaMux) GET(path string, handler HandlerFunc) {
	r.addRoute("GET", path, handler)
}

// POST registers a new POST route with the given path and handler
func (r *LambdaMux) POST(path string, handler HandlerFunc) {
	r.addRoute("POST", path, handler)
}

// PUT registers a new PUT route with the given path and handler
func (r *LambdaMux) PUT(path string, handler HandlerFunc) {
	r.addRoute("PUT", path, handler)
}

// DELETE registers a new DELETE route with the given path and handler
func (r *LambdaMux) DELETE(path string, handler HandlerFunc) {
	r.addRoute("DELETE", path, handler)
}

// Handle processes the incoming API Gateway proxy request and returns the appropriate response
func (r *LambdaMux) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := req.HTTPMethod + " " + req.Path
	node, params := r.tree.Search(path)

	if node != nil && node.Handler != nil {
		req.PathParameters = params
		return node.Handler(ctx, req)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"error": "404 Not Found"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
