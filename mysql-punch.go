package main

import (
  "fmt"
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "flag"
  "sync"
  "time"
  "os"
  "io/ioutil"
  "os/signal"
  "syscall"
  "strings"
)
var (
  user = flag.String("u", "root", "user default to root")
  passwd = flag.String("p", "", "password defaults to blank")
  database = flag.String("d", "test", "database to use for testing")
  query = flag.String("q", "", "file containing list of queries to run per connection")
  servers = flag.String("s", "localhost", "comma delimited list of servers to use") 
  port = flag.String("port", "3306", "mysql server port")
  maxConn = flag.Int("c", 2, "total number of connections to create")
  interval = flag.Int("i", 200, "interval in miliseconds to sleep between sql statements")
  
  
  serverPool []string
  queries []string
  poolIndex = 0
  openConn = 0 // current number of open connections
  indexLock = sync.Mutex{}
  connLock = sync.Mutex{}
)

func connectSQL() (*sql.DB, error) {
  server := serverPool[getNextServer()]
  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", *user, *passwd, server, *port, *database))
  if err != nil {
    return db, err 
  }
  db.SetMaxOpenConns(0)
  return db,err
}


// round robin the servers
func getNextServer() int {
  indexLock.Lock() 
  if poolIndex > len(serverPool)-1 {
    poolIndex = 0
  }
  newID := poolIndex // copy to new var in case index changes before we return it
  poolIndex++
  indexLock.Unlock()
  return newID
}

func readQueries() {
  b, err := ioutil.ReadFile(*query)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  s := fmt.Sprintf("%s", b)
  queries = strings.Split(s, "\n")
}

func spingSQL(i int, errMsg chan error, kill chan int) {
  /*db, err := connectSQL()
  defer db.Close()
  if err != nil {
    errMsg <- err 
    return
  }*/
  for {
    select {
    case _ =  <- kill:
      return
    default:
      for q := range queries {
        
        db, err := connectSQL()
        if err != nil {
          errMsg <- err 
          return
        }
        
        tx, err := db.Begin()
        if err != nil {
          db.Close()
          errMsg <- err
          return
        }
        
        stmt, err := tx.Prepare(queries[q])
        if err != nil {
          db.Close()
          errMsg <- err 
          return
        }
        
        _, err = stmt.Exec()
        if err != nil {
          tx.Rollback()
          db.Close()
          errMsg <- err 
          return
        }
        stmt.Close()
        tx.Commit()
        db.Close()
        
      }
      time.Sleep(time.Millisecond * time.Duration(int64(i)))
    }
  }
}

func addOpenConn() {
  connLock.Lock()
    openConn++
  connLock.Unlock()
}

func removeOpenConn() {
  connLock.Lock()
    openConn--
  connLock.Unlock()
}

func getOpenConn() int {
  connLock.Lock()
    newConn := openConn
  connLock.Unlock()
  return newConn
}

// send kill signal to all mysql routines so connections get closed
func closeSQL(c chan os.Signal, k chan int) {
  <-c 
  fmt.Println("Cleaning up active go routines before exit!")
  for i := 0; i < getOpenConn(); i++ {
    if i >= getOpenConn() {
      return
    }
    k <- 1
  }
  time.Sleep(time.Millisecond * 1000)
  os.Exit(100)
}

func main(){
  flag.Parse()
  readQueries()
  errChan := make(chan error,0)
  killChan := make(chan int,0)
  serverPool = strings.Split(*servers, ",")
  
  fmt.Println("Starting connections...")
  for i := 0; i < *maxConn; i++ {
    go spingSQL(*interval, errChan, killChan)
    addOpenConn()
    time.Sleep(time.Millisecond * 500) // slowly start up connections
  }
  
  fmt.Println("Finished starting connections...")
  fmt.Println("Listening for SEGTERM events...")
  // capture ctrl+c
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go closeSQL(c, killChan)
  
  fmt.Println("Running...")
  for i := 0; i <*maxConn; i++ {
    select {
    case e := <-errChan:
      fmt.Println(e)
      removeOpenConn()
    }
  }
}
