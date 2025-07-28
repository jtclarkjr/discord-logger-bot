package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	sess         *discordgo.Session
	logFile      string
	cache        = make(map[string]cachedMessage)
	cacheMutex   sync.Mutex
	cacheTimeout = 24 * time.Hour
)

type cachedMessage struct {
	Content   string
	Author    string
	Channel   string
	Timestamp time.Time
}

func main() {
	logFile = filepath.Join(".", "mod_logs.txt")

	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN not set")
	}

	var err error
	sess, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Register handlers
	sess.AddHandler(messageCreateHandler)
	sess.AddHandler(messageDeleteHandler)

	// Set intents to receive message events and content
	sess.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent

	// Start cache cleanup goroutine
	go func() {
		for {
			time.Sleep(time.Hour)
			cleanupCache()
		}
	}()

	http.HandleFunc("/bot/mod-logger/on", startBotHandler)
	http.HandleFunc("/bot/mod-logger/off", stopBotHandler)

	go func() {
		fmt.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

// messageCreateHandler logs messages and caches them for deletion tracking
func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.GuildID == "" {
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		channel = &discordgo.Channel{Name: "Unknown"}
	}

	cacheMutex.Lock()
	cache[m.ID] = cachedMessage{
		Content:   m.Content,
		Author:    m.Author.String(),
		Channel:   channel.Name,
		Timestamp: time.Now(),
	}
	cacheMutex.Unlock()

	logEntry := fmt.Sprintf("[%s] | [MESSAGE] | Channel: %s | Author: %s | Content: %s\n",
		time.Now().Format(time.RFC3339),
		channel.Name,
		m.Author.String(),
		escapeNewlines(m.Content),
	)
	appendLog(logEntry)
}

// messageDeleteHandler logs deleted messages using the cache
func messageDeleteHandler(s *discordgo.Session, d *discordgo.MessageDelete) {
	channel, err := s.State.Channel(d.ChannelID)
	channelName := "Unknown"
	if err == nil {
		channelName = channel.Name
	}

	cacheMutex.Lock()
	cached, ok := cache[d.ID]
	if ok {
		delete(cache, d.ID)
	}
	cacheMutex.Unlock()

	var logEntry string
	if ok {
		logEntry = fmt.Sprintf("[%s] | [DELETED] | Channel: %s | Author: %s | Content: %s\n",
			time.Now().Format(time.RFC3339),
			channelName,
			cached.Author,
			escapeNewlines(cached.Content),
		)
	} else {
		logEntry = fmt.Sprintf("[%s] | [DELETED] | Channel: %s | Author: Unknown | Content: Unknown (not cached)\n",
			time.Now().Format(time.RFC3339),
			channelName,
		)
	}
	appendLog(logEntry)
}

// appendLog writes a log entry to the log file and stdout
func appendLog(entry string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Failed to open log file:", err)
		return
	}
	defer f.Close()
	f.WriteString(entry)
	fmt.Print(entry)
}

// escapeNewlines replaces newlines with \n for log formatting
func escapeNewlines(s string) string {
	return string([]byte(s))
}

// cleanupCache removes old cached messages
func cleanupCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	now := time.Now()
	for id, msg := range cache {
		if now.Sub(msg.Timestamp) > cacheTimeout {
			delete(cache, id)
		}
	}
}

// startBotHandler starts the Discord bot session
func startBotHandler(w http.ResponseWriter, r *http.Request) {
	if sess != nil && sess.State != nil && sess.State.Ready.Version != 0 {
		fmt.Fprintln(w, "Bot is already running.")
		return
	}
	err := sess.Open()
	if err != nil {
		http.Error(w, "Failed to start bot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Bot started successfully.")
}

// stopBotHandler stops the Discord bot session
func stopBotHandler(w http.ResponseWriter, r *http.Request) {
	if sess == nil || sess.State == nil || sess.State.Ready.Version == 0 {
		fmt.Fprintln(w, "Bot is not running.")
		return
	}
	err := sess.Close()
	if err != nil {
		http.Error(w, "Failed to stop bot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Bot stopped successfully.")
}
