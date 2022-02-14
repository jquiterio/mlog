/*
 * @file: mlog.go
 * @author: Jorge Quitério
 * @copyright (c) 2022 Jorge Quitério
 * @license: MIT
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	mem_retTime := os.Getenv("LOG_RETENTION_PERIOD")
	if mem_retTime == "" {
		mem_retTime = "960h"
	}
	dur, err := time.ParseDuration(mem_retTime)
	if err != nil {
		panic(err)
	}
	mem := NewMem(dur)

	go func() {
		ticker := time.NewTicker(mem.retTime)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()
				mem.items.Range(func(k, v interface{}) bool {
					item := v.(*item)
					if item.expires > 0 && now > item.expires {
						mem.items.Delete(k)
					}
					return true
				})
			case <-mem.close:
				return
			}
		}
	}()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, mem.getAll())
	})

	e.GET("/:col", func(c echo.Context) error {
		col := c.Param("col")
		return c.JSON(200, mem.getLogsByCollection(col))
	})

	e.POST("/", func(c echo.Context) error {
		// {"collection": "test","source": "test", "msg": {"something": "data", "otherthing": "data"}}
		var rlog Log
		if err := c.Bind(&rlog); err != nil {
			fmt.Println("bind error: ", err.Error())
			return c.JSON(400, echo.Map{
				"error": err.Error(),
			})
		}
		if rlog.Collection == "" {
			rlog.Collection = "no_collection"
		}
		if rlog.Msg == nil {
			return c.JSON(400, map[string]string{"error": "msg is required"})
		}
		log := NewLog(rlog.Collection, rlog.CorrelationId, rlog.Src, rlog.Msg)
		mem.Set(*log)
		return c.JSON(201, log)
	})

	e.POST("/q/:col", func(c echo.Context) error {
		col := c.Param("col")
		var q Query
		if err := c.Bind(&q); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		res := mem.doQuery(q.Q, col)
		return c.JSON(200, res)
	})

	// midleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}
	s := &http.Server{
		Addr:    "127.0.0.1:" + port,
		Handler: e,
	}
	fmt.Println("listening on port: ", s.Addr)
	s.ListenAndServe()
}

func (m *Mem) getAll() (logs []Log) {
	m.items.Range(func(key, value interface{}) bool {
		logs = append(logs, value.(*item).log)
		return true
	})
	return logs
}

func (mem *Mem) getLogsByCollection(col string) []Log {
	var logs []Log
	mem.items.Range(func(key, value interface{}) bool {
		item := value.(*item)
		if item.log.Collection == col {
			logs = append(logs, item.log)
		}
		return true
	})
	return logs
}
