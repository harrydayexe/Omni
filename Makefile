all: OmniRead OmniWrite OmniView

OmniRead:
	go build -o tmp/ ./cmd/OmniRead

OmniWrite:
	go build -o tmp/ ./cmd/OmniWrite

OmniView:
	go build -o tmp/ ./cmd/OmniView

Test:
	go test ./...

Cover:
	go test -cover ./...

CoverageReport:
	-go test -coverprofile=tmp/c.out ./...
	go tool cover -html="tmp/c.out" 
