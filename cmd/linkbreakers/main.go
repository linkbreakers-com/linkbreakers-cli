package main

import (
	"os"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
