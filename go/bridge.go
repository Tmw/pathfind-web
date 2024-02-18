package main

import (
	"encoding/json"
	"fmt"

	"syscall/js"
)

var (
	defaultGridWidth  = 40
	defaultGridHeight = 30
	defaultTileSize   = 30
	state             State
)
func initState() {
	state = NewState(defaultTileSize)
	state.Grid.Resize(defaultGridWidth, defaultGridHeight)
}

func main() {
	initState()

	// scope all functions under window.bridge
	bridge := map[string]interface{}{
		"solve":     makeSolveFunc(),
		"clearPath": makeClearPathFunc(),
		"reset":     makeResetFunc(),
		"getState":  makeGetStateFunc(),
		"mouseMove": makeMouseMoveFunc(),
		"mouseDown": makeMouseDownFunc(),
		"mouseUp":   makeMouseUpFunc(),
		"resize":    makeResizeFunc(),
	}

	js.Global().Set("bridge", js.ValueOf(bridge))

	// ensure Go doesn't exit before we have a chance to call it from JS.
	<-make(chan struct{})
}

type SolveRequest struct {
	GridWidth  int          `json:"gridWidth"`
	GridHeight int          `json:"gridHeight"`
	Grid       [][]TileType `json:"grid"`
}

func makeSolveFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		state.Grid.Solve()
		return nil
	})
}

func makeClearPathFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		state.Grid.ClearPath()
		return nil
	})
}

func makeResetFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		initState()
		return nil
	})
}

func makeGetStateFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		r, err := json.Marshal(state)
		if err != nil {
			fmt.Println("error marshalling json.", err)
		}

		return string(r)
	})
}

func makeResizeFunc() js.Func {
	type ResizeRequest struct {
		Width    int `json:"width"`
		Height   int `json:"height"`
		TileSize int `json:"tile_size"`
	}

	return js.FuncOf(func(this js.Value, args []js.Value) any {
		var req ResizeRequest

		if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
			fmt.Println("error parsing json", err)
			return nil
		}

		state.TileSize = req.TileSize
		state.Grid.Resize(req.Width, req.Height)

		return nil
	})
}

func makeMouseMoveFunc() js.Func {
	type MouseMoveRequest struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	return js.FuncOf(func(this js.Value, args []js.Value) any {
		var req MouseMoveRequest

		if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
			fmt.Println("error parsing json", err)
			return nil
		}

		state.OnMouseMove(req.X, req.Y)
		return nil
	})
}

func makeMouseDownFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		state.OnMouseDown()
		return nil
	})
}

func makeMouseUpFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		state.OnMouseUp()
		return nil
	})
}
