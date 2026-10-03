# Hasil implementasi — modul `nbtreatyin` (NB Treaty In)

> Cabang `modul/nbtreatyin/implementasi` dari `dev` 51d9b9cc · 2026-10-03.
> Status: ✅ lulus (ada bukti uji/kode) · 🟡 sebagian · ⛔ tidak dibangun (alasan) · 📄 butir dokumen ·
> 🔒 `[terbuka]` (menunggu pihak lain). "Bukti" menunjuk uji (`*_test.go`, `*.test.ts`) atau kode.
> Uji `-tags db` (Oracle sungguhan) **belum pernah dijalankan** — env skema uji tidak tersedia; AC yang
> buktinya hanya uji db ditandai 🟡.

## 1 · Ringkasan

Dihitung dari tabel bab 2-3 (rentang dijabarkan; tiap nomor tepat satu kali).

| | ✅ | 🟡 | ⛔ | 🔒 | 📄 | jumlah |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| spec.md | 75 | 12 | 5 | 0 | 4 | **96** |
| spec-penyimpanan-relasional | 39 | 14 | 11 | 1 | 1 | **66** |

⚠️ `issues/00-PETA-AC-PENYIMPANAN.md` menyebut **67** AC NB di satu tempat dan **66** di tempat lain; bab 6
spec penyimpanan memuat **66** butir bernomor — angka itulah yang dicacah di sini.

## 2 · spec.md — 96 acceptance criteria

