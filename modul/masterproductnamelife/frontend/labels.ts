// Label modul Master Product Name Life - paket 10.
//
// ⛔ VERBATIM dari korpus `D:\XML\RNM_BRD\Master Product Name Life\` (satu tag per baris; nomor baris mentah =
// nomor `sed -e 's/></>\n</g'`). Nomor di komentar adalah baris TAG yang memuat teksnya, dan `labels.test.ts`
// membuka korpus di baris itu. Teks tombol = `pyModes.pyLabel` (yang dirender `pxButton`/`pxLink`), label medan
// = `pyLabelFieldValue`, judul/kolom = `pyValue`, tombol dialog = `pySubmitLabel`/`pyCancelLabel` FlowAction -
// BUKAN `pyLabelPreview` (pratinjau perancang).
//
// Label korpus berbahasa Inggris (kecuali yang memang berbahasa Indonesia di korpus). Pesan server VERBATIM
// (`Product Name Empty`, `Plan tidak boleh sama`, ...) tampil apa adanya dari jawaban server.
//
// Teks yang TIDAK ada di korpus ditandai `[tidak ada di korpus]` beserta alasannya; dijaga `labels.test.ts`.

/** Menu - `M_NAV_MENU.LABEL` sesudah slot 961. */
export const MENU_MPNL = {
  /**
   * `[tidak ada di korpus]` - nama tampilan: folder korpus `D:\XML\RNM_BRD\Master Product Name Life` tanpa kata
   * "Master" (keputusan work owner 03-10-2026, nama tampilan saja; migrasi 961). Kode modul tetap.
   */
  kelompok: 'Product Name Life',
} as const


/** Grid daftar produk - `InboxProductName` mode daftar (wadah b71246, RD `BrowseProduct_Life`). */
export const GRID_MPNL = {
  /** `InboxProductName.xml` b71783 `<pyLabelFieldValue>`. */
  labelSelAdd: 'End Period',
  /** `InboxProductName.xml` b71865 `<pyLabel>`. */
  add: 'Add',
  /** `InboxProductName.xml` b71863 `<pyTooltip>`. */
  tooltipAdd: 'Add New Data',
  /** `InboxProductName.xml` b72403 `<pyValue>`. */
  kolomId: 'ID',
  /** `InboxProductName.xml` b72551 `<pyValue>`. */
  kolomCeding: 'Ceding',
  /** `InboxProductName.xml` b72707 `<pyValue>`. */
  kolomTreatyNumber: 'Treaty Number',
  /** `InboxProductName.xml` b72866 `<pyValue>`. */
  kolomTreatyName: 'Treaty Name',
  /** `InboxProductName.xml` b73025 `<pyValue>`. */
  kolomCreateOp: 'Create Operator',
  /** `InboxProductName.xml` b73184 `<pyValue>`. */
  kolomUpdateOp: 'Last Updated Operator',
  /** `InboxProductName.xml` b74798 `<pyLabel>`. */
  view: 'View',
} as const

