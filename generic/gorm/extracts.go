package gorm

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/tochukaso/golang-study/generic/util"
)

func ExtractTargets(entity interface{}, columns []*structs.Field) []string {
	cs := make([]string, 0, len(columns))
	dtoStructs := structs.New(entity)
	for _, v := range columns {
		f := dtoStructs.Field(v.Name())
		cs = append(cs, extractGormColumn(f.Name(), f.Tag("gorm")))
	}
	return cs
}

func ExtractUpdateColumns(m interface{}) []string {
	val := reflect.ValueOf(m).Elem()

	res := []string{}

	for i := 0; i < val.NumField(); i++ {
		vf := val.Field(i)

		// フィールドがポインター型で、なおかつ、nilの場合、更新対象外なので除外する
		if vf.Kind() == reflect.Ptr && vf.IsNil() {
			continue
		}

		tf, tags := extractGormTag(val, i)

		// `gorm:"-"`のタグが付いているカラムは更新対象外
		if _, ok := tags["-"]; ok {
			continue
		}

		res = append(res, extractColumn(tags, tf.Name))
	}

	return res
}

func extractGormColumn(fieldName string, tag string) string {
	tags := parseTagSetting(tag, ";")

	// `gorm:"-"`のタグが付いているカラムは更新対象外
	if _, ok := tags["-"]; ok {
		return ""
	}

	return extractColumn(tags, fieldName)
}

func extractColumn(tags map[string]string, name string) string {
	specificColumn := tags["column"]
	if specificColumn == "" {
		return util.ConvSnakeFromCamel(name)
	}
	return specificColumn
}

func extractGormTag(val reflect.Value, i int) (reflect.StructField, map[string]string) {
	tf := val.Type().Field(i)
	tag := tf.Tag

	tags := parseTagSetting(tag.Get("gorm"), ";")
	return tf, tags
}

func parseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names[j])-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(strings.ToUpper(values[0]))

		if len(values) >= 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = k
		}
	}

	return settings
}
