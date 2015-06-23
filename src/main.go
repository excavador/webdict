package main

import (
	"webdict"
)

func main() {
	api := webdict.NewApi("/dictionary")
	api.Run()
}
