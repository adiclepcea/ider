package main

import (
	"fmt"
	"time"
  "sync"
  "os"
  "log"
)

var chanIndex chan uint
var wg sync.WaitGroup
const MAX = 100000
const machineIndex = 1

func init(){
  chanIndex = make(chan uint,1)
}

func main() {
  f, err:=os.Create("out.txt")
  if err!=nil{
    log.Println(err.Error)
    return
  }
  defer f.Close()

  go generateIndex()

  values:=make([]int64,MAX)

  for i:=0; i<MAX; i+=1000 {
    wg.Add(1)
    go func(from int){
			defer wg.Done()
      for j:=from; j<from+1000; j++{
        values[j] = GenerateId()
      }
    }(i)
  }
  wg.Wait()
  for i:=0; i<MAX; i++ {
    f.Write([]byte(fmt.Sprintf("%d\n",values[i])))
  }
  f.Sync();
}

func GenerateId() int64 {
	c:=<-chanIndex
	return time.Now().Unix()*int64(100000000)+int64(c)*10000+int64(machineIndex)
}

func generateIndex() uint {
  var curVal uint = 0
  var when int64 = 0
  for {
    now:=time.Now().Unix()
		if now==when{
			if curVal>9998 {
				time.Sleep(1*time.Millisecond)
			}
		} else {
      when=now
      curVal=0
    }
    curVal++
    chanIndex<-curVal
  }
}
