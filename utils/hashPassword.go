package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 生成密码的哈希值
func HashPassword(password string) (string, error) {
	// 使用bcrypt生成哈希值，其中14是加密的成本因子，可以根据需要调整
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, hashedPassword string) error {
	// 将哈希值和明文密码进行比较，如果匹配则验证通过
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
