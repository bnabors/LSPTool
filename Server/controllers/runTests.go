/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controller

import (
	"../models"
	"./helpers"
)

func RunTests(o models.TestOptions) (*models.TestResults, error) {
	routeResults := models.TestResults{P2P: []*models.RouteResult{}, P2MP: []*models.RouteResult{}}

	builder := controllerHelper.TestResultBuilder{}
	builder.Init(o.LspGroups)
	defer builder.Dispose()

	for _, testLsp := range o.P2P {
		if !testLsp.IsSelected {
			continue
		}
		if err := builder.TryRunTest(testLsp.LspItem); err != nil {
			return nil, err
		}
	}

	for _, testLsp := range o.P2MP {
		if !testLsp.IsSelected {
			continue
		}
		if err := builder.TryRunTest(testLsp.LspItem); err != nil {
			return nil, err
		}
	}

	for _, testLsp := range o.P2P {
		builder.AddLspBand(testLsp.LspItem)
	}

	for _, testLsp := range o.P2MP {
		builder.AddLspBand(testLsp.LspItem)
	}

	for _, testLsp := range o.P2P {
		if testLsp.IsSelected {
			routeResults.P2P = append(routeResults.P2P, builder.GetRouteResult(testLsp.LspItem))
		}
	}

	for _, testLsp := range o.P2MP {
		if testLsp.IsSelected {
			routeResults.P2MP = append(routeResults.P2MP, builder.GetRouteResult(testLsp.LspItem))
		}
	}

	routeResults.GroupRouters = builder.GetGroupRouters()

	return &routeResults, nil
}

func GetLspItemTestResult(lspItem models.LspItem, lspGroups []*models.LspGroup) (*models.RouteResult, error) {
	builder := controllerHelper.TestResultBuilder{}
	builder.Init(lspGroups)
	defer builder.Dispose()

	if err := builder.TryRunTest(lspItem); err != nil {
		return nil, err
	}
	builder.AddLspBand(lspItem)

	return builder.GetRouteResult(lspItem), nil
}
