/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import BaseTab from "./BaseTab"
import React from "react"
import {Diagram, BandsTable, RoutersTable} from "./LspTabContent/Content"

export default class LspsTab extends BaseTab {
    clicked(routerId) {
        this.props.onClicked(routerId);
    }

    print(){
        this.props.printLsp(this.props.data);
    }

    render() {
        return (
            <div id={this.props.data.id} className="tabcontent">
                <div className="dia-content">
                    <div className="result-tab-header">
                        <div>LSPs Diagram</div>
                        <div className="right-float">
                            <button className="app-control-gradient" onClick={this.print.bind(this)}>Print</button>
                        </div>
                    </div>
                    <div className="dia-wrapper">
                        <Diagram data={this.props.data.content.diagram}/>
                    </div>
                </div>
                <div className="lps-content">
                    <RoutersTable data={this.props.data.content.routes} onClicked={this.clicked.bind(this)}/>
                    <BandsTable data={this.props.data.content.bands}/>
                </div>
            </div>
        )
    }
}