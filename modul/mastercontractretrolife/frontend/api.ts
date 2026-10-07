// modul/mastercontractretrolife/frontend/api.ts - panggilan backend modul Master Contract Retro Life, satu
// fungsi per rute (`backend/handlers/rute_mcrl.go`, `rute_tulis.go`). Klien HTTP-nya `inti/klien.ts`.
//
// ⛔ Angka (share, komisi, batas proteksi) TEKS sepanjang jalan - tidak pernah
// `Number` (ADR-0003). Kosong = kosong (`""`), bukan nol.
// ⛔ Identitas baris baru tidak pernah dikirim klien: baru = POST ke induk,
// ubah = PUT ke `/{id}` (ADR-0006).

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_MCRL = '/api/master-contract-retro-life'

// ---------------------------------------------------------------------------
// Bentuk jawaban - SAMA dengan `MarshalJSON` di `backend/models/mcrl_entitas.go`.
// Tanggal `YYYY-MM-DD`, `tglUpdate` `YYYY-MM-DD HH:MM:SS`, kosong = "".
// ---------------------------------------------------------------------------

/** Satu baris `TREATYYEAR_LIFE` - grid `BrowseTreatyYear_Life_RD`. */
export interface TahunTreaty {
  id: string
  /** Kolom `TREATYYEAR` - label layar **TRANSACTION YEAR**. */
  treatyYear: string
  underwritingYear: string
  startDate: string
  endDate: string
  userId: string
  tglUpdate: string
}

/** Satu baris `TREATYCONTRACT_LIFE` - grid `BrowseTreatyContract_Life_RD`. */
export interface Kontrak {
  id: string
  idTreatyYear: string
  reinsTypeId: string
  reinsTypeName: string
  treatyStartDate: string
  treatyEndDate: string
  userId: string
  tglUpdate: string
  /** `MAXIMUM LIMIT (IDR)`. */
  idr: string
  /** `MAXIMUM LIMIT (USD)`. */
  usd: string
  /** `MINIMUM LIMIT (IDR)`. */
  bIdr: string
  /** `MINIMUM LIMIT (USD)`. */
  bUsd: string
  idrSelisih: string
  usdSelisih: string
}

/** Satu baris `TREATYREINSURER_LIFE` - grid `BrowseDetailTreatyReisurerLife_RD`. */
export interface Reinsurer {
  id: string
  treatyYearId: string
  treatyContractId: string
  reinsTypeId: string
  reinsTypeName: string
  reinsurerId: string
  reinsurerName: string
  pctShare: string
  /** Kolom `COMMISION` - label layar **(%) DISCOUNT**. */
  komisi: string
  ovrComm: string
  userId: string
  tglUpdate: string
}

/** Satu baris `TREATYSECURITYREINSURER_LIFE` - grid `BrowseSecurityReinsurer_Life_RD`. */
export interface SecurityReinsurer {
  id: string
  treatyYearId: string
  treatyContractId: string
  treatyReinsurerId: string
  reinsurerId: string
  reinsurerName: string
  pctShare: string
  userId: string
  tglUpdate: string
}

/** Satu baris `TREATYBUSINESS_LIFE` - grid `BrowseTreatyBusiness_Life_RD`. */
export interface Business {
  id: string
  treatyYearId: string
  treatyYear: string
  treatyContractId: string
  reinsTypeId: string
  reinsTypeName: string
  bizCode: string
  bizName: string
  /** Isi autocomplete `R/I RATE` - dipakai `View Rate` sebagai `idusedby`. */
  riRateId: string
  /** `R/I RATE` - nama tabel rate, teks apa adanya (Pertanyaan A, RALAT R7). */
  riRate: string
  userId: string
  tglUpdate: string
}

/** Master `REINSURANCETYPE` (`.Flag = 1`) - dropdown `REINS TYPE`; tampil `.Note`. */
export interface JenisReasuransi {
  id: string
  note: string
}

