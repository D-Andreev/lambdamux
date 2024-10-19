package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

var routes = []struct {
	method string
	path   string
}{
	{"POST", "/pet"},
	{"PUT", "/pet"},
	{"GET", "/pet/findByStatus"},
	{"GET", "/pet/findByTags"},
	{"GET", "/pet/:petId"},
	{"POST", "/pet/:petId"},
	{"DELETE", "/pet/:petId"},
	{"POST", "/pet/:petId/uploadImage"},
	{"GET", "/store/inventory"},
	{"POST", "/store/order"},
	{"GET", "/store/order/:orderId"},
	{"DELETE", "/store/order/:orderId"},
	{"POST", "/user"},
	{"POST", "/user/createWithList"},
	{"GET", "/user/login"},
	{"GET", "/user/logout"},
	{"GET", "/user/:username"},
	{"PUT", "/user/:username"},
	{"DELETE", "/user/:username"},
}

func setupLambdaHTTPRouter() *Router {
	router := NewRouter()
	for _, route := range routes {
		router.addRoute(route.method, route.path, lambdahttpCreateHandler(route.method, route.path))
	}
	return router
}

func setupGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	for _, route := range routes {
		r.Handle(route.method, route.path, ginCreateHandler(route.method, route.path))
	}
	return r
}

func lambdahttpCreateHandler(method, path string) HandlerFunc {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		responseBody := map[string]interface{}{
			"message": "Handled " + method + " request for " + path,
		}
		if len(req.PathParameters) > 0 {
			responseBody["params"] = req.PathParameters
		}
		jsonBody, _ := json.Marshal(responseBody)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}
}

func ginCreateHandler(method, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseBody := map[string]interface{}{
			"message": "Handled " + method + " request for " + path,
		}
		params := c.Params
		if len(params) > 0 {
			paramMap := make(map[string]string)
			for _, param := range params {
				paramMap[param.Key] = param.Value
			}
			responseBody["params"] = paramMap
		}
		c.JSON(200, responseBody)
	}
}

func setupFiberRouter() *fiber.App {
	app := fiber.New()
	for _, route := range routes {
		app.Add(route.method, route.path, fiberCreateHandler(route.method, route.path))
	}
	return app
}

func fiberCreateHandler(method, path string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		responseBody := map[string]interface{}{
			"message": "Handled " + method + " request for " + path,
		}
		params := c.AllParams()
		if len(params) > 0 {
			responseBody["params"] = params
		}
		return c.JSON(responseBody)
	}
}

func setupLmdRouter() *lmdrouter.Router {
	router := lmdrouter.NewRouter("")
	for _, route := range routes {
		router.Route(route.method, route.path, lmdCreateHandler(route.method, route.path))
	}
	return router
}

func lmdCreateHandler(method, path string) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		responseBody := map[string]interface{}{
			"message": "Handled " + method + " request for " + path,
		}
		if len(req.PathParameters) > 0 {
			responseBody["params"] = req.PathParameters
		}
		jsonBody, _ := json.Marshal(responseBody)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}
}

var benchmarkRequests = []events.APIGatewayProxyRequest{
	{HTTPMethod: "POST", Path: "/pet"},
	{HTTPMethod: "PUT", Path: "/pet"},
	{HTTPMethod: "GET", Path: "/pet/findByStatus"},
	{HTTPMethod: "GET", Path: "/pet/findByTags"},
	{HTTPMethod: "GET", Path: "/pet/123", PathParameters: map[string]string{"petId": "123"}},
	{HTTPMethod: "POST", Path: "/pet/456", PathParameters: map[string]string{"petId": "456"}},
	{HTTPMethod: "DELETE", Path: "/pet/789", PathParameters: map[string]string{"petId": "789"}},
	{HTTPMethod: "POST", Path: "/pet/101/uploadImage", PathParameters: map[string]string{"petId": "101"}},
	{HTTPMethod: "GET", Path: "/store/inventory"},
	{HTTPMethod: "POST", Path: "/store/order"},
	{HTTPMethod: "GET", Path: "/store/order/202", PathParameters: map[string]string{"orderId": "202"}},
	{HTTPMethod: "DELETE", Path: "/store/order/303", PathParameters: map[string]string{"orderId": "303"}},
	{HTTPMethod: "POST", Path: "/user"},
	{HTTPMethod: "POST", Path: "/user/createWithList"},
	{HTTPMethod: "GET", Path: "/user/login"},
	{HTTPMethod: "GET", Path: "/user/logout"},
	{HTTPMethod: "GET", Path: "/user/johndoe", PathParameters: map[string]string{"username": "johndoe"}},
	{HTTPMethod: "PUT", Path: "/user/janedoe", PathParameters: map[string]string{"username": "janedoe"}},
	{HTTPMethod: "DELETE", Path: "/user/bobsmith", PathParameters: map[string]string{"username": "bobsmith"}},
}

func assertResponse(b *testing.B, resp events.APIGatewayProxyResponse, req events.APIGatewayProxyRequest) {
	b.Helper()
	assert.Equal(b, 200, resp.StatusCode)
	if len(req.PathParameters) > 0 {
		var body map[string]interface{}
		err := json.Unmarshal([]byte(resp.Body), &body)
		assert.NoError(b, err)
		params, ok := body["params"].(map[string]interface{})
		assert.True(b, ok, "params should be a map")
		for key, expectedValue := range req.PathParameters {
			actualValue, exists := params[key]
			assert.True(b, exists, "param %s should exist", key)
			assert.Equal(b, expectedValue, actualValue, "param %s should match", key)
		}
	}
}

func BenchmarkLambdaHTTPRouter(b *testing.B) {
	router := setupLambdaHTTPRouter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := benchmarkRequests[i%len(benchmarkRequests)]
		resp, err := router.Handle(context.Background(), req)
		assert.NoError(b, err)
		assertResponse(b, resp, req)
	}
}

func BenchmarkAWSLambdaGoAPIProxyWithGin(b *testing.B) {
	r := setupGinRouter()
	adapter := ginadapter.New(r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := benchmarkRequests[i%len(benchmarkRequests)]
		resp, err := adapter.ProxyWithContext(context.Background(), req)
		assert.NoError(b, err)
		assertResponse(b, resp, req)
	}
}

func BenchmarkAWSLambdaGoAPIProxyWithFiber(b *testing.B) {
	app := setupFiberRouter()
	adapter := fiberadapter.New(app)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := benchmarkRequests[i%len(benchmarkRequests)]
		resp, err := adapter.ProxyWithContext(context.Background(), req)
		assert.NoError(b, err)
		assertResponse(b, resp, req)
	}
}

func BenchmarkLmdRouter(b *testing.B) {
	router := setupLmdRouter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := benchmarkRequests[i%len(benchmarkRequests)]
		resp, err := router.Handler(context.Background(), req)
		assert.NoError(b, err)
		assertResponse(b, resp, req)
	}
}
