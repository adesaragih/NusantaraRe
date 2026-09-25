# Tiket - Klaim Non-Jiwa

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Claim Fac In | **15** | 14 | 0 | 0 | 0 | 1 |
| Claim Prop | **15** | 15 | 0 | 0 | 0 | 0 |
| Claim Non Prop | **38** | 0 | 0 | 0 | 0 | 38 |
| **Jumlah** | **68** | **29** | **0** | **0** | **0** | **39** |

---

# Claim Fac In

Jumlah tiket: **15**

## Claim Fac In - 01 - Registrasi klaim, pengikatan polis, dan daur hidup kasus

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR)
**Menutup:** AC 4 · 5 · 6 · 7 · 8 · 9 · 10 · 11 · 12 *(9 AC)* — US 1–5

#### Hasil & nilai pengguna

Hari ini Sebuah klaim fakultatif masuk **belum bisa dibuat**. Tahapan yang boleh dilaluinya juga belum ditetapkan, sehingga tidak ada yang mencegah kasus melompat ke tahap yang tidak sah.

Sesudah tiket ini, Penilai dapat **mendaftarkan klaim**, mengikatnya ke polis yang benar, dan melihat kasus berjalan melalui **tahapan yang ditetapkan satu berkas alur** — ⛔ tahap di luar itu tidak dibuat.

#### Area codebase

- Lapisan layanan klaim: pembuatan kasus dan pengikatan polis
- Mesin tahapan kasus — daftar tahap tertutup

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Tahapan kasus | satu berkas alur; tahap di luarnya tidak dibuat |
| Pengikatan polis | pembacaan data polis dan kutipan dari modul penawaran |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0014** — belum menyentuh tiket ini; wewenang di tiket 08

#### Acceptance criteria

- [ ] **AC 4** — perilaku yang ditiru adalah perilaku **salinan PRODUKSI**; membandingkan terhadap salinan pengembangan **gagal**
- [ ] **AC 5 · 6 · 7** — tahapan kasus mengikuti berkas alur; tahap tambahan **ditolak**
- [ ] **AC 8–12** — registrasi dan pengikatan polis

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | ⚠️ menahan pengikatan polis yang lengkap |

#### Perintah verifikasi

1. Buat satu klaim, periksa ia terikat ke polis yang benar.
2. Coba pindahkan kasus ke tahap yang **tidak ada** di berkas alur — ⛔ **ditolak**.
3. Tutup lalu buka kembali kasus — ⭐ tahapnya **tidak mundur**.

## Claim Fac In - 02 - Klasifikasi lini produk — 13 penggolong hidup, 36 tidak dialihkan

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi)
**Menutup:** AC 13 · 14 · 15 · 16 · 17 · 18 *(6 AC)* — US 34–36

#### Hasil & nilai pengguna

Hari ini Penggolongan lini usaha **membaca medan data kutipan yang tidak pernah disalin** ke objek kerja. ⛔ `[terverifikasi]` Akibatnya cabang **MBU** dan **Travel** tidak pernah terbit — diam-diam, tanpa galat.

Sesudah tiket ini, Setiap lini usaha **tergolong benar**, termasuk MBU dan Travel yang dulu tak pernah terbit. ⭐ Penggolongan membaca **data kutipan yang lengkap**.

#### Area codebase

- Lapisan layanan: penggolongan lini usaha
- Penyalinan data kutipan dari kasus induk

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penggolong lini usaha | 49 rule penggolong; ⭐ **13 dipakai hidup**, ⛔ **36 tidak dipakai sama sekali** |
| Sumber data kutipan | salinan data kutipan ke objek kerja |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0003** — dua kaki dasar klasifikasi

#### Acceptance criteria

- [ ] **AC 13–18** — klasifikasi lini produk
- [ ] ⭐ **13 penggolong** dialihkan; ⛔ **36 tidak dibangun**
- [ ] ⭐ Cabang **MBU** dan **Travel** **terbit** bila datanya memenuhi — ⚠️ inilah cacat lama yang ditutup

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi rule peran | ⛔ tidak menyentuh tiket ini |

#### Perintah verifikasi

1. Jalankan satu klaim tiap lini usaha yang hidup — ⭐ ketiga belasnya tergolong benar.
2. Jalankan klaim lini **MBU** dan **Travel** — ⭐ cabangnya **terbit**; ⛔ di sistem lama tidak.
3. Cari pemanggilan salah satu dari 36 penggolong yang tidak dialihkan — ⛔ **nihil**.

## Claim Fac In - 03 - Estimasi nilai kerugian dan validasinya

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi)
**Menutup:** AC 19 · 20 · 21 · 22 · 23 · 24 · 25 · 26 · 27 · 28 · 29 · 30 *(12 AC)* — US 7 · 8

#### Hasil & nilai pengguna

Hari ini Estimasi kerugian **belum punya tempat menggantung** pada lini FAC, dan aturan validasinya belum ditetapkan.

Sesudah tiket ini, Penilai dapat **mencatat estimasi per item objek**, dengan mata uang, risiko sendiri, dan konversinya — dan sistem **menolak** estimasi yang melanggar aturan.

#### Area codebase

- Lapisan layanan klaim: estimasi
- Validasi nilai estimasi terhadap nilai pertanggungan

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Estimasi | daftar estimasi di dalam item objek |
| Validasi | rule validasi masukan estimasi |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit atas perubahan nilai

#### Acceptance criteria

- [ ] **AC 19–30** — estimasi dan validasinya
- [ ] ⭐ Estimasi menggantung pada **item objek**, bukan pada klaim

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **12** | ⚠️ **10 dari 13 kolom estimasi tanpa penulis di korpus** — diduga diisi lewat layar, ⛔ belum terbukti | ⚠️ **menahan jalur TULIS**, tidak menahan jalur BACA |

#### Perintah verifikasi

1. Catat estimasi pada satu item objek — ⭐ tersimpan dan terbaca kembali.
2. Catat estimasi melampaui nilai pertanggungan — ⭐ perilakunya sesuai AC.
3. Hapus item objek — ⭐ estimasinya **ikut terhapus**.

## Claim Fac In - 04 - Pembagian klaim per treaty, dan aturan keseragaman berbagi

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 03 (estimasi)
**Menutup:** AC 31 · 32 · 33 *(3 AC)* — US 8

#### Hasil & nilai pengguna

Hari ini Pembagian klaim antar treaty **belum punya tabel** pada lini FAC, dan aturan keseragaman persentase berbagi **belum ada penegaknya**.

Sesudah tiket ini, Pembagian klaim tercatat per treaty per mata uang, dan ⭐ **persentase berbagi seragam untuk satu treaty di dalam satu item objek** — sehingga daftar ringkasnya dapat diturunkan tanpa tabel kedua.

#### Area codebase

- Lapisan layanan klaim: pembagian per treaty
- ⭐ Penegak aturan keseragaman — ⛔ **di lapisan layanan**, bukan basis data

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pembagian klaim | daftar pembagian di dalam item objek |
| Daftar ringkas per treaty | ⛔ **bukan tabel** — himpunan bagian, selisih **NIHIL** |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit

#### Acceptance criteria

- [ ] **AC 31 · 32 · 33** — pembagian klaim per treaty
- [ ] ⭐ **Semua baris ber-treaty sama di dalam satu item objek punya persentase berbagi yang sama**
- [ ] ⭐ Menyunting persentase **mengenai seluruh baris treaty itu sekaligus** — ⛔ bukan satu baris

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **9** | ⚠️ Penegak aturan keseragaman belum ditetapkan — dijaga saat tulis, atau diperiksa berkala | ⚠️ **menahan bentuk penegakannya**, tidak menahan aturannya |

#### Perintah verifikasi

1. Sunting persentase berbagi satu baris — ⭐ **seluruh baris treaty itu ikut berubah**.
2. Coba simpan dua baris treaty sama dengan persentase berbeda — ⛔ **ditolak**.
3. Turunkan daftar ringkas per treaty — ⭐ **selisih NIHIL** terhadap pembagian penuh.

## Claim Fac In - 05 - Penyesuaian nilai klaim, pembagian dan quota share di atasnya

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 03 (estimasi) · 04 (pembagian)
**Menutup:** AC 34 · 35 · 36 · 37 · 38 *(5 AC)* — US 9 · 10

#### Hasil & nilai pengguna

Hari ini Baris penyesuaian **belum punya tempat**, dan pembagian serta quota share di atasnya belum terpisah dari pembagian tingkat item objek.

Sesudah tiket ini, Penilai dapat **mengajukan penyesuaian nilai klaim**, dengan pembagian dan quota share yang melekat **pada penyesuaian itu** — ⭐ terpisah dari pembagian tingkat item objek.

#### Area codebase

- Lapisan layanan klaim: penyesuaian
- Pembagian dan quota share tingkat penyesuaian

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penyesuaian | daftar penyesuaian di dalam item objek, ditambah bahan pertimbangan komite |
| Pembagian atas penyesuaian | dua daftar terpisah, ditulis **langkah bertetangga** di berkas yang sama |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit

#### Acceptance criteria

- [ ] **AC 34–38** — penyesuaian nilai klaim
- [ ] ⭐ Pembagian tingkat **penyesuaian** terpisah dari pembagian tingkat **item objek**
- [ ] ⭐ Enam medan bahan pertimbangan komite tersimpan **bersama penyesuaian**, bukan tabel sendiri

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **13** | Daftar lokasi di dalam retro fakultatif — larik bersarang | tidak menahan |

#### Perintah verifikasi

1. Buat satu penyesuaian, isi pembagian dan quota share-nya — ⭐ keduanya tersimpan terpisah.
2. Hapus penyesuaian — ⭐ keduanya **ikut terhapus**.
3. Bandingkan pembagian tingkat penyesuaian dengan tingkat item objek — ⭐ **baris berbeda**.

## Claim Fac In - 06 - Uang, ketelitiannya, dan batas transaksi

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 05 (penyesuaian)
**Menutup:** AC 39 · 40 · 41 · 42 · 65 · 66 · 67 · 114 *(8 AC)* — US 7 · 29

#### Hasil & nilai pengguna

Hari ini Nilai uang **berubah karena urutan pemanggilan**, dan tidak ada batas transaksi yang menjamin sekelompok tulisan selesai bersama atau gagal bersama.

Sesudah tiket ini, Angka uang **berhenti berubah karena urutan**, konversi mata uang punya ketelitian yang ditetapkan, dan sekelompok tulisan **selesai bersama atau gagal bersama**.

#### Area codebase

- Lapisan layanan: perhitungan uang dan konversi
- Batas transaksi pada penulisan berkelompok

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Perhitungan uang | penjumlahan bruto, bagian Nusantara Re, dan konversi ke rupiah |
| Batas transaksi | ⚠️ rule SQL lama membawa penyelesaian transaksi **di tengah** pekerjaan pemanggil |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit

#### Acceptance criteria

- [ ] **AC 39–42 · 114** — uang dan ketelitiannya
- [ ] **AC 65 · 66 · 67** — transaksi
- [ ] ⭐ Kegagalan separuh jalan **tidak meninggalkan data separuh tersimpan**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

#### Perintah verifikasi

1. Hitung total penyesuaian dua kali dengan urutan pemanggilan berbeda — ⭐ **hasilnya sama**.
2. Paksa gagal di tengah sekelompok tulisan — ⭐ **tak satu pun baris tersimpan**.
3. Periksa konversi ke rupiah — ⭐ ketelitiannya sesuai AC.

## Claim Fac In - 07 - Layar klaim — apa yang tampil dan apa yang dapat disunting

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi) · 03 (estimasi) · 05 (penyesuaian)
**Menutup:** AC 99 · 100 · 101 · 102 *(4 AC)* — US 6 · 7

#### Hasil & nilai pengguna

Hari ini Penilai **belum punya layar** untuk melihat klaim, objeknya, estimasinya, dan penyesuaiannya dalam satu tempat.

Sesudah tiket ini, Penilai melihat klaim **utuh dalam satu layar** — objek, item, estimasi, pembagian, penyesuaian — dan dapat menyunting yang memang boleh disunting.

#### Area codebase

- Antarmuka klaim
- Lapisan layanan: pembacaan klaim utuh

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Layar klaim | rule layar klaim fakultatif masuk |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0012** — wewenang eksplisit, bukan efek samping layar

#### Acceptance criteria

- [ ] **AC 99–102** — layar
- [ ] ⛔ Layar **tidak menjadi penjaga wewenang** — ia hanya menyembunyikan; penegakan di tiket 08

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **11** | ⚠️ **29 medan bergantung jenis objek** — menjadi kolom, atau dibaca dari polis | ⚠️ **menahan bentuk layar objek** |

#### Perintah verifikasi

1. Buka satu klaim — ⭐ objek, item, estimasi, pembagian, dan penyesuaian tampil.
2. Coba sunting medan yang tidak boleh disunting — ⭐ layar menolaknya.
3. ⭐ Panggil lapisan layanan **langsung** untuk menyunting medan itu — ⛔ **juga ditolak** *(tiket 08)*.

## Claim Fac In - 08 - Wewenang — ditegakkan di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi) · 05 (penyesuaian)
**Menutup:** AC 72 · 73 · 74 *(3 AC)* — US 19 · 20 · 21

#### Hasil & nilai pengguna

Hari ini ⛔ `[terverifikasi]` Korpus **tidak memuat satu pun rule otorisasi**, dan 13 medan privilese seluruhnya kosong. Siapa pun yang dapat memanggil lapisan layanan dapat mengerjakan apa pun.

Sesudah tiket ini, Wewenang **ditegakkan di lapisan layanan** — ⛔ bukan di layar. Percobaan yang tidak berwenang **ditolak dan terekam**, bukan diabaikan diam-diam.

#### Area codebase

- Lapisan layanan: penegakan wewenang
- Jejak audit atas percobaan yang ditolak

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Wewenang | ⛔ **nihil di korpus** — 13 medan privilese kosong; ADR-U-0014 mencatatnya ABSENT |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0014** — wewenang ditegakkan di lapisan layanan
- **ADR-U-0002** — peran ditegakkan di lapisan layanan, bukan visibilitas layar
- **ADR-U-0012** — wewenang eksplisit

#### Acceptance criteria

- [ ] **AC 72 · 73 · 74** — wewenang
- [ ] ⭐ Penolakan terjadi **di lapisan layanan**
- [ ] ⭐ Percobaan yang ditolak **terekam**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | ⛔⛔ **Isi daftar jabatan dan susunan jenjang BELUM ADA** — korpus tidak memuatnya | ⛔⛔ **MENAHAN PEMBANGUNAN tiket ini.** ⭐ Tiketnya tetap ditulis, jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum daftar jabatan diberikan work owner |

#### Perintah verifikasi

1. Panggil lapisan layanan sebagai pengguna **tanpa wewenang** — ⛔ **ditolak**.
2. Periksa jejak audit — ⭐ percobaan yang ditolak **tercatat**.
3. ⛔ Sembunyikan tombolnya di layar lalu panggil layanan langsung — ⛔ **tetap ditolak**.

## Claim Fac In - 09 - Dokumen akseptasi dan pencetakannya

**Status:** ready-for-agent
**Blocked by:** 05 (penyesuaian) · 06 (uang)
**Menutup:** AC 51 · 52 *(2 AC)* — US 28 · 32

#### Hasil & nilai pengguna

Hari ini Dokumen akseptasi **belum dapat terbit**, sehingga penilai tidak punya bukti tertulis atas penyesuaian yang disetujui.

Sesudah tiket ini, Dokumen akseptasi **tercetak dengan angka yang benar**, dan ⭐ **tidak tercetak dua kali** untuk akseptasi yang sama.

#### Area codebase

- Lapisan layanan: penerbitan dokumen
- Penanda dokumen sudah tercetak

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pencetakan akseptasi | rule pencetak dokumen akseptasi bermata-uang-banyak |
| Penjaga cetak ganda | penanda dokumen sudah dicetak, diuji di sembilan gerbang |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit

#### Acceptance criteria

- [ ] **AC 51 · 52** — dokumen
- [ ] ⭐ Dokumen tercetak **sekali** per akseptasi; percobaan kedua **ditolak**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

#### Perintah verifikasi

1. Terbitkan dokumen akseptasi — ⭐ angkanya cocok dengan penyesuaian.
2. Terbitkan ulang untuk akseptasi yang sama — ⛔ **ditolak**.
3. Periksa penanda tercetak — ⭐ terisi sesudah cetak pertama.

## Claim Fac In - 10 - Penyimpanan berkas lampiran

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi)
**Menutup:** AC 53 · 54 · 55 · 56 *(4 AC)* — US 5

#### Hasil & nilai pengguna

Hari ini Berkas lampiran klaim — surat, foto, laporan survei — **belum punya tempat simpan**, dan tautannya ke klaim belum ditetapkan.

Sesudah tiket ini, Penilai dapat **melampirkan berkas** pada klaim, dan berkas itu **tetap dapat dibuka** sesudah kasus ditutup.

#### Area codebase

- Lapisan layanan: penyimpanan berkas
- Tautan berkas ke klaim

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penyimpanan berkas | rule unggah dan pengambilan tautan berkas |
| Token penyimpanan | rule pengambil token penyimpanan |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit atas unggahan

#### Acceptance criteria

- [ ] **AC 53–56** — penyimpanan berkas
- [ ] ⭐ Berkas **tetap terbuka** sesudah kasus ditutup

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

#### Perintah verifikasi

1. Lampirkan satu berkas, tutup kasus, buka kembali — ⭐ berkas **masih terbuka**.
2. Lampirkan berkas bernama sama dua kali — ⭐ perilakunya sesuai AC.
3. Periksa jejak audit unggahan — ⭐ tercatat siapa dan kapan.

## Claim Fac In - 11 - Efek keluar — kasir, surel, dan konversi

**Status:** ready-for-agent
**Blocked by:** 06 (uang) · 09 (dokumen)
**Menutup:** AC 57 · 58 · 59 · 60 · 61 · 62 · 63 · 64 *(8 AC)* — US 30 · 31 · 33

#### Hasil & nilai pengguna

Hari ini Instruksi pembayaran, surel pemberitahuan, dan konversi klaim **belum terkirim ke mana pun**. ⚠️ Dan di sistem lama, **penjaga ganda-bayar bersandar pada gerbang yang belum dapat dipastikan bacanya**.

Sesudah tiket ini, Instruksi pembayaran terkirim ke kasir **tepat satu kali**, surel terkirim **hanya di lingkungan produksi**, dan ⭐ **pemberitahuan galat terkirim HANYA ketika pengiriman gagal**.

#### Area codebase

- Lapisan layanan: efek keluar
- ⭐ Penjaga ganda-bayar — **eksplisit**, dalam satu transaksi

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pengiriman ke kasir | rule panggil layanan kasir; dua panggilan luar, keduanya bergerbang lingkungan produksi |
| Penjaga ganda-bayar | ⚠️ gerbang lama memakai tanda sama dengan **tunggal** — ⛔ bacanya belum pasti |
| Pemberitahuan galat | ⚠️ gerbang lama **bendera mati** ⇒ terkirim **setiap kali**; medannya salah eja |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit atas efek keluar

#### Acceptance criteria

- [ ] **AC 57–64** — efek keluar
- [ ] ⭐ Pembayaran terkirim **tepat satu kali**; panggilan kedua **ditolak**
- [ ] ⭐ Penanda terkirim tersimpan **dalam transaksi yang sama** dengan pengirimannya
- [ ] ⭐ Pemberitahuan galat terkirim **hanya pada kegagalan**
- [ ] ⭐ Surel terkirim **hanya di lingkungan produksi**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **26** | ⚠️ Tanda sama dengan **tunggal** pada penjaga ganda-bayar — pembandingan atau penugasan | ⛔ **tidak menahan** — penjaga baru dibuat eksplisit; jawabannya untuk **memeriksa data lama** |
| **27** | Apakah pemberitahuan galat lama benar-benar sampai ke seseorang | tidak menahan |

#### Perintah verifikasi

1. Kirim satu instruksi pembayaran — ⭐ terkirim, penanda terisi.
2. Kirim ulang untuk akseptasi yang sama — ⛔ **ditolak**.
3. Paksa kegagalan pengiriman — ⭐ pemberitahuan galat **terkirim**.
4. Kirim yang berhasil — ⛔ pemberitahuan galat **TIDAK terkirim**.
5. Jalankan di lingkungan bukan produksi — ⛔ surel **tidak terkirim**.

## Claim Fac In - 12 - Jejak audit dan kronologi klaim

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi) · 06 (uang)
**Menutup:** AC 68 · 69 · 70 · 71 *(4 AC)* — US 24 · 25 · 26

#### Hasil & nilai pengguna

Hari ini Perubahan pada klaim **tidak meninggalkan jejak yang dapat dibaca**, sehingga tidak ada yang dapat menjawab siapa mengubah apa dan kapan.

Sesudah tiket ini, Setiap perubahan penting meninggalkan **catatan kronologi**, dan ⛔ **catatan itu tidak ikut terhapus** ketika klaimnya dihapus.

#### Area codebase

- Lapisan layanan: penulisan kronologi
- ⚠️ Perilaku hapus — ⛔ **JANGAN berantai**

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Kronologi klaim | daftar usulan pandangan pada data klaim; dua penulis |
| Penulis kedua | rule transformasi data penyisip kronologi |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit

#### Acceptance criteria

- [ ] **AC 68–71** — jejak audit
- [ ] ⭐ Catatan kronologi mencatat **akun**, **jabatan saat itu**, keputusan, dan waktu
- [ ] ⛔⛔ Menghapus klaim **TIDAK menghapus** kronologinya — ⭐ jejak yang ikut terhapus berhenti menjadi jejak

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar jabatan | ⚠️ menahan **pencatatan jabatan**, tidak menahan pencatatan akun |

#### Perintah verifikasi

1. Ubah satu nilai klaim — ⭐ kronologi bertambah satu baris.
2. Hapus klaimnya — ⛔ kronologinya **tetap ada**.
3. Naikkan jabatan pengubahnya — ⭐ catatan lama **tidak berubah**.

## Claim Fac In - 13 - Penyerahan ke komite — kontrak muatan

**Status:** ready-for-agent
**Blocked by:** 05 (penyesuaian) · 06 (uang) · 08 (wewenang)
**Menutup:** AC 43 · 44 · 45 · 46 · 47 · 48 · 49 · 50 *(8 AC)* — US 13 · 14

#### Hasil & nilai pengguna

Hari ini Penyesuaian di atas kewenangan penilai **belum dapat diserahkan ke komite**, dan tidak ada kesepakatan tentang **apa yang diserahkan**.

Sesudah tiket ini, Penilai dapat **menyerahkan penyesuaian ke komite**, dan ⭐ **muatan yang diserahkan lengkap** — termasuk data kutipan **utuh**, sehingga penggolongan lini usaha di sisi komite tidak pincang.

#### Area codebase

- Lapisan layanan klaim: penyerahan ke komite
- ⭐ Kontrak muatan — apa yang disalin dan seberapa lengkap

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pembentukan kasus komite | rule pembuat nomor komite; menyemai jumlah jenjang dan giliran mulai |
| Jalur satu jenjang | rule kirim tutup klaim dan kirim tolak klaim — ⭐ keduanya **berkomite satu** |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0014** — wewenang ditegakkan di lapisan layanan

#### Acceptance criteria

- [ ] **AC 43–50** — penyerahan ke komite
- [ ] ⭐ **Data kutipan disalin UTUH**, ⛔ bukan daftar medan bernama — ⚠️ daftar bernama itulah yang melahirkan cacat MBU dan Travel
- [ ] ⭐ Jalur **tutup klaim** dan **tolak klaim** membentuk komite **satu jenjang**
- [ ] ⛔ Tiket ini **berhenti di kontrak muatan** — ⭐ perilaku komite ada di putaran sisi komite

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar jabatan dan susunan jenjang | ⛔⛔ **MENAHAN** penentuan jumlah jenjang |

#### Perintah verifikasi

1. Serahkan satu penyesuaian ke komite — ⭐ kasus komite lahir dengan muatan lengkap.
2. Periksa data kutipan di sisi komite — ⭐ **utuh**, bukan dua medan.
3. Serahkan lewat jalur tutup klaim — ⭐ komitenya **satu jenjang**.

## Claim Fac In - 14 - Migrasi data lama, paritas yang ditegaskan, dan penamaan ulang

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · seluruh tiket 01–13
**Menutup:** AC 75 · 76 · 77 · 78 · 79 · 80 · 81 · 82 · 92 · 93 · 94 · 95 · 96 · 97 · 98 · 105 · 106 · 107 · 108 · 109 · 110 · 111 · 112 · 113 *(24 AC)* — US 27 · 36

#### Hasil & nilai pengguna

Hari ini Data klaim lama **belum berpindah**, nama kolom yang menyesatkan **belum diluruskan**, dan perilaku yang sengaja **dipertahankan sama** belum punya uji yang membuktikannya.

Sesudah tiket ini, Data lama **berpindah dengan benar**, nama yang menyesatkan **diganti dengan aturan yang tertulis**, dan ⭐ **paritas yang sengaja dipertahankan punya uji** sehingga tidak berubah diam-diam.

#### Area codebase

- Migrasi data klaim lama
- Penamaan ulang kolom dan aturannya
- Uji paritas

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penamaan menyesatkan | bab aturan menamai ulang pada spec |
| Paritas | perilaku yang sengaja ditiru apa adanya |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

#### ADR terkait

- **ADR-U-0007** — jejak audit atas migrasi

#### Acceptance criteria

- [ ] **AC 75–78** — migrasi
- [ ] **AC 79–82** — nama kolom dan pembacaan
- [ ] **AC 92–98** — paritas yang ditegaskan
- [ ] **AC 105–113** — sisa yang ditegaskan
- [ ] ⭐ Tiap perilaku paritas punya **uji yang membuktikannya** — ⛔ supaya tidak 'terperbaiki' diam-diam oleh pengembang berikutnya

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **29** | Cacah baris lama terdampak — `[data DBA]` | ⚠️ menahan **cacah**, tidak menahan cara migrasinya |
| **1** | Kolom tabel data kutipan — `[data DBA]` | ⚠️ menahan pemetaan kolom lama |

#### Perintah verifikasi

1. Migrasikan satu klaim lama — ⭐ seluruh tingkatnya terbentuk benar.
2. Jalankan uji paritas — ⭐ perilaku yang sengaja sama **tetap sama**.
3. Cari nama kolom yang menyesatkan di skema baru — ⛔ **nihil**.

## Claim Fac In - uru - Urutan tiket — Claim Fac In

⛔ **Ini bukan tiket.** Berkas ini hanya menjelaskan **tiket mana harus selesai sebelum tiket mana**.
Isinya **tidak menambah keputusan apa pun**.

Sumber: `claim-facin\spec.md` · `STRUKTUR-TABEL-CLAIM-FACIN.md` · `RELASI-TABEL-CLAIM-FACIN.md`.

