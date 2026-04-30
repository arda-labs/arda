import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component } from '@angular/core';
import { TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';

@Component({
  selector: 'app-banking-reference-page',
  standalone: true,
  imports: [CommonModule, TableModule, Tag],
  templateUrl: './banking-reference.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class BankingReferencePage {
  readonly capabilities = [
    { name: 'Tiền tệ & tỷ giá', code: 'CURRENCY, FX_RATE', owner: 'Treasury/Core Banking' },
    { name: 'Lịch làm việc, ngày nghỉ', code: 'BUSINESS_CALENDAR', owner: 'Operations' },
    { name: 'Ngân hàng, chi nhánh, SWIFT/NAPAS', code: 'BANK, BANK_BRANCH, PAYMENT_NETWORK', owner: 'Payments' },
    { name: 'Sản phẩm & kênh giao dịch', code: 'BANKING_PRODUCT, SERVICE_CHANNEL, PRODUCT_CHANNEL_RULE', owner: 'Product/Digital Banking' },
    { name: 'Loại giấy tờ, nghề nghiệp, ngành kinh tế', code: 'ID_TYPE, OCCUPATION, ECONOMIC_SECTOR', owner: 'KYC' },
    { name: 'Nhóm rủi ro, AML, sanction source', code: 'RISK_GRADE, AML_LIST_SOURCE', owner: 'Risk & Compliance' },
    { name: 'Biểu phí, thuế, hạn mức chuẩn', code: 'FEE_TYPE, TAX_CODE, LIMIT_PROFILE', owner: 'Product/Finance' },
  ];
}
