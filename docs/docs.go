// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://example.com/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://http://81.177.220.96/",
            "email": "lyoshabura@gmail.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/calendar": {
            "get": {
                "description": "Возвращает события календаря на указанный день.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "calendar"
                ],
                "summary": "Получить список событий",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Дата в формате YYYY-MM-DD",
                        "name": "date",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/nocsqlcpg.StudentEventSwagger"
                            }
                        }
                    },
                    "400": {
                        "description": "неверные данные",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/students": {
            "get": {
                "description": "Получить все записи студентов",
                "produces": [
                    "application/json"
                ],
                "summary": "Получить список студентов",
                "responses": {
                    "200": {
                        "description": "Список студентов",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/nocsqlcpg.StudentSwagger"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "description": "Добавляет нового студента в базу данных.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "students"
                ],
                "summary": "Создать студента",
                "parameters": [
                    {
                        "description": "Данные студента",
                        "name": "student",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/nocsqlcpg.CreateStudentSwagger"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/nocsqlcpg.StudentSwagger"
                        }
                    },
                    "400": {
                        "description": "неверные данные",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "nocsqlcpg.CreateStudentSwagger": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "orderCost": {
                    "type": "integer"
                },
                "orderDay": {
                    "type": "integer"
                },
                "orderTime": {
                    "$ref": "#/definitions/nocsqlcpg.OrderTime"
                },
                "school": {
                    "type": "string"
                },
                "sclass": {
                    "type": "string"
                }
            }
        },
        "nocsqlcpg.EventDate": {
            "type": "object",
            "properties": {
                "infinityModifier": {
                    "type": "integer"
                },
                "time": {
                    "type": "string"
                },
                "valid": {
                    "type": "boolean"
                }
            }
        },
        "nocsqlcpg.OrderCheck": {
            "type": "object",
            "properties": {
                "bool": {
                    "type": "boolean"
                },
                "valid": {
                    "type": "boolean"
                }
            }
        },
        "nocsqlcpg.OrderTime": {
            "type": "object",
            "properties": {
                "microseconds": {
                    "type": "integer"
                },
                "valid": {
                    "type": "boolean"
                }
            }
        },
        "nocsqlcpg.StudentEventSwagger": {
            "type": "object",
            "properties": {
                "eventDate": {
                    "$ref": "#/definitions/nocsqlcpg.EventDate"
                },
                "id": {
                    "type": "integer"
                },
                "orderCheck": {
                    "$ref": "#/definitions/nocsqlcpg.OrderCheck"
                },
                "orderCost": {
                    "type": "integer"
                },
                "orderTime": {
                    "$ref": "#/definitions/nocsqlcpg.OrderTime"
                },
                "studentID": {
                    "type": "integer"
                }
            }
        },
        "nocsqlcpg.StudentSwagger": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "orderCost": {
                    "type": "integer"
                },
                "orderDay": {
                    "type": "integer"
                },
                "orderTime": {
                    "$ref": "#/definitions/nocsqlcpg.OrderTime"
                },
                "school": {
                    "type": "string"
                },
                "sclass": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{"http", "https"},
	Title:            "Student and Calendar API",
	Description:      "API для управления студентами и событиями в календаре.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
