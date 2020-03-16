
test:
	@echo "> Launch tests ..."
	CGO_ENABLED=0 go test -count=1 -v ./...
