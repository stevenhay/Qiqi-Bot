package role

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stevenhay/qiqi_bot/src/config"
	"github.com/stevenhay/qiqi_bot/src/log"
)

type DefaultRole struct {
	defaultRole *discordgo.Role
}

func NewDefaultRole() *DefaultRole {
	defaultRole := &DefaultRole{}
	defaultRole.load()
	return defaultRole
}

func (r *DefaultRole) load() error {
	r.defaultRole = &discordgo.Role{
		Name: "None",
		ID:   "",
	}

	id := config.Get().DefaultRole.ID
	name := config.Get().DefaultRole.Name

	//id, name, err := db.Get().GetDefaultRole()
	//if err != nil {
	//	return err
	//}

	r.defaultRole.Name = name
	r.defaultRole.ID = id
	return nil
}

func (r *DefaultRole) GetDefaultRole() *discordgo.Role {
	return r.defaultRole
}

func (r *DefaultRole) SetDefaultRole(role *discordgo.Role, author, authorID string) {
	r.defaultRole = role

	config.Get().DefaultRole = config.DefaultRole{
		Name: role.Name,
		ID:   role.ID,
	}
	config.Get().Save()
}

func (r *DefaultRole) ApplyDefaultRole(s *discordgo.Session, m *discordgo.Member) {
	AddUserRole(s, m, r.defaultRole)
	log.LogInfo(s, s.State.User.Username, "Assigned Role", "Assigned default role "+r.defaultRole.Name+" to "+m.User.Username)
}
