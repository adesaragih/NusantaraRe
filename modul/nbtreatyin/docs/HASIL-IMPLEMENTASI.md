# Hasil implementasi — modul `nbtreatyin` (NB Treaty In)

> **Putaran 3** · konsolidasi paket putaran 3 (nusare-03, R1–R6) di cabang integrasi
> `modul/nbtreatyin/implementasi` (`724ef8eb`), dirangkum paket **P3K** pada 04-10-2026 (cabang
> `modul/nbtreatyin/p3r-hasil`). Keputusan work owner putaran 3: `PROMPT-NB-TREATY-IN-PUTARAN-3.md` bab 2 (F1–F8,
> U1, U2, K11–K13, K18); keputusan putaran 2 (K1–K10, K14–K17) tetap berlaku. Setiap bukti di bawah diperiksa ulang
> ke kode dan uji cabang itu — nomor baris dibangkitkan dari berkasnya, bukan disalin dari laporan paket.
>
> Status: ✅ lulus (ada bukti uji/kode) · 🟡 sebagian (penahannya disebut) · ⛔ tidak dibangun, dengan alasan
> (a) tidak terjangkau / nol efek, (b) keputusan work owner tertulis, atau (c) data hanya di JSON di luar K8 ·
> 📄 butir dokumen · 🔒 `[terbuka]`, menunggu pihak lain.
>
> "Bukti" menunjuk uji (`berkas:baris NamaUji`) atau kode (`berkas:baris`). Jalur relatif terhadap
> `modul/nbtreatyin/backend/` kecuali diawali `frontend/`, `docs/`, `inti/`, atau `modul/`.
>
> ⛔ Uji bertag `db` (Oracle sungguhan) **belum pernah dijalankan**: skema uji K11 kosong (PERMINTAAN C9). AC yang
> buktinya hanya uji db ditandai 🟡 dengan penahan **K11**. **Setiap** AC 🟡-K11 punya uji db yang memeriksa yang
> dituntut AC-nya (putaran 3 bab 3.3; ditulis dan dikompilasi `go vet -tags db`, belum dijalankan).

## Isi

1. Ringkasan
2. `spec.md` — 96 AC
3. `spec-penyimpanan-relasional.md` — 66 AC
4. Perubahan sejak putaran 2
5. Penyimpangan sadar
6. Rule yang tetap tidak dibangun
7. Jalur Proporsional dan NonProporsional — dari layar sampai tersimpan
8. Tabel: delapan `CREATE TABLE` dan tabel warisan
9. Status tiket 00–23
10. Perintah verifikasi

## 1 · Ringkasan

Dihitung dari tabel bab 2–3 (rentang dijabarkan; setiap nomor tepat satu kali).

| | ✅ | 🟡 | ⛔ | 🔒 | 📄 | jumlah |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `spec.md` | **80** | 10 | 2 | 0 | 4 | **96** |
| `spec-penyimpanan-relasional.md` | **53** | 12 | 0 | 0 | 1 | **66** |
| *putaran 2 — `spec.md`* | *80* | *10* | *2* | *0* | *4* | *96* |
| *putaran 2 — penyimpanan* | *53* | *12* | *0* | *0* | *1* | *66* |
| *putaran 1 — `spec.md`* | *75* | *12* | *5* | *0* | *4* | *96* |
| *putaran 1 — penyimpanan* | *39* | *14* | *11* | *1* | *1* | *66* |

Hitungan status **tidak berubah** sejak putaran 2: pekerjaan putaran 3 menutup keputusan F1–F8 dan temuan audit
silang di dalam AC yang sudah ✅ (bukti diperkuat, bunyi di-RALAT), dan memperkuat uji db AC 🟡-K11 — penahan
K11/K12 di luar jangkauan agen tetap. Rinciannya bab 4. Cacah penyimpanan = nomor AC berurutan **1–66**; sisipan
**`20b`** dicacah bersama nomor 20 (`docs/issues/00-PETA-AC-PENYIMPANAN.md`).

Penahan setiap AC yang belum ✅:

| Penahan | `spec.md` | penyimpanan |
| --- | --- | --- |
| **K11** — skema uji Oracle (PERMINTAAN C9; `make test-db` tanpa modul ini, A4): uji db ditulis, belum dijalankan | 🟡 17, 23, 29, 31, 72, 83, 89 | 🟡 1, 6, 9, 12, 13, 45, 46, 49, 50, 51, 55 |
| **K12** — pemetaan peran IAM | 🟡 12, 45 (bersama K7 (b) untuk `DateofSurvey`) | — |
| **F7** — urutan resmi pemuatan dokumen lama (WO/DBA; langkah 2 = uji-kering di skema uji, K11) | 🟡 68 | 🟡 59 |
| **keputusan WO tertulis** — tetap ⛔ | ⛔ 57 (K8 butir 4: treaty keluar JSON `M_TREATY_OUT`), ⛔ 66 ((a) `isFOR` nol efek, bukti XML) | — |
| **K13** — migrasi 968 tercatat di `POOLDATA` oleh pihak lain | — (penahan migrasi, bukan AC; bab 8) | — |

Rule terjangkau (176, `docs/alat/status.json`): **103 dibangun** · 73 tidak dibangun = **(a) 36**, **(b) 35**,
**(c) 2** (bab 6).

## 2 · `spec.md` — 96 acceptance criteria

