package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"trader-web-api/dbcontroller"
	utilites "trader-web-api/utilities"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func fromTitleToUrl(title string) string {
	return strings.ReplaceAll(title, " ", "-")
}

var saveCounter uint8

func handleBase64Image(content string) string {
	condition := regexp.MustCompile(`data:image\/([a-zA-Z]*);base64,([^\"]*)`)
	result := condition.FindAllStringSubmatch(content, -1)
	for i := range result {
		var fileName = strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.Itoa(int(saveCounter)) + "." + result[i][1]
		saveCounter++
		var fileNameWithKey = "SERVER_PUBLIC_KEY" + fileName
		// replace data to url to save to database
		content = strings.Replace(content, result[i][0], fileNameWithKey, 1)
		writeToFile("./public/images/"+fileName, result[i][2])
	}
	return content
}

func handleHashTag(reqTags []string) {
	tags, err := dbcontroller.GetHashTag()
	if err == mongo.ErrNoDocuments {
		dbcontroller.AddHashTag(reqTags)
	} else {
		var isTagChange = false
		// Check if need to add new hashtag to DB
		for _, value := range reqTags {
			if !utilites.Contains(tags.Tags, value) {
				isTagChange = true
				tags.Tags = append(tags.Tags, value)
			}
		}
		if isTagChange {
			dbcontroller.UpdateHashTag(tags)
		}
	}
}

func validateUrl(title string) string {
	var url = fromTitleToUrl(title)
	_, err := dbcontroller.FindTutByUrl(url)
	if err == mongo.ErrNoDocuments {
		return url
	} else {
		return ""
	}
}

type TutorialBody struct {
	Id      string   `json:"Id"`
	Title   string   `json:"Title"`
	Content string   `json:"Content"`
	Tags    []string `json:"Tag"`
}

/* Upload new post */
func uploadNewPost(w http.ResponseWriter, r *http.Request) {
	var tutBody TutorialBody
	json.NewDecoder(r.Body).Decode(&tutBody)

	var newUrl = validateUrl(tutBody.Title)
	if newUrl != "" {
		content := handleBase64Image(tutBody.Content)
		dbcontroller.CreateNewPost(tutBody.Title, content, tutBody.Tags, newUrl)
		handleHashTag(tutBody.Tags)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getAllTutorial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, _ := dbcontroller.FindAllTutorial()
	json.NewEncoder(w).Encode(result)
}
func getDetailTutorial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := mux.Vars(r)["url"]
	result, err := dbcontroller.FindTutByUrl(url)
	if err == nil {
		result.Content = strings.ReplaceAll(result.Content, "SERVER_PUBLIC_KEY", "http://localhost:5000/images/")
		json.NewEncoder(w).Encode(result)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getAllHashTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tags, _ := dbcontroller.GetHashTag()
	json.NewEncoder(w).Encode(tags)
}

func updateTutorial(w http.ResponseWriter, r *http.Request) {
	var tutBody TutorialBody
	json.NewDecoder(r.Body).Decode(&tutBody)

	id, _ := primitive.ObjectIDFromHex(tutBody.Id)
	tut := dbcontroller.Tutorial{
		Id:    id,
		Title: tutBody.Title,
		Tag:   tutBody.Tags,
	}
	tut.Content = handleBase64Image(tutBody.Content)
	handleHashTag(tutBody.Tags)

	err := dbcontroller.UpdateTutorial(tut)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func deleteTutorial(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		panic(err1)
	}
	err := dbcontroller.DeleteTutorial(objID)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
