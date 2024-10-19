package lambdamux

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"

	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
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
	{"GET", "/pet/:petId/medical-history"},
	{"POST", "/pet/:petId/vaccination"},
	{"GET", "/store/order/:orderId/tracking"},
	{"PUT", "/store/order/:orderId/status"},
	{"GET", "/user/:username/preferences"},
	{"POST", "/user/:username/address"},
	{"GET", "/pet/:petId/appointments"},
	{"POST", "/pet/:petId/appointment"},
	{"PUT", "/pet/:petId/appointment/:appointmentId"},
	{"DELETE", "/pet/:petId/appointment/:appointmentId"},
	{"GET", "/store/products"},
	{"GET", "/store/product/:productId"},
	{"POST", "/store/product"},
	{"PUT", "/store/product/:productId"},
	{"DELETE", "/store/product/:productId"},
	{"GET", "/user/:username/orders"},
	{"POST", "/user/:username/review"},
	{"GET", "/user/:username/review/:reviewId"},
	{"PUT", "/user/:username/review/:reviewId"},
	{"DELETE", "/user/:username/review/:reviewId"},
	{"GET", "/clinic/:clinicId"},
	{"POST", "/clinic"},
	{"PUT", "/clinic/:clinicId"},
	{"DELETE", "/clinic/:clinicId"},
	{"GET", "/clinic/:clinicId/staff"},
	{"POST", "/clinic/:clinicId/staff"},
	{"GET", "/clinic/:clinicId/staff/:staffId"},
	{"PUT", "/clinic/:clinicId/staff/:staffId"},
	{"DELETE", "/clinic/:clinicId/staff/:staffId"},
	{"GET", "/clinic/:clinicId/appointments"},
	{"POST", "/clinic/:clinicId/appointment/:appointmentId/reschedule"},
}

var allParams = []string{
	"petId", "orderId", "username", "appointmentId", "productId", "reviewId", "clinicId", "staffId",
}

func setupLambdaMux() *LambdaMux {
	router := NewLambdaMux()
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
		params := make(map[string]string)
		for _, param := range allParams {
			if value := c.Param(param); value != "" {
				params[param] = value
			}
		}
		if len(params) > 0 {
			responseBody["params"] = params
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
		params := make(map[string]string)
		for _, param := range allParams {
			if value := c.Params(param); value != "" {
				params[param] = value
			}
		}
		if len(params) > 0 {
			responseBody["params"] = params
		}
		return c.JSON(responseBody)
	}
}

func setupChiRouter() *chi.Mux {
	r := chi.NewRouter()
	for _, route := range routes {
		// Convert :param to {param} for Chi router
		chiPath := route.path
		for _, param := range allParams {
			chiPath = strings.Replace(chiPath, ":"+param, "{"+param+"}", -1)
		}
		r.MethodFunc(route.method, chiPath, chiCreateHandler(route.method, route.path))
	}
	return r
}

