package restapireceiver

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/hgokhale/restapireceiver/internal/metadata"
)

var errConfigNotRestAPIConfig = errors.New("config was not a restapi receiver config")

// NewFactory creates a new receiver factory for SNMP
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability))
}

// createDefaultConfig creates a config with defaults
func createDefaultConfig() component.Config {
	return &Config{
		ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
	}
}

// createMetricsReceiver creates the metric receiver
func createMetricsReceiver(
	_ context.Context,
	params receiver.CreateSettings,
	config component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	recvConfig, ok := config.(*Config)
	if !ok {
		return nil, errConfigNotRestAPIConfig
	}

	if err := adjustConfigAndValidate(recvConfig); err != nil {
		return nil, fmt.Errorf("failed to validate added config defaults: %w", err)
	}

	restapiScraper := newScraper(params.Logger, recvConfig, params)
	scraper, err := scraperhelper.NewScraper(metadata.Type.String(), restapiScraper.scrape, scraperhelper.WithStart(restapiScraper.start))
	if err != nil {
		return nil, err
	}

	return scraperhelper.NewScraperControllerReceiver(&recvConfig.ControllerConfig, params, consumer, scraperhelper.AddScraper(scraper))
}

// adjustConfigAndValidate adds any missing config parameters that have defaults
func adjustConfigAndValidate(cfg *Config) error {
	//TODO adjust configs if needed
	return component.ValidateConfig(cfg)
}