/** Form sisi umum - halaman `ProductName` (PARITAS §3.1). */
export const UMUM_MPNL = {
  /** `InboxProductName.xml` b2934 `<pyTitle>` - judul kartu sisi umum. */
  judul: 'TREATY NAME',
  /** `InboxProductName.xml` b3620 `<pyLabelFieldValue>`. */
  productName: 'Product Name',
  /** `InboxProductName.xml` b3894 `<pyLabelFieldValue>`. */
  productCode: 'Product Code',
  /** `InboxProductName.xml` b4075 `<pyLabelFieldValue>`. */
  ceding: 'Ceding',
  /** `InboxProductName.xml` b5222 `<pyLabel>`. */
  chooseCeding: 'Choose Ceding Name',
  /** `InboxProductName.xml` b4463 `<pyLabelFieldValue>`. */
  sob: 'SOB',
  /** `InboxProductName.xml` b5563 `<pyLabel>`. */
  chooseSob: 'Choose SOB',
  /** `InboxProductName.xml` b7104 `<pyLabelFieldValue>`. */
  deduction: 'Deduction (%)',
  /** `InboxProductName.xml` b7398 `<pyLabelFieldValue>`. */
  riRisk: 'R/I Risk Name',
  /** `InboxProductName.xml` b8205 `<pyLabel>`. */
  chooseRiRisk: 'Choose R/I Risk',
  /** `InboxProductName.xml` b8719 `<pyLabelFieldValue>`. */
  treatyName: 'Treaty Name',
  /** `InboxProductName.xml` b8900 `<pyLabelFieldValue>`. */
  treatyNumber: 'Treaty Number',
  /** `InboxProductName.xml` b10661 `<pyLabelFieldValue>`. */
  causeOfLoss: 'Cause Of Loss',
  /** `InboxProductName.xml` b11360 `<pyLabel>`. */
  chooseCause: 'Choose Cause of Loss',
  /** `InboxProductName.xml` b47312 `<pyCheckboxCaption>`. */
  onRetention: 'On Retention',
} as const

/** Form sisi inward - halaman `ProductNameInward` (PARITAS §3.2). */
export const INWARD_MPNL = {
  /** `InboxProductName.xml` b16621 `<pyTitle>` - judul kartu sisi inward. */
  judul: 'INWARD',
  /** `InboxProductName.xml` b17097 `<pyLabelFieldValue>`. */
  policyHolder: 'Policy Holder',
  /** `InboxProductName.xml` b17827 `<pyLabel>`. */
  choosePolicyHolder: 'Choose Policy Holder',
  /** `InboxProductName.xml` b18341 `<pyLabelFieldValue>`. */
  insured: 'Insured',
  /** `InboxProductName.xml` b18834 `<pyLabelFieldValue>`. */
  addendumNo: 'Addendum No.',
  /** `InboxProductName.xml` b19432 `<pyLabelFieldValue>`. */
  addendum: 'Addendum',
  /** `InboxProductName.xml` b20155 `<pyLabelFieldValue>`. */
  amandementNo: 'Amandement No.',
  /** `InboxProductName.xml` b20760 `<pyLabelFieldValue>`. */
  amandement: 'Amandement',
  /** `InboxProductName.xml` b21781 `<pyLabelFieldValue>`. */
  maxExpiredClaim: 'Max Notification Claim Expired',
  /** `InboxProductName.xml` b21969 `<pyLabelFieldValue>`. */
  begin: 'Begin Date',
  /** `InboxProductName.xml` b22304 `<pyLabelFieldValue>`. */
  stnc: 'STNC',
  /** `InboxProductName.xml` b22494 `<pyLabelFieldValue>`. */
  cedingRetention: 'Ceding Retention (%)',
  /** `InboxProductName.xml` b22657 `<pyLabelFieldValue>`. */
  cedingLimit: 'Ceding\'s Limit',
  /** `InboxProductName.xml` b22915 `<pyLabelFieldValue>`. */
  brokerage: 'Brokerage Fee (%)',
  /** `InboxProductName.xml` b23078 `<pyLabelFieldValue>`. */
  minAge: 'Minimum Age (Years)',
  /** `InboxProductName.xml` b23265 `<pyLabelFieldValue>`. */
  maxAge: 'Maximum Age (Years)',
  /** `InboxProductName.xml` b23452 `<pyLabelFieldValue>`. */
  expiryAge: 'Expiry Age (Years)',
  /** `InboxProductName.xml` b23635 `<pyLabelFieldValue>`. */
  extraPremi: 'Extra Premium',
  /** `InboxProductName.xml` b23796 `<pyLabelFieldValue>`. */
  minSumInsured: 'Min Sum Insured',
  /** `InboxProductName.xml` b23956 `<pyLabelFieldValue>`. */
  maxSumInsured: 'Max Sum Insured',
  /** `InboxProductName.xml` b24214 `<pyLabelFieldValue>`. */
  maxSumReasured: 'Max Sum Reasured',
  /** `InboxProductName.xml` b24674 `<pyLabelFieldValue>`. */
  rnmShare: 'Nusantara Re Share (%)',
  /** `InboxProductName.xml` b24880 `<pyValue>`. */
  ofSumReasured: 'OF SUM REASURED',
  /** `InboxProductName.xml` b25236 `<pyLabelFieldValue>`. */
  rnmLimit: 'Nusantara Re\'s Limit',
  /** `InboxProductName.xml` b25398 `<pyLabelFieldValue>`. */
  premiumFactor: 'Premium Factor (%)',
  /** `InboxProductName.xml` b25611 `<pyLabelFieldValue>`. */
  payment: 'Premium Payment Method',
  /** `InboxProductName.xml` b25924 `<pyLabelFieldValue>`. */
  subjectTo: 'Subject To',
  /** `InboxProductName.xml` b26092 `<pyLabelFieldValue>`. */
  annuityInterest: 'Annuity Interest (%)',
  /** `InboxProductName.xml` b26306 `<pyLabelFieldValue>`. */
  premiumRefundFactor: 'Premium Refund Factor (%)',
  /** `InboxProductName.xml` b27062 `<pyLabelFieldValue>`. */
  maxDataReceive: 'Max Production Data Receive',
  /** `InboxProductName.xml` b27250 `<pyLabelFieldValue>`. */
  mature: 'Expired Date',
  /** `InboxProductName.xml` b27960 `<pyLabelFieldValue>`. */
  birthday: 'Birthday',
  /** `InboxProductName.xml` b28140 `<pyLabelFieldValue>`. */
  currency: 'Currency',
  /** `InboxProductName.xml` b28742 `<pyLabel>`. */
  chooseCurrency: 'Choose Currency',
  /** `InboxProductName.xml` b29256 `<pyLabelFieldValue>`. */
  extraMortality: 'Extra Mortality (%)',
  /** `InboxProductName.xml` b29442 `<pyLabelFieldValue>`. */
  maxContract: 'Max Contract (year)',
  /** `InboxProductName.xml` b29626 `<pyLabelFieldValue>`. */
  proportionalTable: 'Proportional Table',
} as const

