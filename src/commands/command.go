package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/stevenhay/qiqi_bot/src/role"
)

type Handler interface {
	CanPerform(s *discordgo.Session, m *discordgo.Member) bool
	Do(s *discordgo.Session, msg *discordgo.MessageCreate)
}

type Command struct {
	Handler Handler
}

type CommandHandler struct {
	session *discordgo.Session
}

var registeredCommands map[string]Command

func registerCommands(roleHandler *role.RoleHandler) {
	registeredCommands = make(map[string]Command)

	registeredCommands["role"] = Command{
		Handler: RoleCommand{roleHandler},
	}
}

func NewCommandHandler(session *discordgo.Session, roleHandler *role.RoleHandler) *CommandHandler {
	registerCommands(roleHandler)
	return &CommandHandler{session}
}

func (c *CommandHandler) PerformCommand(msg *discordgo.MessageCreate) {
	if msg.Author.ID == c.session.State.User.ID {
		return
	}

	if !strings.HasPrefix(msg.Content, "!") {
		return
	}

	cmd, ok := registeredCommands[strings.Split(strings.TrimPrefix(msg.Content, "!"), " ")[0]]
	if !ok {
		return
	}

	if !cmd.Handler.CanPerform(c.session, msg.Member) {
		return
	}

	cmd.Handler.Do(c.session, msg)
}
