package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	suite.Run(t, new(BalancerSuite))
}

type BalancerSuite struct {
	suite.Suite
}

func (s *BalancerSuite) TestBalancer() {

	healthChecker := &HealthChecker{}
	healthChecker.healthyServers = []string{"4", "5", "6"}

	balancer := &Balancer{}
	balancer.healthChecker = healthChecker

	index1 := balancer.GetServerIndex("/check")
	index1secondTime := balancer.GetServerIndex("/check")
	index2 := balancer.GetServerIndex("/check2")
	index3 := balancer.GetServerIndex("/check4")

	assert.Equal(s.T(), index1, 0)
	assert.Equal(s.T(), index1, index1secondTime)
	assert.Equal(s.T(), index3, 1)
	assert.Equal(s.T(), index2, 2)
}

func (s *BalancerSuite) TestHealthChecker() {
	healthChecker := &HealthChecker{}
	healthChecker.health = func(s string) bool {
		if s == "1" {
			return false
		} else {
			return true
		}
	}

	healthChecker.serversPool = []string{"1", "2", "3"}
	healthChecker.healthyServers = []string{"4", "5", "6"}
	healthChecker.checkInterval = 1 * time.Second

	healthChecker.StartHealthCheck()

	time.Sleep(2 * time.Second)

	assert.Equal(s.T(), healthChecker.healthyServers[0], "2")
	assert.Equal(s.T(), healthChecker.healthyServers[1], "3")
	assert.Equal(s.T(), len(healthChecker.healthyServers), 2)
}
