// Label layar Endorsement Life - VERBATIM korpus `D:\XML\RNM_BRD\Endorsement Life\`.
//
// ⛔ Setiap kunci punya bukti baris XML di `labels.test.ts` (dibuka dan dicocokkan), ATAU terdaftar
// `BUKAN_KORPUS` di sana dengan alasannya. Label yang dikarang tanpa salah satunya menggagalkan uji.
// Nomor baris = `sed -e 's/></>\n</g' <berkas> | grep -n '<teks>'` (kepala `docs/PARITAS-LAYAR-DAN-AKSI.md`).

/** Menu - nama FOLDER korpus = `M_NAV_MENU.LABEL` baris `endorsementlife` (900). */
export const MENU_EDM = {
  kelompok: 'Endorsement Life',
} as const

/** Teks bawaan modul tanpa padanan korpus (bahasa Inggris, seperti label korpusnya). */
export const UMUM_EDM = {
  kosong: 'No data.',
  tutup: 'Close',
  kembali: 'Back',
  memuat: 'Loading…',
  memeriksa: 'Checking…',
  rinci: 'Details',
  terkunci: 'Locked: the policy and the endorsement type cannot change once the case is created. Decline the case and create a new one.',
  menyimpan: 'Saving…',
  tandaiHapus: 'Mark for deletion',
  mengunggah: 'Uploading…',
  pilihBerkas: 'CSV file',
  barisDibaca: 'Rows read',
  barisDitolak: 'Rows rejected',
  barisDisimpan: 'Rows added',
  kolomDiabaikan: 'Ignored columns (not saved)',
  pesanTerpotong: 'Only the first rejections are listed.',
  kolomBaris: 'Row',
  kolomKolom: 'Column',
  kolomPesan: 'Message',
  memutuskan: 'Submitting…',
} as const

/** `Section/InboxEndorsementLife.xml` - halaman awal (harness portal `InboxEndorsementLife`). */
export const INBOX_EDM = {
  judul: 'Endorsement Life', // b5505
  createAddendum: 'Create Addendum', // b6620
  kolomCaseId: 'Case ID', // b8418
  kolomEndorsementNo: 'Endorsement No', // b8556
  kolomType: 'Type', // b8692
  kolomEdmType: 'EDM Type', // b8828
  kolomPolicyNo: 'Policy No', // b8964
  kolomSob: 'SOB', // b9100
  kolomCeding: 'Ceding', // b9236
  kolomPolicyHolder: 'Policy Holder', // b9372
  kolomMarketingName: 'Marketing Name', // b9508
  kolomCreateDate: 'Create Date', // b9644
  kolomCreateOperator: 'Create Operator', // b9780
  kolomEdmTypeBatal: 'EDM Type Batal', // b9916
  kolomStatus: 'Status', // b10052
} as const

/** `Section/InputEDMLife.xml` - kepala kasus (seluruhnya baca-saja). */
export const KASUS_EDM = {
  judul: 'Endorsement Life Detail', // b817
  productName: 'Product Name', // b3399
  productNameId: 'Product Name ID', // b3589
  type: 'Type', // b3812
  reinsuranceSystem: 'Reinsurance System', // b4143
  classOfBusiness: 'Class of Business', // b4473
  sob: 'SOB', // b4658
  policyHolder: 'Policy Holder', // b4897
  premiumMethod: 'Premium Method', // b5976
  ceding: 'Ceding', // b6339
  marketingOfficer: 'Marketing Officer', // b6578
  edmType: 'EDM Type', // b6910
  description: 'Description', // b7697
} as const

/** `Section/EndorsmentLife_Section.xml` - label yang `InputEDMLife` ambil dari deskripsi properti. */
export const BUAT_EDM = {
  policyNo: 'Policy No', // b1075 (`InputEDMLife.xml` b3126 memakai deskripsi properti yang sama)
  edmType: 'EDM Type', // b1347
  description: 'Description', // b2140
  edmDate: 'EDM Date', // b3160
  submit: 'Submit', // b4226
} as const

/** `Section/ShowLifePremiumSummary_EDM.xml` - tombol yang dipindah ke `InputEDMLife` (RALAT R02). */
export const POLIS_LAMA_EDM = {
  viewOldPolicy: 'View Old Policy', // b64965 (QR), b65522 (QP), b66083 (TR), b66640 (TP)
} as const

