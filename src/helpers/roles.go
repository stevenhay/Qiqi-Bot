package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/stevenhay/qiqi_bot/src/config"
	"github.com/stevenhay/qiqi_bot/src/log"
)

func GetRoleForName(s *discordgo.Session, roleName string) (*discordgo.Role, error) {
	allRoles, err := s.GuildRoles(config.Get().GuildID)
	if err != nil {
		log.LogError(s, s.State.User.ID, "'GetRoleForName' failed!", err.Error())
		return nil, err
	}

	for _, r := range allRoles {
		if strings.EqualFold(r.Name, roleName) {
			return r, nil
		}
	}

	return nil, fmt.Errorf("GetRoleForName: role not found")
}

func GetRoleForID(s *discordgo.Session, roleID string) (*discordgo.Role, error) {
	allRoles, err := s.GuildRoles(config.Get().GuildID)
	if err != nil {
		log.LogError(s, s.State.User.ID, "'GetRoleForID' failed!", err.Error())
		return nil, err
	}

	for _, r := range allRoles {
		if r.ID == roleID {
			return r, nil
		}
	}

	return nil, fmt.Errorf("GetRoleForID: role not found")
}

func RoleAtToID(roleAtString string) string {
	idx := strings.IndexRune(roleAtString, '&')
	if idx == -1 {
		logrus.WithTime(time.Now()).Warn("Role `" + roleAtString + "` did not contain an '&' symbol or wasn't a role @")
		return ""
	}

	return roleAtString[idx+1 : len(roleAtString)-1]
}

func RoleIDToAt(roleID string) string {
	return fmt.Sprintf("<@&%s>", roleID)
}
