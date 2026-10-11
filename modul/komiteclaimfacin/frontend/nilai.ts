// Nilai layar Komite Claim Fac In - fungsi murni (diuji tanpa DOM). Disalin dari Komite Claim Non Prop
// (`modul/komiteclaimnonprop/frontend/nilai.ts`), tidak diimpor (batas modul): angka SELALU 4 desimal
// (`pxNumber DecimalPlaces=4`; work owner 11-10-2026 "TAMPILAN 4 DIGIT BELAKANG KOMA") pemisah Indonesia, tanggal
// `dd-mm-yyyy` (+ ` HH:mm`).

import { formatNumber } from '../../../inti/frontend/lib/format'
import { keInputTanggal } from '../../../inti/frontend/lib/tanggalInput'
import type { IsianLayar, JenisNilai, Keputusan } from './api'

/** Desimal tampilan angka. */
export const DESIMAL_TAMPIL = 4

/** Hasil `formatNumber` yang berupa angka (titik ribuan, koma desimal opsional). */
const ANGKA_TAMPIL = /^-?\d{1,3}(\.\d{3})*(,\d+)?$/

/** Angka 4 desimal: `formatNumber` (dibulatkan pada digit) lalu nol ekor dilengkapi; bukan angka = apa adanya. */
function angka4(s: string): string {
  const t = formatNumber(s, DESIMAL_TAMPIL)
  if (!ANGKA_TAMPIL.test(t)) return t
  const [bulat, pecahan = ''] = t.split(',')
  return `${bulat},${pecahan.padEnd(DESIMAL_TAMPIL, '0')}`
}

/** pxCheckbox hanya-baca: "true" / "1" dicentang. */
export function tercentang(v: string): boolean {
  const s = (v ?? '').trim()
  return s === 'true' || s === '1'
}

/** Nilai halaman -> tampilan menurut jenisnya. Tidak terbaca = apa adanya. */
export function tampil(v: string, jenis: JenisNilai | ''): string {
  const s = (v ?? '').trim()
  if (s === '') return ''
  switch (jenis) {
    case 'angka':
      return angka4(s)
    case 'tanggal':
    case 'tanggalJam': {
      const iso = keInputTanggal(s)
      if (iso === '') return s
      const dmy = `${iso.slice(8, 10)}-${iso.slice(5, 7)}-${iso.slice(0, 4)}`
      if (jenis === 'tanggal') return dmy
      const jam = /[T ](\d{2}):(\d{2})/.exec(s)
      return jam ? `${dmy} ${jam[1]}:${jam[2]}` : dmy
    }
    case 'centang': // pxCheckbox hanya-baca (layar merender kotak centang nonaktif; teks ini untuk sel tanpa kotak)
      return tercentang(s) ? '✓' : ''
    default:
      return s
  }
}

/** Kode `.AcceptStatus` (`models.KeputusanSetuju` / `KeputusanTolak`: 1 Approve, 2 Reject). */
export const KEPUTUSAN = { setuju: '1', tolak: '2' } as const

/** Pesan wajib isi bawaan Pega (`models.PesanKosong`). */
export const PESAN_KOSONG = 'Value cannot be blank'

/** Validasi klien (LS45 `REQ=true`) - sama dengan `models.PeriksaIsian` di server; server tetap menegakkan. */
export function periksa(k: Keputusan): Partial<Record<keyof Keputusan, string>> {
  const out: Partial<Record<keyof Keputusan, string>> = {}
  if (k.acceptStatus !== KEPUTUSAN.setuju && k.acceptStatus !== KEPUTUSAN.tolak) out.acceptStatus = PESAN_KOSONG
  if (k.comment.trim() === '') out.comment = PESAN_KOSONG
  return out
}

/** Dua Propose dapat diubah: tampil (TT2) dan terbuka (tingkat 1). */
export function usulTerbuka(isian: Pick<IsianLayar, 'tampilUsul' | 'terbuka'>): boolean {
  return isian.tampilUsul && isian.terbuka
}

/** Isian yang dikirim: dua Propose yang tersembunyi / nonaktif tidak membawa nilai baru. */
export function isianKirim(k: Keputusan, isian: Pick<IsianLayar, 'tampilUsul' | 'terbuka'>): Keputusan {
  const usul = usulTerbuka(isian)
  return {
    acceptStatus: k.acceptStatus,
    comment: k.comment,
    usulTutup: usul ? k.usulTutup : false,
    usulCadang: usul ? k.usulCadang : false,
  }
}
