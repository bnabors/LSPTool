/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react";
import Checkbox from "../Common/Checkbox"

export default class LspItem extends React.Component {
    constructor(props) {
        super(props);
        this.state = {selected: false};
    }

    selectionChanged(state) {
        this.setSelected(state);
    }

    setSelected(state) {
        this.setState({
            selected: state
        });
    }

    getState() {
        return this.state.selected;
    }

    render() {
        if (this.props.can_select === true) {
            return (
                <div className="checkbox-wrapper">
                    <Checkbox id={this.props.id} content={this.props.name} checked={this.state.selected} onClick={this.selectionChanged.bind(this)} />
                </div>
            );
        } else {
            return (
                <div className="item-unselectable" id={this.props.id}>{this.props.name}</div>
            );
        }
    }
}
