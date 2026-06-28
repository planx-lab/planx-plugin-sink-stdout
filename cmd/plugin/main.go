package main

import (
	"github.com/planx-lab/planx-plugin-sink-stdout/internal/plugin"
	"github.com/planx-lab/planx-sdk-go/sdk"
)

func main() {
	sdk.Serve(sdk.Plugin{
		ID:          "sink-stdout",
		Version:     "1.0.0",
		DisplayName: "Stdout Sink",
		Description: "Acceptance Test Sink (Stdout)",
		Summary:     "Writes each batch to stdout (test/demonstration only).",
		Components: []sdk.ComponentSpec{{
			ID:          "sink",
			Kind:        sdk.KindSink,
			DisplayName: "Stdout Sink",
			Sink:        plugin.New,
		}},
	})
}
