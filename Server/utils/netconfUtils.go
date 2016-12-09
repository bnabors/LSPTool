/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package utils

import (
	"time"

	"../config"
	"github.com/Juniper/go-netconf/netconf"
)

func CreateSession(address string) (*netconf.Session, error) {
	user, password := config.LspConfig.User, config.LspConfig.Password

	var timeout = time.Duration(config.LspConfig.SSHConnectionTimout) * time.Second
	session, err := netconf.DialSSHTimeout(address, netconf.SSHConfigPassword(user, password), timeout)

	return session, err
}

func MakeNetconfRequest(session *netconf.Session, request string) (*netconf.RPCReply, error) {
	reply, err := session.Exec(netconf.RawMethod(request))

	return reply, err
}