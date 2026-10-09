// Nilai tampilan kasus Komite Claim Life - fungsi murni (dipindah dari halaman Inbox Komite, yang kini menjadi tabel
// komite di bawah inbox Claim Life - perintah work owner 09-10-2026: menu komite dihapus).

import type { BarisInboxKomite } from './api'

/** Sel kosong ditandai (ADR-U-0027). */
export function selKomite(nilai: string): string {
  return nilai.trim() === '' ? '—' : nilai
}

/** Tingkat berjalan dari seluruh tingkat — `KomiteCount`/`KomiteLoop`. */
export function tingkatKomite(b: BarisInboxKomite): string {
  return `${String(b.tingkatBerjalan)} / ${String(b.komiteLoop)}`
}
