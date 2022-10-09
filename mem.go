/*
 * @file: inmem.go
 * @author: Jorge Quitério
 * @copyright (c) 2022 Jorge Quitério
 * @license: MIT
 */

package main

import (
	"sync"
	"time"
)

// var M *Mem
// var D *Disk

// func init() {
// 	M = NewMem(time.Duration(2) * time.Minute)
// }

type Mem struct {
	retTime    time.Duration
	items      sync.Map
	colections []string
	close      chan struct{}
}

func (m *Mem) Close() {
	close(m.close)
}

type item struct {
	log     Log
	expires int64
}

func NewMem(keeTime time.Duration) *Mem {
	m := &Mem{
		close:   make(chan struct{}),
		retTime: keeTime,
	}
	return m
}

func (m *Mem) Get(key string) *Log {
	i, ok := m.items.Load(key)
	if !ok {
		return nil
	}
	item := i.(*item)
	return &item.log
}

func (m *Mem) Set(log Log) {
	m.items.Store(log.ID, &item{
		log:     log,
		expires: time.Now().Add(1 * time.Minute).UnixNano(),
	})
}

func (m *Mem) GetCollections() []string {
	return m.colections
}

func (m *Mem) CollectionExist(col string) bool {
	return contains(m.colections, col)
}

func contains(l []string, s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}
