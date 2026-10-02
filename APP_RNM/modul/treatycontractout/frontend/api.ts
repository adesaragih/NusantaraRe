// modul/treatycontractout/frontend/api.ts - panggilan backend modul Treaty Contract Out, satu
// fungsi per endpoint. Klien HTTP-nya `inti/klien.ts` (refactor bentuk B,
// 30-09-2026: dipecah dari `services/api.ts` tanpa mengubah satu panggilan pun).

import { minta, mintaFormulir } from '../../../inti/frontend/klien'

// ---------------------------------------------------------------------------
// TREATY CONTRACT OUT — tiket 02: master jenis reasuransi non-life.
//
// ⚠️ Bagian ini ditambahkan ADITIF oleh sesi Treaty Contract Out. Bentuknya
// dari `internal/services/tco_jenisreasuransi.go` (`JenisReasuransi`).
// ---------------------------------------------------------------------------

/**
 * Satu jenis reasuransi dari master `REINSURANCETYPE`, SUDAH tersaring
 * server (12 awalan blacklist, Flag active, Type 1/2/3 — RD dominan yang
 * dipakai 11 grid klausul). Nama kunci mengikuti properti Pega `.ID`,
 * `.Note`, `.Type`.
 */
export interface JenisReasuransiTreaty {
  /** ⛔ TEKS: kode, bukan bilangan. `00007` harus tetap `00007`. */
  id: string
  note: string
  tipe: string
}

/** Jawaban `GET /api/treaty-contract-out/jenis-reasuransi`. */
export interface DaftarJenisReasuransiTreaty {
  daftar: JenisReasuransiTreaty[]
  total: number
}

/**
 * Daftar jenis reasuransi non-life — SATU tempat untuk layar kontrak maupun
 * seluruh grid klausul (tiket 02 AC). Master kosong menjawab **503** dengan
 * pesan yang menyebut masternya (ADR-0015); pesannya tampil apa adanya.
 */
export async function ambilJenisReasuransiTreaty(): Promise<DaftarJenisReasuransiTreaty> {
  return minta<DaftarJenisReasuransiTreaty>('/api/treaty-contract-out/jenis-reasuransi')
}

/**
 * `AturanKlausul.pilihanReins` baris anak Treaty Limit — SAMA dengan
 * `models.PilihanReinsAnakTreatyLimit` di backend (dijaga uji).
 */
export const PILIHAN_REINS_ANAK_TREATY_LIMIT = 'anak-treaty-limit'

/**
 * Pilihan ReinsType baris anak SEMUA grid `Show Child` — `TreatyContractSetReinsTypeList` atas nama
 * ReinsType baris induk (`CARIDESCFACIN ← .ReinsTypeName`) [keputusan work owner 02-10-2026].
 */
export async function ambilJenisReasuransiAnakTreatyLimit(namaInduk: string): Promise<DaftarJenisReasuransiTreaty> {
  return minta<DaftarJenisReasuransiTreaty>('/api/treaty-contract-out/jenis-reasuransi/anak-treaty-limit', {
    kueri: { namaInduk },
  })
}

// ---------------------------------------------------------------------------
// TREATY CONTRACT OUT — tiket 03: tahun treaty dan master grup treaty.
// Bentuknya dari `internal/services/tco_tahun.go` dan `tco_gruptreaty.go`.
// ---------------------------------------------------------------------------

/** Satu grup treaty dari master `TREATYGROUP` (`.ID`, `.TreatyGroupName`), urutan ID DESC. */
export interface GrupTreaty {
  id: string
  treatyGroupName: string
}

/** `GET /api/treaty-contract-out/grup-treaty`; master kosong menjawab 503 (ADR-0015). */
export async function ambilGrupTreaty(): Promise<{ daftar: GrupTreaty[]; total: number }> {
  return minta<{ daftar: GrupTreaty[]; total: number }>('/api/treaty-contract-out/grup-treaty')
}

/**
 * Satu tahun treaty (`TREATYYEAR`). Nama kunci mengikuti `InputTreatyYear.*`.
 * Tanggal TEKS `YYYY-MM-DD` (kosong = kosong); `tglUpdate` `YYYY-MM-DD HH:MM:SS`.
 * `proportion` menyimpan pilihan "Reinsurance Type" (`.ID` master) — OQ-TCO-04.
 */
