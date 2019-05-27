package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/signintech/gopdf"
)

// QSOsData is QSO logs data
type QSOsData struct {
	ID      int    `josn:"id"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	HisCall string `json:"his_callsign"`
	Mode    string `json:"mode"`
}

func main() {
	bytes, err := ioutil.ReadFile("qsos.json")
	if err != nil {
		log.Fatal(err)
	}
	// json decode
	var qsos []QSOsData
	if err := json.Unmarshal(bytes, &qsos); err != nil {
		log.Fatal(err)
	}

	// print debug data
	for _, qso := range qsos {
		fmt.Printf("%d: Date=%s %s, HisCall=%s, Mode=%s\n", qso.ID, qso.Date, qso.Time, qso.HisCall, qso.Mode)
	}

	pdf := gopdf.GoPdf{}
	// 10.0 cm x 14.8 cm = 283.5 pt x 419.5 pt
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{H: 283.5, W: 419.5}})

	drawNewPage(&pdf)

	pdf.WritePdf("qslcards.pdf")
	//pdf.Line()
}

// draw box
func drawNewPage(pdf *gopdf.GoPdf) {
	pdf.AddPage()

	pdf.SetLineWidth(2)

	// draw his callsign box
	x := 230.0
	y := 50.0
	h := 25.0
	w := 15.0
	for i := 0.0; i < 6.0; i += 1.0 {
		pdf.RectFromUpperLeftWithStyle(x+i*20, y, w, h, "single")
	}

	// draw QSO data form
	x = 50.0
	y = 100.0
	w = 100
	pdf.RectFromUpperLeftWithStyle(x, y, w, h, "single")
	h2 := 40.0
	for i := 0.0; i < 3; i += 1.0 {
		pdf.RectFromUpperLeftWithStyle(x+i*w/3, y+h, w/3, h2, "single")
	}
}
