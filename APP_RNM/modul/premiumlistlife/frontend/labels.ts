// Label modul PremiumList Life — tiket 01.
//
// ⛔ VERBATIM dari korpus `D:\XML\RNM_BRD\PremiumList Life\`, dengan path +
// baris di sebelahnya. Label yang diketik dari ingatan adalah label yang
// bergeser, dan yang bergeser tidak berbunyi: keduanya benar menurut dirinya
// sendiri.

// `LABEL_MENU_PREMIUMLIST` (butir menu `PremiumList`) DIBUANG - menu datar,
// keputusan work owner 30-09-2026: tombol modul berlabel `M_NAV_MENU.LABEL`.

/**
 * Kolom kotak masuk — `ReportDefinition/InboxPremiumList.xml`.
 *
 * ⛔ SEBELAS dari tiga belas. Dua yang tidak ada:
 *
 *   `.pzInsKey` b961 `Instance Handle Key` — pengenal internal Pega, bukan
 *     kolom untuk manusia.
 *   `A.KetentuanUnderwriting` b891 — **tidak punya kolom** di migrasi
 *     050–056 mana pun. Ia TIDAK ditampilkan, bukan diisi teks kosong:
 *     sel kosong di layar terbaca "memang kosong", bukan "kami tidak punya
 *     datanya".
 */
export const KOLOM_INBOX_POLIS = {
  /** b721 `pyFieldLabel`. */
  caseId: 'Case ID',
  /** b735. */
  tglCreate: 'Create Date/Time',
  /** b750. */
  createOpName: 'Create Operator Name',
  /** b764. */
  statusWork: 'Work Status',
  /** b778. */
  cedingCoName: 'CedingCoName',
  /** b794. */
  policyHolderName: 'PolicyHolderName',
  /** b807 — ⚠️ jalur propertinya `PremiumListSummary`, kolomnya di detail. */
  plNumber: 'PL_NUMBER',
  /** b821. */
  riSlipRnm: 'RISLIPRNM',
  /** b836. */
  type: 'Type',
  /** b849. */
  marketingName: 'MarketingName',
  /** b863. */
  sobName: 'SobName',
  /** b878. */
  dateReceived: 'DateReceived',
} as const

/**
 * Tombol portal — `Section/PremiumList.xml`.
 *
 * ⚠️ DUA baris untuk satu label, dan keduanya nyata: `pyLabel` pada kendali
 * (b3273, b3921) dan `pyRuleName pyButtonLabel` pada rule labelnya (b16964,
 * b16472). Yang dipakai layar adalah `pyLabel`.
 */
export const TOMBOL_POLIS = {
  /** b3273 `pyLabel` → `runActivity` b3283 → `CreateInputLife` b3291. */
  inputOffer: 'Input Offer',
  /** b3921 `pyLabel` → `runActivity` b3931 → `CreateInputLife` b3939. */
  inputPremium: 'Input Premium',
} as const

/**
 * Keputusan penawaran — VERBATIM `pyExpression` konektor
 * `InputPolicyHolder.xml`.
 *
 * ⛔ `Reject` hanya ada pada `Decision2` (b2306), yaitu sesudah Input Premium
 * Detail. Di tahap penawaran ia TIDAK punya jalur — lihat ralat AC 2 tiket 01.
 */
export const KEPUTUSAN_POLIS = {
  confirm: 'Confirm',
  reject: 'Reject',
  decline: 'Decline',
} as const

/**
 * Judul kolom grid Premium List Detail — `Section/PL_Detail_Sec.xml`.
 *
 * ⛔ NOL judul yang dikarang. Section itu tidak memuat satu pun `pyCaption`
 * untuk medannya: Pega menampilkan label properti masing-masing, yang tidak
 * ikut ter-ekspor ke korpus. Karena itu yang di bawah adalah NAMA KOLOM apa
 * adanya, dirapikan hanya dengan mengganti garis bawah menjadi spasi —
 * bukan terjemahan yang kami karang sendiri.
 *
 * ⚠️ Nama kolomnya sendiri DATANG DARI SERVER; peta ini hanya memperindah
 * yang datang. Kolom yang tidak ada di sini tetap tampil, dengan namanya.
 */
