package radix

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type InsertTestCase struct {
	id        int
	input     []string
	keys      []string
	fullItems []string
}

type SearchTestCase struct {
	id               int
	input            []string
	search           string
	output           string
	params           map[string]string
	notFoundExpected bool
}

func TestInsert(t *testing.T) {
	testCases := []InsertTestCase{
		{id: 1, input: []string{"water"}, keys: []string{"water"}, fullItems: []string{"water"}},
		{
			id: 2, input: []string{"water", "water"}, keys: []string{"water"}, fullItems: []string{"water"},
		},
		{
			id: 3, input: []string{"water", "slow"}, keys: []string{"slow", "water"},
			fullItems: []string{"slow", "water"},
		},
		{
			id:        4,
			input:     []string{"water", "slow", "slower"},
			keys:      []string{"er", "slow", "water"},
			fullItems: []string{"slow", "slower", "water"},
		},
		{
			id:        5,
			input:     []string{"water", "slow", "slower", "wash"},
			keys:      []string{"er", "sh", "slow", "ter", "wa"},
			fullItems: []string{"slow", "slower", "wash", "water"},
		},
		{
			id:        6,
			input:     []string{"water", "slow", "slower", "wash", "washer"},
			keys:      []string{"er", "er", "sh", "slow", "ter", "wa"},
			fullItems: []string{"slow", "slower", "wash", "washer", "water"},
		},
		{
			id:        7,
			input:     []string{"water", "slow", "slower", "wash", "washer", "wasnt"},
			keys:      []string{"er", "er", "h", "nt", "s", "slow", "ter", "wa"},
			fullItems: []string{"slow", "slower", "wash", "washer", "wasnt", "water"},
		},
		{
			id:        8,
			input:     []string{"water", "slow", "slower", "wash", "washer", "wasnt", "watering"},
			keys:      []string{"er", "er", "h", "ing", "nt", "s", "slow", "ter", "wa"},
			fullItems: []string{"slow", "slower", "wash", "washer", "wasnt", "water", "watering"},
		},
		{
			id: 9,
			input: []string{
				"alligator", "alien", "baloon", "chromodynamic", "romane", "romanus", "romulus", "rubens", "ruber",
				"rubicon", "rubicundus", "all", "rub", "ba",
			},
			keys: []string{
				"al", "an", "ba", "chromodynamic", "e", "e", "ic", "ien", "igator", "l", "loon", "ns", "om", "on", "r",
				"r", "ub", "ulus", "undus", "us",
			},
			fullItems: []string{
				"alien", "all", "alligator", "ba", "baloon", "chromodynamic", "romane", "romanus", "romulus", "rub",
				"rubens",
				"ruber", "rubicon", "rubicundus",
			},
		},
	}

	for _, tc := range testCases {
		tree := NewNode("", false)
		for _, j := range tc.input {
			tree.Insert(j)
		}
		allKeys := tree.GetAllNodeValues()
		assert.Equal(t, tc.keys, allKeys)
		allItems := tree.GetAllCompleteItems()
		assert.Equal(t, tc.fullItems, allItems)
	}

}

func TestInsertWithRandomIds(t *testing.T) {
	randomIds := generateUUIDs()
	sort.Strings(randomIds)
	testCases := []InsertTestCase{
		{input: randomIds, fullItems: randomIds},
	}

	for _, tc := range testCases {
		tree := NewNode("", false)
		for _, j := range tc.input {
			tree.Insert(j)
		}
		allItems := tree.GetAllCompleteItems()
		assert.Equal(t, tc.fullItems, allItems)
	}
}

