/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react"
import IcmpTableRow from "./IcmpTableRow"
import * as ReactDOM from "react/lib/ReactDOM";

export default class IcmpTable extends React.Component {
    refreshRow(options) {
        this.props.onRefreshPing(options);
    }

    refreshAll() {
        this.props.onRefreshPings({id: this.props.data.id});
    }

    componentDidMount(){
        /* lazy load
        if(this.props.data.content == null || this.props.data.content.length < 1){
            this.refreshAll();
        }
        */
        let table = ReactDOM.findDOMNode(this.refs.table);
        if(table){
            table.setAttribute("cellspacing", "0px");
        }
    }

    render() {
        let rows = (<div></div>);
        if (this.props.data.content != null) {
            rows = this.props.data.content.map(function (item) {
                return (<IcmpTableRow data={item} onRefresh={this.refreshRow.bind(this)}/>)
            }, this)
        }

        return (
            <table ref="table" className="table-content">
                <thead>
                <tr>
                    <th>From Device</th>
                    <th>Destination Device</th>
                    <th>Destination IP</th>
                    <th>Loss</th>
                    <th>Average</th>
                    <th>Max</th>
                    <th>SDT Dev</th>
                    <th className="router-result-buttons">
                        <div className="router-result-buttons">
                            <a id={this.props.data.id} onClick={this.refreshAll.bind(this)}>Refresh All</a>
                        </div>
                    </th>
                </tr>
                </thead>
                <tbody>
                {rows}
                </tbody>
            </table>
        )
    }
}
