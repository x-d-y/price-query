package httpRequest

import (
	"context"
	"net/http"
	"time"
)

type Requester struct {
	client  *http.Client
	options interface{} // 后期要改成日志 + 可能还有钩子
	timeout time.Duration
}

func NewClient() *Requester {
	httpClient := &http.Client{}
	return &Requester{
		client:  httpClient,
		options: nil,
	}
}

func (r *Requester) ContextGet(ctx context.Context, url string, header http.Header) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if header != nil {
		request.Header = header
	}

	response, rErr := r.client.Do(request)
	if rErr != nil {
		return nil, rErr
	}
	return response, nil
}
