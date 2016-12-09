/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import BaseTab from "./BaseTab"
import React from "react"
import {IcmpTable} from "./IcmpTabContent/Content"

export default class IcmpTab extends BaseTab {
    refreshRow(options) {
        this.props.commands.refreshPing({
            id: this.props.data.id,
            contentId: options.contentId,
            routerStartId: options.routerStartId,
            routerFinishId: options.routerFinishId
        });
    }

    refreshAll() {
        this.props.commands.refreshPings({id: this.props.data.id});
    }

    render() {
        return (
            <div id={this.props.data.id} className="tabcontent">
                <div className="result-tab-header">ICMP</div>
                <div className="icmp-tab">
                    <IcmpTable data={this.props.data} onRefreshPing={this.refreshRow.bind(this)}
                               onRefreshPings={this.refreshAll.bind(this)}/>
                </div>
            </div>
        )
    }
}
