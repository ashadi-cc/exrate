package helper

import (
	"xrate/common/crate"
	"xrate/common/crate/provider"
	"xrate/config"
)

//CurrentConverterClient returns current converter client from os.environment
func CurrentConverterClient() (provider.Client, error) {
	client, err := crate.Open(config.GetClient().ApiProvider)
	if err != nil {
		return nil, err
	}

	return client, nil
}
