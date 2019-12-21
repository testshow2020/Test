package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// query parameters startDate , numOfSessions , days
func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Error while parsing form"}"}`))
		return
	}

	scheduleForm := &ScheduleForm{
		NumOfSessions: r.Form.Get("numOfSessions"),
		Days:          r.Form.Get("days"),
		StartDate:     r.Form.Get("startDate")}

	if !scheduleForm.Validate() {
		resErr, _ := json.Marshal(scheduleForm.Errors)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "query parameters is not valid", "errors": %s }`, resErr)))
		return
	}

	w.WriteHeader(http.StatusOK)
	resData, _ := json.Marshal(scheduleForm.ValidatedData.getSchedule())
	w.Write([]byte(resData))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not a valid url"}`))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/schedule", get).Methods(http.MethodGet)
	r.HandleFunc("/", notFound)
	log.Fatal(http.ListenAndServe(":8080", r))
}
