package controller

import (
	"fmt"

	baseresponse "github.com/dieg0code/rag-diary/base_response"
	"github.com/dieg0code/rag-diary/diary/dto"
	"github.com/dieg0code/rag-diary/diary/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DiaryControllerImpl struct {
	DiaryService service.DiaryService
}

// CreateDiary implements DiaryController.
func (d *DiaryControllerImpl) CreateDiary(c *gin.Context) {
	createDiaryReq := dto.CreateDiaryDTO{}

	err := c.ShouldBindJSON(&createDiaryReq)
	if err != nil {
		logrus.WithError(err).Error("cannot bind json")
		res := baseresponse.BaseResponse[string]{
			Code:   400,
			Status: "Bad Request",
			Msg:    fmt.Sprintf("cannot bind json: %v", err.Error()),
			Data:   "",
		}

		c.JSON(400, res)
		return
	}

	err = d.DiaryService.CreateDiary(createDiaryReq)
	if err != nil {
		logrus.WithError(err).Error("cannot create diary")
		res := baseresponse.BaseResponse[string]{
			Code:   500,
			Status: "Internal Server Error",
			Msg:    fmt.Sprintf("cannot create diary: %v", err.Error()),
			Data:   "",
		}

		c.JSON(500, res)
		return
	}

	res := baseresponse.BaseResponse[string]{
		Code:   201,
		Status: "Created Successfully",
		Msg:    "diary created successfully",
		Data:   "",
	}

	c.JSON(201, res)
}

// DeleteDiary implements DiaryController.
// func (d *DiaryControllerImpl) DeleteDiary(c *gin.Context) {
// 	diaryID := c.Param("id")

// 	id, err := strconv.Atoi(diaryID)
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot convert id to int")
// 		res := baseresponse.BaseResponse[string]{
// 			Code:   400,
// 			Status: "Bad Request",
// 			Msg:    fmt.Sprintf("cannot convert id to int: %v", err.Error()),
// 			Data:   "",
// 		}

// 		c.JSON(400, res)
// 		return
// 	}

// 	err = d.DiaryService.DeleteDiary(id)
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot delete diary")
// 		res := baseresponse.BaseResponse[string]{
// 			Code:   500,
// 			Status: "Internal Server Error",
// 			Msg:    fmt.Sprintf("cannot delete diary: %v", err.Error()),
// 			Data:   "",
// 		}

// 		c.JSON(500, res)
// 		return
// 	}

// 	res := baseresponse.BaseResponse[string]{
// 		Code:   200,
// 		Status: "OK",
// 		Msg:    "diary deleted successfully",
// 		Data:   "",
// 	}

// 	c.JSON(200, res)
// }

// GetAllDiaries implements DiaryController.
// func (d *DiaryControllerImpl) GetAllDiaries(c *gin.Context) {

// 	diaries, err := d.DiaryService.GetAllDiaries()
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot get all diaries")
// 		res := baseresponse.BaseResponse[string]{
// 			Code:   500,
// 			Status: "Internal Server Error",
// 			Msg:    fmt.Sprintf("cannot get all diaries: %v", err.Error()),
// 			Data:   "",
// 		}

// 		c.JSON(500, res)
// 		return
// 	}

// 	res := baseresponse.BaseResponse[[]*dto.DiaryDTO]{
// 		Code:   200,
// 		Status: "OK",
// 		Msg:    "diaries fetched successfully",
// 		Data:   diaries,
// 	}

// 	c.JSON(200, res)
// }

// GetDiary implements DiaryController.
// func (d *DiaryControllerImpl) GetDiary(c *gin.Context) {

// 	diaryID := c.Param("id")

// 	id, err := strconv.Atoi(diaryID)
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot convert id to int")
// 		res := baseresponse.BaseResponse[string]{
// 			Code:   400,
// 			Status: "Bad Request",
// 			Msg:    fmt.Sprintf("cannot convert id to int: %v", err.Error()),
// 			Data:   "",
// 		}

// 		c.JSON(400, res)
// 		return
// 	}

// 	diary, err := d.DiaryService.GetDiary(id)
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot get diary")
// 		res := baseresponse.BaseResponse[string]{
// 			Code:   500,
// 			Status: "Internal Server Error",
// 			Msg:    fmt.Sprintf("cannot get diary: %v", err.Error()),
// 			Data:   "",
// 		}

// 		c.JSON(500, res)
// 		return
// 	}

// 	res := baseresponse.BaseResponse[*dto.DiaryDTO]{
// 		Code:   200,
// 		Status: "OK",
// 		Msg:    "diary fetched successfully",
// 		Data:   diary,
// 	}

// 	c.JSON(200, res)
// }

// RAGResponse implements DiaryController.
func (d *DiaryControllerImpl) RAGResponse(c *gin.Context) {

	clientIp := c.ClientIP()

	query := dto.SemanticQueryWithHistoryDTO{}

	err := c.ShouldBindJSON(&query)
	if err != nil {
		logrus.WithError(err).Error("cannot bind json")
		res := baseresponse.BaseResponse[string]{
			Code:   400,
			Status: "Bad Request",
			Msg:    fmt.Sprintf("cannot bind json: %v", err.Error()),
			Data:   "",
		}

		c.JSON(400, res)
		return
	}

	err = d.DiaryService.SaveUserMessage(query.Query, clientIp)
	if err != nil {
		logrus.WithError(err).Error("error saving user message")
	}

	response, err := d.DiaryService.RAGResponse(query)
	if err != nil {
		logrus.WithError(err).Error("cannot get rag response")
		res := baseresponse.BaseResponse[string]{
			Code:   500,
			Status: "Internal Server Error",
			Msg:    fmt.Sprintf("cannot get rag response: %v", err.Error()),
			Data:   "",
		}

		c.JSON(500, res)
		return
	}

	res := baseresponse.BaseResponse[string]{
		Code:   200,
		Status: "OK",
		Msg:    "rag response fetched successfully",
		Data:   response,
	}

	c.JSON(200, res)
}

// SemanticSearch implements DiaryController.
func (d *DiaryControllerImpl) SemanticSearch(c *gin.Context) {

	query := dto.SemanticQueryDTO{}

	err := c.ShouldBindJSON(&query)
	if err != nil {
		logrus.WithError(err).Error("cannot bind json")
		res := baseresponse.BaseResponse[string]{
			Code:   400,
			Status: "Bad Request",
			Msg:    fmt.Sprintf("cannot bind json: %v", err.Error()),
			Data:   "",
		}

		c.JSON(400, res)
		return
	}

	response, err := d.DiaryService.SematicSearch(query.Query)
	if err != nil {
		logrus.WithError(err).Error("cannot get semantic search")
		res := baseresponse.BaseResponse[string]{
			Code:   500,
			Status: "Internal Server Error",
			Msg:    fmt.Sprintf("cannot get semantic search: %v", err.Error()),
			Data:   "",
		}

		c.JSON(500, res)
		return
	}

	res := baseresponse.BaseResponse[string]{
		Code:   200,
		Status: "OK",
		Msg:    "semantic search fetched successfully",
		Data:   response,
	}

	c.JSON(200, res)
}

func NewDiaryControllerImpl(diaryService service.DiaryService) DiaryController {
	return &DiaryControllerImpl{
		DiaryService: diaryService,
	}
}
