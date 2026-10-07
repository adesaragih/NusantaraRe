# Kepatuhan — sheet *NB Treaty In Prop* lawan implementasi `nbtreatyin`

*Disusun 04-10-2026 (paket R9 `p3r-sheetprop`), atas perintah work owner: "Baca dan pelajari sheet NB Treaty In Prop
pada `Diagram-Skema-Tabel-NusantaraRe.xlsx`, buat sesuai yang di sheet NB Treaty In Prop."* Isi sheet dibaca utuh
(koordinat, sel gabungan, warna; berkas Excel tidak disunting). Setiap pernyataan sheet dipetakan ke kode, migrasi, atau
uji yang memenuhinya.

**Status:** **SESUAI** — dipenuhi, bukti disebut · **SESUAI (n/b)** — pernyataan tentang jalur EDM atau penulis lain;
NB tidak menulisnya, dan kode membuktikan NB tidak menyentuhnya · **SELISIH → diperbaiki** — berbeda sebelum paket ini,
diperbaiki dengan uji merah-lalu-hijau di paket ini.

Jalur relatif terhadap `modul/nbtreatyin/backend/` kecuali diawali `docs/` atau `inti/`. Uji bertag `db` ditulis dan
dikompilasi (`go vet -tags db`), **belum dijalankan** (K11).

⛔ Keputusan WO terbaru yang dipegang dokumen ini: keputusan sementara *"`T_GENERAL_POLIS` tabel bersama FacIn — 320
menjadi `ALTER`"* **dibatalkan** (commit `ea90fec6`..`b81e1ec1`). `T_GENERAL_POLIS_TREATY` tetap tabel Treaty In sendiri lewat
`CREATE TABLE` di migrasi 320 — tepat delapan `CREATE TABLE`. Tabrakan nama dengan migrasi `nbfacin` 182 (PERMINTAAN C10)
selesai 05-10-2026: tabel Treaty In bernama `T_GENERAL_POLIS_TREATY`.

## Ringkasan

| | Jumlah |
| --- | ---: |
| Baris kepatuhan (pernyataan sheet B5–B122, dikelompokkan per kesatuan makna) | **59** |
| SESUAI (termasuk F103 dengan RALAT F8 yang disetujui WO) | **51** |
| SESUAI (n/b) — jalur EDM / penulis lain | **5** |
| SELISIH → diperbaiki | **3** (F17, F20, J69) |
| SELISIH tidak diperbaiki | **0** |

Sisa yang dicatat, bukan selisih terhadap sheet: nilai berdesimal **lebih dari sepuluh** (`ShareValue` 24 dan
`PremiumSpreaded` 20 desimal di dokumen lama, rancangan §4q5.5) tetap dibulatkan Oracle di desimal kesebelas. Sheet F20
menuntut skala *minimal* 9; skala tetap yang memuat 24 desimal menekan sisi kiri koma di bawah 17 digit sentinel
`99999999999999999.99` (spec-penyimpanan bab *Presisi*).

## 1 · Akar — `T_WORK_POLIS` (B5–B7)

| Sel | Pernyataan sheet | Pemenuhan | Status |
| --- | --- | --- | --- |
| B5 | `T_WORK_POLIS` akar · lintas-lini · 1 baris per work object | tabel milik premiumlistlife (migrasi 050/059), **tidak dibuat** modul ini; satu baris per kasus `repository/kasus.go` `SisipKasus` (`sqlSisipKerja`, `LINI` = `models.LiniKasus`) | SESUAI |
| B6 | diambil dari `pyWorkPage.pzInsKey` | `T_WORK_POLIS.ID` = pengenal kasus `NB-<n>` (`IDKasusBerikut`); padanan `pzInsKey` = `models.KunciInstans` (`models/kasus.go:108`) dipakai `IDPEGA` riwayat; kasus lama: `IDPEGA` = `pzInsKey` dokumen lama (`SetelKolomDatarLama`) | SESUAI |
| B7 | tabel sama dengan akar sheet PremiumList; baris NB dan EDM sejajar, tidak saling menunjuk | nol kunci tamu NB↔EDM di `T_WORK_POLIS` (spec-penyimpanan AC 7) | SESUAI |

