# Hasil implementasi — modul `nbtreatyin` (NB Treaty In)

> Putaran 2 · konsolidasi paket P1–P9 di cabang integrasi `modul/nbtreatyin/implementasi` (`71842cb0`), dirangkum
> paket P10 pada 04-10-2026. Setiap bukti di bawah diperiksa ulang ke kode dan uji cabang itu — bukan disalin dari
> laporan paket.
>
> Status: ✅ lulus (ada bukti uji/kode) · 🟡 sebagian (penahannya disebut) · ⛔ tidak dibangun, dengan alasan
> (a) tidak terjangkau / nol efek, (b) keputusan work owner tertulis, atau (c) data hanya di JSON di luar K8
> (PROMPT putaran 2 bab 4) · 📄 butir dokumen · 🔒 `[terbuka]`, menunggu pihak lain.
>
> "Bukti" menunjuk uji (`berkas:baris NamaUji`) atau kode (`berkas:baris`). Jalur relatif terhadap
> `modul/nbtreatyin/backend/` kecuali diawali `frontend/`, `docs/`, atau `inti/`.
>
> ⛔ Uji bertag `db` (Oracle sungguhan) **belum pernah dijalankan**: skema uji K11 kosong. AC yang buktinya hanya
> uji db ditandai 🟡 dengan penahan **K11**. Bila uji db-nya **belum ditulis** (atau tidak memeriksa yang dituntut
> AC), hal itu disebut terang — itu pekerjaan yang masih dapat dikerjakan (bab 9).

## Isi

1. Ringkasan
2. `spec.md` — 96 AC
3. `spec-penyimpanan-relasional.md` — 66 AC
4. Perubahan sejak putaran 1
5. Penyimpangan sadar
6. Rule yang tetap tidak dibangun
7. Jalur Proporsional dan NonProporsional — dari layar sampai tersimpan
8. Tabel: delapan `CREATE TABLE` dan tabel warisan
9. Status tiket 00–23 dan pekerjaan yang masih dapat dikerjakan
10. Perintah verifikasi

## 1 · Ringkasan

Dihitung dari tabel bab 2–3 (rentang dijabarkan; setiap nomor tepat satu kali).

| | ✅ | 🟡 | ⛔ | 🔒 | 📄 | jumlah |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `spec.md` | **80** | 10 | 2 | 0 | 4 | **96** |
| `spec-penyimpanan-relasional.md` | **51** | 14 | 0 | 0 | 1 | **66** |
| *putaran 1 — `spec.md`* | *75* | *12* | *5* | *0* | *4* | *96* |
| *putaran 1 — penyimpanan* | *39* | *14* | *11* | *1* | *1* | *66* |

Cacah penyimpanan = nomor AC berurutan **1–66**; sisipan **`20b`** (koreksi presisi 23-09-2026) dicacah bersama
nomor 20 (`docs/issues/00-PETA-AC-PENYIMPANAN.md`, RALAT P6 03-10-2026: 67 → 66).

Penahan setiap 🟡:

| Penahan | `spec.md` | penyimpanan |
| --- | --- | --- |
| **K11** — uji db ditulis, belum dijalankan | 23, 29, 31, 68 | 1, 6, 13, 46, 49, 50, 51, 55 |
| **K11** + uji db **belum ditulis / tidak memeriksa yang dituntut** (bab 9) | 17, 72, 83, 89 | 9, 12, 45 |
| **K12** — pemetaan peran IAM | 12, 45 (bersama K7) | — |
| keputusan work owner (PERMINTAAN F) | 68 (F7, pemuatan di produksi — bersama K11) | 59 (F3) |
| belum ada RALAT / bukti dokumen (bab 9) | — | 15, 25 |

## 2 · `spec.md` — 96 acceptance criteria

