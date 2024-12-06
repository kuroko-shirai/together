package music_server

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

// Структура для хранения состояния воспроизведения
type PlaybackState struct {
	CurrentTrack string `json:"current_track"`
	IsPlaying    bool   `json:"is_playing"`
}

// Сервер для управления воспроизведением музыки
type MusicServer struct {
	listener net.Listener
	state    PlaybackState
	mu       sync.Mutex
}

// Создать новый сервер
func New(config *Config) (*MusicServer, error) {
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	return &MusicServer{
		listener: listener,
		state: PlaybackState{
			CurrentTrack: "",
			IsPlaying:    false,
		},
	}, nil
}

// Запустить сервер
func (s *MusicServer) Run() {
	fmt.Println("Сервер запущен на порту 8080")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// Обработать подключение клиента
func (s *MusicServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Отправить текущее состояние воспроизведения клиенту
	s.mu.Lock()
	state, err := json.Marshal(s.state)
	s.mu.Unlock()

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(state)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Обработать команды от клиента
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		cmd := string(buf[:n])
		fmt.Println("Получена команда:", cmd)

		switch cmd {
		case "play":
			s.mu.Lock()
			s.state.IsPlaying = true
			s.mu.Unlock()
		case "pause":
			s.mu.Lock()
			s.state.IsPlaying = false
			s.mu.Unlock()
		case "next":
			s.mu.Lock()
			s.state.CurrentTrack = "Трек 2"
			s.mu.Unlock()
		case "prev":
			s.mu.Lock()
			s.state.CurrentTrack = "Трек 1"
			s.mu.Unlock()
		}

		// Отправить обновленное состояние воспроизведения клиенту
		s.mu.Lock()
		state, err = json.Marshal(s.state)
		s.mu.Unlock()

		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = conn.Write(state)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (s *MusicServer) Stop() {
	s.listener.Close()
}
