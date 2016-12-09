/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package utils

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"

	"../config"
	"../log"
	"../models"

	"errors"

	_ "github.com/go-sql-driver/mysql"
)

func ReadFromFile(fileName string) []byte {
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		lspLogger.Errorln(err.Error())
	}

	return bytes
}

func WriteToFile(fileName string, bytes []byte) {
	file, err := os.Create(fileName)

	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.Write(bytes)
	writer.Flush()
}

func WriteToFileAsXML(fileName string, data interface{}) {
	xmlOut, err := xml.MarshalIndent(data, "", "    ")

	if err != nil {
		lspLogger.Errorf("error: %v", err)
	}

	WriteToFile(fileName, xmlOut)
}

func ConvertToJson(data interface{}) (result []byte, err error) {
	result, err = json.Marshal(&data)

	if err != nil {
		lspLogger.Error(err)
		return
	}
	lspLogger.Debugln(string(result))

	return
}

// GetPhysicalName returns physical name of interface by logical name
func GetPhysicalName(logicalName string) string {
	var lastIndex = strings.LastIndex(logicalName, ".")

	if lastIndex == -1 {
		return logicalName
	}

	return logicalName[:lastIndex]
}

/*main page*/
func GetRoutersFromMysql() (result []*models.Router, err error) {
	result = []*models.Router{}

	db, err := sql.Open("mysql", config.LspConfig.MysqlConnectionsString)
	if err != nil {
		lspLogger.Error(err)
		return
	}

	defer db.Close()

	rows, err := db.Query(config.LspConfig.MysqlRouterQuery)
	if err != nil {
		lspLogger.Error(err)
		return
	}

	var (
		id       int
		name     string
		ip4      string
		puttyIP4 string
	)

	for {
		if rows.Next() {
			rowErr := rows.Scan(&id, &name, &ip4, &puttyIP4)
			if rowErr != nil {
				return result, rowErr
			}
			result = append(result, &models.Router{
				Id:      id,
				Name:    name,
				Ip:      ip4,
				PuttyIp: puttyIP4,
			})
		} else {
			break
		}
	}
	return
}

func GetRoutersFromMysqlByID(routersIdList []string) (result []*models.Router, err error) {
	result = []*models.Router{}

	db, err := sql.Open("mysql", config.LspConfig.MysqlConnectionsString)
	if err != nil {
		lspLogger.Error(err)
		return
	}

	defer db.Close()

	var routersList = ""

	for _, routerId := range routersIdList {
		if routersList != "" {
			routersList += ", "
		}

		routersList += routerId
	}

	rows, err := db.Query(config.LspConfig.MysqlRouterQuery + " WHERE Id IN (" + routersList + ")")
	if err != nil {
		lspLogger.Error(err)
		return
	}

	var (
		id       int
		name     string
		ip4      string
		puttyIP4 string
	)

	for {
		if rows.Next() {
			rowErr := rows.Scan(&id, &name, &ip4, &puttyIP4)
			if rowErr != nil {
				return result, rowErr
			}
			result = append(result, &models.Router{
				Id:      id,
				Name:    name,
				Ip:      ip4,
				PuttyIp: puttyIP4,
			})
		} else {
			break
		}
	}
	return
}

func GetRouterFromMysqlByID(idRouter string) (*models.Router, error) {
	db, err := sql.Open("mysql", config.LspConfig.MysqlConnectionsString)
	if err != nil {
		lspLogger.Error(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(config.LspConfig.MysqlRouterQuery + " WHERE Id=" + idRouter)
	if err != nil {
		lspLogger.Error(err)
		return nil, err
	}

	var (
		id       int
		name     string
		ip4      string
		puttyIP4 string
	)

	if rows.Next() {
		rowErr := rows.Scan(&id, &name, &ip4, &puttyIP4)
		if rowErr != nil {
			return nil, rowErr
		}
		return &models.Router{
			Id:      id,
			Name:    name,
			Ip:      ip4,
			PuttyIp: puttyIP4,
		}, nil
	}
	return nil, errors.New("Router does not exists")
}
