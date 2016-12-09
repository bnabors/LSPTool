/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react"
import * as ReactDOM from "react/lib/ReactDOM";

export default class RoutersTable extends React.Component {
    clicked(routeId) {
        this.props.onClicked(routeId)
    }

    componentDidMount(){
        let table = ReactDOM.findDOMNode(this.refs.table);
        if(table){
            table.setAttribute("cellspacing", "0px");
        }
    }

    render() {
        let rows = this.props.data.map(function (item) {
            return (
                <tr>
                    <td>{item.type}</td>
                    <td><RouterRow data={item.router} onClicked={this.clicked.bind(this)}/></td>
                </tr>
            );
        }, this);

        return (
            <div className="lsps-tab-routers">
                <div className="table-header-div">Routers</div>
                <table ref="table" className="table-content">
                    <thead>
                    <tr>
                        <th>Router Type</th>
                        <th>Name</th>
                    </tr>
                    </thead>
                    <tbody>
                    {rows}
                    </tbody>
                </table>
            </div>
        );
    }
}

class RouterRow extends React.Component {
    clicked() {
        this.props.onClicked(this.props.data.id);
    }

    render() {
        return (<a onClick={this.clicked.bind(this)}>{this.props.data.name}</a>)
    }
}
