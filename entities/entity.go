package entities

type Entity interface {
	Fields() ([]string, []interface{})
	TableName() string
}

// GetFieldNames return all field names from entity
func GetFieldNames(e Entity) []string {
	fieldName, _ := e.Fields()
	return fieldName
}

// GetScanFields return scan fields from entity
func GetScanFields(e Entity, names []string) []interface{} {
	fieldNames, fieldMaps := e.Fields()

	n := len(fieldMaps)
	if len(names) < n {
		n = len(names)
	}

	fields := make([]interface{}, 0, n)
	for _, n := range names {
		for i, fieldName := range fieldNames {
			if fieldName == n {
				fields = append(fields, fieldMaps[i])
				break
			}
		}
	}

	return fields
}
