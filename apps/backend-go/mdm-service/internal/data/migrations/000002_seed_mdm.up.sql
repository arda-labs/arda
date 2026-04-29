INSERT INTO code_sets (code, name, description, is_system)
VALUES
  ('ADMIN_UNIT_LEVEL', 'Cấp đơn vị hành chính', 'Các cấp hành chính đang được MDM Service hỗ trợ', true),
  ('ADMIN_UNIT_TYPE', 'Loại đơn vị hành chính', 'Tỉnh, thành phố, phường, xã, đặc khu và cấp huyện legacy', true),
  ('MDM_STATUS', 'Trạng thái dùng chung', 'Trạng thái chuẩn cho dữ liệu MDM/master data', true),
  ('SYSTEM_PARAMETER_GROUP', 'Nhóm tham số hệ thống', 'Nhóm cấu hình hệ thống', true),
  ('CURRENCY', 'Tiền tệ', 'Danh mục tiền tệ sử dụng trong tài chính kế toán và thanh toán', true),
  ('COUNTRY', 'Quốc gia', 'Danh mục quốc gia phục vụ KYC, giao dịch quốc tế và báo cáo', true),
  ('BANK_ACCOUNT_TYPE', 'Loại tài khoản ngân hàng', 'Các loại tài khoản nghiệp vụ ngân hàng', true),
  ('PAYMENT_METHOD', 'Phương thức thanh toán', 'Phương thức thanh toán chuẩn toàn hệ thống', true),
  ('PAYMENT_CHANNEL', 'Kênh giao dịch', 'Kênh thực hiện giao dịch ngân hàng', true),
  ('TRANSACTION_TYPE', 'Loại giao dịch', 'Phân loại giao dịch tài chính ngân hàng', true),
  ('CUSTOMER_SEGMENT', 'Phân khúc khách hàng', 'Phân khúc khách hàng phục vụ sản phẩm, phí và hạn mức', true),
  ('RISK_RATING', 'Xếp hạng rủi ro', 'Nhóm rủi ro phục vụ KYC/AML và phê duyệt', true),
  ('DOCUMENT_TYPE', 'Loại giấy tờ', 'Loại giấy tờ định danh và hồ sơ pháp lý', true),
  ('INTEREST_RATE_TYPE', 'Loại lãi suất', 'Phân loại lãi suất cho sản phẩm tiền gửi, vay và đầu tư', true),
  ('COLLATERAL_TYPE', 'Loại tài sản bảo đảm', 'Phân loại tài sản bảo đảm cho tín dụng', true),
  ('FEE_TYPE', 'Loại phí', 'Danh mục phí dùng chung cho sản phẩm và giao dịch', true)
ON CONFLICT DO NOTHING;

