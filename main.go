package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"steam-offerid-searcher/searcher"
)

type result struct {
	apiKey         string
	successOfferId string
	offer          searcher.Offer
}

func main() {
	apiKey := flag.String("key", "", "[Required] Provide api key to search trade offer")
	message := flag.String("m", "", "[Required] Provide message to identify your trade offer")
	failedID := flag.Int("id", 0, "[Required] offer id from search")
	workers := flag.Int("w",  5, "Workers count")
	flag.Parse()

	if *apiKey == "" || *message == "" {
		flag.Usage()
		logrus.Fatalf("Please provide required flags\n")
	}

	offerRequests := []searcher.SearchOfferRequest{
		{
			Message:  *message,
			ApiKey:   *apiKey,
			FailedId: *failedID,
		},
	}

	var (
		found = make(chan searcher.Offer)

		res = make([]result, len(offerRequests))
	)


	for i, req := range offerRequests {
		logrus.WithField("message", req.Message).
			WithField("startId", req.FailedId).
			Info("start searching")

		s := searcher.NewSearcher(req.ApiKey, req.Message, 100, 2000)

		offer, err := s.FindSuccessOffer(req, workers, )
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
