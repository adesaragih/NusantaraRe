// Label modul Master Contract Retro Life - paket 9.
//
// ⛔ VERBATIM dari korpus `D:\XML\RNM_BRD\Master Contract Retro Life\`, dengan berkas + baris di
// sebelahnya (berkas korpus satu tag per baris; nomor baris mentah = nomor `sed -e 's/></>\n</g'`).
// Nomor yang disebut adalah baris TAG yang memuat teksnya (bukan baris awal sel), dan
// `labels.test.ts` membuka korpus di baris itu. Teks tombol = `pyModes.pyLabel` (keterangan yang
// dirender kontrol `pxButton`), bukan `pyLabelPreview`.
//
// Label korpus sudah berbahasa Inggris. Pesan server VERBATIM korpus (`Data Berhasil di Hapus`,
// `Copied to all reins types.`, ...) tampil apa adanya dari jawaban server - tidak diterjemahkan.
//
// Teks yang TIDAK ada di korpus ditandai `[tidak ada di korpus]` beserta alasannya (penyimpangan
// sadar atau AC tiket); dijaga `labels.test.ts` supaya tidak ada yang menyamar sebagai label XML.

/** Menu - `M_NAV_MENU.LABEL` isi awal 900 = nama folder korpus. */
export const MENU_MCRL = {
  /** Nama folder korpus `D:\XML\RNM_BRD\Master Contract Retro Life`. */
  kelompok: 'Master Contract Retro Life',
} as const

/** Kosakata bersama beberapa panel. */
export const UMUM_MCRL = {
  /**
   * `[tidak ada di korpus]` - keterangan ikon `pxIconCancel` kerangka harness popup
   * (`Harness/InboxRetroLimitReinsurers.xml` b694, `InboxRetroLifeReinsurersList.xml` b692,
   * `InboxSecurityReinsurerLife.xml` b697, `InboxBusinessLifeReinsurers.xml` b691): kontrol bawaan
   * Pega tanpa teks yang menutup harness popup.
   */
  tutup: 'Close',
  /** `[tidak ada di korpus]` - grid tanpa baris (ADR-U-0027: kosong dinyatakan). */
  kosong: 'No items',
} as const

/** Halaman awal - `Section/GridRetrocessionLife.xml` + `InputRetrocessionLife.xml` + `InputDtlRetrocessionLife.xml`. */
export const TAHUN_MCRL = {
  /** `GridRetrocessionLife.xml` b1082 `<pyValue>` (sel b1017, format `Heading 1`). */
  judul: 'MASTER CONTRACT RETRO LIFE',
  /** `InputRetrocessionLife.xml` b8927 `<pyLabelFieldValue>` - label SEL tombol b8888 (`pyIncludeLabel` true). */
  labelSelAdd: 'End Period',
  /** `InputRetrocessionLife.xml` b9007 `<pyLabel>` - teks tombol b8888 (`NewInputTreatyYear_Life_Act`). */
  add: 'Add',
  /** `InputRetrocessionLife.xml` b9005 `<pyTooltip>` tombol b8888. */
  tooltipAdd: 'Add New Data',
  /** `InputRetrocessionLife.xml` b9588 `<pyValue>`. */
  kolomId: 'ID',
  /** `InputRetrocessionLife.xml` b9750 `<pyValue>` → `.UNDERWRITINGYEAR`. */
  kolomUnderwritingYear: 'UNDERWRITING YEAR',
  /** `InputRetrocessionLife.xml` b9922 `<pyValue>` → `.TREATYYEAR`. */
  kolomTransactionYear: 'TRANSACTION YEAR',
  /** `InputRetrocessionLife.xml` b10094 `<pyValue>`. */
  kolomStartDate: 'START DATE',
  /** `InputRetrocessionLife.xml` b10266 `<pyValue>`. */
  kolomEndDate: 'END DATE',
  /** `InputRetrocessionLife.xml` b11774 `<pyLabel>` - tombol b11645 (`SetTreatyYearLife_Act`). */
  edit: 'Edit',
  /** `InputRetrocessionLife.xml` b11772 `<pyTooltip>` tombol b11645. */
  tooltipEdit: 'Edit Data',
  /** `InputRetrocessionLife.xml` b12117 `<pyLabel>` - tombol b11988 (`showHarness` `InboxRetroLimitReinsurers`). */
  reinsType: 'ReinsType',
  /** `InputDtlRetrocessionLife.xml` b892 `<pyValue>`. */
  inputNewData: 'Input New Data',
  /** `InputDtlRetrocessionLife.xml` b1834 `<pyLabelFieldValue>` (ro, tampil bila `NOTBLANK`). */
  formId: 'ID',
  /** `InputDtlRetrocessionLife.xml` b2018 `<pyLabelFieldValue>`. */
  formUnderwritingYear: 'UNDERWRITING YEAR',
  /** `InputDtlRetrocessionLife.xml` b2277 `<pyLabelFieldValue>`. */
  formTransactionYear: 'TRANSACTION YEAR',
  /** `InputDtlRetrocessionLife.xml` b2522 `<pyLabelFieldValue>`. */
  formStartDate: 'START DATE',
  /** `InputDtlRetrocessionLife.xml` b2812 `<pyLabelFieldValue>`. */
  formEndDate: 'END DATE',
  /** `InputDtlRetrocessionLife.xml` b3633 `<pyLabelFieldValue>` (ro). */
  formModifiedDate: 'Modified Date',
  /** `InputDtlRetrocessionLife.xml` b3819 `<pyLabelFieldValue>` (ro). */
  formInputor: 'Inputor',
  /** `InputDtlRetrocessionLife.xml` b5059 `<pyLabel>` - tombol b4944 (`SaveTreatyYearLife_Act`). */
  save: 'Save',
  /** `InputDtlRetrocessionLife.xml` b5336 `<pyLabel>` - tombol b5221 (`CancelActivityTreatyContract`). */
  cancel: 'Cancel',
} as const

