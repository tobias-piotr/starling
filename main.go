/*
Copyright © 2024 Piotr Tobiasz
*/
package main

import (
	"starling/cmd"
	_ "starling/cmd/server"
	_ "starling/cmd/worker"
)

// @title Starling
// @version 0.1.0
// @description Smart travel assistant
// @BasePath /sl
func main() {
	cmd.Execute()
}
