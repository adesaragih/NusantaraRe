---
status: aktif
golongan: pelestarian
---

# 40: Pengisi kontrak melihat nilai versi sebelumnya berdampingan dengan versi yang sedang disusun

*Asal: `DAFTAR-PEKERJAAN.md` `P-44` · `ADR-0048` butir 2 · `GRL-03`, `GRL-10` · `SPEC-INVARIAN.md` `INV-58`, `INV-59`.*

**What to build:** **PK** yang menyusun versi penyesuaian melihat **nilai versi dasarnya berdampingan**
dengan nilai yang sedang ia isi — dan nilai lama itu **di-SELECT lewat `ID_VERSI_KONTRAK_DASAR`,
bukan disalin** ke kolom mana pun.

**Kenapa begini:** Sistem lama menyalin **seluruh halaman clipboard** versi dasar ke dalam sebuah
pohon bernama `OLDDATA` di dalam `JSONDATA` versi baru. Salinan itu melahirkan **dua sumber kebenaran
untuk satu fakta**: ketika versi dasarnya dibetulkan, salinannya **tidak ikut berubah**, dan selisih
yang dihitung terhadap salinan basi **tidak menghasilkan galat** — ia menghasilkan angka yang salah
dengan tenang. `ADR-0048` butir 2 membalikkannya: **tidak ada tabel `OLDDATA`**; sisi lama di-`SELECT`
ke versi dasar lewat rujukannya.

**Persyaratan:** `ADR-0048` butir 2 · `INV-59` (tidak ada nilai yang disalin dari entitas lain; hilir
diberi rujukan) · `INV-58` (tidak ada kolom turunan yang disimpan, kecuali ia fakta terbukukan yang
membawa penunjuk asalnya) · `GRL-11` (turunan *"versi berlaku"* **dihitung, bukan disimpan**)

**Tidak termasuk:** **Kolom `ID_VERSI_KONTRAK_DASAR` itu sendiri** — **tiket `01`** yang
membangunnya, beserta constraint keterisiannya. Irisan ini **membacanya**.
**Baris selisih** — tiket `06`. Yang di sini **penyandingan untuk dibaca orang**; yang di sana
**besaran berselisih yang dibukukan**.
**Pohon `ActualValue`** — **dibuang seluruhnya** (`GRL-14`); premi aktual menjadi nilai versinya
sendiri. Pernyataan keputusan, bukan pekerjaan tertunda.
**Layar penyandingnya** — `L-4`; yang dibangun **jalur bacanya**.

**Jalur gagal:** Sebuah kolom pada `VERSI_KONTRAK` yang memuat salinan nilai versi dasar ->
**tidak ada**, dan ketiadaannya diperiksa dengan sapuan atas `KAMUS-KOLOM.md` · Versi dasar dibetulkan
sesudah versi baru dibuat, lalu nilai lama dibaca -> **menampilkan nilai yang sudah dibetulkan**,
bukan yang basi · Versi pertama sebuah kontrak — dasarnya kosong -> **jalur bacanya menyatakan tidak
ada versi sebelumnya**, bukan galat dan bukan nilai kosong yang menyerupai nol.

**Uji:** **Negatif:** minta nilai versi sebelumnya untuk versi pertama; pastikan hasilnya **pernyataan
ketiadaan**, bukan nol. Sapu `KAMUS-KOLOM.md` untuk kolom salinan.
**Positif — dan ia yang membuktikan SELECT benar-benar SELECT:** buat versi penyesuaian, lalu
**betulkan sebuah nilai pada versi dasarnya**, lalu baca nilai lamanya. Hasilnya **harus nilai yang
sudah dibetulkan**. Sebuah rancangan yang diam-diam menyalin lulus setiap uji negatif di atas — sebab
salinan dan rujukan **memberi jawaban yang sama sampai salah satunya berubah**.

**Menggantikan:** pohon `OLDDATA` di dalam `JSONDATA`. Golongannya **PELESTARIAN** — kemampuan
melihat nilai lama berdampingan memang sudah ada — tetapi **caranya berubah**, dan bunyi lamanya
disebut di sini supaya yang sudah membacanya tahu apa yang bergeser.

**Blocked by:** 14 · **01**

**Dasar:**
```
EVIDENCED(TreatyIn.OLDDATA@ekspor-2026-09 - salinan penuh halaman clipboard versi dasar di dalam JSONDATA)
        DECIDED(ADR-0048, GRL-03, GRL-10, GRL-11, GRL-14, INV-58, INV-59)
```

- [ ] jalur baca nilai versi dasar berdiri dan memakai `ID_VERSI_KONTRAK_DASAR` — diperiksa di rencana kueri, bukan diandaikan
- [ ] sapuan atas `KAMUS-KOLOM.md` membuktikan **nol** kolom salinan nilai versi dasar
- [ ] uji positif lulus: pembetulan pada versi dasar **terlihat** pada pembacaan berikutnya
- [ ] versi pertama menghasilkan **pernyataan ketiadaan**, bukan nol
- [ ] ketiadaan tabel `OLDDATA` dan pembuangan `ActualValue` tertulis sebagai **pernyataan keputusan** yang menyebut `ADR-0048` dan `GRL-14`
