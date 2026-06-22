module github.com/gregtwallace/brother-cert

go 1.26.0

require (
	github.com/peterbourgon/ff/v4 v4.0.0-beta.1
	software.sslmate.com/src/go-pkcs12 v0.7.1
)

require golang.org/x/crypto v0.52.0 // indirect

replace github.com/gregtwallace/brother-cert/cmd/brother-cert => /pkg/cmd/brother-cert

replace github.com/gregtwallace/brother-cert/pkg/app => /pkg/app

replace github.com/gregtwallace/brother-cert/pkg/printer => /pkg/printer
