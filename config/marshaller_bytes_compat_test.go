package config

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestMarshallerBytesPreservesOnlyRequestedGlobalSection(t *testing.T) {
	glob := Global{
		LogLevel:   "debug",
		TproxyPort: 12345,
	}
	m := Marshaller{IndentSpace: 2, IgnoreZero: true}
	if err := m.MarshalSection("global", reflect.ValueOf(glob), 0); err != nil {
		t.Fatal(err)
	}
	encoded := m.Bytes()
	if !bytes.HasPrefix(encoded, []byte("global {\n")) ||
		!strings.Contains(string(encoded), "log_level:") ||
		!strings.Contains(string(encoded), "tproxy_port:") ||
		strings.Contains(string(encoded), "subscription {") ||
		strings.Contains(string(encoded), "dns {") ||
		strings.Contains(string(encoded), "routing {") {
		t.Fatalf("global-only section was altered: %q", encoded)
	}
	// An exported snapshot of the output must not permit accidental mutation
	// of the marshaller's unexported buffer by future callers.
	encoded[0] ^= 1
	if !bytes.HasPrefix(m.Bytes(), []byte("global {\n")) {
		t.Fatal("mutating returned Bytes changed the underlying config")
	}

	empty := Marshaller{IgnoreZero: true}
	if err := empty.MarshalSection("global", reflect.ValueOf(Global{}), 0); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(empty.Bytes()), "tproxy_port:") ||
		strings.Contains(string(empty.Bytes()), "allow_insecure:") {
		t.Fatal("zero-valued fields were emitted despite IgnoreZero")
	}
}
