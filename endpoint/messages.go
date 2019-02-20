package endpoint

type NotificationMsg struct {
	Data string `json:"data"`
}

type Notification struct {
	Date   uint64 `json:"date"`
	Data   string `json:"data"`
	ToAddr string `json:"toAddr"`
}

type GetNotificationMsg struct {
	IdAddr string `json:"idAddr"`
	// SignedPacket string `json:"signedPacket"`
	// ProofKSign ProofClaim `json:"proofKSign"`
}
