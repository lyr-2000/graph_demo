package main

import (
	"aa/bfs/xcontainer"
	"context"
	"log"
	"sync"
	"time"
)

// import (
// 	"bfs"
// )
// 队列结构
var q xcontainer.Queue = xcontainer.New()

type Node struct {
	Func func(c context.Context, cancelFunc func()) string
	_x   string
}
type Graph map[string]Node

type RunSignals struct {
	Msg  chan string   //用来通知前端的信息
	Done chan []string //下一圈节点的 Nodeid
}

func MapSize(u *sync.Map) int {
	cnt := 0
	u.Range(func(k interface{}, v interface{}) bool {
		cnt++
		return true
	})
	return cnt
}
func _runNode(c context.Context, cancel func(), u Graph, nodeId string) {
	x, ok := u[nodeId]
	if ok {
		x.Func(c, cancel)
	}
}

var (
	mu sync.Mutex
)

func eventLoop(c context.Context, cancelFunc func(), u Graph) error {
	mu.Lock()
	cnt := q.Len()
	watchTask := make(chan struct{}, cnt)
	func() {
		if cnt <= 0 {
			return
		}
		// for v, _ := <- range watchTask {
		for i := 0; i < cnt; i++ {
			hnodeId, _ := q.Front().Value.(string)
			q.Remove()
			go func() {
				defer func() {
					watchTask <- struct{}{}
				}()
				_runNode(c, cancelFunc, u, hnodeId)
			}()
		}
	}()
	mu.Unlock()
	for i := 0; i < cnt; i++ {
		<-watchTask
		go eventLoop(c, cancelFunc, u)
	}
	return nil

}
func RunGraph(c context.Context, cancelFunc func(), u Graph, _startNode string) string {
	_, exists := graph[_startNode]
	if exists {
		q.Add(_startNode)
		// wg.Add(1)
		// c, cancel := context.WithCancel(context.TODO())
		go func() {
			eventLoop(c, cancelFunc, u)
			// log.Println("所有任务均执行完成，正式退出")

		}()

		// wg.Wait()
	}

	return "NULL"
}

// var (
// 	wg sync.WaitGroup
// )
var msg sync.Map
var graph = Graph{
	"start": Node{
		Func: func(c context.Context, cancelFunc func()) string {
			// defer wg.Done()
			msg.Store("result", "result[执行耗费时间较长的api节点,提前返回]")

			cancelFunc()
			time.Sleep(time.Second * 10)
			q.Add("js1")
			// q.Add("js3")
			// wg.Add(2)
			// next <- []string{"js2", "js3"}
			log.Println("node1 执行结束,10秒")
			return "ok"
		},
		_x: "",
	},
	"js1": Node{
		Func: func(c context.Context, cancelFunc func()) string {
			msg.Store("result", "执行耗费时间较长的api节点,提前返回")

			cancelFunc()
			time.Sleep(time.Second * 10)
			q.Add("js2")
			q.Add("js3")
			// wg.Add(2)
			// next <- []string{"js2", "js3"}
			log.Println("node1 执行结束,10秒")
			return "ok"
		},
		_x: "",
	},
	"js2": Node{
		Func: func(c context.Context, cancelFunc func()) string {
			// defer wg.Done()
			// 执行耗时操作
			time.Sleep(time.Second * 7)
			log.Println("node2 执行结束, 8秒")
			return "ok"
		},
		_x: "",
	},
	"js3": Node{
		Func: func(c context.Context, cancelFunc func()) string {
			// defer wg.Done()
			log.Println("node3 执行结束")
			return "ok"
		},
		_x: "",
	},
}

func main() {

	log.Println("hello-world")
	// msgs := make(chan string, 1)
	go func() { //模拟一次 http请求
		c, cancel := context.WithCancel(context.TODO())
		RunGraph(c, cancel, graph, "start")

		select {
		case <-c.Done():
			// return "node ok"
			msg.Range(func(k, v interface{}) bool {
				log.Println(k, "msg get := ", v)
				return true
			})
		case <-time.After(3 * time.Second):
			// return "运行超时"
			log.Println("!!@@@运行超时")
		}
	}()
	// log.Println(<-msgs)
	log.Println("--- 程序启动 ---- ")
	// select {}
	i := 0
	for {
		time.Sleep(time.Second)
		i++
		log.Printf("运行 %d 秒", i)
	}
}
