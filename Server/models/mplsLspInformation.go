/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package models

import (
	"encoding/xml"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/log"
)

type MplsLspInformation struct {
	XMLName  xml.Name  `xml:"mpls-lsp-information"`
	MplsLsps []MplsLsp `xml:"rsvp-session-data>rsvp-session>mpls-lsp"`
}

type MplsLsp struct {
	LspState           string      `xml:"lsp-state"`
	Name               string      `xml:"name"`
	DestinationAddress string      `xml:"destination-address"`
	MplsP2mpName       string      `xml:"mpls-p2mp-name"`
	MplsLspPath        MplsLspPath `xml:"mpls-lsp-path"`
}

type MplsLspPath struct {
	ReceivedRro string `xml:"received-rro"`
	Bandwidth   string `xml:"bandwidth"`
}

type MplsLspInfo struct {
	Name        string
	ReceivedRro []string
	Bandwidth   string
}

func (data MplsLsp) IsP2MP() bool {
	return data.MplsP2mpName != ""
}

func (data MplsLspInformation) FilterMplsLspInfoByEgress(egressAddress string) []MplsLsp {
	result := make([]MplsLsp, 0)

	for _, mplsLsp := range data.MplsLsps {
		if mplsLsp.DestinationAddress == egressAddress {
			result = append(result, mplsLsp)
		}
	}

	return result
}

func ParseMplsLspInformation(xmlText []byte) MplsLspInformation {
	result := MplsLspInformation{}

	err := xml.Unmarshal(xmlText, &result)
	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}

	return result
}
