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

	"github.com/WOWLABS/LSPTool/Server/config"
	"github.com/WOWLABS/LSPTool/Server/log"
	"github.com/WOWLABS/LSPTool/Server/models"
	"github.com/WOWLABS/LSPTool/Server/utils"
)

func Ping(icmpInfo *models.IcmpInfo) (models.PingResult, error) {

	var requestPattern = `<ping>
	<count>%d</count> 
	<size>%d</size> 
	<rapid/> 
	<source>%s</source> 
	<host>%s</host> 
</ping>`

	var request = fmt.Sprintf(requestPattern, config.LspConfig.PingCount, config.LspConfig.PingSize, icmpInfo.InterfaceIpSource, icmpInfo.InterfaceIpDest)

	commandDescription := "command ping from: " + icmpInfo.Source.Name + " to: " + icmpInfo.Destination.Name + " request: " + request
	lspLogger.Infoln(commandDescription)

	reply, err := utils.SshSessionManager.DoNetconfRequest(icmpInfo.Source.GetAddress(), request)
	if err != nil {
		lspLogger.Error(err, request)
		return models.PingResult{}, errors.New(err.Error() + "\r\n Information: " + commandDescription)
	}

	return models.ParsePing([]byte(reply.Data)), nil
}
