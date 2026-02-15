package cloudflare

import (
	"net/http"
	"reflect"
	"testing"
	"time"
	"unsafe"

	cf "github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/assert"
)

func TestCreateClientConfiguresTimeoutAndRetryPolicy(t *testing.T) {
	client, err := createClient("test-token")
	if err != nil {
		t.Fatalf("failed to create client: %s", err)
	}

	httpClient := mustGetUnexportedField[*http.Client](t, client, "httpClient")
	retryPolicy := mustGetUnexportedField[cf.RetryPolicy](t, client, "retryPolicy")

	if httpClient == nil {
		t.Fatal("expected http client to be configured")
	}

	assert.Equal(t, cloudflareClientTimeout, httpClient.Timeout)
	assert.Equal(t, cloudflareMaxRetries, retryPolicy.MaxRetries)
	assert.Equal(t, time.Duration(cloudflareMinRetryDelaySeconds)*time.Second, retryPolicy.MinRetryDelay)
	assert.Equal(t, time.Duration(cloudflareMaxRetryDelaySeconds)*time.Second, retryPolicy.MaxRetryDelay)
}

func mustGetUnexportedField[T any](t *testing.T, value interface{}, name string) T {
	t.Helper()

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		t.Fatalf("expected pointer value, got %v", v.Kind())
	}

	elem := v.Elem()
	field := elem.FieldByName(name)
	if !field.IsValid() {
		t.Fatalf("missing field %s", name)
	}

	field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	return field.Interface().(T)
}
