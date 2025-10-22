package parse

import (
	"fmt"
)

type Scope struct {
	names map[Ident]uint64
}

func (scope *Scope) tryAddDep(id Ident, deps *map[uint64]bool) bool {
	i, ok := scope.names[id]
	if !ok {
		fmt.Printf("Undeclared Identifier: %s\n", id.toString())
		return false
	}
	(*deps)[i] = true
	return true
}

type ParseNode struct {
	deps map[uint64]bool
	visited bool
	temp bool

	ast_node ParseUnit
}

type ParseUnit interface {
	GetName() Ident
	GetChildren() []ParseUnit
	ParseDepGetter
}

type ParseOrder struct {
	scope Scope
	nodes_underlying []ParseNode
	nodes_sorted []uint64
	nodes_sorted_index int
}

func (po *ParseOrder) topSortVisit(id uint64) bool {
	n := &po.nodes_underlying[id]
	if n.visited {
		return true
	}
	if n.temp { // cycle !
		fmt.Printf("Node %s is in a cycle\n", n.ast_node.GetName().toString())
		return false
	}
	n.temp = true
	for dep_id := range n.deps {
		po.topSortVisit(dep_id)
	}
	// fmt.Printf("adding node %s\n", n.ast_node.GetName().toString())
	n.visited = true
	po.nodes_sorted[po.nodes_sorted_index] = id
	po.nodes_sorted_index += 1
	return true
}

func (po *ParseOrder) addNames(unit ParseUnit) bool {
	ident := unit.GetName()
	_, dupes := po.scope.names[ident]

	if dupes {
		fmt.Printf("Duplicate Identifier: %s\n", ident.toString())
	}

	my_id := len(po.nodes_underlying)
	po.nodes_underlying = append(po.nodes_underlying, ParseNode{
		deps: make(map[uint64]bool),
		ast_node: unit,
	})
	po.scope.names[ident] = uint64(my_id)

	for _, c := range unit.GetChildren() {
		dupes = po.addNames(c) || dupes
	}
	return dupes
}

func GetParseOrder(c *Contract) ParseOrder {
	po := ParseOrder{
		scope: Scope{
			names: make(map[Ident]uint64),
		},
	}
	dupes := false
	for _, s := range c.spaces {
		dupes = po.addNames(&s) || dupes
	}
	for _, a := range c.agents {
		dupes = po.addNames(&a) || dupes
	}
	for _, p := range c.paths {
		dupes = po.addNames(&p) || dupes
	}

	if dupes {
		fmt.Println("there were duplicate identifiers")
		// todo idk
	}

	for _, n := range po.nodes_underlying {
		n.ast_node.GetDeps(&n.deps, &po.scope)
	}

	// po.printDeps()

	// topological sort
	po.nodes_sorted_index = 0
	po.nodes_sorted = make([]uint64, len(po.nodes_underlying))

	for i, n := range po.nodes_underlying {
		if !n.visited {
			if !po.topSortVisit(uint64(i)) {
				po.nodes_sorted = po.nodes_sorted[:po.nodes_sorted_index]
				return po
			}
		}
	}

	po.printDepsOrdered()

	return po
}

// func (po *ParseOrder) printDeps() {
// 	for i, n := range po.nodes_underlying {
// 		fmt.Printf("%d %s:\t", i, n.ast_node.GetName().toString())
// 		for dep, _ := range n.deps {
// 			fmt.Printf("%d ", dep)
// 		}
// 		fmt.Println("")
// 	}
// }

func (po *ParseOrder) printDepsOrdered() {
	fmt.Printf("\n%+v\n", po.nodes_sorted)
	for _, i := range po.nodes_sorted {
		n := po.nodes_underlying[i]
		fmt.Printf("%d %s:\t", i, n.ast_node.GetName().toString())
		for dep, _ := range n.deps {
			fmt.Printf("%d ", dep)
		}
		fmt.Println("")
	}
}
