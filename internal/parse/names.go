package parse

import (
	"strings"
	"fmt"
)

type Scope struct {
	names map[Ident]uint64
}

func (scope *Scope) tryAddDep(id Ident, deps *map[uint64]bool) bool {
	i, ok := scope.get(id)
	if !ok {
		return false
	}
	deps[i] = true
	return true
}

type ParseNode struct {
	deps map[uint64]bool
	visited bool
	temp bool

	ast_node *ParseUnit
}

type ParseUnit interface {
	GetName() Ident
	GetChildren() []ParseUnit
	GetDeps(deps *map[uint64]bool, scope *Scope) bool
}

type ParseOrder struct {
	scope Scope
	nodes_underlying []ParseNode
	nodes_sorted []uint64
}

func (po *ParseOrder) addNames(unit *ParseUnit) bool {
	ident := unit.GetName()
	_, ok := po.scope.names[Ident]

	my_id := len(po.nodes_underlying)
	po.nodes_underlying = append(po.nodes_underlying, ParseNode{
		deps: make(map[uint64]bool)
		ast_node: unit
	})
	po.scope.names[Ident] = my_id

	for _, c := unit.GetChildren() {
		ok &= po.addNames(c)
	}
	return ok
}

func GetParseOrder(c *Contract) ParseOrder {
	po := ParseOrder{
		scope: Scope{
			names: make(map[Ident]uint64)
		}
	}
	ok := true
	for _, s := c.spaces {
		ok &= po.addNames(s)
	}
	for _, a := c.agents {
		ok &= po.addNames(a)
	}
	for _, p := c.paths {
		ok &= po.addNames(p)
	}

	if !ok {
		fmt.Print("there were duplicate identifiers")
		// todo idk
	}

	po.nodes_sorted = make([]uint64, len(po.nodes_underlying))

	for _, n := po.nodes_underlying {
		n.ast_node.GetDeps(&n.deps, po.scope)
	}

	// topological sort
	// todo
}
