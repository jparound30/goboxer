package gobox

const (
	applicationJson = "application/json"
)

func buildFieldsQueryParams(fields []string) string {
	var params = ""
	if fieldsLen := len(fields); fieldsLen != 0 {
		buffer := make([]byte, 0, 1024)
		buffer = append(buffer, "fields="...)
		for index, v := range fields {
			buffer = append(buffer, v...)
			if index != fieldsLen-1 {
				buffer = append(buffer, ',')
			}
		}
		params = string(buffer)
	}
	return params
}
