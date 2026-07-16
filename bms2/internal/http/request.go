package http

import (
	"fmt"
	"net/http"
	"strconv"
)

func ParseIntPathParam(r *http.Request, name string) (int64, error) {
	value := r.PathValue(name)

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", name, err)
	}

	return id, nil
}
