/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

import "./svg.css"
import React from "react";
//import detectBrowsers from "../../../../../common/detectBrowsers"

const DIRECTION = {RIGHT: 0, LEFT: 1, DOWN: 2};
const SHIFT_TEXT = 37;
const SHIFT_TEXT_CENTER = 10;
const RADIUS = 25;
const LENGTH = 300;
const MAX_STEP = 3;

const INITIAL_POSITION_X = 200;
const INITIAL_POSITION_X_SHORT = 100;

const INITIAL_POSITION_Y = 70;
const INITIAL_POSITION_LABEL_X = 10;
const INITIAL_POSITION_LABEL_Y = INITIAL_POSITION_Y;

const ADDITIONAL_X_LABEL_RIGHT = 175;
const ADDITIONAL_X_LABEL_RIGHT_SHORT = 75;

const SUMMARY_STEP = 18;
const SMALL_SHIFT_TEXT = 7;

export default class Diagram extends React.Component {

    static renderLine(option) {
        let sx = option.sx;
        let sy = option.sy;
        let ex = option.ex;
        let ey = option.ey;
        switch (option.direction) {
            case DIRECTION.DOWN:
                sy += RADIUS;
                ey -= RADIUS;
                break;
            case DIRECTION.RIGHT:
                sx += RADIUS;
                ex -= RADIUS;
                break;
            case DIRECTION.LEFT:
                sx -= RADIUS;
                ex += RADIUS;
                break;
        }
        return (<line x1={sx + "px"} y1={sy + "px"} x2={ex + "px"} y2={ey + "px"} className="line-style"/>)
    }

    static renderCircle(x, y) {
        return (<circle cx={x + "px"} cy={y + "px"} r={RADIUS + "px"} className="circle-style"/>)
    }

    static renderTextRouter(x, y, name) {
        return (<text className="text-center" x={x + "px"} y={(y - RADIUS - 2) + "px"}>{name}</text>)
    }

    static renderTextCenter(option, name, className) {
        let x = (option.ex + option.sx) / 2;
        let y = (option.ey + option.sy) / 2;
        let baseClassName = "text-center";
        if (option.direction == DIRECTION.DOWN) {
            baseClassName = "text-left";
            x += SHIFT_TEXT_CENTER;
            y += SHIFT_TEXT_CENTER * 3;
        } else {
            y -= SHIFT_TEXT_CENTER;
        }


        let finishClassName = className ? baseClassName + " " + className : baseClassName;
        return (<text className={finishClassName} x={x + "px"} y={y + "px"}>{name}</text>)
    }

    static renderText(x, y, className, text) {
        return (<text className={className} x={x + "px"} y={y + "px"}>{text}</text>)
    }

    static textRouters(options, text1, text2) {
        let res = [];
        switch (options.direction) {
            case (DIRECTION.LEFT):
                if (text2) {
                    res.push(Diagram.renderText(options.ex, options.ey - SHIFT_TEXT, "text-center text-bold", text2));
                }
                if (text1) {
                    res.push(Diagram.renderText(options.sx, options.sy - SHIFT_TEXT, "text-center text-bold", text1));
                }
                break;
            case (DIRECTION.RIGHT):
                if (text1) {
                    res.push(Diagram.renderText(options.sx, options.sy - SHIFT_TEXT, "text-center text-bold", text1));
                }
                if (text2) {
                    res.push(Diagram.renderText(options.ex, options.ey - SHIFT_TEXT, "text-center text-bold", text2));
                }
                break;
            case (DIRECTION.DOWN):
                if (text1) {
                    res.push(Diagram.renderText(options.sx, options.sy - SHIFT_TEXT, "text-center text-bold", text1));
                }
                if (text2) {
                    res.push(Diagram.renderText(options.ex, options.ey - SHIFT_TEXT, "text-center text-bold", text2));
                }
                break;
        }
        return res;
    }