| AC | Status | Bukti / alasan |
| ---: | :---: | --- |
| 1 | ✅ | `models.Disetujui` (`models/tangga.go:127`); `models/tangga_test.go:12` TestIsApprovedNolDitolak; `handlers/alur_test.go:155` TestAdminMenolakDiselesaikanDitolak |
| 2 | ✅ | `models/tangga_test.go:18` TestIsApprovedSelainNolDisetujuiTermasukKosong |
| 3-4 | ✅ | `models/tangga_test.go:26` TestIsApprovedDibandingkanSebagaiTeks — teks, menurut DecisionTable `isApproved` (bukan When bernama sama) |
| 5 | ✅ | `models/tangga_test.go:41` TestTanggaAdminMenolakSelesaiDitolak; `handlers/alur_test.go:155` (Resolved-Rejected; kirim ulang 409) |
| 6 | ✅ | `models/tangga_test.go:51` TestTanggaAtasanMenolakKembaliKeAdmin; `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali |
| 7 | ✅ | `models/tangga_test.go:63` TestTanggaAdminMenyetujuiNaikKeSecHead; `handlers/alur_test.go:202` |
| 8 | ✅ | `models/tangga_test.go:73` TestTanggaSecHeadMenyetujuiSelaluNaikKeDeptHead — **K2**: WO lama berlaku; batas 200 juta `CekLimitTreatyAcc_Act` tidak dibangun (b), pertentangannya di tiket 03 |
| 9 | ✅ | `models/tangga_test.go:84` TestTanggaDeptHeadMenyetujuiSelesai; `handlers/alur_test.go:202` (alur penuh sampai selesai) |
| 10 | ✅ | `models/tangga_test.go:95` TestPosisiBuanganTidakAda; `When/isClaimTreaty` (Decision5 tanpa connector masuk) tidak dibangun (a) |
| 11 | ✅ | antrean = workbasket (`services/layanan.go:80` `anggota`); `handlers/alur_test.go:139` TestBukanAnggotaAntreanDitolak; `handlers/portal_test.go:49` TestGerbangDaftarPortal (wadah `Section/SFAPortal_OpportunitiesList`, P8) |
| 12 | 🟡 | wewenang = keanggotaan antrean (`inti.Pelaku.Peran`), nol pemeriksaan identitas di kode; kedua belas tempat identitas XML terdaftar di `models.DaftarTempat` (`models/peran_tempat.go:90`), pemetaannya `PemetaanPeranTempat` **kosong** (`:123`, K16) ⇒ tertunda — penahan **K12 (IAM)** |
| 13 | ✅ | peran dari `inti.Pelaku.Peran` (workbasket akun); pemetaan tempat = konstanta `models/peran_tempat.go` (K16 — migrasi 330 dihapus); nol kolom telepon: `When/IsSPVTreaty1`, `When/IsTreaty1` (`OperatorID.pyTelephone`) tidak dibangun (b) K5 |
| 14 | ✅ | `services/layanan.go:80` `anggota` = `PunyaPeran(nama)`; `handlers/alur_test.go:139` (urutan peran terbalik tetap anggota); `When/IsUW` (nomor urut `pyWorkBasketList(1)`) tidak dibangun (b) |
| 15 | ✅ | RALAT K8: kontrak dari view (`repository.DetailKontrak`, `repository/acuan.go:144`); JSON master **hanya** `repository.MasterXOLDariJSON` (`repository/masterxol.go:74`), medan dari daftar `models/masterxol.go`; `repository/masterxol_test.go:39` TestUraiMasterXOLHanyaMedanK8 |
| 16 | ✅ | nol penulisan JSON: penulisan hanya ke 8 tabel diagram + `T_WORK_POLIS` + `HISTORYAKSEPTASIPEGA` + `HISTORYAKSEPTASIPRODUCTION` (+ `GENERATE_SEQUENCE_NUMBER` lewat penomor inti, bab 5 butir 24); `handlers/nonprop_test.go:82` TestNonPropPilihBisnisHitungSimpanBacaKembali (halaman master tidak tersimpan); `SaveJsonPolisTreatyIn_Act`, `SavePolisTreatyIn_SQL`, `SaveTreatyIn` tidak dibangun (b) |
| 17 | 🟡 | 33 kolom RD dibaca dari view (`repository/acuan.go:34` `KolomRDDetail`); uji lawan view sungguhan **belum ditulis** (bab 9) dan tidak dapat dijalankan (K11) |
| 18 | ✅ | medan uang dibaca apa adanya (`TO_CHAR` TM9, `repository/kolom.go:147`), tanpa hitung ulang saat dibuka; `handlers/nonprop_test.go:82` (NetPremium tetap 2700 sesudah simpan) |
| 19-21 | ✅ | `models/penggolong_test.go:185` TestPenggolongBerhentiDiBarisPertama, `:193` TestPenggolongBawaanUnknown, `:145` TestKe128KodeSamaDenganSistemLama |
| 22 | ✅ | syarat When yang DIJALANKAN — `docs/INVENTARIS-XML.md` bab 11 |
| 23 | 🟡 | `NUMBER(38,8)`; `repository/kolom_test.go:21` TestPecahAngkaEksak; pulang-pergi `repository/polis_db_test.go:64` (`830.82191780804`) — **K11** |
| 24 | ✅ | nol pembulatan di repository (`pecahAngka` eksak, `repository/kolom_test.go:21`); pembulatan hanya di penyajian (`frontend/sajian.ts`, `frontend/sajian.test.ts`) |
| 25 | ✅ | `apd` di seluruh paket, nol `float` di kode `backend/`; `models/dokumenlama_test.go:314` TestAngkaJSONTidakLewatFloat, `models/nilaipega_test.go:11` TestTeksSkalarJSONTanpaFloat; penjaga `inti/backend/penjaga/migrasi_test.go:179` TestKolomUangDesimalDanNolJSON |
| 26 | ✅ | RALAT K3: `Deduction1/2` = uang (`models/katalog.go:183` `kUang`; `frontend/medan.test.ts` "AC 85 + K3"); `BROKERAGE`/`RNM_SHARE` view: **nol pembaca** di 176 rule terjangkau (`docs/alat/pemakai.py`; tiket 07 RALAT P4) ⇒ tidak dibaca, tak satu pun diperlakukan sebagai uang |
| 27-28 | ✅ | `models/hitung_test.go:89` TestPajakBrokerageInclusiveDibagi1022, `:99`, `:134` TestTypeTaxHurufKecilMenggeserBalance; `models/setppnpph_test.go:20` TestSetPPNPPHNilaiXML |
| 29 | 🟡 | satu transaksi (`services/gudang.go:98` `DalamTransaksi`); `handlers/alur_test.go:355` TestSatuTransaksiPembatalanUtuh (tiruan); Oracle `repository/polis_db_test.go:187` — **K11** |
| 30 | ✅ | `db.Qualify` di setiap SQL; penjaga `inti/backend/penjaga/lintasaplikasi_test.go:541` TestNolNamaTabelTelanjangDiQuery |
| 31 | 🟡 | indeks unik fungsi `CASE (NOPOLIS, PRODKE)` (`migrations/320_t_general_polis.sql:109`) + penomor `FOR UPDATE`; `repository/polis_db_test.go:166` TestNomorPolisSekaliDanUnik — **K11** |
| 32-33 | ✅ | kolom `DATE` (DDL 320–327); satu format tukar `YYYY-MM-DD[ HH:MI:SS]` (`repository/kolom_test.go:75` menolak format kedua); tampilan satu format `frontend/sajian.test.ts` (AC 33) |
| 34-35 | ✅ | `models.PraprosesAdmin` (`models/polisaturan.go:36`); `models/tangga_test.go:125` TestPraprosesAdminTanggalKosongHariIni. `SystemSetOneYear_DT` hanya saat admin mengubah StartDate (`models/tanggal_test.go:10`, `handlers/tanggal_test.go:16`) — tidak menyentuh pengisian tanggal kosong (spec §9.2 butir 12) |
| 36-38 | ✅ | `handlers/alur_test.go:435` TestPilihBisnis (422, nol simpan); `handlers/nonprop_test.go:207` TestNonPropMasterTidakAdaMenghentikanPilihBisnis (master tidak ada / rusak → 422) |
| 39-41 | ✅ | `HISTORYAKSEPTASIPEGA.OPERATORID` = identitas login, `USERNAME` = nama tampilan; `HISTORYAKSEPTASIPRODUCTION.PIC` = nama tampilan — `handlers/alur_test.go:155` |
| 42 | ✅ | `models.PraprosesAtasan` (`models/polisaturan.go:94`) memakai nama tampilan (P33); `NamaTampilan` tanpa jatuh-balik ke login (`repository/riwayat.go:146`) |
| 43 | ✅ | satu baris riwayat per submit — `handlers/alur_test.go:202` (5 submit → 5 baris) |
| 44 | ✅ | K5: NBStatus dari data (nama tampilan / nama posisi) — `models/tangga_test.go:116` TestTeksNBStatus, `models/layar_test.go:215` TestPascaAdminTolakNBStatus, `handlers/alur_test.go:155`, `:202` |
| 45 | 🟡 | 25 medan wajib dua layar realisasi + `ListSuggest` ditegakkan saat Submit DAN Save (`models.MedanWajibBerlaku`, `models/layar.go:173`; `handlers/alur_test.go:314`); `ProductionDate` lewat tempat berperan — penahan **K12** (mekanisme: `handlers/alur_test.go:534`); `DateofSurvey` ⛔ (b) **K7**: layar survei tidak dibangun, tidak ada tabel diagram grilling |
| 46-47 | ✅ | `models/layar_test.go:45` TestMedanWajibAtasanTanpaEnamMedanAdmin |
| 48 | ✅ | `handlers/alur_test.go:314` TestMedanWajibMenahanKirimDanSimpan (Submit dan Save) |
| 49-51 | ✅ | `models.GabungMasukanLayar` (`models/layar.go:275`, daftar izin per posisi); turunan dihitung ulang server; `handlers/alur_test.go:369` TestMedanTerkunciAtasanDanTurunanAdmin, `models/layar_test.go:137` |
| 52 | ✅ | RALAT P2: atasan mengisi Approval + Suggest (+ ProductionDate bersyarat tempat); `models/layar_test.go:167` TestAtasanHanyaMengisiPutusanCatatanDanTanggalProduksiBersyarat, `handlers/alur_test.go:369`, `frontend/medan.test.ts` "AC 52" |
| 53 | ✅ | elemen `1=2` / `NEVER` tidak dibangun — `frontend/medan.test.ts` "elemen mati" |
| 54 | ✅ | `repository.DaftarMataUang` `CURRENCY <> 'ITL'` (`repository/acuan.go:435`; RD `BrowseCurrencyTreatyIn_RD` filter B) |
| 55-56 | ✅ | kode tampil apa adanya; medan berpilihan "associated values" (tidak ada di korpus) = isian teks, diterima dan disimpan (tiket 11) |
| 57 | ⛔ | (b) **K8 butir 4**: treaty keluar seluruhnya JSON `M_TREATY_OUT` — `RDBList/BrowseTreatyOut`: `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`; tidak dibaca |
| 58 | ✅ | nol penulisan treaty keluar / master; jalur XOL/NonProp dibangun ke `T_POLIS_XOL`, `T_POLIS_XOL_LAYER`, `T_POLIS_INSTALMENT(_DETAIL)`, `T_POLIS_SPREADING` — `handlers/nonprop_test.go:82` |
| 59 | ✅ | RALAT P4: `TreatyRealizationCheckDuplicate` (`Protect Duplicate Policy; data is similar to …`) — `handlers/alur_test.go:336` TestPolisSerupaMenahan, `handlers/logika_test.go:151` TestCekPolisSerupaHanyaSaatAdminMenyetujui |
| 60 | ✅ | penjaga `inti/backend/penjaga/lintasaplikasi_test.go:491` TestHandlersTidakMengimporRepository |
| 61-63 | ✅ | rule yatim, pembongkar JSON (RALAT K8: empat tetap tidak dimigrasi, dua dimigrasi baca-saja), `When/isApproved` — `docs/INVENTARIS-XML.md` bab 2, `docs/alat/status.json` |
| 64-65 | ✅ | kolom `isApprovedtoDeptHead` tidak dibuat; `LetterNo` bukan penanda arah (tidak dibaca `models.Langkah`) — `repository/kolom_test.go:155` TestKatalogSepakatDenganDDL |
| 66 | ⛔ | (a) RALAT P4: `isFOR` dikirim **kosong** oleh `Flow/InputRealizationTreatyIn` Utility1; `SaveJsonPolisTreatyIn_Act` langkah 3 (EDM) `//`, langkah 5 (POLICY) When tak dicentang — nol efek; NB hanya jalur polis baru |
| 67 | ✅ | `models.KodeLiniJiwa` (`models/penggolong.go:115`); `models/penggolong_test.go:208` TestLiniJiwaEnamBelasKode |
| 68 | 🟡 | pemuat dokumen lama (tiket 22) menulis lewat `SimpanHalaman`, berkas dibuka `BacaHalaman` — `repository/lama_db_test.go:74` — **K11**; pemuatan di produksi menunggu keputusan WO (PERMINTAAN F7) |
| 69 | ✅ | `models.PecahDokumenLama` (`models/dokumenlama.go:379`): EndDate kosong = StartDate — `models/dokumenlama_test.go:93` |
| 70 | 📄 | penyimpangan dicatat di kode, tiket, dan bab 5 berkas ini |
| 71 | ✅ | K4: catatan ke `HISTORYAKSEPTASIPRODUCTION` (TGL_INP, PIC, AKSES_LOGIN) dan dibaca balik untuk layar — `handlers/alur_test.go:155`, `models/usulan_test.go:76` TestBarisCatatanDariRiwayatProduksi |
| 72 | 🟡 | `DaftarRiwayat` `ORDER BY TGL_TRANSFER, ROWID` (`repository/riwayat.go:51`); uji db **belum ditulis** (bab 9) — K11 |
| 73 | ✅ | `models.RakitNomorPolis` (`models/polisaturan.go:252`); `models/tangga_test.go:216` TestRakitNomorPolis; `handlers/alur_test.go:202` (bentuk) |
| 74 | ✅ | nomor kedua = nomor pertama — `handlers/alur_test.go:202`; `SetelNomorPolis` hanya bila `NOPOLIS IS NULL` (`repository/polis.go:241`) |
| 75 | ✅ | `models/tangga_test.go:162` TestNonProporsionalDitandai |
| 76 | ✅ | `models/tangga_test.go:184` TestHasFacOutHanyaDuaKode |
| 77 | ✅ | halaman tidak membawa pesan antar permintaan; validasi kirim dihitung atas salinan |
| 78 | ✅ | `SetValidateInstallment` langkah 1 (`Page-Clear-Messages`) tidak dibangun — `models/angsuran.go:93` |
| 79 | ✅ | setiap rumus berkutip `PropertiesValue` langkahnya (`models/hitung.go:7`, `models/nonprop.go:15`); `models/hitung_test.go`, `models/nonprop*_test.go` |
| 80 | ✅ | layar Dept Head dapat disimpan dan disubmit — `handlers/alur_test.go:202`, `handlers/nonprop_test.go:82` |
| 81 | ✅ | 12 tempat terdaftar, pemetaan kosong ⇒ tertunda — `handlers/alur_test.go:505` TestTempatBerperanTidakDitebak, `models/peran_tempat_test.go:19`, `:37` |
| 82 | ✅ | arah `MUNCUL`/`KECUALI` tidak ditebak — `models/peran_tempat_test.go:95` TestTempatTampilMenurutArah, `frontend/tempat.test.ts` |
| 83 | 🟡 | gagal `CatatRiwayat`/`CatatUsulan` membatalkan submit (`handlers/alur_test.go:355`, tiruan); lawan Oracle: uji db **belum ditulis** (bab 9) — K11 |
| 84 | ✅ | K6: tabel keputusan kosong = disetujui (`models/tangga_test.go:18`); layar: submit tanpa Approval 422 di tiga jenjang (`handlers/logika_test.go:124` TestSubmitTanpaApprovalDitolakDiSetiapJenjang) |
| 85 | ✅ | `KotakMedan` kode mata uang di samping angka — `frontend/medan.test.ts` "AC 85 + K3" |
| 86 | ✅ | RALAT K14: format sel Section (`pyDecimalPlaces`, `pySeparators`) — `frontend/sajian.ts`, `frontend/sajian.test.ts`, `frontend/medan.test.ts` |
| 87 | ✅ | KEPUTUSAN-RONDE-12 butir 7: muatan 4 medan sesudah commit, gagal tidak membatalkan — `handlers/alur_test.go:662` TestKonversiSesudahSelesai; sambungan `[terbuka]` (PERMINTAAN C1) |
| 88 | ✅ | rule tak terjangkau tidak dibangun — `docs/INVENTARIS-XML.md` bab 1–2 |
| 89 | 🟡 | kolom view yang hilang = galat (`repository/kolom_test.go:229` TestKolomViewHilangAdalahGalat); keberadaan kolom lawan view sungguhan: uji db **belum ditulis** (bab 9) — K11 |
| 90 | ✅ | tabel berejaan ganda (`HISTORYAKSEPTASIPEGA` dst.) lewat `Qualify` satu skema |
| 91 | ✅ | nol peran karangan; uji memakai peran `UJI-` (`handlers/alur_test.go:505`, `:534`) |
| 92 | ✅ | berkas menunggu posisi — `handlers/alur_test.go:139`, `handlers/portal_test.go:49` |
| 93-95 | 📄 | butir dokumen spec |
| 96 | ✅ | nol nama orang di artefak: riwayat cabang ditulis ulang (K1, 0 dari 195 baris tersisa — PROMPT putaran 2 bab 1); fixture `UJI-` (P9 A1) |

