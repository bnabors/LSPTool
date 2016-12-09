/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controllerHelper

import (
	"strconv"
	"strings"

	"../../models"
)

type HostTestResult struct {
	Router       models.Router
	PfeStatistic models.PfeStatistic
	RouterType   string
}

type LinkDirectTestResult struct {
	InterfaceInfo        models.IRouterStatistics
	LogicalInterfaceName string
	Destination          string
}

type LinkTestResult struct {
	RtDestination string
	Forward       LinkDirectTestResult
	Backward      LinkDirectTestResult
}

type GroupTestResult struct {
	GroupId       string
	IngressIp     string
	DiagramResult models.DiagramResult
	IcmpResult    models.TestResult
	LspBands      []*models.LspBand
	Hosts         []*HostTestResult
	Links         []*LinkTestResult
}

func (gtr *GroupTestResult) BuildDiagrammResult() {
	gtr.DiagramResult = models.DiagramResult{Paths: []*models.DiaPath{}}

	for index := 1; index < len(gtr.Hosts); index++ {

		var hostStart = *(gtr.Hosts[index-1])
		var hostFinish = *(gtr.Hosts[index])
		var link = *(gtr.Links[index-1])

		var baseIp, destinationBack, destination, _ = getAdresses(link.RtDestination, link.Backward.Destination, link.Forward.Destination)

		var backwardStatistics = link.Backward.InterfaceInfo.GetTrafficStatistics()
		var forwardStatistics = link.Forward.InterfaceInfo.GetTrafficStatistics()

		router1 := models.DiaRouter{
			Id:            hostStart.Router.Name,
			Name:          hostStart.Router.Name,
			Ip:            destinationBack,
			Interface:     link.Backward.LogicalInterfaceName,
			InputBytes:    backwardStatistics.InputBytes,
			OutputBytes:   backwardStatistics.OutputBytes,
			InputPackets:  backwardStatistics.InputPackets,
			OutputPackets: backwardStatistics.OutputPackets,
		}

		router2 := models.DiaRouter{
			Id:            hostFinish.Router.Name,
			Name:          hostFinish.Router.Name,
			Ip:            destination,
			Interface:     link.Forward.LogicalInterfaceName,
			InputBytes:    forwardStatistics.InputBytes,
			OutputBytes:   forwardStatistics.OutputBytes,
			InputPackets:  forwardStatistics.InputPackets,
			OutputPackets: forwardStatistics.OutputPackets,
		}

		gtr.DiagramResult.Paths = append(gtr.DiagramResult.Paths, &models.DiaPath{Router1: router1, Router2: router2, BaseIp: baseIp})
	}
}

func (gtr *GroupTestResult) BuildIcmpResult() (err error) {
	var routers = []*models.Router{}

	for _, host := range gtr.Hosts {
		routers = append(routers, &host.Router)
	}

	gtr.IcmpResult, err = BuildIcmpResultByRouters("icmp_"+gtr.GroupId, routers)
	return
}

func getAdresses(rtDestination string, addressLeft string, addressRight string) (baseIp string, destinationBack string, destination string, err error) {

	baseIp = rtDestination
	destinationBack = addressLeft
	destination = addressRight

	var address = rtDestination[:strings.LastIndex(rtDestination, "/")]
	var maskStr = rtDestination[strings.LastIndex(rtDestination, "/")+1:]

	var octets = strings.Split(address, ".")
	var octetsLeft = strings.Split(addressLeft, ".")
	var octetsRight = strings.Split(addressRight, ".")

	mask, err := strconv.Atoi(maskStr)
	if err != nil {
		return
	}

	var maskOctets []string

	var div = mask / 8
	var mod = mask % 8

	for pos := 0; pos < 4; pos++ {
		if pos < div {
			maskOctets = append(maskOctets, octets[pos])
			continue
		}

		if pos == div {
			if mod == 0 {
				break
			}

			oct, errConv := strconv.Atoi(octets[pos])
			if errConv != nil {
				err = errConv
				return
			}

			var adressBits = uint(8 - mod)
			var maskOct = (255 >> adressBits) << adressBits
			maskOctets = append(maskOctets, strconv.Itoa(oct&maskOct))
			break
		}
	}

	baseIp = strings.Join(maskOctets, ".") + "/" + maskStr
	destinationBack = strings.Join(octetsLeft[div:], ".")
	destination = strings.Join(octetsRight[div:], ".")

	if div > 0 {
		destinationBack = "." + destinationBack
		destination = "." + destination
	}

	return
}
