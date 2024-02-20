package main

//AdminAPI ReInitialize Static APIs
func AdminAPI(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {
	switch RequestJSON["endpoint"] {
	default:
		ResponseJSON["status"] = 404
		ResponseJSON["message"] = "UnKnown Admin API Call"
		break
	}
}
