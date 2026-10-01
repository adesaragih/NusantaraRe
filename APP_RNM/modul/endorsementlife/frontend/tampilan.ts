// Pembantu tampilan MURNI modul Endorsement Life.
//
// ⛔ Angka tidak pernah menjadi `Number`: sel angka lewat `formatNumber` (digit, bukan float).

import { DESIMAL_TAK_DIBATASI, formatDate, formatNumber } from '../../../inti/frontend/lib/format'
import { OPSI_EDM_TYPE, OPSI_TYPE_CEDING } from './labels'

/** Baris per halaman grid - SAMA dengan `services.UkuranHalaman` (`pyGridPaginator`). */
export const UKURAN_HALAMAN_EDM = 20

/** Sel kosong ditandai, bukan dibiarkan kosong (ADR-U-0027). */
export function sel(v: string | undefined): string {
  const t = (v ?? '').trim()
  return t === '' ? '—' : t
}

/** Sel tanggal `YYYY-MM-DD` → `DD-MM-YYYY`. */
export function selTanggal(v: string | undefined): string {
  return sel(formatDate(v ?? ''))
}

/** Sel waktu `YYYY-MM-DD HH:MM:SS` → `DD-MM-YYYY HH:MM:SS`; bentuk lain apa adanya. */
export function selWaktu(v: string | undefined): string {
  const m = /^(\d{4})-(\d{2})-(\d{2})[ T](\d{2}:\d{2}:\d{2})/.exec((v ?? '').trim())
  return m ? `${m[3]}-${m[2]}-${m[1]} ${m[4]}` : sel(v)
}

/** Sel angka - digit apa adanya, pemisah ribuan, nol ekor dipangkas; negatif tetap bertanda. */
export function selAngka(v: string | undefined): string {
  const t = (v ?? '').trim()
  return t === '' ? '—' : formatNumber(t, DESIMAL_TAK_DIBATASI)
}

/** Label `EDM Type` (`1` Perubahan Data, `3` Batal); nilai lain apa adanya. */
export function labelEdmType(v: string): string {
  return OPSI_EDM_TYPE[v] ?? sel(v)
}

/** Label `Reinsurance System` dari kode `TYPE_CEDING`; kode tak dikenal apa adanya. */
export function labelTypeCeding(v: string): string {
  return OPSI_TYPE_CEDING[v] ?? sel(v)
}

/** Pilihan hapus yang dikirim bersama `Save` - `models.PilihanHapus`. */
export interface PilihanHapusEDM {
  pilih: readonly string[]
  semua: boolean
  kecuali: readonly string[]
}

export const PILIHAN_KOSONG: PilihanHapusEDM = { pilih: [], semua: false, kecuali: [] }

/**
 * `DELETE ALL` b13607 → `SelectAllEdmLife_act` b259: sakelar `CARI1` - seluruh baris ikut nilai
 * barunya, centang satu per satu sebelumnya dibuang.
 */
export function alihSemua(p: PilihanHapusEDM): PilihanHapusEDM {
  return p.semua ? PILIHAN_KOSONG : { pilih: [], semua: true, kecuali: [] }
}

function alih(daftar: readonly string[], id: string): string[] {
  return daftar.includes(id) ? daftar.filter((x) => x !== id) : [...daftar, id]
}

/** Kotak centang `.EdmBatal` b15753 satu baris; sesudah `DELETE ALL` ia mengisi pengecualian. */
export function alihBaris(p: PilihanHapusEDM, id: string): PilihanHapusEDM {
  return p.semua ? { ...p, kecuali: alih(p.kecuali, id) } : { ...p, pilih: alih(p.pilih, id) }
}

export function tercentang(p: PilihanHapusEDM, id: string): boolean {
  return p.semua ? !p.kecuali.includes(id) : p.pilih.includes(id)
}

/**
 * Kotak centang hidup: kasus `Perubahan Data` (grid b11899; grid Batal b17500 tanpa kotak) yang
 * terbuka dan belum disimpan (`Save` mati sesudahnya, b37200), atas peserta `Old` saja (R30).
 */
export function bolehCentang(
  p: { edmStatus: string },
  k: { status: string; sudahSimpan: boolean; edmType: string },
): boolean {
  return k.status === '' && !k.sudahSimpan && k.edmType === '1' && p.edmStatus === 'Old'
}
