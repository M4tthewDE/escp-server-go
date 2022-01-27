.PHONY: lint 

lint:
	golangci-lint run . internal/api/... internal/db/...