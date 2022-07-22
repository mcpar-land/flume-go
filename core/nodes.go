package core

import (
	"encoding/json"
	"fmt"
)

type Blueprint map[string]Node

func BlueprintFromJson(src string) (Blueprint, error) {
	var b Blueprint
	err := json.Unmarshal([]byte(src), &b)
	return b, err
}

func (b Blueprint) Root() (*Node, error) {
	var root *Node
	for _, n := range b {
		if root != nil {
			root = &n
		} else {
			return nil, fmt.Errorf("Blueprint has more than one root node")
		}
	}
	if root == nil {
		return nil, fmt.Errorf("Blueprint does not have a root node")
	}
	return root, nil
}

func (b Blueprint) GetNode(id string) (*Node, error) {
	if node, ok := b[id]; ok {
		return &node, nil
	}
	return nil, fmt.Errorf("Cannot find node %v", id)
}

func (b Blueprint) GetConnections(
	c NodeConnectionsMap,
) (map[string][]*Node, error) {
	nodes := map[string][]*Node{}
	for portName, connections := range c {
		cnPtrs := []*Node{}
		for _, connection := range connections {
			if node, err := b.GetNode(connection.NodeId); err == nil {
				cnPtrs = append(cnPtrs, node)
			} else {
				return nodes, err
			}
		}
		nodes[portName] = cnPtrs
	}
	return nodes, nil
}

type Node struct {
	Id          string          `json:"id"`
	X           float64         `json:"x"`
	Y           float64         `json:"y"`
	Type        string          `json:"type"`
	Width       float64         `json:"width"`
	Root        bool            `json:"root"`
	Connections NodeConnections `json:"connections"`
	InputData   NodeData        `json:"inputData"`
}

type NodeData map[string]map[string]any

type NodeConnections struct {
	Inputs  NodeConnectionsMap `json:"inputs"`
	Outputs NodeConnectionsMap `json:"outputs"`
}

type NodeConnectionsMap map[string][]NodeConnection

type NodeConnection struct {
	NodeId   string `json:"nodeId"`
	PortName string `json:"portName"`
}
