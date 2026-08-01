package plugin

import (
	"context"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

// DB connectors (postgres, sqlserver) emit a gob-encoded DBBatch under the
// shared wire name "planx.io/dbbatch.DBBatch". The SDK's batch codec decodes
// into a concrete registered type, so a sink must register a structurally
// identical type under the same wire name to receive DB rows — otherwise gob
// fails with "name not registered for interface".
//
// These types mirror planx-plugin-{postgres,sqlserver}/internal/dbbatch. They
// are NOT imported (these are separate repos / modules); the gob wire name is
// the contract. See ADR-005 (batch identity) and the dbbatch package docs.
//
// NOTE: this is the contained workaround. The durable fix is to make the SDK's
// batch contract truly bytes-opaque so sinks need not know source formats.
type DBRow struct {
	Types []byte
	Vals  []string
}

type DBBatch struct {
	Columns []string
	Rows    []DBRow
}

func init() {
	gob.RegisterName("planx.io/dbbatch.DBBatch", DBBatch{})
	gob.RegisterName("planx.io/dbbatch.DBRow", DBRow{})
}

type Sink struct{}

func New() sdk.SinkSPI {
	return &Sink{}
}

func (s *Sink) Init(ctx context.Context, cfg []byte) error {
	return nil
}

func (s *Sink) WriteBatch(batch sdk.Batch) error {
	// Pretty-print DB rows when the upstream source is a DB connector;
	// otherwise fall back to the generic dump.
	if db, ok := batch.(DBBatch); ok && len(db.Columns) > 0 {
		for _, row := range db.Rows {
			parts := make([]string, len(db.Columns))
			for i, col := range db.Columns {
				val := ""
				if i < len(row.Vals) {
					val = row.Vals[i]
				}
				parts[i] = fmt.Sprintf("%s=%s", col, val)
			}
			fmt.Printf("[SINK] row: %s\n", strings.Join(parts, " "))
		}
		return nil
	}
	fmt.Printf("[SINK] Received Batch: %v\n", batch)
	return nil
}

func (s *Sink) Close() error {
	return nil
}
