/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package models

import (
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/config"
)

/* Main page */

type Router struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Ip      string `json:"ip"`
	ProxyIp string `json:"proxyip"`
}

func (router Router) GetAddress() string {
	if router.ProxyIp != "" && config.LspConfig.UseProxy {
		return router.ProxyIp
	}
	return router.Ip
}

type RouterContainer struct {
	IngressRouters []*Router `json:"ingress_routers"`
	EgressRouters  []*Router `json:"egress_routers"`
}

/* Determine LSPs */

type DetermineOptions struct {
	Ingress Router `json:"ingress"`
	Egress  Router `json:"egress"`
}

type LspItem struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IngressIp      string `json:"ingressIp"`
	IngressProxyip string `json:"ingressProxyip"`
	EgressIp       string `json:"egressIp"`
	EgressProxyip  string `json:"egressProxyip"`
	GroupId        string `json:"groupId"`
	Bandwidth      string `json:"bandwidth"`
}

type LspGroup struct {
	Id          string   `json:"id"`
	ReceivedRro []string `json:"receivedRro"`
}

type LspGroupRouters struct {
	Id            string   `json:"id"`
	RoutersIdList []string `json:"routersIdList"`
}

type LspCollection struct {
	P2P       []*LspItem  `json:"p2p"`
	P2MP      []*LspItem  `json:"p2mp"`
	Down      []*LspItem  `json:"down"`
	LspGroups []*LspGroup `json:"lspGroup"`
}

/* Run tests */

type TestLspOptions struct {
	LspItem    LspItem `json:"lsp"`
	IsSelected bool    `json:"selected"`
}

type TestOptions struct {
	Ingress   Router            `json:"ingress"`
	Egress    Router            `json:"egress"`
	P2P       []*TestLspOptions `json:"p2p"`
	P2MP      []*TestLspOptions `json:"p2mp"`
	LspGroups []*LspGroup       `json:"lspGroup"`
}

type LspBand struct {
	Name      string `json:"name"`
	Bandwidth string `json:"bandwidth"`
}

type LspRouter struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type LspRoute struct {
	Type   string    `json:"type"`
	Router LspRouter `json:"router"`
}

type LspContent struct {
	Bands   []*LspBand    `json:"bands"`
	Routes  []*LspRoute   `json:"routes"`
	Diagram DiagramResult `json:"diagram"`
}

type TestResult struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Type    byte        `json:"type"`
	Content interface{} `json:"content"`
	IsError bool        `json:"isError"`
}

type RouteResult struct {
	Id      string        `json:"id"`
	Names   []string      `json:"names"`
	Results []*TestResult `json:"results"`
	LspItem LspItem       `json:"lsp"`
}

type TestResults struct {
	P2P          []*RouteResult     `json:"p2p"`
	P2MP         []*RouteResult     `json:"p2mp"`
	GroupRouters []*LspGroupRouters `json:"groupRouters"`
}

type DiaRouter struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Ip            string `json:"ip"`
	Interface     string `json:"interface"`
	InputBytes    string `json:"inputb"`
	OutputBytes   string `json:"outputb"`
	InputPackets  string `json:"inputp"`
	OutputPackets string `json:"outputp"`
}

type DiaPath struct {
	Router1 DiaRouter `json:"router1"`
	Router2 DiaRouter `json:"router2"`
	BaseIp  string    `json:"baseIp"`
}

type DiagramResult struct {
	Paths []*DiaPath `json:"paths"`
}

type IcmpResult struct {
	Id                string `json:"id"`
	FromDevice        string `json:"fromDevice"`
	DestinationDevice string `json:"destDevice"`
	DestinationIp     string `json:"destIp"`
	Loss              string `json:"loss"`
	Max               string `json:"max"`
	Average           string `json:"avg"`
	StdDev            string `json:"stdDev"`
	RouterStartId     string `json:"routerStartId"`
	RouterFinishId    string `json:"routerFinishId"`
	IsError           bool   `json:"isError"`
}
type IcmpContent struct {
	IcmpResults []*IcmpResult `json:"icmps"`
	IsError     bool
}

type StatisticsValue struct {
	Header  string `json:"header"`
	Value   string `json:"value"`
	IsError bool   `json:"isError"`
}

type Statistics struct {
	Title  string             `json:"title"`
	Values []*StatisticsValue `json:"values"`
}

type RouterStatisticsContent struct {
	Id            string                     `json:"id"`
	Name          string                     `json:"name"`
	Statistics    []*Statistics              `json:"statistics"`
	SubInterfaces []*RouterStatisticsContent `json:"sub_interfaces"`
	IsError       bool                       `json:"isError"`
}

func (obj *RouterStatisticsContent) DetectErrors() {
	obj.IsError = false

	for _, statistic := range obj.Statistics {
		for _, statisticValue := range statistic.Values {
			if statisticValue.IsError {
				obj.IsError = true
			}
		}
	}

	for _, subInterface := range obj.SubInterfaces {
		subInterface.DetectErrors()

		obj.IsError = obj.IsError || subInterface.IsError
	}
}

/* Base Request|Response */

type RequestModel struct {
	Ingress         Router             `json:"ingress"`
	Egress          Router             `json:"egress"`
	LspItem         LspItem            `json:"lsp"`
	LspNames        []string           `json:"names"`
	LspGroupRouters []*LspGroupRouters `json:"groupRouters"`
	Options         string             `json:"options"`
}
type ResponseModel struct {
	LspItem LspItem     `json:"lsp"`
	Result  interface{} `json:"result"`
}

type ResponseErrorModel struct {
	Error string `json:"error"`
}

/* Reload All */

type ReloadAllOptions struct {
	LspGroups []*LspGroup `json:"lspGroups"`
}

/* Refresh|ClearRefresh Router|Interface */

type RefreshRouterOptions struct {
	Router     string    `json:"router"`
	Interfaces []*string `json:"interfaces"`
}

type RouterInterface struct {
	Interface      string           `json:"interface"`
	InnerInterface *RouterInterface `json:"innerInterface"`
}

func (op RouterInterface) GetLeafName() string {
	if op.InnerInterface == nil {
		return op.Interface
	}
	return op.InnerInterface.GetLeafName()
}

type RefreshRouterInterfaceOptions struct {
	Router    string           `json:"router"`
	Interface *RouterInterface `json:"interface"`
}

func (op RefreshRouterInterfaceOptions) GetInterfaceName() string {
	if op.Interface == nil {
		return ""
	}

	return op.Interface.GetLeafName()
}

type RefreshRouterInterfaceResult struct {
	Router    string                  `json:"id"`
	Interface *RouterInterface        `json:"interface"`
	Result    RouterStatisticsContent `json:"result"`
}

/* Refresh ICMP */

type RefreshPingsOptions struct {
	Id string `json:"id"`
}
type RefreshPingsResult struct {
	Id      string        `json:"id"`
	Content []*IcmpResult `json:"result"`
}

type RefreshPingOptions struct {
	Id             string `json:"id"`
	IcmpId         string `json:"conentId"`
	RouterStartId  string `json:"routerStartId"`
	RouterFinishId string `json:"routerFinishId"`
}

type RefreshPingResult struct {
	Id     string     `json:"id"`
	Result IcmpResult `json:"result"`
}
