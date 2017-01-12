/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package models

import (
	"encoding/xml"
	"strings"

	"github.com/WOWLABS/LSPTool/Server/log"
)

type OspfNeighborInformation struct {
	XMLName       xml.Name       `xml:"ospf-neighbor-information"`
	OspfNeighbors []OspfNeighbor `xml:"ospf-neighbor"`
}

type OspfNeighbor struct {
	NeighborAddress string `xml:"neighbor-address"`
	InterfaceName   string `xml:"interface-name"`
	NeighborId      string `xml:"neighbor-id"`
}

func ParseOspfNeighborInformation(xmlText []byte) OspfNeighborInformation {
	result := OspfNeighborInformation{}

	err := xml.Unmarshal(xmlText, &result)
	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}

	return result
}

func (data OspfNeighborInformation) GetOspfNeighbor(interfaceName string) *OspfNeighbor {
	var pattern = strings.ToLower(interfaceName)

	for _, info := range data.OspfNeighbors {
		if strings.ToLower(info.InterfaceName) == pattern {
			return &info
		}
	}

	return nil
}
