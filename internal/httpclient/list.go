package httpclient

import "context"

// listEnvelope 는 OpenDART list 응답 공통 형태 (status/message + list).
type listEnvelope[T any] struct {
	Envelope
	List []T `json:"list"`
}

// GetList 는 JSON list 응답을 디코드해 list 만 반환한다.
// status 검사(013→ErrNoData, 그 외→*APIError)는 GetJSON 이 수행한다.
func GetList[T any](ctx context.Context, c *Client, path string, params map[string]string) ([]T, error) {
	var resp listEnvelope[T]
	if err := c.GetJSON(ctx, path, params, &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}
