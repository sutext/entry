package main

import (
	"log"

	"sutext.github.io/entry/server"
)

func main() {
	// 创建服务器实例
	s := server.New()
	
	// 启动服务器
	log.Println("Starting server...")
	if err := s.Serve(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
