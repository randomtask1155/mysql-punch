package main

import (
  "testing"
  "fmt"
)



func TestRunTests(t *testing.T) {
   tgetNextServer(t)
}

func tgetNextServer(t *testing.T) {
  poolIndex = 0
  counter := make(map[int]int)
  serverPool = make([]string,0)
  serverPool = append(serverPool, "server1")
  serverPool = append(serverPool, "server2")
  serverPool = append(serverPool, "server3")
  serverPool = append(serverPool, "server4")
  for i := 0; i < (len(serverPool))*100; i++ {
    n := getNextServer()
    _, ok := counter[n]
      if ok {
        counter[n]++
      } else {
        counter[n] = 1
      }
    }
    for _,v := range counter {
      if v != ((len(serverPool))*100)/4 {
        t.Fatal(fmt.Errorf("%d is not equal to %d", v, ((len(serverPool))*100)/4))
      }
    }
  }