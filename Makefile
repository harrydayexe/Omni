all: OmniRead OmniWrite OmniView

OmniRead:
	go build -o tmp/ ./cmd/OmniRead
	tmp/OmniRead

OmniWrite:
	go build -o tmp/ ./cmd/OmniWrite
	tmp/OmniWrite

OmniView:
	go build -o tmp/ ./cmd/OmniView
	tmp/OmniView