/** Grid `LIEN CLAUSE` - `ProductName.LienClause` (RALAT R10). */
export const LIEN_MPNL = {
  /** `InboxProductName.xml` b12201 `<pyValue>`. */
  judul: 'LIEN CLAUSE (Potongan Manfaat Klaim)',
  /** `InboxProductName.xml` b12741 `<pyValue>`. */
  usia: 'Usia saat Klaim',
  /** `InboxProductName.xml` b12890 `<pyValue>`. */
  manfaat: '% Manfaat yang dibayarkan',
} as const

/** Grid `DOCUMENT CLAIM` - `ProductName.DocumentClaim`. */
export const DOKUMEN_MPNL = {
  /** `InboxProductName.xml` b14601 `<pyValue>`. */
  judul: 'DOCUMENT CLAIM',
  /** `InboxProductName.xml` b15125 `<pyValue>`. */
  documentList: 'Document List',
} as const

/** Grid `PLAN LIST` - `ProductName.PlanList`. */
export const PLAN_MPNL = {
  /** `InboxProductName.xml` b31557 `<pyValue>`. */
  judul: 'PLAN LIST',
  /** `InboxProductName.xml` b31845 `<pyValue>`. */
  planName: 'Plan Name',
  /** `InboxProductName.xml` b31994 `<pyValue>`. */
  bussines: 'Bussines',
  /** `InboxProductName.xml` b32143 `<pyValue>`. */
  benefit: 'Benefit',
  /** `InboxProductName.xml` b32296 `<pyValue>`. */
  riRate: 'R/I Rate',
  /** `InboxProductName.xml` b32823 `<pyLabel>`. */
  add: 'Add',
  /** `InboxProductName.xml` b34113 `<pyLabel>`. */
  viewRate: 'View Rate',
  /** `InboxProductName.xml` b34589 `<pyLabel>`. */
  chooseRiRate: 'Choose R/I Rate',
  /** `InboxProductName.xml` b35075 `<pyLabel>`. */
  delete: 'Delete',
} as const

