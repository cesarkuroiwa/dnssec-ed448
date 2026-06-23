package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/netip"

	"github.com/cloudflare/circl/sign/ed448"

	// NOTE: the repo lives at codeberg.org/cesarkuroiwa/dns (branch ed448),
	// but its go.mod declares the module path as codeberg.org/miekg/dns,
	// so that is the path we must import. go.mod has a `replace` directive
	// pointing this path at the cloned ed448 branch.
	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

func main() {
	const zone = "example.org."

	dns.RegisterED448Verifier(func(pub dns.ED448PublicKey, message, sig []byte) bool {
		return ed448.Verify(ed448.PublicKey(pub), message, sig, "")
	})

	// Generate an Ed448 key pair.
	pub, priv, err := ed448.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("generate key: %v", err)
	}

	// Build DNSKEY RR
	dnskey := dns.NewDNSKEY(zone, dns.ED448)
	dnskey.PublicKey = base64.StdEncoding.EncodeToString(pub)

	// Random RRSET to be signed
	addr, _ := netip.ParseAddr("1.1.1.1")
	rrset := []dns.RR{
		&dns.A{
			Hdr: dns.Header{Name: zone, TTL: 3600, Class: dns.ClassINET},
			A:   rdata.A{Addr: addr},
		},
	}

	// Build the RRSIG RR
	sig := dns.NewRRSIG(zone, dns.ED448, dnskey.KeyTag())

	// Sign the RRset.
	if err := sig.Sign(priv, rrset, &dns.SignOption{}); err != nil {
		log.Fatalf("sign: %v", err)
	}
	fmt.Println("signed RRSIG:")
	fmt.Println(" ", sig.String())

	// Verify the signature against the DNSKEY.
	if err := sig.Verify(dnskey, rrset, &dns.SignOption{}); err != nil {
		log.Fatalf("verify failed: %v", err)
	}
	fmt.Println("verify: OK")
}
