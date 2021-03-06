/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package command

import (
	"errors"
	"fmt"

	"github.com/WOWLABS/LSPTool/Server/log"
	"github.com/WOWLABS/LSPTool/Server/models"
	"github.com/WOWLABS/LSPTool/Server/utils"
)

func LoadRouteInfo(address string, destination string, table string) (models.RouteInformation, error) {
	var requestPattern = `<get-route-information>
	<destination>%s</destination>
	<table>%s</table>
</get-route-information>
`
	commandDescription := "command getRouteInfo from: " + address + " to: " + destination
	lspLogger.Infoln(commandDescription)

	var request = fmt.Sprintf(requestPattern, destination, table)

	reply, err := utils.SshSessionManager.DoNetconfRequest(address, request)
	if err != nil {
		lspLogger.Error(err, request)
		return models.RouteInformation{}, errors.New(err.Error() + "\r\n Information: " + commandDescription)
	}

	lspLogger.Debug(reply.Data)

	return models.ParseRouteInformation([]byte(reply.Data)), nil
}
