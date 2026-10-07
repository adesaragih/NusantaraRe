// Aturan halaman awal Claim Prop (keputusan work owner 07-10-2026): dua tab, Process dan Resolve. Tab Process memilih
// workbasket: Admin = Assignment2 "Outstanding Claim" (worklist pembuat, `Flow_TreatyIn`), Teknik = Assignment1 "Input
// Acceptation" (workbasket `ReasKlaimTeknik`). Pembuatan klaim (Start1 -> Assignment2) hanya dari workbasket Admin.

import type { JenisDaftar } from '../api'

export type TabInbox = 'proses' | 'selesai'
export type Workbasket = 'admin' | 'teknik'

export const TAB_INBOX: readonly TabInbox[] = ['proses', 'selesai']
export const WORKBASKET: readonly Workbasket[] = ['admin', 'teknik']

/** Daftar yang diminta ke server untuk tab dan workbasket terpilih. */
export function jenisDaftar(tab: TabInbox, wb: Workbasket): JenisDaftar {
  if (tab === 'selesai') return 'selesai'
  return wb === 'admin' ? 'saya' : 'workbasket'
}

/** Tombol Add Claim tampil hanya di tab Process, workbasket Admin. */
export function bolehTambahKlaim(tab: TabInbox, wb: Workbasket): boolean {
  return tab === 'proses' && wb === 'admin'
}
