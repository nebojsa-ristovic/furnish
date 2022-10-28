package util

import (
	"bufio"
	"os"
	"strings"

	"github.com/fatih/color"
)

var reader = bufio.NewReader(os.Stdin)

func ReadConfirmation(print, confirm string) bool {
	color.Yellow(print)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(input) == strings.ToLower(confirm)
}
