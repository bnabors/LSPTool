/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controllerHelper

import (
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/commands"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
)

// BuildIcmpResult make icmp result between two routers
func BuildIcmpResult(routerLeft models.Router, routerRight models.Router) (models.IcmpResult, error) {
	var pingResult, err = command.Ping(routerLeft, routerRight)
	if err != nil {
		return models.IcmpResult{}, err
	}
	return pingResult.ToIcmpResult(routerLeft, routerRight), nil
}

// BuildIcmpResultByRouters make icmp results by route
func BuildIcmpResultByRouters(optionsId string, routers []*models.Router) (models.TestResult, error) {
	content := []*models.IcmpResult{}

	isError := false
	for index := 1; index < len(routers); index++ {

		var routerLeft = routers[index-1]
		var routerRight = routers[index]

		var icmpResult, err = BuildIcmpResult(*routerLeft, *routerRight)
		if err != nil {
			return models.TestResult{}, err
		}

		content = append(content, &icmpResult)
		isError = isError || icmpResult.IsError
	}

	result := models.TestResult{
		Id:      optionsId,
		Name:    "ICMP",
		Type:    3,
		IsError: isError,
		Content: content}

	return result, nil
}