## 3 · `spec-penyimpanan-relasional.md` — 66 acceptance criteria

| AC | Status | Bukti / alasan |
| ---: | :---: | --- |
| 1 | 🟡 | indeks unik fungsi `CASE` (`migrations/320_t_general_polis.sql:109`); `repository/polis_db_test.go:166` — **K11** |
| 2-5 | ✅ | `SisipKasus` PRODKE 0, `OLD_POLIS_ID` NULL + `UNIQUE` (`migrations/320_t_general_polis.sql:106`), PK bersama `T_WORK_POLIS` — `repository/kolom_test.go:310` TestTabelDanKolomMengikutiDiagramGrilling |
| 6 | 🟡 | generasi tertutup = ada penerus (`repository.syaratTerbuka`, `repository/polis.go:135`; kolom `TGL_TUTUP` dibuang P1) → `ErrGenerasiTertutup`; `repository/polis_db_test.go:187` — **K11** |
| 7 | ✅ | nol kunci tamu NB↔EDM di `T_WORK_POLIS` |
| 8, 10 | ✅ | `NOURUT` + `UNIQUE (induk, NOURUT)` di setiap tabel anak (DDL 322–327) — `repository/kolom_test.go:310` |
| 9 | 🟡 | hapus-sisip menomori ulang 1..n (`repository.SimpanHalaman`, `repository/polis.go:67`); uji db `repository/polis_db_test.go:64` menghapus baris tetapi **tidak membaca kolom `NOURUT`** atas tiga angsuran (bab 9) — K11 |
| 11 | ✅ | NB: tidak ada pemasangan antar generasi (milik EDM) |
| 12 | 🟡 | kode tetap `VARCHAR2` (katalog golongan kode); uji db memeriksa `BizCode "006"` dan `BusinessOldId "01"`, **bukan** `GROUP_PANEL` (bab 9) — K11 |
| 13 | 🟡 | `BusinessOldId "01"` pulang-pergi `repository/polis_db_test.go:64` — **K11** |
| 14 | ✅ | `models/penggolong_test.go` (006/01 → FireStyle2); `handlers/alur_test.go:435` (atas kolom tersimpan) |
| 15 | 🟡 | Oracle menyimpan `''` sebagai NULL (RALAT tabel bab 6); dibaca kembali `""` (`repository/kolom_test.go:93` TestNilaiBaca), `""` ≠ `"0"`; bunyi baru AC belum ditulis (bab 9) |
| 16 | ✅ | `models.Disetujui`; `models/tangga_test.go:12`, `:18` |
| 17-18 | ✅ | `repository/kolom_test.go:43` TestNilaiTulisKosongJadiNULL |
| 19-20 (+20b) | ✅ | `NUMBER(38,8)` (DDL); pengikatan tanpa pemotongan (`repository/kolom_test.go:21`, 11 desimal diterima); pembulatan pada skala 8 diperiksa `repository/lama_db_test.go:74` (`592629512.88000028`, K11) |
| 21-22 | ✅ | `models.BacaTanggalLama` (`models/nilaipega.go:80`); `models/dokumenlama_test.go:14` TestBacaTanggalLama, `:35` TestTanggalAmbiguTidakDitebak (K15) |
| 23-24 | ✅ | kolom `DATE`; nol FLOAT (penjaga `inti/backend/penjaga/migrasi_test.go:179`) |
| 25 | 🟡 | port rumus meniru pembandingan Pega — persen lawan 100 dan tanda lawan 0 (`models/hitung.go:196`, `:221`, `:300`–`:417`; `models/angsuran.go:71`, `:117`); belum ada bukti tertulis bahwa tak satu pun membandingkan **dua nilai uang**, atau RALAT AC (bab 9) |
| 26 | ✅ | tanpa `LAYER*` di `T_GENERAL_POLIS` (diagram F26) — `repository/kolom_test.go:310`; `handlers/alur_test.go:435` (dibaca balik dari view) |
| 27 | ✅ | `T_POLIS_QUOTATION` = 10 medan diagram J37 + 6 RALAT berbukti XML (`docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 2) — `repository/kolom_test.go:310` |
| 28 | ✅ | `T_POLIS_CEDING` `CEDING_CO_ID` + `CEDING_CO_NAME` (`migrations/322_t_polis_ceding.sql`, diagram R43), di bawah `T_POLIS_QUOTATION` (O39) |
| 29-30 | ✅ | disalin apa adanya (`models/dokumenlama_test.go:93`); ceding dihapus per baris |
| 31, 33 | ✅ | `models.PeriksaBentukSimpan` (`models/katalog.go:430`); `models/layar_test.go:12` TestBentukProporsionalMenolakXOLDanRincian |
| 32 | ✅ | akibat AC 33 |
| 34-35 | ✅ | `LAYER*` hanya di `T_POLIS_XOL_LAYER`; `DEDUCTION` uang (DDL 326/327); `handlers/nonprop_test.go:82` |
| 36 | ✅ | `models/layar_test.go:34` TestSpreadingBagiRataPresisiSepuluh |
| 37 | ✅ | persen `NUMBER(38,8)` |
| 38 | ✅ | RALAT K3: `DEDUCTION1/2` uang (`models/katalog.go:183`); `TOTAL_SHARE_PERCENTAGE_*` persen turunan (`models.HitungTotalSpreading`, `models/angsuran.go:200`), tak berkolom — penjaga repo melarang `TOTAL_` di migrasi (`docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 1e) |
| 39 | ✅ | K4: tabel warisan ditulis, tidak dibuat (`MODUL.md` *Tabel warisan*); `repository/kolom_test.go:243` TestSQLRiwayatProduksiMengikutiInsertViewSuggest |
| 40-41 | ✅ | `models/usulan_test.go:12` TestUsulanBelumTersimpanMenurutSaveViewSuggest, `:66` TestApprovalPerBarisHanyaAcceptReject; db `repository/polis_db_test.go:227` (K11) |
| 42 | ✅ | `APPROVAL` per baris terpisah dari `T_GENERAL_POLIS.IS_APPROVED` — `models/usulan_test.go:76` |
| 43-44 | ✅ | `AKSES_LOGIN` = identitas login, `PIC` = nama tampilan — `models/usulan_test.go:12`, `handlers/alur_test.go:155` |
| 45 | 🟡 | satu transaksi (`DalamTransaksi`), pembatalan diuji atas tiruan (`handlers/alur_test.go:355`); uji db kegagalan **tabel anak** belum ditulis (bab 9) — K11 |
| 46 | 🟡 | satu transaksi per tindakan, nol `COMMIT` di repository; `repository/polis_db_test.go:187` — **K11** |
| 47 | ✅ | penjaga `TestNolNamaTabelTelanjangDiQuery`; `inti/backend/penjaga/migrasi_test.go:54` TestSetiapPernyataanSahDanBerskema |
| 48 | ✅ | nol stored procedure (penomor menulis SQL); `repository/lama_test.go:16` TestSQLPemuatLamaBerskemaTanpaCommit |
| 49-51 | 🟡 | `repository/polis_db_test.go:64` TestPulangPergiHalamanLewatKatalog, `repository/lama_db_test.go:120` TestPemuatLamaNonProporsionalBersarang — **K11** |
| 52-54 | ✅ | `models/dokumenlama_test.go:93` (datar), `:178` (bersarang), `:211` TestUjiPemecahMencakupDuaBentuk |
| 55 | 🟡 | pemecah membawa nilai utuh (`models/dokumenlama_test.go:93`); kolom `592629512.88000028` di `repository/lama_db_test.go:74` — **K11** |
| 56 | ✅ | `repository/lama_test.go:43` TestPemuatTanpaJalurTulisTerpisah |
| 57 | ✅ | RALAT K17: penampung = berkas CSV — `models/laporanlama_test.go:24` TestLaporanMedanTakDikenalBerkasCSV |
| 58 | ✅ | `models/laporanlama_test.go:65`, `models/dokumenlama_test.go:256` |
| 59 | 🟡 | mekanisme: jumlah dicetak, kode keluar ≠ 0 (`models/laporanlama_test.go:107`); nol baris atas data nyata menunggu keputusan WO **F3** (22 pola medan tanpa kolom + `SuggestList` lama) |
| 60-61 | ✅ | riwayat cabang bersih (K1); inventaris disamarkan |
| 62 | ✅ | fixture `UJI-` |
| 63 | 📄 | migrasi dijalankan manusia (work owner) |
| 64 | ✅ | 79 medan diagram F10 = 69 kolom katalog + `NOPOLIS` + 9 tak berkolom bersebab (4 `LAYER*` F26, 4 `Total*` turunan, `isApprovedtoDeptHead` spec AC 64) + 7 kolom `json_polis` — `docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 1, 10 |
| 65 | ✅ | tanpa penjaga sinkron |
| 66 | ✅ | penjaga `TestHandlersTidakMengimporRepository`; `repository` tidak mengimpor `services` |

## 4 · Perubahan sejak putaran 1

### 4a · `spec.md`

| AC | Putaran 1 | Sekarang | Sebab |
| ---: | :---: | :---: | --- |
| 8 | ✅ | ✅ | bukti tetap; WO lama dikukuhkan **K2** (batas 200 juta tidak dibangun) |
| 12 | 🟡 | 🟡 | penahan: semula *"`M_NBTRIN_PERAN_TEMPAT` kosong"* → konstanta `models/peran_tempat.go` kosong (**K16**, migrasi 330 dihapus); 12 tempat terdaftar (P9) |
| 13, 81 | ✅ | ✅ | bukti *"tabel 330"* → konstanta kode (**K16**) |
| 15 | ✅ | ✅ | bunyi diperluas **RALAT K8**: JSON master baca-saja di satu fungsi (P5) |
| 26 | 🟡 | ✅ | **K3**: `Deduction1/2` uang (P2); `BROKERAGE`/`RNM_SHARE` nol pemakai (P4) |
| 45 | 🟡 | 🟡 | alasan `DateofSurvey`: semula *"layar survei tidak dibangun"* → **K7** (b), tidak ada tabel diagram |
| 52 | 🟡 | ✅ | RALAT P2 dari XML: atasan mengisi Approval + Suggest (+ ProductionDate bersyarat), bukan "tujuh medan" |
| 57 | ⛔ | ⛔ | alasan semula *"P29"* → **K8 butir 4** (b), SQL `BrowseTreatyOut` dikutip |
| 58 | ✅ | ✅ | diperluas: jalur XOL/NonProp dibangun (K8, P5) |
| 59 | ⛔ | ✅ | RALAT P4: `TreatyRealizationCheckDuplicate` dibangun dan diuji |
| 62 | ✅ | ✅ | bunyi baru RALAT K8: empat tetap, dua dimigrasi baca-saja |
| 66 | ⛔ | ⛔ | alasan (a) baru (P4): `isFOR` dikirim kosong, nol efek |
| 68 | ⛔ | 🟡 | pemuat dokumen lama dibangun (P7); uji db K11 |
| 69 | ⛔ | ✅ | pemuat dokumen lama (P7) |
| 71 | ✅ | ✅ | tempat simpan semula *"`T_POLIS_SUGGEST`"* → `HISTORYAKSEPTASIPRODUCTION` (**K4**, migrasi 328 dihapus, P1) |
| 84 | 🟡 | ✅ | **K6**: submit tanpa Approval 422 (P4) |
| 86 | ✅ | ✅ | bunyi baru **K14**: format dari setelan sel Section (P2) |
| 96 | ✅ | ✅ | peringatan *"commit 9456900e masih memuat nama"* gugur: riwayat ditulis ulang (**K1**) |

### 4b · `spec-penyimpanan-relasional.md`

| AC | Putaran 1 | Sekarang | Sebab |
| ---: | :---: | :---: | --- |
| 6 | 🟡 | 🟡 | penanda tutup semula *"`TGL_TUTUP IS NULL`"* → ada penerus via `OLD_POLIS_ID` (kolom dibuang, P1) |
| 21-22 | ⛔ | ✅ | pemuat dokumen lama (P7) |
| 27 | ✅ | ✅ | semula *"26 kolom"* → 10 medan diagram + 6 RALAT (P1) |
| 38 | 🟡 | ✅ | **K3** (P2, RALAT P9): `DEDUCTION1/2` uang; `TOTAL_*` turunan |
| 39 | ✅ | ✅ | semula *"tidak ditulis (syarat BusinessFac F)"* → ditulis (**K4**, P1) |
| 40-41 | ⛔ | ✅ | **K4**: catatan ditulis ke `HISTORYAKSEPTASIPRODUCTION` (P1; kotak tiket 19 RALAT P9) |
| 42-44 | ✅ | ✅ | semula *"`T_POLIS_SUGGEST.IS_APPROVED` / `OPERATOR_ID` / `OPERATOR_NAME`"* → `APPROVAL` / `AKSES_LOGIN` / `PIC` (**K4**) |
| 52-54, 56-58 | ⛔ | ✅ | pemuat dokumen lama (P7); AC 57 RALAT **K17** (CSV, bukan `T_POLIS_MEDAN_LAIN`) |
| 55 | ⛔ | 🟡 | pemuat dibangun (P7); kolom diperiksa uji db, K11 |
| 59 | 🔒 | 🟡 | mekanisme dibangun (P7); nol baris menunggu **F3** |
| 64 | 🟡 | ✅ | dicocokkan lawan diagram F10 (P1) |

### 4c · Bunyi putaran 1 yang gugur di berkas ini

Dikutip supaya pembaca berkas lama tahu apa yang tidak berlaku lagi: *"`M_NBTRIN_PERAN_TEMPAT` kosong"* (AC 12),
*"tabel 330"* (AC 13, 81), *"T_POLIS_SUGGEST + OperatorID"* (AC 71), *"T_POLIS_SUGGEST.IS_APPROVED terpisah"* (AC 42
penyimpanan), *"P29"* sebagai alasan jalur yang kini dibangun (AC 15, 57, 58; penyimpangan #3), *"pemuat tiket 22
belum dibangun"* (AC 68, 69; penyimpanan 21–22, 52–59), *"`T_POLIS_MEDAN_LAIN` (329)"*, dan penyimpangan #10
*"Daftar portal — semua kasus terbuka"* (kini bab 5 butir 11).

## 5 · Penyimpangan sadar

### 5a · Berdasar keputusan tertulis

| # | Hal | XML | Yang dibangun | Dasar |
| ---: | --- | --- | --- | --- |
| 1 | Sec Head menyetujui | ≤ 200 juta & bernomor → selesai; > batas → Dept Head (`CekLimitTreatyAcc_Act`, Decision13/8) | selalu naik ke Dept Head | **K2**, AC 8 |
| 2 | Tombol Generate/Submit atasan | identitas `<ID-operator-1>` (`DetailDeptHeadTreatyIn_UW`) | menurut posisi kasus (Dept Head menerbitkan nomor); enam tombol Submit terdaftar sebagai tempat (`models.DaftarTempat`) tetapi tidak dibaca layanan | P13, AC 12, tiket 05 |
| 3 | Sumber data kontrak | JSON master (`FetchMasterTreatyIn`, `BrowseTreatyIn`) | view `TREATYINDETAILJOINEDM`; jalur NonProp/XOL membaca JSON `M_TREATY_IN`/`_EDM` **baca-saja** di satu fungsi (`repository.MasterXOLDariJSON`) | P29 + **K8** (`[penyimpangan sadar]` atas P29) |
| 4 | Master jalur Proporsional (`TreatyIn.RNMShareP`, `BrokeragePercentP`, `SpreadingTotalPct` …) | JSON master | langkahnya dilewati (`models.MasterTersedia`) | (c) — di luar ukuran K8 (F1) |
| 5 | Hari tutup buku pra-proses | `>25` tertanam (`InputPolicyTreatyInPre_Act` 9) | hari dari `TANGGAL_CLOSING` (`models/polisaturan.go:137`; `handlers/logika_test.go:67`) | preseden WO PremiumList Life (`modul/premiumlistlife/backend/repository/ambangperiode_test.go:48` TestNolAmbangTutupBukuTertanam) |
| 6 | NBStatus | nama orang tertanam di connector | nama posisi tujuan / nama tampilan | **K5**, P40, AC 44 |
| 7 | OperatorName atasan | `pyUserIdentifier` | nama tampilan | P33, AC 42 |
| 8 | Riwayat akseptasi | tanpa skema, COMMIT per RDB | skema eksplisit, satu transaksi | AC 29, 30, 83 |
| 9 | Arasapas | stub kelas Data | bentuk kelas Work, 4 medan, sesudah commit | KEPUTUSAN-RONDE-12 butir 7 |
| 10 | `GetOldIDBusiness_SQL` | tanpa ORDER BY | `ORDER BY ID` (`repository/acuan.go:356`) | hasil tetap |
| 11 | Daftar portal *(RALAT #10 putaran 1: "semua kasus terbuka")* | satu grid `GetListOpportunity` di wadah `pyWorkGroup!='ReasLife' && pyWorkBasketList(2)=='ReasTreatyInAdmin'` lalu `!IsOperatorLife`; `UserIdentifier = InputParam.CARI38` (ditulis nol rule) | admin: semua kasus terbuka; Sec/Dept Head: hanya posisi yang dipegang; tanpa posisi tangga 403; cari hanya filter G (pengenal kasus), batas 500; klausa `pyWorkGroup` tidak dibangun (K12); filter C `.Name` tidak dibangun (ditulis nol rule, tak berkolom) | P8, P9; AC 11, 14 |
| 12 | Submit tanpa Approval | tombol Submit hanya untuk 1/0 | 422 di tiga jenjang | **K6**, AC 84 |
| 13 | Catatan usulan — syarat `Quotation.BusinessFac=="F"` | treaty tak pernah menulis | selalu ditulis ke `HISTORYAKSEPTASIPRODUCTION` | **K4** (grilling ID-31, AC 39) |
| 14 | `HISTORYAKSEPTASIPRODUCTION.DIV` | `OperatorID.pyOrgDivision` | NULL — tanpa sumber di `inti.Pelaku` | butir terbuka (PERMINTAAN B4) |
| 15 | Pemecah dokumen lama | ID-4: seam repository | fungsi murni `models.PecahDokumenLama`; penulisan lewat antarmuka repository yang sama (ID-3) | RALAT ID-4 (P7) |
| 16 | Medan tak dikenal dokumen lama | tabel penampung (ID-27) | berkas CSV `POLIS_ID,JALUR,NILAI` | **K17** |
| 17 | Cap waktu ` GMT` dokumen/master lama | teks | jam dinding Asia/Jakarta; tanggal ambigu tidak ditebak → laporan galat | **K15**; `GeneratePolicyNoTreaty_Act` 5.3 |
| 18 | Tampilan tanggal | beberapa format (`dd/MM/yyyy`, `Date-Short-Custom-YYYY`, …) | satu format (`formatDate` inti) | AC 33 (P2) |
| 19 | NonProp — halaman master | disimpan di blob kasus | dibaca ulang saat dibuka; tampilan dari penanda terkini (`models.TampilanMasterNonProp`) | K8 (spec RALAT K8 butir 1) |
| 20 | NonProp — master tidak ada saat pilih bisnis | `catch … oLog.error`, halaman kosong | 422, nol simpanan | K8, AC 36–38 |
| 21 | NonProp — `TreatyRealizationCheckXOLList` | berjalan bila `IsNewPolicyNonProp` 1 | dilewati untuk `ProportionalType = 'Proportional'` | spec-penyimpanan AC 33 |
| 22 | Breakdown spreading | `BreakDownSpreading_Act` | tidak dibangun | **K9** |
| 23 | Laporan survei historis | popup + `SetSurveyReport_Act` | tidak dibangun | **K7** |
| 24 | `GENERATE_SEQUENCE_NUMBER` | ditulis procedure `PROC_GENERATE_SEQUENCE_NUMBER`; diagram F103/F118 menyebutnya **"dibaca saja"** | ditulis SQL oleh penomor bersama `inti/backend/penomor` (padanan procedure — spec-penyimpanan AC 48 nol procedure) | ⚠️ dicatat untuk ditinjau (PERMINTAAN F8) |

### 5b · `[menunggu konfirmasi WO]` — sudah dibangun dengan tafsiran tertulis

| # | Hal | Tafsiran yang dibangun | Rujukan |
| ---: | --- | --- | --- |
| W1 | **K4 (1)** catatan usulan per jenjang | XML menulis hanya pasca-submit admin; di sini setiap submit (admin, Sec Head, Dept Head), transaksi yang sama | `models/usulan.go` butir 2; tiket 10 bab P9; PERMINTAAN F2 |
| W2 | **K4 (2)** `TGL_INP` | XML `hh` (jam 12) lalu `HH24`; di sini jam 24 apa adanya | `models/usulan.go` butir 3; F2 |
| W3 | **K4 (3)** `NOURUT` | XML `.pxListSubscript`; di sini MAX+1 per IDPEGA di bawah kunci kasus — terbukti sama (`handlers/alur_test.go:273`) | `models/usulan.go` butir 4; F2 |
| W4 | Tafsiran **K8** | "medan master" = setiap medan yang **dibaca rule terjangkau jalur NB NonProp** (berkutip langkah di `models/masterxol.go`), melampaui daftar harfiah K8 | tiket 01 bab P9; PERMINTAAN F1 |
| W5 | **P3** pemilih SOB disimpan saat klik | Pega memegang hasil `SearchHierarkiSourceBizAgent_PostDT` di clipboard sampai Save/Submit; di sini `POST /kasus/{id}/pilih-sumber-bisnis` menyimpan dalam satu transaksi (layar tidak boleh mengirim `Quotation.*`) | `services/sumberbisnis.go:41`; PERMINTAAN F4 |
| W6 | **P5** keanehan Pega ditiru apa adanya | `InsertToTreatyXOLList`: `DueTo` selalu kosong, `Currency` induk dari layer terakhir, potongan kumulatif (`models/nonprop.go:133`); `InsertToTreatyXOLListRetroShare`: CARI terbawa antar layer, 2.3.7.4 membaca net/potongan layer (`:220`); `TreatyNonPropSetSpreading` `PremiumSpreaded = local.netpremi` (`:340`); `InputDetailNonProp` 16.2.2/17/19 (`models/nonprop_detail.go:111`); preACT 18 (`:242`) | PERMINTAAN F5 |

## 6 · Rule yang tetap tidak dibangun

Sumber: `docs/alat/status.json` (176 rule terjangkau: **109 dibangun**, termasuk yang sebagian, dan **67 tidak
dibangun**). Lima belas baris status.json menyebut dasarnya tanpa huruf (a)/(b)/(c); hurufnya ditetapkan di sini
dari teks baris itu (ditandai \*).

### (a) Tidak terjangkau dari titik masuk nyata, atau efeknya dibaca nol rule — 29

| Kelompok | Rule |
| --- | --- |
| Cek penawaran ganda lama (langkah `//`, pemanggil tanpa parameter) | `Activity/CheckDuplicateOffer`, `RDBList/GetCountClaim`\*, `ReportDefinition/BrowseTREATY_IN`\* |
| Revisi master (`revisionstate==1`, tak pernah dari NB) / nol efek | `RDBList/GetCurrentDate`\*, `Activity/TreatyInInputVis`\* (ditimpa `SetTreatyIn_Act` 3), `Activity/ConvertHistoryDate` (komentar master dibaca nol rule) |
| Varian EDM jalur NonProp (syarat mustahil sesudah NonProp langkah 10) | `Activity/InputPolicyTreatyEDMDetail_NP`, `Activity/InsertToTreatyOutXOLList`, `Section/DetailPolicyTreatyInNonProportionalEDM`, `Section/Installments_ReadOnly`, `FlowAction/Installments_ReadOnly` |
| Tombol bertampil `1=2` | `Activity/TreatyInNonSetTotal` |
| Lampiran (efek hanya ke `Protection_Act`, yang tidak ada di korpus; BusinessFac "T") | `Activity/SetCategoryAttach`, `Activity/InputParamUploadReas_act`, `DataTransform/setCategoryAttachment_DT`, `RDBList/AttachmentLife`, `RDBList/CategoryAttach_SQL` |
| Pemilih SOB — cabang yang ditulis/dibaca nol rule | `Activity/SearchHierarkiSourceBizAgentTreatyIn_Act`, `DataTransform/btnCedingCO_DT`, `RDBList/BrowseClientEmail_SQL`, `ReportDefinition/BrowseAgentNusaRe_RD`, `ReportDefinition/BrowseCedingCo_RD` |
| Hasil dibuang | `RDBList/GenerateNoPolicy` |
| When yang tak pernah benar / tak mengubah perilaku | `When/IsClaim` (`CLM-`), `When/isClaimTreaty` (Decision5 tanpa connector masuk), `When/isSellingModeB2B`, `When/isSellingModeB2BB2C`, `When/isSellingModeB2C`, `When/pyIsIpadOrDesktop` |

