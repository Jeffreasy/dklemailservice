package api

import (
	"dklautomationgo/logger"
	"dklautomationgo/services"
	"dklautomationgo/utils"
	"net/http"
)

// TelegramBotHandler bevat de API endpoints voor de Telegram bot
type TelegramBotHandler struct {
	telegramBotService *services.TelegramBotService
}

// NewTelegramBotHandler maakt een nieuwe TelegramBotHandler
func NewTelegramBotHandler(telegramBotService *services.TelegramBotService) *TelegramBotHandler {
	return &TelegramBotHandler{
		telegramBotService: telegramBotService,
	}
}

// RegisterRoutes registreert alle routes voor de TelegramBotHandler
func (h *TelegramBotHandler) RegisterRoutes(router *http.ServeMux) {
	// Secured routes
	router.HandleFunc("/api/v1/telegrambot/config", h.HandleJWTMiddleware(h.GetConfig))
	router.HandleFunc("/api/v1/telegrambot/send", h.HandleJWTMiddleware(h.SendMessage))
	router.HandleFunc("/api/v1/telegrambot/commands", h.HandleJWTMiddleware(h.GetCommands))
}

// GetConfig haalt de configuratie van de Telegram bot op
func (h *TelegramBotHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check of de service bestaat
	if h.telegramBotService == nil {
		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"enabled":  false,
			"message":  "Telegram bot service is niet geactiveerd",
			"chatId":   "",
			"commands": []string{},
		})
		return
	}

	// Haal basisgegevens op
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"enabled": true,
		"message": "Telegram bot service is actief",
		"chatId":  h.telegramBotService.GetChatID(),
	})
}

// SendMessage stuurt een bericht via de Telegram bot
func (h *TelegramBotHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check of de service bestaat
	if h.telegramBotService == nil {
		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Telegram bot service is niet geactiveerd",
		})
		return
	}

	// Parse message uit request body
	type MessageRequest struct {
		Message string `json:"message"`
	}

	var req MessageRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Stuur het bericht
	err := h.telegramBotService.SendMessage(req.Message)
	if err != nil {
		logger.Error("Fout bij verzenden Telegram bericht", "error", err)
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Fout bij verzenden bericht: " + err.Error(),
		})
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Bericht succesvol verzonden",
	})
}

// GetCommands haalt de beschikbare commando's op
func (h *TelegramBotHandler) GetCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check of de service bestaat
	if h.telegramBotService == nil {
		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success":  false,
			"message":  "Telegram bot service is niet geactiveerd",
			"commands": []interface{}{},
		})
		return
	}

	// Haal de commando's op
	commands := h.telegramBotService.GetCommands()

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"message":  "Commando's succesvol opgehaald",
		"commands": commands,
	})
}

// HandleJWTMiddleware is een helper functie voor JWT middleware
func (h *TelegramBotHandler) HandleJWTMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return HandleJWTMiddleware(handler)
}
