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

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/config"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/sessions"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/utils"
)

func Ping(sm *sessions.SessionsManager, source models.Router, host models.Router) (models.PingResult, error) {

	var requestPattern = `<ping>
	<count>%d</count> 
	<size>%d</size> 
	<rapid/> 
	<source>%s</source> 
	<host>%s</host> 
</ping>`

	var request = fmt.Sprintf(requestPattern, config.LspConfig.PingCount, config.LspConfig.PingSize, source.Ip, host.Ip)

	commandDescription := "command ping from: " + source.Name + " to: " + host.Name + " request: " + request
	lspLogger.Infoln(commandDescription)

	session, err := sm.GetSession(source.GetAddress())
	if err != nil {
		lspLogger.Error(err, request)
		return models.PingResult{}, errors.New(err.Error() + "\r\n Information: " + commandDescription)
	}

	reply, err := utils.MakeNetconfRequest(session, request)
	if err != nil {
		lspLogger.Error(err, request)
		return models.PingResult{}, errors.New(err.Error() + "\r\n Information: " + commandDescription)
	}

	return models.ParsePing([]byte(reply.Data)), nil
}
