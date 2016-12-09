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

	"github.com/Juniper/24287_WOW_LSP_GOLANG/log"
)

type InterfaceInformation struct {
	XMLName           xml.Name           `xml:"interface-information"`
	Name              string             `xml:"physical-interface>name"`
	Bandwidth         string             `xml:"physical-interface>speed"`
	LastFlapped       string             `xml:"physical-interface>interface-flapped"`
	StatsLastCleared  string             `xml:"physical-interface>statistics-cleared"`
	TrafficStatistics TrafficStatistics  `xml:"physical-interface>traffic-statistics"`
	InputErrorList    InputErrorList     `xml:"physical-interface>input-error-list"`
	OutputErrorList   OutputErrorList    `xml:"physical-interface>output-error-list"`
	QueueCounters     []QueueCounterInfo `xml:"physical-interface>queue-counters>queue"`
	PcsStatistics     PcsStatistics      `xml:"physical-interface>ethernet-pcs-statistics"`
}

type TrafficStatistics struct {
	InputBytes    string `xml:"input-bytes"`
	OutputBytes   string `xml:"output-bytes"`
	InputPackets  string `xml:"input-packets"`
	OutputPackets string `xml:"output-packets"`
}

type InputErrorList struct {
	Errors         string `xml:"input-errors"`
	Drops          string `xml:"input-drops"`
	FramingsErrors string `xml:"framing-errors"`
	Runts          string `xml:"input-runts"`
	Discards       string `xml:"input-discards"`
	L3Incompletes  string `xml:"input-l3-incompletes"`
}

type OutputErrorList struct {
	Errors             string `xml:"output-errors"`
	Drops              string `xml:"output-drops"`
	CarrierTransitions string `xml:"carrier-transitions"`
	Collisions         string `xml:"output-collisions"`
	CRCErrors          string `xml:"hs-link-crc-errors"`
	MTUErrors          string `xml:"mtu-errors"`
}

type PcsStatistics struct {
	BitErrors     string `xml:"bit-error-seconds"`
	ErroredBlocks string `xml:"errored-blocks-seconds"`
}

type QueueCounterInfo struct {
	Name             string `xml:"forwarding-class-name"`
	TotalDropPackets string `xml:"queue-counters-total-drop-packets"`
}

func (obj InterfaceInformation) ToRouterStatisticsContent(router *Router) RouterStatisticsContent {
	var statistics = []*Statistics{}

	var rows = []*StatisticsValue{}
	AddValueRow(&rows, obj.Bandwidth, "")
	statistics = append(statistics, &Statistics{Title: "Bandwidth", Values: rows})

	rows = []*StatisticsValue{}
	AddLastFlappedErrorRow(&rows, obj.LastFlapped)
	statistics = append(statistics, &Statistics{Title: "Last Flapped", Values: rows})

	rows = []*StatisticsValue{}
	AddValueRow(&rows, obj.StatsLastCleared, "")
	statistics = append(statistics, &Statistics{Title: "Stats Last Cleared", Values: rows})

	rows = []*StatisticsValue{}
	AddValueRow(&rows, "Input Bytes", obj.TrafficStatistics.InputBytes)
	AddValueRow(&rows, "Ouput Bytes", obj.TrafficStatistics.OutputBytes)
	AddValueRow(&rows, "Input packets", obj.TrafficStatistics.InputPackets)
	AddValueRow(&rows, "Ouput packets", obj.TrafficStatistics.OutputPackets)
	statistics = append(statistics, &Statistics{Title: "Traffic Statistics", Values: rows})

	rows = []*StatisticsValue{}
	AddErrorRow(&rows, "Errors", obj.InputErrorList.Errors)
	AddErrorRow(&rows, "Drops", obj.InputErrorList.Drops)
	AddErrorRow(&rows, "Framing Errors", obj.InputErrorList.FramingsErrors)
	AddErrorRow(&rows, "Runts", obj.InputErrorList.Runts)
	AddErrorRow(&rows, "Policied Discards", obj.InputErrorList.Discards)
	AddErrorRow(&rows, "L3 Incompletes", obj.InputErrorList.L3Incompletes)
	statistics = append(statistics, &Statistics{Title: "Input Errors", Values: rows})

	rows = []*StatisticsValue{}
	AddErrorRow(&rows, "Errors", obj.OutputErrorList.Errors)
	AddErrorRow(&rows, "Drops", obj.OutputErrorList.Drops)
	AddErrorRow(&rows, "Carrier Transitions", obj.OutputErrorList.CarrierTransitions)
	AddErrorRow(&rows, "Collisions", obj.OutputErrorList.Collisions)
	AddErrorRow(&rows, "CRC Errors", obj.OutputErrorList.CRCErrors)
	AddErrorRow(&rows, "MTU Errors", obj.OutputErrorList.MTUErrors)
	statistics = append(statistics, &Statistics{Title: "Output Errors", Values: rows})

	rows = []*StatisticsValue{}
	obj.addQueueCounterValue(&rows, "Best-Effort")
	obj.addQueueCounterValue(&rows, "Video-Backup")
	obj.addQueueCounterValue(&rows, "Business-Low")
	obj.addQueueCounterValue(&rows, "Business-Medium")
	obj.addQueueCounterValue(&rows, "Business-High")
	obj.addQueueCounterValue(&rows, "Video-Primary")
	obj.addQueueCounterValue(&rows, "VoIP")
	obj.addQueueCounterValue(&rows, "Network-Control")
	statistics = append(statistics, &Statistics{Title: "Ingress/Egress Queue Counters", Values: rows})

	rows = []*StatisticsValue{}
	AddErrorRow(&rows, "Bit Errors", obj.PcsStatistics.BitErrors)
	AddErrorRow(&rows, "Errored blocks", obj.PcsStatistics.ErroredBlocks)
	statistics = append(statistics, &Statistics{Title: "PCS statistics", Values: rows})

	rows = []*StatisticsValue{}
	AddErrorRow(&rows, "RX", "0")
	AddErrorRow(&rows, "TX", "0")
	statistics = append(statistics, &Statistics{Title: "CRC/Align errors", Values: rows})

	rows = []*StatisticsValue{}
	AddValueRow(&rows, "Input", "None")
	AddValueRow(&rows, "Output", "None")
	statistics = append(statistics, &Statistics{Title: "Configured Policer", Values: rows})

	return RouterStatisticsContent{
		Id:            strings.TrimSpace(obj.GetName()),
		Name:          obj.Name,
		Statistics:    statistics,
		SubInterfaces: []*RouterStatisticsContent{}}
}

func (obj InterfaceInformation) GetName() string {
	return obj.Name
}

func (obj InterfaceInformation) GetTrafficStatistics() TrafficStatistics {
	return obj.TrafficStatistics
}

func (obj InterfaceInformation) addQueueCounterValue(rows *[]*StatisticsValue, header string) {
	for _, queueCounter := range obj.QueueCounters {
		if strings.TrimSpace(queueCounter.Name) == header {
			AddErrorRow(rows, header, queueCounter.TotalDropPackets)
			return
		}
	}

	AddValueRow(rows, header, "None")
}

func ParseInterfaceInformation(xmlText []byte) InterfaceInformation {
	result := InterfaceInformation{}

	err := xml.Unmarshal(xmlText, &result)
	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}

	return result
}
