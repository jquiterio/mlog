/*
 * @file: data.go
 * @author: Jorge Quitério
 * @copyright (c) 2022 Jorge Quitério
 * @license: MIT
 */

package main

import (
	"encoding/json"
	"time"

	"github.com/jquiterio/uuid"
)

type Msg map[string]interface{}

type Log struct {
	ID            string    `json:"id"`
	Time          time.Time `json:"time"`
	Src           string    `json:"source"`
	Collection    string    `json:"collection"`
	CorrelationId string    `json:"Correlation_id"` //  correlation Id
	Msg           Msg       `json:"msg"`
}

// type Collection struct {
// 	Name string `json:"name"`
// 	Logs []Log  `json:"logs"`
// }

//var Colletions = map[string]*Collection{}
//var Collections = make(map[string]*Collection)

func NewLog(col, correlId, src string, msg Msg) *Log {
	return &Log{
		ID:            uuid.NewV4().String(),
		Time:          time.Now(),
		Src:           src,
		Collection:    col,
		CorrelationId: correlId,
		Msg:           msg,
	}
}

// Get Map from Json Byte
func (l *Log) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = l.ID
	m["time"] = l.Time
	m["src"] = l.Src
	m["coll"] = l.Collection
	m["correl_id"] = l.CorrelationId
	m["msg"] = l.Msg
	return m
}

func (l *Log) ToJson() []byte {
	b, err := json.Marshal(l)
	if err != nil {
		return []byte{}
	}
	return b
}

// func NewCollection(name string) *Collection {
// 	return &Collection{
// 		Name: name,
// 	}
// }

// // Map Json Byte to map[string]interface{}
// func (c *Collection) ToMap(log Log) map[string]interface{} {
// 	var m map[string]interface{}
// 	if err := json.Unmarshal(log.Msg, &m); err != nil {
// 		return map[string]interface{}{}
// 	}
// 	return m
// }

// func (c *Collection) AddLog(log Log) {
// 	if log.Collection == c.Name {
// 		c.Logs = append(c.Logs, log)
// 	}
// }
