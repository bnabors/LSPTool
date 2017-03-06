/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./main.css"
import React from "react";
import RouterList from "./RouterList/RouterList";
import LspPanel from "./Lsp/LspPanel"
import ResultPage from "./Results/ResultPage"
import RoutesPanel from "./Results/RoutesPanel"
import PrintingPage from "./Printing/PrintingPage"
import CssLoader from "./Common/CssLoader"
import MessageWindow from "./Common/MessageWindow"
import Copyright from "./Copyright/Copyright"
import {showResult, fetchRouters, determineLsps, runTests, sendInfo} from "../actions/bundleActions"
import {connect} from "react-redux";

@connect((store) => {
    return {
        routers: store.routers,
        LSPs: store.LSPs,
        results: store.results,
        errors: store.errors
    }
})

export default class MainPage extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            result: false,
        }
    }

    printSummary() {

    }

    componentWillMount() {
        this.props.dispatch(fetchRouters());
    }

    determineLsps(data) {
        this.props.dispatch(determineLsps(data));
    }

    detectSelection(selected) {
        let hasSelected = false;
        let selectedItem = selected.p2p.find(function (item) {
            if (item.selected) {
                return item;
            }
        });
        hasSelected = hasSelected || selectedItem !== undefined;
        selectedItem = selected.p2mp.find(function (item) {
            if (item.selected) {
                return item;
            }
        });
        return hasSelected || selectedItem !== undefined;
    }

    runTests() {
        let selected = this.refs.lspPanel.getSelected();
        if (!this.detectSelection(selected)) {
            this.props.dispatch(sendInfo("Select LSP(s) and run tests"));
            return
        }
        this.props.dispatch(runTests(selected));

        this.setState({
            result: true
        })
    }

    selectOtherLSPs() {
        this.setState({
            result: false
        });
        this.props.dispatch(showResult(null));
    }

    onTabSelected(route) {
        this.props.dispatch(showResult(route));
    }

    sendMessage(message) {
        this.props.dispatch(sendInfo(message));
    }

    render() {
        let actionButton = (<div></div>);
        let contentPanel = (<div></div>);
        if (!this.state.result) {
            if (this.props.LSPs.fetched) {
                contentPanel = (<LspPanel ref="lspPanel" p2pLSPs={this.props.LSPs.p2pLSPs}
                                          p2mpLSPs={this.props.LSPs.p2mpLSPs} downLSPs={this.props.LSPs.downLSPs}
                                          lspGroups={this.props.LSPs.lspGroups}
                                          ingress={this.props.LSPs.routers.ingress}
                                          egress={this.props.LSPs.routers.egress}/>);
                actionButton = (<button className="app-control-gradient" onClick={this.runTests.bind(this)}
                                        title="Run tests on selected LSPs">Run tests on selected LSPs</button>);
            } else if (this.props.LSPs.fetching) {
                contentPanel = (<CssLoader cwidth="49" cheight="49"/>);
            }
        } else {
            contentPanel = (
                <RoutesPanel data={this.props.results.testResults} onSelected={this.onTabSelected.bind(this)}
                             sendMessage={this.sendMessage.bind(this)}/>);
            actionButton = (
                <button className="app-control-gradient" onClick={this.selectOtherLSPs.bind(this)}
                        disabled={this.props.results.fetching} title="Select other LSPs">Select other
                    LSPs</button>);
        }

        return (
            <div className="app-div">
                <div className="app-panel app-main-panel">
                    <div className="right-border"></div>
                    <div className="container">
                        <RouterList determineLsps={this.determineLsps.bind(this)}
                                    sendMessage={this.sendMessage.bind(this)}
                                    ingress_routers={this.props.routers.ingress_routers}
                                    egress_routers={this.props.routers.egress_routers}
                                    isEnabled={!this.state.result && !this.props.LSPs.fetching}/>
                        <div className="scrollable">
                            {contentPanel}
                        </div>
                        <div className="footer">
                            <div className="separator bottom"></div>
                            <div className="bottom">
                                {actionButton}
                            </div>
                        </div>
                    </div>
                </div>
                <div className="app-result-lsp app-main-content">
                    <ResultPage />
                </div>
                <MessageWindow />
                <PrintingPage />
            </div>
        )
    }
}
