package person

type PersonSave struct {
    Name string `json:"name"`
    Surname string `json:"surname"`
    Patronymic string `json:"patronymic"`
}

type PersonGet struct {
    Id int `json:"id"`
}

type PersonUpdate struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Surname string `json:"surname"`
    Patronymic string `json:"patronymic"`
}

type PersonDelete struct {
    Id int `json:"id"`
}

type PersonRecover struct {
    Id int `json:"id"`
}