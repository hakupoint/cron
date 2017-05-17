package crontab

import (
  "sort"
  "time"
)
type FuncJob interface {
  update()
}
type Job func()
func (j Job) update() {
  j()
}
type crontab struct {
  runing  bool
  entries []*Entry
  add     chan *Entry
  stop    chan struct{}
}
type Entry struct {
  interval time.Duration
  Next     time.Time
  Job      FuncJob
}
type EntryArray []*Entry
func (e EntryArray) Len() int {
  return len(e)
}
func (e EntryArray) Less(i, j int) bool {
  return e[i].Next.Before(e[j].Next)
}
func (e EntryArray) Swap(i, j int) {
  e[i], e[j] = e[j], e[i]
}
func NewCrontab() *crontab {
  return &crontab{
    runing:  false,
    entries: nil,
    add:     make(chan *Entry),
    stop:    make(chan struct{}),
  }
}
func (c *crontab) AddTask(t uint, cmd func()) {
  entry := &Entry{
    interval: time.Duration(t),
    Job:      Job(cmd),
  }
  if !c.runing {
    c.entries = append(c.entries, entry)
    return
  }
  c.add <- entry
}
func (c *crontab) Start() {
  c.runing = true
  go c.run()
}
func (c *crontab) run() {
  for _, e := range c.entries {
    e.Next = time.Now().Add(time.Second * time.Duration(e.interval))
  }
  for {
    sort.Sort(EntryArray(c.entries))
    next := c.entries[0].Next
    select {
    case <-time.After(next.Sub(time.Now())):
      for _, e := range c.entries {
        if e.Next == next {
          go e.Job.update()
          e.Next = time.Now().Add(time.Second * time.Duration(e.interval))
        }
      }
    case newentry := <-c.add:
      newentry.Next = time.Now().Add(time.Second * time.Duration(newentry.interval))
      c.entries = append(c.entries, newentry)
    case <-c.stop:
      return
    }
  }
}
func (c *crontab) Stop() {
  c.stop <- struct{}{}
}
