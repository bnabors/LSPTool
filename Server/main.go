/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package main

import (
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/config"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
)

func main() {
	config.ReadConfig()

	lspLogger.Initialize(config.LspConfig.LogFile, config.LspConfig.LogLevel)
	lspLogger.Infoln("Server has started")

	StartServer(config.LspConfig.ServerIP, config.LspConfig.ServerPort, config.LspConfig.ConnectionCount)
}
