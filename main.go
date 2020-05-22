package main

import (
	"hekv/redis"
)

func main() {
	/*
		bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
		bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
		opts := gorocksdb.NewDefaultOptions()
		opts.SetBlockBasedTableFactory(bbto)
		opts.SetCreateIfMissing(true)
		db, err := gorocksdb.OpenDb(opts, "/path/to/db")
		defer db.Close()

		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		wo := gorocksdb.NewDefaultWriteOptions()
		ro := gorocksdb.NewDefaultReadOptions()
		err = db.Put(wo, []byte("foo"), []byte("bar"))
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		value, err := db.Get(ro, []byte("foo"))
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		defer value.Free()
		fmt.Printf("%+v\n", value.Data())
	*/
	//for i := 0; i < 101; i++ {
	//	fmt.Printf("%02d\n", i)
	//}
	server := redis.CreateServer()
	server.Run()
}
