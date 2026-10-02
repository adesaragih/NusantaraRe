// modul/masterproductnamelife/frontend/api.ts - panggilan backend modul Master Product Name Life, satu fungsi
// per rute (`backend/handlers/rute_*.go`). Klien HTTP-nya `inti/klien.ts`.
//
// ⛔ Angka (uang, persen, usia, hari) TEKS sepanjang jalan - tidak pernah `Number` (ADR-0003). Kosong = kosong
// (`""`), bukan nol. Tanggal `YYYY-MM-DD`.
// ⛔ Identitas produk baru tidak pernah dikirim klien: baru = POST, ubah = PUT ke `/{id}` (ADR-0006).

import { BATAS_WAKTU_MS, kegagalanDari, minta, mintaFormulir, rakitURL, unduhBerkasBeridentitas } from '../../../inti/frontend/klien'
import { simpanBlob } from '../../../inti/frontend/lib/simpanBlob'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_MPNL = '/api/master-product-name-life'

// ---------------------------------------------------------------------------
// Bentuk jawaban - SAMA dengan `backend/models/mpnl_produk.go`.
// ---------------------------------------------------------------------------

/** Satu baris grid daftar (`BrowseProduct_Life`, urut ID). */
export interface RingkasanProduk {
  id: string
  ceding: string
  treatyNumber: string
  /** `Treaty Name` - kunci JSON `INWARDNAME`. */
  inwardName: string
  createOp: string
  updateOp: string
}

/** Halaman `ProductName` - sisi umum. */
export interface ProdukUmum {
  productName: string
  ceding: string
  cedingId: string
  sobName: string
  sobId: string
  /** `Deduction (%)` - kunci JSON `RICOMM`. */
  riComm: string
  riRisk: string
  riRiskId: string
  /** `Treaty Name`. */
  inwardName: string
  treatyNumber: string
  cause: string
  causeId: string
  policyHolder: string
  policyHolderName: string
  createOp: string
  updateOp: string
  /** Checkbox `On Retention` - kunci JSON `IsORS`. */
  isOrs: boolean
  /** `Comment` popup `SaveProductName_Confirm`. */
  comment: string
  // Medan layar mati (`1=2`) - dari JSON lama, dipakai `Generate`; server tidak menerimanya dari klien.
  type: string
  typeCeding: string
  grup: string
  productCode: string
  productType: string
  productTypeId: string
  riRate: string
  riRateId: string
  riCommId: string
  outwardName: string
  outwardNameId: string
  outwardRate: string
  outwardRateId: string
  outwardComm: string
  outwardCommId: string
  benefit: string
  benefitId: string
}

/** Halaman `ProductNameInward` - sisi inward. */
export interface ProdukInward {
  id: string
  productId: string
  policyHolder: string
  policyHolderName: string
  insured: string
  addendumNo: string
  addendumWord: string
  amandementNo: string
  amandementSchd: string
  maxExpiredClaim: string
  begin: string
  stnc: string
  cedingRetentionNum: string
  cedingLimit: string
  brokerage: string
  minAge: string
  maxAge: string
  expiryAge: string
  extraPremi: string
  minSumInsured: string
  maxSumInsured: string
  maxSumReasured: string
  rnmShare: string
  rnmLimitNum: string
  premiumFactor: string
  payment: string
  subjectTo: string
  annuityInterest: string
  premiumRefundFactor: string
  maxDataReceive: string
  mature: string
  birthday: string
  currency: string
  currencyId: string
  extraMortality: string
  maxContract: string
  proportionalTable: string
  // Kunci view tanpa medan form - dari JSON lama.
  ceding: string
  treatyNumber: string
  inwardTreatyNm: string
  cedingRetentionPct: string
  cedingLimitXpn: string
  rnmLimitPct: string
  lienClause: string
  months: string
}

/** `asli` - kunci baris JSON lama yang tidak dikelola layar; dibawa bolak-balik apa adanya. */
interface Asli {
  asli?: string
}

export interface BarisLien extends Asli {
  usia: string
  manfaat: string
}

export interface BarisDokumen extends Asli {
  document: string
}

export interface BarisPlan extends Asli {
  plan: string
  planId: string
  name: string
  benefit: string
  riRate: string
  riRateId: string
}

export interface BarisFinUW extends Asli {
  minInsured: string
  maxInsured: string
  employee: string
  nonEmployee: string
}

