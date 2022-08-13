package graph

import (
	"context"
	"fmt"
)

//将信号和 异常组合在一起
type GraphTravelError struct {
	Signal error       //用来通知上层调用 的逻辑
	Err    error       // 执行用户代码出现的异常
	Value  interface{} //其他信息
	At     *Node       // 发生异常的节点
}

type Node struct {
	NodeId         string                 `json:"nodeId"`
	RawJsonContent string                 `json:"-"`    // 数据库前端存的数据，有可能有其他字段，可以拿到字段并且解析
	Attr           map[string]interface{} `json:"attr"` // 节点的属性，尽量不要在节点上面挂载数据或者变量，导致污染严重
	// IsBegin        bool                   // 前端发出申请的时候，从某个节点开始，于是从这个节点出发
	Type string `json:"type"`
}
type Edge struct {
	Cond           bool // 是否通过
	CondExpr       string
	FlowProperties string //1为同意，否则为拒绝
	From           *Node
	To             *Node
	RawJsonContent string
}

func (u *Edge) String() string {
	return fmt.Sprintf("edge{%v --> %v}", u.From.NodeId, u.To.NodeId)
}
func newGraph() *Graph {
	g := new(Graph)
	g.allNodes = map[string]*Node{}
	g.outEdges = map[string][]*Edge{}
	g.inEdges = map[string][]*Edge{}
	g._dfsVisitedNode = map[string]struct{}{}
	// g._dfsBeginNode = nil
	return g
}

type Graph struct {
	// allNodes []*Node
	allNodes map[string]*Node
	outEdges map[string][]*Edge // 出边
	inEdges  map[string][]*Edge //入边
	//下面变量禁止修改
	_dfsVisitedNode       map[string]struct{} //用来判断环状结构，出现环立刻返回
	_dfsBeginNode         *Node               //[当前请求开始节点] 用小写，防止外部调用修改这个变量
	_dfsCurrentTravelNode *Node               //当前方法遍历到的这个节点 【禁止修改变量，当出现panic的时候用来获取出异常在哪个点】
	_dfsTrace             []*Node             //深度遍历时候记录节点访问的路径
}

func (g *Graph) FindNodeByNodeId(nodeid string) *Node {
	enode := g.allNodes[nodeid]
	return enode
}
func (u *Graph) GetPreviousEdges(n *Node) []*Edge {
	return u.inEdges[n.NodeId]
}

// func EnsureNoCycle(g *Graph, beginNode *Node) error {
// 	if beginNode == nil || g == nil {
// 		return nil
// 	}
// 	nnext := g.outEdges[beginNode.NodeId]
// 	var mp = map[string] struct{}{}
// 	var dfs_cycle func(g *Graph , n *Node) {

// 	}
// 	for _, v := range nnext {

