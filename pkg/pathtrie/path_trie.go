// Copyright © 2021 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Inspired by https://github.com/akitasoftware/akita-libs/blob/main/path_trie/path_trie.go
package pathtrie

import (
	"strings"

	"github.com/apiclarity/speculator/pkg/utils"
)

type TrieNode struct {
	Children PathToTrieNode

	// Name of the path segment corresponding to this node.
	// E.g. if this node represents /v1/foo/bar,
	// the Name would be "bar" and the FullPath would be "/v1/foo/bar".
	Name string

	// FullPath includes the node's name and uniquely identifies the node in the tree.
	FullPath string

	// PathParamCounter counts the amount of path params in the FullPath
	PathParamCounter int

	// Value of the full path
	Value interface{}
}

type PathToTrieNode map[string]*TrieNode

type PathTrie struct {
	Trie          PathToTrieNode
	PathSeparator string
}

type ValueMergeFunc func(existing, newV *interface{})

// Create a PathTrie with "/" as the path separator.
func New() PathTrie {
	return NewWithPathSeparator("/")
}

// Create a PathTrie with a user-supplied path separator.
func NewWithPathSeparator(pathSeparator string) PathTrie {
	return PathTrie{
		Trie:          make(PathToTrieNode),
		PathSeparator: pathSeparator,
	}
}

// Insert val at path, with path segments separated by PathSeparator.
// Returns true if a new path was created, false if an existing path
// was overwritten.
//
func (pt *PathTrie) Insert(path string, val interface{}) bool {
	return pt.InsertMerge(path, val, func(existing, newV *interface{}) {
		*existing = *newV
	})
}

// Insert val at path, with path segments separated by PathSeparator.
// Returns true if a new path was created, false if an existing path
// was overwritten.
//
// The merge function is responsible for updating the existing value
// with the new value.
func (pt *PathTrie) InsertMerge(path string, val interface{}, merge ValueMergeFunc) (isNewPath bool) {
	trie := pt.Trie
	isNewPath = true
	// TODO: what about path that ends with pt.PathSeparator is it different ?
	segments := strings.Split(path, pt.PathSeparator)

	// Traverse the Trie along path, inserting nodes where necessary.
	for idx, segment := range segments {
		isLastSegment := idx == len(segments)-1
		if node, ok := trie[segment]; ok {
			if isLastSegment {
				// If this is the last path segment, then this is the node to update.
				// If node value is not empty it means that an existing path is overwritten
				isNewPath = utils.IsNil(node.Value)
				merge(&node.Value, &val)
			} else {
				// Otherwise, continue descending.
				trie = node.Children
			}
		} else {
			newNode := pt.createPathTrieNode(segments, idx, isLastSegment, val)
			trie[segment] = newNode
			trie = newNode.Children
		}
	}

	return isNewPath
}

func (pt *PathTrie) createPathTrieNode(segments []string, idx int, isLastSegment bool, val interface{}) *TrieNode {
	fullPathSegments := segments[:idx+1]
	node := &TrieNode{
		Children: make(PathToTrieNode),
		Name:     segments[idx],
		FullPath: strings.Join(fullPathSegments, pt.PathSeparator),
	}
	node.PathParamCounter = countPathParam(fullPathSegments)
	if isLastSegment {
		node.Value = val
	}

	return node
}

func countPathParam(segments []string) int {
	count := 0

	for _, segment := range segments {
		if utils.IsPathParam(segment) {
			count += 1
		}
	}

	return count
}

// GetValue returns the given node path value, nil if node is not found.
func (pt *PathTrie) GetValue(path string) interface{} {
	node := pt.getNode(path)
	if node == nil {
		return nil
	}

	return node.Value
}

// GetPathAndValue returns the given node full path and value, nil if node is not found.
func (pt *PathTrie) GetPathAndValue(path string) (string, interface{}, bool) {
	node := pt.getNode(path)
	if node == nil {
		return "", nil, false
	}

	return node.FullPath, node.Value, true
}

func (pt *PathTrie) getNode(path string) *TrieNode {
	segments := strings.Split(path, pt.PathSeparator)

	nodes := pt.Trie.getMatchNodes(segments, 0)

	if len(nodes) == 0 {
		return nil
	}

	if len(nodes) == 1 {
		return nodes[0]
	}

	// if multiple nodes found, return the node with less path params segments
	return getMostAccurateNode(nodes, path, len(segments))
}

func (trie PathToTrieNode) getMatchNodes(segments []string, idx int) []*TrieNode {
	var nodes []*TrieNode

	isLastSegment := idx == len(segments)-1

	for _, node := range trie {
		// Check for node segment match
		if !node.isNameMatch(segments[idx]) {
			continue
		}

		// If this is the last path segment, then return node if it holds a value.
		if isLastSegment {
			if node.Value != nil {
				nodes = append(nodes, node)
			}
			continue
		}

		// Otherwise, continue descending.
		newNodes := node.Children.getMatchNodes(segments, idx+1)
		if len(newNodes) > 0 {
			nodes = append(nodes, newNodes...)
		}
	}

	return nodes
}

// getMostAccurateNode returns the node with less path params segments.
func getMostAccurateNode(nodes []*TrieNode, path string, segmentsLen int) *TrieNode {
	var retNode *TrieNode
	minPathParamSegmentsCount := segmentsLen + 1

	for _, node := range nodes {
		if node.isFullPathMatch(path) {
			// return exact match
			return node
		}

		// TODO: if node.PathParamCounter == minPathParamSegmentsCount
		if node.PathParamCounter < minPathParamSegmentsCount {
			// found more accurate node
			minPathParamSegmentsCount = node.PathParamCounter
			retNode = node
		}
	}

	return retNode
}

func (node *TrieNode) isNameMatch(segment string) bool {
	if utils.IsPathParam(node.Name) {
		return true
	}

	if node.Name == segment {
		return true
	}

	return false
}

func (node *TrieNode) isFullPathMatch(path string) bool {
	return node.FullPath == path
}
