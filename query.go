/*
 * @file: query.go
 * @author: Jorge Quitério
 * @copyright (c) 2022 Jorge Quitério
 * @license: MIT
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/elgs/jsonql"
)

type Query struct {
	Q      string `json:"q"`
	result []Log
}

func (q *Query) Result() []Log {
	return q.result
}

func NewQuery(q string, col string) interface{} {

	var logs []Log
	logs = append(logs, query(col)...)

	b, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return []Log{}
	}
	fmt.Println("QUERY JSON: ", string(b))
	parser, err := jsonql.NewStringQuery(string(b))
	if err != nil {
		return []Log{}
	}
	v, e := parser.Query(q)
	if e != nil {
		return nil
	}
	return v
}

func query(col string) []Log {
	l := []Log{}
	mem := M
	mem.items.Range(func(key, value interface{}) bool {
		item := value.(*item)
		if item.log.Collection == col {
			l = append(l, item.log)
		}
		return true
	})
	return l
}
