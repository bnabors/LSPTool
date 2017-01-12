/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/WOWLABS/LSPTool/Server/log"
)

type (
	LspToolConfig struct {
		User                     string   `json:"user"`
		Password                 string   `json:"password"`
		ServerIP                 string   `json:"serverIp"`
		ServerPort               int      `json:"serverPort"`
		ConnectionCount          int      `json:"connectionCount"`
		MysqlConnectionsString   string   `json:"mysqlConnectionsString"`
		MysqlRouterQuery         string   `json:"mysqlRouterQuery"`
		UseProxy                 bool     `json:"useProxy"`
		LogLevel                 string   `json:"logLevel"`
		LogFile                  string   `json:"logFile"`
		IngressRouterNames       []string `json:"ingressRouterNames"`
		EgressRouterNames        []string `json:"egressRouterNames"`
		SSHConnectionTimout      int      `json:"sshConnectionTimout"`
		PingCount                int      `json:"pingCount"`
		PingSize                 int      `json:"pingSize"`
		PingLossPercentThreshold int      `json:"pingLossPercentThreshold"`
		PingAvgThreshold         float64  `json:"pingAvgThreshold"`
		PingMaxThreshold         float64  `json:"pingMaxThreshold"`
		PingSTDDevThreshold      float64  `json:"pingSTDDevThreshold"`
	}
)

var LspConfig LspToolConfig

func ReadConfig() {
	data, _ := ioutil.ReadFile("lsp.config.json")
	//lspLogger.CheckError(err)
	//lspLogger.Debugln(string(data))

	var lspConfig LspToolConfig

	if err := json.Unmarshal(data, &lspConfig); err != nil {
		lspLogger.Error(err)
	}

	if lspConfig.UseProxy {
		lspConfig.MysqlRouterQuery = "Select element_id,hostname,loopbackip, proxyIp4 from router"
	} else {
		lspConfig.MysqlRouterQuery = "Select element_id,hostname,loopbackip from bgpmanager.element"
	}

	LspConfig = lspConfig

	return
}
