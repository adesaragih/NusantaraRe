// Disalin dari `modul/claimnonprop/frontend/pages/inbox.ts` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Aturan halaman awal (pola Claim Prop, keputusan work owner 07-10 dan 08-10-2026): dua tab, Process dan Resolve. Tab
// Process bawaan = worklist pembuat (Assignment1 Input Register + Assignment7 Input Estimasi, `Register_Flow` ke worklist
// pembuat). Switch Teknik menampilkan Assignment3 "Choose Surveyor" (workbasket `ReasKlaimTeknik`) dan hanya dapat
// dinyalakan anggota workbasket itu. Pembuatan klaim (Start -> Assignment1) hanya saat switch Teknik mati. Tabel komite
// di bawah inbox (modul `komiteclaimfacin` tanpa menu, prompt tahap 2 §2 butir 2; pola Claim Non Prop).

import type { HakPelaku, JenisDaftar } from '../api'

export type TabInbox = 'proses' | 'selesai'

export const TAB_INBOX: readonly TabInbox[] = ['proses', 'selesai']

/** Daftar yang diminta ke server untuk tab dan posisi switch Teknik. */
export function jenisDaftar(tab: TabInbox, teknik: boolean): JenisDaftar {
  if (tab === 'selesai') return 'selesai'
  return teknik ? 'workbasket' : 'saya'
}

/** Switch Teknik dapat dinyalakan hanya bila akun memegang workbasket ReasKlaimTeknik (Choose Surveyor). */
export function switchTeknikAktif(hak: HakPelaku | null): boolean {
  return hak?.workbasketSurveyor === true
}

/**
 * Tabel komite di bawah inbox: hanya bagi anggota roster komite FACIN (`HakPelaku.komite`); tanpa switch, tidak ikut tab
 * Process / Resolve.
 */
export function tabelKomiteTampil(hak: HakPelaku | null): boolean {
  return hak?.komite === true
}

/** Tombol Add Claim tampil hanya di tab Process saat switch Teknik mati. */
export function bolehTambahKlaim(tab: TabInbox, teknik: boolean): boolean {
  return tab === 'proses' && !teknik
}
