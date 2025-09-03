package v1

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ming-0x0/scaffold/adapter/grpc/responder"
	portalv1 "github.com/ming-0x0/scaffold/api/gen/go/portal/v1"
	"github.com/ming-0x0/scaffold/domain"
	"github.com/ming-0x0/scaffold/shared/domainerror"
	"github.com/ming-0x0/scaffold/shared/utils"
)

type AuthHandler struct {
	portalv1.UnimplementedPortalAuthServer
	userRepo       domain.UserRepositoryInterface
	errorResponder responder.ErrorResponderInterface
}

func NewAuthHandler(
	userRepo domain.UserRepositoryInterface,
	errorResponder responder.ErrorResponderInterface,
) *AuthHandler {
	return &AuthHandler{
		userRepo:       userRepo,
		errorResponder: errorResponder,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *portalv1.LoginRequest) (*portalv1.LoginResponse, error) {
	user, err := h.userRepo.TakeByConditions(ctx, map[string]any{
		"username": req.Username,
	}, h.userRepo.PreloadAssociations(ctx))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, h.errorResponder.WrapMsg("Không tìm thấy tài khoản!", err)
		}

		return nil, h.errorResponder.Wrap(err)
	}

	checkPasswordHash := utils.CheckPasswordHash(req.Password, user.Password)
	if !checkPasswordHash {
		return nil, h.errorResponder.WrapCode("Sai mật khẩu!", domainerror.InvalidArgument)
	}

	if user.Status != domain.UserStatusActive {
		return nil, h.errorResponder.WrapCode("Không thể đăng nhập do tài khoản của bạn đang bị khóa!", domainerror.FailedPrecondition)
	}

	return &portalv1.LoginResponse{
		Token: user.Email,
	}, nil
}
