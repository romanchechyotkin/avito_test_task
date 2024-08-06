//go:build integration

package service

import (
	"testing"

	"github.com/romanchechyotkin/avito_test_task/pkg/utest"
)

var log, cfg, pg, prepareErr = utest.Prepare()

func TestMain(t *testing.M) {
	t.Run()
}
