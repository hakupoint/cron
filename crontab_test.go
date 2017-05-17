package crontab

import (
  "testing"
  "time"
)

func Test_Cron(t *testing.T){
  sumA, sumB := 0, 0
  cron := NewCrontab()
  cron.AddTask(1, func(){sumA++})
  cron.AddTask(1, func(){sumB++})
  cron.Start()
  time.Sleep(time.Second * 10)
  if sumA != 9 || sumB != 9{
    t.Fail()
  }
}