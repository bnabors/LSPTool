/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package command

import (
	"bytes"
	"strings"
	"time"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/config"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"

	"golang.org/x/crypto/ssh"
)

// GetPfeStatistic возвращает PFE статистику
func GetPfeStatistic(host models.Router) (models.PfeStatistic, error) {
	lspLogger.Infoln("command getPfeStatistic router: " + host.Name)

	getPfeStatisticCommamd := "show pfe statistics traffic"
	client, err := createSSHClient(config.LspConfig.User, config.LspConfig.Password, host)
	if err != nil {
		return models.PfeStatistic{}, err
	}
	defer client.Close()

	commandResult, err := runCommand(client, getPfeStatisticCommamd)
	if err != nil {
		return models.PfeStatistic{}, err
	}
	result := parseStatistic(commandResult)
	return result, nil
}

// ClearPfeStatistic очищает PFE статистику
func ClearPfeStatistic(host models.Router) error {
	clearPfeStatisticCommamd := "clear pfe statistics traffic"
	client, err := createSSHClient(config.LspConfig.User, config.LspConfig.Password, host)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = runCommand(client, clearPfeStatisticCommamd)
	return err
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

func createSSHClient(user string, password string, router models.Router) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: time.Duration(config.LspConfig.SSHConnectionTimout) * time.Second,
	}
	client, err := ssh.Dial("tcp", router.GetAddress(), config)
	if err != nil {
		lspLogger.Errorf("SshCommand Error: %v", err)
		return nil, err
	}
	return client, nil
}

func runCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		lspLogger.Errorf("SshCommand Error: %v", err)
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		lspLogger.Errorf("SshCommand Error: %v", err)
		return "", err
	}
	return b.String(), nil
}

func sshTest() {
	router := models.Router{Name: "r1", Ip: "172.31.0.1", PuttyIp: "127.0.0.1:2001"}

	var pfeStatistics, _ = GetPfeStatistic(router)
	var content = pfeStatistics.ToRouterStatisticsContent(&router)

	lspLogger.Debug(content)
}
