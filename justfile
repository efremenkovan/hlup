test:
    #!/usr/bin/env sh
    if command -v gotestsum >/dev/null 2>&1; then
    	gotestsum ./...
    else
    	go test ./...
    fi

lint:
    golangci-lint run ./...

fmt:
    golangci-lint fmt "./..."