/** Panel kontrak - `Section/InputRetroLimitReinsurers.xml` (harness `InboxRetroLimitReinsurers`). */
export const KONTRAK_MCRL = {
  /** b666 `<pyValue>`. */
  judul: 'Reins Type',
  /** b1093 `<pyLabelFieldValue>` (ro `InputTreatyContract.IDTREATYYEAR`). */
  idTreatyYear: 'ID Treaty Year',
  /** b2108 `<pyLabelFieldValue>` (ro). */
  formInputor: 'Inputor',
  /** b2537 `<pyLabelFieldValue>` (ro). */
  formModifiedDate: 'Modified Date',
  /** b3278 `<pyLabelFieldValue>` - dropdown `REINSTYPEID`, tampil `.Note`. */
  formReinsType: 'REINS TYPE',
  /** b3642 `<pyLabelFieldValue>` (ro). */
  formTreatyStart: 'TREATY START',
  /** b3932 `<pyLabelFieldValue>` (ro). */
  formTreatyEnd: 'TREATY END',
  /** b4222 `<pyLabelFieldValue>` → `B_IDR`. */
  formMinIdr: 'MINIMUM LIMIT (IDR)',
  /** b4486 `<pyLabelFieldValue>` → `IDR`. */
  formMaxIdr: 'MAXIMUM LIMIT (IDR)',
  /** b4709 `<pyLabelFieldValue>` → `B_USD`. */
  formMinUsd: 'MINIMUM LIMIT (USD)',
  /** b4974 `<pyLabelFieldValue>` → `USD`. */
  formMaxUsd: 'MAXIMUM LIMIT (USD)',
  /** b5889 `<pyLabel>` - tombol b5773 (`SaveTreatyLimit_Act`). */
  save: 'Save',
  /** b6176 `<pyLabel>` - tombol b6066 (`CancelActivityTreatyContract`). */
  cancel: 'Cancel',
  /** b8711 `<pyLabel>` - tombol b8596 (`NewInputTreatyLimit_Life`). */
  add: 'Add',
  /** b9370 `<pyValue>`. */
  kolomId: 'ID',
  /** b9519 `<pyValue>`. */
  kolomReinsType: 'REINS TYPE',
  /** b9681 `<pyValue>`. */
  kolomTreatyStart: 'TREATY START',
  /** b9843 `<pyValue>`. */
  kolomTreatyEnd: 'TREATY END',
  /** b10009 `<pyValue>`. */
  kolomMinIdr: 'MINIMUM LIMIT (IDR)',
  /** b10175 `<pyValue>`. */
  kolomMaxIdr: 'MAXIMUM LIMIT (IDR)',
  /** b10341 `<pyValue>`. */
  kolomMinUsd: 'MINIMUM LIMIT (USD)',
  /** b10507 `<pyValue>`. */
  kolomMaxUsd: 'MAXIMUM LIMIT (USD)',
  /** b12738 `<pyLabel>` - tombol b12621 (`SetRetroListLife_Act`). */
  edit: 'Edit',
  /** b13113 `<pyLabel>` - tombol b12999 (`showHarness` `InboxBusinessLifeReinsurers`). */
  businessList: 'Business List',
  /** b14148 `<pyLabel>` - tombol b14034 (`showHarness` `InboxRetroLifeReinsurersList`). */
  reinsurerList: 'Reinsurer List',
  /** b15275 `<pyLabel>` - tombol b15161 (`DeleteTreatyLimit_Act`). */
  delete: 'Delete',
} as const