> ⛔⛔ **RALAT — 2026-09-20.** Kalimat di bawah ini **SALAH**, dan dikutip utuh di sini alih-alih
> dihapus:
>
> *"⚠️⚠️ **Nomor di folder ini BERSAMBUNG ke putaran sisi komite.** `[keputusan work owner]`
> **K6** — komite lini FAC dan lini PROP disimpan di tabel yang sama, jadi **tiketnya satu
> folder**. ⛔ **Tiket sisi komite belum ditulis**; ia akan mulai dari **15**. ⭐ Jangan memakai
> nomor 15 ke atas untuk apa pun yang lain."*
>
> ⭐ **Yang benar:** tiket sisi komite ada di **folder modulnya sendiri** —
> `komite-claim-facin\issues\` — bernomor **00–12**.
> ⭐ **Nomor 15 ke atas di folder ini BEBAS dipakai.**
>
> ⚠️ **Kenapa kalimat lama keliru:** ia mencampur dua hal. ⭐ **Tabelnya** memang dipakai bersama
> *(K6)*, ⛔ **tetapi tiket memerikan perilaku sebuah MODUL**, dan modul komite punya **spec
> sendiri dengan 100 AC sendiri**. ⭐ **Tiket tinggal bersama spec yang ia tutup** — sehingga
> `Menutup: AC …` tidak pernah ambigu.
>
> ⛔ **Rujukan lintas-modul WAJIB berawalan foldernya** — contoh: `komite-claim-facin\issues\00`.
> ⚠️ Tiket `00` ada di **kedua** folder dan **bukan tiket yang sama**.

---

#### Lima belas tiket sisi klaim

| # | Judul | Blocked by |
| --- | --- | --- |
| ⭐ **00** | **PREFACTOR** — dua tabel baru dan kolom induk kedua, **expand–contract** | ⭐ **None** |
| **01** | Registrasi klaim, pengikatan polis, dan daur hidup kasus | 00 |
| **02** | Klasifikasi lini produk — 13 penggolong hidup, 36 tidak dialihkan | 00 · 01 |
| **03** | Estimasi nilai kerugian dan validasinya | 00 · 01 |
| **04** | Pembagian klaim per treaty, dan aturan keseragaman berbagi | 00 · 03 |
| **05** | Penyesuaian nilai klaim, pembagian dan quota share di atasnya | 00 · 03 · 04 |
| **06** | Uang, ketelitiannya, dan batas transaksi | 00 · 05 |
| **07** | Layar klaim | 01 · 03 · 05 |
| ⛔ **08** | **Wewenang** — ditegakkan di lapisan layanan | 01 · 05 |
| **09** | Dokumen akseptasi dan pencetakannya | 05 · 06 |
| **10** | Penyimpanan berkas lampiran | 01 |
| **11** | Efek keluar — kasir, surel, konversi | 06 · 09 |
| **12** | Jejak audit dan kronologi klaim | 01 · 06 |
| **13** | Penyerahan ke komite — **kontrak muatan** | 05 · 06 · 08 |
| **14** | Migrasi data lama, paritas, dan penamaan ulang | 00 · seluruh 01–13 |

---

#### ⭐ Rantai terdalam

```
00 ──▶ 03 ──▶ 04 ──▶ 05 ──▶ 06 ──▶ 09 ──▶ 11 ──▶ 14
```

⭐ **Delapan tingkat.** ⛔ Tiket **14** menunggu seluruhnya, sebab ia menguji **paritas** atas
perilaku yang dibangun tiket-tiket sebelumnya.

#### ⭐ Yang dapat berjalan bersamaan

| Sesudah selesai | Dapat mulai bersamaan |
| --- | --- |
| **00** | 01 · **10** *(berkas lampiran hanya butuh registrasi)* |
| **01** | 02 · 03 · 10 |
| **05** | 06 · 07 · 08 |
| **06** | 09 · 12 |

---

#### ⛔ Satu tiket yang TERTAHAN, dan sebabnya

| Tiket | Penahan | Sifat penahan |
| --- | --- | --- |
| ⛔ **08 · Wewenang** | **butir 6** — ⛔ **isi daftar jabatan dan susunan jenjang belum ada** | ⛔ **MENAHAN PEMBANGUNAN.** ⭐ Tiketnya **tetap ditulis** dan jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum work owner memberikan daftarnya |

⚠️ **Tiket 13 ikut tertahan sebagian** — ⭐ kontrak muatannya dapat ditulis, ⛔ tetapi **jumlah
jenjang** menunggu butir 6 yang sama.

---

#### ⭐⭐ Liputan Acceptance Criteria

`claim-facin\spec.md` memuat **114 AC**, bernomor **1–114**, ⛔ tanpa nomor hilang.

| | Jumlah |
| --- | ---: |
| ⭐ tertutup **tepat satu kali** | ⭐ **109** |
| tertutup **lebih dari satu kali** | ⛔ **0** |
| ⛔ **tidak tertutup** | ⚠️ **5** |
| **TOTAL** | **114** |

⭐ **Dihitung DUA CARA** — *(a)* dari baris `Menutup:` tiap tiket; *(b)* sisiran pola nomor AC pada
seluruh berkas tiket. ✅ **Keduanya sepakat.**

##### ⚠️ Lima AC yang TIDAK tertutup — dan kenapa

⛔ **Bukan kelalaian.** ⭐ Kelimanya berasal dari bab spec berjudul **"Butir yang belum punya sasaran
uji"** — ⭐ spec **sendiri** menyatakan mereka tidak punya sasaran uji.

| AC | Isinya | Kenapa tidak dapat menjadi tiket |
| --- | --- | --- |
| **87** | Dua penunjuk bernama nyaris sama, diisi dari dua sumber berbeda; mana yang mana **tidak terbaca** | ⛔ tidak ada perilaku yang dapat diuji sebelum artinya diketahui |
| **88** | Apakah rule di luar ekspor dapat melempar ke mesin pembangkit tiket **tidak terbukti** | ⛔ buktinya ada **di luar korpus** |
| **89** | Arti sebuah nilai pada tanda tangan parameter **belum diketahui** — ⛔ **dan tidak dipakai sebagai dasar apa pun** | ⭐ tidak menyentuh perilaku |
| **90** | Isi penampung generik bernomor **tidak terbaca** dari rule yang memakainya | ⛔ artinya ditentukan halaman pembawanya, per pemakaian |
| **91** | **327 catatan pengembang belum diuji** dengan membuka rule-nya | ⛔ pekerjaan **pembacaan korpus**, bukan pembangunan |

⭐ **Kelimanya dibawa ke daftar `[terbuka]`, bukan ke tiket.** ⚠️ Bila salah satunya kelak terjawab
dan **ternyata menyentuh perilaku**, ⛔ ia **wajib melahirkan tiket baru** — dan nomornya diambil
dari deret sesudah tiket sisi komite.

---

#### Bukti berkas lain tidak disentuh

| Berkas / folder | Keadaan |
| --- | --- |
| ⛔ `claim-life\` · `komite-claim-life\` · `premiumlist-life\` · `endorsement-life\` | ✅ **NOL disentuh** |
| `claim-facin\spec.md` · `STRUKTUR` · `RELASI` | ✅ md5 tidak berubah — hanya dibaca |
| `komite-claim-facin\` · `claim-prop\` · `komite-claim-prop\` · `docs\adr\` · `CLAUDE.md` | ✅ **NOL disunting** |
| korpus | ✅ md5 tidak berubah |

⛔ Kode **NOL** · DDL **NOL** · `CREATE TABLE` **NOL** · jalur berkas Go **NOL** · nomor baris XML
**NOL** · tiket sisi komite **NOL**.

# Claim Prop

Jumlah tiket: **15**

## Claim Prop - 01 - Registrasi klaim dan validasi nomor polis treaty

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR)
**Menutup:** AC 102 · 103 · 104 · 108 · 109 · 110 · 111 · 112 *(8 AC)* — US 1–4

#### Hasil & nilai pengguna

Claim Admin dapat mendaftarkan klaim baru dengan nomor polis treaty, dan sistem menolak sejak awal
nomor polis yang kosong, salah format, atau tidak ada di master treaty. Klaim ganda pada Date of Loss
yang sama **benar-benar tertahan**, bukan sekadar diperingati — sementara klaim pengganti atas klaim
yang sudah ditolak tetap dapat dibuat.

#### Area codebase

Endpoint registrasi klaim · validasi masuk · lookup master polis · layar registrasi.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CheckNoPolicy.xml` | 2 | pesan format salah; ⚠️ arah gerbang **terbalik** — pesan muncul saat nomor polis **tidak** mengandung penanda format |
| `Activity/CheckNoPolicy.xml` | 4 | menarik data polis lewat `RDBList/GetDataNopolisTreatyin.xml` |
| `RDBList/GetDataNopolisTreatyin.xml` | — | membaca tabel `policyjson`; ⚠️ tanpa prefiks schema |
| `RDBList/SetPolicyTreatyProp.xml` | — | menarik nomor polis, nomor offer, kelompok treaty, sumber bisnis, tahun UW, kuartal; ⚠️ alias berbohong + literal dipalsukan jadi kolom + tabel tanpa prefiks |
| `Activity/SetMasterID.xml` | 3 | pemanggil `SetPolicyTreatyProp` |
| `RDBList/CekPolicyNumber_SQL.xml` | — | mencari id di `pooldata.TREATYINPRODUCTION` menurut nomor polis |
| `Activity/CheeckNoRNM_Act.xml` | 4 | ⚠️ arah **terbalik** — pencarian jalan saat nomor polis **tidak** kosong |
| `Activity/CheeckNoRNM_Act.xml` | 5 | menolak saat hasil pencarian **nol baris** |
| `Activity/CheckNopolicy_Act.xml` | 3 · 4 · 5 | nomor polis kosong ditolak; dua step berikutnya menuliskan syarat yang sama dalam dua bentuk berlawanan |
| `Activity/CheckDateDOL_Act.xml` | 20 | ⚠️ loop riwayat klaim; flag prakondisi **`false`** — pembebasan satu nomor polis yang di-hardcode **mati** |
| `Activity/CheckDateDOL_Act.xml` | 20.1 | menandai duplikat saat DOL sama dan nomor klaim berbeda; menyimpan kunci klaim lama |
| `Activity/CheckReportDate_Act.xml` | 11 · 14 | **BARU 2026-09-19 (ronde 6)** — ⚠️ pemeriksa Report Date terhadap **Start/End Date Policy** keduanya **ber-remark**, dan catatan pengembangnya berbunyi *"matikan protek dalam periode polis"*. **Proteksi ini dimatikan dengan sengaja** |
| `Activity/CheckDateReceived_Act.xml` | 10 · 13 | **BARU** — pola yang sama untuk **Date Received**: kedua pemeriksa terhadap periode polis **ber-remark** |
| `Activity/CheckPeriodPolicy_Act.xml` | 1 | **BARU** — ⚠️ langkah 1 ber-remark padahal ia **satu-satunya pengisi** `Local.CompareDate` dan `Local.EqualDate`; langkah 2 dan 4 tetap membacanya. `[terbuka]` butir 7 spec |
| `Activity/SetEndDate_Act.xml` | 8 | **BARU** — pemanggilan `CheckPeriodPolicy_Act` dari sini **ber-remark**; rule itu tetap dirujuk dari dua Section |
| `Activity/AddAdjustment_Act.xml` · `Section/InputAcceptation.xml` · `Section/OutstandingClaim.xml` | — | **BARU** — ⚠️ properti `.ClaimData.PeriodPolicyTBA` dan pesan *"Policy period is TBA. Please verify the dates."*. **Konsep periode polis "belum pasti" belum pernah tercatat.** `[terbuka]` butir 8 spec |
| `Activity/MakeLowercase_Act.xml` | 1 | **BARU** — lokasi kerugian dan uraian laporan **ditimpa versi huruf kecil**; nilai asli tidak disimpan |

> **Catatan lingkup tambalan ronde 6.** Tambahan di atas **tidak** menyangkut `CheckNoPolicy`.
> Butir *"pencarian polis rangkap"* yang sempat digantung ronde 4 sudah **GUGUR**
> `[data work owner 2026-09-19]`: langkah 3 dan 4 rule itu ber-remark, jadi pencariannya **tidak
> rangkap**. Baris `CheckNoPolicy` di tabel ini **tidak disentuh**.

> **Acceptance criteria tiket ini TIDAK berubah.** Keenam baris di atas adalah **bukti tambahan**
> dan dua butir `[terbuka]` ringan yang **tidak memblokir** — perilaku yang ditiru adalah perilaku
> yang berjalan.
| `Activity/CheckDateDOL_Act.xml` | 21 · 22 | mengambil data klaim lama lewat `ReportDefinition/RejectedClaim_RD` |
| `Activity/CheckDateDOL_Act.xml` | 23.1 | **pengecualian** — klaim lama ada di tabel reject → penanda dikembalikan ke nol |
| `Activity/CheckDateDOL_Act.xml` | 24 · 25 · 26 | peringatan tampil; penanda kesalahan diset `1` — ⚠️ **tidak pernah lolos ambang `>1`** |

#### ADR terkait

**ADR-U-0003** (uang non-float, lewat data polis).

#### Acceptance criteria

- [ ] `[terverifikasi]` Nomor polis treaty wajib berformat yang dikenali; ⚠️ arah gerbangnya terbalik di Pega — aturan sebenarnya **nomor polis harus mengandung penanda format treaty** *(AC 108 spec)*
- [ ] `[terverifikasi]` Data polis ditarik dari master saat klaim didaftarkan — nomor polis, nomor offer, kelompok treaty, sumber bisnis, tahun underwriting, kuartal *(AC 109 spec)*
- [ ] ⚠️ Nama kolom dibuat sesuai isinya dan schema selalu eksplisit. **Alasan menyimpang:** `RDBList/SetPolicyTreatyProp.xml` menamai kelompok treaty sebagai *nama bisnis*, menyisipkan literal sebagai kolom, dan mengeja tabelnya **tanpa prefiks** padahal `RDBList/CekPolicyNumber_SQL.xml` mengeja tabel yang sama **dengan** prefiks *(AC 110 spec)*
- [ ] ⚠️ `[terverifikasi]` Nomor polis yang **tidak ada di master treaty** ditolak dengan pesan yang menyebut sebabnya; penolakan dipicu saat hasil pencarian **nol baris**. **Jebakan:** arah gerbang pencariannya **terbalik** di Pega — pencarian berjalan justru saat nomor polis **tidak** kosong; kedua gerbangnya flag `true`, jadi berlaku *(AC 111 spec)*
- [ ] `[terverifikasi]` Nomor polis **kosong** ditolak terpisah dari nomor polis tidak dikenal; langkah lanjutan hanya berjalan saat nomor polis tidak kosong *(AC 112 spec)*
- [ ] ⚠️ Klaim dengan Date of Loss sama pada polis yang sama **MEMBLOKIR penyimpanan** dan menampilkan peringatan. **Alasan menyimpang:** niat memblokir tertulis di keterangan step 25 (*"make disable submit if there is similar date"*), tetapi step itu memasang nilai `1` sedangkan kedua pembacanya menguji `> 1` — peringatan tampil dan penyimpanan tetap lolos *(AC 102 spec)*
- [ ] `[terverifikasi]` Pengecualian klaim yang sudah ditolak tetap berlaku — tidak ada blokir dan tidak ada peringatan; alur klaim pengganti tidak terganggu *(AC 103 spec)*
- [ ] ⚠️ Pembebasan satu nomor polis yang di-hardcode **tidak dimigrasikan**. **Alasan menyimpang:** literal identitas dunia nyata di dalam kode, sejenis dengan empat nama orang yang dibuang di AC 55; flag prakondisinya `false` sehingga pembebasan itu sudah mati hari ini *(AC 104 spec)*

#### Catatan `[terbuka]` ringan — TIDAK memblokir

⚠️ `[terbuka]` Pembebasan satu nomor polis di `Activity/CheckDateDOL_Act.xml` step **20** adalah
**bekas pengecualian bisnis yang sengaja dipasang lalu dimatikan**, bukan sisa saringan uji coba —
karena itu layak ditanyakan kembali. Masih dikehendaki? Bila ya, jadikan baris data seperti keputusan
AC 55. Pemilik: **work owner**. Bawa apa adanya; **default yang ditulis adalah tidak dimigrasikan**.

#### Perintah verifikasi

```
jalankan test "nomor polis kosong -> ditolak"
jalankan test "nomor polis format salah -> ditolak"
jalankan test "nomor polis tidak ada di master -> ditolak dengan sebab"
jalankan test "DOL sama pada polis sama -> penyimpanan DITOLAK"
jalankan test "DOL sama tetapi klaim lama sudah ditolak -> penyimpanan LOLOS, tanpa peringatan"
cari literal nomor polis di kode                          -> nihil
```

## Claim Prop - 02 - Pemeriksaan kelunasan premi dan pembebasannya lewat data proteksi

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis)
**Menutup:** AC 113 · 114 · 115 · 116 · 117 *(5 AC)* — US 5–6

#### Hasil & nilai pengguna

Klaim tidak diproses di atas polis yang preminya belum dibayar. Claim Admin melihat sebabnya secara
eksplisit, dan kasus yang sudah disetujui manajemen tetap dapat berjalan lewat data proteksi —
tanpa perlu mematikan pemeriksaannya.

#### Area codebase

Validasi akseptasi · integrasi baca ke schema `arasapas` · aturan pembebasan proteksi.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `RDBList/CekLunasPremi_Sql.xml` | — | menjumlahkan nilai bertanda atas `arasapas.invoice` dan `detail_invoice`; ⚠️ pasangan tabelnya **tanpa prefiks schema** |
| `Activity/CekPremiLunas_Act.xml` | 2 | mengisi parameter untuk prefix `CLMP-` |
| `Activity/CekPremiLunas_Act.xml` | 3 | mengisi parameter untuk prefix `CLM-` — ⚠️ **tidak ada cabang untuk `CLMNP-`** |
| `Activity/CekPremiLunas_Act.xml` | 4 | menjalankan pencarian **tanpa gerbang** |
| `Activity/CekPremiLunas_Act.xml` | 5 | ⚠️ menambal koma-ke-titik atas hasil jumlah premi |
| `Activity/CekPremiLunas_Act.xml` | 6 | menandai belum lunas dan menyiapkan pesan bila hasilnya lebih besar dari nol |
| `Activity/CekPremiLunas_Act.xml` | 7 | gerbang induk — pembebasan **hanya dievaluasi ketika premi belum lunas** |
| `Activity/CekPremiLunas_Act.xml` | 7.1 | memanggil `RDBList/CekProteksiKlaim.xml` |
| `Activity/CekPremiLunas_Act.xml` | 7.2 | mengembalikan status menjadi lunas bila proteksi ditemukan |
| `Activity/CekPremiLunas_Act.xml` | 8 | menampilkan pesan bila status tetap belum lunas |
| `RDBList/CekProteksiKlaim.xml` | — | menyaring `pooldata.openproteksi_edm`; ⚠️ dua kode di-hardcode di dalam SQL |
| `Activity/GetDtlPaymentPremi_act.xml` | — | jalur rincian pembayaran premi |

#### ADR terkait

**ADR-U-0003** (uang non-float — nilai premi tidak boleh melewati floating point).

#### Acceptance criteria

- [ ] `[terverifikasi]` Premi polis diperiksa lunas sebelum akseptasi; bila belum lunas, akseptasi ditahan dan pesan sebabnya ditampilkan. ⚠️ **"Lunas" berarti jumlah bertanda ≤ 0, bukan = 0** *(AC 113 spec)*
- [ ] ⚠️ Kedua tabel premi diprefiks schema eksplisit. **Alasan menyimpang:** di Pega `arasapas.invoice` diprefiks sementara pasangannya `detail_invoice` tidak — dalam satu perintah yang sama *(AC 114 spec)*
- [ ] ⚠️ Perilaku pemeriksaan premi untuk prefix `CLMNP-` dibawa apa adanya dan **ditandai**, bukan diam-diam dilengkapi. **Catatan `[terbuka]`:** belum dapat digolongkan — apakah Go menambah cabangnya belum dijawab work owner, jadi apakah ini menyimpang pun belum dapat diputuskan. Yang terbukti: Pega mengisi parameter hanya untuk `CLMP-` dan `CLM-`, sedangkan pencariannya berjalan tanpa gerbang — sehingga untuk `CLMNP-` pencarian berjalan dengan parameter yang tidak pernah diisi *(AC 115 spec)*
- [ ] `[terverifikasi]` Pemeriksaan premi dapat dibebaskan lewat data proteksi; pembebasan **hanya dievaluasi ketika premi belum lunas** *(AC 116 spec)*
- [ ] ⚠️ Dua kode proteksi menjadi **nilai bernama**, bukan literal di dalam perintah SQL. **Alasan menyimpang:** di Pega keduanya dipaku ke dalam teks SQL sehingga tidak dapat diubah tanpa menyentuh rule *(AC 117 spec)*

#### Catatan `[terbuka]` ringan — TIDAK memblokir

⚠️ `[terbuka]` Apakah prefix `CLMNP-` memang dikecualikan dari pemeriksaan premi, atau cabangnya
tertinggal? Pemilik: **work owner**. Bawa apa adanya; perilaku yang ditiru adalah perilaku yang
berjalan.

#### Perintah verifikasi

```
jalankan test "premi belum lunas -> akseptasi ditahan dengan pesan sebabnya"
jalankan test "jumlah bertanda = 0 -> dianggap LUNAS"
jalankan test "jumlah bertanda < 0 -> dianggap LUNAS"
jalankan test "premi belum lunas + proteksi terbuka -> akseptasi LOLOS"
jalankan test "premi sudah lunas -> proteksi TIDAK dibaca"
periksa kedua tabel premi berprefiks schema               -> tidak ada yang telanjang
cari literal kode proteksi di dalam SQL                   -> nihil
```

## Claim Prop - 03 - Cause of Loss dua tingkat dan penandaan katastrofa

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis)
**Menutup:** AC 96 · 97 · 98 · 118 · 119 · 120 · 121 · 122 *(8 AC)* — US 7–9

#### Hasil & nilai pengguna

Claim Admin mengisi penyebab kerugian dari master dua tingkat yang bertaut, sehingga penyebab
tercatat seragam dan dapat dilaporkan. Klaim tidak dapat disimpan tanpa penyebab. Klaim yang berasal
dari satu peristiwa besar dapat ditandai katastrofa dan dikelompokkan di bawah satu event, sehingga
akumulasi kerugian per peristiwa terlihat.

#### Area codebase

Master Cause of Loss (induk, anak, relasi lini bisnis) · entitas event katastrofa · validasi simpan.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `ReportDefinition/BrowseVMCauseOfLoss_RD.xml` | — | tingkat **induk** — id induk, deskripsi, id lama |
| `ReportDefinition/BrowseVDCauseOfLoss_RD.xml` | — | tingkat **anak**, disaring oleh id induk yang sedang dipilih; ⚠️ memuat satu rujukan kelas salah ketik |
| `Activity/CNMInsertCauseOfLoss_act.xml` | — | perawatan master induk; ⚠️ menerima muatan JSON lewat parameter bernama seolah kolom teks |
| `Activity/CNMInsertDetailCauseOfLoss_act.xml` | — | perawatan master anak; ⚠️ pola parameter yang sama |
| `RDBList/GetLBUID_SQL.xml` | — | **tingkat ketiga** — relasi penyebab kerugian ke lini bisnis; ⚠️ kedua tabel **tanpa prefiks schema**, digabung dengan koma gaya lama |
| `Activity/ProteksiData_act.xml` | 14 | gerbang wajib — simpan ditolak saat penyebab kerugian kosong; flag `true`, arah normal |
| `Activity/SaveCatasrtope_Act.xml` | 3 | master event katastrofa disimpan lewat `Obj-Save`; ⚠️ step SQL-nya (step 2) **di-remark** → tidak ditulis |
| `Activity/SetDefNonCatastrope_Act.xml` | 1 · 2 | penanda katastrofa ya/tidak dan jenis non-katastrofa |
| `Activity/InputCatastrope.xml` · `Activity/SetCatastrope_act.xml` · `Activity/SetEditCatastrope.xml` | — | pembuatan dan penyuntingan event |
| `ReportDefinition/GetCatastrope_RD.xml` | — | pemilihan event yang sudah ada |

#### ADR terkait

Tidak ada ADR yang mengikat langsung; keputusan skema mengikuti **tiket 00**.

#### Acceptance criteria

- [ ] `[terverifikasi]` Cause of Loss dipilih dari **master dua tingkat yang bertaut** — tingkat anak disaring oleh id induk yang sedang dipilih *(AC 118 spec)*
- [ ] `[terverifikasi]` Cause of Loss dua tingkat dihubungkan kunci induk, ditambah relasi ke lini bisnis; **wajib** sebelum simpan *(AC 97 spec)*
- [ ] `[terverifikasi]` Tingkat ketiga — relasi Cause of Loss ke lini bisnis — tersedia. ⚠️ Kedua tabelnya diprefiks schema eksplisit dengan `JOIN` yang tertulis; di Pega keduanya telanjang dan digabung koma gaya lama *(AC 119 spec)*
- [ ] ⚠️ Parameter procedure pembaruan Cause of Loss **dinamai sesuai isinya**. **Alasan menyimpang:** kedua procedure Pega menerima muatan JSON lewat parameter bernama seolah kolom teks, sehingga tanda tangannya menyesatkan pembaca *(AC 98 spec)*
- [ ] ⚠️ `[terverifikasi]` Simpan **ditolak** bila Cause of Loss belum diisi. **Jebakan:** rule yang sama menyiapkan **dua pesan kembar** berbunyi sama dan memakai **penanda kesalahan kedua** berambang `>1` — pola yang sama dengan penanda yang dipecah di tiket 00 (AC 105); jangan menyatukan keduanya tanpa membaca kodenya *(AC 121 spec)*
- [ ] `[terverifikasi]` Katastrofa punya **entitas event sendiri** dan satu event **mengelompokkan banyak klaim** *(AC 96 spec)*
- [ ] ⚠️ Penanda katastrofa **ya/tidak** dicatat terpisah dari event-nya, dengan **satu ejaan dan satu enum**. **Alasan menyimpang:** di Pega nama propertinya bahasa Indonesia sedangkan nilainya bahasa Inggris, dan modul mengeja konsep ini dalam **empat bentuk** — `Catastrope` (82 kemunculan) · `Catastrophe` (16) · `Catasrtope` (8) · `Catastrofe` (3 berkas) *(AC 122 spec)*
- [ ] ⚠️ Rujukan kelas salah ketik di `ReportDefinition/BrowseVDCauseOfLoss_RD.xml` **tidak dibawa**; penamaan mengikuti satu bentuk *(AC 120 spec)*

#### Catatan `[terbuka]` ringan — TIDAK memblokir

⚠️ `[terbuka]` `ReportDefinition/BrowseVDCauseOfLoss_RD.xml` memuat rujukan kelas yang kurang satu
huruf, berdampingan dengan ejaan benar di berkas yang sama. Hampir pasti salah ketik lama; perbaikannya
sepele. Pemilik: **pemilik export Pega**.

⚠️ Aturan induk berlaku: `Activity/SaveCatasrtope_Act.xml` step **2** di-remark → **step itu tidak
ditulis sama sekali**; penyimpanan mengikuti jalur step 3.

#### Perintah verifikasi

```
jalankan test "pilih penyebab induk -> daftar anak tersaring oleh induk itu"
jalankan test "simpan tanpa Cause of Loss -> DITOLAK"
jalankan test "dua klaim menunjuk satu event katastrofa -> keduanya terkelompok"
jalankan test "penanda katastrofa tersimpan sebagai enum tertutup"
cari ejaan katastrofa di kode                             -> tepat satu bentuk
cari rujukan kelas salah ketik                            -> nihil
```

## Claim Prop - 04 - Klasifikasi lini bisnis — satu kunci, satu mekanisme

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis)
**Menutup:** AC 71 · 72 · 73 · 74 *(4 AC)*

#### Hasil & nilai pengguna

Klaim terklasifikasi ke lini bisnis yang benar dengan **satu** mekanisme, bukan dua yang saling
bertentangan. Klasifikasi itu menentukan template teks objek pertanggungan yang dilihat Claim Admin,
sehingga tiket **05 (Insured Interest)** bergantung padanya.

#### Area codebase

Pemetaan kelompok treaty ke kelas lini bisnis · template teks objek pertanggungan.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/SetValueToClaim_Act.xml` | 8 | penautan klasifikasi ke kelompok treaty |
| `Activity/SetValueToClaim_Act.xml` | 9 · 10 · 11 · 12 | empat template teks objek pertanggungan menurut kelas |
| `RDBList/GetTreatyGroupID.xml` | — | sumber nilai kelompok treaty |
| `When/` (61 berkas) | — | ⭐ **52 dari 61** menguji properti pada halaman yang **tidak ada** pada objek kerja Claim Prop — **mati di sini, hidup di Claim Fac In**. ⛔ Teks lama: ~~sisa impor model Fac In~~ *(diralat 2026-09-19)* |
| `When/IsCustomBonds.xml` | — | berkas terbaru; **menambah** barisan sisa itu, bukan mengurangi |

#### ADR terkait

Tidak ada ADR yang mengikat langsung.

#### Acceptance criteria

- [ ] `[terverifikasi]` + `[data DBA]` Klasifikasi lini bisnis memakai **satu kunci: kelompok treaty**; tidak ada mekanisme kedua *(AC 71 spec)*
- [ ] ⭐ `[keputusan work owner]` **2026-09-19** Ke-61 rule klasifikasi lini **dipindahkan sekali sebagai satu himpunan**; **hidup-matinya tidak ikut dipindahkan**. Di Claim Prop **52 dari 61** menguji halaman yang tidak ada di sini sehingga **tidak pernah bernilai benar**; di Claim Fac In halaman itu ada dan rule yang sama **hidup**. **Satu salinan, dua nasib.** `[data DBA]` data klaim nyata memuat nama dan kode bisnis, bukan properti yang diuji rule-rule itu. ✅ **Penggolongannya SUDAH diputuskan: penerapan aturan berdiri**, bukan penyimpangan — label `**Alasan menyimpang:**` **dicabut** *(AC 72 spec)*
      > ⛔ **RALAT 2026-09-19** — teks lamanya **dikutip utuh, tidak dihapus**: *"⚠️ Rule klasifikasi warisan **tidak dimigrasikan**. **Alasan menyimpang:** 52 dari 61 rule `When` menguji properti pada halaman yang tidak ada pada objek kerja Claim Prop; dihidupkan pun semuanya bernilai salah. … ⚠️ `[terbuka]` **penggolongannya belum diputuskan** — lihat catatan di bawah; label `**Alasan menyimpang:**` dibiarkan apa adanya sampai dijawab."*
- [ ] ⚠️ **Cacat template dibawa apa adanya, dan itu keputusan sadar.** Marine Cargo memakai template proyek konstruksi yang sama dengan keranjang Aneka, **termasuk typo pada label**. **Catatan paritas:** ini **paritas, bukan penyimpangan** — dicatat di sini agar bila kelak diperbaiki, perbaikannya dicatat sebagai penyimpangan sadar tersendiri *(AC 73 spec)*
- [ ] ⚠️ `[terverifikasi]` Kode pada tabel **kelompok treaty** dan kode bernilai sama pada tabel **jenis treaty** adalah **dua enum berbeda**; test yang menyatukannya **gagal**. **Jebakan:** kode yang sama pernah dibaca sebagai jenis treaty dan ternyata salah — dua tabel kode, dua arti *(AC 74 spec)*

#### ✅ Catatan `[terbuka]` ringan — **SUDAH DIJAWAB 2026-09-19**

> ✅ **`[keputusan work owner]` 2026-09-19 — butir ini DITUTUP.** Penggolongannya **penerapan
> aturan berdiri**. ⭐ **Sebabnya: pertanyaan di bawah terbukti salah pertanyaan.** Selama
> pilihannya dibingkai *"dimigrasikan atau tidak"*, ketiga aturan induk memang tidak menjawabnya.
> Keputusan **A5-6 Claim Fac In** — *rule dipisahkan dari kehidupannya* — membingkainya ulang:
> rule-nya **dipindahkan**; yang tidak dipindahkan adalah **anggapan bahwa ia hidup di sini**.
> ⛔ **Seluruh teks di bawah dibiarkan apa adanya sebagai jejak.**

⚠️ `[terbuka]` **Penggolongan AC 72 belum dapat diputuskan dari aturan berdiri.** Ketiga aturan induk
menjawab tiga hal lain: gerbang yang flag-nya mati (WHEN tidak ditulis, STEP tetap ditulis) · langkah
ber-remark (STEP tidak ditulis) · **elemen UI** selalu-salah (elemen tidak dibuat). **Tidak satu pun
menyebut sebuah rule utuh** yang selalu bernilai salah karena menguji halaman yang tidak ada.

Pertanyaannya: tidak memigrasikan 52 dari 61 rule `When` itu **penerapan aturan berdiri** — sejenis
dengan elemen selalu-salah yang tidak dibuat — atau **penyimpangan tersendiri** yang wajib dihitung?
Pemilik: **work owner**. **Tidak memblokir** — perilaku yang ditiru sama saja dalam kedua bacaan,
karena rule-rule itu bernilai salah di kedua dunia. Yang bergantung padanya hanya **penggolongan**,
bukan hasil kerjanya.

##### ✅ Jalur catatan pengembang — **sudah diperiksa, tertutup** *(ronde 4, 2026-09-19)*

`[terverifikasi]` Ronde 3 menduga ada *"tiga rule `When` yang catatannya menyebut perubahan
penggolongan"* yang perlu diperiksa. Sensus 100% atas `pyMemo` menemukan **empat**, dan
**tidak satu pun mengubah pembacaan tiket ini**:

| Rule (`pxInsName`) | Catatan pengembang | Hasil |
| --- | --- | --- |
| `ASM-FW-GISFW-DATA!ISANEKA` | *"ganti IsBonding jadi IsBondingAndCustomBonds"* | ✅ **sudah dikerjakan** — rujukan yang berlaku `IsBondingAndCustomBonds`; `IsBonding` polos hanya tersisa di dalam teks catatan itu sendiri |
| `ASM-FW-GISFW-WORK!ISEDMADJSHARECEDANT` | *"save as ganti value"* | ⛔ bukan klasifikasi lini bisnis |
| `@BASECLASS!ISPEGAPROD` | *"ganti"* | ⛔ penanda lingkungan, bukan klasifikasi |
| `ASM-FW-GCNMFW-WORK-PNC!ISPA_PNC` | *"tambah pyWorkPage dan hapus policy"* | ⛔ kelas `WORK-PNC` — termasuk ke-52 rule yang **mati di modul ini** *(diralat 2026-09-19; teks lama: ~~sisa impor~~)* |

⛔ **Nol AC berubah.** Dicatat supaya jalur ini tidak diperiksa ulang.

##### ✅ Lima rule `When` terakhir — **sapu bersih selesai**, nol yang hidup *(ronde 6, 2026-09-19)*

`[terverifikasi]` Kelima rule `When` yang belum pernah dibuka telah dibaca satu per satu:

| Rule (`pxInsName`) | Halaman yang diuji | Hidup? |
| --- | --- | --- |
| `ASM-FW-GISFW-DATA!ISMBD` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |
| `ASM-FW-GISFW-DATA!ISMBU` | `pyWorkPage.OfferFacIn.QuotationData.BusinessType` | ⛔ tidak |
| `ASM-FW-GISFW-DATA!ISMARINEHULL` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |
| `@BASECLASS!ISBILLBOARDNEONSYARIAH` | `pyWorkPage.Quotation.BusinessType` + `BusinessCode` | ⛔ tidak |
| `@BASECLASS!ISMAINTENANCE` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |

Kelimanya menguji **halaman yang tidak ada** pada objek kerja Claim Prop — termasuk kelompok 52
rule sisa impor. ⛔ **AC 72 tidak berubah**, dan penggolongannya tetap `[terbuka]` di tangan work
owner.

⚠️ **Jebakan parser yang hampir mengenai:** `isMaintenance` tampak **tanpa kondisi** bila dibaca
dari medan *viewer*, yang berisi pola kosong. Kondisi sebenarnya ada di `pyConditionString`.
Aturan *"isi viewer adalah sisa basi"* berlaku dan terbukti lagi.

#### Perintah verifikasi

```
jalankan test "kelompok treaty properti -> template Property"
jalankan test "kelompok treaty motor -> template Motor Vehicle"
jalankan test "kelompok treaty marine cargo -> template proyek konstruksi (paritas)"
jalankan test "kelompok treaty lain -> keranjang Aneka"
jalankan test "kode kelompok treaty dan kode jenis treaty tidak saling tertukar"
cari rule klasifikasi warisan yang dimigrasikan           -> nihil
```

## Claim Prop - 05 - Insured Interest — objek pertanggungan dan TSI

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis) · 04 (klasifikasi lini bisnis)
**Menutup:** AC 41 · 42 · 43 *(3 AC)* — US 10–13

#### Hasil & nilai pengguna

Claim Admin mendaftar objek pertanggungan beserta nilai TSI-nya, masing-masing dengan mata uang dan
kursnya sendiri, dan melihat total per mata uang beserta totalnya dalam IDR. Dengan itu dasar
perhitungan klaim jelas dan plafonnya terlihat sebelum estimasi disusun.

⚠️ **Baca namanya baik-baik.** "Interest" di modul ini adalah **objek pertanggungan dengan TSI**,
**bukan bunga finansial**. Salah baca di sini akan merusak seluruh perhitungan di hilir.

#### Area codebase

Entitas objek pertanggungan · subtotal TSI per mata uang · template teks per kelas lini bisnis ·
layar Insured Interest.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CountTotalInsterest_Act.xml` | — | menghitung total TSI; **nol rate, nol jumlah hari, nol basis 360/365** — bukan bunga |
| `Activity/SetValueToClaim_Act.xml` | 9 · 10 · 11 · 12 | template teks objek pertanggungan menurut kelas lini bisnis (lihat tiket 04) |
| `Activity/ProteksiData_act.xml` | — | gerbang wajib — daftar objek pertanggungan harus terisi sebelum simpan |
| `Section/InputAcceptation_Est.xml` | — | tampilan daftar objek dan subtotalnya |

