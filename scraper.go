package restapireceiver

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"time"
)

// restapiScraper handle scraping of metrics
type restapiScraper struct {
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

	timestamp := pcommon.NewTimestampFromTime(time.Now().UTC())

	metrics := pmetric.NewMetrics()

	// create cluster
	clusterResourceMetrics := metrics.ResourceMetrics().AppendEmpty()
	clusterResourceMetrics.Resource().Attributes().PutStr(attr_cluster_name, clusterName)
	clusterScopeMetrics := clusterResourceMetrics.ScopeMetrics().AppendEmpty()
	clusterScopeMetrics.Scope().SetName(scope_name)
	clusterScopeMetrics.Scope().SetVersion(scope_version)

	// create node1
	node1ResourceMetrics := metrics.ResourceMetrics().AppendEmpty()
	node1ResourceMetrics.Resource().Attributes().PutStr(attr_cluster_name, clusterName)
	node1ResourceMetrics.Resource().Attributes().PutStr(attr_node_name, node1Name)
	node1ResourceMetrics.Resource().Attributes().PutStr(attr_ip, node1IP)
	node1ScopeMetrics := node1ResourceMetrics.ScopeMetrics().AppendEmpty()
	node1ScopeMetrics.Scope().SetName("otelcol/restapireceiver")
	node1ScopeMetrics.Scope().SetVersion("0.0.1")

	// create node2
	node2ResourceMetrics := metrics.ResourceMetrics().AppendEmpty()
	node2ResourceMetrics.Resource().Attributes().PutStr(attr_cluster_name, clusterName)
	node2ResourceMetrics.Resource().Attributes().PutStr(attr_node_name, node2Name)
	node2ResourceMetrics.Resource().Attributes().PutStr(attr_ip, node2IP)
	node2ScopeMetrics := node2ResourceMetrics.ScopeMetrics().AppendEmpty()
	node2ScopeMetrics.Scope().SetName("otelcol/restapireceiver")
	node2ScopeMetrics.Scope().SetVersion("0.0.1")

	// add cluster metrics
	newMetric := clusterScopeMetrics.Metrics().AppendEmpty()
	newMetric.SetName(total_capacity_metric_name)
	newMetric.SetUnit(unit_kilo_bytes)
	newMetric.SetEmptyGauge()
	dp := newMetric.Gauge().DataPoints().AppendEmpty()
	dp.SetIntValue(clustercapacity)
	dp.SetTimestamp(timestamp)

	newMetric = clusterScopeMetrics.Metrics().AppendEmpty()
	newMetric.SetName(used_capacity_metric_name)
	newMetric.SetUnit(unit_kilo_bytes)
	newMetric.SetEmptyGauge()
	dp = newMetric.Gauge().DataPoints().AppendEmpty()
	dp.SetIntValue(clusterusage)
	dp.SetTimestamp(timestamp)

	// add node1 metrics
	newMetric = node1ScopeMetrics.Metrics().AppendEmpty()
	newMetric.SetName(total_capacity_metric_name)
	newMetric.SetUnit(unit_kilo_bytes)
	newMetric.SetEmptyGauge()
	dp = newMetric.Gauge().DataPoints().AppendEmpty()
	dp.SetIntValue(node1capacity)
	dp.SetTimestamp(timestamp)

	newMetric = node1ScopeMetrics.Metrics().AppendEmpty()
	newMetric.SetName(used_capacity_metric_name)
	newMetric.SetUnit(unit_kilo_bytes)
	newMetric.SetEmptyGauge()
	dp = newMetric.Gauge().DataPoints().AppendEmpty()
	dp.SetIntValue(node1usage)
	dp.SetTimestamp(timestamp)

	// add node2 metrics
	newMetric = node2ScopeMetrics.Metrics().AppendEmpty()
	newMetric.SetName(total_capacity_metric_name)
	newMetric.SetUnit(unit_kilo_bytes)
	newMetric.SetEmptyGauge()
	dp = newMetric.Gauge().DataPoints().AppendEmpty()
	dp.SetIntValue(node2capacity)
	dp.SetTimestamp(timestamp)

	newMetric = node2ScopeMetrics.Metrics().AppendEmpty()
	newMetric.SetName(used_capacity_metric_name)
	newMetric.SetUnit(unit_kilo_bytes)
	newMetric.SetEmptyGauge()
	dp = newMetric.Gauge().DataPoints().AppendEmpty()
	dp.SetIntValue(node2usage)
	dp.SetTimestamp(timestamp)

	return metrics
}
