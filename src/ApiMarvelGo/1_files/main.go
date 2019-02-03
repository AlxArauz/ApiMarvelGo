package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	PrivateKey = "752215bd2cbff133741a8da5df1be8233983344a"
	PublicKey  = "49b955850cf4d2d3856a7cb5d08bc753"
)

//ESTRUCTURA DE DATOS
type JsonStructResponse struct {
	/*Code            int    `json:"code"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
	Etag            string `json:"etag"`*/
	Data struct {
		/*Offset  int `json:"offset"`
		Limit   int `json:"limit"`
		Total   int `json:"total"`
		*/
		Count   int `json:"count"`
		Results []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Modified    string `json:"modified"`
			Comics      struct {
				Available int `json:"available"`
				Items     []struct {
					Name string `json:"name"`
				} `json:"items"`
				/*Returned int `json:"returned"`*/
			} `json:"comics"`
			Series struct {
				Available int `json:"available"`
				Items     []struct {
					Name string `json:"name"`
				} `json:"items"`
				/*Returned int `json:"returned"`*/
			} `json:"series"`
			Stories struct {
				Available int `json:"available"`
				Items     []struct {
					Name string `json:"name"`
					Type string `json:"type"`
				} `json:"items"`
				/*Returned int `json:"returned"`*/
			} `json:"stories"`
			Events struct {
				Available int `json:"available"`
				Items     []struct {
					Name string `json:"name"`
				} `json:"items"`
				/*Returned int `json:"returned"`*/
			} `json:"events"`
		} `json:"results"`
	} `json:"data"`
}

func main() {

	time := time.Now().String()
	ts := fmt.Sprintf("%x", md5.Sum([]byte(time)))
	hash := fmt.Sprintf("%x", md5.Sum([]byte(ts+PrivateKey+PublicKey)))

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Get Api Marvel Go EdTeam \n")

	for {
		//OPCIONES DEL MENU
		fmt.Println("---SELECCIONE---")
		fmt.Println("0) Exit")
		fmt.Println("1) Search by Name")
		fmt.Println("2) List  \n")
		fmt.Print("-> ")
		Opcion, _ := reader.ReadString('\n')
		// convert CRLF to LF
		Opcion = strings.Replace(Opcion, "\r\n", "", -1)

		if strings.Compare("0", Opcion) == 0 {
			os.Exit(0)
		} else if strings.Compare("1", Opcion) == 0 {
			fmt.Println("Write Name")
			fmt.Print("-> ")
			Name, _ := reader.ReadString('\n')
			Name = strings.Replace(Name, "\r\n", "", -1)
			if len(strings.TrimSpace(Name)) == 0 {
				fmt.Println("\n No ha introducido ningún valor para la búsqueda \n")
			} else {
				SearchHeroByName(Name, ts, hash)
			}

		} else if strings.Compare("2", Opcion) == 0 {
			GetHeroList(ts, hash)
		}

	}

}

//BUSQUEDA POR NOMBRE
func SearchHeroByName(Name string, ts string, hash string) {
	hero, _ := UrlEncoded(Name)
	response, err := http.Get("http://gateway.marvel.com/v1/public/characters?name=" + hero + "&ts=" + ts + "&apikey=" + PublicKey + "&hash=" + hash)

	if err != nil {
		fmt.Println(string(err.Error()))
	} else {
		ResponseJson, _ := ioutil.ReadAll(response.Body)

		textBytes := []byte(ResponseJson)
		var JSP JsonStructResponse
		json.Unmarshal(textBytes, &JSP)

		if len(JSP.Data.Results) == 0 {
			fmt.Println("\n Record not found \n")
		} else {
			DrawRecords(JSP, 1)
		}
	}

}

//OBTENER LISTA DE HEROES
func GetHeroList(ts string, hash string) {

	response, err := http.Get("http://gateway.marvel.com/v1/public/characters?ts=" + ts + "&apikey=" + PublicKey + "&hash=" + hash)

	if err != nil {
		fmt.Println(string(err.Error()))
	} else {
		ResponseJson, _ := ioutil.ReadAll(response.Body)

		textBytes := []byte(ResponseJson)

		var JSP JsonStructResponse
		json.Unmarshal(textBytes, &JSP)

		if len(JSP.Data.Results) == 0 {
			fmt.Println("\n Record not found \n")
		} else {
			DrawRecords(JSP, 2)
		}
	}

}

//IMPRIMIR REGISTROS
func DrawRecords(JSP JsonStructResponse, Op int) {
	sort.Slice(JSP.Data.Results, func(i, j int) bool {
		return strings.ToLower(JSP.Data.Results[i].Name) < strings.ToLower(JSP.Data.Results[j].Name)
	})

	for i := 0; i < len(JSP.Data.Results); i++ {
		if i == 0 {
			if Op == 1 {
				fmt.Println("\n Data of the Hero")
			} else {
				fmt.Println(string("\n List of Heros sort by Name"))
			}
			fmt.Println("Records : " + strconv.Itoa(JSP.Data.Count))
		}
		fmt.Println("++++++++++++++")
		fmt.Println("Nro: " + strconv.Itoa(i+1))
		fmt.Println("Identifier: " + strconv.Itoa(JSP.Data.Results[i].ID))
		fmt.Println("Name: " + JSP.Data.Results[i].Name)
		fmt.Println("Description: " + JSP.Data.Results[i].Description)
		fmt.Println("Modified: " + JSP.Data.Results[i].Modified)
		fmt.Println("Comics: " + strconv.Itoa(JSP.Data.Results[i].Comics.Available))
		for c := 0; c < len(JSP.Data.Results[i].Comics.Items); c++ {
			fmt.Println("Comic Name: " + JSP.Data.Results[i].Comics.Items[c].Name)
			if c == 2 {
				break
			}
		}
		fmt.Println("Series: " + strconv.Itoa(JSP.Data.Results[i].Series.Available))
		for s := 0; s < len(JSP.Data.Results[i].Series.Items); s++ {
			fmt.Println("Serie Name: " + JSP.Data.Results[i].Series.Items[s].Name)
			if s == 2 {
				break
			}
		}
		fmt.Println("Stories: " + strconv.Itoa(JSP.Data.Results[i].Stories.Available))
		for t := 0; t < len(JSP.Data.Results[i].Stories.Items); t++ {
			fmt.Println("Stories Name: " + JSP.Data.Results[i].Stories.Items[t].Name)
			if t == 2 {
				break
			}
		}
		fmt.Println("Events: " + strconv.Itoa(JSP.Data.Results[i].Events.Available))
		for e := 0; e < len(JSP.Data.Results[i].Events.Items); e++ {
			fmt.Println("Events Name: " + JSP.Data.Results[i].Events.Items[e].Name)
			if e == 2 {
				break
			}
		}

		fmt.Println("================= \n")
	}
}

//encodeuricomponent
func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
