package test

import (
	"context"
	proto "goodsSrv/proto/gen"
	"testing"
)

func TestGetCB(t *testing.T) {
	s := proto.NewCategoryBrandClient(SrvInit())
	defer SrvClose()
	b, err := s.GetCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id: 1,
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", b)
}

func TestListCB(t *testing.T) {
	s := proto.NewCategoryBrandClient(SrvInit())
	defer SrvClose()
	b, err := s.CategoryBrandList(context.Background(), &proto.CategoryBrandInfoRequest{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", b)
}

func TestCRUD(t *testing.T) {
	s := proto.NewCategoryBrandClient(SrvInit())
	defer SrvClose()
	create, err := s.CreateCategoryBrand(context.Background(), &proto.CreateCategoryBrandInfo{
		CategoryId: 90,
		BrandId:    50,
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", create)
	_, err = s.UpdateCategoryBrand(context.Background(), &proto.UpdateCategoryBrandInfo{
		Id:         create.Id,
		CategoryId: 90,
		BrandId:    49,
	})
	if err != nil {
		t.Error(err)
	}
	_, err = s.DeleteCategoryBrand(context.Background(), &proto.DeleteCategoryBrandInfo{
		Id: create.Id,
	})
	if err != nil {
		t.Error(err)
	}
}