    static textIps(options, text1, text2) {
        let res = [];
        switch (options.direction) {
            case (DIRECTION.LEFT):
                if (text1) {
                    res.push(Diagram.renderText(options.sx - SHIFT_TEXT, options.sy - SMALL_SHIFT_TEXT, "text-right text-decorator-green", text1));
                }
                if (text2) {
                    res.push(Diagram.renderText(options.ex + SHIFT_TEXT, options.ey - SMALL_SHIFT_TEXT, "text-left text-decorator-green", text2));
                }
                break;
            case (DIRECTION.RIGHT):
                if (text1) {
                    res.push(Diagram.renderText(options.sx + SHIFT_TEXT, options.sy - SMALL_SHIFT_TEXT, "text-left text-decorator-green", text1));
                }
                if (text2) {
                    res.push(Diagram.renderText(options.ex - SHIFT_TEXT, options.ey - SMALL_SHIFT_TEXT, "text-right text-decorator-green", text2));
                }
                break;
            case (DIRECTION.DOWN):
                if (options.nextDirection === DIRECTION.LEFT) {
                    if (text1) {
                        res.push(Diagram.renderText(options.sx + SHIFT_TEXT, options.sy - SMALL_SHIFT_TEXT, "text-left", text1));
                    }
                    if (text2) {
                        res.push(Diagram.renderText(options.ex + SHIFT_TEXT, options.ey - SMALL_SHIFT_TEXT, "text-left", text2));
                    }
                } else {
                    if (text1) {
                        res.push(Diagram.renderText(options.sx - SHIFT_TEXT, options.sy - SMALL_SHIFT_TEXT, "text-right", text1));
                    }
                    if (text2) {
                        res.push(Diagram.renderText(options.ex - SHIFT_TEXT, options.ey - SMALL_SHIFT_TEXT, "text-right", text2));
                    }
                }
                break;
        }
        return res;
    }

    static textInterfaces(options, text1, text2) {
        let res = [];
        switch (options.direction) {
            case (DIRECTION.LEFT):
                if (text1) {
                    res.push(Diagram.renderText(options.sx - SHIFT_TEXT, options.sy + SMALL_SHIFT_TEXT * 2, "text-right", text1));
                }
                if (text2) {
                    res.push(Diagram.renderText(options.ex + SHIFT_TEXT, options.ey + SMALL_SHIFT_TEXT * 2, "text-left", text2));
                }
                break;
            case (DIRECTION.RIGHT):
                if (text1) {
                    res.push(Diagram.renderText(options.sx + SHIFT_TEXT, options.sy + SMALL_SHIFT_TEXT * 2, "text-left", text1));
                }
                if (text2) {
                    res.push(Diagram.renderText(options.ex - SHIFT_TEXT, options.ey + SMALL_SHIFT_TEXT * 2, "text-right", text2));
                }
                break;
            case (DIRECTION.DOWN):
                if (options.nextDirection === DIRECTION.LEFT) {
                    if (text1) {
                        res.push(Diagram.renderText(options.sx + SHIFT_TEXT, options.sy + SMALL_SHIFT_TEXT * 2, "text-left", text1));
                    }
                    if (text2) {
                        res.push(Diagram.renderText(options.ex + SHIFT_TEXT, options.ey + SMALL_SHIFT_TEXT * 2, "text-left", text2));
                    }
                } else {
                    if (text1) {
                        res.push(Diagram.renderText(options.sx - SHIFT_TEXT, options.sy + SMALL_SHIFT_TEXT * 2, "text-right", text1));
                    }
                    if (text2) {
                        res.push(Diagram.renderText(options.ex - SHIFT_TEXT, options.ey + SMALL_SHIFT_TEXT * 2, "text-right", text2));
                    }
                }
                break;
        }
        return res;
    }

    static checkCrds(x, y) {
        if (x < 0) {
            Diagram.minx = x;
        } else {
            Diagram.maxx = Math.max(Diagram.maxx, x);
        }
        if (y < 0) {
            Diagram.miny = y;
        } else {
            Diagram.maxy = Math.max(Diagram.maxy, y);
        }

    }

    static numberToLocal(stringValue) {
        let number = parseFloat(stringValue);
        if (number) {
            return number.toLocaleString("de-DE", {maximumFractionDigits: 0}); //de-DE with point delimiter
        }
        return stringValue;
    }

    static textSummaryRouter(x, y, router, className, drawUp) {
        let res = [];

        if (drawUp === true) {
            res.push(Diagram.renderText(x, y, className, Diagram.numberToLocal(router.outputp)));

            y -= SUMMARY_STEP;
            res.push(Diagram.renderText(x, y, className, Diagram.numberToLocal(router.inputp)));

            y -= SUMMARY_STEP;
            res.push(Diagram.renderText(x, y, className, Diagram.numberToLocal(router.outputb)));

            y -= SUMMARY_STEP;
            res.push(Diagram.renderText(x, y, className, Diagram.numberToLocal(router.inputb)));

        } else {
            res.push(Diagram.renderText(x, y, className, Diagram.numberToLocal(router.inputb)));

            y += SUMMARY_STEP;
            res.push(Diagram.renderText(x, y, className, Diagram.numberToLocal(router.outputb)));

            y += SUMMARY_STEP;
            res.push(Diagram.renderText(x, y, className, Diagram.numberToLocal(router.inputp)));

            y += SUMMARY_STEP;
            res.push(Diagram.renderText(x, y, className, Diagram.numberToLocal(router.outputp)));
        }
        Diagram.checkCrds(x, y);
        return res;
    }

