package anilist

import (
	"AnimeGUI/verniy"
)

func CustomMediaListGroupFieldEntries(field verniy.MediaListField, fields ...verniy.MediaListField) verniy.MediaListGroupField {
	str := []string{string(field)}
	for _, f := range fields {
		str = append(str, string(f))
	}
	params := map[string]interface{}{"sort": "UPDATED_TIME_DESC"}
	return verniy.MediaListGroupField(verniy.FieldObject("entries", params, str...))
}
