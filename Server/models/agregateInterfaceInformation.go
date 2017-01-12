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

	"github.com/WOWLABS/LSPTool/Server/helpers"
	"github.com/WOWLABS/LSPTool/Server/log"
)

type (
	AgregateInterfaceInformation struct {
		XMLName                 xml.Name `xml:"interface-information"`
		Bundle                  string   `xml:"physical-interface>name"`
		Bandwidth               string   `xml:"physical-interface>speed"`
		LastFlapped             string   `xml:"physical-interface>interface-flapped"`
		StatsLastCleared        string   `xml:"physical-interface>statistics-cleared"`
		ConfiguredPolicerInput  string   `xml:"logical-interface>policer-input"`
		ConfiguredPolicerOutput string   `xml:"logical-interface>policer-output"`

		TrafficStatistics TrafficStatistics `xml:"physical-interface>traffic-statistics"`

		InputErrors  InputError  `xml:"physical-interface>input-error-list"`
		OutputErrors OutputError `xml:"physical-interface>output-error-list"`

		QueueCounters []QueueCounter `xml:"physical-interface>queue-counters"`

		LogicalInterfaces []AggregateLogicalnterface `xml:"physical-interface>logical-interface"`

		SubInterface []InterfaceInformation
	}

	InputError struct {
		Errors           string `xml:"input-errors"`
		Drops            string `xml:"input-drops"`
		FramingErrors    string `xml:"framing-errors"`
		Runts            string `xml:"input-runts"`
		PoliciedDiscards string `xml:"input-discards"`
	}

	OutputError struct {
		Errors             string `xml:"output-errors"`
		Drops              string `xml:"output-drops"`
		CarrierTransitions string `xml:"carrier-transitions"`
		MTUErrors          string `xml:"mtu-errors"`
	}

	QueueCounter struct {
		Name   string                 `xml:"interface-cos-short-summary>intf-cos-queue-type"`
		Values []QueueCounterKeyValue `xml:"queue"`
	}

	QueueCounterKeyValue struct {
		Name  string `xml:"forwarding-class-name"`
		Value string `xml:"queue-counters-total-drop-packets"`
	}

	AggregateLogicalnterface struct {
		AddressFamilies      []AggregateAddressFamily `xml:"address-family"`
		LogicalInterfaceName string                   `xml:"name"`
		SubInterfaceNames    []string                 `xml:"lag-traffic-statistics>lag-link>name"`
	}

	AggregateAddressFamily struct {
		AddressFamilyName string `xml:"address-family-name"`
		LocalIp           string `xml:"interface-address>ifa-local"`
	}
)

func ParseAgregateInterfaceInformation(xmlText []byte) AgregateInterfaceInformation {
	result := AgregateInterfaceInformation{
		Bundle:                  "None",
		Bandwidth:               "None",
		LastFlapped:             "None",
		StatsLastCleared:        "None",
		ConfiguredPolicerInput:  "None",
		ConfiguredPolicerOutput: "None",

		TrafficStatistics: TrafficStatistics{
			InputBytes:    "None",
			InputPackets:  "None",
			OutputBytes:   "None",
			OutputPackets: "None",

			InputBytesPerSecond:    "None",
			OutputBytesPerSecond:   "None",
			InputPacketsPerSecond:  "None",
			OutputPacketsPerSecond: "None",
		},
		SubInterface: []InterfaceInformation{},
		InputErrors: InputError{
			Errors:           "None",
			Drops:            "None",
			FramingErrors:    "None",
			Runts:            "None",
			PoliciedDiscards: "None",
		},
		OutputErrors: OutputError{
			Errors:             "None",
			Drops:              "None",
			CarrierTransitions: "None",
			MTUErrors:          "None",
		},
		QueueCounters: []QueueCounter{},
	}

	/*result.IngressQueueCounter.Values = append(result.IngressQueueCounter.Values, QueueCounterKeyValue{Name: "Best-Effort", Value: "None"})
	result.IngressQueueCounter.Values = append(result.IngressQueueCounter.Values, QueueCounterKeyValue{Name: "Video-Backup", Value: "None"})
	result.IngressQueueCounter.Values = append(result.IngressQueueCounter.Values, QueueCounterKeyValue{Name: "Business-Low", Value: "None"})
	result.IngressQueueCounter.Values = append(result.IngressQueueCounter.Values, QueueCounterKeyValue{Name: "Business-Med", Value: "None"})
	result.IngressQueueCounter.Values = append(result.IngressQueueCounter.Values, QueueCounterKeyValue{Name: "Business-Hig", Value: "None"})
	result.IngressQueueCounter.Values = append(result.IngressQueueCounter.Values, QueueCounterKeyValue{Name: "Video-Primar", Value: "None"})
	result.IngressQueueCounter.Values = append(result.IngressQueueCounter.Values, QueueCounterKeyValue{Name: "VoIP", Value: "None"})
	result.IngressQueueCounter.Values = append(result.IngressQueueCounter.Values, QueueCounterKeyValue{Name: "Network-Contr", Value: "None"})

	result.EgressQueueCounter.Values = append(result.EgressQueueCounter.Values, QueueCounterKeyValue{Name: "Best-Effort", Value: "None"})
	result.EgressQueueCounter.Values = append(result.EgressQueueCounter.Values, QueueCounterKeyValue{Name: "Video-Backup", Value: "None"})
	result.EgressQueueCounter.Values = append(result.EgressQueueCounter.Values, QueueCounterKeyValue{Name: "Business-Low", Value: "None"})
	result.EgressQueueCounter.Values = append(result.EgressQueueCounter.Values, QueueCounterKeyValue{Name: "Business-Med", Value: "None"})
	result.EgressQueueCounter.Values = append(result.EgressQueueCounter.Values, QueueCounterKeyValue{Name: "Business-Hig", Value: "None"})
	result.EgressQueueCounter.Values = append(result.EgressQueueCounter.Values, QueueCounterKeyValue{Name: "Video-Primar", Value: "None"})
	result.EgressQueueCounter.Values = append(result.EgressQueueCounter.Values, QueueCounterKeyValue{Name: "VoIP", Value: "None"})
	result.EgressQueueCounter.Values = append(result.EgressQueueCounter.Values, QueueCounterKeyValue{Name: "Network-Contr", Value: "None"})*/

	err := xml.Unmarshal(xmlText, &result)
	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}

	return result
}

