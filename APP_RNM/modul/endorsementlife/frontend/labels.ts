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
  edmDate: 'EDM Date', // b3160
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
