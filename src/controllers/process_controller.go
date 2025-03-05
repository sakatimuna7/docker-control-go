package controllers

import (
	"context"
	"docker-control-go/src/configs"
	"docker-control-go/src/constant"
	"docker-control-go/src/helpers"
	logger "docker-control-go/src/log"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Constant keys
const (
	UserIDKey   = "userID"
	UserRoleKey = "userRole"
)

// type PM2Process struct {
// 	Name   string  `json:"name"`
// 	PMID   int     `json:"pm_id"`
// 	Status string  `json:"pm2_env.status"`
// 	CPU    float64 `json:"monit.cpu"`
// 	Memory float64 `json:"monit.memory"`
// }

type ProcessList struct {
	Processes []map[string]interface{} `json:"processes"`
}

// GetProcessList retrieves the list of PM2 processes
func GetProcessList(ctx context.Context, userID string, userRole string) ([]map[string]interface{}, error) {
	output, err := exec.Command("pm2", "jlist").Output()
	if err != nil {
		return nil, err
	}

	var processList ProcessList
	if err := json.Unmarshal(output, &processList.Processes); err != nil {
		return nil, err
	}
	var filteredProcesses []map[string]interface{}
	for _, process := range processList.Processes {
		identityName := fmt.Sprintf("pm2:%v", process["name"])

		// Cek izin user terhadap container ini
		permittedActions := map[string]bool{
			"read":    false,
			"start":   false,
			"stop":    false,
			"restart": false,
		}

		// Jika user adalah admin, beri akses penuh
		if userRole == "admin" {
			for action := range permittedActions {
				permittedActions[action] = true
			}
		} else {
			// Cek apakah user memiliki izin read
			allowedRead, _ := configs.Enforcer.Enforce(userID, identityName, "read")
			if !allowedRead {
				continue // Jika tidak punya izin read, skip container ini
			}

			// Cek izin lainnya
			actions := []string{"read", "start", "restart", "stop"}
			for _, action := range actions {
				allowed, _ := configs.Enforcer.Enforce(userID, identityName, action)
				permittedActions[action] = allowed
			}

			// Set izin read ke true karena user bisa melihat container ini
			permittedActions["read"] = true
		}

		name := process["name"]

		pm2Env, ok := process["pm2_env"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type for pm2_env")
		}
		env, ok := pm2Env["env"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type for env")
		}

		monit, ok := process["monit"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type for monit")
		}

		pmID := pm2Env["pm_id"]
		restartTime := pm2Env["restart_time"]
		status := pm2Env["status"]
		mode := env["exec_mode"]
		execPath := env["pm_exec_path"]
		uptime := pm2Env["pm_uptime"]
		version := pm2Env["version"]
		pid := process["pid"]
		cpu := monit["cpu"]
		memory := monit["memory"]
		watch := pm2Env["watch"]

		filteredProcesses = append(filteredProcesses, map[string]interface{}{
			"name":             name,
			"pm_id":            pmID,
			"pid":              pid,
			"mode":             mode,
			"exec_path":        execPath,
			"uptime":           uptime,
			"restart_time":     restartTime,
			"status":           status,
			"version":          version,
			"cpu":              cpu,
			"memory":           memory,
			"watch":            watch,
			"permitted_action": permittedActions,
		})
	}

	return filteredProcesses, nil
}

// HandlePM2Command executes PM2 commands with proper permission validation
func HandlePM2Command(c *websocket.Conn, command string, ctx context.Context) {
	var cmd struct {
		Action       string `json:"action"`       // "start", "stop", "restart", etc.
		IdentityName string `json:"identityName"` // Process name or ID
	}
	if err := json.Unmarshal([]byte(command), &cmd); err != nil {
		sendWSMessagePM2(c, "failed", "Invalid command format", cmd.Action, cmd.IdentityName)
		return
	}
	// 🔹 Ambil user ID dan user Role dari context
	userID, ok := ctx.Value(constant.UserIDKey).(string)
	if !ok {
		sendWSMessagePM2(c, "failed", "Unauthorized: Missing userID", cmd.Action, cmd.IdentityName)
		return
	}

	userRole, _ := ctx.Value(constant.UserRoleKey).(string) // Jika kosong, asumsi bukan admin

	if userRole == "admin" {
		// Admin selalu diizinkan
		exec.Command("pm2", cmd.Action, cmd.IdentityName).Run()
		sendWSMessage(c, "success", "Command executed successfully", cmd.Action, cmd.IdentityName)
		return
	}

	allowed, err := configs.Enforcer.Enforce(userID, "pm2", cmd.IdentityName)
	if err != nil {
		sendWSMessage(c, "failed", err.Error(), cmd.Action, cmd.IdentityName)
		return
	}

	if !allowed {
		sendWSMessage(c, "failed", "Permission denied", cmd.Action, cmd.IdentityName)
		return
	}

	exec.Command("pm2", cmd.Action, cmd.IdentityName).Run()
	sendWSMessage(c, "success", "Command executed successfully", cmd.Action, cmd.IdentityName)
}

// PM2Controller handles WebSocket connections for PM2 process monitoring
func PM2Controller(c *websocket.Conn, ctx context.Context) {
	userID, ok := ctx.Value(constant.UserIDKey).(string)
	if !ok {
		log.Println("❌ Error: userID not found in context")
		c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Unauthorized: Missing userID"}`))
		return
	}

	userRole, _ := ctx.Value(constant.UserRoleKey).(string)

	if userRole != "admin" {
		allowed, err := configs.Enforcer.Enforce(userRole, "resource:pm2", "read")
		if err != nil {
			log.Println("❌ Error checking permission:", err)
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Internal Server Error"}`))
			return
		}

		if !allowed {
			log.Println("❌ Access denied for user:", userID)
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Forbidden"}`))
			return
		}
	}

	// 🔹 3. Channel untuk menerima perintah dari client
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}

			// 🔥 Tangani perintah dari client
			HandlePM2Command(c, string(message), ctx)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			processList, err := GetProcessList(ctx, userID, userRole)

			if err != nil {
				log.Println("❌ Error fetching process list:", err)
				c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Failed to fetch process list"}`))
				continue
			}

			response, _ := json.Marshal(processList)
			c.WriteMessage(websocket.TextMessage, response)
		case <-ctx.Done():
			log.Println("WebSocket connection closed")
			return
		}
	}
}

