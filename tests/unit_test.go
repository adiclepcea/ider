package ider_test

import (
	"ider"
	"testing"
	"time"
)

func TestIndexCreationShouldFail(t *testing.T) {
	_, err := ider.NewIder(10000)
	if err == nil {
		t.Fatal("An error should be returned for the creation of an ider with machineIndex above 9999, got nil")
	}
}

func TestIndexCreationOneIder(t *testing.T) {
	var index1 int64
	var index2 int64

	idProvider, err := ider.NewIder(1)
	if err != nil {
		t.Errorf("No error expected while creating an ider, got %s", err.Error())
		t.FailNow()
	}

	before := time.Now().Unix()

	index1 = idProvider.GenerateID()
	time.Sleep(1 * time.Millisecond)
	index2 = idProvider.GenerateID()

	after := time.Now().Unix()

	//if we were at the second border, we requery
	if before != after {
		index1 = idProvider.GenerateID()
		time.Sleep(1 * time.Millisecond)
		index2 = idProvider.GenerateID()
	}

	if index1/100000000 != index2/100000000 {
		t.Errorf("The time shoule be the same for two consecutive ids generated in the same second. Got %d and %d\n", index1/100000000, index2/100000000)
		t.FailNow()
	}

	if index2-index1 != 10000 {
		t.Errorf("Two consecutive ids generated in the same second should be 10000 unit appart. Got %d and %d\n", index1, index2)
		t.FailNow()
	}

}

func TestIndexCreationTwoIders(t *testing.T) {
	var index1 int64
	var index2 int64

	idProvider1, err := ider.NewIder(1)
	if err != nil {
		t.Errorf("No error expected while creating an ider, got %s", err.Error())
		t.FailNow()
	}

	idProvider2, err := ider.NewIder(2)
	if err != nil {
		t.Errorf("No error expected while creating an ider, got %s", err.Error())
		t.FailNow()
	}

	before := time.Now().Unix()

	index1 = idProvider1.GenerateID()
	index2 = idProvider2.GenerateID()

	after := time.Now().Unix()

	//if we were at the second border, we requery
	if before != after {
		index1 = idProvider1.GenerateID()
		index2 = idProvider2.GenerateID()
	}

	if index1/100000000 != index2/100000000 {
		t.Errorf("The time shoule be the same for two consecutive ids generated in the same second. Got %d and %d\n", index1/100000000, index2/100000000)
		t.FailNow()
	}

	if index2-index1 != 1 {
		t.Errorf("Two consecutive ids generated in the same second on two iders with consecutive machineIndexes should be 1 unit appart. Got %d and %d\n", index1, index2)
		t.FailNow()
	}
}

/*
func main() {
	f, err := os.Create("out.txt")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer f.Close()

	go generateIndex()

	values := make([]int64, MAX)

	for i := 0; i < MAX; i += 1000 {
		wg.Add(1)
		go func(from int) {
			defer wg.Done()
			for j := from; j < from+1000; j++ {
				values[j] = GenerateId()
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < MAX; i++ {
		f.Write([]byte(fmt.Sprintf("%d\n", values[i])))
	}
	f.Sync()
}
*/
