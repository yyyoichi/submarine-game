package handler

import (
	"context"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/yyyoichi/submarine-game/internal/core"
	v2 "github.com/yyyoichi/submarine-game/internal/gen/api/v2"
	"github.com/yyyoichi/submarine-game/internal/gen/api/v2/apiv2connect"
	"github.com/yyyoichi/submarine-game/internal/services/battle"
	"github.com/yyyoichi/submarine-game/internal/services/matching"
	"github.com/yyyoichi/submarine-game/internal/store"
)

type V2Handler struct {
	matchingService *matching.MatchingService
	battleService   *battle.BattleService
	apiv2connect.MatchingServiceHandler
	apiv2connect.BattleServiceHandler
}

func NewV2(s *store.Store) V2Handler {
	return V2Handler{
		matchingService: matching.New(s),
		battleService:   battle.New(s),
	}
}

func (h *V2Handler) Join(ctx context.Context, req *connect.Request[v2.JoinRequest]) (*connect.Response[v2.JoinResponse], error) {
	output, err := h.matchingService.Join()
	if err != nil {
		return nil, err
	}
	resp := &v2.JoinResponse{
		PlayerId: output.PlayerId,
		GameId:   output.GameId,
	}
	if output.Matched {
		err := h.battleService.NewGame(output.GameId, [2]string{output.EnemyId, output.PlayerId})
		if err != nil {
			return nil, err
		}
	}
	return &connect.Response[v2.JoinResponse]{
		Msg: resp,
	}, nil
}

func (h *V2Handler) WaitEnemy(ctx context.Context, req *connect.Request[v2.WaitEnemyRequest], stream *connect.ServerStream[v2.WaitEnemyResponse]) error {
	tick := time.NewTicker(time.Duration(1 * time.Second))
	defer tick.Stop()
	var resp = &v2.WaitEnemyResponse{
		PlayerId: req.Msg.PlayerId,
	}
	ch := h.matchingService.Wait(ctx, req.Msg.PlayerId)
	for {
		select {
		case <-ctx.Done():
		case <-tick.C:
			if err := stream.Send(resp); err != nil {
				return err
			}
		case gameIdEnemyId := <-ch:
			resp.GameId = gameIdEnemyId[0]
			if err := stream.Send(resp); err != nil {
				return err
			}
		}
	}
}

