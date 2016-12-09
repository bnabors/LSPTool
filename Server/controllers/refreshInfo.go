/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controller

import (
	"strconv"
	"strings"

	"errors"

	"../commands"
	"../models"
	"../sessions"
	"../utils"
	"./helpers"
)

// RefreshInterfaceInfo - get updated interface information
func RefreshInterfaceInfo(router models.Router, interfaceName string) (models.RouterStatisticsContent, error) {
	if interfaceName == "pfe" {
		return refreshPfeStatistic(router)
	}

	sessionsManager := sessions.SessionsManager{}
	defer sessionsManager.CloseAllSessions()

	var statistic models.IRouterStatistics
	var err error
	if strings.Index(interfaceName, "ae") == 0 {
		statistic, err = command.LoadAggregateInterfaceInfo(&sessionsManager, router.GetAddress(), interfaceName)
	} else {
		statistic, err = command.LoadInterfaceInfo(&sessionsManager, router.GetAddress(), interfaceName)
	}

	if err != nil {
		return models.RouterStatisticsContent{}, err
	}
	return statistic.ToRouterStatisticsContent(&router), nil
}

// ClearAndRefreshInterfaceInfo - clear and get updated interface inforamtion
func ClearAndRefreshInterfaceInfo(router models.Router, interfaceName string) (models.RouterStatisticsContent, error) {
	if interfaceName == "pfe" {
		return clearAndRefreshPfeStatistic(router)
	}

	if err := command.ClearInterfacesStatistics(router.GetAddress(), interfaceName); err != nil {
		return models.RouterStatisticsContent{}, err
	}

	return RefreshInterfaceInfo(router, interfaceName)
}

// RefreshRouter - get updated router information. Parallel
func RefreshRouter(router models.Router, interfaceNames []*string) (models.TestResult, error) {
	return getStatisticsContent(router, interfaceNames, RefreshInterfaceInfo)
}

// ClearAndRefreshRouter - clear and get updated router information. Parallel
func ClearAndRefreshRouter(router models.Router, interfaceNames []*string) (models.TestResult, error) {
	return getStatisticsContent(router, interfaceNames, ClearAndRefreshInterfaceInfo)
}

func getStatisticsContent(router models.Router, interfaceNames []*string, statisticsContentFunc func(router models.Router, name string) (models.RouterStatisticsContent, error)) (models.TestResult, error) {
	var contents = make([]models.RouterStatisticsContent, len(interfaceNames))
	isError := false
	for index, interfaceName := range interfaceNames {
		routerStatistic, err := statisticsContentFunc(router, *interfaceName)
		if err != nil {
			return models.TestResult{}, err
		}
		routerStatistic.DetectErrors()
		isError = isError || routerStatistic.IsError
		contents[index] = routerStatistic
	}
	return models.TestResult{Id: strconv.Itoa(router.Id), Name: router.Name, Type: 2, Content: contents, IsError: isError}, nil
}

// refreshPfeStatistic - get updated interface information
func refreshPfeStatistic(router models.Router) (models.RouterStatisticsContent, error) {
	statistic, err := command.GetPfeStatistic(router)
	if err != nil {
		return models.RouterStatisticsContent{}, err
	}

	return statistic.ToRouterStatisticsContent(&router), nil
}

// clearAndRefreshPfeStatistic - clear and get updated interface inforamtion
func clearAndRefreshPfeStatistic(router models.Router) (models.RouterStatisticsContent, error) {
	if err := command.ClearPfeStatistic(router); err != nil {
		return models.RouterStatisticsContent{}, err
	}

	return refreshPfeStatistic(router)
}

// RefreshPing get icmp result between two routers
func RefreshPing(requestModel models.RequestModel, options models.RefreshPingOptions) (models.IcmpResult, error) {
	var routers, err = utils.GetRoutersFromMysqlByID([]string{options.RouterStartId, options.RouterFinishId})

	if err != nil {
		return models.IcmpResult{}, err
	}

	if len(routers) != 2 {
		return models.IcmpResult{}, errors.New("Need two routers")
	}

	var routerLeft = routers[0]
	var routerRight = routers[1]

	return controllerHelper.BuildIcmpResult(*routerLeft, *routerRight)
}

// RefreshPings get icmp results by route
func RefreshPings(optionsId string, requestModel models.RequestModel) (models.TestResult, error) {
	var groupRouters *models.LspGroupRouters = nil

	for _, res := range requestModel.LspGroupRouters {
		if (*res).Id == requestModel.LspItem.GroupId {
			groupRouters = res
			break
		}
	}

	if groupRouters == nil {
		return models.TestResult{}, errors.New("Unable to get the router group")
	}

	var routers = []*models.Router{}

	for _, routerId := range (*groupRouters).RoutersIdList {

		var router, err = utils.GetRouterFromMysqlByID(routerId)
		if err != nil {
			return models.TestResult{}, err
		}

		routers = append(routers, router)
	}

	return controllerHelper.BuildIcmpResultByRouters(optionsId, routers)
}