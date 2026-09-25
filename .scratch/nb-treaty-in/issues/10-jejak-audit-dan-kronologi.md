# 10: Jejak audit dan kronologi — identitas akses terpisah dari nama tampilan

**Status:** ready-for-agent
**Blocked by:** 03
**Menutup:** AC 39 · 40 · 41 · 42 · 43 · 44 · 71 · 72 *(8 AC)* — US 15 · 16 · 20 · 40 · 41 · 42

## Hasil & nilai pengguna

Hari ini riwayat akseptasi mencatat siapa memutuskan apa. ⛔ Tetapi penampung identitas operator
**tidak pernah diisi** dari aplikasi — `[terverifikasi]` **nol dari 1.151** naskah SQL di 21 modul
menyebutnya sebagai kolom. ⚠️ Dan medan "nama operator" diisi dari **dua sumber berbeda** di dua
tahap berdekatan: pengenal akun di satu tahap, nama tampilan di tahap lain.

Sesudah tiket ini, setiap tindakan mencatat **identitas akses login**-nya terpisah dari **nama
tampilan**-nya, dan ⭐ keduanya tidak lagi tertukar — jejaknya tetap benar walau nama tampilan
seseorang berubah.

## Area codebase

- Lapisan repository: penulisan riwayat akseptasi
- Lapisan service: kronologi dan catatan pengguna

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pemisahan yang sudah benar | `RDBList\InsertViewSuggest_SQL.xml` — penampung akses-login diisi dari pengenal akun, terpisah dari penampung nama tampilan |
| ⛔ Pengisian yang **bug** | `DataTransform\DeptHeadTreatyInUW_preDT.xml` — medan nama operator diisi dari **pengenal akun** |
| Pengisian yang benar | `DataTransform\InputPolicyTreatyIn_preDT.xml` — dari **nama tampilan** |
| Mekanisme kronologi | `DataTransform\AddToListCommentsPolicyTreatyIn_DT.xml` — empat medan: catatan, penanda persetujuan, tanggal, operator |
| ⛔ Nama orang di dalam teks pesan | `DataTransform\DeptHeadTreatyIn_UW_postDT.xml` — **1** kemunculan; pesan lain di berkas sama mengambilnya dari data |

## ADR terkait

- ⭐ **ADR-0007** — jejak audit merekam siapa + kapan untuk setiap transisi **dan setiap jalur balik**

## Acceptance criteria

- [ ] **AC 39** — penampung identitas operator diisi dari **identitas akses login**
- [ ] **AC 40** — penampung nama diisi dari **nama tampilan**
- [ ] **AC 41** — penampung identitas operator **terisi** pada setiap penulisan riwayat
- [ ] **AC 42** — medan nama operator diisi dari **nama tampilan di setiap tahap**, termasuk tahap
      jenjang ketiga
- [ ] **AC 43** — setiap perpindahan tahap menulis **satu baris riwayat**
- [ ] **AC 44** — pemberitahuan menyebut nama orang **dari data**; ⛔ tidak tertanam di dalam teks
- [ ] **AC 71** — catatan pengguna tersimpan bersama tanggal dan operatornya
- [ ] **AC 72** — riwayat dapat dibaca **berurutan waktu**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **16** | berapa baris yang penampung identitas operatornya **sudah terisi** sekarang | tidak menahan — ⚠️ bila bukan nol, ada penulis **di luar aplikasi** yang belum terpeta |

## Perintah verifikasi

1. Lakukan satu tindakan di tiap jenjang — ⭐ tiga baris riwayat, **masing-masing berisi kedua
   identitas**.
2. Ubah nama tampilan seorang pengguna — ⭐ riwayat lama **tidak berubah**.
3. Baca pemberitahuan — ⭐ namanya **dari data**, dan berubah mengikuti data.

## Catatan

⚠️ `[penyimpangan sadar]` Mengisi penampung identitas operator adalah **perbaikan jejak audit**,
bukan peniruan — kolomnya selama ini kosong. ⭐ Dan pengisian nama operator dari pengenal akun
adalah **bug**, ⛔ bukan perbedaan maksud antar tahap *(P33)*.