### (b) Keputusan work owner tertulis — 36

| Keputusan | Rule |
| --- | --- |
| **K2** — Sec Head selalu ke Dept Head | `Activity/CekLimitTreatyAcc_Act`, `When/ToTREATYDEPTHEAD` |
| **K7** + bab 0 butir 11 — survei historis, tidak ada tabel diagram | `Activity/SetSurveyReport_Act`, `Activity/ConcatSlipOfferNo_Act`, `FlowAction/InputHistoricalSurveyReport`, `FlowAction/InputHistoricalSurveyReportUW`, `Harness/HistoricalSurveyReport`, `Harness/HistoricalSurveyReportUW`, `Section/HistoricalSurveyReportDtl`, `Section/HistoricalSurveyReportDtlUW`, `Section/InputHistoricalSurveyReportDtl`, `Section/InputHistoricalSurveyReportDtlUW` |
| **K8 butir 4** — treaty keluar, JSON `M_TREATY_OUT` | `Activity/InputPolicyTreatyOutDetail_NonProp`\*, `Activity/InputPolicyTreatyOutDetail_preACT`\*, `Activity/SetValueRetro_Act`\*, `Activity/TreatyNonPropOutSetSpreading`\*, `Harness/BusinessAndSOBListRetro`\*, `Section/BusinessAndSOBListRetro`\*, `Section/DetailPolicyTreatyOutNonProportional`\*, `RDBList/BrowseTreatyOut`\*, `RDBList/BrowseTreatyOutDetail`\*, `ReportDefinition/BrowseTreatyOutDetail`\* |
| **K9** + KEPUTUSAN-RONDE-12 butir 3/3b — breakdown spreading | `Activity/BreakDownSpreading_Act`, `RDBList/GetBreakDownSpread_SQL` |
| **AC 16** (nol tulis JSON) + AC 48 (nol procedure) + diagram F11 | `Activity/SaveJsonPolisTreatyIn_Act`, `RDBList/SavePolisTreatyIn_SQL`, `RDBList/SaveTreatyIn`\* (juga (a): hanya `revisionstate==1`) |
| **P44 / AC 53** — elemen mati | `Section/SFAPortal_OpportunitiesList_Header`, `ReportDefinition/crmOpportunitiesList` (juga (a)) |
| **K5 / K12 / AC 12–14, 81** — pemeriksaan identitas/atribut operator tanpa padanan | `When/IsSPVCreate`, `When/IsSPVTreaty1`, `When/IsTreaty1`, `When/IsUW`, `When/IsNotAdmin`, `When/IsOperatorLife`, `When/crmCreateOpportunity` |