    static textSummary(options, router1, router2) {
        let res = [];
        switch (options.direction) {
            case (DIRECTION.LEFT):
                if (router1) {
                    res.push(Diagram.textSummaryRouter(options.sx - SHIFT_TEXT, options.sy + SHIFT_TEXT + SMALL_SHIFT_TEXT, router1, "text-right", false));
                }
                if (router2) {
                    res.push(Diagram.textSummaryRouter(options.ex + SHIFT_TEXT, options.ey + SHIFT_TEXT + SMALL_SHIFT_TEXT, router2, "text-left", false));
                }
                break;
            case (DIRECTION.RIGHT):
                if (router1) {
                    res.push(Diagram.textSummaryRouter(options.sx + SHIFT_TEXT, options.sy + SHIFT_TEXT + SMALL_SHIFT_TEXT, router1, "text-left", false));
                }
                if (router2) {
                    res.push(Diagram.textSummaryRouter(options.ex - SHIFT_TEXT, options.ey + SHIFT_TEXT + SMALL_SHIFT_TEXT, router2, "text-right", false));
                }
                break;
            case (DIRECTION.DOWN):
                if (options.nextDirection === DIRECTION.LEFT) {
                    if (router1) {
                        res.push(Diagram.textSummaryRouter(options.sx + SHIFT_TEXT, options.sy + SHIFT_TEXT + SMALL_SHIFT_TEXT, router1, "text-left", false));
                    }
                    if (router2) {
                        res.push(Diagram.textSummaryRouter(options.ex + SHIFT_TEXT, options.ey + SHIFT_TEXT + SMALL_SHIFT_TEXT, router2, "text-left", false));
                    }
                } else {
                    if (router1) {
                        res.push(Diagram.textSummaryRouter(options.sx - SHIFT_TEXT, options.sy + SHIFT_TEXT + SMALL_SHIFT_TEXT, router1, "text-right", false));
                    }
                    if (router2) {
                        res.push(Diagram.textSummaryRouter(options.ex - SHIFT_TEXT, options.ey + SHIFT_TEXT + SMALL_SHIFT_TEXT, router2, "text-right", false));
                    }
                }
                break;
        }
        return res;
    }

    static textSummaryLabel(options) {

        if (options.direction !== DIRECTION.DOWN) {
            return [];
        }

        options.labelYs.push(options.ey);
        return Diagram.textSummaryLabelRender(INITIAL_POSITION_LABEL_X, options.ey, false);
    }

    static textSummaryLabelRender(x, y, isRight = false, isShort = false) {
        y += SHIFT_TEXT + SMALL_SHIFT_TEXT;
        let res = [];

        let txtClass = "text-left text-label";
        if (isRight === true) {
            txtClass = "text-right text-label";
            x += isShort === true ? ADDITIONAL_X_LABEL_RIGHT_SHORT : ADDITIONAL_X_LABEL_RIGHT;
        }
        res.push(Diagram.renderText(x, y, txtClass, "Byte In"));

        y += SUMMARY_STEP;
        res.push(Diagram.renderText(x, y, txtClass, "Byte Out"));

        y += SUMMARY_STEP;
        res.push(Diagram.renderText(x, y, txtClass, "Packet In"));

        y += SUMMARY_STEP;
        res.push(Diagram.renderText(x, y, txtClass, "Packet Out"));

        Diagram.checkCrds(x, y);

        return res
    }

    static caclulateOptions(options, lastId, nextId) {

        if (lastId === nextId) {
            options.render = false;
            return options;
        }
        options.render = true;

        if (options.step >= MAX_STEP) {
            let direction = null;
            let nextDirection = options.nextDirection;
            options.step = 1; //1 элемент уже нарисовали

            switch (options.direction) {
                case DIRECTION.RIGHT:
                    direction = DIRECTION.DOWN;
                    nextDirection = DIRECTION.LEFT;
                    options.step++; //чтобы пройти по условию MAX_STEP - вниз рисуем 1 раз
                    break;
                case DIRECTION.LEFT:
                    direction = DIRECTION.DOWN;
                    nextDirection = DIRECTION.RIGHT;
                    options.step++; //чтобы пройти по условию MAX_STEP - вниз рисуем 1 раз
                    break;
                case DIRECTION.DOWN:
                    direction = options.nextDirection;
                    nextDirection = options.nextDirection;
                    break;
            }

            options.direction = direction;
            options.nextDirection = nextDirection;
        }
        options.sx = options.ex;
        options.sy = options.ey;
        if (options.direction === DIRECTION.RIGHT) {
            options.ex += LENGTH;
        } else if (options.direction === DIRECTION.LEFT) {
            options.ex -= LENGTH;
        } else if (options.direction === DIRECTION.DOWN) {
            options.ey += LENGTH * 3 / 4;
        }
        options.step++;
        return options;
    }

