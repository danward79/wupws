# Weather Underground - Personal Weather Station Upload

Weather Underground (WU) mines weather data from multiple sources. These include personal weather stations (PWS). This go library provides an interface to the WU PWS upload api.

Weather Underground provide a simple web HTTP GET API. See [here](http://wiki.wunderground.com/index.php/PWS_-_Upload_Protocol). Also see [here](./api.md) for a transcription of the API.

In addition to the standard parameters, the this library provides the ability to use degrees Celcius. Instead of Fahrenheit, with conversion to match the API automatically. The additional parameters are:
  - tempc - This will convert a Celcius parameter to Fahrenheit for tempf
  - indoortempc - This will convert a Celcius parameter to Fahrenheit for indoortempc
  - barohpa - This will convert an atomospheric pressure parameter to inHg for baromin
