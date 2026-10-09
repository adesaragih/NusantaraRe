// Disalin dari `modul/claimprop/frontend/pages/inbox.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Aturan halaman awal Claim Prop (keputusan work owner 07-10 dan 08-10-2026): dua tab, Process dan Resolve. Tab Process
// bawaan = worklist pembuat (Assignment2 "Outstanding Claim", `Flow_TreatyIn` `ToCurrentOperator` - tanpa cek workbasket,
// XML apa adanya). Switch Teknik menampilkan Assignment1 "Input Acceptation" (workbasket `ReasKlaimTeknik`) dan hanya
// dapat dinyalakan anggota workbasket itu. Pembuatan klaim (Start1 -> Assignment2) hanya saat switch Teknik mati. Tabel
// komite di bawah inbox (menu Komite Claim Non Prop disembunyikan, perintah work owner 09-10-2026).

import type { HakPelaku, JenisDaftar } from '../api'

export type TabInbox = 'proses' | 'selesai'

export const TAB_INBOX: readonly TabInbox[] = ['proses', 'selesai']

/** Daftar yang diminta ke server untuk tab dan posisi switch Teknik. */
export function jenisDaftar(tab: TabInbox, teknik: boolean): JenisDaftar {
  if (tab === 'selesai') return 'selesai'
  return teknik ? 'workbasket' : 'saya'
}

/** Switch Teknik dapat dinyalakan hanya bila akun memegang workbasket ReasKlaimTeknik. */
export function switchTeknikAktif(hak: HakPelaku | null): boolean {
  return hak?.workbasketTeknik === true
}

/**
 * Tabel komite di bawah inbox: hanya bagi pemegang workbasket yang tercantum di roster EMAILKOMITE NONPROP; tanpa switch,
 * tidak ikut tab (pola Claim Prop).
 */
export function tabelKomiteTampil(hak: HakPelaku | null): boolean {
  return hak?.komite === true
}

/** Tombol Add Claim tampil hanya di tab Process saat switch Teknik mati. */
export function bolehTambahKlaim(tab: TabInbox, teknik: boolean): boolean {
  return tab === 'proses' && !teknik
}
