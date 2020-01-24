package main

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MetricServerReqDur = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: serverNamespace,
			Subsystem: "requests",
			Name:      "duration_ms",
			Help:      "http server requests duration(ms).",
			Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
		},
		[]string{"path", "caller"},
	)
	MetricServerReqCodeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: serverNamespace,
			Subsystem: "requests",
			Name:      "code_total",
			Help:      "http server requests error count.",
		},
		[]string{"path", "caller", "code"},
	)
)

func init() {
	prometheus.MustRegister(MetricServerReqDur)
	prometheus.MustRegister(MetricServerReqCodeTotal)
}

const (
	serverNamespace = "http_server"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		MetricServerReqDur.WithLabelValues("/path", "user").Observe(float64(time.Second))
		MetricServerReqCodeTotal.WithLabelValues("/path", "user", strconv.FormatInt(200, 10)).Inc()
		writer.Write([]byte("hello world"))
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
