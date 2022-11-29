package minecraft

import (
	"fmt"
	"time"

	"github.com/yinshaoze/BSProxy/config"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net/packet"
)

func generateKickMessage(s *config.ConfigProxyService, name packet.String) chat.Message {
	return chat.Message{
		Color: chat.White,
		Extra: []chat.Message{
			{Bold: true, Color: chat.Red, Text: "BS"},
			{Bold: true, Text: "Proxy"},
			{Text: " - "},
			{Bold: true, Color: chat.Gold, Text: "Connection Rejected\n"},

			{Text: "Your connection request is refused by BSProxy.\n"},
			{Text: "Reason: "},
			{Color: chat.LightPurple, Text: "You don't have permission to access this service.\n"},
			{Text: "Please contact the Administrators for help.\n\n"},

			{
				Color: chat.Gray,
				Text: fmt.Sprintf("Timestamp: %d | Player Name: %s | Service: %s\n",
					time.Now().UnixMilli(), name, s.Name),
			},
		},
	}
}

func generatePlayerNumberLimitExceededMessage(s *config.ConfigProxyService, name packet.String) chat.Message {
	return chat.Message{
		Color: chat.White,
		Extra: []chat.Message{
			{Bold: true, Color: chat.Red, Text: "BS"},
			{Bold: true, Text: "Proxy"},
			{Text: " - "},
			{Bold: true, Color: chat.Gold, Text: "Connection Rejected\n"},

			{Text: "Your connection request is refused by BSProxy.\n"},
			{Text: "Reason: "},
			{Color: chat.LightPurple, Text: "Service online player number limitation exceeded.\n"},
			{Text: "Please contact the Administrators for help.\n\n"},

			{
				Color: chat.Gray,
				Text: fmt.Sprintf("Timestamp: %d | Player Name: %s | Service: %s\n",
					time.Now().UnixMilli(), name, s.Name),
			},
			{Text: "GitHub: "},
			{
				Color: chat.Aqua, UnderLined: true,
				Text:       "https://github.com/yinshaoze/BSProxy",
				ClickEvent: chat.OpenURL("https://github.com/yinshaoze/BSProxy"),
			},
		},
	}
}
