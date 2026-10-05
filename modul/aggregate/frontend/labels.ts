// Label modul Aggregate - bahasa Inggris. Caption mengikuti section Pega folder korpus `Aggregate`:
// `GridDasbordAgg` (daftar: Add Data, Tanggal Input, Ceding Code, Ceding Name, Treaty Type, As At, UW Year, Delete),
// `ShowAggregateList` (Master ID, Template, Upload CSV, Save, Close, kolom grid TempCSV), `ChooseMasterID` (kolom
// popup). Pesan penolakan Save datang dari backend apa adanya (`SaveAggregate_Act`).

import { formatBulat } from './aturan'

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 911, keputusan work owner 04-10-2026). */
export const MENU_AG = { kelompok: 'Aggregate' } as const

export const AG = {
  judul: 'Aggregate',
  addData: 'Add Data',
  cari: 'Search ceding name or code',
  tombolCari: 'Search',
  memuat: 'Loading aggregate data…',
  kosong: 'No aggregate data yet.',
  tidakCocok: 'No aggregate data matches the search.',
  jumlah: (n: number) => `${formatBulat(n)} uploads`,
  halaman: (h: number, total: number) => `Page ${h} of ${total}`,
  sebelumnya: '‹ Previous',
  berikutnya: 'Next ›',
  /** Kepala kolom nomor grid berkepala dua baris. */
  no: 'No',

  tanggalInput: 'Tanggal Input',
  cedingCode: 'Ceding Code',
  cedingName: 'Ceding Name',
  treatyType: 'Treaty Type',
  asAt: 'As At',
  uwYear: 'UW Year',
  jumlahBaris: 'Rows',
  aksi: 'Action',
  hapus: 'Delete',
  petunjukBaris: 'Double-click to see the details',

  judulHapus: 'Confirm Delete',
  kalimatHapus: (n: number) => `Delete this aggregate data (${formatBulat(n)} rows)? This cannot be undone.`,
  menghapus: 'Deleting…',
  batal: 'Cancel',
  terhapus: (n: number) => `${formatBulat(n)} rows deleted.`,

  judulRincian: 'Aggregate Detail',
  memuatRincian: 'Loading details…',
  tutup: 'Close',

  chart: 'RNM Value (USD)',
  chartKosong: 'No aggregate data for this As At.',
  memuatChart: 'Loading chart…',
  jejak: 'Chart level',
  semuaCeding: 'All Cedings',
  totalNilai: 'Total RNM Value (USD)',
  labelTingkat: { ceding: 'Number of Cedings', treatyType: 'Number of Treaty Types', coverage: 'Number of Coverages' },
  petunjukBuka: { ceding: 'Click a ceding to see its treaty types', treatyType: 'Click a treaty type to see its coverages' },

  masterId: 'Master ID',
  template: 'Template',
  uploadCsv: 'Upload CSV',
  save: 'Save',
  close: 'Close',
  membaca: 'Reading CSV…',
  menyimpan: 'Saving…',
  pratinjauKosong: 'Upload a CSV file to see its rows here.',
  masterKosong: 'No Master ID selected yet.',

  treatyContractName: 'Treaty Contract Name',
  reinsuranceType: 'Reinsurance Type',
  ceding: 'Ceding',
  sob: 'SOB',
  treatyGroup: 'Treaty Group',
  rnmShare: 'RNM Share',
  treatyYear: 'Treaty Year',

  judulMaster: 'Choose Master ID',
  cariMaster: 'Type a ceding name and press Enter',
  pilih: 'Select',
  submit: 'Submit',
  mencari: 'Searching…',
  masterTidakAda: 'No treaty found for this ceding.',
  dipilih: (n: number) => `${n} selected`,
} as const

