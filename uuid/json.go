/*
 * @file: json.go
 * @author: Jorge Quitério
 * @copyright (c) 2021 Jorge Quitério
 * @license: MIT
 */

package uuid

import (
	"bytes"
	"encoding/json"
	"errors"
)

func (u *UUID) UnmarshalBinary(b []byte) (err error) {
	if len(b) != 16 {
		err = errors.New("uuid must be 16 bites long")
	}
	copy(u[:], b)
	return
}

func (u UUID) MarshalBinary() ([]byte, error) {
	return u.Byte(), nil
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UUID) UnmarshalJSON(b []byte) (err error) {
	err = json.Unmarshal(b, u)
	return
}

func (nu NullUUID) MarshalJSON() ([]byte, error) {
	if !nu.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nu.UUID)
}

func (nu *NullUUID) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte("null")) {
		nu.UUID, nu.Valid = Nil, false
		return nil
	}
	err = json.Unmarshal(b, &nu.UUID)
	if err != nil {
		return
	}
	nu.Valid = true
	return
}
