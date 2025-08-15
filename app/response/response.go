package response

import (
	message "app/app/messsage"
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Regexp definitions
var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
var wordBarrierRegex = regexp.MustCompile(`([a-z_0-9])([A-Z])`)

type conventionalMarshallerFromPascal struct {
	Value any
}

func convertToCamelCase(marshalled []byte) []byte {
	return keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			// Empty keys are valid JSON, only lowercase if we do not have an empty key.
			if len(match) > 2 {
				// Convert to camel case
				converted := bytes.ToLower(wordBarrierRegex.ReplaceAll(
					match,
					[]byte(`${1}_${2}`),
				))

				// Remove underscores and capitalize the following letter
				var result []byte
				underscore := false
				for i := 1; i < len(converted)-1; i++ {
					if converted[i] == '_' {
						underscore = true
					} else {
						if underscore {
							result = append(result, byte(unicode.ToUpper(rune(converted[i]))))
							underscore = false
						} else {
							result = append(result, converted[i])
						}
					}
				}
				result = append([]byte{converted[0]}, result...)
				result = append(result, converted[len(converted)-1])
				return result
			}
			return match
		},
	)
}

func (c conventionalMarshallerFromPascal) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(c.Value)
	if err != nil {
		return nil, err
	}
	naming := viper.GetString("HTTP_JSON_NAMING")

	val := reflect.TypeOf(c.Value)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		field, ok := val.FieldByName("json")
		if ok {
			if field.Tag.Get("naming") != "" {
				naming = field.Tag.Get("naming")
			}
		}
	}

	var converted []byte
	switch naming {
	case "snake_case":
		// https://gist.github.com/Rican7/39a3dc10c1499384ca91
		converted = keyMatchRegex.ReplaceAllFunc(
			marshalled,
			func(match []byte) []byte {
				return bytes.ToLower(wordBarrierRegex.ReplaceAll(
					match,
					[]byte(`${1}_${2}`),
				))
			},
		)
	case "camel_case":
		converted = convertToCamelCase(marshalled)
	case "pascal_case":
		return marshalled, nil
	default:
		return marshalled, nil
	}

	return converted, nil
}

type Response struct {
	Code     string            `json:"code"`
	Message  string            `json:"message"`
	Data     any               `json:"data,omitempty"`
	Paginate *ResponsePaginate `json:"paginate,omitempty"`
}

type ResponsePaginate struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

func defaultJSON(ctx *gin.Context, code int, msgID string, data any, paginate *ResponsePaginate, params ...map[string]string) error {
	// Simplified without i18n for now
	msg := msgID
	if msg == "" {
		msg = "Unknown error"
	}

	ctx.JSON(code, conventionalMarshallerFromPascal{Response{
		Message:  msg,
		Code:     strconv.Itoa(code),
		Data:     data,
		Paginate: paginate,
	}})
	return nil
}

func JSON(ctx *gin.Context, code int, data any) error {
	ctx.JSON(code, conventionalMarshallerFromPascal{data})
	return nil
}

// Success 200 success
func Success(ctx *gin.Context, data any) error {
	return defaultJSON(ctx, http.StatusOK, message.Success, data, nil)
}

func StatusContinue(ctx *gin.Context, message string, data any) error {
	return defaultJSON(ctx, http.StatusAccepted, message, data, nil)
}

// Paginate 200 success
func SuccessWithPaginate(ctx *gin.Context, data any, size, page, total int) error {
	paginate := &ResponsePaginate{
		Page:  page,
		Size:  size,
		Total: total,
	}
	return defaultJSON(ctx, http.StatusOK, message.Success, data, paginate)
}

// BadRequest 400 other and external error
func BadRequest(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusBadRequest, message, data, nil)
}

// Unauthorized 401 un authentication
func Unauthorized(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusUnauthorized, message, data, nil)
}

// Forbidden 403 No permission
func Forbidden(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusForbidden, message, data, nil)
}

// ValidateFailed 412 Validate error
func ValidateFailed(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusPreconditionFailed, message, data, nil)
}

// InternalServerError 500 internal error
func InternalServerError(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusInternalServerError, message, data, nil)
}

// NotImplemented 501 not implemented
func NotImplemented(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusNotImplemented, message, data, nil)
}
