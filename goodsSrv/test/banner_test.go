package test

import (
	"context"
	"fmt"
	proto "srv/goodsSrv/proto/gen"
	"testing"
)

func TestBannerCRUD(t *testing.T) {
	c = SrvInit()
	defer SrvClose()
	s := proto.NewBannerClient(c)
	res, err := s.CreateBanner(context.Background(), &proto.CreateBannerInfo{
		Image: "imgx",
		Url:   "urlx",
		Index: 0,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(res)

	_, err = s.UpdateBanner(context.Background(), &proto.UpdateBannerInfo{
		ID:    res.ID,
		Image: "img_u",
		Url:   "",
		Index: 0,
	})
	if err != nil {
		t.Error(err)
	}
	_, err = s.DeleteBanner(context.Background(), &proto.DeleteBannerInfo{
		Id: res.ID,
	})
	if err != nil {
		t.Error(err)
	}

	rs, err := s.BannerList(context.Background(), &proto.BannerInfoRequest{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rs.Total)
	for _, r := range rs.Data {
		fmt.Println(r)
	}
}
