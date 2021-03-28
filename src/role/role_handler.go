package role

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stevenhay/qiqi_bot/src/config"
	"github.com/stevenhay/qiqi_bot/src/log"
)

type RoleHandler struct {
	DefaultRole *DefaultRole
	ColourRoles *ColourRoles
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{NewDefaultRole(), NewColourRoles()}
}

func AddUserRole(s *discordgo.Session, member *discordgo.Member, role *discordgo.Role) {
	err := s.GuildMemberRoleAdd(config.Get().GuildID, member.User.ID, role.ID)
	if err != nil {
		log.LogError(s, s.State.User.Username, "Failed to assign role", err.Error())
		return
	}
}

func RemoveUserRole(s *discordgo.Session, member *discordgo.Member, role *discordgo.Role) {
	err := s.GuildMemberRoleRemove(config.Get().GuildID, member.User.ID, role.ID)
	if err != nil {
		log.LogError(s, s.State.User.Username, "Failed to renove role", err.Error())
		return
	}
}