func TestSearch(t *testing.T) {
	testCases := []SearchTestCase{
		{id: 0, input: []string{"water"}, search: "non-existing-item", output: ""},
		{id: 1, input: []string{"water"}, search: "water", output: "water"},
		{
			id: 2, input: []string{"water", "water"}, search: "water", output: "water",
		},
		{
			id: 3, input: []string{"water", "slow"}, output: "slow", search: "slow",
		},
		{
			id:     4,
			input:  []string{"water", "slow", "slower", "wash", "washer", "wasnt", "watering"},
			search: "wasnt",
			output: "wasnt",
		},
		{
			id: 5,
			input: []string{
				"alligator", "alien", "baloon", "chromodynamic", "romane", "romanus", "romulus", "rubens", "ruber",
				"rubicon", "rubicundus", "all", "rub", "ba",
			},
			search: "chromodynamic",
			output: "chromodynamic",
		},
		{
			id: 6,
			input: []string{
				"alligator", "alien", "baloon", "chromodynamic", "romane", "romanus", "romulus", "rubens", "ruber",
				"rubicon", "rubicundus", "all", "rub", "ba",
			},
			search: "rubicundus",
			output: "rubicundus",
		},
		{
			id: 7,
			input: []string{
				"alligator", "alien", "baloon", "chromodynamic", "romane", "romanus", "romulus", "rubens", "ruber",
				"rubicon", "rubicundus", "all", "rub", "ba",
			},
			search: "rub",
			output: "rub",
		},
		{
			id: 8,
			input: []string{
				"alligator", "alien", "baloon", "chromodynamic", "romane", "romanus", "romulus", "rubens", "ruber",
				"rubicon", "rubicundus", "all", "rub", "ba",
			},
			search: "ba",
			output: "ba",
		},
	}

	for _, tc := range testCases {
		tree := NewNode("", false)
		for _, j := range tc.input {
			tree.Insert(j)
		}

		node, _ := tree.Search(tc.search)
		if tc.output == "" {
			assert.Nil(t, node, fmt.Sprintf("Failed test id: %d\n", tc.id))
		} else {
			assert.NotNil(t, node, fmt.Sprintf("Failed test id: %d\n", tc.id))
			assert.Equal(t, tc.output, node.fullValue, fmt.Sprintf("Failed test id: %d\n", tc.id))
		}
	}
}

func TestSearchWithRandomIds(t *testing.T) {
	randomIds := generateUUIDs()
	sort.Strings(randomIds)
	var testCases []*SearchTestCase
	for i, randId := range randomIds {
		testCases = append(
			testCases, &SearchTestCase{
				id: i, input: randomIds, search: randId, output: randId,
			},
		)
	}

	for _, tc := range testCases {
		tree := NewNode("", false)
		for _, j := range tc.input {
			tree.Insert(j)
		}

		node, _ := tree.Search(tc.search)

		assert.NotNil(t, node)
		assert.Equal(t, tc.output, node.fullValue)
	}
}

func generateUUIDs() []string {
	var res []string
	for i := 0; i < 1000; i++ {
		res = append(res, uuid.New().String())
	}
	return res
}

type HashmapImplementation struct {
	items map[string]string
}

func NewHashmapImplementation() *HashmapImplementation {
	return &HashmapImplementation{
		items: make(map[string]string),
	}
}

func (hs *HashmapImplementation) Insert(input string) {
	hs.items[input] = input
}

func (hs *HashmapImplementation) Search(input string) string {
	return hs.items[input]
}

