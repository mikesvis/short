package subnet

import "net"

func ValidateSubnet(clientIP, subnet string) bool {
	// Настройка открыта для всех - можно
	if subnet == "" {
		return true
	}

	// Настройка закрыта по маске, но X-Real-IP пустой - нельзя
	if clientIP == "" {
		return false
	}

	// Ошибка при парсинге - нельзя
	_, cidr, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}

	// IP не входит в доверенную сеть - нельзя
	if !cidr.Contains(net.ParseIP(clientIP)) {
		return false
	}

	// МОЖНА!
	return true
}
