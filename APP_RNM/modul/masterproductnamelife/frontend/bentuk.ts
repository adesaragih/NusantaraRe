// Logika murni layar Master Product Name Life - aturan XML yang Pega jalankan DI KLIEN (DataTransform dan
// activity `refresh` tanpa RDB). Tanpa React, tanpa jaringan; diuji `bentuk.test.ts`.
//
// ⛔ Angka tetap TEKS (ADR-0003): `CountMaxSumReasured_Act` dihitung per digit lewat `inti/lib/desimal`, tidak
// pernah lewat `Number`.

import { desimalSah, jumlahDesimal } from '../../../inti/frontend/lib/desimal'
import type { Opsi } from '../../../inti/frontend/components/ui/dasar'
import type { Produk, ProdukInward, ProdukUmum } from './api'
import { PEMBAYARAN_MPNL } from './labels'

function umumKosong(): ProdukUmum {
  return {
    productName: '', ceding: '', cedingId: '', sobName: '', sobId: '', riComm: '', riRisk: '', riRiskId: '',
    inwardName: '', treatyNumber: '', cause: '', causeId: '', policyHolder: '', policyHolderName: '',
    createOp: '', updateOp: '', isOrs: false, comment: '',
    type: '', typeCeding: '', grup: '', productCode: '', productType: '', productTypeId: '', riRate: '', riRateId: '',
    riCommId: '', outwardName: '', outwardNameId: '', outwardRate: '', outwardRateId: '', outwardComm: '',
    outwardCommId: '', benefit: '', benefitId: '',
  }
}

function inwardKosong(): ProdukInward {
  return {
    id: '', productId: '', policyHolder: '', policyHolderName: '', insured: '', addendumNo: '', addendumWord: '',
    amandementNo: '', amandementSchd: '', maxExpiredClaim: '', begin: '', stnc: '', cedingRetentionNum: '',
    cedingLimit: '', brokerage: '', minAge: '', maxAge: '', expiryAge: '', extraPremi: '', minSumInsured: '',
    maxSumInsured: '', maxSumReasured: '', rnmShare: '', rnmLimitNum: '', premiumFactor: '', payment: '',
    subjectTo: '', annuityInterest: '', premiumRefundFactor: '', maxDataReceive: '', mature: '', birthday: '',
    currency: '', currencyId: '', extraMortality: '', maxContract: '', proportionalTable: '',
    ceding: '', treatyNumber: '', inwardTreatyNm: '', cedingRetentionPct: '', cedingLimitXpn: '', rnmLimitPct: '',
    lienClause: '', months: '',
  }
}

/** `Add` b71865 → `NewProductLife`: 1 b234 Page-New `ProductName`, 3 b579 hapus `ProductNameInward` - halaman kosong. */
export function produkBaru(): Produk {
  return {
    id: '', umum: umumKosong(), inward: inwardKosong(), lienClause: [], documentClaim: [], planList: [],
    financialUnderwriting: [], underwritingLimit: [], outwardList: [], commentList: [],
  }
}

/**
 * `Copy` b59854 → DataTransform `CopyProduct`: 1 b144 `ProductName.ID := ""`, 2 b173 `ProductNameInward.ID := ""`;
 * selebihnya halaman tetap. Server mewarisi medan milik server dari produk asal (`salinanDari`).
 */
export function salinProduk(p: Produk): Produk {
  return {
    ...p,
    id: '',
    inward: { ...p.inward, id: '', productId: '' },
    salinanDari: p.salinanDari !== undefined && p.id === '' ? p.salinanDari : p.id,
    hitungOutward: p.hitungOutward ?? false,
  }
}

/** `SetTreatyName_Act` 2 b436: `INWARDNAME = PRODUCTNAME + " " + POLICYHODERNAME` (tanpa pemangkasan, seperti Pega). */
export function namaTreaty(productName: string, policyHolderName: string): string {
  return productName + ' ' + policyHolderName
}

/** Desimal dengan koma sebagai pemisah diterima seperti server (`12,5`); kosong = `0` (`local.* = 0`, b239). */
function angka(v: string): string | null {
  const t = v.trim()
  if (t === '') return '0'
  const titik = t.includes(',') && !t.includes('.') ? t.replace(',', '.') : t
  return desimalSah(titik) ? titik : null
}

function negasi(v: string): string {
  if (v.startsWith('-')) return v.slice(1)
  if (v.startsWith('+')) return '-' + v.slice(1)
  return '-' + v
}

