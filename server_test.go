package httputil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mangohow/httputil/logger"
	"testing"
)

func TestHttpServer(t *testing.T) {
	log := logger.FakeLogger{}

	server := NewServer()

	const (
		Success uint32 = iota
		Failed
	)

	messages := map[uint32]string{
		Success: "request success",
		Failed:  "request failed",
	}

	getMessage := func(code uint32) string {
		return messages[code]
	}

	SetMessager(MessagerFunc(getMessage))
	server.Get("/", func(ctx *gin.Context) *Response {
		return NewResponseOK(0, "hello,world")
	})

	server.AddBeforeServerCloseHandlers(func() {
		fmt.Println("test")
	})

	server.Run(":8080", log)

	server.Wait()
}
