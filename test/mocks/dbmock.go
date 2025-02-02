package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// MockDB implement DBTX interface
type MockDB struct {
	ExecFunc     func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryFunc    func(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRowFunc func(context.Context, string, ...interface{}) pgx.Row
}

func (m *MockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if m.ExecFunc == nil {
		return pgconn.CommandTag{}, fmt.Errorf("ExecFunc not implemented")
	}
	return m.ExecFunc(ctx, sql, args...)
}

func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if m.QueryFunc == nil {
		return nil, fmt.Errorf("QueryFunc not implemented")
	}
	return m.QueryFunc(ctx, sql, args...)
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if m.QueryRowFunc == nil {
		return &ErrorRow{Err: fmt.Errorf("QueryRowFunc not implemented")}
	}
	return m.QueryRowFunc(ctx, sql, args...)
}

// MockRows implement pgx.Rows
type MockRows struct {
	data    [][]interface{}
	idx     int
	scanErr error
}

func NewMockRows(data [][]interface{}) *MockRows {
	return &MockRows{data: data, idx: -1}
}

func NewErrorMockRows(err error) *MockRows {
	return &MockRows{scanErr: err}
}

func (m *MockRows) Next() bool {
	m.idx++
	return m.idx < len(m.data)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	if m.scanErr != nil {
		return m.scanErr
	}

	if m.idx < 0 || m.idx >= len(m.data) {
		return pgx.ErrNoRows
	}
	row := m.data[m.idx]
	if len(row) != len(dest) {
		return fmt.Errorf("mismatched columns: expected %d, got %d", len(row), len(dest))
	}

	for i, val := range row {
		if err := assignValue(dest[i], val); err != nil {
			return fmt.Errorf("column %d: %w", i, err)
		}
	}
	return nil
}

func (m *MockRows) Close()     {}
func (m *MockRows) Err() error { return nil }

type MockRow struct {
	Data []interface{}
}

func (r *MockRow) Scan(dest ...interface{}) error {
	if len(r.Data) != len(dest) {
		return fmt.Errorf("mismatched columns: expected %d, got %d", len(r.Data), len(dest))
	}
	for i, val := range r.Data {
		if err := assignValue(dest[i], val); err != nil {
			return fmt.Errorf("column %d: %w", i, err)
		}
	}
	return nil
}

// ErrorRow for emulation scan func
type ErrorRow struct {
	Err error
}

// Scan implements pgx.Row.
func (r ErrorRow) Scan(dest ...any) error {
	return r.Err
}

func assignValue(dest interface{}, src interface{}) error {
	switch d := dest.(type) {
	case *pgtype.Time:
		switch v := src.(type) {
		case pgtype.Time:
			*d = v
		case time.Time:
			*d = pgtype.Time{
				Microseconds: v.UnixMicro(),
				Valid:        true,
			}
		case int64:
			*d = pgtype.Time{
				Microseconds: v,
				Valid:        true,
			}
		default:
			return fmt.Errorf("unsupported type for Time: %T", src)
		}

	case *int64:
		switch v := src.(type) {
		case int64:
			*d = v
		case int:
			*d = int64(v)
		default:
			return fmt.Errorf("unsupported type for int64: %T", src)
		}

	case *int16:
		switch v := src.(type) {
		case int16:
			*d = v
		case int:
			*d = int16(v)
		default:
			return fmt.Errorf("unsupported type for int16: %T", src)
		}

	case *string:
		switch v := src.(type) {
		case string:
			*d = v
		default:
			return fmt.Errorf("unsupported type for string: %T", src)
		}

	default:
		return fmt.Errorf("unsupported destination type: %T", dest)
	}
	return nil
}
