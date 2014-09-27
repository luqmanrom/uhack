package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sort"
	"strings"
)

var dataPath = flag.String("data", "", "Path to data")
var indexPath = flag.String("output", "data/ir/tok_index.txt", "Output index path")

type tokFreq struct {
	tok  string
	freq int
}

type ByFreq []tokFreq

func (tf ByFreq) Len() int           { return len(tf) }
func (tf ByFreq) Less(i, j int) bool { return tf[i].freq > tf[j].freq }
func (tf ByFreq) Swap(i, j int)      { tf[i], tf[j] = tf[j], tf[i] }

func parseTokensFromLine(line string) []string {
	split := strings.Split(line, ";")
	return strings.Split(split[0], ",")
}

func main() {
	flag.Parse()
	if *dataPath == "" {
		log.Fatal("Need a path")
	}

	file, err := os.Open(*dataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	outFile, err := os.Create(*indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	token_freqs := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := parseTokensFromLine(scanner.Text())
		for _, t := range tokens {
			token_freqs[t]++
		}
	}

	token_freqs_slice := make([]tokFreq, 0)
	for token, freq := range token_freqs {
		token_freqs_slice = append(token_freqs_slice, tokFreq{tok: token, freq: freq})
	}

	sort.Sort(ByFreq(token_freqs_slice))
	out := bufio.NewWriter(outFile)
	if len(token_freqs_slice) > 1000 {
		token_freqs_slice = token_freqs_slice[:1000]
	}
	for _, v := range token_freqs_slice {
		log.Print("Here")
		out.Write([]byte(v.tok + "\n"))
	}
	out.Flush()
}
