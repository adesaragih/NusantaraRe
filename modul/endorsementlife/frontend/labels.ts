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
  memuat: 'Loading…',
  memeriksa: 'Checking…',
  rinci: 'Details',
  terkunci: 'Locked: the policy and the endorsement type cannot change once the case is created. Decline the case and create a new one.',
  menyimpan: 'Saving…',
  tandaiHapus: 'Mark for deletion',
  mengunggah: 'Uploading…',
  pilihBerkas: 'CSV file',
  barisDisimpan: 'Rows added',
  pesanTerpotong: 'Only the first rejections are listed.',
  kolomBaris: 'Row',
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
  // Sel Work Status kasus terbuka (STATUSS kosong) - "Input EDM Life", nama tahap assignment `InputEDMLife` berspasi
  // (keputusan work owner 08-10-2026, tampilan saja; database tetap kosong = terbuka).
  statusTerbuka: 'Input EDM Life',
  /** Kasus tertutup di kotak masuk (09-10-2026, seperti PremiumList) - bukan korpus. */
  tutup: 'Case closed — cannot be opened',
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

/** Summary Comparison sebelum/sesudah endorsement (permintaan work owner 08-10-2026) - bukan korpus. */
export const PERBANDINGAN_EDM = {
  judul: 'Summary Comparison',
  kolomBaris: 'Item',
  kolomPeserta: 'Participants',
  before: 'Old policy (before)',
  deleted: 'Deleted / Cancelled',
  added: 'Added (CSV)',
  after: 'After endorsement',
  difference: 'Difference (endorsement premium)',
  kosong: 'No participants to compare.',
  /** Judul tabel rekap `T_PREMIUM_LIST_SUMMARY` (tampilan 08-10-2026; keterangan dihapus 08-10-2026). */
  judulRekap: 'Summary to be posted on Confirm',
  judulKeputusan: 'Decision',
} as const

/** Judul bagian popup Details peserta (keputusan work owner 08-10-2026, tampilan) - bukan korpus. */
export const RINCI_EDM = {
  bagianPeserta: 'Participant',
  bagianPeriode: 'Period',
  bagianPertanggungan: 'Coverage',
  bagianPremi: 'Premium',
  bagianSpreading: 'Retro Spreading',
} as const

/** `Section/EndorsmentLife_Section.xml` - label yang `InputEDMLife` ambil dari deskripsi properti. */
export const BUAT_EDM = {
  policyNo: 'Policy No', // b1075 (`InputEDMLife.xml` b3126 memakai deskripsi properti yang sama)
  /** Keterangan Policy No bernomor endorsement (10-10-2026) - `{polis}`, `{versi}` diganti. Bukan korpus. */
  infoAkhiran: 'Endorsement number entered — using policy {polis} ({versi}).',
  versiTerakhir: 'latest version /{akhiran}',
  versiNB: 'new business version',
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
  /** Kolom tampilan rekap View Old Policy (keputusan work owner 10-10-2026) - bukan korpus. */
  endorsementNo: 'Endorsement No',
  newBusiness: 'New Business',
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
} as const

/**
 * Tombol `Save` b37202 (→ `SetPremi_EDM`) DIGANTI NAMA "Process Endorsement", dipindah ke bawah layar, dan
 * didahului konfirmasi (keputusan work owner 08-10-2026) - bukan korpus. `{n}` diganti angka oleh
 * `kalimatKonfirmasiProses`.
 */
