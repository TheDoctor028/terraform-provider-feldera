TERRAFORM_D=${HOME}/.terraform.d
TERRAFORM_PLUGINS=${TERRAFORM_D}/plugins
VERSION=1.0.0
ARCH=$(shell uname -o | perl -ne 'print lc')_$(shell uname -m)


gen-api-v0:
	oapi-codegen -package v0 openapi.yaml > api.gen.go

build-local:
	go build -o ${TERRAFORM_PLUGINS}/terraform.local/local/feldera/${VERSION}/${ARCH}/terraform-provider-feldera_v${VERSION}
