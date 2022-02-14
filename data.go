/*
 * @file: data.go
 * @author: Jorge Quitério
 * @copyright (c) 2022 Jorge Quitério
 * @license: MIT
 */

package main

import (
	"time"

	"github.com/jquiterio/uuid"
)

type Log struct {
	ID            string                 `json:"id"`
	Time          string                 `json:"time"`
	Src           string                 `json:"source"`
	Collection    string                 `json:"collection"`
	CorrelationId string                 `json:"correlation_id"` //  correlation Id
	Msg           map[string]interface{} `json:"msg"`
}

func NewLog(col, correlId, src string, msg map[string]interface{}) *Log {
	return &Log{
		ID:            uuid.NewV4().String(),
		Time:          time.Now().Format(time.RFC3339),
		Src:           src,
		Collection:    col,
		CorrelationId: correlId,
		Msg:           msg,
	}
}
