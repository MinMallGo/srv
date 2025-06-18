package test

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	proto "srv/goodsSrv/proto/gen"
	"testing"
)

func TestAllCategories(t *testing.T) {
	s := proto.NewCategoryClient(SrvInit())
	defer SrvClose()
	categories, err := s.GetAllCategories(context.Background(), &empty.Empty{})
	if err != nil {
		t.Error(err)
	}
	t.Log(categories.Data)
}

func TestCategory(t *testing.T) {
	s := proto.NewCategoryClient(SrvInit())
	defer SrvClose()
	category, err := s.GetSubCategory(context.Background(), &proto.CategoryInfoRequest{
		Level: 2,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(category.Data)
}

func TestCreateCategory(t *testing.T) {
	s := proto.NewCategoryClient(SrvInit())
	defer SrvClose()
	//category, err := s.CreateCategory(context.Background(), &proto.CreateCategoryInfo{
	//	Name:             "消息",
	//	ParentCategoryID: 0,
	//	Level:            2,
	//	IsTab:            false,
	//})
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//t.Log(category)
	//_, err = s.UpdateCategory(context.Background(), &proto.UpdateCategoryInfo{
	//	ID:               category.ID,
	//	Name:             "xxh__)_*改",
	//	ParentCategoryID: 1,
	//	Level:            1,
	//	IsTab:            false,
	//})
	//if err != nil {
	//	return
	//}

	_, err := s.DeleteCategory(context.Background(), &proto.DeleteCategoryInfo{
		Id: 92,
	})
	if err != nil {
		t.Error(err)
	}
}
