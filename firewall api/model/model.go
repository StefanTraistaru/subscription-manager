package model

type User struct {
    Username  string `json:"username"`
    FirstName string `json:"firstname"`
    LastName  string `json:"lastname"`
    Password  string `json:"password"`
    Token     string `json:"token"`
}

type Subscription struct {
    ID      string        `json:"id"`
    Name    string        `json:"name"`
    Price   string        `json:"price"`
    Details string        `json:"details"`
    Date_d  string        `json:"date_d"`
    Date_m  string        `json:"date_m"`
    Date_y  string        `json:"date_y"`
}

type ResponseResult struct {
    Error  string         `json:"error"`
    Result string         `json:"result"`
    Data   []Subscription `json:"data"`
}