func (h *V2Handler) Logs(ctx context.Context, req *connect.Request[v2.LogsRequest]) (*connect.Response[v2.LogsResponse], error) {
	output, err := h.battleService.GetLogs(ctx, &battle.GetLogsInput{
		GameId:   req.Msg.GameId,
		PlayerId: req.Msg.PlayerId,
	})
	if err != nil {
		return nil, err
	}
	resp := &v2.LogsResponse{
		RequireAction:       output.RequireAction,
		RequireDeployAction: output.RequireDeployAction,
		GameIsOver:          output.GameOver != nil,
		NumTurn:             int32(output.NumTurn),
		Timeout:             output.Timeout.UnixMilli(),
		MillSecondPerTurn:   output.TimeoutDurationMSec,
		BoardWidth:          int32(output.Game.OceanMap.W),
		ActionLogs:          make([]*v2.LogsResponse_TurnAction, output.NumTurn),
		Sectors:             make([]*v2.LogsResponse_SectorStatus, output.Game.OceanMap.H*output.Game.OceanMap.W),
	}
	if output.GameOver != nil {
		resp.Win = output.GameOver.Winner == req.Msg.PlayerId
		switch output.GameOver.Reason {
		case core.TorpedoHit:
			resp.GameOverReason = v2.GameOverReason_GAME_OVER_REASON_TORPEDO_HIT
		case core.MineHit:
			resp.GameOverReason = v2.GameOverReason_GAME_OVER_REASON_MINE_HIT
		case core.Timeout:
			resp.GameOverReason = v2.GameOverReason_GAME_OVER_REASON_TIMEOUT
		default:
			resp.GameOverReason = v2.GameOverReason_GAME_OVER_REASON_UNSPECIFIED
		}
	}
	for _, src := range output.Actions {
		dist := &v2.LogsResponse_Action{
			At:        int32(src.At),
			Form:      int32(src.From),
			To:        int32(src.To),
			Direction: int32(src.Direction),
			Turn:      int32(src.Turn),
			Me:        src.PlayerId == req.Msg.PlayerId,
		}
		switch src.ActionResult {
		case core.FullSpeedAhead:
			dist.Result = v2.ActionResult_ACTION_RESULT_FULL_SPEED_AHEAD
		case core.HardToStarboard:
			dist.Result = v2.ActionResult_ACTION_RESULT_HARD_TO_STARBOARD
		case core.Hit:
			dist.Result = v2.ActionResult_ACTION_RESULT_HIT
		default:
			dist.Result = v2.ActionResult_ACTION_RESULT_UNSPECIFIED
		}
		switch src.T {
		case core.MoveAction:
			dist.Type = v2.ActionType_ACTION_TYPE_MOVE
		case core.TorpedoFireAction:
			dist.Type = v2.ActionType_ACTION_TYPE_FIIRE_TORPEDO
		case core.MineTriggerAction:
			dist.Type = v2.ActionType_ACTION_TYPE_TRIGGER_MINE
		default:
			dist.Type = 0
		}
		if resp.ActionLogs[src.Turn] == nil {
			resp.ActionLogs[src.Turn] = &v2.LogsResponse_TurnAction{}
		}
		if dist.Me {
			resp.ActionLogs[src.Turn].Me = dist
		} else {
			resp.ActionLogs[src.Turn].Enemy = dist
		}
	}
	for _, sector := range output.Game.OceanMap.Sectors() {
		resp.Sectors[sector] = &v2.LogsResponse_SectorStatus{
			Sector:        int32(sector),
			Island:        slices.Contains(output.Game.Islands, sector),
			SelfOccupied:  output.Prev.At == sector,
			EnableActions: make([]v2.ActionType, 0, len(output.SectorActionsMap[sector])),
		}
		for _, actionType := range output.SectorActionsMap[sector] {
			var dist v2.ActionType
			switch actionType {
			case core.MoveAction:
				dist = v2.ActionType_ACTION_TYPE_MOVE
			case core.TorpedoFireAction:
				dist = v2.ActionType_ACTION_TYPE_FIIRE_TORPEDO
			case core.MineTriggerAction:
				dist = v2.ActionType_ACTION_TYPE_TRIGGER_MINE
			}
			if dist != 0 {
				resp.Sectors[sector].EnableActions = append(resp.Sectors[sector].EnableActions, dist)
			}
		}
	}
	return &connect.Response[v2.LogsResponse]{
		Msg: resp,
	}, nil
}

// 初回の行動する
func (h *V2Handler) Deploy(ctx context.Context, req *connect.Request[v2.DeployRequest]) (*connect.Response[v2.DeployResponse], error) {
	input := &battle.DeploySubmarineAndMinesInput{
		GameId:   req.Msg.GameId,
		PlayerId: req.Msg.PlayerId,
		At:       int8(req.Msg.At),
		Mines:    make([]int8, len(req.Msg.Mines)),
	}
	for i, s := range req.Msg.Mines {
		input.Mines[i] = int8(s)
	}
	err := h.battleService.DeploySubmarineAndMines(ctx, input)
	if err != nil {
		return nil, err
	}
	return &connect.Response[v2.DeployResponse]{}, nil
}

// 行動する
func (h *V2Handler) Action(ctx context.Context, req *connect.Request[v2.ActionRequest]) (*connect.Response[v2.ActionResponse], error) {
	input := &battle.ActionInput{
		GameId:   req.Msg.GameId,
		PlayerId: req.Msg.PlayerId,
		At:       int8(req.Msg.At),
	}
	var err error
	switch req.Msg.Type {
	case v2.ActionType_ACTION_TYPE_MOVE:
		err = h.battleService.Move(ctx, input)
	case v2.ActionType_ACTION_TYPE_FIIRE_TORPEDO:
		err = h.battleService.FireTorpedo(ctx, input)
	case v2.ActionType_ACTION_TYPE_TRIGGER_MINE:
		err = h.battleService.TriggerMine(ctx, input)
	}
	if err != nil {
		return nil, err
	}
	return &connect.Response[v2.ActionResponse]{}, nil
}

// 相手の行動を待機する
func (h *V2Handler) Wait(ctx context.Context, req *connect.Request[v2.WaitRequest], stream *connect.ServerStream[v2.WaitResponse]) error {
	done, err := h.battleService.WaitTurn(ctx, &battle.WaitTurnInput{GameId: req.Msg.GameId, PlayerId: req.Msg.PlayerId})
	if err != nil {
		return err
	}
	resp := &v2.WaitResponse{}
	tick := time.NewTicker(time.Duration(1 * time.Second))
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
		case <-tick.C:
			if err := stream.Send(resp); err != nil {
				return err
			}
		case <-done:
			return nil
		}
	}
}