func sendWSMessagePM2(c *websocket.Conn, status, message, action, identityName string) {
	response := map[string]string{
		"status":       status, // "in_progress", "success", "failed"
		"message":      message,
		"action":       action,
		"identityName": identityName,
	}
	data, _ := json.Marshal(response)
	c.WriteMessage(websocket.TextMessage, data)
}

// fiber
// Struct untuk request body
type NewDeamonPayload struct {
	ExecPath   string `json:"exec_path"`
	DeamonName string `json:"name"`
}

// PM2NewDeamon membuat daemon baru menggunakan PM2
func PM2NewDeamon(c *fiber.Ctx) error {
	// Ambil user ID dan user Role dari context
	userID := c.Locals("userID").(string)
	userRole := c.Locals("userRole").(string)

	var req NewDeamonPayload
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse request: ", err)
		return helpers.ErrorResponse(c, 400, "Invalid request", err)
	}

	// Cek izin dengan Casbin jika bukan admin
	if userRole != "admin" {
		allowed, err := configs.Enforcer.Enforce(userID, "pm2", "new-deamon")
		if err != nil {
			return helpers.ErrorResponse(c, 500, "Internal Server Error", err)
		}

		if !allowed {
			return helpers.ErrorResponse(c, 403, "Forbidden", nil)
		}
	}

	// Validasi ExecPath dengan `which`
	whichCmd := exec.Command("which", req.ExecPath)
	whichOutput, err := whichCmd.Output()
	if err != nil || len(strings.TrimSpace(string(whichOutput))) == 0 {
		logger.Log.Error("Invalid ExecPath: ", req.ExecPath)
		return helpers.ErrorResponse(c, 400, "Invalid ExecPath: command not found", nil)
	}

	// Gunakan path yang dikembalikan oleh `which`
	execPath := strings.TrimSpace(string(whichOutput))

	// Jalankan perintah PM2
	cmd := exec.Command("pm2", "start", execPath, "--name", req.DeamonName)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Error("Failed to create daemon: ", string(output), err)
		return helpers.ErrorResponse(c, 500, "Failed to create daemon", err)
	}

	return helpers.SuccessResponse(c, 201, "Daemon created successfully", string(output))
}