export interface TahunTreaty {
  id: string
  treatyYear: string
  underwritingYear: string
  treatyGroupId: string
  treatyGroupName: string
  proportion: string
  startDate: string
  endDate: string
  userId: string
  tglUpdate: string
}

/** Satu halaman daftar tahun treaty, terbaru dahulu (`.ID DESC`). */
export interface HalamanTahunTreaty {
  baris: TahunTreaty[]
  total: number
  halaman: number
  ukuran: number
}

/** Badan simpan. `id` kosong = baru (POST); terisi = perbarui (PUT /{id}). */
export interface TahunTreatyMasuk {
  id: string
  treatyYear: string
  underwritingYear: string
  treatyGroupId: string
  treatyGroupName: string
  proportion: string
  startDate: string
  endDate: string
}

/** `GET /api/treaty-contract-out/tahun?halaman=&ukuran=`. */
export async function ambilTahunTreaty(halaman = 1, ukuran = 20): Promise<HalamanTahunTreaty> {
  return minta<HalamanTahunTreaty>('/api/treaty-contract-out/tahun', { kueri: { halaman, ukuran } })
}

/** `GET /api/treaty-contract-out/tahun/{id}`; 404 bila tidak ada. */
export async function ambilSatuTahunTreaty(id: string): Promise<TahunTreaty> {
  return minta<TahunTreaty>(`/api/treaty-contract-out/tahun/${encodeURIComponent(id)}`)
}

/**
 * Simpan tahun treaty — tombol `Save` (`InputDtlTreatyContact.xml` b10332).
 *
 * ⛔ Identitas baris baru TIDAK dikirim (AC 5): `id` kosong → POST, server
 * menerbitkannya dari sequence. `id` terisi → PUT /{id}, SELURUH medan
 * tertimpa (AC 8). 409 = periode + grup sudah dipakai baris lain (AC 73);
 * 422 = gerbang wajib isi / periode terbalik (AC 9); pesannya tampil apa adanya.
 */
export async function simpanTahunTreaty(masuk: TahunTreatyMasuk): Promise<TahunTreaty> {
  if (masuk.id === '') {
    return minta<TahunTreaty>('/api/treaty-contract-out/tahun', { metode: 'POST', badan: masuk })
  }
  return minta<TahunTreaty>(`/api/treaty-contract-out/tahun/${encodeURIComponent(masuk.id)}`, {
    metode: 'PUT',
    badan: masuk,
  })
}

// ---------------------------------------------------------------------------
// Treaty Contract Out tiket 12 — lampiran tahun treaty (FITUR BARU).
// ---------------------------------------------------------------------------

/** Status lampiran — kosakata kami, fiturnya tidak ada di Pega. */
export type StatusLampiranTahun = 'terkirim' | 'tertunda' | 'gagal'

/** Satu lampiran tahun treaty. Kunci berkas di penyimpanan TIDAK dikirim backend. */
export interface LampiranTahun {
  id: string
  idTreatyYear: string
  fileName: string
  fileMimeType: string
  category: string
  userId: string
  /** tco4: dibaca dari ID warisan (stempel `YYYYMMDDHH24MISSFF3`). Ukuran berkas tidak disimpan tabel warisan. */
  tglUpload: string
  status: StatusLampiranTahun
  /** Cacah percobaan efek unggah terakhir. */
  percobaan: number
  /** Galat efek unggah terakhir; kosong bila terkirim. */
  galat: string
}

/** Jawaban unggah / ulangi. `peringatan` terisi bila antrean tidak dapat dijalankan. */
export interface HasilLampiranTahun {
  lampiran: LampiranTahun
  peringatan: string
}

/** Satu rekam yang tidak sejalan dengan penyimpanan (AC 61). */
export interface TemuanSelarasLampiran {
  lampiranId: string
  fileName: string
  masalah: string
  perbaikan: 'ulangi' | 'hapus'
}

function jalurLampiranTahun(tahunID: string): string {
  return `/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/lampiran`
}

/** Master kategori `CATEGORY_ATTACH_REAS`; kosong = 503 dari backend. */
export async function ambilKategoriLampiranTCO(): Promise<string[]> {
  const j = await minta<{ daftar: string[] | null }>('/api/treaty-contract-out/kategori-lampiran')
  return j.daftar ?? []
}

