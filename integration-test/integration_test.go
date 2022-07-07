package integrationtest_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/suite"
)

const (
	attempts      = 20
	host          = "devops_server:8080"
	serverAddress = "http://" + host
	healthPath    = serverAddress + "/healthz"
	valuePath     = "/value/"
	updatePath    = "/update/"
	key           = "a"
)

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}
		log.Printf("Integration tests: url %s is not available, attempt: %d", healthPath, attempts)
		time.Sleep(time.Second)
		attempts--
	}
	return err
}

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}
	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func TestIteration1(t *testing.T) {
	// Implement_9
	suite.Run(t, new(CryptoSuite))
}
