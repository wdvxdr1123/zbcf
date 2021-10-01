package user

import (
	"fmt"
	"strconv"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/kv"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/wdvxdr1123/zbcf/util"
)

var db = kv.New("__cf_user__")

func init() {
	engine := zero.New()
	engine.OnCommand("绑定", zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		var cmd extension.CommandModel
		_ = ctx.Parse(&cmd)

		unquote, err := strconv.Unquote(cmd.Args)
		if err != nil {
			unquote = cmd.Args
		}

		id := []byte(strconv.FormatInt(ctx.Event.UserID, 10))
		_ = db.Put(id, []byte(unquote))

		user := GetUserInfo(ctx.Event.UserID)
		if user != nil {
			ctx.SendChain(
				message.At(ctx.Event.UserID),
				message.Text(fmt.Sprintf(`绑定用户成功!
用户名: %s
Rank: %s
Rating: %d`, user.Handle, user.Rank, user.Rating)),
			)
		} else {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败!"))
		}
	})
}

func GetUserInfo(user int64) *Info {
	handle, err := db.Get([]byte(strconv.FormatInt(user, 10)))
	if err != nil || handle == nil {
		return nil
	}
	var rsp InfoResponse
	err = util.NewCodeforcesAPI("/user.info").Params(util.H{
		"handles": string(handle),
	}).Decode(&rsp)
	if err != nil {
		return nil
	}
	for i, r := range rsp.Result {
		if r.Handle == string(handle) {
			return &rsp.Result[i]
		}
	}
	return nil
}
