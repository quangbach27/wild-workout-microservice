package cqrs

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type commandLoggingDecorator[C any] struct {
	baseHandler CommandHandler[C]
	logger      *logrus.Entry
}

func (commandLogging commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	logger := commandLogging.logger.WithFields(logrus.Fields{
		"command":      generateActionName(cmd),
		"command_body": fmt.Sprintf("%#v", cmd),
	})

	logger.Debug("Executing command")
	defer func() {
		if err == nil {
			logger.Info("Command executed successfully")
		} else {
			logger.WithError(err).Error("Command execution failed")
		}
	}()

	return commandLogging.baseHandler.Handle(ctx, cmd)
}

type queryLoggingDecorator[Q any, R any] struct {
	baseHandler QueryHandler[Q, R]
	logger      *logrus.Entry
}

func (queryLogging queryLoggingDecorator[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {

	logger := queryLogging.logger.WithFields(logrus.Fields{
		"query":      generateActionName(query),
		"query_body": fmt.Sprintf("%#v", query),
	})

	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query executed successfully")
		} else {
			logger.WithError(err).Error("Query execution failed")
		}
	}()

	return queryLogging.baseHandler.Handle(ctx, query)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
