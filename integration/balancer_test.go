package integration

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

type IntegrationSuite struct {
	suite.Suite
}

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

type RespBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *IntegrationSuite) TestBalancer() {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		s.T().Skip("Integration test is not enabled")
	}

	// test server1
	server1, err := client.Get(fmt.Sprintf("%s/check", baseAddress))
	assert.NoError(s.T(), err)
	server1Header := server1.Header.Get("lb-from")
	assert.Equal(s.T(), "server1:8080", server1Header)

	// test server2
	server2, err := client.Get(fmt.Sprintf("%s/check4", baseAddress))
	assert.NoError(s.T(), err)
	server2Header := server2.Header.Get("lb-from")
	assert.Equal(s.T(), "server2:8080", server2Header)

	// test server3
	server3, err := client.Get(fmt.Sprintf("%s/check2", baseAddress))
	assert.NoError(s.T(), err)
	server3Header := server3.Header.Get("lb-from")
	assert.Equal(s.T(), "server3:8080", server3Header)

	// test repeated request
	server1Repeat, err := client.Get(fmt.Sprintf("%s/check", baseAddress))
	assert.NoError(s.T(), err)
	server1RepeatHeader := server1Repeat.Header.Get("lb-from")
	assert.Equal(s.T(), server1Header, server1RepeatHeader)
}

func (s *IntegrationSuite) BenchmarkBalancer(b *testing.B) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
	  s.T().Skip("Integration test is not enabled")
	}
 
	for i := 0; i < b.N; i++ {
	  _, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	  assert.NoError(s.T(), err)
	}
 }