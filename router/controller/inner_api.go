package controller

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/zzzgydi/zbyai/model"
	"github.com/zzzgydi/zbyai/router/utils"
	"github.com/zzzgydi/zbyai/service/thread"
)

func PostAuth(c *gin.Context) {
	ReturnSuccess(c, nil)
}

func PostThreadCreate(c *gin.Context) {
	req := &PostThreadCreateRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		ReturnBadRequest(c, err)
		return
	}

	if req.Query == "" {
		ReturnBadRequest(c, fmt.Errorf("query is empty"))
		return
	}

	user := utils.GetUser(c)
	if user == nil {
		ReturnBadRequest(c, fmt.Errorf("unauthorized"))
		return
	}

	var logger *slog.Logger
	trace := utils.GetTraceLogger(c)
	if trace != nil {
		logger = trace.Logger
		trace.SetBizRequest(req)
	}

	thread_, threadRun, err := thread.CreateThread(user.Id, req.Query, logger)
	if err != nil {
		ReturnServerError(c, err)
		return
	}

	ret := map[string]any{
		"id":    thread_.Id,
		"runId": threadRun.Id,
	}
	if trace != nil {
		trace.SetBizResponse(ret)
	}
	ReturnSuccess(c, ret)
}

func PostThreadAppend(c *gin.Context) {
	req := &PostThreadAppendRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		ReturnBadRequest(c, err)
		return
	}

	if req.Query == "" {
		ReturnBadRequest(c, fmt.Errorf("query is empty"))
		return
	}

	user := utils.GetUser(c)
	if user == nil {
		ReturnBadRequest(c, fmt.Errorf("unauthorized"))
		return
	}

	var logger *slog.Logger
	trace := utils.GetTraceLogger(c)
	if trace != nil {
		logger = trace.Logger
		trace.SetBizRequest(req)
	}

	thread_, threadRun, err := thread.AppendThread(user.Id, req.Id, req.Query, logger)
	if err != nil {
		ReturnServerError(c, err)
		return
	}

	ret := map[string]any{
		"id":    thread_.Id,
		"runId": threadRun.Id,
	}
	if trace != nil {
		trace.SetBizResponse(ret)
	}
	ReturnSuccess(c, ret)
}

func PostThreadRewrite(c *gin.Context) {
	req := &PostThreadRewriteRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		ReturnBadRequest(c, err)
		return
	}

	if req.Id == "" {
		ReturnBadRequest(c, fmt.Errorf("id is required"))
		return
	}

	if req.RunId == 0 {
		ReturnBadRequest(c, fmt.Errorf("run id is required"))
		return
	}

	user := utils.GetUser(c)
	if user == nil {
		ReturnBadRequest(c, fmt.Errorf("unauthorized"))
		return
	}

	var logger *slog.Logger
	trace := utils.GetTraceLogger(c)
	if trace != nil {
		logger = trace.Logger
		trace.SetBizRequest(req)
	}

	thread_, threadRun, err := thread.RewriteThread(user.Id, req.Id, req.RunId, logger)
	if err != nil {
		ReturnServerError(c, err)
		return
	}

	ret := map[string]any{
		"id":    thread_.Id,
		"runId": threadRun.Id,
	}
	if trace != nil {
		trace.SetBizResponse(ret)
	}
	ReturnSuccess(c, ret)
}

// SSE
func PostThreadStream(c *gin.Context) {
	req := &PostThreadStreamRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		ReturnBadRequest(c, err)
		return
	}

	if req.Id == "" {
		ReturnBadRequest(c, fmt.Errorf("id is empty"))
		return
	}

	var logger *slog.Logger
	if trace := utils.GetTraceLogger(c); trace != nil {
		logger = trace.Logger
		trace.SetBizRequest(req)
	}

	thread.StreamThread(req.Id, req.RunId, 0, c, logger)
}

func PostThreadDetail(c *gin.Context) {
	req := &PostThreadDetailRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		ReturnBadRequest(c, err)
		return
	}

	if req.Id == "" {
		ReturnBadRequest(c, fmt.Errorf("id is empty"))
		return
	}

	user := utils.GetUser(c)
	if user == nil {
		ReturnBadRequest(c, fmt.Errorf("unauthorized"))
		return
	}

	if trace := utils.GetTraceLogger(c); trace != nil {
		trace.SetBizRequest(req)
	}

	result, err := thread.DetailThread(user.Id, req.Id)
	if err != nil {
		ReturnServerError(c, err)
		return
	}

	ReturnSuccess(c, result)
}

func PostThreadDelete(c *gin.Context) {
	req := &PostThreadDeleteRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		ReturnBadRequest(c, err)
		return
	}

	if req.Id == "" {
		ReturnBadRequest(c, fmt.Errorf("id is empty"))
		return
	}

	user := utils.GetUser(c)
	if user == nil {
		ReturnBadRequest(c, fmt.Errorf("unauthorized"))
		return
	}

	if trace := utils.GetTraceLogger(c); trace != nil {
		trace.SetBizRequest(req)
	}

	err := thread.DeleteThread(user.Id, req.Id)
	if err != nil {
		ReturnServerError(c, err)
		return
	}

	ReturnSuccess(c, nil)
}

func PostListThread(c *gin.Context) {
	user := utils.GetUser(c)
	if user == nil {
		ReturnBadRequest(c, fmt.Errorf("unauthorized"))
		return
	}

	// TODO: list thread
	if user.AuthType == model.AUTH_NONE {
		ReturnBadRequest(c, fmt.Errorf("unauthorized"))
		return
	}

	threads, err := thread.ListThread(user.Id)
	if err != nil {
		ReturnServerError(c, err)
		return
	}

	ReturnSuccess(c, threads)
}
