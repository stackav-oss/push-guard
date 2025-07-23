package config

import (
	"encoding/base64"
	"fmt"
)

func DecodeConfigString(variable string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(variable)
	if err != nil {
		fmt.Printf("Error decoding Base64: %v\n", err)
		return ""
	}
	return string(decodedBytes)
}

// The following variables will be set by ldflags
var PushGuardVersion string
var Disclaimer string
var LogCollectorURL string
var ProtectedBranches string
var ProtocolAndDomainAllowList string
var DirectoryAllowList string
var DirectoryRegexAllowList string
