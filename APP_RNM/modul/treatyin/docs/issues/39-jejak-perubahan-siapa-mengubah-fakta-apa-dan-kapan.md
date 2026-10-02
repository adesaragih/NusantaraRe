---
status: aktif
golongan: baru
---

# 39: Pemeriksa jejak melihat siapa mengubah fakta apa dan kapan, untuk kontrak mana pun

*Asal: `DAFTAR-PEKERJAAN.md` `P-47` · `SPEC-MODEL-DATA.md` §10.21 · `ADR-0045` · `ADR-0044` · `KTV-D`.*

**What to build:** **PJ** membuka kontrak mana pun dan melihat **siapa mengubah ruas apa, dari nilai
apa ke nilai apa, kapan, dan di bawah peran apa**. Jejaknya terisi pada **setiap penulisan**, bukan
hanya pada perpindahan persetujuan.

Artefak: entitas `JEJAK_PERUBAHAN` dengan kedelapan atributnya, dan mekanisme pengisiannya.

**PEMBUAT PERTAMA** untuk `JEJAK_PERUBAHAN`.

**Kenapa begini:** **Jejak perubahan adalah bukti bahwa pembekuan benar-benar terjadi.** `ADR-0036`
menetapkan angka yang disetujui dibekukan, dan modul Adjustment mengambil data lama lewat `SELECT`
alih-alih menyalinnya. Keduanya hanya dapat dipercaya bila ada cara **membuktikan data lama tidak
berubah di antara dua pembacaan**. Tanpa jejak, pembekuan adalah **janji**, bukan fakta — dan selisih
yang dibukukan Adjustment berdiri di atas janji.

Di sistem lama satu-satunya jejak adalah daftar komentar, ditulis **hanya pada perpindahan
persetujuan**. Penyimpanan biasa — **termasuk yang mengubah angka** — tidak meninggalkan jejak siapa
maupun kapan, dan **tidak ada kolom waktu sama sekali** di tabel kontrak maupun tabel addendum.

**Persyaratan:** **`ADR-0045`** — isi minimalnya *"siapa, kapan, apa yang berubah dari nilai apa ke
nilai apa, dan **di bawah peran apa**"* · **`KTV-D`** — `PERAN_PELAKU` adalah **POTRET**: teks, final,
**tidak dinormalisasi ulang** ketika entitas peran kelak lahir · `ADR-0041` (satu fakta, satu penulis) ·
`INV-01`

**Tidak termasuk:** **Catatan persetujuan** — benda yang **berbeda**, dan `ADR-0045` memisahkannya
tegas: catatan adalah narasi manusia yang **dapat disunting**, jejak adalah fakta mesin yang **tidak
dapat**. Menggabungkannya menghasilkan yang terburuk dari keduanya — catatan yang bisa disunting dan
karena itu tidak membuktikan apa-apa. `CATATAN_PERSETUJUAN` milik batch 2.
**Entitas peran dan penugasan bertanggal** — **`F-16`**, lubang terbuka: `ADR-0044` menempatkan
keduanya di gelombang 1 dan §10 tidak memuat satu pun. **Ia tidak menahan irisan ini**, sebab `KTV-D`
sudah memutuskan bentuk `PERAN_PELAKU`, dan bentuk itu **tidak berubah** apa pun jadinya entitas peran
nanti.
**Layar pemeriksa jejak** — `L-4`: tidak ada spesifikasi layar di mana pun, dan lubang itu milik batch
lapisan aplikasi. Yang dibangun di sini **jalur bacanya**, bukan tampilannya.

**Jalur gagal:** Sebuah penulisan yang **tidak** meninggalkan baris jejak -> **kriteria selesai tidak
terpenuhi** · Baris jejak yang disunting sesudah ditulis -> **ditolak**; entitas ini **tambah-saja** ·
`PERAN_PELAKU` kosong -> ditolak · Baris jejak dihapus -> ditolak.

**Uji:** **Negatif:** `UPDATE` atas baris jejak; `DELETE` atas baris jejak; simpan jejak tanpa peran;
lakukan penulisan lewat jalur yang **tidak** menuliskan jejak dan pastikan jalurnya **tidak ada**.
**Positif — dan ia yang menangkap jejak yang terlalu sempit:** ubah **satu ruas yang bukan bagian
persetujuan** — misalnya `LINGKUP_WILAYAH` — dan pastikan jejaknya **tetap tercatat**. Sistem lama
mencatat hanya pada perpindahan persetujuan; mekanisme yang meniru bentuk lama lulus ketiga uji
negatif di atas dan **melewatkan justru perubahan yang paling sering terjadi**.
**Positif kedua:** dua perubahan pada **detik yang sama** menghasilkan **dua baris**. `WAKTU_PERUBAHAN`
bertipe `DATE` beresolusi detik, dan entitas ini **sengaja tanpa kunci alami** justru supaya keadaan
itu sah.

**Menggantikan:** daftar komentar sistem lama sebagai satu-satunya jejak. **Golongannya BARU** —
jejak perubahan sebagai fakta mesin memang belum pernah ada — sehingga `CARA MENYALAKANNYA` wajib, dan
**hasilnya tidak punya pembanding sama sekali** di data lama.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemeriksa jejak (PJ), bulanan**, atas cacah baris jejak per kontrak. Angka nol pada kontrak yang jelas berubah adalah tanda mekanismenya tidak terpasang di salah satu jalur |
| 3 | **ambang berangka** — jejak dinyatakan lengkap ketika **setiap** penulisan atas `VERSI_KONTRAK` dan anaknya menghasilkan sedikitnya satu baris jejak, diperiksa atas contoh **100 penulisan berturut-turut** dengan hasil **100** |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In, bersama `PJ` |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(ADR-0045 konteks - daftar komentar sebagai satu-satunya jejak; nol kolom waktu di tabel kontrak dan addendum)
        DECIDED(ADR-0045, ADR-0044, KTV-D, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(F-16)
```

- [ ] `JEJAK_PERUBAHAN` berdiri dengan **kedelapan** atributnya sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] entitas ini **tambah-saja**: `UPDATE` dan `DELETE` atasnya ditolak, keduanya diuji terpisah
- [ ] `PERAN_PELAKU` tersimpan sebagai **teks beku**, dan `KTV-D` dirujuk di dalam berkas DDL-nya sebagai pernyataan keputusan
- [ ] kedua uji positif lulus — ruas non-persetujuan tetap berjejak, **dan** dua perubahan sedetik menghasilkan dua baris
- [ ] tidak ada jalur penulisan yang melewati jejak — diperiksa dengan sapuan, hasilnya dicetak
- [ ] `F-16` tercatat di `ASUMSI-CLEAR.md` sebagai lubang yang **tidak menahan** tiket ini, beserta sebabnya
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka
