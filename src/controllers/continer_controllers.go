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

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var dockerClient *client.Client

// Inisialisasi Docker Client
func InitDockerClient() {
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error initializing Docker client:", err)
	}
}

// Mengambil daftar container yang sedang berjalan real-time
func GetRunningContainersWS(c *websocket.Conn, ctx context.Context) {
	fmt.Println("Client connected for Docker events")
	logger.Log.Info("Client connected for Docker events")

	// 🔹 1. Kirim daftar container saat client pertama kali connect
	sendContainerList(c, ctx)

	// 🔹 2. Dengarkan event perubahan container
	eventFilter := filters.NewArgs()
	eventFilter.Add("type", "container")

	msgs, errs := dockerClient.Events(ctx, events.ListOptions{Filters: eventFilter})

	// 🔹 3. Channel untuk menerima perintah dari client
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}

			// 🔥 Tangani perintah dari client
			handleContainerCommand(c, string(message), ctx)
		}
	}()

	for {
		select {
		case _, ok := <-msgs:
			if !ok {
				fmt.Println("Docker event channel closed")
				logger.Log.Info("Docker event channel closed")
				return
			}

			// 🔥 Perbarui daftar container ke client
			sendContainerList(c, ctx)

		case err := <-errs:
			if err != nil {
				fmt.Println("Error receiving Docker event:", err)
				logger.Log.Error("Error receiving Docker event:", err)
				c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Error receiving Docker events"}`))
			}
			return
		}
	}
}

// 🔹 Kirim daftar container ke client
func sendContainerList(c *websocket.Conn, ctx context.Context) {
	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		fmt.Println("Error updating container list:", err)
		logger.Log.Error("Error updating container list:", err)
		c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Failed to update containers"}`))
		return
	}

	// Ambil user ID dan user Role dari context
	userID, ok := ctx.Value(constant.UserIDKey).(string)
	if !ok {
		log.Println("❌ Error: userID not found in context")
		c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Unauthorized: Missing userID"}`))
		return
	}

	userRole, _ := ctx.Value(constant.UserRoleKey).(string) // Jika kosong, asumsi bukan admin

	// Mapping data
	var mappedContainers []map[string]interface{}
	for _, container := range containers {
		containerID := container.ID
		containerName := ""
		if len(container.Names) > 0 {
			containerName = container.Names[0]
		}

		// Cek izin user terhadap container ini
		permittedActions := map[string]bool{
			"read":    false,
			"running": false,
			"pause":   false,
			"restart": false,
			"stop":    false,
		}

		// Jika user adalah admin, beri akses penuh
		if userRole == "admin" {
			for action := range permittedActions {
				permittedActions[action] = true
			}
		} else {
			// Cek apakah user memiliki izin read
			allowedRead, _ := configs.Enforcer.Enforce(userID, containerID, "read")
			if !allowedRead {
				continue // Jika tidak punya izin read, skip container ini
			}

			// Cek izin lainnya
			actions := []string{"running", "pause", "restart", "stop"}
			for _, action := range actions {
				allowed, _ := configs.Enforcer.Enforce(userID, containerID, action)
				permittedActions[action] = allowed
			}

			// Set izin read ke true karena user bisa melihat container ini
			permittedActions["read"] = true
		}

		// Format data container
		mappedContainers = append(mappedContainers, map[string]interface{}{
			"container_name":   containerName,
			"container_id":     containerID,
			"image":            container.Image,
			"image_id":         container.ImageID,
			"ports":            container.Ports,
			"created":          container.Created,
			"state":            container.State,
			"status":           container.Status,
			"permitted_action": permittedActions,
		})
	}

	// **Jika tidak ada container yang boleh dilihat user, kirimkan array kosong (`[]`)**
	if len(mappedContainers) == 0 {
		c.WriteMessage(websocket.TextMessage, []byte(`[]`))
		return
	}

	// Kirim data yang telah diformat ke client
	data, _ := json.Marshal(mappedContainers)
	c.WriteMessage(websocket.TextMessage, data)
}

// 🔹 Tangani perintah dari client
func handleContainerCommand(c *websocket.Conn, command string, ctx context.Context) {
	var cmd struct {
		Action      string `json:"action"`      // "start", "stop", "restart"
		ContainerID string `json:"containerId"` // ID container
	}

	if err := json.Unmarshal([]byte(command), &cmd); err != nil {
		sendWSMessage(c, "failed", "Invalid command format", cmd.Action, cmd.ContainerID)
		return
	}

	// 🔹 Kirim status "in_progress" sebelum menjalankan perintah
	sendWSMessage(c, "in_progress", "Processing request...", cmd.Action, cmd.ContainerID)

	var err error
	switch cmd.Action {
	case "start":
		err = dockerClient.ContainerStart(ctx, cmd.ContainerID, container.StartOptions{})
	case "stop":
		err = dockerClient.ContainerStop(ctx, cmd.ContainerID, container.StopOptions{})
	case "restart":
		err = dockerClient.ContainerRestart(ctx, cmd.ContainerID, container.StopOptions{})
	default:
		sendWSMessage(c, "failed", "Unknown command", cmd.Action, cmd.ContainerID)
		return
	}

	// 🔹 Kirim status "success" atau "failed" berdasarkan hasil eksekusi
	if err != nil {
		sendWSMessage(c, "failed", fmt.Sprintf("Error: %s", err.Error()), cmd.Action, cmd.ContainerID)
	} else {
		sendWSMessage(c, "success", "Command executed successfully", cmd.Action, cmd.ContainerID)
		sendContainerList(c, ctx) // Kirim daftar container terbaru setelah perubahan
	}
}

// 🔹 Fungsi untuk mengirim pesan status ke WebSocket
func sendWSMessage(c *websocket.Conn, status, message, action, containerID string) {
	response := map[string]string{
		"status":      status, // "in_progress", "success", "failed"
		"message":     message,
		"action":      action,
		"containerId": containerID,
	}
	data, _ := json.Marshal(response)
	c.WriteMessage(websocket.TextMessage, data)
}

func GetRunningConainters(c *fiber.Ctx) error {
	// 🔹 1. Kirim daftar container saat client pertama kali connect
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
		All: true, // Ambil semua container (running & stopped)
	})

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to get running containers", err)
	}

	// return c.JSON(fiber.Map{"message": "get running containers", "data": containers, "error": err})
	return helpers.SuccessResponse(c, 200, "get running containers", containers)
}

// WebSocket handler untuk event real-time
func DockerEventsWS(c *websocket.Conn, ctx context.Context) {
	fmt.Println("Client connected for Docker events")

	msgs, errs := dockerClient.Events(context.Background(), events.ListOptions{})

	for {
		select {
		case event := <-msgs:
			jsonData, _ := json.Marshal(event)

			if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				fmt.Println("Error sending message:", err)
				return
			}

		case err := <-errs:
			fmt.Println("Error receiving Docker event:", err)
			return
		}
	}
}
