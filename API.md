## Personal Weather Station Upload Protocol

*This is a direct copy from WU PWS Upload API page [here](http://wiki.wunderground.com/index.php/PWS_-_Upload_Protocol), which is formatted for MD*

Here is how to create a Personal Weather Station update to wunderground.com:

### Background
To upload a weather condition, you make a standard HTTP GET request with the ID, PASSWORD and weather conditions as GET parameters

### URL
Here is the URL used in the uploading (if you browse here without parameters you will get a brief usage):
```
http://weatherstation.wunderground.com/weatherstation/updateweatherstation.php
```

### GET parameters
NOT all fields need to be set, the _required_ elements are:
- action
    - ID
    - PASSWORD
    - dateutc
- IMPORTANT all fields must be url escaped
  - reference http://www.w3schools.com/tags/ref_urlencode.asp
- example

```
  2001-01-01 10:32:35
   becomes
  2000-01-01+10%3A32%3A35
```

- if the weather station is not capable of producing a timestamp, our system will accept "now". Example:

```
dateutc=now
```

- list of fields:

```
action [action=updateraw] -- always supply this parameter to indicate you are making a weather observation upload
ID [ID as registered by wunderground.com]
PASSWORD [PASSWORD registered with this ID, case sensative]
dateutc - [YYYY-MM-DD HH:MM:SS (mysql format)] In Universal Coordinated Time (UTC) Not local time
winddir - [0-360 instantaneous wind direction]
windspeedmph - [mph instantaneous wind speed]
windgustmph - [mph current wind gust, using software specific time period]
windgustdir - [0-360 using software specific time period]
windspdmph_avg2m  - [mph 2 minute average wind speed mph]
winddir_avg2m - [0-360 2 minute average wind direction]
windgustmph_10m - [mph past 10 minutes wind gust mph ]
windgustdir_10m - [0-360 past 10 minutes wind gust direction]
```

```
humidity - [% outdoor humidity 0-100%]
dewptf - [F outdoor dewpoint F]
```

```
tempf - [F outdoor temperature]
 * for extra outdoor sensors use temp2f, temp3f, and so on
 ```

 ```
rainin - [rain inches over the past hour)] -- the accumulated rainfall in the past 60 min
dailyrainin - [rain inches so far today in local time]
```

```
baromin - [barometric pressure inches]
```

```
weather - [text] -- metar style (+RA)
clouds - [text] -- SKC, FEW, SCT, BKN, OVC
```

```
soiltempf - [F soil temperature]
 * for sensors 2,3,4 use soiltemp2f, soiltemp3f, and soiltemp4f
soilmoisture - [%]
* for sensors 2,3,4 use soilmoisture2, soilmoisture3, and soilmoisture4
```

```
leafwetness  - [%]
+ for sensor 2 use leafwetness2
```

```
solarradiation - [W/m^2]
UV - [index]
```

```
visibility - [nm visibility]
```

```
indoortempf - [F indoor temperature F]
indoorhumidity - [% indoor humidity 0-100]
```

- Pollution Fields:

```
AqNO - [ NO (nitric oxide) ppb ]
AqNO2T - (nitrogen dioxide), true measure ppb
AqNO2 - NO2 computed, NOx-NO ppb
AqNO2Y - NO2 computed, NOy-NO ppb
AqNOX - NOx (nitrogen oxides) - ppb
AqNOY - NOy (total reactive nitrogen) - ppb
AqNO3 -NO3 ion (nitrate, not adjusted for ammonium ion) UG/M3
AqSO4 -SO4 ion (sulfate, not adjusted for ammonium ion) UG/M3
AqSO2 -(sulfur dioxide), conventional ppb
AqSO2T -trace levels ppb
AqCO -CO (carbon monoxide), conventional ppm
AqCOT -CO trace levels ppb
AqEC -EC (elemental carbon) – PM2.5 UG/M3
AqOC -OC (organic carbon, not adjusted for oxygen and hydrogen) – PM2.5 UG/M3
AqBC -BC (black carbon at 880 nm) UG/M3
AqUV-AETH  -UV-AETH (second channel of Aethalometer at 370 nm) UG/M3
AqPM2.5 - PM2.5 mass - UG/M3
AqPM10 - PM10 mass - PM10 mass
AqOZONE - Ozone - ppb
```

```
softwaretype - [text] ie: WeatherLink, VWS, WeatherDisplay
```

### Example URL
Here is an example of a full URL:
```
http://weatherstation.wunderground.com/weatherstation/updateweatherstation.php?ID=KCASANFR5&PASSWORD=XXXXXX&dateutc=2000-01-01+10%3A32%3A35&winddir=230&windspeedmph=12&windgustmph=12&tempf=70&rainin=0&baromin=29.1&dewptf=68.2&humidity=90&weather=&clouds=&softwaretype=vws%20versionxx&action=updateraw
```
 - *NOTE:* not all fields need to be set

### Response Text
The response from an HTTP GET request contains some debugging data.
RESPONSES and MEANINGS:
response

```
"success"
```

- the observation was ingested successfully
response

```
"INVALIDPASSWORDID|Password and/or id are incorrect"
```

 - invalid user data entered in the ID and PASSWORD GET parameters

response

```
<b>RapidFire Server</b><br><br>
<b>usage</b><br>
```

- the minimum GET parameters ID, PASSWORD, action, and dateutc were not set

### RapidFire Updates
RapidFire Updates allow you to update weather station conditions at a frequency up to once observation every 2.5 seconds. Web site visitors will see these observations change in real-time on the wunderground.com site.
- A real-time update should look almost like the standard update.
- However, the server to request is:
  - rtupdate.wunderground.com, not weatherstation.wunderground.com
- And, please add to the query string:
  - &realtime=1&rtfreq=2.5
  - where rtrfreq is the frequency of updates in seconds.
- here is an example:

```
http://rtupdate.wunderground.com/weatherstation/updateweatherstation.php?ID=KCASANFR5&PASSWORD=XXXXXX&dateutc=2000-01-01+10%3A32%3A35&winddir=230&windspeedmph=12&windgustmph=12&tempf=70&rainin=0&baromin=29.1&dewptf=68.2&humidity=90&weather=&clouds=&softwaretype=vws%20versionxx&action=updateraw&realtime=1&rtfreq=2.5
```
- We haven't decided whether you should also send standard updates if the user is uploading in real-time. For now, I think we are leaning toward either only sending standard updates, or only sending real-time updates, not requiring that both are sent when the user is in real-time mode. That'll be simpler.

### Automatically generate a Weather Underground station ID
To automatically register a station on Weather Underground please use the following URL to create a station id.
- All fields followed by an equal sign "=" need to be filled in by you.
  - Lat
  - Lon
  - Neighborhood
  - Stationtype
  - Email
  - Password (if you are using an email address that already has an account on wunderground, you will have to include "&password=xxxxxx" at the end of the URL.)

```
http://api.wunderground.com/api/c39ea9c961f1a516/pwssignup/view.json?lat=43.592501&lon=-70.786009&neighborhood=YourStationName&stationtype=Davis&email=abc@gmail.net&allownews=1
```