/** Panel reinsurer - `Section/InputSecurityLifeReinsurers.xml` (harness `InboxRetroLifeReinsurersList`). */
export const REINSURER_MCRL = {
  /** b637 `<pyValue>`. */
  judul: 'Reinsurer List',
  /** b1066 `<pyLabelFieldValue>`. */
  idTreatyYear: 'ID Treaty Year',
  /** b1280 `<pyLabelFieldValue>` - ⚠️ nilainya `InputBusinessLife.TREATYCONTRACTID` (ID kontrak). */
  idReinsType: 'ID Reins Type',
  /** b1483 `<pyLabelFieldValue>`. */
  reinsType: 'Reins Type',
  /** b2539 `<pyLabelFieldValue>` (ro). */
  formInputor: 'Inputor',
  /** b3782 `<pyLabelFieldValue>` - autocomplete `BrowseCedingCoLife_RD`. */
  formReinsurerName: 'REINSURER NAME',
  /** b4821 `<pyLabelFieldValue>`. */
  formShare: '(%) SHARE',
  /** b5111 `<pyLabelFieldValue>` → kolom `COMMISION`. */
  formDiscount: '(%) DISCOUNT',
  /** b5403 `<pyLabelFieldValue>`. */
  formOvrComm: '(%) OVR COMM',
  /** b5770 `<pyLabel>` - tombol b5656 (`SaveSecurityLife_Act`). */
  save: 'Save',
  /** b6050 `<pyLabel>` - tombol b5942 (`CancelActivityTreatyContract`). */
  cancel: 'Cancel',
  /** b8797 `<pyLabel>` - tombol b8689 (`NewInputSecurityLife_Act`). */
  add: 'Add',
  /** b9481 `<pyValue>`. */
  kolomId: 'ID',
  /** b9625 `<pyValue>`. */
  kolomReinsurerName: 'REINSURER NAME',
  /** b9781 `<pyValue>`. */
  kolomShare: '(%) SHARE',
  /** b9939 `<pyValue>`. */
  kolomDiscount: '(%) DISCOUNT',
  /** b10097 `<pyValue>`. */
  kolomOvrComm: '(%) OVR COMM',
  /** b10253 `<pyValue>`. */
  kolomInputor: 'INPUTOR',
  /** b10408 `<pyValue>`. */
  kolomUpdateDate: 'UPDATE DATE',
  /** b12216 `<pyLabel>` - tombol b12099 (`SetSecurityLife_Act`). */
  edit: 'Edit',
  /** b12555 `<pyLabel>` - tombol b12441 (`showHarness` `InboxSecurityReinsurerLife`). */
  securityReinsurer: 'Security Reinsurer',
  /** b13532 `<pyLabel>` - tombol b13422 (`DeleteSecurityLife_Act`). */
  delete: 'Delete',
  /** b14404 `<pyValue>` (`Total Share --&gt;&gt;`). */
  totalShare: 'Total Share -->>',
  /** `[tidak ada di korpus]` - tiket 05 AC 19: kontrak dengan total ≠ 100 ditandai mencolok; tidak memblokir. */
  totalBukan100: 'Total share is not 100%',
} as const

