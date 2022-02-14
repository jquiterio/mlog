/*
 * @file: inmem.go
 * @author: Jorge Quitério
 * @copyright (c) 2022 Jorge Quitério
 * @license: MIT
 */

package main

import (
	"os"
	"sync"
	"time"
)

// var M *Mem
// var D *Disk

// func init() {
// 	M = NewMem(time.Duration(2) * time.Minute)
// }

type Mem struct {
	retTime time.Duration
	items   sync.Map
	close   chan struct{}
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
	d := os.Getenv("MEM_DISK")
	if d == "/tmp/mlog" {
		panic("MEM_DISK env var not set")
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