/** `Refresh` b1023 — daftar lampiran satu tahun treaty. */
export async function ambilLampiranTahun(tahunID: string): Promise<LampiranTahun[]> {
  const j = await minta<{ daftar: LampiranTahun[] | null }>(jalurLampiranTahun(tahunID))
  return j.daftar ?? []
}

/** `Add attachment` b578 — multipart lewat `mintaFormulir` (identitas ikut). */
export async function unggahLampiranTahun(
  tahunID: string,
  berkas: File,
  kategori: string,
): Promise<HasilLampiranTahun> {
  const isi = new FormData()
  isi.append('berkas', berkas)
  isi.append('kategori', kategori)
  return mintaFormulir<HasilLampiranTahun>(jalurLampiranTahun(tahunID), isi)
}

/** `Delete` b3897. Penghapusan di penyimpanan menyusul lewat outbox. */
export async function hapusLampiranTahun(tahunID: string, lampiranID: string): Promise<{ peringatan: string }> {
  return minta<{ peringatan: string }>(
    `${jalurLampiranTahun(tahunID)}/${encodeURIComponent(lampiranID)}`,
    { metode: 'DELETE' },
  )
}

/** Unggah ulang lampiran yang tertunda / gagal / kehilangan berkasnya (AC 58, 61). */
export async function ulangiLampiranTahun(tahunID: string, lampiranID: string): Promise<HasilLampiranTahun> {
  return minta<HasilLampiranTahun>(
    `${jalurLampiranTahun(tahunID)}/${encodeURIComponent(lampiranID)}/ulangi`,
    { metode: 'POST' },
  )
}

/** Keselarasan rekam dengan penyimpanan (AC 61). */
export async function periksaSelarasLampiran(tahunID: string): Promise<TemuanSelarasLampiran[]> {
  const j = await minta<{ temuan: TemuanSelarasLampiran[] | null }>(`${jalurLampiranTahun(tahunID)}/selaras`)
  return j.temuan ?? []
}

/** Jalur isi satu lampiran — tautan nama berkas b3470 (`TreatyOutDownloadOne`). */
export function jalurIsiLampiran(tahunID: string, lampiranID: string): string {
  return `${jalurLampiranTahun(tahunID)}/${encodeURIComponent(lampiranID)}/isi`
}

/** Jalur arsip seluruh lampiran terkirim — `Download All` b2659. */
export function jalurSemuaLampiran(tahunID: string): string {
  return `${jalurLampiranTahun(tahunID)}/semua`
}

// ---------------------------------------------------------------------------
// Treaty Contract Out tiket 04 — kontrak treaty di dalam tahun treaty.
// ---------------------------------------------------------------------------

/** Satu kontrak — kolom `TREATYCONTRACT`. Tanggal YYYY-MM-DD. */
export interface KontrakTreaty {
  id: string
  idTreatyYear: string
  reinsTypeId: string
  /** Nama dari master `REINSURANCETYPE` — server yang mengisinya, bukan klien. */
  reinsTypeName: string
  treatyStartDate: string
  treatyEndDate: string
  userId: string
  tglUpdate: string
}

/** Badan simpan kontrak; `id` kosong = kontrak baru (POST). */
export interface KontrakMasuk {
  id: string
  reinsTypeId: string
  treatyStartDate: string
  treatyEndDate: string
}

function jalurKontrakTahun(tahunID: string): string {
  return `/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/kontrak`
}

/** Grid kontrak satu tahun treaty (`BrowseTreatyContract_RD`). */
export async function ambilKontrakTahun(tahunID: string): Promise<KontrakTreaty[]> {
  const j = await minta<{ daftar: KontrakTreaty[] | null }>(jalurKontrakTahun(tahunID))
  return j.daftar ?? []
}

/** `Save` b3618 — POST bila baru, PUT /{id} bila ubah. */
export async function simpanKontrakTahun(tahunID: string, masuk: KontrakMasuk): Promise<KontrakTreaty> {
  if (masuk.id === '') {
    return minta<KontrakTreaty>(jalurKontrakTahun(tahunID), { metode: 'POST', badan: masuk })
  }
  return minta<KontrakTreaty>(`${jalurKontrakTahun(tahunID)}/${encodeURIComponent(masuk.id)}`, {
    metode: 'PUT',
    badan: masuk,
  })
}

