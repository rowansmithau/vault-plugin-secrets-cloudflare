package cloudflare

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/vault/sdk/logical"
)

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func WithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{Header: make(http.Header), rt: rt}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}

const (
	cloudflareClientTimeout        = 60 * time.Second
	cloudflareMaxRetries           = 5
	cloudflareMinRetryDelaySeconds = 1
	cloudflareMaxRetryDelaySeconds = 60
)

func createClient(token string) (*cloudflare.API, error) {
	client := &http.Client{
		Timeout: cloudflareClientTimeout,
	}
	return cloudflare.NewWithAPIToken(
		token,
		cloudflare.HTTPClient(client),
		cloudflare.UsingRetryPolicy(cloudflareMaxRetries, cloudflareMinRetryDelaySeconds, cloudflareMaxRetryDelaySeconds),
	)
}

func (b *backend) client(ctx context.Context, s logical.Storage) (*cloudflare.API, error) {
	conf, err := b.readConfigToken(ctx, s)
	if err != nil {
		return nil, err
	}
	return createClient(conf.Token)
}