/** Master `AGENT` (`BrowseCedingCoLife_RD`) - autocomplete `REINSURER NAME`. */
export interface MasterReinsurer {
  id: string
  clientName: string
}

/** Master `BUSINESS` (`BrowseBusinessLife_RD`) - autocomplete `BUSINESS NAME`. */
export interface MasterBusiness {
  id: string
  note: string
  oldId: string
}

/** Jawaban daftar sederhana (`jawabanDaftar` di handlers). */
export interface Daftar<T> {
  daftar: T[]
  total: number
}

/** `GET …/tahun/{id}/kontrak` - kepala panel + grid. */
export interface JawabanKontrak {
  tahun: TahunTreaty
  daftar: Kontrak[]
}

/** `GET …/kontrak/{id}/reinsurer` - total share dihitung server (desimal persis). */
export interface JawabanReinsurer {
  kontrak: Kontrak
  daftar: Reinsurer[]
  totalShare: string
  /** Tiket 05: total ≠ 100 DITANDAI, tidak memblokir. */
  totalBukan100: boolean
}

/** `GET …/reinsurer/{id}/security` - eksposur per ID baris security (OQ-MCRL-08). */
export interface JawabanSecurity {
  induk: Reinsurer
  daftar: SecurityReinsurer[]
  eksposur: Record<string, string>
  /** Total share security - seperti Reinsurer List (04-10-2026). Opsional: server lama tidak mengirimnya. */
  totalShare?: string
  totalBukan100?: boolean
}

/** `GET …/kontrak/{id}/business`. */
export interface JawabanBusiness {
  kontrak: Kontrak
  daftar: Business[]
}

/** Cacah baris yang ikut terhapus - `models.Dampak`. */
export interface Dampak {
  security: number
  reinsurer: number
  business: number
}

/** Empat jenis baris yang dapat dihapus; tahun treaty TIDAK (tahun abadi). */
export type JenisHapus = 'kontrak' | 'reinsurer' | 'security' | 'business'

/** `GET …/{jenis}/{id}/dampak-hapus` - isi popup konfirmasi. */
export interface JawabanDampak {
  jenis: JenisHapus
  id: string
  dampak: Dampak
}

/** Jawaban hapus - `pesan` VERBATIM korpus (`Data Berhasil di Hapus`, dst.). */
export interface HasilHapus {
  pesan: string
  terhapus: Dampak
}

/** `GET …/business/{id}/salin-semua` - pratinjau `Copy to all Reinstype`. */
export interface PratinjauSalin {
  business: Business
  reinsTypeId: string
  sasaran: Kontrak[]
}

/** `POST …/business/{id}/salin-semua` - pesan VERBATIM `Copied to all reins types.`. */
export interface HasilSalin {
  pesan: string
  jumlah: number
  baru: Business[]
}

// ---------------------------------------------------------------------------
// Badan simpan - SAMA dengan `services.*Masuk`.
// ---------------------------------------------------------------------------

export interface TahunMasuk {
  id: string
  treatyYear: string
  underwritingYear: string
  startDate: string
  endDate: string
}

export interface KontrakMasuk {
  id: string
  reinsTypeId: string
  idr: string
  usd: string
  bIdr: string
  bUsd: string
}

export interface ReinsurerMasuk {
  id: string
  reinsurerName: string
  reinsurerId: string
  pctShare: string
  komisi: string
  ovrComm: string
}

export interface SecurityMasuk {
  id: string
  reinsurerName: string
  reinsurerId: string
  pctShare: string
}

export interface BusinessMasuk {
  id: string
  bizCode: string
  bizName: string
  riRateId: string
  riRate: string
}

/** Metode + jalur satu simpan. */
export interface RencanaSimpan {
  metode: 'POST' | 'PUT'
  jalur: string
}

const e = encodeURIComponent

