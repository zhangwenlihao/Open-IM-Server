package friend

import "github.com/xiaoyiEdu/Open-IM-Server/pkg/errs"

func (m *ApplyToAddFriendReq) Check() error {
	if m.GetToUserID() == "" {
		return errs.ErrArgs.Wrap("get toUserID is empty")
	}
	return nil
}
