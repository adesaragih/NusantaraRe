# 08: Baris adjustment, adjuster/consultant, dan pagar nilai

**Status:** dibangun 07-10-2026 — pemeliharaan Adjuster / Consultant di luar lingkup (OQ-CP-04) — RALAT 07-10-2026 (semula `ready-for-agent`)
**Blocked by:** 00 (PREFACTOR) · 06 (loss allocation dan spreading) · 07 (estimasi)
**Menutup:** AC 44 · 45 · 46 · 47 · 48 · 49 · 95 · **126** · **127** *(9 AC)* — US 24–29

## Hasil & nilai pengguna

Claim Admin menambahkan baris adjustment sebagai **unit keputusan** klaim — satu klaim dapat memuat
beberapa keputusan berbeda, masing-masing dengan statusnya sendiri dan data bank penerimanya. Bila
validasi gagal, **tidak ada baris yatim** yang tertinggal. Kelebihan terhadap estimasi ditampilkan
dengan label yang jujur: yang memblokir memang memblokir, yang sekadar peringatan terlihat sebagai
peringatan.

## ⚠️ Baris beku — MENEGAKKAN, **perilaku BARU**

`[keputusan work owner]` 2026-09-19 — **beku adalah SIFAT TURUNAN sebuah baris penyesuaian:**
ia beku **bila ada kasus komite yang menunjuknya**, **berjalan maupun selesai**. ⛔ **Tidak ada
kolom penanda** — keadaan itu **dihitung**, bukan dibaca dari kolom.

**Yang ditegakkan tiket ini, dua jalur:**

| Jalur | Perilaku |
| --- | --- |
| **setiap jalur sunting** baris penyesuaian | memeriksa ada-tidaknya kasus komite yang menunjuk baris itu, dan **menolak bila ada** — **selamanya** |
| **setiap jalur hapus satu-baris** | idem — **menolak bila ada** |

⚠️ **Penegakannya DI LAPISAN LAYANAN, bukan di layar.** Menyembunyikan tombol saja tidak cukup:
permintaan yang datang langsung ke layanan **tetap harus ditolak**. Sejalan dengan **ADR-0014**.

⚠️ **Penolakan komite tidak mencairkannya.** Perbaikan atas isi baris dilakukan dengan **baris
penyesuaian baru**, bukan dengan menyunting baris lama.

⛔ **Hapus KLAIM bukan urusan tiket ini** — itu **tiket 00**. Tiket ini hanya menahan sunting dan
hapus **satu-baris**.

⚠️ `[terverifikasi]` **Pega tidak punya kunci ini** — nol pemeriksaan, nol pesan, nol penanganan.
**Dibangun, bukan dimigrasikan.**

## Area codebase