export interface BarisUWLimit extends Asli {
  minInsured: string
  maxInsured: string
  minAge: string
  maxAge: string
  medical: string
  description: string
}

export interface BarisOutward extends Asli {
  reinsTypeId: string
  reinsTypeName: string
  transactionYear: string
  treatyContractId: string
  underwritingYear: string
  ovrComm: string
}

export interface BarisKomentar extends Asli {
  /** `@CurrentDateTime()` Pega, apa adanya (`YYYYMMDDTHHMMSS.SSS GMT`). */
  date: string
  operatorName: string
  isApproved: string
  suggest: string
}

/** Satu produk utuh - `GET …/produk/{id}`. */
export interface Produk {
  id: string
  umum: ProdukUmum
  inward: ProdukInward
  lienClause: BarisLien[]
  documentClaim: BarisDokumen[]
  planList: BarisPlan[]
  financialUnderwriting: BarisFinUW[]
  underwritingLimit: BarisUWLimit[]
  outwardList: BarisOutward[]
  commentList: BarisKomentar[]
  /** Permintaan simpan saja: `Copy` b59854 - produk asal yang medan milik server-nya diwarisi. */
  salinanDari?: string
  /** Permintaan simpan saja: checkbox `On Retention` b47312 diubah - `OutwardList` dihitung ulang server. */
  hitungOutward?: boolean
}

/** Satu baris master pemilih (`ID` / `Name`). */
export interface NilaiMaster {
  id: string
  nama: string
}

/** Jenis pemilih master - SAMA dengan `models.JenisMaster`. */
export type JenisMaster = 'ceding' | 'sob' | 'pemegang-polis' | 'mata-uang' | 'ri-risk' | 'ri-rate' | 'penyebab'

/** Satu jenis plan `PRODUCT_TYPE_LIFE` - autocomplete `Plan Name` b33121. */
export interface JenisPlan {
  id: string
  coverName: string
  business: string
  benefit: string
}

/** Status lampiran - SAMA dengan `models.Status*`. */
export type StatusLampiran = 'terunggah' | 'gagal' | 'belum'

/** Satu lampiran `M_ATTACHMENTPRODUCTNAME`. */
export interface Lampiran {
  id: string
  produkId: string
  category: string
  fileName: string
  /** Ekstensi huruf kecil (`InsertGoogleStorage_Act` 3 b566). */
  fileMimeType: string
  userName: string
  storageId: string
  status: StatusLampiran
  galat?: string
}

/** Jawaban daftar sederhana (`jawabanDaftar` di handlers). */
export interface Daftar<T> {
  daftar: T[]
  total: number
}

const e = encodeURIComponent

// ---------------------------------------------------------------------------
// Baca.
// ---------------------------------------------------------------------------

/** Grid halaman awal `InboxProductName` (mode daftar). */
export async function ambilDaftarProduk(): Promise<Daftar<RingkasanProduk>> {
  return minta<Daftar<RingkasanProduk>>(`${PREFIX_MPNL}/produk`)
}

/** Tombol `View` b74798 - satu produk, kedua sisi dan ketujuh daftar. */
export async function ambilProduk(id: string): Promise<Produk> {
  return minta<Produk>(`${PREFIX_MPNL}/produk/${e(id)}`)
}

// ---------------------------------------------------------------------------
// Tulis (paket 3-9).
// ---------------------------------------------------------------------------

/** `Save` b59041 → `SaveProductName_Act`: baru = POST (tanpa id), ubah = PUT `/{id}`. */
export async function simpanProduk(p: Produk): Promise<Produk> {
  if (p.id === '') return minta<Produk>(`${PREFIX_MPNL}/produk`, { metode: 'POST', badan: p })
  return minta<Produk>(`${PREFIX_MPNL}/produk/${e(p.id)}`, { metode: 'PUT', badan: p })
}

/** Nama berkas `Generate` - `FileName=SeeDetail` b1335. */
export const BERKAS_GENERATE = 'SeeDetail.csv'

/**
 * `Generate` b60122 → `GenerateUpload_Act`: CSV dari isi form saat itu (POST ber-badan, jadi bukan
 * `unduhBerkasBeridentitas` yang hanya GET; preseden `ambilIsiDokumen` Claim Life). Galat = amplop `{galat}`.
 */