    static getMax(max, test1, test2) {
        return Math.max(Math.max(max, test1), test2);
    }

    static maxx = 0;
    static maxy = 0;
    static minx = 0;
    static miny = 0;

    static generate(routes) {
        let initialX = routes.length < 6 ? INITIAL_POSITION_X_SHORT : INITIAL_POSITION_X;

        let options = {
            sx: initialX,
            ex: initialX,
            sy: INITIAL_POSITION_Y,
            ey: INITIAL_POSITION_Y,
            direction: DIRECTION.RIGHT,
            nextDirection: DIRECTION.RIGHT,
            step: 0,
            render: false,
            labelYs: []
        };
        Diagram.maxx = options.sx;
        Diagram.maxy = options.sy;

        let lastRouterId = undefined;
        let lines = [];
        let res = [];

        res.push(Diagram.textSummaryLabelRender(INITIAL_POSITION_LABEL_X, INITIAL_POSITION_LABEL_Y, false));
        options.labelYs.push(INITIAL_POSITION_LABEL_Y);

        for (let i = 0; i < routes.length; i++) {
            let route = routes[i];

            options = Diagram.caclulateOptions(options, lastRouterId, route.router1.id);
            Diagram.maxx = Diagram.getMax(Diagram.maxx, options.ex, options.ex);
            Diagram.maxy = Diagram.getMax(Diagram.maxy, options.ey, options.ey);

            if (options.render === true) {
                lines.push(Diagram.renderLine(options));
                res.push(Diagram.renderTextCenter(options, route.baseIp, "text-decorator-green"));

                res.push(Diagram.renderCircle(options.sx, options.sy));
                res.push(Diagram.renderCircle(options.ex, options.ey));

                res.push(Diagram.textRouters(options, route.router1.name, route.router2.name));
                res.push(Diagram.textIps(options, route.router1.ip, route.router2.ip));
                res.push(Diagram.textInterfaces(options, route.router1.interface, route.router2.interface));
                res.push(Diagram.textSummary(options, route.router1, route.router2));

                res.push(Diagram.textSummaryLabel(options));

                options.step++;
            } else {
                options = Diagram.caclulateOptions(options, route.router1.id, route.router2.id);
                Diagram.maxx = Diagram.getMax(Diagram.maxx, options.ex, options.ex);
                Diagram.maxy = Diagram.getMax(Diagram.maxy, options.ey, options.ey);

                if (options.render === true) {
                    lines.push(Diagram.renderLine(options));
                    res.push(Diagram.renderTextCenter(options, route.baseIp, "text-decorator-green"));

                    res.push(Diagram.renderCircle(options.ex, options.ey));

                    res.push(Diagram.textRouters(options, undefined, route.router2.name));
                    res.push(Diagram.textIps(options, route.router1.ip, route.router2.ip));
                    res.push(Diagram.textInterfaces(options, route.router1.interface, route.router2.interface));
                    res.push(Diagram.textSummary(options, route.router1, route.router2));

                    res.push(Diagram.textSummaryLabel(options));
                }
            }

            lastRouterId = route.router2.id;
        }

        let maxx = Diagram.maxx;
        for (let i = 0; i < options.labelYs.length; i++) {
            res.push(Diagram.textSummaryLabelRender(maxx, options.labelYs[i], true, routes.length < 3));
        }

        let sx = Diagram.minx;
        let sy = Diagram.miny;
        let ex = Diagram.maxx + INITIAL_POSITION_LABEL_X;
        let ey = Diagram.maxy + 30;

        let svgStyle = {
            minWidth: (ex / 2) + 'px',
            maxWidth: ex + 'px',
            minHeight: (ey / 2) + 'px',
            maxHeight: ey + 'px',
            height: ey + 'px',
            width: ex + 'px',
        };

        return (<svg className="dia" width={ex + "px"} height={ey + "px"} x="0px" y="0px"
                     style={svgStyle}>{lines}{res}</svg>)
    }

    render() {
        return Diagram.generate(this.props.data.paths);
    }
}
