package influxdb

// MetricWriter represents component able to send metrics data
// to somewhere
type MetricWriter interface {
	Write(data interface {
		ForEachMetric(func(name string, value int64, tags map[string]string))
	}) error
}
