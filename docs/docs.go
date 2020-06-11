// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-06-11 16:49:41.783949 +0300 EEST m=+0.042024293

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "RDXP3",
            "url": "https://confluence.taxibeat.com/display/TEAM404/RDXP3+-+TechnoMules"
        },
        "license": {
            "name": "BEAT Mobility Services"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/hecho_transito": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bollobas"
                ],
                "summary": "Get return all the traffic incidents",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "start date (epoch time)",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "end date (epoch time)",
                        "name": "to",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit Value",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset Value",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/view.TrafficIncident"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/view.ErrorSwagger"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/view.ErrorSwagger"
                        }
                    }
                }
            }
        },
        "/stats_operador": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bollobas"
                ],
                "summary": "Get return all the operator stats",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit Value",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset Value",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/view.OperatorStats"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/view.ErrorSwagger"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/view.ErrorSwagger"
                        }
                    }
                }
            }
        },
        "/viajes_agregados": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bollobas"
                ],
                "summary": "Get All the aggregated rides",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "start date (epoch time)",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "end date (epoch time)",
                        "name": "to",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit Value",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset Value",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/view.AggregatedTrips"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/view.ErrorSwagger"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/view.ErrorSwagger"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "view.AggregatedTrips": {
            "type": "object",
            "properties": {
                "accesibilidad": {
                    "type": "number"
                },
                "dist_pasajero": {
                    "type": "number"
                },
                "dist_pasajero_eod": {
                    "type": "number"
                },
                "dist_solicitud": {
                    "type": "number"
                },
                "dist_solicitud_eod": {
                    "type": "number"
                },
                "dist_vacio": {
                    "type": "number"
                },
                "dist_vacio_eod": {
                    "type": "number"
                },
                "fecha": {
                    "type": "string"
                },
                "fin_eod": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "id_proveedor": {
                    "type": "string"
                },
                "inicio_eod": {
                    "type": "integer"
                },
                "multiplicador_eod": {
                    "type": "number"
                },
                "operador_mujer": {
                    "type": "number"
                },
                "tiempo_pasajero": {
                    "type": "number"
                },
                "tiempo_pasajero_eod": {
                    "type": "integer"
                },
                "tiempo_solicitud": {
                    "type": "number"
                },
                "tiempo_solicitud_eod": {
                    "type": "number"
                },
                "tiempo_vacio": {
                    "type": "number"
                },
                "tiempo_vacio_eod": {
                    "type": "number"
                },
                "tot_veh_disp": {
                    "type": "integer"
                },
                "tot_veh_viaje": {
                    "type": "integer"
                },
                "tot_viajes": {
                    "type": "integer"
                }
            }
        },
        "view.ErrorSwagger": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "view.OperatorStats": {
            "type": "object",
            "properties": {
                "cant_viajes": {
                    "type": "integer"
                },
                "edad": {
                    "type": "string"
                },
                "genero": {
                    "type": "integer"
                },
                "horas_conectado": {
                    "type": "string"
                },
                "horas_viaje": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "id_operador": {
                    "type": "string"
                },
                "ingreso_totales": {
                    "type": "string"
                },
                "tiempo_registro": {
                    "type": "integer"
                }
            }
        },
        "view.TrafficIncident": {
            "type": "object",
            "properties": {
                "distancia_viaje": {
                    "type": "string"
                },
                "hecho_trans": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "licencia": {
                    "type": "string"
                },
                "placa": {
                    "type": "string"
                },
                "tiempo_hecho": {
                    "type": "string"
                },
                "tiempo_viaje": {
                    "type": "string"
                },
                "ubicacion": {
                    "type": "string"
                }
            }
        }
    },
    "tags": [
        {
            "name": "bollobas"
        }
    ]
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{"http"},
	Title:       "Bollobas",
	Description: "Bollobas microservice is responsible for any analytics that go through Beat's backend platform.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
