package helpers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/stevenhay/qiqi_bot/src/config"
	"github.com/stevenhay/qiqi_bot/src/log"
)

func MemberHasAllPermissions(s *discordgo.Session, m *discordgo.Member, permissions []int64) bool {
	allRoles, err := s.GuildRoles(config.Get().GuildID)
	if err != nil {
		log.LogError(s, s.State.User.Username, "Permissions: failed to get roles", err.Error())
		return false
	}

	roleMap := make(map[string]*discordgo.Role)
	for _, role := range allRoles {
		roleMap[role.ID] = role
	}

	okPerms := make([]int64, 0, len(permissions))
	for _, roleId := range m.Roles {
		role, ok := roleMap[roleId]
		if !ok {
			log.LogWarn(s, s.State.User.Username, "Permissions: failed to find role", fmt.Sprintf("role %s doesn't seem to exist on this server", role.Name))
			continue
		}

		for _, permissionId := range permissions {
			if (role.Permissions & permissionId) == permissionId {
				okPerms = append(okPerms, permissionId)
			}
		}

		if len(okPerms) == len(permissions) {
			return true
		}
	}

	return false
}

func MemberHasAnyPermissions(s *discordgo.Session, m *discordgo.Member, permissions []int64) bool {
	allRoles, err := s.GuildRoles(config.Get().GuildID)
	if err != nil {
		log.LogError(s, s.State.User.Username, "Permissions: failed to get roles", err.Error())
		return false
	}

	roleMap := make(map[string]*discordgo.Role)
	for _, role := range allRoles {
		roleMap[role.ID] = role
	}

	for _, roleId := range m.Roles {
		role, ok := roleMap[roleId]
		if !ok {
			log.LogWarn(s, s.State.User.Username, "Permissions: failed to find role", fmt.Sprintf("role %s doesn't seem to exist on this server", role.Name))
			continue
		}

		for _, permissionId := range permissions {
			if (role.Permissions & permissionId) == permissionId {
				return true
			}
		}
	}

	return false
}
