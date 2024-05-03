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

	"github.com/Nurs005/gql-yt/graph/model"
)

const GRAPHURI = "https://gateway-arbitrum.network.thegraph.com/api/cfdf070cf997c79bcae014e7ab2bee7b/subgraphs/id/JCNWRypm7FYwV8fx5HhzZPSFaMxgkPuw4TnR3Gpi81zk"
const GRAPHV2 = "https://gateway-arbitrum.network.thegraph.com/api/cfdf070cf997c79bcae014e7ab2bee7b/subgraphs/id/C2zniPn45RnLDGzVeGZCx2Sw3GXrbc9gL4ZfL8B8Em2j"

func FetchFromTheGraph(where *model.AccountFilter) ([]*model.Account, error) {
	url := GRAPHURI
	method := "POST"
	message := fmt.Sprintf("{\"query\":\"query MyQuery ($id: Bytes!) {\\r\\n  accounts(where: {id: $id}) {\\r\\n    borrows {\\r\\n      amountUSD\\r\\n    }\\r\\n    liquidations {\\r\\n      amountUSD\\r\\n    }\\r\\n  }\\r\\n}\",\"variables\":{\"id\":\"%v\"}}", *where.ID)

	payload := strings.NewReader(message)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(body)
		log.Fatal(err)
	}

	jsonData, _ := parseJson(body, string(*where.ID))
	dataFromV2, _ := fetchFromV2(where)
	jsonData.Borrows = append(jsonData.Borrows, dataFromV2.Borrows...)
	jsonData.Liquidations = append(jsonData.Liquidations, dataFromV2.Liquidations...)
	newData := addRating([]*model.Account{jsonData})
	return newData, nil
}

func parseJson(data []byte, id string) (*model.Account, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err)
		return nil, err
	}

	account := &model.Account{
		ID:      id,
		Raiting: "0",
	}
	accountsData, ok := result["data"].(map[string]interface{})["accounts"].([]interface{})
	if !ok {
		return nil, errors.New("no accounts data")
	}

	for _, accountData := range accountsData {
		accountMap := accountData.(map[string]interface{})
		borrowsData, ok := accountMap["borrows"].([]interface{})
		if !ok {
			fmt.Println("no borrows")
			continue
		}

		for _, borrowData := range borrowsData {
			borrowMap := borrowData.(map[string]interface{})
			amountUSD, ok := borrowMap["amountUSD"].(string)
			if !ok {
				fmt.Println("no amountUSD in borrows")
				continue
			}
			borrow := &model.Borrow{
				AmountUsd: amountUSD,
			}
			account.Borrows = append(account.Borrows, borrow)
		}

		liquidationsData, ok := accountMap["liquidations"].([]interface{})
		if !ok {
			fmt.Println("no liquidations")
			continue
		}

		for _, liquidationData := range liquidationsData {
			liquidationMap := liquidationData.(map[string]interface{})
			amountUSD, ok := liquidationMap["amountUSD"].(string)
			if !ok {
				fmt.Println("no amountUSD in liquidations")
				continue
			}
			liquidation := &model.Liquidate{
				AmountUsd: amountUSD,
			}
			account.Liquidations = append(account.Liquidations, liquidation)
		}
	}

	return account, nil
}

func addRating(account []*model.Account) []*model.Account {
	for _, acc := range account {
		borrowsLen := len(acc.Borrows)
		fmt.Println(borrowsLen)

		liquidateLen := len(acc.Liquidations)
		fmt.Println(liquidateLen)
		var rating float64 = 0

		if borrowsLen == 0 {
			if borrowsLen == 0 {
				acc.Raiting = "0"
			}
			return account
		} else if liquidateLen == 0 {
			if liquidateLen == 0 {
				acc.Raiting = "5"
			}
			return account
		}
		if borrowsLen > 0 && liquidateLen > 0 {
			if borrowsLen > liquidateLen {
				rating = (float64(borrowsLen) / float64(liquidateLen)) * 100
				newRating := strconv.FormatFloat(rating/100*5, 'f', 2, 64)
				if rating/100*5 > 5.0 {
					acc.Raiting = "5"
					return account
				}
				acc.Raiting = newRating
				return account
			} else if liquidateLen > borrowsLen {
				rating := float64(borrowsLen) / float64(liquidateLen) * 100
				newRating := strconv.FormatFloat(math.Max(0.5, (1-rating/100)*5), 'f', 2, 64)
				acc.Raiting = newRating
				return account
			}
		}
		if borrowsLen == liquidateLen {
			acc.Raiting = "2,5"
			return account
		}
	}

	return account
}

func fetchFromV2(where *model.AccountFilter) (*model.Account, error) {
	url := GRAPHV2
	method := "POST"
	message := fmt.Sprintf("{\"query\":\"query MyQuery ($id: Bytes!) {\\r\\n  accounts(where: {id: $id}) {\\r\\n    borrows {\\r\\n      amountUSD\\r\\n    }\\r\\n    liquidations {\\r\\n      amountUSD\\r\\n    }\\r\\n  }\\r\\n}\",\"variables\":{\"id\":\"%v\"}}", *where.ID)

	payload := strings.NewReader(message)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(body)
		log.Fatal(err)
	}

	jsonData, _ := parseJson(body, string(*where.ID))
	return jsonData, nil
}