| AC | Status | Bukti / alasan |
| ---: | :---: | --- |
| 1 | ✅ | `models.Disetujui` (`models/tangga.go:127`); `models/tangga_test.go:12` TestIsApprovedNolDitolak; `handlers/alur_test.go:155` TestAdminMenolakDiselesaikanDitolak |
| 2 | ✅ | `models/tangga_test.go:18` TestIsApprovedSelainNolDisetujuiTermasukKosong |
| 3-4 | ✅ | `models/tangga_test.go:26` TestIsApprovedDibandingkanSebagaiTeks — teks, menurut DecisionTable `isApproved` (bukan When bernama sama) |
| 5 | ✅ | `models/tangga_test.go:41` TestTanggaAdminMenolakSelesaiDitolak; `handlers/alur_test.go:155` TestAdminMenolakDiselesaikanDitolak (Resolved-Rejected; kirim ulang 409) |
| 6 | ✅ | `models/tangga_test.go:51` TestTanggaAtasanMenolakKembaliKeAdmin; `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali |
| 7 | ✅ | `models/tangga_test.go:63` TestTanggaAdminMenyetujuiNaikKeSecHead; `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali |
| 8 | ✅ | `models/tangga_test.go:73` TestTanggaSecHeadMenyetujuiSelaluNaikKeDeptHead — **K2**: batas 200 juta `CekLimitTreatyAcc_Act` tidak dibangun (b), pertentangannya di tiket 03 |
| 9 | ✅ | `models/tangga_test.go:84` TestTanggaDeptHeadMenyetujuiSelesai; `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali (alur penuh sampai selesai) |
| 10 | ✅ | `models/tangga_test.go:95` TestPosisiBuanganTidakAda; `When/isClaimTreaty` (Decision5 tanpa connector masuk) tidak dibangun (a) |
| 11 | ✅ | antrean = workbasket (`anggota`, `services/layanan.go:80`); `handlers/alur_test.go:139` TestBukanAnggotaAntreanDitolak; `handlers/portal_test.go:49` TestGerbangDaftarPortal; grid portal = 6 kolom `GetListOpportunity` (kolom buatan Position/No Polis dibuang R6 W6) — `handlers/portal_test.go:116` TestDaftarPortalSesuaiGetListOpportunity |
| 12 | 🟡 | wewenang = keanggotaan antrean (`inti.Pelaku.Peran`), nol pemeriksaan identitas di kode; kedua belas tempat identitas XML terdaftar di `models.DaftarTempat` (`models/peran_tempat.go:90`), pemetaannya `PemetaanPeranTempat` **kosong** (`models/peran_tempat.go:123`, K16) ⇒ tertunda — penahan **K12 (IAM)** |
| 13 | ✅ | peran dari `inti.Pelaku.Peran`; pemetaan tempat = konstanta `models/peran_tempat.go` (K16); nol kolom telepon: `When/IsSPVTreaty1`, `When/IsTreaty1` (`OperatorID.pyTelephone`) tidak dibangun (b) K5 |
| 14 | ✅ | `anggota` = `PunyaPeran(nama)`; `handlers/alur_test.go:139` TestBukanAnggotaAntreanDitolak (urutan peran terbalik tetap anggota); `When/IsUW` (nomor urut `pyWorkBasketList(1)`) tidak dibangun (b) |
| 15 | ✅ | RALAT K8: kontrak dari view (`repository.DetailKontrak`, `repository/acuan.go:141`; popup `repository.DaftarBisnis` RD `BrowseTreatyJoinEDM`, `repository/acuan.go:231` — R5); JSON master **hanya** `repository.MasterXOLDariJSON` (`repository/masterxol.go:74`) dengan daftar medan **tertutup** (F1): `models/masterxol_test.go:14` TestDaftarMedanMasterXOLTertutup, `repository/masterxol_test.go:162` TestUraiMasterXOLHanyaMedanDaftarTertutup, `repository/masterxol_penjaga_test.go:199` TestMasterXOLDariJSONHanyaMengembalikanHasilUrai, `repository/masterxol_penjaga_test.go:108` TestKolomDokumenHanyaDiMasterXOLDariJSON |
| 16 | ✅ | nol penulisan JSON: penulisan hanya ke 8 tabel diagram + `T_WORK_POLIS` + `HISTORYAKSEPTASIPEGA` + `HISTORYAKSEPTASIPRODUCTION` (+ `GENERATE_SEQUENCE_NUMBER` lewat penomor inti, **F8 disetujui WO**); `handlers/nonprop_test.go:82` TestNonPropPilihBisnisHitungSimpanBacaKembali (halaman master tidak tersimpan); `SaveJsonPolisTreatyIn_Act`, `SavePolisTreatyIn_SQL`, `SaveTreatyIn` tidak dibangun (b) |
| 17 | 🟡 | 33 kolom RD dibaca dari view (`KolomRDDetail`, `repository/acuan.go:34`; ke-33 `pyUIFields` `BrowseTreatyJoinEDM` = `BrowseTreatyInDetail`, spec §5.1 RALAT W1); `repository/kontrak_db_test.go:135` TestDetailKontrakMengisiKe33MedanRD, `repository/kontrak_db_test.go:229` TestDaftarBisnisRDBrowseTreatyJoinEDM — **K11** |
| 18 | ✅ | medan uang dibaca apa adanya (`TO_CHAR` TM9, `repository/kolom.go:147`), tanpa hitung ulang saat dibuka; `handlers/nonprop_test.go:82` TestNonPropPilihBisnisHitungSimpanBacaKembali; kiriman layar tidak menimpa uang master NonProp `handlers/masukanlayar_test.go:343` TestKirimanLayarTidakMenimpaUangMasterNonProp |
| 19-21 | ✅ | `models/penggolong_test.go:185` TestPenggolongBerhentiDiBarisPertama, `models/penggolong_test.go:193` TestPenggolongBawaanUnknown, `models/penggolong_test.go:145` TestKe128KodeSamaDenganSistemLama |
| 22 | ✅ | syarat When yang DIJALANKAN — `docs/INVENTARIS-XML.md` bab 11 |
| 23 | 🟡 | `NUMBER(38,10)` (diagram Prop F20, RALAT 04-10-2026); `repository/kolom_test.go:21` TestPecahAngkaEksak; Oracle `repository/pulangpergi_db_test.go:250` TestUangPresisiPenuhTanpaPembulatanRepository (38 digit utuh; galat lama `-592629512.880000276` utuh; 11 desimal dibulatkan Oracle di desimal ke-11, bukan ke-2) — **K11** |
| 24 | ✅ | nol pembulatan di repository (`pecahAngka` eksak, `repository/kolom_test.go:21` TestPecahAngkaEksak); pembulatan hanya di penyajian (`frontend/sajian.ts`, `frontend/sajian.test.ts`) |
| 25 | ✅ | `apd` di seluruh paket, nol `float` di `backend/`; `models/dokumenlama_test.go:345` TestAngkaJSONTidakLewatFloat, `models/nilaipega_test.go:11` TestTeksSkalarJSONTanpaFloat; penjaga `inti/backend/penjaga/migrasi_test.go:179` TestKolomUangDesimalDanNolJSON |
| 26 | ✅ | RALAT K3: `Deduction1/2` = uang (`models/katalog.go:193`; `frontend/medan.test.ts:175` "AC 85 + K3"); `BROKERAGE`/`RNM_SHARE` view: **nol pembaca** di rule terjangkau (`docs/alat/pemakai.py`; tiket 07 RALAT P4) |
| 27-28 | ✅ | `models/hitung_test.go:89` TestPajakBrokerageInclusiveDibagi1022, `models/hitung_test.go:99` TestPajakBrokerageSelainInclusiveApaAdanya, `models/hitung_test.go:134` TestTypeTaxHurufKecilMenggeserBalance; `models/setppnpph_test.go:20` TestSetPPNPPHNilaiXML |
| 29 | 🟡 | satu transaksi (`DalamTransaksi`, `services/gudang.go:98`); tiruan `handlers/alur_test.go:355` TestSatuTransaksiPembatalanUtuh; Oracle `repository/pulangpergi_db_test.go:301` TestKegagalanDiTengahTidakMenyisakanBarisDiTabelManaPun (kegagalan sesudah seluruh tulisan ⇒ nol baris di kesembilan tabel), `repository/polis_db_test.go:204` TestGenerasiTertutupDitolakDanPembatalanUtuh — **K11** |
| 30 | ✅ | `db.Qualify` di setiap SQL; penjaga `inti/backend/penjaga/lintasaplikasi_test.go:541` TestNolNamaTabelTelanjangDiQuery |
| 31 | 🟡 | indeks unik fungsi `CASE (NOPOLIS, PRODKE)` (`migrations/320_t_general_polis_treaty.sql:117`) + penomor `FOR UPDATE`; `repository/nomorpolis_db_test.go:102` TestNomorPolisDariDeretTanpaBentrok (deret penomor bersama, F8), `repository/polis_db_test.go:170` TestNomorPolisSekaliDanUnik — **K11** |
| 32-33 | ✅ | kolom `DATE` (DDL 320–327); satu format tukar (`repository/kolom_test.go:75` TestNilaiTulisMenolakMasukanRusak); tampilan satu format `frontend/sajian.test.ts:51` "AC 33" |
| 34-35 | ✅ | `models.PraprosesAdmin` (`models/polisaturan.go:36`); `models/tangga_test.go:125` TestPraprosesAdminTanggalKosongHariIni. `SystemSetOneYear_DT` hanya saat admin mengubah StartDate (`models/tanggal_test.go:10` TestSetahunDari29FebruariJatuhPada28Februari, `handlers/tanggal_test.go:16` TestUbahTanggalMulaiMengisiTanggalAkhirSetahun) |
| 36-38 | ✅ | `handlers/alur_test.go:435` TestPilihBisnis (422, nol simpan); `handlers/nonprop_test.go:207` TestNonPropMasterTidakAdaMenghentikanPilihBisnis; popup bisnis tanpa simpan `handlers/daftarbisnis_test.go:149` TestDaftarBisnisTidakMenyimpan |
| 39-41 | ✅ | `HISTORYAKSEPTASIPEGA.OPERATORID` = identitas login, `USERNAME` = nama tampilan; `HISTORYAKSEPTASIPRODUCTION.PIC` = nama tampilan — `handlers/alur_test.go:155` TestAdminMenolakDiselesaikanDitolak |
| 42 | ✅ | `models.PraprosesAtasan` (`models/polisaturan.go:94`) memakai nama tampilan (P33); `NamaTampilan` tanpa jatuh-balik ke login (`repository/riwayat.go:146`) |
| 43 | ✅ | satu baris riwayat per submit — `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali (5 submit → 5 baris) |
| 44 | ✅ | K5: NBStatus dari data — `models/tangga_test.go:116` TestTeksNBStatus, `models/layar_test.go:215` TestPascaAdminTolakNBStatus, `handlers/alur_test.go:155` TestAdminMenolakDiselesaikanDitolak |
| 45 | 🟡 | medan wajib dua layar realisasi + `ListSuggest` ditegakkan saat Submit DAN Save (`models.MedanWajibBerlaku`, `models/layar.go:176`; `handlers/alur_test.go:314` TestMedanWajibMenahanKirimDanSimpan; wajib mengikuti wadah tampil `models/layar_test.go:102` TestMedanWajibIkutWadahTampil); `ProductionDate` lewat tempat berperan — penahan **K12** (mekanisme `handlers/alur_test.go:534` TestTanggalProduksiMengikutiPemetaanTempat); `DateofSurvey` ⛔ (b) **K7** |
| 46-47 | ✅ | `models/layar_test.go:45` TestMedanWajibAtasanTanpaEnamMedanAdmin |
| 48 | ✅ | `handlers/alur_test.go:314` TestMedanWajibMenahanKirimDanSimpan (Submit dan Save) |
| 49-51 | ✅ | `models.GabungMasukanLayar` (`models/layar.go:406`): daftar izin per posisi **dan** syarat tampil sel/wadah (R6 W4), hasil tombol Enable/Disable hanya dari DT server (W3), kolom hanya-baca dari server (W5); `handlers/alur_test.go:369` TestMedanTerkunciAtasanDanTurunanAdmin, `models/layar_test.go:137` TestGabungMasukanAtasanHanyaMedanTerbuka, `models/layar_masukan_test.go:21` TestMedanAdminTersembunyiTidakDiterima, `handlers/masukanlayar_test.go:378` TestEnableDisableHanyaDariTombolXOL, `handlers/masukanlayar_test.go:132` TestKolomHanyaBacaSpreadingTidakDariLayar, `handlers/masukanlayar_test.go:416` TestAngsuranDariActionSetServer |
| 52 | ✅ | **RALAT R6 (W2)**: atasan mengisi Approval + Suggest (+ ProductionDate bersyarat) dan — polis NonProp baru — grid `SpreadingRiskList` NonProp bila `FacultativeShare` 0/''; `models/layar_test.go:167` TestAtasanHanyaMengisiPutusanCatatanDanTanggalProduksiBersyarat, `handlers/masukanlayar_test.go:52` TestAtasanMenyuntingSpreadingNonProp, `handlers/masukanlayar_test.go:84` TestAtasanSpreadingTerkunciBilaFakultatifAtauProporsional, `frontend/medan.test.ts:120` "AC 52", `frontend/components/DetailNonProp.test.tsx:55` "W2" |
| 53 | ✅ | elemen `1=2` / `NEVER` tidak dibangun — `frontend/medan.test.ts:112` "elemen mati"; termasuk grid lama S11 popup (`BrowseTreatyInDetail`, wadah `1=2`, R5) `handlers/daftarbisnis_test.go:168` TestRuteDaftarBisnisLamaTidakAda |
| 54 | ✅ | `repository.DaftarMataUang` `CURRENCY <> 'ITL'` (`repository/acuan.go:427`; RD `BrowseCurrencyTreatyIn_RD` filter B) |
| 55-56 | ✅ | kode tampil apa adanya; medan berpilihan "associated values" (tidak ada di korpus) = isian teks, diterima dan disimpan (tiket 11) |
| 57 | ⛔ | (b) **K8 butir 4**: treaty keluar seluruhnya JSON `M_TREATY_OUT` — `RDBList/BrowseTreatyOut`: `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}`; tidak dibaca |
| 58 | ✅ | nol penulisan treaty keluar / master; jalur XOL/NonProp dibangun ke `T_POLIS_XOL`, `T_POLIS_XOL_LAYER`, `T_POLIS_INSTALMENT(_DETAIL)`, `T_POLIS_SPREADING` — `handlers/nonprop_test.go:82` TestNonPropPilihBisnisHitungSimpanBacaKembali |
| 59 | ✅ | RALAT P4: `TreatyRealizationCheckDuplicate` — `handlers/alur_test.go:336` TestPolisSerupaMenahan, `handlers/logika_test.go:154` TestCekPolisSerupaHanyaSaatAdminMenyetujui |
| 60 | ✅ | penjaga `inti/backend/penjaga/lintasaplikasi_test.go:491` TestHandlersTidakMengimporRepository |
| 61-63 | ✅ | rule yatim, pembongkar JSON (RALAT K8: empat tetap tidak dimigrasi, dua dimigrasi baca-saja), `When/isApproved` — `docs/INVENTARIS-XML.md` bab 2, `docs/alat/status.json` |
| 64-65 | ✅ | kolom `isApprovedtoDeptHead` tidak dibuat; `LetterNo` bukan penanda arah — `repository/kolom_test.go:155` TestKatalogSepakatDenganDDL |
| 66 | ⛔ | (a) RALAT P4: `isFOR` dikirim **kosong** oleh `Flow/InputRealizationTreatyIn` Utility1; `SaveJsonPolisTreatyIn_Act` langkah 3 (EDM) `//`, langkah 5 (POLICY) When tak dicentang — nol efek; NB hanya jalur polis baru |
| 67 | ✅ | `models.KodeLiniJiwa` (`models/penggolong.go:115`); `models/penggolong_test.go:208` TestLiniJiwaEnamBelasKode |
| 68 | 🟡 | pemuat dokumen lama (tiket 22) menulis lewat `SimpanHalaman`, berkas dibuka `BacaHalaman` — `repository/lama_db_test.go:74` TestPemuatLamaMenulisLewatAntarmukaSama (**K11**); pemuatan produksi = urutan **F7** (`MODUL.md` bab *Migrasi*: migrasi WO → uji-kering K11 → F3 tuntas → pemuatan WO/DBA) |
| 69 | ✅ | `models.PecahDokumenLama` (`models/dokumenlama.go:521`): EndDate kosong = StartDate — `models/dokumenlama_test.go:70` TestPecahDokumenProporsionalDatar |
| 70 | 📄 | penyimpangan dicatat di kode, tiket, dan bab 5 berkas ini |
| 71 | ✅ | K4 (+F2 disetujui WO): catatan ke `HISTORYAKSEPTASIPRODUCTION` (TGL_INP, PIC, AKSES_LOGIN) dan dibaca balik — `handlers/alur_test.go:155` TestAdminMenolakDiselesaikanDitolak, `models/usulan_test.go:76` TestBarisCatatanDariRiwayatProduksi |
| 72 | 🟡 | `DaftarRiwayat` `ORDER BY TGL_TRANSFER, ROWID` (`repository/riwayat.go:45`); `repository/riwayat_db_test.go:54` TestDaftarRiwayatBerurutWaktu — **K11** |
| 73 | ✅ | `models.RakitNomorPolis` (`models/polisaturan.go:252`); `models/tangga_test.go:216` TestRakitNomorPolis; `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali |
| 74 | ✅ | nomor kedua = nomor pertama — `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali; `SetelNomorPolis` hanya bila `NOPOLIS IS NULL` (`repository/polis.go:241`) |
| 75 | ✅ | `models/tangga_test.go:162` TestNonProporsionalDitandai |
| 76 | ✅ | `models/tangga_test.go:184` TestHasFacOutHanyaDuaKode |
| 77 | ✅ | halaman tidak membawa pesan antar permintaan; validasi kirim dihitung atas salinan |
| 78 | ✅ | `SetValidateInstallment` langkah 1 (`Page-Clear-Messages`) tidak dibangun — `models/angsuran.go:95` |
| 79 | ✅ | setiap rumus berkutip `PropertiesValue` langkahnya (`models/hitung.go`, `models/nonprop.go`); uji berharapan XML untuk setiap Activity/When/DecisionTable dibangun (`docs/AUDIT-SILANG-PUTARAN-3.md` bab 3; mis. `models/hitung_onp_test.go:33` TestCountResult1OnpPctDanAmount) |
| 80 | ✅ | layar Dept Head dapat disimpan dan disubmit — `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali, `handlers/nonprop_test.go:82` TestNonPropPilihBisnisHitungSimpanBacaKembali; grup treaty kosong menahan submit Dept Head seperti XML (`handlers/masukanlayar_test.go:479` TestGrupTreatyKosongMenahanSubmitDeptHead) |
| 81 | ✅ | 12 tempat terdaftar, pemetaan kosong ⇒ tertunda — `handlers/alur_test.go:505` TestTempatBerperanTidakDitebak, `models/peran_tempat_test.go:19` TestPemetaanPeranTempatKosongSampaiIAMMenjawab, `models/peran_tempat_test.go:37` TestDuaBelasTempatTerdaftar |
| 82 | ✅ | arah `MUNCUL`/`KECUALI` tidak ditebak — `models/peran_tempat_test.go:95` TestTempatTampilMenurutArah, `frontend/tempat.test.ts` |
| 83 | 🟡 | gagal `CatatRiwayat`/`CatatUsulan` membatalkan submit (`handlers/alur_test.go:355` TestSatuTransaksiPembatalanUtuh, tiruan); Oracle `repository/riwayat_db_test.go:106` TestGagalCatatRiwayatMembatalkanSubmit — **K11** |
| 84 | ✅ | K6: tabel keputusan kosong = disetujui (`models/tangga_test.go:18` TestIsApprovedSelainNolDisetujuiTermasukKosong); submit tanpa Approval 422 di tiga jenjang (`handlers/logika_test.go:127` TestSubmitTanpaApprovalDitolakDiSetiapJenjang) |
| 85 | ✅ | `KotakMedan` kode mata uang di samping angka — `frontend/medan.test.ts:175` "AC 85 + K3" |
| 86 | ✅ | RALAT K14: format sel Section (`pyDecimalPlaces`, `pySeparators`) — `frontend/sajian.ts`, `frontend/sajian.test.ts`, `frontend/medan.test.ts`; format grid master NonProp (R6 W6) `frontend/nonprop.test.ts` |
| 87 | ✅ | KEPUTUSAN-RONDE-12 butir 7: muatan 4 medan sesudah commit, gagal tidak membatalkan — `handlers/alur_test.go:662` TestKonversiSesudahSelesai; sambungan `[terbuka]` (PERMINTAAN C1) |
| 88 | ✅ | rule tak terjangkau tidak dibangun — `docs/INVENTARIS-XML.md` bab 1–2; rantai per rule (a) di `docs/AUDIT-SILANG-PUTARAN-3.md` bab 4 |
| 89 | 🟡 | kolom view yang hilang = galat (`repository/kolom_test.go:232` TestKolomViewHilangAdalahGalat); lawan view sungguhan `repository/kontrak_db_test.go:99` TestViewKontrakMemuatSetiapKolomYangDibaca (melewati bila DBA tidak menyediakan view di skema uji — PERMINTAAN C4, C9) — **K11** |
| 90 | ✅ | tabel berejaan ganda (`HISTORYAKSEPTASIPEGA` dst.) lewat `Qualify` satu skema |
| 91 | ✅ | nol peran karangan; uji memakai peran `UJI-` (`handlers/alur_test.go:505` TestTempatBerperanTidakDitebak, `handlers/alur_test.go:534` TestTanggalProduksiMengikutiPemetaanTempat) |
| 92 | ✅ | berkas menunggu posisi — `handlers/alur_test.go:139` TestBukanAnggotaAntreanDitolak, `handlers/portal_test.go:49` TestGerbangDaftarPortal |
| 93-95 | 📄 | butir dokumen spec |
| 96 | ✅ | nol nama orang di artefak: riwayat cabang ditulis ulang (K1); fixture `UJI-`; identitas XML disamarkan `<ID-operator-N>` (`docs/alat/inventaris.py`) |

## 3 · `spec-penyimpanan-relasional.md` — 66 acceptance criteria

| AC | Status | Bukti / alasan |
| ---: | :---: | --- |
| 1 | 🟡 | indeks unik fungsi `CASE` (`migrations/320_t_general_polis_treaty.sql:117`); `repository/polis_db_test.go:170` TestNomorPolisSekaliDanUnik (dua baris bernomor-generasi sama ditolak Oracle) — **K11** |
| 2-5 | ✅ | `SisipKasus` PRODKE 0, `OLD_POLIS_ID` NULL + `UNIQUE` (`migrations/320_t_general_polis_treaty.sql:12`), PK bersama `T_WORK_POLIS` — `repository/kolom_test.go:313` TestTabelDanKolomMengikutiDiagramGrilling |
| 6 | 🟡 | generasi tertutup = ada penerus (`repository.syaratTerbuka`, `repository/polis.go:135`) → `ErrGenerasiTertutup`; `repository/polis_db_test.go:204` TestGenerasiTertutupDitolakDanPembatalanUtuh (perubahan tidak tersimpan, dibaca langsung) — **K11** |
| 7 | ✅ | nol kunci tamu NB↔EDM di `T_WORK_POLIS` |
| 8, 10 | ✅ | `NOURUT` + `UNIQUE (induk, NOURUT)` di setiap tabel anak (DDL 322–327) — `repository/kolom_test.go:313` TestTabelDanKolomMengikutiDiagramGrilling |
| 9 | 🟡 | hapus-sisip menomori ulang 1..n (`SimpanHalaman`, `repository/polis.go:67`); `repository/penyimpanan_db_test.go:85` TestNourutDinomoriUlangSesudahBarisKeduaDihapus — **K11** |
| 11 | ✅ | NB: tidak ada pemasangan antar generasi (milik EDM) |
| 12 | 🟡 | kode tetap `VARCHAR2`; `repository/penyimpanan_db_test.go:110` TestKodeBernolDepanUtuhDiKolom (`GROUP_PANEL` = `006`, `DATA_TYPE` `VARCHAR2`) — **K11** |
| 13 | 🟡 | `BusinessOldId "01"` pulang-pergi `repository/polis_db_test.go:68` TestPulangPergiHalamanLewatKatalog; kolom `BUSINESS_OLD_ID` dibaca langsung `repository/penyimpanan_db_test.go:110` TestKodeBernolDepanUtuhDiKolom — **K11** |
| 14 | ✅ | `models/penggolong_test.go:161` TestPenggolongKasusTangan (006/01 → FireStyle2); `handlers/alur_test.go:435` TestPilihBisnis (atas kolom tersimpan) |
| 15 | ✅ | RALAT P11: `""` tersimpan NULL dan terbaca `""`, `"0"` tetap `'0'` — `repository/kolom_test.go:43` TestNilaiTulisKosongJadiNULL, `repository/kolom_test.go:93` TestNilaiBaca; db `repository/penyimpanan_db_test.go:140` TestIsApprovedKosongDanNolTetapBerbeda (K11) |
| 16 | ✅ | `models.Disetujui`; `models/tangga_test.go:12` TestIsApprovedNolDitolak, `models/tangga_test.go:18` TestIsApprovedSelainNolDisetujuiTermasukKosong |
| 17-18 | ✅ | `repository/kolom_test.go:43` TestNilaiTulisKosongJadiNULL |
| 19-20 (+20b) | ✅ | ⛔ RALAT 04-10-2026 (diagram Prop F20 *"skala MINIMAL 9 desimal"*): `NUMBER(38,10)` (DDL 320, 323–327; `repository/kolom_test.go:155` TestKatalogSepakatDenganDDL, `repository/kolom_test.go:423` TestSkalaUangPersenMinimalSembilanDiSemuaTabel); pengikatan tanpa pemotongan (`repository/kolom_test.go:21` TestPecahAngkaEksak); galat lama utuh `592629512.880000276` diperiksa `repository/lama_db_test.go:74` TestPemuatLamaMenulisLewatAntarmukaSama (K11) |
| 21-22 | ✅ | `models.BacaTanggalLama` (`models/nilaipega.go:80`); `models/dokumenlama_test.go:18` TestBacaTanggalLama, `models/dokumenlama_test.go:39` TestTanggalAmbiguTidakDitebak (K15) |
| 23-24 | ✅ | kolom `DATE`; nol FLOAT (penjaga `inti/backend/penjaga/migrasi_test.go:179` TestKolomUangDesimalDanNolJSON) |
| 25 | ✅ | RALAT P11 (AC 25 + ID-20): nol pembandingan dua nilai uang hidup — `models/pembandingan_uang_test.go:97` TestPortTidakMembandingkanDuaNilaiUang, `models/pembandingan_uang_test.go:156` TestTandaUangLawanNolEksakSepertiXML |
| 26 | ✅ | tanpa `LAYER*` di `T_GENERAL_POLIS_TREATY` (diagram F26) — `repository/kolom_test.go:313` TestTabelDanKolomMengikutiDiagramGrilling; `handlers/alur_test.go:435` TestPilihBisnis (dibaca balik dari view) |
| 27 | ✅ | `T_POLIS_QUOTATION` = 10 medan diagram J37 + 6 RALAT berbukti XML (`docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 2) — `repository/kolom_test.go:313` TestTabelDanKolomMengikutiDiagramGrilling |
| 28 | ✅ | `T_POLIS_CEDING` `CEDING_CO_ID` + `CEDING_CO_NAME` (`migrations/322_t_polis_ceding.sql`, diagram R43) |
| 29-30 | ✅ | disalin apa adanya (`models/dokumenlama_test.go:70` TestPecahDokumenProporsionalDatar); ceding dihapus per baris |
| 31, 33 | ✅ | `models.PeriksaBentukSimpan` (`models/katalog.go:440`); `models/layar_test.go:12` TestBentukProporsionalMenolakXOLDanRincian; AC 33 juga dasar syarat ketiga `PerluCekDaftarXOL` (`models/nonprop_detail.go:452`; `models/nonprop_master_test.go:68` TestPerluCekDaftarXOL, bab 5 butir 30) |
| 32 | ✅ | akibat AC 33 |
| 34-35 | ✅ | `LAYER*` hanya di `T_POLIS_XOL_LAYER`; `DEDUCTION` uang (DDL 326/327); `handlers/nonprop_test.go:82` TestNonPropPilihBisnisHitungSimpanBacaKembali |
| 36 | ✅ | `models/layar_test.go:34` TestSpreadingBagiRataPresisiSepuluh |
| 37 | ✅ | persen `NUMBER(38,10)` (RALAT 04-10-2026, diagram F20); bagi rata NB presisi 10 utuh di kolom: `repository/pulangpergi_db_test.go:344` TestSpreadingBagiRataPresisiSepuluhUtuhDiKolom (J69, K11) |
| 38 | ✅ | RALAT K3: `DEDUCTION1/2` uang; `TOTAL_SHARE_PERCENTAGE_*` persen turunan (`HitungTotalSpreading`, `models/angsuran.go:197`), tak berkolom |
| 39 | ✅ | K4: tabel warisan ditulis, tidak dibuat (`MODUL.md` *Tabel warisan*); `repository/kolom_test.go:246` TestSQLRiwayatProduksiMengikutiInsertViewSuggest |
| 40-41 | ✅ | `models/usulan_test.go:12` TestUsulanBelumTersimpanMenurutSaveViewSuggest, `models/usulan_test.go:66` TestApprovalPerBarisHanyaAcceptReject; dokumen lama `models/usulanlama_test.go:31` TestSuggestListLamaDisalinMenurutSaveViewSuggest; db `repository/polis_db_test.go:268` TestRiwayatProduksiPulangPergi (K11) |
| 42 | ✅ | `APPROVAL` per baris terpisah dari `T_GENERAL_POLIS_TREATY.IS_APPROVED` — `models/usulan_test.go:76` TestBarisCatatanDariRiwayatProduksi |
| 43-44 | ✅ | `AKSES_LOGIN` = identitas login, `PIC` = nama tampilan — `models/usulan_test.go:12` TestUsulanBelumTersimpanMenurutSaveViewSuggest, `handlers/alur_test.go:155` TestAdminMenolakDiselesaikanDitolak |
| 45 | 🟡 | satu transaksi; pembatalan atas tiruan (`handlers/alur_test.go:355` TestSatuTransaksiPembatalanUtuh); Oracle `repository/penyimpanan_db_test.go:175` TestGagalTulisTabelAnakMembatalkanInduk — **K11** |
| 46 | 🟡 | satu transaksi per tindakan, nol `COMMIT` di repository; `repository/pulangpergi_db_test.go:301` TestKegagalanDiTengahTidakMenyisakanBarisDiTabelManaPun, `repository/polis_db_test.go:204` TestGenerasiTertutupDitolakDanPembatalanUtuh — **K11** |
| 47 | ✅ | penjaga `TestNolNamaTabelTelanjangDiQuery`; `inti/backend/penjaga/migrasi_test.go:54` TestSetiapPernyataanSahDanBerskema |
| 48 | ✅ | nol stored procedure (penomor menulis SQL, F8); `repository/lama_test.go:16` TestSQLPemuatLamaBerskemaTanpaCommit |
| 49 | 🟡 | `repository/pulangpergi_db_test.go:157` TestPolisProporsionalSeluruhMedanPulangPergi (SETIAP kolom katalog lima tabel bentuk proporsional), `repository/polis_db_test.go:68` TestPulangPergiHalamanLewatKatalog — **K11** |
| 50 | 🟡 | `repository/pulangpergi_db_test.go:174` TestPolisNonProporsionalSeluruhLayerPulangPergi (dua XOL × tiga layer), `repository/lama_db_test.go:186` TestPemuatLamaNonProporsionalBersarang — **K11** |
| 51 | 🟡 | `repository/pulangpergi_db_test.go:197` TestUrutanBarisAnakMenurutNourut (12 baris, urutan teks ≠ urutan tulis) — **K11** |
| 52-54 | ✅ | `models/dokumenlama_test.go:70` TestPecahDokumenProporsionalDatar, `models/dokumenlama_test.go:133` TestPecahDokumenNonProporsionalBersarang, `models/dokumenlama_test.go:167` TestUjiPemecahMencakupDuaBentuk |
| 55 | 🟡 | pemecah membawa nilai utuh (`models/dokumenlama_test.go:70` TestPecahDokumenProporsionalDatar); kolom `592629512.880000276` **utuh** (RALAT 04-10-2026, `NUMBER(38,10)`; semula `592629512.88000028`) di `repository/lama_db_test.go:74` TestPemuatLamaMenulisLewatAntarmukaSama — **K11** |
| 56 | ✅ | `repository/lama_test.go:50` TestPemuatTanpaJalurTulisTerpisah |
| 57 | ✅ | **RALAT F3** (putaran 3): CSV = arsip audit `POLIS_ID,JALUR,NILAI,KEPUTUSAN` — `models/laporanlama_test.go:34` TestLaporanArsipMedanTanpaKolomBerkasCSV |
| 58 | ✅ | `models/laporanlama_test.go:99` TestLaporanGalatBerkasCSVDanTanggalAmbiguDihitung, `models/dokumenlama_test.go:287` TestDokumenBergalatTidakDimuatDanSebabnyaDisebut |
| 59 | 🟡 | **RALAT F3**: "nol medan belum diputuskan". Mekanisme ✅ — kode keluar ≠ 0 hanya bila ada medan `BELUM DIPUTUSKAN` / galat (`models/laporanlama_test.go:144` TestRingkasanSelesaiHanyaBilaNolGalatDanNolBelumDiputuskan); seluruh 378 jalur daun panduan bentuk dokumen sudah diputuskan (`models/dokumenlama_test.go:367` TestPanduanBentukDokumenNolMedanBelumDiputuskan). Panduan basi (diagram R45) — atas data nyata baru pasti pada uji-kering pemuat: penahan **F7** langkah 2 (skema uji, K11) |
| 60-61 | ✅ | riwayat cabang bersih (K1); inventaris disamarkan |
| 62 | ✅ | fixture `UJI-` |
| 63 | 📄 | migrasi dijalankan manusia (work owner; K13, K18) |
| 64 | ✅ | 79 medan diagram F10 = 69 kolom katalog + `NOPOLIS` + 9 tak berkolom bersebab + 7 kolom `json_polis`; putaran 3 menambah `EDM_TYPE` (F3) ⇒ 70 kolom katalog medan `PolicyTreatyIn` — `docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 1, 10; `repository/kolom_test.go:313` TestTabelDanKolomMengikutiDiagramGrilling |
| 65 | ✅ | tanpa penjaga sinkron |
| 66 | ✅ | penjaga `TestHandlersTidakMengimporRepository`; `repository` tidak mengimpor `services` |

## 4 · Perubahan sejak putaran 2

Status AC **tidak berubah** (spec ✅80/🟡10/⛔2/📄4; penyimpanan ✅53/🟡12/📄1). Yang berubah: bunyi AC (RALAT),
bukti, dan penahan. Status "lama" = HASIL putaran 2 (`27910561`).

### 4a · `spec.md`

| AC | Putaran 2 | Sekarang | Sebab |
| ---: | :---: | :---: | --- |
| 11, 92 | ✅ | ✅ | portal: kolom buatan *Position*/*No Polis*, spanduk NBStatus layar kasus dibuang; tautan `.Name` dipindah ke "Offer No" (R6 W6, bab 5 butir 33) |
| 15 | ✅ | ✅ | **F1 disetujui WO**: daftar medan master XOL **tertutup** + uji penjaga pembaca JSON tunggal (`23db3bc9`); popup bisnis kini view `TREATYINDETAILJOINEDM` (R5 W1) |
| 16, 48 | ✅ | ✅ | **F8 disetujui WO**: `GENERATE_SEQUENCE_NUMBER` ditulis lewat `inti/backend/penomor` — bukan lagi "⚠️ dicatat untuk ditinjau" |
| 17, 89 | 🟡 | 🟡 | **RALAT spec §5.1 (W1)**: RD yang dipakai `BrowseTreatyJoinEDM` (bukan `BrowseTreatyInDetail`); uji db `TestDaftarBisnisRDBrowseTreatyJoinEDM` ditambah (R5); penahan tetap K11 |
| 23 | 🟡 | 🟡 | harapan uji db lama *"`830.82191780804`"* mustahil di `NUMBER(38,8)` → `TestUangPresisiPenuhTanpaPembulatanRepository` (38 digit; skala 8) — nusare-03; ⛔ **p3r-sheetprop**: skala **10** (`NUMBER(38,10)`, diagram F20) — harapan kini `830.8219178081` dan galat lama utuh — K11 |
| 29 | 🟡 | 🟡 | uji db baru kegagalan di tengah ⇒ nol baris di **kesembilan** tabel (`TestKegagalanDiTengahTidakMenyisakanBarisDiTabelManaPun`) — K11 |
| 31 | 🟡 | 🟡 | uji db baru deret penomor bersama tanpa bentrok (`TestNomorPolisDariDeretTanpaBentrok`) — K11 |
| 36-38 | ✅ | ✅ | popup bisnis = `POST /kasus/{id}/bisnis` tanpa simpan, hanya admin & bukan XOL Retro (R5) |
| 49-51 | ✅ | ✅ | server menerima hanya sel yang **tampil dan terbuka** (W4), hasil Enable/Disable hanya dari DT (W3), kolom hanya-baca grid dari server (W5); `POST /hitung` hanya aksi sel terbuka di layar posisi (audit P1) |
| 52 | ✅ | ✅ | **RALAT AC 52 (R6 W2)** — bunyi lama P2: *"…atau menemukan `DueTo`, `FlagPPH`, `No Offer Slip`, atau medan lain dapat diisi, gagal"* → atasan juga menyunting grid `SpreadingRiskList` NonProp bila `FacultativeShare` 0/'' (US 11, §5.11 ikut di-RALAT) |
| 53 | ✅ | ✅ | grid lama S11 popup (`1=2`) kini juga tidak dibangun (R5) |
| 59 | ✅ | ✅ | tidak berubah (penanda — jangan tertukar dengan penyimpanan AC 59) |
| 68 | 🟡 | 🟡 | penahan *"pemuatan di produksi menunggu keputusan WO (PERMINTAAN F7)"* → **F7 diputuskan**: urutan resmi di `MODUL.md`; tetap 🟡 sampai langkah 4 (pemuatan WO/DBA) |
| 71 | ✅ | ✅ | **F2 disetujui WO**: tiga jenjang, `TGL_INP` 24 jam, `NOURUT` repository |
| 79, 88 | ✅ | ✅ | audit silang: +17 uji Go/+13 vitest berharapan XML; rantai bukti 30 rule (a) (`AUDIT-SILANG-PUTARAN-3.md` bab 3–4) |
| 80 | ✅ | ✅ | `FetchTreatyGroupOldID` langkah 3 kini dibangun: pesan *"Cannot fetch Treaty Group ID, Contact IT"* menahan Submit Dept Head (R6, 7.4) |
| 86 | ✅ | ✅ | format grid master NonProp per sel (R6 W6) |

### 4b · `spec-penyimpanan-relasional.md`

| AC | Putaran 2 | Sekarang | Sebab |
| ---: | :---: | :---: | --- |
| 1, 6 | 🟡 | 🟡 | uji db dikuatkan (nusare-03): dua baris bentrok ditolak Oracle; perubahan generasi tertutup tidak tersimpan — K11 |
| 31, 33 | ✅ | ✅ | AC 33 dinyatakan dasar syarat ketiga `PerluCekDaftarXOL` (R6, tiket 19) |
| 40-44 | ✅ | ✅ | `SuggestList` dokumen lama disalin ke `HISTORYAKSEPTASIPRODUCTION` (F3); F2 disetujui |
| 46 | 🟡 | 🟡 | + `TestKegagalanDiTengahTidakMenyisakanBarisDiTabelManaPun` — K11 |
| 49-51 | 🟡 | 🟡 | semula satu uji pulang-pergi sebagian medan → tiga uji SETIAP kolom / setiap layer / urutan `NOURUT` (`repository/pulangpergi_db_test.go`) — K11 |
| 57 | ✅ | ✅ | **RALAT AC 57 + ID-27 (F3)** — bunyi lama K17: *"test yang menemukan medan tak dikenal tidak tertulis di berkas itu beserta nilainya gagal"* → CSV = **arsip audit** `POLIS_ID,JALUR,NILAI,KEPUTUSAN`; uji `TestLaporanMedanTakDikenalBerkasCSV` → `TestLaporanArsipMedanTanpaKolomBerkasCSV` |
| 59 | 🟡 | 🟡 | **RALAT AC 59 + ID-27 (F3)** — bunyi lama K17: *"berkas itu wajib nol baris data"* → "nol medan **belum diputuskan**"; penahan **F3 → F7** langkah 2 (uji-kering, K11) |
| 64 | ✅ | ✅ | kolom baru `T_GENERAL_POLIS_TREATY.EDM_TYPE` (F3, bab 8) |
| 19-20, 20b, 37 | ✅ | ✅ | **RALAT 04-10-2026 (p3r-sheetprop, diagram Prop F20)**: `NUMBER(38,8)` → `NUMBER(38,10)`; AC 19 kembali *"tanpa kehilangan satu digit pun"* untuk `592629512.880000276`; 20b "> 8" → "> 10" desimal |
| 55 | 🟡 | 🟡 | pembulatan "desimal kesembilan" → "kesebelas" (`NUMBER(38,10)`) — K11 |
| 6 | 🟡 | 🟡 | + `UPDATE` kolom datar pemuat lama hanya generasi terbuka (diagram F17; `repository/lama_test.go:82` TestUbahGeneralPolisHanyaGenerasiTerbuka) — K11 |

### 4c · RALAT dokumen putaran 3 (ringkas)

| Dokumen | Bunyi lama → bunyi baru | Paket |
| --- | --- | --- |
| `spec.md` AC 52, US 11, §5.11 | atasan hanya Approval/Suggest/ProductionDate → + grid spreading NonProp bersyarat `FacultativeShare` | R6 |
| `spec.md` §5.1, §9.2 butir 6 | RD `BrowseTreatyInDetail` (33 medan) → `BrowseTreatyJoinEDM`; *"Tipe tabel `TREATYINDETAIL` belum tercakup"* → gugur, tabel tidak dibaca | R5 |
| `spec-penyimpanan-relasional.md` ID-27, AC 57, 59 | penampung CSV wajib nol baris → arsip audit, nol medan belum diputuskan | R1 |
| `rancangan-tabel-datar-treaty-in.md` §4sexies, §4sx.5, §4bis.4 | kolom `EDM_TYPE`; F8 "dibaca saja" → ditulis lewat penomor | R1, R3 |
| `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md` butir 5 poin 2 | penanda `SUMBER='PEGA'` → gugur (F6) | R1 |
| `PERBANDINGAN-KOLOM-DIAGRAM.md` bab 1c/1d/5/9/11 | `SetPPNPPH` langkah 3 → 4; `GeneratePolicyNoTreaty_Act` langkah 23 `//`; butir terbuka medan lama → diputuskan | R1 |
| tiket 01, 10, 11, 12, 19, 22 | F1/F2 label; pemilih SOB tanpa simpan (F4); popup W1; W2–W6; 7.4; F3 per medan | R1–R6 |
| `docs/alat/status.json` | `TreatyInInputVis` langkah 3 seolah berjalan → keluar di langkah 1 (`1==1`→6); `GetCurrentDate` "parameter viewstate dan ID" → `pyPassCurrentParameterPage=true`; 5 aktivitas tak terpicu → (a) | P3K |
| `MODUL.md`, `PERMINTAAN-TIM-INTI.md` C3 | *"tabel `TREATYINDETAIL` (grid popup)"* → gugur | P3K |
| `spec-penyimpanan-relasional.md` ID-14, ID-15, AC 19, 20, 20b, 55, bab *Presisi*; `rancangan-tabel-datar-treaty-in.md` §4q5.5, §4sx.6; `PERBANDINGAN-KOLOM-DIAGRAM.md` | *"`NUMBER(38,8)` … desimalnya tetap delapan"* → **`NUMBER(38,10)`** (diagram sheet NB Treaty In Prop F20 *"skala MINIMAL 9 desimal (P29)"*, J69 presisi 10); kepatuhan butir demi butir `docs/KEPATUHAN-SHEET-NB-TREATY-IN-PROP.md` | p3r-sheetprop |