export const JUDUL_KOLOM_PESERTA: Record<string, string> = {
  BEGIN_DATE: 'Begin Date',
  ENTRY_AGE: 'Entry Age',
  PLAN: 'Plan',
  STNC: 'STNC',
  EXPIRED_DATE: 'Expired Date',
  CURRENT_AGE: 'Current Age',
  PERIOD_MM: 'Period MM',
  POLICY_HOLDER: 'Policy Holder',
  EFFECTIVE_DATE: 'Effective Date',
  WPC: 'WPC',
  CERTIFICATE_NO: 'Certificate No',
  CURRENCY: 'Currency',
  SUM_REASURED: 'Sum Reasured',
  DEDUCTION: 'Deduction',
  BROKERAGE_FEE: 'Brokerage Fee',
  SUM_AT_RISK_RETRO: 'Sum At Risk Retro',
  RETROCEDED_SHARE: 'Retroceded Share',
  SUM_INSURED: 'Sum Insured',
  SHARE_NUSANTARA_RE_GROSS: 'Share Nusantara Re Gross',
  EM_PERCENT: 'EM Percent',
  SUM_AT_RISK_GROSS: 'Sum At Risk Gross',
  GROSS_PREMIUM: 'Gross Premium',
  CEDING_RETENTION: 'Ceding Retention',
  RATE: 'Rate',
  FACTOR: 'Factor',
  RISK: 'Risk',
  NET_PREMIUM: 'Net Premium',
  DEDUCTION_REFUND: 'Deduction Refund',
  BROKERAGE_FEE_REFUND: 'Brokerage Fee Refund',
  GROSS_PREMIUM_REFUND: 'Gross Premium Refund',
  NET_PREMIUM_REFUND: 'Net Premium Refund',
  RI_ADMIN_FEE_REFUND_RETRO: 'RI Admin Fee Refund Retro',
  SHARE_RETRO: 'Share Retro',
  GROSS_PREMIUM_REFUND_RETRO: 'Gross Premium Refund Retro',
  NET_PREMIUM_REFUND_RETRO: 'Net Premium Refund Retro',
  RI_ADMIN_FEE_RETRO: 'RI Admin Fee Retro',
  GROSS_PREMIUM_RETRO: 'Gross Premium Retro',
  NET_PREMIUM_RETRO: 'Net Premium Retro',
}

/** Label layar Premium List Detail — tiket 03. */
export const DETAIL_POLIS = {
  judul: 'Premium List Detail',
  nomor: 'PL_NUMBER',
  belumBernomor: 'Belum bernomor',
  terbitkan: 'Generate PL Number',
  /**
   * ⛔ Kalimatnya MENYEBUT apa yang harus dikerjakan lebih dahulu. "Tidak
   * dapat dinomori" membuat orang mencari kerusakan yang tidak ada.
   */
  perluPeserta:
    'Polis ini belum punya baris peserta, sehingga nomornya belum punya ' +
    'tempat tersimpan. Unggah rincian peserta lebih dahulu.',
  sudahBernomor: 'Polis ini sudah bernomor; tombolnya tidak perlu ditekan lagi.',
} as const

/**
 * Label layar unggah CSV peserta — tiket 04.
 *
 * ⚠️ Pesan penolakan TIDAK ada di sini: ia datang dari server, VERBATIM dari
 * `ValidasiUploadPL_act.xml`. Menyalinnya ke layar membuat dua sumber untuk
 * satu kalimat, dan yang satu akan tertinggal.
 */
