package router

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := NewRouter()

	router.POST("/pet", createHandler("POST", "/pet"))
	router.PUT("/pet", createHandler("PUT", "/pet"))
	router.GET("/pet/findByStatus", createHandler("GET", "/pet/findByStatus"))
	router.GET("/pet/findByTags", createHandler("GET", "/pet/findByTags"))
	router.GET("/pet/:petId", createHandler("GET", "/pet/:petId"))
	router.POST("/pet/:petId", createHandler("POST", "/pet/:petId"))
	router.DELETE("/pet/:petId", createHandler("DELETE", "/pet/:petId"))
	router.POST("/pet/:petId/uploadImage", createHandler("POST", "/pet/:petId/uploadImage"))
	router.GET("/store/inventory", createHandler("GET", "/store/inventory"))
	router.POST("/store/order", createHandler("POST", "/store/order"))
	router.GET("/store/order/:orderId", createHandler("GET", "/store/order/:orderId"))
	router.DELETE("/store/order/:orderId", createHandler("DELETE", "/store/order/:orderId"))
	router.POST("/user", createHandler("POST", "/user"))
	router.POST("/user/createWithList", createHandler("POST", "/user/createWithList"))
	router.GET("/user/login", createHandler("GET", "/user/login"))
	router.GET("/user/logout", createHandler("GET", "/user/logout"))
	router.GET("/user/:username", createHandler("GET", "/user/:username"))
	router.PUT("/user/:username", createHandler("PUT", "/user/:username"))
	router.DELETE("/user/:username", createHandler("DELETE", "/user/:username"))

	testCases := []struct {
		id             int
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   map[string]interface{}
		expectedParams map[string]string
	}{
		{
			1,
			"POST /pet",
			"POST",
			"/pet",
			200,
			map[string]interface{}{"message": "Handled POST request for /pet"},
			nil,
		},
		{
			2,
			"PUT /pet",
			"PUT",
			"/pet",
			200,
			map[string]interface{}{"message": "Handled PUT request for /pet"},
			nil,
		},
		{
			3,
			"GET /pet/findByStatus",
			"GET",
			"/pet/findByStatus",
			200,
			map[string]interface{}{"message": "Handled GET request for /pet/findByStatus"},
			nil,
		},
		{
			4,
			"GET /pet/findByTags",
			"GET",
			"/pet/findByTags",
			200,
			map[string]interface{}{"message": "Handled GET request for /pet/findByTags"},
			nil,
		},
		{
			5,
			"GET /pet/:petId",
			"GET",
			"/pet/123",
			200,
			map[string]interface{}{
				"message": "Handled GET request for /pet/:petId",
				"params":  map[string]string{"petId": "123"},
			},
			map[string]string{"petId": "123"},
		},
		{
			6,
			"POST /pet/:petId",
			"POST",
			"/pet/456",
			200,
			map[string]interface{}{
				"message": "Handled POST request for /pet/:petId",
				"params":  map[string]string{"petId": "456"},
			},
			map[string]string{"petId": "456"},
		},
		{
			7,
			"DELETE /pet/:petId",
			"DELETE",
			"/pet/789",
			200,
			map[string]interface{}{
				"message": "Handled DELETE request for /pet/:petId",
				"params":  map[string]string{"petId": "789"},
			},
			map[string]string{"petId": "789"},
		},
		{
			8,
			"POST /pet/:petId/uploadImage",
			"POST",
			"/pet/101/uploadImage",
			200,
			map[string]interface{}{
				"message": "Handled POST request for /pet/:petId/uploadImage",
				"params":  map[string]string{"petId": "101"},
			},
			map[string]string{"petId": "101"},
		},
		{
			9,
			"GET /store/inventory",
			"GET",
			"/store/inventory",
			200,
			map[string]interface{}{"message": "Handled GET request for /store/inventory"},
			nil,
		},
		{
			10,
			"POST /store/order",
			"POST",
			"/store/order",
			200,
			map[string]interface{}{"message": "Handled POST request for /store/order"},
			nil,
		},
		{
			11,
			"GET /store/order/:orderId",
			"GET",
			"/store/order/1001",
			200,
			map[string]interface{}{
				"message": "Handled GET request for /store/order/:orderId",
				"params":  map[string]string{"orderId": "1001"},
			},
			map[string]string{"orderId": "1001"},
		},
		{
			12,
			"DELETE /store/order/:orderId",
			"DELETE",
			"/store/order/1002",
			200,
			map[string]interface{}{
				"message": "Handled DELETE request for /store/order/:orderId",
				"params":  map[string]string{"orderId": "1002"},
			},
			map[string]string{"orderId": "1002"},
		},
		{
			13,
			"POST /user",
			"POST",
			"/user",
			200,
			map[string]interface{}{"message": "Handled POST request for /user"},
			nil,
		},
		{
			14,
			"POST /user/createWithList",
			"POST",
			"/user/createWithList",
			200,
			map[string]interface{}{"message": "Handled POST request for /user/createWithList"},
			nil,
		},
		{
			15,
			"GET /user/login",
			"GET",
			"/user/login",
			200,
			map[string]interface{}{"message": "Handled GET request for /user/login"},
			nil,
		},
		{
			16,
			"GET /user/logout",
			"GET",
			"/user/logout",
			200,
			map[string]interface{}{"message": "Handled GET request for /user/logout"},
			nil,
		},
		{
			17,
			"GET /user/:username",
			"GET",
			"/user/johndoe",
			200,
			map[string]interface{}{
				"message": "Handled GET request for /user/:username",
				"params":  map[string]string{"username": "johndoe"},
			},
			map[string]string{"username": "johndoe"},
		},
		{
			18,
			"PUT /user/:username",
			"PUT",
			"/user/janedoe",
			200,
			map[string]interface{}{
				"message": "Handled PUT request for /user/:username",
				"params":  map[string]string{"username": "janedoe"},
			},
			map[string]string{"username": "janedoe"},
		},
		{
			19,
			"DELETE /user/:username",
			"DELETE",
			"/user/testuser",
			200,
			map[string]interface{}{
				"message": "Handled DELETE request for /user/:username",
				"params":  map[string]string{"username": "testuser"},
			},
			map[string]string{"username": "testuser"},
		},
		{20, "Not Found", "GET", "/nonexistent", 404, map[string]interface{}{"error": "404 Not Found"}, nil},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d: %s", tc.id, tc.name), func(t *testing.T) {
			req := events.APIGatewayProxyRequest{
				HTTPMethod: tc.method,
				Path:       tc.path,
			}
			resp, err := router.Handle(context.Background(), req)

			assert.NoError(t, err, "Test case %d: %s - Unexpected error", tc.id, tc.name)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Test case %d: %s - Status code mismatch", tc.id, tc.name)

			var bodyMap map[string]interface{}
			err = json.Unmarshal([]byte(resp.Body), &bodyMap)
			assert.NoError(t, err, "Test case %d: %s - Failed to unmarshal response body", tc.id, tc.name)

			for key, expectedValue := range tc.expectedBody {
				actualValue, exists := bodyMap[key]
				assert.True(t, exists, "Test case %d: %s - Key '%s' not found in response body", tc.id, tc.name, key)
				if exists {
					assert.Equal(
						t,
						fmt.Sprintf("%v", expectedValue),
						fmt.Sprintf("%v", actualValue),
						"Test case %d: %s - Value mismatch for key '%s'", tc.id, tc.name, key,
					)
				}
			}
		})
	}
}

func createHandler(method, path string) HandlerFunc {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		responseBody := map[string]interface{}{
			"message": "Handled " + method + " request for " + path,
		}
		if len(req.PathParameters) > 0 {
			responseBody["params"] = req.PathParameters
		}
		jsonBody, err := json.Marshal(responseBody)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}
}