INSERT INTO code_items (code_set_id, code, name, value, sort_order, is_system)
SELECT cs.id, item.code, item.name, item.value, item.sort_order, true
FROM code_sets cs
JOIN (
  VALUES
    ('ADMIN_UNIT_LEVEL', 'PROVINCE', 'Tỉnh/Thành phố', 'PROVINCE', 10),
    ('ADMIN_UNIT_LEVEL', 'WARD', 'Phường/Xã/Đặc khu', 'WARD', 20),
    ('ADMIN_UNIT_LEVEL', 'DISTRICT_LEGACY', 'Quận/Huyện legacy', 'DISTRICT_LEGACY', 30),
    ('ADMIN_UNIT_TYPE', 'TINH', 'Tỉnh', 'TINH', 10),
    ('ADMIN_UNIT_TYPE', 'THANH_PHO', 'Thành phố trực thuộc trung ương', 'THANH_PHO', 20),
    ('ADMIN_UNIT_TYPE', 'PHUONG', 'Phường', 'PHUONG', 30),
    ('ADMIN_UNIT_TYPE', 'XA', 'Xã', 'XA', 40),
    ('ADMIN_UNIT_TYPE', 'DAC_KHU', 'Đặc khu', 'DAC_KHU', 50),
    ('MDM_STATUS', 'ACTIVE', 'Đang hoạt động', 'ACTIVE', 10),
    ('MDM_STATUS', 'INACTIVE', 'Ngưng hoạt động', 'INACTIVE', 20),
    ('MDM_STATUS', 'MERGED', 'Đã sáp nhập', 'MERGED', 30),
    ('MDM_STATUS', 'DELETED', 'Đã xóa', 'DELETED', 40),
    ('SYSTEM_PARAMETER_GROUP', 'SECURITY', 'Bảo mật', 'SECURITY', 10),
    ('SYSTEM_PARAMETER_GROUP', 'UI', 'Giao diện', 'UI', 20),
    ('SYSTEM_PARAMETER_GROUP', 'NOTIFICATION', 'Thông báo', 'NOTIFICATION', 30),
    ('SYSTEM_PARAMETER_GROUP', 'INTEGRATION', 'Tích hợp', 'INTEGRATION', 40),
    ('SYSTEM_PARAMETER_GROUP', 'FINANCE', 'Tài chính', 'FINANCE', 50),
    ('SYSTEM_PARAMETER_GROUP', 'BANKING', 'Ngân hàng', 'BANKING', 60),
    ('SYSTEM_PARAMETER_GROUP', 'ACCOUNTING', 'Kế toán', 'ACCOUNTING', 70),
    ('SYSTEM_PARAMETER_GROUP', 'RISK', 'Rủi ro', 'RISK', 80),
    ('SYSTEM_PARAMETER_GROUP', 'COMPLIANCE', 'Tuân thủ', 'COMPLIANCE', 90),
    ('SYSTEM_PARAMETER_GROUP', 'LIMIT', 'Hạn mức', 'LIMIT', 100),
    ('SYSTEM_PARAMETER_GROUP', 'PRICING', 'Giá/phí', 'PRICING', 110),
    ('CURRENCY', 'VND', 'Vietnamese Dong', 'VND', 10),
    ('CURRENCY', 'USD', 'US Dollar', 'USD', 20),
    ('CURRENCY', 'EUR', 'Euro', 'EUR', 30),
    ('CURRENCY', 'JPY', 'Japanese Yen', 'JPY', 40),
    ('COUNTRY', 'VN', 'Việt Nam', 'VN', 10),
    ('COUNTRY', 'US', 'United States', 'US', 20),
    ('BANK_ACCOUNT_TYPE', 'CASA', 'Tài khoản thanh toán', 'CASA', 10),
    ('BANK_ACCOUNT_TYPE', 'SAVINGS', 'Tài khoản tiết kiệm', 'SAVINGS', 20),
    ('BANK_ACCOUNT_TYPE', 'LOAN', 'Tài khoản vay', 'LOAN', 30),
    ('BANK_ACCOUNT_TYPE', 'GL', 'Tài khoản sổ cái', 'GL', 40),
    ('PAYMENT_METHOD', 'CASH', 'Tiền mặt', 'CASH', 10),
    ('PAYMENT_METHOD', 'BANK_TRANSFER', 'Chuyển khoản', 'BANK_TRANSFER', 20),
    ('PAYMENT_METHOD', 'CARD', 'Thẻ', 'CARD', 30),
    ('PAYMENT_METHOD', 'QR', 'QR', 'QR', 40),
    ('PAYMENT_CHANNEL', 'BRANCH', 'Quầy giao dịch', 'BRANCH', 10),
    ('PAYMENT_CHANNEL', 'ATM', 'ATM', 'ATM', 20),
    ('PAYMENT_CHANNEL', 'INTERNET_BANKING', 'Internet Banking', 'INTERNET_BANKING', 30),
    ('PAYMENT_CHANNEL', 'MOBILE_BANKING', 'Mobile Banking', 'MOBILE_BANKING', 40),
    ('PAYMENT_CHANNEL', 'POS', 'POS', 'POS', 50),
    ('TRANSACTION_TYPE', 'DEPOSIT', 'Nộp tiền', 'DEPOSIT', 10),
    ('TRANSACTION_TYPE', 'WITHDRAWAL', 'Rút tiền', 'WITHDRAWAL', 20),
    ('TRANSACTION_TYPE', 'TRANSFER', 'Chuyển tiền', 'TRANSFER', 30),
    ('TRANSACTION_TYPE', 'PAYMENT', 'Thanh toán', 'PAYMENT', 40),
    ('TRANSACTION_TYPE', 'FEE', 'Thu phí', 'FEE', 50),
    ('TRANSACTION_TYPE', 'INTEREST', 'Lãi', 'INTEREST', 60),
    ('CUSTOMER_SEGMENT', 'RETAIL', 'Cá nhân', 'RETAIL', 10),
    ('CUSTOMER_SEGMENT', 'SME', 'Doanh nghiệp SME', 'SME', 20),
    ('CUSTOMER_SEGMENT', 'CORPORATE', 'Doanh nghiệp lớn', 'CORPORATE', 30),
    ('CUSTOMER_SEGMENT', 'VIP', 'Khách hàng ưu tiên', 'VIP', 40),
    ('RISK_RATING', 'LOW', 'Rủi ro thấp', 'LOW', 10),
    ('RISK_RATING', 'MEDIUM', 'Rủi ro trung bình', 'MEDIUM', 20),
    ('RISK_RATING', 'HIGH', 'Rủi ro cao', 'HIGH', 30),
    ('DOCUMENT_TYPE', 'ID_CARD', 'Căn cước công dân', 'ID_CARD', 10),
    ('DOCUMENT_TYPE', 'PASSPORT', 'Hộ chiếu', 'PASSPORT', 20),
    ('DOCUMENT_TYPE', 'BUSINESS_LICENSE', 'Giấy phép kinh doanh', 'BUSINESS_LICENSE', 30),
    ('INTEREST_RATE_TYPE', 'FIXED', 'Cố định', 'FIXED', 10),
    ('INTEREST_RATE_TYPE', 'FLOATING', 'Thả nổi', 'FLOATING', 20),
    ('COLLATERAL_TYPE', 'REAL_ESTATE', 'Bất động sản', 'REAL_ESTATE', 10),
    ('COLLATERAL_TYPE', 'VEHICLE', 'Phương tiện', 'VEHICLE', 20),
    ('COLLATERAL_TYPE', 'DEPOSIT', 'Tiền gửi', 'DEPOSIT', 30),
    ('COLLATERAL_TYPE', 'GUARANTEE', 'Bảo lãnh', 'GUARANTEE', 40),
    ('FEE_TYPE', 'ACCOUNT_MAINTENANCE', 'Phí duy trì tài khoản', 'ACCOUNT_MAINTENANCE', 10),
    ('FEE_TYPE', 'TRANSFER', 'Phí chuyển tiền', 'TRANSFER', 20),
    ('FEE_TYPE', 'CARD_ANNUAL', 'Phí thường niên thẻ', 'CARD_ANNUAL', 30),
    ('FEE_TYPE', 'EARLY_SETTLEMENT', 'Phí tất toán trước hạn', 'EARLY_SETTLEMENT', 40)
) AS item(set_code, code, name, value, sort_order)
  ON cs.code = item.set_code
