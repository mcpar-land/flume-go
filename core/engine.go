package core

import (
	"context"
	"fmt"
)

type Engine struct {
	Blueprint
	NodeHandlers
}

// memoizer?

func (e Engine) ResolveRoot(input NodeData, ctx context.Context) (NodeData, error) {
	root, err := e.Blueprint.Root()
	if err != nil {
		return nil, err
	}
	return e.RecusiveResolveNode(root, ctx, 1)
}

const MAX_DEPTH int = 10_000

func (e Engine) RecusiveResolveNode(n *Node, ctx context.Context, depth int) (NodeData, error) {
	if depth > MAX_DEPTH {
		return nil, fmt.Errorf("Max node depth exceeded! Check for a loop in your nodes.")
	}
	inputConnections, err := e.Blueprint.GetConnections(n.Connections.Inputs)
	if err != nil {
		return nil, err
	}
	input := NodeData{}
	for k, conns := range inputConnections {
		res, err := e.RecusiveResolveNode(conns[0], ctx, depth+1)
		if err != nil {
			return nil, err
		}
		input[k] = res[k]
	}
	return e.ResolveNode(n, input, ctx)
}

func (e Engine) ResolveNode(n *Node, input NodeData, ctx context.Context) (NodeData, error) {
	if n == nil {
		return nil, fmt.Errorf("nil node supplied to GetHandler")
	}
	if handler, err := e.GetHandler(n); err == nil {
		res, err := (*handler)(n, input, ctx)
		if err != nil {
			return nil, fmt.Errorf("Err in node type %v: %v", n.Type, err)
		}
		return res, nil
	} else {
		return nil, err
	}
}

func (e *Engine) MemoizeCall(n *Node, input, output NodeData) {

}

func (e Engine) GetHandler(n *Node) (*NodeHandler, error) {
	if n == nil {
		return nil, fmt.Errorf("nil node supplied to GetHandler")
	}
	if h, ok := e.NodeHandlers[n.Type]; ok {
		return &h, nil
	}
	return nil, fmt.Errorf("No handler for node type '%v'", n.Type)
}

type NodeHandler func(
	node *Node, input NodeData, ctx context.Context,
) (NodeData, error)

type NodeHandlers map[string]NodeHandler
