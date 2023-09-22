package models

type General struct {
	Process           string
	Directory         string
	UserAbbreviation  string
	OrgCode           string
	Mailpartner       string
	Seller            string
	Plantcode         string
	UnloadingPoint    string
	UnitOfMeasurement string
	StartingNumber    int
	NumberOfFiles     int
	NumberOfPositions int
	Currency          string
	Delivery          string
	//Invoice           string
	//OrderReply        string
	ConsigneeLevel string
	InvoiceeLevel  string
}