func (obj AgregateInterfaceInformation) ToRouterStatisticsContent() RouterStatisticsContent {
	var statistics = []*Statistics{}

	var rows = []*StatisticsValue{}
	AddValueRow(&rows, obj.Bundle, "")
	statistics = append(statistics, &Statistics{Title: "Bundle", Values: rows})

	rows = []*StatisticsValue{}
	AddValueRow(&rows, obj.Bandwidth, "")
	statistics = append(statistics, &Statistics{Title: "Bandwidth", Values: rows})

	rows = []*StatisticsValue{}
	AddLastFlappedErrorRow(&rows, obj.LastFlapped)
	statistics = append(statistics, &Statistics{Title: "Last Flapped", Values: rows})

	rows = []*StatisticsValue{}
	AddValueRow(&rows, obj.StatsLastCleared, "")
	statistics = append(statistics, &Statistics{Title: "Stats Last Cleared", Values: rows})

	rows = []*StatisticsValue{}
	AddValueRow(&rows, "BPS In", helpers.ParceNumberAndLocalize(obj.TrafficStatistics.InputBytesPerSecond))
	AddValueRow(&rows, "BPS Out", helpers.ParceNumberAndLocalize(obj.TrafficStatistics.OutputBytesPerSecond))
	AddValueRow(&rows, "PPS In", helpers.ParceNumberAndLocalize(obj.TrafficStatistics.InputPacketsPerSecond))
	AddValueRow(&rows, "PPS Out", helpers.ParceNumberAndLocalize(obj.TrafficStatistics.OutputPacketsPerSecond))
	statistics = append(statistics, &Statistics{Title: "Traffic Statistics", Values: rows})

	rows = []*StatisticsValue{}
	AddErrorRow(&rows, "Errors", obj.InputErrors.Errors)
	AddErrorRow(&rows, "Drops", obj.InputErrors.Drops)
	AddErrorRow(&rows, "Framing Errors", obj.InputErrors.FramingErrors)
	AddErrorRow(&rows, "Runts", obj.InputErrors.Runts)
	AddErrorRow(&rows, "Policied Discards", obj.InputErrors.PoliciedDiscards)
	statistics = append(statistics, &Statistics{Title: "Input Errors", Values: rows})

	rows = []*StatisticsValue{}
	AddErrorRow(&rows, "Errors", obj.OutputErrors.Errors)
	AddErrorRow(&rows, "Drops", obj.OutputErrors.Drops)
	AddErrorRow(&rows, "Carrier Transitions", obj.OutputErrors.CarrierTransitions)
	AddErrorRow(&rows, "MTU Errors", obj.OutputErrors.MTUErrors)
	statistics = append(statistics, &Statistics{Title: "Output Errors", Values: rows})

	rows = []*StatisticsValue{}
	AddValueRow(&rows, "Input", obj.ConfiguredPolicerInput)
	AddValueRow(&rows, "Output", obj.ConfiguredPolicerOutput)
	statistics = append(statistics, &Statistics{Title: "Configured Policer", Values: rows})

	rows = []*StatisticsValue{}
	for _, value := range obj.QueueCounters {
		AddValueRow(&rows, value.Name+" counters", "Dropped Packets")
		for _, queueValue := range value.Values {
			AddErrorRow(&rows, queueValue.Name, queueValue.Value)
		}
		AddValueRow(&rows, "", "")
	}
	statistics = append(statistics, &Statistics{Title: "QueueCounter", Values: rows})

	subInterfaces := []*RouterStatisticsContent{}
	for _, subInterface := range obj.SubInterface {
		var subInterfaceStatistic = subInterface.ToRouterStatisticsContent()
		subInterfaces = append(subInterfaces, &subInterfaceStatistic)
	}
	result := RouterStatisticsContent{
		Id:            strings.TrimSpace(obj.GetName()),
		Name:          obj.GetName(),
		Statistics:    statistics,
		SubInterfaces: subInterfaces,
	}

	return result
}

func (obj AgregateInterfaceInformation) GetName() string {
	return obj.Bundle
}

func (obj AgregateInterfaceInformation) GetTrafficStatistics() TrafficStatistics {
	return obj.TrafficStatistics
}

func (obj AgregateInterfaceInformation) GetLocalIp(logicalInterfaceName string) string {
	for _, logicalInterface := range obj.LogicalInterfaces {
		if strings.TrimSpace(logicalInterface.LogicalInterfaceName) != strings.TrimSpace(logicalInterfaceName) {
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

func (obj AgregateInterfaceInformation) GetSubInterfaceNames() []string {
	result := []string{}

	for _, logicalInterface := range obj.LogicalInterfaces {
		if logicalInterface.SubInterfaceNames == nil {
			continue
		}

		for _, subInterfaceName := range logicalInterface.SubInterfaceNames {
			result = append(result, subInterfaceName)
		}
	}

	return result
}
