# RINGKASAN — Claim Non Prop

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (bertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`); `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**, terurai ke `pengetahuan/ddl/` (49 berkas); dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).
> Disusun 18 September 2026. **Tidak ada satu pun kalimat analisis baru di sini** — seluruh isinya rujukan ke artefak yang sudah ada di folder ini.
> Berlaku hanya untuk keadaan sistem pada ekspor di atas. Perubahan sesudahnya tidak tercermin; deteksinya lewat sapuan ulang (ADR-0012).

Sesi grilling modul Claim Non Prop **selesai 18 September 2026**. Frontier yang dapat dibuka dari XML sendiri habis.

> **Register pertanyaan ditutup seluruhnya 19 September 2026.** `_selesai/OPEN-QUESTIONS.md` kini **aktif 0 · selesai 57 · dihapus 2**. Bagian di bawah yang menyebut butir "masih terbuka", "menunggu", atau "sembilan butir DEFERRED" adalah keadaan **18 September** dan sudah tidak berlaku — dipertahankan sebagai catatan perjalanan, bukan sebagai keadaan. Yang masih berjalan hanyalah **verifikasi**. **REQ-032 pun turun jadi verifikasi** pada tanggal yang sama — batas 30 byte ditetapkan sebagai aturan tetap. Empat REQ menahan **pelaksanaan** tiket, bukan keputusan rancangan: REQ-018, REQ-033, REQ-021, REQ-037.

---

## 1. Apa yang sudah diputuskan — 28 ADR

Berkas lengkap di `docs/adr/`.

| ADR | Keputusan |
|---|---|
| 0001 | Aggregate root adalah Klaim; satu klaim satu kejadian kerugian |
| 0002 | Klaim yang sudah ditutup tidak dapat dibuka kembali |
| 0003 | Presisi tinggi sepanjang rantai, pembulatan hanya di tepi — **accepted** 18 Sep lewat AK-1 |
| 0004 | Tambalan per-case tidak dimigrasi; angkanya dipindahkan sebagai data |
| 0005 | Claim dan Komite satu unit cutover, dengan shadow-run sebagai bukti paritas |
| 0006 | RBAC sistem baru dirancang dari nol, bukan diwarisi |
| 0007 | Setiap nilai uang disimpan berpasangan; ambang kewenangan dibandingkan terhadap nilai IDR |
| 0008 | Nilai hasil suntingan manual bertahan terhadap hitung ulang, dan selalu terlihat |
| 0009 | Keputusan alur disimpan sebagai field terstruktur; komentar tidak pernah dibaca mesin |
| 0010 | Layer dan Retensi Cedant adalah dua entitas terpisah |
| 0011 | Menjalankan ulang perhitungan alokasi harus menghasilkan keadaan yang identik |
| 0012 | Tambalan baru dideteksi lewat sapuan ulang, bukan lewat register |
| 0013 | Perhitungan dijalankan sebagai turunan dari data, bukan akibat penekanan tombol |
| 0014 | Aturan penguraian teks ke angka, dan perlakuan atas kurs yang tidak ada |
| 0015 | Akseptasi adalah satu entitas dengan keadaan, bukan dua tabel |
| 0016 | Data pegawai masuk lewat API atau salinan tersinkron, bukan database link |
| 0017 | Data aplikasi hanya diubah lewat aplikasi |
| 0018 | Klaim terdampak penjaga tanggal yang salah dimigrasi apa adanya, dan didaftar |
| 0019 | NULL bukan nol; keadaan perhitungan disimpan di tingkat baris |
| 0020 | Fakultatif dan Treaty dua entitas berbeda; sistem ini hanya memiliki Treaty |
| 0021 | Pengenal polis dan klaim diberi nama menurut isinya; `CASEID` dipensiunkan |
| 0022 | Masa berlaku treaty disimpan sebagai tanggal, batasnya inklusif |
| 0023 | Data akseptasi disimpan dalam satu bentuk kanonik; sistem hilir diberi view, bukan salinan |
| 0024 | Satu akseptasi per klaim per layer per mata uang |
| 0025 | Tanggal tutup buku disimpan bertanggal berlaku, bukan satu baris tanpa riwayat |
| 0026 | Batas kepemilikan mengikuti nama class: `-Work-` dan `-Data-` dimiliki, `-Int-` tidak |
| 0027 | Dokumen klaim dirujuk, tidak disimpan ulang |
| 0028 | Basis data tujuan Oracle; skema baru bersebelahan dengan `POOLDATA`, satu pintu tulis ditegakkan lewat grant |

**27 accepted, 1 proposed** — ADR-0003 naik, ADR-0028 turun menunggu konfirmasi.

---

## 2. Apa yang ditemukan — 7 FINDING

