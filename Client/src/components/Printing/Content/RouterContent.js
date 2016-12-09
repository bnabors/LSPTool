/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./print-content.css"
import React from "react";
import {InterfaceContent} from "../../Results/Table/Tabs/RouterTabContent/Content"

export default class RouterContent extends React.Component {
    initializeState(tabState, data) {
        let state = tabState;
        if (!state.state.innerStates) {
            state.state.innerStates = [];
        }
        data.content.forEach(function (item, index, array) {
            let innerState = state.state.innerStates.find(function (element, indx) {
                if (element.id === item.id) {
                    return state;
                }
            });

            if (!innerState) {
                innerState = {
                    id: item.id,
                    expanded: true,
                    errorsOnly: false,
                    childs: []
                };
                state.state.innerStates.push(innerState);
            }
        });

        return state;
    }

    render() {
        let states = this.initializeState({state: {state: null}}, this.props.data);
        let content = this.props.data.content.map(function (item) {
            let itemState = states.state.innerStates.find(function (element, index) {
                if (element.id === item.id) {
                    return element;
                }
            });
            return <InterfaceContent data={item} itemState={itemState}/>
        }, this);
        return (
            <div className="print-block">
                <div className="block-header">{this.props.data.name + " Interface Statistics"}</div>
                <div className="block-body">
                    {content}
                </div>
            </div>
        )
    }
}
