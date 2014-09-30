package openid

import (
	"fmt"
	"net/http"
	"encoding/json"
	"appengine"
	"appengine/datastore"
	"io/ioutil"
)

type SearchTerms struct {
	Searchterm  string
	Uid 	    string
}

type Users struct {
	Username  string
	Password  string
	Userid int64
	Email string
	Time int64
	Url string
	Type string
	Name string
	Contacts string
}

type JsonReply struct {
	Status	   string
	Uid		   string
	Sid		   string
}

func init() {
	http.HandleFunc("vidao.com/",hello)
    http.HandleFunc("vidao.com/search/testsearchuser", testsearch)
    http.HandleFunc("vidao.com/search/user", search)
}

func hello (w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "Coming Soon")
}

func testsearch (w http.ResponseWriter, r *http.Request){
const guestbookForm = `
		<html>
			<body>
				<form action="/search/user" method="post">
				<div>Username:<input type="text" name="searchterm"></div><br>
				<div><input type="submit" value="Search"></div>
				</form>
			</body>
		</html>
		`
		fmt.Fprint(w, guestbookForm)
}

func sendJson (w *http.ResponseWriter, r *http.Request, message string, uid string, sid string){
	
	var jsreply JsonReply
	jsreply.Status=message
	jsreply.Sid=sid
	jsreply.Uid=uid
	js, err := json.Marshal(jsreply)
	if err != nil {
		http.Error((*w), err.Error(), http.StatusInternalServerError)
		return
	}
	(*w).Header().Add("Content-Type", "application/json")
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("X-Requested-With","XMLHttpRequest")
	(*w).Write(js)
    return
}

func sendResult (w *http.ResponseWriter, r *http.Request, result []Users){
	
	for i := range result {
        result[i].Password = "Censored"
    }
    
	js, err := json.Marshal(result)
	if err != nil {
		http.Error((*w), err.Error(), http.StatusInternalServerError)
		return
	}
	(*w).Header().Add("Content-Type", "application/json")
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
	(*w).Write(js)
    return
}

func search (w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	
    var inputsubstring SearchTerms
   
    c := appengine.NewContext(r)
    
    if r.Header.Get("Content-Type") == "application/json;charset=UTF-8"{
		jsonBinInfo, err := ioutil.ReadAll(r.Body)
	
		r.Body.Close()
		if err != nil {
			fmt.Fprintln(w,err)
		}
	

		var jsreply SearchTerms
	
		Jerr := json.Unmarshal(jsonBinInfo, &jsreply)
		if Jerr != nil {
			fmt.Fprintln(w,Jerr)
		}

		inputsubstring.Uid= jsreply.Uid
		inputsubstring.Searchterm= jsreply.Searchterm
		if inputsubstring.Searchterm == ""{
			sendJson (&w, r,"Empty String", "-1", "0")
			return
		}

	}else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded"{
	
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		inputsubstring.Uid = r.FormValue("uid")
		inputsubstring.Searchterm = r.FormValue("searchterm")


	}else{
		sendJson (&w, r,"Invalid Input Format", "0", "0")
		return
	}
	

	qClient := datastore.NewQuery("User").
                Filter("Username >=", inputsubstring.Searchterm).
                Filter("Username <=", inputsubstring.Searchterm+"\ufffd").Order("Username").
                Limit(5)
                
     var client []Users
     if _, err := qClient.GetAll(c, &client); client==nil{ 
     	sendJson (&w, r,"No User Found", "-1", "0")
     	return
     }else if err!=nil{
     	fmt.Fprint(w,err)
     }else{
		sendResult (&w, r, client)
     }
}

