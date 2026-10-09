// Label tab INSTALLMENT (Non-Proportional) — `Section/TreatyInTabsNonProportional.xml`
// tab `Installment` + rincian baris `Section/Installments.xml` (FlowAction
// `Installments`, `pyEditingMode = expandPane`), dan gambar Pega
// `38-nonprop-menu-installment.png`.

export const ANGSURAN = {
  judul: 'Installment',
  /** `TreatyIn.InstallmentNo` — `pxTextInput`, label `Installment`. */
  jumlah: 'Installment',
  /** Tombol — `TreatyInSetValueInstallment(Installment, status=update)`. */
  perbaruiNilai: 'Update Value',
  /** Grid `TreatyIn.Installment` — satu kolom. */
  mataUang: 'Currency',
  /** Medan `.PctTotal` / `.AmountTotal` Section Installments. */
  totalPersen: '% Total',
  total: 'Total',
  /** Grid `TreatyIn.TotalInstallmentNP`. */
  totalAngsuran: 'Total Installment Amount',
  nilai: 'Value',
  /** Tombol — `TreatyInNPSetTotal(type=installment)`. */
  perbaruiTotal: 'Update Total',
  tanpaBaris: 'No items',
  rincian: 'Rincian baris',
} as const

/** Grid `.InstallmentList` Section Installments — urut sel ekspor / gambar 38. */
export const KOLOM_RINCIAN_ANGSURAN = ['Installment', 'Due Date', 'WPC (In Days)', 'Payment Date', '% Installment', 'Amount'] as const
