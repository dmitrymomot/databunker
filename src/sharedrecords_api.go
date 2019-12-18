package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/tidwall/gjson"
)

func (e mainEnv) newSharedRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userTOKEN := ps.ByName("token")
	event := audit("create shared record by user token", userTOKEN, "token", userTOKEN)
	defer func() { event.submit(e.db) }()

	if enforceUUID(w, userTOKEN, event) == false {
		return
	}
	if e.enforceAuth(w, r, event) == false {
		return
	}
	records, err := getJSONPostData(r)
	if err != nil {
		returnError(w, r, "failed to decode request body", 405, err, event)
		return
	}
	fields := ""
	session := ""
	partner := ""
	appName := ""
	expiration := e.conf.Policy.Max_shareable_record_retention_period
	if value, ok := records["fields"]; ok {
		if reflect.TypeOf(value) == reflect.TypeOf("string") {
			fields = value.(string)
		}
	}
	if value, ok := records["session"]; ok {
		if reflect.TypeOf(value) == reflect.TypeOf("string") {
			session = value.(string)
		}
	}
	if value, ok := records["partner"]; ok {
		if reflect.TypeOf(value) == reflect.TypeOf("string") {
			partner = value.(string)
		}
	}
	if value, ok := records["expiration"]; ok {
		if reflect.TypeOf(value) == reflect.TypeOf("string") {
			expiration = setExpiration(e.conf.Policy.Max_shareable_record_retention_period, value.(string))
		} else {
			returnError(w, r, "failed to parse expiration field", 405, err, event)
			return
		}
	}
	if value, ok := records["app"]; ok {
		if reflect.TypeOf(value) == reflect.TypeOf("string") {
			appName = strings.ToLower(value.(string))
			if len(appName) > 0 && isValidApp(appName) == false {
				returnError(w, r, "unknown app name", 405, nil, event)
			}
		} else {
			// type is different
			returnError(w, r, "failed to parse app field", 405, nil, event)
		}
	}
	if len(expiration) == 0 {
		// using default expiration time for record
		expiration = "1m"
	}
	recordUUID, err := e.db.saveSharedRecord(userTOKEN, fields, expiration, session, appName, partner)
	if err != nil {
		fmt.Println(err)
		returnError(w, r, err.Error(), 405, err, event)
		return
	}
	event.Record = userTOKEN
	event.Msg = "Generated " + recordUUID
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, `{"status":"ok","record":%q}`, recordUUID)
}

func (e mainEnv) getRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	record := ps.ByName("record")
	event := audit("get record by record", record, "record", record)
	defer func() { event.submit(e.db) }()

	if enforceUUID(w, record, event) == false {
		return
	}
	recordInfo, err := e.db.getSharedRecord(record)
	if err != nil {
		fmt.Printf("%d access denied for : %s\n", http.StatusForbidden, record)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Access denied"))
		return
	}
	var resultJSON []byte
	if len(recordInfo.token) > 0 {
		event.Record = recordInfo.token
		event.App = recordInfo.appName
		fmt.Printf("displaying fields: %s, user token: %s\n", recordInfo.fields, recordInfo.token)

		if len(recordInfo.appName) > 0 {
			resultJSON, err = e.db.getUserApp(recordInfo.token, recordInfo.appName)
		} else {
			resultJSON, err = e.db.getUser(recordInfo.token)
		}
		if err != nil {
			returnError(w, r, "internal error", 405, err, event)
			return
		}
		if resultJSON == nil {
			returnError(w, r, "not found", 405, err, event)
			return
		}
		fmt.Printf("Full user json: %s\n", resultJSON)
		if len(recordInfo.fields) > 0 {
			raw := make(map[string]interface{})
			//var newJSON json
			allFields := parseFields(recordInfo.fields)
			for _, f := range allFields {
				if f == "token" {
					raw["token"] = recordInfo.token
				} else {
					value := gjson.Get(string(resultJSON), f)
					//fmt.Printf("result %s -> %s\n", f, value)
					/*
						var raw2 map[string]interface{}
						err = json.Unmarshal([]byte(value.String()), &raw2)
						if err != nil {
							fmt.Printf("Err: %s\n", err)
						}
					*/
					raw[f] = value.Value()
				}
			}
			resultJSON, _ = json.Marshal(raw)
		}
	}
	//fmt.Fprintf(w, "<html><head><title>title</title></head>")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	var str string

	if len(recordInfo.appName) > 0 {
		str = fmt.Sprintf(`{"status":"ok","app":"%s","data":%s}`,
			recordInfo.appName, resultJSON)
	} else if len(recordInfo.session) > 0 {
		str = fmt.Sprintf(`{"status":"ok","session":"%s","data":%s}`,
			recordInfo.appName, resultJSON)
	} else {
		str = fmt.Sprintf(`{"status":"ok","data":%s}`, resultJSON)
	}

	fmt.Printf("result: %s\n", str)
	w.Write([]byte(str))
}