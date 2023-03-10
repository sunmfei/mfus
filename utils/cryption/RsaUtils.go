package cryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/pkcs12"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

// ReadFromPem 从Pem文件中读取秘钥
func ReadFromPem(pemFile string) ([]byte, error) {
	buffer, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(buffer)
	return block.Bytes, nil
}

// ReadFromPfx 从pfx文件中读取公私密钥（需要安装golang.org/x/crypto/pkcs12）
func ReadFromPfx(pfxFile, password string, usePKCS8 bool) ([]byte, []byte) {
	buffer, err := ioutil.ReadFile(pfxFile)
	if err != nil {
		panic(err)
	}

	privateKeyInterface, certificate, err := pkcs12.Decode(buffer, password)
	if err != nil {
		panic(err)
	}

	privateKey := privateKeyInterface.(*rsa.PrivateKey)
	publicKey := certificate.PublicKey.(*rsa.PublicKey)

	var (
		privateKeyBuffer []byte
		publicKeyBuffer  []byte
	)
	if usePKCS8 {
		privateKeyBuffer, err = x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			panic(err)
		}
		publicKeyBuffer, err = x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			panic(err)
		}
	} else {
		privateKeyBuffer = x509.MarshalPKCS1PrivateKey(privateKey)
		publicKeyBuffer = x509.MarshalPKCS1PublicKey(publicKey)
	}
	return publicKeyBuffer, privateKeyBuffer
}

// ReadPublicKeyFromCrt 从crt中读取公钥
func ReadPublicKeyFromCrt(crtFile string, usePKCS8 bool) ([]byte, error) {
	buffer, err := ioutil.ReadFile(crtFile)
	if err != nil {
		return nil, err
	}
	certDERBlock, _ := pem.Decode(buffer)
	certificate, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := certificate.PublicKey.(*rsa.PublicKey)

	var publicKeyBuffer []byte
	if usePKCS8 {
		publicKeyBuffer, err = x509.MarshalPKIXPublicKey(publicKey)
	} else {
		publicKeyBuffer = x509.MarshalPKCS1PublicKey(publicKey)
	}
	if err != nil {
		return nil, err
	}
	return publicKeyBuffer, nil
}

// WriteToPem 将秘钥写入Pem文件
func WriteToPem(isPrivateKey bool, buffer []byte, pemFile string) error {
	var _type string
	if isPrivateKey {
		_type = "RSA PRIVATE KEY"
	} else {
		_type = "RSA PUBLIC KEY"
	}

	block := &pem.Block{
		Type:  _type, //这个字符串随便写
		Bytes: buffer,
	}

	file, err := os.Create(pemFile)
	if err != nil {
		return err
	}
	return pem.Encode(file, block)
}

// Pkcs1ToPkcs8 Pkcs1转换为Pkcs8
func Pkcs1ToPkcs8(isPrivateKey bool, buffer []byte) []byte {
	var (
		oid  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
		info interface{}
	)
	if isPrivateKey {
		val := struct {
			Version    int
			Algo       []asn1.ObjectIdentifier
			PrivateKey []byte
		}{}
		val.Version = 0
		val.Algo = []asn1.ObjectIdentifier{oid}
		val.PrivateKey = buffer
		info = val
	} else {
		val := struct {
			Algo      pkix.AlgorithmIdentifier
			BitString asn1.BitString
		}{}
		val.Algo.Algorithm = oid
		val.Algo.Parameters = asn1.NullRawValue
		val.BitString.Bytes = buffer
		val.BitString.BitLength = 8 * len(buffer)
		info = val
	}

	b, err := asn1.Marshal(info)
	if err != nil {
		panic(err)
	}
	return b
}

// Pkcs8ToPkcs1 Pkcs8转换为Pkcs1
func Pkcs8ToPkcs1(isPrivateKey bool, buffer []byte) []byte {
	if isPrivateKey {
		val := struct {
			Version    int
			Algo       pkix.AlgorithmIdentifier
			PrivateKey []byte
		}{}
		_, err := asn1.Unmarshal(buffer, &val)
		if err != nil {
			panic(err)
		}
		return val.PrivateKey
	} else {
		val := struct {
			Algo      pkix.AlgorithmIdentifier
			BitString asn1.BitString
		}{}

		_, err := asn1.Unmarshal(buffer, &val)
		if err != nil {
			panic(err)
		}
		return val.BitString.Bytes
	}
}

