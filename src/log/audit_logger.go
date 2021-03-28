package log

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/stevenhay/qiqi_bot/src/config"
)

type LogLevel int

const (
	logLevelInfo = iota
	logLevelWarn
	logLevelError
)

func LogInfo(s *discordgo.Session, author, title, message string) {
	logrus.WithTime(time.Now()).Infof("LogInfo={author=%s, title=%s, description=%s}", author, title, message)
	_, err := s.ChannelMessageSendEmbed(config.Get().AuditChannelID, generateEmbed(logLevelInfo, author, title, message))
	if err != nil {
		msg := fmt.Sprintf("Failed to send embed (%s), falling back to text: {time=%s, author=%s, title=%s, description=%s}", err, time.Now().Local().Format(time.RFC3339), author, title, message)
		logrus.WithTime(time.Now()).Error(err)

		_, err = s.ChannelMessageSend(config.Get().AuditChannelID, msg)
		logrus.WithTime(time.Now()).Error(err)
	}
}

func LogError(s *discordgo.Session, author, title, message string) {
	logrus.WithTime(time.Now()).Errorf("LogError={author=%s, title=%s, description=%s}", author, title, message)
	_, err := s.ChannelMessageSendEmbed(config.Get().ErrorChannelID, generateEmbed(logLevelError, author, title, message))
	if err != nil {
		msg := fmt.Sprintf("Failed to send embed (%s), falling back to text: {time=%s, author=%s, title=%s, description=%s}", err, time.Now().Local().Format(time.RFC3339), author, title, message)
		logrus.WithTime(time.Now()).Error(err)

		_, err = s.ChannelMessageSend(config.Get().ErrorChannelID, msg)
		logrus.WithTime(time.Now()).Error(err)
	}
}

func LogWarn(s *discordgo.Session, author, title, message string) {
	logrus.WithTime(time.Now()).Warnf("LogWarn{author=%s, title=%s, description=%s}", author, title, message)
	_, err := s.ChannelMessageSendEmbed(config.Get().ErrorChannelID, generateEmbed(logLevelWarn, author, title, message))
	if err != nil {
		msg := fmt.Sprintf("Failed to send embed (%s), falling back to text: {time=%s, author=%s, title=%s, description=%s}", err, time.Now().Local().Format(time.RFC3339), author, title, message)
		logrus.WithTime(time.Now()).Error(err)

		_, err = s.ChannelMessageSend(config.Get().ErrorChannelID, msg)
		logrus.WithTime(time.Now()).Error(err)
	}
}

func generateEmbed(loglevel LogLevel, author, title, message string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{}
	embed.Timestamp = time.Now().Local().Format(time.RFC3339)
	embed.Title = title
	switch loglevel {
	case logLevelInfo:
		embed.Color = 0x45f5ec
	case logLevelWarn:
		embed.Color = 0xffa600
	case logLevelError:
		embed.Color = 0xb40000
	}
	embed.Description = message
	embed.Footer = &discordgo.MessageEmbedFooter{Text: "Authored by " + author}
	return embed
}
