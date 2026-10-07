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
  // Kolom nomor baris tabel peserta View Old Policy (07-10-2026) - tampilan saja, bukan kolom korpus.
  nomor: 'No',
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
  efekDilewati: 'Not sent to Arasapas: this is not the production environment.',
  efekGagal: 'The endorsement is saved, but sending it onward failed and needs attention:',
} as const

/** `Section/InboxEndorsementLife.xml` - halaman awal (harness portal `InboxEndorsementLife`). */
export const INBOX_EDM = {
  judul: 'Endorsement Life', // b5505
  createAddendum: 'Create Addendum', // b6620
  kolomCaseId: 'Case ID', // b8418
  kolomEndorsementNo: 'Endorsement No', // b8556
  kolomType: 'Type', // b8692
  kolomEdmType: 'EDM Type', // b8828
  // b8964 korpus "Policy No" → "PL_NUMBER" (istilah kotak masuk PremiumList Life) - keputusan work owner 07-10-2026;
  // isi tetap OLD_POLICY_NO (PL Number polis yang di-endorse). Label "Policy No" layar lain tidak berubah.
  kolomPolicyNo: 'PL_NUMBER',
  kolomCeding: 'Ceding', // b9236
  kolomCreateDate: 'Create Date', // b9644
  kolomCreateOperator: 'Create Operator', // b9780
  // b10052 korpus "Status" → "Work Status" (istilah kotak masuk PremiumList Life) - keputusan work owner 07-10-2026.
  kolomStatus: 'Work Status',
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
  // Judul bagian layar kasus dirapikan (keputusan work owner 07-10-2026; pola Premium List Detail) - bukan korpus.
  bagianPolis: 'Policy Data',
  bagianPihak: 'Parties',
  bagianEndorsement: 'Endorsement',
  peserta: 'Participants',
} as const

/**
 * Nama `Premium Method` (`PRO_RATE_TYPE`) - 1 Single, 2 Annually, 3 Others; sama dengan opsi Premium Payment
 * Method PremiumList Life (keputusan work owner 07-10-2026). Kode lain tampil apa adanya.
 */
export const OPSI_METODE_PREMI: Readonly<Record<string, string>> = {
  '1': 'Single',
  '2': 'Annually',
  '3': 'Others',
}

/** `Section/EndorsmentLife_Section.xml` - label yang `InputEDMLife` ambil dari deskripsi properti. */
export const BUAT_EDM = {
  policyNo: 'Policy No', // b1075 (`InputEDMLife.xml` b3126 memakai deskripsi properti yang sama)
  edmType: 'EDM Type', // b1347
  description: 'Description', // b2140
  edmDate: 'EDM Date', // b3160
  // Sub-jenis EDM Type (OQ-EDM-005 ditutup, keputusan work owner 07-10-2026).
  edmTypePerubahanData: 'EDM Type Perubahan Data', // b1645
  edmTypeBatal: 'EDM Type Batal', // b1937
  pilih: 'Choose',
  submit: 'Submit', // b4226
  // Tampilan dirapikan (keputusan work owner 07-10-2026): footer Cancel + Submit, teks bantuan Submit nonaktif.
  cancel: 'Cancel',
  bantuan: {
    policyNo: 'Fill in Policy No to continue.',
    edmType: 'Choose an EDM Type (Perubahan Data or Batal) to continue.',
    memeriksa: 'Checking the policy…',
    menunggu: 'Waiting for the policy check.',
    pesan: 'Resolve the messages above before submitting.',
    penolakan: 'Change the input to submit again.',
    perubahanData: 'Choose an EDM Type Perubahan Data to continue.',
    batal: 'Choose an EDM Type Batal to continue.',
    mengirim: 'Submitting…',
  },
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
 * Sub-jenis EDM Type (OQ-EDM-005, keputusan work owner 07-10-2026) - opsi dari layar Pega, DISIMPAN sebagai teks
 * label persis (`models.OpsiEdmTypePerubahanData` / `OpsiEdmTypeBatal`, migrasi 483).
 */
export const OPSI_EDM_TYPE_PERUBAHAN_DATA: readonly string[] = [
  'Perubahan data (No Polis)',
  'Perubahan data (No Sertifikat)',
  'Perubahan data (Nama Peserta)',
  'Perubahan data (DOB)',
  'Perubahan data (Begin & Expired Date)',
  'Perubahan data (UP & Premi)',
]
/** Opsi `EDM Type Batal` - tanpa Choose; yang pertama bawaan terpilih. */
export const OPSI_EDM_TYPE_BATAL: readonly string[] = ['Batal Internal', 'Batal Eksternal']

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
