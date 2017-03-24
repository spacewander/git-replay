package main

import (
	"github.com/Shopify/go-lua"
	"reflect"
	"strings"
)

var (
	vm *lua.State

	curCommit *CommitInfo
)

func init() {
	vm = lua.NewState()
	lua.OpenLibraries(vm)
	injectGitReplayApi(vm)
}

func dumpPlainStringStructToLuaTable(L *lua.State, obj interface{}) {
	ref := reflect.Indirect(reflect.ValueOf(obj))
	refType := ref.Type()
	for i := 0; i < ref.NumField(); i++ {
		refField := ref.Field(i)
		refFieldType := refType.Field(i)
		if refField.Type().Kind() == reflect.String {
			L.PushString(refFieldType.Name)
			L.PushString(refField.String())
			L.RawSet(-3)
		}
	}
}

func GetCommit(L *lua.State) int {
	debugLogger.Println("call get_commit")
	L.CreateTable(0, 16)
	dumpPlainStringStructToLuaTable(L, curCommit)

	if curCommit.tags != nil {
		// array-like subtable for tags
		L.PushString("tags")
		L.CreateTable(len(curCommit.tags), 0)
		for i, tagInfo := range curCommit.tags {
			// Note that Lua index starts from one
			L.PushNumber(float64(i + 1))
			L.CreateTable(8, 0)
			dumpPlainStringStructToLuaTable(L, tagInfo)
			L.RawSet(-3)
		}
		L.RawSet(-3)
	}
	return 1
}

func PrintToStoryView(L *lua.State) int {
	nargs := L.Top()
	debugLogger.Println("call print with nargs: ", nargs)
	args := []string{}
	for i := 1; i <= nargs; i++ {
		luaType := L.TypeOf(i)
		switch luaType {
		case lua.TypeTable:
			fallthrough
		case lua.TypeNumber:
			fallthrough
		case lua.TypeString:
			if s, ok := lua.ToStringMeta(L, i); ok {
				args = append(args, s)
			}
		case lua.TypeNil:
			args = append(args, "nil")
		case lua.TypeBoolean:
			if b := L.ToBoolean(i); b {
				args = append(args, "true")
			} else {
				args = append(args, "false")
			}
			// ignore invalid input
		}
	}
	storyView.Show(strings.Join(args, ""))
	return 0
}

func injectGitReplayApi(L *lua.State) {
	L.CreateTable(8, 0)
	// override lua print
	L.PushGoFunction(PrintToStoryView)
	L.SetGlobal("print")
	// inject apis
	L.PushGoFunction(PrintToStoryView)
	L.SetField(-2, "display")
	L.PushGoFunction(GetCommit)
	L.SetField(-2, "get_commit")

	// inject global api namespace
	L.Global("package")
	L.Field(-1, "loaded")
	L.PushValue(-3)
	L.SetField(-2, "git_replay")
	L.Pop(2)
	L.SetGlobal("git_replay")
}

func PlayWithCommitInfo(script string, commit *CommitInfo) error {
	curCommit = commit
	return lua.DoFile(vm, script)
}
