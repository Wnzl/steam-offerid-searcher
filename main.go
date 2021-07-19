package main

import (
	"github.com/sirupsen/logrus"
	"steam-offerid-searcher/searcher"
)

type result struct {
	apiKey         string
	successOfferId string
	offer          searcher.Offer
}

func main() {
	offerRequests := []searcher.SearchOfferRequest{
		{
			Message:  "",
			ApiKey:   "",
			FailedId: 0,
		},
	}

	res := make([]result, len(offerRequests))
	for i, req := range offerRequests {
		logrus.WithField("message", req.Message).
			WithField("startId", req.FailedId).
			Info("start searching")

		s := searcher.NewSearcher(req.ApiKey, req.Message, 100, 2000)

		//todo split to goroutines with separate proxy
		offer, err := s.FindSuccessOffer(req)
		if err != nil {
			logrus.WithError(err).Error("searching success offer")
			continue
		}

		if offer.Message == req.Message {
			res[i] = result{
				apiKey:         req.ApiKey,
				successOfferId: offer.Tradeofferid,
				offer:          offer,
			}
			println(offer.Tradeofferid)
			continue
		}

		logrus.Error("successful offer not found")
	}

	for _, found := range res {
		logrus.WithField("successful_offer_id", found.successOfferId).
			WithField("apiKey", found.apiKey).
			Info("Offer found")
	}
}
