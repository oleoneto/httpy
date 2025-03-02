package assertions

/*

func IsPassing(debugFunc func(...any), expectations Expectations, response Response) bool {
	var isPassing bool

	// ----------------------------------
	// Status
	// ----------------------------------
	if expectations.HTTP.Status != nil {
		if response.response.StatusCode != *expectations.HTTP.Status {
			debugFunc(fmt.Sprintf("expected %v but got %v", expectations.HTTP.Status, response.response.StatusCode))
			isPassing = false
		}
	}

	// ----------------------------------
	// Headers
	// ----------------------------------
	if expectations.HTTP.Headers != nil {
		if !reflect.DeepEqual(response.Headers, expectations.HTTP.Headers) {
			debugFunc(fmt.Sprintf("expected %v but got %v", expectations.HTTP.Headers, response.Headers))
			isPassing = false
		}
	}

	// ----------------------------------
	// Body
	// ----------------------------------
	if expectations.HTTP.Body != nil {
		if !reflect.DeepEqual(response.Body, expectations.HTTP.Body) {
			debugFunc(fmt.Sprintf("expected %v but got %v", expectations.HTTP.Body, response.Body))
			isPassing = false
		}
	}

	return isPassing
}
*/
