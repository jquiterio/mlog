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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, getAll())
	})

	e.GET("/:col", func(c echo.Context) error {
		col := c.Param("col")
		return c.JSON(200, getLogsByCollection(col))
	})

	e.POST("/", func(c echo.Context) error {
		// {"collection": "test","source": "test", "msg": {"something": "data", "otherthing": "data"}}
		var rlog Log
		if err := c.Bind(&rlog); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if rlog.Src == "" {
			return c.JSON(400, map[string]string{"error": "source is required"})
		}
		if rlog.Collection == "" {
			rlog.Collection = "no_collection"
		}
		if rlog.Msg == nil {
			return c.JSON(400, map[string]string{"error": "msg is required"})
		}
		log := NewLog(rlog.Collection, rlog.CorrelationId, rlog.Src, rlog.Msg)
		M.Set(*log)
		return c.JSON(201, log)
	})

	e.POST("/q/:col", func(c echo.Context) error {
		col := c.Param("col")
		var q Query
		if err := c.Bind(&q); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fmt.Println("Query: ", q.Q)
		res := NewQuery(q.Q, col)
		return c.JSON(200, res)

	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	s := &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}
	fmt.Println("Listening on port " + port)
	s.ListenAndServe()
}

func getAll() (logs []Log) {
	M.items.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %v, value: %v\n", key, value)
		logs = append(logs, value.(*item).log)
		return true
	})
	return logs
}

func getLogsByCollection(col string) []Log {
	var logs []Log
	M.items.Range(func(key, value interface{}) bool {
		item := value.(*item)
		if item.log.Collection == col {
			logs = append(logs, item.log)
		}
		return true
	})
	return logs
}
