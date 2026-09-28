// Label modul PremiumList Life — tiket 01.
//
// ⛔ VERBATIM dari korpus `D:\XML\RNM_BRD\PremiumList Life\`, dengan path +
// baris di sebelahnya. Label yang diketik dari ingatan adalah label yang
// bergeser, dan yang bergeser tidak berbunyi: keduanya benar menurut dirinya
// sendiri.

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

/** Hasil penggolong `Decision3` — b1658 dan b1807. */
export const PENGGOLONG_POLIS = {
  premium: 'Premium',
  offer: 'Offer',
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
