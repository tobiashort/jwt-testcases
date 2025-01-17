package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type JWT struct {
	Header    map[string]any
	Payload   map[string]any
	Signature string
}

func ParseEncodedJWT(encodedJWT string) (JWT, error) {
	jwt := JWT{}
	parts := strings.Split(encodedJWT, ".")
	if len(parts) != 3 {
		return jwt, fmt.Errorf("Expected 3 parts delimited by a '.', but got %d", len(parts))
	}
	jwt.Header = make(map[string]any)
	jwt.Payload = make(map[string]any)
	jwt.Signature = parts[2]
	header, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return jwt, err
	}
	err = json.Unmarshal(header, &jwt.Header)
	if err != nil {
		return jwt, err
	}
	payload, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return jwt, err
	}
	err = json.Unmarshal(payload, &jwt.Payload)
	if err != nil {
		return jwt, err
	}
	return jwt, nil
}
