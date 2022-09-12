package influx

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/michaelpeterswa/aqi-api-influx-connector/internal/requests"
)

type InfluxConn struct {
	Conn influxdb2.Client
}

func InitInflux(endpoint, token string) *InfluxConn {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(endpoint, token)
	return &InfluxConn{
		Conn: client,
	}
}

func (ic *InfluxConn) Close() {
	ic.Conn.Close()
}

func (ic *InfluxConn) WriteAQIPoint(r *requests.AQIResponse) {
	point := influxdb2.NewPointWithMeasurement("aqi").
		AddTag("primary_pollutant", r.PrimaryPollutant).
		AddField("aqi", r.AQI).
		AddField("level", r.Level).
		SetTime(time.Now())

	writeAPI := ic.Conn.WriteAPI("main", "aqi")
	writeAPI.WritePoint(point)
	writeAPI.Flush()
}
