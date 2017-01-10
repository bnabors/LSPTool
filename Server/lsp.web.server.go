/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"golang.org/x/net/netutil"

	"sync"

	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/controllers"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/log"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/models"
	"github.com/Juniper/24287_WOW_LSP_GOLANG/Server/utils"
)

var (
	sessionSync     sync.Mutex
	activeSessions  int = 0
	connectionCount int
)

func StartServer(ip string, port int, maxConnectionsCount int) {

	connectionCount = maxConnectionsCount

	//main page
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/print", printHandler)

	//resources
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("site/src"))))

	//get

	//post
	http.HandleFunc("/fetchRouters", getHandleFunc(getRoutersHandler))
	http.HandleFunc("/determineLsps", getHandleFunc(determineLspsHandler))
	http.HandleFunc("/runTests", getHandleFunc(runTestsHandler))
	http.HandleFunc("/refreshLsp", getHandleFunc(refreshLspHandler))
	http.HandleFunc("/refreshRouterInfo", getHandleFunc(refreshRouterInfoHandler))
	http.HandleFunc("/refreshRouterInterface", getHandleFunc(refreshRouterInterfaceHandler))
	http.HandleFunc("/clearRefreshRouterInfo", getHandleFunc(clearRefreshRouterInfoHandler))
	http.HandleFunc("/clearRefreshRouterInterface", getHandleFunc(clearRefreshRouterInterfaceHandler))
	http.HandleFunc("/refreshPings", getHandleFunc(refreshPingsHandler))
	http.HandleFunc("/refreshPing", getHandleFunc(refreshPingHandler))

	//server
	address := ip + ":" + strconv.Itoa(port)
	startUnlimitedServer(address)
}

func startLimitedServer(address string, maxConnectionCount int) {
	listner, err := net.Listen("tcp", address)
	if err != nil {
		lspLogger.Error(err)
		return
	}

	//limit of number simultaneous connections
	listner = netutil.LimitListener(listner, maxConnectionCount)
	err = http.Serve(listner, nil)
	if err != nil {
		lspLogger.Error(err)
	}
}

func startUnlimitedServer(address string) {
	err := http.ListenAndServe(address, nil)
	if err != nil {
		lspLogger.Error(err)
	}
}

func getHandleFunc(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionSync.Lock()
		if activeSessions >= connectionCount {
			sendErrorResponse(w, "Connection Error: The number of active connections exceeds the maximum value")
			sessionSync.Unlock()
			return
		}
		activeSessions++
		sessionSync.Unlock()

		handler(w, r)

		sessionSync.Lock()
		activeSessions--
		sessionSync.Unlock()
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("main request")
	data, err := ioutil.ReadFile("site/index.html")

	if err != nil {
		lspLogger.Error(err)
		return
	}

	w.Write(data)
	lspLogger.Infoln("main response")
}
func printHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("print request")
	data, err := ioutil.ReadFile("site/print.html")

	if err != nil {
		lspLogger.Error(err)
		return
	}

	w.Write(data)
	lspLogger.Infoln("print response")
}

func getRoutersHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("getRouters request")

	res, err := controller.GetRouters()
	if err != nil {
		sendErrorResponse(w, "MySQL Error: \""+err.Error()+"\"")
		return
	}

	sendResponse(w, res)
	lspLogger.Infoln("getRouters response")
}

func determineLspsHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("determineLsps request")

	r.ParseForm()
	var t models.DetermineOptions
	for key, _ := range r.Form {
		lspLogger.Debugln(key)
		err := json.Unmarshal([]byte(key), &t)
		if err != nil {
			sendErrorResponse(w, err.Error())
			return
		}
	}

	lspLogger.Debug(t.Ingress)
	lspLogger.Debug(t.Egress)

	res, err := controller.DetermineLsps(t)
	if err != nil {
		sendErrorResponse(w, "MPLS Error: \""+err.Error()+"\"")
		return
	}

	sendResponse(w, res)
	lspLogger.Infoln("determineLsps response")
}

func runTestsHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("runTests request")
	r.ParseForm()
	lspLogger.Debug(r.Form)

	var t models.TestOptions
	for key, _ := range r.Form {
		lspLogger.Debugln(key)
		err := json.Unmarshal([]byte(key), &t)
		if err != nil {
			sendErrorResponse(w, err.Error())
			return
		}
	}
	lspLogger.Debug(t.P2P)
	lspLogger.Debug(t.P2MP)

	res, err := controller.RunTests(t)

	if err != nil {
		sendErrorResponse(w, "MPLS Error: \""+err.Error()+"\"")
		return
	}

	sendResponse(w, res)

	lspLogger.Infoln("runTests response")
}

func getRefreshRequestParams(r *http.Request, requestModel *models.RequestModel) error {
	r.ParseForm()
	lspLogger.Debug(r.Form)

	for key, _ := range r.Form {
		lspLogger.Debugln(key)
		err := json.Unmarshal([]byte(key), &requestModel)
		if err != nil {
			lspLogger.Errorln(err.Error())
			return err
		}
	}

	return nil
}

func refreshLspHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("refreshLsp request")
	var requestModel = models.RequestModel{}
	err := getRefreshRequestParams(r, &requestModel)
	if err != nil {
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(requestModel)

	var options = models.ReloadAllOptions{}
	err = json.Unmarshal([]byte(requestModel.Options), &options)
	if err != nil {
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(options.LspGroups)

	routeResult, err := controller.RefreshLsp(requestModel, options.LspGroups)
	if err != nil {
		sendErrorResponse(w, "MPLS Error: \""+err.Error()+"\"")
		return
	}

	responseModel := models.ResponseModel{LspItem: requestModel.LspItem, Result: *routeResult}

	sendResponse(w, responseModel)
	lspLogger.Infoln("refreshLsp response")
}

func refreshRouterInfoHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("refreshRouterInfo request")

	var requestModel = models.RequestModel{}
	err := getRefreshRequestParams(r, &requestModel)
	if err != nil {
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(requestModel)

	var options = models.RefreshRouterOptions{}
	err = json.Unmarshal([]byte(requestModel.Options), &options)
	if err != nil {
		lspLogger.Errorln(err.Error())
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(options)

	router, err := utils.GetRouterFromMysqlByID(options.Router)
	if err != nil {
		sendErrorResponse(w, "MySql Error: \""+err.Error()+"\"")
		return
	}
	testResult, err := controller.RefreshRouter(*router, options.Interfaces)
	if err != nil {
		sendErrorResponse(w, "MPLS Error: \""+err.Error()+"\"")
		return
	}

	sendResponse(w, models.ResponseModel{
		LspItem: requestModel.LspItem,
		Result:  testResult,
	})
	lspLogger.Infoln("refreshRouterInfo response")
}

func refreshRouterInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("refreshRouterInterface request")
	var requestModel = models.RequestModel{}
	err := getRefreshRequestParams(r, &requestModel)
	if err != nil {
		lspLogger.Errorln(err.Error())
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(requestModel)

	var options = models.RefreshRouterInterfaceOptions{}
	err = json.Unmarshal([]byte(requestModel.Options), &options)
	if err != nil {
		lspLogger.Errorln(err.Error())
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(options)

	router, err := utils.GetRouterFromMysqlByID(options.Router)
	if err != nil {
		sendErrorResponse(w, "MySql Error: \""+err.Error()+"\"")
		return
	}
	content, err := controller.RefreshInterfaceInfo(*router, options.GetInterfaceName())
	if err != nil {
		sendErrorResponse(w, "MPLS Error: \""+err.Error()+"\"")
		return
	}

	sendResponse(w, models.ResponseModel{
		LspItem: requestModel.LspItem,
		Result: models.RefreshRouterInterfaceResult{
			Router:    options.Router,
			Interface: options.Interface,
			Result:    content,
		},
	})

	lspLogger.Infoln("refreshRouterInterface response")
}

func clearRefreshRouterInfoHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("clearRefreshRouterInfo request")
	var requestModel = models.RequestModel{}
	err := getRefreshRequestParams(r, &requestModel)
	if err != nil {
		lspLogger.Errorln(err.Error())
		sendErrorResponse(w, err.Error())
		return
	}

	lspLogger.Debug(requestModel.Options)

	var options = models.RefreshRouterOptions{}
	err = json.Unmarshal([]byte(requestModel.Options), &options)
	if err != nil {
		lspLogger.Errorln(err.Error())
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(options)

	router, err := utils.GetRouterFromMysqlByID(options.Router)
	if err != nil {
		sendErrorResponse(w, "MySql Error: \""+err.Error()+"\"")
		return
	}
	testResult, err := controller.ClearAndRefreshRouter(*router, options.Interfaces)
	if err != nil {
		sendErrorResponse(w, "MPLS Error: \""+err.Error()+"\"")
		return
	}

	sendResponse(w, models.ResponseModel{
		LspItem: requestModel.LspItem,
		Result:  testResult,
	})

	lspLogger.Infoln("clearRefreshRouterInfo response")
}

func clearRefreshRouterInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("clearRefreshRouterInfo request")

	var requestModel = models.RequestModel{}
	err := getRefreshRequestParams(r, &requestModel)
	if err != nil {
		lspLogger.Errorln(err.Error())
		sendErrorResponse(w, err.Error())
		return
	}

	lspLogger.Debug(requestModel.Options)

	var options = models.RefreshRouterInterfaceOptions{}
	err = json.Unmarshal([]byte(requestModel.Options), &options)
	if err != nil {
		lspLogger.Errorln(err.Error())
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(options)

	router, err := utils.GetRouterFromMysqlByID(options.Router)
	if err != nil {
		sendErrorResponse(w, "MySql Error: \""+err.Error()+"\"")
		return
	}
	content, err := controller.ClearAndRefreshInterfaceInfo(*router, options.GetInterfaceName())
	if err != nil {
		sendErrorResponse(w, "MPLS Error: \""+err.Error()+"\"")
		return
	}

	sendResponse(w, models.ResponseModel{
		LspItem: requestModel.LspItem,
		Result: models.RefreshRouterInterfaceResult{
			Router:    options.Router,
			Interface: options.Interface,
			Result:    content,
		},
	})

	lspLogger.Infoln("clearRefreshRouterInfo response")
}

func refreshPingsHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("refreshPings request")

	var requestModel = models.RequestModel{}
	err := getRefreshRequestParams(r, &requestModel)
	if err != nil {
		lspLogger.Errorln(err.Error())
		return
	}
	lspLogger.Debug(requestModel)

	var options = models.RefreshPingsOptions{}
	err = json.Unmarshal([]byte(requestModel.Options), &options)
	if err != nil {
		sendErrorResponse(w, err.Error())
		return
	}

	testResult, err := controller.RefreshPings(options.Id, requestModel, options)
	if err != nil {
		sendErrorResponse(w, err.Error())
		return
	}
	responseModel := models.ResponseModel{
		LspItem: requestModel.LspItem,
		Result:  testResult}

	sendResponse(w, responseModel)
	lspLogger.Infoln("refreshPings response")
}

func refreshPingHandler(w http.ResponseWriter, r *http.Request) {
	lspLogger.Infoln("refreshPing request")

	var requestModel = models.RequestModel{}
	var err = getRefreshRequestParams(r, &requestModel)
	if err != nil {
		lspLogger.Errorln(err.Error())
		return
	}

	lspLogger.Debug(requestModel.Options)

	var options = models.RefreshPingOptions{}
	err = json.Unmarshal([]byte(requestModel.Options), &options)
	if err != nil {
		sendErrorResponse(w, err.Error())
		return
	}
	lspLogger.Debug(options)

	result, err := controller.RefreshPing(requestModel, options)
	if err != nil {
		sendErrorResponse(w, err.Error())
		return
	}
	responseModel := models.ResponseModel{LspItem: requestModel.LspItem, Result: models.RefreshPingResult{Id: options.Id, Result: result}}

	sendResponse(w, responseModel)
	lspLogger.Infoln("refreshPing response")
}

func sendErrorResponse(w http.ResponseWriter, errorValue string) {
	lspLogger.Errorln(errorValue)
	sendResponse(w, models.ResponseErrorModel{Error: errorValue})
}

func sendResponse(w http.ResponseWriter, data interface{}) {
	b, err := utils.ConvertToJson(&data)
	if err != nil {
		lspLogger.Error(err)
		return
	}
	lspLogger.Debug(string(b))
	w.Write(b)
}