Berkas lengkap `FINDING-001` … `FINDING-007`. Seluruhnya laporan kondisi sistem lama, netral, tanpa penilaian atas siapa pun.

| # | Temuan | Keadaan |
|---|---|---|
| 001 | Ambang kewenangan Komite dibandingkan tanpa konversi mata uang | BELUM TERBUKTI — menunggu REQ-011 |
| 002 | Alur bercabang berdasarkan identitas orang | TERBUKTI DARI XML |
| 003 | Baris Retensi Cedant hadir di `SpreadingRisk` saat rule tanpa penyaring membacanya | KEHADIRAN **MUNGKIN**, bergantung urutan tindakan pengguna |
| 004 | Hasil suntingan manual alokasi tidak terlindungi | BERSYARAT — berlaku bila hipotesis F2 (`.IsEditClaim` tanpa pembaca) benar |
| 005 | Tanggal kerugian dibandingkan terhadap akhir treaty dalam format berbeda | TERBUKTI DARI XML |
| 006 | Kurs yang tidak ditemukan dikembalikan sebagai `1` | TERBUKTI DARI DDL — besarannya belum terukur |
| 007 | Dua rumus premi reinstatement bekerja atas nilai yang berbeda | TERBUKTI DARI XML |

**Satu pola menghubungkan 004, 005, dan 007**: nilai dibaca dalam keadaan yang tidak dimaksudkan, dan tidak ada yang menjaganya.

---

## 3. Apa yang menunggu orang — lima kelompok

Sumber: `_selesai/OPEN-QUESTIONS.md` bagian A-MANUSIA dan `ASK-AKUNTANSI.md`. Setiap butir di sini sudah dibuktikan **tidak** terjawab oleh kolom mana pun di tabel work Pega maupun DDL `POOLDATA`.

| Kelompok | Butir | Apa yang ditanyakan |
|---|---|---|
| **Admin Pega** | A1, A2, A12, A15 | nama workbasket tujuan Akseptasi · nama peran/access group sebenarnya · alamat `prweb` mana yang sah untuk produksi · adakah node terpisah untuk bisnis syariah |
| **Pemilik proses** | A3, A6b, A7 | boleh/tidaknya satu orang Registrasi sekaligus Akseptasi · arti tiap nilai `PaymentType` · apakah `CloseClaimMD` memang boleh menutup klaim tanpa Komite |
| **Akuntansi** | A8 + `ASK-AKUNTANSI.md` (6 pertanyaan) | desimal nilai antara · bolehkah pembulatan mengubah angka lama · perhitungan lama yang keliru: dipertahankan atau diperbaiki · toleransi shadow-run · ambang nilai penghenti migrasi · tarif 2,5% / 2% / 2,2% masih berlaku · tiga butir lama (ambang Komite 15 atau 30, `TotalUR` selalu nol, cut-off tanggal 25) |
| **HRD** | A10b, A11b | apakah `VINCENTVERNANDO_1` akun uji atau akun bisnis nyata · status kepegawaian kelima nama di FINDING-002 |
| **Manajemen** | A13 | siapa berwenang menyetujui pengecualian pembekuan tambalan (ADR-0012) |

**A8 tidak punya jalur sendiri** — ia terjawab hanya lewat `ASK-AKUNTANSI.md` no. 5. **A12 dan E11 satu hal dilihat dari dua sisi.**

---

## 4. Apa yang menunggu basis data — 30 REQ

Register lengkap `ORACLE-REQUESTS.md` (**Aktif 29 · Selesai 1**). Daftar tarikan objek `pengetahuan/PULL-LIST.csv` (74 objek). Hasil yang sudah masuk `pengetahuan/SCHEMA-ACTUAL.csv` (698 baris).

### BLOCKER — 11

| REQ | Isi |
|---|---|
| 001 | Definisi `Rule-Obj-Property` di schema PegaRULES (`PR4_*`) — **akses belum ada** |
| 002 | `ALL_TAB_COLUMNS` untuk 44 tabel bisnis |
| 011 | Kuantifikasi Adjustment non-IDR vs ambang kewenangan Komite — memutuskan FINDING-001 |
| 012 | Isi `Data-Admin-DB-Table` untuk class `ASM-FW-%` — **akses belum ada** |
| 016 | Profil case buatan `VINCENTVERNANDO_1`, sekaligus menguji ADR-0023 |
| 017 | DDL 14 objek yang dirujuk dari dalam procedure/view tetapi tidak disertakan (`M_CURRENCYSTANDARD` dll.) |
| 018 | Cacah `PYID` ganda di tabel work **dan** baris ganda `(CASEID, TypeLoss, Currency)` di `OS_AKSEPTASI_KLAIM` |
| 021 | Grant `DATAPEGA` ke `POOLDATA`, dan adakah source `POOLDATA` yang menulis ke `DATAPEGA.PC_*` |
| 023 | Siapa menulis dan membaca ~60 kolom skalar `OS_AKSEPTASI_KLAIM` yang tidak diisi prosedurnya |
| 024 | `PYREOPENCOUNT > 0` — mengonfirmasi atau membatalkan ADR-0002 |
| 029 | Akseptasi ber-`STS_SUBJECTIVITY='1'` yang sudah punya catatan pembayaran |

