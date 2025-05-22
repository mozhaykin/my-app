//go:build integration

package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/app"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclient"
)

// Run test: make integration-test

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	profile *httpclient.Client
}

func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()

	s.profile = httpclient.New("localhost:8080")

	go func() {
		err := app.Run()
		panic(err)
	}()

	time.Sleep(time.Second)
}

func (s *Suite) TearDownSuite() {}

func (s *Suite) SetupTest() {}

func (s *Suite) TearDownTest() {}
