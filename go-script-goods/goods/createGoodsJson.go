package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Excel struct {
	Filename string
	colMap   map[string]int
	file     *excelize.File
	rowsNum  int
	colsNum  int
	rows     [][]string
}

type LevelOne map[string]LevelTwo

type LevelTwo map[string][]Goods

type Goods map[string]GoodsInfoStruct

type GoodsInfoStruct struct {
	ImgUrl    string `json:"img_url"`
	Brand     string `json:"brand"`
	GoodsName string `json:"goods_name"`
	Layer     string `json:"layer"`
	Similar   string `json:"similar"`
}

func OpenExcel(filename, sheet string) Excel {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
		return Excel{}
	}
	e := Excel{}
	e.Filename = filename
	e.file = xlsx
	e.initColMap(sheet)
	return e
}

func (e *Excel) initColMap(sheet string) {
	rows, err := e.file.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
		return
	}
	e.rowsNum = len(rows)
	e.rows = rows
	if len(e.rows) <= 0 {
		log.Fatal(errors.New("Excel has no data! "))
		return
	}
	e.colMap = make(map[string]int)
	for index, colCell := range e.rows[0] {
		if colCell == "" {
			continue
		}
		e.colsNum += 1
		e.colMap[colCell] = index
	}
}

func (e *Excel) getCol(colName string) int {
	value, ok := e.colMap[colName]
	if ok {
		return value
	}
	log.Fatal("Col name no exist! ")
	return -1
}

func (e *Excel) CreateJson() {
	mainJson := LevelOne{}
	for i, row := range e.rows {
		if i == 0 {
			continue
		}
		levelOne := row[e.getCol("一级分类")]
		if _, ok := mainJson[levelOne]; !ok {
			mainJson[levelOne] = LevelTwo{}
		}
		levelTwo := row[e.getCol("二级分类")]
		if _, ok := mainJson[levelOne][levelTwo]; !ok {
			mainJson[levelOne][levelTwo] = []Goods{}
		}
		barcode := row[e.getCol("条形码")]
		brand := row[e.getCol("品牌")]
		goodsName := row[e.getCol("商品名字")]
		layer := row[e.getCol("支持层数")]
		similar := row[e.getCol("近似商品组")]
		goods := Goods{}
		goods[barcode] = GoodsInfoStruct{
			ImgUrl:    barcode + ".png",
			Brand:     brand,
			GoodsName: goodsName,
			Layer:     layer,
			Similar:   similar,
		}
		mainJson[levelOne][levelTwo] = append(mainJson[levelOne][levelTwo], goods)

		if _, ok := mainJson["相似类"]; !ok {
			mainJson["相似类"] = LevelTwo{}
		}
		if similar != "" {
			levelSimilar := similar
			if _, ok := mainJson["相似类"][levelSimilar]; !ok {
				mainJson["相似类"][levelSimilar] = []Goods{}
			}
			mainJson["相似类"][levelSimilar] = append(mainJson["相似类"][levelSimilar], goods)
		}
	}

	jsonStr, err := json.Marshal(mainJson)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonStr))
}

func main() {
	filename := flag.String("filename", "", "Excel file path")
	sheet := flag.String("sheet", "", "Sheet name")
	flag.Parse()
	e := OpenExcel(*filename, *sheet)
	e.CreateJson()
}