### (c) Data hanya di JSON, di luar pengecualian K8 — 2

| Rule | Alasan |
| --- | --- |
| `Activity/FetchMasterTreatyIn` | master jalur Proporsional (`TreatyInputPctCommSpreading` langkah 1); `RiCommOgp` dibaca dari kolom view, medan master lain hanya JSON — di luar ukuran K8 (F1) |
| `RDBList/BrowseTreatyInDetailJoinEDM` | JSON `M_TREATY_IN_DETAIL_EDM` — tabel di luar dua tabel pengecualian K8 |

Bagian rule yang **dibangun sebagian** (langkah lain tidak, dengan alasan per langkah di status.json):
`AgentSourceBizTreatyIn_Act`, `CalculatePremi_Act` (c), `InputPolicyTreatyInDetail_preACT` (langkah 9–10, 13 (c)),
`InputQuotation_PreAct`, `SetTreatyIn_Act`, `TreatyInputPctCommSpreading` (c), `SearchHierarkiSourceBizAgent_PostDT`,
`btnSOB_DT`, `ReportDefinition/GetListOpportunity`, `Section/DetailPolicyTreatyIn`, `Section/SFAPortal_OpportunitiesList`,
`Section/SourceHierarki`. `Protection_Act` varian treaty **tidak ada di korpus** (`docs/INVENTARIS-XML.md` bab 12) —
tidak dapat dibangun.

