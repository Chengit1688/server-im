package rpc

import (
	"context"
	"im/config"
	"im/pkg/common"
	"im/pkg/logger"
	"im/pkg/pb"
	"im/pkg/util"

	"google.golang.org/grpc"
)

var (
	connectClient pb.ConnectClient
)

func GetConnectClient() pb.ConnectClient {
	if connectClient == nil {
		initConnectClient()
	}
	return connectClient
}

func initConnectClient() {
	conn, err := grpc.DialContext(context.TODO(), config.Config.Server.ConnectRPCAddr, grpc.WithInsecure())
	if err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), err)
		return
	}
	connectClient = pb.NewConnectClient(conn)
}

func SendMessageToConnectServer(operationID string, mt common.MessageType, message interface{}, userIDs ...string) {
	data, err := util.JsonMarshal(message)
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "json marshal error, err:", err)
		return
	}

	var msg common.Message
	msg.Type = mt
	msg.Data = data

	var req pb.Message
	req.OperationID = operationID
	req.UserIDList = userIDs
	req.Data, _ = util.JsonMarshal(&msg)

	logger.Sugar.Infow(operationID, "func", util.GetSelfFuncName(), "msg", req)
	_, _ = GetConnectClient().SendMessage(context.Background(), &req)
}

func BroadcastMessageToConnectServer(operationID string, mt common.MessageType, message interface{}) {
	data, err := util.JsonMarshal(message)
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "json marshal error, err:", err)
		return
	}

	var msg common.Message
	msg.Type = mt
	msg.Data = data

	var req pb.Message
	req.OperationID = operationID
	req.Data, _ = util.JsonMarshal(&msg)

	logger.Sugar.Infow(operationID, "func", util.GetSelfFuncName(), "msg", req)
	_, _ = GetConnectClient().SendMessage(context.Background(), &req)
}
