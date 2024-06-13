package restapireceiver

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"time"
)

// restapiScraper handle scraping of metrics
type restapiScraper struct {
	client    *HttpClientHelper
	logger    *zap.Logger
	cfg       *Config
	settings  receiver.CreateSettings
	startTime pcommon.Timestamp
}

// newScraper creates and initializes restapiScraper
func newScraper(logger *zap.Logger, cfg *Config, settings receiver.CreateSettings) *restapiScraper {
	return &restapiScraper{
		logger:   logger,
		cfg:      cfg,
		settings: settings,
	}
}

// start gets the Client ready
func (s *restapiScraper) start(_ context.Context, _ component.Host) error {
	//TODO implement start
	return nil
}

// scrape collects and creates OTEL metrics from a SNMP environment
func (s *restapiScraper) scrape(_ context.Context) (pmetric.Metrics, error) {
	s.logger.Warn("scrape() not Implemented")
	return getDummyMetrics(), nil
	//return pmetric.NewMetrics(), fmt.Errorf("Not implemented")
}

func getDummyMetrics() pmetric.Metrics {

	// metadata
	scope_name, scope_version := "otelcol/restapireceiver", "0.0.1"

	total_capacity_metric_name, used_capacity_metric_name := "total_capacity", "used_capacity"
	attr_cluster_name, attr_node_name, attr_ip := "cluster_name", "node_name", "ip"
	unit_kilo_bytes := "KiBy"

	// data
	clusterName := "cluster1"
	node1Name := "node1"
	node1IP := "10.10.1.20"
	node2Name := "node2"
	node2IP := "10.10.1.21"
	var node1capacity, node2capacity int64 = 2 * 1024 * 1024 * 1024, 1024 * 1024 * 1024 // in kb
	clustercapacity := node1capacity + node2capacity
	var node1usage, node2usage int64 = 67 * 1024 * 1024, 25 * 1024 * 1024 // in kb
	clusterusage := node1usage + node2usage

	timestamp := time.Now().UTC()

	builder := NewMetricsBuilder()

	// create cluster
	clusterAttrs := map[string]any{
		attr_cluster_name: clusterName,
	}
	clusterRb, err := builder.GetOrCreateResource(clusterAttrs, scope_name, scope_version)
	if err != nil {
		fmt.Printf("failed to create cluster: %v", err.Error())
	} else {
		// add cluster metrics
		clusterRb.AddGaugeMetricInt(total_capacity_metric_name, unit_kilo_bytes, clustercapacity, timestamp)
		clusterRb.AddGaugeMetricInt(used_capacity_metric_name, unit_kilo_bytes, clusterusage, timestamp)
	}
	// create node1
	node1attrs := map[string]any{
		attr_cluster_name: clusterName,
		attr_node_name:    node1Name,
		attr_ip:           node1IP,
	}
	node1Rb, err := builder.GetOrCreateResource(node1attrs, scope_name, scope_version)
	if err != nil {
		fmt.Printf("failed to create node1: %v", err.Error())
	} else {
		// add node1 metrics
		node1Rb.AddGaugeMetricInt(total_capacity_metric_name, unit_kilo_bytes, node1capacity, timestamp)
		node1Rb.AddGaugeMetricInt(used_capacity_metric_name, unit_kilo_bytes, node1usage, timestamp)
	}

	// create node1
	node2attrs := map[string]any{
		attr_cluster_name: clusterName,
		attr_node_name:    node2Name,
		attr_ip:           node2IP,
	}
	node2Rb, err := builder.GetOrCreateResource(node2attrs, scope_name, scope_version)
	if err != nil {
		fmt.Printf("failed to create node2: %v", err.Error())
	} else {
		// add node1 metrics
		node2Rb.AddGaugeMetricInt(total_capacity_metric_name, unit_kilo_bytes, node2capacity, timestamp)
		node2Rb.AddGaugeMetricInt(used_capacity_metric_name, unit_kilo_bytes, node2usage, timestamp)
	}

	return builder.GetMetrics()
}
