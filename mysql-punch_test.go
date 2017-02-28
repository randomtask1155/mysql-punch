package main

import (
  "testing"
  "fmt"
  "os"
  "path"
  "time"
  "io/ioutil"
)



func TestRunTests(t *testing.T) {
   tgetNextServer(t)
   treadQueries(t)
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
  
func treadQueries(t *testing.T) {
  *query = path.Join(os.TempDir(), fmt.Sprintf("mysql-punch-readq-tempfile-%d", time.Now().Unix()) )
  data := []byte("select 1\nselect 1\n\n")
  err := ioutil.WriteFile(*query, data, 0666)
  if err != nil {
    t.Fatal(err)
  }
  readQueries()
  
  for i := range queries {
    if queries[i] == "" {
      t.Fatal(fmt.Errorf("Null found in queries"))
    }
  }
  
}