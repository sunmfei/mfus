package commonUtils

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
	"time"
)

// mathConfig 生成图形化算术验证码配置
func mathConfig() *base64Captcha.DriverMath {
	mathType := &base64Captcha.DriverMath{
		Height:          50,
		Width:           100,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine,
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return mathType
}

// digitConfig 生成图形化数字验证码配置
func digitConfig() *base64Captcha.DriverDigit {
	digitType := &base64Captcha.DriverDigit{
		Height:   50,
		Width:    100,
		Length:   5,
		MaxSkew:  0.45,
		DotCount: 80,
	}
	return digitType
}

// stringConfig 生成图形化字符串验证码配置
func stringConfig() *base64Captcha.DriverString {
	stringType := &base64Captcha.DriverString{
		Height:          100,
		Width:           50,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          5,
		Source:          "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM",
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return stringType
}

// chineseConfig 生成图形化汉字验证码配置
// "设想,你在,处理,消费者,的音,频输,出音,频可,能无,论什,么都,没有,任何,输出,或者,它可,能是,单声道,立体声,或是,环绕立,体声的,不想要,的值"
func chineseConfig(str string) *base64Captcha.DriverChinese {
	chineseType := &base64Captcha.DriverChinese{
		Height:          50,
		Width:           100,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowSlimeLine,
		Length:          2,
		Source:          str,
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return chineseType
}

// autoConfig 生成图形化数字音频验证码配置
func autoConfig() *base64Captcha.DriverAudio {
	chineseType := &base64Captcha.DriverAudio{
		Length:   4,
		Language: "zh",
	}
	return chineseType
}

// CreateCode audio 音频验证码、string 字符串+数字验证码、math 算术运算验证码、chinese 纯汉字验证码、digit 纯数字验证码
func (b Base64Code) CreateCode(typ string) (string, string, error) {
	var driver base64Captcha.Driver
	// switch case分支中的方法为目录3的配置
	// switch case分支中的方法为目录3的配置
	// switch case分支中的方法为目录3的配置
	switch typ {
	case "audio":
		driver = autoConfig()
	case "string":
		driver = stringConfig()
	case "math":
		driver = mathConfig()
	case "chinese":
		driver = chineseConfig(b.World)
	case "digit":
		driver = digitConfig()
	}
	if driver == nil {
		panic("生成验证码的类型没有配置")
	}
	// 创建验证码并传入创建的类型的配置，以及存储的对象
	c := base64Captcha.NewCaptcha(driver, b.Store)
	id, b64s, err := c.Generate()
	return id, b64s, err
}

// VerifyCaptcha 校验验证码
// Verify(id, VerifyValue, true) 中的 true参数
// 当为 true 时，校验 传入的id 的验证码，校验完 这个ID的验证码就要在内存中删除
// 当为 false 时，校验 传入的id 的验证码，校验完 这个ID的验证码不删除
func (b Base64Code) VerifyCaptcha(id, VerifyValue string) bool {
	// result 为步骤1 创建的图片验证码存储对象
	return b.Store.Verify(id, VerifyValue, true)
}
func (b Base64Code) VerifyCaptchaBool(id, VerifyValue string, bl bool) bool {
	// result 为步骤1 创建的图片验证码存储对象
	return b.Store.Verify(id, VerifyValue, bl)
}

// GetCodeAnswer 获取验证码答案
// Get(codeId, false) 中的 false 参数
// 当为 true 时，根据ID获取完验证码就要删除这个验证码
// 当为 false 时，根据ID获取完验证码不删除
func (b Base64Code) GetCodeAnswer(codeId string) string {
	// result 为步骤1 创建的图片验证码存储对象
	return b.Store.Get(codeId, false)
}

func (b Base64Code) GetCodeAnswerBool(codeId string, bl bool) string {
	// result 为步骤1 创建的图片验证码存储对象
	return b.Store.Get(codeId, bl)
}

type Base64Code struct {
	Store base64Captcha.Store
	World string
}

func NewBase64Code(store base64Captcha.Store, world string) *Base64Code {

	if store == nil {
		store = base64Captcha.NewMemoryStore(20240, 3*time.Minute)
	}

	return &Base64Code{
		store,
		world,
	}
}
