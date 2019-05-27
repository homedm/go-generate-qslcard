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
	err = pdf.AddTTFFont("migu-1m", "../../../../../.fonts/migu-1m-regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	err = pdf.SetFont("migu-1m", "", 14)
	if err != nil {
		log.Fatal(err)
	}

	drawNewPage(&pdf)
	writeHisCallsign(&pdf, "JJ1HGP")

	pdf.WritePdf("qslcards.pdf")
	//pdf.Line()
}

// draw box
func drawNewPage(pdf *gopdf.GoPdf) {
	pdf.AddPage()

	pdf.SetLineWidth(2)

	// draw his callsign box
	x := 240.0
	y := 50.0
	h1 := 25.0
	w1 := 15.0
	for i := 0.0; i < 6.0; i += 1.0 {
		pdf.RectFromUpperLeftWithStyle(x+i*20, y, w1, h1, "single")
	}

	// draw QSO data form
	x = 50.0
	y = 100.0
	w1 = 300.0
	h1 = 60.0
	pdf.RectFromUpperLeftWithStyle(x, y, w1, h1, "single")
	h2 := 20.0
	pdf.SetLineWidth(1)
	pdf.Line(x, y+h2, x+w1, y+h2)

}

func writeHisCallsign(pdf *gopdf.GoPdf, call string) {
	x := 245.0
	y := 70.0
	for i, c := range call {
		pdf.SetX(x + float64(i*20))
		pdf.SetY(y)
		pdf.Text(string(c))
	}
}
