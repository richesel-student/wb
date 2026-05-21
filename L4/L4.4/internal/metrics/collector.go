package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {

	memAllocBytes          prometheus.Gauge
	memTotalAllocBytes     prometheus.Counter
	memSysBytes           prometheus.Gauge
	gcCyclesTotal         prometheus.Counter
	gcLastTimeSeconds     prometheus.Gauge
	gcPauseTotalNs        prometheus.Counter
	goroutinesCount       prometheus.Gauge
	heapObjectsCount      prometheus.Gauge


	ticker *time.Ticker
}


func NewCollector() *Collector {
	c := &Collector{
		memAllocBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gc_monitor_mem_alloc_bytes",
				Help: "Memory in use (Alloc)",
			},
		),
		memTotalAllocBytes: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "gc_monitor_mem_total_alloc_bytes_total",
				Help: "Total bytes allocated (TotalAlloc)",
			},
		),
		memSysBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gc_monitor_mem_sys_bytes",
				Help: "Memory requested from OS (Sys)",
			},
		),
		gcCyclesTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "gc_monitor_gc_cycles_total",
				Help: "Number of garbage collections (NumGC)",
			},
		),
		gcLastTimeSeconds: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gc_monitor_gc_last_time_seconds",
				Help: "Time of the last GC in Unix timestamp (LastGC)",
			},
		),
		gcPauseTotalNs: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "gc_monitor_gc_pause_total_ns",
				Help: "Total GC pause nanoseconds (PauseTotalNs)",
			},
		),
		goroutinesCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gc_monitor_goroutines_count",
				Help: "Number of active goroutines",
			},
		),
		heapObjectsCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gc_monitor_heap_objects_count",
				Help: "Number of heap objects (HeapObjects)",
			},
		),
	}

	c.ticker = time.NewTicker(5 * time.Second)
	go c.updateMetrics()

	return c
}


func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.memAllocBytes.Describe(ch)
	c.memTotalAllocBytes.Describe(ch)
	c.memSysBytes.Describe(ch)
	c.gcCyclesTotal.Describe(ch)
	c.gcLastTimeSeconds.Describe(ch)
	c.gcPauseTotalNs.Describe(ch)
	c.goroutinesCount.Describe(ch)
	c.heapObjectsCount.Describe(ch)
}


func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.memAllocBytes.Collect(ch)
	c.memTotalAllocBytes.Collect(ch)
	c.memSysBytes.Collect(ch)
	c.gcCyclesTotal.Collect(ch)
	c.gcLastTimeSeconds.Collect(ch)
	c.gcPauseTotalNs.Collect(ch)
	c.goroutinesCount.Collect(ch)
	c.heapObjectsCount.Collect(ch)
}


func (c *Collector) updateMetrics() {
	var lastTotalAlloc uint64 = 0
	var lastNumGC uint32 = 0
	var lastPauseTotalNs uint64 = 0

	for range c.ticker.C {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)


		c.memAllocBytes.Set(float64(memStats.Alloc))

		
		c.memTotalAllocBytes.Add(float64(memStats.TotalAlloc - lastTotalAlloc))
		lastTotalAlloc = memStats.TotalAlloc

		c.memSysBytes.Set(float64(memStats.Sys))


		c.gcCyclesTotal.Add(float64(memStats.NumGC - lastNumGC))
		lastNumGC = memStats.NumGC

		c.gcLastTimeSeconds.Set(float64(memStats.LastGC) / 1e9)

	
		c.gcPauseTotalNs.Add(float64(memStats.PauseTotalNs - lastPauseTotalNs))
		lastPauseTotalNs = memStats.PauseTotalNs

		c.goroutinesCount.Set(float64(runtime.NumGoroutine()))
		c.heapObjectsCount.Set(float64(memStats.HeapObjects))
	}
}


func (c *Collector) Close() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
}