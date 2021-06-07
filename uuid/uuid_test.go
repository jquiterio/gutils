/*
 * @file: uuid_test.go
 * @author: Jorge Quitério
 * @copyright (c) 2021 Jorge Quitério
 * @license: MIT
 */

package uuid

import (
	"reflect"
	"strings"
	"testing"
)

func TestNew3(t *testing.T) {
	u := NewV3(&NamespaceURL, []byte("http://golang.go"))
	u2 := NewV3(&NamespaceURL, []byte("http://github.com"))
	if u == Nil {
		t.Error("Error generating V5 UUID")
	}
	if u.Equal(u2) {
		t.Errorf("%v cannot be the same as %v", u, u2)
	}
	tostr := u.String()
	if reflect.TypeOf(tostr).Kind() != reflect.String {
		t.Errorf("uuid not converted to a string")
	}
	if len(tostr) != 36 {
		t.Errorf("%v is not 36 long", tostr)
	}
	if u.Version() != 3 {
		t.Errorf("%v not == 3", u.Version())
	}
}

func TestNewV4(t *testing.T) {
	u := NewV4()
	u2 := NewV4()
	if u == Nil {
		t.Error("Error generating V4 UUID")
	}
	if u.Equal(u2) {
		t.Errorf("%v cannot be the same as %v", u, u2)
	}
	tostr := u.String()
	if reflect.TypeOf(tostr).Kind() != reflect.String {
		t.Errorf("uuid not converted to a string")
	}
	if len(tostr) != 36 {
		t.Errorf("%v is not 36 long", tostr)
	}
	if u.Version() != 4 {
		t.Errorf("%v not == 4", u.Version())
	}
}

func TestNew5(t *testing.T) {
	u := NewV5(&NamespaceURL, []byte("golang.go"))
	u2 := NewV5(&NamespaceURL, []byte("github.com"))
	if u == Nil {
		t.Error("Error generating V5 UUID")
	}
	if u.Equal(u2) {
		t.Errorf("%v cannot be the same as %v", u, u2)
	}
	tostr := u.String()
	if reflect.TypeOf(tostr).Kind() != reflect.String {
		t.Errorf("uuid not converted to a string")
	}
	if len(tostr) != 36 {
		t.Errorf("%v is not 36 long", tostr)
	}
	if u.Version() != 5 {
		t.Errorf("%v not == 5", u.Version())
	}
}

func TestParse(t *testing.T) {

	uuids := []string{
		"efb65913-a881-4006-bce4-9fc33be32ddd",
		"efb65913a8814006bce49fc33be32ddd",
		"{efb65913-a881-4006-bce4-9fc33be32ddd}",
		"urn:uuid:efb65913-a881-4006-bce4-9fc33be32ddd",
		"urn:uuid:efb65913a8814006bce49fc33be32ddd",
	}
	for _, u := range uuids {
		uuid, err := Parse(u)
		if err != nil {
			t.Errorf("%v is invalid", uuid.String())
		}
	}
}

func TestFormat(t *testing.T) {
	uuids := []string{
		"efb65913-a881-4006-bce4-9fc33be32ddd",
		"efb65913a8814006bce49fc33be32ddd",
		"{efb65913-a881-4006-bce4-9fc33be32ddd}",
		"urn:uuid:efb65913-a881-4006-bce4-9fc33be32ddd",
		"urn:uuid:efb65913a8814006bce49fc33be32ddd",
	}
	for _, u := range uuids {
		s := Format(u)
		if s == "" || len(s) != 36 {
			t.Errorf("%v is invalid", u)
		}
	}
}

func TestURN(t *testing.T) {
	u := New()
	s := u.URN()
	t.Log(s)
	if !strings.Contains(s, "urn:uuid:") && len(s) != 45 {
		t.Errorf("%v is not a valid urn format uuid", s)
	}
}

func TestVariant(t *testing.T) {
	u := New()
	if u.Variant() != VarRFC4122 {
		t.Errorf("%v is not RC4122 compliant", u.Variant())
	}
}
