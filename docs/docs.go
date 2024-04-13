// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "always returns 200",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/connect/{userId}": {
            "get": {
                "description": "promotes an http request to a websocket connection",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "opens a connection",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/join": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "makes user join a room",
                "parameters": [
                    {
                        "description": "room data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.PostJoinRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/leave": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "makes user leave a room",
                "parameters": [
                    {
                        "description": "room data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.PostLeaveRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/room": {
            "post": {
                "description": "opens a new room and joins the specified driver",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "endpoints"
                ],
                "summary": "opens a new room",
                "parameters": [
                    {
                        "description": "room data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.PostRoomRequest"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "api.PostJoinRequest": {
            "type": "object",
            "properties": {
                "driver": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "api.PostLeaveRequest": {
            "type": "object",
            "properties": {
                "driver": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "api.PostRoomRequest": {
            "type": "object",
            "properties": {
                "driver": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Chat Engine API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
