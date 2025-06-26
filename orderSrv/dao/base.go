package dao

import (
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"runtime"
)

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// HandleError 检查 DB 操作错误及影响行数是否符合预期。
// affect == 0 表示忽略行数检查。
func HandleError(res *gorm.DB, expectAffect int) error {
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		zap.L().Error("gorm error",
			zap.String("caller", callerFunc(2)), // skip=2 是更常见的调用层级
			zap.Error(res.Error),
		)
		return status.Error(codes.Internal, "数据库错误")
	}

	if expectAffect > 0 && res.RowsAffected != int64(expectAffect) {
		zap.L().Warn("affected rows mismatch",
			zap.String("caller", callerFunc(2)),
			zap.Int("expected", expectAffect),
			zap.Int64("actual", res.RowsAffected),
		)
		return status.Error(codes.Internal, "操作失败")
	}

	return nil
}

// callerFunc 获取调用者函数名（skip=0 当前函数，1 调用者...）
func callerFunc(skip int) string {
	if pc, _, _, ok := runtime.Caller(skip); ok {
		if f := runtime.FuncForPC(pc); f != nil {
			return f.Name()
		}
	}
	return "unknown"
}
