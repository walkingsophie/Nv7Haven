package base

import "github.com/Nv7-Github/sevcord/v2"

func (b *Base) Error(ctx sevcord.Ctx, err error) {
	if err != nil {
		ctx.Acknowledge()
		ctx.Respond(sevcord.NewMessage("").AddEmbed(
			sevcord.NewEmbed().
				Title("Error").
				Description("```" + err.Error() + "```"),
		))
	}
}
