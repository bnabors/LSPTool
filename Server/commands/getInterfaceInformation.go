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
	"strings"

	"github.com/WOWLABS/LSPTool/Server/log"
	"github.com/WOWLABS/LSPTool/Server/models"
	"github.com/WOWLABS/LSPTool/Server/utils"
)

const (
	requestPattern = `<get-interface-information> 
	<extensive/> 
	<interface-name>%s</interface-name>
</get-interface-information>
`
)

// LoadInterfaceInfo returns information about interface with name is interfaceName in router by address (possible nil)
func LoadInterfaceInfo(address string, interfaceName string) (models.IRouterStatistics, error) {
	commandDescription := "command loadInterfaceInfo address: " + address + " interfaceName: " + interfaceName
	lspLogger.Infoln(commandDescription)

	var request = fmt.Sprintf(requestPattern, interfaceName)

	reply, err := utils.SshSessionManager.DoNetconfRequest(address, request)
	if err != nil {
		lspLogger.Error(err, request)
		return nil, errors.New(err.Error() + "\r\n Information: " + commandDescription)
	}

	lspLogger.Debug(reply.Data)

	result := models.ParseInterfaceInformation([]byte(reply.Data))

	utils.ConvertToJson(result)
	return result, nil
}

// LoadAggregateInterfaceInfo returns information about agregate interface with name is interfaceName in router by address (possible nil)
func LoadAggregateInterfaceInfo(address string, interfaceName string) (models.IRouterStatistics, error) {
	commandDescription := "command loadAggregateInterfaceInfo address: " + address + " interfaceName: " + interfaceName
	lspLogger.Infoln(commandDescription)

	var request = fmt.Sprintf(requestPattern, interfaceName)

	reply, err := utils.SshSessionManager.DoNetconfRequest(address, request)
	if err != nil {
		lspLogger.Error(err, request)
		return nil, errors.New(err.Error() + "\r\n Information: " + commandDescription)
	}
	result := models.ParseAgregateInterfaceInformation([]byte(reply.Data))

	subInterfaceNames := result.GetSubInterfaceNames()

	result.SubInterface, err = getSubInterfaceInfo(address, subInterfaceNames)

	utils.ConvertToJson(result)
	return result, err
}

func getSubInterfaceInfo(address string, subInterfaceNames []string) ([]models.InterfaceInformation, error) {
	result := make([]models.InterfaceInformation, len(subInterfaceNames))
	for index, logicalName := range subInterfaceNames {
		subInterfaceName := strings.TrimSpace(utils.GetPhysicalName(logicalName))
		interfaceInfo, err := LoadInterfaceInfo(address, subInterfaceName)
		if err != nil {
			lspLogger.Error(err)
			return nil, err
		}
		result[index] = interfaceInfo.(models.InterfaceInformation)
	}

	return result, nil
}
