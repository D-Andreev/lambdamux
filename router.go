package router

import (
	"context"
	"net/http"

	"github.com/D-Andreev/lambdamux/internal/radix"
	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Router struct {
	tree *radix.Node
}

func NewRouter() *Router {
	return &Router{
		tree: radix.NewNode("", false),
	}
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.addRoute("GET", path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.addRoute("POST", path, handler)
}

func (r *Router) PUT(path string, handler HandlerFunc) {
	r.addRoute("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.addRoute("DELETE", path, handler)
}

func (r *Router) addRoute(method, path string, handler HandlerFunc) {
	fullPath := method + " " + path
	r.tree.InsertWithHandler(fullPath, handler)
}

func (r *Router) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