export async function unduhGenerate(p: Produk): Promise<void> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jawab = await fetch(rakitURL(`${PREFIX_MPNL}/produk/generate`), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...headerIdentitas() },
      body: JSON.stringify(p),
      signal: kendali.signal,
    })
    if (!jawab.ok) throw kegagalanDari(jawab.status, await jawab.text())
    simpanBlob(await jawab.blob(), BERKAS_GENERATE)
  } finally {
    clearTimeout(jam)
  }
}

// ---------------------------------------------------------------------------
// Pemilih master (paket 2, 6).
// ---------------------------------------------------------------------------

/** Saran autocomplete paling banyak sekian baris (`batas`); grid pemilih membaca seluruh hasil RD. */
export const BATAS_SARAN = 20

/** Tombol `Choose*` / autocomplete - `Search` diubah huruf besar di server (`SearchPolicyHolder_act` b236). */
export async function cariMaster(jenis: JenisMaster, cari: string, batas?: number): Promise<Daftar<NilaiMaster>> {
  return minta<Daftar<NilaiMaster>>(`${PREFIX_MPNL}/master/${e(jenis)}`, { kueri: { cari, batas } })
}

/** Autocomplete `Plan Name` b33121 (RD `BrowseProductTypeLife_RD`). */
export async function cariPlan(cari: string): Promise<Daftar<JenisPlan>> {
  return minta<Daftar<JenisPlan>>(`${PREFIX_MPNL}/master-plan`, { kueri: { cari } })
}

/** Satu baris dialog `View Rate` - kolom ID, USEDBY, GENDER, CONTRACT, AGE, RATE (teks apa adanya). */
export interface BarisRate {
  id: string
  usedBy: string
  gender: string
  contract: string
  age: string
  rate: string
}

/**
 * `View Rate` b34113 - view `RATE_LIFE` baca saja (K1 keputusan work owner 01-10-2026, OQ-MPNL-03), disaring
 * `RIRATEID` baris plan. `terpotong` = view memuat lebih dari 500 baris (`pyMaxRecords` 500, seperti Pega).
 */
export async function ambilRate(riRateId: string): Promise<Daftar<BarisRate> & { terpotong: boolean }> {
  return minta<Daftar<BarisRate> & { terpotong: boolean }>(`${PREFIX_MPNL}/rate`, { kueri: { riRateId } })
}

// ---------------------------------------------------------------------------
// Lampiran (paket 8).
// ---------------------------------------------------------------------------

const lampiran = (produkID: string): string => `${PREFIX_MPNL}/produk/${e(produkID)}/lampiran`

/** `Refresh` b65270 / `View` → `LoadAttachmentProdName`. */
export async function ambilLampiran(produkID: string): Promise<Daftar<Lampiran>> {
  return minta<Daftar<Lampiran>>(lampiran(produkID))
}

/** `Add attachment` b64747 → `Attach` b24 (`ProductNameSaveAttachment`) - multipart `berkas`. */
export async function unggahLampiran(produkID: string, berkas: File | null): Promise<Lampiran> {
  const isi = new FormData()
  if (berkas !== null) isi.append('berkas', berkas)
  return mintaFormulir<Lampiran>(lampiran(produkID), isi)
}

/** Kirim ulang lampiran yang gagal (tiket 09). */
export async function ulangiLampiran(produkID: string, id: string): Promise<Lampiran> {
  return minta<Lampiran>(`${lampiran(produkID)}/${e(id)}/ulangi`, { metode: 'POST' })
}

/** Tautan nama berkas b68903 (`DownloadAttProdName_Act`). */
export async function unduhLampiran(produkID: string, l: Lampiran): Promise<void> {
  return unduhBerkasBeridentitas(`${lampiran(produkID)}/${e(l.id)}/unduh`, l.fileName)
}

/** `Download All` b67657 - zip lampiran produk ini (RALAT R15). */
export async function unduhSemuaLampiran(produkID: string): Promise<void> {
  return unduhBerkasBeridentitas(`${lampiran(produkID)}/unduh-semua`, `lampiran-${produkID}.zip`)
}

/** `View Office Online` b69291 - stub: server menjawab 503 berkalimat (OQ-MPNL-11). */
export async function lihatOffice(produkID: string, id: string): Promise<void> {
  await minta<unknown>(`${lampiran(produkID)}/${e(id)}/office`)
}

/** `Delete` b69714 (`DeleteAttacProdName_act`). */
export async function hapusLampiran(produkID: string, id: string): Promise<void> {
  await minta<unknown>(`${lampiran(produkID)}/${e(id)}`, { metode: 'DELETE' })
}
