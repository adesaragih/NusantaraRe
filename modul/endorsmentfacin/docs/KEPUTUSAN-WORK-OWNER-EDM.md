# Register keputusan work owner — modul Endorsment Fac In (`endorsmentfacin`)

> Satu baris per keputusan, nomor `W-EDM-nn` berurutan, tidak dipakai ulang. Keputusan agent yang menunggu
> konfirmasi tetap di laporan implementasi (A01–A23); bila dijawab, jawabannya dicatat di sini dan laporan
> itu diberi penunjuk.

| # | Tanggal | Pertanyaan | Keputusan | Akibat |
| --- | --- | --- | --- | --- |
| W-EDM-01 | 11-10-2026 | Irisan pertama | **I0 + I1 + I2**: pendaftaran modul + menu slot 966, form *Create Addendum*, tombol periksa gerbang (E04) dan tanggal (E05). Nol tulis DB selain slot menu | I3 (buat kasus `EDM-`) dan I4 (Old Data) menyusul |
| W-EDM-02 | 11-10-2026 | Tempat dan penomoran kasus `EDM-` (untuk I3) | **`T_WORK_POLIS`** (LINI `FAC`, ID `EDM-<n>`, sejalan ADR-0041) + **sequence baru** di migrasi **260**, angka awal diisi work owner (pola NB migrasi 181 `{NB_MULAI}`) | Nomor 260 dicadangkan untuk itu |
| W-EDM-03 | 11-10-2026 | Wadah daftar pengecualian nomor polis (pengganti literal `CheckEDMPolisDate` langkah 2 dan `SetValueToEDMWork` langkah 5) | **Tabel konfigurasi** baru, isinya diisi DBA (ADR-0035). Tabel kosong = tanpa pengecualian | Migrasi 261; nilai nomor polis tidak pernah masuk kode, uji, atau dokumen |
| W-EDM-04 | 11-10-2026 | Keputusan agent A11–A23 | **Ditunda** — tidak memblokir alur masuk; dibahas saat irisan premi/pembayaran (E12–E15, E22) | Status A11–A23 tetap "menunggu konfirmasi" |
| W-EDM-05 | 11-10-2026 | A30 — polis Life (pola RNML) ditolak "Invalid business" karena `Quotation.BusinessOldId` tidak terisi di portal | **"ini fokus ke edm facin, life tidak termasuk"** (dijawab di sesi asisten `nusantarare-0b`, diteruskan ke sesi utama) | A30 ditutup tanpa perubahan kode; penolakan polis RNML = perilaku yang dimaksud. Form tetap menyembunyikan medan untuk `Life` seperti Pega (tidak berdampak). Jalur Life E18/E19/E20 di luar fokus modul ini sampai diputuskan lain |
| W-EDM-06 | 11-10-2026 | A25 — zona jam `FACINPRODUCTION.BEGINDATE/ENDDATE` | **"zona waktu samain kaya NB facin"** (dijawab di sesi `nusantarare-0b`) | Jam dinding DATE dibaca **WIB**, konvensi penulis kolom itu (`modul/nbfacin/backend/services/binding.go:54` `tanggalFacInOffer`: `t.In(WIB)`). A25 diralat |
| W-EDM-07 | 11-10-2026 | Klep `OPENPROTEKSI_EDM` gerbang 1/3/4/5 (semula dianggap kosong) | **"openproteksi_edm dipakai dong"** (di sesi `nusantarare-0b`), disertai berkas baru `D:\migrasi\RNM\DDL\OPENPROTEKSI_EDM.txt` dan `BrowseOpenProteksiEdm_RD.xml` | Klep dibaca dari tabel warisan `POOLDATA.OPENPROTEKSI_EDM` (baca saja, nol migrasi) menurut saringan report itu; A40–A41 |
| W-EDM-08 | 11-10-2026 | Angka awal nomor `EDM-` dan izin tulis `T_WORK_POLIS` (I3) | **"EDM-00001 AJA, NNTI INI BISA DIRUBAH LAGI KAN?"** dan **"SETUJU UTK menulis ke tabel T_WORK_POLIS"** (di sesi `nusantarare-0b`) | Migrasi 260 `SEQ_WORK_POLIS_EDM` START WITH 1; ID = `EDM-` + 5 digit (A42). ⚠️ Kasus EDM Pega sudah melewati nomor 10000 (mis. `EDM-13445` di fixture NB) → bentuk 5 digit akan sama dengan ID Pega mulai `EDM-10000`; mengubahnya = migrasi sequence baru / kode Go |
| W-EDM-09 | 11-10-2026 | LINI kasus `EDM-` di `T_WORK_POLIS` — LINI `FAC` bocor ke portal, kotak masuk, dan buka kasus NB (NB hanya menyaring `LINI='FAC'`: `modul/nbfacin/backend/repository/portal.go:44`, `kotakmasuk.go:87`) | **LINI `FACEDM`** | Mengganti "LINI FAC" di W-EDM-02; NB tidak diubah |
| W-EDM-10 | 11-10-2026 | `STATUS_WORK` awal | **`Pending-Policy`** (sama NB) | — |
| W-EDM-11 | 11-10-2026 | `POSITION` awal | **`ReasFacInMarketing`** (workbasket `Assignment7`, pola NB POSITION = workbasket) | — |
| W-EDM-12 | 11-10-2026 | Gerbang 2 tanpa `GetListEdm` | **Periksa kasus sistem baru**: kasus `EDM-` LINI `FACEDM` yang belum `Resolved-*` atas polis yang sama menolak; kasus EDM era Pega dicatat "tidak diperiksa"; Create boleh | `lengkap` tidak lagi diturunkan G2 |

*Tanpa nama orang, tanpa nomor polis, tanpa data nasabah.*
