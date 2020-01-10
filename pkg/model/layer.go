package model

import (
	"gorgonia.org/tensor"
	g "gorgonia.org/gornonia"
)

// Layer in a network.
type Layer interface {
	// Fwd is a foward pass through the layer.
	Fwd(x *g.Node) (*g.Node, error)

	// Learnables returns all learnable nodes within this layer.
	Learnables() g.Nodes
}

// ActivationFn is an activation funciton for a layer.
type ActivationFn func(x *g.Node) (*g.Node, error)
