package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/GregoryDosh/metrotransit"
	log "github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"github.com/valyala/fasthttp"
)

var departuresType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Departure",
	Fields: graphql.Fields{
		"actual": &graphql.Field{
			Type: graphql.Boolean,
		},
		"block_number": &graphql.Field{
			Type: graphql.Int,
		},
		"departure_text": &graphql.Field{
			Type: graphql.String,
		},
		"departure_time": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"gate": &graphql.Field{
			Type: graphql.String,
		},
		"map": &graphql.Field{
			Type: graphql.String,
		},
		"route": &graphql.Field{
			Type: graphql.String,
		},
		"route_direction": &graphql.Field{
			Type: graphql.String,
		},
		"terminal": &graphql.Field{
			Type: graphql.String,
		},
		"vehicle_heading": &graphql.Field{
			Type: graphql.Int,
		},
		"vehicle_latitude": &graphql.Field{
			Type: graphql.Float,
		},
		"vehicle_longitude": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

var detailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Details",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"code": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"latitude": &graphql.Field{
			Type: graphql.Float,
		},
		"longitude": &graphql.Field{
			Type: graphql.Float,
		},
		"zone_id": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"location_type": &graphql.Field{
			Type: graphql.Int,
		},
		"wheelchair_boarding": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var stopType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Stop",
	Fields: graphql.Fields{
		"departures": &graphql.Field{
			Type: graphql.NewList(departuresType),
		},
		"stop_details": &graphql.Field{
			Type: detailsType,
		},
		"stop_id": &graphql.Field{
			Type: graphql.Int,
		},
		"update_time": &graphql.Field{
			Type: graphql.String,
		},
		"full_map": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"departures": &graphql.Field{
			Type:        stopType,
			Description: "Gets NexTrip departures for a given stop_id.",
			Args: graphql.FieldConfigArgument{
				"stop_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"map_scale": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"map_width": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"map_height": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: resolveDepartures,
		},
	},
})

type newStopType struct {
	Departures []newDepartureType `json:"departures"`
	Details    newDetailsType     `json:"stop_details"`
	FullMap    string             `json:"full_map"`
	StopID     int                `json:"stop_id"`
	UpdateTime time.Time          `json:"update_time"`
}

type newDepartureType struct {
	Actual           bool    `json:"actual"`
	BlockNumber      int     `json:"block_number"`
	DepartureText    string  `json:"departure_text"`
	DepartureTime    string  `json:"departure_time"`
	Description      string  `json:"description"`
	Gate             string  `json:"gate"`
	Map              string  `json:"map,omitempty"`
	Route            string  `json:"route"`
	RouteDirection   string  `json:"route_direction"`
	Terminal         string  `json:"terminal"`
	VehicleHeading   int     `json:"vehicle_heading"`
	VehicleLatitude  float64 `json:"vehicle_latitude"`
	VehicleLongitude float64 `json:"vehicle_longitude"`
}

type newDetailsType struct {
	ID                 int64   `json:"id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	ZoneID             string  `json:"zone_id"`
	URL                string  `json:"url"`
	LocationType       int64   `json:"location_type"`
	WheelchairBoarding int64   `json:"wheelchair_boarding"`
}

func mapHelper(stopLat float64, stopLong float64, busLatLong [][]float64, scale int, size string) string {
	if googleMapsAPIKey != "" {
		baseURL := fmt.Sprintf("maps.googleapis.com/maps/api/staticmap?markers=color:blue%%7Clabel:S%%7C%[1]f,%[2]f&scale=%[3]d", stopLat, stopLong, scale)
		busMarker := ""
		sizeParam := ""
		keyParm := fmt.Sprintf("&key=%s", googleMapsAPIKey)
		for _, bus := range busLatLong {
			busLat, busLong := bus[0], bus[1]
			if busLat != 0 && busLong != 0 {
				busMarker += fmt.Sprintf("&markers=%%7Ccolor:green%%7Clabel:B%%7C%[1]f,%[2]f", busLat, busLong)
			}
		}
		if size != "" {
			sizeParam = fmt.Sprintf("&size=%s", size)
		}

		return baseURL + busMarker + keyParm + sizeParam
	}
	return ""
}

func translateDepartureType(o *metrotransit.Stop, scale int, size string) *newStopType {
	newDepartureList := []newDepartureType{}
	allBusLatLong := [][]float64{}
	for _, oldDepart := range o.Departures {
		busLatLong := [][]float64{
			{oldDepart.VehicleLatitude, oldDepart.VehicleLongitude},
		}
		allBusLatLong = append(allBusLatLong, []float64{oldDepart.VehicleLatitude, oldDepart.VehicleLongitude})
		newDepartureList = append(newDepartureList, newDepartureType{
			Actual:           oldDepart.Actual,
			BlockNumber:      oldDepart.BlockNumber,
			DepartureText:    oldDepart.DepartureText,
			DepartureTime:    oldDepart.DepartureTime.String(),
			Description:      oldDepart.Description,
			Gate:             oldDepart.Gate,
			Map:              mapHelper(o.Details.Latitude, o.Details.Longitude, busLatLong, scale, size),
			Route:            oldDepart.Route,
			RouteDirection:   oldDepart.RouteDirection,
			Terminal:         oldDepart.Terminal,
			VehicleHeading:   oldDepart.VehicleHeading,
			VehicleLatitude:  oldDepart.VehicleLatitude,
			VehicleLongitude: oldDepart.VehicleLongitude,
		})
	}

	return &newStopType{
		StopID:     o.StopID,
		UpdateTime: o.UpdateTime,
		Departures: newDepartureList,
		FullMap:    mapHelper(o.Details.Latitude, o.Details.Longitude, allBusLatLong, scale, size),
		Details: newDetailsType{
			ID:                 o.Details.ID,
			Code:               o.Details.Code,
			Name:               o.Details.Name,
			Description:        o.Details.Description,
			Latitude:           o.Details.Latitude,
			Longitude:          o.Details.Longitude,
			ZoneID:             o.Details.ZoneID,
			URL:                o.Details.URL,
			LocationType:       o.Details.LocationType,
			WheelchairBoarding: o.Details.WheelchairBoarding,
		},
	}
}

func resolveDepartures(params graphql.ResolveParams) (interface{}, error) {
	scale, isOK := params.Args["map_scale"].(int)
	if !isOK {
		scale = 1
	}
	width, isOK := params.Args["map_width"].(int)
	if !isOK {
		width = 0
	}
	height, isOK := params.Args["map_height"].(int)
	if !isOK {
		height = 0
	}
	size := ""
	if width != 0 && height != 0 {
		size = fmt.Sprintf("%dx%d", width, height)
	}
	idQuery, isOK := params.Args["stop_id"].(int)
	if isOK {
		s, err := env.GetDepartures(idQuery)
		if err != nil {
			return nil, err
		}
		return translateDepartureType(s, scale, size), nil
	}
	return nil, errors.New("stop_id not found")
}

// define schema, with our rootQuery and rootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		log.Errorf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func graphqlHandler(ctx *fasthttp.RequestCtx) {
	result := executeQuery(string(ctx.QueryArgs().Peek("query")), schema)
	pretty := string(ctx.QueryArgs().Peek("pretty"))

	enc := json.NewEncoder(ctx)
	enc.SetEscapeHTML(false)

	if pretty != "" {
		enc.SetIndent("", "    ")
	}
	if err := enc.Encode(result); err != nil {
		log.Error(err)
	}
}
