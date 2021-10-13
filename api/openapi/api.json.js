var api_spec =
{
  "openapi": "3.0.0",
  "servers": [
    {
      "url": "http://ssup2ket.com",
      "description": "ssup2ket auth service"
    }
  ],
  "info": {
    "description": "ssup2ket auth service",
    "version": "1.0.0",
    "title": "ssup2ket auth service",
    "contact": {
      "email": "supsup5642@gmail.com"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    }
  },
  "components": {
    "schemas": {
      "TokenCreate": {
        "type": "object",
        "required": [
          "loginId",
          "password"
        ],
        "properties": {
          "loginId": {
            "type": "string"
          },
          "password": {
            "type": "string"
          }
        }
      },
      "TokenInfo": {
        "type": "object",
        "required": [
          "token",
          "issuedAt",
          "expiresAt"
        ],
        "properties": {
          "token": {
            "type": "string"
          },
          "issuedAt": {
            "type": "string",
            "format": "date-time"
          },
          "expiresAt": {
            "type": "string",
            "format": "date-time"
          }
        }
      },
      "UserCreate": {
        "type": "object",
        "required": [
          "loginId",
          "password",
          "role",
          "phone",
          "email"
        ],
        "properties": {
          "loginId": {
            "type": "string"
          },
          "password": {
            "type": "string"
          },
          "role": {
            "$ref": "#/components/schemas/UserRole"
          },
          "phone": {
            "type": "string"
          },
          "email": {
            "type": "string"
          }
        }
      },
      "UserUpdate": {
        "type": "object",
        "required": [
          "password",
          "role",
          "phone",
          "email"
        ],
        "properties": {
          "password": {
            "type": "string"
          },
          "role": {
            "$ref": "#/components/schemas/UserRole"
          },
          "phone": {
            "type": "string"
          },
          "email": {
            "type": "string"
          }
        }
      },
      "UserInfo": {
        "type": "object",
        "required": [
          "id",
          "loginId",
          "role",
          "phone",
          "email"
        ],
        "properties": {
          "id": {
            "type": "string"
          },
          "loginId": {
            "type": "string"
          },
          "role": {
            "$ref": "#/components/schemas/UserRole"
          },
          "phone": {
            "type": "string"
          },
          "email": {
            "type": "string"
          }
        }
      },
      "UserInfoList": {
        "type": "object",
        "required": [
          "users",
          "metadata"
        ],
        "properties": {
          "users": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/UserInfo"
            }
          },
          "metadata": {
            "$ref": "#/components/schemas/ListMeta"
          }
        }
      },
      "UserRole": {
        "type": "string",
        "enum": [
          "admin",
          "user"
        ]
      },
      "ListMeta": {
        "type": "object",
        "required": [
          "limit",
          "offset",
          "total"
        ],
        "properties": {
          "limit": {
            "type": "integer"
          },
          "offset": {
            "type": "integer"
          },
          "total": {
            "type": "integer"
          }
        }
      },
      "ErrorInfo": {
        "type": "object",
        "required": [
          "code",
          "message"
        ],
        "properties": {
          "code": {
            "type": "string"
          },
          "message": {
            "type": "string"
          }
        }
      }
    },
    "parameters": {
      "UserID": {
        "name": "UserID",
        "in": "path",
        "required": true,
        "schema": {
          "type": "string"
        }
      },
      "Offset": {
        "name": "Offset",
        "in": "query",
        "required": false,
        "schema": {
          "type": "integer",
          "minimum": 0
        }
      },
      "Limit": {
        "name": "Limit",
        "in": "query",
        "required": false,
        "schema": {
          "type": "integer",
          "default": 50
        }
      }
    }
  },
  "paths": {
    "/tokens": {
      "post": {
        "tags": [
          "token"
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/TokenCreate"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/TokenInfo"
                }
              }
            }
          },
          "400": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "401": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      }
    },
    "/users": {
      "get": {
        "parameters": [
          {
            "$ref": "#/components/parameters/Offset"
          },
          {
            "$ref": "#/components/parameters/Limit"
          }
        ],
        "tags": [
          "user"
        ],
        "responses": {
          "200": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/UserInfoList"
                }
              }
            }
          },
          "400": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "401": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "404": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      },
      "post": {
        "tags": [
          "user"
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UserCreate"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/UserInfo"
                }
              }
            }
          },
          "400": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "409": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      }
    },
    "/users/{UserID}": {
      "parameters": [
        {
          "$ref": "#/components/parameters/UserID"
        }
      ],
      "get": {
        "tags": [
          "user"
        ],
        "responses": {
          "200": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/UserInfo"
                }
              }
            }
          },
          "400": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "401": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "404": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      },
      "put": {
        "tags": [
          "user"
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UserUpdate"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": ""
          },
          "400": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "401": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "404": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      },
      "delete": {
        "tags": [
          "user"
        ],
        "responses": {
          "200": {
            "description": ""
          },
          "401": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "404": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      }
    },
    "/users/me": {
      "get": {
        "tags": [
          "user"
        ],
        "responses": {
          "200": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/UserInfo"
                }
              }
            }
          },
          "400": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "401": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "404": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      },
      "put": {
        "tags": [
          "user"
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UserUpdate"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": ""
          },
          "400": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "401": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "404": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      },
      "delete": {
        "tags": [
          "user"
        ],
        "responses": {
          "200": {
            "description": ""
          },
          "404": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "401": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          },
          "500": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorInfo"
                }
              }
            }
          }
        }
      }
    }
  }
}