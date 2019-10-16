package novels

import (
	"crawlnovel/pkg/api/model"
	"crawlnovel/pkg/common/db"
	"errors"
	"fmt"
	"time"
)

const FileExt = "crawnovel"

func SaveBook(row map[int]interface{}) error {
	dto := model.Book{}
	dto.Name = fmt.Sprintf("%v", row[0])
	dto.Author = fmt.Sprintf("%v", row[1])
	dto.Image = fmt.Sprintf("%v", row[2])
	dto.Url = fmt.Sprintf("%v", row[3])
	dto.Desc = fmt.Sprintf("%v", row[4])
	dto.Status = 0
	dto.CreateTime = time.Now()
	dto.LastUpdateTime = time.Now()
	db := db.GetGormDB()
	created := db.Save(&dto)
	if created.RowsAffected > 0 {
		return nil
	}
	return errors.New("add book error")

}