export const UNGGAH_CSV = {
  judul: 'Upload CSV Premium List Detail',
  pilihBerkas: 'Berkas CSV',
  tinjau: 'Tinjau',
  simpan: 'Simpan permanen',
  /**
   * ⛔ Aturan pemisah DINYATAKAN DI MUKA, bukan hanya saat menolak. Korpus
   * menyebutnya enam kali di nama langkahnya ("SEPARATOR MENGGUNAKAN TITIK")
   * dan tidak pernah di pesan yang dilihat pemakai — sehingga orang baru tahu
   * aturannya setelah berkasnya ditolak.
   */
  aturanPemisah:
    'Angka memakai TITIK sebagai pemisah desimal, dan tanpa pemisah ribuan. ' +
    'Contoh: 1234567.89 — bukan 1,234,567.89 dan bukan 1.234.567,89. ' +
    'Tanggal berformat dd/mm/yyyy.',
  perbaikiDulu:
    'Perbaiki dulu baris yang ditolak, lalu tinjau ulang. Selama masih ada ' +
    'penolakan, tidak ada satu baris pun yang disimpan.',
  kolomBaris: 'Baris',
  kolomKolom: 'Kolom',
  kolomPesan: 'Pesan',
  kolomSebab: 'Sebab',
} as const

/**
 * Label layar rekap — tiket 05a bagian 2, `Section/ShowLifePremiumSummary.xml`.
 */
export const SUMMARY_POLIS = {
  judul: 'Summary Premium Life',
  submit: 'Submit',
  tanpaRekap: 'Belum ada rekap: polis ini belum punya baris peserta.',
  /** Tiket 05b — `finishAssignment` → `END52`, `pyWorkStatus` b947. */
  ditutup: 'Kasus ditutup: Resolved-Completed.',
  /** Tiket 06 — efek keluar sesudah simpan. */
  efekDilewati: 'Kiriman ke Arasapas dilewati: lingkungan ini bukan produksi.',
  efekGagal: 'Rekap tersimpan, tetapi efek keluar gagal:',
  efekTerantre: 'Kegagalannya terantre untuk dicoba ulang.',
  efekTidakTerantre: 'Sebagian kegagalan TIDAK dapat diantre — hubungi tim operasi.',
  /** ⛔ `.COB` baris rekap tidak ditetapkan rule mana pun di korpus. */
  cobTanpaSumber:
    'Kolom COB tidak terisi: tidak satu pun rule PremiumList Life menetapkan ' +
    '.COB pada baris rekap mata uang.',
} as const

/** Satu kolom grid rekap: judul VERBATIM dan properti sumbernya. */
export interface KolomRekap {
  judul: string
  kolom: string
}

/**
 * Grid rekap per `Type`, VERBATIM.
 *
 * `[terverifikasi]` Empat layout ber-`pyContainerVisibleWhen` `.Type = 'QR'`
 * / `'QP'` / `'TP'` / `'TR'` di `ShowLifePremiumSummary.xml` - judul grid
 * `Summary` di b4018 / b8761 / b15196 / b20690 (nomor baris =
 * `sed -e 's/></>
</g'`); judul dan properti kolom disalin urut dokumen.
 *
 * ⚠️ KEANEHAN WARISAN, disalin apa adanya: pada QP, TP, dan TR judul
 * `PREMIUM DEDUCTION` berdiri di atas properti `.COMMISSION`. Pada QR judul
 * `DEDUCTION` berdiri di atas `.DEDUCTION` dan kolom `.COMMISSION` tidak ada.
 */
