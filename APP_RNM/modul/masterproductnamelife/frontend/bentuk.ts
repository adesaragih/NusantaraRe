// Logika murni layar Master Product Name Life - aturan XML yang Pega jalankan DI KLIEN (DataTransform dan
// activity `refresh` tanpa RDB). Tanpa React, tanpa jaringan; diuji `bentuk.test.ts`.
//
// ⛔ Angka tetap TEKS (ADR-0003): `CountMaxSumReasured_Act` dihitung per digit lewat `inti/lib/desimal`, tidak
// pernah lewat `Number`.

import { desimalSah, jumlahDesimal } from '../../../inti/frontend/lib/desimal'
import type { Opsi } from '../../../inti/frontend/components/ui/dasar'
import type { HasilSalinLama, Produk, ProdukInward, ProdukLama, ProdukUmum, StatusSalinLama } from './api'
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

/**
 * `Document List` b15286 (grid `DOCUMENT CLAIM`) - daftar PromptList properti `.Document` kelas
 * `ASM-FW-GISFW-Data-UnderwritingLimit` (`Rule-Obj-Property`, ruleset GISFW 01-01-91, `pyTableOption` PromptList,
 * `pyPromptTableList` 16 baris; `pyStandardValue` = `pyLocalizedValue`). Rule ini tidak ada di korpus ekspor (OQ-MPNL-05);
 * XML-nya dikirim work owner 03-10-2026: "untuk document list productname pilih dari list ini". Urutan = urutan rule.
 */
export const PILIHAN_DOKUMEN_KLAIM: readonly string[] = [
  'Sertifikat peserta (Participant certificate)',
  'Copy identitas diri KTP/SIM/Paspor (Copy of ID card/Driving license/Passport)',
  'Copy kartu keluarga (Copy of family card)',
  'Copy sertifikat kematian (Copy of death certificate)',
  'Copy bukti pembayaran klaim (Copy of claim payment receipt)',
  'Copy legalisir rincian biaya perawatan dari rumah sakit (Legalized copy of hospital treatment cost details)',
  'Copy legalisir kwitansi biaya perawatan dari rumah sakit (Legalized copy of hospital payment receipts)',
  'Surat pernyataan meninggal oleh dokter/rumah sakit (Doctor/Hospital death statement letter)',
  'Surat keterangan meninggal oleh polisi (Police Statement for death)',
  'Surat keterangan meninggal karena kecelakaan oleh polisi (Police Statement for accidental death)',
  'Surat keterangan kepolisian untuk klaim akibat kecelakaan (Police Statement for accident claim)',
  'Surat diagnosa dari dokter/rumah sakit (Doctor/Hospital diagnosis letter)',
  'Surat pernyataan kesehatan / SPK (Health declaration form)',
  'Formulir klaim dari perusahaan asuransi (Insurance claim form)',
  'Laporan resume medis dokter/rumah sakit tentang perawatan/pembedahan peserta (Medical summary report from doctor/hospital regarding treatment/surgery)',
  'Lain-lain (Others)',
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

/**
 * Dropdown master (pengganti tombol `Choose*`, keputusan work owner 02-10-2026): baris yang dirender paling banyak.
 * Server diminta satu baris lebih untuk mengetahui daftar terpotong - master besar (`CLIENT` ratusan ribu baris)
 * disaring lewat `Search` di dalam dropdown, bukan dimuat utuh.
 */
export const BATAS_DROPDOWN = 200

/** Potong jawaban server (`batas` = BATAS_DROPDOWN + 1): `lebih` = masih ada baris yang tidak dirender. */
export function potongPilihan<T>(daftar: readonly T[]): { tampil: T[]; lebih: boolean } {
  return { tampil: daftar.slice(0, BATAS_DROPDOWN), lebih: daftar.length > BATAS_DROPDOWN }
}

/** Indeks aktif sesudah panah/Home/End (`langkah` ±1 / ±Infinity), dijepit ke daftar; daftar kosong = -1. */
export function geserAktif(aktif: number, langkah: number, n: number): number {
  if (n <= 0) return -1
  return Math.min(n - 1, Math.max(0, aktif + langkah))
}

/**
 * Mode lihat - angka seperti layar Pega (`Ceding's Limit` 250.000.000): pemisah ribuan titik, desimal koma. Dihitung
 * dari TEKS kanonik (`-?digit[.digit]`), tanpa float; teks lain (data lama) tampil apa adanya.
 */
export function tampilAngka(teks: string): string {
  const m = /^(-?)(\d+)(?:\.(\d+))?$/.exec(teks.trim())
  if (m === null) return teks
  const bulat = (m[2] ?? '').replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  return `${m[1] ?? ''}${bulat}${m[3] !== undefined ? `,${m[3]}` : ''}`
}

/** Mode lihat - tanggal `YYYY-MM-DD` tampil `DD/MM/YYYY` (`Begin Date` 01/08/2023 di Pega); teks lain apa adanya. */
export function tampilTanggal(teks: string): string {
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(teks.trim())
  return m === null ? teks : `${m[3]}/${m[2]}/${m[1]}`
}

/** Copy Old - `Search` menyaring ID, Product Name, Ceding, Treaty Number, Treaty Name (tanpa membedakan huruf). */
export function saringLama(daftar: readonly ProdukLama[], kata: string): ProdukLama[] {
  const k = kata.trim().toUpperCase()
  if (k === '') return [...daftar]
  return daftar.filter((d) => [d.id, d.productName, d.ceding, d.treatyNumber, d.inwardName].some((v) => v.toUpperCase().includes(k)))
}

/** Copy Old - ID yang dikirim `Process Copy`: terpilih DAN boleh disalin, urutan daftar. */
export function idBolehDisalin(daftar: readonly ProdukLama[], terpilih: ReadonlySet<string>): string[] {
  return daftar.filter((d) => d.bolehDisalin && terpilih.has(d.id)).map((d) => d.id)
}

/** Copy Old - cacah hasil `Process Copy` per status. */
export function ringkasSalin(hasil: readonly HasilSalinLama[]): Record<StatusSalinLama, number> {
  const r: Record<StatusSalinLama, number> = { disalin: 0, sudahAda: 0, ditolak: 0, gagal: 0 }
  for (const h of hasil) r[h.status]++
  return r
}

