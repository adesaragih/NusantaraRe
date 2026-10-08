// Label modul Komite Claim Prop. Label LAYAR KASUS (judul, medan, judul grid, kolom, isian, tombol) datang dari server
// VERBATIM Section ShowTransfer (`models.Layar`). Berkas ini hanya memuat teks halaman awal dan perangkat yang tidak
// punya label di XML (ditandai).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 900) = nama folder korpus. */
export const MENU_KCP = { kelompok: 'Komite Claim Prop' } as const

/**
 * Daftar kerja penyetuju (worklist assignment "KomiteRouter"). Worklist itu tidak punya section inbox di korpus;
 * kolom mengikuti konvensi daftar kerja Komite Claim Life (prompt §8) - `[tidak ada di korpus]`.
 */
export const KCP = {
  judul: 'Komite Claim Prop',
  memuat: 'Loading…',
  kosong: 'No committee case is waiting for your decision.',
  kolomKasus: 'Case ID',
  kolomTgl: 'Update Date/Time',
  kolomStatus: 'Work Status',
  kolomKlaim: 'Claim',
  kolomTingkat: 'Level',
  kolomJabatan: 'Committe Name',
  kolomNilai: 'Adjustment RNM',
  kolomMataUang: 'Currency',
  kolomStatusBaris: 'Status',
  /** `[tidak ada di korpus]` - layar dibuka bukan oleh pemegang (mis. kasus sudah berpindah tingkat). */
  hanyaLihat: 'Read only - this case is waiting for another committee member.',
  memproses: 'Submitting…',
  /** `.AcceptStatus` `pyNoSelectionText` (ShowTransfer). */
  pilih: 'Choose',
  /** `[tidak ada di korpus]` - View more details: menu Claim Prop tidak dipegang akun ini (teks sama dengan Claim Prop). */
  berkasTakTerpasang: 'This file cannot be opened here: the menu of its module is not assigned to your account.',
} as const
