package main

//CallStaticAPIsv1 ReInitialize Static APIs
func CallStaticAPIsv1(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {
	Log("---------------------API Version: " + InterfaceToString(RequestJSON["api_version"]) + "---------------------")
	switch RequestJSON["endpoint"] {
	default:
		ResponseJSON["status"] = 404
		ResponseJSON["message"] = "UnKnown API Call"
		break
	}
}
