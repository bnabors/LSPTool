/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package models

import (
	"encoding/xml"

	"../log"
)

type RouteInformation struct {
	XMLName xml.Name `xml:"route-information"`
	Rt      Rt       `xml:"route-table>rt"`
}

type Rt struct {
	RtDestination string  `xml:"rt-destination"`
	RtEntry       RtEntry `xml:"rt-entry"`
}

type RtEntry struct {
	ProtocolName string `xml:"protocol-name"`
	Age          string `xml:"age"`
	Nh           Nh     `xml:"nh"`
}

type Nh struct {
	Via string `xml:"via"`
}

func (data RouteInformation) GetInterfaceName() string {
	return data.Rt.RtEntry.Nh.Via
}

func ParseRouteInformation(xmlText []byte) RouteInformation {
	result := RouteInformation{}

	err := xml.Unmarshal(xmlText, &result)
	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}

	return result
}
