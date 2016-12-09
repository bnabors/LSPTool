/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react";
import LspItem from "./LspItem";
import Checkbox from "../Common/Checkbox"

export default class LspList extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            allSelected: false
        }
    }

    getSelected() {
        let selected = [];
        for (let i = 0; i < this.props.data.length; i++) {
            let item = this.props.data[i];
            selected.push({lsp: item, selected: this.refs[item.id].getState()});
        }
        return selected;
    }

    setChildState(state) {
        for (let i = 0; i < this.props.data.length; i++) {
            this.refs[this.props.data[i].id].setSelected(state);
        }
    }

    selectAll() {
        let value = !this.state.allSelected;
        this.setChildState(value);
        this.setState({
            allSelected: value
        })
    }

    itemSelected() {
        let selectedCount = this.getSelected().length
        if (selectedCount === 0) {
            this.setState({
                allSelected: false
            });
        } else if (selectedCount === this.props.data.length) {
            this.setState({
                allSelected: false
            });
        }
    }

    render() {

        if (this.props.data === undefined) {
            debugger;
            return (<div>error</div>)
        }

        let result = this.props.data.map(function (lsp) {
            return (
                <LspItem name={lsp.name} id={lsp.id} can_select={this.props.can_select}
                         selectionChanged={this.itemSelected.bind(this)}
                         ref={lsp.id}>
                    {lsp.name}
                </LspItem>
            )
        }, this);

        let header = (<div></div>);
        if (this.props.can_select) {
            header = (
                <div className="lsp-header">
                    <Checkbox id={this.props.id} onClick={this.selectAll.bind(this)}
                              checked={this.state.allSelected} content={this.props.name}/>
                </div>
            )
        } else {
            header = (
                <div className="lsp-header-unselectable">
                    <label>{this.props.name}</label>
                </div>
            )
        }

        return (
            <div className="lsp-list">
                {header}
                <div className="lsp-body">
                    {result}
                </div>
            </div>
        )
    }
}
