/*
 * @file: sql.go
 * @author: Jorge Quitério
 * @copyright (c) 2021 Jorge Quitério
 * @license: MIT
 */

package uuid

import (
	"database/sql/driver"
	"errors"
)

func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case UUID:
		*u = src
		return nil
	case []byte:
		if len(src) != 16 {
			u.Scan(string(src))
		}
		copy((*u)[:], src)
		return nil
	case string:
		*u, _ = Parse(src)
		return nil
	}
	return errors.New("uuid: errors while converting text to UUID")
}

func (nu NullUUID) Value() (d driver.Value, err error) {
	if !nu.Valid {
		return
	}
	return nu.UUID.Value()
}

func (nu *NullUUID) Scan(src interface{}) error {
	if src == nil {
		nu.UUID, nu.Valid = Nil, false
		return nil
	}
	nu.Valid = true
	return nu.UUID.Scan(src)
}
