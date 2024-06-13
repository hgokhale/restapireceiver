package restapireceiver

import (
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"sort"
	"strings"
	"time"
)

type metric struct {
	name       string
	unit       string
	metricType pmetric.MetricType
	value      pcommon.Value
}
type resource struct {
	attributes map[string]any
	metrics    []metric
}

type MetricsBuilder struct {
	metrics        pmetric.Metrics
	resourceLookup map[string]*ResourceBuilder // resource key to ResourceMetrics map
}

func NewMetricsBuilder() *MetricsBuilder {
	return &MetricsBuilder{
		metrics:        pmetric.NewMetrics(),
		resourceLookup: make(map[string]*ResourceBuilder),
	}
}

type ResourceBuilder struct {
	pmetric.ResourceMetrics
}

func (builder MetricsBuilder) GetMetrics() pmetric.Metrics {
	return builder.metrics
}

func (builder *MetricsBuilder) GetOrCreateResource(resourceAttributes map[string]any, scope_name, scope_version string) (*ResourceBuilder, error) {
	key := generateResourceKey(resourceAttributes)
	if key == "" {
		return nil, fmt.Errorf("No attributes were provided for resource")
	}
	if val, ok := builder.resourceLookup[key]; ok {
		return val, nil
	}
	rm := builder.metrics.ResourceMetrics().AppendEmpty()
	err := rm.Resource().Attributes().FromRaw(resourceAttributes)
	if err != nil {
		return nil, err
	}
	sm := rm.ScopeMetrics().AppendEmpty() // single scope metric at 0th position
	sm.Scope().SetName(scope_name)
	sm.Scope().SetVersion(scope_version)

	rmb := &ResourceBuilder{rm}
	builder.resourceLookup[key] = rmb
	return rmb, nil
}

func (rb *ResourceBuilder) AddGaugeMetricDouble(metricName, unit string, value float64, timestamp time.Time) {
	dp := rb.createGaugeMetricDatapoint(metricName, unit)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))
	dp.SetDoubleValue(value)
}

func (rb *ResourceBuilder) AddGaugeMetricInt(metricName, unit string, value int64, timestamp time.Time) {
	dp := rb.createGaugeMetricDatapoint(metricName, unit)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))
	dp.SetIntValue(value)
}

func (rb *ResourceBuilder) createGaugeMetricDatapoint(metricName, unit string) *pmetric.NumberDataPoint {
	newMetric := rb.ResourceMetrics.ScopeMetrics().At(0).Metrics().AppendEmpty()
	newMetric.SetName(metricName)
	newMetric.SetUnit(unit)
	g := newMetric.SetEmptyGauge()
	dp := g.DataPoints().AppendEmpty()
	return &dp
}

func generateResourceKey(attrs map[string]any) string {
	if attrs == nil || len(attrs) == 0 {
		return ""
	}
	// Extract keys from the map
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}

	// Sort the keys
	sort.Strings(keys)

	// Build the result string
	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("_")
		}
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%v", attrs[k]))
	}

	return sb.String()
}
