// Pilihan peserta layar Register yang BERTAHAN antar pencarian - keputusan work owner 10-10-2026: "saat user ceklis
// 2 orang dari PL A, lalu 4 orang dari PL B, claim membuat 2 case (CLMLF) berdasarkan PL number".
//
// Murni (tanpa React): kunci pilihan, alih centang, pengelompokan tampilan per PL + mata uang, dan badan permintaan.
// ⛔ Pengelompokan di sini hanya TAMPILAN dan perkiraan jumlah case; server membentuk ulang case dari baris sumber
// (services.KelompokkanPeserta) - mata uang dan polis sumber yang menang.

import type { CalonPeserta, KelompokDaftar } from './api'

/** Satu peserta terpilih - salinan baris hasil pencarian saat dicentang. */
export type PesertaTerpilih = Pick<
  CalonPeserta,
  'nomorPremiList' | 'nomorPolis' | 'nomorSertifikat' | 'namaTertanggung' | 'mataUang'
> & { idPega?: string; tanggalLahir?: string }

/**
 * Kunci pilihan = PL + polis + sertifikat.
 *
 * ⛔ Polis ikut, bukan PL + sertifikat saja: satu premium list dapat memuat beberapa polis, dan nomor sertifikat
 * hanya unik di dalam polisnya - dua peserta bersertifikat "006" dari polis berbeda akan tercentang berbarengan bila
 * kuncinya tanpa polis. PL ikut sejak pilihan bertahan antar pencarian PL lain.
 */
export function kunciPilihan(p: PesertaTerpilih): string {
  return `${p.nomorPremiList}|${p.nomorPolis}|${p.nomorSertifikat}`
}

/** Centang / lepas satu peserta; daftar lama tidak diubah. */
export function alihPilihan(daftar: readonly PesertaTerpilih[], p: PesertaTerpilih): PesertaTerpilih[] {
  const k = kunciPilihan(p)
  if (daftar.some((x) => kunciPilihan(x) === k)) return daftar.filter((x) => kunciPilihan(x) !== k)
  const { nomorPremiList, nomorPolis, nomorSertifikat, namaTertanggung, mataUang, idPega, tanggalLahir } = p
  return [...daftar, { nomorPremiList, nomorPolis, nomorSertifikat, namaTertanggung, mataUang, idPega, tanggalLahir }]
}

/** Apakah baris hasil pencarian sudah terpilih (centang tetap tampil saat PL dicari lagi). */
export function sudahDipilih(daftar: readonly PesertaTerpilih[], p: PesertaTerpilih): boolean {
  const k = kunciPilihan(p)
  return daftar.some((x) => kunciPilihan(x) === k)
}

/** Satu kelompok tampilan = satu case yang akan dibuat. */
export interface KelompokPilihan {
  nomorPremiList: string
  mataUang: string
  peserta: PesertaTerpilih[]
}

/** Kelompok tampilan per PL number + mata uang, urut kemunculan pertama. */
export function kelompokkanPilihan(daftar: readonly PesertaTerpilih[]): KelompokPilihan[] {
  const hasil: KelompokPilihan[] = []
  for (const p of daftar) {
    let k = hasil.find((x) => x.nomorPremiList === p.nomorPremiList && x.mataUang === p.mataUang)
    if (k === undefined) {
      k = { nomorPremiList: p.nomorPremiList, mataUang: p.mataUang, peserta: [] }
      hasil.push(k)
    }
    k.peserta.push(p)
  }
  return hasil
}

/** "{n} claim case(s) will be created". */
export function kalimatJumlahCase(n: number, pola: string): string {
  return pola.split('{n}').join(String(n))
}

/**
 * Badan `kelompok` untuk POST /api/klaim-life: satu kelompok per PL + polis (nomor polis dicocokkan server dengan
 * POLICY_NO baris sumber). Mata uang ikut sebagai keterangan; server mengabaikannya.
 */
export function kelompokPermintaan(daftar: readonly PesertaTerpilih[]): KelompokDaftar[] {
  const hasil: KelompokDaftar[] = []
  for (const p of daftar) {
    let k = hasil.find((x) => x.nomorPremiList === p.nomorPremiList && x.nomorPolis === p.nomorPolis)
    if (k === undefined) {
      k = { nomorPremiList: p.nomorPremiList, nomorPolis: p.nomorPolis, mataUang: p.mataUang, sertifikat: [] }
      hasil.push(k)
    }
    k.sertifikat.push(p.nomorSertifikat)
  }
  return hasil
}

/** Hasil baca polis untuk panel Data Polis. */
export interface HasilPolisPeserta<P> {
  polis: P | null
  /** true = header tidak ditemukan untuk peserta itu (atau IDPEGA-nya kosong). */
  tidakDitemukan: boolean
  /** IDPEGA yang dicari ("" bila peserta tanpa IDPEGA). */
  idPega: string
}

/**
 * Membaca header polis peserta dari IDPEGA-nya (keputusan work owner 10-10-2026: "ambil ke tabel T_PREMIUM_LIST
 * berdasarkan idpega peserta nya saja"). ⛔ Tanpa fallback NO_POLIS / PL / POLICY_NO - sama dengan backend
 * `services/kuncipolis.go`. Tanpa peserta (idPega undefined) = tidak ada yang dibaca, tanpa catatan.
 */
export async function bacaPolisPeserta<P>(
  ambil: (idPega: string) => Promise<P>,
  idPega: string | undefined,
): Promise<HasilPolisPeserta<P>> {
  if (idPega === undefined) return { polis: null, tidakDitemukan: false, idPega: '' }
  const id = idPega.trim()
  if (id === '') return { polis: null, tidakDitemukan: true, idPega: '' }
  try {
    return { polis: await ambil(id), tidakDitemukan: false, idPega: id }
  } catch {
    return { polis: null, tidakDitemukan: true, idPega: id }
  }
}
