package main

import (
	"Food-Planner-Food-Management/model"
	"Food-Planner-Food-Management/provider"
	"Food-Planner-Food-Management/utils"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Printf("dfsdf")

	err := ImportDataFromFileExcel()
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("Import success")
	}
}

func GetDataExcel() ([]model.FoodSchemaModel , error){
	data := make([]model.FoodSchemaModel , 0)
	f, err := excelize.OpenFile("food_data.xlsx")

	if err != nil {
		return nil , err
	}

	rows, err := f.GetRows("Sheet1")
	for index, row := range rows {
		if index == 0 {
			continue
		}
		typeMaterial , errConv := strconv.Atoi(row[1])
		if errConv != nil {
			fmt.Println(errConv.Error())
			continue
		}
		materialLevel := model.MaterialLevelSchemaModel{
			Protein: utils.ConvertStringToInt(row[2]),
			Fiber:   utils.ConvertStringToInt(row[3]),
			Canxi:   utils.ConvertStringToInt(row[4]),
			Fat:     utils.ConvertStringToInt(row[5]),
			Starch:  utils.ConvertStringToInt(row[6]),
		}
		item := model.FoodSchemaModel{
			Name:          row[0],
			TypeFood:      typeMaterial,
			CreatedDate:   time.Now(),
			MaterialLevel: materialLevel,
			Image:         "",
			Mode:          "",
			Status: 1,
		}
		data = append(data , item)
	}

	return data , nil
}

func ImportDataFromFileExcel() error {
	dataExcel , errExcel := GetDataExcel()

	if errExcel != nil {
		return errExcel
	}

	db := provider.ConnectPostgres()
	if db == nil {
		return errors.New("...fail to connect postgres")
	}

	defer func() {
		errClose := db.Close()
		if errClose != nil {
			fmt.Println(errClose.Error())
		}
	}()

	sqlParams := make([]string , 0)
	sqlvalues := make([]interface{} , 0)

	for index , row := range dataExcel {
		valueLoop := index *5
		sqlParams = append(sqlParams , fmt.Sprintf(" ($%d, $%d, CURRENT_TIMESTAMP, $%d, $%d, $%d)",
			1 + valueLoop,2 + valueLoop,3 + valueLoop, 4 + valueLoop , 5 + valueLoop))
		sqlvalues = append(sqlvalues , row.Name , row.Status ,row.TypeFood , row.MaterialLevel, row.Image)
	}

	sqlQuery := fmt.Sprintf(`INSERT INTO public."food-info"(
	name, status, created_date, type_food, material_level,image) VALUES %s ; ` , strings.Join(sqlParams , ","))

	rowsQuery :=db.QueryRow(sqlQuery , sqlvalues...)
	if rowsQuery.Err() != nil {
		return rowsQuery.Err()
	}

	return nil
}