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

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/config"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"

	"github.com/Juniper/go-netconf/netconf"

	"golang.org/x/crypto/ssh"
)

type (
	// SSHSessionManager - type of session manager for pfe statistics
	SSHSessionManager struct {
		routerSSHClients map[string]*ssh.Client
		mu               sync.Mutex

		sessions map[string]*netconf.Session
		mutex    sync.Mutex
	}
)

var (
	SshSessionManager = NewSSHSessionManager()
)

func (sm SSHSessionManager) getClientOnRouter(router models.Router) (*ssh.Client, error) {
	client, hasClient := sm.routerSSHClients[router.Ip]
	if hasClient {
		return client, nil
	}
	return sm.createSSHClientOnRouter(router)
}

func (sm SSHSessionManager) createSSHClientOnRouter(router models.Router) (*ssh.Client, error) {
	cfg := &ssh.ClientConfig{
		User: config.LspConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.LspConfig.Password),
		},
		Timeout: time.Duration(config.LspConfig.SSHConnectionTimout) * time.Second,
	}
	address := router.GetAddress()

	var finalAddress = ""
	if config.LspConfig.UseProxy {
		finalAddress = address
	} else if strings.Contains(address, ":") {
		finalAddress = address
	} else {
		finalAddress = address + ":22"
	}

	client, err := ssh.Dial("tcp", finalAddress, cfg)
	if err != nil {
		lspLogger.Errorf("SshCommand Error: %v", err)
		return nil, err
	}

	sm.routerSSHClients[router.Ip] = client
	return client, nil
}

// RunSSHCommand - try get or open session on router and run command
func (sm SSHSessionManager) RunSSHCommand(router models.Router, command string) (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	//first try remove client
	result, err := sm.tryRunSSHCommand(router, command)
	if err == nil {
		return result, nil
	}

	delete(sm.routerSSHClients, router.Ip)

	//second try with error
	return sm.tryRunSSHCommand(router, command)
}

func NewSSHSessionManager() SSHSessionManager {
	return SSHSessionManager{
		routerSSHClients: map[string]*ssh.Client{},
		mu:               sync.Mutex{},

		sessions: map[string]*netconf.Session{},
		mutex:    sync.Mutex{},
	}
}

func (sm SSHSessionManager) tryRunSSHCommand(router models.Router, command string) (string, error) {
	client, err := sm.getClientOnRouter(router)
	if err != nil {
		return "", err
	}

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

func (sm SSHSessionManager) GetSession(address string) (session *netconf.Session, err error) {
	if sm.sessions == nil {
		sm.sessions = map[string]*netconf.Session{}
	}

	session, ok := sm.sessions[address]
	if ok {
		lspLogger.Infoln("ssh session is already open: " + address)
		return
	}

	session, err = CreateSession(address)

	if session != nil {
		sm.sessions[address] = session
		lspLogger.Infoln("create ssh session: " + address)
	}

	return
}

func (sm SSHSessionManager) CloseSession(address string) {
	session, ok := sm.sessions[address]

	if !ok {
		return
	}

	if session != nil {
		session.Close()
	}

	delete(sm.sessions, address)

	lspLogger.Infoln("close ssh session: " + address)
}

func CreateSession(address string) (*netconf.Session, error) {
	user, password := config.LspConfig.User, config.LspConfig.Password

	var timeout = time.Duration(config.LspConfig.SSHConnectionTimout) * time.Second

	var finalAddress = ""
	if config.LspConfig.UseProxy {
		finalAddress = address
	} else if strings.Contains(address, ":") {
		finalAddress = address
	} else {
		finalAddress = address + ":22"
	}

	session, err := netconf.DialSSHTimeout(finalAddress, netconf.SSHConfigPassword(user, password), timeout)

	return session, err
}

func (sm SSHSessionManager) DoNetconfRequest(address string, request string) (*netconf.RPCReply, error) {

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	//first try
	result, err := sm.tryDoNetconfRequest(address, request)
	if err == nil {
		return result, nil
	}

	sm.CloseSession(address)

	//second try with error
	return sm.tryDoNetconfRequest(address, request)
}

func (sm SSHSessionManager) tryDoNetconfRequest(address string, request string) (*netconf.RPCReply, error) {
	session, err := sm.GetSession(address)
	if err != nil {
		lspLogger.Error(err, request)
		return nil, errors.New(err.Error() + "\r\n Information: " + request)
	}

	return session.Exec(netconf.RawMethod(request))
}
