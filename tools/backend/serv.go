package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/fzzy/radix/redis"
)

var redisAddr = flag.String("redis", "172.16.235.139:6379", "Redis Server address (IP:Port)")
var indexPath = flag.String("index", "data/ir/tok_index.txt", "Index path")
var frontendPath = flag.String("frontend", "frontend/saver/app", "Path to frontend")
var redisClient *redis.Client
var junkRegexp = regexp.MustCompile("[^0-9a-z]")
var validToken map[string]bool

type totalMessage struct {
	Store     string `json:"store"`
	TotalCost int64  `json:"totalCost"`
}

type rpcMessage struct {
	Totals    []*totalMessage `json:"totals"`
	BadTokens []string        `json:"badTokens"`
}

func getTokens(q string) ([]string, []string) {
	var good []string
	var bad []string
	split := strings.Split(strings.ToLower(q), " ")
	for i := 0; i < len(split); i++ {
		t := string(junkRegexp.ReplaceAll([]byte(split[i]), nil))
		if t == "" {
			continue
		}
		if validToken[t] {
			good = append(good, t)
		} else {
			bad = append(bad, t)
		}
	}
	return good, bad
}

func getValidTokens(r io.Reader) map[string]bool {
	tokens := make(map[string]bool)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		tokens[scanner.Text()] = true
	}
	return tokens
}

func getStoresAndPrices(queried []string) ([]*totalMessage, error) {
	storePrices := make(map[string][]int64)
	for _, q := range queried {
		resp := redisClient.Cmd("lrange", q, 0, 10000)
		if resp.Err != nil {
			return nil, resp.Err
		}
		stores, err := resp.List()
		if err != nil {
			return nil, err
		}
		for _, store := range stores {
			resp = redisClient.Cmd("get", q + "," + store)
			if resp.Err != nil {
				return nil, resp.Err
			}
			val, err := resp.Int64()
			if err != nil {
				return nil, err
			}
			storePrices[store] = append(storePrices[store], val)
		}
	}
	var finalCosts []*totalMessage
	for k, v := range storePrices {
		if len(v) < len(queried) {
			continue
		}
		var total totalMessage
		total.Store = k
		for _, c := range v {
			total.TotalCost += c
		}
		finalCosts = append(finalCosts, &total)
	}
	return finalCosts, nil
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	log.Printf("Got query: %s", q)
	if q == "" {
		http.Error(w, "Must provide query", http.StatusBadRequest)
		return
	}
	tokens, badTokens := getTokens(q)
	if len(tokens) == 0 {
		if len(badTokens) == 0 {
			http.Error(w, "Must provide a valid query", http.StatusBadRequest)
			return
		}
		resp := rpcMessage{Totals: nil, BadTokens: badTokens}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
			return
		}
		w.Write(jsonResp)
		return
	}
	totalMessages, err := getStoresAndPrices(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	resp := rpcMessage{Totals: totalMessages, BadTokens: badTokens}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With")
	w.Write(jsonResp)
}

func main() {
	flag.Parse()
	var err error

	redisClient, err = redis.Dial("tcp", *redisAddr)
	defer redisClient.Close()
	if err != nil {
		log.Fatal("Could not connect to redis")
	}

	indexFile, err := os.Open(*indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer indexFile.Close()
	validToken = getValidTokens(indexFile)

	http.HandleFunc("/rpc", rpcHandler)
	http.Handle("/", http.FileServer(http.Dir(*frontendPath)))
	http.ListenAndServe(":8080", nil)
}
