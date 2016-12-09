/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import React from "react"
import {clearMessage} from "../../actions/bundleActions"
//import * as ReactDOM from "react/lib/ReactDOM";
import {connect} from "react-redux";
import $ from "jquery"

const ERROR_HEADER = "Error";
const INFO_HEADER = "Information";
const COPY_HEADER = "Copyright";

@connect((store) => {
    return {
        errors: store.errors
    }
})

export default class MessageWindow extends React.Component {

    componentDidMount() {
        let self = this;
        $(document).keydown(function (e) {
            if (e.which === 27) self.onClose();
        })
    }

    onClose() {
        this.props.dispatch(clearMessage());
    }

    getMessage() {
        if (this.props.errors.hasError) {
            return <p>{this.props.errors.error ? this.props.errors.error.toString() : ""}</p>;
        }
        if (this.props.errors.hasInfo) {
            return <p>{this.props.errors.info ? this.props.errors.info.toString() : ""}</p>;
        }
        if(this.props.errors.hasCopy) {
            return this.props.errors.copy ? this.props.errors.copy : <p>""</p>;
        }
        return "";
    }

    getHeader() {
        if (this.props.errors.hasError) {
            return ERROR_HEADER;
        }
        if (this.props.errors.hasInfo) {
            return INFO_HEADER;
        }
        if (this.props.errors.hasCopy) {
            return COPY_HEADER;
        }
        return "";
    }

    getClass() {
        if (this.props.errors.hasError) {
            return "modal-header error";
        }
        if (this.props.errors.hasInfo) {
            return "modal-header info";
        }
        if (this.props.errors.hasCopy) {
            return "modal-header info";
        }
        return "";
    }

    getVisibility() {
        if (this.props.errors.hasError || this.props.errors.hasInfo || this.props.errors.hasCopy) {
            return "modal visible"
        }
        return "hidden"
    }

    handleKeyDown() {

    }

    render() {
        return (
            <div className={this.getVisibility()} onKeyDown={this.onClose.bind(this)}>
                <div className="modal-content">
                    <div className={this.getClass()}>
                        <span className="close" onClick={this.onClose.bind(this)}>×</span>
                        <div className="modal-header-name">{this.getHeader()}</div>
                    </div>
                    <div className="modal-body">
                        {this.getMessage()}
                    </div>
                </div>
            </div>
        )
    }
}