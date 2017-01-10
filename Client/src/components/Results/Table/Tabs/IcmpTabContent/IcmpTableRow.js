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
            icmpInfo: this.props.data.icmpInfo,
        });
    }

    render() {
        return (
            <tr id={this.props.data.id}>
                <td>{this.props.data.icmpInfo.source.name}</td>
                <td>{this.props.data.icmpInfo.dest.name}</td>
                <td>{this.props.data.icmpInfo.sourceIp}</td>
                <td>{this.props.data.icmpInfo.destIp}</td>
                <td className={this.props.data.loss.error ? "error" : ""}>{this.props.data.loss.value}</td>
                <td className={this.props.data.avg.error ? "error" : ""}>{this.props.data.avg.value}</td>
                <td className={this.props.data.max.error ? "error" : ""}>{this.props.data.max.value}</td>
                <td className={this.props.data.stdDev.error ? "error" : ""}>{this.props.data.stdDev.value}</td>
                <td>
                    <div className="router-result-buttons">
                        <a id={this.props.data.id} onClick={this.refreshRow.bind(this)}>Refresh</a>
                    </div>
                </td>
            </tr>
        );
    }
}
