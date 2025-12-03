// Package main is the entry point of the application.
// It initializes configuration and starts the HTTP server.
package main

import "github.com/mengbin92/example/cmd"

// main is the entry point of the application.
// It delegates to cmd.Execute() to start the server.
func main() {
	cmd.Execute()
}