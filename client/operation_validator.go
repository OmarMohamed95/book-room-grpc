package client

// Validate the incoming operation
func Validate(operation string) {
	for o := range getValidOperationsMap() {
		if operation == o {
			return
		}
	}

	panic("Provided operation is invalid")
}
