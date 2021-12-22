package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

const (
	privatePem = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA12xyGxqk340BIJpmO7+U2UXy4hLh0z2otvDDHCIL9+khD69J
6iBZ2rZDZmKNIGuYk7jVWdinvHWMx1/egHXiS9YNreP094PcxoQMGlpaxx8PVRbQ
RwlfgiFYv4V/BUkbZbp2ImedYpT1TXPMYOf7ehHThjCJxRG5By1X9oZP8BUSB41L
6kcG9CPDs/kOu0cm3TH9lL90MjQzhGjWxtWxgURMSxuY0g844oKy6lu1snES/XkA
W8Ew52hg6013g9WYqEZbamI+5ybci7ZcjVmQEBW5K5GzP7kkg5eWI8zue+iD04ts
+GLOvrigXwUZTN/gBbP2t8CV/TXxaEz1ycP7NQIDAQABAoIBAQCNro5xkxmCyftG
1SWOAFfGesHevGp4A2KWR00bkKzsdRDAfxoO+Q3/0cYMbZ7CBuIHrhXsDJltUqav
bjcpp96Y4ASJLJctvzUR/0DyiCCSO7Ra0zFStzBwprv24rcC4+03/+W4hQs3Dh8z
vIhb59c2lsjWxc+mpBxcQw9KXVXv4MW/KMTg06sIInpj83l6Gt4pLmVA+n+0xP61
LRGBmPTgFhvyMV0XEVQyAWtZtJpvEkoAWo4H3pCRGwl1dBnPYkzGxBEaJgpFeWga
ZAr/Hn4NihaB5qy6oA6HsJMrxQY0Ad/K5rsndawhz3lj/od51UbD/lowoo2sGh7R
0bJDbJGhAoGBAP+4plF4By21mivBBO96M05RWexp6LNN5cjPhCr87VkcyvC7O5rV
QezdHDQXCl6ofbao8z1U7uOkvoK1mU0VYFFg2CDLKFpw7iRrL//xky05Kl6fToVL
oGRKyWY6xVGCdHeCv6rgDdCAbtbWsrHgu2ZUA73ohz3S7bK7gTo253MJAoGBANeo
jWbydoOwvZIjTMx1G/4zbcmOMTvsCn+fqWYrsFjchkHuq3CbVInPqgpfpRtEoFFF
/F9eh2rAND5W+mbiZFCghTuHhxOax5Wyxz7yZWux3Qu1z/pT4WPcqDrrbA7iIqPM
PYlKctKRrYSj5frf1DJeLoBcHd+csQIHkABFvDXNAoGAJeZG+BIS9kpQ9CUiRx/U
VMonyqsTqudjo/RlgT2FK8zhovYM6nCq2aEXmfzEM61DHHxDuJZK5YA4IAUsGEmP
wd/ZiFqzu1u7X7hnH8a86lnrlqSDrau8tMCEwtr4/ZCZFFFTeM7GHV27j6m4SDan
b44KE+5PhPEq+29gwrD6cokCgYEAuMJj2oXxeSNrVg8+FZBjWiYPcfWLQq4X1H0i
MTFO1OKhd00VvdSl2ad7I4YLus/Rla+i5sXiuFdQqvPzdT+R9+1+F6El3WrmgN74
ino638gy+3xZYTqJx/dcfZYCLsIYMUKimcOZmcNK6G9Ocd9fOYOszTWeNlxU3ctC
2Kjl9SUCgYAt3cBY+vnqAu6dQ9TKEn3xo7qaPB382ZHjhvr4bwH62puaAxn7Vqln
4vBwKjElmmaMIbRFVU4M0KoVJnLJdpmLXoRSLKB6VrqjBkOcot8uVapdiUmUNL2w
t7fU9JqNK8CgQMp4gJFIUjv9LBJzje1GbXjuRoXtBQRjQ0WprxLLZw==
-----END RSA PRIVATE KEY-----`
)

// GenerateRSAKey 生成RSA私钥和公钥，保存到文件中
func GenerateRSAKey(bits int) {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: X509PrivateKey}
	//将数据保存到文件
	_ = pem.Encode(privateFile, &privateBlock)

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Id", Bytes: X509PublicKey}
	//保存到文件
	_ = pem.Encode(publicFile, &publicBlock)
}

// RSA 加签
func signRSA(plainText []byte) []byte {
	buf := []byte(privatePem)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//对密文进行加签
	hashed := sha256.Sum256(plainText)
	signText, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	//返回明文
	return signText
}

//RSA加密
func EncryptRSA(plainText []byte, pemTtr string) []byte {
	buf := []byte(pemTtr)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		panic(err)
	}
	//返回密文
	return cipherText
}

//RSA解密
func DecryptRSA(cipherText []byte) []byte {
	buf := []byte(privatePem)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	//返回明文
	return plainText
}

func BcryptRSA(cipherText string) string {
	cipher, _ := base64.StdEncoding.DecodeString(cipherText)
	// 解密
	plainText := DecryptRSA(cipher)
	return string(plainText)
}

func SignRSA(data []byte) string {
	sign := signRSA(data)
	return base64.StdEncoding.EncodeToString(sign)
}
