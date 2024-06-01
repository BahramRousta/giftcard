package exceptions

const (
	InternalServerError      = "خطای داخلی رخ داده است"
	AuthenticationError      = "خطا در احراز هویت"
	StatusForbidden          = "خطای دسترسی"
	StatusBadRequest         = "درخواست نامعتبر"
	RecordNotFound           = "سفارش یافت نشد"
	InvalidInput             = "ورودی نامعتبر"
	InvalidCreateOrderInput  = "داده ورودی برای ایجاد درخواست نامعتبر می باشد"
	InvalidConfirmOrderInput = "داده ورودی برای تایید درخواست می باشد"
	RequiredOrderID          = "شماره پیگیری سفارش الزامیست"
	EmptyProductList         = "لیست سفارش نمی تواند خالی باشد"
	PageSizeError            = "تعداد لیست درخواست ها باید بین ۵ تا ۵۰ باشد"
	ProductIDError           = "آی دی محصول مورد نیاز است"
	DBError                  = "خطای دیتابیس"
)
