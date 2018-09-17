package suite

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Scenario struct {
	Description string          `yaml:"description"`
	Tests       map[string]Test `yaml:"tests"`
	ZoneData    ZoneData        `yaml:"zonedata"`
}

type Test struct {
	Spec        MultiString
	Description string
	Host        string
	Helo        string
	MailFrom    string
	Result      MultiString
	Explanation string
}

type MultiString []string

func (this *MultiString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var foo interface{}

	err := unmarshal(&foo)
	if err != nil {
		errors.Wrap(err, "spec")
	}

	switch x := foo.(type) {
	case []interface{}:
		for _, s := range x {
			(*this) = append((*this), fmt.Sprint(s))
		}
	case string:
		(*this) = []string{x}
	}

	return nil
}

type ZoneData map[string][]Single

type Single struct {
	Type      uint16
	String    string
	Array     []interface{}
	IsString  bool
	IsSpecial bool
}

func (this Single) Text(fqdn string) string {
	if this.IsSpecial {
		return this.String
	}
	if this.IsString {
		return fmt.Sprintf("%s 10 IN %s %s", fqdn, dns.TypeToString[this.Type], this.String)
	}
	var repr []string
	for _, a := range this.Array {
		repr = append(repr, fmt.Sprint(a))
	}
	return fmt.Sprintf("%s 10 IN %s %s", fqdn, dns.TypeToString[this.Type], strings.Join(repr, " "))
}

func (this *ZoneData) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if *this == nil {
		(*this) = make(ZoneData)
	}

	var ms yaml.MapSlice

	err := unmarshal(&ms)

	if err != nil {
		return errors.Wrap(err, "zonedata")
	}

	// fqdns
	for _, mi := range ms {
		fqdn := mi.Key.(string)
		zdss := mi.Value.([]interface{})

		//
		for _, z := range zdss {

			switch x := z.(type) {
			case yaml.MapSlice:
				for _, zd := range x {
					key := strings.ToUpper(zd.Key.(string))

					var type_key uint16

					for t, s := range dns.TypeToString {
						if s == key {
							type_key = t
						}
					}

					if type_key == dns.TypeNone {
						return errors.Errorf("unknown type: %s", key)
					}

					var result Single = Single{Type: type_key}
					switch key {
					case "TXT", "SPF":
						switch x := zd.Value.(type) {
						case string:
							result.Array = []interface{}{fmt.Sprintf("\"%s\"", x)}
						case []interface{}:
							result.Array = x
						}
					default:
						switch x := zd.Value.(type) {
						case string:
							result.String = x
							result.IsString = true
						case []interface{}:
							result.Array = x
						}
					}
					(*this)[fqdn] = append((*this)[fqdn], result)
				}

			case string:
				(*this)[fqdn] = append((*this)[fqdn], Single{String: x, IsSpecial: true})
			}

		}
	}

	return nil
}

func LoadFromFile(path string) ([]Scenario, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return Load(f)
}

func Load(r io.Reader) ([]Scenario, error) {

	scanner := bufio.NewScanner(r)
	scanner.Split(splitYAMLDocument)
	var output []Scenario

	for scanner.Scan() {
		var foo Scenario
		err := yaml.Unmarshal(scanner.Bytes(), &foo)

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		output = append(output, foo)
	}
	return output, nil
}
