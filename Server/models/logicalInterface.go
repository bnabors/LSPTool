package models

import (
	"strings"
)

type LogicalInterface struct {
	AddressFamilies      []AddressFamily `xml:"address-family"`
	LogicalInterfaceName string          `xml:"name"`
}

type AddressFamily struct {
	AddressFamilyName string `xml:"address-family-name"`
	LocalIp           string `xml:"interface-address>ifa-local"`
}

func GetLocalIp(logicalInterfaces []LogicalInterface, interfaceName string) string {
	for _, logicalInterface := range logicalInterfaces {
		if strings.TrimSpace(logicalInterface.LogicalInterfaceName) != strings.TrimSpace(interfaceName) {
			continue
		}

		for _, addresFamily := range logicalInterface.AddressFamilies {
			if strings.TrimSpace(addresFamily.AddressFamilyName) != "inet" {
				continue
			}

			return strings.TrimSpace(addresFamily.LocalIp)
		}
	}

	return ""
}
