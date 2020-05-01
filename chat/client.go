package main

import "default-cource-work/chat/client"

func main()  {
	client := client.NewClient()
	client.Start()
}