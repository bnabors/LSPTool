/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package utils

import (
	"bytes"
	"strings"
	"sync"
	"time"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/config"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
	"golang.org/x/crypto/ssh"
)

type (
	// SSHSessionManager - type of session manager for pfe statistics
	SSHSessionManager struct {
		routerSSHClients map[string]*ssh.Client
		mu               sync.Mutex
	}
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
