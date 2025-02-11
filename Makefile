all: OmniRead OmniWrite OmniView

OmniRead:
	go build -o tmp/ ./cmd/OmniRead

OmniWrite:
	go build -o tmp/ ./cmd/OmniWrite

OmniAuth:
	go build -o tmp/ ./cmd/OmniAuth

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

sqlc:
	sqlc generate

OmniRead-Image: OmniRead
	docker build -t raspidb.local:5000/harrydayexe/omniread -f ./cmd/OmniRead/Dockerfile .
	docker push raspidb.local:5000/harrydayexe/omniread

OmniWrite-Image: OmniWrite
	docker build -t raspidb.local:5000/harrydayexe/omniwrite -f ./cmd/OmniWrite/Dockerfile .
	docker push raspidb.local:5000/harrydayexe/omniwrite

OmniAuth-Image: OmniAuth
	docker build -t raspidb.local:5000/harrydayexe/omniauth -f ./cmd/OmniAuth/Dockerfile .
	docker push raspidb.local:5000/harrydayexe/omniauth
