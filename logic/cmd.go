package logic

import (
	"io/fs"
	"log"
	"sync"
)

type Data struct {
	Name    string
	Path    *string
	Info    fs.FileInfo
	Matched bool
}

type SafeEntries struct {
	mu sync.Mutex
	v  []Data
}

func (c *SafeEntries) Set(newDatas []Data) {
	c.mu.Lock()
	c.v = newDatas
	c.mu.Unlock()
}

func (c *SafeEntries) Add(newData Data) {
	c.mu.Lock()
	c.v = append(c.v, newData)
	c.mu.Unlock()
}
func (c *SafeEntries) AddAll(newDatas []Data) {
	c.mu.Lock()
	c.v = append(c.v, newDatas...)
	c.mu.Unlock()
}
func (c *SafeEntries) Value() *[]Data {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &c.v
}

func (c *SafeEntries) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.v)
}

func (c *SafeEntries) Get(index int) (*Data, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.v) > index && index >= 0 {
		return &c.v[index], nil
	}
	return &Data{}, nil
}

type SafeCounter struct {
	mu sync.Mutex
	v  int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	c.v++
	c.mu.Unlock()
}
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}

func hasArgsCount[T []string](count int, commandArguments *T) bool {
	res := !bool(count < len(*commandArguments) || count > len(*commandArguments))
	if !res {
		log.Printf("%v arguments are required\r\n", count-1)
	}

	return res
}