## 7 · Jalur Proporsional dan NonProporsional — dari layar sampai tersimpan

Satu halaman portal `nbtreatyin-portal` (`frontend/pages/PortalNBTreatyIn.tsx`) dan satu layar kasus
(`frontend/pages/LayarKasus.tsx`); rute `/api/nb-treaty-in` (`handlers/rute.go:45`–`:57`). Padanan setiap tombol →
activity: `docs/INVENTARIS-XML.md` bab 13 (`docs/alat/tombol.json`).

| Tahap | Proporsional | NonProporsional / XOL |
| --- | --- | --- |
| Buat kasus | `POST /kasus` → `services.BuatKasus` (`services/layanan.go:144`): `T_WORK_POLIS` (`NB-<SEQ_WORK_POLIS>`) + `T_GENERAL_POLIS`, `BusinessFac` "T", posisi `ReasTreatyInAdmin` | sama |
| Buka layar | `GET /kasus/{id}` → `BukaKasus` (`:193`): pra-proses admin (`PraprosesAdmin`, hari tutup buku dari `TANGGAL_CLOSING`), atasan (`PraprosesAtasan`); `LAYER*` dibaca balik dari view lewat `TREATY_IN_ID` | + `IsNewPolicyNonProp` 1 → `TreatyRealizationCheckXOLList` (`services/nonprop.go:118`): master dibaca ulang (`MasterXOLDariJSON`), `InsertToTreatyXOLList(RetroShare)`; master hilang → pesan `Error fetching XolList` |
| Pilih bisnis | `POST /pilih-bisnis` → `PilihBisnis` (`services/tindakan.go:274`): view `TREATYINDETAILJOINEDM` (`DetailKontrak`), `BUSINESS` (`GetOldIDBusiness_SQL`), komisi `RIONR` (`KomisiKontrak`, preACT 17), preACT 1–6, 8, 11, 14, 15, 19; gagal baca → 422 | + preACT 16/18 → `pilihBisnisNonProp` (`services/nonprop.go:89`): `InputPolicyTreatyInDetail_NonProp`, `InsertToTreatyXOLList(RetroShare)`, `TreatyNonPropSetSpreading`, `TreatySetReinstatement`, `SetReinstatementPct`, PPN/PPH layer XOL; master tidak ada → 422 |
| Pemilih SOB (ClaimType `XOL Retro`) | `GET /sumber-bisnis` (`AGENT`), `POST /pilih-sumber-bisnis` → `Quotation.SourceOfBusiness` | sama |
| Hitung (refresh sel) | `POST /hitung` urutan action set XML (`Count*_Act` → `CountOGPONP_Act`, `CalculatePremi`, `SetPPNPPH`, `FillPaymentInstallment`, `CountSpreading`, `SystemSetOneYear` …) — tidak menyimpan | `CountSpreading` grid `SpreadingRiskList` NonProp; rantai uang Proporsional tidak terpicu (wadah tersembunyi) |
| Simpan draf | `PUT /kasus/{id}` (admin) → `SimpanDraf` (`services/tindakan.go:247`): medan wajib (`MedanWajibBerlaku`), `PeriksaBentukSimpan`, `SimpanHalaman` | sama; `TreatyXOLList` tidak diterima dari layar |
| Tersimpan di | `T_GENERAL_POLIS`, `T_POLIS_QUOTATION`, `T_POLIS_CEDING`, `T_POLIS_INSTALMENT`, `T_POLIS_SPREADING` (XOL dan rincian angsuran **ditolak**, AC 31–33) | + `T_POLIS_INSTALMENT_DETAIL`, `T_POLIS_XOL`, `T_POLIS_XOL_LAYER`; halaman master **tidak** disimpan |
| Submit admin | `POST /kirim` → `Kirim` (`services/tindakan.go:359`): DT `isApproved`; Approval wajib (K6); `TreatyRealizationCheckDuplicate` (`TREATYINPRODUCTION`, IsApproved 1, ClaimType bukan XOL); riwayat `HISTORYAKSEPTASIPEGA`; catatan `HISTORYAKSEPTASIPRODUCTION` (K4); NBStatus — satu transaksi | sama |
| Sec Head / Dept Head | Sec Head setuju → Dept Head (K2); tolak → admin. Dept Head: `POST /nomor-polis` → `TerbitkanNomor` (`services/tindakan.go:537`; `KODE_PRODUKSI` NONLIFE, `GENERATE_SEQUENCE_NUMBER`, sekali per berkas) → submit → `Resolved-Completed` → konversi Arasapas (KEPUTUSAN-RONDE-12 butir 7) | sama |
| Uji ujung ke ujung | `handlers/alur_test.go:202`, `:435`; `handlers/logika_test.go`; db `repository/polis_db_test.go:64` (K11) | `handlers/nonprop_test.go:82` (pilih bisnis → hitung → simpan → buka ulang → tiga jenjang → nomor → selesai), `:207`, `:226` |
| Belum / tidak | master Proporsional dari JSON (c, F1); `BreakDownSpreading` (K9); survei (K7); `CekLimitTreatyAcc_Act` (K2); `Protection_Act` (tak ada di korpus); 4 tempat `ProductionDate` (K12) | treaty keluar / XOL Retro `ReinsuranceListTONP` (K8 butir 4); varian EDM NonProp (a); PPN/PPh rincian angsuran tak berkolom (dihitung, tidak disimpan) |

