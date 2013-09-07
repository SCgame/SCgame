package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandDispatching(t *testing.T) {
	g := BuildGame()
	assert.Equal(t, g.runCommand("UNKNOWN_COMMAND 0 1 2"), "ERROR 1 INVALID COMMAND")
	g.registerCommand("UNKNOWN_COMMAND", 3, func(a []string) string {
		return fmt.Sprintf("Function called with %s %s and %s", a[0], a[1], a[2])
	})
	assert.Equal(t, g.runCommand("UNKNOWN_COMMAND 0 1 2"), "Function called with 0 1 and 2")
	assert.Equal(t, g.runCommand("UNKNOWN_COMMAND 1 2"), "ERROR 2 INVALID NUMBER OF ARGUMENTS")
	assert.Equal(t, g.runCommand("UNKNOWN_COMMAND"), "ERROR 2 INVALID NUMBER OF ARGUMENTS")
}