/**
 * Satu `Save` Pega ber-upsert menjadi DUA rute: baris tanpa `id` = POST ke
 * induknya, baris ber-`id` = PUT ke `/{segmen}/{id}`.
 *
 * `induk` null = tahun treaty (akar, tanpa induk).
 */
export function rencanaSimpan(segmen: string, induk: { segmen: string; id: string } | null, id: string): RencanaSimpan {
  if (id !== '') return { metode: 'PUT', jalur: `${PREFIX_MCRL}/${segmen}/${e(id)}` }
  if (induk === null) return { metode: 'POST', jalur: `${PREFIX_MCRL}/${segmen}` }
  return { metode: 'POST', jalur: `${PREFIX_MCRL}/${induk.segmen}/${e(induk.id)}/${segmen}` }
}

async function simpan<T>(r: RencanaSimpan, badan: unknown): Promise<T> {
  return minta<T>(r.jalur, { metode: r.metode, badan })
}

// ---------------------------------------------------------------------------
// Baca.
// ---------------------------------------------------------------------------

/** Grid halaman awal - `BrowseTreatyYear_Life_RD`, urut `.ID ASC`. */
export async function ambilTahun(): Promise<Daftar<TahunTreaty>> {
  return minta<Daftar<TahunTreaty>>(`${PREFIX_MCRL}/tahun`)
}

/** Tombol `ReinsType` - grid kontrak tahun itu. */
export async function ambilKontrak(tahunId: string): Promise<JawabanKontrak> {
  return minta<JawabanKontrak>(`${PREFIX_MCRL}/tahun/${e(tahunId)}/kontrak`)
}

/** Tombol `Reinsurer List` - grid reinsurer + total share kontrak itu. */
export async function ambilReinsurer(kontrakId: string): Promise<JawabanReinsurer> {
  return minta<JawabanReinsurer>(`${PREFIX_MCRL}/kontrak/${e(kontrakId)}/reinsurer`)
}

/** Tombol `Security Reinsurer` - grid security reinsurer itu. */
export async function ambilSecurity(reinsurerId: string): Promise<JawabanSecurity> {
  return minta<JawabanSecurity>(`${PREFIX_MCRL}/reinsurer/${e(reinsurerId)}/security`)
}

/** Tombol `Business List` - grid business kontrak itu. */
export async function ambilBusiness(kontrakId: string): Promise<JawabanBusiness> {
  return minta<JawabanBusiness>(`${PREFIX_MCRL}/kontrak/${e(kontrakId)}/business`)
}

/** Dropdown `REINS TYPE`; master kosong = 503 berkalimat yang menyebut masternya. */
export async function ambilJenisReasuransi(): Promise<Daftar<JenisReasuransi>> {
  return minta<Daftar<JenisReasuransi>>(`${PREFIX_MCRL}/jenis-reasuransi`)
}

/** Autocomplete `REINSURER NAME` / `SECURITY REINSURER NAME`; `""` = daftar awal. */
export async function cariMasterReinsurer(cari: string): Promise<Daftar<MasterReinsurer>> {
  return minta<Daftar<MasterReinsurer>>(`${PREFIX_MCRL}/master-reinsurer`, { kueri: { cari } })
}

/** Autocomplete `BUSINESS NAME`. */
export async function cariMasterBusiness(cari: string): Promise<Daftar<MasterBusiness>> {
  return minta<Daftar<MasterBusiness>>(`${PREFIX_MCRL}/master-business`, { kueri: { cari } })
}

/**
 * Autocomplete `R/I RATE` - view `RATE_LIFE_SUMMARY`, baca saja (K1 keputusan work owner 01-10-2026,
 * OQ-MCRL-13). View tak terbaca = 503 berkalimat yang menyebut view-nya; pesannya tampil apa adanya.
 */
export async function cariRingkasanRate(cari: string): Promise<Daftar<{ id: string; usedBy: string }>> {
  return minta<Daftar<{ id: string; usedBy: string }>>(`${PREFIX_MCRL}/ringkasan-rate`, { kueri: { cari } })
}