## 2 · `T_GENERAL_POLIS_TREATY` (F9–F33)

| Sel | Pernyataan sheet | Pemenuhan | Status |
| --- | --- | --- | --- |
| F9 | 1:1 · SHARED PK · satu baris per GENERASI | `migrations/320_t_general_polis_treaty.sql:111-112` PK `ID` + FK `ID` → `T_WORK_POLIS (ID)`; `CREATE TABLE` (bukan `ALTER`), `migrasi_test.go` `TestMigrasi320TanpaBlokPLSQL`, `repository/kolom_test.go:313` TestTabelDanKolomMengikutiDiagramGrilling | SESUAI |
| F10 | `PolicyTreatyIn` — 79 medan skalar tingkat atas | 70 kolom katalog medan `PolicyTreatyIn` + `NOPOLIS` + 9 medan tak berkolom berbukti (4 `LAYER*` F26, 4 `Total*` turunan, `isApprovedtoDeptHead` P36) ≥ 79 (batas bawah, B119) — `docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 1e; kolom di luar 79 hanya dengan RALAT berbukti XML (bab 1c) | SESUAI |
| F11 | `json_polis` — 7 kolom yang sudah datar | `NOPOLIS`, `PRODKE`, `NOENDORS`, `IDPEGA`, `TGL_INPUT`, `USERNAME`, `TGL_PROD` (320; PERBANDINGAN bab 1a) | SESUAI |
| F12 | kunci alami `(NOPOLIS, PRODKE)` · NB = PRODKE 0 | indeks unik `320:112` (bentuk `CASE` supaya draf tanpa nomor tidak bentrok); `PRODKE NUMBER(10) DEFAULT 0 NOT NULL` (`320:27`); `repository/kasus.go:39` `sqlSisipGenerasi` menulis `PRODKE` 0; daftar portal `g.PRODKE = 0`; uji db `repository/polis_db_test.go:170` TestNomorPolisSekaliDanUnik (K11) | SESUAI |
| F13 | `OLD_POLIS_ID` → `T_WORK_POLIS.ID` · nullable · penunjuk generasi sebelumnya | `320:29` (nullable) dan `320:108` FK ke `T_WORK_POLIS (ID)` | SESUAI |
| F14–F15 | pengganti `OldData`; sistem lama mencari `order by PRODKE desc` + `substr(nopolis,1,24)` | nol `SUBSTR` atas nomor polis di modul; `OldData` dokumen lama tidak dimigrasi (`models/medan_abaikan_lama.json`); penunjuk = FK | SESUAI |
| F16 | `UNIQUE (OLD_POLIS_ID)` · `UNIQUE (NOPOLIS, PRODKE)` | `320:109` `UQ_GENERAL_POLIS_OLD`; `320:112` `UQ_GENERAL_POLIS_NOPOLIS` | SESUAI |
| F17 | ⛔ baris generasi lampau TIDAK BOLEH disunting (P58) | `repository/polis.go:135` `syaratTerbuka` di `tulisInduk`, `SetelNomorPolis`, `sqlPindahGenerasi`; ⚠️ **sebelum paket ini** `UPDATE` kolom datar pemuat lama (`repository/lama.go` `sqlSetelKolomDatarLama`) **tanpa** syarat itu → kini `UPDATE … g … WHERE g.ID = :5 AND NOT EXISTS (… s.OLD_POLIS_ID = g.ID)`; uji `repository/lama_test.go:82` TestUbahGeneralPolisHanyaGenerasiTerbuka (merah lalu hijau); db `repository/polis_db_test.go:204` TestGenerasiTertutupDitolakDanPembatalanUtuh (K11) | SELISIH → diperbaiki |
| F18 | bentuk aktif: `PROPORTIONAL_TYPE = 'Proportional'` | `T_POLIS_QUOTATION.PROPORTIONAL_TYPE` (321); `models.JenisProporsional = "Proportional"` (`models/hitung.go:35`); bentuk dijaga `models.PeriksaBentukSimpan` | SESUAI |
| F19 | tipe kolom: dokumen lama semuanya teks; konversi SEKALI saat masuk | `repository/kolom.go` `nilaiTulis`/`ekspresiTulis` (satu-satunya konversi); pemuat lama menulis lewat `SimpanHalaman` yang sama (spec-penyimpanan AC 56, `repository/lama_test.go:50` TestPemuatTanpaJalurTulisTerpisah) | SESUAI |
| F20 | uang · persen → angka presisi tetap, **skala MINIMAL 9 desimal** (P29: galat lama diikuti apa adanya, tidak dibulatkan) | ⚠️ **sebelum paket ini** `NUMBER(38,8)` (skala 8; `592629512.880000276` tersimpan `592629512.88000028`) → kini **`NUMBER(38,10)`** di seluruh kolom uang/persen delapan tabel (migrasi 320, 323–327 dibangkitkan `docs/alat/skema.py` `TIPE_DESIMAL`); uji `repository/kolom_test.go:155` TestKatalogSepakatDenganDDL dan `repository/kolom_test.go:423` TestSkalaUangPersenMinimalSembilanDiSemuaTabel (merah lalu hijau); db `repository/pulangpergi_db_test.go:250` TestUangPresisiPenuhTanpaPembulatanRepository, `repository/lama_db_test.go:74` TestPemuatLamaMenulisLewatAntarmukaSama (galat lama utuh, K11); penjaga inti `presisiSah` + `NUMBER(38,10)` (PERMINTAAN A5); RALAT spec-penyimpanan ID-14/15, AC 19/20/20b/55, rancangan §4sx.6 | SELISIH → diperbaiki |
| F21 | tanggal → `DATE`; dua format masuk `YYYYMMDD` dan cap waktu Pega bersufiks GMT | `models/nilaipega.go:80` `BacaTanggalLama` (`polaCapWaktuGMT` `:61`); uji `models/dokumenlama_test.go:18` TestBacaTanggalLama, `:39` TestTanggalAmbiguTidakDitebak | SESUAI |
| F22 | cacah → bilangan bulat: `NOURUT` · `PRODKE` · `INSTALLMENT_NO` | `NOURUT NUMBER(5)`, `PRODKE`/`INSTALLMENT_NO NUMBER(10)`; `repository/kolom_test.go:62` TestCacahBilanganBulat | SESUAI |
| F23 | kode dan penanda tetap TEKS — `GroupPanel "006"`, `BusinessOldId "01"` nol depan utuh | golongan `kKode`/`kPenanda` → `VARCHAR2`; `models/penggolong_test.go:201` TestPenggolongMembandingkanTeksBukanBilangan; db `repository/penyimpanan_db_test.go:110` TestKodeBernolDepanUtuhDiKolom (K11) | SESUAI |
| F24 | `IsApproved` dibandingkan sebagai TEKS; `""` sah dan BERBEDA dari `"0"` | `kPenanda("IsApproved")`; `models/tangga_test.go:12` TestIsApprovedNolDitolak; `repository/kolom_test.go:93` TestNilaiBaca; db `repository/penyimpanan_db_test.go:140` TestIsApprovedKosongDanNolTetapBerbeda (K11) | SESUAI |
| F25 | teks kosong `""` di kolom angka/tanggal → NULL, bukan 0 | `nilaiTulis` (`repository/kolom.go:74`); `repository/kolom_test.go:43` TestNilaiTulisKosongJadiNULL | SESUAI |
| F26–F27 | `LAYER*` DICORET dari tabel ini (pantulan `pxResults(1)`) | nol kolom `LAYER*` di 320 (`TestTabelDanKolomMengikutiDiagramGrilling`); dibaca ulang dari view lewat `TREATY_IN_ID` | SESUAI |
| F28 | setiap tabel anak punya `NOURUT`, kunci pasangan antar generasi | `NOURUT NUMBER(5) NOT NULL` + `UNIQUE (induk, NOURUT)` di 322–327 | SESUAI |
| F29 | NB (PRODKE 0): baris boleh dihapus, `NOURUT` dinomori ulang rapat 1..n | `repository/polis.go:67` `SimpanHalaman` hapus-lalu-sisip seluruh anak, `sisipAnak` `NOURUT = i+1` (`:192`); db `repository/penyimpanan_db_test.go:85` TestNourutDinomoriUlangSesudahBarisKeduaDihapus (K11) | SESUAI |
| F30–F31 | EDM: baris tidak bisa dihapus, `NOURUT` lama terbawa; pembatalan P56 menolkan nilai | jalur EDM (`edmtreatyin`); NB tidak menulis generasi `PRODKE > 0` — pemuat menolak dokumen endorsemen (`models.ErrGenerasiEndorsemen`, `models/dokumenlama.go:115`) | SESUAI (n/b) |
| F32 | generasi n+1 wajib memuat setiap `NOURUT` generasi n — ditolak di services | NB hanya melahirkan generasi 0 (tanpa pendahulu); penegakan saat generasi n+1 ditulis milik services EDM | SESUAI (n/b) |
| F33 | di NB `OLD_POLIS_ID` selalu KOSONG ⇒ tabel proyeksi NOL BARIS | `sqlSisipGenerasi` tidak mengisi `OLD_POLIS_ID`; tabel proyeksi tidak dibuat (tepat delapan `CREATE TABLE`) | SESUAI |

## 3 · `T_POLIS_QUOTATION` → `T_POLIS_CEDING` (G34–R50)

| Sel | Pernyataan sheet | Pemenuhan | Status |
| --- | --- | --- | --- |
| G34, J35 | 1:1 · `POLIS_ID` UNIK | `321_t_polis_quotation.sql` PK `POLIS_ID` + FK `T_GENERAL_POLIS_TREATY (ID)` | SESUAI |
| J36–J37 | `QuotationData` — 10 medan skalar (`ProportionalType` … `MarketingName`) | sepuluh kolom ada (`models/katalog.go:222` `TabelQuotation`); enam kolom tambahan masing-masing RALAT berbukti rule NB terjangkau (PERBANDINGAN bab 2; bab 0 butir 12) | SESUAI |
| J38 | `GroupPanel` + `BusinessOldId` masukan penggolong `BusinessType_DeT` — 36 baris, 128 kode, bawaan UNKNOWN | `models/penggolong.go` (`TabelPenggolong` 36 baris, `JenisUsahaTakDikenal = "UNKNOWN"`); `models/penggolong_test.go:145` TestKe128KodeSamaDenganSistemLama, `:193` TestPenggolongBawaanUnknown, `:185` TestPenggolongBerhentiDiBarisPertama | SESUAI |
| O39, R40 | `T_POLIS_CEDING` 1:N `QUOTATION_ID` | `322_t_polis_ceding.sql` FK `QUOTATION_ID` → `T_POLIS_QUOTATION (POLIS_ID)`; ditagih `TestTabelDanKolomMengikutiDiagramGrilling` | SESUAI |
| R41–R43 | `CedingCoList()` — 2 medan; `NOURUT`; `CEDING_CO_ID ← .CedingCo`, `CEDING_CO_NAME ← .CedingCoName` | `models/katalog.go:255` `TabelCeding`; `UQ_POLIS_CEDING_NOURUT` | SESUAI |
| R44–R46 | `CedingCoList` ada di tingkat polis DAN di `OldData`; data guide basi; `CedingCoName` tingkat polis berakhiran `"; "` | pemuat membaca `PolicyTreatyIn.QuotationData.CedingCoList` dari dokumen (bukan data guide); `OldData` tidak dimigrasi; ekor `"; "` utuh — db `repository/lama_db_test.go:74` (`UJI-CEDING A; UJI-CEDING B; `, K11) | SESUAI |
| R47 | `CEDING_CO_NAME` dan `CEDING_CO` di `T_GENERAL_POLIS_TREATY` DISALIN APA ADANYA, tidak dirangkai ulang | kolom `VARCHAR2(4000)` di 320; nol kode perangkai dari `T_POLIS_CEDING` (PERBANDINGAN bab 1b); `models/dokumenlama_test.go:70` TestPecahDokumenProporsionalDatar | SESUAI |
| R48 | bentuk gabungan dan tabel bisa TIDAK SINKRON, tanpa penjaga | nol penjaga sinkron, sama dengan sistem lama | SESUAI |
| R49 | `CEDING_CO` tingkat polis daftar yang digabung | `CEDING_CO` golongan kode `VARCHAR2(4000)`; `TREATYINPRODUCTION.CEDINGCOID` ditulis modul EDM (F87) | SESUAI |
| R50 | penghapusan ceding `@replaceAll` teks gabungan TIDAK ditiru; cukup hapus barisnya | nol `@replaceAll` ceding; baris ceding dihapus lewat hapus-sisip `SimpanHalaman` | SESUAI |

## 4 · `T_POLIS_INSTALMENT`, BREAKDOWN, `T_POLIS_SPREADING`, XOL (G51–J72)

| Sel | Pernyataan sheet | Pemenuhan | Status |
| --- | --- | --- | --- |
| G51, J52 | `T_POLIS_INSTALMENT` 1:N `POLIS_ID` · SATU TINGKAT | `323_t_polis_instalment.sql` FK `POLIS_ID` → `T_GENERAL_POLIS_TREATY` | SESUAI |
| J53–J54 | `ListInstallment()` — 10 medan; `NOURUT` | `models/katalog.go:263` `TabelAngsuran` (10 medan + 4 RALAT `PPN`/`PPh`/`PaymentTotalAfterPPN`/`Tax` dibaca preACT 18.3.4.1, PERBANDINGAN bab 4) | SESUAI |
| J55 | memegang TOTAL — `PaymentTotal` | kolom `PAYMENT_TOTAL` | SESUAI |
| J56 | proporsional BERHENTI di sini — tanpa sarang `InstallmentList` | `models.PeriksaBentukSimpan` (`models/katalog.go:443`) menolak rincian angsuran polis `Proportional` (dipanggil `SimpanHalaman`); `models/layar_test.go:12` TestBentukProporsionalMenolakXOLDanRincian | SESUAI |
| G57, J58–J64 | ⛔ `T_POLIS_BREAKDOWN_SPREAD` DIBATALKAN 23-09-2026 | tidak dibuat; `BreakDownSpreading_Act` (CountSpreading_Act langkah 6) tidak dibangun (`models/angsuran.go` komentar `CountSpreading`; KEPUTUSAN-RONDE-12 butir 3/3b); `BreakDownSpreadList` dokumen lama dibuang berbukti | SESUAI |
| G65, J66–J68 | `T_POLIS_SPREADING` 1:N — `SpreadingRiskList()` 9 medan; `NOURUT` | `325_t_polis_spreading.sql`; `models/katalog.go:303` `TabelSpreading` (9 medan) | SESUAI |
| J69 | P60 — NB: `100 / jumlah baris` presisi 10 · ditiru apa adanya | hitungan: `models/angsuran.go:31` `presisiBagiRata = 10`, `models/layar_test.go:34` TestSpreadingBagiRataPresisiSepuluh (`33.3333333333`); ⚠️ **sebelum paket ini** kolom `SHARE_PERCENTAGE`/`CLAIM_PERCENTAGE` `NUMBER(38,8)` membulatkannya ke `33.33333333` saat disimpan → kini `NUMBER(38,10)` (bersama F20); db `repository/pulangpergi_db_test.go:344` TestSpreadingBagiRataPresisiSepuluhUtuhDiKolom (K11) | SELISIH → diperbaiki |
| G70, J71–J72 | `T_POLIS_XOL` · `T_POLIS_XOL_LAYER` NOL BARIS pada proporsional — tidak ada barisnya | `models.PeriksaBentukSimpan` menolak `TreatyXOLList` polis `Proportional`; `TestBentukProporsionalMenolakXOLDanRincian` | SESUAI |

## 5 · `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (G73–J78)

