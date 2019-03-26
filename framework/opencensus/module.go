package opencensus

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/systemendpoint/domain"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

var (
	registerOnce = new(sync.Once)
	// KeyArea is the key to represent the current flamingo area
	KeyArea, _ = tag.NewKey("area")
)

// View helps to register opencensus views with the default "area" tag
func View(name string, m stats.Measure, aggr *view.Aggregation, tagKeys ...tag.Key) error {
	return view.Register(&view.View{
		Name:        name,
		Measure:     m,
		Aggregation: aggr,
		TagKeys:     append([]tag.Key{KeyArea}, tagKeys...),
	})
}

type correlationIDInjector struct {
	next http.RoundTripper
}

// RoundTrip a request
func (rt *correlationIDInjector) RoundTrip(req *http.Request) (*http.Response, error) {
	if span := trace.FromContext(req.Context()); span != nil {
		req.Header.Add("X-Correlation-ID", span.SpanContext().TraceID.String())
	}

	return rt.next.RoundTrip(req)
}

// Module registers the opencensus module which in turn enables jaeger & co
type Module struct {
	Endpoint     string `inject:"config:opencensus.jaeger.endpoint"`
	ServiceName  string `inject:"config:opencensus.serviceName"`
	ServiceAddr  string `inject:"config:opencensus.serviceAddr"`
	JaegerEnable bool   `inject:"config:opencensus.jaeger.enable"`
}

// Configure the opencensus Module
func (m *Module) Configure(injector *dingo.Injector) {
	registerOnce.Do(func() {
		// For demoing purposes, always sample.
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.NeverSample()})
		http.DefaultTransport = &correlationIDInjector{next: &ochttp.Transport{Base: http.DefaultTransport}}

		if m.JaegerEnable {
			// generate a random IP in 127.0.0.0/8 network to trick jaegers clock skew detection
			// todo fix this?
			randomIP := fmt.Sprintf("127.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255))

			// Register the Jaeger exporter to be able to retrieve
			// the collected spans.
			exporter, err := jaeger.NewExporter(jaeger.Options{
				CollectorEndpoint: m.Endpoint,
				Process: jaeger.Process{
					ServiceName: m.ServiceName,
					Tags: []jaeger.Tag{
						jaeger.StringTag("ip", randomIP),
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			trace.RegisterExporter(exporter)
		}

		{
			exporter, err := prometheus.NewExporter(prometheus.Options{})
			if err != nil {
				log.Fatal(err)
			}
			view.RegisterExporter(exporter)
			injector.BindMap((*domain.Handler)(nil), "/metrics").ToInstance(exporter)
		}
	})
}

// DefaultConfig for opencensus module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"opencensus": config.Map{
			"jaeger.endpoint": "http://localhost:14268/api/traces",
			"jaeger.enable":   false,
			"serviceName":     "flamingo",
			"serviceAddr":     ":13210",
			"tracing": config.Map{
				"sampler": config.Map{
					"whitelist":        config.Slice{},
					"blacklist":        config.Slice{},
					"allowParentTrace": true,
				},
			},
		},
	}
}