ON CONFLICT DO NOTHING;

INSERT INTO area_types (code, name, description, allow_hierarchy, status)
VALUES
  ('BUSINESS_REGION', 'Khu vực kinh doanh', 'Vùng kinh doanh/quản lý theo nhu cầu nghiệp vụ', true, 'ACTIVE'),
  ('DELIVERY_ZONE', 'Khu vực giao hàng', 'Vùng phục vụ giao vận hoặc vận hành', true, 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO system_parameters (
  key, name, group_code, value_type, value_text, value_number, value_boolean, value_json,
  default_value, is_system, description, status
)
VALUES
  ('DEFAULT_LANGUAGE', 'Ngôn ngữ mặc định', 'UI', 'STRING', 'vi', NULL, NULL, '{}', 'vi', true, 'Ngôn ngữ mặc định của hệ thống', 'ACTIVE'),
  ('DEFAULT_TIMEZONE', 'Múi giờ mặc định', 'UI', 'STRING', 'Asia/Ho_Chi_Minh', NULL, NULL, '{}', 'Asia/Ho_Chi_Minh', true, 'Múi giờ mặc định cho ngày giờ hiển thị', 'ACTIVE'),
  ('SESSION_TIMEOUT_MINUTES', 'Thời gian hết phiên', 'SECURITY', 'NUMBER', '', 30, NULL, '{}', '30', true, 'Thời gian hết phiên đăng nhập tính theo phút', 'ACTIVE'),
  ('ENABLE_MAINTENANCE_MODE', 'Bật chế độ bảo trì', 'SECURITY', 'BOOLEAN', '', NULL, false, '{}', 'false', true, 'Chặn các thao tác không cần thiết khi hệ thống bảo trì', 'ACTIVE'),
  ('FINANCE_DEFAULT_CURRENCY', 'Tiền tệ mặc định', 'FINANCE', 'STRING', 'VND', NULL, NULL, '{}', 'VND', true, 'Tiền tệ mặc định cho kế toán và báo cáo tài chính', 'ACTIVE'),
  ('BANKING_BUSINESS_DATE_MODE', 'Chế độ ngày làm việc ngân hàng', 'BANKING', 'STRING', 'SYSTEM_DATE', NULL, NULL, '{}', 'SYSTEM_DATE', true, 'Cách xác định ngày nghiệp vụ mặc định cho giao dịch ngân hàng', 'ACTIVE')
ON CONFLICT DO NOTHING;