| AC | Status | Bukti / alasan |
| ---: | :---: | --- |
| 1 | ✅ | `models.Disetujui`; `tangga_test.go` TestTanggaAdminMenolakSelesaiDitolak, handlers/alur_test TestAdminMenolakDiselesaikanDitolak |
| 2 | ✅ | `tangga_test.go` (kosong = disetujui) |
| 3 | ✅ | pembandingan teks `isApproved != "0"`; uji " 0"/"0.0" |
| 4 | ✅ | DecisionTable, bukan When (komentar `tangga.go`; uji kosong) |
| 5 | ✅ | TestAdminMenolakDiselesaikanDitolak (Resolved-Rejected) |
| 6 | ✅ | TestTanggaAtasanMenolakKembaliKeAdmin; handlers/alur_test TestTanggaPenuhDanNomorPolisSekali |
| 7 | ✅ | TestTanggaAdminMenyetujuiNaikKeSecHead |
| 8 | ✅ | TestTanggaSecHeadMenyetujuiSelaluNaikKeDeptHead — ⚠️ WO diikuti, bertentangan dengan `CekLimitTreatyAcc_Act` (RALAT tiket 03) |
| 9 | ✅ | TestTanggaDeptHeadMenyetujuiSelesai; handlers/alur_test alur penuh |
| 10 | ✅ | TestPosisiBuanganTidakAda |
| 11 | ✅ | antrean = workbasket; TestBukanAnggotaAntreanDitolak |
| 12 | 🟡 | wewenang = keanggotaan antrean (peran); 12 tempat identitas → tiket 05 (`M_NBTRIN_PERAN_TEMPAT` kosong) |
| 13 | ✅ | peran di `inti.Pelaku.Peran` / tabel 330; nol kolom telepon |
| 14 | ✅ | `anggota()` menurut nama; uji urutan peran terbalik |
| 15 | ✅ | `repository.DetailKontrak` membaca view TREATYINDETAILJOINEDM; nol JSONDATA |
| 16 | ✅ | nol penulisan JSON (SaveJsonPolis diganti relasional) |
| 17 | 🟡 | 33 kolom RD dibaca dari view (`KolomRDDetail`); belum diuji lawan Oracle |
| 18 | ✅ | medan uang dibaca apa adanya (TO_CHAR TM9), tanpa hitung ulang saat buka |
| 19-21 | ✅ | `penggolong_test.go` (128 nilai, first-match, UNKNOWN) |
| 22 | ✅ | syarat When yang dijalankan (INVENTARIS bab 11) |
| 23 | 🟡 | NUMBER(38,8) presisi penuh; pulang-pergi hanya di uji db |
| 24 | ✅ | nol pembulatan di repository (`pecahAngka` eksak) |
| 25 | ✅ | apd di seluruh paket; penjaga repo |
| 26 | 🟡 | DEDUCTION1/2 bergolongan persen (WO P29) — rumus XML memakainya sebagai jumlah, diport apa adanya (pertentangan dicatat, tiket 07); RNM_SHARE/BROKERAGE view tidak dipakai (`[terbuka]`) |
| 27-28 | ✅ | hitung_test TestPajakBrokerage* |
| 29 | 🟡 | urutan layanan: handlers TestSatuTransaksiPembatalanUtuh (gudang tiruan); Oracle: polis_db_test (belum dijalankan) |
| 30 | ✅ | `db.Qualify` di setiap SQL; penjaga TestNolNamaTabelTelanjangDiQuery |
| 31 | 🟡 | indeks unik (NOPOLIS, PRODKE) + penomor FOR UPDATE; uji db saja |
| 32-33 | ✅ | kolom DATE; satu format `YYYY-MM-DD[ HH:MI:SS]` (kolom_test) |
| 34-35 | ✅ | `PraprosesAdmin` (hari ini) |
| 36-38 | ✅ | handlers TestPilihBisnis (422, nol simpan) |
| 39-41 | ✅ | riwayat OPERATORID login, USERNAME nama tampilan (handlers/alur_test TestAdminMenolakDiselesaikanDitolak) |
| 42 | ✅ | `PraprosesAtasan` nama tampilan (P33); `NamaTampilan` tanpa jatuh-balik ke ID login |
| 43 | ✅ | satu baris riwayat per submit (5 submit → 5 baris) |
| 44 | ✅ | NBStatus dari data (nama tampilan / nama posisi) — RALAT tiket 10 |
| 45 | 🟡 | 25 medan di 2 layar dibangun + ProductionDate (tiket 05) + DateofSurvey (layar survei tidak dibangun) = 27 |
| 46-47 | ✅ | layar_test TestMedanWajibAtasanTanpaEnamMedanAdmin |
| 48 | ✅ | TestMedanWajibMenahanKirimDanSimpan — Submit dan Save (WO) |
| 49-51 | ✅ | `GabungMasukanLayar` (daftar izin) + turunan dihitung ulang server; TestMedanTerkunciAtasanDanTurunanAdmin |
| 52 | 🟡 | 6 medan dapat diisi atasan (DueTo, FlagPPH, NoOfferSlip, IsApproved, Suggest, ProductionDate) — spec menyebut 7 |
| 53 | ✅ | elemen `1=2`/`NEVER` tidak dibangun (medan.test.ts) |
| 54 | ✅ | `DaftarMataUang` `CURRENCY <> 'ITL'` |
| 55-56 | ✅ | kode tampil apa adanya; medan berpilihan "associated" = isian teks bebas |
| 57 | ⛔ | treaty keluar: seluruh bagiannya JSON master (P29) — tidak dapat dibaca dari view |
| 58 | ✅ | nol penulisan treaty keluar |
| 59 | ⛔ | RALAT: `CheckDuplicateOffer` langkah 1-4 `//`, parameter kosong — peringatan tidak pernah menyala di XML |
| 60 | ✅ | penjaga TestHandlersTidakMengimporRepository |
| 61-63 | ✅ | INVENTARIS bab 2 status; nol pemanggil |
| 64-65 | ✅ | kolom tidak dibuat (kolom_test TestKatalogSepakatDenganDDL) |
| 66 | ⛔ | `isFOR` milik SaveJsonPolis (tidak dibangun); EDM di luar NB |
| 67 | ✅ | `KodeLiniJiwa` L1..L16 + uji IsLife |
| 68 | ⛔ | pemuat dokumen lama (tiket 22) belum dibangun |
| 69 | ⛔ | idem (migrasi) |
| 70 | 📄 | penyimpangan dicatat di kode + bab 4 berkas ini |
| 71 | ✅ | T_POLIS_SUGGEST + OperatorID (layar_test, handlers/alur_test) |
| 72 | 🟡 | `DaftarRiwayat` ORDER BY TGL_TRANSFER; uji db belum |
| 73 | ✅ | `RakitNomorPolis` + handlers/alur_test nomor polis (bentuk) |
| 74 | ✅ | TerbitkanNomor kedua mengembalikan nomor yang sama; SetelNomorPolis `NOPOLIS IS NULL` |
| 75 | ✅ | `PraprosesAdmin` IsNewPolicyNonProp |
| 76 | ✅ | TestPascaAdminTolakNBStatus (HasFacOut) |
| 77 | ✅ | halaman tidak membawa pesan antar permintaan; validasi kirim dihitung atas salinan |
| 78 | ✅ | `SetValidateInstallment` langkah 1 tidak dibangun |
| 79 | ✅ | rumus dari PropertiesValue (komentar per langkah) |
| 80 | ✅ | layar Dept Head dapat disimpan (handlers/alur_test alur penuh) |
| 81 | ✅ | tempat tertunda tanpa baris tabel 330 (TestTempatBerperanTidakDitebak); pemetaan sendiri menunggu IAM |
| 82 | ✅ | ARAH tidak ditebak; dua arah = tertunda |
| 83 | 🟡 | GagalDi CatatRiwayat membatalkan submit (tiruan); lawan Oracle belum dijalankan |
| 84 | 🟡 | tabel keputusan: kosong = disetujui (diuji); layar: Approval wajib dan tombol hanya untuk 1/0 (XML) → submit kosong 422 — dicatat (spec RALAT) |
| 85 | ✅ | `KotakMedan` kode mata uang di samping angka (medan.test.ts) |
| 86 | ✅ | format tidak ditebak — angka apa adanya |
| 87 | ✅ | RALAT oleh KEPUTUSAN-RONDE-12 butir 7: muatan 4 medan, sesudah commit, gagal tidak membatalkan (TestKonversiSesudahSelesai); sambungan `[terbuka]` |
| 88 | ✅ | rule tak terjangkau tidak dibangun (INVENTARIS) |
| 89 | 🟡 | kolom view yang hilang = galat (TestKolomViewHilangAdalahGalat); belum diverifikasi lawan instance |
| 90 | ✅ | `HISTORYAKSEPTASIPEGA` lewat Qualify (POOLDATA) |
| 91 | ✅ | nol peran karangan; uji memakai UJI- |
| 92 | ✅ | kasus menunggu posisi |
| 93-95 | 📄 | butir dokumen spec |
| 96 | ✅ | penyamaran inventaris diperluas (nama di When, nomor polis) — ⚠️ commit 9456900e masih memuat 3 nama dan 2 nomor polis di INVENTARIS; diperbaiki 40460975 — sebelum push perlu squash/penulisan ulang riwayat (keputusan manusia) |

