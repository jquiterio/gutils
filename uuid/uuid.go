/*
 * @file: uuid.go
 * @author: Jorge Quitério
 * @copyright (c) 2021 Jorge Quitério
 * @license: MIT
 */

package uuid

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"math/rand"
	"regexp"
)

// as in https://datatracker.ietf.org/doc/html/rfc4122#section-4.1.1
const (
	VarNCS       byte = 0x80
	VarRFC4122   byte = 0x40
	VarMicrosoft byte = 0x20
	VarFuture    byte = 0x00
)

// UUID as in RFC 4122 https://datatracker.ietf.org/doc/html/rfc4122
type UUID [16]byte

type NullUUID struct {
	UUID  UUID
	Valid bool
}

// Nil as in https://datatracker.ietf.org/doc/html/rfc4122#section-4.1.7
var Nil = UUID{}

// as in https://datatracker.ietf.org/doc/html/rfc4122#appendix-C
var (
	NamespaceDNS, _  = ParseHex("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = ParseHex("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = ParseHex("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = ParseHex("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

const hexPattern = "^(urn\\:uuid\\:)?\\{?([a-z0-9]{8})-([a-z0-9]{4})-" +
	"([1-5][a-z0-9]{3})-([a-z0-9]{4})-([a-z0-9]{12})\\}?$"

// ParseHex Crates uuid from hex formated string
func ParseHex(s string) (u *UUID, err error) {
	var re = regexp.MustCompile(hexPattern)
	sm := re.FindStringSubmatch(s)
	if sm == nil {
		err = errors.New("invalid uuid")
		return
	}
	h := sm[2] + sm[4] + sm[5] + sm[6]
	b, err := hex.DecodeString(h)
	if err != nil {
		return
	}
	u = new(UUID)
	copy(u[:], b)
	return
}

func Parse(b []byte) (u *UUID, err error) {
	if len(b) != 16 {
		err = errors.New("invalid UUID")
		return
	}
	u = new(UUID)
	copy(u[:], b)
	return
}

// New create a new V4 UUID
func New() UUID {
	return NewV4()
}

func (u *UUID) setHash(h hash.Hash, ns, name []byte) {
	h.Write(ns[:])
	h.Write(name)
	copy(u[:], h.Sum([]byte{})[:16])
}

func NewV3(ns *UUID, name []byte) UUID {
	if ns == nil {
		panic("invalid namespace")
	}
	var u UUID
	u.setHash(md5.New(), ns[:], name)
	u.setVariant(VarRFC4122)
	u.setVersion(3)
	return u
}

// New create a new V4 UUID
func NewV4() UUID {
	var u UUID
	_, _ = rand.Read(u[:])
	u.setVariant(VarRFC4122)
	u.setVersion(4)
	return u
}

func NewV5(ns *UUID, name []byte) UUID {
	var u UUID
	u.setHash(sha1.New(), ns[:], name)
	u.setVariant(VarRFC4122)
	u.setVersion(5)
	return u
}

func (u *UUID) setVariant(v byte) {
	switch v {
	case VarNCS:
		u[8] = (u[8] | VarNCS) & 0xBF
	case VarRFC4122:
		u[8] = (u[8] | VarRFC4122) & 0x7F
	case VarMicrosoft:
		u[8] = (u[8] | VarMicrosoft) & 0x3F
	}
}

// Variant output the UUID variant
func (u UUID) Variant() byte {
	switch {
	case (u[8] >> 7) == 0x00:
		return VarNCS
	case (u[8] >> 6) == 0x02:
		return VarRFC4122
	case (u[8] >> 5) == 0x06:
		return VarMicrosoft
	case (u[8] >> 5) == 0x07:
		fallthrough
	default:
		return VarFuture
	}
}

// Equal checks if current UUID is equals to another one
func (u *UUID) Equal(ub UUID) bool {
	return bytes.Equal(u[:], ub[:])
}

// Compare compares two uuids
func Compare(ua, ub UUID) bool {
	return bytes.Equal(ua[:], ub[:])
}

func CompareAll(us ...UUID) bool {
	u1 := us[0]
	for _, u := range us {
		if u1 != u {
			return false
		}
	}
	return true
}

func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func (u *UUID) Byte() []byte {
	return u[:]
}

func (u *UUID) setVersion(v byte) {
	u[6] = (u[6] & 0xF) | (v << 4)
}

func (u *UUID) Version() uint8 {
	return uint8(u[6] >> 4)
}

func FromString(s string) (UUID, error) {
	var u UUID
	err := u.UnmarshalText([]byte(s))
	return u, err
}

func (u *UUID) UnmarshalText(b []byte) (err error) {
	err = u.decode(b)
	return
}

func (u *UUID) decode(b []byte) (err error) {
	u, err = ParseHex(string(b))
	return err
}
