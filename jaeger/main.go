package main

import (
	"context"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	tracer, closer, err := jaegercfg.Configuration{
		ServiceName: "golang-test",
		Headers: &jaeger.HeadersConfig{
			TraceContextHeaderName: "trace-id",
		},
		//Sampler: &jaegercfg.SamplerConfig{
		//	Type:  jaeger.SamplerTypeConst,
		//	Param: 1,
		//},
		// 采样率
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "probabilistic",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831", // 替换host
		},
	}.NewTracer()
	//closer, err := cfg.InitGlobalTracer(
	//	"serviceName",
	//)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	var ctx = context.TODO()
	span1, ctx := opentracing.StartSpanFromContext(ctx, "span_1")
	time.Sleep(time.Second / 2)
	span11, _ := opentracing.StartSpanFromContext(ctx, "span_1-1")
	time.Sleep(time.Second / 2)
	span11.Finish()
	span1.Finish()
}
