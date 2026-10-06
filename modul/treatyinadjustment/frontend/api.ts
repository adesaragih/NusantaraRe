// modul/treatyinadjustment/frontend/api.ts - panggilan backend modul Treaty In
// Adjustment (`backend/handlers/rute_treaty_in_adjustment.go`).
//
// ⛔ Seluruh rute BACA. Jalur simpan belum dibangun; tombol tulis layar dimatikan.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_TREATYINADJUSTMENT = '/api/treaty-in-adjustment'

/** Kepala kontrak - lapisan BEKU yang seluruh versinya bagi (ADR-0040). */
export interface Kontrak {
  id: number
  nomorKontrakWarisan: string
  sifatProporsi: string
  tanggalMulai: string
  tanggalBerakhir: string
}

/**
 * Satu baris rantai versi - SAMA dengan `models.Versi`.
 *
 * `nomorUrutVersi` null = baris warisan yang belum dinomori ulang (tiket 10),
 * BUKAN nol. `idVersiDasar` null = versi PERTAMA kontrak itu, dan itu keadaan
 * yang benar, bukan data yang hilang.
 */
export interface Versi {
  id: number
  idKontrak: number
  nomorUrutVersi: number | null
  keadaanSiklusHidup: string
  jenisAddendum: string
  sifatMaterialAddendum: string
  tanggalBerlakuAddendum: string
  idVersiDasar: number | null
  namaKontrak: string
}

/** `GET /api/treaty-in-adjustment/kontrak` - kepala seluruh kontrak. */
export async function ambilKontrak(): Promise<Kontrak[]> {
  return minta<Kontrak[]>(`${PREFIX_TREATYINADJUSTMENT}/kontrak`)
}

/** `GET /api/treaty-in-adjustment/kontrak/{id}/versi` - rantai versi satu kontrak. */
export async function ambilRantaiVersi(idKontrak: number): Promise<Versi[]> {
  return minta<Versi[]>(`${PREFIX_TREATYINADJUSTMENT}/kontrak/${idKontrak}/versi`)
}

/**
 * Panel Attachment — dari `POOLDATA.M_ATTACHMENTTREATY_2`, tabel WARISAN
 * 43 baris. Nol tabel baru, nol migrasi.
 */
export interface BarisLampiranWarisan {
  id: string
  kodeKategori: string
  namaKategori: string
  namaBerkas: string
  jenisMime: string
  idSimpanan: string
  diunggah: string
  pengunggah: string
}

/**
 * Satu baris panel Attachment: Category + Count.
 *
 * ⛔ `dipastikan` menyatakan apakah pasangan kode↔nama ini TERBUKTI. Empat
 * kode (`00003` `00004` `00008` `00009`) punya nol baris dan namanya nihil
 * di korpus kedua modul; layar menandainya alih-alih menebak.
 */
export interface BarisKategoriLampiran {
  kode: string
  nama: string
  cacah: number
  dipastikan: boolean
}

export interface LampiranKontrak {
  kategori: BarisKategoriLampiran[]
  berkas: BarisLampiranWarisan[]
}

/**
 * `GET /api/treaty-in-adjustment/kontrak-warisan/{id}/lampiran`.
 *
 * ⛔ Pengenalnya TEKS — `M_ATTACHMENTTREATY_2.TREATYID` adalah
 * `VARCHAR2(100)` milik sistem lama, bukan `idKontrak` model baru.
 */
export async function ambilLampiran(masterID: string): Promise<LampiranKontrak> {
  return minta<LampiranKontrak>(
    `${PREFIX_TREATYINADJUSTMENT}/kontrak-warisan/${encodeURIComponent(masterID)}/lampiran`,
  )
}

/** Satu baris panel History — dari `T_VIEW_COMMENT`, BACA SAJA. */
export interface BarisRiwayatWarisan {
  tanggal: string
  operator: string
  disetujui: string
  catatan: string
}

/** `GET /api/treaty-in-adjustment/kontrak-warisan/{id}/riwayat`. */
export async function ambilRiwayat(masterID: string): Promise<BarisRiwayatWarisan[]> {
  return minta<BarisRiwayatWarisan[]>(
    `${PREFIX_TREATYINADJUSTMENT}/kontrak-warisan/${encodeURIComponent(masterID)}/riwayat`,
  )
}

// ===========================================================================
// LAYAR ADJUSTMENT — `TREATY_IN_EDM` + `M_TREATY_IN_EDM`, BACA SAJA.
// ===========================================================================

/**
 * Satu baris grid daftar — SAMA dengan `models.BarisPenyesuaian`.
 *
 * ⛔ Nilainya APA ADANYA dari kolom tabel warisan. `jenisPenyesuaian` dan
 * `jenisMaterial` adalah KODE; teks pilihannya tidak ada di ekspor.
 */
export interface BarisPenyesuaian {
  id: string
  idAsal: string
  jenisPenyesuaian: string
  jenisMaterial: string
  namaKontrak: string
  sifatProporsi: string
  asalBisnis: string
  cedant: string
  tanggalMulai: string
  tanggalBerakhir: string
  posisi: string
  statusAkseptasi: string
}

/**
 * Satu halaman di dalam dokumen — `TreatyIn` (New) atau `TreatyIn.OLDDATA`
 * (Old).
 *
 * ⛔ `medan` hanya memuat kunci yang ADA di dokumen. `kunci in medan` yang
 * salah berarti "tidak ada di sistem lama" — BUKAN "kosong".
 */
export interface SisiPenyesuaian {
  medan: Record<string, string>
  larik: Record<string, Record<string, string>[]>
}

export interface Penyesuaian {
  id: string
  idAsal: string
  baru: SisiPenyesuaian
  lama: SisiPenyesuaian
}

/** `GET /api/treaty-in-adjustment/penyesuaian-warisan` — grid daftar. */
export async function ambilDaftarPenyesuaian(): Promise<BarisPenyesuaian[]> {
  return minta<BarisPenyesuaian[]>(`${PREFIX_TREATYINADJUSTMENT}/penyesuaian-warisan`)
}

/**
 * `GET /api/treaty-in-adjustment/penyesuaian-warisan/satu?id=…`.
 *
 * ⛔ Lewat parameter kueri: pengenalnya BERGARIS MIRING (`1000080/R02`).
 */
export async function ambilPenyesuaian(id: string): Promise<Penyesuaian> {
  return minta<Penyesuaian>(
    `${PREFIX_TREATYINADJUSTMENT}/penyesuaian-warisan/satu?id=${encodeURIComponent(id)}`,
  )
}