func chiCreateHandler(method, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseBody := map[string]interface{}{
			"message": "Handled " + method + " request for " + path,
		}

		params := make(map[string]string)
		for _, param := range allParams {
			if value := chi.URLParam(r, param); value != "" {
				params[param] = value
			}
		}

		if len(params) > 0 {
			responseBody["params"] = params
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseBody)
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
	{HTTPMethod: "GET", Path: "/pet/404/medical-history", PathParameters: map[string]string{"petId": "404"}},
	{HTTPMethod: "POST", Path: "/pet/505/vaccination", PathParameters: map[string]string{"petId": "505"}},
	{HTTPMethod: "GET", Path: "/store/order/606/tracking", PathParameters: map[string]string{"orderId": "606"}},
	{HTTPMethod: "PUT", Path: "/store/order/707/status", PathParameters: map[string]string{"orderId": "707"}},
	{HTTPMethod: "GET", Path: "/user/alicesmith/preferences", PathParameters: map[string]string{"username": "alicesmith"}},
	{HTTPMethod: "POST", Path: "/user/bobdoe/address", PathParameters: map[string]string{"username": "bobdoe"}},
	{HTTPMethod: "GET", Path: "/pet/808/appointments", PathParameters: map[string]string{"petId": "808"}},
	{HTTPMethod: "POST", Path: "/pet/909/appointment", PathParameters: map[string]string{"petId": "909"}},
	{HTTPMethod: "PUT", Path: "/pet/1010/appointment/2020", PathParameters: map[string]string{"petId": "1010", "appointmentId": "2020"}},
	{HTTPMethod: "DELETE", Path: "/pet/1111/appointment/2121", PathParameters: map[string]string{"petId": "1111", "appointmentId": "2121"}},
	{HTTPMethod: "GET", Path: "/store/products"},
	{HTTPMethod: "GET", Path: "/store/product/3030", PathParameters: map[string]string{"productId": "3030"}},
	{HTTPMethod: "POST", Path: "/store/product"},
	{HTTPMethod: "PUT", Path: "/store/product/4040", PathParameters: map[string]string{"productId": "4040"}},
	{HTTPMethod: "DELETE", Path: "/store/product/5050", PathParameters: map[string]string{"productId": "5050"}},
	{HTTPMethod: "GET", Path: "/user/charlielee/orders", PathParameters: map[string]string{"username": "charlielee"}},
	{HTTPMethod: "POST", Path: "/user/davidwang/review", PathParameters: map[string]string{"username": "davidwang"}},
	{HTTPMethod: "GET", Path: "/user/evebrown/review/6060", PathParameters: map[string]string{"username": "evebrown", "reviewId": "6060"}},
	{HTTPMethod: "PUT", Path: "/user/frankgreen/review/7070", PathParameters: map[string]string{"username": "frankgreen", "reviewId": "7070"}},
	{HTTPMethod: "DELETE", Path: "/user/gracewu/review/8080", PathParameters: map[string]string{"username": "gracewu", "reviewId": "8080"}},
	{HTTPMethod: "GET", Path: "/clinic/9090", PathParameters: map[string]string{"clinicId": "9090"}},
	{HTTPMethod: "POST", Path: "/clinic"},
	{HTTPMethod: "PUT", Path: "/clinic/1212", PathParameters: map[string]string{"clinicId": "1212"}},
	{HTTPMethod: "DELETE", Path: "/clinic/1313", PathParameters: map[string]string{"clinicId": "1313"}},
	{HTTPMethod: "GET", Path: "/clinic/1414/staff", PathParameters: map[string]string{"clinicId": "1414"}},
	{HTTPMethod: "POST", Path: "/clinic/1515/staff", PathParameters: map[string]string{"clinicId": "1515"}},
	{HTTPMethod: "GET", Path: "/clinic/1616/staff/1717", PathParameters: map[string]string{"clinicId": "1616", "staffId": "1717"}},
	{HTTPMethod: "PUT", Path: "/clinic/1818/staff/1919", PathParameters: map[string]string{"clinicId": "1818", "staffId": "1919"}},
	{HTTPMethod: "DELETE", Path: "/clinic/2020/staff/2121", PathParameters: map[string]string{"clinicId": "2020", "staffId": "2121"}},
	{HTTPMethod: "GET", Path: "/clinic/2222/appointments", PathParameters: map[string]string{"clinicId": "2222"}},
	{HTTPMethod: "POST", Path: "/clinic/2323/appointment/2424/reschedule", PathParameters: map[string]string{"clinicId": "2323", "appointmentId": "2424"}},
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

func BenchmarkLambdaMux(b *testing.B) {
	router := setupLambdaMux()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := benchmarkRequests[i%len(benchmarkRequests)]
		resp, err := router.Handle(context.Background(), req)
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

func BenchmarkAWSLambdaGoAPIProxyWithChi(b *testing.B) {
	r := setupChiRouter()
	adapter := chiadapter.New(r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := benchmarkRequests[i%len(benchmarkRequests)]
		resp, err := adapter.ProxyWithContext(context.Background(), req)
		assert.NoError(b, err)
		assertResponse(b, resp, req)
	}
}
