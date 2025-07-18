// Package infrastructure содержит инфраструктурные инструменты
package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// TODO: сделать для batch и в какую-то общую помойку скидывать
type PgTracer struct {
	IsProd bool
}

type TraceData struct {
	Command   string
	Args      []any
	StartTime time.Time
}

type ctxTraceKey string

const keyStartTime ctxTraceKey = "traceData"

func (t PgTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	start := time.Now()
	return context.WithValue(ctx, keyStartTime, TraceData{
		Command:   data.SQL,
		Args:      data.Args,
		StartTime: start,
	})
}

func (t PgTracer) TraceQueryEnd(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	tData, ok := ctx.Value(keyStartTime).(TraceData)
	if !ok {
		return
	}
	duration := time.Since(tData.StartTime)
	if !t.IsProd {
		fmt.Printf("[PGX TRACE] SQL: %s | Args: %v | Duration: %v | Err: %v\n",
			tData.Command, tData.Args, duration, data.Err)
	}
}
