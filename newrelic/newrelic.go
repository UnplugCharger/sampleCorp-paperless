package newrelic

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/qwetu_petro/backend/utils"
)

func RelicApp(conf utils.Config) (*newrelic.Application, error) {

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(conf.NewRelicAppName),
		newrelic.ConfigLicense(conf.NewRelicLicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		fmt.Println(conf.NewRelicAppName)
		fmt.Println(conf.NewRelicLicenseKey)
		panic(err)
	}

	return app, nil

}