func TestRouter(t *testing.T) {
	input := []string{
		"GET /",
		"GET /contact",
		"GET /api/widgets",
		"POST /api/widgets",
		"POST /api/widgets/:id",
		"POST /api/widgets/:id/parts",
		"POST /api/widgets/:id/parts/:partId/update",
		"POST /api/widgets/:id/parts/:partId/delete",
		"POST /:id",
		"POST /:id/admin",
		"POST /:id/image",
		"DELETE /users/:id",
		"DELETE /users/:id",
		"DELETE /users/:id/admin",
		"DELETE /images/:id",
		"GET /products/:category/:id",
		"PUT /customers/:customerId/orders/:orderId",
		"PATCH /articles/:articleId/comments/:commentId",
		"GET /search/:query/page/:pageNumber",
		"POST /upload/:fileType/:userId",
	}
	testCases := []SearchTestCase{
		{
			id:     1,
			input:  input,
			search: "GET /",
			output: "GET /",
			params: map[string]string{},
		},
		{
			id:     2,
			input:  input,
			search: "GET /contact",
			output: "GET /contact",
			params: map[string]string{},
		},
		{
			id:     3,
			input:  input,
			search: "GET /api/widgets",
			output: "GET /api/widgets",
			params: map[string]string{},
		},
		{
			id:     4,
			input:  input,
			search: "POST /api/widgets/123",
			output: "POST /api/widgets/:id",
			params: map[string]string{"id": "123"},
		},
		{
			id:     5,
			input:  input,
			search: "POST /api/widgets/123/parts",
			output: "POST /api/widgets/:id/parts",
			params: map[string]string{"id": "123"},
		},
		{
			id:               6,
			input:            input,
			search:           "POST /api/widgets/123/parts/123",
			output:           "POST /api/widgets/:id/parts/:partId",
			params:           map[string]string{"id": "123", "partId": "123"},
			notFoundExpected: true,
		},
		{
			id:     7,
			input:  input,
			search: "POST /api/widgets/123/parts/123/update",
			output: "POST /api/widgets/:id/parts/:partId/update",
			params: map[string]string{"id": "123", "partId": "123"},
		},
		{
			id:     8,
			input:  input,
			search: "POST /api/widgets/123/parts/123/delete",
			output: "POST /api/widgets/:id/parts/:partId/delete",
			params: map[string]string{"id": "123", "partId": "123"},
		},
		{
			id:     9,
			input:  input,
			search: "POST /123",
			output: "POST /:id",
			params: map[string]string{"id": "123"},
		},
		{
			id:     10,
			input:  input,
			search: "POST /123/admin",
			output: "POST /:id/admin",
			params: map[string]string{"id": "123"},
		},
		{
			id:     11,
			input:  input,
			search: "POST /123/image",
			output: "POST /:id/image",
			params: map[string]string{"id": "123"},
		},
		{
			id:               12,
			input:            input,
			search:           "POST /123/images",
			output:           "",
			params:           map[string]string{},
			notFoundExpected: true,
		},
		{
			id:               13,
			input:            input,
			search:           "GET /nonexistent",
			output:           "",
			params:           map[string]string{},
			notFoundExpected: true,
		},
		{
			id:               14,
			input:            input,
			search:           "POST /api/widgets/123/nonexistent",
			output:           "",
			params:           map[string]string{},
			notFoundExpected: true,
		},
		{
			id:               15,
			input:            input,
			search:           "PUT /api/widgets/123",
			output:           "",
			params:           map[string]string{},
			notFoundExpected: true,
		},
		{
			id:               16,
			input:            input,
			search:           "POST /api/widgets/123/parts/456/unknown",
			output:           "",
			params:           map[string]string{},
			notFoundExpected: true,
		},
		{
			id:     18,
			input:  input,
			search: "POST /api/widgets/very-long-slug-with-dashes",
			output: "POST /api/widgets/:id",
			params: map[string]string{"id": "very-long-slug-with-dashes"},
		},
		{
			id:     19,
			input:  input,
			search: "POST /api/widgets/123-456/parts/789-abc/update",
			output: "POST /api/widgets/:id/parts/:partId/update",
			params: map[string]string{"id": "123-456", "partId": "789-abc"},
		},
		{
			id:     20,
			input:  input,
			search: "POST /api/widgets/123_456/parts/789_abc/delete",
			output: "POST /api/widgets/:id/parts/:partId/delete",
			params: map[string]string{"id": "123_456", "partId": "789_abc"},
		},
		{
			id:     21,
			input:  input,
			search: "POST /complex-slug-with-numbers-123",
			output: "POST /:id",
			params: map[string]string{"id": "complex-slug-with-numbers-123"},
		},
		{
			id:     22,
			input:  input,
			search: "POST /UpperCaseSlug/admin",
			output: "POST /:id/admin",
			params: map[string]string{"id": "UpperCaseSlug"},
		},
		{
			id:               23,
			input:            input,
			search:           "GET /api/widgets/",
			output:           "",
			notFoundExpected: true,
		},
		{
			id:               24,
			input:            input,
			search:           "DELETE /users",
			output:           "",
			notFoundExpected: true,
		},
		{
			id:     25,
			input:  input,
			search: "DELETE /users/123",
			output: "DELETE /users/:id",
			params: map[string]string{"id": "123"},
		},
		{
			id:     26,
			input:  input,
			search: "DELETE /users/123",
			output: "DELETE /users/:id",
			params: map[string]string{"id": "123"},
		},
		{
			id:     27,
			input:  input,
			search: "DELETE /users/123/admin",
			output: "DELETE /users/:id/admin",
			params: map[string]string{"id": "123"},
		},
		{
			id:     28,
			input:  input,
			search: "DELETE /images/123",
			output: "DELETE /images/:id",
			params: map[string]string{"id": "123"},
		},
		{
			id:     29,
			input:  input,
			search: "DELETE /images/123",
			output: "DELETE /images/:id",
			params: map[string]string{"id": "123"},
		},
		{
			id:     30,
			input:  input,
			search: "GET /products/electronics/laptop-123",
			output: "GET /products/:category/:id",
			params: map[string]string{"category": "electronics", "id": "laptop-123"},
		},
		{
			id:     31,
			input:  input,
			search: "PUT /customers/cust-456/orders/order-789",
			output: "PUT /customers/:customerId/orders/:orderId",
			params: map[string]string{"customerId": "cust-456", "orderId": "order-789"},
		},
		{
			id:     32,
			input:  input,
			search: "PATCH /articles/art-101/comments/comment-202",
			output: "PATCH /articles/:articleId/comments/:commentId",
			params: map[string]string{"articleId": "art-101", "commentId": "comment-202"},
		},
		{
			id:     33,
			input:  input,
			search: "GET /search/golang/page/2",
			output: "GET /search/:query/page/:pageNumber",
			params: map[string]string{"query": "golang", "pageNumber": "2"},
		},
		{
			id:     34,
			input:  input,
			search: "POST /upload/image/user-303",
			output: "POST /upload/:fileType/:userId",
			params: map[string]string{"fileType": "image", "userId": "user-303"},
		},
	}

	for _, tc := range testCases {
		tree := NewNode("", false)
		for _, j := range tc.input {
			tree.Insert(j)
		}
		result, params := tree.Search(tc.search)
		if tc.notFoundExpected {
			assert.Nil(t, result, fmt.Sprintf("Test id %d failed: expected nil result, but got %v", tc.id, result))
		} else {
			assert.NotNil(t, result, fmt.Sprintf("Test id %d failed: expected non-nil result, but got nil", tc.id))
			assert.NotNil(t, params, fmt.Sprintf("Test id %d failed: expected non-nil params, but got nil", tc.id))
			assert.Equal(t, tc.output, result.fullValue, fmt.Sprintf("Test id %d failed: expected output %s, but got %s", tc.id, tc.output, result.fullValue))
			assert.Equal(t, tc.params, params, fmt.Sprintf("Test id %d failed: expected params %v, but got %v", tc.id, tc.params, params))
		}
	}
}

