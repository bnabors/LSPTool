/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package command

import (
	"fmt"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/utils"
)

// ClearInterfacesStatistics очищает статистику интерфеса interfaceName, для роутера по адресу address
func ClearInterfacesStatistics(address string, interfaceName string) error {
	var requestPattern = `<clear-interfaces-statistics>
	<interface-name>%s</interface-name>
</clear-interfaces-statistics>`

	lspLogger.Infoln("command clearInterfacesStatistics address: " + address + " interfaceName: " + interfaceName)

	var request = fmt.Sprintf(requestPattern, interfaceName)

	session, err := utils.CreateSession(address)
	if err != nil {
		lspLogger.Error(err)
		return err
	}
	defer session.Close()

	_, err = utils.MakeNetconfRequest(session, request)
	if err != nil {
		lspLogger.Error(err)
		return err
	}

	return nil
}
