/*
Copyright © 2024 Piotr Tobiasz
*/
package main

import (
	"starling/cmd"
	_ "starling/cmd/server"
	_ "starling/cmd/worker"
)

func main() {
	cmd.Execute()
}
