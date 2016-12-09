/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react";
import * as ReactDOM from "react/lib/ReactDOM";

export default class HeaderItem extends React.Component {

    clickHandler() {
        this.props.onClicked(this.props.data);
    }

    componentWillReceiveProps(newProps) {
        let element = ReactDOM.findDOMNode(this.refs.test);
        if (element) {
            let id = this.props.id;
            if (this.props.id !== newProps.id) {
                element.className = element.className.replace(" active", "");
                id = newProps.id;
            }
            if (id === newProps.selectedId) {
                if (element.className.indexOf("active") < 0) {
                    element.className += " active";
                }
            } else {
                element.className = element.className.replace(" active", "");
            }
        }
    }

    render() {
        let className = "tablinks" + (this.props.data.isError ? " error" : " no-error");
        return (
            <li>
                <a ref="test" className={className}
                   onClick={this.clickHandler.bind(this)} title={this.props.data.name}>
                    <span className={this.props.data.isError ? " error" : " no-error"}>{this.props.data.name}</span>
                    <div/>
                </a>
            </li>
        )
    }
}
