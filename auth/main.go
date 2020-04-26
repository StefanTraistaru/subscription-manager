package main

import (
    "github.com/dgrijalva/jwt-go"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "strconv"
    "strings"
    "github.com/go-redis/redis"
    "github.com/twinj/uuid"

    "github.com/rs/cors"
    "github.com/gorilla/mux"
    "encoding/json"
    "io/ioutil"
)

type User struct {
    ID uint64           `json:"id"`
    Username string     `json:"username"`
    Password string     `json:"password"`
    Phone string        `json:"phone"`
}

type TokenDetails struct {
    AccessToken  string
    RefreshToken string
    AccessUuid   string
    RefreshUuid  string
    AtExpires    int64
    RtExpires    int64
}

type Todo struct {
    UserID uint64 `json:"user_id"`
    Title string `json:"title"`
}

type AccessDetails struct {
    AccessUuid string
    UserId   uint64
}

var user = User{
    ID:            1,
    Username: "username",
    Password: "password",
    Phone: "49123454322", //this is a random number
}

var client *redis.Client

func init() {
    //Initializing redis
    dsn := os.Getenv("REDIS_DSN")
    if len(dsn) == 0 {
        dsn = "redis:6379"
    }
    client = redis.NewClient(&redis.Options{
        Addr: dsn, //redis port
    })
    _, err := client.Ping().Result()
    if err != nil {
        panic(err)
    }
}

func main() {
    // router.POST("/login", Login)
    // log.Fatal(router.Run(":5001"))

    router := mux.NewRouter()

    router.HandleFunc("/login", Login).Methods("POST")
    router.HandleFunc("/logout", Logout).Methods("POST")
    router.HandleFunc("/todo", CreateTodo).Methods("POST")

    http.ListenAndServe(":5001", cors.AllowAll().Handler(router))
    log.Println("Listening on port 5001...")
}

func Login(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responseError(w, err.Error(), http.StatusBadRequest)
        return
    }

    u := &User{}
    err = json.Unmarshal(data, u)
    if err != nil {
        responseError(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    //compare the user from the request, with the one we defined:
    if user.Username != u.Username || user.Password != u.Password {
        responseError(w, "Please provide valid login details", http.StatusBadRequest)
        return
    }

    ts, err := CreateToken(user.ID)
    if err != nil {
        responseError(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    saveErr := CreateAuth(user.ID, ts)
    if saveErr != nil {
        responseError(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    tokens := map[string]string{
        "access_token":  ts.AccessToken,
        "refresh_token": ts.RefreshToken,
     }

    responseJSON(w, tokens)


}

func CreateToken(userid uint64) (*TokenDetails, error) {
    // var err error
    // //Creating Access Token
    // os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
    // atClaims := jwt.MapClaims{}
    // atClaims["authorized"] = true
    // atClaims["user_id"] = userid
    // atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
    // // atClaims["exp"] = time.Now().Add().Unix()
    // at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
    // token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
    // if err != nil {
    //     return "", err
    // }
    // return token, nil

    td := &TokenDetails{}
    td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
    td.AccessUuid = uuid.NewV4().String()

    td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
    td.RefreshUuid = uuid.NewV4().String()

    var err error
    //Creating Access Token
    os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
    atClaims := jwt.MapClaims{}
    atClaims["authorized"] = true
    atClaims["access_uuid"] = td.AccessUuid
    atClaims["user_id"] = userid
    atClaims["exp"] = td.AtExpires
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
    td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
    if err != nil {
        return nil, err
    }
    //Creating Refresh Token
    os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
    rtClaims := jwt.MapClaims{}
    rtClaims["refresh_uuid"] = td.RefreshUuid
    rtClaims["user_id"] = userid
    rtClaims["exp"] = td.RtExpires
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
    td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
    if err != nil {
        return nil, err
    }
    return td, nil
}

func CreateAuth(userid uint64, td *TokenDetails) error {
    at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
    rt := time.Unix(td.RtExpires, 0)
    now := time.Now()

    errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
    if errAccess != nil {
        return errAccess
    }
    errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
    if errRefresh != nil {
        return errRefresh
    }
    return nil
}

func ExtractToken(r *http.Request) string {
    bearToken := r.Header.Get("Authorization")
    //normally Authorization the_token_xxx
    strArr := strings.Split(bearToken, " ")
    if len(strArr) == 2 {
        fmt.Println("ExtractToken(): " + strArr[1])
        return strArr[1]
    }
    return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
    tokenString := ExtractToken(r)
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        //Make sure that the token method conform to "SigningMethodHMAC"
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("ACCESS_SECRET")), nil
    })
    fmt.Println("VerifyToken(): ")
    fmt.Printf("%+v\n",token)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    return token, nil
}

func TokenValid(r *http.Request) error {
    token, err := VerifyToken(r)
    if err != nil {
        return err
    }
    if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
        return err
    }
    return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
    token, err := VerifyToken(r)
    if err != nil {
        return nil, err
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if ok && token.Valid {
        fmt.Println("ExtractTokenMetadata: ")
        fmt.Printf("%+v\n",claims)
        accessUuid, ok := claims["access_uuid"].(string)
        if !ok {
            return nil, err
        }
        userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
        if err != nil {
            return nil, err
        }
        return &AccessDetails{
            AccessUuid: accessUuid,
            UserId:   userId,
        }, nil
    }
    return nil, err
}

func FetchAuth(authD *AccessDetails) (uint64, error) {
    userid, err := client.Get(authD.AccessUuid).Result()
    if err != nil {
        return 0, err
    }
    fmt.Println("FetchAuth(): " + userid)
    userID, _ := strconv.ParseUint(userid, 10, 64)
    return userID, nil
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
    // var td *Todo
    // if err := c.ShouldBindJSON(&td); err != nil {
    //    responseError(w, "invalid json", http.StatusUnprocessableEntity)
    //    return
    // }
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responseError(w, err.Error(), http.StatusBadRequest)
        return
    }

    td := &Todo{}
    err = json.Unmarshal(data, td)
    if err != nil {
        responseError(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }


    tokenAuth, err := ExtractTokenMetadata(r)
    fmt.Println("CreateTodo(): ")
    fmt.Printf("%+v\n",tokenAuth)
    if err != nil {
       responseError(w, "unauthorized", http.StatusUnauthorized)
       return
    }
    userId, err := FetchAuth(tokenAuth)
    if err != nil {
       responseError(w, "unauthorized", http.StatusUnauthorized)
       return
    }
    td.UserID = userId

    //you can proceed to save the Todo to a database
    //but we will just return it to the caller here:
    responseJSON(w, td)
}

func DeleteAuth(givenUuid string) (int64,error) {
    deleted, err := client.Del(givenUuid).Result()
    if err != nil {
       return 0, err
    }
    return deleted, nil
}

func Logout(w http.ResponseWriter, r *http.Request) {
    au, err := ExtractTokenMetadata(r)
    if err != nil {
        responseError(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    deleted, delErr := DeleteAuth(au.AccessUuid)
    if delErr != nil || deleted == 0 { //if any goes wrong
        responseError(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    responseError(w, "Successfully logged out", http.StatusOK)
}


func responseError(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}