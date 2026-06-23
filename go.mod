module github.com/cesarkuroiwa/ed448

go 1.26.1

replace codeberg.org/miekg/dns => codeberg.org/cesarkuroiwa/dns v0.6.81-ed448

require (
	codeberg.org/miekg/dns v0.6.81
	github.com/cloudflare/circl v1.6.4
)

require (
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
)
