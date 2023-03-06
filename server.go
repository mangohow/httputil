package httputil

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mangohow/httputil/logger"
	"github.com/mangohow/httputil/proc"
	"net/http"
	"sync"
)

type HttpHandler func(ctx *gin.Context) *Response

// Decorate 装饰器
func Decorate(handlers []HttpHandler) []gin.HandlerFunc {
	hds := make([]gin.HandlerFunc, len(handlers))
	for i := 0; i < len(handlers); i++ {
		h := handlers[i]
		hds[i] = func(ctx *gin.Context) {
			r := h(ctx)
			if r != nil {
				ctx.JSON(r.HttpStatus, &r.R)
			}

			putResponse(r)
		}
	}

	return hds
}

type BeforeCloseHandler func()

type GinServer struct {
	*gin.Engine
	server    *http.Server
	startOnce sync.Once

	beforeCloseHandlers []BeforeCloseHandler
	mutex               sync.Mutex

	ctx context.Context

	log logger.Logger
}

func NewServer() *GinServer {
	return &GinServer{
		Engine: gin.Default(),
	}
}

func NewWithEngine(r *gin.Engine) *GinServer {
	return &GinServer{
		Engine: r,
	}
}

func (s *GinServer) Get(relativePath string, handlers ...HttpHandler) {
	s.GET(relativePath, Decorate(handlers)...)
}

func (s *GinServer) Post(relativePath string, handlers ...HttpHandler) {
	s.POST(relativePath, Decorate(handlers)...)
}

func (s *GinServer) Put(relativePath string, handlers ...HttpHandler) {
	s.PUT(relativePath, Decorate(handlers)...)
}

func (s *GinServer) Delete(relativePath string, handlers ...HttpHandler) {
	s.DELETE(relativePath, Decorate(handlers)...)
}

func (s *GinServer) Head(relativePath string, handlers ...HttpHandler) {
	s.HEAD(relativePath, Decorate(handlers)...)
}

func (s *GinServer) HttpServer() *http.Server {
	return s.server
}

func (s *GinServer) AddBeforeServerCloseHandlers(handlers ...BeforeCloseHandler) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.beforeCloseHandlers = append(s.beforeCloseHandlers, handlers...)
}

func (s *GinServer) Run(addr string, logger logger.Logger) context.Context {
	var ctx context.Context

	s.startOnce.Do(func() {
		if addr == "" {
			addr = ":8080"
		}
		server := &http.Server{
			Addr:    addr,
			Handler: s.Engine,
		}
		s.server = server
		s.log = logger

		ctx = proc.SetupSignalHandler(logger)
		s.ctx = ctx

		go func() {
			<-ctx.Done()
			if err := server.Shutdown(ctx); err != nil {
				logger.Errorf("server shutdown :%v", err)
			}
		}()

		go func() {
			logger.Infof("http server listen at %s", addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(err)
			}

			logger.Info("http server stopping...")
		}()
	})

	return ctx
}

func (s *GinServer) Wait() {
	<-s.ctx.Done()
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := 0; i < len(s.beforeCloseHandlers); i++ {
		h := s.beforeCloseHandlers[i]
		h()
	}
}