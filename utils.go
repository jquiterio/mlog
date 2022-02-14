/*
 * @file: utils.go
 * @author: Jorge Quitério
 * @copyright (c) 2022 Jorge Quitério
 * @license: MIT
 */

package main

import (
	"fmt"
	"os"
	"sync"
)

func (d *Disk) getMuttex(col string) (mut *sync.Mutex) {
	d.mut.Lock()
	defer d.mut.Unlock()
	var ok bool
	if mut, ok = d.mutts[col]; !ok {
		mut = &sync.Mutex{}
		d.mutts[col] = mut
	}
	return mut
}

func writeFile(dir string, j []byte) error {
	f, err := os.OpenFile(dir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", dir, err)
	}
	defer f.Close()

	if _, err := f.Write(j); err != nil {
		return fmt.Errorf("error writing file %s: %w", dir, err)
	}
	return nil
}
