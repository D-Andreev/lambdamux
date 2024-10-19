package radix

import (
	"context"
	"log/slog"
	"slices"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Node struct {
	edges      []*Node // sorted in ascending order
	isComplete bool
	value      string
	fullValue  string // used only when returning a node from Search method, otherwise it's not populated
	isParam    bool
	paramNames []string
	Handler    func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

// NewNode creates a new node
func NewNode(value string, isComplete bool) *Node {
	isParam := strings.Contains(value, ":")
	paramNames := getParamNames(value)

	return &Node{
		edges:      []*Node{},
		isComplete: isComplete,
		value:      value,
		isParam:    isParam,
		paramNames: paramNames,
		Handler:    nil,
	}
}

// getParamNames returns the names of the parameters in the given value
func getParamNames(value string) []string {
	var paramNames []string
	segments := strings.Split(value, "/")
	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			paramNames = append(paramNames, segment[1:])
		}
	}
	return paramNames
}

// InsertWithHandler inserts a new node in the tree with a handler
func (n *Node) InsertWithHandler(
	input string,
	handler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error),
) {
	node := n.Insert(input)
	if node == nil {
		return
	}
	node.Handler = handler
}

// Insert inserts a new node in the tree
func (n *Node) Insert(input string) *Node {
	node := n
	search := input
	for {
		// The key is exhausted and the deepest possible node found.
		if len(search) == 0 {
			node.isComplete = true
			return node
		}

		parent := node
		node = node.getEdge(search[0], false)

		// No matching edge was found. Just create new edge.
		if node == nil {
			node = NewNode(search, true)
			parent.addEdge(node)
			return node
		}

		commonPrefix := getCommonPrefix(search, node.value)

		// The current node's value is a fullValue of the search string. Go deeper with remainder.
		if commonPrefix == len(node.value) {
			search = search[commonPrefix:]
			continue
		}

		// Split node
		child := NewNode(search[:commonPrefix], false)

		// Check for param name conflict
		if node.isParam && child.isParam && !slicesEqual(child.paramNames, node.paramNames) {
			slog.Error("Route param conflict",
				"path", input,
				"conflicting_param", node.value,
				"existing_param", node.paramNames,
				"new_route", input,
			)
			return nil
		}
		parent.updateEdge(search[0], child)
		node.value = node.value[commonPrefix:]
		node.isParam = strings.Contains(node.value, ":")
		node.paramNames = getParamNames(node.value)
		child.addEdge(node)
		search = search[commonPrefix:]
		newNode := NewNode(search, true)
		if len(search) == 0 {
			child.isComplete = true
		} else {
			child.addEdge(newNode)
		}

		return newNode
	}
}

// slicesEqual checks if two slices of strings are equal
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	return slices.Equal(a, b)
}

// Search gets an item from the tree
func (n *Node) Search(input string) (*Node, map[string]string) {
	node := n
	search := input
	fullPath := ""
	params := map[string]string{}
	for {
		// If search string is empty, we've found the deepest possible node
		if len(search) == 0 {
			if !node.isComplete {
				return nil, nil
			}
			node.fullValue = fullPath
			return node, params
		}

		// Get the child node for the first character of the search string
		node = node.getEdge(search[0], true)

		// No match
		if node == nil {
			return nil, nil
		}

		// Handle param search
		if node.isParam {
			searchSegments := strings.Split(search, "/")
			nodeSegments := strings.Split(node.value, "/")

			if len(nodeSegments) > len(searchSegments) {
				return nil, nil
			}

			paramIndex := 0
			for i := range nodeSegments {
				if i < len(searchSegments) {
					if paramIndex > 0 && paramIndex > len(node.paramNames)-1 {
						break
					}
					if nodeSegments[i] != searchSegments[i] && !strings.HasPrefix(nodeSegments[i], ":") {
						return nil, nil
					}
					if strings.HasPrefix(nodeSegments[i], ":") {
						if paramIndex < len(node.paramNames) {

							params[node.paramNames[paramIndex]] = searchSegments[i]
							searchSegments[i] = nodeSegments[i]
							paramIndex++
						}
					}
				}
			}

			search = strings.Join(searchSegments, "/")
			fullPath += node.value

			if len(search) > len(node.value) {
				search = search[len(node.value):]
			} else {
				search = ""
			}
			continue
		}

		// Find the common prefix between the search string and the node's value
		commonPrefix := getCommonPrefix(search, node.value)
		fullPath += search[:commonPrefix]

		// If the common prefix length equals the node's value length, continue searching
		if commonPrefix == len(node.value) {
			search = search[commonPrefix:]
			continue
		}

		// No match
		return nil, nil
	}
}

// getFirstMatchIdx performs a binary search to find the index of the first edge
// that matches or exceeds the given label.
func (n *Node) getFirstMatchIdx(label byte) int {
	num := len(n.edges)
	return sort.Search(
		num, func(i int) bool {
			return n.edges[i].value[0] >= label
		},
	)
}

// addEdge inserts a new edge while maintaining sorted order
func (n *Node) addEdge(e *Node) {
	idx := n.getFirstMatchIdx(e.value[0])
	n.edges = append(n.edges, e)
	copy(n.edges[idx+1:], n.edges[idx:])
	n.edges[idx] = e
}

// updateEdge updates an existing edge with a new node
func (n *Node) updateEdge(label byte, node *Node) {
	idx := n.getFirstMatchIdx(label)
	if idx < len(n.edges) && n.edges[idx].value[0] == label {
		n.edges[idx] = node
		return
	}

	panic("We're trying to replace a missing node. This should never happen.")
}

// getEdge returns the edge that matches the given label
func (n *Node) getEdge(label byte, matchParam bool) *Node {
	idx := n.getFirstMatchIdx(label)

	if idx < len(n.edges) && n.edges[idx].value[0] == label {
		return n.edges[idx]
	}

	if !matchParam {
		return nil
	}

	// There was no exact match, so check for slug edges
	for _, e := range n.edges {
		nodeSegments := strings.Split(e.value, "/")
		if len(nodeSegments) == 0 {
			continue
		}
		if strings.Contains(nodeSegments[0], ":") {
			return e
		}
	}

	return nil
}

// getCommonPrefix returns the length of the longest common prefix between two strings
func getCommonPrefix(k1, k2 string) int {
	maxL := len(k1)
	if l := len(k2); l < maxL {
		maxL = l
	}
	var i int
	for i = 0; i < maxL; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

// GetAllCompleteItems returns all complete items in the tree
func (n *Node) GetAllCompleteItems() []string {
	var result []string

	result = dfsCompleteItems(n, "", result)
	sort.Strings(result)
	return result
}

// GetAllNodeValues returns all node values in the tree
func (n *Node) GetAllNodeValues() []string {
	var result []string

	result = dfs(n, result)
	sort.Strings(result)
	return result
}

// dfsCompleteItems performs a depth-first search to find all complete items
func dfsCompleteItems(node *Node, prefix string, result []string) []string {
	if node.isComplete {
		result = append(result, prefix)
	}

	for _, edge := range node.edges {
		result = dfsCompleteItems(edge, prefix+edge.value, result)
	}

	return result
}

// dfs performs a depth-first search to find all keys in the tree
func dfs(node *Node, result []string) []string {
	if len(node.value) > 0 {
		result = append(result, node.value)
	}

	for _, child := range node.edges {
		result = dfs(child, result)
	}

	return result
}