export const PROSES_EDM = {
  processEndorsement: 'Process Endorsement',
  memproses: 'Processing…',
  judulKonfirmasi: 'Process Endorsement?',
  hapusN: '{n} participant(s) will be marked Delete and their amounts reversed.',
  hapusNol: 'No participants are marked Delete.',
  baruN: '{n} new participant(s) will be kept.',
  batalN: 'All {n} existing participant(s) will be marked Batal and their amounts reversed.',
  batalSemua: 'All existing participants will be marked Batal and their amounts reversed.',
  rekap: 'The summary will be recalculated. This can only be done once.',
  batal: 'Cancel',
  lanjut: 'Continue',
  /** Langkah FINAL (08-10-2026) - SAMA dengan `models.PesanSudahDiproses` server. */
  sudahDiproses: 'This endorsement has been processed. Participants can no longer be added or deleted.',
  /** Proteksi tanpa perubahan (08-10-2026) - SAMA dengan `models.PesanTanpaPerubahan` server. */
  tanpaPerubahan: 'No changes to process. Mark at least one participant for Delete or add participants from CSV.',
} as const

/** Unggah CSV - wadah `InputEDMLife.xml` b8698, popup `ViewCSVResult_LifeEDM`, modal `UploadCSV_LifeEndorsement`. */
export const UNGGAH_EDM = {
  uploadCsv: 'Upload CSV', // b8973 → `UploadCSV_LifeEndorsement`
  viewUpload: 'View Upload', // b9340 → `ViewCSVResult_LifeEDM`
  addCsvData: 'Add CSV Data', // b10405 → `SaveCSVEDMLife`
  judulModal: 'Endorsement Life - Upload CSV', // `FlowAction/UploadCSV_LifeEndorsement.xml` b191
  generateDataDetail: 'Generate Data Detail', // `ViewCSVResult_LifeEDM.xml` b14322 → `GenerateDataDtlLife_act`
  /** Templat CSV + daftar kolom wajib (keputusan work owner 10-10-2026) - bukan korpus. */
  downloadTemplate: 'Download Template',
  /** Tinjauan unggahan untuk user awam (disederhanakan 10-10-2026) - `{n}`, `{total}`, `{kolom}` diganti. */
  ringkasMasalah: '{n} of {total} rows have a problem. Fix them in your file and upload again.',
  ringkasLolos: 'All {total} rows are valid. Click Add CSV Data to save them.',
  kolomAsing: 'These columns are not used and will be ignored: {kolom}',
  catatanBaris: 'Row numbers match the rows in your spreadsheet.',
  kolomMasalah: 'Problem',
  petunjukTemplat: 'Download the template to see the correct format.',
} as const

/** Keputusan - `Section/ConfirmSection.xml` (wadah `InputEDMLife.xml` b35518 `.IsJsonPolis=1`). */
export const PUTUSAN_EDM = {
  status: 'Status', // b496 radio `EmailTypePL`
  comment: 'Comment', // b829 `.Description`
  kolomDate: 'Date', // b1984
  kolomPic: 'PIC', // b2125
  kolomStatus: 'Status', // b2263
  kolomComment: 'Comment', // b2399
} as const

/**
 * Tombol `Submit` b37494 (Confirm) / b38109 (Decline) DIGANTI NAMA "Confirm" + popup konfirmasi (keputusan work owner
 * 08-10-2026, pola Confirm PremiumList) - bukan korpus. Perilaku keputusan tidak berubah.
 */
export const KONFIRMASI_PUTUSAN_EDM = {
  confirm: 'Confirm',
  judulAccept: 'Confirm endorsement?',
  isiAccept: "This will create the endorsement number and finish the case. You can't undo this.",
  judulDecline: 'Decline endorsement?',
  isiDecline: "This will decline the endorsement and close the case. You can't undo this.",
  batal: 'Cancel',
  ya: 'Yes',
} as const

/**
 * Popup hasil Confirm (Accept) - menggantikan `ConfirmSubmitEDM` (`Thank you for Submit !` b526, `No. Endorsement`
 * b654, `Close` b1366) atas permintaan work owner 09-10-2026, pola "PL Number created" PremiumList - bukan korpus.
 */
export const HASIL_EDM = {
  judul: 'Endorsement confirmed',
  noEndorsement: 'Endorsement No',
  edmType: 'EDM Type',
  subType: 'Sub-type',
  caseId: 'Case ID',
  ok: 'OK',
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