/**
 * Kolom daftar `Plan Name` (`InboxProductName.xml` b33121 `pxAutoComplete` - dropdown sejak 03-10-2026, RD `BrowseProductTypeLife_RD`):
 * `pyAdditionalFields` b33213 `.ID`, b33247 `.CoverName`, b33280 `.Business`, b33313 `.Benefit` - semuanya `pyShow` true.
 * Labelnya label kolom RD sumber (`BrowseProductTypeLife_RD.xml` `pyFieldLabel`).
 */
export const SARAN_PLAN_MPNL = {
  /** `BrowseProductTypeLife_RD.xml` b568 `<pyFieldLabel>`. */
  kolomId: 'ID',
  /** `BrowseProductTypeLife_RD.xml` b642 `<pyFieldLabel>`. */
  kolomCoverName: 'CoverName',
  /** `BrowseProductTypeLife_RD.xml` b583 `<pyFieldLabel>`. */
  kolomBusiness: 'Business',
  /** `BrowseProductTypeLife_RD.xml` b613 `<pyFieldLabel>`. */
  kolomBenefit: 'Benefit',
} as const

/** Grid `FINANCIAL UNDERWRITING` - `ProductName.FinancialUnderwritingList`. */
export const FINUW_MPNL = {
  /** `InboxProductName.xml` b37148 `<pyValue>`. */
  judul: 'FINANCIAL UNDERWRITING',
  /** `InboxProductName.xml` b37670 `<pyValue>`. */
  minInsured: 'Min Insured',
  /** `InboxProductName.xml` b37818 `<pyValue>`. */
  maxInsured: 'Max Insured',
  /** `InboxProductName.xml` b37966 `<pyValue>`. */
  employee: 'Employee',
  /** `InboxProductName.xml` b38115 `<pyValue>`. */
  nonEmployee: 'Non-Employee',
  /** `InboxProductName.xml` b38425 `<pyLabel>`. */
  add: 'Add',
  /** `InboxProductName.xml` b40024 `<pyLabel>`. */
  delete: 'Delete',
} as const

/** Grid `UNDERWRITING LIMIT` - `ProductName.UnderwritingLimitList`. */
export const UWLIMIT_MPNL = {
  /** `InboxProductName.xml` b42075 `<pyValue>`. */
  judul: 'UNDERWRITING LIMIT',
  /** `InboxProductName.xml` b42594 `<pyValue>`. */
  minInsured: 'Min Insured',
  /** `InboxProductName.xml` b42742 `<pyValue>`. */
  maxInsured: 'Max Insured',
  /** `InboxProductName.xml` b42890 `<pyValue>`. */
  minAge: 'Min Age',
  /** `InboxProductName.xml` b43038 `<pyValue>`. */
  maxAge: 'Max Age',
  /** `InboxProductName.xml` b43187 `<pyValue>`. */
  medical: 'Medical',
  /** `InboxProductName.xml` b43336 `<pyValue>`. */
  description: 'Description',
  /** `InboxProductName.xml` b43643 `<pyLabel>`. */
  add: 'Add',
  /** `InboxProductName.xml` b45682 `<pyLabel>`. */
  delete: 'Delete',
} as const

/** Grid komentar - `ProductName.CommentList` (wadah b61044, baca-saja). */
export const KOMENTAR_MPNL = {
  /** `InboxProductName.xml` b62095 `<pyValue>`. */
  date: 'Date',
  /** `InboxProductName.xml` b62246 `<pyValue>`. */
  pic: 'PIC',
  /** `InboxProductName.xml` b62399 `<pyValue>`. */
  comment: 'Comment',
} as const

