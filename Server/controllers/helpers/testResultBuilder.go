/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controllerHelper

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/commands"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/utils"
)

type TestResultBuilder struct {
	AllRouters       []*models.Router
	GroupTestResults []*GroupTestResult
	LspGroups        []*models.LspGroup
}

func (trb *TestResultBuilder) Init(lspGroups []*models.LspGroup) {
	allRouters, _ := utils.GetRoutersFromMysql()
	trb.AllRouters = allRouters

	trb.LspGroups = lspGroups

	trb.GroupTestResults = []*GroupTestResult{}
}

func (trb *TestResultBuilder) TryRunTest(lspItem models.LspItem) error {

	var testResult *GroupTestResult
	var needTests bool
	testResult, needTests = trb.getOrCreateGroupTestResults(lspItem)

	if !needTests {
		return nil
	}

	var routerAddress = testResult.IngressIp
	var currentIp, err = trb.runHostTests(routerAddress, true, false, testResult)
	if err != nil {
		return err
	}

	var lspGroup = trb.getLspGroupById(lspItem.GroupId)

	for index, nextIp := range lspGroup.ReceivedRro {
		// Get forwards route info
		neighborId, rtDestination, forwardTestResult, localIp, err := trb.runLinkDirectTests(currentIp, nextIp, testResult)
		if err != nil {
			return err
		}

		var isEgress = index == len(lspGroup.ReceivedRro)-1
		currentIp, err = trb.runHostTests(neighborId, false, isEgress, testResult)
		if err != nil {
			return err
		}

		// Get backwards route info
		_, _, backwardTestResult, _, err := trb.runLinkDirectTests(currentIp, localIp, testResult)
		if err != nil {
			return err
		}

		linkTestResult := LinkTestResult{RtDestination: rtDestination, Forward: forwardTestResult, Backward: backwardTestResult}
		testResult.Links = append(testResult.Links, &linkTestResult)

		routerAddress = neighborId
	}

	testResult.BuildDiagrammResult()
	testResult.BuildIcmpResult()

	return nil
}

func (trb *TestResultBuilder) AddLspBand(lspItem models.LspItem) {
	for _, res := range trb.GroupTestResults {
		if res.GroupId == lspItem.GroupId {
			res.LspBands = append(res.LspBands, &models.LspBand{Name: lspItem.Name, Bandwidth: lspItem.Bandwidth})
			return
		}
	}
}

func (trb *TestResultBuilder) GetRouteResult(lspItem models.LspItem) *models.RouteResult {
	var res = trb.getGroupTestResults(lspItem)

	// LSP
	result := models.RouteResult{Results: []*models.TestResult{}, Id: lspItem.Id, Names: []string{lspItem.Name}, LspItem: lspItem}

	lspContent := models.LspContent{Bands: []*models.LspBand{}, Routes: []*models.LspRoute{}, Diagram: res.DiagramResult}

	for _, lspBand := range res.LspBands {
		lspContent.Bands = append(lspContent.Bands, lspBand)
	}

	for _, host := range res.Hosts {
		router := host.Router
		lspRouter := models.LspRouter{Id: strconv.Itoa(router.Id), Name: router.Name}
		lspRoute := models.LspRoute{Type: host.RouterType, Router: lspRouter}
		lspContent.Routes = append(lspContent.Routes, &lspRoute)
	}

	result.Results = append(result.Results, &models.TestResult{Id: "lspsres_" + lspItem.Id, Name: "LSP (s)", Type: 1, Content: &lspContent})

	// Routers
	for index, host := range res.Hosts {
		routeTestContent := []models.RouterStatisticsContent{}

		router := host.Router
		routerPfeStatistic := host.PfeStatistic

		var isError = false

		var pfeStatistic = routerPfeStatistic.ToRouterStatisticsContent(&router)
		pfeStatistic.DetectErrors()
		isError = isError || pfeStatistic.IsError
		routeTestContent = append(routeTestContent, pfeStatistic)

		if index < len(res.Links) {
			var routerStatistic = (*res.Links[index]).Forward.InterfaceInfo.ToRouterStatisticsContent(&router)
			routerStatistic.DetectErrors()
			isError = isError || routerStatistic.IsError
			routeTestContent = append(routeTestContent, routerStatistic)
		}

		routerTestResult := models.TestResult{Id: strconv.Itoa(router.Id), Name: router.Name, Type: 2, Content: routeTestContent, IsError: isError}
		result.Results = append(result.Results, &routerTestResult)
	}

	// ICMP
	result.Results = append(result.Results, &res.IcmpResult)

	return &result
}

