package controllers

import (
	"main/database"
	m "main/models"
	c "main/constants"
	//"main/utils"
	//"bytes"
	"encoding/json"
	//"fmt"
	//"github.com/satori/go.uuid"
	"net/http"
	//"strconv"
	//"os"
	//"strings"
	 "fmt"
	 "errors"
	 //"io/ioutil"
	 //"log"
)

// /api/v1/auth?code=<code>
func VerifyAuthCodeAndLogin(w http.ResponseWriter, r *http.Request){
	code := r.URL.Query().Get("code")
	if len(code) == 0{//TODO:code pattern validity check
		respondWithError(w, http.StatusBadRequest, "Code Is invalid")
		return
	}

	url := c.VarifyAuthCode + "&code=" + code + "&access_token=AA|" + c.FacebookAppID + "|" + c.FacebookSecretKey
	method := "GET"

	client := &http.Client {
	    CheckRedirect: func(req *http.Request, via []*http.Request) error {
	    	return http.ErrUseLastResponse
	    },
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
	   	fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not make request to authorization service api1")
		return
	}
	if res.StatusCode != 200{
		respondWithError(w, http.StatusBadRequest, "Invalid verification code")
		return	
	}
	var by_code map[string]string
	if err = json.NewDecoder(res.Body).Decode(&by_code); err != nil{
		respondWithError(w, http.StatusBadRequest, "Could not parse json from authorization service api1")
		return
	}
	user, err := Authorize(by_code["access_token"])
	var admin string
	if isAdmin(user){
		admin = "true"
	}else{
		admin  = "false"
	}	
	respondWithJSON(w,http.StatusOK, map[string]interface{}{"message":OK, "access_token" : by_code["access_token"], "user":user, "admin":admin})//if user is empty, then register
}

//api/v1/register
func Register(w http.ResponseWriter, r *http.Request) {
	var user m.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(user.Name) == 0 || len(user.Phone) == 0{
		respondWithError(w, http.StatusBadRequest, "Empty field is not allowed")
		return
	}
	db := database.GetDB().Create(&user)
	err = db.Error
	rowsNum := db.RowsAffected
	if err != nil || rowsNum == 0 {
		respondWithError(w, http.StatusBadRequest, "Unknown error." + err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "Successfully created.", "user": user})
}
func isAdmin(user m.User)bool{
	//return true
	var admin m.Admin
	if database.GetDB().Where("user_id = ?", user.Id).Find(&admin).RecordNotFound(){
		return false
	}
	return true
}
func Authorize(access_token string) (m.User, error){
	//return m.User{},nil
	url := c.GetPhoneByToken + "?access_token="+ access_token
	method := "GET"
	client := &http.Client {
	   	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	      	return http.ErrUseLastResponse
	    },
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200{
		return m.User{}, errors.New("bad request to fb auth service")
	}
	defer res.Body.Close()
	var verfy_token map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&verfy_token); err != nil{
		return m.User{}, errors.New("can't parse json")
	}
	user := m.User{Phone : (verfy_token["phone"].(map[string]string))["number"]} 
	if !database.GetDB().Where(&user).Take(&user).RecordNotFound(){
		return m.User{}, errors.New("user is not registered")
	}
	return user,nil
}