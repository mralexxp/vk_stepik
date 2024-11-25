package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
	"sync"
	"time"
)

const (
	checkTick = 10 * time.Millisecond
)

func main() {
	ExecutePipeline()
}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}

	// Поднимаем все JOB в отдельных горутинах с надстройкой:
	for i, work := range jobs {
		if i == 0 {
			wg.Add(1)
			go generator(work, in, out, doneCh, wg)
			continue
		}
		wg.Add(1)
		if i%2 == 0 {
			go worker(work, in, out, wg)
		} else {
			go worker(work, out, in, wg)
		}
	}
LOOP:
	for {
		select {
		case <-doneCh:
			close(in)
			close(out)
			break LOOP
		default:
			time.Sleep(checkTick)
			continue
		}
	}

	wg.Wait()
}

func worker(work job,
	in chan interface{},
	out chan interface{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	// Основная работа
	work(in, out)

}

// Функция первой работы с закрытием канала
func generator(work job, in chan interface{}, out chan interface{}, doneCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	work(in, out)

	close(doneCh)
}

// SingleHash TODO: Отдельная функция на дженериках?
func SingleHash(in chan interface{}, out chan interface{}) {
	const op = "SingleHash"
	data := <-in

	switch data.(type) {
	case int:
		dataInt := data.(int)
		hash := md5.Sum([]byte(fmt.Sprintf("%v", dataInt)))
		md5Data := fmt.Sprintf("%x", hash)                                // cfcd208495d565ef66e7dff9f98764da
		crcMd5Data := crc32.ChecksumIEEE([]byte(md5Data))                 // 502633748
		crcData := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v", dataInt))) // 4108050209
		result := strconv.Itoa(int(crcData)) + "~" + strconv.Itoa(int(crcMd5Data))
		out <- result
	default:
		close(out)
		panic(fmt.Sprintf("%s: unknown type: %T", op, data))
	}

}
func MultiHash(in chan interface{}, out chan interface{}) {
	const op = "MultiHash"

	data := <-in

	var result string

	if data, ok := data.(string); ok {
		for th := 0; th <= 5; th++ {
			step := crc32.ChecksumIEEE([]byte(strconv.Itoa(th) + data))
			result += strconv.Itoa(int(step))

		}
	}

	out <- result
}
func CombineResults(in chan interface{}, out chan interface{}) {}
