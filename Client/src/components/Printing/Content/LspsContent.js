/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./print-content.css"
import React from "react";
import {Diagram, BandsTable, RoutersTable} from "../../Results/Table/Tabs/LspTabContent/Content"

export default class LspsContent extends React.Component {
    render() {
        return (
            <div className="print-block">
                <div className="block-header">LSPs Diagram</div>
                <div className="block-body">
                    <Diagram data={this.props.data.content.diagram}/>
                </div>
                <div className="block-header">LSPs</div>
                <div className="block-body">
                    <div className="inline arial-table">
                        <RoutersTable data={this.props.data.content.routes}/>
                        <BandsTable data={this.props.data.content.bands}/>
                    </div>
                </div>
            </div>
        )
    }
}