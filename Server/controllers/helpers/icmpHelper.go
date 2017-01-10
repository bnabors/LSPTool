/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controllerHelper

import (
	"strconv"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/commands"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/config"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/helpers"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
)

// BuildIcmpResult make icmp result between two routers
func BuildIcmpResult(icmpInfo *models.IcmpInfo) (models.IcmpResult, error) {
	var pingResult, err = command.Ping(icmpInfo)
	if err != nil {
		return models.IcmpResult{}, err
	}
	return pingResult.ToIcmpResult(icmpInfo), nil
}

// BuildIcmpResultByRouters make icmp results by route
func BuildIcmpResultByRouters(optionsId string, icmpInfos []*models.IcmpInfo) (models.TestResult, error) {
	content := []*models.IcmpResult{}

	isError := false

	for _, icmpInfo := range icmpInfos {

		var icmpResult, err = BuildIcmpResult(icmpInfo)
		if err != nil {
			return models.TestResult{}, err
		}

		content = append(content, &icmpResult)
		isError = isError || icmpResult.IsError
	}

	result := models.TestResult{
		Id:          optionsId,
		Name:        "ICMP",
		Description: " - each ping test is " + strconv.Itoa(config.LspConfig.PingCount) + " count at " + helpers.ParceNumberAndLocalize(strconv.Itoa(config.LspConfig.PingSize)) + " bytes",
		Type:        3,
		IsError:     isError,
		Content:     content}

	return result, nil
}
