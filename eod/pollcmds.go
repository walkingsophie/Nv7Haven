package eod

import (
	"fmt"
	"strings"
)

func (b *EoD) suggestCmd(suggestion string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	if dat.combCache == nil {
		dat.combCache = make(map[string]comb)
	}
	comb, exists := dat.combCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You haven't combined anything!")
		return
	}

	err := b.createPoll(poll{
		Channel:   dat.votingChannel,
		Guild:     m.GuildID,
		Kind:      pollCombo,
		Value1:    comb.elem1,
		Value2:    comb.elem2,
		Value3:    suggestion,
		Value4:    m.Author.ID,
		Data:      make(map[string]interface{}),
		Upvotes:   0,
		Downvotes: 0,
	})
	if rsp.Error(err) {
		return
	}
	rsp.Resp("Suggested " + comb.elem1 + " + " + comb.elem2 + " = " + suggestion + " ✨")
}

func (b *EoD) markCmd(elem string, mark string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem))
		return
	}

	if el.Creator == m.Author.ID {
		b.mark(m.GuildID, elem, mark, "")
		rsp.Resp(fmt.Sprintf("You have signed **%s**! 🖋️", el.Name))
		return
	}

	err := b.createPoll(poll{
		Channel: dat.votingChannel,
		Guild:   m.GuildID,
		Kind:    pollSign,
		Value1:  el.Name,
		Value2:  mark,
		Value3:  el.Comment,
		Value4:  m.Author.ID,
	})
	if rsp.Error(err) {
		return
	}
	rsp.Resp(fmt.Sprintf("Suggested a note for **%s** 🖊️", el.Name))
}

func (b *EoD) imageCmd(elem string, image string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem))
		return
	}

	if el.Creator == m.Author.ID {
		b.image(m.GuildID, elem, image, "")
		rsp.Resp(fmt.Sprintf("You added an image to **%s**! 📷", el.Name))
		return
	}

	err := b.createPoll(poll{
		Channel: dat.votingChannel,
		Guild:   m.GuildID,
		Kind:    pollImage,
		Value1:  el.Name,
		Value2:  image,
		Value3:  el.Image,
		Value4:  m.Author.ID,
	})
	if rsp.Error(err) {
		return
	}
	rsp.Resp(fmt.Sprintf("Suggested an image for **%s** 📷", el.Name))
}