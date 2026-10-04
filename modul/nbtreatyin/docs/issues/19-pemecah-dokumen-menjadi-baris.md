# 19: Pemecah dokumen menjadi baris

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."* ⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, 232 di antaranya tampil di layar. Hasilnya di **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**. Data guide turun derajat jadi **penambal**, sah hanya untuk panjang maksimum per medan.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*
~~**Blocked by:** **16** · **17** · **18** · ⛔ `[data DBA]` **daftar kolom lengkap** — panduan bentuk dokumen dari DBA terbukti **basi**~~ ⛔ **gugur 23-09-2026 sore**
**Menutup:** NB AC **26–44** *(19 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-21..ID-31

## Hasil & nilai pengguna

Dokumen polis lama dipecah menjadi baris tabel: data umum, kuotasi, ceding, angsuran, penyebaran,
lapisan, dan riwayat usulan. Sesudah tiket ini, bentuk data polis **dapat dinyatakan** — bukan
potret halaman kerja.

## Yang dibangun

Pemecah dokumen dan penyusun baris, beserta aturan isi tiap tabel:

- **data umum** — medan skalar tingkat atas ditambah kolom datar dari tabel dokumen lama;
  ⛔ penanda lapisan **tidak** disimpan di sini, sebab di sistem lama ia pantulan baris pertama saja
- **ceding** — satu baris per ceding, id dan nama terpisah; ⭐ bentuk gabungan di data umum
  **disalin apa adanya**, tidak pernah dirangkai ulang
- **angsuran** — ⭐ **satu tingkat** pada bentuk proporsional
- **penyebaran** — pembagian rata menurut cacah baris, presisi sepuluh
- **lapisan** — pemegang penanda lapisan; ⚠️ potongan di sini adalah **uang**, bukan persen
- **riwayat usulan** — tabel yang sudah datar, ⛔ **tidak dibuat ulang**

⭐ **Penampung medan tak dikenal** disediakan — medan dokumen yang tidak dikenal **disimpan**,
bukan dibuang.

## Batas — yang TIDAK termasuk

⛔ Uji pulang-pergi dan dua bentuk dokumen — tiket **21**.
⛔ Pemuatan dokumen lama secara massal — tiket **22**.
⛔ Pohon objek pertanggungan — milik modul fakultatif, **di luar lingkup**.

## Cara mengujinya

Lewat seam `repository`. ⚠️ **Penampung medan tak dikenal wajib kosong** sebelum pekerjaan
dinyatakan selesai — isinya yang tidak kosong berarti sensus medan belum lengkap.

## Acceptance criteria

- [ ] **AC 26** — penanda lapisan **tidak** ada di tabel data umum
- [ ] **AC 27** — tabel kuotasi menyimpan sepuluh medan, termasuk jenis kontrak dan kode bisnis
- [ ] **AC 28** — satu baris per ceding, id dan nama terpisah
- [ ] **AC 29** — bentuk gabungan di data umum tersimpan **persis** seperti di dokumen
- [ ] **AC 30** — menghapus satu ceding menghapus **barisnya**, bukan menyunting teks gabungan
- [ ] **AC 31** — polis proporsional tersimpan **tanpa** baris rincian angsuran bertingkat
- [ ] **AC 32** — polis proporsional tersimpan **tanpa** baris lapisan
- [ ] **AC 33** — menyimpan baris lapisan pada polis proporsional **ditolak**
- [ ] **AC 34** — penanda lapisan tersimpan di tabel lapisan, bukan di data umum
- [ ] **AC 35** — potongan pada lapisan bertipe **uang**
- [ ] **AC 36** — penyebaran pada polis baru memakai presisi **sepuluh**
- [ ] **AC 37** — tiga medan bagian tersimpan sebagai **persentase**
- [ ] **AC 38** — medan potongan dan total bagian tersimpan sebagai **persentase**
- [ ] **AC 39** — tabel riwayat usulan **tidak dibuat ulang**
- [ ] **AC 40** — keterangan usulan dipotong pada batas panjangnya
- [ ] **AC 41** — keputusan usulan diturunkan dari penandanya, per baris
- [ ] **AC 42** — penanda per baris usulan tersimpan **terpisah** dari penanda tingkat polis
- [ ] **AC 43** — kolom pelaku terisi dari identitas login
- [ ] **AC 44** — kolom penanggung jawab terisi dari nama tampilan

## ⛔ Kenapa tiket ini `blocked`

~~`[data DBA]` **Panduan bentuk dokumen dari DBA terbukti basi.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."*

⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, **232** di antaranya tampil di layar. Sumbernya: **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**. Sapuan itu menutup dua titik buta yang sebelumnya meleset empat kali: **dua konvensi tag yang berlawanan**, dan **rujukan relatif di dalam loop**.

⚠️ Data guide turun derajat jadi **penambal** — sah untuk panjang maksimum per medan, dan untuk 24 nama yang ditulis sistem sehingga tidak tersapu aturan *(`OldData` · `ProdKe` · `EDMNo` · `OldPolicyNo` · `BranchCode`)*.

⛔ **394 tetap batas bawah, bukan total.** Penampung medan tak dikenal tetap wajib, dan wajib **kosong** sebelum pekerjaan dinyatakan selesai.

⭐ **Yang dapat dikerjakan lebih dulu tanpa menunggu:** kerangka pemecah, penampung medan tak
dikenal, dan seluruh AC yang tidak bergantung pada daftar kolom — ⚠️ tetapi **pekerjaan tidak
dapat dinyatakan selesai** sampai penampung itu kosong.