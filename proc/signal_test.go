package proc

import (
	"github.com/mangohow/httputil/logger"
	"testing"
)

func TestSetupSignalHandler(t *testing.T) {
	log := logger.FakeLogger{}

	ctx := SetupSignalHandler(log)

	<-ctx.Done()
}