#### ADR terkait

**ADR-U-0003** (uang non-float — nilai TSI dan kurs).

#### Acceptance criteria

- [ ] ⚠️ `[terverifikasi]` "Interest" di modul ini adalah **Insured Interest (objek pertanggungan dengan TSI)**, bukan bunga finansial — nol rate, nol jumlah hari, nol basis 360/365 di seluruh rule perhitungannya *(AC 41 spec)*
- [ ] `[terverifikasi]` Daftar objek pertanggungan **wajib terisi** sebelum simpan *(AC 42 spec)*
- [ ] `[terverifikasi]` Total TSI **per mata uang** dan total dalam IDR tersedia *(AC 43 spec)*

#### Perintah verifikasi

```
jalankan test "simpan tanpa objek pertanggungan -> DITOLAK"
jalankan test "dua objek beda mata uang -> subtotal per mata uang benar, total IDR benar"
jalankan test "kurs per objek dipakai, bukan kurs klaim"
cari rate/jumlah hari/basis 360/365 di perhitungan interest -> nihil
```

## Claim Prop - 06 - Loss allocation dan spreading — share wajib tepat 100 %

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis) · 05 (Insured Interest / TSI)
**Menutup:** AC 29 · 30 · 31 · 32 · 33 *(5 AC)* — US 20–23, 30–32

#### Hasil & nilai pengguna

Nilai kerugian terbagi ke beberapa treaty menurut persentase share, disegmentasi per mata uang, dan
**tidak ada nilai yang hilang atau terhitung dua kali** — total share ditolak baik saat kurang dari
100 % maupun lebih. Finance mendapat jaminan bahwa jumlah seluruh baris spreading sama persis dengan
nilai induknya.

#### Area codebase

Entitas loss allocation · spreading tingkat klaim · spreading quota share · validasi total share.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CountSpreading_Act.xml` | 6.2.3 | **satu-satunya** pemeriksaan total share di 329 berkas — hanya menjaga sisi **lebih dari** 100 |
| `Activity/SetTreatyNameSpreading_Act.xml` | 10 | ⚠️ step **di-remark** → tidak ditulis sama sekali |
| `Activity/AddAdjustment_Act.xml` | 5 | loss allocation disalin sebagai **snapshot** ke baris adjustment saat baris dibuat |

⚠️ `[terverifikasi]` Enam properti spreading berbagi satu class di Pega; di Oracle mereka **tabel
berbeda dengan peran berbeda** (lihat tiket 00, AC 3). Jangan disatukan hanya karena class-nya sama.

#### ADR terkait

**ADR-U-0003** (uang non-float) · **ADR-U-0011** (unit keputusan = baris `AdjustmentList`).

#### Acceptance criteria

- [ ] `[terverifikasi]` Loss allocation membagi nilai kerugian ke **treaty**, disegmentasi **per mata uang** — tanpa dimensi tahun maupun coverage *(AC 29 spec)*
- [ ] ⚠️ `[keputusan work owner]` **Total share wajib tepat 100 %**, ditolak pada kurang dari 100 % **dan** lebih dari 100 %. **Alasan menyimpang:** Pega hanya menjaga sebelah — pemeriksaan tunggalnya tidak pernah menangkap kurang dari 100 *(AC 30 spec)*
- [ ] `[terverifikasi]` **Alokasi sisa pembulatan tidak diperlukan** — karena total 100 % dan presisi disimpan penuh, penjumlahan baris spreading tepat sampai digit terakhir *(AC 31 spec)*
- [ ] Spreading di tingkat klaim terpisah dari spreading pada baris adjustment *(AC 32 spec)*
- [ ] Loss allocation disalin sebagai **snapshot** ke baris adjustment saat baris dibuat *(AC 33 spec)*

#### Perintah verifikasi

```
jalankan test "total share 99,99 % -> DITOLAK"
jalankan test "total share 100,01 % -> DITOLAK"
jalankan test "total share tepat 100 % -> diterima"
jalankan test "dua mata uang -> masing-masing disegmentasi dan masing-masing 100 %"
jalankan test "jumlah baris spreading == nilai induk sampai digit terakhir"
jalankan test "snapshot loss allocation pada baris adjustment tidak berubah saat sumbernya berubah"
```

## Claim Prop - 07 - Estimasi klaim per treaty

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 05 (Insured Interest / TSI) · 06 (loss allocation dan spreading)
**Menutup:** AC 34 · 35 · 36 · 37 · 38 · 39 · 40 *(7 AC)* — US 14–19

#### Hasil & nilai pengguna

Claim Admin melihat estimasi terbagi per treaty sesuai penempatan risiko, masing-masing dengan mata
uang dan kursnya sendiri. Ia diperingatkan bila total estimasi melampaui TSI atau plafon cash call,
dan menghapus baris membuat seluruh subtotal ikut menyesuaikan — jadi angkanya tidak pernah
menggantung.

#### Area codebase

Entitas baris estimasi · subtotal estimasi per mata uang · validasi tanggal · peringatan plafon.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/AddEstimation_Act.xml` | — | baris estimasi dibentuk di clipboard; lookup kurs |
| `Activity/CountEstimation_Act.xml` | — | bercabang eksplisit satu-mata-uang versus multi-mata-uang |
| `Activity/CurencyEstimation_Act.xml` | — | mata uang ditetapkan **per baris** |
| `Activity/DeleteEstimation_Act.xml` | — | penghapusan baris dan penyesuaian subtotal |
| `Activity/CheckDateDOL_Act.xml` | — | batas bawah tanggal estimasi = Date of Loss |
| `Activity/AddAdjustment_Act.xml` | 5 | ⚠️ menulis **potret** nilai estimasi ke baris adjustment; **nol rule lain memperbaruinya** |
| `Activity/CountValueADJTreaty_Act.xml` | 16 · 17 | satu-satunya pembaca potret itu |

#### ADR terkait

**ADR-U-0003** (uang non-float).

#### Acceptance criteria

- [ ] `[terverifikasi]` Baris estimasi dibentuk **per treaty**, bersumber dari hasil loss allocation *(AC 34 spec)*
- [ ] `[terverifikasi]` Tanggal estimasi wajib **antara Date of Loss dan hari ini**, inklusif di kedua ujung, zona `Asia/Jakarta` *(AC 35 spec)*
- [ ] Total estimasi melebihi TSI → peringatan *(AC 36 spec)*
- [ ] Total estimasi melebihi plafon cash call → peringatan *(AC 37 spec)*
- [ ] Menghapus baris estimasi menyesuaikan seluruh subtotal *(AC 38 spec)*
- [ ] ⚠️ Dua nama kolom Pega **tidak dibawa apa adanya** ke nama kolom baru. **Alasan menyimpang:** satu kolom bernama persen berisi **nilai uang**, satu kolom bernama jenis kerugian berisi **identitas treaty** — nama yang berbohong membuat kode tidak terbaca *(AC 39 spec)*
- [ ] ⚠️ **Pagar nilai mengikat estimasi TERKINI, bukan potret.** Potret tetap disimpan untuk jejak audit tetapi **tidak** dipakai sebagai pembanding. **Alasan menyimpang:** di Pega nilainya potret yang ditulis sekali saat baris adjustment dibuat dan **nol rule lain memperbaruinya**, sehingga pagar dapat mengikat angka yang sudah basi *(AC 40 spec)*

#### Perintah verifikasi

```
jalankan test "estimasi bertanggal sebelum Date of Loss -> DITOLAK"
jalankan test "estimasi bertanggal besok -> DITOLAK"
jalankan test "estimasi bertanggal tepat Date of Loss -> diterima"
jalankan test "total estimasi > TSI -> peringatan tampil"
jalankan test "total estimasi > plafon cash call -> peringatan tampil"
jalankan test "hapus satu baris -> seluruh subtotal menyesuaikan"
jalankan test "estimasi berubah sesudah baris adjustment dibuat -> pagar memakai nilai TERKINI"
jalankan test "potret estimasi tetap tersimpan di jejak audit"
```

## Claim Prop - 08 - Baris adjustment, adjuster/consultant, dan pagar nilai

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 06 (loss allocation dan spreading) · 07 (estimasi)
**Menutup:** AC 44 · 45 · 46 · 47 · 48 · 49 · 95 · **126** · **127** *(9 AC)* — US 24–29

#### Hasil & nilai pengguna

Claim Admin menambahkan baris adjustment sebagai **unit keputusan** klaim — satu klaim dapat memuat
beberapa keputusan berbeda, masing-masing dengan statusnya sendiri dan data bank penerimanya. Bila
validasi gagal, **tidak ada baris yatim** yang tertinggal. Kelebihan terhadap estimasi ditampilkan
dengan label yang jujur: yang memblokir memang memblokir, yang sekadar peringatan terlihat sebagai
peringatan.

#### ⚠️ Baris beku — MENEGAKKAN, **perilaku BARU**

`[keputusan work owner]` 2026-09-19 — **beku adalah SIFAT TURUNAN sebuah baris penyesuaian:**
ia beku **bila ada kasus komite yang menunjuknya**, **berjalan maupun selesai**. ⛔ **Tidak ada
kolom penanda** — keadaan itu **dihitung**, bukan dibaca dari kolom.

**Yang ditegakkan tiket ini, dua jalur:**

| Jalur | Perilaku |
| --- | --- |
| **setiap jalur sunting** baris penyesuaian | memeriksa ada-tidaknya kasus komite yang menunjuk baris itu, dan **menolak bila ada** — **selamanya** |
| **setiap jalur hapus satu-baris** | idem — **menolak bila ada** |

⚠️ **Penegakannya DI LAPISAN LAYANAN, bukan di layar.** Menyembunyikan tombol saja tidak cukup:
permintaan yang datang langsung ke layanan **tetap harus ditolak**. Sejalan dengan **ADR-U-0014**.

⚠️ **Penolakan komite tidak mencairkannya.** Perbaikan atas isi baris dilakukan dengan **baris
penyesuaian baru**, bukan dengan menyunting baris lama.

⛔ **Hapus KLAIM bukan urusan tiket ini** — itu **tiket 00**. Tiket ini hanya menahan sunting dan
hapus **satu-baris**.

⚠️ `[terverifikasi]` **Pega tidak punya kunci ini** — nol pemeriksaan, nol pesan, nol penanganan.
**Dibangun, bukan dimigrasikan.**

#### Area codebase

Entitas baris adjustment · data bank · adjuster/consultant · pagar nilai terhadap estimasi ·
spreading pada baris adjustment.

#### Rule Pega sumber

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

#### ADR terkait

**ADR-U-0011** (unit keputusan = baris `AdjustmentList`) · **ADR-U-0003** (uang non-float).

#### Acceptance criteria

- [ ] `[terverifikasi]` Unit keputusan adalah **baris adjustment**, masing-masing dengan statusnya sendiri *(AC 44 spec)*
- [ ] `[terverifikasi]` Adjuster dan consultant wajib terisi **hanya sebelum baris adjustment pertama**; baris ke-2 dan seterusnya tidak menuntutnya. Ditiru apa adanya *(AC 45 spec)*
- [ ] ⚠️ Validasi gagal **tidak meninggalkan baris adjustment yatim**. **Alasan menyimpang:** di Pega step keluar mendahului step rollback, sehingga baris kosong tertinggal di klaim *(AC 46 spec)*
- [ ] Baris adjustment menyimpan nama bank, id bank, nomor rekening, dan kode SWIFT *(AC 47 spec)*
- [ ] ⚠️ **RISIKO DITERIMA SADAR — dipertahankan, bukan diperbaiki.** Perhitungan spreading berjalan untuk **semua** payment type; saringan payment type tidak ditulis. **Alasan menyimpang:** konsekuensi langsung aturan induk — gerbang yang flag-nya mati tidak ditulis. Diterima `[keputusan work owner]`, bukan kelalaian migrasi *(AC 48 spec)*
- [ ] ⚠️ Pesan dan penanda blokir memakai **satu besaran yang sama, dalam IDR**, dibandingkan terhadap **estimasi terkini** (tiket 07). **Kelebihan total memblokir penyimpanan; kelebihan per baris menjadi peringatan berlabel jelas.** **Alasan menyimpang:** di Pega step 16 dan step 17 membandingkan besaran berbeda — mata uang asal versus IDR — dan step 16 baris pertama mati, sehingga pengguna melihat teks error tanpa ada yang memblokir. Kelebihan per baris memang tidak boleh memblokir: baris salvage wajib bernilai negatif, sehingga satu baris positif bisa melebihi estimasi sementara hasil bersihnya tidak *(AC 49 spec)*
- [ ] `[terverifikasi]` **Adjuster dan Consultant satu master, dua peran** — master tanpa kolom tipe. Satu klaim maksimal satu adjuster dan satu consultant *(AC 95 spec)*
- [ ] ⚠️ Baris yang **ada kasus komitenya** — berjalan maupun selesai — **tidak dapat diubah, selamanya**; **setiap** jalur sunting menolaknya **di lapisan layanan**, bukan hanya di layar. Test yang berhasil menyunting lewat layanan langsung **gagal**, dan test yang berhasil menyunting **sesudah komite menolak** juga **gagal**. **Alasan menyimpang:** di Pega tidak ada pemeriksaan apa pun — baris boleh berubah sesudah diserahkan *(AC 126 spec)*
- [ ] ⚠️ Baris beku **tidak dapat dihapus satu per satu**; upaya menghapusnya **ditolak** di lapisan layanan. Perbaikan memakai **baris penyesuaian baru**. **Alasan menyimpang:** di Pega baris boleh hilang sesudah diserahkan, tanpa pesan apa pun *(AC 127 spec)*

#### Perintah verifikasi

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

## Claim Prop - 09 - Deductible — satu rumus, satu makna

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 05 (Insured Interest / TSI) · 08 (baris adjustment)
**Menutup:** AC 50 · 51 · 52 · 53 *(4 AC)* — US 33–35

#### Hasil & nilai pengguna

Deductible dihitung dengan **satu** rumus yang disepakati, atas basis yang dipilih eksplisit, dan
dikurangkan pada titik yang sama setiap kali. Claim Admin berhenti menebak versi mana yang sedang
berlaku.

⚠️ Korpus memuat **dua rumus deductible yang bertentangan**. Tiket ini menutup pertentangan itu
dengan memilih versi 2024.

#### Area codebase

Perhitungan deductible · pemilihan basis · urutan pengurangan terhadap share ceding.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CountValueADJTreaty_Act.xml` | 6 | penyelarasan mata uang sebelum deductible dari TSI |
| `Activity/CountValueADJTreaty_Act.xml` | 4 · 5 · 7 · 8 · 9 | tiga jenis basis risiko individual dan cabang per tipe |

⚠️ `[terverifikasi]` Properti bernama sama berarti **berbeda** di dua rule lama — sumber langsung
pertentangan rumus.

#### ADR terkait

**ADR-U-0003** (uang non-float).

#### Acceptance criteria

- [ ] `[keputusan work owner]` Deductible = **MAX(persentase × basis, nilai flat)** *(AC 50 spec)*
- [ ] `[terverifikasi]` Basis dipilih antara **TSI** dan **nilai klaim** *(AC 51 spec)*
- [ ] `[keputusan work owner]` Rumus yang berlaku adalah versi **2024** — deductible dikurangkan **sesudah** share ceding diterapkan. Versi 2022 **ditinggalkan** *(AC 52 spec)*
- [ ] ⚠️ Properti bernama sama diberi **satu makna saja**. **Alasan menyimpang:** di dua rule lama nama yang sama berarti dua hal berbeda, sehingga pembaca tidak dapat tahu rumus mana yang sedang berjalan *(AC 53 spec)*

#### Perintah verifikasi

```
jalankan test "persentase x basis > flat -> deductible = hasil persentase"
jalankan test "flat > persentase x basis -> deductible = flat"
jalankan test "basis TSI dan basis nilai klaim menghasilkan angka berbeda sesuai pilihan"
jalankan test "deductible dikurangkan SESUDAH share ceding (versi 2024)"
cari implementasi rumus deductible kedua                  -> nihil
```

## Claim Prop - 10 - Wewenang — satu sumber, ditegakkan di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR)
**Menutup:** AC 54 · 55 · 56 · 57 · 58 · 59 · 60 · 61 *(8 AC)* — US 55–58

#### Hasil & nilai pengguna

Tingkat wewenang setiap pelaku dibaca dari **satu** sumber data, sehingga tidak ada dua daftar yang
berbeda. Pergantian pemegang jabatan cukup dengan mengubah baris tabel — tanpa menyentuh kode. Dan
yang paling penting: gerbang wewenang **benar-benar menolak**, bukan sekadar menyembunyikan tombol.

⚠️⚠️ `[terverifikasi]` **Pega tidak punya satu pun gerbang wewenang blocking di sisi server:** nol
`pyValidateActivity` di 12 dari 12 FlowAction; nol pemeriksaan peran di Activity, 418 ekspresi gerbang
di 36 dari 36 Section, 19 ReportDefinition, dan 11 Harness. Gerbang penyerahan komite sepenuhnya
*client-side*. **Membangunnya di Go bukan penyimpangan — ini mengisi lubang yang memang tidak pernah
ditutup.**

#### Area codebase

Roster komite · resolusi tingkat wewenang · penegakan di lapisan layanan · batas nilai Direktur Utama.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `DataTransform/InsertChronology_DT.xml` | 1.1.4 | tingkat dasar diset **tanpa syarat** lebih dulu |
| `DataTransform/InsertChronology_DT.xml` | 1.1.5 · 1.1.6 · 1.1.7 | ⚠️ tiga cabang bernama orang — **nilai dasar + tiga tingkat komite**, bukan empat tingkat sejajar |
| `ReportDefinition/FilterEmailKomiteWithLimit.xml` | — | roster; penyaringnya **tiga** — batas bawah pita, status klaim, status aktif |
| `RDBList/GetLimitDirekturUtama_SQL.xml` | — | ⚠️ mencari menurut **nama orang** yang di-hardcode |
| `Activity/AttachmentProtect_ACT.xml` | 5 | batas nilai Direktur Utama menggerbangi **kewajiban dokumen**, bukan tangga persetujuan |
| `Activity/AddKomiteTreatyChild_ACT.xml` · `Activity/SetKomiteTreaty_ACT.xml` | — | pemakai roster |

#### ADR terkait

**ADR-U-0014** (keputusan komite ditegakkan di lapisan layanan — ⚠️ cakupan asli ADR-nya Komite Claim
Life; yang dipakai di sini **prinsipnya**).

#### Acceptance criteria

- [ ] ⚠️ `[keputusan work owner]` Tingkat wewenang dibaca dari **satu sumber**: roster komite. Tidak ketemu → tingkat dasar Claim Admin. **Alasan menyimpang:** di Pega tingkat dasar diset tanpa syarat lalu ditimpa tiga cabang bernama orang *(AC 54 spec)*
- [ ] ⚠️ `[keputusan work owner]` **Empat nama orang yang di-hardcode dibuang** — tiga di transform jejak audit dan satu di rule batas Direktur Utama. Perilaku tidak berubah, tetapi pergantian pemegang jabatan cukup lewat baris tabel. **Alasan menyimpang:** identitas dunia nyata dipaku ke dalam kode; orangnya pindah jabatan → jejak audit salah label tanpa peringatan *(AC 55 spec)*
- [ ] `[keputusan work owner]` Batas nilai **Direktur Utama** adalah wewenang **terpisah**, bukan tingkat kelima tangga komite *(AC 56 spec)*
- [ ] ⚠️ Seluruh gerbang wewenang **ditegakkan di lapisan layanan**; test yang menemukan gerbang hanya di UI **gagal**. **Catatan pengisian lubang:** ini **bukan penyimpangan** — spec menyatakannya terang di Implementation Decision 7 (*"Ini bukan penyimpangan — melainkan mengisi lubang yang memang tidak pernah ditutup"*). Pega tidak punya satu pun gerbang blocking di sisi server, jadi tidak ada perilaku yang disimpangi *(AC 57 spec)*
- [ ] Test: panggil endpoint penyerahan komite langsung, tanpa UI → **ditolak** bila syarat tidak terpenuhi *(AC 58 spec)*
- [ ] `[terverifikasi]` Jumlah tingkat komite = **COUNT baris roster aktif** yang pita limitnya mencakup nilai klaim, dihitung saat penyerahan — **tidak** dari konstanta *(AC 59 spec)*
- [ ] Batas **atas** pita roster **tidak** dijadikan penyaring; penyaring tetap batas bawah + status klaim + status aktif. **Catatan paritas:** bila batas atas ikut menyaring, hanya satu pita lolos dan tangga berjenjang hilang *(AC 60 spec)*
- [ ] ✅ `[data DBA]` Kunci pencocokan pelaku ke roster = **kolom operator id pada roster**, dikonfirmasi sama dengan user id Pega. Kolom id pengguna baru **batal**; kolom email tidak dipakai sebagai kunci. **Nol perubahan skema** untuk bagian ini *(AC 61 spec)*

#### Catatan `[terbuka]` ringan — TIDAK memblokir

⚠️ `[terbuka]` Roster memuat **enam kolom status per jenis aksi** — adjustment, adjuster, registrasi,
reject, salvage, survey — tetapi **nol** dipakai menyaring oleh rule mana pun (sensus 329 berkas).
Apakah wewenang per jenis aksi memang dimaksudkan berlaku, atau keenam kolom itu warisan yang tidak
terpakai? Pemilik: **Finance + work owner**. Bawa apa adanya; yang ditiru adalah tiga penyaring yang
berjalan.

⚠️ **Jebakan nama:** salah satu nama kolom itu juga muncul sebagai properti muatan layanan konversi
non-life. Tabrakan nama, bukan pemakaian — jangan salah simpulkan kolom roster itu terpakai.

#### Perintah verifikasi

```
jalankan test "pelaku tidak ada di roster -> tingkat dasar Claim Admin"
jalankan test "pelaku ada di roster -> tingkat dari kolom tingkat roster"
jalankan test "panggil endpoint kirim-komite langsung tanpa hak -> DITOLAK di layanan"
jalankan test "jumlah tingkat komite = jumlah baris roster yang lolos, bukan konstanta"
jalankan test "batas Direktur Utama tidak menambah tingkat tangga komite"
cari nama orang yang di-hardcode di kode dan SQL          -> nihil
```

## Claim Prop - 11 - Penyerahan ke Komite dan penutupan klaim

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 08 (baris adjustment) · 10 (wewenang)
**Menutup:** AC 62 · 63 · 64 · 65 · 66 · 67 · 68 · 69 · 70 · **125** *(10 AC)* — US 36–44

#### Hasil & nilai pengguna

Claim Admin menyerahkan baris adjustment ke Komite untuk diputuskan, dan penyerahan ditolak bila data
bank atau lampiran wajib belum lengkap — sehingga Komite tidak pernah memutus di atas data setengah
jadi. Klaim yang tidak dibayar tetap punya jejak keputusan lewat jalur penutupan tanpa pembayaran, dan
klaim tidak dapat ditutup selagi masih ada keputusan yang menggantung.

⚠️ **Batas konteks.** Tiket ini membangun **batas** menuju Komite Claim Prop, bukan isinya. Tangga
persetujuan, routing per tingkat, dan penomoran akseptasi milik konteks sebelah.

##### ⚠️ Penyerahan MEMBEKUKAN baris penyesuaian induknya — **perilaku BARU**

`[keputusan work owner]` 2026-09-19 — **penyerahan ke komite melahirkan kasus komite, dan
keberadaan kasus itulah yang membekukan baris penyesuaian yang diserahkan.** Sejak saat itu baris
tersebut **tidak dapat diubah, selamanya**, dan **tidak dapat dihapus satu per satu**. **Penolakan
komite tidak mencairkannya**; perbaikan memakai **baris penyesuaian baru**.

⛔ **NOL PENANDA DIPASANG — tidak ada yang perlu ditulis.** Bekunya **diturunkan**, bukan disimpan:
sebuah baris beku **bila ada kasus komite yang menunjuknya**, berjalan maupun selesai. ⛔ **Jangan
membuat kolom penanda beku**; acuan tunggal nama kolom tetap `STRUKTUR-TABEL-CLAIM-PROP.md`.

⚠️ **Penegakannya bukan di sini.** Tiket ini hanya **melahirkan kasusnya**; yang menolak sunting
dan hapus satu-baris adalah **tiket 08**, dan yang menolak hapus klaim adalah **tiket 00**.

⚠️ `[terverifikasi]` **Pega tidak punya kunci ini** — nol pemeriksaan, nol pesan kesalahan, nol
penanganan untuk baris penyesuaian yang berubah atau hilang sesudah diserahkan. **Dibangun, bukan
dimigrasikan.** Lingkup: **Claim Prop** dan **Claim Non Prop** *(Non Prop menyusul)*; **Claim —
Life tidak termasuk**.

#### Area codebase

Penyerahan kasus komite · penjaga kelengkapan · jalur penutupan tanpa pembayaran · penutupan klaim.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/AddKomiteTreatyChild_ACT.xml` | — | membuat kasus anak komite; ⚠️ penunjuknya **indeks posisi**, muatan digemukkan 24 field sebagai snapshot |
| `Activity/SetKomiteTreaty_ACT.xml` | — | roster kasus komite |
| `Activity/AttachmentProtect_ACT.xml` | 5 | lampiran wajib bergantung jenis pembayaran dan batas nilai |
| `Activity/CloseClaimProp.xml` | 5 | ⚠️ memo penutupan ditulis ke slot yang **tidak pernah dibaca** (ditangani tiket 14) |
| `Activity/SaveOutstanding_Act.xml` | 39 | ⚠️ data dikirim ke sistem luar **tanpa menunggu komite** |
| `ConnectREST/SendAcceptationToKasir.xml` | — | ⭐ **BARU 2026-09-19 (ronde 4)** — pengiriman ke Kasir **berautentikasi**; bentuknya **profil autentikasi bernama**. AC 69 (*penutupan ditolak bila pengiriman Kasir belum berhasil*) karenanya bergantung pada kredensial yang **harus tersedia saat runtime** — rinciannya di **tiket 13** |
| `Activity/SetKomiteNo_Act.xml` | 1 | ⭐ **BARU 2026-09-19 (ronde 6)** — ⚠️ **nomor komite tidak disimpan, melainkan DIIRIS**: karakter ke-19 sampai ke-30 dari **kunci internal Pega** kasus anak **terakhir**. Dua kerapuhan sekaligus — bergantung **panjang/format kunci internal**, dan bergantung **posisi terakhir** dalam daftar. Menguatkan AC 63 (*rujukan memakai ID stabil, bukan indeks posisi*) |

⚠️ `[terverifikasi]` **Nol rule pembatal di 329 berkas.**

#### ADR terkait

**ADR-U-0011** (unit keputusan = baris `AdjustmentList`) · **ADR-U-0014** (penegakan di lapisan layanan).

