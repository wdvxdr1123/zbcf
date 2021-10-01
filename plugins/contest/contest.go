package contest

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/kv"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/wdvxdr1123/zbcf/util"
)

type Contest struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Phase               string `json:"phase"`
	Frozen              bool   `json:"frozen"`
	DurationSeconds     int    `json:"durationSeconds"`
	StartTimeSeconds    int    `json:"startTimeSeconds"`
	RelativeTimeSeconds int    `json:"relativeTimeSeconds"`
}

type List struct {
	Status string    `json:"status"`
	Result []Contest `json:"result"`
}

var (
	db            = kv.New("__cf_contest__")
	engine        = zero.New()
	reminderDBKey = []byte("reminder")
)

func upcoming() []Contest {
	var l List
	err := util.NewCodeforcesAPI(`/contest.list`).Params(util.H{
		"gym": "false",
	}).Decode(&l)
	if err != nil {
		return nil
	}
	const sevendaysago = -int(time.Hour * 24 * 7 / time.Second)
	var validContest []Contest
	for _, r := range l.Result {
		// Skip all finished Contest
		if r.Phase == "FINISHED" {
			continue
		}
		if r.RelativeTimeSeconds >= sevendaysago {
			validContest = append(validContest, r)
		}
	}
	sort.Slice(validContest, func(i, j int) bool {
		return validContest[i].StartTimeSeconds < validContest[j].StartTimeSeconds
	})
	return validContest
}

func init() {
	engine.OnCommandGroup([]string{
		`最近赛事`,
		`upcoming Contest`,
	}, zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		result := upcoming()
		if result == nil {
			ctx.Send("获取近期比赛失败!")
			return
		}

		var msg message.Message
		self := ctx.Event.SelfID
		for _, r := range result {
			startTime := time.Unix(int64(r.StartTimeSeconds), 0).Format(`2006-01-02 15:04:05`)
			var format = fmt.Sprintf(`比赛ID: %d
比赛名: %s
开始时间: %s
比赛时长: %.2f Hours`, r.Id, r.Name, startTime, float64(r.DurationSeconds)/3600.0)

			msg = append(msg, message.CustomNode("CodeforcesBot", self, format))
		}

		ctx.SendGroupForwardMessage(ctx.Event.GroupID, msg)
	})

	go func() {
		b, _ := db.Get(reminderDBKey)
		_ = gob.NewDecoder(bytes.NewReader(b)).Decode(&reminderList)
		t := time.NewTicker(time.Minute)
		for {
			<-t.C
			reminderMutex.Lock()
			j := 0
			now := time.Now()
			for i, r := range reminderList {
				st := time.Unix(int64(r.Contest.StartTimeSeconds), 0)
				if st.Sub(now) < time.Hour*30 {
					broadcast(r)
					continue
				}
				if i != j {
					reminderList[j] = r
				}
				j++
			}
			if j != len(reminderList) {
				reminderList = reminderList[:j]
				saveReminder()
			}
			reminderMutex.Unlock()
		}
	}()

	engine.OnRegex(`提醒(我|\[CQ:at,qq=[0-9]+\])参加比赛([0-9]+)`, zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
		var cmd extension.RegexModel
		_ = ctx.Parse(&cmd)
		var u int64
		if cmd.Matched[1] == "我" {
			u = ctx.Event.UserID
		} else {
			at := cmd.Matched[1]
			at = strings.TrimPrefix(at, "[CQ:at,qq=")
			at = strings.TrimSuffix(at, "]")
			u, _ = strconv.ParseInt(at, 10, 64)
		}
		cid, _ := strconv.ParseInt(cmd.Matched[2], 10, 64)
		gid := ctx.Event.GroupID
		var rm *reminder
		reminderMutex.Lock()
		defer reminderMutex.Unlock()
		// check remider list
		for _, r := range reminderList {
			if r.Group == gid && int64(r.Contest.Id) == cid {
				rm = r
				goto ok
			}
		}
		// check upcoming Contest
		for _, c := range upcoming() {
			fmt.Println(c)
			if int64(c.Id) == cid {
				rm = &reminder{
					Contest: c,
					Group:   gid,
					User:    []int64{},
				}
				reminderList = append(reminderList, rm)
				goto ok
			}
		}
		// no such Contest
		ctx.Send("无法获取比赛信息!")
		return
	ok:
		rm.User = append(rm.User, u)
		saveReminder()
		ctx.SendChain(message.At(u), message.Text("添加成功!"))
	})
}

func broadcast(r *reminder) {
	defer func() { recover() }()
	var msg message.Message
	msg = append(msg, message.Text(fmt.Sprintf(`Codeforces小助手提醒您:
您预定的比赛 %s 即将开始!
`, r.Contest.Name)))
	for _, u := range r.User {
		msg = append(msg, message.At(u))
	}
	zero.RangeBot(func(_ int64, ctx *zero.Ctx) bool {
		ctx.SendGroupMessage(r.Group, msg)
		return true
	})
}