/**
 * `CountMaxSumReasured_Act` (onChange `Ceding's Limit` b22725 dan `Max Sum Insured` b24024):
 * 1 b239 `local.MaxSumInsured = 0`, `local.CedingLimit = 0`, lalu nilai medan; 2 b436
 * `MAXSUMREASURED = local.MaxSumInsured - local.CedingLimit`. Isian yang bukan angka = `null` (medan dibiarkan).
 */
export function hitungMaxSumReasured(maxSumInsured: string, cedingLimit: string): string | null {
  const a = angka(maxSumInsured)
  const b = angka(cedingLimit)
  if (a === null || b === null) return null
  return jumlahDesimal([a, negasi(b)]).total
}

/** Pilihan `Premium Payment Method` b25611 - kode dan teks dari `GenerateUpload_Act` b1141 (OQ-MPNL-05). */
export const PILIHAN_PEMBAYARAN: readonly Opsi[] = [
  { value: '1', label: PEMBAYARAN_MPNL.annual },
  { value: '2', label: PEMBAYARAN_MPNL.semiAnnual },
  { value: '3', label: PEMBAYARAN_MPNL.quarterly },
  { value: '4', label: PEMBAYARAN_MPNL.monthly },
]

/** `Premium Factor (%)` b25398 - visibilitas `OTHER ProductNameInward.PAYMENT==3`. */
export function tampilPremiumFactor(payment: string): boolean {
  return payment.trim() === '3'
}

/** `View Office Online` b69291 - visibilitas `OTHER .pyFileMimeType = xls/xlsx/doc/docx/ppt/pptx`. */
export function tampilViewOffice(ekstensi: string): boolean {
  return ['xls', 'xlsx', 'doc', 'docx', 'ppt', 'pptx'].includes(ekstensi)
}

/**
 * `.Date` grid History b62561 (`pxDateTime`): stempel `@CurrentDateTime()` Pega `YYYYMMDDTHHMMSS.SSS GMT`
 * tampil sebagai tanggal-jam LOKAL `DD-MM-YYYY HH:mm` (format tanggal seragam NFR-14 ditambah jam); teks
 * lain tampil apa adanya. `selisihMenit` = zona lokal terhadap GMT (bawaan: zona peramban). Audit 02-10-2026:
 * dulu stempel mentah yang tampil.
 */
export function waktuPega(teks: string, selisihMenit = -new Date().getTimezoneOffset()): string {
  const m = /^(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})(\d{2})(?:\.\d+)?\s*GMT$/.exec(teks.trim())
  if (!m) return teks
  const n = (i: number) => Number(m[i] ?? '')
  const d = new Date(Date.UTC(n(1), n(2) - 1, n(3), n(4), n(5), n(6)) + selisihMenit * 60_000)
  const p = (x: number) => String(x).padStart(2, '0')
  return `${p(d.getUTCDate())}-${p(d.getUTCMonth() + 1)}-${d.getUTCFullYear()} ${p(d.getUTCHours())}:${p(d.getUTCMinutes())}`
}

/** `View Rate` b34113 - visibilitas `OTHER .RIRATE!=''`. */
export function tampilViewRate(riRate: string): boolean {
  return riRate !== ''
}

/** `pyRDLPageSize` b71371 - 10 baris per halaman grid produk. */
export const UKURAN_HALAMAN_MPNL = 10

/** Potongan satu halaman grid (1-berbasis). */
export function potongHalaman<T>(semua: readonly T[], halaman: number): T[] {
  const awal = (halaman - 1) * UKURAN_HALAMAN_MPNL
  return semua.slice(awal, awal + UKURAN_HALAMAN_MPNL)
}

/** Halaman yang sah sesudah daftar berubah panjang. */
export function jepitHalaman(halaman: number, total: number): number {
  return Math.min(Math.max(1, halaman), Math.max(1, Math.ceil(total / UKURAN_HALAMAN_MPNL)))
}

/** `CopyFinancialWriting` 1 b225 / `CopyUnderWritingLimit` 1 b226: baris salinan ditambah di akhir daftar. */
export function salinBaris<T extends { asli?: string }>(daftar: readonly T[], i: number): T[] {
  const b = daftar[i]
  if (b === undefined) return [...daftar]
  const salin: T = { ...b }
  delete salin.asli
  return [...daftar, salin]
}
