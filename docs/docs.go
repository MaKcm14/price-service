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
        "/api/markets": {
            "get": {
                "description": "this endpoint provides getting the current markets that are supported by the service",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service-Info"
                ],
                "summary": "markets getting",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.SupportedMarkets"
                        }
                    }
                }
            }
        },
        "/products/filter/markets": {
            "get": {
                "description": "this endpoint provides filtering products from marketplaces without any specified filtration",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Common-Filters"
                ],
                "summary": "common filtering",
                "parameters": [
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "iphone+11",
                        "description": "the exact query string",
                        "name": "query",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "enum": [
                                "wildberries",
                                "megamarket"
                            ],
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "megamarket+wildberries",
                        "description": "the list of the markets using for search",
                        "name": "markets",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "the num of products' sample",
                        "name": "sample",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "popular",
                            "pricedown",
                            "priceup",
                            "newly"
                        ],
                        "type": "string",
                        "default": "popular",
                        "description": "the type of products' sample sorting",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "enum": [
                            0,
                            1
                        ],
                        "type": "integer",
                        "default": 1,
                        "description": "the flag that defines 'Should image links be parsed?'",
                        "name": "no-image",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "min",
                            "max"
                        ],
                        "type": "string",
                        "default": "min",
                        "description": "the amount of the products in response's sample",
                        "name": "amount",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chttp.ProductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    },
                    "502": {
                        "description": "Bad Gateway",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    }
                }
            }
        },
        "/products/filter/price/best-price": {
            "get": {
                "description": "this endpoint provides filtering products from marketplaces with the best and minimum price",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Price-Filters"
                ],
                "summary": "best price filtering",
                "parameters": [
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "iphone+11",
                        "description": "the exact query string",
                        "name": "query",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "enum": [
                                "wildberries",
                                "megamarket"
                            ],
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "megamarket+wildberries",
                        "description": "the list of the markets using for search",
                        "name": "markets",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "the num of products' sample",
                        "name": "sample",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "popular",
                            "pricedown",
                            "priceup",
                            "newly"
                        ],
                        "type": "string",
                        "default": "popular",
                        "description": "the type of products' sample sorting",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "enum": [
                            0,
                            1
                        ],
                        "type": "integer",
                        "default": 1,
                        "description": "the flag that defines 'Should image links be parsed?'",
                        "name": "no-image",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "min",
                            "max"
                        ],
                        "type": "string",
                        "default": "min",
                        "description": "the amount of the products in response's sample",
                        "name": "amount",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chttp.ProductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    },
                    "502": {
                        "description": "Bad Gateway",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    }
                }
            }
        },
        "/products/filter/price/best-price/async": {
            "post": {
                "description": "this endpoint provides filtering products from marketplaces with the best and minimum price in async mode",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Price-Filters"
                ],
                "summary": "async best price filtering",
                "parameters": [
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "iphone+11",
                        "description": "the exact query string",
                        "name": "query",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "enum": [
                                "wildberries",
                                "megamarket"
                            ],
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "megamarket+wildberries",
                        "description": "the list of the markets using for search",
                        "name": "markets",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "the num of products' sample",
                        "name": "sample",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "popular",
                            "pricedown",
                            "priceup",
                            "newly"
                        ],
                        "type": "string",
                        "default": "popular",
                        "description": "the type of products' sample sorting",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "enum": [
                            0,
                            1
                        ],
                        "type": "integer",
                        "default": 1,
                        "description": "the flag that defines 'Should image links be parsed?'",
                        "name": "no-image",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "min",
                            "max"
                        ],
                        "type": "string",
                        "default": "min",
                        "description": "the amount of the products in response's sample",
                        "name": "amount",
                        "in": "query"
                    },
                    {
                        "description": "the headers that need to be included into the async response",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/chttp.ExtraHeaders"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    }
                }
            }
        },
        "/products/filter/price/exact-price": {
            "get": {
                "description": "this endpoint provides filtering products from marketplaces with price in range (exact-price, exact-price * 1.05 (+5%))",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Price-Filters"
                ],
                "summary": "exact price filtering",
                "parameters": [
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "iphone+11",
                        "description": "the exact query string",
                        "name": "query",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "the value of exact price",
                        "name": "price",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "enum": [
                                "wildberries",
                                "megamarket"
                            ],
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "megamarket+wildberries",
                        "description": "the list of the markets using for search",
                        "name": "markets",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "the num of products' sample",
                        "name": "sample",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "popular",
                            "pricedown",
                            "priceup",
                            "newly"
                        ],
                        "type": "string",
                        "default": "popular",
                        "description": "the type of products' sample sorting",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "enum": [
                            0,
                            1
                        ],
                        "type": "integer",
                        "default": 1,
                        "description": "the flag that defines 'Should image links be parsed??'",
                        "name": "no-image",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "min",
                            "max"
                        ],
                        "type": "string",
                        "default": "min",
                        "description": "the amount of the products in response's sample",
                        "name": "amount",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chttp.ProductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    },
                    "502": {
                        "description": "Bad Gateway",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    }
                }
            }
        },
        "/products/filter/price/price-range": {
            "get": {
                "description": "this endpoint provides filtering products from marketplaces with specified price range",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Price-Filters"
                ],
                "summary": "price range filtering",
                "parameters": [
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "iphone+11",
                        "description": "the exact query string",
                        "name": "query",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "description": "the price range's lower bound: less than price_up",
                        "name": "price_down",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "the price range's upper bound: more than price_down",
                        "name": "price_up",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minLength": 1,
                        "type": "array",
                        "items": {
                            "enum": [
                                "wildberries",
                                "megamarket"
                            ],
                            "type": "string"
                        },
                        "collectionFormat": "ssv",
                        "example": "megamarket+wildberries",
                        "description": "the list of the markets using for search",
                        "name": "markets",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "the num of products' sample",
                        "name": "sample",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "popular",
                            "pricedown",
                            "priceup",
                            "newly"
                        ],
                        "type": "string",
                        "default": "popular",
                        "description": "the type of products' sample sorting",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "enum": [
                            0,
                            1
                        ],
                        "type": "integer",
                        "default": 1,
                        "description": "the flag that defines 'Should image links be parsed?'",
                        "name": "no-image",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "min",
                            "max"
                        ],
                        "type": "string",
                        "default": "min",
                        "description": "the amount of the products in response's sample",
                        "name": "amount",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chttp.ProductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    },
                    "502": {
                        "description": "Bad Gateway",
                        "schema": {
                            "$ref": "#/definitions/chttp.ResponseErr"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "chttp.ExtraHeaders": {
            "type": "object",
            "properties": {
                "headers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/chttp.Header"
                    }
                }
            }
        },
        "chttp.Header": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "chttp.ProductResponse": {
            "type": "object",
            "properties": {
                "samples": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/entities.ProductSample"
                    }
                }
            }
        },
        "chttp.ResponseErr": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "entities.Currency": {
            "type": "string",
            "enum": [
                "rub"
            ],
            "x-enum-varnames": [
                "RUB"
            ]
        },
        "entities.MarketView": {
            "type": "object",
            "properties": {
                "emoji": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "entities.Price": {
            "type": "object",
            "properties": {
                "base_price": {
                    "type": "integer"
                },
                "discount": {
                    "type": "integer"
                },
                "discount_price": {
                    "type": "integer"
                }
            }
        },
        "entities.Product": {
            "type": "object",
            "properties": {
                "brand": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "$ref": "#/definitions/entities.Price"
                },
                "related_links": {
                    "$ref": "#/definitions/entities.ProductLink"
                },
                "supplier": {
                    "type": "string"
                }
            }
        },
        "entities.ProductLink": {
            "type": "object",
            "properties": {
                "image_link": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "entities.ProductSample": {
            "type": "object",
            "properties": {
                "currency": {
                    "$ref": "#/definitions/entities.Currency"
                },
                "main_products_sample": {
                    "type": "string"
                },
                "market": {
                    "type": "string"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Product"
                    }
                }
            }
        },
        "entities.SupportedMarkets": {
            "type": "object",
            "properties": {
                "markets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.MarketView"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