/** Grid peserta `InputEDMLife.xml` b11899 (EdmType 1) / b17500 (EdmType 3). */
export const GRID_EDM = {
  policyNo: 'POLICY NO', // b12035
  policyHolder: 'POLICY HOLDER', // b12177
  certificateNo: 'CERTIFICATE NO', // b12320
  nameOfInsured: 'NAME OF INSURED', // b12463
  sex: 'SEX', // b12602
  dateOfBirth: 'DATE OF BIRTH', // b12738
  entryAge: 'ENTRY AGE', // b12874
  plan: 'PLAN', // b13010
  beginDate: 'BEGIN DATE', // b13146
  effectiveDate: 'EFFECTIVE DATE', // b13282
  expiredDate: 'EXPIRED DATE', // b13417
} as const

/** Tombol simpan `InputEDMLife.xml` - kotak centang `.EdmBatal` b15753 tanpa teks di korpus. */
export const SIMPAN_EDM = {
  deleteAll: 'DELETE ALL', // b13607 → `SelectAllEdmLife_act`
  save: 'Save', // b37202 → `SetPremi_EDM`
} as const

/** Unggah CSV - wadah `InputEDMLife.xml` b8698, popup `ViewCSVResult_LifeEDM`, modal `UploadCSV_LifeEndorsement`. */
export const UNGGAH_EDM = {
  uploadCsv: 'Upload CSV', // b8973 → `UploadCSV_LifeEndorsement`
  viewUpload: 'View Upload', // b9340 → `ViewCSVResult_LifeEDM`
  addCsvData: 'Add CSV Data', // b10405 → `SaveCSVEDMLife`
  judulModal: 'Endorsement Life - Upload CSV', // `FlowAction/UploadCSV_LifeEndorsement.xml` b191
  generateDataDetail: 'Generate Data Detail', // `ViewCSVResult_LifeEDM.xml` b14322 → `GenerateDataDtlLife_act`
} as const

/** Keputusan - `Section/ConfirmSection.xml` (wadah `InputEDMLife.xml` b35518 `.IsJsonPolis=1`). */
export const PUTUSAN_EDM = {
  status: 'Status', // b496 radio `EmailTypePL`
  comment: 'Comment', // b829 `.Description`
  kolomDate: 'Date', // b1984
  kolomPic: 'PIC', // b2125
  kolomStatus: 'Status', // b2263
  kolomComment: 'Comment', // b2399
  submit: 'Submit', // `InputEDMLife.xml` b37494 (Confirm) dan b38109 (Decline)
} as const

/** `Section/ConfirmSubmitEDM.xml` - FlowAction `SetJsonPolisEDMLife_Confirm`. */
export const TERIMA_EDM = {
  terimaKasih: 'Thank you for Submit !', // b526
  noEndorsement: 'No. Endorsement', // b654 `.PremiumListSummary.PL_NUMBER_EDM`
  close: 'Close', // b1366 → `finishAssignment`
} as const

/**
 * Opsi radio `Status` - properti `EmailTypePL` tidak diekspor (R14, OQ-EDM-006); label dari
 * `Activity/AddHistorySuggest.xml` b353 `@if(.EmailTypePL=1,"Accept","Decline")`. Nilai `7` tidak ditawarkan.
 */
export const OPSI_KEPUTUSAN: Readonly<Record<string, string>> = {
  '1': 'Accept',
  '2': 'Decline',
}

/** Label opsi `EDM Type` - spec §5 `[keputusan work owner]` (opsi properti tidak diekspor). */
export const OPSI_EDM_TYPE: Readonly<Record<string, string>> = {
  '1': 'Perubahan Data',
  '3': 'Batal',
}

/**
 * Domain `TYPE_CEDING` (`Reinsurance System`) - STRUKTUR PremiumList bab `T_PREMIUM_LIST`,
 * keputusan tiket 00 Endorsement: `1`=QS, `2`=SURPLUS, `3`=QS+SURPLUS, `4`=XOL.
 */
export const OPSI_TYPE_CEDING: Readonly<Record<string, string>> = {
  '1': 'QS',
  '2': 'SURPLUS',
  '3': 'QS+SURPLUS',
  '4': 'XOL',
}
