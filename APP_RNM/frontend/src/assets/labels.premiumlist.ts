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