## 5 · Penyimpangan sadar

### 5a · Berdasar keputusan tertulis

| # | Hal | XML | Yang dibangun | Dasar |
| ---: | --- | --- | --- | --- |
| 1 | Sec Head menyetujui | ≤ 200 juta & bernomor → selesai; > batas → Dept Head (`CekLimitTreatyAcc_Act`) | selalu naik ke Dept Head | **K2**, AC 8 |
| 2 | Tombol Generate/Submit atasan | identitas `<ID-operator-1>` | menurut posisi kasus; enam tombol Submit terdaftar sebagai tempat tetapi tidak dibaca layanan | P13, AC 12, tiket 05 |
| 3 | Sumber data kontrak | JSON master (`FetchMasterTreatyIn`, `BrowseTreatyIn`) | view `TREATYINDETAILJOINEDM`; jalur NonProp/XOL membaca JSON `M_TREATY_IN`/`_EDM` baca-saja di satu fungsi | P29 + **K8** + **F1** (disetujui WO 04-10-2026) |
| 4 | Master jalur Proporsional (`RNMShareP`, `BrokeragePercentP`, `SpreadingTotalPct` …) | JSON master | langkahnya dilewati (`models.MasterTersedia`) | (c) — di luar ukuran K8 (F1) |
| 5 | Hari tutup buku pra-proses | `>25` tertanam (`InputPolicyTreatyInPre_Act` 9) | hari dari `TANGGAL_CLOSING` (`handlers/logika_test.go:67` TestTanggalProduksiDariTanggalClosing) | preseden WO PremiumList Life |
| 6 | NBStatus | nama orang tertanam di connector | nama posisi tujuan / nama tampilan | **K5**, AC 44 |
| 7 | OperatorName atasan | `pyUserIdentifier` | nama tampilan | P33, AC 42 |
| 8 | Riwayat akseptasi | tanpa skema, COMMIT per RDB | skema eksplisit, satu transaksi | AC 29, 30, 83 |
| 9 | Arasapas | stub kelas Data | bentuk kelas Work, 4 medan, sesudah commit | KEPUTUSAN-RONDE-12 butir 7 |
| 10 | `GetOldIDBusiness_SQL` | tanpa ORDER BY | `ORDER BY ID` | hasil tetap |
| 11 | Daftar portal | satu grid `GetListOpportunity` di wadah `pyWorkGroup!='ReasLife' && pyWorkBasketList(2)=='ReasTreatyInAdmin'` | admin: semua kasus terbuka; Sec/Dept Head: posisi yang dipegang; klausa `pyWorkGroup` tidak dibangun (K12, C6) | P8, P9; AC 11, 14 |
| 12 | Submit tanpa Approval | tombol Submit hanya untuk 1/0 | 422 di tiga jenjang | **K6**, AC 84 |
| 13 | Catatan usulan — syarat `BusinessFac=="F"` | treaty tak pernah menulis | selalu ditulis ke `HISTORYAKSEPTASIPRODUCTION` | **K4** |
| 14 | `HISTORYAKSEPTASIPRODUCTION.DIV` | `OperatorID.pyOrgDivision` | NULL — tanpa sumber di `inti.Pelaku` | PERMINTAAN B4 |
| 15 | Pemecah dokumen lama | ID-4: seam repository | fungsi murni `models.PecahDokumenLama` | RALAT ID-4 (P7) |
| 16 | Medan dokumen lama tanpa kolom | tabel penampung (ID-27) | diputuskan per medan: berkolom (`EDM_TYPE`), disalin (`SuggestList`), atau dibuang berbukti; CSV = arsip audit | **K17** + **F3** |
| 17 | Cap waktu ` GMT` dokumen/master lama | teks | jam dinding Asia/Jakarta; tanggal ambigu → laporan galat | **K15** |
| 18 | Tampilan tanggal | beberapa format | satu format (`formatDate` inti) | AC 33 |
| 19 | NonProp — halaman master | disimpan di blob kasus | dibaca ulang saat dibuka (`models.TampilanMasterNonProp`) | K8 |
| 20 | NonProp — master tidak ada saat pilih bisnis | `catch … oLog.error` | 422, nol simpanan | K8, AC 36–38 |
| 21 | NonProp — `TreatyRealizationCheckXOLList` | berjalan bila `IsNewPolicyNonProp` 1 | dilewati untuk `ProportionalType = 'Proportional'` — lihat butir 30 | spec-penyimpanan AC 33 |
| 22 | Breakdown spreading | `BreakDownSpreading_Act` | tidak dibangun | **K9** |
| 23 | Laporan survei historis | popup + `SetSurveyReport_Act` | tidak dibangun | **K7** |
| 24 | `GENERATE_SEQUENCE_NUMBER` | ditulis procedure `PROC_GENERATE_SEQUENCE_NUMBER`; diagram F103/F118 "dibaca saja" | ditulis SQL oleh `inti/backend/penomor` (nol procedure, AC 48) | **F8 — disetujui WO 04-10-2026** (RALAT catatan rancangan §4bis.4) |
| 25 | **K4 (1)** catatan usulan per jenjang | XML menulis hanya pasca-submit admin | setiap submit (admin, Sec Head, Dept Head), transaksi yang sama | **F2 — disetujui WO 04-10-2026** |
| 26 | **K4 (2)** `TGL_INP` | XML `hh` (jam 12) | jam 24 | **F2 — disetujui WO 04-10-2026** |
| 27 | **K4 (3)** `NOURUT` | `.pxListSubscript` | MAX+1 per IDPEGA di bawah kunci kasus — terbukti sama (`handlers/alur_test.go:273` TestNourutUsulanSamaDenganSubskripSuggestList) | **F2 — disetujui WO 04-10-2026** |
| 28 | Keanehan Pega jalur XOL | lihat PERMINTAAN G1–G17 | **ditiru apa adanya** | **F5 — disetujui WO 04-10-2026**; ditinjau Finance/Product lewat bagian G |
| 29 | Pemilih Source Of Business | hasil `SearchHierarkiSourceBizAgent_PostDT` dipegang clipboard sampai Save/Submit | **ikuti XML**: klik = pencarian tanpa simpan; layar memegang pilihan, server mencocokkannya ulang dengan RD `BrowseAgentHierarkiList_RD` saat Save/Submit (`terimaSumberBisnis` → `models.TerimaKirimanTerkunci`, `services/sumberbisnis.go`; `handlers/sumberbisnis_test.go:108` TestKlikBarisSumberBisnisTidakMenyimpan, `handlers/sumberbisnis_test.go:226` TestNilaiSumberBisnisPalsuDitolak) | **F4 — keputusan WO 04-10-2026** |
| 30 | `PerluCekDaftarXOL` syarat ketiga `QuotationData.ProportionalType != "Proportional"` | `InputPolicyTreatyInPre_Act` langkah 10 hanya `IsNewPolicyNonProp==1` dan `EDMType=="3"` | baris XOL tidak dibuat bagi polis yang AC 33 akan tolak (penanda `IsNewPolicyNonProp` tidak pernah diturunkan pra-proses) — `models/nonprop_master_test.go:68` TestPerluCekDaftarXOL | **spec-penyimpanan AC 33** (`[terverifikasi]`, ID-35, hasil grilling; tiket 19 catatan R6) |

