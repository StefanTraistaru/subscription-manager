package controller

import (
    "encoding/json"
    "log"
    "io/ioutil"
    "net/http"
    "../db"
    "../model"

    jwt "github.com/dgrijalva/jwt-go"
    "gopkg.in/mgo.v2/bson"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/mux"
)

////////////////////////////////////////////////////
// Handler functions
////////////////////////////////////////////////////

func Register(w http.ResponseWriter, r *http.Request) {
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
    collection, session, err := db.GetDBCollection()
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

func Login(w http.ResponseWriter, r *http.Request) {
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
    collection, session, err := db.GetDBCollection()
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

func Delete(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: delete")

    // Read body
    params := mux.Vars(r)
    id := params["username"]

    // Get connection to DB
    collection, session, err := db.GetDBCollection()
    if err != nil {
        responseError(w, "Cannot get connection to DB", err, http.StatusInternalServerError)
        return
    }
    defer session.Close()

    // Query DB for user
    err = collection.Remove(bson.M{"username": id})
    if err != nil {
        if err.Error() == "not found" {
            responseJSON(w, "Invalid user", "Invalid user")
            return
        }
        responseError(w, "Query DB", err, http.StatusInternalServerError)
        return
    }

    responseJSON(w, "Deletion successful", "Deletion successful")
}


////////////////////////////////////////////////////
// Util functions
////////////////////////////////////////////////////

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