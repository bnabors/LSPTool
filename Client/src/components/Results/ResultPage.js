/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react";
import {connect} from "react-redux";
import {
    refreshAll,
    refreshRouter,
    refreshInterface,
    clearRefreshRouter,
    clearRefreshInterface,
    refreshPings,
    refreshPing,
    setSate,
    sendError,
    openPrintDialog,
    printContent
} from "../../actions/bundleActions"
import Table from "./Table/Table"
//import Loader from "../Common/Loader"
import CssLoader from "../Common/CssLoader"

@connect((store) => {
    return {
        result: store.currentResult,
        LSPs: store.LSPs,
        results: store.results,
        states: store.states,
        errors: store.errors
    }
})

export default class ResultPage extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            commands: {
                reloadAll: this.reloadAll.bind(this),
                printSummary: this.printSummary.bind(this),
                refreshRouter: this.refreshRouter.bind(this),
                clearRefreshRouter: this.clearRefreshRouter.bind(this),
                refreshRouterInterface: this.refreshRouterInterface.bind(this),
                clearRefreshRouterInterface: this.clearRefreshRouterInterface.bind(this),
                refreshPings: this.refreshPings.bind(this),
                refreshPing: this.refreshPing.bind(this),
                printLsp: this.printLsp.bind(this),
            }
        }
    }

    reloadAll() {
        this.sendRequest({lspGroups: this.props.LSPs.lspGroups}, refreshAll);
    }

    printSummary() {
        this.props.dispatch(openPrintDialog(this.props.result.currentResult));
    }

    refreshRouter(options) {
        this.sendRequest(options, refreshRouter);
    }

    clearRefreshRouter(options) {
        this.sendRequest(options, clearRefreshRouter);
    }

    refreshRouterInterface(options) {
        this.sendRequest(options, refreshInterface);
    }

    clearRefreshRouterInterface(options) {
        this.sendRequest(options, clearRefreshInterface);
    }

    refreshPings(options) {
        this.sendRequest(options, refreshPings);
    }

    refreshPing(options) {
        this.sendRequest(options, refreshPing);
    }

    sendRequest(options, action) {
        let data = {
            lsp: this.props.result.currentResult.lsp,
            ingress: this.props.LSPs.routers.ingress,
            egress: this.props.LSPs.routers.egress,
            names:this.props.result.currentResult.names,
            groupRouters: this.props.results.testResults.groupRouters,
            options: JSON.stringify(options).toString()
        };
        console.log(JSON.stringify(data));
        if (action && typeof(action) === "function") {
            this.props.dispatch(action(data))
        }
    }

    printLsp(data){
        this.props.dispatch(printContent(data))
    }

    onStateChanged(state) {
        this.props.dispatch(setSate(state))
    }

    componentWillReceiveProps(newProps) {
        let resultState = null;
        if (this.props.result !== newProps.result) {
            if (this.props.result.currentResult !== null) {
                let findId = this.props.result.currentResult.lsp.id;
                resultState = this.props.states.states.find(function (item, index) {
                    if (item.lspId === findId) {
                        return item;
                    }
                });
            } else if (newProps.states.activeState !== null) {
                if (newProps.result.currentResult !== null) {
                    let findId = newProps.result.currentResult.lsp.id;
                    resultState = this.props.states.states.find(function (item, index) {
                        if (item.lspId === findId) {
                            return item;
                        }
                    });
                }
            }
        }
        if (resultState) {
            if (this.refs.table) {
                this.refs.table.updateState();
            }
            this.setState({states: resultState});
        }
    }


    render() {
        if (this.props.result.currentResult == null) {
            return (<div className="no-results">Select LSPs and run tests</div>)
        }

        return (
            <div className="app-result">
                <Table ref="table"
                       data={this.props.result.currentResult} states={this.props.states.activeState}
                       onStateChanged={this.onStateChanged.bind(this)}
                       commands={this.state.commands}
                />
                <div className={this.props.result.fetching ? "large-loader" : "hidden"}>
                    <CssLoader cwidth="101" cheight="101"/>
                </div>
            </div>
        )
    }
}