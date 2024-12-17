package entity

type GetCypherPayload struct {
	Token string
}

type RotatePayload struct {
	Max             int
	BatchSize       int
	DayDifference   int
	MaxAsyncProcess int
	MsDelayEachJob  int
}
