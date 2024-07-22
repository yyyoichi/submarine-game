package game

import (
	"fmt"
	"sync"
	"time"

	apiv1 "github.com/yyyoichi/submarine-game/internal/gen/api/v1"
)

type Game struct {
	Id        string
	Users     [2]string
	Island    [2]int64
	createdAt time.Time

	mu        sync.RWMutex
	NextUser  string
	Winner    string
	histories []history
}

type history struct {
	user string
	camp uint32
	t    apiv1.ActionType
	at   time.Time
}

func (g *Game) GetHistory(me string) *apiv1.HistoryResponse {
	g.mu.RLock()
	defer g.mu.RUnlock()

	resp := &apiv1.HistoryResponse{
		Island:    g.Island[:],
		MyTurn:    me == g.NextUser,
		Winner:    g.Winner,
		Histories: make([]*apiv1.History, len(g.histories)),
		Timeout:   g.getTimeout().UnixMilli(),
	}
	switch {
	case g.Winner != "":
		if me == g.Winner {
			resp.Description = "🏆Win!!"
		} else {
			resp.Description = "💣Lose.."
		}
	case resp.MyTurn:
		resp.Description = "🪖My turn"
	default:
		resp.Description = "👀Enemy turn"
	}

	for i, hist := range g.histories {
		switch hist.user {
		// 操作ユーザの履歴
		case me:
			resp.Histories[i] = &apiv1.History{
				UserId: hist.user,
				Turn:   int64(i),
				Camp:   hist.camp,
				Type:   hist.t,
			}
			if i == 0 || i == 1 {
				resp.Histories[i].Description = fmt.Sprintf("📍Placed in '%d'", hist.camp)
				continue
			}
			if hist.camp == uint32(apiv1.ActionType_ACTION_TYPE_MOVE) {
				resp.Histories[i].Description = fmt.Sprintf("🌊Moved to '%d'", hist.camp)
				continue
			}
			if hist.camp == uint32(apiv1.ActionType_ACTION_TYPE_LEAVE) {
				resp.Histories[i].Description = "⚠️You escaped"
				continue
			}
			resp.Histories[i].Description = fmt.Sprintf("💣Torpedo fired on '%d'", hist.camp)
			continue

		// 対戦相手の履歴
		default:
			mask := hist.mask()
			resp.Histories[i] = &apiv1.History{
				UserId: mask.user,
				Turn:   int64(i),
				Camp:   mask.camp,
				Type:   mask.t,
			}
			if i == 0 || i == 1 {
				resp.Histories[i].Description = "📍Placed in '?'"
				continue
			}
			if hist.camp == uint32(apiv1.ActionType_ACTION_TYPE_MOVE) {
				b := g.histories[i-2].camp
				a := hist.camp
				resp.Histories[i].Description = fmt.Sprintf("🌊Moved '%s'", calcDirection(b, a))
				continue
			}
			if hist.camp == uint32(apiv1.ActionType_ACTION_TYPE_LEAVE) {
				resp.Histories[i].Description = "✨Enemy escaped"
				continue
			}
			resp.Histories[i].Description = fmt.Sprintf("💣Torpedo fired on '%d'", hist.camp)
		}
	}
	return resp
}

func (g *Game) Action(user string, camp uint32, action apiv1.ActionType) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.Winner != "" {
		return ErrGameIsOver
	}
	if g.NextUser != user {
		return ErrIsnotYourTurn
	}
	if g.getTimeout().Add(time.Duration(500) * time.Millisecond).Before(time.Now()) {
		// timeoutよりも500milsec大きいときにタイムアウト判定
		g.leave(user)
		return ErrTimeout
	}
	if camp > campSize {
		return ErrOutOfCampSize
	}

	if len(g.histories) < 2 {
		// place
		g.histories = append(g.histories, history{
			user: user,
			camp: camp,
			t:    apiv1.ActionType_ACTION_TYPE_PLACE,
			at:   time.Now(),
		})
		g.changeTurn(user)
		return nil
	}

	switch action {
	case apiv1.ActionType_ACTION_TYPE_MOVE, apiv1.ActionType_ACTION_TYPE_BOMB:
		if d := calcDirection(g.histories[len(g.histories)-2].camp, camp); d == "" {
			return fmt.Errorf("%w: '%d'", ErrInvalidCamp, camp)
		}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidAction, action)
	}

	g.histories = append(g.histories, history{
		user: user,
		camp: camp,
		t:    action,
		at:   time.Now(),
	})
	g.changeTurn(user)
	return nil
}

func (g *Game) Leave(user string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.leave(user)
}

func (g *Game) leave(user string) {
	g.histories = append(g.histories, history{
		user: user,
		camp: 0,
		t:    apiv1.ActionType_ACTION_TYPE_LEAVE,
		at:   time.Now(),
	})
	if g.Users[0] == user {
		g.Winner = g.Users[1]
	} else {
		g.Winner = g.Users[0]
	}
}

func (g *Game) changeTurn(user string) {
	if g.Users[0] == user {
		g.NextUser = g.Users[1]
	} else {
		g.NextUser = g.Users[0]
	}
}

func (g *Game) getTimeout() time.Time {
	if len(g.histories) == 0 {
		return g.createdAt.Add(time.Duration(30) * time.Second)
	}
	return g.histories[len(g.histories)-1].at.Add(time.Duration(30) * time.Second)
}

func (h history) mask() history {
	h.user = ""
	if h.t == apiv1.ActionType_ACTION_TYPE_MOVE {
		// 移動の場合場所はマスクする
		h.camp = 0
	}
	return h
}

func calcDirection(b, a uint32) string {
	switch {
	case b-a == lineSize:
		return "north"
	case a-b == lineSize:
		return "south"
	case b-a == 1:
		return "west"
	case a-b == 1:
		return "east"
	}
	return ""
}
