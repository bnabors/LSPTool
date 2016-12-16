/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./checkbox.css"
import React from "react"

export default class Checkbox extends React.Component{
    onClick(){
        let element = this.refs.checkbox;
        if(element){
            let state = !element.checked;
            this.props.onClick(state);
        }
    }

    render(){
        return(
            <div className="checkbox-container">
                <input ref="checkbox" type="checkbox" value="None" id={this.props.id} name="check" checked={this.props.checked} />
                <label className="innerCheckbox" htmlFor={this.props.id} onClick={this.onClick.bind(this)}/>
                <div className="checkbox-text-container" htmlFor={this.props.id} onClick={this.onClick.bind(this)} title={this.props.content}>
                    {this.props.content}
                </div>
            </div>
        )
    }
}