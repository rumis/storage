package processor

import (
	"context"
	"time"

	"github.com/rumis/storage/v2/locker"
	"github.com/rumis/storage/v2/meta"
	"github.com/rumis/storage/v2/srepo"
)

// RepoReadProcessor DB数据读取
type RepoOneReadProcessor struct {
	locker locker.Locker
	reader meta.RepoReader
	writer meta.RepoInserter
	next   ReadProcessor
}

// Read
func (r *RepoOneReadProcessor) Read(ctx context.Context, in interface{}, out interface{}, exp time.Duration) error {
	zero, ok := out.(meta.Zero)
	if !ok {
		return meta.EI_ZeroNotImplement
	}
	var fn meta.QueryExprHandler
	switch query := in.(type) {
	case meta.Query:
		fn = query.Query(ctx)
	case meta.QueryExprHandler:
		fn = query
	default:
		return meta.EI_QueryNotImplement
	}
	err := r.reader(ctx, out, fn)
	if err != nil && r.next == nil {
		return err
	}
	if zero.Zero() && r.next == nil {
		return meta.EI_Zero
	}
	err = r.next.Read(ctx, in, out, exp)
	if err != nil {
		return err
	}
	return nil
}

// Write
func (r *RepoOneReadProcessor) Write(ctx context.Context, in interface{}, exp time.Duration) error {
	_, err := r.writer(ctx, in)
	return err
}

// NewRepoOneReadProcessor
func NewRepoOneReadProcessor(next ReadProcessor, tableName string, columns []string) ReadProcessor {
	return &RepoOneReadProcessor{
		next:   next,
		reader: srepo.NewSealMysqlOneReader(srepo.WithDB(srepo.SealR()), srepo.WithColumns(columns), srepo.WithName(tableName)),
		writer: srepo.NewSealMysqlOneInserter(srepo.WithDB(srepo.SealW()), srepo.WithColumns(columns), srepo.WithName(tableName)),
		locker: locker.NewLocker(),
	}
}
