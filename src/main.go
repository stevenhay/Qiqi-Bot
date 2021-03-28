package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/stevenhay/qiqi_bot/src/commands"
	"github.com/stevenhay/qiqi_bot/src/config"
	"github.com/stevenhay/qiqi_bot/src/role"
)

var (
	token          string
	commandHandler *commands.CommandHandler
	roleHandler    *role.RoleHandler
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}

	discord.AddHandler(messageCreate)
	discord.AddHandler(guildMemberAdd)
	discord.AddHandler(messageReactionAdd)
	discord.AddHandler(messageReactionRemove)
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers

	config.Load()

	roleHandler = role.NewRoleHandler()
	commandHandler = commands.NewCommandHandler(
		discord,
		roleHandler,
	)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	fmt.Println("Qiqi Bot running!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandHandler.PerformCommand(m)
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	roleHandler.DefaultRole.ApplyDefaultRole(s, m.Member)
}

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	if r.MessageID == roleHandler.ColourRoles.GetMessageID() {
		roleHandler.ColourRoles.HandleReactionAdd(s, r)
	}
}

func messageReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == s.State.User.ID {
		return
	}

	if r.MessageID == roleHandler.ColourRoles.GetMessageID() {
		roleHandler.ColourRoles.HandleReactionRemove(s, r)
	}
}
