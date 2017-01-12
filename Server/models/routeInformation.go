/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package models

import (
	"encoding/xml"

	"github.com/WOWLABS/LSPTool/Server/log"
)

type RouteInformation struct {
	XMLName xml.Name `xml:"route-information"`
	Rt      Rt       `xml:"route-table>rt"`
}

type Rt struct {
	RtDestination string    `xml:"rt-destination"`
	RtEntries     []RtEntry `xml:"rt-entry"`
}

type RtEntry struct {
	ProtocolName  string      `xml:"protocol-name"`
	Age           string      `xml:"age"`
	Nhs           []Nh        `xml:"nh"`
	CurrentActive *ExistField `xml:"current-active"`
}

type Nh struct {
	SelectedNextHop *ExistField `xml:"selected-next-hop"`
	Via             string      `xml:"via"`
}

type ExistField struct {
}

func (data RouteInformation) GetInterfaceName() string {
	for _, rtEntry := range data.Rt.RtEntries {
		if rtEntry.CurrentActive == nil {
			continue
		}

		for _, nh := range rtEntry.Nhs {
			if nh.SelectedNextHop == nil {
				continue
			}

			return nh.Via
		}
	}

	return ""
}

func ParseRouteInformation(xmlText []byte) RouteInformation {
	result := RouteInformation{}

	err := xml.Unmarshal(xmlText, &result)
	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}

	return result
}