/**
 * Tanggal akhir bawaan tahun treaty BARU (belum ber-ID) — aturan yang SAMA
 * dengan kontrak (`models.AkhirKontrakBawaanTCO`, mulai + 1 tahun kalender),
 * dihitung server [keputusan work owner 30-09-2026].
 */
export async function ambilAkhirBawaanTahun(mulai: string): Promise<string> {
  const j = await minta<{ endDate: string }>('/api/treaty-contract-out/tahun/akhir-bawaan', { kueri: { mulai } })
  return j.endDate
}

// ---------------------------------------------------------------------------
// Treaty Contract Out tiket 05 — reinsurer pada kombinasi kontrak.
// ---------------------------------------------------------------------------

/** Kunci gabungan anak kontrak (tahun teks, grup, jenis). */
export interface KombinasiTreaty {
  treatyYear: string
  treatyGroupId: string
  treatyGroupName: string
  reinsTypeId: string
  reinsTypeName: string
}

/** Satu reinsurer. ⛔ `pctShare` dan `ricomm` TEKS desimal — tidak pernah `Number`. */
export interface ReinsurerTreaty {
  id: string
  treatyYear: string
  treatyGroupId: string
  treatyGroupName: string
  reinsTypeId: string
  reinsTypeName: string
  reinsurerId: string
  clientId: string
  name: string
  ricomm: string
  pctShare: string
  iuDate: string
  userId: string
  startDate: string
  endDate: string
  statusOn: string
  stdRating: string
  operatorName: string
  tglUpdate: string
}

/** Grid reinsurer + `Total Share -->>`. */
export interface DaftarReinsurer {
  daftar: ReinsurerTreaty[]
  total: number
  /** Teks desimal persis dari server. */
  totalShare: string
  kombinasi: KombinasiTreaty
}

/** Badan simpan — hanya medan yang tampil di form. */
export interface ReinsurerMasuk {
  id: string
  reinsurerId: string
  /** Teks; koma atau titik desimal — dinormalkan server. */
  pctShare: string
  ricomm: string
  stdRating: string
}

/** Satu pilihan master reinsurer (`BrowseAgentReinsSOA_RD`). */
export interface ReinsurerMaster {
  id: string
  clientName: string
  clientId: string
}

function jalurReinsurer(tahunID: string, kontrakID: string): string {
  return `/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/kontrak/${encodeURIComponent(kontrakID)}/reinsurer`
}

/** `Reinsurer List` b11308 — reinsurer kombinasi kontrak + total share. */
export async function ambilReinsurerKombinasi(tahunID: string, kontrakID: string): Promise<DaftarReinsurer> {
  const j = await minta<DaftarReinsurer>(jalurReinsurer(tahunID, kontrakID))
  return { ...j, daftar: j.daftar ?? [] }
}

/** `Save` b11405 — POST bila baru, PUT /{id} bila ubah. */
export async function simpanReinsurerKombinasi(
  tahunID: string,
  kontrakID: string,
  masuk: ReinsurerMasuk,
): Promise<{ reinsurer: ReinsurerTreaty; totalShare: string }> {
  const jalur = jalurReinsurer(tahunID, kontrakID)
  if (masuk.id === '') {
    return minta(jalur, { metode: 'POST', badan: masuk })
  }
  return minta(`${jalur}/${encodeURIComponent(masuk.id)}`, { metode: 'PUT', badan: masuk })
}

/** Pemilih `Reinsurer` b8226 — master aktif yang namanya memuat teks. */
export async function cariReinsurerMaster(teks: string): Promise<ReinsurerMaster[]> {
  const j = await minta<{ daftar: ReinsurerMaster[] | null }>('/api/treaty-contract-out/reinsurer-master', {
    kueri: { cari: teks },
  })
  return j.daftar ?? []
}

// ---------------------------------------------------------------------------
// Treaty Contract Out tiket 07 — business pada kombinasi kontrak.
// ---------------------------------------------------------------------------

/** Satu baris bisnis; `isActive` '1' aktif / '0' nonaktif, `aktif` turunan server. */
export interface BusinessTreaty {
  id: string
  isActive: string
  aktif: boolean
  treatyYear: string
  treatyYearId: string
  treatyGroupId: string
  treatyGroupName: string
  reinsTypeId: string
  reinsTypeName: string
  bizCode: string
  bizName: string
  userId: string
  tglUpdate: string
}

