package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	GaugeVecApiDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiDuration",
		Help: "api耗时单位ms",
	}, []string{"WSorAPI"})
	GaugeVecApiMethod = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiCount",
		Help: "各种网络请求次数",
	}, []string{"method"})
	GaugeVecApiError = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiErrorCount",
		Help: "请求api错误的次数type: api/ws",
	}, []string{"type"})
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(GaugeVecApiMethod, GaugeVecApiDuration, GaugeVecApiError)
}

func MwPrometheusHttp(c *gin.Context) {
	start := time.Now()
	method := c.Request.Method
	GaugeVecApiMethod.WithLabelValues(method).Inc()

	c.Next()
	// after request
	end := time.Now()
	d := end.Sub(start) / time.Millisecond
	GaugeVecApiDuration.WithLabelValues(method).Set(float64(d))
}

func JsonError(c *gin.Context, msg interface{}) {
	GaugeVecApiError.WithLabelValues("API").Inc()
	var ms string
	switch v := msg.(type) {
	case string:
		ms = v
	case error:
		ms = v.Error()
	default:
		ms = ""
	}
	c.JSON(200, gin.H{"ok": false, "msg": ms})
}

func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api := router.Group("/api").Use(MwPrometheusHttp)
	api.GET("/test", func(c *gin.Context) {
		time.Sleep(time.Second)
		c.JSON(200, "success")
	})
	if err := router.Run(":3001"); err != nil {
		panic(err)
	}
}
