/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package command

import (
	"errors"
	"strings"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/utils"
)

var (
	sshSessionManager = utils.NewSSHSessionManager()
)

func GetPfeStatistic(host models.Router) (models.PfeStatistic, error) {
	lspLogger.Infoln("command getPfeStatistic router: " + host.Name)

	getPfeStatisticCommamd := "show pfe statistics traffic"
	commandResult, err := sshSessionManager.RunSSHCommand(host, getPfeStatisticCommamd)
	if err != nil {
		return models.PfeStatistic{}, errors.New(err.Error() + "\r\n Information: command " + getPfeStatisticCommamd)
	}
	result := parseStatistic(commandResult)
	return result, nil
}

func ClearPfeStatistic(host models.Router) error {
	lspLogger.Infoln("command clearPfeStatistic router: " + host.Name)
	clearPfeStatisticCommamd := "clear pfe statistics traffic"
	_, err := sshSessionManager.RunSSHCommand(host, clearPfeStatisticCommamd)
	if err != nil {
		return errors.New(err.Error() + "\r\n Information: command " + clearPfeStatisticCommamd)
	}
	return nil
}

func parseStatistic(commandResult string) models.PfeStatistic {
	result := models.PfeStatistic{}
	temp := strings.Split(commandResult, "\n")
	for _, line := range temp {
		if strings.Contains(line, "Timeout") {
			result.Timeout = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Truncated key") {
			result.TruncatedKey = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Bits to test") {
			result.BitsToTest = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Data error") {
			result.DataError = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Stack underflow") {
			result.StackUnderflow = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Stack overflow") {
			result.StackOverflow = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Extended discard") {
			result.ExtendedDiscard = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Invalid interface") {
			result.InvalidInterface = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Info cell drops") {
			result.InfoCellDrops = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Fabric drops") {
			result.FabricDrops = getValueFromLine(line)
			continue
		}
		if strings.Contains(line, "Output MTU") {
			result.OutputMTU = getValueFromLine(line)
			continue
		}

	}
	return result
}

func getValueFromLine(line string) string {
	pair := strings.Split(line, ":")
	if len(pair) != 2 || pair[1] == "" {
		return ""
	}

	return strings.TrimSpace(pair[1])
}