### 5b · Penyesuaian sadar baru putaran 3 — tanpa keputusan WO khusus, dicatat untuk ditinjau

F1–F8 seluruhnya **diputuskan WO 04-10-2026** (bab 5a butir 3, 16, 24–29; PERMINTAAN F). Penyesuaian di bawah lahir
dari penerapan keputusan itu dan temuan audit; masing-masing tercatat juga sebagai butir terbuka di PERMINTAAN
(bagian H).

| # | Hal | XML | Yang dibangun | Rujukan |
| ---: | --- | --- | --- | --- |
| 31 | **R2** — salinan pilihan SOB ke `PolicyTreatyIn.QuotationData` saat diterima | PostDT menulis `Quotation.*` saja; penyalin ke `QuotationData` hanya pra-proses (`InputPolicyTreatyIn_preDT` 14, preACT 14.9, `GeneratePolicyNoTreaty_Act` 10); apakah `opener.location.reload` menjalankan ulang pra-proses tidak ada di korpus | pilihan yang diterima langsung disalin ke `QuotationData` (`SalinKeQuotationData` lewat `TerapkanPilihanSumberBisnis`, `models/sumberbisnis.go`) agar `SetPPNPPH` langkah 1 memakai nilai layar (`handlers/sumberbisnis_test.go:342` TestPilihanDipegangMenggerakkanPPNPPH); alternatif harfiah: salin hanya pada pra-proses berikut | PERMINTAAN H1 — konfirmasi WO |
| 32 | **R6 W5** — dua penyesuaian grid | kolom hanya-baca spreading dihitung `CountSpreading_Act` 4.1 hanya saat sel %Share berubah; urutan beberapa refresh dalam satu permintaan | (1) baris spreading yang dihapus tanpa perubahan %Share dihitung ulang (baris tak berkunci); (2) Installment yang berubah bersama sel uang dihitung dengan BalanceDueTo akhir (`models/layar.go:512`, `models/layar.go:602`) | tiket 11 R6 W5 |
| 33 | **R6 W6** — tautan `.Name` portal | kolom "Name" `.Name` pxLink → `openWorkByHandle(.pzInsKey)`; `.Name` ditulis nol rule ⇒ tautan selalu kosong | kolom "Name" tidak dirender; `openWorkByHandle` dipasang di sel "Offer No" (`.TextNoQuotation`) — tanpa ini berkas tak dapat dibuka dari portal | tiket 11 RALAT tautan `.Name` |
| 34 | ~~**R6** — anggota baris spreading bukan sel~~ ✅ **ditutup R7** | `Currency`, `CurrencyID`, `SplitRNMSharePct`, `TreatyName` bukan sel grid `SpreadingRiskList`; penulisnya hanya `TreatyInputPctCommSpreading` (TreatyName) dan `TreatyNonPropSetSpreading` 3/4.1/6 — bukan `CountSpreading_Act` | ~~masih diterima dari kiriman layar bersama barisnya~~ → **tidak diterima dari layar**: nilai baris server yang mengikuti barisnya; baris Add / nilai karangan kosong (`models/layar.go` `pasangAnggotaServer`; `handlers/masukanlayar_test.go` TestAnggotaBarisSpreadingDariServer) | PERMINTAAN H2 — selesai |
| 35 | **R7** — tombol `Choose` popup bisnis (pola F4) | `Choose` hanya di baris grid aktif `BrowseTreatyJoinEDM` (`SetValue_Act(ID=.ID)`) | `idDetail` di luar RD yang dijalankan ulang dengan saringan kasus → 422, nol simpanan (`models.PeriksaPilihanBisnis`; `handlers/daftarbisnis_test.go` TestPilihBisnisDiLuarDaftarPopupDitolak) | PERMINTAAN H2 — selesai |
| 36 | **R7** — satu pola kiriman terkunci | Enable / Disable Input Type (`TreatyEnableDisableInput`) dan Select Source Of Business (`SearchHierarkiSourceBizAgent_PostDT`) — hasil tombol dipegang layar sampai Save/Submit | satu fungsi `models.TerimaKirimanTerkunci` (kiriman = halaman server ditimpa jalur yang dikirim; diterima hanya bila sama dengan hasil hitung ulang; satu jalur 422 `jawabKiriman`). Enable / Disable kini mencocokkan KEDUA medan sekaligus (tombol menulis keduanya) | tinjauan standards P3 C4 |
| 37 | **R7** — `AKSES_LOGIN` SuggestList dokumen lama | `AddToListCommentsPolicyTreatyIn_DT` menulis `.Suggest/.IsApproved/.Date/.OperatorName`; dataguide tanpa `OperatorID` baris | `AKSES_LOGIN` selalu NULL; `OperatorID` baris (bila ada) = BELUM DIPUTUSKAN (`models/usulanlama.go`; TestOperatorIDBarisUsulanTidakDipetakan) | RALAT tiket 22 (tinjauan spec P3 (c)3) |

