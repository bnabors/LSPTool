/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package sessions

import (
	"github.com/Juniper/go-netconf/netconf"

	"../log"
	"../utils"
)

type SessionsManager struct {
	Sessions map[string]*netconf.Session
}

func (obj *SessionsManager) GetSession(address string) (session *netconf.Session, err error) {
	if obj.Sessions == nil {
		obj.Sessions = map[string]*netconf.Session{}
	}

	session, ok := obj.Sessions[address]
	if ok {
		lspLogger.Infoln("ssh session is already open: " + address)
		return
	}

	session, err = utils.CreateSession(address)

	if session != nil {
		obj.Sessions[address] = session
		lspLogger.Infoln("create ssh session: " + address)
	}

	return
}

func (obj *SessionsManager) CloseSession(address string) {
	session, ok := obj.Sessions[address]
	if !ok {
		return
	}

	if session != nil {
		session.Close()
	}

	delete(obj.Sessions, address)

	lspLogger.Infoln("close ssh session: " + address)
}

func (obj *SessionsManager) CloseAllSessions() {
	for address := range obj.Sessions {
		obj.CloseSession(address)
	}
}
