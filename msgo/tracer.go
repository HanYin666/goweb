package msgo

import (
	tracer2 "github.com/HanYin666/goweb/tracer"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
)

func Tracer(serviceName string, samplerConfig *config.SamplerConfig, reporter *config.ReporterConfig, options ...config.Option) MiddlewareFunc {

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			//接收jaeger的信息 解析上下文
			// 使用 opentracing.GlobalTracer() 获取全局 Tracer
			tracer, closer, spanContext, _ := tracer2.CreateTracerHeader(serviceName, ctx.R.Header, samplerConfig, reporter, options...)
			defer closer.Close()
			//
			// 生成依赖关系，并新建一个 span、
			// 这里很重要，因为生成了  References []SpanReference 依赖关系
			startSpan := tracer.StartSpan(ctx.R.URL.Path, ext.RPCServerOption(spanContext))
			defer startSpan.Finish()
			// 记录 tag
			// 记录请求 Url
			ext.HTTPUrl.Set(startSpan, ctx.R.URL.Path)
			// Http Method
			ext.HTTPMethod.Set(startSpan, ctx.R.Method)
			// 记录组件名称
			ext.Component.Set(startSpan, "Msgo-Http")

			// 在 header 中加上当前进程的上下文信息
			ctx.R = ctx.R.WithContext(opentracing.ContextWithSpan(ctx.R.Context(), startSpan))
			next(ctx)
			// 继续设置 tag
			ext.HTTPStatusCode.Set(startSpan, uint16(ctx.StatusCode))
		}
	}
}
