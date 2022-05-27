package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min-gateway/storage"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/util"
	"math"
	"net/http"
	"time"
)

type TokenService struct {
	abstract.Service
	JwtSignatureKey string `value:"${Authorization.SignatureKey}"`
	JwtExpiresIn    int64  `value:"${Authorization.Ttl}"`

	// Storage *storage.StorageService `autowired:"true"`
	Values *storage.ValuesService `autowired:"true"`
}

func (this *TokenService) Generate(storeId string, storePrefix string, singleton bool, authorizeType string, moreData map[string]interface{}) (accessToken string, expireAt int64, refreshToken string, refreshExpireAt int64, err error) {

	tokenId := util.GenerateUUID().String()
	refreshId := util.GenerateUUID().String()

	accessToken, expireAt, err = this.Create(tokenId, "")

	jwtExpiresIn := int64(math.MaxInt32)

	if this.JwtExpiresIn != -1 {
		jwtExpiresIn = this.JwtExpiresIn
	}

	if err == nil {
		refreshToken, refreshExpireAt, err = this.Create(refreshId, tokenId)

		if err == nil {
			/*min.Db.Insert(&mapper.SystemTokens{
				Id: tokenId,
				Expiretime: expireAt,
				UserId: storeId,
				RefreshId: refreshId,
			})*/

			if singleton {
				tokens := map[string]mapper.SystemTokens{}

				this.Values.GetAll(&tokens)

				for key, token := range tokens {
					if token.UserId == storeId && authorizeType == token.AuthorizeType {
						this.Values.Remove(key)
					}
				}
			}

			err = this.Values.Set(storePrefix+tokenId, mapper.SystemTokens{
				Id:            tokenId,
				Expiretime:    expireAt,
				AuthorizeType: authorizeType,
				UserId:        storeId,
				RefreshId:     refreshId,
				More:          moreData,
			}, jwtExpiresIn)
			return
		}
	}

	return "", 0, "", 0, err
}

func (this *TokenService) CustomGenerate(storeId string, storePrefix string, moreData map[string]interface{}, expiresIn int64) (expireAt int64, err error) {

	expireAt = int64(math.MaxInt32)

	if expiresIn != -1 {
		expireAt = time.Now().Unix() + expiresIn
	}

	err = this.Values.Set(storePrefix+storeId, mapper.SystemTokens{
		Id:         storeId,
		Expiretime: expireAt,
		UserId:     storeId,
		More:       moreData,
	}, expiresIn)
	return
}

func (this *TokenService) Auth(tokenString string, prefix string) bool {

	now := time.Now().Unix()
	claim, err := this.Parse(tokenString)
	token := mapper.SystemTokens{}

	if err != nil {
		return false
	}

	if now > claim.ExpiresAt {
		return false
	}
	has, err := this.Values.Get(prefix+":"+claim.Id, &token)

	if has && now <= token.Expiretime {
		return true
	}

	return false
}

func (this *TokenService) CustomAuth(tokenString string, prefix string) (ok bool, moreData map[string]interface{}, err error) {
	token := mapper.SystemTokens{}
	has, err := this.Values.Get(prefix+":"+tokenString, &token)
	if err != nil {
		return false, nil, err
	}
	if has {
		return true, token.More, nil
	}

	return false, nil, nil
}

func (this *TokenService) Fetch(tokenId string, prefix string) (storeId string, moreData map[string]interface{}, err error) {

	/*token := mapper.SystemTokens{Id:tokenId}

	has, err := min.Db.Get(&token)*/

	token := mapper.SystemTokens{}

	has, err := this.Values.Get(prefix+":"+tokenId, &token)

	if has {
		storeId = token.UserId
		moreData = token.More
	}
	return
}

/**
刷新token过期时间
*/
func (this *TokenService) Refresh(tokenId string, refreshId string, storePrefix string) (accessToken string, expireAt int64, refreshToken string, refreshExpireAt int64, err error) {
	expireAt = int64(math.MaxInt32)

	if this.JwtExpiresIn != -1 {
		expireAt = time.Now().Unix() + this.JwtExpiresIn
	}

	/*token := mapper.SystemTokens{
		Id:tokenId,
		RefreshId: refreshId,
	}

	has, err := min.Db.Get(&token)*/

	token := mapper.SystemTokens{}

	has, err := this.Values.Get(storePrefix+tokenId, &token)

	if !(has && token.RefreshId == refreshId) {
		has = false
	}

	if !has {
		return "", 0, "", 0, types.RuntimeError{Message: "未找到可用access token", ErrorCode: http.StatusNotFound}
	}

	accessToken, expireAt, err = this.Create(tokenId, "")

	if err == nil {
		newRefreshId := util.GenerateUUID().String()
		refreshToken, refreshExpireAt, err = this.Create(newRefreshId, tokenId)

		if err == nil {
			/*_, err = min.Db.
			Where("id = ?", tokenId).
			And("refresh_id = ?", refreshId).
			Cols("expiretime,refresh_id").
			Update(&mapper.SystemTokens{
				Expiretime: expireAt,
				RefreshId: newRefreshId,
			})*/

			jwtExpiresIn := int64(math.MaxInt32)

			if this.JwtExpiresIn != -1 {
				jwtExpiresIn = this.JwtExpiresIn
			}

			err = this.Values.Set(storePrefix+tokenId, mapper.SystemTokens{
				Id:         tokenId,
				Expiretime: expireAt,
				UserId:     token.UserId,
				RefreshId:  newRefreshId,
			}, jwtExpiresIn)
			return
		}
	}

	return
}

/**
解析jwt-token
*/
func (this *TokenService) Parse(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(this.JwtSignatureKey), nil
	})

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	return &jwt.StandardClaims{}, err
}

/**
生成jwt-token
*/
func (this *TokenService) Create(id string, issuer string) (token string, expireAt int64, err error) {
	expireAt = int64(math.MaxInt32)

	if this.JwtExpiresIn != -1 {
		expireAt = time.Now().Unix() + this.JwtExpiresIn
	}

	generated := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        id,
		ExpiresAt: expireAt,
		Issuer:    issuer,
	})

	token, err = generated.SignedString([]byte(this.JwtSignatureKey))

	return
}

func (this *TokenService) Delete(tokenId string, storePrefix string) bool {
	/*token := mapper.SystemTokens{Id:tokenId}
	_, _ = min.Db.Id(tokenId).Delete(&token)*/

	this.Values.Remove(storePrefix + tokenId)

	return true
}
