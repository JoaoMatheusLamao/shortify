package fieldmapper

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/clbanning/mxj"
	"github.com/gin-gonic/gin"
)

// Mapper is a handler function that processes the input and returns the paths
func Mapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input map[string]interface{}

		contentType := c.GetHeader("Content-Type")
		if err := bindInput(c, contentType, &input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		paths := parseInput(input, "")
		c.JSON(http.StatusOK, gin.H{"paths": paths})
	}
}

func bindInput(c *gin.Context, contentType string, input *map[string]interface{}) error {
	switch contentType {
	case "application/json":
		return c.ShouldBindJSON(input)
	case "application/xml":
		body, err := c.GetRawData()
		if err != nil {
			return err
		}
		mv, err := mxj.NewMapXml(body)
		if err != nil {
			return err
		}
		*input = mv.Old()
		return nil
	default:
		return fmt.Errorf("unsupported content type")
	}
}

// parseInput parses the input JSON or XML and returns a list of paths
func parseInput(data map[string]interface{}, prefix string) []string {
	var paths []string
	for key, value := range data {
		fullPath := buildFullPath(prefix, key)
		paths = append(paths, parseValue(value, fullPath)...)
	}
	sort.Strings(paths)
	return paths
}

func buildFullPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func parseValue(value interface{}, fullPath string) []string {
	switch v := value.(type) {
	case map[string]interface{}:
		return parseInput(v, fullPath)
	case []interface{}:
		return parseArray(v, fullPath)
	default:
		return []string{fullPath}
	}
}

func parseArray(array []interface{}, fullPath string) []string {
	if len(array) == 0 {
		return []string{fullPath + "[]"}
	}
	if itemMap, ok := array[0].(map[string]interface{}); ok {
		return parseInput(itemMap, fullPath+"[]")
	}
	return []string{fullPath + "[]"}
}
