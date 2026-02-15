package cloudflare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleEntryToMap(t *testing.T) {
	role := &cloudflareRoleEntry{PolicyDocument: "policy"}

	respData, err := roleEntryToMap(role)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"policy_document": "policy"}, respData)
}
