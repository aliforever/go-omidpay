package omidpay

import "errors"

const (
	defaultPrefix           = "https://ref.sayancard.ir"
	tokenEndpoint           = "/ref-payment/RestServices/mts/generateTokenWithNoSign/"
	verificationEndpoint    = "/ref-payment/RestServices/mts/verifyMerchantTrans/"
	paymentURL              = "https://say.shaparak.ir/_ipgw_/MainTemplate/payment/"
	redirectionFormTemplate = `<!doctypehtml><html lang=en><meta charset=UTF-8><title>Pay...</title><style>.text-center{text-align:center}.mt-2{margin-top:2em}.spinner{margin:100px auto 0;width:70px;text-align:center}.spinner>div{width:18px;height:18px;background-color:#333;border-radius:100%;display:inline-block;-webkit-animation:sk-bouncedelay 1.4s infinite ease-in-out both;animation:sk-bouncedelay 1.4s infinite ease-in-out both}.spinner .bounce1{-webkit-animation-delay:-.32s;animation-delay:-.32s}.spinner .bounce2{-webkit-animation-delay:-.16s;animation-delay:-.16s}@-webkit-keyframes sk-bouncedelay{0%,100%,80%{-webkit-transform:scale(0)}40%{-webkit-transform:scale(1)}}@keyframes sk-bouncedelay{0%,100%,80%{-webkit-transform:scale(0);transform:scale(0)}40%{-webkit-transform:scale(1);transform:scale(1)}}</style><body onload=submitForm()><div class=spinner><div class=bounce1></div><div class=bounce2></div><div class=bounce3></div></div><form action={{.PaymentUrl}} class="mt-2 text-center"method=POST><p>در حال انتقال به درگاه پرداخت<p>در صورتیکه بعد از <span id=countdown>5</span> ثانیه... وارد درگاه پرداخت نشدید کلیک کنید</p><input name=token type=hidden value={{.PaymentToken}}> <input name=language type=hidden value=fa> <button type=submit>ورود به درگاه پرداخت</button></form><script>var seconds=5;function submitForm(){document.forms[0].submit()}function countdown(){(seconds-=1)<=0?submitForm():(document.getElementById("countdown").innerHTML=seconds,window.setTimeout("countdown()",5000))}countdown()</script>`
)

var statusCodes = map[string]string{
	"erSucceed":                      "سرویس با موفقیت اجراء شد.",
	"erAAS_UseridOrPassIsRequired":   "کد کاربری و رمز الزامی هست.",
	"erAAS_InvalidUseridOrPass":      "کد کاربری یا رمز صحیح نمی باشد.",
	"erAAS_InvalidUserType":          "نوع کاربر صحیح نمی‌باشد.",
	"erAAS_UserExpired":              "کاربر منقضی شده است.",
	"erAAS_UserNotActive":            "کاربر غیر فعال هست.",
	"erAAS_UserTemporaryInActive":    "کاربر موقتا غیر فعال شده است.",
	"erAAS_UserSessionGenerateError": "خطا در تولید شناسه لاگین",
	"erAAS_UserPassMinLengthError":   "حداقل طول رمز رعایت نشده است.",
	"erAAS_UserPassMaxLengthError":   "حداکثر طول رمز رعایت نشده است.",
	"erAAS_InvalidUserCertificate":   "برای کاربر فایل سرتیفکیت تعریف نشده است.",
	"erAAS_InvalidPasswordChars":     "کاراکترهای غیر مجاز در رمز",
	"erAAS_InvalidSession":           "شناسه لاگین معتبر نمی‌باشد ",
	"erAAS_InvalidChannelId":         "کانال معتبر نمی‌باشد.",
	"erAAS_InvalidParam":             "پارامترها معتبر نمی‌باشد.",
	"erAAS_NotAllowedToService":      "کاربر مجوز سرویس را ندارد.",
	"erAAS_SessionIsExpired":         "شناسه الگین معتبر نمی‌باشد.",
	"erAAS_InvalidData":              "داده‌ها معتبر نمی‌باشد.",
	"erAAS_InvalidSignature":         "امضاء دیتا درست نمی‌باشد.",
	"erAAS_InvalidToken":             "توکن معتبر نمی‌باشد.",
	"erAAS_InvalidSourceIp":          "آدرس آی پی معتبر نمی‌باشد.",

	"erMts_ParamIsNull":                        "پارمترهای ورودی خالی می‌باشد.",
	"erMts_UnknownError":                       "خطای ناشناخته",
	"erMts_InvalidAmount":                      "مبلغ معتبر نمی‌باشد.",
	"erMts_InvalidBillId":                      "شناسه قبض معتبر نمی‌باشد.",
	"erMts_InvalidPayId":                       "شناسه پرداخت معتبر نمی‌باشد.",
	"erMts_InvalidEmailAddLen":                 "طول ایمیل معتبر نمی‌باشد.",
	"erMts_InvalidGoodsReferenceIdLen":         "طول شناسه خرید معتبر نمی‌باشد.",
	"erMts_InvalidMerchantGoodsReferenceIdLen": "طول شناسه خرید پذیرنده معتبر نمی‌باشد.",
	"erMts_InvalidMobileNo":                    "فرمت شماره موبایل معتبر نمی‌باشد.",
	"erMts_InvalidPorductId":                   "طول یا فرمت کد محصول معتبر نمی‌باشد.",
	"erMts_InvalidRedirectUrl":                 "طول یا فرمت آدرس صفحه رجوع معتبر نمی‌باشد.",
	"erMts_InvalidReferenceNum":                "طول یا فرمت شماره رفرنس معتبر نمی‌باشد.",
	"erMts_InvalidRequestParam":                "پارامترهای درخواست معتبر نمی‌باشد.",
	"erMts_InvalidReserveNum":                  "طول یا فرمت شماره رزرو معتبر نمی‌باشد.",
	"erMts_InvalidSessionId":                   "شناسه الگین معتبر نمی‌باشد.",
	"erMts_InvalidSignature":                   "طول یا فرمت امضاء دیتا معتبر نمی‌باشد.",
	"erMts_InvalidTerminal":                    "کد ترمینال معتبر نمی‌باشد.",
	"erMts_InvalidToken":                       "توکن معتبر نمی‌باشد.",
	"erMts_InvalidTransType":                   "نوع تراکنش معتبر نمی‌باشد.",
	"erMts_InvalidUniqueId":                    "کد یکتا معتبر نمی‌باشد.",
	"erMts_InvalidUseridOrPass":                "رمز یا کد کاربری معتبر نمی باشد.",
	"erMts_RepeatedBillId":                     "پرداخت قبض تکراری می باشد.",
	"erMts_AASError":                           "کد کاربری و رمز الزامی هست.",
	"erMts_SCMError":                           "خطای سرور مدیریت کانال",

	"erScm_erOrgTransNotExists":    "تراکنش مورد نظر یافت نشد.",
	"erScm_InvalidReferenceNum":    "شماره ارجاءه معتبر نمی باشد.",
	"erScm_OrgTransReversedBefore": "تراکنش برگشت خورده است",
}

var (
	MissingParams error = errors.New("missing_parameters")
	MissingRefNum error = errors.New("missing_ref_num")
)
