package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xianggong/m2svis/server/modules/database"
)

func cycleCountByInstType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	data := database.GetCycleCountByInstType(traceName)

	// Encode to JSON and write
	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func cycleCountByExecUnit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	data := database.GetCycleCountByExecUnit(traceName)

	// Encode to JSON and write
	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func cycleCountByCU(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	data := database.GetCycleCountByCU(traceName)

	// Encode to JSON and write
	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}
