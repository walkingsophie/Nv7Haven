package discord

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type redditResp struct {
	Data redditData
}

type redditData struct {
	Children []previewData
	After    string
}

type previewData struct {
	Data meme
}

type meme struct {
	URL       string
	Title     string
	Permalink string
}

func (b *Bot) memes(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "meme") {
		if len(b.memedat) == 0 { // first time since startup
			s.ChannelMessageSend(m.ChannelID, "Sorry, this is the first time someone has asked for a meme since the server started. It may take a moment for us to download the memes from reddit.")
			success := b.loadMemes(m)
			if !success {
				b.dg.ChannelMessageSend(m.ChannelID, "Failed to lead memes")
			}
		}
		if (time.Now().Sub(b.memerefreshtime)).Hours() >= 1 { // its been an hour
			go b.loadMemes(m)
		}

		// send message
		unique := false
		var randnum int
		if len(b.memecache[m.GuildID]) == len(b.memedat) {
			b.memecache[m.GuildID] = make(map[int]empty, 0)
		}
		for !unique {
			randnum = rand.Intn(len(b.memedat))
			unique = true
			_, exists := b.memecache[m.GuildID]
			if !exists {
				b.memecache[m.GuildID] = make(map[int]empty, 0)
				unique = true
				break
			}
			_, unique = b.memecache[m.GuildID][randnum]
			unique = !unique
		}
		b.memecache[m.GuildID][randnum] = empty{}
		meme := b.memedat[randnum]
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			URL:   meme.Permalink,
			Type:  discordgo.EmbedTypeImage,
			Title: meme.Title,
			Image: &discordgo.MessageEmbedImage{
				URL: meme.URL,
			},
		})
		if b.handle(err, m) {
			return
		}
	}

	if strings.HasPrefix(m.Content, "cmeme") {
		if len(b.memedat) == 0 { // first time since startup
			s.ChannelMessageSend(m.ChannelID, "Sorry, this is the first time someone has asked for a meme since the server started. It may take a moment for us to download the memes from reddit.")
			success := b.loadMemes(m)
			if !success {
				b.dg.ChannelMessageSend(m.ChannelID, "Failed to lead memes")
			}
		}
		if (time.Now().Sub(b.memerefreshtime)).Hours() >= 1 { // its been an hour
			go b.loadMemes(m)
		}

		// send message
		unique := false
		var randnum int
		if len(b.cmemecache[m.GuildID]) == len(b.cmemedat) {
			b.cmemecache[m.GuildID] = make(map[int]empty, 0)
		}
		for !unique {
			randnum = rand.Intn(len(b.memedat))
			unique = true
			_, exists := b.cmemecache[m.GuildID]
			if !exists {
				b.cmemecache[m.GuildID] = make(map[int]empty, 0)
				unique = true
				break
			}
			_, unique = b.cmemecache[m.GuildID][randnum]
			unique = !unique
		}
		b.cmemecache[m.GuildID][randnum] = empty{}
		meme := b.cmemedat[randnum]
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			URL:   meme.Permalink,
			Type:  discordgo.EmbedTypeImage,
			Title: meme.Title,
			Image: &discordgo.MessageEmbedImage{
				URL: meme.URL,
			},
		})
		if b.handle(err, m) {
			return
		}
	}

	if strings.HasPrefix(m.Content, "pmeme") {
		if len(b.memedat) == 0 { // first time since startup
			s.ChannelMessageSend(m.ChannelID, "Sorry, this is the first time someone has asked for a meme since the server started. It may take a moment for us to download the memes from reddit.")
			success := b.loadMemes(m)
			if !success {
				b.dg.ChannelMessageSend(m.ChannelID, "Failed to lead memes")
			}
		}
		if (time.Now().Sub(b.memerefreshtime)).Hours() >= 1 { // its been an hour
			go b.loadMemes(m)
		}

		// send message
		unique := false
		var randnum int
		if len(b.pmemecache[m.GuildID]) == len(b.pmemedat) {
			b.pmemecache[m.GuildID] = make(map[int]empty, 0)
		}
		for !unique {
			randnum = rand.Intn(len(b.pmemedat))
			unique = true
			_, exists := b.pmemecache[m.GuildID]
			if !exists {
				b.pmemecache[m.GuildID] = make(map[int]empty, 0)
				unique = true
				break
			}
			_, unique = b.pmemecache[m.GuildID][randnum]
			unique = !unique
		}
		b.pmemecache[m.GuildID][randnum] = empty{}
		meme := b.pmemedat[randnum]
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			URL:   meme.Permalink,
			Type:  discordgo.EmbedTypeImage,
			Title: meme.Title,
			Image: &discordgo.MessageEmbedImage{
				URL: meme.URL,
			},
		})
		if b.handle(err, m) {
			return
		}
	}
}

func (b *Bot) loadMemes(m *discordgo.MessageCreate) bool {
	b.memerefreshtime = time.Now()
	b.memecache = make(map[string]map[int]empty, 0)
	b.cmemecache = make(map[string]map[int]empty, 0)
	b.pmemecache = make(map[string]map[int]empty, 0)
	var suc bool
	b.memedat, suc = b.downloadMeme(m, "memes")
	if !suc {
		return false
	}
	b.cmemedat, suc = b.downloadMeme(m, "cleanmemes")
	if !suc {
		return false
	}
	b.pmemedat, suc = b.downloadMeme(m, "ProgrammerHumor")
	if !suc {
		return false
	}
	return true
}

func (b *Bot) downloadMeme(m *discordgo.MessageCreate, subreddit string) ([]meme, bool) {
	children := make([]previewData, 0)
	after := ""
	for len(children) < 200 {
		// Download
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://reddit.com/r/"+subreddit+"/hot.json?after="+after, nil)
		if b.handle(err, m) {
			return nil, false
		}
		req.Header.Set("User-Agent", "Nv7 Bot")
		res, err := client.Do(req)
		if b.handle(err, m) {
			return nil, false
		}
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if b.handle(err, m) {
			return nil, false
		}

		// Process
		var dat redditResp
		err = json.Unmarshal(data, &dat)
		if b.handle(err, m) {
			return nil, false
		}
		children = append(children, dat.Data.Children...)
		after = dat.Data.After
	}

	memedat := make([]meme, len(children))
	for i, val := range children {
		memedat[i] = val.Data
		memedat[i].Permalink = "https://reddit.com" + memedat[i].Permalink
	}
	return memedat, true
}