/** Satu baris section `ViewRate` (`Rate List`) - kolom ID, USEDBY, GENDER, CONTRACT, AGE, RATE. */
export interface BarisRate {
  id: string
  usedBy: string
  gender: string
  contract: string
  age: string
  rate: string
}

/**
 * Tombol `View Rate` - `idusedby` = `RIRATEID`, view `RATE_LIFE` baca saja (K1, OQ-MCRL-13).
 * `terpotong` = view memuat lebih dari 500 baris (`BrowseRateLife_RD` `pyMaxRecords` 500, seperti Pega).
 */
export async function ambilRate(idusedby: string): Promise<Daftar<BarisRate> & { terpotong: boolean }> {
  return minta<Daftar<BarisRate> & { terpotong: boolean }>(`${PREFIX_MCRL}/rate`, { kueri: { idusedby } })
}

// ---------------------------------------------------------------------------
// Tulis.
// ---------------------------------------------------------------------------

/** `Save` form `Input New Data`. */
export async function simpanTahun(m: TahunMasuk): Promise<TahunTreaty> {
  return simpan<TahunTreaty>(rencanaSimpan('tahun', null, m.id), m)
}

/** `Save` form kontrak (`SaveTreatyLimit_Act`). */
export async function simpanKontrak(tahunId: string, m: KontrakMasuk): Promise<Kontrak> {
  return simpan<Kontrak>(rencanaSimpan('kontrak', { segmen: 'tahun', id: tahunId }, m.id), m)
}

/** `Save` form reinsurer (`SaveSecurityLife_Act`). */
export async function simpanReinsurer(kontrakId: string, m: ReinsurerMasuk): Promise<Reinsurer> {
  return simpan<Reinsurer>(rencanaSimpan('reinsurer', { segmen: 'kontrak', id: kontrakId }, m.id), m)
}

/** `Save` form security (`SaveSecurityReinsurerLife_Act`). */
export async function simpanSecurity(reinsurerId: string, m: SecurityMasuk): Promise<SecurityReinsurer> {
  return simpan<SecurityReinsurer>(rencanaSimpan('security', { segmen: 'reinsurer', id: reinsurerId }, m.id), m)
}

/** `Save` form business (`SaveBusinessLife_Act`). */
export async function simpanBusiness(kontrakId: string, m: BusinessMasuk): Promise<Business> {
  return simpan<Business>(rencanaSimpan('business', { segmen: 'kontrak', id: kontrakId }, m.id), m)
}

/** Pratinjau `Copy to all Reinstype` - nol tulisan (penyimpangan 5). */
export async function pratinjauSalinSemua(businessId: string): Promise<PratinjauSalin> {
  return minta<PratinjauSalin>(`${PREFIX_MCRL}/business/${e(businessId)}/salin-semua`)
}

/** Konfirmasi `Copy to all Reinstype` - sasaran yang DILIHAT dikirim; server menolak (409) bila sudah lain. */
export async function salinSemua(businessId: string, sasaran: readonly string[]): Promise<HasilSalin> {
  return minta<HasilSalin>(`${PREFIX_MCRL}/business/${e(businessId)}/salin-semua`, {
    metode: 'POST',
    badan: { sasaran },
  })
}

/** Isi popup konfirmasi hapus - cacah anak yang ikut terhapus. */
export async function ambilDampakHapus(jenis: JenisHapus, id: string): Promise<JawabanDampak> {
  return minta<JawabanDampak>(`${PREFIX_MCRL}/${jenis}/${e(id)}/dampak-hapus`)
}

/** `Yes` di popup - cacahan yang dilihat dikirim; server menolak (409) bila sudah lain. */
export async function hapus(jenis: JenisHapus, id: string, dampak: Dampak): Promise<HasilHapus> {
  return minta<HasilHapus>(`${PREFIX_MCRL}/${jenis}/${e(id)}`, { metode: 'DELETE', badan: { dampak } })
}