func (trb *TestResultBuilder) GetGroupRouters() []*models.LspGroupRouters {
	var result = []*models.LspGroupRouters{}

	for _, res := range trb.GroupTestResults {
		var routersIdList = []string{}

		for _, host := range res.Hosts {
			routersIdList = append(routersIdList, strconv.Itoa(host.Router.Id))
		}

		var groupRouters = models.LspGroupRouters{Id: res.GroupId, RoutersIdList: routersIdList}
		result = append(result, &groupRouters)
	}

	return result
}

func (trb *TestResultBuilder) getGroupTestResults(lspItem models.LspItem) *GroupTestResult {
	for _, res := range trb.GroupTestResults {
		if res.GroupId == lspItem.GroupId {
			return res
		}
	}

	return &GroupTestResult{}
}

func (trb *TestResultBuilder) getOrCreateGroupTestResults(lspItem models.LspItem) (res *GroupTestResult, needTests bool) {
	for _, res := range trb.GroupTestResults {
		if res.GroupId == lspItem.GroupId {
			return res, false
		}
	}

	newRes := GroupTestResult{GroupId: lspItem.GroupId, IngressIp: lspItem.IngressIp}

	trb.GroupTestResults = append(trb.GroupTestResults, &newRes)

	return &newRes, true
}

func (trb *TestResultBuilder) getLspGroupById(lspGroupId string) *models.LspGroup {
	for _, lspGroup := range trb.LspGroups {
		if lspGroup.Id == lspGroupId {
			return lspGroup
		}
	}

	return &models.LspGroup{}
}

func (trb *TestResultBuilder) getRouterByAddress(address string) *models.Router {
	for _, router := range trb.AllRouters {
		if router.Ip == address {
			return router
		}
	}

	return &models.Router{}
}

func (trb *TestResultBuilder) runHostTests(routerAddress string, isIngress bool, isEgress bool, res *GroupTestResult) (string, error) {
	var routerType = "Transit Router"

	switch {
	case isIngress:
		routerType = "Ingress"
	case isEgress:
		routerType = "Egress"
	}

	router := trb.getRouterByAddress(routerAddress)
	var pfeStatistic, err = command.GetPfeStatistic(*router)
	if err != nil {
		return "", err
	}

	host := HostTestResult{Router: *router, RouterType: routerType, PfeStatistic: pfeStatistic}
	res.Hosts = append(res.Hosts, &host)

	return router.GetAddress(), nil
}

func (trb *TestResultBuilder) runLinkDirectTests(currentIp string, nextIp string, res *GroupTestResult) (neighborId string, rtDestination string, result LinkDirectTestResult, localIp string, err error) {
	result = LinkDirectTestResult{}

	rtDestination, interfaceName, interfaceInfo, err := trb.getInterfaceInfo(currentIp, nextIp)
	if err != nil {
		return
	}

	if interfaceInfo == nil {
		err = errors.New("cannot get interface information")
		return
	}

	localIp = interfaceInfo.GetLocalIp(interfaceName)

	result.LogicalInterfaceName = interfaceName
	result.InterfaceInfo = interfaceInfo

	nextRouter, err := command.GetOspfNeighbor(currentIp, interfaceName)
	if err != nil {
		return
	}

	if nextRouter == nil {
		err = errors.New("cannot get ospf neighbor information")
		return
	}

	neighborId = (*nextRouter).NeighborId
	result.Destination = (*nextRouter).NeighborAddress

	return
}

func (trb *TestResultBuilder) getInterfaceInfo(fromIp string, toIp string) (destination string, interfaceName string, interfaceInfo models.IRouterStatistics, err error) {
	routeInfo, err := command.LoadRouteInfo(fromIp, toIp, "inet.0")
	if err != nil {
		return
	}

	destination = routeInfo.Rt.RtDestination
	interfaceName = routeInfo.GetInterfaceName()

	var physicalInterfaceName = utils.GetPhysicalName(interfaceName)

	if strings.Index(interfaceName, "ae") == 0 {
		interfaceInfo, err = command.LoadAggregateInterfaceInfo(fromIp, physicalInterfaceName)
	} else {
		interfaceInfo, err = command.LoadInterfaceInfo(fromIp, physicalInterfaceName)
	}

	return
}
