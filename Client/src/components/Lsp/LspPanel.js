/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "../main.css"
import "./main.css"
import React from "react";
import LspList from "./LspList";

export default class LspPanel extends React.Component {

    constructor(props){
        super(props);
    }

    getSelected() {
        return {
            ingress: this.props.ingress,
            egress: this.props.egress,
            p2p: this.refs.p2p.getSelected(),
            p2mp: this.refs.p2mp.getSelected(),
            lspGroup: this.props.lspGroups
        };
    }

    render(){
        return (
            <div className="lsp-panel">
                <div>
                    <div className="routers-header">Select LSPs</div>
                </div>
                <div>
                    <LspList ref="p2p" data={this.props.p2pLSPs} can_select={true} name={"P2P LSPs"} />
                    <LspList ref="p2mp" data={this.props.p2mpLSPs} can_select={true} name={"P2MP LSPs"} />
                    <LspList data={this.props.downLSPs} can_select={false} name={"LSPs that are down"} />
                </div>
            </div>
        )
    }
}