package suite

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/otaku/pretty"
)

var test_yaml = `
# Sample scenarios from pyspf test suite.
description: >-
  check trailing dot with redirect and exp
tests:
  traildot1:
    spec:        8.1
    description: Check that trailing dot is accepted for domains
    host:        192.168.218.40
    helo:        msgbas2x.cos.example.com
    mailfrom:    test@example.com
    result:      pass
  traildot2:
    spec:        8.1
    description: Check that trailing dot is not removed from explanation
    helo:        msgbas2x.cos.example.com
    host:        192.168.218.40
    mailfrom:    test@exp.example.com
    result:      fail
    explanation: This is a test.
zonedata:
  example.com.d.spf.example.com:
    - SPF:  v=spf1 redirect=a.spf.example.com
  a.spf.example.com:
    - SPF:  >-
            v=spf1 mx:example.com include:o.spf.example.com -exists:%{s}.S.bl.spf.example.com
            exists:%{s}.S.%{i}.AI.spf.example.com ~all
  o.spf.example.com:
    - SPF:  v=spf1 ip4:192.168.144.41 ip4:192.168.218.40 ip4:192.168.218.41
  msgbas1x.cos.example.com:
    - A:    192.168.240.36
  example.com:
    - A:    192.168.90.76
    - SPF:  v=spf1 redirect=%{d}.d.spf.example.com.
    - MX:   [10, msgbas1x.cos.example.com]
  exp.example.com:
    - SPF:  v=spf1 exp=msg.example.com. -all
  msg.example.com:
   - TXT:  This is a test.
---
description: test empty MX
tests:
  emptyMX:
    helo:        mail.example.com
    host:        1.2.3.4
    mailfrom:    ''
    result:      neutral
zonedata:
  mail.example.com:
    - MX:   [0, '']
    - SPF:  v=spf1 mx
`

var ss []Scenario

func init() {
	var err error
	ss, err = Load(strings.NewReader(test_yaml))

	if err != nil {
		panic(err)
	}
}

func TestParsing(t *testing.T) {
	pretty.Println(ss)
}

func TestComplete(t *testing.T) {
	var err error
	_, err = LoadFromFile("testdata/rfc7208-tests.yml")

	assert.NoError(t, err)
}