#### Acceptance criteria

- [ ] `[keputusan work owner]` Penyerahan membuat rekam kasus komite yang membawa **penunjuk stabil** ke baris adjustment *(AC 62 spec)*
- [ ] ⚠️ Rujukan memakai **ID stabil**, bukan indeks posisi. **Alasan menyimpang:** di Pega penunjuknya indeks posisi, sehingga muatan harus digemukkan 24 field sebagai snapshot supaya rujukan tidak tergeser *(AC 63 spec)*
- [ ] `[terverifikasi]` Penyerahan **ditolak** bila data bank belum lengkap *(AC 64 spec)*
- [ ] `[terverifikasi]` Penyerahan **ditolak** bila lampiran wajib belum lengkap; daftar lampiran wajib **bergantung jenis pembayaran** *(AC 65 spec)*
- [ ] ⚠️ **Dua jalur ke Komite** — penyerahan adjustment dan penutupan tanpa pembayaran — diberi **saling-kunci**. **Alasan menyimpang:** di Pega keduanya menghasilkan kasus anak berkelas sama dan **keduanya dapat aktif pada klaim yang sama**, tanpa apa pun yang mencegahnya *(AC 66 spec)*
- [ ] ⚠️ Jalur penutupan tanpa pembayaran mengambil roster dari **data**. **Alasan menyimpang:** di Pega jalur itu memakai roster tertanam satu orang *(AC 67 spec)*
- [ ] `[terverifikasi]` Penutupan **ditolak** bila masih ada baris adjustment yang belum diputus *(AC 68 spec)*
- [ ] `[terverifikasi]` Penutupan **ditolak** bila pengiriman ke Kasir belum berhasil *(AC 69 spec)*
- [ ] ⚠️ **RISIKO DITERIMA SADAR — dipertahankan, bukan diperbaiki.** Klaim yang **ditolak** komite **tetap** memiliki data di sistem luar (Arasapas); pengiriman terjadi tanpa menunggu komite dan **nol rule pembatal** ada di 329 berkas. **Alasan menyimpang:** paritas atas keputusan `[keputusan work owner]`; rekonsiliasinya urusan manual di luar lingkup modul ini *(AC 70 spec)*
- [ ] ⚠️ Sesudah penyerahan tercatat, baris penyesuaian yang diserahkan **beku** — dan bekunya **diturunkan dari adanya kasus komite yang menunjuknya**, bukan dari kolom penanda. ⛔ Test yang mencari **kolom penanda beku gagal**; test yang menemukan baris terserahkan **tanpa kasus komite** juga **gagal**. **Alasan menyimpang:** di Pega baris boleh berubah atau hilang sesudah diserahkan, tanpa satu pun pemeriksaan *(AC 125 spec)*

#### Perintah verifikasi

```
jalankan test "kirim ke komite tanpa data bank -> DITOLAK"
jalankan test "kirim ke komite tanpa lampiran wajib -> DITOLAK"
jalankan test "lampiran wajib berbeda menurut jenis pembayaran"
jalankan test "jalur adjustment aktif -> jalur tutup-tanpa-bayar DITOLAK (saling-kunci)"
jalankan test "tutup klaim dengan baris adjustment belum diputus -> DITOLAK"
jalankan test "tutup klaim sebelum kirim Kasir berhasil -> DITOLAK"
jalankan test "rujukan baris adjustment memakai ID stabil, tidak tergeser saat baris lain dihapus"
jalankan test "klaim ditolak komite -> data di sistem luar TETAP ADA (risiko diterima sadar)"
```

## Claim Prop - 12 - Dokumen — PLA, DLA, dan Acceptance Note

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 08 (baris adjustment)
**Menutup:** AC 84 · 85 · 86 · 87 · 88 *(5 AC)* — US 45–49

#### Hasil & nilai pengguna

Claim Admin mencetak PLA pada tahap outstanding, DLA per reinsurer per baris adjustment, dan
Acceptance Note otomatis saat akseptasi disimpan. Semua dokumen tersimpan dan dapat dibuka kembali
lewat tautan, sehingga tidak perlu mencetak ulang — dan nomor dokumen **tidak terbakar** saat
pembuatan gagal (dijamin tiket 00).

#### Area codebase

Pembuatan PDF · penyimpanan berkas dan metadata · tautan berbatas waktu · gerbang anti-cetak-ganda.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/TryMakePLA_Act.xml` | 8 · 11 · 14 · 26 | PLA; nomor diambil lalu hanya ke clipboard, step lanjut dapat melempar exception, **nol `Obj-Save`** |
| `Activity/PrintDLATreatyIn.xml` | 7 | DLA; menarik nomor lewat generator DLA |
| `Activity/PrintDLATreatyIn.xml` | 15.11.1 | cabang yang **hidup** — diikuti apa adanya, termasuk pola halaman sementara berindeks tetap |
| `Activity/PrintDLATreatyIn.xml` | 15.11.2 · 15.11.3 | ⚠️ cabang salvage dan adjuster fee **di-remark** → **tidak dimigrasikan** |
| `Activity/PrintFileAcceptance.xml` | — | Acceptance Note tercetak saat akseptasi disimpan |
| `RDBList/GetTokenStorage_SQL.xml` · `RDBList/Insert_T_Storage_SQL.xml` | — | token → unggah → simpan metadata |
| `Activity/CloseClaimProp.xml` | — | ⚠️ **tidak** memeriksa nomor PLA maupun DLA — pencetakan bukan gerbang |

#### ADR terkait

**ADR-U-0010** (penyimpanan berkas tetap di Google Storage) · **ADR-U-0006** (penomoran lewat stored
procedure).

#### Acceptance criteria

- [ ] `[keputusan work owner]` **DLA = Definite Loss Advise** — satu PDF **per reinsurer per baris adjustment**, dengan saudara retro *(AC 84 spec)*
- [ ] `[terverifikasi]` Ketiga dokumen **disimpan** ke penyimpanan berkas dengan metadata di database; URL berbatas waktu dengan jalur penyegaran *(AC 85 spec)*
- [ ] `[terverifikasi]` **Pencetakan bukan gerbang** — penutupan klaim tidak memeriksa nomor dokumen. Ditiru apa adanya *(AC 86 spec)*
- [ ] `[terverifikasi]` Gerbang **anti-cetak-ganda** dipertahankan *(AC 87 spec)*
- [ ] ⚠️ **RISIKO DITERIMA SADAR — dipertahankan, bukan diperbaiki.** Cabang **salvage** (step 15.11.2) dan **adjuster fee** (step 15.11.3) pada DLA di-remark → **tidak dimigrasikan**, tanpa pengecualian; **step 15.11.1 diikuti apa adanya**, termasuk pola halaman sementara berindeks tetap. **Konsekuensinya diterima:** DLA untuk salvage dan adjuster fee menampilkan **nilai klaim**, sama seperti produksi hari ini. **Alasan menyimpang:** aturan induk — langkah ber-remark tidak ditulis sama sekali *(AC 88 spec)*

#### Perintah verifikasi

```
jalankan test "cetak PLA -> berkas tersimpan, metadata tercatat, tautan dapat dibuka"
jalankan test "cetak DLA -> satu PDF per reinsurer per baris adjustment"
jalankan test "simpan akseptasi -> Acceptance Note tercetak otomatis"
jalankan test "pembuatan PDF gagal -> nomor dokumen TIDAK terbakar"
jalankan test "cetak dua kali -> gerbang anti-cetak-ganda menahan"
jalankan test "tutup klaim tanpa cetak PLA/DLA -> LOLOS (paritas)"
jalankan test "DLA salvage menampilkan nilai klaim (risiko diterima sadar)"
```

## Claim Prop - 13 - Efek keluar — Kasir, Arasapas, konversi non-life, dan email komite

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 08 (baris adjustment)
**Menutup:** AC 89 · 90 · 91 · 92 · 93 · 94 · 99 · 100 · 101 *(9 AC)* — US 50–54

#### Hasil & nilai pengguna

Data akseptasi sampai ke Kasir, klaim tersinkron ke sistem inti reinsurance non-life, dan email
pemberitahuan sampai ke anggota komite yang berwenang — **dan bila salah satunya gagal, pekerjaan itu
tidak hilang diam-diam**. Claim Admin melihat status pengiriman, dan Finance mendapat nomor kasus
Kasir untuk rekonsiliasi.

⚠️⚠️ `[terverifikasi]` Hari ini **tidak ada jaring pengaman sama sekali**: nol `Exit-Activity`, nol
`Page-Set-Messages`, nol rollback, dan **39 dari 39** penangan exception kosong — termasuk pada kedua
langkah pemanggilan REST. Tabel log yang ada adalah **jejak, bukan antrean**: nol kolom status/retry,
satu INSERT dan satu SELECT di seluruh korpus, nol UPDATE, nol job pemroses ulang.

#### Area codebase

Outbox transaksional · klien Kasir · klien konversi non-life · pengirim email · resolusi endpoint.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/HitServiceToKasir_Act.xml` | — | 14 langkah tanpa penanganan kegagalan; nol retry |
| `Activity/HitServiceToKasir_Act.xml` | 9.8 · 13.5 | mencatat ke tabel log **empat kolom**: muatan JSON, kunci internal Pega, nomor akseptasi, dan pesan respons — ⚠️ **nomor kasus Kasir dibuang** |
| `ConnectREST/SendAcceptationToKasir.xml` | — | memetakan seluruh JSON ke satu properti tanpa menyebut sub-properti |
| `ConnectREST/SendAcceptationToKasir.xml` | pengaturan Connect-REST | ⭐ **BARU 2026-09-19 (ronde 4)** — **satu-satunya** rule berautentikasi di 329 berkas: `pyUseAuthentication = true`, dan kredensialnya diambil dari **profil autentikasi bernama** (`pyAuthProfileSelectionType`), **bukan** header yang ditulis di badan rule. Proksi `NO_AUTH`. Catatan pengembangnya: *"tambah auth"* |
| kelima `ConnectREST/` lain di Claim Prop | pengaturan Connect-REST | ⭐ **BARU** — `GETDTLPAYMENTCLAIM` · `KONVERSIKLAIMNONLIFE` · `SERVICEGOOGLE` · `GETPAYATTACHMENT` · `GETPREMIUMPAIDONTREATYIN` semuanya `pyUseAuthentication = **false**` — ⛔ **jangan** diberi autentikasi |
| `ConnectREST/ServiceGoogle.xml` | pengaturan Connect-REST | ⭐ **BARU 2026-09-19 (ronde 6)** — **batas waktu respons `300 000` milidetik (5 menit)**, catatan pengembang *"buat jadi 300rb"*. Satu-satunya batas waktu eksplisit yang terbaca di keenam Connect-REST |
| `Activity/SetPayableTreaty_Act.xml` | 4 · 5 | ⭐ **BARU** — ⚠️ `Obj-Save` dan `Obj-Refresh-And-Lock` **ber-remark**, jadi rule ini **tidak pernah menyimpan**. Catatan pengembangnya *"BUKA PROTEKSI SALVAGE KASIR"*. Termasuk daftar **jalur simpan mati** `[data work owner 2026-09-19]` |
| `ConnectREST/KonversiKlaimNonLife.xml` | — | sinkronisasi ke sistem inti non-life |
| `Activity/KonversiKlaim_Act.xml` | — | pemanggil konversi |
| `Activity/SendEmailKlaim.xml` · `Activity/SendEmailKlaimRejectClose.xml` | — | pemberitahuan email |
| `RDBList/GetEmailCeding_SQL.xml` | — | alamat ceding lewat fungsi di schema `gl` |
| keenam `ConnectREST/` | — | endpoint **tidak pernah literal**; berasal dari tabel tautan layanan dikunci sepasang kategori |

⚠️ `[terverifikasi]` Kelima *production level* Pega bernilai identik — pemisahan lingkungan sepenuhnya
bergantung isi tabel, dan satu cabang digerbangi **nama node aplikasi**.

#### ADR terkait

**ADR-U-0015** (efek keluar wajib berhasil — transactional outbox) · **ADR-U-0013** (resolusi endpoint
lewat tabel tautan layanan) · **ADR-U-0010** (penyimpanan berkas).

⚠️ **Catatan keamanan — dicatat, BUKAN untuk ditindaklanjuti oleh tiket ini.** Kredensial autentikasi
Kasir ikut beredar di dalam berkas ekspor korpus. Ia **sebaiknya diganti sesudah migrasi**. Itu
urusan **tim Kasir**; dicatat di sini semata supaya tidak terlewat. ⛔ Nilainya tidak pernah disalin
ke berkas mana pun di repositori ini.

#### Acceptance criteria

- [ ] ⚠️ Efek keluar memakai **outbox transaksional** — dicatat dalam transaksi yang sama dengan perubahan data, dikirim, dan **diulang** bila gagal. **Alasan menyimpang:** Pega tidak punya retry, tidak punya rollback, dan 39 dari 39 penangan exception-nya kosong; kegagalan hilang tanpa jejak *(AC 89 spec)*
- [ ] Test: gagalkan klien luar → pekerjaan **tetap ada di outbox** dan terkirim pada percobaan berikutnya *(AC 90 spec)*
- [ ] `[terverifikasi]` Endpoint **di-resolve saat runtime** dari tabel tautan layanan; **nol URL literal** di kode *(AC 91 spec)*
- [ ] ⚠️ Pemisahan lingkungan **tidak** bergantung pada identitas node aplikasi. **Alasan menyimpang:** di Pega satu cabang digerbangi nama node, sehingga perilaku berubah bila node diganti nama atau ditambah *(AC 92 spec)*
- [ ] ✅ `[data DBA]` **Tidak ada bug pada pembacaan kode respons Kasir.** Kasir mengirim **kedua ejaan sekaligus dengan nilai identik**, sehingga teks status sukses memang muncul selama ini. Label cacat dicabut *(AC 93 spec)*
- [ ] Status pengiriman ke Kasir dapat dilihat pengguna *(AC 94 spec)*
- [ ] ⭐ **BARU 2026-09-19 (ronde 4)** — `[terverifikasi]` Klien Kasir **berautentikasi**, dan kredensialnya berasal dari **konfigurasi runtime bernama** — bukan dipaku di kode, bukan di berkas sumber. Bentuk di Pega: profil autentikasi bernama yang dirujuk rule Connect-REST *(⛔ belum ada nomor AC spec — `spec.md` belum ditambal)*
- [ ] ⭐ **BARU** — `[terverifikasi]` **Hanya klien Kasir** yang berautentikasi. Kelima klien luar lain **tidak**; memberi mereka autentikasi adalah penyimpangan yang tidak diminta
- [ ] Test: kredensial Kasir salah/kosong → pengiriman **gagal terang-terangan** dan pekerjaan **tetap di outbox** (bukan gagal diam-diam)
- [ ] ⚠️ `[data DBA]` **Nomor kasus Kasir disimpan** bersama jejak pengiriman. **Alasan menyimpang:** Pega membuangnya — sensus 329 berkas menemukannya di **nol berkas**, sehingga rantai rekonsiliasi dengan Kasir putus *(AC 99 spec)*
- [ ] ⚠️ `[data DBA]` Kode respons Kasir **dibandingkan sebagai string, atau dinormalisasi eksplisit**. **Alasan menyimpang:** Kasir mengirim nilai bertipe string sedangkan Pega membandingkannya sebagai angka dan bergantung pada konversi diam-diam; Go tidak melakukan konversi itu *(AC 100 spec)*
- [ ] `[data DBA]` Parser menerima **kedua ejaan** kode dan pesan respons, **dan mencatat di log** bila yang datang hanya ejaan yang tidak diharapkan *(AC 101 spec)*

#### Perintah verifikasi

```
jalankan test "klien Kasir gagal -> pekerjaan tetap di outbox"
jalankan test "outbox diproses ulang -> pekerjaan terkirim"
jalankan test "nomor kasus Kasir tersimpan bersama jejak pengiriman"
jalankan test "kode respons string '1' -> sukses"
jalankan test "kode respons ' 1' dan '01' -> tidak lolos diam-diam sebagai sukses"
jalankan test "hanya ejaan tak diharapkan yang datang -> tercatat di log"
cari URL literal di kode                                  -> nihil
cari ketergantungan pada nama node aplikasi               -> nihil
```

## Claim Prop - 14 - Jejak audit klaim — empat cacat yang diperbaiki

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 10 (wewenang)
**Menutup:** AC 77 · 78 · 79 · 80 · 81 · 82 · 83 *(7 AC)* — US 59–62

#### Hasil & nilai pengguna

Auditor dapat merekonstruksi klaim dari awal sampai akhir: setiap aksi meninggalkan entri berisi
aksi, pelaku, waktu, dan tingkat wewenang; riwayat tampil urut waktu; **semua** peran terekam tanpa
kecuali; dan memo penutupan klaim benar-benar sampai ke jejak.

⚠️ Keempat cacat di bawah **membalik** perilaku Pega. Itu disengaja dan disetujui — tanpa perbaikan
ini jejak audit tidak dapat dipercaya sebagai bukti.

#### Area codebase

Entitas jejak audit · enum tingkat wewenang · migrasi entri warisan · pengurutan riwayat.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `DataTransform/InsertChronology_DT.xml` | 1 | ⚠️ digerbangi jabatan pelaku, dan **seluruh penulisan bersarang di bawahnya** → aksi peran itu **nol entri** |
| `DataTransform/InsertChronology_DT.xml` | 1.1.1 | membaca slot memo pada halaman kronologi |
| `DataTransform/InsertChronology_DT.xml` | 1.1.4–1.1.7 | tingkat wewenang ditulis ke properti berawalan `Is…` yang **bukan boolean** |
| `Activity/SethistoryKlaimTreaty.xml` | — | penulis kedua; **tidak** punya pengecualian jabatan, dan mengeja tingkat terbawah **berbeda** |
| `Activity/CloseClaimProp.xml` | 5 | ⚠️ memo penutupan ditulis ke slot yang **tidak pernah dibaca** |
| `Activity/MakeLowercase_Act.xml` | 1 | ⭐ **BARU 2026-09-19 (ronde 6)** — ⚠️ **lokasi kerugian dan uraian laporan ditimpa versi huruf kecilnya**; **nilai asli tidak disimpan di mana pun**. Teks bebas yang diketik pengguna hilang bentuk aslinya sebelum sempat tercatat |
| `Activity/RemoveLossAlloction_act.xml` | 2 | ⭐ **BARU** — menulis kronologi *"Delete Loss Allocation"*; satu jenis peristiwa jejak audit yang belum terdaftar |
| `Activity/AddInterest_act.xml` | — | ⭐ **BARU** — catatan pengembangnya *"Untuk History Claim"*, penulis kronologi ketiga di luar kedua penulis yang sudah tercatat |

⚠️ `[terverifikasi]` Sensus **329 berkas berlingkup halaman kronologi**: slot yang dibaca muncul
**34 kali di 27 berkas**; slot yang ditulis step 5 muncul **satu kali di satu berkas** — yaitu langkah
yang salah itu. **Lingkup halaman wajib disebut:** grep nama telanjang memberi 315 dan 57, karena
slot itu juga dipakai sebagai parameter umum di rule SQL lain.

#### ADR terkait

**ADR-U-0014** (tingkat wewenang — sumbernya dari tiket 10).

#### Acceptance criteria

- [ ] `[terverifikasi]` Setiap aksi meninggalkan entri berisi **aksi, pelaku, waktu, dan tingkat wewenang** *(AC 77 spec)*
- [ ] ⚠️ Tingkat wewenang disimpan sebagai **enum tertutup**, dan **nama pelaku di kolom terpisah**. **Alasan menyimpang:** di Pega satu properti berawalan `Is…` menyimpan tingkat wewenang, bukan boolean — namanya berbohong dan isinya tidak dapat disaring *(AC 78 spec)*
- [ ] ⚠️ Migrasi memetakan **kedua ejaan** tingkat terbawah ke satu nilai enum, dan memindahkan nama pelaku ke kolom pelaku. **Alasan menyimpang:** dua penulis mengeja tingkat terbawah berbeda, dan korpus sendiri mengakalinya dengan pencocokan substring *(AC 79 spec)*
- [ ] ⚠️ **Semua peran terekam tanpa kecuali**; test yang menemukan satu peran tidak berjejak **gagal**. **Alasan menyimpang:** di Pega satu jabatan dikecualikan dan seluruh penulisan bersarang di bawah gerbang itu, sehingga aksinya menghasilkan nol entri *(AC 80 spec)*
- [ ] ⚠️ **Memo penutupan klaim tercatat di riwayat.** **Alasan menyimpang:** di Pega memo ditulis ke slot yang tidak pernah dibaca, sehingga alasan penutupan hilang seluruhnya *(AC 81 spec)*
- [ ] ⚠️ `[data DBA]` Riwayat ditampilkan **urut dari kolom tanggal**, bukan urutan baris. **Alasan menyimpang:** Pega menampilkan **urutan baris tersimpan** — grid jejak audit tidak punya konfigurasi pengurutan sama sekali, baik pada `Section/InputAcceptation.xml` maupun `Section/OutstandingClaim.xml`, sehingga mengurutkan menurut tanggal **mengubah apa yang dilihat pengguna**. Ini **penyimpangan sadar**. **Bukti:** satu baris contoh memuat entri terakhir bertanggal lebih lambat dari sebelas lainnya, jadi urutan tersimpan memang bukan kronologis *(AC 82 spec)*
- [ ] `[terverifikasi]` Jejak audit mencatat **jenis aksi**, bukan nilai sebelum/sesudah — batas yang diwarisi dan **tidak** diperluas dalam spec ini *(AC 83 spec)*

#### Perintah verifikasi

```
jalankan test "aksi oleh tiap peran -> semuanya menghasilkan entri jejak"
jalankan test "tutup klaim dengan memo -> memo tercatat dan terbaca kembali"
jalankan test "entri bertanggal tidak urut -> riwayat tampil urut tanggal"
jalankan test "tingkat wewenang tersimpan sebagai enum tertutup, bukan teks bebas"
jalankan test "migrasi: dua ejaan tingkat terbawah -> satu nilai enum, nama ke kolom pelaku"
cari nama pelaku yang tercampur ke kolom tingkat          -> nihil
```

## Claim Prop - 15 - Tiga prefix klaim dan varian Syariah

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 13 (efek keluar)
**Menutup:** AC 75 · 76 *(2 AC)* — US 66

#### Hasil & nilai pengguna

Modul tetap melayani ketiga lini yang selama ini ditanganinya, berikut varian Syariah-nya — tidak ada
lini yang tertinggal saat pindah ke Go. Klaim Syariah diperlakukan identik dengan induknya, sehingga
tidak ada aturan perhitungan kembar yang harus dirawat dua kali.

⚠️ Tiket ini adalah **tiket paritas lintas-lini**: ia menyapu seluruh jalur yang sudah dibangun tiket
sebelumnya dan memastikan ketiganya berlaku untuk semua prefix — bukan hanya untuk prefix yang
kebetulan dipakai saat pengembangan.

#### Area codebase

Penentuan prefix klaim · tiga jalur simpan proyeksi OS akseptasi · penomoran per prefix.

#### Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CekPremiLunas_Act.xml` | 2 · 3 | cabang per prefix — ⚠️ hanya dua dari tiga prefix punya cabang (lihat tiket 02) |
| `RDBList/SaveOSClaim_SQL.xml` | — | jalur simpan proyeksi OS akseptasi — klaim |
| `RDBList/SaveOSKlaimTreaty_SQL.xml` | — | jalur simpan proyeksi — treaty |
| `RDBList/SaveDataToOsAkseptasiNP.xml` | — | jalur simpan proyeksi — treaty non proporsional |
| `RDBList/GenerateNoCLMTreatyIn.xml` | — | penomoran menerima kode dan tipe sebagai parameter |

#### ADR terkait

**ADR-U-0006** (penomoran lewat stored procedure) · **ADR-U-0015** (efek keluar wajib berhasil).

#### Acceptance criteria

- [ ] `[terverifikasi]` Modul melayani **tiga prefix klaim**, dan varian Syariah diperlakukan **identik** dengan induknya *(AC 75 spec)*
- [ ] `[terverifikasi]` **Tiga jalur simpan proyeksi OS akseptasi terpisah** dipertahankan sesuai jenis *(AC 76 spec)*

#### Perintah verifikasi

```
jalankan test "klaim tiap prefix -> terdaftar, ternomori, dan tersimpan"
jalankan test "klaim varian Syariah -> perlakuan identik dengan induknya"
jalankan test "tiap jenis -> masuk ke jalur simpan proyeksi OS yang benar"
jalankan test "regresi lintas prefix: registrasi, estimasi, adjustment, komite, dokumen, efek keluar"
```

# Claim Non Prop

Jumlah tiket: **38**

## Claim Non Prop - 01 - Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis

---
status: selesai
---



> **Ditinjau ulang 19 September 2026 sesudah gerbang penamaan dicabut.**
> Kedelapan pengenal berkas ini diperiksa terhadap aturan final (SPEC bagian 16): terpanjang `KLAIMNP_PERAN_PEMASANGAN` **24 byte**, di atas 30 byte **nol**. Tidak ada nama yang berubah.
> Dua nama dicatat terbuka sebagai bukan-kata-utuh — `KLAIMNP` dan sufiks `_APP` — beserta alasan keduanya tidak diganti. Lihat kepala `ddl-usulan/00_SKEMA_DAN_AKUN.sql`.

