package telegram

import (
	"context"
	"fmt"
	"net"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotModel interface {
	Processor(ctx context.Context)
	Bot() *tgbotapi.BotAPI
}

type MODE int

const (
	DownloadMode MODE = iota + 1
	RemoveMode
)

type UserState struct {
	IsAdmin bool
	Mode    MODE
}

func welcomeMSG(chatID int64) string {
	welcome := strings.Builder{}
	welcome.WriteString(fmt.Sprintf("Доступ разрешен для: %d", int(chatID)))
	welcome.WriteString("\n")
	welcome.WriteString(ip())
	welcome.WriteString("\n")
	welcome.WriteString("\n")
	welcome.WriteString("Вставьте ссылку для отправки ее в очередь на скачивание")
	welcome.WriteString("\n")
	welcome.WriteString("Или выберите одну из опций ниже")

	return welcome.String()
}

func ip() string {
	result := strings.Builder{}

	result.WriteString("Host ip: ")
	ips, err := func() ([]net.IP, error) {
		ips, err := getLocalIPs()
		if err != nil {
			return nil, fmt.Errorf("FindIP(): %w", err)
		}

		return ips, nil
	}()
	if err != nil {
		return fmt.Errorf("error into search ip. Reason: %w", err).Error()
	}

	for i, v := range ips {
		result.WriteString(v.String())
		if i > 0 && i < len(ips)-1 {
			result.WriteString(", ")
		}
	}
	result.WriteString(" ")
	result.WriteString("\n")
	return result.String()
}

func getLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips, nil
}
