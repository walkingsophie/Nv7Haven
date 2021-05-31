package eod

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

func (b *EoD) giveCmd(elem string, giveTree bool, user string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv, exists := dat.invCache[user]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		rsp.Resp(fmt.Sprintf("Element **%s** doesn't exist!", elem))
		return
	}

	msg, suc := giveElem(dat.elemCache, giveTree, elem, &inv)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
		return
	}

	dat.invCache[user] = inv
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)

	rsp.Resp("Successfully gave element **" + el.Name + "**!")
}

func giveElem(elemCache map[string]element, giveTree bool, elem string, out *map[string]empty) (string, bool) {
	el, exists := elemCache[strings.ToLower(elem)]
	if !exists {
		return elem, false
	}
	if giveTree {
		for _, parent := range el.Parents {
			if len(strings.TrimSpace(parent)) == 0 {
				continue
			}
			_, exists := (*out)[strings.ToLower(parent)]
			if !exists {
				msg, suc := giveElem(elemCache, giveTree, parent, out)
				if !suc {
					return msg, false
				}
			}
		}
	}
	(*out)[strings.ToLower(el.Name)] = empty{}
	return "", true
}

func (b *EoD) calcTreeCmd(elem string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	txt, suc, msg := calcTree(dat.elemCache, elem)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
		return
	}
	if len(txt) <= 2000 {
		rsp.Message("Sent path in DMs!")
		rsp.DM(txt)
		return
	}
	rsp.Message("The path was too long! Sending it as a file in DMs!")

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := bytes.NewBufferString(txt)
	b.dg.ChannelFileSend(channel.ID, "path.txt", buf)
}

// Treecalc
func calcTree(elemCache map[string]element, elem string) (string, bool, string) {
	t := tree{
		text:      "",
		elemCache: elemCache,
		calced:    make(map[string]empty),
		num:       1,
	}
	suc, msg := t.addElem(elem)
	if len(t.text) > 2000 {
		return t.rawTxt, suc, msg
	}
	return t.text, suc, msg
}

type tree struct {
	text      string
	rawTxt    string
	elemCache map[string]element
	calced    map[string]empty
	num       int

	its int
}

func (t *tree) addElem(elem string) (bool, string) {
	t.its++
	if t.its == 500 {
		t.its = 0
		runtime.GC()
	}

	_, exists := t.calced[strings.ToLower(elem)]
	if !exists {
		el, exists := t.elemCache[strings.ToLower(elem)]
		if !exists {
			return false, elem
		}
		if len(el.Parents) == 1 {
			el.Parents = append(el.Parents, el.Parents[0])
		}
		for _, parent := range el.Parents {
			if len(strings.TrimSpace(parent)) == 0 {
				continue
			}
			suc, msg := t.addElem(parent)
			if !suc {
				return false, msg
			}
		}
		perf := "%d. "
		params := make([]interface{}, len(el.Parents))
		for i, val := range el.Parents {
			if i == 0 {
				perf += "%s"
			} else {
				perf += " + %s"
			}
			params[i] = interface{}(t.elemCache[strings.ToLower(val)].Name)
		}
		params = append([]interface{}{t.num}, params...)
		params = append(params, el.Name)
		if len(el.Parents) >= 2 {
			t.text += fmt.Sprintf(perf+" = **%s**\n", params...)
			t.rawTxt += fmt.Sprintf(perf+" = %s\n", params...)
			t.num++
		}
		t.calced[strings.ToLower(elem)] = empty{}
	}
	return true, ""
}
