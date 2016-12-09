/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

package lspLogger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
)

var logLevels = map[string]int{
	"none":  0,
	"error": 1,
	"info":  2,
	"debug": 3,
}

var (
	logLevel int
	logger   *log.Logger
	outIo    *bufio.Writer
	objSync  sync.Mutex
)

// Initialize инициализирует lspLogger
func Initialize(fileName string, level string) {
	var ok bool
	logLevel, ok = logLevels[level]
	if !ok {
		logLevel = 0
	}
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logLevel = 0
		fmt.Printf("Error opening file: %v", err)
		return
	}

	log.SetOutput(f)
	log.SetPrefix("Critical Error ")
	log.SetFlags(log.LstdFlags | log.Llongfile)
	outIo = bufio.NewWriter(f)
	logger = log.New(outIo, "Log ", log.LstdFlags)
}

// Errorln записывает ошибку
func Errorln(message string) {
	if logLevel < logLevels["error"] {
		return
	}
	stack := string(debug.Stack())
	objSync.Lock()
	logger.Printf("ERROR: %q \r\n StackTrace: %v \r\n", message, stack)
	outIo.Flush()
	objSync.Unlock()
}

// Error записывает ошибку
func Error(v ...interface{}) {
	Errorln(fmt.Sprint(v))
}

// Errorf записывает ошибку c форматированием
func Errorf(format string, v ...interface{}) {
	Errorln(fmt.Sprintf(format, v))
}

// Infoln записывает входящие запросы и исходящие ответы
func Infoln(message string) {
	if logLevel < logLevels["info"] {
		return
	}
	objSync.Lock()
	logger.Printf("INFO: %q \r\n", message)
	outIo.Flush()
	objSync.Unlock()
}

// Info записывает ошибку
func Info(v ...interface{}) {
	Infoln(fmt.Sprint(v))
}

// Infof записывает входящие запросы и исходящие ответы c форматированием
func Infof(format string, v ...interface{}) {
	Infoln(fmt.Sprintf(format, v))
}

// Debug записывает ошибку
func Debug(v ...interface{}) {
	Debugln(fmt.Sprint(v))
}

// Debugln записывает отладочное сообщение
func Debugln(message string) {
	if logLevel < logLevels["debug"] {
		return
	}
	objSync.Lock()
	logger.Printf("DEBUG: %q \r\n", message)
	outIo.Flush()
	objSync.Unlock()
}

// Debugf записывает отладочное сообщение c форматированием
func Debugf(format string, v ...interface{}) {
	Debugln(fmt.Sprintf(format, v))
}

// CheckError проверить на ошибку и записать
func CheckError(e error) {
	if e != nil {
		Errorln(e.Error())
	}
}
