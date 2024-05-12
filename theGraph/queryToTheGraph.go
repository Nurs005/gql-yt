package thegraph

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/Nurs005/gql-yt/graph/model"
)

var URLS = []string{"https://gateway-arbitrum.network.thegraph.com/api/cfdf070cf997c79bcae014e7ab2bee7b/subgraphs/id/JCNWRypm7FYwV8fx5HhzZPSFaMxgkPuw4TnR3Gpi81zk", "https://gateway-arbitrum.network.thegraph.com/api/cfdf070cf997c79bcae014e7ab2bee7b/subgraphs/id/C2zniPn45RnLDGzVeGZCx2Sw3GXrbc9gL4ZfL8B8Em2j"}
var wg sync.WaitGroup

func FetchFromTheGraph(id *model.AccountFilter) (response *model.Account, err error) {
	dataFrom := model.Account{
		ID:      *id.ID,
		Raiting: "0",
	}
	for _, val := range URLS {
		wg.Add(1)
		go func(vv string, ids string) {
			defer wg.Done()
			d, _ := fetchV3(ids, vv)
			dataFrom.Borrows = append(dataFrom.Borrows, d.Borrows...)
			dataFrom.Liquidations = append(dataFrom.Liquidations, d.Liquidations...)
		}(val, *id.ID)
	}
	wg.Wait()
	addRating(&dataFrom)
	return &dataFrom, nil
}

func fetchV3(id string, url string) (*model.Account, error) {
	msg := fmt.Sprintf("{\"query\":\"query MyQuery ($id: Bytes!) {\\r\\n  accounts(where: {id: $id}) {\\r\\n    borrows {\\r\\n      amountUSD\\r\\n    }\\r\\n    liquidations {\\r\\n      amountUSD\\r\\n    }\\r\\n  }\\r\\n}\",\"variables\":{\"id\":\"%v\"}}", id)

	r := strings.NewReader(msg)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, r)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "__cf_bm=N29JxLLfRWmIoOoW7ZP6k4fshmrNrYkf.hnLRFzNnlU-1714527200-1.0.1.1-KxMFJ6deDWLZwLZT3IudKsIdICFQc0abi1liVUndBoCy.MsVp6f4m5ja_R86.xc__5famndB85LEMkH1BcFN1g")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bdy, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(bdy)
		log.Fatal(err)
	}

	jsnDta, _ := parse(bdy, id)

	return jsnDta, nil
}

func parse(d []byte, id string) (*model.Account, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(d, &result); err != nil {
		log.Fatal(err)
		return nil, err
	}

	acc := &model.Account{
		ID:      id,
		Raiting: "0",
	}
	accsData, ok := result["data"].(map[string]interface{})["accounts"].([]interface{})
	if !ok {
		return nil, errors.New("no accounts data")
	}

	var brsData []interface{}
	var liqsData []interface{}

	for _, accData := range accsData {
		accMap := accData.(map[string]interface{})
		brsData, ok = accMap["borrows"].([]interface{})
		if !ok {
			fmt.Println("no borrows")
			continue
		}
		liqsData, ok = accMap["liquidations"].([]interface{})
		if !ok {
			fmt.Println("no liquidations")
			continue
		}

	}

	for _, brData := range brsData {
		brMap := brData.(map[string]interface{})
		amUSD, ok := brMap["amountUSD"].(string)
		if !ok {
			fmt.Println("no amountUSD in borrows")
			continue
		}
		br := &model.Borrow{
			AmountUsd: amUSD,
		}
		acc.Borrows = append(acc.Borrows, br)
	}

	for _, liqData := range liqsData {
		liqMap := liqData.(map[string]interface{})
		amUSD, ok := liqMap["amountUSD"].(string)
		if !ok {
			fmt.Println("no amountUSD in liquidations")
			continue
		}
		liq := &model.Liquidate{
			AmountUsd: amUSD,
		}
		acc.Liquidations = append(acc.Liquidations, liq)
	}

	return acc, nil
}

func addRating(a *model.Account) {
	brLen := len(a.Borrows)
	fmt.Println(brLen)

	liqLen := len(a.Liquidations)
	fmt.Println(liqLen)
	// переменная рейтинга
	var ri float64 = 0

	if brLen == 0 && liqLen == 0 {
		if brLen == 0 {
			a.Raiting = "0"
		}
		return
	}
	if liqLen == 0 && brLen > 0 {
		a.Raiting = "2.5"
		return
	}
	if brLen > 0 && liqLen > 0 {
		ri = (float64(brLen) / float64(liqLen)) * 100
		strRi := strconv.FormatFloat(math.Max(0.5, (ri/100)*2.5), 'g', -1, 64)
		if ri/100*2.5 > 2.5 {
			a.Raiting = "2.5"
			return
		}
		a.Raiting = strRi
		return
	}
	if brLen == liqLen {
		a.Raiting = "2,5"
		return
	}
}
