/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controller

import (
	"strings"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/config"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/models"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/utils"
)

func GetRouters() (*models.RouterContainer, error) {
	routers, err := utils.GetRoutersFromMysql()
	if err != nil {
		lspLogger.Error(err)
		return nil, err
	}
	ingress := []*models.Router{}
	egress := []*models.Router{}
	for _, router := range routers {
		if checkIngress(*router) {
			ingress = append(ingress, router)
		}
		if checkEgress(*router) {
			egress = append(egress, router)
		}
	}

	return &models.RouterContainer{IngressRouters: ingress, EgressRouters: egress}, nil
}

func checkIngress(router models.Router) bool {
	for _, name := range config.LspConfig.IngressRouterNames {
		if strings.Contains(router.Name, name) {
			return true
		}
	}
	return false
}
func checkEgress(router models.Router) bool {
	for _, name := range config.LspConfig.EgressRouterNames {
		if strings.Contains(router.Name, name) {
			return true
		}
	}
	return false
}
