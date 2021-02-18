package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	tcm "github.com/gurbos/tcmodels"
	"gorm.io/gorm"
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

	var productLines []tcm.ProductLine
	tx := dbConn.Model(tcm.ProductLine{}).Find(&productLines) // Query list of product line info
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	// Populate list of product line representations with product line info
	plr := make([]ProductLineRep, len(productLines), len(productLines))
	for i, elem := range productLines {
		plr[i].Set(elem)
	}

	jbuff, err := json.Marshal(plr) // Encode list of product line representations to json
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jbuff)
}

func CardSetHandler(w http.ResponseWriter, r *http.Request) {
	var ds DataSourceName = GetDataSource()
	dbConn := DBConnection(ds.DSNString(), logger.Silent)
	vars := mux.Vars(r)
	// plName := vars["productLine"]
	var plInfo tcm.ProductLine
	var setInfoList []tcm.SetInfo

	err := dbConn.Model(tcm.ProductLine{}).Where("Name = ?", vars["productLine"]).First(&plInfo).Error // Query product line info
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(404)
			ebuff, _ := json.Marshal(APIError{ProductLineNotFoundErr})
			w.Write([]byte(ebuff))
		}
	} else {
		tx := dbConn.Preload("ProductLine").Model(tcm.SetInfo{}).Where("product_line_id = ?", plInfo.ID).Find(&setInfoList) // Query product line set info list
		if tx.Error != nil {
			log.Fatal(tx.Error)
		}

		cardSetRepList := MakeCardSetRepList(setInfoList, plInfo)
		jbuff, err := json.Marshal(cardSetRepList)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		w.Write(jbuff)
	}
}