/** Caption kolom grid TempCSV `ShowAggregateList`, menurut kolom Oracle. */
export const LABEL_KOLOM: Readonly<Record<string, string>> = {
  ASSESMENT_ZONE: 'Assesment Zone',
  M_TREATY_ID: 'Master Treaty',
  TREATY_TYPE: 'Treaty Type',
  COVERAGE: 'Coverage',
  AS_AT: 'As at',
  UW_YEAR: 'UW Year',
  TREATYYEAR: 'Treaty Year',
  CEDING_CODE: 'Ceding Code',
  CEDING_NAME: 'Ceding Names',
  CURRENCY: 'Currency',
  TO_USD: 'To_USD',
  NOR_BUILDINGS: 'NoR Buildings',
  BUILDINGS: 'IA Buildings',
  NOR_STOCKS: 'NoR Stocks',
  STOCKS: 'IA Stocks',
  NOR_MACHINERY: 'NoR Machinery',
  MACHINERY: 'IA Machinery',
  NOR_OTHER_CONTENTS: 'NoR Other Contents',
  OTHER_CONTENTS: 'IA Other Contents',
  NOR_CONSEQUENTIAL_LOSS: 'NoR Consequential Loss',
  CONSEQUENTIAL_LOSS: 'IA Consequential Loss',
  NOR_RESIDENTIAL: 'NoR Residential',
  RESIDENTIAL: 'IA Residential',
  NOR_COMMERCIAL: 'NoR Commercial',
  COMMERCIAL: 'IA Commercial',
  NOR_INDUSTRIAL: 'NoR Industrial',
  INDUSTRIAL: 'IA Industrial',
  NOR_AGRICULTURE: 'NoR Agriculture',
  AGRICULTURE: 'IA Agriculture',
  NOR_MISCELLANEOUS: 'NoR Miscellaneous',
  MISCELLANEOUS: 'IA Miscellaneous',
  NOR_UTILITIES: 'NoR Utilities',
  UTILITIES: 'IA Utilities',
  TOTAL_NO_OF_RISK: 'Total No of Risk',
  TOTAL_IN_AMOUNT: 'Total In Amount',
  TOTAL_IN_AMOUNT_IN_USD: 'Total In Amount (USD)',
  RNM_SHARE: 'RNM Share',
  RNM_VALUE: 'RNM Value',
  RNM_VALUE_IN_USD: 'RNM Value (USD)',
  REMARK: 'Remark',
  ID: 'ID',
}

/**
 * Kepala grid dua baris (perintah work owner 05-10-2026: "lakukan grouping sama seperti bordereaux"): pasangan
 * NoR (jumlah risiko) dan IA (nilai) dikelompokkan per kategori, lalu Total dan RNM. `[grup, judul bawah]`; kolom yang
 * tidak ada di sini tanpa grup - judul `LABEL_KOLOM`-nya menempati dua baris. Caption Pega lengkap tetap di
 * `LABEL_KOLOM` (dipakai sebagai `title` sel kepala).
 */
export const GRUP_KOLOM: Readonly<Record<string, readonly [string, string]>> = {
  NOR_BUILDINGS: ['Buildings', 'NoR'],
  BUILDINGS: ['Buildings', 'IA'],
  NOR_STOCKS: ['Stocks', 'NoR'],
  STOCKS: ['Stocks', 'IA'],
  NOR_MACHINERY: ['Machinery', 'NoR'],
  MACHINERY: ['Machinery', 'IA'],
  NOR_OTHER_CONTENTS: ['Other Contents', 'NoR'],
  OTHER_CONTENTS: ['Other Contents', 'IA'],
  NOR_CONSEQUENTIAL_LOSS: ['Consequential Loss', 'NoR'],
  CONSEQUENTIAL_LOSS: ['Consequential Loss', 'IA'],
  NOR_RESIDENTIAL: ['Residential', 'NoR'],
  RESIDENTIAL: ['Residential', 'IA'],
  NOR_COMMERCIAL: ['Commercial', 'NoR'],
  COMMERCIAL: ['Commercial', 'IA'],
  NOR_INDUSTRIAL: ['Industrial', 'NoR'],
  INDUSTRIAL: ['Industrial', 'IA'],
  NOR_AGRICULTURE: ['Agriculture', 'NoR'],
  AGRICULTURE: ['Agriculture', 'IA'],
  NOR_MISCELLANEOUS: ['Miscellaneous', 'NoR'],
  MISCELLANEOUS: ['Miscellaneous', 'IA'],
  NOR_UTILITIES: ['Utilities', 'NoR'],
  UTILITIES: ['Utilities', 'IA'],
  TOTAL_NO_OF_RISK: ['Total', 'No of Risk'],
  TOTAL_IN_AMOUNT: ['Total', 'In Amount'],
  TOTAL_IN_AMOUNT_IN_USD: ['Total', 'In Amount (USD)'],
  RNM_SHARE: ['RNM', 'Share'],
  RNM_VALUE: ['RNM', 'Value'],
  RNM_VALUE_IN_USD: ['RNM', 'Value (USD)'],
}