| Sel | Pernyataan sheet | Pemenuhan | Status |
| --- | --- | --- | --- |
| G73, J74 | 1:N `IDPEGA` · 15 kolom · SUDAH datar | tabel warisan, tidak dibuat (`MODUL.md` *Tabel warisan*); `repository/usulan.go` `sqlSisipUsulan` 15 kolom; `repository/kolom_test.go:246` TestSQLRiwayatProduksiMengikutiInsertViewSuggest | SESUAI |
| J75–J76 | diambil dari `PolicyTreatyIn.SuggestList`; ditulis `SaveViewSuggest → InsertViewSuggest_SQL` | `models/usulan.go` (pemetaan `SaveViewSuggest`), `repository/usulan.go` `CatatUsulan` (pengganti `InsertViewSuggest_SQL`); `SuggestList` dokumen lama disalin (F3, `SalinUsulanLama`) | SESUAI |
| J77 | `.Suggest → KETERANGAN (substr 0,3990)` · `.IsApproved → APPROVAL ("1"=Accept, "0"=Reject)` per baris | `SUBSTR(:11, 1, models.PanjangKeterangan)`; `models/usulan_test.go:66` TestApprovalPerBarisHanyaAcceptReject; db `repository/polis_db_test.go:268` TestRiwayatProduksiPulangPergi (K11) | SESUAI |
| J78 | `AKSES_LOGIN ← OperatorID.pyUserIdentifier` (P4) · `PIC ←` nama tampilan (P33) | `models/usulan_test.go:12` TestUsulanBelumTersimpanMenurutSaveViewSuggest | SESUAI |

