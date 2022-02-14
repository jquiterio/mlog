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

var M *Mem

func init() {
	M = NewMem(time.Duration(24) * time.Hour)
}

type Mem struct {
	items sync.Map
	close chan struct{}
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
		close: make(chan struct{}),
	}
	d := os.Getenv("MEM_DISK")
	if d == "/tmp/mlog" {
		panic("MEM_DISK env var not set")
	}
	disk, err := NewDisk(d)
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(keeTime)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()
				m.items.Range(func(k, v interface{}) bool {
					item := v.(*item)
					if item.expires < now {
						disk.WriteLog(item.log)
						m.items.Delete(k)
					}
					return true
				})
			case <-m.close:
				return
			}
		}
	}()
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
		expires: time.Now().UnixNano() + time.Hour.Nanoseconds(),
	})
}
