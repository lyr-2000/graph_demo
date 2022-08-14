package graph

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func Test_callback(t *testing.T) {
	var orderId = 1
	var userId = "lyr1"
	var req = new(CallbackRequest)
	var callbackJsNode = func(c context.Context, g *Graph, n *Node, _ *Edge) (err error) {
		if n.Type != "script" {
			return nil
		}
		iuserId := c.Value("userId") //通过 context 获取userId，尽量避免这样做
		fmt.Printf("%+v, %v, %v\n", orderId, userId, iuserId)
		return nil
	}
	req.OnVisitNode = append(req.OnVisitNode, callbackJsNode)
	fmt.Printf("%+v\n", req)
}

func Test_read_edges(t *testing.T) {
	var s = `b --> c
c --> d
b --> e
b --> f
e --> d
f --> d
	`
	edges := readEdge_forDebug(s)
	t.Logf("edges = %+v\n", edges)

}

func Test_build_graph_byEdges(t *testing.T) {
	var s = `start-->nodea
nodea-->js
nodea-->user1
nodea-->user3
start-->user2
start-->a
user1-->start
a-->b
b-->c
# c-->start
c-->end
# b-->start
`
	edges := readEdge_forDebug(s)
	for _, v := range edges {
		if v.To.NodeId == "end" {
			//结束
			v.To.Type = "end"
		}
		if strings.HasPrefix(v.To.NodeId, "user") {
			//审核节点
			v.To.Type = "userTask"
		}
		if strings.HasPrefix(v.To.NodeId, "js") {
			//脚本节点
			v.To.Type = "js"
		}
	}
	// t.Logf("edges = %+v\n", edges)
	g := buildGraphByEdges_debug(edges)
	// g := buildGraph_debug_by_Expr(s)
	req := new(CallbackRequest)
	req.OnVisitNode = append(req.OnVisitNode, func(c Context, g *Graph, n *Node, _ *Edge) error {
		t.Logf("current Node is %+v\n", n.NodeId)
		return nil
	})
	req.OnReachDestNode = append(req.OnReachDestNode, func(c Context, g *Graph, n *Node, _ *Edge) error {
		t.Logf("reach destination node := %+v\n", n.NodeId)
		return nil
	})
	req.OnReachEndNode = append(req.OnReachEndNode, func(c Context, g *Graph, n *Node, _ *Edge) error {
		t.Logf("end !!%+v\n", n.NodeId)
		return nil
	})
	node := g.FindNodeByNodeId("start")
	err := TravelNode(context.TODO(), g, node, req)
	if err != nil {
		t.Logf("err = %+v\n", err)
	}
}

func Test_build_graph_user_to_start(t *testing.T) {
	var s = `user1-->a
a-->js
js-->b
js-->userTask1
b-->start
userTask1-->xyz
# userTask是审核节点,所有要先审核,后面的流程要卡在这不动

`
	//  user1 用户审核，到 a节点，a调用脚本，之后再回到 start 节点 [排他网关的时候]
	edges := readEdge_forDebug(s)

	t.Logf("all edges = [%+v]\n", edges)
	for _, v := range edges {
		if v.To.NodeId == "start" {
			v.To.Type = "start"
		}
		if v.To.NodeId == "end" {
			//结束
			v.To.Type = "end"
		}
		if strings.HasPrefix(v.To.NodeId, "user") {
			//审核节点
			v.To.Type = "userTask"
		}
		if strings.HasPrefix(v.To.NodeId, "js") {
			//脚本节点
			v.To.Type = "js"
		}
	}
	// t.Logf("edges = %+v\n", edges)
	g := buildGraphByEdges_debug(edges)
	// g := buildGraph_debug_by_Expr(s)
	req := new(CallbackRequest)
	req.OnVisitNode = append(req.OnVisitNode, func(c Context, g *Graph, n *Node, _ *Edge) error {
		t.Logf("current Node is %+v\n", n.NodeId)
		return nil
	})
	req.OnReachDestNode = append(req.OnReachDestNode, func(c Context, g *Graph, n *Node, e *Edge) error {
		t.Logf("reach destination node := %+v,edges=%v\n", n.NodeId, e)
		return nil
	})
	req.OnReachEndNode = append(req.OnReachEndNode, func(c Context, g *Graph, n *Node, _ *Edge) error {
		t.Logf("end !!%+v\n", n.NodeId)
		return nil
	})
	req.OnStart = append(req.OnStart, func(c Context, g *Graph, n *Node, _ *Edge) error {
		t.Logf("startNode !!%+v\n", n.NodeId)
		return nil
	})
	req.OnBegin = append(req.OnBegin, func(c Context, g *Graph, n *Node, _ *Edge) error {
		t.Logf("onBegin =%+v\n", n.NodeId)
		return nil
	})
	node := g.FindNodeByNodeId("user1")
	err := TravelNode(context.TODO(), g, node, req)
	if err != nil {
		t.Logf("err = %+v\n", err)
	}
}
