package services

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"trader-web-api/dtos"
	"trader-web-api/models"
	"trader-web-api/repositories"
	"trader-web-api/utils"
)

type TutorialService interface {
	//GetUserInfo(userID uint) *dtos.GetUserInfoResponse
	GetAllPost() []dtos.TutorialJson
	CreateNewPost(tut dtos.TutorialJson) bool
	GetDetailTutorial(url string) dtos.TutorialResponse
	UpdateTutorial(tut dtos.TutorialJson) dtos.Meta
	DeleteTutorial(id uint) dtos.Meta
	GetAllHashTag() []string
}

type tutorialServiceImpl struct {
	tutRepo     repositories.TutorialRepository
	hashtagRepo repositories.HashTagRepository
}

func newTutorialService(tutRepo repositories.TutorialRepository, hashtagRepo repositories.HashTagRepository) TutorialService {
	return &tutorialServiceImpl{
		tutRepo:     tutRepo,
		hashtagRepo: hashtagRepo,
	}
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
		utils.WriteToFile("./public/images/"+fileName, result[i][2])
	}
	return content
}
func fromTitleToUrl(title string) string {
	return strings.ReplaceAll(title, " ", "-")
}
func (r *tutorialServiceImpl) CreateNewPost(tutJson dtos.TutorialJson) bool {
	// validate url
	var newUrl = fromTitleToUrl(tutJson.Title)
	_, err := r.tutRepo.FindTutByUrl(newUrl)
	if err != nil {
		newUrl = ""
	}

	if newUrl != "" {
		tutJson.Content = handleBase64Image(tutJson.Content)
		var temp = models.Tutorial{
			ID:          tutJson.Id,
			Title:       tutJson.Title,
			Content:     tutJson.Content,
			Url:         newUrl,
			Tags:        strings.Join(tutJson.Tags, repositories.SPLITKEY),
			AuthorID:    tutJson.AuthorID,
			Description: tutJson.Description,
		}
		err = r.tutRepo.CreateNewPost(temp)
		if err != nil {
			return false
		}
		// handle hashtag
		tags, err := r.hashtagRepo.Find()
		if err != nil {
			r.hashtagRepo.Create(tutJson.Tags)
			log.Println("Can not get HashTag")
		} else {
			var isTagChange = false
			// Check if need to add new hashtag to DB
			for _, value := range tutJson.Tags {
				if !utils.Contains(tags, value) {
					isTagChange = true
					tags = append(tags, value)
				}
			}
			if isTagChange {
				log.Println("update hashtag in DB: ", tags)
				r.hashtagRepo.Update(strings.Join(tags, repositories.SPLITKEY))
			}
		}
		return true
	} else {
		return false
	}
}

func (r *tutorialServiceImpl) GetAllHashTag() []string {
	result, err := r.hashtagRepo.Find()
	if err != nil {
		return nil
	} else {
		return result
	}
}
func (r *tutorialServiceImpl) GetAllPost() []dtos.TutorialJson {
	tutArray, err := r.tutRepo.FindAllTutorial()
	if err != nil {
		// write error
	}
	var result []dtos.TutorialJson
	for _, value := range tutArray {
		var temp = dtos.TutorialJson{
			Id:          value.ID,
			Title:       value.Title,
			Url:         value.Url,
			Content:     value.Content,
			AuthorID:    value.AuthorID,
			Description: value.Description,
		}
		result = append(result, temp)
	}
	return result
}

func (r *tutorialServiceImpl) GetDetailTutorial(url string) dtos.TutorialResponse {
	result, err := r.tutRepo.FindTutByUrl(url)
	var res = dtos.TutorialResponse{}
	// ID ==0 means there was not row effect
	if err != nil || result.ID == 0 {
		res.Data = dtos.TutorialJson{}
		res.Meta = *dtos.InternalServerErrorMeta
		return res
	}

	var jsonResult = dtos.TutorialJson{
		Id:          result.ID,
		Title:       result.Title,
		Content:     result.Content,
		Tags:        strings.Split(result.Tags, repositories.SPLITKEY),
		Url:         result.Url,
		AuthorID:    result.AuthorID,
		Description: result.Description,
	}

	res.Data = jsonResult
	res.Meta = *dtos.SuccessMeta
	return res
}

func (r *tutorialServiceImpl) UpdateTutorial(tut dtos.TutorialJson) dtos.Meta {
	var meta dtos.Meta
	err := r.tutRepo.UpdateTutorial(models.Tutorial{
		ID:          tut.Id,
		Title:       tut.Title,
		Content:     tut.Content,
		Url:         tut.Url,
		Tags:        strings.Join(tut.Tags, repositories.SPLITKEY),
		AuthorID:    tut.AuthorID,
		Description: tut.Description,
	})
	if err == nil {
		meta.Code = http.StatusOK
	} else {
		meta.Code = http.StatusInternalServerError
		meta.Message = "Cannot update Tutorial. <br> Error:" + err.Error()
	}
	return meta
}

func (r *tutorialServiceImpl) DeleteTutorial(id uint) dtos.Meta {
	var meta dtos.Meta
	err := r.tutRepo.DeleteTutorial(id)
	if err == nil {
		meta.Code = http.StatusOK
	} else {
		meta.Code = http.StatusInternalServerError
		meta.Message = "Cannot delete Tutorial. <br> Error:" + err.Error()
	}
	return meta
}
