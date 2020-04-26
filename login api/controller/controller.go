package controller

import (
    // "context"
    "encoding/json"
    "log"
    "fmt"
    "../db"
    // "../model"
    "../model"
    "io/ioutil"
    "net/http"

    jwt "github.com/dgrijalva/jwt-go"
    "go.mongodb.org/mongo-driver/bson"
    "golang.org/x/crypto/bcrypt"

    // "gopkg.in/mgo.v2"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: register")

    // Read body of the request
    var user model.User
    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(body, &user)
    if err != nil {
        responseError(w, "Cannot unmarshal body", err, http.StatusBadRequest)
        return
    }

    // Get connection to DB
    collection, session, err := db.GetDBCollection2()
    if err != nil {
        responseError(w, "Cannot get connection to DB", err, http.StatusInternalServerError)
        return
    }
    defer session.Close()

    // Register user in DB
    var result model.User
    err = collection.Find(bson.M{"username": user.Username}).One(&result)
    if err != nil {
        if err.Error() == "not found" {
            hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
            if err != nil {
                responseError(w, "Cannot hash password", err, http.StatusInternalServerError)
                return
            }

            user.Password = string(hash)
            err = collection.Insert(user)
            if err != nil {
                responseError(w, "Cannot create user", err, http.StatusInternalServerError)
                return
            }

            responseJSON(w, "Registration successful", "Registration successful")
            return
        }

        responseError(w, "Cannot query the DB", err, http.StatusInternalServerError)
        return
    }

    responseJSON(w, "Username already exists", "Username already exists")
    return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: login")

    // Read body of the request
    var user model.User
    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(body, &user)
    if err != nil {
        responseError(w, "Cannot unmarshal body", err, http.StatusBadRequest)
        return
    }

    // Get connection to DB
    collection, session, err := db.GetDBCollection2()
    if err != nil {
        responseError(w, "Cannot get connection to DB", err, http.StatusInternalServerError)
        return
    }
    defer session.Close()

    var result model.User

    // Query DB for user
    err = collection.Find(bson.M{"username": user.Username}).One(&result)
    if err != nil {
        if err.Error() == "not found" {
            responseJSON(w, "Invalid username", "Invalid password")
            return
        }
        responseError(w, "Query DB", err, http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
    if err != nil {
        responseJSON(w, "Invalid password", "Invalid password")
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username":  result.Username,
        "firstname": result.FirstName,
        "lastname":  result.LastName,
    })

    tokenString, err := token.SignedString([]byte("secret"))
    if err != nil {
        responseError(w, "Cannot generate token", err, http.StatusInternalServerError)
        return
    }

    responseJSON(w, "Login successful", tokenString)
}

// TODO: Keep this?
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    tokenString := r.Header.Get("Authorization")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method")
        }
        return []byte("secret"), nil
    })
    var result model.User
    var res model.ResponseResult
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        result.Username = claims["username"].(string)
        result.FirstName = claims["firstname"].(string)
        result.LastName = claims["lastname"].(string)

        json.NewEncoder(w).Encode(result)
        return
    } else {
        res.Error = err.Error()
        json.NewEncoder(w).Encode(res)
        return
    }

}


func responseError(w http.ResponseWriter, logMessage string, err error, code int) {
    log.Println("Error: " + logMessage)
    log.Println(err.Error())

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    var res model.ResponseResult
    res.Error = err.Error()
    json.NewEncoder(w).Encode(res)
}

func responseJSON(w http.ResponseWriter, logMessage string, result string) {
    log.Println(logMessage)
    w.Header().Set("Content-Type", "application/json")
    var res model.ResponseResult
    res.Result = result
    json.NewEncoder(w).Encode(res)
}