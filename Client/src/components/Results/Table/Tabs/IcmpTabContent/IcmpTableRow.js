/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react"

export default class IcmpTableRow extends React.Component {
    refreshRow() {
        this.props.onRefresh({
            contentId: this.props.data.id,
            routerStartId: this.props.data.routerStartId,
            routerFinishId: this.props.data.routerFinishId
        });
    }

    render() {
        let className = this.props.data.isError ? "error" : "";
        return (
            <tr id={this.props.data.id}>
                <td className={className}>{this.props.data.fromDevice}</td>
                <td className={className}>{this.props.data.destDevice}</td>
                <td className={className}>{this.props.data.destIp}</td>
                <td className={className}>{this.props.data.loss}</td>
                <td className={className}>{this.props.data.avg}</td>
                <td className={className}>{this.props.data.max}</td>
                <td className={className}>{this.props.data.stdDev}</td>
                <td>
                    <div className="router-result-buttons">
                        <a id={this.props.data.id} onClick={this.refreshRow.bind(this)}>Refresh</a>
                    </div>
                </td>
            </tr>
        );
    }
}