// 	}
// }
func (u *Graph) Range(callback func(n *Node) error) error {
	for _, node := range u.allNodes {
		err := callback(node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Graph) GetBeginNode() *Node {
	return u._dfsBeginNode
}

func (u *Graph) GetNextGroupNode(curNode *Node) []*Node {
	nextNodesEdges := u.outEdges[curNode.NodeId]
	var list []*Node
	for _, v := range nextNodesEdges {
		if v.Cond {
			list = append(list, v.To)
		}
	}
	return list
}

//这一轮的终点，遇到就停止
func IsCurrentDestinationNode(n *Node) bool {

	return true
}

type Context = context.Context

// type CallbackFunc func() error
type CallbackRequest struct {
	OnStart         []func(Context, *Graph, *Node) error //这个可能没啥用，先放着
	OnBegin         []func(Context, *Graph, *Node) error //当前请求的开始节点
	OnPanic         []func(Context, *Graph, *Node, error) error
	OnVisitNode     []func(Context, *Graph, *Node) error // 遍历到当前节点
	OnReachDestNode []func(Context, *Graph, *Node) error //遍历到当前的停止节点
	OnReachEndNode  []func(Context, *Graph, *Node) error //遍历到 end节点
}

func IsStartNode(u *Node) bool {
	if u == nil {
		return false
	}
	return u.Type == "start"
}
func IsEndNode(u *Node) bool {
	return u.Type == "end"
}

//reach the destination
func IsReachDestNode(u *Node) bool {
	//如果是审核节点，就需要停下来进行审核 【 如果 是end节点，也要停下来，  然后最后判断整个流程是否结束】
	return u.Type == "end" || u.Type == "userTask" //用户审核 ，或者是某个特殊的网关，需要卡在这
}
func TravelNode(c Context, g *Graph, begin *Node, callback *CallbackRequest) error {
	g._dfsBeginNode = begin
	err := travelNode0(c, g, callback)
	if err == nil {
		return nil
	}
	switch err {
	case SkipCurrentNodeCallback:
		//
	case Done:
	default:
		return err
	}
	return err
}

//获取 深度遍历时候的 trace
func GetDfsTrace(g *Graph) []*Node {
	return g._dfsTrace
}

// func (u *Graph) TravelNode
func travelNode0(c context.Context, g *Graph, callback *CallbackRequest) (err error) {
	// var currentNode *Node

	defer func() {
		if erri := recover(); erri != nil {
			if err0, iserr := erri.(error); iserr {
				err = err0
			} else {
				err = fmt.Errorf("%v", erri)
			}
			for _, callback := range callback.OnPanic {
				err = callback(c, g, g._dfsCurrentTravelNode, err)
				if err == nil {
					break // 消费了这个error, 如果还返回error，就进入下一个
				}
			}
		}
	}()
	//递归调用

	beginNode := g._dfsBeginNode
	g._dfsCurrentTravelNode = beginNode //进入下一个节点的时候要修改变量
	if beginNode == nil {
		return fmt.Errorf("no begin node")
	}
	if IsStartNode(beginNode) {
		//startNode,触发回调
	onStartLoop:
		for _, onStart := range callback.OnStart {

			if err := onStart(c, g, beginNode); err != nil {
				switch err {
				case Done: //不做任何事情
					break onStartLoop
				case SkipCurrentNodeCallback: //跳过所有的回调
					break onStartLoop
				default:
					return err
				}
			}
		}
	}
	if IsEndNode(beginNode) {
	onEndLoop:
		for _, cb := range callback.OnReachEndNode {
			if err := cb(c, g, beginNode); err != nil {
				switch err {
				case Done: //不做任何事情
					break onEndLoop
				case SkipCurrentNodeCallback: //跳过所有的回调
					break onEndLoop
				default:
					return err
				}
			}
		}
		//开始就是 endNode ,就直接结束
		return Done
	}
onBeginLoop:
	for _, onbegin := range callback.OnBegin {
		if err := onbegin(c, g, beginNode); err != nil {
			switch err {
			case Done: //不做任何事情
				break onBeginLoop
			case SkipCurrentNodeCallback: //跳过所有的回调
				break onBeginLoop
			default:
				return err
			}
		}
	}
	return dfs(c, g, beginNode, 0, callback)
	// is end ,and callback on end
}
func (g *Graph) SetVisitedNode(n *Node, isVisited bool) {
	if isVisited {
		g._dfsVisitedNode[n.NodeId] = struct{}{}
	} else {
		delete(g._dfsVisitedNode, n.NodeId)
	}

}
func (g *Graph) IsVisited(n *Node) bool {
	if n == nil {
		return false
	}
	_, ok := g._dfsVisitedNode[n.NodeId]
	return ok
}

// 访问节点时候用到的哨兵异常
var (
	CycleNodeError          = fmt.Errorf("cycleNode 环状结构,流程异常")
	StackOverflow           = fmt.Errorf("stackoverflow")
	Done                    = fmt.Errorf("done")
	SkipCurrentNodeCallback = fmt.Errorf("skip current callback") //到达了停止节点，例如审核节点，当前的路径就停止了，等待下一轮请求

	MaxTravelStack uint8 = 100
)

func dfs(c context.Context, g *Graph, n *Node, step uint8 /*最大256层*/, callback *CallbackRequest) (err error) {
	if n == nil {
		return nil
	}
	if step >= MaxTravelStack {
		// 递归栈超过 MaxTravelStack 层，说明程序有严重问题，例如环状结果，需要退出
		return StackOverflow
	}
	if g.IsVisited(n) { //环形检测
		return CycleNodeError //检测环形节点
	}
	g.SetVisitedNode(n, true)
	//当前遍历到的节点
	g._dfsCurrentTravelNode = n
	g._dfsTrace = append(g._dfsTrace, n)
	traceLen := len(g._dfsTrace)
	defer func() {
		g.SetVisitedNode(n, false) //remove mark visited
		g._dfsTrace = g._dfsTrace[:traceLen-1]
	}()
	nextGroups := g.GetNextGroupNode(n)
	// RunNextNode:
	for _, vnextNode := range nextGroups { //如果 是 排他网关，只会数组长度为1， 如果是并行网关，将返回所有符合条件的后续节点

		if IsReachDestNode(vnextNode) {
		callback_destNodes:
			for _, cb := range callback.OnReachDestNode { // stop node

				if err := cb(c, g, vnextNode); err != nil {
					switch err {
					case Done: //不做任何事情
						break callback_destNodes
					case SkipCurrentNodeCallback: //跳过所有的回调
						break callback_destNodes
					default:
						return err
					}
				}
			}

			// continue RunNextNode
		}
		if IsEndNode(vnextNode) {

			for _, cb := range callback.OnReachEndNode { // end node
				if err := cb(c, g, vnextNode); err != nil {
					switch err {
					case Done:
						break
					case SkipCurrentNodeCallback:
						break
					default:
						return err
					}
				}

			}

			//已经是结束节点了，直接返回
			return nil
		}
		//先序
	Loop_applyVisited:
		for _, onVisitedFunc := range callback.OnVisitNode {
			if err := onVisitedFunc(c, g, vnextNode); err != nil {
				switch err {
				case Done:
					break Loop_applyVisited
				case SkipCurrentNodeCallback:
					break Loop_applyVisited
				default:
					return err
				}
			}
		}

		//后序
		if err = dfs(c, g, vnextNode, step+1, callback); err != nil {
			return err
		}
	}

	return nil
}
