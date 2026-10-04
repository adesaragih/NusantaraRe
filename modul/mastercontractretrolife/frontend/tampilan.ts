// Pembantu tampilan MURNI modul Master Contract Retro Life - dipakai kelima grid dan kelima form.
//
// ⛔ Angka tidak pernah menjadi `Number`: sel angka lewat `formatNumber` (digit, bukan float).

import { DESIMAL_TAK_DIBATASI, formatDate, formatNumber } from '../../../inti/frontend/lib/format'
import { keInputTanggal } from '../../../inti/frontend/lib/tanggalInput'
import { pelakuStub } from '../../../inti/frontend/store/sesi'

/**
 * Baris per halaman - `pyRDLPageSize` 10 di kelima grid (`InputRetrocessionLife.xml` b9441,
 * `InputRetroLimitReinsurers.xml` b9230, `InputSecurityLifeReinsurers.xml` b9333,
 * `InputSecurityReinsurerLife.xml` b8823, `InputBusinessLifeReinsurers.xml` b8979), penomoran
 * `pyGridPaginator` (b9192, b8995, b8956, b8446, b8752).
 */
export const UKURAN_HALAMAN_MCRL = 10

/** Baris halaman ke-`halaman` (mulai 1). */
export function potongHalaman<T>(daftar: readonly T[], halaman: number, ukuran = UKURAN_HALAMAN_MCRL): T[] {
  const mulai = (Math.max(1, halaman) - 1) * ukuran
  return daftar.slice(mulai, mulai + ukuran)
}

/** Halaman yang sah sesudah daftar berubah (mis. baris terakhir halaman terhapus). */
export function jepitHalaman(halaman: number, total: number, ukuran = UKURAN_HALAMAN_MCRL): number {
  const akhir = Math.max(1, Math.ceil(total / ukuran))
  return Math.min(Math.max(1, halaman), akhir)
}

/** Sel kosong ditandai, bukan dibiarkan kosong (ADR-U-0027). */
export function sel(v: string): string {
  return v.trim() === '' ? '—' : v
}

/** Sel tanggal `YYYY-MM-DD` → `DD-MM-YYYY`. */
export function selTanggal(v: string): string {
  return sel(formatDate(v))
}

/** Sel/medan waktu `YYYY-MM-DD HH:MM:SS` → `DD-MM-YYYY HH:MM:SS`; bentuk lain apa adanya. */
export function selWaktu(v: string): string {
  const m = /^(\d{4})-(\d{2})-(\d{2})[ T](\d{2}:\d{2}:\d{2})/.exec(v.trim())
  return m ? `${m[3]}-${m[2]}-${m[1]} ${m[4]}` : sel(v)
}

/** Sel angka (share, komisi, batas) - digit apa adanya, pemisah ribuan, nol ekor dipangkas. */
export function selAngka(v: string): string {
  return sel(formatNumber(v, DESIMAL_TAK_DIBATASI))
}

/**
 * Tanggal isian → `YYYY-MM-DD`. Teks yang bukan tanggal dikirim APA ADANYA supaya server yang
 * menolaknya dengan kalimat yang menyebut medannya - layar tidak menebak.
 */
export function keTanggalKabel(v: string): string {
  const t = v.trim()
  if (t === '') return ''
  const iso = keInputTanggal(t)
  return iso === '' ? t : iso
}

/** `@CurrentDateTime()` Pega sebagai `YYYY-MM-DD HH:MM:SS` waktu setempat. */
export function waktuKini(d: Date): string {
  const p = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

/**
 * `OperatorID.pyUserName` - medan `Inputor` form baru (`DATASHOW2 = 0`) dan nilai `USERID` yang
 * disalin activity `Set*`/`New*`. Yang dikenal layar adalah akun pelaku; server menulis `USERID`
 * dari pelaku yang sama, bukan dari isian.
 */
export function operatorKini(): string {
  return pelakuStub()?.akunID ?? ''
}