/** Panel security - `Section/InputSecurityReinsurerLife.xml` (harness `InboxSecurityReinsurerLife`). */
export const SECURITY_MCRL = {
  /** b659 `<pyValue>`. */
  judul: 'Security Reinsurer',
  /** b1088 `<pyLabelFieldValue>`. */
  idReinsurer: 'ID Reinsurer',
  /** b1306 `<pyLabelFieldValue>`. */
  reinsurerName: 'Reinsurer Name',
  /** b1522 `<pyLabelFieldValue>`. */
  pctShare: 'PCT Share (%)',
  /** b2603 `<pyLabelFieldValue>` (ro). */
  formInputor: 'Inputor',
  /** b3856 `<pyLabelFieldValue>` - autocomplete `BrowseCedingCoLife_RD`. */
  formSecurityReinsurerName: 'SECURITY REINSURER NAME',
  /** b4898 `<pyLabelFieldValue>` - persen DARI share induk (tiket 06). */
  formShare: '(%) SHARE',
  /** b5266 `<pyLabel>` - tombol b5151 (`SaveSecurityReinsurerLife_Act`). */
  save: 'Save',
  /** b5537 `<pyLabel>` - tombol b5432 (`CancelActivityTreatyContract`). */
  cancel: 'Cancel',
  /** b8287 `<pyLabel>` - tombol b8179 (form baru KOSONG - OQ-MCRL-12). */
  add: 'Add',
  /** b8978 `<pyValue>`. */
  kolomId: 'ID',
  /** b9122 `<pyValue>`. */
  kolomReinsurerName: 'REINSURER NAME',
  /** b9278 `<pyValue>`. */
  kolomShare: '(%) SHARE',
  /**
   * `[tidak ada di korpus]` - tiket 06 AC 23/`Layar menampilkan share mentah dan eksposur efektif
   * berdampingan, dengan label yang membedakan keduanya` (OQ-MCRL-08).
   */
  kolomEksposur: 'EXPOSURE TO TREATY (%)',
  /** b9434 `<pyValue>`. */
  kolomInputor: 'INPUTOR',
  /** b9589 `<pyValue>`. */
  kolomUpdateDate: 'UPDATE DATE',
  /** b10893 `<pyLabel>` - tombol b10779 (`SetSecurityReinsurerLife_Act`). */
  edit: 'Edit',
  /** b11214 `<pyLabel>` - tombol b11097 (`DeleteSecurityReinsurerLife_Act`). */
  delete: 'Delete',
} as const

/** Panel business - `Section/InputBusinessLifeReinsurers.xml` (harness `InboxBusinessLifeReinsurers`). */
export const BUSINESS_MCRL = {
  /** b684 `<pyValue>`. */
  judul: 'Business List',
  /** b1103 `<pyLabelFieldValue>`. */
  idTreatyYear: 'ID Treaty Year',
  /** b1320 `<pyLabelFieldValue>` - ⚠️ nilainya ID kontrak. */
  idReinsType: 'ID Reins Type',
  /** b1535 `<pyLabelFieldValue>`. */
  reinsType: 'Reins Type',
  /** b2573 `<pyLabelFieldValue>` (ro). */
  formInputor: 'Inputor',
  /** b2982 `<pyLabelFieldValue>` (ro). */
  formModifiedDate: 'Modified Date',
  /** b3712 `<pyLabelFieldValue>` (ro, diisi autocomplete `BUSINESS NAME`). */
  formBusinessCode: 'BUSINESS CODE',
  /** b3898 `<pyLabelFieldValue>` - autocomplete `BrowseBusinessLife_RD`. */
  formBusinessName: 'BUSINESS NAME',
  /** b4351 `<pyLabelFieldValue>` - autocomplete `BrowseRateLifeSummary` (view `RATE_LIFE_SUMMARY`, K1). */
  formRiRate: 'R/I RATE',
  /** b4849 `<pyLabel>` - tombol b4732, tampil bila `RIRATEID != ''` (`localAction ViewRate`). */
  viewRateForm: 'View Rate',
  /** b5846 `<pyLabel>` - tombol b5739 (`SaveBusinessLife_Act`, aksi `pyBehaviors` lama). */
  save: 'Save',
  /** b6064 `<pyLabel>` - tombol b5954 (`CancelActivityTreatyContract`). */
  cancel: 'Cancel',
  /** b8599 `<pyLabel>` - tombol b8484 (`NewInputBusinessLife_Act`). */
  add: 'Add',
  /** b9128 `<pyValue>`. */
  kolomBusinessName: 'BUSINESS NAME',
  /** b9286 `<pyValue>`. */
  kolomRiRate: 'R/I RATE',
  /** b9440 `<pyValue>`. */
  kolomInputor: 'INPUTOR',
  /** b9595 `<pyValue>`. */
  kolomUpdateDate: 'UPDATE DATE',
  /** b10938 `<pyLabel>` - tombol b10824 (`SetBusinessListLife_Act`). */
  edit: 'Edit',
  /** b11264 `<pyLabel>` - tombol b11154 (`DeleteRowBusiness`). */
  delete: 'Delete',
  /** b11581 `<pyLabel>` - tombol baris b11469 (`SetParamRateTable` + `localAction ViewRateTable`). */
  viewRate: 'View Rate',
  /** b12135 `<pyLabel>` - tombol b12020 (`SaveBusinessToAllLife_Act`). */
  copyToAll: 'Copy to all Reinstype',
} as const

