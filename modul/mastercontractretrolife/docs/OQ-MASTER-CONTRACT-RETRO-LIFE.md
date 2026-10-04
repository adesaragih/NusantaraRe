# Pertanyaan terbuka — Master Contract Retro Life

> Register OQ modul ini (awalan `OQ-MCRL-`), dibuka 30-09-2026 sesi implementasi. OQ lintas proyek tetap di
> `discovery/open-questions.md`. ⛔ Hanya work owner (atau pemilik yang disebut) yang menutup OQ; asisten mencatat **bawaan**
> yang dibangun sampai jawaban datang, dan bawaan itu dapat dibalik tanpa migrasi skema.
> Bukti setiap butir: `PARITAS-LAYAR-DAN-AKSI.md` dan `RALAT-DEV-30-09-2026.md`.

| OQ | Pertanyaan | Bawaan sampai dijawab | Pemilik | Status |
| --- | --- | --- | --- | --- |
| OQ-MCRL-01 | Gerbang tahun (tahun tanggal mulai kontrak = tahun treaty) di Pega ternyata **mati** (`SaveSecurityLife_Act` 4·5·6, `SaveSecurityReinsurerLife_Act` 4·5·6, `SaveBusinessLife_Act` 6 — `PRE=false`). Tetap ditegakkan sebagai penyimpangan sadar, atau ikut Pega? | **ikut Pega, tidak ditegakkan** (K7; bawaan brief gelombang 2 §10, diperiksa 01-10-2026 — jawaban belum ada). Catatan: tanggal kontrak adalah salinan read-only tanggal tahun (R5), jadi yang berbeda hanya terhadap `TREATYYEAR` yang diketik | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — gerbang tahun **ikut Pega** (tidak ditegakkan) |
| OQ-MCRL-02 | PK dan FK kelima tabel warisan tidak ada di DEV, padahal `ddl-tables-from-dba.md` menyebutnya sudah dipasang. Dipasang di lingkungan lain, atau memang belum pernah? | nol DDL; kaskade dan keunikan `ID` di Go (K1–K3) | DBA / work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — ⏳ **konfirmasi menyusul: DBA** |
| OQ-MCRL-03 | `SetErrorMessageReinsurer` memeriksa halaman non-life `InputTreatyReinsurer`, bukan medan layar life — di Pega batas 0..100 tidak pernah mengena. Ditegakkan, atau ikut perilaku efektif Pega? | **ditegakkan** 0..100 untuk `PCTSHARE`, `COMMISION` reinsurer dan `PCTSHARE` security (preseden OQ-TCO-17); `OVR_COMM` tidak | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — batas 0–100 **ditegakkan** |
| OQ-MCRL-04 | Teks pesan `SetErrorMessageBetween` (Rule-Message) tidak ikut diekspor. | kalimat Go sendiri berbahasa Inggris yang menyebut kolomnya: *"… must be between 0 and 100"* (sama dengan Treaty Contract Out) | pemilik ekspor Pega | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — ⏳ **konfirmasi menyusul: pemilik ekspor Pega (teks pesan)** |
| OQ-MCRL-05 | Objek fisik kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY` (sumber autocomplete `R/I RATE`, RD `BrowseRateLifeSummary`) tidak terbukti di korpus maupun katalog yang tercatat. | ⚠️ **digabung ke OQ-MCRL-13** (30-09-2026, paket 1): sebelum objeknya ditentukan, pembacaannya pun menunggu persetujuan | DBA | ✅ **ditutup 01-10-2026 — K1 keputusan work owner 01-10-2026**: objek fisik = view DEV `RATE_LIFE_SUMMARY` (6 kolom), dibaca saja `ID`, `USEDBY` |
| OQ-MCRL-06 | `Copy to all Reinstype` selalu `INSERT` baris baru tanpa pemeriksaan dobel — menjalankannya dua kali menggandakan business di kontrak sasaran. Ditambah penjaga dobel? | **ikut Pega** (selalu sisip); pratinjau menyebut kontrak sasaran dan jumlahnya sebelum konfirmasi | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MCRL-07 | Laporan kontrak total share ≠ 100 (tiket 11) tidak punya layar maupun tombol di Pega. Dibuat layar (dan tombol di mana)? | rute API baca saja `GET /api/master-contract-retro-life/laporan/total-share-bukan-100`; nol layar | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — laporan **hanya API**, tanpa layar |
| OQ-MCRL-08 | Kolom eksposur efektif security (tiket 06) tidak ada di grid Pega. Label dan letaknya? | kolom tambahan `EXPOSURE TO TREATY (%)` di grid security, sesudah `(%) SHARE` (dibangun paket 9+10 `a3c07bc`) | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MCRL-09 | Batas bawah > batas atas (tiket 02 AC 10) tidak diperiksa Pega. Ditolak? | **ditolak** 422, pesan menyebut mata uangnya | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MCRL-10 | Tanggal kontrak adalah salinan tanggal tahun (R5). Di Pega salinan itu tidak diperbarui saat tanggal tahun diubah. Diperbarui ikut induk? | **diperbarui** dalam transaksi yang sama (K4) | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MCRL-11 | Sel paginator grid reinsurer/security berlabel sisa `Tambah` membawa aksi `NewTreatyReinsurerDetail_Act` + `SetOutputParam_DT` yang hanya menyembunyikan form dan pesan. Dibawa sebagai tombol? | **tidak** — paginator dibangun sebagai penomoran halaman 10 baris (`pyRDLPageSize` 10; paket 9+10 `a3c07bc`) | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MCRL-13 | Sumber tabel rate layar business — autocomplete `R/I RATE` (kelas `…-RATE_LIFE_SUMMARY`) dan section `Rate List` (kelas `…-M_RATE_LIFE` = view atas `M_RATE_LIFE.JSONDATA`) — adalah view atas JSON produk rate. Penjaga Claim Life `TestMasterViewTidakDisentuh` (OQ-M7) menyatakan membacanya **menuntut persetujuan manusia**; izinnya terikat pada satu pembaca Claim Life, dan pencocokan nama substring menangkap kedua kelas. Disetujui untuk modul ini (kolom apa, objek fisik ringkasan apa)? | **tidak dibaca**: `GET …/ringkasan-rate` dan `GET …/rate` menjawab **503 berkalimat** (bukan daftar kosong); tombol `View Rate` dan medan `R/I RATE` tetap ada (XML). Akibatnya business **baru** tidak dapat disimpan (wajib-isi `RIRATEID`/`RIRATE`, `SaveBusinessLife_Act` b589); `Edit`, `Delete`, dan `Copy to all Reinstype` baris yang sudah ada tetap berjalan | work owner (+ pemilik penjaga Claim Life) | ✅ **ditutup 01-10-2026 — K1 keputusan work owner 01-10-2026** (*"izinkan membaca view rate, baca-saja"*), dibangun `d0c7b0a` *(riwayat: terbuka — **01-10-2026 gelombang 2 butir A1:** diperiksa, izin work owner belum tercatat di berkas ini; kedua rute tetap 503 berkalimat (brief gelombang 2 §3 A1, §10) — ✅ **ditutup 01-10-2026 — K1 keputusan work owner 01-10-2026** (*"izinkan membaca view rate, baca-saja"*): `RATE_LIFE_SUMMARY` dan `RATE_LIFE` dibaca saja, kolom RD saja; kedua rute 200; penjaga Claim Life `TestMasterViewTidakDisentuh` sudah dipersempit ke `modul/claimlife/` (`e13ad9e`))* |
| OQ-MCRL-12 | `Add` di layar security memanggil `NewInputSecurityLife_Act` yang mengosongkan halaman **reinsurer**, sehingga form security tampil berisi nilai terakhir (dan `Save` dapat menimpa baris yang tadi di-`Edit`). | form security **dikosongkan** saat `Add` (baris baru) — dibangun paket 9+10 (`a3c07bc`) | work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |
| OQ-MCRL-14 | Medan `Inputor` form baru di Pega menampilkan `OperatorID.pyUserName` (nama tampilan operator). Aplikasi baru mengenal akun pelaku (`X-Pelaku`), belum nama tampilannya. Sumber nama tampilan (login/direktori) disediakan bersama? | layar menampilkan **akun pelaku**; `USERID` tetap ditulis server dari pelaku, bukan dari isian (dibuka paket 9, 01-10-2026) | tim inti / work owner | ✅ **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** |

## Terjawab dari korpus (bukan dari work owner)

| Pertanyaan lama | Jawaban | Bukti |
| --- | --- | --- |
| Pertanyaan A ronde 2 — isi `RIRATE` | nama tabel rate (`USEDBY`); `RIRATEID` = ID tabel rate | R7 |
| Pertanyaan D ronde 2 — `HASIL12` | penampung ID baris yang dihapus | R8 |
| Pertanyaan B ronde 2 — validasi mati | wajib-isi hidup sudah mencakup medan halaman yang disunting; yang mati membaca halaman lain (sisa salin-tempel) — kecuali gerbang tahun (OQ-MCRL-01) | R9 |

## Keputusan work owner 01-10-2026 (`PROMPT-LANJUTAN-TIGA-MODUL-LIFE-KEPUTUSAN-OQ.md` §2)

*"Ikuti rekomendasi"* — 12 OQ ditutup dengan bawaan yang sudah dibangun; OQ-MCRL-13 + OQ-MCRL-05 diputuskan berbeda (K1, view rate
dibaca saja, `d0c7b0a`). Status lama dikutip:

| OQ | Status sebelum 01-10-2026 |
| --- | --- |
| OQ-MCRL-01 | terbuka |
| OQ-MCRL-02 | terbuka |
| OQ-MCRL-03 | terbuka |
| OQ-MCRL-04 | terbuka |
| OQ-MCRL-06 | terbuka |
| OQ-MCRL-07 | terbuka |
| OQ-MCRL-08 | terbuka |
| OQ-MCRL-09 | terbuka |
| OQ-MCRL-10 | terbuka |
| OQ-MCRL-11 | terbuka |
| OQ-MCRL-12 | terbuka |
| OQ-MCRL-14 | terbuka |

**Konfirmasi menyusul** (ditutup dengan bawaan; pihak yang disebut dapat membaliknya tanpa migrasi skema): OQ-MCRL-02 — DBA
(PK/FK kelima tabel warisan); OQ-MCRL-04 — pemilik ekspor Pega (teks `SetErrorMessageBetween`).

Masih terbuka: **nol** OQ modul ini.

## Keputusan work owner 01-10-2026 atas laporan K1 — penyimpangan sadar

Work owner 01-10-2026: *"ikuti rekomendasi semua"* atas laporan lanjutan keputusan OQ (commit `d0c7b0a`, `628c148`, `a64828f`, `9af0be9`).

| OQ | Hal | Keputusan | Status |
| --- | --- | --- | --- |
| OQ-MCRL-15 | Server menolak `RIRATEID` pilihan baru yang tidak ada di view `RATE_LIFE_SUMMARY`. Pega tidak memeriksanya. | **dipertahankan** — pemilih hanya menawarkan rate yang ada; pemeriksaan mencegah data rusak | ditutup 01-10-2026 |

## OQ-MCRL-02 diperbarui 01-10-2026

PK dan FK lima tabel tidak ada di DEV. Jawaban baru — keputusan work owner 01-10-2026: *"untuk modul Treaty Contract Retro Life aku izinkan kamu ALTER TABLE yang digunakan, perbaiki logic relasi datanya dan logic aplikasinya"*. Bawaan *"nol DDL; kaskade dan keunikan di Go"* **diganti**: PK, FK kaskade, dan FK penjaga salinan dipasang lewat migrasi 100 (lihat `RALAT-DEV-30-09-2026.md` §Ralat 01-10-2026). Ditutup.
