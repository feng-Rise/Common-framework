package RpcProxyMode

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"testing"
	"time"

	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func TestServer_Start(t *testing.T) {
	s := NewServer()
	s.Register(&UserServiceServer{})
	s.Start(":8081")
}

type UserServiceServer struct {
}

func (u *UserServiceServer) Name() string {
	return "user-service"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		Name: "feng",
	}, nil
}

func TestABC(t *testing.T) {

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
		ServiceName: "fish_test",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	span := opentracing.StartSpan("go-grpc-web")
	traceID := span.Context().(jaeger.SpanContext).TraceID()
	fmt.Println(traceID)
	time.Sleep(time.Second)
	defer span.Finish()

}
