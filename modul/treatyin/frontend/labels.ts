// Label layar modul Treaty In.
//
// ⛔ Nama himpunan di sini adalah nama ENTITAS spec (`KAMUS-KOLOM.md` §10.22),
// bukan label korpus Pega: keenam tabel acuan TIDAK ADA di sistem lama - ADR-0038
// membuatnya menggantikan delapan kolom bernama bahaya di kepala kontrak.

import type { HimpunanAcuan } from './api'

export const MENU_TREATYIN = {
  acuan: 'Treaty In — Tabel Acuan',
} as const

export const LABEL_HIMPUNAN: Record<HimpunanAcuan, string> = {
  'mata-uang': 'Mata Uang',
  'jenis-potongan': 'Jenis Potongan',
  'kelas-bisnis': 'Kelas Bisnis',
  'kelompok-treaty': 'Kelompok Treaty',
  bahaya: 'Bahaya',
  'jenis-reasuransi': 'Jenis Reasuransi',
}

export const ACUAN_TREATYIN = {
  kolomKode: 'Kode',
  kolomNama: 'Nama',
  kolomAktif: 'Aktif',
  kolomInduk: 'Induk',
  // Dua kalimat, dua medan `Kosong` yang berbeda: `pesan` menyatakan KEADAAN
  // (tabelnya memang belum berisi), `petunjuk` menyatakan SIAPA yang akan
  // mengisinya. Digabung jadi satu paragraf, yang kedua terbaca sebagai alasan
  // kosongnya - padahal ia jadwal, bukan sebab.
  kosong: 'Tabel acuan ini belum berisi.',
  kosongPetunjuk: 'Pemindahan isinya dari sistem lama adalah tiket 44.',
} as const
