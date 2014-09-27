package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fzzy/radix/redis"
)

var dataPath = flag.String("data", "", "Path to data")
var indexPath = flag.String("index", "data/ir/tok_index.txt", "Index path")
var redisAddr = flag.String("redis", "172.16.235.139:6379", "Redis Server address (IP:Port)")

type tokenPair struct {
	tok   string
	price int64
}

type rawItems []*tokenPair

func (r rawItems) Len() int           { return len(r) }
func (r rawItems) Less(i, j int) bool { return r[i].tok < r[j].tok }
func (r rawItems) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func getTokens(r io.Reader) map[string]bool {
	tokens := make(map[string]bool)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		tokens[scanner.Text()] = true
	}
	return tokens
}

func genTokenPrices(r io.Reader, valid map[string]bool) (rawItems, error) {
	scanner := bufio.NewScanner(r)
	var items rawItems
	for scanner.Scan() {
		pair := strings.Split(scanner.Text(), ";")
		price, err := strconv.ParseInt(pair[1], 10, 64)
		if err != nil {
			return nil, err
		}
		for _, tok := range strings.Split(pair[0], ",") {
			if !valid[tok] {
				continue
			}
			items = append(items, &tokenPair{tok: tok, price: price})
		}
	}
	return items, nil
}

func inList(item string, list []string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func pushCombinedPrices(raw rawItems, c *redis.Client, storeName string) error {
	curTok := ""
	start := 0
	for i := 0; i < len(raw); i++ {
		if raw[i].tok == curTok {
			continue
		}
		if start != i {
			avg := int64(0)
			for _, item := range raw[start:i] {
				avg += item.price
			}
			avg /= int64(i - start)
			resp := c.Cmd("set", curTok + "," + storeName, avg)
			if resp.Err != nil {
				return resp.Err
			}

			resp =  c.Cmd("lrange", curTok, 0, 10000)
			if resp.Err != nil {
				return resp.Err
			}

/*			if resp.Type == redis.NilReply {
				resp = c.Cmd("lpush", curTok, storeName)
				if resp.Err != nil {
					return resp.Err
				}
			} else {*/
//				stores, err := resp.List()
//				if err != nil {
//					return err
//				}
//				if !inList(storeName, stores) {
					resp = c.Cmd("lpush", curTok, storeName)
//				}
		//	}
		}
		start = i
		curTok = raw[i].tok
	}
	return nil
}

func main() {
	flag.Parse()
	if *dataPath == "" {
		log.Fatal("Need a path")
	}

	redisClient, err := redis.Dial("tcp", *redisAddr)
	defer redisClient.Close()
	if err != nil {
		log.Fatal("Could not connect to redis")
	}

	dataFile, err := os.Open(*dataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dataFile.Close()

	indexFile, err := os.Open(*indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer indexFile.Close()
	tokens := getTokens(indexFile)

	raw, err := genTokenPrices(dataFile, tokens)
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(raw)
	split := strings.Split(*dataPath, "/")
	store_txt := split[len(split)-1]
	storeName := store_txt[:len(store_txt)-4]
	err = pushCombinedPrices(raw, redisClient, storeName)
	if err != nil {
		log.Fatal(err)
	}
}