## 6 · Rule yang tetap tidak dibangun

Sumber: `docs/alat/status.json` (176 rule terjangkau): **103 dibangun** (termasuk yang sebagian) dan **73 tidak
dibangun** — **(a) 36**, **(b) 35**, **(c) 2**. Cara hitung (putaran 3, menyamakan selisih audit bab 2): huruf
**pertama** di baris status.json. `crmOpportunitiesList` ("(a) + (b)") dihitung (a); `SaveTreatyIn` ("(b) … juga
(a)") dihitung (b). Selisih dengan putaran 2 (109 / 29 / 36 / 2): `crmOpportunitiesList` (b) → (a); R5:
`ReportDefinition/BrowseTreatyInDetail` dibangun → (a); P3K: lima aktivitas **tidak terpicu** dibangun → (a).

### (a) Tidak terjangkau dari titik masuk nyata, tidak terpicu, atau efeknya dibaca nol rule — 36

Rantai pemanggilan dan nol pembaca per rule: `docs/AUDIT-SILANG-PUTARAN-3.md` bab 4 (30 rule) dan status.json.

| Kelompok | Rule |
| --- | --- |
| Cek penawaran ganda lama (langkah `//`, pemanggil tanpa parameter) | `Activity/CheckDuplicateOffer`, `RDBList/GetCountClaim`, `ReportDefinition/BrowseTREATY_IN` |
| Revisi master (`revisionstate==1`, tak pernah dari NB) / nol efek | `RDBList/GetCurrentDate`, `Activity/TreatyInInputVis`, `Activity/ConvertHistoryDate` |
| Varian EDM jalur NonProp (syarat mustahil sesudah NonProp langkah 10) | `Activity/InputPolicyTreatyEDMDetail_NP`, `Activity/InsertToTreatyOutXOLList`, `Section/DetailPolicyTreatyInNonProportionalEDM`, `Section/Installments_ReadOnly`, `FlowAction/Installments_ReadOnly` |
| Tombol bertampil `1=2` / grid berwadah `1=2` | `Activity/TreatyInNonSetTotal`, `ReportDefinition/BrowseTreatyInDetail` (grid lama S11 popup, R5) |
| ⭐ **Tidak terpicu** — sel pemicu ber-`pyReadOnly` / grid `readOnly` (audit 7.4; fungsi port dipertahankan berujian, `POST /hitung` menolaknya 409 di setiap posisi, `services/aksiposisi.go`) | `Activity/CountRiCommOgp_act`, `Activity/CountRiCommOnp_act`, `Activity/CountOverridingCommOgp_Act`, `Activity/CountOverridingCommOnp_Act` (sel `.ResultOgp1/Onp1/Ogp2/Onp2` `DetailDeptHeadTreatyIn_UW`), `Activity/CountPctInstallment_Act` (sel `.Premium` grid S45) |
| Lampiran (efek hanya ke `Protection_Act`, tidak ada di korpus; BusinessFac "T") | `Activity/SetCategoryAttach`, `Activity/InputParamUploadReas_act`, `DataTransform/setCategoryAttachment_DT`, `RDBList/AttachmentLife`, `RDBList/CategoryAttach_SQL` |
| Pemilih SOB — cabang yang ditulis/dibaca nol rule | `Activity/SearchHierarkiSourceBizAgentTreatyIn_Act`, `DataTransform/btnCedingCO_DT`, `RDBList/BrowseClientEmail_SQL`, `ReportDefinition/BrowseAgentNusaRe_RD`, `ReportDefinition/BrowseCedingCo_RD` |
| Hasil dibuang | `RDBList/GenerateNoPolicy` |
| Metadata kolom / elemen mati (juga (b) P44/AC 53) | `ReportDefinition/crmOpportunitiesList` |
| When yang tak pernah benar / tak mengubah perilaku | `When/IsClaim`, `When/isClaimTreaty`, `When/isSellingModeB2B`, `When/isSellingModeB2BB2C`, `When/isSellingModeB2C`, `When/pyIsIpadOrDesktop` |

`SetValidateInstallment_Act` dan `CountNetPremi_act` tetap **dibangun**: keduanya terjangkau sebagai sub-panggilan
(`CountOGPONP_Act` langkah 10 / 7), walau sel pemicu langsungnya tertutup.

### (b) Keputusan work owner tertulis — 35

| Keputusan | Rule |
| --- | --- |
| **K2** — Sec Head selalu ke Dept Head | `Activity/CekLimitTreatyAcc_Act`, `When/ToTREATYDEPTHEAD` |
| **K7** + bab 0 butir 11 — survei historis | `Activity/SetSurveyReport_Act`, `Activity/ConcatSlipOfferNo_Act`, `FlowAction/InputHistoricalSurveyReport`, `FlowAction/InputHistoricalSurveyReportUW`, `Harness/HistoricalSurveyReport`, `Harness/HistoricalSurveyReportUW`, `Section/HistoricalSurveyReportDtl`, `Section/HistoricalSurveyReportDtlUW`, `Section/InputHistoricalSurveyReportDtl`, `Section/InputHistoricalSurveyReportDtlUW` |
| **K8 butir 4** — treaty keluar, JSON `M_TREATY_OUT` | `Activity/InputPolicyTreatyOutDetail_NonProp`, `Activity/InputPolicyTreatyOutDetail_preACT`, `Activity/SetValueRetro_Act`, `Activity/TreatyNonPropOutSetSpreading`, `Harness/BusinessAndSOBListRetro`, `Section/BusinessAndSOBListRetro`, `Section/DetailPolicyTreatyOutNonProportional`, `RDBList/BrowseTreatyOut`, `RDBList/BrowseTreatyOutDetail`, `ReportDefinition/BrowseTreatyOutDetail` |
| **K9** + KEPUTUSAN-RONDE-12 butir 3/3b — breakdown spreading | `Activity/BreakDownSpreading_Act`, `RDBList/GetBreakDownSpread_SQL` |
| **AC 16** (nol tulis JSON) + AC 48 + diagram F11 | `Activity/SaveJsonPolisTreatyIn_Act`, `RDBList/SavePolisTreatyIn_SQL`, `RDBList/SaveTreatyIn` (juga (a)) |
| **P44 / AC 53** — elemen mati | `Section/SFAPortal_OpportunitiesList_Header` |
| **K5 / K12 / AC 12–14, 81** — identitas/atribut operator tanpa padanan | `When/IsSPVCreate`, `When/IsSPVTreaty1`, `When/IsTreaty1`, `When/IsUW`, `When/IsNotAdmin`, `When/IsOperatorLife`, `When/crmCreateOpportunity` |

### (c) Data hanya di JSON, di luar pengecualian K8 — 2

| Rule | Alasan (diperiksa ulang audit bab 5) |
| --- | --- |
| `Activity/FetchMasterTreatyIn` | master jalur Proporsional (`TreatyInputPctCommSpreading` langkah 1); `SpreadingTotalPct`/`SpreadingType*`, `RNMShareP`, `BrokeragePercentP` tidak ada di 39 kolom view dan tidak di daftar F1 |
| `RDBList/BrowseTreatyInDetailJoinEDM` | JSON `M_TREATY_IN_DETAIL_EDM` (`INSTALLMENT(n)`, peka huruf) — tabel di luar dua tabel pengecualian K8 |

Rule **dibangun sebagian** (langkah lain tidak, dengan alasan per langkah di status.json): `AgentSourceBizTreatyIn_Act`,
`CalculatePremi_Act` (c), `InputPolicyTreatyInDetail_preACT`, `SetTreatyIn_Act`, `TreatyInputPctCommSpreading` (c),
`ReportDefinition/GetListOpportunity`, `Section/SFAPortalOpportunitiesHeader`, dan beberapa Section/DT layar.
`Protection_Act` varian treaty **tidak ada di korpus** (`docs/INVENTARIS-XML.md` bab 12).

## 7 · Jalur Proporsional dan NonProporsional — dari layar sampai tersimpan

Satu halaman portal `nbtreatyin-portal` (`frontend/pages/PortalNBTreatyIn.tsx`) dan satu layar kasus
(`frontend/pages/LayarKasus.tsx`); rute `/api/nb-treaty-in` (`handlers/rute.go:45`–`:57`). Padanan
setiap tombol → activity: `docs/INVENTARIS-XML.md` bab 13 (`docs/alat/tombol.json`). Nol menu/halaman di luar Pega
(audit bab 6).

| Tahap | Proporsional | NonProporsional / XOL |
| --- | --- | --- |
| Buat kasus | `POST /kasus` → `BuatKasus` (`services/layanan.go:144`): `T_WORK_POLIS` (`NB-<SEQ_WORK_POLIS>`) + `T_GENERAL_POLIS_TREATY`, `BusinessFac` "T", posisi `ReasTreatyInAdmin` | sama |
| Buka layar | `GET /kasus/{id}` → `BukaKasus` (`services/layanan.go:193`): pra-proses admin/atasan; `LAYER*` dibaca balik dari view | + `IsNewPolicyNonProp` 1 (dan bukan `EDMType` 3, bukan `Proportional` — bab 5 butir 30) → `TreatyRealizationCheckXOLList`: master dibaca ulang, `InsertToTreatyXOLList(RetroShare)` |
| Popup pilih bisnis (tombol *Choose Business*, admin, ClaimType ≠ `XOL Retro`) | `POST /kasus/{id}/bisnis` → `DaftarBisnis` (`services/layanan.go:362`): RD `BrowseTreatyJoinEDM` atas view `TREATYINDETAILJOINEDM`, filter H `PROPORTIONTYPE = QuotationData.ProportionalType` (kosong diabaikan), urut `TREATYID`, ≤ 500; tanpa simpan; judul "Business And SOB List", 25 kolom XML, saring/urut per kolom, 50 per halaman (R5) | sama |
| Pilih bisnis | `POST /kasus/{id}/pilih-bisnis` → `PilihBisnis` (`services/tindakan.go:312`): `DetailKontrak`, `BUSINESS` (`GetOldIDBusiness_SQL`), komisi `RIONR`, preACT 1–6, 8, 11, 14, 15, 19; gagal baca → 422 | + preACT 16/18 → `pilihBisnisNonProp` (`services/nonprop.go:89`); master tidak ada → 422 |
| Pemilih SOB (ClaimType `XOL Retro`) — **F4** | `GET /sumber-bisnis` (pohon `BrowseAgentHierarkiList_RD`); `POST /kasus/{id}/pilih-sumber-bisnis` = **pencarian tanpa simpan** (jawab nilai `SearchHierarkiSourceBizAgent_PostDT`); layar memegang pilihan; **disimpan bersama Save/Submit** sesudah dicocokkan ulang di server (`terimaSumberBisnis`); hanya `SOURCE_OF_BUSINESS` berkolom | sama |
| Hitung (refresh sel) | `POST /kasus/{id}/hitung` → `Hitung` (`services/tindakan.go:231`): hanya aksi sel yang **terbuka** di layar posisi (`aksiTerbuka`, `services/aksiposisi.go:48`; lainnya 409) — tidak menyimpan | + `CountSpreading` grid `SpreadingRiskList` NonProp, juga bagi atasan bila `FacultativeShare` 0/'' (W2) |
| Simpan draf | `PUT /kasus/{id}` (admin) → `SimpanDraf` (`services/tindakan.go:285`): masukan disaring sel tampil+terbuka (W3/W4/W5), pilihan SOB dicocokkan, medan wajib, `PeriksaBentukSimpan`, `SimpanHalaman` | sama; `TreatyXOLList`/`ListInstallment` tidak diterima dari layar |
| Tersimpan di | `T_GENERAL_POLIS_TREATY`, `T_POLIS_QUOTATION`, `T_POLIS_CEDING`, `T_POLIS_INSTALMENT`, `T_POLIS_SPREADING` (XOL dan rincian angsuran **ditolak**, AC 31–33) | + `T_POLIS_INSTALMENT_DETAIL`, `T_POLIS_XOL`, `T_POLIS_XOL_LAYER`; halaman master tidak disimpan |
| Submit admin | `POST /kasus/{id}/kirim` → `Kirim` (`services/tindakan.go:407`): DT `isApproved`; Approval wajib (K6); `TreatyRealizationCheckDuplicate`; riwayat `HISTORYAKSEPTASIPEGA`; catatan `HISTORYAKSEPTASIPRODUCTION`; NBStatus — satu transaksi | sama |
| Sec Head / Dept Head | Sec Head setuju → Dept Head (K2); tolak → admin. Dept Head: `POST /kasus/{id}/nomor-polis` → `TerbitkanNomor` (`services/tindakan.go:593`; penomor bersama, F8) → submit (ditahan pesan grup treaty kosong, 7.4) → `Resolved-Completed` → konversi Arasapas | + atasan menyunting grid spreading NonProp bersyarat (W2) |
| Riwayat | `GET /kasus/{id}/riwayat` (panel History; dasar AC 72) | sama |
| Uji ujung ke ujung | `handlers/alur_test.go:202` TestTanggaPenuhDanNomorPolisSekali, `handlers/alur_test.go:435` TestPilihBisnis; `handlers/daftarbisnis_test.go:72` TestDaftarBisnisMenyaringJenisProporsiKasus; `handlers/sumberbisnis_test.go:198` TestSaveMenyimpanSumberBisnisYangCocok; db `repository/pulangpergi_db_test.go:157` TestPolisProporsionalSeluruhMedanPulangPergi (K11) | `handlers/nonprop_test.go:82` TestNonPropPilihBisnisHitungSimpanBacaKembali, `handlers/nonprop_test.go:226` TestNonPropCekDaftarXOLSaatDibuka; `handlers/masukanlayar_test.go:52` TestAtasanMenyuntingSpreadingNonProp; db `repository/pulangpergi_db_test.go:174` TestPolisNonProporsionalSeluruhLayerPulangPergi (K11) |
| Belum / tidak | master Proporsional dari JSON (c); `BreakDownSpreading` (K9); survei (K7); `CekLimitTreatyAcc_Act` (K2); `Protection_Act` (tak ada di korpus); 4 tempat `ProductionDate` (K12) | treaty keluar / XOL Retro `ReinsuranceListTONP` (K8 butir 4); varian EDM NonProp (a); PPN/PPh rincian angsuran tak berkolom (dibuang, F3) |

Rute lama yang **dihapus** putaran 3: `GET /bisnis` (R5, kini `POST /kasus/{id}/bisnis`;
`handlers/daftarbisnis_test.go:168` TestRuteDaftarBisnisLamaTidakAda). `POST /kasus/{id}/pilih-sumber-bisnis` tetap
ada tetapi tidak lagi menyimpan (F4).

## 8 · Tabel: delapan `CREATE TABLE` dan tabel warisan

Tepat **delapan** `CREATE TABLE` di migrasi modul (`backend/migrations/320`–`327`), sama dengan diagram grilling
`Diagram-Skema-Tabel-NusantaraRe.xlsx` (bab 0 butir 11), ditagih `repository/kolom_test.go:313` TestTabelDanKolomMengikutiDiagramGrilling
(nama, jumlah, kolom persis; *"TEPAT delapan CREATE TABLE"*) — **tidak dilemahkan** putaran 3. Nol tabel baru.
Perbandingan kolom: `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

⭐ **Kepatuhan sheet *NB Treaty In Prop* (perintah WO 04-10-2026 *"buat sesuai yang di sheet NB Treaty In Prop"*):**
setiap pernyataan sheet (B5–B122) dipetakan ke kode/migrasi/uji di `docs/KEPATUHAN-SHEET-NB-TREATY-IN-PROP.md`.
Selisih yang diperbaiki: **(1) F20** — uang dan persen `NUMBER(38,8)` → **`NUMBER(38,10)`** di seluruh delapan tabel
(skala minimal 9; 10 supaya bagi rata spreading NB presisi 10, J69, tersimpan utuh; 28 digit di depan koma), dibangkitkan
`docs/alat/skema.py` (`TIPE_DESIMAL`), ditagih `repository/kolom_test.go:155` TestKatalogSepakatDenganDDL dan
`repository/kolom_test.go:423` TestSkalaUangPersenMinimalSembilanDiSemuaTabel; penjaga inti `TestNolNumberTanpaPresisi`
menerima `NUMBER(38,10)` (satu entri `presisiSah`, commit `inti:` tersendiri — PERMINTAAN A5). **(2) F17** — `UPDATE`
kolom datar pemuat lama kini juga hanya menyentuh generasi terbuka (`repository/lama_test.go:82`
TestUbahGeneralPolisHanyaGenerasiTerbuka). `T_GENERAL_POLIS_TREATY` tetap tabel Treaty sendiri lewat `CREATE TABLE` di
migrasi 320 (perintah WO 04-10-2026; keputusan "tabel bersama FacIn" dibatalkan) — tepat delapan `CREATE TABLE`.

| # | Migrasi | Tabel | Kolom | Sheet `NB Treaty In Prop` | Sheet `NB Treaty In NonProp` | Relasi | PERBANDINGAN |
| ---: | --- | --- | ---: | --- | --- | --- | --- |
| 1 | `320_t_general_polis_treaty.sql` | `T_GENERAL_POLIS_TREATY` | 81 | F9–F33, ringkasan K108 | F9–F33, ringkasan K123 | 1:1 shared PK `T_WORK_POLIS`; `UNIQUE (OLD_POLIS_ID)`, `UNIQUE (NOPOLIS, PRODKE)` (indeks fungsi `CASE`) | bab 1 |
| 2 | `321_t_polis_quotation.sql` | `T_POLIS_QUOTATION` | 17 | G34, J35–J38, K109 | G34, J35–J38, K124 | 1:1 `POLIS_ID` | bab 2 |
| 3 | `322_t_polis_ceding.sql` | `T_POLIS_CEDING` | 5 | O39, R40–R50, K110 | O39, R40–R50, K125 | 1:N `QUOTATION_ID` | bab 3 |
| 4 | `323_t_polis_instalment.sql` | `T_POLIS_INSTALMENT` | 17 | G51, J52–J56, K111 | G51, J52–J55, K126 | 1:N `POLIS_ID` | bab 4 |
| 5 | `324_t_polis_instalment_detail.sql` | `T_POLIS_INSTALMENT_DETAIL` | 14 | — | O56, R57–R61, K127 | 1:N `INSTALMENT_ID` · NonProp saja | bab 5 |
| 6 | `325_t_polis_spreading.sql` | `T_POLIS_SPREADING` | 12 | G65, J66–J69, K113 | G70, J71–J74, K129 | 1:N `POLIS_ID` | bab 6 |
| 7 | `326_t_polis_xol.sql` | `T_POLIS_XOL` | 16 | J71–J72 (nol baris) | G75, J76–J80, K130 | 1:N `POLIS_ID` · NonProp saja | bab 7 |
| 8 | `327_t_polis_xol_layer.sql` | `T_POLIS_XOL_LAYER` | 20 | J71–J72 (nol baris) | O81, R82–R87, K131 | 1:N `XOL_ID` · NonProp saja | bab 8 |

*Kolom* = definisi kolom di `CREATE TABLE` (termasuk kunci dan `NOURUT`).

**Perubahan kolom putaran 3 — satu:** `T_GENERAL_POLIS_TREATY.EDM_TYPE` (`VARCHAR2(16)`, `migrations/320_t_general_polis_treaty.sql:48`;
katalog `models/katalog.go:125`), medan `PolicyTreatyIn.EDMType` dokumen lama. Bukti XML: dibaca
prasyarat `Activity/InputPolicyTreatyInPre_Act` **langkah 10** (`[.PolicyTreatyIn.EDMType=="3"]` → lewati
`Call TreatyRealizationCheckXOLList`), terjangkau dari FlowAction `InboxPolicyTreatyIn`/`DeptHeadTreatyIn_UW`
(`models.PerluCekDaftarXOL`). RALAT rancangan §4sexies (rancangan §4.1 *penentu bentuk* sudah memuat `EDM_TYPE`);
keputusan WO F3 butir (a). 320 kini 81 kolom = 8 kunci + 73 katalog `models.TabelGeneralPolis`.

**K18 — ✅ SELESAI 05-10-2026:** perintah work owner 05-10-2026: tabel induk NB Treaty In diganti nama `T_GENERAL_POLIS_TREATY` (migrasi `320_t_general_polis_treaty`); `T_GENERAL_POLIS` tetap milik `nbfacin` dan tidak dipakai modul ini. Dikunci `migrasi_test.go` TestMigrasiMemakaiTGeneralPolisTreatyBukanTabelFacIn (menggantikan TestPraTerbangMenolakTGeneralPolisBentukFacIn). *Riwayat* — tabrakan nama `T_GENERAL_POLIS` (PERMINTAAN C10), penahan migrasi: `POOLDATA.T_GENERAL_POLIS` sudah ada
(FacIn, migrasi `182_t_general_polis` di luar repo, 7 kolom). Modul tidak mengganti nama dan tidak menyentuh tabel itu.
Penjaganya pra-terbang inti `praTerbangBentuk` (`inti/backend/migrasi/migrasi.go:353`): `-migrate`
berhenti sebelum satu pernyataan pun dikirim dengan galat *"… sudah ada … tetapi BENTUKNYA BERBEDA …"*. Dikunci tanpa
Oracle di `backend/migrasi_test.go`: `migrasi_test.go:125` TestMigrasi320TeruraiPraTerbangPersisKolomGeneralPolis,
`migrasi_test.go:150` TestPraTerbangMenolakTGeneralPolisBentukFacIn (5 kolom lebih, 79 kurang),
`migrasi_test.go:181` TestMigrasi320TanpaBlokPLSQL (koreksi WO: tanpa PL/SQL). Temuan untuk tim inti: `KolomCreateTable`
membaca satu kolom per baris (C10).

**Migrasi yang dijalankan WO** (agen tidak menjalankan apa pun): 320–327 sesudah C10 diputuskan; 968 (slot menu)
**sudah tercatat** di `POOLDATA` 03-10-2026 oleh pihak lain (**K13**, belum diketahui).

Tabel lama yang disentuh — **tidak dibuat, tidak diubah strukturnya**:

| Tabel | Akses | Dasar diagram / keputusan |
| --- | --- | --- |
| `T_WORK_POLIS` (+ `SEQ_WORK_POLIS`) | tulis/baca | akar B5, milik premiumlistlife (migrasi 050/057/059) |
| `HISTORYAKSEPTASIPEGA` | tulis/baca | Prop F98–F99 / NonProp F113: `InsertHistoryAkseptasiPega_Sql` |
| `HISTORYAKSEPTASIPRODUCTION` | tulis/baca | Prop J74–J76 / NonProp J89–J92: `SaveViewSuggest → InsertViewSuggest_SQL` (K4); + salinan `SuggestList` dokumen lama (F3) |
| `GENERATE_SEQUENCE_NUMBER` | tulis lewat `inti/backend/penomor` | diagram "dibaca saja" — F8 disetujui WO (bab 5 butir 24) |
| `TANGGAL_CLOSING`, `KODE_PRODUKSI`, `CURRENCY`, `BUSINESS`, `REINSURANCETYPE`, `TREATYGROUP`, `MARKETINGOFFICER`, `CLIENT`, `AGENT`, `M_LOGIN_GO`, view `TREATYINDETAILJOINEDM`, `TREATYINPRODUCTION` | baca | RD/RDB terjangkau; `TREATYINPRODUCTION` ditulis modul EDM. *(Putaran 2 juga menyebut tabel `TREATYINDETAIL` — gugur: grid popup lama S11 tidak dibangun, R5.)* |
| `M_TREATY_IN`, `M_TREATY_IN_EDM` | baca (JSON, satu fungsi) | K8, F1 |
| `JSON_POLIS` | baca (pemuat dokumen lama saja) | tiket 22 |

## 9 · Status tiket 00–23

Aturan penandaan: **selesai** bila setiap AC tiketnya ✅, 📄, ⛔ dengan alasan terbukti, atau tertahan **hanya** oleh
pihak luar (K11, K12, K13, F7 WO/DBA, tim inti). Target putaran 3: semua yang dapat dikerjakan tanpa pihak luar
selesai.

| Tiket | Status | AC belum ✅ dan penahannya |
| --- | --- | --- |
| 00 | selesai | 0 AC; butir DBA (PERMINTAAN C7, C8) |
| 01 | selesai | 17, 89 — K11; 57 ⛔ (b) K8 butir 4; F1 diputuskan; popup W1 dibangun |
| 02 | selesai | 31 — K11 |
| 03 | selesai | — |
| 04 | selesai | — (klausa `pyWorkGroup` K12, PERMINTAAN C6) |
| 05 | **needs-info** | 12 — **K12** (IAM) |
| 06 | selesai | 66 ⛔ (a) |
| 07 | selesai | 23 — K11 |
| 08 | selesai | 29, 83 — K11 |
| 09 | selesai | — |
| 10 | selesai | 72 — K11; F2 diputuskan; DIV (B4); C8 |
| 11 | selesai | 45 — K12 (`ProductionDate`), K7 (b) (`DateofSurvey`); W3–W6, F4 dibangun |
| 12 | selesai | — (W2, 7.4 dibangun) |
| 13 | selesai | — (sambungan Arasapas C1) |
| 14 | selesai | — |
| 15 | selesai | 68 — F7 (+K11); 70, 93–95 📄 |
| 16 | selesai | peny. 1, 6 — K11 |
| 17 | selesai | peny. 9 — K11 |
| 18 | selesai | peny. 12, 13 — K11 |
| 19 | selesai | — (F2, B4, C8; `PerluCekDaftarXOL` berdasar AC 33) |
| 20 | selesai | peny. 45, 46 — K11 |
| 21 | selesai | peny. 49–51 — K11 |
| 22 | selesai | peny. 55 — K11; peny. 59 — F7 (uji-kering K11); F3, F6, F7 diputuskan; C5, C7 |
| 23 | selesai | peny. 63 📄 |

Hasil: **23 tiket selesai**, **05 needs-info** (K12). **Pekerjaan yang masih dapat dikerjakan tanpa pihak luar: tidak
ada.** Sisa penahan: K11 (skema uji, C9; `make test-db` A4), K12 (IAM), K13 (968 di `POOLDATA`), K18/C10 (tabrakan
`T_GENERAL_POLIS_TREATY`), F7 (pemuatan produksi oleh WO/DBA), konfirmasi WO butir H PERMINTAAN.

## 10 · Perintah verifikasi

| Perintah | Hasil (P3K, worktree `wt-nbtr-p3r-hasil`, sesudah sunting dokumen) |
| --- | --- |
| `go test ./inti/...` | 14 paket `ok`; gagal **hanya** `inti/backend/penjaga` `TestNolAlamatLayananDiKode` — butuh berkas `.env` lokal yang tidak ada di worktree (lingkungan, sama dengan putaran 2); penjaga `MODUL.md` (`modulmd`, `rentang`, `strukturkolom`, `migrasi`) lulus |
| `go test ./modul/nbtreatyin/...` | lulus — `backend`, `handlers`, `models`, `repository`, `tiruan` `ok` (termasuk `TestMigrasi320*`, `TestPraTerbangMenolakTGeneralPolisBentukFacIn`, `TestTabelDanKolomMengikutiDiagramGrilling`) |
| `go test -tags=db ./modul/nbtreatyin/...` | **tidak dijalankan** — K11 kosong |

Paket putaran 3 masing-masing melaporkan `go vet ./...`, `go vet -tags db ./modul/nbtreatyin/...`, `go test ./...`
(merah hanya baseline claimlife + `TestNolAlamatLayananDiKode` tanpa `.env`), `npm run typecheck` dan vitest modul
hijau di worktree-nya. Verifikasi akhir empat perintah dari akar repo atas tip integrasi: dirangkum orkestrator.
