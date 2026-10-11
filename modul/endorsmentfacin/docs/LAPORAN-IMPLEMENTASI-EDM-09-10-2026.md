# Laporan implementasi EDM — putaran 09-10-2026 (E13 pembayaran, E14 sebagian)

> Lanjutan `LAPORAN-IMPLEMENTASI-EDM-01-10-2026.md` (putaran itu sudah di-commit work owner, `23110826`).
> Aturan yang berlaku tetap `PROMPT-LANJUT-IMPLEMENT-TANPA-TANYA-EDM.md`. **Belum di-commit** —
> commit, push, dan PR diurus work owner. Korpus: `D:\migrasi\RNM\Endorsment Fac In\` (READ-ONLY).
> Label mengikuti `CLAUDE.md` §4. Nol nomor polis dan nol nama orang di berkas ini dan di kode modul:
> `grep -rniE 'RNM-[A-Z0-9]'` → nol, `TestRegistryTanpaLiteralIdentitas` hijau.
> Sejak 01-10 tidak ada jawaban work owner atas butir §4 laporan sebelumnya. Itu dicek di
> `modul/nbfacin/docs/` (`_HANDOFF-EDM.md` bertanggal 21-09).

## 1. Ringkasan

| Tiket | Status | Isi putaran ini |
| --- | --- | --- |
| **E13** | sebagian | `CountPaymentEdm_Act` kini diport dari langkah 1 sampai blok 15, kecuali dua query: langkah 2 `SetRIComm_Act` (query master pajak; E17 dilewati) dan langkah 6 `GetEdmProdKe_Act` (satu `RDB-List`; hasilnya masukan `HasilProdKe`). Rinciannya: pembayaran per mata uang blok 13.2/14.4/15.5; `CountPaymentEdmTSIObj_Act` 1.4/1.5/2.2/langkah 3 (cicilan, `NetPremium`); langkah 5 RI slip; langkah 7 `CARI24`. **Rekonsiliasi kasus nyata blok 13** terhadap nilai tersimpan Pega |
| **E14** | sebagian | `DataTransform/CountPremiEDM_DT` (langkah 10) diport: `premiedmdt.go`. Lima berkas rumus lain belum |

## 2. Keputusan agent — menunggu konfirmasi

Nomor melanjutkan A11–A18 (laporan 01-10).

| # | Keputusan | Bukti |
| --- | --- | --- |
| A19 | Lokal Pega bertipe **`double`** diport sebagai **desimal eksak** (CLAUDE.md §7). Akibatnya, hasil berbeda dari Pega di digit belakang. | `[terverifikasi]` `pyLocalParameters` `CountPaymentEdm_Act`: 7 lokal `double` (comm, newpremi, oldpayment, premiins, preminew, prorateend, ricommins); TSIObj 0. Terukur pada kasus nyata: `EDMOldPayment` Pega `11242.674235452056` (17 digit), sedangkan hitungan desimal memberi `11242.67423545205479467564`; keduanya sama sesudah dibulatkan ke double (`TestPembayaranKasusNyataBlok13`) |
| A20 | Lokal activity dan Param DataTransform **bertahan antar-iterasi loop** `[dugaan]`. Baris mata uang tanpa pasangan OldData membaca nilai yang terbawa dari baris sebelumnya, yang bermata uang lain. Ini diport apa adanya (K-046) dan ditandai `BarisTerbawa`. Lokal yang belum pernah diisi → `ErrLokalBelumTerisi` | korpus 13.2/15.5 tanpa penolan, 14.4.4 tiga medan, TSIObj 1.4.4 empat medan |
| A21 | Batas port pembayaran: nilai baris sesudah langkah 2 (`SetRIComm_Act`) dan langkah 10 adalah **masukan**. Langkah 10 dijalankan pemanggil lewat `HitungPremiEDMDT` sebelum `HitungPorsiPembayaran`, supaya blok 13–15 tetap teruji sendiri | `premiedmdt.go` (kepala berkas) |
| A22 | `@divide` (bukan `@Math.divide`; A37 tidak mencakupnya) hanya diterima bila hasilnya **eksak** pada skalanya; selain itu `ErrModeDivideBelumTerverifikasi`. Sikapnya sama dengan paket `medis` | `CountPremiEDM_DT` cabang EdmType 1/2 |
| A23 | `HASIL2` kosong → galat (`cari24`). Literal `"0"` menjadi `models.HasilProdKePolisNBRNW`, dengan arti `[dugaan]` dari deskripsi langkah 7 | `CountPaymentEdm_Act` langkah 7 |

## 3. Temuan korpus

1. **Pega menghitung sebagian premi pembayaran dengan `double`** (A19). Ini pertanyaan untuk work owner/ADR:
   apakah sistem baru wajib meniru presisi `double` itu demi rekonsiliasi eksak, atau selisih di digit
   belakang diterima. Selama belum diputuskan, kode tetap desimal eksak. Test hanya *mengukur*
   selisihnya dengan `strconv.ParseFloat` sebagai alat ukur, bukan sebagai nilai aplikasi.
2. **`SumTotalPayment` bukan lokal double**. Nilainya selisih desimal dua medan tersimpan yang sudah membawa
   pembulatan double: `13566.16033544548 − 11242.674235452056 = 2323.486099993424` = nilai Pega.
3. **Medan absen dibaca 0 oleh Pega** `[dugaan]`. `Deduction2` tidak tersimpan di baris baru maupun lama fixture,
   tetapi `EDMOldDeduction2` tersimpan `"0"`. Kode produksi tetap menolak nilai kosong; test memasok 0.
4. **Langkah 5 mengubah gerbang**: RI slip (`EdmTypeNew==4`) menyetel `Type = 0` sebelum gerbang, sehingga
   Adj TSI dinilai sebagai Adj Ref No (`KasusSesudahRISlip`).
5. **Ketidaksimetrisan yang diport apa adanya (K-046)**:
   - netto **menambahkan** PPh+PPN;
   - `PaymentTotal` cicilan TSIObj **tanpa** −Deduction2, sedangkan blok 13–15 dengan −Deduction2;
   - `Deduction2ins` 14.4.6 tanpa bagian lama;
   - 14.4.3 dan TSIObj 1.4.3 tidak menulis `EDMOldDeduction2`;
   - `OldPremi` DataTransform (CARI24 0) tanpa Deduction2;
   - `EDMOldPremi` DataTransform EdmType 4 diisi Premium lama, walau nama Param-nya "OldPremiProrate".
6. **DataTransform**:
   - kedua `WHEN` tanggal dapat benar bersamaan, dan B (tanggal polis lama + `EdmDate` lama) menimpa A;
   - cabang EdmType 1 berjalan untuk **semua** lini;
   - kasus nyata NB-15 **ditolak**, karena tanggal polis lama (17:00 GMT, tidak dinormalisasi langkah 4) melawan tanggal baru 05:00 GMT memberi n,5 hari (A09).
7. **Tidak diport**: `CheckData.CARI1/2/3` (TSIObj 1.4.3), karena hanya muncul di berkas itu sendiri.

## 4. Menunggu work owner

1. **A19 — presisi `double` Pega** (temuan 3.1). Ini menentukan apakah rekonsiliasi eksak E22 mungkin dicapai.
2. **Satuan `@DateTimeDifference(…,D)` atas pecahan hari** (A09, tetap terbuka). Kini butir ini juga
   menghalangi DataTransform untuk data nyata.
3. **Mode pembulatan `@divide`** (A22).
4. **Nilai awal lokal/Param Pega** dan medan absen (temuan 3.3), yang menentukan perilaku baris mata uang baru.
5. Butir §4 laporan 01-10 yang belum dijawab tetap berlaku (A11–A18, tafsiran A18, OQ-037/040/046,
   pendaftaran modul untuk E16).

## 5. Uji

| | Hasil |
| --- | --- |
| `go vet ./modul/endorsmentfacin/...` | bersih |
| `go test ./modul/endorsmentfacin/...` | ok |
| `go test ./...` (APP_RNM, full suite) | 121 paket ok. Gagal **di luar modul ini**, tidak disentuh: `inti/backend/penjaga` (`premiumlistlife/.../rute_datapolis.go` mengimpor repository; label menu); `claimlife/backend/models` (`kodestatus`, `satutype` — menunjuk berkas `bordereaux.go`, `usulan.go`, `polis_datapolis.go` dll. milik modul lain); `TestKolomTakDibawaHanyaAdaDiKatalog` (`nbfacin/.../loader`). `grep endorsmentfacin` atas keluaran kedua penjaga → 0, dan tidak satu pun berkas yang disebut ada di modul ini |
| Rekonsiliasi kasus nyata (`edm-fire-1.json`, fixture NB-15 dibaca di tempat) | blok 13: 6 medan eksak, 6 medan sama-sesudah-double (2 di antaranya terbukti tidak eksak, A19), `SumTotalPayment` cocok rumus di kedua sisi; DataTransform: ditolak sesuai A09 |
| Uji mutasi | `pembayaran.go` 12/12 tertangkap; `premiedmdt.go` 8/8 tertangkap, ditambah 1 mutan ekuivalen (`WHEN A` OR→AND: setiap kali berbeda, `WHEN B` ikut benar dan menimpa); sebelum refaktor 2 mutan ekuivalen tambahan tercatat (1.4.1 nol `oldpremi`, 15.5.4 medan vs lokal) |
| Code review dua sumbu | kelompok pembayaran (E13): temuan diperbaiki, rinciannya di §6; kelompok DataTransform (E14): lihat §6 |

## 6. Code review

**Pembayaran (E13).** Diperbaiki:
- langkah 5 diport (`KasusSesudahRISlip`);
- data clump `(comm, brok, ded2, pph, ppn)` dijadikan `komponen`, dan prolog per baris digabung (`prologBaris`, `perBaris`, `isiSelisih`, `nolkanLama`);
- perbandingan share tanpa trik desimal;
- label `[terverifikasi]` dan "4 dari 7 lokal double";
- `CARI1-3` dicatat "tidak diport";
- `FlagEDMCurrency` didokumentasikan (false = tidak ditulis);
- literal `HASIL2` dipindah ke `models`;
- pengelompokan `RICommision` di test rekonsiliasi diralat (lokal `ricommins` double).

Tidak diubah, beserta alasannya:
- kunci map `lokal` tidak dijadikan konstanta. Sensus reviewer, dua cara, menunjukkan ke-14 kunci sama dengan ejaan korpus.
- `ErrLokalBelumTerisi` di 14.4.3/1.4.3 dipertahankan (korpus diam).

**DataTransform (E14) + langkah 5.** Reviewer spec tidak menemukan penyimpangan rumus terhadap korpus
(23 `pyDisabled` dihitung dua cara; `PassParameterPage=false` → Param mulai kosong). Diperbaiki:
- langkah 10 kini **dirangkai** di tempatnya: `LangkahSepuluh` di `HitungPorsiPembayaran`, sesudah langkah 4/9
  dan sebelum 11. Urutannya terbukti oleh `TestRangkaianLangkahSepuluh`: NetPremium dihitung dari keluaran
  DataTransform, lalu blok 13 menimpa;
- label dan perintah audit: angka 23, `[terverifikasi]` WHEN A/B, cabang CARI24, dan Deduction2. Rujukan relatif
  `.Policy.Payment.RICommision` diturunkan ke `[dugaan]`, karena halaman primernya belum terbukti `pyWorkPage`;
- loop lewat `perBaris`, `cari24` dihitung sekali;
- baris pertama tanpa pasangan → `ErrLokalBelumTerisi`, kini didokumentasikan dan diuji;
- pesan galat menyebut "lokal/Param";
- sisa `pyPropertiesValue` di FOR_EACH dicatat `[dugaan]`;
- langkah 5 tidak menulis balik `Type` (dicatat);
- literal RI slip dipindah ke `models` (`EdmTypeNewRISlip`, `TypeRISlip`). Dua cara bandingnya **setia korpus**:
  teks berkutip di gerbang E04, angka di langkah 5.

Tidak diubah: enam `time.Time` di `MasukanPremiEDMDT` (data clump; setiap medan memang halaman berbeda).

⚠️ Kedua reviewer standar melaporkan `go test ./inti/backend/penjaga/...` GAGAL di luar modul ini
(`modul/premiumlistlife/.../rute_datapolis.go` dan label menu). Itu di luar jatah, jadi tidak disentuh.

## 7. Berkas (belum di-commit)

Baru:
- `backend/services/pembayaran.go`, `premiedmdt.go` (+ test masing-masing);
- laporan ini.

Diubah:
- `backend/services/porsipembayaran.go`: `CacahMataUang` diganti baris pembayaran;
- `pembagian.go`: `kurangEksak` dan tipe `rumus`;
- `predikatkasus.go`: `KasusSesudahRISlip`;
- `gerbangtolak.go`: literal RI slip ke `models`;
- `backend/models/penawaran.go`: `HasilProdKePolisNBRNW`, `EdmTypeNewRISlip`, `TypeRISlip`;
- test `porsipembayaran_test.go`, `rekonsiliasi_test.go`, `gerbangtolak_test.go`;
- baris `Status:` E13 dan E14.
