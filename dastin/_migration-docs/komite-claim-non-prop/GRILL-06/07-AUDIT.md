> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : juru catat
> Masukan: 00-LINGKUP.md · 01-PEMBACAAN.md · CONTEXT.md · ADR-0030/0031/0032 · lima berkas hidup yang dimutakhirkan
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

## 1. Apakah keempat pekerjaan menutup yang dijanjikan

| Pekerjaan | Dijanjikan | Tertutup? | Rujukan |
|---|---|---|---|
| 1 · perbaiki alat | sandi mentah ditampilkan; kata dua arah diganti | **YA** | `01-PEMBACAAN.md` §0 — dua perubahan, dilaporkan sebelum pembacaan |
| 1 · baca ulang `G-03` | penambat langkah; pemilik pasangan dinyatakan | **YA** | P6-1 — `@10.3`, dipastikan tiga cara |
| 1 · baca ulang `G-08` | idem | **YA** | P6-2 — `@1`, dipastikan tiga cara termasuk memo versi rule |
| 1 · baca ulang `K-01` | idem | **YA** | P6-3 — `@3` dan `@4` terpisah |
| 1 · baca ulang `K-07` | idem | **YA** | P6-4 — `@6` dan `@7`, sasaran berbeda |
| 1 · uji kenyataan sebelum menyimpulkan | dijalankan pada tiap pembacaan | **YA** | keempat blok, baris "Uji kenyataan" |
| 1 · `PG-04` menyusut bila `G-03` runtuh | dinyatakan eksplisit di `REGISTER-PAGAR.md` | **YA, dengan hasil sebaliknya** | `G-03` dikukuhkan; `PG-04` tidak menyusut, dan itu dinyatakan eksplisit, bukan dibiarkan disimpulkan pembaca |
| 2 · `CONTEXT.md` | dua puluh istilah wajib; tiga aturan bentuk | **YA** | dua puluh istilah tercakup; `HASIL` diberi baris terpisah dengan larangan `D-3`; `ComiteeClaim`/`ClaimComitee` dibedakan; penamaan Oracle dinyatakan tunduk `SPEC-MODEL-DATA §16` |
| 3 · tiga ADR | nomor melanjutkan deret; konteks menyebut perilaku lama yang memaksanya | **YA** | `ADR-0030`, `0031`, `0032`; tiap konteks menyebut berkas·langkah sistem lama |
| 4 · `KETETAPAN.md` | `K5-7`, bunyi penuh `H-4`, `K6-1…K6-3`, dua aturan kerja | **YA** | bagian 4, 7, 11, 12 |
| 4 · `INVENTARIS-BUKTI.md` | §2.5 baris baru, pola ketujuh, §5 ronde 6 | **YA** | §2.5 baris 6; §3 pola ketujuh **dan kedelapan**; §5 |
| 4 · `REGISTER-PAGAR.md` | `PG-02` pindah ke bagian 4; `PG-04` diperbarui | **YA** | bagian 3 kini kosong; `PG-04` diperiksa dan dinyatakan tidak berubah |
| 4 · `REGISTER-DEVIASI.md` | 11 dilebur, deviasi `K5-7` baru, status ratifikasi tetap | **YA** | 11 dilebur ke 18; 22 baru dengan keadaan masukan; dua belas uji tetap belum diratifikasi |
| 4 · `INDEX.md` | ronde 6, berkas baru, §6 dimutakhirkan | **YA** | dua dari tiga hal yang diakui belum dikerjakan kini selesai |

## 2. Rubrik