/** Grid bisnis satu kombinasi — aktif maupun nonaktif. */
export interface DaftarBusiness {
  daftar: BusinessTreaty[]
  total: number
  kombinasi: KombinasiTreaty
}

/** Badan simpan — `Business Name` (kode) dan `Active`. */
export interface BusinessMasuk {
  id: string
  bizCode: string
  isActive: string
}

/** Satu pilihan master bisnis. */
export interface BusinessMaster {
  id: string
  note: string
}

function jalurBusiness(tahunID: string, kontrakID: string): string {
  return `/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/kontrak/${encodeURIComponent(kontrakID)}/business`
}

/** `Business List` b10842. */
export async function ambilBusinessKombinasi(tahunID: string, kontrakID: string): Promise<DaftarBusiness> {
  const j = await minta<DaftarBusiness>(jalurBusiness(tahunID, kontrakID))
  return { ...j, daftar: j.daftar ?? [] }
}

/** `Save` b6966 — POST bila baru, PUT /{id} bila ubah. */
export async function simpanBusinessKombinasi(
  tahunID: string,
  kontrakID: string,
  masuk: BusinessMasuk,
): Promise<BusinessTreaty> {
  const jalur = jalurBusiness(tahunID, kontrakID)
  if (masuk.id === '') {
    return minta(jalur, { metode: 'POST', badan: masuk })
  }
  return minta(`${jalur}/${encodeURIComponent(masuk.id)}`, { metode: 'PUT', badan: masuk })
}

/** `Delete` b4826 — pesan "Data Dengan ID … Berhasil di Hapus". */
export async function hapusBusinessKombinasi(tahunID: string, kontrakID: string, id: string): Promise<string> {
  const j = await minta<{ pesan: string }>(`${jalurBusiness(tahunID, kontrakID)}/${encodeURIComponent(id)}`, {
    metode: 'DELETE',
  })
  return j.pesan
}

/** Pemilih `Business Name` b6241. */
export async function ambilBusinessMaster(): Promise<BusinessMaster[]> {
  const j = await minta<{ daftar: BusinessMaster[] | null }>('/api/treaty-contract-out/business-master')
  return j.daftar ?? []
}

// ---------------------------------------------------------------------------
// Treaty Contract Out tiket 08 — klausul: satu tabel, 25 jenis.
// ---------------------------------------------------------------------------

/** Aturan satu jenis dari server — satu tabel kebenaran untuk form. */
export interface AturanKlausul {
  jenis: string
  anak: boolean
  subjenis: string
  medan: string[]
  wajib: string[]
  turunan: string[] | null
  ditahan: string
  /** Tiket 11: form menuntut kurs berlaku (`NewTreatyArr*`). */
  berkurs: boolean
  /** Tiket 11: '' | 'RpKeUsd' (`HitungRpUsd_depan`) | 'DuaArah' (`CalculateTSIExcludeTreaty`). */
  konversi: string
  sumber: string
  /** Sumber pilihan ReinsTypeID: '' = daftar induk tiket 02; `PILIHAN_REINS_ANAK_TREATY_LIMIT`. */
  pilihanReins?: string
  /** Satu baris per tahun (Minimum LOL, Max Coins Panel, Minimum LOL MB) — `Add` hilang begitu ada baris. */
  satuBaris?: boolean
}

/** Satu baris grid jenis (`BrowseTreatyDesc_RD`) beserta aturannya. */
export interface JenisKlausul {
  id: string
  descName: string
  isXol: string
  statusAktif: string
  aturan: AturanKlausul[]
  catatan: string
}

/** Satu klausul. ⛔ Nilai medan TEKS — uang dan persen tidak pernah `Number`. */
export interface Klausul {
  id: string
  treatyYear: string
  treatyYearId: string
  treatyGroupId: string
  treatyDescId: string
  treatyDescName: string
  reinsTypeId: string
  reinsTypeName: string
  parentReinsTypeId: string
  subjenis: string
  medan: Record<string, string>
  kurs: string
  userId: string
  tglUpdate: string
}

/** Grid satu jenis (induk atau anak satu induk). */
export interface DaftarKlausul {
  daftar: Klausul[]
  total: number
  totalPct: string
  peringatan: string
}

/** Badan simpan satu klausul. */
export interface KlausulMasuk {
  id: string
  descId: string
  anak: boolean
  subjenis: string
  parentReinsTypeId: string
  medan: Record<string, string>
}

