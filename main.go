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
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	HisCall string `json:"his_callsign"`
	Mode    string `json:"mode"`
}

// QSLCard config
type QSLCard struct {
	CardSize struct {
		H float64 `json:"H"`
		W float64 `json:"W"`
	} `json:"CardSize"`
	StationData struct {
		Call string `json:"Call"`
		QRA  string `json:"QRA"`
	} `json:"StationData"`
	UrCallSign struct {
		Size     int    `json:"Size"`
		Name     string `json:"Name"`
		Location string `json:"Location"`
	} `json:"UrCallSign"`
	Body struct {
		Size     int    `json:"Size"`
		Name     string `json:"Name"`
		Location string `json:"Location"`
	} `json:"Body"`
}

func main() {
	qsos, err := readLogs("qsos.json")
	if err != nil {
		log.Fatal(err)
	}

	setting, err := readConfig("qslcard.json")
	if err != nil {
		log.Fatal(err)
	}

	// init pdf
	pdf := gopdf.GoPdf{}
	// 10.0 cm x 14.8 cm = 283.5 pt x 419.5 pt
	config := gopdf.Config{PageSize: gopdf.Rect{H: setting.CardSize.H, W: setting.CardSize.W}}

	pdf.Start(config)

	// print debug data
	for _, qso := range qsos {
		fmt.Printf("%d: Date=%s %s, HisCall=%s, Mode=%s\n", qso.ID, qso.Date, qso.Time, qso.HisCall, qso.Mode)

		drawNewPage(&pdf, setting)
		writeHisCallsign(&pdf, setting, qso.HisCall)
	}

	// output
	pdf.WritePdf("qslcards.pdf")
}

// draw box
func drawNewPage(pdf *gopdf.GoPdf, setting QSLCard) {
	pdf.AddPage()

	// set font
	fontName := setting.Body.Name
	fontLocation := setting.Body.Location
	fontSize := setting.Body.Size

	err := pdf.AddTTFFont(fontName, fontLocation)
	if err != nil {
		log.Fatal(err)
	}

	err = pdf.SetFont(fontName, "", fontSize)
	if err != nil {
		log.Fatal(err)
	}

	// draw QSO data form
	y := 150.0
	w1, err := pdf.MeasureTextWidth("00    01    2019    00:00    599    144    SSB")
	w1 *= 1.3
	x := (setting.CardSize.W - w1) / 2
	if err != nil {
		log.Fatal(err)
	}
	h1 := 60.0
	setLoc(pdf, x, y-float64(fontSize))
	pdf.Text("To")
	pdf.SetLineWidth(2)
	pdf.RectFromUpperLeftWithStyle(x, y, w1, h1, "single")
	pdf.Line(x, y+float64(fontSize)*1.2, x+w1, y+float64(fontSize)*1.2)
	pdf.SetLineWidth(1)
	setLoc(pdf, x, y+float64(fontSize))
	pdf.Text(string("  DAY    MONTH    YEAR     TIME      RST      BAND    MODE"))

	w2, err := pdf.MeasureTextWidth("  DAY  ")
	if err != nil {
		log.Fatal(err)
	}
	x += w2
	pdf.Line(x, y, x, y+h1)

	w2, err = pdf.MeasureTextWidth("  MONTH  ")
	if err != nil {
		log.Fatal(err)
	}
	x += w2
	pdf.Line(x, y, x, y+h1)

	w2, err = pdf.MeasureTextWidth("  YEAR  ")
	if err != nil {
		log.Fatal(err)
	}
	x += w2
	pdf.Line(x, y, x, y+h1)

	w2, err = pdf.MeasureTextWidth("   TIME   ")
	if err != nil {
		log.Fatal(err)
	}
	x += w2
	pdf.Line(x, y, x, y+h1)

	w2, err = pdf.MeasureTextWidth("   RST   ")
	if err != nil {
		log.Fatal(err)
	}
	x += w2
	pdf.Line(x, y, x, y+h1)

	w2, err = pdf.MeasureTextWidth("   BAND   ")
	if err != nil {
		log.Fatal(err)
	}
	x += w2
	pdf.Line(x, y, x, y+h1)

	// draw other QSO Data form
	// Rig
	x = (setting.CardSize.W - w1) / 2
	y += h1 + float64(fontSize)*1.2
	pdf.Line(x, y, x+w1/2-10, y)
	setLoc(pdf, x, y-2)
	pdf.Text("Rig")
	// ANT
	pdf.Line(x+w1/2+10, y, x+w1, y)
	setLoc(pdf, x+w1/2+10, y-2)
	pdf.Text("ANT")
	// RMKS
	y += float64(fontSize)
	pdf.Line(x, y, x+w1, y)
	setLoc(pdf, x, y-2)
	pdf.Text("Remarks")
}

func writeHisCallsign(pdf *gopdf.GoPdf, setting QSLCard, call string) {
	// set font
	fontName := setting.UrCallSign.Name
	fontLocation := setting.UrCallSign.Location
	fontSize := setting.UrCallSign.Size

	err := pdf.AddTTFFont(fontName, fontLocation)
	if err != nil {
		log.Fatal(err)
	}

	err = pdf.SetFont(fontName, "", fontSize)
	if err != nil {
		log.Fatal(err)
	}

	// calc your call sign position
	l, err := pdf.MeasureTextWidth(call)
	if err != nil {
		log.Fatal(err)
	}

	x := setting.CardSize.W/2 - l/2
	y := setting.CardSize.H / 4

	setLoc(pdf, x, y)
	pdf.Text(string(call))
}

func readLogs(file string) ([]QSOsData, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	// json decode
	var qsos []QSOsData
	if err := json.Unmarshal(bytes, &qsos); err != nil {
		log.Fatal(err)
	}
	return qsos, nil
}

func readConfig(file string) (QSLCard, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	// json decode
	var setting QSLCard
	if err := json.Unmarshal(bytes, &setting); err != nil {
		log.Fatal(err)
	}
	return setting, nil
}

func setLoc(pdf *gopdf.GoPdf, x float64, y float64) {
	pdf.SetX(x)
	pdf.SetY(y)
}