export const GRID_REKAP: Record<string, readonly KolomRekap[]> = {
  QR: [
    { judul: 'COB', kolom: 'COB' },
    { judul: 'PL NUMBER', kolom: 'PL_NUMBER' },
    { judul: 'CURRENCY', kolom: 'CURRENCY' },
    { judul: 'PREMIUM', kolom: 'PREMIUM' },
    { judul: 'DEDUCTION', kolom: 'DEDUCTION' },
    { judul: 'BROKERAGE FEE', kolom: 'BROKERAGE_FEE' },
    { judul: 'RI ADMIN FEE', kolom: 'RI_ADMIN_FEE' },
    { judul: 'TAX', kolom: 'TAX' },
    { judul: 'PROF COMM', kolom: 'PROF_COMM' },
    { judul: 'CLAIM', kolom: 'CLAIM' },
    { judul: 'BALANCE', kolom: 'BALANCE' },
  ],
  QP: [
    { judul: 'COB', kolom: 'COB' },
    { judul: 'PL NUMBER', kolom: 'PL_NUMBER' },
    { judul: 'CURRENCY', kolom: 'CURRENCY' },
    { judul: 'PREMIUM', kolom: 'PREMIUM' },
    { judul: 'PREMIUM DEDUCTION', kolom: 'COMMISSION' },
    { judul: 'BROKERAGE FEE', kolom: 'BROKERAGE_FEE' },
    { judul: 'OVR COMM', kolom: 'OVR_COMM' },
    { judul: 'TAX', kolom: 'TAX' },
    { judul: 'GROSS_PREMIUM_REFUND', kolom: 'GROSS_PREMIUM_REFUND' },
    { judul: 'DEDUCTION REFUND', kolom: 'DEDUCTION_REFUND' },
    { judul: 'RI ADMIN FEE REFUND', kolom: 'RI_ADMIN_FEE_REFUND' },
    { judul: 'BROKERAGE FEE REFUND', kolom: 'BROKERAGE_FEE_REFUND' },
    { judul: 'TAX REFUND', kolom: 'TAX_REFUND' },
    { judul: 'CLAIM_AMOUNT', kolom: 'CLAIM_AMOUNT' },
    { judul: 'NET_PREMIUM_REFUND', kolom: 'NET_PREMIUM_REFUND' },
    { judul: 'BALANCE', kolom: 'BALANCE' },
  ],
  TP: [
    { judul: 'COB', kolom: 'COB' },
    { judul: 'PL NUMBER', kolom: 'PL_NUMBER' },
    { judul: 'CURRENCY', kolom: 'CURRENCY' },
    { judul: 'SHARE RETRO', kolom: 'SHARE_RETRO' },
    { judul: 'GROSS PREMIUM RETRO', kolom: 'GROSS_PREMIUM_RETRO' },
    { judul: 'PREMIUM DEDUCTION', kolom: 'COMMISSION' },
    { judul: 'BROKERAGE FEE RETRO', kolom: 'BROKERAGE_FEE_RETRO' },
    { judul: 'DISCOUNT PREMIUM RETRO', kolom: 'DISCOUNT_PREMIUM_RETRO' },
    { judul: 'RI ADMIN FEE RETRO', kolom: 'RI_ADMIN_FEE_RETRO' },
    { judul: 'TAX', kolom: 'TAX' },
    { judul: 'PROF COMM', kolom: 'PROF_COMM' },
    { judul: 'CLAIM', kolom: 'CLAIM' },
    { judul: 'BALANCE', kolom: 'BALANCE' },
  ],
  TR: [
    { judul: 'COB', kolom: 'COB' },
    { judul: 'PL NUMBER', kolom: 'PL_NUMBER' },
    { judul: 'CURRENCY', kolom: 'CURRENCY' },
    { judul: 'SHARE RETRO', kolom: 'SHARE_RETRO' },
    { judul: 'GROSS PREMIUM REFUND RETRO', kolom: 'GROSS_PREMIUM_REFUND_RETRO' },
    { judul: 'PREMIUM DEDUCTION', kolom: 'COMMISSION' },
    { judul: 'BROKERAGE FEE REFUND RETRO', kolom: 'BROKERAGE_FEE_REFUND_RETRO' },
    { judul: 'DISCOUNT PREMIUM REFUND RETRO', kolom: 'DISCOUNT_PREMIUM_REFUND_RETRO' },
    { judul: 'RI ADMIN FEE REFUND RETRO', kolom: 'RI_ADMIN_FEE_REFUND_RETRO' },
    { judul: 'TAX', kolom: 'TAX' },
    { judul: 'PROF COMM', kolom: 'PROF_COMM' },
    { judul: 'CLAIM', kolom: 'CLAIM' },
    { judul: 'BALANCE', kolom: 'BALANCE' },
  ],
}
