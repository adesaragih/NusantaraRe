// `LABEL_MENU_KOMITE` (butir menu `Inbox Komite`) DIBUANG - menu datar,
// keputusan work owner 30-09-2026: tombol modul berlabel `M_NAV_MENU.LABEL`.
//
// ⭐ 12-10-2026 (work owner "perbaiki tampilannya"): SELURUH teks layar kasus komite berbahasa Inggris, label Title Case
// (sejajar layar Claim Life). Teks VERBATIM korpus tetap (ditandai di bawah).

/**
 * Label modul Komite Claim Life — tiket 01.
 *
 * ⚠️ Tiga kolom pertama Inbox Komite DIPINJAM dari kolom Pega-standar
 * `InboxPremiumList` (`pyFieldLabel` b721/b735/b764): worklist `KomiteRouter`
 * tidak punya section inbox sendiri di korpus. Sisanya kolom komite
 * (klaim induk, tingkat, nilai + `CURRENCY`, status baris) — kata layarnya
 * sama dengan tabel Committee Inbox di bawah inbox Claim Life.
 */
export const KOLOM_INBOX_KOMITE = {
  /** dipinjam `InboxPremiumList` b721. */
  kasusId: 'Case ID',
  /** dipinjam b735 — di sini waktu pembaruan terakhir work object. */
  tglUpdate: 'Update Date/Time',
  /** dipinjam b764. */
  statusWork: 'Work Status',
  nomorKlaim: 'Parent Claim',
  tingkat: 'Level',
  nilaiKlaim: 'Claim Amount',
  mataUang: 'Currency',
  statusBaris: 'Row Status',
} as const

export const KASUS_KOMITE = {
  judul: 'Committee Case',
  kembali: 'Back to Claim Life Inbox',
  /** `Section/ShowTransfer.xml` b29209 (12-10-2026 lanjutan: urutan dan judul kotak mengikuti ShowTransfer). */
  tangga: 'List of Committee',
  /**
   * Penyerahan dari Claim Life (pop-up Send Claim to Committee, 12-10-2026). ShowTransfer b27637 / b27836:
   * `Submission Date`, `Remarks` (kotak tanpa judul di Pega; judul kartu kosakata kami).
   */
  penyerahan: 'Submission',
  penyerahanTanggal: 'Submission Date',
  penyerahanPIC: 'PIC',
  penyerahanRemarks: 'Remarks',
  urut: 'Level',
  /**
   * Tiket 09 — judul VERBATIM grid `.KomiteList` `Section/ShowTransfer.xml`
   * b29727: `Committee` (`.IDKomite` = jabatan), `Status` (`.KomiteAproval`),
   * `Date Approve`, `Comment`.
   */
  committee: 'Committee',
  status: 'Status',
  dateApprove: 'Date Approve',
  comment: 'Comment',
  /** `KOMITE_OPERATORID` — pengenal akun; kosakata kami. */
  anggota: 'Member',
  /**
   * OQ-K-05 (GILIRAN-17) — keputusan tingkat sebelum langkah 5.1 (Reject di
   * tingkat akhir) menimpanya. Kosakata kami: Pega tidak menampilkan riwayat
   * yang ditimpanya. Kolom tampil hanya bila ada tingkat yang tertimpa.
   */
  asli: 'Before Final Reject',
  tanpaTangga: 'No committee members found.',
  eskalasi: 'Escalation History',
  giliranAnda: 'Your turn to decide',
  /** Chip "Waiting for level N" — N = tingkat pertama yang masih menunggu. */
  menungguTingkat: 'Waiting for level',
  tanggaSelesai: 'Ladder Completed',
} as const

/**
 * Layar keputusan — tiket 02, `Section/ShowTransfer.xml`.
 *
 * ⚠️ Pilihan `.AcceptStatus` (b32607, `pyListSource associated`) milik aturan
 * properti yang TIDAK diekspor; `1 = Approve`, `2 = Reject` adalah
 * `[keputusan work owner]`, bukan teks korpus. Kalimat konfirmasi dan tombol
 * VERBATIM.
 */
