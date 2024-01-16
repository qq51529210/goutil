package media

import (
	"context"
	"goutil/onvif/wsdl/ver10/schema"
)

// type getProfilesReq struct {
// 	soap.Envelope[any, struct {
// 		XMLName string `xml:"tt:GetProfiles"`
// 	}]
// }

// type getProfilesRes struct {
// 	soap.Envelope[any, getProfilesResponse]
// }

// type getProfilesResponse struct {
// 	Profiles []*schema.Profile
// }

// This operation gets the device system date and time.
// The device shall support the return of the daylight
// saving setting and of the manual system date and time
// (if applicable) or indication of NTP time (if applicable)
// through the GetProfiles command.
// A device shall provide the UTCDateTime information.
func GetProfiles(ctx context.Context, url string) ([]*schema.Profile, error) {
	// // 消息
	// var _req getProfilesReq
	// _req.Envelope.Attr = envelopeAttr
	// var _res getProfilesRes
	// // 请求
	// err := soap.Do(ctx, url, &_req, &_res)
	// if err != nil {
	// 	return nil, err
	// }
	// // 错误
	// if _res.Body.Fault != nil {
	// 	return nil, _res.Body.Fault
	// }
	// // 成功
	// return _res.Body.Data.Profiles, nil
	return nil, nil
}