/** Jawaban simpan; `peringatan` = "Please make sure spreading is 100%" bila berlaku. */
export interface HasilKlausul {
  klausul: Klausul
  totalPct: string
  peringatan: string
}

/** Satu pilihan occupation / clause. */
export interface PilihanKlausul {
  id: string
  nama: string
}

/** Jenis klausul per `isXol`; layar memakai '0' saja (grid `Treaty Desc` — For XOL dibuang 30-09-2026). */
export async function ambilJenisKlausul(isXol: string): Promise<JenisKlausul[]> {
  const j = await minta<{ daftar: JenisKlausul[] | null }>('/api/treaty-contract-out/jenis-klausul', {
    kueri: { isXol },
  })
  return j.daftar ?? []
}

/** Grid satu jenis — `induk` '00' untuk baris induk, jenis reasuransi induk untuk anak. */
export async function ambilKlausul(tahunID: string, descId: string, induk: string): Promise<DaftarKlausul> {
  const j = await minta<DaftarKlausul>(`/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/klausul`, {
    kueri: { descId, induk },
  })
  return { ...j, daftar: j.daftar ?? [] }
}

/** `Save` per jenis — POST bila baru, PUT /{id} bila ubah. */
export async function simpanKlausul(tahunID: string, masuk: KlausulMasuk): Promise<HasilKlausul> {
  const jalur = `/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/klausul`
  if (masuk.id === '') {
    return minta(jalur, { metode: 'POST', badan: masuk })
  }
  return minta(`${jalur}/${encodeURIComponent(masuk.id)}`, { metode: 'PUT', badan: masuk })
}

/** Pemilih ExclutionTreaty — `occupation` (BrowseOccupationFIRE_RD) atau `clause` (BrowseFireClauseFacIn_RD). */
export async function cariPilihanKlausul(
  master: 'occupation' | 'clause' | 'occupation-limitmb',
  cari: string,
): Promise<PilihanKlausul[]> {
  const j = await minta<{ daftar: PilihanKlausul[] | null }>(
    `/api/treaty-contract-out/klausul-pilihan/${master}`,
    { kueri: { cari } },
  )
  return j.daftar ?? []
}

// ---------------------------------------------------------------------------
// Treaty Contract Out tiket 06 — security di bawah reinsurer.
// ---------------------------------------------------------------------------

/** Satu security. ⛔ `pctShare` TEKS — tidak pernah `Number`. */
export interface SecurityReinsurer {
  id: string
  thnTreaty: string
  reasId: string
  reasSecurity: string
  clientName: string
  pctShare: string
  topId: string
  tpTreaty: string
  userId: string
}

/** Grid security seorang reinsurer. */
export interface DaftarSecurity {
  daftar: SecurityReinsurer[]
  total: number
  reinsurer: ReinsurerTreaty
}

/** Badan simpan — `Security Name` (ID dari pemilih) dan `%Share`. */
export interface SecurityMasuk {
  id: string
  reasSecurity: string
  pctShare: string
}

function jalurSecurity(tahunID: string, kontrakID: string, reinsurerID: string): string {
  return `${jalurReinsurer(tahunID, kontrakID)}/${encodeURIComponent(reinsurerID)}/security`
}

/** Grid `SelectSecurityReinsurer` (THN_TREATY + REAS_ID). */
export async function ambilSecurity(tahunID: string, kontrakID: string, reinsurerID: string): Promise<DaftarSecurity> {
  const j = await minta<DaftarSecurity>(jalurSecurity(tahunID, kontrakID, reinsurerID))
  return { ...j, daftar: j.daftar ?? [] }
}

/** `Save` b20246 — POST bila baru, PUT /{id} bila ubah. */
export async function simpanSecurity(
  tahunID: string,
  kontrakID: string,
  reinsurerID: string,
  masuk: SecurityMasuk,
): Promise<SecurityReinsurer> {
  const j = jalurSecurity(tahunID, kontrakID, reinsurerID)
  if (masuk.id === '') {
    return minta(j, { metode: 'POST', badan: masuk })
  }
  return minta(`${j}/${encodeURIComponent(masuk.id)}`, { metode: 'PUT', badan: masuk })
}

