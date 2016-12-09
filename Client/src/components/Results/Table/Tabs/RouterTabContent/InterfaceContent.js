/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react"
import {SwitchButton} from "../../../../Common/Buttons"
import * as ReactDOM from "react/lib/ReactDOM";
import StatisticsTable from "./StatisticsTable"

export default class InterfaceContent extends React.Component {
    constructor(props) {
        super(props);
        let state = this.initializeState(props.itemState, props.data);
        let expanded = state ? state.expanded : false;
        let errorsOnly = state ? state.errorsOnly : false;
        this.state = {
            state: state,
            expanded: expanded,
            errorsOnly: errorsOnly
        }
    }

    initializeState(state, data) {
        if (data.sub_interfaces == null || data.sub_interfaces.length < 1) {
            return state;
        }

        if (!state.childs) {
            state.childs = [];
        }

        data.sub_interfaces.forEach(function (item, index, array) {
            let childState = state.childs.find(function (element, indx) {
                if (element.id === item.id) {
                    return state;
                }
            });

            if (!childState) {
                childState = {
                    id: item.id,
                    expanded: state.expanded,
                    errorsOnly: state.errorsOnly,
                    childs: []
                };
                state.childs.push(childState);
            }
        });

        return state;
    }

    getState() {
        let state = this.state.state;
        state.expanded = this.state.expanded;
        state.errorsOnly = this.state.errorsOnly;

        if (state.childs != null && state.childs.length > 0) {
            for (let i = 0; i < state.childs.length; i++) {
                let element = this.refs[state.childs[i].id];
                if (element) {
                    state.childs[i] = element.getState();
                }
            }
        }
        return state;
    }

    expandInterface(errorsOnly) {
        this.setState({
            expanded: true,
            errorsOnly: errorsOnly
        });

        if (this.props.data.sub_interfaces != null) {
            let items = this.props.data.sub_interfaces;
            if (items != null && items.length > 0) {
                for (let i = 0; i < items.length; i++) {
                    let element = this.refs[items[i].id];
                    if (element) {
                        element.expandInterface(errorsOnly);
                    }
                }
            }
        }

        this.updateTable(errorsOnly);
    }

    openCloseStats() {
        let state = !this.state.expanded;
        this.setState({expanded: state});
    }

    filterFailed() {
        this.setState({errorsOnly: true});
        this.updateTable(true);
    }

    filterAll() {
        this.setState({errorsOnly: false});
        this.updateTable(false);
    }

    updateTable(state) {
        /*
         let element = ReactDOM.findDOMNode(this.refs.tbl);
         if (element) {
         if (state === true) {
         element.setAttribute("errorsOnly", "");
         } else {
         element.removeAttribute("errorsOnly");
         }
         }
         */
    }

    refresh() {
        this.refreshInnerInterface(null);
    }

    clearRefresh() {
        this.clearRefreshInnerInterface(null);
    }

    componentWillReceiveProps(newProps) {
        if (newProps.itemState !== this.props.itemState) {
            let state = this.initializeState(newProps.itemState, newProps.data);
            let expanded = state ? state.expanded : false;
            let errorsOnly = state ? state.errorsOnly : false;
            this.setState({
                state: state,
                expanded: expanded,
                errorsOnly: errorsOnly
            });
        }
    }

    refreshInnerInterface(options) {
        let data = {interface: this.props.data.id, innerInterface: options};
        if (this.props.onRefresh) {
            this.props.onRefresh(data);
        }
    }

    clearRefreshInnerInterface(options) {
        let data = {interface: this.props.data.id, innerInterface: options};
        if (this.props.onClearRefresh) {
            this.props.onClearRefresh(data);
        }
    }

    findError(stat) {
        for (let i = 0; i < stat.values.length; i++) {
            let item = stat.values[i];
            if (item.isError) {
                return true;
            }
        }
        return false;
    }

    render() {
        let stats = [];
        if (this.props.data.statistics != null && this.props.data.statistics.length > 0) {
            /*
             stats = this.props.data.statistics.map(function (item) {
             return <StatisticsTable data={item} errorsOnly={this.state.errorsOnly}/>
             }, this);
             */

            for (let i = 0; i < this.props.data.statistics.length; i++) {
                stats.push(<StatisticsTable data={this.props.data.statistics[i]} errorsOnly={this.state.errorsOnly}/>);
                if (i + 1 < this.props.data.statistics.length - 1) {
                    let className = this.findError(this.props.data.statistics[i]) ? "error" : "no-error";
                    stats.push(
                        <tbody className={className}>
                        <tr>
                            <td>
                                <div className="table-separator"/>
                            </td>
                        </tr>
                        </tbody>
                    );
                }
            }
        }

        let innerInterfaces = [];
        if (this.props.data.sub_interfaces != null && this.props.data.sub_interfaces.length > 0) {
            innerInterfaces = this.props.data.sub_interfaces.map(function (item) {
                let childState = this.state.state.childs.find(function (element, indx) {
                    if (element.id === item.id) {
                        return element;
                    }
                });
                return <InterfaceContent ref={item.id} data={item} is_inner="true"
                                         itemState={childState}
                                         onRefresh={this.refreshInnerInterface.bind(this)}
                                         onClearRefresh={this.clearRefreshInnerInterface.bind(this)}/>
            }, this);
        }

        let className = "stat-interface";
        className += this.props.is_inner ? " stat-inner" : "";
        className += this.state.expanded ? " expanded" : "";

        return (
            <div className={className}>
                <div className="stat-header">
                    <div className="plus-minus-button">
                        <SwitchButton isSwitch={this.state.expanded}
                                      onClick={this.openCloseStats.bind(this)}/>
                    </div>
                    <div className={this.props.data.isError ? "error" : "no-error"}>
                        <div
                            className={this.props.data.isError ? "name-wrap error" : "name-wrap no-error"}>{this.props.data.name}</div>
                    </div>
                    <div className="router-result-buttons">
                        <div>
                            <a className={this.state.errorsOnly ? "selected switch-btn" : "switch-btn"}
                               id={this.props.data.id}
                               onClick={this.filterFailed.bind(this)}>Failed</a>
                        </div>
                        <div>
                            <a className={this.state.errorsOnly ? "switch-btn" : "selected switch-btn"}
                               id={this.props.data.id}
                               onClick={this.filterAll.bind(this)}>All</a>
                        </div>
                        <div className="space30"/>
                        <div>
                            <a id={this.props.data.id} onClick={this.refresh.bind(this)}>Refresh</a>
                        </div>
                        <div>
                            <a id={this.props.data.id} onClick={this.clearRefresh.bind(this)}>Clear/Refresh</a>
                        </div>
                    </div>
                </div>
                <div ref="content" className={this.state.expanded ? "stat-content expanded" : "stat-content hidden"}>
                    <div className="stat-content-table">
                        <table className={"router-result-stats-table" + (this.state.errorsOnly ? " errorsOnly" : "")}
                               ref="tbl">
                            {stats}
                        </table>
                    </div>
                    {innerInterfaces}
                </div>
            </div>
        )
    }
}
