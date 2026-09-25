> Modul  : Treaty In Adjustment · ronde TDA
> Dibuat : 2026-09-24
> Sifat  : skenario paritas. **Fase ini hampir tidak berlaku untuk ronde ini**, dan alasannya
>          ditulis sebagai kalimat, bukan diisi pengisi.

# Paritas ronde TDA

## 1. Kenapa ronde ini hampir tidak menghasilkan kelas paritas

`GRL-18` **tidak menggerakkan satu angka pun.** Ia mengatur siapa yang menulis sebuah nilai kode dan
sampai kapan — bukan bagaimana sebuah besaran dihitung. Nilai `TreatyIn.EDMState` warisan dibawa
apa adanya (`GRL-13`, ADR-0042), jadi tidak ada yang dihitung ulang dan tidak ada selisih yang dapat
muncul.

Keenam TDA lain ditutup oleh putusan yang **sudah punya kelas paritasnya sendiri** di ronde masing-
masing, dan tidak diulang di sini.

## 2. Satu hal yang tetap harus masuk daftar harapan run paritas

Bukan kelas cacat — **pengecualian pembacaan**, sekeluarga dengan `PC-4` ronde C:

> `TreatyIn.EDMState` **warisan tunduk penuh pada kriteria identik** ADR-0043. Ia angka yang
> **dipindahkan**, bukan mekanisme yang diberikan.

Ini perlu dinyatakan justru karena `GRL-18` menyempitkan himpunan nilai yang sah **untuk versi
baru** menjadi dua. Pembaca laporan run paritas yang mengetahui penyempitan itu dapat menyimpulkan
bahwa baris warisan bernilai di luar himpunan tersebut adalah **kegagalan migrasi**.

**Ia bukan.** Ia data lama yang sah dipindahkan apa adanya, dan justru kombinasi yang tidak lagi
ditawarkan layar itulah yang paling perlu bertahan utuh — sebab tidak ada cara lain
memunculkannya kembali.

| Yang diperiksa | Harapan |
|---|---|
| cacah baris per nilai `TreatyIn.EDMState` | **identik** antara sumber dan hasil migrasi, untuk **setiap** nilai — termasuk nilai yang tidak ada di lima kombinasi sah lewat layar |
| baris ber-`TreatyIn.EDMMaterialType` terisi | **identik**; nilainya dibawa beserta asal-usulnya (`GRL-06`, ADR-0042) |

## 3. Yang belum dapat ditulis

Populasi dalam angka — belum ada. Tidak ada instans basis data yang terjangkau dan izin kueri
baca-saja belum turun.

Bentuknya ditulis sekarang justru karena itu: **bentuk harus mendahului angka**, atau uji terima
berubah jadi latihan menjelaskan (ADR-0043).
