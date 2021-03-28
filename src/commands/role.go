package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/stevenhay/qiqi_bot/src/helpers"
	"github.com/stevenhay/qiqi_bot/src/log"
	"github.com/stevenhay/qiqi_bot/src/role"
)

const (
	defaultRoleErrorUnableToGetGuildRolesFmt string = "I was unable to get the guild roles because of %s"
	defaultRoleSetRoleMessageFmt             string = "Okay! Set default role to %s!"
	defaultRoleSetRoleAuditMessageFmt        string = "Default role changed from %s to %s"
	defaultRoleUnableToFindRole              string = "I was unable to find a role with that name"
	defaultRoleAlreadySetToThatFmt           string = "The default role is already %s!"
	defaultRoleCurrentRoleFmt                string = "The default role is currently set to %s"
)

type RoleCommand struct {
	roleHandler *role.RoleHandler
}

func (d RoleCommand) CanPerform(s *discordgo.Session, m *discordgo.Member) bool {
	if len(m.Roles) == 0 {
		return false
	}

	return helpers.MemberHasAnyPermissions(s, m, []int64{discordgo.PermissionManageRoles, discordgo.PermissionAdministrator})
}

func (d RoleCommand) Do(s *discordgo.Session, msg *discordgo.MessageCreate) {
	message := strings.Split(msg.Content, " ")
	if len(message) < 2 {
		// TODO: help message
		return
	}

	cmd := strings.ToLower(message[1])
	if cmd == "default" {
		d.HandleDefaultCommands(s, msg, message[2:])
	}

	if cmd == "colour-picker" {
		embed := d.generateColourEmbed()
		msg, _ := s.ChannelMessageSendEmbed(msg.ChannelID, embed)
		for _, f := range embed.Fields {
			if f.Name != "\u200b" {
				s.MessageReactionAdd(msg.ChannelID, msg.ID, f.Name)
			}
		}
		d.roleHandler.ColourRoles.SetMessageID(msg.ID)
	}
}

func (d RoleCommand) HandleDefaultCommands(s *discordgo.Session, msg *discordgo.MessageCreate, msgParts []string) {
	if len(msgParts) == 0 || strings.ToLower(msgParts[0]) == "get" {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(defaultRoleCurrentRoleFmt, d.roleHandler.DefaultRole.GetDefaultRole().Name))
	} else if strings.ToLower(msgParts[0]) == "set" {
		var role *discordgo.Role
		var err error

		// check if it's an @role
		id := helpers.RoleAtToID(msgParts[1])
		if id != "" {
			role, err = helpers.GetRoleForID(s, id)
		} else {
			// check for manual role name entered
			role, err = helpers.GetRoleForName(s, strings.Join(msgParts[1:], " "))
		}

		if err != nil {
			log.LogError(s, msg.Author.Username, "Setting default role failed!", fmt.Sprintf(defaultRoleErrorUnableToGetGuildRolesFmt, err.Error()))
			return
		}

		if role == nil {
			s.ChannelMessageSend(msg.ChannelID, defaultRoleUnableToFindRole)
			return
		}

		if role.Name == d.roleHandler.DefaultRole.GetDefaultRole().Name {
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(defaultRoleAlreadySetToThatFmt, role.Name))
			return
		}

		log.LogInfo(s, msg.Author.Username, "Default Role Changed", fmt.Sprintf(defaultRoleSetRoleAuditMessageFmt, d.roleHandler.DefaultRole.GetDefaultRole().Name, role.Name))
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(defaultRoleSetRoleMessageFmt, role.Name))

		d.roleHandler.DefaultRole.SetDefaultRole(role, msg.Author.Username, msg.Author.ID)
	}
}

func (d RoleCommand) generateColourEmbed() *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{}
	embed.Title = "Colour Role Picker"
	embed.Description = "Description"
	embed.Color = 0xb40000

	keys := make([]string, 0)
	for k := range d.roleHandler.ColourRoles.GetEmojiToRoleMapping() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys[:3] {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   k,
			Value:  helpers.RoleIDToAt(d.roleHandler.ColourRoles.GetEmojiToRoleMapping()[k]),
			Inline: true,
		})
	}

	for i := 0; i < 3; i++ {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "\u200b",
			Value:  "\u200b",
			Inline: true,
		})
	}

	for _, k := range keys[3:] {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   k,
			Value:  helpers.RoleIDToAt(d.roleHandler.ColourRoles.GetEmojiToRoleMapping()[k]),
			Inline: true,
		})
	}

	return embed
}
