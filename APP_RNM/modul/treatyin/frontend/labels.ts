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
  kosong: 'Tabel acuan ini belum berisi. Pemindahan isinya dari sistem lama adalah tiket 44.',
  keterangan:
    'Enam himpunan yang dapat bertambah (tiket 15, ADR-0038). Menambah satu baris tidak menyentuh kontrak yang sudah tercatat dan tidak menuntut perubahan skema.',
} as const
