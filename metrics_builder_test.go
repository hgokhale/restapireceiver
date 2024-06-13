package restapireceiver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestGenerateResourceKey(t *testing.T) {
	tests := []struct {
		name     string
		attrs    map[string]any
		expected string
	}{
		{
			name:     "NilAttributes",
			attrs:    nil,
			expected: "",
		},
		{
			name:     "EmptyAttributes",
			attrs:    map[string]any{},
			expected: "",
		},
		{
			name: "SingleAttribute",
			attrs: map[string]any{
				"a": "v1",
			},
			expected: "a:v1",
		},
		{
			name: "MultipleAttributes",
			attrs: map[string]any{
				"b": "v2",
				"a": "v1",
			},
			expected: "a:v1_b:v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := generateResourceKey(tt.attrs)
			assert.Equal(t, tt.expected, key)
		})
	}
}

func TestMetricsBuilder_GetOrCreateResource(t *testing.T) {
	mb := NewMetricsBuilder()
	resourceAttrs := map[string]any{"service": "test-service"}

	rb, err := mb.GetOrCreateResource(resourceAttrs, "scope", "v1")
	assert.NoError(t, err)
	assert.NotNil(t, rb)

	// Check if the same resource is returned for the same attributes
	rb2, err := mb.GetOrCreateResource(resourceAttrs, "scope", "v1")
	assert.NoError(t, err)
	assert.Equal(t, rb, rb2)

	// Check if a new resource is created for different attributes
	resourceAttrs2 := map[string]any{"service": "other-service"}
	rb3, err := mb.GetOrCreateResource(resourceAttrs2, "scope", "v1")
	assert.NoError(t, err)
	assert.NotEqual(t, rb, rb3)
}

func TestResourceBuilder_AddGaugeMetricDouble(t *testing.T) {
	mb := NewMetricsBuilder()
	rb, err := mb.GetOrCreateResource(map[string]any{"service": "test-service"}, "scope", "v1")
	assert.NoError(t, err)

	metricName := "test_metric_double"
	unit := "ms"
	value := 123.456
	timestamp := time.Now()

	rb.AddGaugeMetricDouble(metricName, unit, value, timestamp)

	metrics := rb.ResourceMetrics.ScopeMetrics().At(0).Metrics()
	assert.Equal(t, 1, metrics.Len())
	m := metrics.At(0)
	assert.Equal(t, metricName, m.Name())
	assert.Equal(t, unit, m.Unit())

	dp := m.Gauge().DataPoints().At(0)
	assert.Equal(t, value, dp.DoubleValue())
	assert.Equal(t, pcommon.NewTimestampFromTime(timestamp), dp.Timestamp())
}

func TestResourceBuilder_AddGaugeMetricInt(t *testing.T) {
	mb := NewMetricsBuilder()
	rb, err := mb.GetOrCreateResource(map[string]any{"service": "test-service"}, "scope", "v1")
	assert.NoError(t, err)

	metricName := "test_metric_int"
	unit := "ms"
	value := int64(123)
	timestamp := time.Now()

	rb.AddGaugeMetricInt(metricName, unit, value, timestamp)

	metrics := rb.ResourceMetrics.ScopeMetrics().At(0).Metrics()
	assert.Equal(t, 1, metrics.Len())
	m := metrics.At(0)
	assert.Equal(t, metricName, m.Name())
	assert.Equal(t, unit, m.Unit())

	dp := m.Gauge().DataPoints().At(0)
	assert.Equal(t, value, dp.IntValue())
	assert.Equal(t, pcommon.NewTimestampFromTime(timestamp), dp.Timestamp())
}

func TestMetricsBuilder_GetMetrics(t *testing.T) {
	mb := NewMetricsBuilder()

	resourceAttrs := map[string]any{"service": "test-service"}
	rb, err := mb.GetOrCreateResource(resourceAttrs, "scope", "v1")
	assert.NoError(t, err)

	metricName := "test_metric"
	unit := "ms"
	value := 123.456
	timestamp := time.Now()

	rb.AddGaugeMetricDouble(metricName, unit, value, timestamp)

	metrics := mb.GetMetrics()
	assert.NotNil(t, metrics)

	resourceMetrics := metrics.ResourceMetrics()
	assert.Equal(t, 1, resourceMetrics.Len())
	rm := resourceMetrics.At(0)
	res := rm.Resource()
	attrMap := res.Attributes()
	attrVal, ok := attrMap.Get("service")
	assert.True(t, ok)
	assert.Equal(t, pcommon.NewValueStr("test-service"), attrVal)

	scopeMetrics := rm.ScopeMetrics()
	assert.Equal(t, 1, scopeMetrics.Len())
	sm := scopeMetrics.At(0)

	m := sm.Metrics().At(0)
	assert.Equal(t, metricName, m.Name())
	assert.Equal(t, unit, m.Unit())

	dp := m.Gauge().DataPoints().At(0)
	assert.Equal(t, value, dp.DoubleValue())
	assert.Equal(t, pcommon.NewTimestampFromTime(timestamp), dp.Timestamp())
}
