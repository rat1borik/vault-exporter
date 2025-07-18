package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type countingBatch struct {
	batch  *pgx.Batch
	queued int
	tracer *PgTracer
}

func (cb *countingBatch) Queue(sql string, args ...any) {
	cb.batch.Queue(sql, args...)
	cb.queued++
}

func NewCountingBatch() *countingBatch {
	return &countingBatch{
		batch:  &pgx.Batch{},
		tracer: &PgTracer{},
	}
}

func (cb *countingBatch) Exec(ctx context.Context, db *pgx.Conn) error {
	cb.tracer.TraceQueryStart(ctx, nil, pgx.TraceQueryStartData{})
	br := db.SendBatch(ctx, cb.batch)
	defer br.Close()

	for i := 0; i < cb.queued; i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}

func (cb *countingBatch) ExecTx(ctx context.Context, tx pgx.Tx) error {
	br := tx.SendBatch(ctx, cb.batch)
	defer br.Close()

	for i := 0; i < cb.queued; i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}