> **SELESAI 19 September 2026.** DDL-nya `ddl-usulan/00_SKEMA_DAN_AKUN.sql` — **usulan, untuk dibaca, belum pernah dijalankan**.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-01` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Skema, akun pemilik, akun aplikasi `KLAIMNP_APP`, dan akun hilir `KLAIMNP_HILIR` ada; akun pemilik hanya dipakai saat pemasangan.

Pembuatan skema dan ketiga akun; tablespace; peran dasar.

**Tidak termasuk:** `GRANT` dan `REVOKE` atas objek — belum ada objeknya. → T-20.

**Blocked by:**

- None (can start immediately)


> **DILEPAS 19 September 2026.** ADR-D-CNP-0028 naik ke `accepted` lewat konfirmasi manusia. Gerbang DDL terangkat.
>
> ~~**Gerbang penamaan tidak ikut terangkat.** Sampai REQ-032 kembali: tanpa nama constraint, tanpa nama index, tanpa singkatan.~~ **Dicabut 19 September 2026** — batas 30 byte jadi aturan tetap, dan nama di berkas ini diperiksa terhadap aturan final tanpa satu pun berubah.
>
> **Penegakan ADR-D-CNP-0017 lewat `GRANT` dan `REVOKE` masuk lingkup tiket ini sebagai bagian pekerjaannya**, bukan catatan operasional yang menyusul. Alasannya satu kalimat: **hak yang tidak tertulis tidak dapat diaudit.** Yang dituntut: akun aplikasi mendapat tepat hak yang dipakainya dan tidak lebih; akun hilir hanya `SELECT`; akun pemilik dipakai saat pemasangan lalu tidak dipakai lagi.

**Dasar:** DECIDED(ADR-D-CNP-0028) bagian *Satu pintu tulis*; DECIDED(ADR-D-CNP-0017). Premisnya EVIDENCED: `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` memberi `POOLDATA` hak `ALTER, DELETE, INSERT, UPDATE` atas tabel work Pega.

- [ ] Akun aplikasi dapat membuat sesi dan melihat skema.
- [ ] Akun hilir dapat membuat sesi dan **tidak** melihat satu objek pun (belum ada).
- [ ] Akun pemilik terpisah dari akun aplikasi — bukan akun yang sama dengan nama berbeda.

**Ketidakpastian:** Tidak ada. **REQ-032 turun jadi verifikasi** 19 September 2026: seluruh pengenal dijaga ≤ 30 byte **sebagai aturan tetap**, bukan sebagai pengamanan sementara. Bila versinya 12.2+, yang ada hanyalah kelonggaran yang sengaja tidak dipakai.

#### Hasil

- **Dua tablespace, tiga akun, tiga peran.** Pengenal terpanjang `KLAIMNP_PERAN_PEMASANGAN` = **24 byte** — aman di 12.1 maupun 12.2+, dan itu berlaku di setiap versi Oracle.
- **Gerbang penamaan dipatuhi**: nol nama constraint, nol nama index, nol singkatan. Seluruh pengenal kata utuh. `KLAIMNP_INDEKS` adalah nama *tablespace*, bukan nama index.
- **Nol DML.**
- **Pemisahan pemilik dari aplikasi ditegakkan, bukan disepakati.** Akun aplikasi diberi `QUOTA 0` pada kedua tablespace dan **tidak diberi `CREATE TABLE`** — tanpa itu ia dapat memiliki objek yang tidak pernah masuk DDL dan tidak terlihat siapa pun.
- **Peran bawaan Oracle `CONNECT` dan `RESOURCE` tidak dipakai**, dan dicabut. Isinya berubah antarversi; memakainya berarti hak akun ditentukan versi basis data, bukan ditentukan dokumen.
- **`CREATE DATABASE LINK` dicabut dari ketiganya** (ADR-D-CNP-0016), dengan catatan bahwa instance ini sudah memuat satu link produksi — nasibnya K1, belum diputuskan.
- **Penguncian akun pemilik ditulis tetapi sengaja dibiarkan sebagai komentar**: menjalankannya di tengah pemasangan menghentikan pemasangan itu sendiri. Ia dijalankan sesudah `V00_HAK_AKSES.sql`.
- **Lima blok pemeriksaan penerimaan** ditulis sebagai `SELECT`, menjawab ketiga kriteria terima.

#### Temuan yang lahir dari mengerjakannya

**Hak sistem berakhiran `ANY` mengatasi hak objek, dan itu belum pernah diperiksa.** `V00_HAK_AKSES.sql` mencabut hak `POOLDATA` atas tiap tabel satu per satu; satu akun yang memegang `INSERT ANY TABLE` membatalkan seluruhnya sekaligus, dan tidak ada baris DDL di skema baru yang dapat mencegahnya — pencabutannya di tingkat instance, oleh DBA.

Tidak disebut di ADR-D-CNP-0017, tidak di ADR-D-CNP-0028, tidak di `V00_HAK_AKSES.sql`, tidak di satu pun REQ sebelum hari ini. Didaftarkan **REQ-037**, BLOCKER. Batasnya ditulis di badan ADR-D-CNP-0017 dan di kepala `V00_HAK_AKSES.sql`, bukan hanya dilaporkan di sini.

#### Perbaikan sesudah tinjauan, 19 September 2026

Empat, seluruhnya dari tinjauan atas `00_SKEMA_DAN_AKUN.sql`.

**1. `CREATE SESSION` diperiksa — aman.** Diberikan eksplisit ke ketiga peran (baris 119, 127, 130), **sebelum** `REVOKE CONNECT`. Hak itu datang lewat peran, bukan lewat `CONNECT`, jadi pencabutan `CONNECT` tidak menutup pintu masuk.

**2. Cacat yang ditemukan pemeriksaan itu, dan diperbaiki.** Kesembilan `REVOKE` ditulis sebagai perintah telanjang dengan catatan *"pada instance bersih ia tidak melakukan apa-apa"*. **Itu salah** — Oracle mengangkat `ORA-01951` dan `ORA-01952`, bukan mengabaikan. Pada instance bersih, yaitu keadaan yang justru diharapkan, baris pertama **menggugurkan seluruh pemasangan**. Diganti blok yang mencabut bila ada dan diam bila tidak, menelan hanya kedua galat itu.

**3. Penguncian akun pemilik keluar dari komentar, jadi `Z02_KUNCI_PEMILIK.sql`.** Alasannya sejarah proyek ini: dasar ADR-D-CNP-0024 adalah blok yang dikomentari, dan `PEGA_JSON_OS_AKSEP_KLAIMTNP` selalu `INSERT` tanpa pernah `UPDATE` karena logika upsert-nya dikomentari seluruhnya. Kode yang dimatikan dengan niat dinyalakan nanti tidak pernah dinyalakan. Sapuan atas seluruh berkas pemasangan sekarang menemukan **nol baris perintah yang dikomentari**.

**4. Pengawasan tulis ditambahkan** — `Z01_PENGAWASAN_TULIS.sql`. Skema tidak dapat **mencegah** tulisan dari pemegang hak `ANY`, tetapi dapat **memperlihatkannya**. Syaratnya memakai `SESSION_USER`, bukan `CURRENT_USER`: di dalam prosedur definer's rights, `CURRENT_USER` menjadi pemilik prosedur dan jejaknya hilang — persis mekanisme `PEGA_JSON_OS_AKSEP_KLAIMTNP`. Dijalankan **DBA**, dan `KLAIMNP` sengaja tidak diberi `AUDIT_ADMIN`: pengawasan yang dapat dimatikan oleh yang diawasi bukan pengawasan.

**Akibatnya bagi ADR-D-CNP-0017**: rumusan klaimnya diubah di badannya dari *mencegah* menjadi **"tidak dapat mencegah tulisan dari luar pintu, tetapi tidak ada tulisan dari luar pintu yang tidak meninggalkan jejak"**. Keputusannya tidak berubah.

**REQ-037 diperluas dari satu jalur menjadi enam**, dan satu di antaranya mengubah sifat REQ itu sendiri: pemegang `GRANT ANY PRIVILEGE` dapat memulihkan jalur tulis **besok**, sesudah pemeriksaan hari ini bersih. REQ mengukur satu titik waktu; ADR bicara keadaan yang bertahan. Pembedaan itu tertulis di badan REQ-nya.

**Catatan kerapian yang ditunda**: `V00_HAK_AKSES.sql` bukan view, dan awalan `V` padanya tidak konsisten dengan `V01`–`V09` yang memang view. Tidak diganti nama karena sudah dirujuk ADR-D-CNP-0028, tiket `01`, dan tiket `35`.

## Claim Non Prop - 03 - Penamaan yang tidak nyaris kembar

---
status: selesai
---



> **SELESAI 2026-09-18.** Penggantian nama sudah dikerjakan: KEPUTUSAN_KOMITE, dan sufiks _FK tersisa nol.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.
>
> **SELESAI PENUH 19 September 2026.** Pilahan *"aturannya selesai, namanya draf"* dicabut: namanya kini berlaku juga.
>
> REQ-032 turun jadi verifikasi, dan **daftar singkatan dibatalkan** — tiket `02` mati. Nama constraint tidak lagi menunggu daftar apa pun, karena ia tidak lagi memakai singkatan.
>
> **Aturan penamaan final — empat, berlaku untuk seluruh skema:**
>
> 1. **Nama tabel dan kolom memakai kata utuh bahasa Indonesia**, `UPPER_SNAKE_CASE`, mengikuti glosarium. Tidak ada singkatan buatan. Bila sebuah nama melewati 30 byte, yang diganti **katanya** dengan kata utuh yang lebih pendek — **bukan dipotong jadi singkatan**.
> 2. **Nama constraint dan index tidak mengeja kolom.** Bentuknya `<peran>_<tabel>`, ditambah pembeda numerik bila sebuah tabel punya lebih dari satu constraint sejenis. Isinya dibaca dari katalog.
> 3. **Peran sebagai awalan tetap**: `PK_` `UQ_` `FK_` `CK_` `IX_`, ditambah `V_` untuk view dan **`SQ_` untuk sequence**. Bukan singkatan bahasa — penanda jenis objek yang baku di Oracle. *(`V_` dan `SQ_` ditambahkan 19 September 2026; `SQ_` menyusul bersama 22 sequence yang dibuat hari itu, dan daftarnya diperbarui supaya aturan ini tidak lagi tertinggal dari berkasnya sendiri.)*
> 4. **Tiap nama diperiksa panjangnya saat ditulis**, dan pembangkit **berhenti** bila ada yang melewati 30 byte.
>
> Daftar `_Avoid_` tetap berlaku penuh di nama tabel maupun kolom.
>
> **Dikerjakan 19 September 2026**: 105 objek bernama ditulis ulang di seluruh `ddl-usulan/`; satu nama tabel diganti karena constraint-nya melewati batas — `DAFTAR_KLAIM_PENJAGA_TANGGAL` → **`KLAIM_PENJAGA_TANGGAL`**, dan kolom `ID_DAFTAR` → `ID_PENJAGA_TANGGAL` supaya tidak menyebut kata yang sudah tidak ada di nama tabelnya. Pengenal di atas 30 byte: **nol**. Rinciannya `SPEC-MODEL-DATA.md` bagian 16.

*Asal: `T-14` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada pasangan nama berjarak satu kata di skema yang belum berjalan.

`ADJUSTMENT.STATUS_AKSEPTASI` → **`KEPUTUSAN_KOMITE`**, dinamai menurut isinya: ia hasil keputusan Komite atas Adjustment, bukan keadaan entitas akseptasi. `FK_DPT_KLM_FK` → `FK_KLAIM_PENJAGA_TANGGAL_1`; sufiks `_FK` tersisa: **nol**.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §10.4 butir 4 mendaftar nama nyaris kembar sebagai technical debt sistem lama — `CountTotalInterest_Act` vs `CountTotalInsterest_Act`, `ProtectNilaiClaim` vs `ProteksiNilaiClaim`.

- [x] Tidak ada dua pengenal di skema yang berbeda hanya pada satu kata.
- [x] Tidak ada constraint bersufiks `_FK`.
- [x] Tidak ada singkatan buatan di nama tabel maupun kolom.
- [x] Tidak ada nama constraint yang mengeja kolomnya.
- [x] Pengenal di atas 30 byte: **nol**, diperiksa pembangkit atas 105 objek bernama.

**Ketidakpastian:** Tidak ada. REQ-032 memverifikasi, dan jawabannya tidak mengubah satu nama pun — 30 byte sah di setiap versi.

## Claim Non Prop - 04 - Index penopang setiap foreign key

---
status: selesai
---



> **SELESAI 2026-09-18.** IX_ADJUSTMENT_1 terpasang; pembangkit kini memberi index penopang untuk setiap FK.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.
>
> **SELESAI PENUH 19 September 2026.** Pilahan *"aturannya selesai, namanya draf"* dicabut.
>
> Nama index tidak lagi menunggu apa pun: daftar singkatan dibatalkan dan REQ-032 turun jadi verifikasi. Index kini dinamai `IX_<tabel>_<n>` sesuai aturan 2 di tiket `03`; `IX_ADJ_REKENING` menjadi **`IX_ADJUSTMENT_1`**, penopang `FK_ADJUSTMENT_2`.
>
> Yang **tidak** berubah: kewajibannya. Setiap FK tetap harus punya index atau unique yang kolom pertamanya sama dengan kolom pertama FK itu, dan pembangkit tetap menegakkannya.

*Asal: `T-15` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada FK yang kolomnya bukan awalan sebuah index atau unique.

`IX_ADJUSTMENT_1` untuk `FK_ADJUSTMENT_2`; sapuan ulang seluruh FK oleh pembangkit.

**Blocked by:**

- None (can start immediately)


**Dasar:** DERIVED: di Oracle, FK tanpa index membuat `DELETE` dan `UPDATE` pada tabel induk mengambil kunci di tingkat tabel.

- [x] Setiap FK punya index atau unique yang kolom pertamanya sama dengan kolom pertama FK itu.
- [x] Empat belas FK ke `KLAIM` tetap tertopang `UNIQUE` berawalan `ID_KLAIM` — **tidak boleh hilang** saat kunci disentuh lagi.

**Ketidakpastian:** Tidak ada.

## Claim Non Prop - 05 - Sapuan S1: kolom yang sesungguhnya diterima Arasapas dan kasir

---
status: selesai
---



> **SELESAI 19 September 2026.** Sisi Arasapas dan sisi kasir keduanya terpetakan; hasilnya `SAPUAN-CELAH-DAN-JAVA.md` §5.1.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

> **DIBUKA KEMBALI 2026-09-18 — sapuan ulang menemukan cakupannya kurang, bukan polanya salah.**
> Pola literal tiga-bentuk tidak mengubah satu pun dari 17 parameter `InputParamOs.*`. Yang kurang adalah **sisi kasir**: `HitServiceToKasir_Act` memakai 25 medan bernama urut `TempKasir.CARI1`…`CARI24` dan `CARIDATETIME`, dan tidak satu pun pernah diinventarisasi. Judul tiket ini menyebut **“Arasapas dan kasir”**; yang terpetakan baru Arasapas.
> Lima berkas juga tidak disebut laporan S1: `SaveToOS`, `GetSelisihActual_Act`, `CloseClaimTNonProp`, dan dua layar uji `Hitung_Test` yang mengikat medan langsung ke `InputParamOs.GrossValue`.
> Dicatat sebagai **D36**. Rinci di `SAPUAN-S10-S16.md` §1.

*Asal: `T-32` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** 17 parameter `InputParamOs.*` dan 25 parameter `TempKasir.CARI*` terbaca, beserta asal tiap satunya.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED: `Activity\SaveDataToOSAksep_Act.xml`, `Activity\HitServiceToKasir_Act.xml`. Hasil lengkap di `SPEC-MODEL-DATA.md` bagian 18.

- [ ] Terpenuhi. Temuan yang mengubah T-16: keempat nilai uang dikirim sebagai **selisih**, bukan nilai penuh.

**Ketidakpastian:** `InsertOSKlaimCNP` **TIDAK DITEMUKAN** sebagai nama berkas; `KonversiKlaim_Act.xml` ada dengan **nol** pasangan properti; `ConnectREST\KonversiKlaimNonLife.xml` belum dibaca isinya. → T-37.

#### Hasil

- **Sisi Arasapas**: 17 parameter `InputParamOs.*` — tidak berubah oleh sapu ulang pola tiga-bentuk.
- **Sisi kasir**: **21 medan bernama dari 25 slot** `TempKasir.CARI*`. Nama bisnisnya terbaca dari 480 baris Java yang belum pernah dibaca sebagai kode. `CARI19`–`CARI21` hanya perakit tanggal; `CARIDATETIME` penampung sementara. **D46**.
- **`Nett`, `Deductible`, `KaliDeduct` dikirim tanpa kutip** sebagai angka JSON, dari `String` hasil `.toString()`. Nilai kosong menghasilkan JSON rusak. **D41**.
- **Tujuh nilai tertanam di kode**, termasuk `StsSyariah = 0` dan `CompanyName = "NUSARE"` — jalur pembayaran tidak pernah menyatakan syariah. **D42**.
- Satu hal **tidak** disimpulkan: `CARI20 = @toDecimal(@month(...)) + 1`. Perilaku `@month()` tidak ada di ekspor — diajukan sebagai **REQ-036**, bukan ditebak.

## Claim Non Prop - 06 - Keadaan kasus: nilai mana yang ada, dan siapa yang membacanya

---
status: selesai
---



> **SELESAI 19 September 2026.** Dijalankan di atas indeks utuh 45 tag; hasilnya `SAPUAN-CELAH-DAN-JAVA.md` §5.2.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

> **DIBUKA KEMBALI 2026-09-18 — judulnya diganti, karena pertanyaannya bukan yang semula ditulis.**
> Judul lama, *“Sapuan S2: kelengkapan enum `CNPStatusCase`”*, mengandaikan ada enum yang tinggal dilengkapi. Sapu ulang dengan pola tiga-bentuk membatalkan andaian itu dari dua arah:
>
> 1. **`"CLAIM ACCEPTED"` dan `"CLAIM REJECTED"` muncul nol kali** di 279 berkas — bukan nol sebagai nilai `CNPStatusCase`, melainkan nol **di mana pun**. Dua dari empat nilai yang didaftar memori §7.3 tidak pernah ada. Butir **C2** ditutup sebagai premis gugur atas dasar ini.
> 2. **`CNPStatusCase` tidak pernah dibaca.** Tiga penulisan, nol pembacaan pada empat lapisan. **Enum yang tidak pernah dibaca bukan enum yang lengkap** — ia bukan enum.
>
> Yang tersisa karena itu bukan *melengkapi daftar* melainkan: nilai mana yang benar-benar ada, siapa yang membacanya (bila ada), dan apakah keadaan kasus disimpan di tempat lain. Dicatat sebagai **D35**; penanda yatim lainnya di **D25**.
> Ejaan tersimpan `COMITEE` satu T, sementara `COMMITTEE` dua T muncul 10 kali sebagai teks tampilan.

*Asal: `T-33` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Empat lapisan disapu; **2 nilai ditulis**, 3 berkas, **nol** rule menguji, dan `"CLAIM ACCEPTED"`/`"CLAIM REJECTED"` **nihil di 279 berkas**.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED. Menguatkan `BLUEPRINT.md` §3.1 — semula berlaku untuk lapisan Activity saja.

- [ ] Terpenuhi. Akibatnya: `STATUS_KLAIM` **tanpa `CHECK`** (T-12).

**Ketidakpastian:** Kelengkapan enum hanya dapat dibuktikan dengan membuka folder Komite — **belum diizinkan**.

#### Batas jawaban — ditetapkan 19 September 2026

Folder `Komite Claim Non Prop` **tertutup untuk batch ini**. Jawaban tiket ini karena itu berbunyi:

> *"Tidak ada pembaca `CNPStatusCase` **di folder `Claim Non Prop`**"*

bukan *"tidak ada pembaca"*. **Batas itu ditulis di badan jawabannya**, bukan disimpulkan pembaca nanti.

Batas kedua, dari pemeriksaan cakupan: pernyataan itu juga berbatas pada **tag yang benar-benar dibaca indeks**. Tiket ini dikerjakan **sesudah** celah tag ditutup, supaya batas keduanya tinggal satu — folder, bukan folder *dan* tag.

#### Hasil

- **Dua nilai** yang benar-benar ditulis: `"COMITEE ACCEPTANCE (DEPT. HEAD)"` dan `"INPUT ACCEPTATION CLAIM"`, di tiga tempat.
- **Nol pembacaan, nol perbandingan** — diuji terhadap 45 tag pada 279 berkas, 12 jenis rule, ditambah 480 baris Java.
- **Batas jawaban, dan ia bagian dari jawabannya**: pernyataan ini berlaku **di folder `Claim Non Prop`**. Folder Komite tertutup untuk batch ini, dan ketiga penulisnya berurusan dengan Komite — jadi pembacanya, bila ada, kemungkinan besar di sana.
- `"CLAIM ACCEPTED"` dan `"CLAIM REJECTED"` **nol kemunculan di mana pun** (D35). Dua dari empat nilai memori §7.3 tidak pernah ada.
- **Ia bukan enum yang perlu dilengkapi.** Kolom teks, diisi tiga tempat, tidak mengatur apa pun di modul ini. Diperlakukan menurut kebijakan **E17**: dimigrasi apa adanya, ditandai tak berpemilik, tanpa domain tertutup.
- Ejaan tersimpan **`COMITEE` satu T**; `COMMITTEE` dua T hanya teks tampilan, 10 kali di 7 berkas.

## Claim Non Prop - 07 - Sapuan S3: perilaku per nilai `PaymentType`

---
status: selesai
---



> **SELESAI 2026-09-18.** Dijalankan 18 September 2026 atas empat lapisan; hasilnya `_migration-docs/claim-non-prop/SAPUAN-S3-S9.md` bagian S3.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-34` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tabel nilai → perilaku yang dijaganya, disusun dari setiap kondisi dan ekspresi yang mengujinya.

Empat lapisan: Activity, Data Transform, Section/Harness/FlowAction, pemetaan RDB/REST.

**Tidak termasuk:** **Nama** tiap nilai — itu tetap **A6b**, pertanyaan untuk pemilik proses.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §7.4 menyebut `1` Final dan `2` Partial/interim, **disimpulkan dari pemakaian**. `Rule-Obj-FieldValue` **tidak ada** di ekspor.

- [ ] Setiap nilai yang muncul di penjaga terdaftar beserta berkas dan barisnya. Hasil nihil tetap dicatat beserta pola yang dipakai.

**Ketidakpastian:** Sebaran nilai yang sesungguhnya dipakai diukur **REQ-027**. Bila hasilnya menunjukkan himpunan tertutup, `CHECK`-nya dipasang saat itu — satu `ALTER`, bukan pembongkaran.

#### Hasil sapuan

- **Ketujuh nilai `PaymentType` EVIDENCED**, beserta nilai kosong — delapan kemungkinan, bukan tiga.
- **Penolakan saya atas usulan §4.1 dicabut.** Pola sapuan lama hanya mencakup literal berkutip dan melewatkan perbandingan tanpa kutip di `HitServiceToKasir_Act` baris 7754 dan 7896. Dicatat sebagai **D9**.
- Nilai instruksi bayar menambahkan **tepat satu** besaran menurut `PaymentType`, bukan menjumlahkan seluruhnya. Dicatat sebagai **D15**.
- **Arti** tiap nilai tetap TIDAK DITEMUKAN — definisi propertinya ada di schema PegaRULES, yaitu **REQ-001**.

## Claim Non Prop - 08 - Sapuan S4: adakah padanan IDR untuk biaya penilaian, salvage, dan biaya lain

---
status: selesai
---



> **SELESAI 2026-09-18.** Dijalankan 18 September 2026 atas empat lapisan; hasilnya `SAPUAN-S3-S9.md` bagian S4.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-35` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Jawaban apakah ketiganya memang tanpa pasangan IDR, atau kolomnya yang kurang.

Sapuan seluruh properti bersufiks `IDR` pada class `SpreadingRisk`, `Data-Adjustment`, dan `ClaimData`, dibandingkan terhadap daftar properti bernilai uang. Empat lapisan.

**Blocked by:**

- None (can start immediately)


**Dasar:** DERIVED dari review; DECIDED(ADR-D-CNP-0007) mewajibkan pasangan untuk setiap nilai uang.

- [ ] Daftar properti uang yang punya padanan IDR dan yang tidak, beserta berkas dan baris.

**Ketidakpastian:** **Menahan T-05.** Bila ketiganya perlu pasangan IDR, kolomnya bertambah di `ALOKASI_LAYER` dan `CHECK`-nya ikut terpasang sendiri lewat pembangkit.

#### Hasil sapuan

- **Ya, ada kelas nilai uang yang hidup tanpa padanan IDR sama sekali** — tujuh properti, lima di antaranya masuk hitungan instruksi bayar.
- Penamaan alternatif (`Rupiah`, `Idr`, `_IDR`, `Rp`, `InRp`) ikut disapu: **nol hasil**. Yang ada hanya sufiks `RNM`, yaitu porsi pihak, bukan mata uang.
- **Menyentuh usulan ADR-D-CNP-0029 dan tidak saya selesaikan sendiri.** Dua pembacaan berdiri di atas bukti yang sama; membedakannya butuh jawaban orang. Diajukan sebagai **A16** ke akuntansi.

## Claim Non Prop - 09 - Sapuan S9: apa persisnya yang ditulis `EditXOLAlokasi`

---
status: selesai
---



> **SELESAI 2026-09-18.** Dijalankan 18 September 2026 atas empat lapisan; hasilnya `SAPUAN-S3-S9.md` bagian S9.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-36` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Daftar properti yang benar-benar dapat disunting petugas.

Seluruh `Property-Set` di rule itu, per properti. Empat lapisan.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED: `BLUEPRINT.md` §2.5 — `EditXOLAlokasi` menulis `.IsEditClaim = 1` pada class `ASM-FW-GISFW-Data-SpreadingRisk`, yang di sistem lama memuat baris layer **dan** baris Retensi Cedant.

- [ ] Daftar properti beserta baris; dan pernyataan tegas properti mana yang **tidak** tersentuh.

**Ketidakpastian:** **Menahan T-05.** Memasang `_HITUNG`/`_SUNTING` di kolom yang tidak pernah disunting adalah beban mati; tidak memasangnya di kolom yang disunting berarti membongkar tabel kemudian (ADR-D-CNP-0008).

#### Hasil sapuan

- `EditXOLAlokasi` menulis **tiga** hal saja: pesan galat lokal, penanda `.IsEditClaim = 1`, dan satu catatan kronologi. **Tidak satu pun properti bernilai uang.**
- Ia **gerbang kata sandi**, bukan penyunting. Layar dan harness-nya memuat dua medan saja, `CARI1` dan `CARI2`.
- **Pertanyaan `_HITUNG`/`_SUNTING` tidak terjawab** — bukan karena sapuannya gagal, melainkan karena tertuju pada rule yang keliru. Sasaran penggantinya belum ditetapkan; dicatat di daftar yang tetap tertutup.
- Temuan tambahan: `.IsEditClaim` ditulis tiga kali dan **dibaca nol kali**. Dicatat sebagai **D10**.

## Claim Non Prop - 10 - Sapuan S5, S6, S7, S8, dan dua sumber yang belum pernah dibaca

---
status: selesai
---



