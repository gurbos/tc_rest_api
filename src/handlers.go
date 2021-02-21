package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	tcm "github.com/gurbos/tcmodels"
	"gorm.io/gorm/logger"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {
	var tca TradingCardAPI
	tca.Init()
	jbuff, err := json.Marshal(tca)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(jbuff)
}

func ProductLineHandler(w http.ResponseWriter, r *http.Request) {
	var ds DataSourceName = GetDataSource()               // Get database info
	dbConn := DBConnection(ds.DSNString(), logger.Silent) // Get database connection handle

	w.Header().Set("Content-type", "application/json")

	var jbuff []byte
	var productLines []tcm.ProductLine
	var productLineReps []ProductLineRep

	q := r.URL.Query()
	plNameList, present := q["productLineNames"]
	plNames := strings.Join(plNameList[:], ",")

	switch {
	case !present, len(plNames) == 0:
		err := dbConn.Model(tcm.ProductLine{}).Find(&productLines).Error
		if err == nil {
			productLineReps = make([]ProductLineRep, len(productLines), len(productLines))
			for i, elem := range productLines {
				productLineReps[i].Set(elem)
			}
			jbuff, err = json.Marshal(productLineReps)
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
			w.Write([]byte(""))
		}
	case present:
		err := dbConn.Model(tcm.ProductLine{}).Where("name IN ?", plNameList).Find(&productLines).Error
		if err == nil {
			productLineReps = make([]ProductLineRep, len(productLines), len(productLines))
			for i, elem := range productLines {
				productLineReps[i].Set(elem)
			}
			jbuff, err = json.Marshal(productLineReps)
			w.WriteHeader(200)
		}

	}

	w.Write(jbuff)
}

func CardSetHandler(w http.ResponseWriter, r *http.Request) {
	var ds DataSourceName = GetDataSource()
	dbConn := DBConnection(ds.DSNString(), logger.Silent)
	vars := mux.Vars(r)

	var productLineInfo tcm.ProductLine
	err := dbConn.Model(tcm.ProductLine{}).Where("name = ?", vars["productLine"]).Find(&productLineInfo).Error
	if err != nil {
		log.Fatal(err)
	}

	var jbuff []byte
	var setInfos []tcm.SetInfo
	var cardSetReps []CardSetRep

	q := r.URL.Query()
	setNameList, present := q["setName"]
	setNames := strings.Join(setNameList[:], ",")

	switch {
	case !present, len(setNames) == 0: // Case without query parameters
		err := dbConn.Preload("ProductLine").Model(tcm.SetInfo{}).Where("product_line_id = ?", productLineInfo.ID).Find(&setInfos).Error
		if err == nil {
			cardSetReps = make([]CardSetRep, len(setInfos), len(setInfos))
			for i, elem := range setInfos {
				cardSetReps[i].Set(elem)
			}
			jbuff, err = json.Marshal(cardSetReps)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(200)
		}
	case present: // Case with query parameters
		err := dbConn.Preload("ProductLine").Model(tcm.SetInfo{}).Where("name IN ?", setNameList).Scan(&setInfos).Error
		if err == nil {
			cardSetReps = make([]CardSetRep, len(setInfos), len(setInfos))
			for i, elem := range setInfos {
				cardSetReps[i].Set(elem)
			}
			jbuff, err = json.Marshal(cardSetReps)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(200)
		}
	}
	w.Write(jbuff)
}
