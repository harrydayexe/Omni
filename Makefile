all: OmniRead OmniWrite OmniView

OmniRead:
	go build -o tmp/ ./cmd/OmniRead

OmniWrite:
	go build -o tmp/ ./cmd/OmniWrite

OmniView:
	go build -o tmp/ ./cmd/OmniView

LoadBalancer:
	go build -o tmp/ ./cmd/LoadBalancer

Test:
	go test ./...

Cover:
	go test -cover ./...

CoverageReport:
	-go test -coverprofile=tmp/c.out ./...
	go tool cover -html="tmp/c.out" 

migrate-up: compose-up
	migrate -database "mysql://root:Password1!@tcp(localhost:3306)/omni" -path ./db/migrations up

migrate-down:
	migrate -database "mysql://root:Password1!@tcp(localhost:3306)/omni" -path ./db/migrations down

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down
