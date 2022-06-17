package pkcesupport_test

import (
	"github.com/initialcapacity/oauth-starter/pkg/pkcesupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodeChallenge(t *testing.T) {
	verifier := pkcesupport.CodeVerifier()
	challenge := pkcesupport.CodeChallenge(verifier)
	assert.Equal(t, challenge, pkcesupport.CodeChallenge(verifier))
}
