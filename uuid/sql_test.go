/*
 * @file: sql_test.go
 * @author: Jorge Quitério
 * @copyright (c) 2021 Jorge Quitério
 * @license: MIT
 */

package uuid

import (
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	u := New()

	v, err := u.Value()
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(v).Kind() != reflect.String {
		t.Error("uuid.Value() must be (string, error)")
	}
}

func TestScan(t *testing.T) {
	u := NewV4()
	scanByte := u.Scan(u.Byte())
	scanString := u.Scan(u.String())
	scanUUID := u.Scan(u)
	if scanByte != nil {
		t.Error(scanByte)
	}
	if scanString != nil {
		t.Error(scanString)
	}
	if scanUUID != nil {
		t.Error(scanUUID)
	}
}
