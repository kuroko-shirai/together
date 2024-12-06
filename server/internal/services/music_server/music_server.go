package music_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/kuroko-shirai/task"
	"github.com/kuroko-shirai/together/server/internal/services/music_server/player"
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
	player   player.Player
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
		player: *player.New(),
	}, nil
}

// Запустить сервер
func (s *MusicServer) Run() {
	fmt.Println("Сервер запущен на порту 8080")

	for {
		g := task.WithRecover(
			func(p any, args ...any) {
				log.Println("panic:", p)
			},
		)

		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		g.Do(
			func() func() error {
				return func() error {
					err := s.handle(conn)
					if err != nil {
						return err
					}

					return nil
				}
			}(),
		)

		g.Wait()
	}
}

// Обработать подключение клиента
func (s *MusicServer) handle(conn net.Conn) error {
	defer conn.Close()

	// Отправить текущее состояние воспроизведения клиенту
	s.mu.Lock()
	state, err := json.Marshal(s.state)
	s.mu.Unlock()

	if err != nil {
		return err
	}

	_, err = conn.Write(state)
	if err != nil {
		return err
	}

	// Обработать команды от клиента
	buf := make([]byte, 1024)
	var perr error
	for {
		n, err := conn.Read(buf)
		if err != nil {
			perr = err
			break
		}

		cmd := string(buf[:n])
		fmt.Println("Получена команда:", cmd)

		switch cmd {
		case "play":
			s.mu.Lock()
			s.state.IsPlaying = true
			go s.player.Play("./playlist/track-001.mp3")
			s.mu.Unlock()
		case "pause":
			s.mu.Lock()
			s.state.IsPlaying = false
			go s.player.Pause()
			s.mu.Unlock()
		case "next":
			s.mu.Lock()
			s.state.CurrentTrack = "Трек 2"
			s.mu.Unlock()
		case "prev":
			s.mu.Lock()
			s.state.CurrentTrack = "Трек 1"
			s.mu.Unlock()
		case "panic":
			panic(errors.New("error!!"))
		}

		// Отправить обновленное состояние воспроизведения
		// клиенту
		s.mu.Lock()
		state, err = json.Marshal(s.state)
		s.mu.Unlock()

		if err != nil {
			perr = err
			break
		}

		_, err = conn.Write(state)
		if err != nil {
			perr = err
			break
		}
	}

	return perr
}

func (s *MusicServer) Stop() error {
	return s.listener.Close()
}