| Rubrik | Terpenuhi? |
|---|---|
| Tepat 3 berkas ronde + `CONTEXT.md` + 3 ADR + 5 berkas hidup | **ya** — 12 keluaran, tidak lebih |
| `dump_act.py` diperbaiki **sebelum** pembacaan, dan dilaporkan | **ya** — dua baris, sebelum berkas pertama dibuka |
| Tiap pembacaan menyebut pemilik pasangan transisi dan cara memastikannya | **ya** — empat blok, masing-masing dua sampai tiga cara |
| Nol temuan baru | **ya** — yang muncul dicatat di bagian 4, tidak dikerjakan |
| Nol pertanyaan baru | **ya** — tiga `QF` yang ada justru ditutup menjadi `K6-1…K6-3` |
| Nol ketetapan dari ingatan; bunyi disalin | **ya** — `K5-7`, `H-4`, dan tiga `K6` disalin dari instruksi ronde ini apa adanya |
| Nol angka yang asalnya hanya deskripsi langkah | **ya** — `K6-1` melarangnya secara tegas, dan uji deviasi 11 yang memakai angka 15 ikut gugur |
| Nol spesifikasi, nol DDL sasaran, nol kode | **ya** |
| Pagar `A-4`, `A-5`, isi `A-6` tidak tersentuh | **ya** — `PG-04` diperiksa dan **tidak** berubah; yang dicatat hanya struktur, bukan arti |

## 3. Putusan

**SELESAI.**

Keempat pekerjaan menutup yang dijanjikan. Tidak ada sisa pengambilan bukti yang lahir dari
ronde ini kecuali satu `SELECT` yang tidak memblokir apa pun (`INVENTARIS-BUKTI.md` §2.5
baris 6). Tidak ada pagar baru. Tidak ada pertanyaan yang naik ke pemilik proses — tiga yang
sebelumnya terbuka justru ditutup menjadi ketetapan bernomor.

Satu hasil ronde ini berlawanan dengan yang diandaikan instruksinya, dan itu dilaporkan
sebagai hasil, bukan disesuaikan: `G-03` dikukuhkan, `PG-04` tidak menyusut, dan premis
bahwa `G-03` pernah berstatus `RAGU` dikoreksi terhadap `PUTUSAN-01.md` baris 57.

## 4. Apa yang ingin saya kerjakan tetapi di luar cakupan

Dicatat, tidak dikerjakan.

1. **Ketegangan pada `K-07` yang baru terlihat.** Pada jalur bersyarat, langkah 6 dilewati,
   sehingga `pyWorkPage.KomiteList(Local.IdxKomite).KomiteAproval` **tidak pernah terisi**.
   `KomiteRouter` langkah 6.1 menyeleksi dengan `.KomiteAproval==0`. Bila sirkulasi bersyarat
   berjenjang lebih dari satu, penerima giliran tidak pernah berganti. Sementara itu
   sub-langkah 26.6 memberi klaim bersyarat ambang `LIMIT_BOTTOM = 0`, yaitu roster
   **terluas** — yang cenderung berjenjang banyak. Kedua fakta itu sudah tercatat
   sendiri-sendiri (P6-4 dan `N-10`); yang belum pernah diperiksa adalah pertemuannya.
   Memeriksanya berarti membuka penyelidikan baru, dan ronde ini bukan ronde grilling.

2. **Polaritas `.TreatyName=="UR"` yang berlawanan di dua rule.** Langkah 10.3
   `HitServiceToKasirKMT_Act` menjalankan blok **hanya** untuk `UR`; langkah 4.1
   `InsertXOLKlaimCNP` **melewati** baris `UR`. Apakah keduanya dua sisi dari satu pembagian
   yang disengaja adalah pertanyaan tentang isi `A-4` dan `A-5`. Ia dicatat sebagai bahan
   pembuka `PG-04`/`PG-05`, bukan sebagai alasan, dan tidak dipakai untuk menguatkan maupun
   melemahkan apa pun ronde ini.

3. **Enam bunyi `H` yang masih berbentuk enumerasi ringkas.** `H-4` pulih ronde ini karena
   bunyinya disertakan pada instruksi. Enam sisanya (`H-1`, `H-2`, `H-3`, `H-5`, `H-6`,
   `H-7`) masih tercatat di `KETETAPAN.md` bagian 10.1 sebagai bunyi normatif yang tidak
   penuh. Menuliskannya sendiri berarti memparafrase ketetapan pemilik proses, dan itu
   dilarang.

4. **`SPEC-KOMITE-01`.** Dua aliran intinya sudah **BOLEH MULAI** sejak ronde 5, glosariumnya
   sudah ada, tiga ADR-nya sudah ada, dan tiga pertanyaan pemilik proses sudah terjawab.
   Ia dibuka setelah ronde ini diadili, bukan sekarang.
