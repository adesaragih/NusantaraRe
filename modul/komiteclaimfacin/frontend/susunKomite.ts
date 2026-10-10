// Penyusun tata letak layar komite - fungsi murni (pola Komite Claim Prop / Non Prop, keputusan work owner 09-10-2026:
// kartu berjudul, pasangan label-nilai, tangga dan kartu keputusan DI BAWAH semua rincian). Beda dengan pola Non Prop:
// server Fac In sudah memberi judul kartu per bagian (`models.SusunLayar`), jadi bagian TIDAK dikelompokkan ulang -
// satu bagian satu kartu, urutan server. Yang diatur di sini hanya letak: teks komite (`teksKomite`) kartu terakhir
// rincian, tangga ("List of Committee") di bawahnya bersama kartu keputusan.

import type { Anggota, Bagian, Grid, Layar } from './api'

/** Kunci bagian teks komite (LS38-LS39) - kartu rincian TERAKHIR, tepat sebelum tangga. */
export const KUNCI_TEKS_KOMITE = 'teksKomite'
/** Kunci bagian tangga ("List of Committee", LS41 `.TransferType == 2`) - dirender sebagai langkah, di bawah. */
export const KUNCI_TANGGA = 'tangga'

/** Bagian berisi: punya medan atau grid. Bagian kosong (mis. "Claim Details" tanpa adjustment) dilewati. */
export function bagianBerisi(b: Bagian): boolean {
  return (b.medan?.length ?? 0) > 0 || (b.grid?.length ?? 0) > 0
}

/** Bagian server -> kartu rincian (urutan server, teks komite terakhir) dan tangga (bila ada). */
export function susunBagian(bagian: readonly Bagian[]): { kartu: Bagian[]; tangga?: Bagian } {
  const kartu = bagian.filter((b) => b.kunci !== KUNCI_TANGGA && b.kunci !== KUNCI_TEKS_KOMITE && bagianBerisi(b))
  kartu.push(...bagian.filter((b) => b.kunci === KUNCI_TEKS_KOMITE && bagianBerisi(b)))
  return { kartu, tangga: bagian.find((b) => b.kunci === KUNCI_TANGGA) }
}

/** Baris grid (server boleh mengirim `null` untuk daftar kosong). */
export function barisGrid(g: Grid): Record<string, string>[] {
  return g.baris ?? []
}

/** Rincian expand pane baris ke-`i` (0-based); `null` = baris itu tidak dapat dibuka. */
export function rincianBaris(g: Grid, i: number): Bagian[] | null {
  const r = g.rincian?.[i]
  return r && r.length > 0 ? r : null
}

/** Grid punya setidaknya satu baris yang dapat dibuka (kolom panah ditampilkan). */
export function adaRincian(g: Grid): boolean {
  return barisGrid(g).some((_, i) => rincianBaris(g, i) !== null)
}

/** Kunci pop-up sel `tautan` yang isinya ada di `Layar.modal`; selainnya `null` (sel tanpa tautan). */
export function kunciModal(nilai: string | undefined, modal: Layar['modal']): string | null {
  const k = (nilai ?? '').trim()
  return k !== '' && modal?.[k] !== undefined ? k : null
}

/** Keadaan satu tingkat tangga: disetujui, ditolak, sedang berjalan (menunggu PERTAMA), menunggu. */
export type KeadaanLangkah = 'setuju' | 'tolak' | 'berjalan' | 'menunggu'

export interface Langkah {
  no: string
  jabatan: string
  /** Label "Status" VERBATIM server (Waiting / Approved / Reject). */
  status: string
  keadaan: KeadaanLangkah
  tanggal: string
  komentar: string
}

/**
 * Grid tangga server -> langkah. Baris grid = tangga tersimpan (`kasus.tangga`, kode `KOMITE_APPROVAL` 0 / 1 / 2) lalu
 * perluasan tingkat 1 yang belum disimpan (KCF-02, selalu menunggu); urutan sama (`models.barisTangga`).
 */
export function langkahTangga(tangga: Bagian | undefined, tersimpan: readonly Anggota[] | null): Langkah[] {
  const g = tangga?.grid?.[0]
  if (!g) return []
  let berjalan = false
  return barisGrid(g).map((b, i) => {
    const kode = tersimpan?.[i]?.keputusan ?? '0'
    let keadaan: KeadaanLangkah = 'menunggu'
    if (kode === '1') keadaan = 'setuju'
    else if (kode === '2') keadaan = 'tolak'
    else if (!berjalan) {
      keadaan = 'berjalan'
      berjalan = true
    }
    return {
      no: b['No'] ?? String(i + 1),
      jabatan: b['jabatan'] ?? '',
      status: b['keputusan'] ?? '',
      keadaan,
      tanggal: b['tanggal'] ?? '',
      komentar: b['komentar'] ?? '',
    }
  })
}

/** Lencana keadaan kepala layar. */
export type KeadaanKasus = { jenis: 'anda' } | { jenis: 'tunggu'; jabatan: string } | { jenis: 'selesai' }

/**
 * Giliran anda (pemegang tingkat berjalan), menunggu jabatan tingkat berjalan (`KomiteList(KomiteCount)`, cadangan:
 * baris menunggu pertama), atau selesai (kasus tertutup / tanpa baris menunggu).
 */
export function keadaanKasus(l: Layar): KeadaanKasus {
  if (l.bolehKerja) return { jenis: 'anda' }
  const t = l.kasus.tangga ?? []
  if (l.kasus.statusWork !== '') return { jenis: 'selesai' }
  const a = t.find((x) => x.urut === l.kasus.komiteCount && x.keputusan === '0') ?? t.find((x) => x.keputusan === '0')
  return a ? { jenis: 'tunggu', jabatan: a.jabatan } : { jenis: 'selesai' }
}
