all: OmniRead OmniWrite OmniView

OmniRead:
	go build -o tmp/ ./cmd/OmniRead

OmniWrite:
	go build -o tmp/ ./cmd/OmniWrite

OmniView:
	go build -o tmp/ ./cmd/OmniView