> **SELESAI 2026-09-18.** Dijalankan 18 September 2026; hasilnya `SAPUAN-S3-S9.md` bagian S5, S6, S7, S8, dan dua sumber.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-37` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Empat pertanyaan terbuka terjawab atau tercatat nihil beserta polanya.

- **S5** penulis `AlokasiXOLPaid` — termasuk sebagai sasaran `Page-Copy`, `Page-New`, `RDB-List` ke page, dan hasil Report Definition, bukan hanya `Property-Set`.
- **S6** asal `"Previously Calculated UR"` — **jangan cari literal utuh**; cari potongan `"Previously"` dan `"Calculated"` di dalam ekspresi penyambungan, karena nilainya mungkin dirakit.
- **S7** `Property-Remove` di `CountLossAllocation_act` langkah 10 yang parameter sasarannya kosong — baca langkah itu utuh beserta tetangganya.
- **S8** `IsTreatyIn`, `IsReject`, `IsCloseFile`, `IsAnyAcceptation` — keempatnya belum berlaku pernyataan "bersih".
- **Dua sumber**: `pengetahuan/DDL_Script_ClaimNonProp2.xls` (**belum pernah dibaca**) dan `excludeXML/GetBase64Attachment.xml`.

**Tidak termasuk:** H1/H2, I1/I2, F1/F2 — berada di modul Komite.

**Blocked by:**

- None (can start immediately)


**Dasar:** `_selesai/OPEN-QUESTIONS.md` C5 dan C9; `PENGETAHUAN.md` §15.

- [ ] Tiap sasaran menghasilkan jawaban atau pernyataan nihil **beserta lapisan yang disapu dan pola yang dipakai**, sesuai `PENGETAHUAN.md` §14: yang sah hanya *"tidak ada X di lapisan yang sudah saya baca"*. Temuan yang mengubah putusan lama masuk `_selesai/OPEN-QUESTIONS.md` sebagai D9 dan seterusnya. Hasilnya dibandingkan terhadap `BLUEPRINT.md` §7.5 — tambalan per-case baru adalah temuan tersendiri (ADR-D-CNP-0012).

**Ketidakpastian:** S6 menentukan apakah `RETENSI_CEDANT` perlu **keadaan ketiga**; bila ya, kunci (klaim, mata uang) **belum cukup** dan T-05 berubah.

#### Hasil sapuan

- **S5** — `AlokasiXOLPaid` **tanpa penulis** pada empat lapisan (EVIDENCED-NIHIL). Yang justru terbaca: daftar kolom muatan selisih jalur Komite. **D11**.
- **S6** — `"Previously Calculated UR"` memang **dirakit**, lalu diuji utuh. Kunci (klaim, mata uang) **belum cukup**; menaikkan pentingnya REQ-033. **D12**.
- **S7** — langkah 10 memang kosong: sasaran, deskripsi, prasyarat, dan parameter wajibnya seluruhnya tidak terisi. Ia sisa, bukan perilaku. Migrasi tidak perlu menirunya.
- **S8** — pernyataan “bersih” **tidak** berlaku seragam. `IsTreatyIn` dibaca empat Activity dengan tipe tak konsisten; `IsReject` dan `IsCloseFile` ditulis dan dibaca; `IsAnyAcceptation` **dibaca sepuluh kali dan ditulis nol kali**. **D10**.
- **`DDL_Script_ClaimNonProp2.xls` dibaca** — 39 objek, empat di antaranya menyentuh tiket berjalan: `DIRECTTOKASIR_LOG` (**D14**), `JSON_KLAIM` (**D13**, mengubah sifat REQ-018), `CLAIMXOL2`, `CLAIMREJECTED`.
- **`excludeXML/GetBase64Attachment.xml` TIDAK DIBACA** — ia di dalam folder `Komite Claim Non Prop`, yang tertutup. Saudaranya di folder terbuka dibaca sebagai gantinya, dan perbedaannya tidak diketahui.

## Claim Non Prop - 12 - Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap

---
status: selesai
menunggu-luar: [REQ-018]
---




> **SELESAI 19 September 2026.** `ddl-usulan/01_KLAIM.sql` ditulis ulang lengkap.
>
> **Yang berubah dari versi sebelumnya — dua kolom baru, keduanya konsekuensi keputusan yang turun hari ini:**
>
> | Kolom | Dari | Kenapa di sini |
> |---|---|---|
> | `LINI_USAHA` | SPEC §21.3 (menutup G1 dan A15) | Penentu syariah lama ada dua — node yang mengeksekusi dan nama akun notifikasi — dan **tidak satu pun ikut pindah**. Apa pun jawaban G1, sistem baru perlu penanda yang berdiri sendiri. Disediakan sekarang meski belum ada yang mengisinya, sebagaimana ADR-D-CNP-0025 menyediakan kolom lingkup sebelum ada yang memerlukannya |
> | `PENUTUPAN_LAMA` | SPEC §21.6 (menutup B9 dan A7) | Jalur `CloseClaimMD` **tidak dimigrasi**, tetapi klaim yang pernah lewat sana dimigrasi sebagai tertutup **dan ditandai**. Tanpa kolom ini, keputusan "ditandai" tidak punya tempat, dan asal-usul baris hilang saat jalurnya dibuang |
>
> Keduanya **tanpa CHECK domain**, dan itu disengaja: nilainya belum terbaca, dan menutup domain yang belum diketahui adalah cacat yang sama dengan yang membuat `STATUS_KLAIM` dibiarkan terbuka (D1, D35, SPEC §21.1).
>
> **Nama constraint final**: `PK_KLAIM`, `UQ_KLAIM_1`, `CK_KLAIM_1`, `CK_KLAIM_2`, `CK_KLAIM_3`, `IX_KLAIM_1` — bentuk `<peran>_<tabel>[_n]`, tidak mengeja kolom. Terpanjang **10 byte**.
>
> **Empat aturan dinyatakan jatuh ke aplikasi**, bukan diam-diam tidak ditegakkan — ketetapan `TANGGAL_KEJADIAN`, domain `STATUS_KLAIM` dan `LINI_USAHA`, keterkaitan `PENUTUPAN_LAMA`, dan ketetapan urutan penulisan. Ditulis di kaki berkas DDL-nya.
>
> **REQ-018 tetap menunggui, tidak menahan.** `UQ_KLAIM_1` dipasang sekarang; bila data lama melanggarnya, yang berubah **penanganan migrasi**, bukan kuncinya.

*Asal: `T-03` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Klaim dapat disisipkan; nomor klaim ganda ditolak; masa berlaku treaty terbalik ditolak.

`KLAIM` beserta `PK_KLAIM`, `UQ_KLAIM_1`, `CK_KLAIM_1`, `CK_KLAIM_2`, `CK_KLAIM_3`, `IX_KLAIM_1` — `ddl-usulan/01_KLAIM.sql`.

**Tidak termasuk:** `CHECK` domain untuk `STATUS_KLAIM` — sengaja tidak ada, lihat T-12. Kolom jejak pelaku dipasang di sini tetapi tanpa foreign key; tabel pengguna tidak dirancang (ADR-D-CNP-0006).

**Blocked by:**

- **REQ-018** — **menunggui, tidak menahan** (dari luar papan): kuncinya sudah ditetapkan ADR-D-CNP-0024; yang berubah bila datanya melanggar adalah **penanganan migrasi**, bukan kunci
- ~~`01`~~ *(selesai)* — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis


**Dasar:** DECIDED(ADR-D-CNP-0001, ADR-D-CNP-0021, ADR-D-CNP-0022). EVIDENCED: `BLUEPRINT.md` §2.2 — sapuan 114 activity, **tidak satu pun `Property-Set` menulis `.DateOfLoss`**.

- [x] Klaim dengan nomor yang sudah ada **ditolak** — `UQ_KLAIM_1`.
- [x] Klaim ber-`TREATY_MULAI` sesudah `TREATY_AKHIR` **ditolak** — `CK_KLAIM_1`.
- [x] Klaim dengan `KEADAAN_BARIS` di luar tiga nilai domain **ditolak** — `CK_KLAIM_3`.
- [x] Klaim tanpa `DIBUAT_OLEH` **ditolak** (`NOT NULL`); `DIBUAT_ATAS_NAMA` kosong **diterima** — kosong berarti tidak ada perwakilan, bukan tidak diketahui.

**Ketidakpastian:** **REQ-018** (duplikat `PYID`) belum terjawab. `UQ_KLAIM_1` dipasang sekarang. Bila data lama melanggarnya, yang berubah adalah **penanganan migrasi**, bukan kuncinya — ADR-D-CNP-0024 sudah menetapkan pemilahan tiga kelompoknya.

## Claim Non Prop - 13 - Bentuk lama boleh masuk utuh, di satu tempat saja

---
status: selesai
---




> **SELESAI 19 September 2026, dikerjakan bersama tiket `14`.** Keduanya dua sisi satu seam: apa yang mendarat, dan apa yang tidak terurai darinya.
>
> **Satu kriteria penerimaan DIGANTI — bukan diberi catatan, tetapi ditulis ulang di daftarnya.**
>
> Bunyi lama: *"Muatan yang bukan JSON sah **ditolak**."* Itu **dicabut dan diganti tiga kriteria baru** di bawah. `CK_MIGRASI_PENDARATAN_1` tidak lagi menjaga `MUATAN IS JSON`.
>
> Dua sebab, keduanya terbaca dari sistem lama:
>
> | Sebab | Bukti |
> |---|---|
> | Bentuk muatan lama **tidak dibatasi** | **D43** — `adoptJSONObject` di tiga rule mengadopsi teks `HASIL1` tanpa memeriksa bentuknya, dan salah satunya menulis hasilnya ke halaman objek kerja. Bentuk yang tidak dibatasi tidak dapat dijaga dengan CHECK; ia hanya dapat **ditampung** |
> | Sistem lama **memang menghasilkan JSON rusak** | **D41** — muatan kasir dirakit dengan penyambungan teks; `Nett`, `Deductible`, `KaliDeduct` dikirim tanpa kutip. Nilai kosong menghasilkan `"Nett":,` — JSON tidak sah |
>
> Gerbang `IS JSON` akan menolak **persis muatan yang paling perlu tercatat**, dan menolaknya bertabrakan dengan **AK-4**: tidak ada baris yang dibuang karena "cuma sedikit". Muatan yang ditolak di pintu tidak punya rumah.
>
> **Yang menggantikannya**: muatan mendarat utuh apa pun bentuknya, dan hasil penguraiannya **dicatat sebagai fakta** — `BENTUK_TERURAI` biner, `SEBAB_GAGAL_URAI` terisi bila dan hanya bila gagal. `CK_MIGRASI_PENDARATAN_2` menolak dua keadaan yang mustahil: gagal tanpa sebab, dan berhasil tapi bersebab.
>
> Itu `SPEC-MODEL-DATA.md` §21.8 sebagai tabel.

*Asal: `T-21` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Satu-satunya tempat JSON hidup di skema ini.

`MIGRASI_PENDARATAN` beserta `PK_MIGRASI_PENDARATAN`, `CK_MIGRASI_PENDARATAN_1`, `CK_MIGRASI_PENDARATAN_2`, `IX_MIGRASI_PENDARATAN_1` — `ddl-usulan/21_MIGRASI_PENDARATAN.sql`.

**Blocked by:**

- ~~`01`~~ *(selesai)* — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis


**Dasar:** DECIDED(ADR-D-CNP-0028). EVIDENCED: `BLUEPRINT.md` §13.2, §13.5 — nilai uang lama tersimpan sebagai **teks di dalam JSON**, dan `CLAIMXOL` mengeluarkannya sebagai `varchar2`.

- [x] **Muatan mendarat utuh, apa pun bentuknya** — tidak ada muatan yang ditolak di pintu.
- [x] **Hasil penguraian tercatat sebagai fakta**, bukan dipakai sebagai syarat masuk — `BENTUK_TERURAI`, dan `SEBAB_GAGAL_URAI` terisi bila dan hanya bila gagal (`CK_MIGRASI_PENDARATAN_2`).
- [x] Muatan yang gagal terurai **tetap tersimpan utuh** dan dapat ditunjuk dari `MIGRASI_NILAI_DITOLAK`.
- [x] Tidak ada kolom JSON di tabel kanonik mana pun — diperiksa dengan sapuan atas seluruh skema.

**Ketidakpastian:** Tidak ada REQ yang menahan. REQ-009 memverifikasi — ia memperbanyak daftar medan yang dikenali, dan tabel ini menerima apa pun yang datang entah medannya dikenali atau tidak.

## Claim Non Prop - 14 - Nilai yang tidak dapat diurai tercatat, tidak dibulatkan, tidak dibuang

---
status: selesai
---




> **SELESAI 19 September 2026, dikerjakan bersama tiket `13`.**
>
> **Yang bertambah: kolom `ID_PENDARATAN`, dan itulah seam-nya.** Tanpa kolom itu, "tercatat **beserta asalnya**" (`SPEC-MODEL-DATA.md` §21.8) hanya berlaku sampai tingkat nama sumber — dan nama sumber tidak menunjuk baris mana.
>
> `FK_MIGRASI_NILAI_DITOLAK_1` menegakkan bahwa muatan yang ditunjuk memang ada; `IX_MIGRASI_NILAI_DITOLAK_1` menopangnya sesuai aturan tiket `04`.
>
> **Kolomnya boleh kosong, dan kosong punya arti**: sebagian nilai ditolak berasal langsung dari kolom tabel lama, bukan dari muatan. Kosong berarti **bukan dari muatan** — bukan **tidak diketahui** (ADR-D-CNP-0019). `SUMBER_LAMA` dan `PENGENAL_LAMA` tetap wajib, jadi asalnya selalu terbaca.
>
> **Arah relasinya dari sini ke sana**, karena satu muatan dapat melahirkan banyak nilai ditolak.
>
> `NILAI_MENTAH` tetap `VARCHAR2`, bukan `NUMBER`: mengubahnya jadi angka adalah persis hal yang gagal dilakukan.

*Asal: `T-22` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada angka yang hilang tanpa jejak.

`MIGRASI_NILAI_DITOLAK` beserta `PK_MIGRASI_NILAI_DITOLAK`, `FK_MIGRASI_NILAI_DITOLAK_1`, `IX_MIGRASI_NILAI_DITOLAK_1` — `ddl-usulan/20_MIGRASI_NILAI_DITOLAK.sql`.

**Blocked by:**

- ~~`01`~~ *(selesai)* — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis


**Dasar:** DECIDED(ADR-D-CNP-0014). DECIDED-TEKNIS(AK-4): migrasi berhenti bila ada **satu** baris yang tidak dapat diurai **dan** tidak dapat diselesaikan. Tidak ada baris yang dibuang karena "cuma sedikit".

- [x] Nilai mentah tersimpan **apa adanya**, tanpa pembulatan dan tanpa konversi.
- [x] Setiap baris membawa sebab penolakan dan pengenal barisnya di sistem lama.
- [x] Nilai yang berasal dari sebuah muatan menunjuk **baris muatannya**, bukan hanya nama sumbernya — `ID_PENDARATAN`.

**Ketidakpastian:** Tidak ada REQ. Volume belum terukur — **REQ-031**.

## Claim Non Prop - 15 - Jembatan ke sistem lama berdiri sebagai tabel terpisah

---
status: selesai
menunggu-luar: [REQ-018]
---




> **SELESAI 19 September 2026.** `ddl-usulan/19_MIGRASI_KORELASI.sql` ditulis ulang lengkap.
>
> **Dua hal bertambah, dan keduanya lahir dari membaca ulang janji tiket ini.**
>
> **1. `CK_MIGRASI_KORELASI_1` — sekurangnya satu pengenal lama wajib ada.** Kelima kolom boleh kosong sendiri-sendiri (itu memang kriteria penerimaan kedua), tetapi **kosong semua** menghasilkan baris jembatan yang tidak menjembatani apa pun. Ini satu-satunya keutuhan yang dapat ditegakkan basis data di tabel ini — karena itu dipasang.
>
> **2. `IX_MIGRASI_KORELASI_2` dan `_3` — arah lama → baru.** Versi sebelumnya hanya punya index `(TABEL_TUJUAN, ID_TUJUAN)`, yaitu arah **baru → lama**. Judul tiket ini menjanjikan arah yang **berlawanan**: *"dapat ditunjuk balik lewat pengenal lamanya"*. Tanpa index atas `PZINSKEY_LAMA` dan `PYID_LAMA`, janji itu berjalan lewat pemindaian penuh. Ditambahkan.
>
> **Tiga hal dinyatakan jatuh ke aplikasi**, bukan diam-diam tidak ditegakkan: `UNIQUE` kombinasi pengenal lama (menunggu REQ-018), keutuhan rujukan `TABEL_TUJUAN`/`ID_TUJUAN` (polimorfik — Oracle tidak punya FK polimorfik, dan memalsukannya dengan 22 kolom nullable lebih buruk daripada tidak punya), dan umur baris.
>
> **Istilah `_Avoid_` di nama kolom dijelaskan di badannya**, bukan dibiarkan tanpa alasan: `CASEID` dan `pyID` dipakai sebagai **nama lama bersufiks `_LAMA`**, di tabel yang seluruh tugasnya memang menyimpan nama lama — pengecualian yang sama dengan lapisan view kompatibilitas (ADR-D-CNP-0021). Tidak satu pun muncul di tabel kanonik.
>
> **REQ-018 tetap menunggui, tidak menahan.** Yang berubah bila jawabannya lain: **constraint, bukan kolom** — dan itu sebabnya `UNIQUE`-nya sengaja belum dipasang. Menebaknya sekarang menggugurkan migrasi di tengah jalan, tepat ketika ia paling mahal dibatalkan.

*Asal: `T-23` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Setiap baris kanonik dapat ditunjuk balik lewat pengenal lamanya, tanpa mengurai teks.

`MIGRASI_KORELASI` dengan **lima kolom terpisah**, `PK_MIGRASI_KORELASI`, `CK_MIGRASI_KORELASI_1`, dan tiga index — `IX_MIGRASI_KORELASI_1` (baru → lama), `_2` dan `_3` (lama → baru). `ddl-usulan/19_MIGRASI_KORELASI.sql`.

**Tidak termasuk:** `UNIQUE` atas kombinasi pengenal lama — **ditahan, dengan alasan**: kolom mana yang ada sudah EVIDENCED; **kombinasi mana yang unik** menunggu REQ-018.

**Blocked by:**

- **REQ-018** — **menunggui, tidak menahan** (dari luar papan): `UNIQUE` atas kombinasi pengenal lama sudah **di luar lingkup tiket ini**; kolomnya EVIDENCED, kombinasinya menyusul
- ~~`01`~~ *(selesai)* — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis


**Dasar:** DECIDED(ADR-D-CNP-0005, ADR-D-CNP-0023). EVIDENCED: kelima pengenal lama — `PZINSKEY` dan `PYID` dari DDL tabel work, `CASEID` dari `OS_AKSEPTASI_KLAIM` dan §13.2, `IndexObject` dari §8.4.

- [x] Kelima pengenal terbaca sebagai kolom tersendiri dan dapat di-`JOIN` tanpa penguraian teks.
- [x] Baris yang bukan berasal dari work object Pega **diterima** dengan `PZINSKEY_LAMA` kosong — kosong berarti bukan dari sana, bukan tidak diketahui.
- [x] Baris tanpa **satu pun** pengenal lama **ditolak** — `CK_MIGRASI_KORELASI_1`.
- [x] Pencarian **dari pengenal lama** tertopang index, bukan pemindaian penuh — `IX_MIGRASI_KORELASI_2`, `_3`.

**Ketidakpastian:** **REQ-018**, BLOCKER, OPEN. Yang berubah bila jawabannya lain: **constraint**, bukan kolom.

## Claim Non Prop - 16 - Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs

---
status: selesai
---




> **SELESAI 19 September 2026 — dan dua cacat ditemukan saat mengerjakannya, bukan saat meninjaunya.**
>
> ### Cacat 1 — arah constraint terbalik
>
> Bentuk lama, dipasang di **23 constraint pada 10 tabel**:
>
> ```sql
> CHECK ((X_IDR IS NULL AND KURS IS NULL) OR (X_IDR IS NOT NULL AND KURS IS NOT NULL))
> ```
>
> Ia **menolak baris ber-`KURS` terisi tetapi `X_IDR` kosong** — dan itu persis kriteria penerimaan yang tiket `16` sudah tuliskan sejak awal: *"`KURS` terisi dan `*_IDR` kosong **diterima** — belum dikonversi, bukan salah."* Kriterianya benar; DDL-nya tidak menurutinya.
>
> ### Cacat 2 — saling mengunci, dan ini yang lebih berbahaya
>
> Lima tabel memasang **empat sampai lima** constraint berbentuk itu di atas **satu kolom `KURS` yang sama**. Akibatnya bukan penjumlahan, melainkan perkalian:
>
> > Begitu **satu** nilai IDR diisi, `KURS` menjadi `NOT NULL`. Cabang pertama setiap constraint lain gugur seketika — sehingga **keempat nilai IDR wajib terisi bersama**.
>
> Artinya: mustahil mencatat nilai kerugian dalam rupiah tanpa **sekaligus** mencatat rupiah untuk biaya penilaian, salvage, dan biaya lain. Padahal ADR-D-CNP-0029 baru saja menetapkan kesembilan nama yang hari ini tanpa padanan IDR **boleh berkolom kosong** — dan bentuk ini melarangnya.
>
> **Tidak ada yang pernah memutuskan itu.** Ia akibat bentuk, bukan niat. Dan ia tidak akan ketahuan sampai baris pertama dimasukkan.
>
> ### Yang berlaku sekarang — satu arah
>
> ```sql
> CHECK (X_IDR IS NULL OR KURS IS NOT NULL)
> ```
>
> *"Nilai rupiah tidak dapat lahir tanpa kurs"* (ADR-D-CNP-0029) — **dan tidak lebih dari itu.** Diperbaiki di 23 constraint, 10 berkas.
>
> ### Yang bertambah — asal-usul kurs, akhirnya ditegakkan
>
> ADR-D-CNP-0029 menuntut **lima** hal per nilai uang: nilai asli, kode mata uang, nilai rupiah, kurs yang dipakai, dan **asal-usul kurs itu**. Empat yang pertama sudah ditegakkan; yang kelima tidak pernah — `KURS_SUMBER` dan `KURS_TANGGAL` ada sebagai kolom, tetapi tidak ada yang mewajibkannya.
>
> Itu bukan kerapian. **FINDING-006**: `GETCURRENCYSTANDARD` mengembalikan `1` ketika kurs **tidak ditemukan**, dan angka `1` tidak dapat dibedakan dari kurs yang sah. Yang membedakannya **hanya asal-usulnya**. Kurs tanpa asal-usul mewarisi cacat itu utuh ke sistem baru.
>
> Ditambahkan di **10 tabel**: `KURS IS NULL OR (KURS_SUMBER IS NOT NULL AND KURS_TANGGAL IS NOT NULL)`.

*Asal: `T-04` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Nilai per mata uang tersimpan; nilai IDR tanpa kurs ditolak basis data.

`NILAI_KLAIM_MATA_UANG`, empat besaran uang berpasangan IDR, `UQ_NILAI_KLAIM_MATA_UANG_1`, empat `CK` pasangan IDR–kurs, `CK` domain keadaan.

**Tidak termasuk:** Penularan keadaan ke baris turunan — diuji di T-28 setelah tabel turunannya ada.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-D-CNP-0007, ADR-D-CNP-0014, ADR-D-CNP-0019). EVIDENCED: FINDING-006 — `GETCURRENCYSTANDARD` mengembalikan `1` bila kurs tidak ditemukan, dan `1` tidak dapat dibedakan dari kurs yang sah.

- [x] Baris dengan `NILAI_KERUGIAN_IDR` terisi sementara `KURS` kosong **ditolak** — berlaku untuk keempat besaran.
- [x] Baris dengan `KURS` terisi dan `*_IDR` kosong **diterima** — belum dikonversi, bukan salah. **Sebelumnya ditolak** oleh bentuk constraint lama; itulah cacat 1.
- [x] Baris kedua dengan (klaim, mata uang) sama **ditolak** — `UQ_NILAI_KLAIM_MATA_UANG_1`.
- [x] Nilai nol pada `*_IDR` **diterima** dan dapat dibedakan dari `NULL` lewat kueri — kolomnya nullable, tidak ada default.
- [x] Keempat nilai IDR **tidak saling mengunci**: mengisi satu tidak memaksa tiga lainnya. **Sebelumnya memaksa**; itulah cacat 2.
- [x] `KURS` terisi tanpa `KURS_SUMBER` **ditolak** — ADR-D-CNP-0029 butir kelima, FINDING-006.

**Ketidakpastian:** Tidak ada REQ yang menyentuhnya. Dua cacat yang ditemukan hari ini **bukan ketidakpastian** — keduanya salah tulis kami sendiri, dan keduanya sudah diperbaiki di tempatnya.

## Claim Non Prop - 17 - Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang

---
status: selesai
---



> **DILEPAS DARI KETERGANTUNGAN KOMITE, 19 September 2026.** Folder `Komite Claim Non Prop` dinyatakan **tertutup untuk batch ini** — bukan ditunda. F1/F2 karena itu tinggal permanen sebagai hipotesis sejajar, dan tiket ini **tidak boleh mati menunggunya**.
>
> **Dikerjakan di atas aturan yang sama-sama benar di kedua cabang:**
>
> | Cabang | Tuntutannya | Di tiket ini |
> |---|---|---|
> | F1 — pembacanya di Komite | kopling lintas modul masuk Boundary Contract | entri ditulis **bertanda belum terselesaikan**, arah dan pemilik **dikosongkan**, bukan ditebak |
> | F2 — tidak ada pembaca | sistem baru punya kunci hasil suntingan yang sebenarnya | **dibangun tanpa syarat** — ADR-D-CNP-0008 menuntutnya terlepas dari apakah sistem lama pernah punya jalurnya |
>
> Pemisahan hitung/sunting karena itu **tetap dipasang**, dan daftar kandidatnya sudah berbukti: sepuluh kolom yang ditulis mesin alokasi **dan** salah satu dari `AdjClaimCNP_Act`, `AdjClaimAmount_Act`, `SetActualPremium_ACT` (E15, S15).
>
> **26 penanda yatim golongan 3b** mengikuti kebijakan E17: **dimigrasi apa adanya dan ditandai tak berpemilik** — tidak dibuang, tidak diberi perilaku.


> **SELESAI 19 September 2026.**
>
> **Kriteria kelima tidak pernah punya constraint — sampai hari ini.** *"`DISUNTING_OLEH` kosong sementara `_SUNTING` terisi **ditolak** — suntingan selalu punya pelaku"* tertulis sejak tiket ini dibuat, dan **tidak ada satu pun `CHECK`** yang menegakkannya di `ALOKASI_LAYER` maupun `RETENSI_CEDANT`.
>
> Sebabnya ADR-D-CNP-0008: suntingan manual bertahan **dan terlihat**. Nilai yang menimpa hasil hitung tanpa ada yang bertanggung jawab atasnya bertahan **tanpa terlihat** — dan itu setengah dari yang dijanjikan, yaitu bagian yang justru lebih mudah disalahgunakan.
>
> Dipasang: `CK_ALOKASI_LAYER_4` (lima besaran, termasuk `ALASAN_SUNTING`) dan `CK_RETENSI_CEDANT_4`. `DISUNTING_ATAS_NAMA` tetap boleh kosong — kosong berarti tidak ada perwakilan, bukan tidak diketahui (ADR-D-CNP-0019).
>
> **Kedua tabel juga terkena dua cacat pasangan IDR–kurs** yang ditemukan di tiket `16`; keduanya sudah diperbaiki, dan asal-usul kurs kini ditegakkan di sini juga.
>
> **Empat kriteria lainnya sudah terpenuhi oleh bentuk yang ada**: `_DIPAKAI` adalah kolom `GENERATED ALWAYS AS (COALESCE(_SUNTING, _HITUNG)) VIRTUAL` — hitung ulang menimpa `_HITUNG` tanpa menyentuh `_SUNTING`, dan menulis langsung ke `_DIPAKAI` ditolak Oracle sendiri, bukan oleh aturan yang harus diingat orang.

*Asal: `T-05` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Hitung ulang menimpa `_HITUNG` dan tidak menyentuh `_SUNTING`; `_DIPAKAI` mengikuti tanpa ditulis siapa pun.

`ALOKASI_LAYER` dan `RETENSI_CEDANT`; kolom berpasangan `_HITUNG`/`_SUNTING`/`_DIPAKAI`; parameter yang di-snapshot; `UQ_ALOKASI_LAYER_1`, `UQ_RETENSI_CEDANT_1`.

**Tidak termasuk:** Kolom jenis reasuransi pada `RETENSI_CEDANT` — **sengaja tidak ada**. EVIDENCED: baris Retensi Cedant menulis `TreatyName = "UR"` secara literal melewati lookup (FINDING-003 bagian 2); `"UR"` artefak penyatuan yang dibatalkan ADR-D-CNP-0010, bukan nilai domain.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`08`~~ *(selesai)* — Sapuan S4: adakah padanan IDR untuk biaya penilaian, salvage, dan biaya lain
- ~~`09`~~ *(selesai)* — Sapuan S9: apa persisnya yang ditulis `EditXOLAlokasi`


**Dasar:** DECIDED(ADR-D-CNP-0010, ADR-D-CNP-0008, ADR-D-CNP-0011). EVIDENCED: `BLUEPRINT.md` §2.5 — `.IsEditClaim` **ditulis lima kali, tidak pernah dibaca**; FINDING-004.

- [x] Menjalankan ulang pengisian `_HITUNG` atas baris yang `_SUNTING`-nya terisi **tidak mengubah** nilai yang dibaca lewat `_DIPAKAI`.
- [x] Menulis langsung ke `_DIPAKAI` **ditolak** — ia kolom turunan.
- [x] Baris kedua dengan (klaim, keempat field layer, mata uang) sama **ditolak** — `UQ_ALOKASI_LAYER_1`.
- [x] Baris Retensi Cedant kedua untuk (klaim, mata uang) sama **ditolak** — `UQ_RETENSI_CEDANT_1`.
- [x] `DISUNTING_OLEH` kosong sementara `_SUNTING` terisi **ditolak** — suntingan selalu punya pelaku. **Baru ditegakkan 19 Sep 2026**: `CK_ALOKASI_LAYER_4`, `CK_RETENSI_CEDANT_4`.

**Ketidakpastian:** **T-36** menentukan besaran mana yang benar-benar berpasangan `_HITUNG`/`_SUNTING`; memasangnya di kolom yang tidak pernah disunting adalah beban mati. **T-35** menentukan apakah Biaya Penilaian, Salvage, dan Biaya Lain perlu pasangan IDR di tabel ini. **REQ-033** (sebaran `LayerPart`) menentukan apakah kunci halus pernah melahirkan dua baris di tempat sistem lama melihat satu — itu mengubah **view** di T-16, bukan tabel ini.

## Claim Non Prop - 18 - Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran

---
status: selesai
menunggu-luar: [REQ-018]
---



> **DILEPAS DARI KETERGANTUNGAN KOMITE, 19 September 2026.** Folder Komite **tertutup untuk batch ini**. H1/H2 dan I1/I2 tinggal permanen sebagai hipotesis sejajar; tiket ini **tidak menunggunya**.
>
> **Dikerjakan di atas aturan yang sama-sama benar di kedua cabang:**
>
> - **I** — model **menyediakan tempat** bagi transisi `AcceptanceStatus` yang datang dari luar, **tanpa mengandaikan ia terjadi**. Keadaannya **empat, termasuk kosong**, sesuai temuan **C7**: properti ini dibandingkan sebagai angka (`0`,`1`,`2`) *dan* sebagai teks, dan bentuk teksnya memuat nilai kosong tanpa padanan numerik. Nilai kosong itu justru yang diuji `CloseClaimMD` (D2).
> - **H** — nilai Adjustment tetap dipasang **berpasangan mata uang** dan **tidak mengandaikan** salah satunya sudah rupiah.
>
> Tidak boleh ada rancangan yang mengandaikan salah satu cabang benar.


> **SELESAI 19 September 2026 — dan satu kolom harus dilonggarkan supaya keadaan keempat dapat disimpan.**
>
> `KEPUTUSAN_KOMITE` dipasang **`NOT NULL`**. Itu membuat **keadaan keempat mustahil disimpan**, dan keadaan keempat justru yang paling berakibat.
>
> S12 menemukan properti ini dibandingkan sebagai **angka** (`0`, `1`, `2`) **dan sebagai teks**, dan bentuk teksnya memuat **nilai kosong** yang tidak punya padanan numerik. Keadaannya **empat**, bukan tiga — itu sudah ditulis di C7 dan `SPEC-MODEL-DATA.md` §21.1, tetapi DDL-nya belum menurutinya.
>
> Dan yang kosong bukan keadaan sepele: penjaga `CloseClaimMD` menolak `AcceptanceStatus == "0"`, sementara Adjustment yang **belum pernah dikirim** ke komite bernilai **kosong**, bukan `"0"` — sehingga ia **lolos** (**D2**). Seluruh pembacaan jalur tutup-langsung bergantung pada perbedaan itu.
>
> Sekarang: `NULL` = **belum pernah dikirim ke komite**, berbeda dari `'0'` = **sudah dikirim, belum diputus**. ADR-D-CNP-0019 sebagai kolom, bukan sebagai kalimat.
>
> **Tanpa `CHECK` domain, dan itu tetap keputusan**: domainnya EXTERNAL — modul Komite yang menulisnya, dan folder itu tertutup. Domain tertutup di sini akan menolak nilai yang sah dari seberang.
>
> Kriteria lain sudah terpenuhi bentuk yang ada: `UQ_AKSEPTASI_1` (empat field layer + mata uang — **lebih ketat** daripada yang tiket sebut), `UQ_ADJUSTMENT_1`, `FK_ADJUSTMENT_2` menolak penghapusan rekening yang masih dirujuk, `CK_AKSEPTASI_1` domain tiga keadaan.

*Asal: `T-06` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Akseptasi kedua untuk klaim, layer, dan mata uang yang sama ditolak. Satu klaim dapat memiliki banyak Adjustment bernomor urut tetap.

`AKSEPTASI`, `ADJUSTMENT`, `REKENING_PENERIMA`; `UQ_AKSEPTASI_1`; `UQ_ADJUSTMENT_1`; `FK_ADJUSTMENT_2` beserta `IX_ADJUSTMENT_1`; `CK_AKSEPTASI_1`.

**Tidak termasuk:** `CHECK` domain untuk `KEPUTUSAN_KOMITE` — **sengaja tidak ada**, domainnya EXTERNAL. Tabel Komite dan keanggotaannya tidak dirancang.

**Blocked by:**

- **REQ-018** — **menunggui, tidak menahan** (dari luar papan): ia mengukur baris ganda pada **kunci lama yang lebih longgar**; kunci di tiket ini lebih ketat, jadi jawabannya batas bawah, bukan penentu
- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`07`~~ *(selesai)* — Sapuan S3: perilaku per nilai `PaymentType`


**Dasar:** DECIDED(ADR-D-CNP-0024, ADR-D-CNP-0015, ADR-D-CNP-0023). EVIDENCED: blok yang **dikomentari** di `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` memakai persis kombinasi itu; `OS_AKSEPTASI_KLAIM` **tanpa primary key maupun unique**, dan prosedurnya **selalu `INSERT`**.

- [x] `INSERT` kedua dengan (klaim, layer, mata uang) sama **ditolak** — `UQ_AKSEPTASI_1`, yang menegakkan **keempat** field layer, bukan satu.
- [x] Adjustment dengan nomor urut yang sudah dipakai pada klaim yang sama **ditolak** — `UQ_ADJUSTMENT_1`.
- [x] Menghapus rekening yang masih dirujuk Adjustment **ditolak** — `FK_ADJUSTMENT_2`.
- [x] `KEADAAN_AKSEPTASI` di luar tiga nilai domain **ditolak** — `CK_AKSEPTASI_1`.
- [x] `KEPUTUSAN_KOMITE` bernilai apa pun **diterima** — dan itu keputusan, bukan kelalaian.
- [x] `KEPUTUSAN_KOMITE` **kosong diterima**, dan dapat dibedakan dari `'0'`. **Sebelumnya `NOT NULL`** — keadaan keempat tidak dapat disimpan sama sekali.

**Ketidakpastian:** **I1/I2** — transisi status akseptasi `0 → 1/2` **TIDAK DITEMUKAN** di 279 berkas; ia ditulis satu kali saja sebagai `0`. **H1/H2** — apakah `.ValueAdjustment` sudah IDR di hulu belum diketahui, sehingga kolom nilai Adjustment dipasang berpasangan mata uang dan **tidak mengandaikan** salah satunya. **REQ-018** mengukur baris ganda pada kunci lama, yang lebih longgar dari kunci ini.

## Claim Non Prop - 19 - Masukan mesin alokasi dan penyebaran tersimpan, bukan hanya keluarannya

---
status: selesai
---




> **SELESAI 19 September 2026.**
>
> **Kriteria butir 1 DIUBAH, bukan sekadar diberi catatan penyimpangan.** Bunyi lama: *"nomor urut ganda dalam satu klaim ditolak, **di ketiga tabel**"*. Bunyi baru memisahkan **tabel yang punya kunci alami** dari yang tidak.
>
> Sebabnya: kriteria yang dibiarkan salah akan dibaca orang berikutnya sebagai **pekerjaan yang belum selesai**, dan ia akan menambahkan `NOMOR_URUT` ke `ESTIMASI_AWAL` justru untuk memenuhinya — melemahkan tabel demi memenuhi kalimat.
>
> `PEMBAGIAN_KERUGIAN` dan `PENYEBARAN` memakai `UQ (ID_KLAIM, NOMOR_URUT)`. **`ESTIMASI_AWAL` tidak punya `NOMOR_URUT` sama sekali**, dan kolomnya **tidak ditambahkan**.
>
> Sebabnya: tabel itu punya **kunci alami yang sesungguhnya** — `(klaim, jenis reasuransi, mata uang)`. Kunci alami **lebih kuat** daripada nomor urut: nomor urut hanya melarang dua baris **bernomor sama**; kunci alami melarang dua baris **berarti sama**.
>
> Menambahkan `NOMOR_URUT` di sana justru akan **melemahkan** tabelnya — ia mengizinkan estimasi kedua untuk jenis dan mata uang yang sama asal nomornya berbeda. Yang dijanjikan tiket, baris ganda ditolak, **terpenuhi dan terpenuhi lebih ketat**; yang tidak terpenuhi hanya **bentuk** kuncinya.
>
> **Butir 2 dilengkapi**: `IX_PENYEBARAN_1` atas `(ID_KLAIM, JENIS_REASURANSI, MATA_UANG)`. Kolomnya sudah ada di baris itu — itu yang membuat janji *"tanpa membaca tabel lain"* benar; index ini yang membuatnya murah.

*Asal: `T-08` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Pembagian kerugian, penyebaran, dan estimasi awal tersimpan sebagai data.

`PEMBAGIAN_KERUGIAN`, `PENYEBARAN`, `ESTIMASI_AWAL`.

**Tidak termasuk:** `SpreadingAdjustment` dan `SpreadingAdjustmentQS` — keduanya **view**, bukan tabel. → T-18.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.2 Langkah 2 — `CNPSpreadLoss` adalah **masukan** mesin alokasi. DECIDED(ADR-D-CNP-0011, ADR-D-CNP-0013): hitung ulang idempoten mensyaratkan seluruh masukan tersimpan. `CONTEXT.md` **Penyebaran**.

- [x] **Nomor urut ganda dalam satu klaim ditolak, pada tabel yang tidak punya kunci alami** — `PEMBAGIAN_KERUGIAN` dan `PENYEBARAN`, lewat `UQ (ID_KLAIM, NOMOR_URUT)`.
- [x] **Pada tabel yang punya kunci alami, kunci itu yang ditegakkan, bukan nomor urut** — `ESTIMASI_AWAL`, lewat `UQ_ESTIMASI_AWAL_1 (ID_KLAIM, JENIS_REASURANSI, MATA_UANG)`.
- [x] Baris penyebaran dapat dijumlahkan per jenis reasuransi per mata uang tanpa membaca tabel lain — kolomnya di baris itu, `IX_PENYEBARAN_1` menopangnya.

**Ketidakpastian:** Tidak ada REQ yang menyentuhnya.

## Claim Non Prop - 20 - Objek, kronologi, dan dokumen klaim

---
status: selesai
---




> **SELESAI 19 September 2026 — tanpa perubahan struktur. Ketiga kriteria sudah terpenuhi bentuk yang ada.**
>
> Ditulis begini supaya penutupan ini dapat diperiksa, bukan dipercaya:
>
> | Kriteria | Yang menegakkannya |
> |---|---|
> | URL terisi tanpa kedaluwarsa ditolak, dan sebaliknya | `CK_DOKUMEN_KLAIM_1` — **dua arah, dan di sini dua arah memang benar**: URL tanpa masa berlaku tidak dapat dipakai, masa berlaku tanpa URL tidak menyatakan apa pun. Berbeda dari pasangan IDR–kurs, yang satu arah (tiket `16`) |
> | Dokumen tanpa URL sama sekali diterima | kedua kolom nullable; ADR-D-CNP-0027 — berkasnya di penyimpanan awan, URL diterbitkan saat diperlukan |
> | Kronologi tanpa jenis tindakan ditolak; tanpa keterangan diterima | `JENIS_TINDAKAN NOT NULL`, `KETERANGAN` nullable |
>
> `JENIS_TINDAKAN` **tanpa `CHECK` domain**, dan itu sejalan ADR-D-CNP-0009: yang terstruktur adalah **keputusannya**, dan jenis tindakan baru muncul seiring proses tanpa mengubah skema. `KETERANGAN` **tidak pernah dibaca mesin** — itu pokok ADR-D-CNP-0009, dan sebabnya terbaca di sistem lama: `@contains(.CommentSuggest,"Accepted by Himawan")` membuat satu salah ketik mengubah jalur (E7).
>
> **`ObjectList` tetap lubang tercatat, bukan tabel.** Apa yang membedakannya dari `InterestList` **TIDAK DITEMUKAN DI XML**, dan tabel yang dibuat atas dasar tebakan lebih buruk daripada lubang yang tercatat — `SPEC-MODEL-DATA.md` bagian 10.

*Asal: `T-09` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Objek pertanggungan, kronologi, dan rujukan dokumen tersimpan; URL tanpa masa berlaku ditolak.

`OBJEK_PERTANGGUNGAN`, `KRONOLOGI_KLAIM`, `DOKUMEN_KLAIM` beserta `CK_DOKUMEN_KLAIM_1`.

**Tidak termasuk:** `ObjectList` — **lubang tercatat**, bukan tabel. Apa yang membedakannya dari `InterestList` **TIDAK DITEMUKAN DI XML**.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-D-CNP-0027): berkas di Google Cloud Storage, yang tersimpan hanya metadata dan URL beserta `EXPDATE`. DECIDED(ADR-D-CNP-0009): kronologi terstruktur, keterangan tidak pernah dibaca mesin.

- [x] Dokumen dengan URL terisi tetapi tanpa kedaluwarsa **ditolak**, dan sebaliknya — `CK_DOKUMEN_KLAIM_1`.
- [x] Dokumen tanpa URL sama sekali **diterima** — URL belum pernah diterbitkan.
- [x] Kronologi tanpa jenis tindakan **ditolak**; tanpa keterangan **diterima**.

**Ketidakpastian:** Tidak ada REQ. Lubang `ObjectList` dicatat di `SPEC-MODEL-DATA.md` bagian 10.

## Claim Non Prop - 21 - Tanggal tutup buku dan tarif menjadi data bertanggal berlaku

---
status: selesai
---




> **SELESAI 19 September 2026 — dan baris awalnya menghasilkan berkas DML pertama di proyek ini.**
>
> ⚠ **`ddl-usulan/Z00_ISIAN_AWAL.sql` berisi `INSERT`. Ia satu-satunya berkas DML di seluruh `ddl-usulan/`**, dan dipisahkan justru supaya kekecualian itu **terlihat** — bukan terselip di kaki sebuah berkas DDL tempat tidak ada yang mencarinya. Ke-22 berkas tabel, kesembilan view, dan berkas skema semuanya tetap menyatakan *"TIDAK ADA DML DI BERKAS INI"*, dan pernyataan itu tetap benar.
>
> Empat baris: tutup buku `25` berlingkup global; brokerage 2,5%; PPh 2%; PPN 2,2%.
>
> **Kenapa baris ini di berkas, bukan diketik orang saat pemasangan.** Karena keempat nilai itu **sekarang tertanam di kode** sistem lama, dan seluruh maksud ADR-D-CNP-0025 adalah memindahkannya dari kode ke data. Baris yang diketik orang tidak dapat ditinjau, tidak dapat dibandingkan, dan tidak meninggalkan jejak siapa yang memilih angkanya.
>
> Angka `25` punya **dua sumber kebenaran** di sistem lama — tertanam di `HitServiceToKasir_Act`, **dan** dibaca `PROC_GENERATE_SEQUENCE_NUMBER` dari `POOLDATA.TANGGAL_CLOSING`. Dua sumber untuk satu aturan adalah cacat yang tabel ini tutup, **tetapi hanya bila tabelnya benar-benar terisi**.
>
> `BERLAKU_SEJAK = 1900-01-01` bukan tanggal yang berarti — ia menyatakan *"sejak sebelum data tertua"*, sehingga tidak ada baris lama yang jatuh di luar masa berlaku tarif mana pun. `DIBUAT_OLEH = 'MIGRASI'` karena **tidak ada orang yang memutuskannya**: nilainya dibaca dari kode lama apa adanya, dan menuliskan nama seseorang akan mengarang pelaku untuk keputusan yang tidak pernah diambil siapa pun.
>
> **Ratifikasi tidak menahan.** `ASK-AKUNTANSI` butir 5 memverifikasi apakah ketiga tarif masih berlaku. Bila ternyata tidak, yang berubah **isi baris** — pekerjaan entri, bukan pekerjaan kode. Itu seluruh maksud ADR-D-CNP-0025.

*Asal: `T-10` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Satu sumber kebenaran untuk tanggal tutup buku dan tarif; keduanya terisi baris awalnya.

`TUTUP_BUKU`, `TARIF_BERLAKU`; baris awal: tutup buku `25` berlingkup global; tarif brokerage 2,5%, PPh 2%, PPN 2,2%, seluruhnya berlingkup global dan berlaku sejak sebelum data tertua.

**Tidak termasuk:** Riwayat tarif — **tidak ada bukti tarif pernah berubah** (AK-5), jadi tidak ada riwayat yang dibuat-buat. Faktor `102,2` **tidak disimpan**; ia diturunkan sebagai `(100 + tarif PPN)`.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-D-CNP-0025). DECIDED-TEKNIS(AK-5, AK-6.3). EVIDENCED: angka `25` tertanam di `HitServiceToKasir_Act`; `PROC_GENERATE_SEQUENCE_NUMBER` membacanya dari `POOLDATA.TANGGAL_CLOSING` — **dua sumber kebenaran untuk satu aturan**.

- [x] Dua baris dengan (lingkup, tanggal berlaku) sama **ditolak** — `UQ_TUTUP_BUKU_1`, `UQ_TARIF_BERLAKU_1`.
- [x] Tanggal tutup buku di luar 1–31 **ditolak** — `CK_TUTUP_BUKU_1`.
- [x] Tarif terbaca lewat satu kueri tanpa membaca kode mana pun.
- [x] Keduanya **terisi baris awalnya** — `ddl-usulan/Z00_ISIAN_AWAL.sql`, 1 baris tutup buku dan 3 baris tarif.

**Ketidakpastian:** **A8** dan **ASK-AKUNTANSI no. 5** — apakah ketiga tarif masih berlaku **belum dikonfirmasi**. Yang dibuat tempatnya; nilainya berlabel DECIDED-TEKNIS menunggu ratifikasi, dan **ratifikasi tidak menahan tiket ini**.

## Claim Non Prop - 22 - Koreksi bernilai tercatat menggantikan tambalan di dalam kode

---
status: selesai
---




> **SELESAI 19 September 2026.**
>
> **Butir 4 tidak terpenuhi oleh index yang ada.** Ia menjanjikan *"riwayat koreksi atas **satu kolom tertentu** terbaca lewat satu kueri"*, sementara `IX_KOREKSI_NILAI_1` hanya mencakup `(TABEL_SASARAN, ID_SASARAN)` — tanpa `KOLOM_SASARAN`. Janji itu berjalan lewat pemindaian seluruh riwayat baris. Index diperluas menjadi `(TABEL_SASARAN, ID_SASARAN, KOLOM_SASARAN, DIBUAT_PADA)`; `DIBUAT_PADA` ikut karena riwayat dibaca berurut waktu.
>
> **Satu constraint bertambah**: `CK_KOREKSI_NILAI_1` menolak baris yang `NILAI_SEBELUM` **dan** `NILAI_SESUDAH`-nya kosong. Keduanya boleh kosong sendiri-sendiri — `NILAI_SEBELUM` kosong berarti kolomnya belum pernah terisi, `NILAI_SESUDAH` kosong berarti koreksi ini **membatalkan** nilai (butir 2) — tetapi kosong keduanya **tidak menyatakan apa pun**. Sebuah koreksi yang tidak mengubah apa pun bukan koreksi.
>
> Kriteria lain terpenuhi bentuk yang ada: `ALASAN NOT NULL`, `DIBUAT_OLEH NOT NULL`, dan `UQ_KLAIM_PENJAGA_TANGGAL_1` atas `ID_KLAIM` — satu klaim muncul sekali.

*Asal: `T-11` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Setiap penyimpangan nilai punya nilai sebelum, nilai sesudah, alasan, pelaku, dan waktu.

`KOREKSI_NILAI`, `KLAIM_PENJAGA_TANGGAL`, `IX_KOREKSI_NILAI_1`.

**Tidak termasuk:** Pengisian daftar klaim terdampak — menunggu migrasi.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-D-CNP-0004, ADR-D-CNP-0018). EVIDENCED: `BLUEPRINT.md` §7 — **29 langkah, 18 ekspresi unik, 8 rule** menambal per-case; tambalan terbaru `CLMNP-975` bertanggal **2026-07-16**.

- [x] Koreksi tanpa alasan **ditolak**; tanpa pelaku **ditolak** — keduanya `NOT NULL`.
- [x] Koreksi dengan `NILAI_SESUDAH` kosong **diterima** — koreksi yang membatalkan nilai.
- [x] Koreksi dengan **kedua** nilai kosong **ditolak** — `CK_KOREKSI_NILAI_1`, baru.
- [x] Satu klaim hanya dapat muncul sekali di daftar klaim penjaga tanggal — `UQ_KLAIM_PENJAGA_TANGGAL_1`.
- [x] Riwayat koreksi atas satu kolom tertentu terbaca lewat satu kueri — `IX_KOREKSI_NILAI_1` **diperluas** dengan `KOLOM_SASARAN` dan `DIBUAT_PADA`.

**Ketidakpastian:** Tidak ada REQ. Bila `KOREKSI_NILAI` dan kolom suntingan di tingkat baris berbeda, **`KOREKSI_NILAI` yang berlaku** — `SPEC-MODEL-DATA.md` bagian 20.

## Claim Non Prop - 23 - Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama

---
status: selesai
---




> **SELESAI 19 September 2026 — view-nya ditulis ulang dari satu tabel menjadi sembilan belas.**
>
> Versi sebelumnya hanya menggabungkan `AKSEPTASI`. Akibatnya dua, dan **keduanya diam**:
>
> 1. Baris korelasi untuk **18 tabel lain** tetap muncul, tetapi `KEADAAN_BARIS`-nya `NULL` — bukan karena barisnya belum lengkap, melainkan **karena view-nya tidak melihat tabelnya**.
> 2. Penyaring `KEADAAN_BARIS = 'LENGKAP'` yang dijanjikan butir 2 akan **membuang seluruh 18 tabel itu tanpa sepatah pesan**.
>
> Itu persis cacat yang ADR-D-CNP-0019 larang: **`NULL` diperlakukan sebagai nilai** — dan di alat paritas, baris yang hilang adalah **selisih yang tidak pernah terhitung**.
>
> **Tiga keadaan dibedakan sekarang, dan tidak satu pun `NULL`:**
>
> | Nilai | Artinya |
> |---|---|
> | `LENGKAP` / `MENUNGGU_KURS` / `GAGAL_URAI` | tabelnya punya kolom keadaan, dan inilah isinya — 11 tabel |
> | `TIDAK_BERLAKU` | tabelnya **tidak** punya kolom keadaan — 8 tabel acuan dan catatan. Berbeda dari "belum lengkap" |
> | `TABEL_TIDAK_DIKENALI` | baris korelasi menunjuk tabel di luar daftar. **Tanpa cabang ini ia hilang dari view** |
>
> Kolom `KEBERADAAN` memisahkan pertanyaan kedua dari yang pertama: `'HILANG'` berarti baris kanoniknya tidak ada, apa pun keadaannya — itu **kegagalan migrasi**, bukan baris yang belum lengkap.
>
> **Tiga tabel sengaja di luar**: `MIGRASI_KORELASI`, `MIGRASI_PENDARATAN`, `MIGRASI_NILAI_DITOLAK`. Ketiganya jembatan, bukan sasaran — baris korelasi tidak pernah menunjuk tabel korelasi. Bila ternyata ada, cabang terakhir menangkapnya, dan itu memang yang seharusnya terjadi.

*Asal: `T-19` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Perbandingan baris per baris mungkin dilakukan.

`V_PARITAS_SHADOW`.

**Tidak termasuk:** Pembandingnya sendiri dan toleransinya — AK-3.

**Blocked by:**

- ~~`15`~~ *(selesai)* — Jembatan ke sistem lama berdiri sebagai tabel terpisah


**Dasar:** DECIDED(ADR-D-CNP-0005). EVIDENCED: `BLUEPRINT.md` §8.4 — `IndexObject` adalah **posisi numerik, bukan surrogate key**, sehingga jembatannya hilang begitu baris berpindah ke kunci sendiri.

- [x] Setiap baris korelasi muncul **tepat sekali** — cabangnya disaring `TABEL_TUJUAN`, ke-19 nilainya saling lepas, dan cabang terakhir menampung sisanya.
- [x] Baris ber-`KEADAAN_BARIS` bukan `LENGKAP` dapat **dikeluarkan** lewat satu penyaring, bukan dihitung sebagai selisih — dan **tanpa ikut membuang 18 tabel** yang tidak punya kolom keadaan.
- [x] Baris kanonik yang **hilang** terbaca sebagai `KEBERADAAN = 'HILANG'`, terpisah dari soal keadaan.

**Ketidakpastian:** Baseline pembandingnya sendiri belum tentu benar: `TotalUR` selalu nol dan dua rumus premi pemulihan bekerja atas nilai berbeda. AK-2b menempatkan ketiganya sebagai pengecualian bernama.

## Claim Non Prop - 24 - Uji: batas tanggal inklusif

---
status: menunggu-instance
---



> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-26` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa klaim yang jatuh tepat di hari terakhir masa berlaku treaty tidak tertolak.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-D-CNP-0022). EVIDENCED: FINDING-005 — di sistem lama perbandingannya dilakukan antara dua format teks yang berbeda.

