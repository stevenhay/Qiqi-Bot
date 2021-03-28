package role

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stevenhay/qiqi_bot/src/config"
	"github.com/stevenhay/qiqi_bot/src/helpers"
	"github.com/stevenhay/qiqi_bot/src/io"
	"github.com/stevenhay/qiqi_bot/src/log"
)

type ColourRoles struct {
	MessageID      string            `json:"message_id"`
	EmojiToRole    map[string]string `json:"emoji_to_role"`
	reactionClicks map[string]*reactionClick
}

type reactionClick struct {
	timer            *time.Timer
	lastEmojiClicked discordgo.Emoji
}

func NewColourRoles() *ColourRoles {
	colourRoles := &ColourRoles{}
	colourRoles.load()
	return colourRoles
}

func (c *ColourRoles) load() {
	/*messageID, err := db.Get().GetColourMessageID()
	if err != nil {
		return err
	}*/

	mappingFile, err := os.Open("../config/colour_role_picker.json")
	if err != nil {
		return
	}
	defer mappingFile.Close()

	byteVal, err := ioutil.ReadAll(mappingFile)
	if err != nil {
		return
	}

	var res ColourRoles
	json.Unmarshal([]byte(byteVal), &res)

	// c.MessageID =
	// c.EmojiToRole = make(map[string]string)
	// c.reactionClicks = make(map[string]*reactionClick)
	// for k, v := range (res["emoji_to_role"]).(map[string]interface{}) {
	// 	c.EmojiToRole[k] = v.(string)
	// }

	fmt.Println(res)
	c.MessageID = res.MessageID
	c.EmojiToRole = res.EmojiToRole
	c.reactionClicks = make(map[string]*reactionClick)
}

func (c *ColourRoles) GetMessageID() string {
	return c.MessageID
}

func (c *ColourRoles) SetMessageID(messageID string) error {
	c.MessageID = messageID
	return c.save() //db.Get().SetColourMessageID(messageID)
}

func (c *ColourRoles) save() error {
	return io.SaveConfig("colour_role_picker", c)
}

func (c *ColourRoles) GetEmojiToRoleMapping() map[string]string {
	return c.EmojiToRole
}

func (c *ColourRoles) HandleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if _, ok := c.EmojiToRole[r.Emoji.Name]; !ok {
		go s.MessageReactionsRemoveEmoji(r.ChannelID, r.MessageID, r.Emoji.Name)
		go s.MessageReactionsRemoveEmoji(r.ChannelID, r.MessageID, helpers.EmojiToID(r.Emoji.Name, r.Emoji.ID))
		return
	}

	if reaction, ok := c.reactionClicks[r.UserID]; ok {
		reaction.timer.Reset(1 * time.Second)
		reaction.lastEmojiClicked = r.Emoji
	} else {
		c.reactionClicks[r.UserID] = &reactionClick{
			timer:            time.NewTimer(1 * time.Second),
			lastEmojiClicked: r.Emoji,
		}
		go c.handleReactionAdd(s, r)
	}
}

func (c *ColourRoles) HandleReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if _, ok := c.EmojiToRole[r.Emoji.Name]; !ok {
		return
	}

	role, err := helpers.GetRoleForID(s, c.EmojiToRole[r.Emoji.Name])
	if err != nil {
		log.LogError(s, s.State.User.Username, "ReactionRemove: failed to get role from role id", err.Error())
		return
	}

	member, err := s.GuildMember(config.Get().GuildID, r.UserID)
	if err != nil {
		log.LogError(s, s.State.User.Username, "ReactionRemove: failed to get member from user id", err.Error())
		return
	}

	RemoveUserRole(s, member, role)
}

func (c *ColourRoles) handleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	click := c.reactionClicks[r.UserID]

	<-click.timer.C
	for e := range c.EmojiToRole {
		if e != click.lastEmojiClicked.Name {
			go s.MessageReactionRemove(r.ChannelID, r.MessageID, e, r.UserID)
		}
	}

	member, err := s.GuildMember(config.Get().GuildID, r.UserID)
	if err != nil {
		log.LogError(s, s.State.User.Username, "ReactionAdd: failed to get member from user id", err.Error())
		return
	}

	role, err := helpers.GetRoleForID(s, c.EmojiToRole[click.lastEmojiClicked.Name])
	if err != nil {
		log.LogError(s, s.State.User.Username, "ReactionAdd: failed to get role from role id", err.Error())
		return
	}

	AddUserRole(s, member, role)
	delete(c.reactionClicks, r.UserID)
}
