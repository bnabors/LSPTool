/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package command

import (
	"fmt"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/sessions"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/utils"
)

func LoadRouteInfo(sm *sessions.SessionsManager, address string, destination string, table string) (models.RouteInformation, error) {
	var requestPattern = `<get-route-information>
	<destination>%s</destination>
	<table>%s</table>
</get-route-information>
`
	lspLogger.Infoln("command getRouteInfo from: " + address + " to: " + destination)

	var request = fmt.Sprintf(requestPattern, destination, table)

	session, err := sm.GetSession(address)
	if err != nil {
		lspLogger.Error(err)
		return models.RouteInformation{}, err
	}

	reply, err := utils.MakeNetconfRequest(session, request)
	if err != nil {
		lspLogger.Error(err)
		return models.RouteInformation{}, err
	}

	return models.ParseRouteInformation([]byte(reply.Data)), nil
}