## 6 · Tabel proyeksi (B80–J83)

| Sel | Pernyataan sheet | Pemenuhan | Status |
| --- | --- | --- | --- |
| B80, G80, J81–J83 | `T_POLIS_DIFFERENCE` (dan dua proyeksi lain) NOL BARIS di NB — tidak ada generasi sebelumnya | tidak dibuat dan tidak ditulis modul NB (nol rujukan `T_POLIS_DIFFERENCE*`); tepat delapan `CREATE TABLE` | SESUAI |

## 7 · Tabel lama yang sudah datar (F86–F103)

| Sel | Pernyataan sheet | Pemenuhan | Status |
| --- | --- | --- | --- |
| F86–F89 | `TREATYINPRODUCTION` ditulis `InsertTreatyInProdEDMT_SQL` (modul EDM); CARIn diisi `InsetTreatyInProdAddendum_Act` | NB **tidak menulis** — satu-satunya rujukan `repository/acuan.go:485` `SELECT` cek duplikat (`TreatyRealizationCheckDuplicate`) | SESUAI |
| F90–F95 | 22 kolom uang menyeberang teks; `PCT_SHARE_*` tertukar; empat titik sisip; `DEDUCTION1` dua satuan; `LAYER*` "0" | catatan penulis EDM; NB tidak menulis tabel itu, dan pembacaan cek duplikat NB tidak memakai kolom uang/`DEDUCTION1`/`LAYER*` | SESUAI (n/b) |
| F96 | ditulis TANPA awalan skema — pelanggaran ketetapan 6 | setiap nama tabel modul lewat `db.Qualify` (spec-penyimpanan AC 47) | SESUAI |
| F98–F99 | `HISTORYAKSEPTASIPEGA` · `JSON_POLIS_MONITORING` … ditulis `InsertHistoryAkseptasiPega_Sql` · `INSERTJSONPOLISMONITORING` | `HISTORYAKSEPTASIPEGA` ditulis `repository/riwayat.go` `CatatRiwayat` (pengganti `InsertHistoryAkseptasiPega_Sql`, transaksi submit); `JSON_POLIS_MONITORING`: nol rule korpus NB menulisnya — NB tidak menulis | SESUAI |
| F100 | `procedure PEGA_TREATY_IN` · `procedure InsertUpdateAchievment` | `RDBList/SaveTreatyIn` (`PEGA_TREATY_IN`) tidak terjangkau — `INVENTARIS-XML.md` baris 161 (prasyarat `revisionstate==1` tak pernah dikirim); `InsertUpdateAchievment`: nol rule korpus NB — NB tidak memanggil procedure apa pun (spec-penyimpanan AC 48) | SESUAI (n/b) |
| F101 | `ACHIEVEMENT`: penjaga duplikat 17 kolom termasuk uang TIDAK ditiru | NB tidak menulis `ACHIEVEMENT` | SESUAI (n/b) |
| F102 | `TREATY_IN` memuat jadwal angsuran master BERTINGKAT DUA | dibaca baca-saja hanya jalur NonProp (`repository.MasterXOLDariJSON`, F1); polis Prop tetap satu tingkat (J56) | SESUAI |
| F103 | dibaca saja: `TANGGAL_CLOSING` · `GENERATE_SEQUENCE_NUMBER` · `CURRENCY` · `BUSINESS` · `REINSURANCETYPE` · `TREATYGROUP` · `M_TREATY_IN_EDM` · ⚠️ `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | semuanya dibaca saja (`repository/gudang.go:32-44`, `repository/riwayat.go` `HariClosing`), **kecuali** `GENERATE_SEQUENCE_NUMBER` ditulis lewat `inti/backend/penomor` — RALAT catatan diagram **F8** (disetujui WO, `rancangan-tabel-datar-treaty-in.md` §4bis.4); `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` tidak dibaca — `RDBList/GetCountClaim` kategori (a) (`docs/alat/status.json`) | SESUAI (RALAT F8) |

## 8 · Ringkasan dan catatan sheet (B105–B122)

| Sel | Pernyataan sheet | Pemenuhan | Status |
| --- | --- | --- | --- |
| B105–B117 | 8 tabel hidup NB Prop: `T_WORK_POLIS`, `T_GENERAL_POLIS_TREATY`, `T_POLIS_QUOTATION`, `T_POLIS_CEDING`, `T_POLIS_INSTALMENT`, (`T_POLIS_BREAKDOWN_SPREAD` dibatalkan), `T_POLIS_SPREADING`, `HISTORYAKSEPTASIPRODUCTION`; 3 proyeksi nol baris | lima tabel dibuat modul (320–323, 325) + akar milik premiumlistlife + tabel warisan ditulis; `T_POLIS_INSTALMENT_DETAIL`/`XOL`/`XOL_LAYER` (324, 326, 327) milik sheet NonProp, nol baris di Prop; proyeksi tidak dibuat. Jumlah `CREATE TABLE` modul tepat delapan (bab 0 butir 11) | SESUAI |
| B119–B120 | angka medan = BATAS BAWAH; `SuggestList`, `CedingCoList`, lima property EDM terlewat sapuan | `SuggestList` → `HISTORYAKSEPTASIPRODUCTION`, `CedingCoList` → `T_POLIS_CEDING`; kolom di atas batas bawah hanya lewat RALAT berbukti XML (PERBANDINGAN bab 1c, 2, 4) | SESUAI |
| B122 | sumber sheet | rujukan, tanpa tuntutan | SESUAI |

## 9 · Yang diubah paket ini

| Selisih | Commit | Uji (merah → hijau) |
| --- | --- | --- |
| F20 + J69 — `NUMBER(38,8)` → `NUMBER(38,10)` (8 tabel, `skema.py`, STRUKTUR, katalog, RALAT dokumen) | `nbtreatyin: uang dan persen NUMBER(38,10) …` | `TestKatalogSepakatDenganDDL`, `TestSkalaUangPersenMinimalSembilanDiSemuaTabel`; db: `TestUangPresisiPenuhTanpaPembulatanRepository`, `TestSpreadingBagiRataPresisiSepuluhUtuhDiKolom`, `TestPemuatLamaMenulisLewatAntarmukaSama` |
| F20 — penjaga inti `presisiSah` menerima `NUMBER(38,10)` | `inti: penjaga presisi NUMBER …` | `inti/backend/penjaga` `TestNolNumberTanpaPresisi` |
| F17 — `UPDATE` pemuat lama hanya generasi terbuka | `nbtreatyin: pemuat lama hanya menyunting generasi terbuka …` | `TestUbahGeneralPolisHanyaGenerasiTerbuka` |
