/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package command

import (
	"fmt"

	"../config"
	"../log"
	"../models"
	"../utils"
)

// Ping - пингует заданный адрес
func Ping(source models.Router, host models.Router) (models.PingResult, error) {

	var requestPattern = `<ping>
	<count>%d</count> 
	<size>1000</size> 
	<rapid/> 
	<source>%s</source> 
	<host>%s</host> 
</ping>`

	var request = fmt.Sprintf(requestPattern, config.LspConfig.PingCount, source.Ip, host.Ip)

	lspLogger.Infoln("command ping from: " + source.Name + " to: " + host.Name + " request: " + request)

	session, err := utils.CreateSession(source.GetAddress())
	if err != nil {
		lspLogger.Error(err)
		return models.PingResult{}, err
	}
	defer session.Close()

	reply, err := utils.MakeNetconfRequest(session, request)
	if err != nil {
		lspLogger.Error(err)
		return models.PingResult{}, err
	}

	return models.ParsePing([]byte(reply.Data)), nil
}

func pingTest() {
	sourceRouter := models.Router{Id: 1, Name: "r1", Ip: "172.31.0.1", PuttyIp: "127.0.0.1:2001"}
	hostRouter := models.Router{Id: 2, Name: "r2", Ip: "172.31.0.2", PuttyIp: "127.0.0.1:2002"}

	pingResult, _ := Ping(sourceRouter, hostRouter)
	pingResult.Print()
	lspLogger.Debug(pingResult.ToIcmpResult(sourceRouter, hostRouter))
}
