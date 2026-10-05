// Aturan layar Template Manager - fungsi murni, diuji tanpa DOM.

import type { Slot } from './api'

/** Slot dikelompokkan per menu (grup), urutan kemunculan dipertahankan. */
export function kelompokSlot(slot: readonly Slot[]): { grup: string; slot: Slot[] }[] {
  const hasil: { grup: string; slot: Slot[] }[] = []
  for (const s of slot) {
    const g = hasil.find((x) => x.grup === s.grup)
    if (g === undefined) hasil.push({ grup: s.grup, slot: [s] })
    else g.slot.push(s)
  }
  return hasil
}

/** Saringan teks (grup, nama, kode, nama berkas aktif) dan saringan menu; kosong = semua. */
export function saringSlot(slot: readonly Slot[], kata: string, grup: string): Slot[] {
  const k = kata.trim().toLowerCase()
  return slot.filter(
    (s) =>
      (grup === '' || s.grup === grup) &&
      (k === '' ||
        [s.grup, s.nama, s.kode, s.aktif?.namaBerkas ?? ''].some((t) => t.toLowerCase().includes(k))),
  )
}

/** Ukuran berkas untuk dibaca manusia: `832 B`, `13.5 KB`. */
export function ukuranBerkas(byte: number): string {
  if (byte < 1024) return `${byte} B`
  return `${(byte / 1024).toFixed(1)} KB`
}

/** Nomor versi berikutnya untuk tombol Simpan (riwayat terbaru dulu). */
export function versiBerikut(riwayat: readonly { versi: number }[]): number {
  return riwayat.reduce((m, v) => Math.max(m, v.versi), 0) + 1
}
