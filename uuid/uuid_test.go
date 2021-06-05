/*
 * @file: uuid_test.go
 * @author: Jorge Quitério
 * @copyright (c) 2021 Jorge Quitério
 * @license: MIT
 */

package uuid

import (
	"reflect"
	"testing"
)

func TestNew3(t *testing.T) {
	u := NewV3(NamespaceURL, []byte("golang.go"))
	u2 := NewV3(NamespaceURL, []byte("github.com"))
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
	u := NewV5(NamespaceURL, []byte("golang.go"))
	u2 := NewV5(NamespaceURL, []byte("github.com"))
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
