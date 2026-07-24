package requestid

import "context"

type contextKey string

const key contextKey = "request_id"

func IntoContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key, id)
}

func FromContext(ctx context.Context) string {
	id, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}

	return id
}
