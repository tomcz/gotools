package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

// adapted from https://regex-snippets.com/domain
func TestDomainName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"example.com", true},
		{"subdomain.example.com", true},
		{"sub.domain.example.com", true},
		{"example.co.uk", true},
		{"my-site.com", true},
		{"example123.org", true},
		{"test.example.museum", true},
		{"a.b.c.d.example.com", true},
		{"xn--80akhbyknj4f.com", true},
		{"example", false},
		{".example.com", false},
		{"example.com.", false},
		{"-example.com", false},
		{"example-.com", false},
		{"example..com", false},
		{"example .com", false},
		{"example.c", false},
		{"", false},
		{" ", false},
		{"example.123", false},
		{"exam ple.com", false},
		{"exam_ple.com", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, domainName.MatchString(test.name), test.valid)
		})
	}
}