/** Tombol bawah form (wadah b57867). */
export const TOMBOL_MPNL = {
  /** `InboxProductName.xml` b58770 `<pyLabel>`. */
  close: 'Close',
  /** `InboxProductName.xml` b59041 `<pyLabel>`. */
  save: 'Save',
  /** `InboxProductName.xml` b59489 `<pyLabel>`. */
  edit: 'Edit',
  /** `InboxProductName.xml` b59854 `<pyLabel>`. */
  copy: 'Copy',
  /** `InboxProductName.xml` b60122 `<pyLabel>`. */
  generate: 'Generate',
} as const

/** Dialog `Save` - FlowAction + section `SaveProductName_Confirm`. */
export const SIMPAN_MPNL = {
  /** `SaveProductName_Confirm.xml` b519 `<pyValue>`. */
  tanya: 'Do you want to save the data?',
  /** `SaveProductName_Confirm.xml` b1025 `<pyLabelFieldValue>`. */
  comment: 'Comment',
  /** `SaveProductName_Confirm.xml` b34 `<pySubmitLabel>`. */
  save: 'Save',
  /** `SaveProductName_Confirm.xml` b33 `<pyCancelLabel>`. */
  cancel: 'Cancel',
} as const

/** Dialog `Edit` - FlowAction + section `EditProductName_Confirm` (`SetViewEdit`). */
export const EDIT_MPNL = {
  /** `EditProductName_Confirm.xml` b496 `<pyValue>`. */
  tanya: 'Do you want to Edit the data?',
  /** `EditProductName_Confirm.xml` b19 `<pySubmitLabel>`. */
  edit: 'Edit',
  /** `EditProductName_Confirm.xml` b18 `<pyCancelLabel>`. */
  cancel: 'Cancel',
} as const

/**
 * Tujuh pemilih master `Choose*` - teks sama di ketujuh section (`Ceding_Section` sebagai bukti; uji memeriksa ketujuhnya).
 * Sejak keputusan work owner 02-10-2026 pemilihnya dropdown (`DropdownMaster`): `search` = isian saring, `kolom*` =
 * kepala kolom daftar; `choose` / `submit` / `cancel` - dan ketujuh tombol `choose*` di objek medan - tetap berbukti
 * korpus tetapi tidak dirender.
 */
export const PEMILIH_MPNL = {
  /** `Ceding_Section.xml` b513 `<pyLabelFieldValue>`. */
  search: 'Search',
  /** `Ceding_Section.xml` b1502 `<pyValue>`. */
  kolomId: 'ID',
  /** `Ceding_Section.xml` b1643 `<pyValue>`. */
  kolomName: 'Name',
  /** `Ceding_Section.xml` b2265 `<pyLabel>`. */
  choose: 'Choose',
  /** `ChooseCeding.xml` b32 `<pySubmitLabel>`. */
  submit: 'Submit',
  /** `ChooseCeding.xml` b31 `<pyCancelLabel>`. */
  cancel: 'Cancel',
  /** `RIRate_Section.xml` b1739 `<pyValue>`. */
  kolomRiRateName: 'RIRate Name',
} as const

/** Dialog `View Rate` - FlowAction + section `ViewRate` (view `RATE_LIFE`, K1 01-10-2026). */
export const RATE_MPNL = {
  /** `ViewRate.xml` b843 `<pyValue>`. */
  judul: 'Outward List',
  /** `ViewRate.xml` b1147 `<pyValue>`. */
  kolomId: 'ID',
  /** `ViewRate.xml` b1293 `<pyValue>`. */
  usedby: 'USEDBY',
  /** `ViewRate.xml` b1439 `<pyValue>`. */
  gender: 'GENDER',
  /** `ViewRate.xml` b1585 `<pyValue>`. */
  contract: 'CONTRACT',
  /** `ViewRate.xml` b1731 `<pyValue>`. */
  age: 'AGE',
  /** `ViewRate.xml` b1877 `<pyValue>`. */
  rate: 'RATE',
  /** `ViewRate.xml` b18 `<pySubmitLabel>`. */
  submit: 'Submit',
  /** `ViewRate.xml` b20 `<pyCancelLabel>`. */
  cancel: 'Cancel',
} as const

