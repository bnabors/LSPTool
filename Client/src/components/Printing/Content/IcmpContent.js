/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./print-content.css"
import React from "react";
import {IcmpTable} from "../../Results/Table/Tabs/IcmpTabContent/Content"

export default class IcmpContent extends React.Component {
    render() {
        return (
            <div className="print-block">
                <div className="block-header">ICMP</div>
                <div className="block-body">
                    <div className="hiden-buttons arial-table">
                        <IcmpTable data={this.props.data} />
                    </div>
                </div>
            </div>
        )
    }
}
