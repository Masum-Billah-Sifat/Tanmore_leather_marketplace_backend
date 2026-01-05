package sqlnull

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// STRING → NullString
func String(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// INT → NullInt64
func Int64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

// BOOL → NullBool
func Bool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// TIME → NullTime
func Time(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}

// UUID → NullUUID
func UUID(id uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: id, Valid: true}
}

// new additions here

// STRING PTR → NullString
func StringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	if *s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// INT64 PTR → NullInt64
func Int64Ptr(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

// BOOL PTR → NullBool
func BoolPtr(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

// TIME PTR → NullTime
func TimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// UUID PTR → NullUUID
func UUIDPtr(id *uuid.UUID) uuid.NullUUID {
	if id == nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{UUID: *id, Valid: true}
}

// INT32 PTR → NullInt32
func Int32Ptr(i *int64) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: int32(*i), Valid: true}
}

func Int32(i int64) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(i), Valid: true}
}

// pkg/sqlnull/sqlnull.go
func Int32From32(i int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}

// In pkg/sqlnull/sqlnull.go
func DecimalPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func Float64Ptr(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: *f,
		Valid:   true,
	}
}

func StringOrEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func Int64OrZero(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

func Int32OrZero(ni sql.NullInt32) int32 {
	if ni.Valid {
		return ni.Int32
	}
	return 0
}

// ToStringPtr converts sql.NullString to *string
func ToStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// ToInt64Ptr converts sql.NullInt64 to *int64
func ToInt64Ptr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

// ToInt32Ptr converts sql.NullInt32 to *int32
func ToInt32Ptr(n sql.NullInt32) *int32 {
	if n.Valid {
		return &n.Int32
	}
	return nil
}
