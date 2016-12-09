/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import BaseTab from "./BaseTab"
import {InterfaceContent} from "./RouterTabContent/Content"
import React from "react"

export default class RouterTab extends BaseTab {
    constructor(props) {
        super(props);
        this.state = {
            state: this.initializeState(props.tabState, props.data)
        }
    }

    getInterfaces(){
        let interfaces = [];
        if(this.props.data.content !== undefined && this.props.data.content != null){
            for (let i =0; i < this.props.data.content.length; i++){
                interfaces.push(this.props.data.content[i].id);
            }
        }
        return interfaces;
    }

    refreshAll() {
        this.props.commands.refreshRouter({router: this.props.data.id, interfaces: this.getInterfaces()});
    }

    clearRefreshAll() {
        this.props.commands.clearRefreshRouter({router: this.props.data.id, interfaces: this.getInterfaces()});
    }

    refreshRow(options) {
        this.props.commands.refreshRouterInterface({router: this.props.data.id, interface: options});
    }

    clearRefreshRow(options) {
        this.props.commands.clearRefreshRouterInterface({router: this.props.data.id, interface: options});
    }

    expandFailed() {
        this.expandChid(this.props.data.content, true);
    }

    expandAll() {
        this.expandChid(this.props.data.content, false);
    }

    expandChid(items, errorsOnly) {
        if (items != null && items.length > 0) {
            for (let i = 0; i < items.length; i++) {
                this.refs[items[i].id].expandInterface(errorsOnly);
            }
        }
    }

    initializeState(tabState, data) {
        let state = tabState;
        if (!state.state.innerStates) {
            state.state.innerStates = [];
        }
        data.content.forEach(function (item, index, array) {
            let innerState = state.state.innerStates.find(function (element, indx) {
                if (element.id === item.id) {
                    return state;
                }
            });

            if (!innerState) {
                innerState = {
                    id: item.id,
                    expanded: false,
                    errorsOnly: false,
                    childs: []
                };
                state.state.innerStates.push(innerState);
            }
        });

        return state;
    }

    componentWillReceiveProps(newProps) {
        let state = newProps.tabState;
        state = this.initializeState(state, newProps.data);
        this.setState({state: state});
    }

    getState() {
        let state = this.state.state;
        if (state.state.innerStates != null && state.state.innerStates.length > 0) {
            for (let i = 0; i < state.state.innerStates.length; i++) {
                let element = this.refs[state.state.innerStates[i].id];
                if (element) {
                    state.state.innerStates[i] = element.getState();
                }
            }
        }
        return state;
    }

    render() {
        let state = this.state.state;
        let content = this.props.data.content.map(function (item) {
            let itemState = state.state.innerStates.find(function (element, index) {
                if (element.id === item.id) {
                    return element;
                }
            });
            return <InterfaceContent ref={item.id} data={item}
                                     itemState={itemState}
                                     onRefresh={this.refreshRow.bind(this)}
                                     onClearRefresh={this.clearRefreshRow.bind(this)}/>
        }, this);

        return (
            <div id={this.props.data.id} className="tabcontent">
                <div className="result-tab-header">Interface Statistics</div>
                <div className="stat-table-header">
                    <div></div>
                    <div className="router-result-buttons">
                        <div><a onClick={this.expandFailed.bind(this)}>Expand Failed</a></div>
                        <div><a onClick={this.expandAll.bind(this)}>Expand All</a></div>
                        <div><a onClick={this.refreshAll.bind(this)}>Refresh All</a></div>
                        <div><a onClick={this.clearRefreshAll.bind(this)}>Clear/Refresh All</a></div>
                    </div>
                </div>
                <div className="stat-root">
                    {content}
                </div>
            </div>
        )
    }
}