/** `Delete` b17559 — satu baris menurut ID. */
export async function hapusSecurity(tahunID: string, kontrakID: string, reinsurerID: string, id: string): Promise<string> {
  const j = await minta<{ pesan: string }>(`${jalurSecurity(tahunID, kontrakID, reinsurerID)}/${encodeURIComponent(id)}`, {
    metode: 'DELETE',
  })
  return j.pesan
}

// ---------------------------------------------------------------------------
// Treaty Contract Out tiket 11 — kurs USD → IDR.
// ---------------------------------------------------------------------------

/** Kurs berlaku pada tanggal mulai tahun treaty. ⛔ TEKS — tidak pernah `Number`. */
export interface KursTahun {
  kurs: string
  tanggal: string
  mulai: string
  akhir: string
  currency: string
  idCurrency: string
  quarter: string
  treatyYear: string
  /** Cacah baris berlaku lain yang identik — master memuat baris kembar. Opsional: backend lama tidak mengirimnya. */
  barisMasterKembar?: number
}

/** Hasil konversi; kedua nilai tetap terpisah (AC 49). */
export interface KonversiKurs {
  rp: string
  usd: string
  kurs: string
}

/** `testingKurs` — 422 "Tidak ada Nilai Kurs di Tahun : <tahun>" bila kosong. */
export async function ambilKursTahun(tahunID: string): Promise<KursTahun> {
  return minta<KursTahun>(`/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/kurs`)
}

/** Konversi di SERVER (desimal persis): dari 'Rp' → Usd pada skala '8' / '4'; dari 'Usd' → Rp. */
export async function konversiKurs(tahunID: string, dari: 'Rp' | 'Usd', nilai: string, skala: '4' | '8'): Promise<KonversiKurs> {
  return minta<KonversiKurs>(`/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/kurs/konversi`, {
    kueri: { dari, nilai, skala },
  })
}

// ---------------------------------------------------------------------------
// Treaty Contract Out tiket 10 — kaskade hapus + popup konfirmasi.
// ---------------------------------------------------------------------------

/** Isi popup: yang ikut terhapus. */
export interface DampakHapusTCO {
  kontrak: number
  reinsurer: number
  security: number
  business: number
  /** Klausul yang TETAP (AC 44) — masih dikirim server, tidak lagi ditampilkan popup [keputusan work owner 01-10-2026]. */
  klausulTetap: number
  /** Kontrak lain yang memakai kombinasi yang sama — anak kombinasinya IKUT terhapus (OQ-TCO-21). */
  bersama: number
}

function jalurKontrakHapus(tahunID: string, kontrakID: string): string {
  return `/api/treaty-contract-out/tahun/${encodeURIComponent(tahunID)}/kontrak/${encodeURIComponent(kontrakID)}`
}

/** Pratinjau dampak hapus kontrak. */
export async function ambilDampakHapusKontrak(tahunID: string, kontrakID: string): Promise<DampakHapusTCO> {
  return minta<DampakHapusTCO>(`${jalurKontrakHapus(tahunID, kontrakID)}/dampak-hapus`)
}

/** Ya di popup — jumlah yang dilihat dikirim; server menolak (409) bila sudah lain. */
export async function hapusKontrak(tahunID: string, kontrakID: string, d: DampakHapusTCO): Promise<string> {
  const j = await minta<{ pesan: string }>(jalurKontrakHapus(tahunID, kontrakID), {
    metode: 'DELETE',
    kueri: { reinsurer: String(d.reinsurer), security: String(d.security), business: String(d.business), bersama: String(d.bersama) },
  })
  return j.pesan
}

/** Pratinjau dampak hapus reinsurer (security yang ikut). */
export async function ambilDampakHapusReinsurer(tahunID: string, kontrakID: string, reinsurerID: string): Promise<DampakHapusTCO> {
  return minta<DampakHapusTCO>(`${jalurKontrakHapus(tahunID, kontrakID)}/reinsurer/${encodeURIComponent(reinsurerID)}/dampak-hapus`)
}

/** Ya di popup hapus reinsurer. */
export async function hapusReinsurer(tahunID: string, kontrakID: string, reinsurerID: string, d: DampakHapusTCO): Promise<string> {
  const j = await minta<{ pesan: string }>(`${jalurKontrakHapus(tahunID, kontrakID)}/reinsurer/${encodeURIComponent(reinsurerID)}`, {
    metode: 'DELETE',
    kueri: { security: String(d.security) },
  })
  return j.pesan
}