## 3 · spec-penyimpanan-relasional — acceptance criteria

| AC | Status | Bukti / alasan |
| ---: | :---: | --- |
| 1 | 🟡 | indeks unik fungsi CASE (320); polis_db_test |
| 2-5 | ✅ | `SisipKasus` PRODKE 0, OLD_POLIS_ID NULL+UQ, kunci bersama (DDL + kolom_test) |
| 6 | 🟡 | `tulisInduk ... TGL_TUTUP IS NULL` → ErrGenerasiTertutup; polis_db_test |
| 7 | ✅ | nol kunci tamu NB↔EDM |
| 8, 10 | ✅ | NOURUT + UQ per induk (DDL) |
| 9 | 🟡 | hapus-sisip menomori ulang 1..n; polis_db_test |
| 11 | ✅ | NB: tidak ada pemasangan antar generasi |
| 12-13 | 🟡 | kode VARCHAR2 apa adanya (kolom_test); pulang-pergi di uji db |
| 14 | ✅ | penggolong_test 006/01 → FireStyle2 |
| 15 | 🟡 | Oracle menyimpan `''` sebagai NULL; dibaca kembali `""` |
| 16 | ✅ | `Disetujui` |
| 17-18 | ✅ | kolom_test TestNilaiTulisKosongJadiNULL |
| 19-20 | ✅ | NUMBER(38,8) (Oracle membulatkan pada skala) |
| 21-22 | ⛔ | format dokumen lama (`YYYYMMDD`, ` GMT`) — pemuat tiket 22 belum dibangun |
| 23-24 | ✅ | DATE; nol FLOAT (penjaga) |
| 25 | 🟡 | perbandingan uang di port rumus ditiru dari Pega (persis) |
| 26 | ✅ | tanpa kolom LAYER di T_GENERAL_POLIS |
| 27 | ✅ | T_POLIS_QUOTATION 26 kolom, termasuk OLD_POLICY_NO dan GROUP_PANEL |
| 28 | ✅ | T_POLIS_CEDING CEDING_CO + CEDING_CO_NAME |
| 29-30 | ✅ | disalin apa adanya; baris dihapus utuh |
| 31, 33 | ✅ | `PeriksaBentukSimpan` (layar_test TestBentukProporsionalMenolakXOLDanRincian) |
| 32 | ✅ | akibat AC 33 |
| 34-35 | ✅ | LAYER di T_POLIS_XOL_LAYER; DEDUCTION XOL uang |
| 36 | ✅ | TestSpreadingBagiRataPresisiSepuluh |
| 37 | ✅ | persen NUMBER(38,8) |
| 38 | 🟡 | DEDUCTION1/2 persen (WO); `TOTAL_SHARE_PERCENTAGE_*` turunan, tidak disimpan (penjaga nama TOTAL_) |
| 39 | ✅ | HISTORYAKSEPTASIPRODUCTION tidak dibuat ulang (dan tidak ditulis: syarat BusinessFac "F") |
| 40-41 | ⛔ | hanya berlaku pada HISTORYAKSEPTASIPRODUCTION (tidak ditulis kasus treaty) |
| 42 | ✅ | T_POLIS_SUGGEST.IS_APPROVED terpisah |
| 43-44 | ✅ | OPERATOR_ID login, OPERATOR_NAME nama tampilan |
| 45-46 | 🟡 | satu transaksi `DalamTransaksi`; pembatalan diuji atas tiruan, uji db belum dijalankan |
| 47 | ✅ | penjaga skema eksplisit |
| 48 | ✅ | nol procedure (penomor menulis SQL) |
| 49-51 | 🟡 | polis_db_test pulang-pergi (belum dijalankan) |
| 52-58 | ⛔ | pemuat dokumen lama (tiket 22) belum dibangun |
| 59 | 🔒 | penampung medan tak dikenal — menunggu pemuat |
| 60-61 | ✅ | penyamaran inventaris (lihat AC 96 spec) |
| 62 | ✅ | fixture UJI- |
| 63 | 📄 | migrasi dijalankan manusia |
| 64 | 🟡 | 72 kolom katalog (69 PolicyTreatyIn + 3 halaman kerja) — medan yang hanya dipakai rule tidak terjangkau / JSON (IsEDMInputOnNB, IDNewBisnis, TOTAL_*, LAYER*, isApprovedtoDeptHead) tidak dibuat |
| 65 | ✅ | tanpa penjaga sinkron |
| 66 | ✅ | penjaga arah lapisan |

