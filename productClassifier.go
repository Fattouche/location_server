package main

import (
	"bufio"
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jbrukh/bayesian"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Photography      bayesian.Class = "Photography"
	Audio            bayesian.Class = "Audio"
	Projector        bayesian.Class = "Projector"
	Drones           bayesian.Class = "Drones"
	DJ               bayesian.Class = "DJ"
	Transport        bayesian.Class = "Transport"
	Storage          bayesian.Class = "Storage"
	Electronics      bayesian.Class = "Electronics"
	Party            bayesian.Class = "Party"
	Sports           bayesian.Class = "Sports"
	Instruments      bayesian.Class = "Instruments"
	HomeOfficeGarden bayesian.Class = "HomeOfficeGarden"
	Kids             bayesian.Class = "Kids"
	Travel           bayesian.Class = "Travel"
	Clothing         bayesian.Class = "Clothing"
	textDirectory    string         = "./classification"
)

var classMap = map[string]bayesian.Class{"Photography": Photography, "Audio": Audio, "Projector": Projector, "Drones": Drones, "DJ": DJ, "Transport": Transport, "Storage": Storage, "Electronics": Electronics, "Party": Party, "Sports": Sports, "Instruments": Instruments, "HomeOfficeGarden": HomeOfficeGarden, "Kids": Kids, "Travel": Travel, "Clothing": Clothing}
var classes = make([]bayesian.Class, 0)
var classifier *bayesian.Classifier
var databaseMap map[bayesian.Class][]Object

//Train the model to understand different types of products
func trainModel() *bayesian.Classifier {
	//Essentially caching the result
	if classifier != nil {
		return classifier
	}
	for _, val := range classMap {
		classes = append(classes, val)
	}
	classifier = bayesian.NewClassifier(classes...)
	files, err := ioutil.ReadDir(textDirectory)
	if err != nil {
		log.Println(err)
		return nil
	}
	modelMap := make(map[string][]string)
	//Iterate through all files in directory, only works if they share the name with the bayesian classes.
	//Would be replaced with DB queries in a production system.
	for _, fileInfo := range files {
		fileName := fileInfo.Name()
		fullName := textDirectory + "/" + fileName
		modelMap[fileName] = make([]string, 0)
		file, err := os.Open(fullName)
		if err != nil {
			log.Println(err)
			return nil
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			modelMap[fileName] = append(modelMap[fileName], scanner.Text())
		}
	}
	for key := range modelMap {
		classifier.Learn(modelMap[key], classMap[key])
	}
	return classifier
}

//Probably not scalable to read the entire items table in, could do some sophisticated caching.
func classifyDatabase() (map[bayesian.Class][]Object, error) {
	if databaseMap != nil {
		return databaseMap, nil
	}
	databaseMap = make(map[bayesian.Class][]Object)
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//Probably okay to use a raw query here instead of prepared statement, no params.
	rows, err := db.Query("select item_name, lng, lat from items")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var itemName string
	var lat, lng float64
	//Iterate through each item, classify it and store it in the map.
	for rows.Next() {
		err = rows.Scan(&itemName, &lat, &lng)
		itemName = strings.ToLower(itemName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		predictedClass := classifyWord(itemName)
		obj := Object{itemName, lng, lat}
		databaseMap[predictedClass] = append(databaseMap[predictedClass], obj)
	}
	return databaseMap, nil
}

//Classifies a single word into a product category
func classifyWord(word string) bayesian.Class {
	classifier := trainModel()
	_, likely, _ := classifier.LogScores([]string{word})
	return classes[likely]
}

//Returns the top numItems from the evaluation as a list of strings (item names)
func topItems(numItems int, params Object) []string {
	likelyClass := classifyWord(params.name)
	arrObjects := databaseMap[likelyClass]
	sort.Slice(arrObjects, func(i, j int) bool {
		obj1 := arrObjects[i]
		obj2 := arrObjects[j]
		containsWordObj1 := false
		containsWordObj2 := false
		if strings.Contains(obj1.name, params.name) {
			containsWordObj1 = true
		}
		if strings.Contains(obj2.name, params.name) {
			containsWordObj2 = true
		}
		if containsWordObj1 && !containsWordObj2 {
			return true
		}
		if containsWordObj2 && !containsWordObj1 {
			return false
		}

		di := Distance(obj1.latitude, obj1.longtitude, params.latitude, params.longtitude)
		dj := Distance(obj2.latitude, obj2.longtitude, params.latitude, params.longtitude)
		return di < dj
	})
	topItems := make([]string, 0)
	for index, item := range arrObjects {
		topItems = append(topItems, item.name)
		if index > 20 {
			break
		}
	}
	return topItems
}
