/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package utils

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Juniper/go-netconf/netconf"
	"github.com/WOWLABS/LSPTool/Server/config"
	"github.com/WOWLABS/LSPTool/Server/log"
	"github.com/WOWLABS/LSPTool/Server/models"
	"golang.org/x/crypto/ssh"
)

type (
	SSHRouterSession struct {
		address        string
		sshClient      *ssh.Client
		netconfSession *netconf.Session
		mutex          sync.Mutex
	}

	SSHSessionManager struct {
		routerSessions map[string]*SSHRouterSession
	}
)

var (
	SshSessionManager = NewSSHSessionManager()
)

func NewSSHSessionManager() SSHSessionManager {
	return SSHSessionManager{
		routerSessions: map[string]*SSHRouterSession{},
	}
}

func (sm *SSHSessionManager) RunSSHCommand(router models.Router, command string) (string, error) {
	address := router.GetAddress()
	session := sm.getSession(address)

	session.mutex.Lock()
	defer session.mutex.Unlock()

	//first try remove client
	result, err := session.tryRunSSHCommand(command)
	if err == nil {
		return result, nil
	}

	session.closeSSHClient()

	lspLogger.Infoln("second try executing command : " + command)

	//second try with error
	return session.tryRunSSHCommand(command)
}

func (sm *SSHSessionManager) DoNetconfRequest(address string, request string) (*netconf.RPCReply, error) {
	session := sm.getSession(address)

	session.mutex.Lock()
	defer session.mutex.Unlock()

	//first try
	result, err := session.tryDoNetconfRequest(request)
	if err == nil {
		return result, nil
	}

	session.closeNetconfSession()

	lspLogger.Infoln("second try executing request : " + request)

	//second try with error
	return session.tryDoNetconfRequest(request)
}

func (sm *SSHSessionManager) getSession(address string) *SSHRouterSession {
	session, hasSession := sm.routerSessions[address]
	if hasSession {
		return session
	}

	newSession := SSHRouterSession{address: address, mutex: sync.Mutex{}}
	sm.routerSessions[address] = &newSession

	return &newSession
}

func (rs *SSHRouterSession) tryRunSSHCommand(command string) (string, error) {
	if rs.sshClient == nil {
		err := rs.createSSHClient()
		if err != nil {
			return "", err
		}

		lspLogger.Infoln("create ssh client: " + rs.address)
	} else {
		lspLogger.Infoln("ssh client is already open: " + rs.address)
	}

	if rs.sshClient == nil {
		return "", errors.New("ssh client opening failed: " + rs.address)
	}

	session, err := rs.sshClient.NewSession()
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

func (rs *SSHRouterSession) tryDoNetconfRequest(request string) (*netconf.RPCReply, error) {
	if rs.netconfSession == nil {
		err := rs.createNetconfSession()
		if err != nil {
			return nil, err
		}

		lspLogger.Infoln("create netconf session: " + rs.address)
	} else {
		lspLogger.Infoln("netconf session is already open: " + rs.address)
	}

	if rs.netconfSession == nil {
		return nil, errors.New("netconf session opening failed: " + rs.address)
	}

	return rs.netconfSession.Exec(netconf.RawMethod(request))
}

func (rs *SSHRouterSession) createSSHClient() error {
	cfg := &ssh.ClientConfig{
		User: config.LspConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.LspConfig.Password),
		},
		Timeout: time.Duration(config.LspConfig.SSHConnectionTimout) * time.Second,
	}

	var finalAddress = ""
	if config.LspConfig.UseProxy {
		finalAddress = rs.address
	} else if strings.Contains(rs.address, ":") {
		finalAddress = rs.address
	} else {
		finalAddress = rs.address + ":22"
	}

	client, err := ssh.Dial("tcp", finalAddress, cfg)
	if err != nil {
		lspLogger.Errorf("SshCommand Error: %v", err)
		return err
	}

	rs.sshClient = client

	return nil
}

func (rs *SSHRouterSession) createNetconfSession() error {
	user, password := config.LspConfig.User, config.LspConfig.Password

	var timeout = time.Duration(config.LspConfig.SSHConnectionTimout) * time.Second

	var finalAddress = ""
	if config.LspConfig.UseProxy {
		finalAddress = rs.address
	} else if strings.Contains(rs.address, ":") {
		finalAddress = rs.address
	} else {
		finalAddress = rs.address + ":22"
	}

	session, err := netconf.DialSSHTimeout(finalAddress, netconf.SSHConfigPassword(user, password), timeout)
	if err != nil {
		lspLogger.Errorf("Netconf Error: %v", err)
		return err
	}

	rs.netconfSession = session

	return nil
}

func (rs *SSHRouterSession) closeSSHClient() {
	if rs.sshClient != nil {
		rs.sshClient = nil
	}

	lspLogger.Infoln("close ssh client: " + rs.address)
}

func (rs *SSHRouterSession) closeNetconfSession() {
	if rs.netconfSession != nil {
		rs.netconfSession.Close()
		rs.netconfSession = nil
	}

	lspLogger.Infoln("close netconf session: " + rs.address)
}