## 4 · Penyimpangan dan pertentangan yang dicatat

| # | Hal | XML | Yang dibangun | Dasar |
| ---: | --- | --- | --- | --- |
| 1 | Sec Head menyetujui | ≤ 200 juta & bernomor → selesai; > batas → Dept Head (`CekLimitTreatyAcc_Act`) | selalu Dept Head | AC 8 / P13 `[keputusan work owner]` |
| 2 | Tombol Generate atasan | identitas `<ID-operator-1>` | posisi kasus Dept Head | P13 + tiket 05 (catatan) |
| 3 | Nilai master `TreatyIn.*` | JSON master | langkah ber-master dilewati bila kosong | P29; RNM_SHARE belum boleh dipakai |
| 4 | Pra-proses tanggal produksi | `>25` tertanam | hari dari TANGGAL_CLOSING | preseden WO PremiumList + penjaga |
| 5 | NBStatus | nama orang tertanam di connector | nama posisi tujuan | P40, AC 44, 92 |
| 6 | OperatorName atasan | `pyUserIdentifier` | nama tampilan | P33, AC 42 |
| 7 | Riwayat | tanpa skema, COMMIT per RDB | skema eksplisit, satu transaksi | AC 29, 30, 83 |
| 8 | Arasapas | stub kelas Data | bentuk kelas Work, 4 medan | KEPUTUSAN-RONDE-12 butir 7 |
| 9 | `GetOldIDBusiness_SQL` | tanpa ORDER BY | ORDER BY ID | hasil tetap |
| 10 | Daftar portal | `GetListOpportunity` logika "A" (per pembuat) | semua kasus terbuka | antrean bersama AC 11 |

## 5 · Perintah verifikasi

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` | lulus |
| `go test ./...` | modul nbtreatyin lulus; gagal: 3 paket claimlife (baseline sebelum implementasi) + `TestSetiapPemanggilBukaMemeriksaBolehDilewati` (sensus 16→17, PERMINTAAN A2) |
| `npm run typecheck` | lulus |
| `npm test` | 15 gagal / 8 berkas — 14 baseline + `daftar.menuTabel.test.ts` (PERMINTAAN A1); 11 uji modul lulus |
| `make test-db` | tidak dijalankan — env skema uji tidak tersedia |
