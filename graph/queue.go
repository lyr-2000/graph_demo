package graph

import "time"

type ScriptNode struct {
}

func (u *ScriptNode) Run() string {

	time.Sleep(time.Second * 3)
	return ""
}

type ValidNode struct {
}

func (u *ValidNode) Run() string {

	return "校验不通过"
}
