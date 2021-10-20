package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	apidb "github.com/microgate-io/microgate-lib-go/v1/db"
	mlog "github.com/microgate-io/microgate/v1/log"
)

type DatabaseServiceImpl struct {
	Connection *pgx.Conn
	apidb.UnimplementedDatabaseServiceServer
	protect          *sync.RWMutex
	openTransactions map[string]pgx.Tx
}

func NewDatabaseServiceImpl(conn *pgx.Conn) *DatabaseServiceImpl {
	return &DatabaseServiceImpl{
		Connection:       conn,
		protect:          new(sync.RWMutex),
		openTransactions: map[string]pgx.Tx{},
	}
}

func (s *DatabaseServiceImpl) Begin(ctx context.Context, req *apidb.BeginRequest) (*apidb.BeginResponse, error) {
	id := uuid.New().String()
	tx, err := s.Connection.Begin(ctx)
	if err != nil {
		mlog.Errorw(ctx, "begin failed", "err", err)
		return nil, err
	}
	s.protect.Lock()
	defer s.protect.Unlock()
	s.openTransactions[id] = tx
	return &apidb.BeginResponse{TransactionId: id}, nil
}
func (s *DatabaseServiceImpl) Commit(ctx context.Context, req *apidb.CommitRequest) (*apidb.CommitResponse, error) {
	s.protect.RLock()
	tx, ok := s.openTransactions[req.TransactionId]
	if !ok {
		mlog.Errorw(ctx, "no such transaction", "tx", req.TransactionId)
		return nil, fmt.Errorf("no such transaction:%s", req.TransactionId)
	}
	s.protect.RUnlock()
	err := tx.Commit(ctx)
	if err != nil {
		mlog.Errorw(ctx, "failed to commit", "tx", req.TransactionId)
		return nil, fmt.Errorf("failed to commit:%s", req.TransactionId)
	}
	s.protect.Lock()
	defer s.protect.Unlock()
	delete(s.openTransactions, req.TransactionId)
	return nil, nil
}
func (s *DatabaseServiceImpl) Rollback(ctx context.Context, req *apidb.RollbackRequest) (*apidb.RollbackResponse, error) {
	s.protect.RLock()
	tx, ok := s.openTransactions[req.TransactionId]
	if !ok {
		mlog.Errorw(ctx, "no such transaction", "tx", req.TransactionId)
		return nil, fmt.Errorf("no such transaction:%s", req.TransactionId)
	}
	s.protect.RUnlock()
	err := tx.Rollback(ctx)
	if err != nil {
		mlog.Errorw(ctx, "failed to rollback", "tx", req.TransactionId, "err", err)
		return nil, fmt.Errorf("failed to rollback:%s", req.TransactionId)
	}
	s.protect.Lock()
	defer s.protect.Unlock()
	delete(s.openTransactions, req.TransactionId)
	return nil, nil
}
func (s *DatabaseServiceImpl) Query(context.Context, *apidb.QueryRequest) (*apidb.QueryResponse, error) {
	return nil, nil
}
func (s *DatabaseServiceImpl) Mutate(ctx context.Context, req *apidb.MutationRequest) (*apidb.MutationResponse, error) {
	mlog.Debugw(ctx, "Mutate", "sql", req.Sql)
	s.protect.RLock()
	tx, ok := s.openTransactions[req.TransactionId]
	if !ok {
		mlog.Errorw(ctx, "no such transaction", "tx", req.TransactionId)
		return nil, fmt.Errorf("no such transaction:%s", req.TransactionId)
	}
	s.protect.RUnlock()
	rows, err := tx.Query(ctx, req.Sql)
	if err != nil {
		mlog.Errorw(ctx, "failed to exec", "tx", req.TransactionId, "sql", req.Sql, "err", err)
		return nil, fmt.Errorf("failed to exec:%s", req.TransactionId)
	}
	return &apidb.MutationResponse{Rows: convertRows(ctx, rows)}, nil
}

func convertRows(ctx context.Context, rows pgx.Rows) (list []*apidb.Row) {
	for rows.Next() {
		desc := rows.FieldDescriptions()
		values, err := rows.Values()
		if err != nil {
			mlog.Errorw(ctx, "cannot access values of rows", "err", err)
			return list
		}
		columns := []*apidb.Column{}
		for i, each := range desc {
			c := &apidb.Column{
				Name:         string(each.Name),
				ValueIsValid: true,
			}
			switch each.DataTypeOID {
			case pgtype.TextOID:
				c.Value = &apidb.Column_StringValue{StringValue: values[i].(string)}
			case pgtype.Int2OID, pgtype.Int4OID:
				c.Value = &apidb.Column_Int32Value{Int32Value: values[i].(int32)}
			case pgtype.Int8OID:
				c.Value = &apidb.Column_Int64Value{Int64Value: values[i].(int64)}
			case pgtype.BoolOID:
				c.Value = &apidb.Column_BoolValue{BoolValue: values[i].(bool)}
			default:
				c.ValueIsValid = false
			}
			columns = append(columns, c)
		}
		list = append(list, &apidb.Row{
			Columns: columns,
		})
	}
	rows.Close()
	return
}
