package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
)

const (
	KrakenAPI = "https://api.kraken.com/0/public/"
	host      = "localhost"
	port      = "5432"
	user      = "toto"
	password  = "mysecretpassword"
	dbname    = "mydatabase"
)

// Format de la réponse Get Server Time
type ServerTime struct {
	Errors []string `json:"error"`
	Result struct {
		Unixtime int    `json:"unixtime"`
		Rfc      string `json:"rfc1123"`
	} `json:"result"`
}

// Format de la réponse Get System Status
type ServerStatus struct {
	Errors []string `json:"error"`
	Result struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	} `json:"result"`
}

func main() {

	// Connexion à la base de données
	connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	fmt.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connecté à postgres !")
	}

	// Requête du ServerTime
	request := QueryKraken("Time")
	// Parse the JSON string
	var sTime ServerTime
	errors := json.Unmarshal(request.([]byte), &sTime)
	if errors != nil {
		panic(errors)
	}

	// Requête du ServerStatus
	request = QueryKraken("SystemStatus")
	// Parse the JSON string
	var status ServerStatus
	errors = json.Unmarshal(request.([]byte), &status)
	if errors != nil {
		panic(errors)
	}

	// Requête des pairs
	request = QueryKraken("AssetPairs?pair=XXBTZUSD,XETHXXBT")
	var assetPairs map[string]interface{}
	errors = json.Unmarshal(request.([]byte), &assetPairs)
	if errors != nil {
		panic(errors)
	}
	resultMap := assetPairs["result"].(map[string]interface{})

	// Requête des ticker
	request = QueryKraken("Ticker?pair=XXBTZUSD,XETHXXBT")
	var ticker map[string]interface{}
	errors = json.Unmarshal(request.([]byte), &ticker)
	if errors != nil {
		panic(errors)
	}
	resultMapTicker := ticker["result"].(map[string]interface{})
	var tickers []interface{}
	for _, value := range resultMapTicker {
		pairMapTicker := value.(map[string]interface{})
		tickers = append(tickers, pairMapTicker["a"].([]interface{})...)
	}

	// Date au format D_M_YYYY_Hh_Mm
	currentDate := strconv.Itoa(time.Now().Day()) + "_" + strconv.Itoa(int(time.Now().Month())) + "_" + strconv.Itoa(time.Now().Year()) + "_" + strconv.Itoa(time.Now().Hour()) + "h" + strconv.Itoa(time.Now().Minute()) + "m"

	// Création du string à stocker dans le fichier
	fileContent := "Sauvegarde du " + currentDate + "\n"
	fileContent += "\nL'unixtime est de : " + strconv.Itoa(sTime.Result.Unixtime)
	fileContent += "\nLa date du serveur est : " + sTime.Result.Rfc
	fileContent += "\nLe status du serveur est : " + status.Result.Status
	fileContent += "\nLe timestamp est : " + status.Result.Timestamp + "\n"
	// Itération de la map result
	i := 0
	for _, value := range resultMap {
		pairMap := value.(map[string]interface{})
		fileContent += "\n" + pairMap["altname"].(string) + " :"
		fileContent += "\nBase : " + pairMap["base"].(string)
		fileContent += "\nQuote : " + pairMap["quote"].(string)
		fileContent += "\nCurrent price : " + tickers[i].(string)
		fileContent += "\nFee Volume Currency : " + pairMap["fee_volume_currency"].(string)
		fileContent += "\nOrder min : " + pairMap["ordermin"].(string)
		fileContent += "\nCost min : " + pairMap["costmin"].(string)
		fileContent += "\nStatus : " + pairMap["status"].(string) + "\n"
		i += 3
	}

	// Creation du fichier
	file, err := os.Create("Archive/" + currentDate + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Ecriture du fichier avec notre giga string
	_, err = file.WriteString(fileContent)
	if err != nil {
		panic(err)
	}

	// Sauvegarde du fichier
	err = file.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println("Fichier créé !")

	// Création de la table si elle n'existe pas
	sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs2 ( id SERIAL NOT NULL, altname character varying NOT NULL, base character varying, quote character varying, current_price character varying, fee_volume_currency character varying, ordermin character varying, costmin character varying, status character varying, PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to toto;"
	_, errors = db.Exec(sqlStat)
	if errors != nil {
		fmt.Println(errors)
		panic(err)
	} else {
		fmt.Println("Table OK")
	}

	// Enregistrement des données dans la base de données
	now := time.Now().Unix()
	i = 0
	for _, value := range resultMap {
		pairMap := value.(map[string]interface{})
		// Insertion des données dans la base de données
		sqlStatement := fmt.Sprintf("INSERT INTO pairs2 (id, altname, base, quote, current_price, fee_volume_currency, ordermin, costmin, status) VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", now, pairMap["altname"], pairMap["base"], pairMap["quote"], tickers[i], pairMap["fee_volume_currency"], pairMap["ordermin"], pairMap["costmin"], pairMap["status"])
		_, err := db.Exec(sqlStatement)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		// Eviter les collisions de clé primaire
		now++
		i += 3
	}

	startServer()

}
