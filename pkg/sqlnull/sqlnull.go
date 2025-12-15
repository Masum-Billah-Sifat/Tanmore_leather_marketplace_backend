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