export const KEPUTUSAN_KOMITE = {
  judul: 'Decision',
  /** VERBATIM — label di atas dropdown `.AcceptStatus` (b32572-b32626, wajib). */
  konfirmasi: 'Are you sure to accept this document?',
  judulKonfirmasi: 'Confirm Decision',
  peringatan: 'The decision is recorded for your level and cannot be changed afterwards.',
  keputusan: 'Decision',
  pilih: 'Select Decision',
  setuju: 'Approve',
  tolak: 'Reject',
  /**
   * VERBATIM b32893: `Note` (`.Comment`) — `pyRequired true` / `pyRequiredNew always` b32936-b32940 (WAJIB; server
   * menegakkannya, `ErrNoteKomiteWajib`). Koreksi 12-10-2026: laporan sebelumnya menyebut b31001 `.KomiteComment`
   * `pyRequired false` - itu kolom grid List of Committee, bukan medan keputusan.
   */
  komentar: 'Note',
  wajibNote: 'Note is required.',
  tanyaKonfirmasi: 'Submit this decision?',
  /** VERBATIM b34722. */
  submit: 'Submit',
  /** VERBATIM b33880. */
  cancel: 'Cancel',
  ya: 'Yes',
  mengirim: 'Submitting…',
  wajibPilih: 'Select a decision first.',
} as const

/**
 * Efek keluar dan laporan "perlu intervensi" — tiket 08. ⚠️ Tidak ada di
 * korpus (perilaku baru, ADR-0015); kosakata kami.
 */
export const EFEK_KOMITE = {
  judul: 'Outgoing Effects',
  jenis: 'Effect',
  keadaan: 'State',
  percobaan: 'Attempts',
  sejak: 'Since',
  laporan: 'Needs Intervention Today',
  laporanKosong: "Today's report: no outgoing effects need intervention.",
} as const

/**
 * Eskalasi — tiket 03. ⚠️ Tidak ada di korpus (penyimpangan sadar, ADR-0014);
 * kosakata kami.
 */
export const ESKALASI_KOMITE = {
  tombol: 'Escalate One Level',
  keterangan:
    'Skips the member of the current level who is unavailable. The skipped ' +
    'level records no decision.',
} as const

/**
 * Form of Internal Claim Requisition — `Section/ShowTransfer.xml` (korpus kelvin/KomiteClaimLife (Done)), urutan kotak
 * VERBATIM: Policy Detail b2117, Claim Detail b5253, Medical Advisor and Claim Analyst b10372, Adjustment Detail b13053,
 * Document Claim b22594. Label medan VERBATIM kecuali Title Case (`RECOMMENDATION` b12412, `CLAIM ANALYST` b12599) dan
 * ejaan `Sum Reasured` b16559 → Sum Reinsured (keputusan work owner 12-10-2026, sama dengan layar Claim Life).
 */
export const FORM_KOMITE = {
  judul: 'Form of Internal Claim Requisition',
  polis: 'Policy Detail',
  cedingCompany: 'Ceding Company',
  plan: 'Plan',
  plNumber: 'PL Number',
  policyHolder: 'Policy Holder Name',
  sobName: 'SOB Name',
  klaim: 'Claim Detail',
  namaTertanggung: 'Name of Insured',
  sertifikat: 'Certificate No',
  nomorPolis: 'Policy No',
  umur: 'Age',
  jenisKelamin: 'Sex',
  tanggalLahir: 'Date of Birth',
  mulai: 'Begin Date',
  expired: 'Expired Date',
  efektif: 'Effective Date',
  valuasiMulai: 'Valuation Begin Date',
  valuasiSelesai: 'Valuation Expired Date',
  terimaKlaim: 'Claim Received Date',
  dokumenLengkap: 'Document Complete Date',
  kejadian: 'Date of Loss',
  medis: 'Medical Advisor and Claim Analyst',
  diagnosa: 'Diagnose',
  rekomendasi: 'Recommendation',
  analis: 'Claim Analyst',
  tanpaItem: 'No items',
  adjustment: 'Adjustment Detail',
  mataUang: 'Currency',
  shareRnm: 'Share Nusantara Re',
  cedingRetention: 'Ceding Retention',
  claimRnm: 'Claim RNM',
  claimRetro: 'Claim Retro',
  sumReinsured: 'Sum Reinsured',
  sumInsured: 'Sum Insured',
  shareRetro: 'Share Retro',
  retrocededShare: 'Retroceded Share',
  payableTo: 'Payable To',
  namaBank: 'Name of Bank',
  cabangBank: 'Branch of Bank',
  nomorRekening: 'Account No',
  treatyType: 'Treaty Type',
  dokumen: 'Document Claim',
  refresh: 'Refresh',
  namaFile: 'File Name',
  jenisLampiran: 'Attachment Type',
  tanggalDokumen: 'Date',
  lihat: 'View',
  memuatKlaim: 'Loading claim data…',
  gagalKlaim: 'The parent claim could not be loaded.',
  pesertaHilang: 'The submitted insured row was not found in the parent claim.',
  polisHilang: 'The policy header was not found.',
} as const
