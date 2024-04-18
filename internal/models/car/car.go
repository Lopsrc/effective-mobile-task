package car

import "effectiveM-test-task/internal/models/person"

type Car struct {
	RegNum string `json:"reg_num"`
	Mark string `json:"mark"`
	Model string `json:"model"`
	Year string `json:"year"`
	Owner person.Person `json:"owner"`
}
