package v1

import (
	"context"

	portalv1 "github.com/ming-0x0/scaffold/api/gen/go/portal/v1"
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/internal/shared/repository/condition"
	"github.com/ming-0x0/scaffold/pkg/utils"
)

type BannerHandler struct {
	portalv1.UnimplementedPortalBannerServer
	*Handler
}

func NewBannerHandler(
	handler *Handler,
) *BannerHandler {
	return &BannerHandler{
		Handler: handler,
	}
}

func (h *BannerHandler) GetListBanner(ctx context.Context, req *portalv1.GetListBannerRequest) (*portalv1.GetListBannerResponse, error) {
	bannerRepo := h.adapter.BannerRepository()

	pageData := utils.WithPagination(req.Page, req.Limit)

	conditions := make([]condition.Condition, 0, 5)
	conditions = append(conditions, condition.PreloadAssociations())

	if req.Status != nil {
		conditions = append(conditions, condition.EQ("status", utils.Dereference(req.Status)))
	}

	if req.Name != nil {
		conditions = append(conditions,
			condition.OR(
				condition.LIKE("name_vi", utils.Dereference(req.Name)),
				condition.LIKE("name_en", utils.Dereference(req.Name)),
				condition.LIKE("name_zh", utils.Dereference(req.Name)),
			),
		)
	}

	banners, total, err := bannerRepo.FindPaginatedByConditions(
		ctx,
		pageData,
		conditions...,
	)
	if err != nil {
		return nil, h.errorResponder.Respond(err)
	}

	resourceMap := make(map[int32]*domain.Resource, len(banners))
	if len(banners) > 0 {
		resourceRepo := h.adapter.ResourceRepository()
		resourceIDs := make([]int32, 0, len(banners))
		for _, banner := range banners {
			resourceIDs = append(resourceIDs, banner.ResourceID)
		}

		resources, err := resourceRepo.FindByConditions(ctx, condition.IN("id", utils.ToAnySlice(resourceIDs)))
		if err != nil {
			return nil, h.errorResponder.Respond(err)
		}

		for _, resource := range resources {
			resourceMap[resource.ID] = resource
		}
	}

	return &portalv1.GetListBannerResponse{
		Banners: utils.TransformSlice(banners, func(banner *domain.Banner) *portalv1.Banner {
			return &portalv1.Banner{
				Id:            banner.ID,
				NameVi:        banner.NameVi,
				NameEn:        banner.NameEn,
				NameZh:        banner.NameZh,
				DescriptionVi: banner.DescriptionVi.Ptr(),
				DescriptionEn: banner.DescriptionEn.Ptr(),
				DescriptionZh: banner.DescriptionZh.Ptr(),
				Position:      banner.Position.Ptr(),
				Status:        int32(banner.Status),
				Resource: &portalv1.Resource{
					Id:        resourceMap[banner.ResourceID].ID,
					Type:      int32(resourceMap[banner.ResourceID].Type),
					Url:       resourceMap[banner.ResourceID].Url,
					YoutubeId: resourceMap[banner.ResourceID].YoutubeID.Ptr(),
				},
				Link:         banner.Link.Ptr(),
				ButtonNameVi: banner.ButtonNameVi.Ptr(),
				ButtonNameEn: banner.ButtonNameEn.Ptr(),
				ButtonNameZh: banner.ButtonNameZh.Ptr(),
				HasContent:   banner.HasContent,
			}
		}),
		TotalPage:   utils.CalcTotalPage(total, pageData["limit"]),
		RecordCount: total,
		Page:        pageData["page"],
		Limit:       pageData["limit"],
	}, nil
}

func (h *BannerHandler) GetBanner(ctx context.Context, req *portalv1.GetBannerRequest) (*portalv1.GetBannerResponse, error) {
	bannerRepo := h.adapter.BannerRepository()
	banner, err := bannerRepo.TakeByConditions(ctx, condition.EQ("id", req.BannerId))
	if err != nil {
		return nil, h.errorResponder.Respond(err)
	}

	resourceRepo := h.adapter.ResourceRepository()
	resource, err := resourceRepo.TakeByConditions(ctx, condition.EQ("id", banner.ResourceID))
	if err != nil {
		return nil, h.errorResponder.Respond(err)
	}

	return &portalv1.GetBannerResponse{
		Banner: &portalv1.Banner{
			Id:            banner.ID,
			NameVi:        banner.NameVi,
			NameEn:        banner.NameEn,
			NameZh:        banner.NameZh,
			DescriptionVi: banner.DescriptionVi.Ptr(),
			DescriptionEn: banner.DescriptionEn.Ptr(),
			DescriptionZh: banner.DescriptionZh.Ptr(),
			Position:      banner.Position.Ptr(),
			Status:        int32(banner.Status),
			Resource: &portalv1.Resource{
				Id:        resource.ID,
				Type:      int32(resource.Type),
				Url:       resource.Url,
				YoutubeId: resource.YoutubeID.Ptr(),
			},
			Link:         banner.Link.Ptr(),
			ButtonNameVi: banner.ButtonNameVi.Ptr(),
			ButtonNameEn: banner.ButtonNameEn.Ptr(),
			ButtonNameZh: banner.ButtonNameZh.Ptr(),
			HasContent:   banner.HasContent,
		},
	}, nil
}
