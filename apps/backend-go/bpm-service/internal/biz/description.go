package biz

// DescriptionTemplate represents a reusable text for responses, rejections, etc.
type DescriptionTemplate struct {
	ID      string
	Code    string
	Content string
	Module  string
}

type DescriptionUseCase struct {
	// Add repo
}

func NewDescriptionUseCase() *DescriptionUseCase {
	return &DescriptionUseCase{}
}

// In a real app, this would use a pattern-based approach for dynamic variables
// e.g., "Kính chào {customer_name}, hồ sơ của bạn đã được phê duyệt."
