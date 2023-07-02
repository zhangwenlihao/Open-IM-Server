package api

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/a2r"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/constant"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/mcontext"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/discoveryregistry"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/errs"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/proto/third"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/rpcclient"
)

type ThirdApi rpcclient.Third

func NewThirdApi(discov discoveryregistry.SvcDiscoveryRegistry) ThirdApi {
	return ThirdApi(*rpcclient.NewThird(discov))
}

func (o *ThirdApi) ApplyPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ApplyPut, o.Client, c)
}

func (o *ThirdApi) GetPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetPut, o.Client, c)
}

func (o *ThirdApi) ConfirmPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ConfirmPut, o.Client, c)
}

func (o *ThirdApi) GetHash(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetHashInfo, o.Client, c)
}

func (o *ThirdApi) FcmUpdateToken(c *gin.Context) {
	a2r.Call(third.ThirdClient.FcmUpdateToken, o.Client, c)
}

func (o *ThirdApi) SetAppBadge(c *gin.Context) {
	a2r.Call(third.ThirdClient.SetAppBadge, o.Client, c)
}

func (o *ThirdApi) GetURL(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		a2r.Call(third.ThirdClient.GetUrl, o.Client, c)
		return
	}
	name := c.Query("name")
	if name == "" {
		c.String(http.StatusBadRequest, "name is empty")
		return
	}
	operationID := c.Query("operationID")
	if operationID == "" {
		operationID = "auto_" + strconv.Itoa(rand.Int())
	}
	expires, _ := strconv.ParseInt(c.Query("expires"), 10, 64)
	if expires <= 0 {
		expires = 3600 * 1000
	}
	attachment, _ := strconv.ParseBool(c.Query("attachment"))
	c.Set(constant.OperationID, operationID)
	resp, err := o.Client.GetUrl(mcontext.SetOperationID(c, operationID), &third.GetUrlReq{Name: name, Expires: expires, Attachment: attachment})
	if err != nil {
		if errs.ErrArgs.Is(err) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if errs.ErrRecordNotFound.Is(err) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, resp.Url)
}
