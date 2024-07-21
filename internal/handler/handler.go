package handler

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/yyyoichi/submarine-game/internal/game"
	apiv1 "github.com/yyyoichi/submarine-game/internal/gen/api/v1"
	"github.com/yyyoichi/submarine-game/internal/gen/api/v1/apiv1connect"
)

type Handler struct {
	pg *game.Playground
	mt *game.Matching
	apiv1connect.GameServiceHandler
}

func NewHandler(ctx context.Context) *Handler {
	h := &Handler{
		pg: game.NewPlayground(),
		mt: game.NewMatching(),
	}
	go h.pg.DeleteRunner(ctx)
	return h
}

// 対戦する
func (h *Handler) Join(ctx context.Context, req *connect.Request[apiv1.JoinRequest], stream *connect.ServerStream[apiv1.JoinResponse]) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Minute)
	defer cancel()

	me, users := h.mt.Join()
	if users != nil {
		gm := h.pg.NewGame(*users)
		resp := &apiv1.JoinResponse{GameId: gm.Id, UserId: me}
		if err := stream.Send(resp); err != nil {
			return err
		}
		return nil
	}
	for {
		select {
		case <-ctx.Done():
			h.mt.Leave(me)
			return context.Cause(ctx)

		case <-time.After(time.Duration(500) * time.Millisecond):
			resp := &apiv1.JoinResponse{}
			if gm := h.pg.Found(me); gm != nil {
				resp.GameId = gm.Id
				resp.UserId = me
			}
			if err := stream.Send(resp); err != nil {
				return err
			}
			if resp.GameId != "" {
				return nil
			}
		}
	}
}

// 対戦から離れる
func (h *Handler) Leave(ctx context.Context, req *connect.Request[apiv1.LeaveRequest]) (*connect.Response[apiv1.LeaveResponse], error) {
	gm, _ := h.pg.Use(req.Msg.GameId)
	if gm != nil {
		gm.Leave(req.Msg.UserId)
	}
	h.mt.Leave(req.Msg.UserId)
	return connect.NewResponse(&apiv1.LeaveResponse{}), nil
}

// 行動履歴を取得する
func (h *Handler) History(ctx context.Context, req *connect.Request[apiv1.HistoryRequest]) (*connect.Response[apiv1.HistoryResponse], error) {
	gm, err := h.pg.Use(req.Msg.GameId)
	if err != nil {
		return nil, err
	}
	resp := gm.GetHistory(req.Msg.UserId)
	return connect.NewResponse(resp), nil
}

// 行動する
func (h *Handler) Action(ctx context.Context, req *connect.Request[apiv1.ActionRequest]) (*connect.Response[apiv1.ActionResponse], error) {
	gm, err := h.pg.Use(req.Msg.GameId)
	if err != nil {
		return nil, err
	}
	err = gm.Action(req.Msg.UserId, req.Msg.Camp, req.Msg.Type)
	return connect.NewResponse(&apiv1.ActionResponse{}), err
}

// 相手の行動を待機する
func (h *Handler) Wait(ctx context.Context, req *connect.Request[apiv1.WaitRequest], stream *connect.ServerStream[apiv1.WaitResponse]) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(31)*time.Second)
	defer cancel()

	gm, err := h.pg.Use(req.Msg.GameId)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
		case <-time.After(time.Duration(500) * time.Millisecond):
		}
		done := gm.NextUser == req.Msg.UserId
		resp := &apiv1.WaitResponse{
			Done: done,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}
