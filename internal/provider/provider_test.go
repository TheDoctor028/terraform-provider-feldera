// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"feldera": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	provider, err := testAccProtoV6ProviderFactories["feldera"]()
	assert.Nil(t, err)
	res, err := provider.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestConfigure(t *testing.T) {
	testAccPreCheck(t)

	t.Run("it should return an error if the endpoint is not set", func(t *testing.T) {
		_, err := testAccProtoV6ProviderFactories["feldera"]()

		assert.Nil(t, err)
	})
}
