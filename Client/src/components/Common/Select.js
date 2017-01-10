/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./select.css"
import "react-select/dist/react-select.css"

import React from "react"
import $ from "jquery"
import ReactSelect from "react-select"

export class Select extends React.Component{

    resetSelection(){
        let select = "#"+this.props.id+ " option";
        $(select).prop('selected', function() {
            return this.defaultSelected;
        });
    }

    render(){
        return(
            <div className="custom-select" disabled={this.props.disabled}>
                <select id={this.props.id} className="app-control-gradient"
                        value={this.props.value} onChange={this.props.onChange}
                        disabled={this.props.disabled} >
                    <option id={-1} value={null} selected="selected">Select Router</option>
                    {this.props.content}
                </select>
                <label />
            </div>
        )
    }
}

export class FilterableSelect extends React.Component{

    resetSelection(){

    }

    render(){
        let id = this.props.value && this.props.value != null ? this.props.value.id : -1;
        return(
            <div className="custom-select" disabled={this.props.disabled}>
                <ReactSelect id={this.props.id}
                             placeholder={"Select Router"}
                             searchable={true}
                             clearable={false}
                             value={id}
                             options={this.props.content}
                             onChange={this.props.onChange}
                             disabled={this.props.disabled}>
                </ReactSelect>
            </div>
        )
    }
}
