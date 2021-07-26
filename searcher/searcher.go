package searcher

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const GetTradeOfferURL = "https://api.steampowered.com/IEconService/GetTradeOffer/v1/"

type Searcher struct {
	apiKey  string
	message string

	stepSize int
	maxLimit int
}

func NewSearcher(apiKey string, message string, step int, limit int) *Searcher {
	return &Searcher{
		apiKey:   apiKey,
		message:  message,
		stepSize: step,
		maxLimit: limit,
	}
}

func (s *Searcher) FindSuccessOffer(in ,  found chan Offer) (Offer, error) {
	asc := s.newSearcher(true)
	desc := s.newSearcher(false)
	limit := s.stepSize

	for {
		offer, err := asc(p.FailedId, limit)
		if err != nil {
			return Offer{}, err
		}

		if offer.Message == p.Message && offer.TradeOfferState == 3 {
			return offer, nil
		}

		offer, err = desc(p.FailedId, limit)
		if err != nil {
			return Offer{}, err
		}

		if offer.Message == p.Message && offer.TradeOfferState == 3 {
			return offer, nil
		}

		limit += s.stepSize

		if limit > s.maxLimit {
			break
		}
	}

	return Offer{}, nil
}

func (s *Searcher) newSearcher(asc bool) func(startId int, limit int) (successOffer Offer, err error) {
	iteration := 1
	return func(startId int, limit int) (successOffer Offer, err error) {
		var id int
		for {
			if asc {
				id = startId + iteration
			} else {
				id = startId - iteration
			}

			offer, err := s.getOfferById(id)
			if err != nil {
				return Offer{}, err
			}

			if !offer.isEmpty() {
				logrus.WithField("offerId", id).
					WithField("apiKey", s.apiKey).
					WithField("isEmpty", offer.isEmpty()).
					Info("found not empty offer")
			}

			if offer.Message == s.message && offer.TradeOfferState == 3 {
				successOffer = offer
				logrus.WithField("offerId", id).
					WithField("apiKey", s.apiKey).
					WithField("ascending", asc).
					WithField("iteration number", iteration).
					Info("found successful offer")
				return successOffer, nil
			}

			if iteration > limit {
				break
			}

			iteration++
			time.Sleep(time.Millisecond * 500)
		}

		return
	}
}

func (s *Searcher) getOfferById(id int) (Offer, error) {
	r, err := s.sendGetTradeOfferRequest(id)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	var response GetOfferResponse
	err = json.Unmarshal(r, &response)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	return response.Response.Offer, nil
}

func (s *Searcher) sendGetTradeOfferRequest(id int) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", GetTradeOfferURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	q := req.URL.Query()
	q.Add("tradeofferid", strconv.Itoa(id))
	q.Add("key", s.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("[Do request] err: %s", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("[Read resp body] err: %s", err)
	}

	return bodyBytes, nil
}
