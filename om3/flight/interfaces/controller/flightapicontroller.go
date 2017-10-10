package controller

import (
	"strings"

	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (

	// FlightApiController for cart api
	FlightApiController struct {
		responder.JSONAware `inject:""`
	}
)

func (cc *FlightApiController) SearchFlightsAction(ctx web.Context) web.Response {

	json := `[
      {
         "scheduledDateTime":"2017-9-24 12:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH120",
         "codeShareFlightNumbers":[
            "DE5337",
            "DE8991",
            "DE4494",
            "DE5706"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor"
      },
{
         "scheduledDateTime":"2017-9-24 14:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH121",
         "codeShareFlightNumbers":[
            "DE5339"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor"
      }
]`
	return &web.ContentResponse{
		Status: 200,
		Body:   strings.NewReader(json),
	}
}

func (cc *FlightApiController) SearchAirportsAction(ctx web.Context) web.Response {

	json := `[
	{
		"airportName":"Frankfurt",
		"airport":"FRA"
	},
	{
		"airportName":"Frankfurt",
		"airport":"FRA"
	}
]`
	return &web.ContentResponse{
		Status: 200,
		Body:   strings.NewReader(json),
	}
}

func (cc *FlightApiController) AutosuggestAction(ctx web.Context) web.Response {

	json := `{
"airports":[
	{
		"airportName":"Frankfurt",
		"airport":"FRA"
	},
	{
		"airportName":"Frankfurt",
		"airport":"FRA"
	}
],
	"flights":[
      {
         "scheduledDateTime":"2017-9-24 12:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH120",
         "codeShareFlightNumbers":[
            "DE5337",
            "DE8991",
            "DE4494",
            "DE5706"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor"
      },
{
         "scheduledDateTime":"2017-9-24 14:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH121",
         "codeShareFlightNumbers":[
            "DE5339"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor"
      }
]
}`
	return &web.ContentResponse{
		Status: 200,
		Body:   strings.NewReader(json),
	}
}

func (cc *FlightApiController) SearchFlightsPerAirlineAction(ctx web.Context) web.Response {

	json := `[
{
"airline": "CBB",
"airlineName": "Condor",
"airlineLogo": "logourl",
"flights": [
      {
         "scheduledDateTime":"2017-9-24 12:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH120",
         "codeShareFlightNumbers":[
            "DE5337",
            "DE8991",
            "DE4494",
            "DE5706"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor"
      },
{
         "scheduledDateTime":"2017-9-24 14:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH121",
         "codeShareFlightNumbers":[
            "DE5339"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor"
      }
   ]
},
{
"airline": "LH",
"airlineName": "Lufthansa",
"airlineLogo": "logourl",
"flights": [
      {
         "scheduledDateTime":"2017-9-24 12:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH120",
         "codeShareFlightNumbers":[
            "DE5337",
            "DE8991",
            "DE4494",
            "DE5706"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor"
      },
{
         "scheduledDateTime":"2017-9-24 14:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH121",
         "codeShareFlightNumbers":[
            "DE5339"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor"
      }
   ]
}
]`
	return &web.ContentResponse{
		Status: 200,
		Body:   strings.NewReader(json),
	}
}

func (cc *FlightApiController) SaveFlightAction(ctx web.Context) web.Response {

	json := `{"message": "ok"}`
	return &web.ContentResponse{
		Status: 200,
		Body:   strings.NewReader(json),
	}
}

func (cc *FlightApiController) GetSessionFlightAction(ctx web.Context) web.Response {

	json := `{
         "scheduledDateTime":"2017-9-24 12:30",
         "scheduledDate":"2017-9-24",
         "primaryFlightNumber":"LH120",
         "codeShareFlightNumbers":[
            "DE5337",
            "DE8991",
            "DE4494",
            "DE5706"
         ],
         "airport":"ABB",
         "airportName":"London-Heathrow Airport",
         "cityName":"London",
         "airline":"CBB",
         "airlineName":"Condor",
	"terminal": "2",
	 "destinationRegionCat": "schengen"
      }`
	return &web.ContentResponse{
		Status: 200,
		Body:   strings.NewReader(json),
	}
}