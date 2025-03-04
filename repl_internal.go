package main

import (
	"fmt"
	"internal"
	"time"
)

func main_() {
	cache := internal.NewCache(time.Second * 5)
	fmt.Println("Cache created")
	cache.Add("pikachu", []byte(`{"name": "Pikachu"}`))
	time.Sleep(time.Second * 2)
	fmt.Println(cache.Items)
	cache.Add("doraemon", []byte(`{"name": "Doraemon"}`))
	time.Sleep(time.Second * 2)
	fmt.Println(cache.Items)
	cache.Add("bulbasaur", []byte(`{"name": "Bulbasaur"}`))
	time.Sleep(time.Second * 2)
	fmt.Println(cache.Items)

}
