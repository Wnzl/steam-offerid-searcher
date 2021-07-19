package searcher

type SearchOfferRequest struct {
	Message  string
	ApiKey   string
	FailedId int
}

type GetOfferResponse struct {
	Response struct {
		Offer `json:"offer"`
	} `json:"response"`
}

type Offer struct {
	Tradeofferid    string `json:"tradeofferid"`
	AccountidOther  int    `json:"accountid_other"`
	Message         string `json:"message"`
	ExpirationTime  int    `json:"expiration_time"`
	TradeOfferState int    `json:"trade_offer_state"`
	ItemsToReceive  []struct {
		Appid      int    `json:"appid"`
		Contextid  string `json:"contextid"`
		Assetid    string `json:"assetid"`
		Classid    string `json:"classid"`
		Instanceid string `json:"instanceid"`
		Amount     string `json:"amount"`
		Missing    bool   `json:"missing"`
	} `json:"items_to_receive"`
	IsOurOffer         bool `json:"is_our_offer"`
	TimeCreated        int  `json:"time_created"`
	TimeUpdated        int  `json:"time_updated"`
	FromRealTimeTrade  bool `json:"from_real_time_trade"`
	EscrowEndDate      int  `json:"escrow_end_date"`
	ConfirmationMethod int  `json:"confirmation_method"`
}

func (o *Offer) isEmpty() bool {
	return o.Tradeofferid == ""
}
