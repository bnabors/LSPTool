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

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/commands"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/controllers/helpers"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
)

func RefreshLsp(requestModel models.RequestModel, oldLspGroups []*models.LspGroup) (*models.RouteResult, error) {
	var determineOptions = models.DetermineOptions{Ingress: requestModel.Ingress, Egress: requestModel.Egress}

	lspCollections, err := DetermineLsps(determineOptions)
	if err != nil {
		return nil, err
	}

	if isLspChanged(lspCollections, &requestModel.LspItem, requestModel.LspNames, oldLspGroups) {
		return nil, errors.New("LSP has been changed. Please determine LSPs again.")
	}

	routeResult, err := GetLspItemTestResult(requestModel.LspItem, lspCollections, oldLspGroups)
	if err != nil {
		return nil, err
	}

	routeResult.Names = requestModel.LspNames

	return routeResult, nil
}

func isLspChanged(newLspCollection models.LspCollection, oldLsp *models.LspItem, oldLspNames []string, oldLspGroups []*models.LspGroup) bool {
	var checker = controllerHelper.RouteChecker{NewGroups: newLspCollection.LspGroups, OldGroups: oldLspGroups}

	for _, oldLspName := range oldLspNames {
		newLspP2P := controllerHelper.GetLspByName(oldLspName, newLspCollection.P2P)
		if newLspP2P != nil {
			if checker.IsLspChanged(newLspP2P, oldLsp) {
				return true
			} else {
				continue
			}
		}

		newLspP2MP := controllerHelper.GetLspByName(oldLspName, newLspCollection.P2MP)
		if newLspP2MP != nil {
			if checker.IsLspChanged(newLspP2MP, oldLsp) {
				return true
			} else {
				continue
			}
		}

		return true
	}

	return false
}

// RefreshInterfaceInfo - get updated interface information
func RefreshInterfaceInfo(router models.Router, interfaceName string) (models.RouterStatisticsContent, error) {
	if interfaceName == "pfe" {
		return refreshPfeStatistic(router)
	}

	var statistic models.IRouterStatistics
	var err error
	if isAEInterface(interfaceName) {
		statistic, err = command.LoadAggregateInterfaceInfo(router.GetAddress(), interfaceName)
	} else {
		statistic, err = command.LoadInterfaceInfo(router.GetAddress(), interfaceName)
	}

	if err != nil {
		return models.RouterStatisticsContent{}, err
	}
	return statistic.ToRouterStatisticsContent(), nil
}

// ClearAndRefreshInterfaceInfo - clear and get updated interface information
func ClearAndRefreshInterfaceInfo(router models.Router, interfaceName string) (models.RouterStatisticsContent, error) {
	if interfaceName == "pfe" {
		return clearAndRefreshPfeStatistic(router)
	}

	if isAEInterface(interfaceName) {
		// clear SubInterfaces for AgregateInterface
		statistic, err := RefreshInterfaceInfo(router, interfaceName)
		if err != nil {
			return models.RouterStatisticsContent{}, err
		}
		for _, subInterface := range statistic.SubInterfaces {
			if err := command.ClearInterfacesStatistics(router.GetAddress(), subInterface.Name); err != nil {
				return models.RouterStatisticsContent{}, err
			}
		}
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

func isAEInterface(interfaceName string) bool {
	return strings.Index(interfaceName, "ae") == 0
}

// refreshPfeStatistic - get updated interface information
func refreshPfeStatistic(router models.Router) (models.RouterStatisticsContent, error) {
	statistic, err := command.GetPfeStatistic(router)
	if err != nil {
		return models.RouterStatisticsContent{}, err
	}

	return statistic.ToRouterStatisticsContent(), nil
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
	return controllerHelper.BuildIcmpResult(options.Info)
}

// RefreshPings get icmp results by route
func RefreshPings(optionsId string, requestModel models.RequestModel, options models.RefreshPingsOptions) (models.TestResult, error) {
	return controllerHelper.BuildIcmpResultByRouters(optionsId, options.Infos)
}
