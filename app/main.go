package main

import (
	"bufio"
	"container/heap"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"example.com/datagen"
	"example.com/datasorter"
	"example.com/kconnect"
	"example.com/logging"
)

var (
	buffer         []datasorter.Item
	numReceived    int64
	numChunks      int64
	totalCount     int64
	sortedID       int64
	baseDir        string
	finalOutput    int
	sortingColumns []string
)

func clearBaseDir() error {
	d, err := os.Open(baseDir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if strings.HasSuffix(name, ".log") {
			continue
		}
		err = os.RemoveAll(filepath.Join(baseDir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func sortAndWriteToTxt(slice []datasorter.Item, basedir string, sortingcolumn string, idx int64) {
	// sorting
	start := time.Now()
	switch sortingcolumn {
	case "name":
		datasorter.SortByName(slice)
	case "continent":
		datasorter.SortByContinent(slice)
	default:
		datasorter.SortByID(slice)
	}

	// write
	filename := basedir + "/" + sortingcolumn + "_" + strconv.FormatInt(idx, 10) + ".txt"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	linewriter := bufio.NewWriter(file)

	for _, data := range slice {
		_, _ = linewriter.WriteString(data.Data + "\n")
	}

	linewriter.Flush()
	file.Close()
	elapsed := time.Since(start)
	logging.TimingLogger.Printf("Time for chunk sorting and saving by %-12v (round %-2v) is %s \n", sortingcolumn, idx, elapsed)
}

func procMsg(partition int32, offset int64, key string, msg string) {
	buffer = append(buffer, datasorter.Item{msg, int(partition)})
	if int64(len(buffer)) >= totalCount/numChunks {
		for _, sc := range sortingColumns {
			logging.InfoLogger.Printf("Round %-5v sorting by %-12v \n", sortedID, sc)
			sortAndWriteToTxt(buffer, baseDir, sc, sortedID)
		}
		buffer = buffer[:0]

		// this method is invoked in a synchronized section (via a channel), no need for atomic operation
		if sortedID++; sortedID >= numChunks-1 {
			outCh := make(chan int, len(sortingColumns))
			var wg sync.WaitGroup
			wg.Add(len(sortingColumns))
			for _, sc := range sortingColumns {
				go mergeSortedFilesAndSend(baseDir, sc, outCh, &wg)
			}
			go func() {
				wg.Wait()
				close(outCh)
				p, _ := os.FindProcess(os.Getpid())
				p.Signal(syscall.SIGINT)
			}()
		}
	}
}

// consume messages from the source topic, and save the content into multiple
func consumeSource(kaddr string, basedir string, topic string, count int64, numchunk int64) {
	totalCount = count
	numReceived = 0
	numChunks = numchunk
	sortedID = 0
	baseDir = basedir
	buffer = make([]datasorter.Item, 0, count/numchunk+10)

	// ensure that the base dir exists
	_ = os.MkdirAll(baseDir, 0755)
	// clear baseDir
	clearBaseDir()

	consumer, err := kconnect.InitConsumer(kaddr)
	if err != nil {
		panic(err)
	}
	kconnect.Consume(topic, procMsg, consumer)
	if err := consumer.Close(); err != nil {
		panic(nil)
	}
}

func columnValue(sortingcolumn string, msg string) string {
	slice := strings.Split(msg, ",")
	switch sortingcolumn {
	case "name":
		return slice[1]
	case "continent":
		return slice[3]
	default:
		return slice[0]
	}
}

func sortingHeap(sortingcolumn string, karr []datasorter.Item) heap.Interface {
	var h heap.Interface
	switch sortingcolumn {
	case "name":
		c := datasorter.ByName(karr)
		h = &c
	case "continent":
		c := datasorter.ByContinent(karr)
		h = &c
	default:
		c := datasorter.ByID(karr)
		h = &c
	}
	return h
}

// merge the data in the sorted files according to given sorting column and send the final result to specified topic
func mergeSortedFilesAndSend(basedir string, sortingcolumn string, outCh chan<- int, wg *sync.WaitGroup) {
	start := time.Now()
	matches, _ := filepath.Glob(basedir + "/" + sortingcolumn + "_*.txt")
	k := len(matches)
	files := make([]*os.File, k)
	scanners := make([]*bufio.Scanner, k)
	karr := []datasorter.Item{}
	// open all files that are sorted by the specified column,
	// take the first line of each file and put them into a minHeap.
	for i := 0; i < k; i++ {
		filename := matches[i]
		file, err := os.Open(filename)
		if err == nil {
			files[i] = file
		}
		if files[i] != nil {
			scanners[i] = bufio.NewScanner(files[i])
			if scanners[i].Scan() {
				karr = append(karr, datasorter.Item{scanners[i].Text(), i})
			}
		}
	}

	outputProducer, err := kconnect.InitProducer("127.0.0.1:9092")
	if err != nil {
		panic(err)
	}

	h := sortingHeap(sortingcolumn, karr)

	heap.Init(h)
	for h.Len() > 0 {
		min := heap.Pop(h).(datasorter.Item)
		kconnect.Publish(sortingcolumn, columnValue(sortingcolumn, min.Data), min.Data, outputProducer)
		minpos := min.Source
		if scanners[minpos].Scan() {
			heap.Push(h, datasorter.Item{scanners[minpos].Text(), minpos})
		}
	}
	elapsed := time.Since(start)
	logging.TimingLogger.Printf("[TIME] Time elapsed for external sorting and sending on %-12v is %s \n", sortingcolumn, elapsed)
	outCh <- h.Len()
	wg.Done()
}

func main() {

	start := time.Now()

	function := flag.String("func", "con", "functionality of the process")
	kaddr := flag.String("kaddr", "127.0.0.1:9092", "address of kafka")
	topic := flag.String("topic", "source", "source topic")
	count := flag.Int64("count", 100000000, "total number of messages")
	basedir := flag.String("basedir", "/tmp/go_sorting", "the directory where temporary files should be created")
	numchunk := flag.Int64("numchunk", 10, "total number of sorted file chunks")
	flag.Parse()

	logging.Init(*basedir)
	sortingColumns = []string{"id", "name", "continent"}

	if *function == "con" {
		consumeSource(*kaddr, *basedir, *topic, *count, *numchunk)
	} else {
		datagen.Generate(*count, *topic, *kaddr)
	}
	elapsed := time.Since(start)
	logging.TimingLogger.Printf("Total elapsed time: %s\n", elapsed)
}
