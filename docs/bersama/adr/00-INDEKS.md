# Indeks catatan keputusan arsitektur

**42 catatan**, bernomor 0001-0042, tanpa lompatan dan tanpa nomor ganda.

Catatan keputusan arsitektur mencatat keputusan yang **mengikat kode** dan **mahal dibalik**.
Ia bukan milik satu modul: satu keputusan berlaku pada setiap modul yang menyentuh persoalan yang
sama. Keputusan tingkat modul - nama kolom, satu medan dimigrasi atau tidak, satu tabel dipecah
atau tidak - tetap tinggal di spec modulnya.

Kolom terakhir menunjukkan **berapa modul** yang merujuk catatan itu di dalam spec atau tiketnya.

| # | Judul | Dirujuk |
| ---: | --- | ---: |
| **0001** | Claim — Life dan Komite Life adalah dua konteks terpisah, dihubungkan kontrak child work | 9 |
| **0002** | RBAC Claim — Life memakai tiga peran yang sudah ada; rangkap peran ditolak | 7 |
| **0003** | Uang di Claim — Life tidak direpresentasikan sebagai `float` | 12 |
| **0004** | Alamat layanan keluar menjadi env var, bukan replikasi lookup `M_LINK_SERVICE` | 8 |
| **0005** | `IsPEGAPROD` tidak ditiru sebagai rule; diganti flag lingkungan eksplisit | 10 |
| **0006** | Penomoran klaim Life memanggil `PROC_GENERATE_SEQUENCE_NUMBER`; dua SQL lama tidak dimigrasi | 11 |
| **0007** | Jejak audit merekam siapa + kapan untuk setiap transisi status dan setiap jalur balik | 12 |
| **0008** | Efek keluar Claim — Life dijalankan asinkron, tidak memblokir alur, dengan antre-ulang | 6 |
| **0009** | Seluruh data Claim — Life dipindahkan; tidak ada koeksistensi dua penulis | 11 |
| **0010** | Penyimpanan berkas klaim tetap di Google Storage | 7 |
| **0011** | Unit status Claim — Life adalah **baris `AdjustmentList`**, bukan klaim | 6 |
| **0012** | Wewenang kirim ke Komite bergantung `Type` — dibawa apa adanya, dengan risiko RBAC tercatat | 3 |
| **0013** | Alamat layanan keluar di-resolve runtime dari `M_LINK_SERVICE` — bukan env var | 10 |
| **0014** | Keputusan komite hanya boleh diambil pemilik `KomiteID` pada tingkat berjalan | 6 |
| **0015** | Efek keluar Komite Claim Life **wajib berhasil** — transactional outbox, at-least-once | 11 |
| **0016** | Kolom uang memakai `NUMBER(38,8)` | - |
| **0017** | Daftar medan disusun dari aturan, bukan dari panduan bentuk dokumen | - |
| **0018** | Generasi polis adalah rantai baris, bukan penyuntingan di tempat | - |
| **0019** | Baris anak dipasangkan antar generasi lewat nomor urut, bukan kunci dagang | - |
| **0020** | Tabel proyeksi selisih bukan tabel sumber | - |
| **0021** | Rumus selisih hidup di lapisan layanan, dan nilai lama tidak pernah dihitung ulang | - |
| **0022** | Konversi tipe terjadi sekali saat masuk, bukan tiap kali dibaca | - |
| **0023** | Pemecah dokumen wajib punya penampung medan tak dikenal | - |
| **0024** | Pembatalan adalah generasi bernilai nol, bukan penghapusan | - |
| **0025** | Tabel inti berbagi kunci utama dengan tabel kerja | - |
| **0026** | Tabel kerja adalah akar lintas-lini, bukan tabel di samping | - |
| **0027** | Kolom basis data nullable; kewajiban isi ditegakkan di kode | - |
| **0028** | Sistem hilir membaca tabel, bukan muatan JSON | - |
| **0029** | Batas transaksi dipegang aplikasi; `COMMIT` tidak tertanam di teks SQL | - |
| **0030** | Aturan peran ditetapkan sekali dan berlaku lintas modul | - |
| **0031** | Penghapusan adalah penanda dan nilai balik, bukan hapus fisik | - |
| **0032** | Nama yang menyesatkan dibetulkan, disertai tabel pemetaan nama lama ke nama benar | - |
| **0033** | Nama skema ditulis eksplisit pada setiap query | - |
| **0034** | Uang melintasi batas stored procedure sebagai teks, dan dikembalikan ke desimal di dalam aplikas | - |
| **0035** | Tetapan operasional dibaca dari tabel, bukan ditanam di kode | - |
| **0036** | Medan kosong tidak hadir di dokumen JSON, dan ketidakhadirannya bukan kegagalan | - |
| **0037** | Kolom dipilih dari medan yang benar-benar diisi, bukan dari yang sekadar tampil | - |
| **0038** | Tangga persetujuan berakhir karena keadaan jenjang, bukan karena pencacah | - |
| **0039** | Rujuk atau salin: identitas dirujuk, jabatan disalin, data kutipan disalin utuh | - |
| **0040** | Pengaju dikeluarkan dari jenjang pertama - dan lubang di jenjang berikutnya diterima dengan sada | - |
| **0041** | Endorsemen Life berbagi tabel dengan new business dan dibedakan kolom versi | - |
| **0042** | Skema relasional dirancang baru bila ia tidak pernah ada | - |

---

## Riwayat penambahan

| Tanggal | Nomor | Asal |
| --- | --- | --- |
| 14-15 September 2026 | 0001-0015 | ronde Claim Life dan Komite Life |
| 23 September 2026 | **0016-0017** | ronde keputusan dua belas butir |
| 23 September 2026 | **0018-0024** | spec penyimpanan NB dan EDM Treaty In |
| 23 September 2026 | **0025-0032** | pola lintas modul dari 13 spec |
| 23 September 2026 | **0033-0042** | keputusan tingkat modul yang mengikat lintas modul |

## Cara memilih apa yang menjadi catatan keputusan

Dari **317 calon keputusan** di 13 modul, yang diangkat hanya yang memenuhi **keduanya**:

1. **Mengikat lebih dari satu modul** - atau akan mengikat, ketika modul serupa digarap.
2. **Mahal dibalik** - membalikkannya menuntut perubahan skema, migrasi ulang, atau perubahan
   kontrak dengan sistem luar.

Uji praktisnya: bila hanya satu modul yang punya persoalan itu dan membalikkannya cukup dengan
mengubah kode di satu tempat, ia **bukan** catatan keputusan arsitektur.

## Yang belum tercakup

Lima modul belum melewati grilling, spec, maupun tiket, sehingga keputusannya belum diketahui:
**NB Fakultatif** - **EDM Fakultatif** - **RNW Fakultatif** - **Treaty In** - **Treaty In
Adjustment**. Keluarga Fakultatif sendirian lebih besar daripada seluruh modul yang sudah selesai,
dan pola endorsemen di sana besar kemungkinan memunculkan catatan baru.

Dua modul lain juga belum: **Claim Non Prop** dan **Komite Claim Non Prop**.

*Indeks disusun 23 September 2026.*
