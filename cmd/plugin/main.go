package main

import (
	"github.com/planx-lab/planx-plugin-sink-stdout/internal/plugin"
	"github.com/planx-lab/planx-sdk-go/sdk"
)

func main() {
	sdk.ServeSink(plugin.New)
}
