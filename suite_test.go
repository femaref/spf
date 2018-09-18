package spf

import (
	"flag"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/femaref/spf/suite"
	"github.com/miekg/dns"
)

var ss string

func init() {
	flag.StringVar(&ss, "filter", "", "specify a single test scenario")
}

func TestSuite(t *testing.T) {

	for _, scenario := range scenarios {
		if ss == "" || ss == scenario.Description {
			t.Run(scenario.Description, runner(scenario))
		}
	}
}

func runner(scenario suite.Scenario) func(t *testing.T) {
	return func(t *testing.T) {
		for fqdn, zd := range scenario.ZoneData {
			data := map[uint16][]string{}
			fqdn = NormalizeFQDN(fqdn)
			for _, single := range zd {
				data[single.Type] = append(data[single.Type], single.Text(fqdn))
			}

			dns.HandleFunc(fqdn, zone(data))
			defer dns.HandleRemove(fqdn)
		}
		for name, test := range scenario.Tests {
			ip := net.ParseIP(test.Host)
			if !assert.NotNilf(t, ip, "ip: %s: %s", scenario.Description, name) {
				continue
			}

			domain := test.MailFrom
			at_index := strings.Index(domain, "@")

			if at_index > -1 {
				after_at := at_index + 1
				if after_at == len(domain) {
					domain = ""
				}
				domain = domain[after_at:]
			}

			assert.NotPanicsf(t, func() {
				result, _, err := CheckHostWithResolver(ip, domain, test.Helo, NewLimitedResolver(testResolver, 10, 10))

				r := result.String()
				switch result {
				case Permerror, Temperror:
				default:
					if err != ErrInvalidDomain {
						assert.NoErrorf(t, err, "err: %s: %s\n%s", scenario.Description, name, test.Description)
					}
				}
				assert.Containsf(t, test.Result, r, "result: %s: %s\n%s", scenario.Description, name, test.Description)

			}, "no panic: %s: %s\n%s", scenario.Description, name, test.Description)
		}
	}
}
