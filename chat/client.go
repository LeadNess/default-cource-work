package main

import "default-cource-work/chat/client"

func main()  {
	c := client.NewClient()
	c.Dial(":3333")
	c.Start()
	c.Send()
}