/** Panel lampiran - `InboxProductName` wadah b64133 + FlowAction `ProductNameAttachContent`. */
export const LAMPIRAN_MPNL = {
  /** `InboxProductName.xml` b64747 `<pyLabel>`. */
  add: 'Add attachment',
  /** `InboxProductName.xml` b65270 `<pyLabel>`. */
  refresh: 'Refresh',
  /** `InboxProductName.xml` b66071 `<pyValue>`. */
  peringatanNama: 'Make sure the file name doesn\'t contain forbidden character such as , / | \' "',
  /** `InboxProductName.xml` b66230 `<pyValue>`. */
  peringatanGanti: 'Recommended safe substitute should be . or _',
  /** `InboxProductName.xml` b67657 `<pyLabel>`. */
  downloadAll: 'Download All',
  /** `InboxProductName.xml` b68426 `<pyValue>`. */
  fileName: 'File Name',
  /** `InboxProductName.xml` b69291 `<pyLabel>`. */
  viewOffice: 'View Office Online',
  /** `InboxProductName.xml` b69714 `<pyLabel>`. */
  delete: 'Delete',
  /** `ProductNameAttachContent.xml` b24 `<pySubmitLabel>`. */
  attach: 'Attach',
  /** `ProductNameAttachContent.xml` b22 `<pyCancelLabel>`. */
  cancel: 'Cancel',
} as const

/** Pilihan `Premium Payment Method` b25611 - daftarnya `associated` (tidak ikut ekspor); teks dan kodenya dari
 *  ekspresi `@if` `GenerateUpload_Act.xml` b1141 (`CARI37`). Kode `Single` tidak diketahui (OQ-MPNL-05). */
export const PEMBAYARAN_MPNL = {
  /** `GenerateUpload_Act.xml` b1141 - `PAYMENT==1`. */
  annual: 'Annual',
  /** `GenerateUpload_Act.xml` b1141 - `PAYMENT==2`. */
  semiAnnual: 'Semi Annual',
  /** `GenerateUpload_Act.xml` b1141 - `PAYMENT==3`. */
  quarterly: 'Quarterly',
  /** `GenerateUpload_Act.xml` b1141 - `PAYMENT==4`. */
  monthly: 'Monthly',
} as const

/** Pesan layar dari aturan, bukan dari section. */
export const PESAN_MPNL = {
  /** `DataTransform/CopyProduct.xml` b236 `<pyPropertiesValue>` (berkutip) - `OutputParam.ERRMSG`, wadah b1269 `STSSAVE==99`. */
  copy: 'Data sudah dicopy, silakan melakukan perubahan dan tekan SAVE untuk menyimpan',
} as const

