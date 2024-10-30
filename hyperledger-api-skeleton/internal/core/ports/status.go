package ports

import (
	"context"
	"net/http"
)

type StatusService interface {
	CheckLiveness(ctx context.Context, w http.ResponseWriter, r *http.Request)
	CheckReadiness(ctx context.Context, w http.ResponseWriter, r *http.Request)
}