### Selebihnya — 18 aktif

**PENTING (16)**: ~~REQ-004~~, 005, 006, 007, 008, 009, 013, 014, 015, 019, 022, 025, 026, 027, 028, 030 — *daftar 18 September. **REQ-004 dicabut** 19 September (`DEAD`), dan REQ-009, 019, 028 turun jadi verifikasi. Angka berjalannya di `ORACLE-REQUESTS.md`.*
**PELENGKAP (2)**: REQ-010, 020
**SELESAI (1)**: REQ-003 — terjawab lewat `pengetahuan/ddl/`

Enam di antaranya **mengambil alih** butir yang sebelumnya hidup di `_selesai/OPEN-QUESTIONS.md` — A5, A10a, A4, A9, A11a, A6a — dan menutupnya tanpa bertanya kepada siapa pun: **REQ-006, 016, 024, 025, 026, 027**.

Tujuh butir lain **sengaja belum ditanyakan** sampai recon Oracle masuk — `_selesai/OPEN-QUESTIONS.md` bagian B.

---

## 5. Apa yang belum tersentuh

### 5.1 Modul Komite

Folder `Komite Claim Non Prop` **tidak pernah dibuka** sepanjang sesi ini, atas instruksi. Sembilan butir menunggu di `_selesai/OPEN-QUESTIONS.md` bagian C, ditandai DEFERRED-TO-KOMITE-SESSION.

Yang sudah diketahui tentangnya **tanpa membuka folder itu**: Claim dan Komite duduk di baris-baris tabel yang sama — `PC_ASM_FW_GCNMFW_WORK` — dibedakan oleh `PXOBJCLASS`; penghubung induk–anak adalah kolom `PXCOVERINSKEY` (`BLUEPRINT.md` §8.4b).

### 5.2 Batas kepemilikan menurut ADR-0026

Di luar kepemilikan modul klaim, karena class-nya `-Int-`: `V_POLIS`, `T_STORAGE_IMAGE`, `EMAILKOMITE`, `M_LINK_SERVICE`. **Perlu dipahami, tidak dipindahkan.**

**298 kolom** di DDL tidak punya pasangan properti Pega dan **549 properti** tidak punya kolom (`BLUEPRINT.md` §19, keduanya angka batas atas). Sebagian besar milik modul lain.

### 5.3 Sisa pekerjaan teknis

| Hal | Keadaan |
|---|---|
| `pengetahuan/DDL_Script_ClaimNonProp2.xls` | **belum pernah dibaca** — muncul 2026-09-18 10:57 |
| `pengetahuan/RECON.sql` | perlu ditulis ulang untuk TOAD |
| `pengetahuan/PREFLIGHT.sql` | T1–T3 terjawab dari DDL; T4, T5, T7 masih menunggu |
| ADR-0003 | `proposed`, menunggu jawaban akuntansi |

### 5.4 Yang belum dimulai sama sekali, atas instruksi

Belum ada DDL sistem baru, belum ada Golang, belum ada React.

---

## Peta berkas

`STATUS.md` memuat daftar lengkap beserta golongannya. Yang perlu dibaca lebih dulu, berurut:

1. **`CONTEXT.md`** — glosarium; istilah yang dipakai di seluruh folder
2. **`RINGKASAN.md`** — berkas ini
3. **`BLUEPRINT.md`** — 19 bagian, pemahaman sistem lama selengkapnya
4. `docs/adr/` · `FINDING-00*.md` — keputusan dan temuan, satu berkas satu hal
5. `_selesai/OPEN-QUESTIONS.md` · `ORACLE-REQUESTS.md` · `ASK-AKUNTANSI.md` — yang masih terbuka
6. `pengetahuan/` — bahan yang masuk: sumber DDL, data mentah sapuan, dan alat tarik. Penjelasnya `pengetahuan/README.md`
7. `PENGETAHUAN.md` — seluruh 40 artefak di atas digabung jadi satu bacaan, untuk dibaca berurut atau dicari sekali jalan. **Turunan**: bila isinya berbeda dari berkas aslinya, yang asli yang berlaku

Akar berisi kesimpulan; `pengetahuan/` berisi bahan yang membantahnya. Keduanya sengaja tidak dicampur.
