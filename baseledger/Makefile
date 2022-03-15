PACKAGES=$(shell go list ./... | grep -v '/simulation')

test:
	@go test -mod=readonly $(PACKAGES)