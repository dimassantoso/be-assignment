package usecase

import (
	"billing-engine/pkg/shared"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/golangid/candi/candishared"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"

	"billing-engine/internal/modules/auth/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *authUsecaseImpl) LoginAuth(ctx context.Context, req *domain.RequestLogin) (result domain.ResponseLogin, attemptLeft int, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "AuthUsecase:LoginAuth")
	defer trace.Finish()

	attemptLeft, _ = uc.AttemptLogin(ctx, req.Email)
	if attemptLeft <= 0 {
		return result, attemptLeft, fmt.Errorf("attempt limit reached")
	}

	user, err := uc.repoSQL.AuthRepo().Find(ctx, &domain.FilterAuth{
		Email: req.Email,
	})
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return result, attemptLeft, fmt.Errorf("user not found")
		}
		return
	}

	isPasswordMatch, err := uc.IsPasswordMatch(ctx, req.Password, user.Password)
	if !isPasswordMatch {
		return result, attemptLeft, fmt.Errorf("password not match")
	}

	uc.cache.Delete(ctx, "login_attempt:"+req.Email)

	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	if req.KeepSignIn {
		// keep login for a week
		expiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	claims := candishared.TokenClaim{
		StandardClaims: jwt.StandardClaims{
			Subject:   fmt.Sprintf("%d", result.ID),
			ExpiresAt: expiresAt,
		},
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := generatedToken.SignedString([]byte(shared.GetEnv().JWTSecret))
	if err != nil {
		return
	}

	result = domain.ResponseLogin{
		ID:    user.ID,
		Email: user.Email,
		Token: signedToken,
	}

	return
}

func (uc *authUsecaseImpl) AttemptLogin(ctx context.Context, email string) (left int, err error) {
	trace, _ := tracer.StartTraceWithContext(ctx, "AuthUsecase:AttemptLogin")
	defer trace.Finish()

	cache, _ := uc.cache.Get(ctx, "login_attempt:"+email)
	if cache == nil {
		uc.cache.Set(ctx, "login_attempt:"+email, []byte{1}, time.Minute*5)
		return 4, nil
	}
	uc.cache.Set(ctx, "login_attempt:"+email, []byte{cache[0] + 1}, time.Minute*5)
	return 5 - int(cache[0]+1), nil
}

func (uc *authUsecaseImpl) IsPasswordMatch(ctx context.Context, password, hash string) (bool, error) {
	trace, _ := tracer.StartTraceWithContext(ctx, "AuthUsecase:IsPasswordMatch")
	defer trace.Finish()
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
