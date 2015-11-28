package wupws

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var baseURL = "http://weatherstation.wunderground.com/weatherstation/updateweatherstation.php?"
var dateUTC time.Time

const timeFormat = "2006-01-02 15:04:05"

//Station stores station ID and password.
type Station struct {
	id                string
	password          string
	softwaretype      string
	calculateDewpoint bool
	weatherParameters map[string]string

	lastUpdateWeather time.Time
	lastPush          time.Time

	mutex sync.RWMutex
}

//Allowed weather and pollution reporting parameters
var allowedParameters = make(map[string]struct{})

var apiParameters = []string{"winddir", "windspeedmph", "windgustmph",
	"windgustdir", "windspdmph_avg2m", "winddir_avg2m", "windgustmph_10m",
	"windgustdir_10m", "humidity", "dewptf", "tempf", "rainin", "dailyrainin",
	"baromin", "weather", "clouds", "soiltempf", "soilmoisture", "leafwetness",
	"solarradiation", "visibility", "indoortempf", "indoorhumidity",
	//Pollution Parameters
	"AqNO", "AqNO2T", "AqNO2", "AqNO2Y", "AqNOX", "AqNOY", "AqNO3", "AqSO4",
	"AqSO2", "AqSO2T", "AqCO", "AqCOT", "AqEC", "AqOC", "AqBC", "AqUV", "AqPM2.5",
	"AqPM10", "AqOZONE",
	//Additional Parameters to extend API
	"tempc", "indoortempc", "barohpa"}

//init initialise library
func init() {
	for _, v := range apiParameters {
		allowedParameters[v] = struct{}{}
	}
}

//String interface to show station details
func (s *Station) String() string {
	return fmt.Sprintf("Station Details - id: %s, Last Weather Update: %s,  Last Push: %s", s.id, s.lastUpdateWeather, s.lastPush)
}

//New generates a new pws object
func New(stnID string, stnPassword string, software string, calculateDewpoint bool) *Station {
	s := &Station{id: stnID, password: stnPassword, softwaretype: software, calculateDewpoint: calculateDewpoint}
	log.Println(s)
	return s
}

//UpdateWeather update weather parameters for next PushUpdate
func (s *Station) UpdateWeather(parameters map[string]string) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if parametersOK(parameters) {

		s.weatherParameters = parameters
		s.lastUpdateWeather = time.Now()

		return nil
	}
	return errors.New("Error: Invalid parameter")
}

//check parameters are in WU API
func parametersOK(parameters map[string]string) bool {

	for k := range parameters {
		_, ok := allowedParameters[k]
		if !ok {
			return false
		}
	}

	return true
}

//PushUpdate pushes latest params to WU, if a time is provided it uses tries to use the time if not it will use the current time.
func (s *Station) PushUpdate(date string) error {
	d, err := handleDate(date)
	if err != nil {
		return err
	}

	url := s.buildURL(s.weatherParameters, d)

	_, err = upload(url)
	if err != nil {
		return err
	}

	s.lastPush = time.Now()

	return nil
}

//buildURL generates an update URL in accordance with the API, from parameters passed in.
func (s *Station) buildURL(parameters map[string]string, date string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	u := url.Values{}

	u.Set("ID", s.id)
	u.Set("PASSWORD", s.password)
	u.Set("dateutc", date)

	for k, v := range parameters {

		if k == "tempc" {
			k = "tempf"
			v = celciusToFahrenheit(v)
		}

		if k == "barohpa" {
			k = "baromin"
			v = hpaToInhg(v)
		}

		if k == "indoortempc" {
			k = "indoortempf"
			v = celciusToFahrenheit(v)
		}

		//If v == "" an error has occured in conversion of the parameter or no value has been provided.
		if v != "" {
			u.Set(k, v)
		}
	}

	//Calculate Dewpoint if required.
	if s.calculateDewpoint {
		if parameters["tempf"] != "" && parameters["humidity"] != "" {
			u.Set("dewptf", dewpointFahrenheit(parameters["tempf"], parameters["humidity"]))
		}
	}

	u.Set("softwaretype", s.softwaretype)
	u.Set("action", "updateraw")

	completeURL := fmt.Sprintf("%s%s", baseURL, u.Encode())
	return completeURL
}

//handleDate formats the date to comply with the WU API for the dateutc parameter
func handleDate(date string) (string, error) {
	var d = time.Now().UTC()
	var err error

	if date != "" {
		d, err = time.Parse(timeFormat, date)
		if err != nil {
			return "", err
		}
	}

	return (d.Format(timeFormat)), nil
}

//upload data using a formatted url that complies to the WU API
func upload(urlFinal string) (string, error) {

	resp, err := http.Get(urlFinal)
	if err != nil {
		fmt.Println("upload err:  ", err)
		return "", err
	}
	defer resp.Body.Close()

	return resp.Status, nil
}

//celciusToFahrenheit helper to convert degrees Celcius To Fahrenheit
func celciusToFahrenheit(c string) string {
	if s, err := strconv.ParseFloat(c, 32); err == nil {
		return fmt.Sprintf("%.2f", (s*9/5)+32)
	}
	return ""
}

//hpaToInhg helper to convert hPa (millibar) to InHg inches of mercury
func hpaToInhg(p string) string {
	if s, err := strconv.ParseFloat(p, 32); err == nil {
		//29.92inHg/1013.25hPa(mbar) = 0.02952874414014
		return fmt.Sprintf("%.2f", s*0.02952874414014)
	}
	return ""
}

//fahrenheitToCelcius helper to convert degrees Celcius To Fahrenheit
func fahrenheitToCelcius(f string) string {
	if s, err := strconv.ParseFloat(f, 32); err == nil {
		return fmt.Sprintf("%.2f", (s-32)*5/9)
	}
	return ""
}

//dewpointFahrenheit helper to calculate dewpoint, uses temperature and humidity
func dewpointFahrenheit(t, h string) string {
	return celciusToFahrenheit(dewpointCelcius(fahrenheitToCelcius(t), h))
}

//dewpointCelcius helper to calculate dewpoint, uses temperature and humidity
func dewpointCelcius(t, h string) string {
	var hs float64
	var err error
	if hs, err = strconv.ParseFloat(h, 32); err != nil {
		return ""
	}

	if ts, err := strconv.ParseFloat(t, 32); err == nil {
		k := (math.Log10(hs)-2)/0.4343 + (17.62*ts)/(243.12+ts)
		return fmt.Sprintf("%.2f", 243.12*k/(17.62-k))
	}
	return ""
}

//floatToString to two decimal places
func floatToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

//stringToFloat returns a float
func stringToFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return f, err
	}
	return f, nil
}
