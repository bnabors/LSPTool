/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

export default function detectBrowsers() {
    let result = {
        isOpera:false,
        isFirefox:false,
        isSafari:false,
        isIE:false,
        isChrome:false,
        isBlink:false,
    };
    // Opera 8.0+
    result.isOpera = (!!window.opr && !!opr.addons) || !!window.opera || navigator.userAgent.indexOf(' OPR/') >= 0;
    // Firefox 1.0+
    result.isFirefox = typeof InstallTrigger !== 'undefined';
    // At least Safari 3+: "[object HTMLElementConstructor]"
    result.isSafari = Object.prototype.toString.call(window.HTMLElement).indexOf('Constructor') > 0;
    // Internet Explorer 6-11
    result.isIE = /*@cc_on!@*/false || !!document.documentMode;
    // Edge 20+
    result.isEdge = !result.isIE && !!window.StyleMedia;
    // Chrome 1+
    result.isChrome = !!window.chrome && !!window.chrome.webstore;
    // Blink engine detection
    result.isBlink = (result.isChrome || result.isOpera) && !!window.CSS;

    return result;
}
