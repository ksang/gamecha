package seeker

import (
	"context"
	"net/http"
)

func httpDo(ctx context.Context, req *http.Request, client *http.Client, callback func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	c := make(chan error, 1)
	if client == nil {
		client = http.DefaultClient
	}
	req = req.WithContext(ctx)
	go func() { c <- callback(client.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c
		return ctx.Err()
	case err := <-c:
		return err
	}
}
