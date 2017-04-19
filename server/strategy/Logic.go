package strategy

import (
	"github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Migration/server/model"
)

func MigrationNeeded(predictionData *protocol.PredictionData, server *model.Server) (bool, *protocol.CheckpointRequest) {
	// complete this function with the migration logic
	return false, &protocol.CheckpointRequest{}
}
