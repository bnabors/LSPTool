/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package command

import (
	"github.com/Juniper/24287_WOW_LSP_GOLANG/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/models"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/sessions"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/utils"
)

func GetOspfNeighbor(sm *sessions.SessionsManager, address string, interfaceName string) (*models.OspfNeighbor, error) {
	var request = `<get-ospf-neighbor-information>
</get-ospf-neighbor-information>
`

	lspLogger.Infoln("command getOspfNeighbor address: " + address + " interfaceName: " + interfaceName)

	session, err := sm.GetSession(address)
	if err != nil {
		lspLogger.Error(err)
		return nil, err
	}

	reply, err := utils.MakeNetconfRequest(session, request)
	if err != nil {
		lspLogger.Error(err)
		return nil, err
	}

	answer := models.ParseOspfNeighborInformation([]byte(reply.Data))

	return answer.GetOspfNeighbor(interfaceName), nil
}