## 8 · Tabel: delapan `CREATE TABLE` dan tabel warisan

Tepat **delapan** `CREATE TABLE` di migrasi modul (`backend/migrations/320`–`327`), sama dengan diagram grilling
`Diagram-Skema-Tabel-NusantaraRe.xlsx` (PROMPT putaran 2 bab 0 butir 11), ditagih `repository/kolom_test.go:310`
TestTabelDanKolomMengikutiDiagramGrilling (nama, jumlah, dan kolom persis). Migrasi 328 `T_POLIS_SUGGEST`, 329
`T_POLIS_MEDAN_LAIN`, 330 `M_NBTRIN_PERAN_TEMPAT` **dihapus** (K4, K17, K16). Perbandingan kolom demi kolom
(dipertahankan / dipindah / dibuang + bukti): `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

| # | Migrasi | Tabel | Kolom | Sheet `NB Treaty In Prop` | Sheet `NB Treaty In NonProp` | Relasi | PERBANDINGAN |
| ---: | --- | --- | ---: | --- | --- | --- | --- |
| 1 | `320_t_general_polis.sql` | `T_GENERAL_POLIS` | 80 | F9–F33, ringkasan K108 | F9–F33, ringkasan K123 | 1:1 shared PK `T_WORK_POLIS`; `UNIQUE (OLD_POLIS_ID)`, `UNIQUE (NOPOLIS, PRODKE)` (indeks fungsi `CASE`) | bab 1 |
| 2 | `321_t_polis_quotation.sql` | `T_POLIS_QUOTATION` | 17 | G34, J35–J38, K109 | G34, J35–J38, K124 | 1:1 `POLIS_ID` | bab 2 |
| 3 | `322_t_polis_ceding.sql` | `T_POLIS_CEDING` | 5 | O39, R40–R50, K110 | O39, R40–R50, K125 | 1:N `QUOTATION_ID` | bab 3 |
| 4 | `323_t_polis_instalment.sql` | `T_POLIS_INSTALMENT` | 17 | G51, J52–J56, K111 | G51, J52–J55, K126 | 1:N `POLIS_ID` | bab 4 |
| 5 | `324_t_polis_instalment_detail.sql` | `T_POLIS_INSTALMENT_DETAIL` | 14 | — (J56: proporsional berhenti di induk) | O56, R57–R61, K127 | 1:N `INSTALMENT_ID` · NonProp saja | bab 5 |
| 6 | `325_t_polis_spreading.sql` | `T_POLIS_SPREADING` | 12 | G65, J66–J69, K113 | G70, J71–J74, K129 | 1:N `POLIS_ID` | bab 6 |
| 7 | `326_t_polis_xol.sql` | `T_POLIS_XOL` | 16 | J71–J72 (nol baris) | G75, J76–J80, K130 | 1:N `POLIS_ID` · NonProp saja | bab 7 |
| 8 | `327_t_polis_xol_layer.sql` | `T_POLIS_XOL_LAYER` | 20 | J71–J72 (nol baris) | O81, R82–R87, K131 | 1:N `XOL_ID` · NonProp saja | bab 8 |

*Kolom* = definisi kolom di `CREATE TABLE` (termasuk kunci dan `NOURUT`).

Tabel lama yang disentuh — **tidak dibuat, tidak diubah strukturnya**:

| Tabel | Akses | Dasar diagram / keputusan |
| --- | --- | --- |
| `T_WORK_POLIS` (+ `SEQ_WORK_POLIS`) | tulis/baca | akar B5, milik premiumlistlife (migrasi 050/057/059) |
| `HISTORYAKSEPTASIPEGA` | tulis/baca | Prop F98–F99 / NonProp F113: `InsertHistoryAkseptasiPega_Sql` |
| `HISTORYAKSEPTASIPRODUCTION` | tulis/baca | Prop J74–J76 / NonProp J89–J92: `SaveViewSuggest → InsertViewSuggest_SQL` (K4) |
| `GENERATE_SEQUENCE_NUMBER` | tulis lewat `inti/backend/penomor` | diagram F103/F118 "dibaca saja" — lihat bab 5 butir 24 |
| `TANGGAL_CLOSING`, `KODE_PRODUKSI`, `CURRENCY`, `BUSINESS`, `REINSURANCETYPE`, `TREATYGROUP`, `MARKETINGOFFICER`, `CLIENT`, `AGENT`, `M_LOGIN_GO`, view `TREATYINDETAILJOINEDM`, `TREATYINDETAIL`, `TREATYINPRODUCTION` | baca | daftar "dibaca saja" F103/F118 dan RD/RDB terjangkau; `TREATYINPRODUCTION` ditulis modul EDM |
| `M_TREATY_IN`, `M_TREATY_IN_EDM` | baca (JSON, satu fungsi) | K8 |
| `JSON_POLIS` | baca (pemuat dokumen lama saja) | tiket 22, KEPUTUSAN-RONDE-12 butir 5 |

## 9 · Status tiket 00–23 dan pekerjaan yang masih dapat dikerjakan

Aturan penandaan (PROMPT putaran 2 bab 0 butir 1): **selesai** bila setiap AC tiketnya ✅, 📄, ⛔ dengan alasan
(a)/(b)/(c) terbukti, atau tertahan **hanya** oleh pihak luar yang terbukti (K11 skema uji Oracle, K12 IAM, keputusan
WO yang tercatat di PERMINTAAN F, data DBA, tim inti). Tiket yang masih punya pekerjaan yang dapat dikerjakan tetap
**sebagian**.

| Tiket | Status | AC belum ✅ dan penahannya |
| --- | --- | --- |
| 00 | selesai | 0 AC; butir DBA `[terbuka]` (PERMINTAAN C3, C7) |
| 01 | **sebagian** | 17, 89 — uji db pembacaan view belum ditulis; 57 ⛔ (b) K8 butir 4; F1 |
| 02 | selesai | 31 — K11 |
| 03 | selesai | — |
| 04 | selesai | — (klausa `pyWorkGroup` K12, PERMINTAAN C6) |
| 05 | needs-info | 12 — K12 |
| 06 | selesai | 66 ⛔ (a) |
| 07 | selesai | 23 — K11 |
| 08 | **sebagian** | 83 — uji db kegagalan riwayat belum ditulis; 29 — K11 |
| 09 | selesai | — |
| 10 | **sebagian** | 72 — uji db urutan riwayat belum ditulis; F2; DIV (B4); tipe kolom fisik (C8) |
| 11 | selesai | 45 — K12 (`ProductionDate`), K7 (b) (`DateofSurvey`); spec §9.2 butir 17 (Product & Underwriting) |
| 12 | selesai | — |
| 13 | selesai | — (sambungan Arasapas C1) |
| 14 | selesai | — |
| 15 | selesai | 68 — K11, F7; 70, 93–95 📄 |
| 16 | selesai | 1, 6 — K11 |
| 17 | **sebagian** | 9 — uji db tidak membaca `NOURUT` |
| 18 | **sebagian** | 25 — bukti/RALAT; 15 — bunyi baru RALAT; 12 — uji db `GROUP_PANEL`; 13 — K11 |
| 19 | selesai | — (F2, B4, C8 untuk tabel warisan) |
| 20 | **sebagian** | 45 — uji db kegagalan tabel anak belum ditulis; 46 — K11 |
| 21 | selesai | 49–51 — K11 |
| 22 | selesai | 55 — K11; 59 — F3; F6, F7, C5, C7 |
| 23 | selesai | 63 📄 |

**Pekerjaan yang masih dapat dikerjakan** (tidak tertahan pihak luar):

1. **Uji bertag `db`** (ditulis sekarang, dijalankan begitu K11 ada; pola lewati-bila-tabel-tak-ada seperti
   `repository/polis_db_test.go:227`): (a) pembacaan view `TREATYINDETAILJOINEDM` — 33 kolom `KolomRDDetail` terisi
   dan tak satu pun hilang (spec AC 17, 89; tiket 01); (b) `HISTORYAKSEPTASIPEGA` — `DaftarRiwayat` berurut
   `TGL_TRANSFER` (AC 72; tiket 10) dan kegagalan `CatatRiwayat` membatalkan seluruh submit lawan Oracle (AC 83;
   tiket 08); (c) tiga baris angsuran, hapus yang kedua, kolom `NOURUT` dibaca langsung = 1, 2 (penyimpanan AC 9;
   tiket 17); (d) `GROUP_PANEL "006"` dibaca langsung dari kolom (penyimpanan AC 12; tiket 18); (e) kegagalan
   menulis salah satu **tabel anak** membatalkan induk (penyimpanan AC 45; tiket 20).
2. **Penyimpanan AC 25** (tiket 18): sisir setiap pembandingan nilai di port (`models/hitung.go`, `models/angsuran.go`,
   `models/nonprop*.go`) lawan XML; bila tak satu pun membandingkan dua nilai uang (yang ada: persen lawan 100,
   tanda lawan 0), tulis buktinya / RALAT AC 25 (ID-20); bila ada, bangun bentuk terbulatkan.
3. **Penyimpanan AC 15** (tiket 18): tulis bunyi baru RALAT (Oracle `''` ≡ NULL; dibaca kembali `""`, tidak pernah
   `"0"`) — tabel RALAT bab 6 hanya mencatat temuan, belum bunyi baru.
4. **`docs/alat/status.json`**: lima belas baris "tidak dibangun" tanpa huruf alasan (\* di bab 6) — tambahkan
   (a)/(b) lalu bangkitkan ulang `INVENTARIS-XML.md`.

## 10 · Perintah verifikasi

| Perintah | Hasil (P10, worktree `wt-nbtr-p10-hasil`) |
| --- | --- |
| `go test ./modul/nbtreatyin/...` | lulus (handlers, models, repository, tiruan) |
| `go test ./inti/...` | lulus kecuali `inti/backend/penjaga` `TestNolAlamatLayananDiKode` — butuh berkas `.env` lokal yang tidak ada di worktree baru (lingkungan; bukan perubahan modul) |
| `make test-db` | **tidak dijalankan** — K11 kosong |
| `go vet ./...`, `go test ./...`, `npm run typecheck`, `npm test` | dirangkum orkestrator di laporan akhir putaran 2 (bab 10 PROMPT), dibandingkan baseline bab 1 |