// GenerateRsaKey 生成公私钥
// usePKCS8:是否使用pkcs8
func GenerateRsaKey(usePKCS8 bool) ([]byte, []byte) {
	//生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024) //1024位
	if err != nil {
		panic(err)
	}
	//公钥
	publicKey := privateKey.PublicKey

	var (
		privateKeyBuffer []byte
		publicKeyBuffer  []byte
	)

	if usePKCS8 {
		privateKeyBuffer, err = x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			panic(err)
		}
		publicKeyBuffer, err = x509.MarshalPKIXPublicKey(&publicKey)
		if err != nil {
			panic(err)
		}
	} else {
		privateKeyBuffer = x509.MarshalPKCS1PrivateKey(privateKey)
		publicKeyBuffer = x509.MarshalPKCS1PublicKey(&publicKey)
	}

	return publicKeyBuffer, privateKeyBuffer
}

func parsePkcsKey(buffer []byte, isPrivateKey, usePKCS8 bool) (interface{}, error) {
	var (
		err          error
		keyInterface interface{}
	)

	if isPrivateKey {
		if usePKCS8 {
			keyInterface, err = x509.ParsePKCS8PrivateKey(buffer)
		} else {
			keyInterface, err = x509.ParsePKCS1PrivateKey(buffer)
		}
	} else {
		if usePKCS8 {
			keyInterface, err = x509.ParsePKIXPublicKey(buffer)
		} else {
			keyInterface, err = x509.ParsePKCS1PublicKey(buffer)
		}
	}
	if err != nil {
		return nil, err
	}
	return keyInterface, nil
}

// RsaEncrypt RSA加密
func RsaEncrypt(value string, publicKey []byte, usePKCS8 bool) (string, error) {
	keyInterface, err := parsePkcsKey(publicKey, false, usePKCS8)
	if err != nil {
		return "", err
	}
	rsaPublicKey := keyInterface.(*rsa.PublicKey)
	buffer, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(value))
	if err != nil {
		return "", err
	}

	//以hex格式数值输出
	encryptText := fmt.Sprintf("%x", buffer)
	return encryptText, nil
}

// RsaDecrypt RSA解密
func RsaDecrypt(value string, privateKey []byte, usePKCS8 bool) (string, error) {
	//将hex格式数据转换为byte切片
	valueBytes := []byte(value)
	var buffer = make([]byte, len(valueBytes)/2)
	for i := 0; i < len(buffer); i++ {
		b, err := strconv.ParseInt(value[i*2:i*2+2], 16, 10)
		if err != nil {
			return "", err
		}
		buffer[i] = byte(b)
	}

	keyInterface, err := parsePkcsKey(privateKey, true, usePKCS8)
	if err != nil {
		return "", err
	}
	key := keyInterface.(*rsa.PrivateKey)
	buffer, err = rsa.DecryptPKCS1v15(rand.Reader, key, buffer)
	return string(buffer), nil
}

// Sign RSA签名
func Sign(value string, privateKey []byte, hash crypto.Hash, usePKCS8 bool) (string, error) {
	keyInterface, err := parsePkcsKey(privateKey, true, usePKCS8)
	if err != nil {
		return "", err
	}
	key := keyInterface.(*rsa.PrivateKey)

	var _hash = hash.New()
	if _, err := io.WriteString(_hash, value); err != nil {
		return "", err
	}

	hashed := _hash.Sum(nil)
	result, err := rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
	if err != nil {
		return "", err
	}

	//以hex格式数值输出
	encryptText := fmt.Sprintf("%x", result)
	return encryptText, nil
}

// Verify RSA验证签名
func Verify(value string, publicKey []byte, signature string, hash crypto.Hash, usePKCS8 bool) error {
	//将hex格式数据转换为byte切片
	valueBytes := []byte(signature)
	var buffer = make([]byte, len(valueBytes)/2)
	for i := 0; i < len(buffer); i++ {
		b, err := strconv.ParseInt(signature[i*2:i*2+2], 16, 10)
		if err != nil {
			return err
		}
		buffer[i] = byte(b)
	}

	keyInterface, err := parsePkcsKey(publicKey, false, usePKCS8)
	if err != nil {
		return err
	}

	key := keyInterface.(*rsa.PublicKey)

	var _hash = hash.New()
	if _, err := io.WriteString(_hash, value); err != nil {
		return err
	}

	hashed := _hash.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, buffer)
}