- [ ] Klaim ber-Tanggal Kejadian **tepat di hari terakhir** masa berlaku treaty **diterima**; sehari sesudahnya **ditolak**.

**Ketidakpastian:** **REQ-020** mengukur berapa klaim lama yang terdampak; tidak menahan uji ini.

## Claim Non Prop - 25 - Arsip muatan keluar: apa yang benar-benar dikirim, tersimpan sebagai catatan

---
status: selesai
---




> **SELESAI 19 September 2026 — dan grant-nya bertentangan dengan tiketnya sendiri.**
>
> Butir 1 berbunyi *"baris arsip **tidak dapat** di-`UPDATE` maupun di-`DELETE` oleh akun aplikasi — tulis-sekali ditegakkan hak akses, bukan kesepakatan"*. Tiket `35` juga sudah menyebut pengecualian ini sejak ditulis. **DDL-nya memberi keempat hak**: `GRANT SELECT, INSERT, UPDATE, DELETE`.
>
> Diperbaiki: **`GRANT SELECT, INSERT`** saja.
>
> Sebabnya bukan kehati-hatian. Arsip ini menjawab pertanyaan **"apa yang benar-benar dikirim"**. Arsip yang dapat diubah menjawab pertanyaan lain — *"apa yang sekarang kita katakan telah dikirim"* — dan itu bukan catatan.
>
> Dan ia **bukan turunan yang dapat dihitung ulang**: `GetDataOS` menjumlahkan seluruh baris berkunci sama ber-`STS_REJECT = 0`, lalu `SaveDataToOSAksep_Act` mengurangkan hasilnya. **Baris adalah tambahan, bukan keadaan.** Satu `UPDATE` mengubah setiap selisih yang pernah dihitung sesudahnya.
>
> `IX_ARSIP_MUATAN_KELUAR_1` diperluas dengan **`DITOLAK`**, karena setiap penjumlahan selisih menyaringnya (butir 3).
>
> Butir 4 sudah terpenuhi: kelima besaran `NUMBER(38,2)`. **Ini satu-satunya tempat pembulatan tepi menjadi baris tersimpan** — di mana pun lagi skala kanonik 20 berlaku (ADR-D-CNP-0003).

*Asal: `T-39` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Setiap muatan yang dikirim ke Arasapas dan kasir tercatat apa adanya, dan muatan berikutnya dapat dihitung sebagai selisih terhadapnya.

`ARSIP_MUATAN_KELUAR`: tujuan, kunci (klaim, jenis reasuransi, mata uang), lima besaran uang **berskala 2**, kurs, penanda ditolak, pelaku, waktu kirim. `IX_ARSIP_MUATAN_KELUAR_1`, `CK_ARSIP_MUATAN_KELUAR_1`.

**Tidak termasuk:** Proses pengirimannya — lapisan aplikasi. Muatan dalam bentuk JSON — **tidak ada kolom JSON di sini**; kolomnya bertipe, karena daftar kolomnya sudah terbaca dari S1 (ADR-D-CNP-0028).

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED-TEKNIS(AK-1) — *"bila muatan keluar perlu diarsipkan, arsipnya tabel tersendiri, bukan kolom di tabel bisnis."* Ini pemakaian pertamanya. EVIDENCED: `RDBList\GetDataOS.xml` menjumlahkan seluruh baris berkunci `(CASEID, TypeLoss, Currency)` dengan `STS_REJECT = 0`, dan `SaveDataToOSAksep_Act` mengurangkan hasilnya — **baris adalah tambahan, bukan keadaan**.

- [x] Baris arsip tidak dapat di-`UPDATE` maupun di-`DELETE` oleh akun aplikasi — `GRANT SELECT, INSERT` saja. **Sebelumnya keempat hak diberikan.**
- [x] Dua muatan atas kunci yang sama pada waktu berbeda **keduanya tersimpan** — tidak ada `UNIQUE` atas kuncinya, dan itu disengaja.
- [x] Muatan bertanda ditolak **tidak** ikut terjumlah — `DITOLAK` ada di kolom **dan** di `IX_ARSIP_MUATAN_KELUAR_1`.
- [x] Nilai tersimpan **berskala 2** — `NUMBER(38,2)` pada kelima besaran.

**Ketidakpastian:** Tiga hal yang ditulis terang supaya tidak dibaca sebagai pelanggaran ADR: view kompatibilitas **membaca** arsip sehingga ia tetap view, bukan salinan (ADR-D-CNP-0023 utuh); yang menulis arsip adalah proses pengiriman lewat akun aplikasi (ADR-D-CNP-0017 utuh); dan arsip berisi **apa yang benar-benar dikirim**, bukan apa yang seharusnya — ia catatan, bukan turunan yang dapat dihitung ulang.

## Claim Non Prop - 26 - Premi pemulihan dapat dihitung ulang dari barisnya sendiri

---
status: selesai
---




> **SELESAI 19 September 2026 — tanpa perubahan struktur.** Ketiga kriteria sudah terpenuhi bentuk yang ada; yang bertambah **catatan di badan DDL-nya**, dan itu perlu.
>
> Sistem lama punya **dua** rumus premi pemulihan yang bekerja atas nilai **berbeda** (FINDING-007). Yang diwarisi `CountReinstatement_Act`; yang **tidak** diwarisi `AdjClaimCNP_Act` baris 3109 (AK-2b). Itu tidak boleh tinggal di dokumen saja: kolom tabel ini adalah masukan **rumus pertama**, dan bila kelak ada yang mengisinya dari rumus kedua, **angkanya akan tampak sah dan hasilnya berbeda tanpa satu pun tanda**. Sekarang tertulis di kepala berkasnya.
>
> `LIMIT_LAYER` adalah **penyebut**. Nol di penyebut bukan nilai yang aneh — ia perhitungan yang **tidak dapat dijalankan**. `CK_PREMI_PEMULIHAN_1` menolaknya di baris, bukan menyerahkannya ke kode yang harus ingat memeriksanya.

*Asal: `T-07` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Setiap masukan rumus tersimpan bersama hasilnya; limit layer nol ditolak.

`PREMI_PEMULIHAN` beserta seluruh masukan rumus; `CK_PREMI_PEMULIHAN_1`; `UQ_PREMI_PEMULIHAN_1`.

**Tidak termasuk:** Perhitungannya sendiri — ini lapisan data.

**Blocked by:**

- ~~`17`~~ *(selesai)* — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang


**Dasar:** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.3 — `CNPLimit` adalah **penyebut**. DECIDED-TEKNIS(AK-2b): acuannya `CountReinstatement_Act`; `AdjClaimCNP_Act` baris 3109 **tidak diwarisi** (FINDING-007).

- [x] Baris ber-`LIMIT_LAYER` nol **ditolak** — `CK_PREMI_PEMULIHAN_1`.
- [x] Baris kedua untuk (klaim, layer, mata uang) sama **ditolak** — `UQ_PREMI_PEMULIHAN_1`, yang menegakkan **keempat** field layer.
- [x] Seluruh masukan rumus dapat dibaca dari satu baris, tanpa menyentuh tabel master treaty — tujuh kolom, empat di antaranya `NOT NULL` karena rumus tanpa salah satunya tidak dapat dijalankan sama sekali.

**Ketidakpastian:** AK-2b menempatkan selisih terhadap sistem lama sebagai **pengecualian bernama** pada shadow-run, bukan kegagalan cutover. Besaran selisihnya belum terukur.

## Claim Non Prop - 27 - Domain kolom keadaan

---
status: selesai
---



> **SELESAI 2026-09-18.** RANCANGAN selesai — CHECK sudah ada di DDL. Pemasangannya ikut tiket tabelnya masing-masing; tidak ada pekerjaan tersendiri di sini.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-12` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Kolom keadaan yang domainnya kita tetapkan sendiri terikat `CHECK`; yang domainnya warisan dan belum lengkap **sengaja tidak terikat**.

`CK` `KEADAAN_BARIS` di **11 dari 11** tabel yang punya kolom itu; `CK` `KEADAAN_AKSEPTASI`.

**Tidak termasuk:** `STATUS_KLAIM` dan `KEPUTUSAN_KOMITE` — **tanpa `CHECK`, dan itu keputusan**. Domain tertutup yang isinya belum lengkap akan menolak nilai yang sah pada hari pertama modul Komite tersambung.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`22`~~ *(selesai)* — Koreksi bernilai tercatat menggantikan tambalan di dalam kode


**Dasar:** DECIDED(ADR-D-CNP-0019 lapis 2). EVIDENCED: sapuan S2 empat lapisan — `CNPStatusCase` hanya **2 nilai ditulis**, **nol rule menguji**, dan `"CLAIM ACCEPTED"`/`"CLAIM REJECTED"` **nihil di 279 berkas**.

- [ ] Nilai `KEADAAN_BARIS` di luar `LENGKAP` · `MENUNGGU_KURS` · `GAGAL_URAI` **ditolak**, di kesebelas tabel.
- [ ] `STATUS_KLAIM` bernilai apa pun **diterima**, dan alasannya tercatat.

**Ketidakpastian:** Kelengkapan enum `CNPStatusCase` hanya dapat dibuktikan dengan membuka folder `Komite Claim Non Prop` — **belum diizinkan**.

## Claim Non Prop - 28 - Pasangan IDR–kurs

---
status: selesai
---



> **SELESAI 2026-09-18.** RANCANGAN selesai — 23 CHECK pasangan sudah ada di DDL, dipasang pembangkit. Pemasangannya ikut tiket tabelnya.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-13` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada nilai IDR yang dapat lahir tanpa kurs yang menghasilkannya, di tabel mana pun.

**23 `CHECK`** pasangan, dipasang **oleh pembangkit** atas setiap kolom `*_IDR` di tabel yang punya `KURS`.

**Tidak termasuk:** Penularan keadaan ke baris turunan — itu aturan prosedur, diuji di T-28.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`22`~~ *(selesai)* — Koreksi bernilai tercatat menggantikan tambalan di dalam kode


**Dasar:** DECIDED(ADR-D-CNP-0007, ADR-D-CNP-0014). EVIDENCED: FINDING-006 — `RETURN 1` pada kurs yang tidak ditemukan.

- [ ] Di **setiap** tabel bernilai uang, `*_IDR` terisi tanpa `KURS` **ditolak**.
- [ ] Menambahkan tabel bernilai uang baru tanpa `CHECK`-nya **tidak mungkin** — pembangkit memasangnya sendiri.

**Ketidakpastian:** **T-35** dapat menambah kolom `*_IDR` di `ALOKASI_LAYER`; bila itu terjadi, `CHECK`-nya ikut terpasang tanpa perubahan tangan.

## Claim Non Prop - 29 - Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam

---
status: selesai
menunggu-luar: [REQ-033]
---




> **SELESAI 19 September 2026 — dan view-nya belum menghasilkan selisih sama sekali.**
>
> ### Yang hilang: seluruh pokok tiket ini
>
> Tiket ini berbunyi *"keduanya membaca tabel kanonik **dan arsip muatan keluar**, lalu menghasilkan **selisih**"*. `V_AKSEPTASI_KOMPATIBEL` **tidak menyentuh `ARSIP_MUATAN_KELUAR` sama sekali** — ia mengeluarkan **nilai mutlak**.
>
> Itu bukan kekurangan kecil. Bukti di tiket ini sendiri menunjukkan sistem hilir **tidak menerima keadaan, ia menerima tambahan**: `GetDataOS` menjumlahkan seluruh baris berkunci sama ber-`STS_REJECT = 0`, lalu `SaveDataToOSAksep_Act` **mengurangkan** hasilnya; prosedur di seberang **selalu `INSERT`**.
>
> > **Mengeluarkan nilai mutlak ke jalur yang menjumlahkan tambahan berarti melipatgandakan setiap nilai pada pengiriman kedua.**
>
> Diperbaiki: `V_AKSEPTASI_KOMPATIBEL` kini `nilai sekarang − SUM(arsip WHERE DITOLAK = 0)`, per kunci kasar.
>
> ### Cacat kedua: kedua view mengelompokkan dengan cara berbeda
>
> Keduanya satu janji, bukan dua: **"tidak ada yang hilang diam-diam"**. Janji itu hanya berlaku bila setiap kelompok yang jatuh dari `V01` muncul di `V02`. Keduanya **tidak sama**, di dua tempat:
>
> | | `V01` | `V02` (lama) |
> |---|---|---|
> | Penyambungan ke `ALOKASI_LAYER` | empat field layer + mata uang | **tanpa `LAYER_BAGIAN` dan `LAYER_BAGIAN_JENIS`** |
> | Penyaring | `KEADAAN_BARIS = 'LENGKAP'` | **tidak ada** |
>
> Akibatnya ada kelompok yang **jatuh dari keduanya** — keadaan terburuk: hilang, dan tidak terlihat hilang. Penyambungan, penyaring, dan pengelompokan kini **identik**; yang berbeda hanya `HAVING`-nya, dan keduanya saling melengkapi tepat.
>
> ### Satu keadaan yang tetap tidak tertangkap, dan dicatat terbuka
>
> Kelompok yang **seluruh** barisnya bukan `'LENGKAP'` tidak muncul di `V01` maupun `V02` — penyaringnya membuang barisnya lebih dulu. Itu bukan kelalaian: sebabnya **"belum lengkap"**, bukan **"tidak sepakat"** — sebab berbeda, yang tempatnya di `V_PARITAS_SHADOW` dan uji tiket `37`. Ditulis di kaki `V02` supaya yang membacanya tidak menyimpulkan sendiri bahwa cakupan kedua view itu lengkap.
>
> **REQ-033 tetap menunggui, tidak menahan**: ia mengubah **berapa banyak** kelompok yang ditolak, bukan apakah penolakannya ada.

*Asal: `T-16` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** View kompatibilitas menghasilkan bentuk yang dikenal hilir; kelompok yang besaran per-layernya tidak sepakat **ditolak** dan muncul di view penolakan.

`V_AKSEPTASI_KOMPATIBEL`, `V_AKSEPTASI_DITOLAK`. Keduanya membaca tabel kanonik **dan arsip muatan keluar**, lalu menghasilkan **selisih**.

**Tidak termasuk:** Menegakkan kunci kasar sebagai `UNIQUE` kedua — **ditolak dengan alasan**: itu menjadikan keterbatasan bentuk lama sebagai aturan bisnis sistem baru.

**Blocked by:**