Entitas baris adjustment · data bank · adjuster/consultant · pagar nilai terhadap estimasi ·
spreading pada baris adjustment.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/AddAdjustment_Act.xml` | 5 | baris adjustment dibentuk; snapshot loss allocation dan spreading disalin |
| `Activity/AddAdjustment_Act.xml` | 9 | adjuster dan consultant wajib — syarat di-AND dengan panjang daftar = 1, jadi **hanya berlaku pada baris pertama** |
| `Activity/AddAdjustment_Act.xml` | 11 · 12 | ⚠️ keluar activity **sebelum** step rollback, sehingga baris kosong tertinggal |
| `Activity/CountValueADJTreaty_Act.xml` | 12 | menghitung nilai dalam mata uang asal, lalu nilai dalam IDR |
| `Activity/CountValueADJTreaty_Act.xml` | 16 | penanda yang **memblokir**; baris syarat pertama **mati**, baris kedua (total) hidup |
| `Activity/CountValueADJTreaty_Act.xml` | 17 | pesan yang **tampil**; kedua baris hidup — dan membandingkan besaran **berbeda** dari step 16 |
| `Activity/CountValueADJTreaty_Act.xml` | 2 · 18 | pesan salvage disiapkan; ditegakkan atas baris bertipe salvage bernilai positif |
| `Activity/CountValueADJTreaty_Act.xml` | 19 | membersihkan penanda saat total masih di bawah estimasi |
| `Activity/CountSpreadingADJ_Act.xml` | 6 | ⚠️ saringan payment type **tidak ditulis** — spreading berjalan untuk **semua** payment type |
| `Activity/DeleteAjsutment_Act.xml` | — | penghapusan baris |
| `Activity/SaveAdjusterConsultant_Act.xml` | 5 | ⚠️ step SQL **di-remark** → tidak ditulis sama sekali |
| `Activity/SaveAdjusterConsultant_Act.xml` | **4** | ⭐ **BARU 2026-09-19** — penyimpanan (`Obj-Save`) punya **jalur kegagalan** di keluarga gerbang **kedua**: syarat *"langkah barusan gagal"* → **lompat ke tanda `FAIL`, yaitu langkah 8**. ⚠️ **Langkah 8 bergerbang MATI**, jadi lompatannya sampai ke langkah yang **tidak berbuat apa-apa**, lalu alur berlanjut ke langkah 9. **Kegagalan penyimpanan diam-diam diabaikan.** `[terverifikasi]` · `[terbuka]` — **menunggu work owner**: apakah perilaku itu dipertahankan. ⛔ Tiket ini **tidak menetapkan** apa yang seharusnya terjadi saat gagal |

## ADR terkait

**ADR-0011** (unit keputusan = baris `AdjustmentList`) · **ADR-0003** (uang non-float).

## Acceptance criteria

- [ ] `[terverifikasi]` Unit keputusan adalah **baris adjustment**, masing-masing dengan statusnya sendiri *(AC 44 spec)*
- [ ] `[terverifikasi]` Adjuster dan consultant wajib terisi **hanya sebelum baris adjustment pertama**; baris ke-2 dan seterusnya tidak menuntutnya. Ditiru apa adanya *(AC 45 spec)*
- [ ] ⚠️ Validasi gagal **tidak meninggalkan baris adjustment yatim**. **Alasan menyimpang:** di Pega step keluar mendahului step rollback, sehingga baris kosong tertinggal di klaim *(AC 46 spec)*
- [ ] Baris adjustment menyimpan nama bank, id bank, nomor rekening, dan kode SWIFT *(AC 47 spec)*
- [ ] ⚠️ **RISIKO DITERIMA SADAR — dipertahankan, bukan diperbaiki.** Perhitungan spreading berjalan untuk **semua** payment type; saringan payment type tidak ditulis. **Alasan menyimpang:** konsekuensi langsung aturan induk — gerbang yang flag-nya mati tidak ditulis. Diterima `[keputusan work owner]`, bukan kelalaian migrasi *(AC 48 spec)*
- [ ] ⚠️ Pesan dan penanda blokir memakai **satu besaran yang sama, dalam IDR**, dibandingkan terhadap **estimasi terkini** (tiket 07). **Kelebihan total memblokir penyimpanan; kelebihan per baris menjadi peringatan berlabel jelas.** **Alasan menyimpang:** di Pega step 16 dan step 17 membandingkan besaran berbeda — mata uang asal versus IDR — dan step 16 baris pertama mati, sehingga pengguna melihat teks error tanpa ada yang memblokir. Kelebihan per baris memang tidak boleh memblokir: baris salvage wajib bernilai negatif, sehingga satu baris positif bisa melebihi estimasi sementara hasil bersihnya tidak *(AC 49 spec)*
- [ ] `[terverifikasi]` **Adjuster dan Consultant satu master, dua peran** — master tanpa kolom tipe. Satu klaim maksimal satu adjuster dan satu consultant *(AC 95 spec)*
- [ ] ⚠️ Baris yang **ada kasus komitenya** — berjalan maupun selesai — **tidak dapat diubah, selamanya**; **setiap** jalur sunting menolaknya **di lapisan layanan**, bukan hanya di layar. Test yang berhasil menyunting lewat layanan langsung **gagal**, dan test yang berhasil menyunting **sesudah komite menolak** juga **gagal**. **Alasan menyimpang:** di Pega tidak ada pemeriksaan apa pun — baris boleh berubah sesudah diserahkan *(AC 126 spec)*
- [ ] ⚠️ Baris beku **tidak dapat dihapus satu per satu**; upaya menghapusnya **ditolak** di lapisan layanan. Perbaikan memakai **baris penyesuaian baru**. **Alasan menyimpang:** di Pega baris boleh hilang sesudah diserahkan, tanpa pesan apa pun *(AC 127 spec)*

## Perintah verifikasi

```
jalankan test "baris pertama tanpa adjuster -> DITOLAK"
jalankan test "baris kedua tanpa adjuster -> diterima (paritas)"
jalankan test "validasi gagal -> nol baris adjustment tertinggal"
jalankan test "total adjustment > estimasi -> penyimpanan DIBLOKIR"
jalankan test "satu baris > estimasi tetapi total <= estimasi -> PERINGATAN, penyimpanan lolos"
jalankan test "baris salvage bernilai positif -> ditolak dengan pesan salvage"
jalankan test "spreading berjalan untuk payment type di luar 1|2|5 (risiko diterima sadar)"
jalankan test "satu klaim maksimal satu adjuster dan satu consultant"
```
