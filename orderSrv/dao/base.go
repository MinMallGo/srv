package dao

import (
	"fmt"
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

// HandleError handle住错误并记录错误，然后判断影响的条数和实际受影响的条数是否一致
func HandleError(res *gorm.DB, affect int) error {
	if res.Error != nil {
		// 通过反射获取是哪个方法调用的。
		funcName := "unknown"
		if pc, _, _, ok := runtime.Caller(1); ok {
			if fun := runtime.FuncForPC(pc); fun != nil {
				funcName = fun.Name()
			}
		}

		zap.L().Error(fmt.Sprintf("%s:error :", funcName), zap.Error(res.Error))
		return status.Error(codes.Internal, "内部错误")
	}

	// 如果为0，说明不需要检查这里，比如什么获取列表之类的操作
	if affect != 0 && res.RowsAffected != int64(affect) {
		return status.Error(codes.Internal, "数据库操作失败失败")
	}
	return nil
}
