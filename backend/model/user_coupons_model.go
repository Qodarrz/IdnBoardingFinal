package models

import "time"

type UserCoupon struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	CouponID  int64      `db:"coupon_id"`
	ClaimedAt time.Time  `db:"claimed_at"`
	UsedAt    *time.Time `db:"used_at"`
	Status    string     `db:"status"`
}

func (UserCoupon) TableName() string {
	return "user_coupons"
}
