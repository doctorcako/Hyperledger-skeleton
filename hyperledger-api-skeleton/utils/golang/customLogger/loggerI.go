package customLogger

import (
	"context"
	"time"
)

type Log interface {
	Error(msg ...any)
	ErrorCtx(ctx context.Context, msg ...any)
	Warning(msg ...any)
	WarningCtx(ctx context.Context, msg ...any)
	Info(msg ...any)
	InfoCtx(ctx context.Context, msg ...any)
	Debug(msg ...any)
	DebugCtx(ctx context.Context, msg ...any)
	DebugTime(msg ...any)
	DebugTimeCtx(ctx context.Context, msg ...any)
	GetTimeNow() time.Time
	CalculateDifference(initial time.Time) time.Duration
	WarningReturnCtx(ctx context.Context, msg ...any) context.Context
	InfoReturnCtx(ctx context.Context, msg ...any) context.Context
	DebugReturnCtx(ctx context.Context, msg ...any) context.Context
}