func TestSearchPetStoreAPI(t *testing.T) {
	input := []string{
		"POST /pet",
		"PUT /pet",
		"GET /pet/findByStatus",
		"GET /pet/findByTags",
		"GET /pet/:petId",
		"POST /pet/:petId",
		"DELETE /pet/:petId",
		"POST /pet/:petId/uploadImage",
		"GET /store/inventory",
		"POST /store/order",
		"GET /store/order/:orderId",
		"DELETE /store/order/:orderId",
		"POST /user",
		"POST /user/createWithList",
		"GET /user/login",
		"GET /user/logout",
		"GET /user/:username",
		"PUT /user/:username",
		"DELETE /user/:username",
	}
	testCases := []SearchTestCase{
		{
			id:     1,
			input:  input,
			search: "POST /pet",
			output: "POST /pet",
			params: map[string]string{},
		},
		{
			id:     2,
			input:  input,
			search: "PUT /pet",
			output: "PUT /pet",
			params: map[string]string{},
		},
		{
			id:     3,
			input:  input,
			search: "GET /pet/findByStatus",
			output: "GET /pet/findByStatus",
			params: map[string]string{},
		},
		{
			id:     4,
			input:  input,
			search: "GET /pet/findByTags",
			output: "GET /pet/findByTags",
			params: map[string]string{},
		},
		{
			id:     5,
			input:  input,
			search: "GET /pet/123",
			output: "GET /pet/:petId",
			params: map[string]string{"petId": "123"},
		},
		{
			id:     6,
			input:  input,
			search: "POST /pet/456",
			output: "POST /pet/:petId",
			params: map[string]string{"petId": "456"},
		},
		{
			id:     7,
			input:  input,
			search: "DELETE /pet/789",
			output: "DELETE /pet/:petId",
			params: map[string]string{"petId": "789"},
		},
		{
			id:     8,
			input:  input,
			search: "POST /pet/101/uploadImage",
			output: "POST /pet/:petId/uploadImage",
			params: map[string]string{"petId": "101"},
		},
		{
			id:     9,
			input:  input,
			search: "GET /store/inventory",
			output: "GET /store/inventory",
			params: map[string]string{},
		},
		{
			id:     10,
			input:  input,
			search: "POST /store/order",
			output: "POST /store/order",
			params: map[string]string{},
		},
		{
			id:     11,
			input:  input,
			search: "GET /store/order/202",
			output: "GET /store/order/:orderId",
			params: map[string]string{"orderId": "202"},
		},
		{
			id:     12,
			input:  input,
			search: "DELETE /store/order/303",
			output: "DELETE /store/order/:orderId",
			params: map[string]string{"orderId": "303"},
		},
		{
			id:     13,
			input:  input,
			search: "POST /user",
			output: "POST /user",
			params: map[string]string{},
		},
		{
			id:     14,
			input:  input,
			search: "POST /user/createWithList",
			output: "POST /user/createWithList",
			params: map[string]string{},
		},
		{
			id:     15,
			input:  input,
			search: "GET /user/login",
			output: "GET /user/login",
			params: map[string]string{},
		},
		{
			id:     16,
			input:  input,
			search: "GET /user/logout",
			output: "GET /user/logout",
			params: map[string]string{},
		},
		{
			id:     17,
			input:  input,
			search: "GET /user/johndoe",
			output: "GET /user/:username",
			params: map[string]string{"username": "johndoe"},
		},
		{
			id:     18,
			input:  input,
			search: "PUT /user/janedoe",
			output: "PUT /user/:username",
			params: map[string]string{"username": "janedoe"},
		},
		{
			id:     19,
			input:  input,
			search: "DELETE /user/testuser",
			output: "DELETE /user/:username",
			params: map[string]string{"username": "testuser"},
		},
		{
			id:               20,
			input:            input,
			search:           "GET /nonexistent/path",
			notFoundExpected: true,
		},
	}

	for _, tc := range testCases {
		tree := NewNode("", false)
		for _, j := range tc.input {
			tree.Insert(j)
		}
		result, params := tree.Search(tc.search)
		if tc.notFoundExpected {
			assert.Nil(t, result, fmt.Sprintf("Test id %d failed: expected nil result, but got %v", tc.id, result))
		} else {
			assert.NotNil(t, result, fmt.Sprintf("Test id %d failed: expected non-nil result, but got nil", tc.id))
			assert.NotNil(t, params, fmt.Sprintf("Test id %d failed: expected non-nil params, but got nil", tc.id))
			assert.Equal(t, tc.output, result.fullValue, fmt.Sprintf("Test id %d failed: expected output %s, but got %s", tc.id, tc.output, result.fullValue))
			assert.Equal(t, tc.params, params, fmt.Sprintf("Test id %d failed: expected params %v, but got %v", tc.id, tc.params, params))
		}
	}
}

func TestInsertConflictParams(t *testing.T) {
	tree := NewNode("", false)

	tree.Insert("GET /users/:id")
	n := tree.Insert("GET /users/:username")

	assert.Nil(t, n)
}

func TestSearchConflictStaticAndParam(t *testing.T) {
	tree := NewNode("", false)

	tree.Insert("GET /users/:id")
	tree.Insert("GET /users/history")

	result, _ := tree.Search("GET /users/history")
	assert.NotNil(t, result)
	assert.Equal(t, "GET /users/history", result.fullValue)

	result, _ = tree.Search("GET /users/123")
	assert.NotNil(t, result)
	assert.Equal(t, "GET /users/:id", result.fullValue)
}
