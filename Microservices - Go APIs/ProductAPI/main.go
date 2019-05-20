package main

func main() {

	//port := os.Getenv("PORT")
	//if len(port) == 0 {
	//	port = "5001"
	//}

	server := NewServer()
	server.Run(":" + "5001")
}
