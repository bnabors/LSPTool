/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./print-content.css";
import React from "react";
import {AppLogo, AppHeader} from "../../Common/Logo";
import LspsContent from "./LspsContent";
import {ResultType} from "../../../constants/resultType";
import RouterContent from "./RouterContent";
import IcmpContent from "./IcmpContent";

export default class PrintContent extends React.Component {

    buildBlock(item) {
        switch (item.type) {
            case ResultType.LSPs:
                return (<LspsContent data={item}/>);
            case ResultType.ROUTER:
                return (
                    <RouterContent data={item}/>
                );
            case ResultType.ICPM:
                return (
                    <IcmpContent data={item} onRefreshPings={this.props.onRefreshPings}/>
                );
            default:
                return (<div>ERROR!</div>);
        }
    }

    render() {
        let blocks = this.buildBlock(this.props.content);

        return (
            <div className="print-content">
                <div className="print-header">
                    <AppLogo fill="#000000"/>
                    <AppHeader/>
                </div>
                <div className="print-body">
                    <div className="print-block">
                        <div className="block-body">
                            <div className="content-block">
                                <div>Ingress Router</div>
                                <div>{this.props.routers.ingress.name}</div>
                            </div>
                            <div className="content-block">
                                <div>Egress Router</div>
                                <div>{this.props.routers.egress.name}</div>
                            </div>
                        </div>
                    </div>
                    {blocks}
                </div>
            </div>
        )
    }
}