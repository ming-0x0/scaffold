package v1

import (
	"context"
	"errors"

	portalv1 "github.com/ming-0x0/scaffold/api/gen/go/portal/v1"
	"github.com/ming-0x0/scaffold/internal/adapter/grpc/responder"
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/pkg/domainerror"
	"github.com/ming-0x0/scaffold/pkg/utils"
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
	}, h.userRepo.PreloadAssociations())
	if err != nil {
		var domainErr *domainerror.DomainError
		if errors.As(err, &domainErr) && domainErr.ErrorCode() == domainerror.NotFound {
			return nil, h.errorResponder.RespondMsg("Không tìm thấy tài khoản", err)
		}

		return nil, h.errorResponder.Respond(err)
	}

	checkPasswordHash := utils.CheckPasswordHash(req.Password, user.Password)
	if !checkPasswordHash {
		return nil, h.errorResponder.RespondCode(domainerror.FailedPrecondition, "Sai mật khẩu", nil)
	}

	if user.Status != domain.UserStatusActive {
		return nil, h.errorResponder.RespondCode(domainerror.FailedPrecondition, "Không thể đăng nhập do tài khoản của bạn đang bị khóa", nil)
	}

	return &portalv1.LoginResponse{
		AccessToken: user.Email,
	}, nil
}
