package contextUtils

import (
	"context"
	"github.com/oklog/ulid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"math/rand"
	"net/http"
	"time"
)

type ContextUtils interface {
	NewCtx() context.Context
	NewCtxFromRequest(r *http.Request) context.Context
	GetCorrelationId(ctx context.Context) string
	SetCorrelationId(ctx context.Context, identifier string) context.Context
	GetAuthFromCtx(ctx context.Context) string
	GenerateHeaderFromCtx(ctx context.Context, r *http.Request) *http.Request
	GenerateTraceId(ctx context.Context) string
	InjectTraceId(ctx context.Context, uberTokenId string) context.Context
}

const (
	correlationId string = "correlationId"
	authorization string = "authorization"

	HeaderCorrelationId string = "CorrelationId"
	HeaderAuthorization string = "Authorization"
)

type contextUtils struct{}

func NewContextUtils() ContextUtils {
	return &contextUtils{}
}

// NewCtx - create new contexts
func (c *contextUtils) NewCtx() context.Context {
	ctx := context.TODO()
	ctxId := generateUID()
	ctx = context.WithValue(ctx, correlationId, ctxId)

	return ctx
}

// NewCtxFromRequest - create new context from request information
func (c *contextUtils) NewCtxFromRequest(r *http.Request) context.Context {
	ctx := context.TODO()
	ctxId := generateUID()
	if r.Header.Get(HeaderCorrelationId) != "" {
		ctxId = r.Header.Get(HeaderCorrelationId)
	}
	ctx = context.WithValue(ctx, correlationId, ctxId)

	token := ""
	if r.Header.Get(HeaderAuthorization) != "" {
		token = r.Header.Get(HeaderAuthorization)
	}
	ctx = context.WithValue(ctx, authorization, token)

	var span opentracing.Span
	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span = tracer.StartSpan(r.Method+" "+r.RequestURI, ext.RPCServerOption(spanCtx))
	span.SetTag(correlationId, r.Header.Get(HeaderCorrelationId))

	ctx = opentracing.ContextWithSpan(ctx, span)
	span.Finish()

	return ctx
}

// GetCorrelationId - obtain correlation identifier from context
func (c *contextUtils) GetCorrelationId(ctx context.Context) string {
	var value = ""
	if ctx != nil {
		if ctx.Value(correlationId) != nil {
			valueId, _ := ctx.Value(correlationId).(string)
			value = valueId
		}
	}
	return value
}

// SetCorrelationId - set correlation identifier
func (c *contextUtils) SetCorrelationId(ctx context.Context, identifier string) context.Context {
	return context.WithValue(ctx, correlationId, identifier)
}

// GetAuthFromCtx - obtain auth token from context
func (c *contextUtils) GetAuthFromCtx(ctx context.Context) string {
	var value = ""
	if ctx != nil {
		if ctx.Value(authorization) != nil {
			valueId, _ := ctx.Value(authorization).(string)
			value = valueId
		}
	}
	return value
}

// GenerateHeaderFromCtx - create request header from context information
func (c *contextUtils) GenerateHeaderFromCtx(ctx context.Context, r *http.Request) *http.Request {
	if ctx != nil {
		if ContextId := ctx.Value(correlationId); ContextId != nil {
			ContextIdS, _ := ContextId.(string)
			r.Header.Set(HeaderCorrelationId, ContextIdS)
		}
		if token := ctx.Value(authorization); token != nil {
			tokenS, _ := token.(string)
			r.Header.Set(HeaderAuthorization, tokenS)
		}

		//Uber tracer ID
		span, _ := opentracing.StartSpanFromContext(ctx, r.Method)
		defer span.Finish()
		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, r.URL.Path)
		ext.HTTPMethod.Set(span, r.Method)
		_ = span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header),
		)
	}
	return r
}

// GenerateTraceId - generate new tracer identifier
func (c *contextUtils) GenerateTraceId(ctx context.Context) string {
	req, _ := http.NewRequest(http.MethodPost, "", nil)
	span, _ := opentracing.StartSpanFromContext(ctx, "Event")
	defer span.Finish()
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, "sns -> sqs")
	ext.HTTPMethod.Set(span, http.MethodPost)
	_ = span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	if req.Header.Get("Uber-Trace-Id") != "" {
		return req.Header.Get("Uber-Trace-Id")
	}
	return ""
}

// InjectTraceId -put tracer identifier into context
func (c *contextUtils) InjectTraceId(ctx context.Context, uberTokenId string) context.Context {
	req, _ := http.NewRequest(http.MethodPost, "", nil)
	req.Header.Set("Uber-Trace-Id", uberTokenId)

	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	span := tracer.StartSpan(""+
		"", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx
}

func generateUID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
