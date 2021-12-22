package rsa

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestGenerateRSAKey(t *testing.T) {
	//生成密钥对，保存到文件
	GenerateRSAKey(2048)
}

func TestEncryptRSA(t *testing.T) {
	message := []byte("zqkj123456")
	cipherText := EncryptRSA(message, "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA12xyGxqk340BIJpmO7+U\n2UXy4hLh0z2otvDDHCIL9+khD69J6iBZ2rZDZmKNIGuYk7jVWdinvHWMx1/egHXi\nS9YNreP094PcxoQMGlpaxx8PVRbQRwlfgiFYv4V/BUkbZbp2ImedYpT1TXPMYOf7\nehHThjCJxRG5By1X9oZP8BUSB41L6kcG9CPDs/kOu0cm3TH9lL90MjQzhGjWxtWx\ngURMSxuY0g844oKy6lu1snES/XkAW8Ew52hg6013g9WYqEZbamI+5ybci7ZcjVmQ\nEBW5K5GzP7kkg5eWI8zue+iD04ts+GLOvrigXwUZTN/gBbP2t8CV/TXxaEz1ycP7\nNQIDAQAB\n-----END PUBLIC KEY-----")
	fmt.Println("加密后为：", base64.StdEncoding.EncodeToString(cipherText))
}

func TestDecryptRSA(t *testing.T) {
	pwd := "Q2Jlg4Kk213UHmg9Al97MiSaO9QNt/4LwGEI2WV0ws3yfCKotgVE3Qzwh0OR75lJKispdaVwNKz+GpjS8m76tJhkYxkkANivPX5lzy0jno+pxgLRrWLS03XtIr4oQCh51gfWMWDWh6+1pEE5gymdM6WO9O1GINECD6caucrfosEVKNb3sUn2ETK/H3J0m2gZZgX5g/d4Gh0xyvK3+hcKHq3SdEbHPnhnOyGV6XwHomRcXI8ICJKvB6k3kqKpGDqdFZlNZrp75lx5mY4zF9wnPMIt06E4/Hb8cC2ie5ux1pdY6xeZTg7LSv1ILR8tudeFderiOBoUeIND/7B1D8iUSw=="
	message, _ := base64.StdEncoding.DecodeString(pwd)
	plainText := DecryptRSA(message)
	fmt.Println("解密后为：", string(plainText))
}
