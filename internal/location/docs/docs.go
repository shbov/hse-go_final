// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/drivers": {
            "get": {
                "description": "Поиск водителей по заданным координатам и радиусу",
                "consumes": [
                    "application/json"
                ],
                "summary": "GetDriversByLocation",
                "parameters": [
                    {
                        "type": "number",
                        "description": "Latitude in decimal degrees",
                        "name": "lat",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Longitude in decimal degrees",
                        "name": "lng",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Radius in meters",
                        "name": "radius",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    }
                }
            }
        },
        "/drivers/{driver_id}/location": {
            "post": {
                "description": "Обновление данных о позиции водителя",
                "consumes": [
                    "application/json"
                ],
                "summary": "SetDriverLocation",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "ID of driver",
                        "name": "driver_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Latitude and longitude  in decimal degrees",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.SetDriverLocationBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        }
    },
    "definitions": {
        "requests.SetDriverLocationBody": {
            "type": "object",
            "properties": {
                "lat": {
                    "type": "number"
                },
                "lng": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Location service",
	Description:      "This is a location service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
