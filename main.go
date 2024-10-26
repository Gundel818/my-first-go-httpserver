package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type NetworkCode struct {
	Name         string `json:"name"`
	PLMN         string `json:"plmn"`
	PLMNSHORT    string `json:"plmnshort"`
	UseShortPLMN bool   `json:"useShortPLMN"`
}

type Response struct {
	Message string `json:"message"`
	MCC     string `json:"mcc"`
	MNC     string `json:"mnc,omitempty"`
}

var networkTab []NetworkCode

func getPLMNSHORT(network NetworkCode) string {
	if network.UseShortPLMN {
		return network.PLMN
	}
	return network.PLMNSHORT
}

func searchByPLMN(plmn string) *NetworkCode {
	for _, network := range networkTab {
		if network.PLMN == plmn || getPLMNSHORT(network) == plmn {
			return &network
		}
	}
	return nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mcc := vars["mcc"]
	mnc := vars["mnc"]

	if mnc == "" {
		response := Response{Message: "MNC non fourni", MCC: mcc}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
		return
	}
	plmn := mcc + mnc
	search := searchByPLMN(plmn)
	if search != nil {
		response := Response{Message: search.Name, MCC: mcc, MNC: mnc}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}

	} else {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		response := Response{Message: "PLMN non trouvé", MCC: mcc, MNC: mnc}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	networkTab = append(networkTab, NetworkCode{
		Name:         "Orange France",
		PLMN:         "20801",
		PLMNSHORT:    "2081",
		UseShortPLMN: false,
	})
	networkTab = append(networkTab, NetworkCode{
		Name:         "SFR",
		PLMN:         "20810",
		PLMNSHORT:    "",
		UseShortPLMN: true,
	})
	networkTab = append(networkTab, NetworkCode{
		Name:         "Free",
		PLMN:         "20815",
		PLMNSHORT:    "",
		UseShortPLMN: true,
	})
	networkTab = append(networkTab, NetworkCode{
		Name:         "Bouygues Telecom",
		PLMN:         "20820",
		PLMNSHORT:    "",
		UseShortPLMN: true,
	})

	router := mux.NewRouter()

	router.HandleFunc("/api/{mcc}", apiHandler).Methods("GET")
	router.HandleFunc("/api/{mcc}/{mnc}", apiHandler).Methods("GET")
	router.Use(middleware)

	fmt.Println("Démarrage du serveur sur le port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}
