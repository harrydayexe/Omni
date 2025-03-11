all: OmniRead OmniWrite OmniView

.PHONY: OmniRead
OmniRead: # Build the OmniRead binary
	go build -o tmp/ ./cmd/OmniRead

.PHONY: OmniWrite
OmniWrite: # Build the OmniWrite binary
	go build -o tmp/ ./cmd/OmniWrite

.PHONY: OmniAuth
OmniAuth: # Build the OmniAuth binary
	go build -o tmp/ ./cmd/OmniAuth

.PHONY: OmniView
OmniView: # Build the OmniView binary
	tailwindcss -i "./internal/omniview/templates/custom.css" -o "./internal/omniview/templates/static/style.css"
	go build -o tmp/ ./cmd/OmniView

.PHONY: LoadBalancer
LoadBalancer: # Build the LoadBalancer binary
	go build -o tmp/ ./cmd/LoadBalancer

.PHONY: Test
Test: # Run all tests
	go test ./...

.PHONY: Cover
Cover: # Run all tests with coverage
	go test -coverprofile=tmp/c.out ./...
	go tool cover -html="tmp/c.out" 

.PHONY: sqlc
sqlc: # Generate sqlc code
	sqlc generate

.PHONY: OmniRead-Image
OmniRead-Image: OmniRead # Build the OmniRead docker image
	docker build -t raspidb.local:5000/harrydayexe/omniread -f ./cmd/OmniRead/Dockerfile .
	docker push raspidb.local:5000/harrydayexe/omniread

.PHONY: OmniWrite-Image
OmniWrite-Image: OmniWrite # Build the OmniWrite docker image
	docker build -t raspidb.local:5000/harrydayexe/omniwrite -f ./cmd/OmniWrite/Dockerfile .
	docker push raspidb.local:5000/harrydayexe/omniwrite

.PHONY: OmniAuth-Image
OmniAuth-Image: OmniAuth # Build the OmniAuth docker image
	docker build -t raspidb.local:5000/harrydayexe/omniauth -f ./cmd/OmniAuth/Dockerfile .
	docker push raspidb.local:5000/harrydayexe/omniauth

.PHONY: OmniView-Image
OmniView-Image: OmniView # Build the OmniView docker image
	docker build -t raspidb.local:5000/harrydayexe/omniview -f ./cmd/OmniView/Dockerfile .
	docker push raspidb.local:5000/harrydayexe/omniview

.PHONY: Push-Images
Push-Images: OmniView-Image OmniRead-Image OmniAuth-Image OmniWrite-Image
