package main

import (
	"encoding/json"
	"fmt"
	"net"
)

// Структура для хранения состояния воспроизведения
type PlaybackState struct {
	CurrentTrack string `json:"current_track"`
	IsPlaying    bool   `json:"is_playing"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Получить текущее состояние воспроизведения от сервера
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	var state PlaybackState
	err = json.Unmarshal(buf[:n], &state)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Текущий трек:", state.CurrentTrack)
	fmt.Println("Воспроизведение:", state.IsPlaying)

	// Отправить команды серверу
	for {
		fmt.Print("Введите команду (play, pause, next, prev, exit): ")
		var cmd string
		fmt.Scanln(&cmd)

		if cmd == "exit" {
			return
		}

		_, err = conn.Write([]byte(cmd))
		if err != nil {
			fmt.Println(err)
			return
		}

		// Получить обновленное состояние воспроизведения от сервера
		n, err = conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = json.Unmarshal(buf[:n], &state)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Текущий трек:", state.CurrentTrack)
		fmt.Println("Воспроизведение:", state.IsPlaying)
	}
}
