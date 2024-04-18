package car

type CarSave struct {
	RegNums  []string `json:"regNums"`
}

type CarGet struct {
	RegNum string
}

type CarUpdate struct {
	Id     int
	RegNum string
	Mark   string
	Model  string
	Year   string
	OwnerID int
}

type CarDelete struct {
	Id     int
	RegNum string
}

type CarRecover struct {
	Id     int
	RegNum string
}
