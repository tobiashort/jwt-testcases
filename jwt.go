package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tobiashort/orderedmap"
)

type JWT struct {
	Header    *orderedmap.OrderedMap[string, any]
	Payload   *orderedmap.OrderedMap[string, any]
	Signature string
}

func DecodedJWT(encodedJWT string) (JWT, error) {
	jwt := JWT{}
	parts := strings.Split(encodedJWT, ".")
	if len(parts) != 3 {
		return jwt, fmt.Errorf("Expected 3 parts delimited by a '.', but got %d", len(parts))
	}
	jwt.Header = orderedmap.NewOrderedMap[string, any]()
	jwt.Payload = orderedmap.NewOrderedMap[string, any]()
	jwt.Signature = parts[2]
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return jwt, err
	}
	err = json.Unmarshal(header, &jwt.Header)
	if err != nil {
		return jwt, err
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return jwt, err
	}
	err = json.Unmarshal(payload, &jwt.Payload)
	if err != nil {
		return jwt, err
	}
	return jwt, nil
}

func (jwt JWT) Encode() string {
	header, err := json.Marshal(jwt.Header)
	AssertNil(err)
	payload, err := json.Marshal(jwt.Payload)
	AssertNil(err)
	b64Header := base64.RawURLEncoding.EncodeToString(header)
	b64Payload := base64.RawURLEncoding.EncodeToString(payload)
	return fmt.Sprintf("%s.%s.%s", b64Header, b64Payload, jwt.Signature)
}

func (jwt JWT) String() string {
	header, err := json.MarshalIndent(jwt.Header, "", "  ")
	AssertNil(err)
	payload, err := json.MarshalIndent(jwt.Payload, "", " ")
	AssertNil(err)
	return fmt.Sprintf("%s\n\n%s\n\n%s", header, payload, jwt.Signature)
}
