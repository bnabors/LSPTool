/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package controllerHelper

import (
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
)

type RouteChecker struct {
	NewGroups []*models.LspGroup
	OldGroups []*models.LspGroup
}

func (obj *RouteChecker) IsLspGroupChanged(testedLsps []*models.TestLspOptions, oldLsps []*models.LspItem) bool {
	for _, testLsp := range testedLsps {
		if !testLsp.IsSelected {
			continue
		}

		var oldLsp = GetLspByName(testLsp.LspItem.Name, oldLsps)

		if oldLsp == nil {
			return true
		}

		if obj.IsLspChanged(&testLsp.LspItem, oldLsp) {
			return true
		}
	}

	return false
}

func (obj *RouteChecker) IsLspChanged(newLsp *models.LspItem, oldLsp *models.LspItem) bool {
	newRro := GetRroByGroupId(obj.NewGroups, newLsp.GroupId)
	oldRro := GetRroByGroupId(obj.OldGroups, oldLsp.GroupId)

	return isDifferentRro(newRro, oldRro)
}

func isDifferentRro(rroA []string, rroB []string) bool {
	if len(rroA) != len(rroB) {
		return true
	}

	for pos := 0; pos < len(rroA); pos++ {
		if rroA[pos] != rroB[pos] {
			return true
		}
	}

	return false
}

func GetLspByName(lspName string, lsps []*models.LspItem) *models.LspItem {
	for _, lsp := range lsps {
		if lsp.Name == lspName {
			return lsp
		}
	}

	return nil
}

func GetRroByGroupId(lspGroups []*models.LspGroup, groupId string) []string {
	for _, lspGroup := range lspGroups {
		if (*lspGroup).Id == groupId {
			return (*lspGroup).ReceivedRro
		}
	}

	return []string{}
}