/** Teks yang tidak ada di korpus - masing-masing beralasan. */
export const LAIN_MPNL = {
  /** `[tidak ada di korpus]` - grid tanpa baris (ADR-U-0027: kosong dinyatakan). */
  kosong: 'No items',
  /**
   * `[tidak ada di korpus]` - dialog `View Rate` dipotong `BrowseRateLife_RD` `pyMaxRecords` 500; Pega memotong
   * diam-diam, di sini potongan DINYATAKAN (code review 01-10-2026).
   */
  terpotong: 'Only the first 500 rows are shown.',
  /** `[tidak ada di korpus]` - ikon grid bawaan `pzPegaDefaultGridIcons` (`LIEN CLAUSE` b12373, `DOCUMENT CLAIM` b14760): tambah baris. */
  tambahBaris: 'Add row',
  /** `[tidak ada di korpus]` - ikon grid bawaan `pzPegaDefaultGridIcons`: hapus baris. */
  hapusBaris: 'Delete row',
  /** `[tidak ada di korpus]` - ikon `pxIcon` tanpa tooltip b39697 (`CopyFinancialWriting`) / b45355 (`CopyUnderWritingLimit`). */
  salinBaris: 'Copy row',
  /** `[tidak ada di korpus]` - status lampiran tiket 08 AC 4 (`terunggah`). */
  terunggah: 'Uploaded',
  /** `[tidak ada di korpus]` - status lampiran tiket 08 AC 4 (`gagal`). */
  gagal: 'Failed',
  /** `[tidak ada di korpus]` - status lampiran tiket 08 AC 4 (`belum`). */
  belum: 'Pending',
  /** `[tidak ada di korpus]` - kirim ulang lampiran gagal, tiket 08 AC 3 / tiket 09 AC 4 (`POST …/ulangi`). */
  ulangi: 'Retry',
  /**
   * `[tidak ada di korpus]` - dropdown master (pengganti tombol `Choose*`, keputusan work owner 02-10-2026) memuat
   * paling banyak `BATAS_DROPDOWN` baris; potongan DINYATAKAN, sisanya dicapai lewat `Search`.
   */
  dropdownTerpotong: 'Only the first 200 rows are shown. Type in Search to narrow the list.',
  /** `[tidak ada di korpus]` - tombol di samping `Add` (permintaan work owner 03-10-2026), teks dari permintaannya. */
  copyOld: 'Copy Old',
  /** `[tidak ada di korpus]` - keterangan popup `Copy Old`. */
  copyOldKeterangan:
    'Products in the old (JSON) tables that are not in the new tables yet. Tick the ones to copy, then press Process Copy.',
  /** `[tidak ada di korpus]` - tombol kaki popup `Copy Old` (permintaan work owner 03-10-2026). */
  prosesCopy: 'Process Copy',
  /** `[tidak ada di korpus]` - kotak centang kepala kolom popup `Copy Old` (pembaca layar). */
  pilihSemua: 'Select all',
  /** `[tidak ada di korpus]` - kotak centang satu baris popup `Copy Old` (pembaca layar), diikuti ID produk. */
  pilihBaris: 'Select',
  /** `[tidak ada di korpus]` - cacah baris tercentang popup `Copy Old`, didahului angkanya. */
  dipilih: 'selected',
  /** `[tidak ada di korpus]` - kolom alasan/catatan popup `Copy Old`. */
  kolomCatatan: 'Notes',
  /** `[tidak ada di korpus]` - popup `Copy Old` tanpa baris. */
  copyOldKosong: 'All old products are already in the new tables.',
  /** `[tidak ada di korpus]` - status hasil `Process Copy`. */
  statusDisalin: 'Copied',
  /** `[tidak ada di korpus]` - status hasil `Process Copy`. */
  statusSudahAda: 'Already in the new tables',
  /** `[tidak ada di korpus]` - status hasil `Process Copy`. */
  statusDitolak: 'Cannot be copied',
  /** `[tidak ada di korpus]` - status hasil `Process Copy`. */
  statusGagal: 'Failed',
  /** `[tidak ada di korpus]` - kotak unggah banyak berkas + seret-lepas `Add attachment` (permintaan work owner 03-10-2026). */
  seretBerkas: 'Drag and drop files here, or click to choose files',
  /** `[tidak ada di korpus]` - kemajuan unggah banyak berkas, diikuti urutan dan nama berkas. */
  mengunggah: 'Uploading',
  /** `[tidak ada di korpus]` - membuang satu berkas dari pilihan unggah (belum diunggah). */
  buangPilihan: 'Remove',
} as const
