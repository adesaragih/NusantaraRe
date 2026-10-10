// Mode form kontrak — dipilih tombol di layar daftar.
//
//   ubah   tombol `Edit` (dan `Add` kontrak baru): seluruh medan dan tombol
//          tabel (`Add`/`Delete`) dapat dipakai.
//   lihat  tombol `View`: baca-saja; tombol `Add`/`Delete` tidak dirender —
//          ekspor menjaganya dengan `TreatyIn.ViewState !='1'`.
//
// ⛔ Padanan `TreatyIn.ViewState` sistem lama: `0` Edit, `1` View.
export type ModeForm = 'ubah' | 'lihat'
