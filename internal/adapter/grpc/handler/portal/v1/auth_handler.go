package v1

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	portalv1 "github.com/ming-0x0/scaffold/api/gen/go/portal/v1"
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/pkg/iterator"
	"github.com/ming-0x0/scaffold/pkg/jwt"
	"github.com/ming-0x0/scaffold/pkg/utils"
)

type AuthHandler struct {
	portalv1.UnimplementedPortalAuthServer
	*Handler
}

func NewAuthHandler(
	handler *Handler,
) *AuthHandler {
	return &AuthHandler{
		Handler: handler,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *portalv1.LoginRequest) (*portalv1.LoginResponse, error) {
	userRepo := h.adapter.UserRepository()

	user, err := userRepo.TakeByConditions(
		ctx,
		userRepo.EQ("username", req.Username),
		userRepo.PreloadAssociations(),
	)
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

	portalFuncList := []*domain.Permission{}

	roleRepo := h.adapter.RoleRepository()
	permissionRepo := h.adapter.PermissionRepository()
	if !user.PermissionGroup.FullPermission {
		roles, err := roleRepo.FindByConditions(
			ctx,
			roleRepo.EQ("permission_group_id", user.PermissionGroupID),
			roleRepo.PreloadAssociations(),
		)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			portalFuncList = append(portalFuncList, &role.Permission)
		}
	} else {
		portalFuncList, err = permissionRepo.FindByConditions(
			ctx,
			permissionRepo.PreloadAssociations(),
		)
		if err != nil {
			return nil, err
		}
	}

	portalFuncList = iterator.From(portalFuncList).Unique(func(p1, p2 *domain.Permission) bool {
		return p1.ID == p2.ID
	}).Collect()

	expiredAt := time.Now().Add(time.Hour * 24 * 7)
	tokenID := uuid.New().String()

	accessToken, err := jwt.GenerateHS256JWT(map[string]any{
		"user_id":  user.ID,
		"sub":      user.Username,
		"email":    user.Email,
		"exp":      expiredAt.Unix(),
		"iat":      time.Now().Unix(),
		"token_id": tokenID,
	})
	if err != nil {
		return nil, err
	}

	go func(userID int64, accessToken string, expiredAt time.Time, tokenID string) {
		userToken := &domain.UserToken{
			UserID:    userID,
			Token:     accessToken,
			ExpiredAt: expiredAt,
			TokenID:   tokenID,
		}
		err := h.adapter.UserTokenRepository().Save(ctx, userToken)
		if err != nil {
			return
		}
	}(user.ID, accessToken, expiredAt, tokenID)

	return &portalv1.LoginResponse{
		AccessToken: accessToken,
		Permissions: utils.TransformSlice(portalFuncList, func(p *domain.Permission) *portalv1.Permission {
			return &portalv1.Permission{
				FunctionCode: p.FunctionCode,
			}
		}),
		TokenId: tokenID,
		User: &portalv1.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
		},
	}, nil
}
