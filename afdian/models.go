package afdian

type AfdianWebhookResponse struct {
	Ec   int    `json:"ec"`
	Em   string `json:"em"`
	Data struct {
		Type  string `json:"type"`
		Order struct {
			OutTradeNo     string        `json:"out_trade_no"`
			UserID         string        `json:"user_id"`
			PlanID         string        `json:"plan_id"`
			Month          int           `json:"month"`
			TotalAmount    string        `json:"total_amount"`
			ShowAmount     string        `json:"show_amount"`
			Status         int           `json:"status"`
			Remark         string        `json:"remark"`
			RedeemID       string        `json:"redeem_id"`
			ProductType    int           `json:"product_type"`
			Discount       string        `json:"discount"`
			SkuDetail      []interface{} `json:"sku_detail"`
			AddressPerson  string        `json:"address_person"`
			AddressPhone   string        `json:"address_phone"`
			AddressAddress string        `json:"address_address"`
		} `json:"order"`
	} `json:"data"`
}

type AfdianQueryResponse struct {
	Ec   int    `json:"ec"`
	Em   string `json:"em"`
	Data struct {
		List []struct {
			OutTradeNo     string        `json:"out_trade_no"`
			UserID         string        `json:"user_id"`
			PlanID         string        `json:"plan_id"`
			Month          int           `json:"month"`
			TotalAmount    string        `json:"total_amount"`
			ShowAmount     string        `json:"show_amount"`
			Status         int           `json:"status"`
			Remark         string        `json:"remark"`
			RedeemID       string        `json:"redeem_id"`
			ProductType    int           `json:"product_type"`
			Discount       string        `json:"discount"`
			SkuDetail      []interface{} `json:"sku_detail"`
			PlanTitle      string        `json:"plan_title"`
			UserPrivateID  string        `json:"user_private_id"`
			AddressPerson  string        `json:"address_person"`
			AddressPhone   string        `json:"address_phone"`
			AddressAddress string        `json:"address_address"`
		} `json:"list"`
		TotalCount int `json:"total_count"`
		TotalPage  int `json:"total_page"`
		Request    struct {
			UserID string `json:"user_id"`
			Params string `json:"params"`
			Ts     string `json:"ts"`
			Sign   string `json:"sign"`
		} `json:"request"`
	} `json:"data"`
}

type AfdianProfileResponse struct {
	Ec   int    `json:"ec"`
	Em   string `json:"em"`
	Data struct {
		User struct {
			UserID         string `json:"user_id"`
			Status         int    `json:"status"`
			Name           string `json:"name"`
			Avatar         string `json:"avatar"`
			Cover          string `json:"cover"`
			URLSlug        string `json:"url_slug"`
			IsVerified     int    `json:"is_verified"`
			VerifiedType   int    `json:"verified_type"`
			IsNotRec       int    `json:"is_not_rec"`
			ShowSponsoring int    `json:"show_sponsoring"`
			DefaultPage    int    `json:"default_page"`
			ShowBadge      int    `json:"show_badge"`
			IsReject       int    `json:"is_reject"`
			IsBlock        int    `json:"is_block"`
			HasMark        int    `json:"has_mark"`
			Creator        struct {
				UserID                string `json:"user_id"`
				Status                int    `json:"status"`
				CategoryID            string `json:"category_id"`
				Type                  int    `json:"type"`
				Doing                 string `json:"doing"`
				Detail                string `json:"detail"`
				Pic                   string `json:"pic"`
				CustomPlan            int    `json:"custom_plan"`
				CustomPlanDesc        string `json:"custom_plan_desc"`
				ShowAlbum             int    `json:"show_album"`
				ShowShop              int    `json:"show_shop"`
				HomepageProductCount  int    `json:"homepage_product_count"`
				CanCopyText           int    `json:"can_copy_text"`
				CanCopyPic            int    `json:"can_copy_pic"`
				Watermark             int    `json:"watermark"`
				MonthlyFans           string `json:"monthly_fans"`
				MonthlyIncome         string `json:"monthly_income"`
				PrivacyPublicIncome   int    `json:"privacy_public_income"`
				PrivacyPublicSponsor  int    `json:"privacy_public_sponsor"`
				DiscountOption        int    `json:"discount_option"`
				LimitShowProductCount int    `json:"limit_show_product_count"`
				Category              struct {
					CategoryID string `json:"category_id"`
					Name       string `json:"name"`
				} `json:"category"`
				AlbumCount int `json:"album_count"`
			} `json:"creator"`
			IsSponsoring int `json:"is_sponsoring"`
		} `json:"user"`
	} `json:"data"`
}