/** Section `ViewRate` - `Section/ViewRate.xml` (FlowAction `ViewRate` / `ViewRateTable`). */
export const RATE_MCRL = {
  /** b861 `<pyValue>`. */
  judul: 'Rate List',
  /** b1178 `<pyValue>`. */
  kolomId: 'ID',
  /** b1324 `<pyValue>`. */
  kolomUsedBy: 'USEDBY',
  /** b1470 `<pyValue>`. */
  kolomGender: 'GENDER',
  /** b1616 `<pyValue>`. */
  kolomContract: 'CONTRACT',
  /** b1762 `<pyValue>`. */
  kolomAge: 'AGE',
  /** b1908 `<pyValue>`. */
  kolomRate: 'RATE',
  /**
   * `FlowAction/ViewRate.xml` b19 `<pySubmitLabel>` (`ViewRateTable.xml` b21 sama) - tombol dialog
   * `pyShowFAButtons` true; tanpa post-processing → menutup dialog.
   */
  submit: 'Submit',
  /** `FlowAction/ViewRate.xml` b20 `<pyCancelLabel>` (`ViewRateTable.xml` b22 sama). */
  cancel: 'Cancel',
} as const

/**
 * Popup konfirmasi hapus - `[tidak ada di korpus]` seluruhnya: Pega menghapus TANPA konfirmasi;
 * popup = penyimpangan sadar 3 (tiket 09 AC 32-36). Judul memakai teks tombol `Delete`.
 */
export const HAPUS_MCRL = {
  /** `[tidak ada di korpus]`. */
  pertanyaan: 'Delete this row?',
  /** `[tidak ada di korpus]` - AC 32/33: apa dan berapa yang ikut terhapus. */
  ikutTerhapus: 'Deleted together with it:',
  /** `[tidak ada di korpus]` - AC 36: baris tanpa anak tetap dikonfirmasi, dengan pesan tanpa anak. */
  tanpaAnak: 'No other row is affected.',
  /** `[tidak ada di korpus]`. */
  reinsurer: 'reinsurer row(s)',
  /** `[tidak ada di korpus]`. */
  security: 'security reinsurer row(s)',
  /** `[tidak ada di korpus]`. */
  business: 'business row(s)',
  /** `[tidak ada di korpus]` - dampak dihitung sesaat sebelum dialog tampil. */
  memuat: 'Counting the rows that will be deleted…',
  /** `[tidak ada di korpus]`. */
  ya: 'Yes',
} as const

/**
 * Pratinjau `Copy to all Reinstype` - `[tidak ada di korpus]` seluruhnya: penyimpangan sadar 5
 * (tiket 08 AC 28, 30, dan pratinjau nol baris). Judul memakai teks tombol `Copy to all Reinstype`.
 */
export const SALIN_MCRL = {
  /** `[tidak ada di korpus]` - AC: pratinjau menyebut jenis reasuransi dasar pemilihan. */
  dasar: 'Source reins type (skipped):',
  /** `[tidak ada di korpus]` - AC 28: jumlah baris yang akan terpengaruh. */
  akanDitulis: 'New business rows that will be written:',
  /** `[tidak ada di korpus]` - pratinjau nol baris: diberi tahu, tanpa tombol konfirmasi. */
  nol: 'There is no other reins type in this treaty year - nothing will be copied.',
  /** `[tidak ada di korpus]`. */
  memuat: 'Loading the other reins types…',
  /** `[tidak ada di korpus]`. */
  ya: 'Yes',
} as const
