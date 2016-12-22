/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package command

import (
	"errors"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/utils"
)

func LoadMplsLspInfo(ingressAddress string, egressAddress string) ([]models.MplsLsp, error) {
	var request = `<get-mpls-lsp-information>
    <ingress/>
    <extensive/>
</get-mpls-lsp-information>
`
	commandDescription := "command loadMplsLspInfo ingressAddress: " + ingressAddress + " egressAddress: " + egressAddress
	lspLogger.Infoln(commandDescription)

	reply, err := utils.SshSessionManager.DoNetconfRequest(ingressAddress, request)
	if err != nil {
		lspLogger.Error(err, request)
		return nil, errors.New(err.Error() + "\r\n Information: " + commandDescription)
	}

	mplsLspInformation := models.ParseMplsLspInformation([]byte(reply.Data))

	return mplsLspInformation.FilterMplsLspInfoByEgress(egressAddress), nil
}
