package main

import (
	"context"
	"log"
	"time"
)

// import (
// 	"bfs"
// )
// 队列结构
// var q xcontainer.Queue = xcontainer.New()

type Node struct {
	Func func(c context.Context, ch chan *RunSignals) string
	_x   string
}
type Graph map[string]Node

type RunSignals struct {
	Msg      string //用来通知前端的信息
	Err      error
	NextNode []string //下一圈节点的 Nodeid
	Done     bool
}

type ExitSignals struct {
	Err error
	Msg string
}

func runNode(c context.Context, ch chan *RunSignals, g Graph, nodeName []string) int {
	taskCnt := 0
	for _, v := range nodeName {
		n, exists := g[v]
		if exists {
			taskCnt++
			go func() {
				//这里不能panic,一定要保证节点不能panic,不然状态无法维护
				n.Func(c, ch)
				//这里可以包装一个 接口，不一定要传入 一个 channel,而是一个包装了channel的接口
			}()
		}
	}
	return taskCnt
}

func eventLoop(c context.Context, ch chan *ExitSignals, u Graph, _startNode string) string {
	// go eventLoop(c, ch, u, _startNode)
	signalQueue := make(chan *RunSignals, 1)
	//可以传入更多的信息，这里是最简化版本
	taskCnt := runNode(c, signalQueue, u, []string{_startNode})
	_notice := false
	//因为 eventLoop 是单线程执行的，可以不用加锁
	var send = func(x *ExitSignals) {
		if _notice {
			return
		}
		_notice = true
		ch <- x
	}
	for taskCnt > 0 {

		signal := <-signalQueue
		if signal.Done {
			taskCnt--
		}
		//如果节点执行err了，就不继续执行该分支
		if signal.Err != nil {
			send(&ExitSignals{
				Err: signal.Err,
				Msg: "执行节点出错，退出",
			})
			continue
		}
		if signal.Msg != "" {
			send(&ExitSignals{
				Err: nil,
				Msg: signal.Msg,
			})
		}
		nextChild := signal.NextNode
		if len(nextChild) > 0 {
			taskCnt += runNode(c, signalQueue, u, nextChild)
		}
		// if taskCnt == 0 {
		// 	break
		// }

	}
	log.Println("所有节点均已退出,结束任务!!! eventLoop")
	return ""
}
func RunGraph(c context.Context, ch chan *ExitSignals, u Graph, _startNode string) string {
	go eventLoop(c, ch, u, _startNode)
	return "NULL"
}

// var (
// 	wg sync.WaitGroup
// )

var graph = Graph{
	"start": Node{
		Func: func(c context.Context, ch chan *RunSignals) string {
			// defer wg.Done()
			time.Sleep(time.Millisecond * 300)
			ch <- &RunSignals{
				Msg: "start[进入耗时操作,通知前端返回]!",
			}
			defer func() {
				ch <- &RunSignals{
					Done:     true,
					NextNode: []string{"js1", "js2"},
				}
			}()
			time.Sleep(time.Second * 10)
			log.Println("node1 执行结束,10秒")
			return "ok"
		},
		// _x: "",
	},
	"js1": Node{
		Func: func(c context.Context, ch chan *RunSignals) string {
			// defer wg.Done()
			defer func() {
				ch <- &RunSignals{Done: true}
			}()
			time.Sleep(time.Second * 1)
			log.Println("js1 执行结束,10秒")
			return "ok"
		},
		_x: "",
	},
	"js2": Node{
		Func: func(c context.Context, ch chan *RunSignals) string {
			// defer wg.Done()
			defer func() {
				ch <- &RunSignals{
					Done:     true,
					NextNode: []string{"js3"},
				}
			}()
			time.Sleep(time.Second * 3)
			log.Println("js2 执行结束,3秒,最后一节点，js3")
			return "ok"
		},
		_x: "",
	},
	"js3": Node{
		Func: func(c context.Context, ch chan *RunSignals) string {
			// defer wg.Done()
			defer func() {
				ch <- &RunSignals{Done: true}
			}()
			log.Println("js3 执行结束,")
			return "ok"
		},
		_x: "",
	},
}

func main() {

	go func() { //模拟一次 http请求
		// c, cancel := context.WithCancel(context.TODO())
		ch := make(chan *ExitSignals, 1)
		RunGraph(context.TODO(), ch, graph, "start")
		select {
		case msg := <-ch: //监听协程传回的响应结果
			log.Printf("请求响应 %#v", msg)
		case <-time.After(time.Second * 8): //如果节点执行还是超时，就直接返回
			log.Println("请求超时！！！")
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
