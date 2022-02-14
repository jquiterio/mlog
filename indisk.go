/*
 * @file: indisk.go
 * @author: Jorge Quitério
 * @copyright (c) 2022 Jorge Quitério
 * @license: MIT
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Disk struct {
	mut   sync.Mutex
	mutts map[string]*sync.Mutex
	dir   string
}

func (d *Disk) CollectionOK(col string) bool {
	if f, err := os.Stat(filepath.Join(d.dir, col)); err != nil {
		if f.IsDir() {
			return false
		}
		return false
	}
	return true
}

func NewDisk(dir string) (*Disk, error) {
	dir = filepath.Clean(dir)

	disk := &Disk{
		mutts: make(map[string]*sync.Mutex),
		dir:   dir,
	}
	return disk, useNewOrCreate(dir)
}

func GetDisk(dir string) (*Disk, error) {
	dir = filepath.Clean(dir)
	disk := &Disk{
		mutts: make(map[string]*sync.Mutex),
		dir:   dir,
	}
	return disk, useNewOrCreate(dir)
}

func useNewOrCreate(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		if !f.IsDir() {
			err = os.Remove(dir)
			if err != nil {
				return err
			}
			err = os.Mkdir(dir, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Disk) OK() bool {
	if f, err := os.Stat(d.dir); err != nil {
		if !f.IsDir() {
			return false
		}
		return false
	}
	return true
}

func (d *Disk) WriteLog(log Log) error {
	col := log.Collection
	col = strings.TrimSpace(col)
	if col == "" {
		col = "no_collection"
	}
	if !d.OK() {
		return fmt.Errorf("disk not ok")
	}

	mut := d.getMuttex(col)
	mut.Lock()
	defer mut.Unlock()

	dir := filepath.Join(d.dir, col+".json")
	if err := useNewOrCreate(dir); err != nil {
		return fmt.Errorf("error creating dir %s: %w", dir, err)
	}

	if j, err := json.MarshalIndent(log, "", " "); err != nil {
		return fmt.Errorf("error marshalling log: %w", err)
	} else {
		if err := writeFile(dir, j); err != nil {
			return fmt.Errorf("error writing file %s: %w", dir, err)
		}
	}
	return nil
}

func (d *Disk) GetLog(log Log) error {
	col := log.Collection
	col = strings.TrimSpace(col)
	if col == "" {
		col = "no_collection"
	}
	if !d.CollectionOK(log.Collection) {
		return fmt.Errorf("collection %s not found", col)
	}

	fcol := filepath.Join(d.dir, col+".json")

	//  read file
	if b, err := ioutil.ReadFile(fcol); err != nil {
		return fmt.Errorf("error reading file %s: %w", fcol, err)
	} else {
		if err := json.Unmarshal(b, &log); err != nil {
			return fmt.Errorf("error unmarshalling log: %w", err)
		}
	}
	return nil
}

func (d *Disk) GetAllLogs() []Log {
	var logs []Log
	if !d.OK() {
		return logs
	}
	files, err := ioutil.ReadDir(d.dir)
	if err != nil {
		return logs
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".json") {
			log := Log{}
			if err := d.GetLog(log); err != nil {
				continue
			}
			logs = append(logs, log)
		}
	}
	return logs
}
