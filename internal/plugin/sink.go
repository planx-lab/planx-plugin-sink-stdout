package plugin

import (
	"context"
	"fmt"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

type Sink struct{}

func New() sdk.SinkSPI {
	return &Sink{}
}

func (s *Sink) Init(ctx context.Context, cfg []byte) error {
	return nil
}

func (s *Sink) WriteBatch(batch sdk.Batch) error {
	fmt.Printf("[SINK] Received Batch: %v\n", batch)
	return nil
}

func (s *Sink) Close() error {
	return nil
}
