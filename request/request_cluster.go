package request

import (
	"context"
	"encoding/json"
	"r3/cluster"
	"r3/types"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func ClusterNodeDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (interface{}, error) {

	var req struct {
		Id uuid.UUID `json:"id"`
	}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, cluster.DelNode_tx(ctx, tx, req.Id)
}

func ClusterNodesGet_tx(ctx context.Context, tx pgx.Tx) (interface{}, error) {
	return cluster.GetNodes_tx(ctx, tx)
}

func ClusterNodeSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (interface{}, error) {

	var req struct {
		Id   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, cluster.SetNode_tx(ctx, tx, req.Id, req.Name)
}

func ClusterNodeShutdown_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (interface{}, error) {

	var req struct {
		Id uuid.UUID `json:"id"`
	}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, cluster.CreateEventForNodes_tx(ctx, tx, []uuid.UUID{req.Id},
		"shutdownTriggered", "{}", types.ClusterEventTarget{})
}