- **REQ-033** — **menunggui, tidak menahan** (dari luar papan): ia menentukan apakah agregasi halus ke kasar pernah menemui baris yang tidak sepakat — itu mengubah **aturan penolakan di view**, bukan keberadaan view-nya
- ~~`18`~~ *(selesai)* — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran
- ~~`15`~~ *(selesai)* — Jembatan ke sistem lama berdiri sebagai tabel terpisah
- ~~`25`~~ *(selesai)* — Arsip muatan keluar: apa yang benar-benar dikirim, tersimpan sebagai catatan


**Dasar:** DECIDED(ADR-D-CNP-0023, ADR-D-CNP-0021 pengecualian `CASEID`). EVIDENCED: sapuan S1 — 17 parameter `InputParamOs.*` di `Activity\SaveDataToOSAksep_Act.xml`; `TypeLoss = .TreatyName` adalah **satuan kasar**.

- [x] Kolom dan tipenya sama dengan yang diterima hilir hari ini, termasuk `CASEID`.
- [x] Besaran aditif **dijumlahkan**; Limit Layer, Premi Deposit, persen premi pemulihan, kurs, dan porsi **tidak** — kelimanya justru menjadi penjaga kesepakatan.
- [x] Bila besaran per-layer berbeda antar baris yang dilebur, view **tidak menghasilkan baris**, dan kelompoknya muncul di view penolakan beserta sebabnya — **kini benar-benar berpasangan**: pengelompokan dan penyaring kedua view identik.
- [x] `LayerPart` dan `LayerPartType` tidak muncul di keluaran.
- [x] Keluarannya **selisih terhadap arsip muatan keluar**, bukan nilai mutlak — dan hanya arsip ber-`DITOLAK = 0` yang terhitung. **Sebelumnya tidak ada sama sekali.**

**Ketidakpastian:** **REQ-033** menentukan apakah agregasi halus→kasar pernah menemui baris yang tidak sepakat.

**Temuan S1 sudah TERTUTUP, dan tidak lewat REQ-018.** Sumber `OutOSAcc` ditemukan: `RDBList\GetDataOS.xml` dan `RDBList\GetDataCNPOS.xml` — EVIDENCED:

```sql
select sum(nvl(a.data_json.Value,0)) as "Value", ...
  from OS_AKSEPTASI_KLAIM a
 WHERE CASEID = {pyWorkPage.pzInsKey}
   AND a.data_json.TypeLoss = {InputParamOs.TypeLoss}
   AND a.data_json.Currency = {InputParamOs.Currency}
   AND STS_REJECT = 0
```

Penyaringnya **persis kunci alami ADR-D-CNP-0024**, dan agregatnya `SUM` atas seluruh baris berkunci sama. Jadi baris memang **tambahan**, bukan keadaan — dan itu pasangan wajib dari prosedur yang **selalu `INSERT`**. Naik dari DERIVED ke **EVIDENCED**.

**Satu syarat yang tidak terduga dan ikut mengikat**: penyaringnya memuat **`STS_REJECT = 0`**. Muatan yang ditolak **tidak** ikut dijumlahkan. Karena itu `ARSIP_MUATAN_KELUAR` berkolom `DITOLAK`, dan view menjumlahkan hanya yang tidak ditolak.

Yang tersisa di **REQ-018** hanya rekonsiliasi: apakah jumlah seluruh baris lama benar-benar sama dengan nilai sekarang di produksi. Bila tidak, itu temuan tentang **data lama**, bukan alasan mengubah rancangan.

## Claim Non Prop - 30 - Rekap per klaim per mata uang, dengan Retensi Cedant yang tidak pernah tercampur

---
status: selesai
---




> **SELESAI 19 September 2026.** Kedua kriteria sudah terpenuhi bentuk yang ada; satu kolom bertambah.
>
> **Tiga subkueri terpisah adalah ADR-D-CNP-0010 sebagai SQL.** Di sistem lama Retensi Cedant adalah **baris di dalam daftar alokasi yang sama** — sehingga setiap loop atas `SpreadingRisk` ikut melihatnya **kecuali menyaringnya sendiri**, dan setiap penulis loop harus ingat melakukannya (`BLUEPRINT.md` §2.4). Di sini ia tabel tersendiri dan kolom tersendiri: menjumlahkan alokasi layer **tidak dapat** memuat Retensi Cedant — bukan karena ada penyaring yang harus diingat, melainkan **karena barisnya tidak ada di sana**.
>
> **`CACAH_BELUM_LENGKAP` bertambah.** Sebuah penjumlahan tidak menyatakan apakah baris yang terjumlah sudah selesai; dua kelompok berangka sama dapat berarti hal yang sangat berbeda. ADR-D-CNP-0019 ditegakkan di tingkat baris oleh `KEADAAN_BARIS`, dan **hilang begitu `SUM` dijalankan** — kecuali agregatnya ikut membawanya keluar.
>
> `SUM` atas nol baris menghasilkan `NULL`, bukan `0`, dan itu **dipertahankan**: `NULL` berarti tidak ada baris untuk dijumlah, berbeda dari nol yang berarti sudah dihitung dan hasilnya nol.

*Asal: `T-17` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Penjumlahan alokasi tidak pernah diam-diam memuat Retensi Cedant.

`V_REKAP_KLAIM_MATA_UANG`.

**Blocked by:**

- ~~`17`~~ *(selesai)* — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang


**Dasar:** DECIDED(ADR-D-CNP-0010). EVIDENCED: `BLUEPRINT.md` §2.4 — di sistem lama setiap loop atas `SpreadingRisk` ikut melihat baris Retensi Cedant kecuali menyaringnya sendiri.

- [x] Retensi Cedant muncul sebagai **kolom tersendiri**, bukan bagian dari jumlah alokasi layer.
- [x] Menjumlahkan kolom alokasi saja menghasilkan angka tanpa Retensi Cedant — tanpa penyaring apa pun di sisi pemanggil.
- [x] Kelengkapan baris sumber terbaca dari agregatnya — `CACAH_BELUM_LENGKAP`, baru.

**Ketidakpastian:** Baris Retensi Cedant di sistem lama membawa `TotalClaim` yang **selalu nol** pada cabang mata uang sama (ADR-D-CNP-0010 tambahan). Sistem baru menghitungnya seperti Retensi Cedant biasa (AK-2b), jadi selisih terhadap sistem lama **diharapkan** dan masuk pengecualian bernama shadow-run.

## Claim Non Prop - 31 - Lima PageList lama tersedia sebagai view, tanpa menyimpan hasil yang dapat menyimpang

---
status: selesai
---



> **DILEPAS DARI KETERGANTUNGAN KOMITE, 19 September 2026.** Folder Komite **tertutup untuk batch ini**.
>
> **26 penanda yatim golongan 3b** mengikuti kebijakan **E17**: view menampilkannya **apa adanya** sebagai data yang dimigrasi dan **ditandai tak berpemilik**. Tidak dibuang — membuang butuh kepastian yang tidak kita punya. Tidak diberi perilaku — perilakunya tidak terbaca.


> **SELESAI 19 September 2026.** Kelima view ada dan menghasilkan angka dari sumbernya; satu kolom bertambah di masing-masing.
>
> **`CACAH_BELUM_LENGKAP`** — alasannya sama dengan tiket `30`, dan di sini berlaku untuk lima view sekaligus. Baris yang belum lengkap **tetap ikut dijumlah**, dan itu disengaja: nilai dalam mata uang aslinya sudah benar sejak awal, yang belum ada hanya padanan rupiahnya. Membuangnya akan menghasilkan angka yang salah diam-diam — **cacat yang sama dari arah berlawanan**.
>
> **Tidak ada nilai IDR yang dijumlahkan di view mana pun.** Bila kelak ada yang menambahkannya, kolom ini penjaga yang membuat kesalahan itu terbaca alih-alih tersembunyi.
>
> **Dua catatan sistem lama masuk ke badan berkasnya**, karena keduanya menjelaskan kenapa view lebih baik daripada tabel di tempat ini:
>
> | View | Perilaku lama yang tidak dapat terulang |
> |---|---|
> | `V_PENYEBARAN_AGREGAT` | `CountSpreadingXOL` **menghapus lalu membangun ulang** isinya setiap kali dijalankan. Sebagai view, tidak ada yang tersimpan untuk menyimpang (ADR-D-CNP-0013) |
> | `V_TOTAL_NILAI_PERTANGGUNGAN` | satu rutin dedup Java **membuang baris** dari `.ClaimData.TotalInterestInsured` — data klaim, tanpa pencatatan, tanpa penanda, tanpa jejak (**D40**). Sebagai view, tidak ada baris yang dapat dibuang |
>
> **Penyaring quota share tetap dua nilai**, `'10028'` dan `'10004'`, dan ketidakpastiannya ditulis di badan `V07`: apakah hanya keduanya yang berarti quota share **TIDAK DITEMUKAN**, dan daftar penuhnya ada di `POOLDATA.REINSURANCETYPE` yang **tidak punya satu pun constraint**. Ditulis sebagai literal **di satu tempat** supaya bila nilai ketiga ditemukan, ada tepat satu baris yang perlu diubah.

*Asal: `T-18` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** `ListTotalEstimation`, `SpreadingAdjustment`, `SpreadingAdjustmentQS`, `TotalInterestInsured`, `ListClaimAcceptation` terbaca tanpa tabel penyimpan.

`V_TOTAL_ESTIMASI`, `V_PENYEBARAN_AGREGAT`, `V_PENYEBARAN_AGREGAT_QS`, `V_TOTAL_NILAI_PERTANGGUNGAN`, `V_REKAP_AKSEPTASI`.

**Blocked by:**

- ~~`17`~~ *(selesai)* — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang
- ~~`18`~~ *(selesai)* — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran
- ~~`19`~~ *(selesai)* — Masukan mesin alokasi dan penyebaran tersimpan, bukan hanya keluarannya
- ~~`20`~~ *(selesai)* — Objek, kronologi, dan dokumen klaim


**Dasar:** DECIDED(ADR-D-CNP-0013). EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.4 — kunci agregasi `SpreadingAdjustment` adalah `(Currency, TreatyName)`, yaitu **satuan kasar**; dan `CountSpreadingXOL` **menghapus lalu membangun ulang** isinya setiap kali dijalankan.

- [x] Setiap view menghasilkan angka yang sama dengan penjumlahan langsung atas tabel sumbernya.
- [x] Tidak ada tabel yang menyimpan angka-angka ini.
- [x] Kelengkapan baris sumber terbaca dari tiap agregat — `CACAH_BELUM_LENGKAP`, baru di kelimanya.

**Ketidakpastian:** Penyaring quota share memakai `JENIS_REASURANSI_ID IN ('10028','10004')` — EVIDENCED dari `CountLossAllocation_act` Langkah 10. Apakah hanya dua nilai itu yang berarti quota share **TIDAK DITEMUKAN**; daftar penuhnya ada di `POOLDATA.REINSURANCETYPE`, yang **tidak punya satu pun constraint**.

## Claim Non Prop - 32 - Uji: kunci alami akseptasi

---
status: menunggu-instance
---



> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-24` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa baris ganda tidak dapat lahir.

**Blocked by:**

- ~~`18`~~ *(selesai)* — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran


**Dasar:** DECIDED(ADR-D-CNP-0024). Bentuk uji mengikuti `pengetahuan/PREFLIGHT.sql`: SQL berblok, tiap blok berdiri sendiri, dapat dijalankan dari TOAD. Itu satu-satunya prior art — `MEMORI_PEMAHAMAN.MD` §10.8 mencatat ketiadaan kerangka uji sebagai "aspek yang tidak ada".

- [ ] `INSERT` kedua dengan (klaim, layer, mata uang) sama **ditolak**; `INSERT` dengan `LayerPart` berbeda **diterima**, dan itu perilaku yang dimaksud — bukan celah.

**Ketidakpastian:** **REQ-033**.

## Claim Non Prop - 33 - Uji: NULL bukan nol

---
status: menunggu-instance
---



> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-27` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa nilai yang belum dihitung tidak pernah terbaca sebagai nol.

**Blocked by:**

- ~~`16`~~ *(selesai)* — Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs


**Dasar:** DECIDED(ADR-D-CNP-0019). **Ketidakpastian** Tidak ada.

- [ ] Baris menunggu kurs ber-IDR `NULL`, `KEADAAN_BARIS` bukan `LENGKAP`, dan agregasi atasnya **tidak** memperlakukannya sebagai nol. Baris bernilai nol yang sudah dihitung dapat dibedakan darinya lewat satu kueri.

**Ketidakpastian:** Tidak ada.

## Claim Non Prop - 34 - Uji: paritas

---
status: menunggu-instance
---



> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-31` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa setiap baris hasil migrasi dapat ditunjuk balik ke barisnya di sistem lama, tepat satu.

**Blocked by:**

- ~~`23`~~ *(selesai)* — Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama
- ~~`15`~~ *(selesai)* — Jembatan ke sistem lama berdiri sebagai tabel terpisah


**Dasar:** DECIDED(ADR-D-CNP-0005), lewat ADR-D-CNP-0004. **Ketidakpastian** **REQ-018**; **REQ-015** mengukur status buka/tutup 8 klaim bertambalan.

- [ ] Untuk satu klaim contoh, tiap baris baru punya **tepat satu** pasangan di `MIGRASI_KORELASI`. Kasus uji utama **`CLMNP-975`** — tambalan terbaru, 2026-07-16, EVIDENCED `BLUEPRINT.md` §7.2.

**Ketidakpastian:** **REQ-018**; **REQ-015** mengukur status buka/tutup 8 klaim bertambalan.

## Claim Non Prop - 35 - Satu pintu tulis ditegakkan hak akses, bukan kesepakatan

---
status: selesai
---



> **Gerbang K1 gugur — 19 September 2026. Tiket ini tidak lagi menunggu apa pun dari bagian K.**
> K1 ditutup sebagai **K1a**: ADR-D-CNP-0016 tetap berlaku tanpa revisi, dan modul ini **tidak mengambil data HRD sama sekali**. `V_MST_USER_TEKNIS` tidak dimiliki dan tidak dimigrasi (ADR-D-CNP-0026); database link yang berjalan produksi (`hrdasm.v_hrd_mst@asmd.sinarmas.co.id`, D20) tetap hidup di sistem lama dan bukan urusan tiket ini.
>
> Yang tersisa sebagai catatan, bukan penahan: **REQ-037**. Ia menentukan apakah klaim ADR-D-CNP-0017 berlaku sungguhan di instance tujuan — enam jalur tulis yang melewati hak objek. Rumusan klaim ADR-D-CNP-0017 sudah diperbaiki di badannya (*"tidak ada tulisan dari luar pintu yang tidak meninggalkan jejak"*), dan `Z01_PENGAWASAN_TULIS.sql` adalah bentuknya. Tiket ini menegakkan hak objek; ia tidak menjanjikan pencegahan yang tidak dapat dijanjikan skema.



> **SELESAI 19 September 2026 — dan cacat yang sama dengan tiket `01` terulang di sini.**
>
> **Keempat puluh empat `REVOKE ALL` akan menggugurkan pemasangan di instance bersih.** Di instance yang belum pernah dipakai, tidak satu pun hak itu pernah diberikan — dan Oracle **menolak pencabutan hak yang tidak ada**: `ORA-01927`. Baris pertama gagal, dan pemasangan berhenti di sana.
>
> Ini **kelas cacat yang sama** dengan sembilan `REVOKE` di `00_SKEMA_DAN_AKUN.sql` yang sudah diperbaiki lebih dulu hari ini, dan ia terulang karena keduanya ditulis dengan anggapan yang sama: bahwa mencabut sesuatu yang tidak ada tidak berakibat apa-apa. **Di Oracle, itu keliru — dan sekali ketahuan, seharusnya disapu ke seluruh berkas, bukan diperbaiki di satu tempat.**
>
> Diperbaiki dengan bentuk yang sama: blok PL/SQL yang mencabut bila ada, melewati `ORA-01927`, dan **merambat untuk galat apa pun selain itu** — sehingga kegagalan yang sesungguhnya tetap menghentikan pemasangan.
>
> **Pengecualian arsip kini disebut di sini juga.** Hibahnya tetap di `22_ARSIP_MUATAN_KELUAR.sql` (`SELECT, INSERT` saja, tanpa `UPDATE`, tanpa `DELETE`); yang ditulis di `V00` hanya **faktanya**, supaya pengecualian itu tidak hanya hidup di satu berkas tabel yang tidak dibaca orang saat meninjau hak akses.
>
> **Tiga pemeriksaan ditulis di kaki berkas**, dan yang ketiga menyatakan batasnya: keenam jalur REQ-037 **tidak dapat diperiksa dari sini**, dan nol baris pada pemeriksaan pertama **tidak berarti** tidak ada yang dapat menulis.
>
> **REQ-021 dan REQ-037 tetap menunggui, tidak menahan**: skema baru tidak mengulangi hak yang dipegang `POOLDATA` atas tabel work Pega, apa pun jawabannya.

*Asal: `T-20` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada akun selain akun aplikasi yang dapat menulis ke tabel kanonik; hilir membaca hanya lewat view.

`V00_HAK_AKSES.sql`: **44 `REVOKE ALL`** atas tabel kanonik untuk `POOLDATA` dan `KLAIMNP_HILIR`; `GRANT SELECT` atas view. **Termasuk `ARSIP_MUATAN_KELUAR`**: akun aplikasi memegang `INSERT` dan `SELECT` saja — **tanpa `UPDATE` dan tanpa `DELETE`**, karena arsip itu tulis-sekali.

**Tidak termasuk:** Perintah yang dijalankan tangan — `GRANT` dan `REVOKE` adalah **objek DDL**. Hak yang tidak tertulis di DDL tidak dapat diaudit.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`22`~~ *(selesai)* — Koreksi bernilai tercatat menggantikan tambalan di dalam kode
- ~~`13`~~ *(selesai)* — Bentuk lama boleh masuk utuh, di satu tempat saja
- ~~`15`~~ *(selesai)* — Jembatan ke sistem lama berdiri sebagai tabel terpisah
- ~~`29`~~ *(selesai)* — Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam
- ~~`23`~~ *(selesai)* — Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama


**Dasar:** DECIDED(ADR-D-CNP-0017, ADR-D-CNP-0028). EVIDENCED: `GRANT ALTER, DELETE, INDEX, INSERT, REFERENCES, SELECT, UPDATE ... TO POOLDATA` pada tabel work Pega — hak yang ada akan dipakai cepat atau lambat.

- [x] `POOLDATA` gagal `INSERT`, `UPDATE`, dan `DELETE` ke setiap tabel kanonik — 44 pencabutan, kini **tahan instance bersih**.
- [x] `POOLDATA` berhasil `SELECT` lewat view kompatibilitas — hibahnya di `V01`…`V09`.
- [x] Akun hilir gagal `SELECT` langsung dari tabel kanonik mana pun.
- [x] Akun aplikasi memegang `SELECT` dan `INSERT` saja atas `ARSIP_MUATAN_KELUAR` — tulis-sekali, dan faktanya terbaca di `V00` maupun di berkas tabelnya.

**Ketidakpastian:** **REQ-021** — apakah `POOLDATA` benar-benar menulis ke tabel Pega hari ini belum diketahui. Bila ya, migrasinya memerlukan pokok tersendiri untuk manajemen; itu **tidak menahan tiket ini**, karena skema baru tidak mengulangi hak itu apa pun jawabannya.

## Claim Non Prop - 36 - Uji: pasangan uang dan mata uang

---
status: menunggu-instance
---



> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-25` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa tidak ada nilai uang yang dapat hidup tanpa satuannya, di tabel mana pun.

**Blocked by:**

- ~~`16`~~ *(selesai)* — Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs
- ~~`28`~~ *(selesai)* — Pasangan IDR–kurs


**Dasar:** DECIDED(ADR-D-CNP-0007). **Ketidakpastian** T-35 dapat menambah kolom yang ikut diuji.

- [ ] Nilai uang tanpa mata uang **ditolak**; nilai IDR tanpa kurs **ditolak**; kurs tanpa IDR **diterima**. Diuji di **setiap** tabel bernilai uang, bukan satu.

**Ketidakpastian:** T-35 dapat menambah kolom yang ikut diuji.

## Claim Non Prop - 37 - Uji: penularan keadaan menunggu kurs

---
status: menunggu-instance
---



> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-28` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa ketiadaan kurs merambat ke seluruh turunannya, bukan berhenti di baris asalnya.

**Blocked by:**

- ~~`17`~~ *(selesai)* — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang
- ~~`18`~~ *(selesai)* — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran
- ~~`26`~~ *(selesai)* — Premi pemulihan dapat dihitung ulang dari barisnya sendiri


**Dasar:** DECIDED(ADR-D-CNP-0014). **Ketidakpastian** Penularannya aturan prosedur; basis data menyimpan keadaannya, tidak menegakkan penularannya. Uji ini karena itu menguji **hasil**, bukan mekanismenya.

- [ ] Baris turunan dari nilai yang menunggu kurs **ikut** bertanda menunggu kurs, di seluruh tabel turunannya — alokasi, Retensi Cedant, akseptasi, Adjustment, premi pemulihan.

**Ketidakpastian:** Penularannya aturan prosedur; basis data menyimpan keadaannya, tidak menegakkan penularannya. Uji ini karena itu menguji **hasil**, bukan mekanismenya.

## Claim Non Prop - 38 - Uji: bentuk view kompatibilitas

---
status: menunggu-instance
---



> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-30` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa hilir menerima bentuk yang sama dengan hari ini, dan menerima selisih pada muatan kedua.

**Blocked by:**

- ~~`29`~~ *(selesai)* — Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam
- ~~`05`~~ *(selesai)* — Sapuan S1: kolom yang sesungguhnya diterima Arasapas dan kasir


**Dasar:** EVIDENCED: `SPEC-MODEL-DATA.md` bagian 18. **Ketidakpastian** Pola selisih `InputParamOs.Value - OutOSAcc...` — lihat T-16.

- [ ] Kolom dan tipe yang dihasilkan sama dengan yang diterima hilir hari ini, termasuk `CASEID`. Dibandingkan terhadap 17 parameter `InputParamOs.*` hasil sapuan S1 — **dibaca, bukan dirancang**.
- [ ] **Muatan kedua atas klaim, jenis reasuransi, dan mata uang yang sama menghasilkan selisih, bukan nilai penuh.**
- [ ] Muatan yang ditandai ditolak **tidak** ikut mengurangi muatan berikutnya.

**Ketidakpastian:** Pola selisih `InputParamOs.Value - OutOSAcc...` — lihat T-16.

## Claim Non Prop - 39 - Uji: satu pintu tulis

---
status: menunggu-instance
---



> **MENUNGGU INSTANCE, bukan aktif — 19 September 2026.** Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun. Rancangannya sudah lengkap; yang kurang mesinnya.

*Asal: `T-29` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Terbukti bahwa satu pintu tulis berlaku sebagai hak akses, bukan sebagai kesepakatan.

**Blocked by:**

- `35` — Satu pintu tulis ditegakkan hak akses, bukan kesepakatan


**Dasar:** DECIDED(ADR-D-CNP-0017). **Ketidakpastian** **REQ-021**.

- [ ] Akun selain akun aplikasi **gagal** `INSERT` ke tabel kanonik dan **berhasil** `SELECT` lewat view.

**Ketidakpastian:** **REQ-021**.

## Claim Non Prop - REA - Papan tiket — lapisan data Claim Non Prop

**Tiket: aktif 0 · menunggu instance 8 · tertahan 1 · selesai 29 · mati 1**

> Pencacah ini menghitung **tiket**. Register REQ punya pencacahnya sendiri dengan format yang mirip; keduanya diberi label bendanya supaya tidak pernah tertukar.

#### Yang sebenarnya dapat dimulai

**0 dari 9 tiket di papan.** Sisanya tertahan sesuatu di **luar papan** — keputusan yang belum dikonfirmasi atau REQ yang belum kembali — bukan tertahan tiket lain. Angka ini, bukan "aktif 0", yang menggambarkan keadaan proyek.

| Penahan dari luar papan | Tiket yang dibekukannya | Ditahan langsung |
|---|---|---|
| **lingkup-batch** | 1 | langsung: `11` |

**Tidak ada yang dapat dimulai sekarang.** Empat sapuan sumber yang tidak menyentuh gerbang sudah dikerjakan sampai habis; sisanya menunggu konfirmasi ADR-D-CNP-0028 dan REQ-032.

Diterbitkan 18 September 2026 dari `_migration-docs/claim-non-prop/TICKETS.md`, dibangun ulang dari field `status:` tiap berkas.
Tracker belum terkonfigurasi (`/setup-matt-pocock-skills` belum dijalankan), jadi papannya berkas lokal.

#### Aturan papan

- **Status adalah field di dalam berkas tiket** (`status: aktif | tertahan | selesai | selesai-sebagian`), bukan lokasi foldernya. Folder `_selesai/` dan `_tertahan/` hanya kerapian; pembangkit membaca field, dan **tidak pernah menghapus berkas tiket**.
- **Penahan dari luar papan ditulis dengan nama aslinya** — `ADR-D-CNP-0028`, `REQ-032` — tidak diterjemahkan jadi nomor tiket. Yang menahan dari luar harus terlihat berasal dari luar.
- **Penahan yang sudah selesai dicoret, tidak dihapus**, supaya rantainya tetap terbaca: `~~05~~ *(selesai)*`.
- **Nomor yang lompat bukan kekeliruan.** Penomoran tidak disusun ulang, karena menyusunnya ulang memutus setiap rujukan yang sudah ada.

#### Papan

| # | Asal | Tiket | Ditahan oleh |
|---|---|---|---|
| `11` | T-38 | Tipe desimal di sisi Golang untuk nilai uang | **lingkup-batch** |
| `24` | T-26 | Uji: batas tanggal inklusif | ~~`12`~~ |
| `32` | T-24 | Uji: kunci alami akseptasi | ~~`18`~~ |
| `33` | T-27 | Uji: NULL bukan nol | ~~`16`~~ |
| `34` | T-31 | Uji: paritas | ~~`23`~~ ~~`15`~~ |
| `36` | T-25 | Uji: pasangan uang dan mata uang | ~~`16`~~ ~~`28`~~ |
| `37` | T-28 | Uji: penularan keadaan menunggu kurs | ~~`17`~~ ~~`18`~~ ~~`26`~~ |
| `38` | T-30 | Uji: bentuk view kompatibilitas | ~~`29`~~ ~~`05`~~ |
| `39` | T-29 | Uji: satu pintu tulis | `35` |

#### Menunggui, tidak menahan

Tiket ini **boleh jalan**. Jawaban yang ditunggu mengubah satu hal yang tiketnya sudah menyebut — sebuah constraint, sebuah aturan penolakan, sebuah jalur migrasi — bukan rancangannya. Dipisahkan dari penahan sungguhan supaya papan tidak berteriak serigala.

| # | Menunggu | Yang berubah bila jawabannya lain |
|---|---|---|

#### Dapat dimulai hari pertama

**Tidak ada.** Setiap tiket yang tersisa tertahan sesuatu di luar papan. Itu keadaan yang sah dan bukan kebuntuan kerja: yang dapat dikerjakan tanpa menyentuh gerbang sudah dikerjakan sampai habis.

---

#### Menunggu instance — bukan aktif

Rancangannya lengkap; yang kurang mesinnya. Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun.

| # | Tiket |
|---|---|
| `24` | Uji: batas tanggal inklusif |
| `32` | Uji: kunci alami akseptasi |
| `33` | Uji: NULL bukan nol |
| `34` | Uji: paritas |
| `36` | Uji: pasangan uang dan mata uang |
| `37` | Uji: penularan keadaan menunggu kurs |
| `38` | Uji: bentuk view kompatibilitas |
| `39` | Uji: satu pintu tulis |

#### Tertahan di luar papan

| # | Tiket | Ditahan |
|---|---|---|
| `11` | Tipe desimal di sisi Golang untuk nilai uang | **lingkup-batch** |

#### Mati — dibatalkan, bukan ditunda

Tiket yang premisnya gugur. Berkasnya tidak dihapus: tiket yang pernah ada adalah bukti bahwa premisnya diperiksa, bukan dilewatkan. Nomornya tidak dipakai ulang.

| # | Tiket | Sebab |
|---|---|---|
| `02` | Tabel singkatan tertutup ditulis dan dibekukan | premis gugur — lihat kepala berkasnya di `_mati/` |

#### Selesai

| # | Tiket | Status |
|---|---|---|
| `01` | Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis | selesai |
| `03` | Penamaan yang tidak nyaris kembar | selesai |
| `04` | Index penopang setiap foreign key | selesai |
| `05` | Sapuan S1: kolom yang sesungguhnya diterima Arasapas dan kasir | selesai |
| `06` | Keadaan kasus: nilai mana yang ada, dan siapa yang membacanya | selesai |
| `07` | Sapuan S3: perilaku per nilai `PaymentType` | selesai |
| `08` | Sapuan S4: adakah padanan IDR untuk biaya penilaian, salvage, dan biaya lain | selesai |
| `09` | Sapuan S9: apa persisnya yang ditulis `EditXOLAlokasi` | selesai |
| `10` | Sapuan S5, S6, S7, S8, dan dua sumber yang belum pernah dibaca | selesai |
| `12` | Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap | selesai |
| `13` | Bentuk lama boleh masuk utuh, di satu tempat saja | selesai |
| `14` | Nilai yang tidak dapat diurai tercatat, tidak dibulatkan, tidak dibuang | selesai |
| `15` | Jembatan ke sistem lama berdiri sebagai tabel terpisah | selesai |
| `16` | Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs | selesai |
| `17` | Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang | selesai |
| `18` | Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran | selesai |
| `19` | Masukan mesin alokasi dan penyebaran tersimpan, bukan hanya keluarannya | selesai |
| `20` | Objek, kronologi, dan dokumen klaim | selesai |
| `21` | Tanggal tutup buku dan tarif menjadi data bertanggal berlaku | selesai |
| `22` | Koreksi bernilai tercatat menggantikan tambalan di dalam kode | selesai |
| `23` | Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama | selesai |
| `25` | Arsip muatan keluar: apa yang benar-benar dikirim, tersimpan sebagai catatan | selesai |
| `26` | Premi pemulihan dapat dihitung ulang dari barisnya sendiri | selesai |
| `27` | Domain kolom keadaan | selesai |
| `28` | Pasangan IDR–kurs | selesai |
| `29` | Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam | selesai |
| `30` | Rekap per klaim per mata uang, dengan Retensi Cedant yang tidak pernah tercampur | selesai |
| `31` | Lima PageList lama tersedia sebagai view, tanpa menyimpan hasil yang dapat menyimpang | selesai |
| `35` | Satu pintu tulis ditegakkan hak akses, bukan kesepakatan | selesai |
