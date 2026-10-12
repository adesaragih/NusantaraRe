// Nilai tampilan kasus Komite Claim Life - fungsi murni (dipindah dari halaman Inbox Komite, yang kini menjadi tabel
// komite di bawah inbox Claim Life - perintah work owner 09-10-2026: menu komite dihapus).

import { DESIMAL_TAK_DIBATASI, formatNumber } from '../../../inti/frontend/lib/format'
import type { BarisInboxKomite } from './api'

/** Sel kosong ditandai (ADR-U-0027). */
export function selKomite(nilai: string): string {
  return nilai.trim() === '' ? '—' : nilai
}

/** Tingkat berjalan dari seluruh tingkat — `KomiteCount`/`KomiteLoop`. */
export function tingkatKomite(b: BarisInboxKomite): string {
  return `${String(b.tingkatBerjalan)} / ${String(b.komiteLoop)}`
}

// ——— Layar kasus komite dirapikan (work owner 12-10-2026, "perbaiki tampilannya") ———
// ⛔ TAMPILAN SAJA: kata status dari server (`Setuju`/`Tolak`/`Menunggu`/`Dilewati (eskalasi)`, juga dipakai jejak) TIDAK
// diubah di server; layar menerjemahkannya ke Inggris di sini. Kata asing ditampilkan apa adanya.

/** "YYYY-MM-DD HH:MM[:SS]" (atau RFC3339) → dd/mm/yyyy HH:mm; kosong / tak terbaca = "—" / apa adanya. */
export function waktuKomite(nilai: string): string {
  const t = nilai.trim()
  if (t === '') return '—'
  const m = /^(\d{4})-(\d{2})-(\d{2})[T ](\d{2}):(\d{2})/.exec(t)
  return m === null ? t : `${m[3]}/${m[2]}/${m[1]} ${m[4]}:${m[5]}`
}

/** Claim Amount "2.342.342 IDR" — teks desimal berformat ribuan, TANPA float. */
export function nilaiKlaimKomite(nilai: string, mataUang: string): string {
  if (nilai.trim() === '') return '—'
  return `${formatNumber(nilai, DESIMAL_TAK_DIBATASI)} ${mataUang.trim()}`.trim()
}

const KATA_STATUS_KOMITE: Readonly<Record<string, string>> = {
  Setuju: 'Approved',
  Tolak: 'Rejected',
  Menunggu: 'Waiting',
  'Dilewati (eskalasi)': 'Skipped (Escalation)',
}

/** Kata status tangga / keputusan dari server → kata layar berbahasa Inggris. */
export function kataStatusKomite(kata: string): string {
  const t = kata.trim()
  if (t === '') return '—'
  const kode = /^Kode (.+)$/.exec(t)
  if (kode !== null) return `Code ${kode[1] ?? ''}`
  return KATA_STATUS_KOMITE[t] ?? t
}

export type NadaStatusKomite = 'menunggu' | 'setuju' | 'tolak' | 'lain'

/** Nada chip status (warna) dari kata server. */
export function nadaStatusKomite(kata: string): NadaStatusKomite {
  switch (kata.trim()) {
    case 'Menunggu':
      return 'menunggu'
    case 'Setuju':
      return 'setuju'
    case 'Tolak':
      return 'tolak'
    default:
      return 'lain'
  }
}

const KATA_EFEK_KOMITE: Readonly<Record<string, string>> = {
  antre: 'Queued',
  jalan: 'Running',
  selesai: 'Done',
  'gagal-permanen': 'Failed Permanently',
  tertunda: 'Pending',
  tuntas: 'Complete',
  'perlu intervensi': 'Needs Intervention',
}

/** Keadaan efek keluar dari server → kata layar berbahasa Inggris. */
export function kataEfekKomite(kata: string): string {
  const t = kata.trim()
  return t === '' ? '—' : (KATA_EFEK_KOMITE[t] ?? t)
}

// ——— Form of Internal Claim Requisition (work owner 12-10-2026, "perbaiki tampilannya lagi") ———

/** Tanggal saja (pyFormatType date ShowTransfer): "YYYY-MM-DD[ ...]" → dd/mm/yyyy; kosong = "—", tak terbaca apa adanya. */
export function tanggalKomite(nilai: string | null | undefined): string {
  const t = (nilai ?? '').trim()
  if (t === '') return '—'
  const m = /^(\d{4})-(\d{2})-(\d{2})/.exec(t)
  return m === null ? t : `${m[3]}/${m[2]}/${m[1]}`
}

/**
 * Uang dua desimal seperti layar Pega ("54.000.000,00"): formatNumber atas TEKS (dibulatkan pada digit, tanpa float),
 * lalu desimalnya dilengkapi sampai dua. Kosong / tanpa nilai = "—".
 */
export function uangKomite(nilai: string | null | undefined): string {
  const t = formatNumber(nilai ?? '', 2)
  if (t === '') return '—'
  if (!/^-?[\d.]+(,\d*)?$/.test(t)) return t
  const [bulat, pecahan = ''] = t.split(',')
  return `${bulat ?? ''},${pecahan.padEnd(2, '0')}`
}

/** Teks bebas; kosong = "—". */
export function teksKomite(nilai: string | null | undefined): string {
  const t = (nilai ?? '').trim()
  return t === '' ? '—' : t
}
