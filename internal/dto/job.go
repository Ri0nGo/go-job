package dto

import "go-job/internal/models"

// ----------- node ----------- //

type ReqJob struct {
	models.Job
	Filename string `json:"filename"`
}
