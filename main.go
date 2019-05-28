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
	ID   int `json:"id"`
	Date struct {
		Year  string `json:"year"`
		Month string `json:"month"`
		Day   string `json:"day"`
	} `json:"date"`
	Time    string `json:"time"`
	HisCall string `json:"his_call"`
	Mode    string `json:"mode"`
	RST     string `json:"rst"`
	Band    string `json:"band"`
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
		writeUrStationData(&pdf, setting)
		writeQSOData(&pdf, setting, qso)
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

func writeUrStationData(pdf *gopdf.GoPdf, setting QSLCard) {
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
	l, err := pdf.MeasureTextWidth(setting.StationData.Call)
	if err != nil {
		log.Fatal(err)
	}

	x := setting.CardSize.W/2 - l/2
	y := setting.CardSize.H / 4

	// write your callsign
	setLoc(pdf, x, y)
	pdf.Text(string(setting.StationData.Call))
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

func writeQSOData(pdf *gopdf.GoPdf, setting QSLCard, qso QSOsData) {
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

	year := qso.Date.Year
	day := qso.Date.Day
	month := qso.Date.Month
	time := qso.Time
	call := qso.HisCall
	mode := qso.Mode
	rst := qso.RST
	band := qso.Band

	// write qso date
	x := 40.0
	y := 180.0
	setLoc(pdf, x, y)
	pdf.Text(day)

	x += float64(fontSize * 4)
	setLoc(pdf, x, y)
	pdf.Text(month)

	x += float64(fontSize * 4)
	setLoc(pdf, x, y)
	pdf.Text(year)

	// write qso time
	x += float64(fontSize * 4)
	setLoc(pdf, x, y)
	pdf.Text(time)

	// write qso's RST report
	x += float64(fontSize * 5)
	setLoc(pdf, x, y)
	pdf.Text(rst)

	// write qso Mode
	x += float64(fontSize * 4)
	setLoc(pdf, x, y)
	pdf.Text(mode)

	// write qso band
	x += float64(fontSize * 5)
	setLoc(pdf, x, y)
	pdf.Text(band)

	// write his callsign
	var w float64
	if w, err = pdf.MeasureTextWidth("To "); err != nil {
		log.Fatal(err)
	}
	w1, err := pdf.MeasureTextWidth("00    01    2019    00:00    599    144    SSB")
	if err != nil {
		log.Fatal(err)
	}
	w1 *= 1.3
	y = 150.0
	setLoc(pdf, (setting.CardSize.W-w1)/2+w, y-float64(fontSize))
	pdf.Text(call)
}
