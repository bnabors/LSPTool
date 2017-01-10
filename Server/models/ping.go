/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package models

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/config"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
)

// resultStr xml model of ping result.
type (
	pingResultXML struct {
		XMLName           xml.Name       `xml:"ping-results"`
		Probes            []pingProbeXML `xml:"probe-result"`
		ProbesSent        string         `xml:"probe-results-summary>probes-sent"`
		ResponsesReceived string         `xml:"probe-results-summary>responses-received"`
		PacketLoss        string         `xml:"probe-results-summary>packet-loss"`
		RttMinimum        string         `xml:"probe-results-summary>rtt-minimum"`
		RttMaximum        string         `xml:"probe-results-summary>rtt-maximum"`
		RttAverage        string         `xml:"probe-results-summary>rtt-average"`
		RttStddev         string         `xml:"probe-results-summary>rtt-stddev"`
	}
	pingProbeXML struct {
		ProbeIndex     string      `xml:"probe-index"`
		ProbeResult    interface{} `xml:"probe-success"`
		SequenceNumber string      `xml:"sequence-number"`
		IPAddress      string      `xml:"ip-address"`
		TimeToLive     string      `xml:"time-to-live"`
		ResponseSize   string      `xml:"response-size"`
		Rtt            string      `xml:"rtt"`
	}
)

func (v pingResultXML) toPingResult() PingResult {
	res := PingResult{}
	res.ProbesSent = toIntAndLog(v.ProbesSent)
	res.ResponsesReceived = toIntAndLog(v.ResponsesReceived)
	res.PacketLoss = toIntAndLog(v.PacketLoss)
	res.RttMinimum = toIntAndLog(v.RttMinimum)
	res.RttMaximum = toIntAndLog(v.RttMaximum)
	res.RttAverage = toIntAndLog(v.RttAverage)
	res.RttStddev = toIntAndLog(v.RttStddev)

	res.IsError = len(v.Probes) <= 0 || toInt(v.Probes[0].SequenceNumber) > 1
	if res.IsError {
		return res
	}

	for i := 0; i < len(v.Probes)-1; i++ {
		if toInt(v.Probes[i+1].SequenceNumber)-toInt(v.Probes[i].SequenceNumber) <= 2 {
			continue
		}
		res.IsError = true
		break
	}

	return res
}

func toInt(v string) int {
	integer, _ := strconv.Atoi(strings.TrimSpace(v))
	return integer
}
func toIntAndLog(v string) int {
	integer, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		lspLogger.Infof("Parse Error: %v", err)
	}
	return integer
}

//ParsePing make PingResult from xml text
func ParsePing(xmlText []byte) PingResult {
	var result pingResultXML

	err := xml.Unmarshal(xmlText, &result)

	if err != nil {
		lspLogger.Errorf("Error: \n %v", err)
		return PingResult{}
	}

	return result.toPingResult()
}

type (

	// PingResult integer result.
	PingResult struct {
		ProbesSent        int
		ResponsesReceived int
		PacketLoss        int
		RttMinimum        int
		RttMaximum        int
		RttAverage        int
		RttStddev         int
		IsError           bool
	}
)

// Print write result into log
func (res PingResult) Print() {
	lspLogger.Debugf("ProbesSent: %v\n", res.ProbesSent)
	lspLogger.Debugf("ResponsesReceived: %v\n", res.ResponsesReceived)
	lspLogger.Debugf("PacketLoss: %v\n", res.PacketLoss)
	lspLogger.Debugf("RttMinimum: %v\n", res.RttMinimum)
	lspLogger.Debugf("RttMaximum: %v\n", res.RttMaximum)
	lspLogger.Debugf("RttAverage: %v\n", res.RttAverage)
	lspLogger.Debugf("RttStddev: %v\n", res.RttStddev)
}

// ToIcmpResult - make models for page view
func (res PingResult) ToIcmpResult(icmpInfo *IcmpInfo) IcmpResult {
	average := float64(res.RttAverage) / 1000
	isError := res.IsError || res.PacketLoss > config.LspConfig.PingLossPercentThreshold || average > config.LspConfig.PingAvgThreshold

	return IcmpResult{
		Id:      "icmp_" + strconv.Itoa(icmpInfo.Source.Id) + "_" + strconv.Itoa(icmpInfo.Destination.Id),
		Info:    icmpInfo,
		Loss:    strconv.FormatFloat(float64(res.PacketLoss), 'f', -1, 32) + "%",
		Max:     strconv.FormatFloat(float64(res.RttMaximum)/1000, 'f', -1, 32),
		Average: strconv.FormatFloat(average, 'f', -1, 32),
		StdDev:  strconv.FormatFloat(float64(res.RttStddev)/1000, 'f', -1, 32),
		IsError: isError,
	}
}
