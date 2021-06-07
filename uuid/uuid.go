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
	"hash"
	"io"
	"math/rand"
	"regexp"
	"strings"
)

// as in https://datatracker.ietf.org/doc/html/rfc4122#section-4.1.1
const (
	VarNCS       byte = 0x80
	VarRFC4122   byte = 0x40
	VarMicrosoft byte = 0x20
	VarFuture    byte = 0x00
)

//const hexPattern = `({?|urn:uuid:)([aA-fF0-9]{8})(-?)([aA-fF0-9]{4})(-?)([1-5][aA-fF0-9]{3})(-?)([aA-fF0-9]{4})(-?)([aA-fF0-9]{12})(}?)$`

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
	NamespaceDNS, _  = Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = Parse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = Parse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = Parse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

var re = regexp.MustCompile(`({?|urn:uuid:)([aA-fF0-9]{8})(-?)([aA-fF0-9]{4})(-?)([1-5][aA-fF0-9]{3})(-?)([aA-fF0-9]{4})(-?)([aA-fF0-9]{12})(}?)$`)

// const hexPattern = "^(urn\\:uuid\\:)?\\{?([a-z0-9]{8})-([a-z0-9]{4})-" +
// 	"([1-5][a-z0-9]{3})-([a-z0-9]{4})-([a-z0-9]{12})\\}?$"

func (u UUID) Valid() bool {
	if len(u) != 16 {
		return false
	}
	return true
}

func Valid(uuid string) bool {
	s := re.FindString(uuid)
	if s == "" {
		return false
	}
	return true
}

// Format turn string any uuid format into canonical string
func Format(uuid string) string {
	if !Valid(uuid) {
		return ""
	}
	rplc := []string{"urn", ":", "uuid", "{", "}"}
	for _, r := range rplc {
		uuid = strings.ReplaceAll(uuid, r, "")
	}
	if (uuid[8] != '-' || uuid[14] != '-' || uuid[18] != '-' || uuid[23] != '-') && len(uuid) == 32 {
		uuid = uuid[:8] + "-" + uuid[8:12] + "-" + uuid[12:16] + "-" + uuid[16:20] + "-" + uuid[20:]
	}
	return uuid
}

// Parse return a valid UUID from string. It returns UUID{} (Nil) if not able to parse
// func Parse(uuid string) UUID {
// 	uuid = Format(uuid)
// 	md := re.FindStringSubmatch(uuid)
// 	if md == nil {
// 		return Nil
// 	}
// 	hash := md[2] + md[3] + md[4] + md[5] + md[6]
// 	b, err := hex.DecodeString(hash)
// 	if err != nil {
// 		return Nil
// 	}
// 	var u UUID
// 	copy(u[:], b)
// 	return u
// }

func Parse(str string) (uuid UUID, err error) {
	// if len(str) != 36 || str[8:9] != "-" || str[13:14] != "-" || str[18:19] != "-" || str[23:24] != "-" {
	// 	return uuid, errors.New("Invalid uuid")
	// }
	str = Format(str)
	b, err := hex.DecodeString(str[0:8] + str[9:13] + str[14:18] + str[19:23] + str[24:])
	if err != nil {
		return uuid, errors.New("error decoding uuid")
	}
	_, err = io.ReadFull(bytes.NewBuffer(b), uuid[:])
	return uuid, err
}

func (u UUID) setHash(hash hash.Hash, ns, name []byte) {
	hash.Write(ns[:])
	hash.Write(name)
	copy(u[:], hash.Sum([]byte{})[:16])
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

func New() UUID {
	return NewV4()
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

func (u UUID) setVariant(v byte) {
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
func (u UUID) Equal(ub UUID) bool {
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

// func (u *UUID) String() string {
// 	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
// }

func (u UUID) String() string {
	b := make([]byte, 36)
	b[8], b[13], b[18], b[23] = '-', '-', '-', '-'
	hex.Encode(b[0:8], u[0:4])
	hex.Encode(b[9:13], u[4:6])
	hex.Encode(b[14:18], u[6:8])
	hex.Encode(b[19:23], u[8:10])
	hex.Encode(b[24:], u[10:])
	return string(b)
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

func (u UUID) URN() string {
	if !u.Valid() {
		return ""
	}
	var tb [45]byte
	copy(tb[:], "urn:uuid:")
	encodeHex(tb[9:], u)
	return string(tb[:])
}

func encodeHex(b []byte, u UUID) {
	hex.Encode(b[:], u[:4])
	b[8] = '-'
	hex.Encode(b[9:13], u[4:6])
	b[13] = '-'
	hex.Encode(b[14:18], u[6:8])
	b[18] = '-'
	hex.Encode(b[19:23], u[8:10])
	b[23] = '-'
	hex.Encode(b[24:], u[10:])
}
