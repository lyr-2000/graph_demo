package graph

import (
	"bufio"
	"bytes"
	"strings"
)

func newIdNode(id string) *Node {
	return &Node{NodeId: id}
}
func buildGraphByEdges_debug(edges []*Edge) *Graph {
	g := newGraph()
	for _, edge := range edges {
		from := edge.From
		to := edge.To
		g.allNodes[from.NodeId] = from
		g.allNodes[to.NodeId] = to
		linkEdge := &Edge{From: from, To: to, Cond: true}
		g.outEdges[from.NodeId] = append(g.outEdges[from.NodeId], linkEdge)
		g.inEdges[to.NodeId] = append(g.outEdges[from.NodeId], linkEdge)
		//记录出边和入边
	}
	return g
}

func readEdge_forDebug(s string) []*Edge {
	// var a, b string
	var list []*Edge
	// fmt.Sscanf(s, "%c --> %c", &a, &b)
	buf := bytes.NewBufferString(s)
	reader := bufio.NewReader(buf)
	var nodes = map[string]*Node{}
	for {
		// var line string
		// _, err := fmt.Fscanf(buf, "%s\n", &line)
		// _, err := fmt.Fscanf(buf, "(%s-->%s)\n", &a, &b)
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, "#") {
			continue
		}
		two := strings.Split(line, "-->")
		if len(two) != 2 {
			// log.Fatalln("cannot parse your expr")
			continue
		}
		a_ := strings.Trim(two[0], "\n ")
		b_ := strings.Trim(two[1], "\n ")
		// if strings.HasPrefix(a_, "#") {
		// 	continue
		// }
		var a, b *Node
		var ok bool
		a, ok = nodes[a_]
		if !ok {
			a = newIdNode(a_)
		}
		b, ok = nodes[b_]
		if !ok {
			b = newIdNode(b_)
		}
		list = append(list, &Edge{From: a, To: b})
	}
	return list
}

func buildGraph_debug_by_Expr(s string) *Graph {
	/*
		a --> b
		 b --> c
			b --> d
			b --> e
		e --> f
	*/
	edges := readEdge_forDebug(s)
	g := buildGraphByEdges_debug(edges)
	return g
}

// func buildGraph_debug_by_Expr(s string) *Graph {
// 	/*
// 		a --> b
// 		 b --> c
// 			b --> d
// 			b --> e
// 		e --> f
// 	*/
// 	edges := readEdge(s)
// 	g := buildGraphByEdges_debug(edges)
// 	return g
// }

// func travel()
