# Hasil implementasi — modul `edmtreatyin` (EDM Treaty In)

> 06-10-2026, sesudah tinjauan kode 06-10-2026. Status setiap acceptance criteria dua spec, dinilai dari kode dan uji
> modul saat ini (`go test ./modul/edmtreatyin/...` dan `npx vitest run modul/edmtreatyin` lulus, bab 7). Peta AC →
> tiket → sumber: `docs/REKONSILIASI-AC.md`; koreksi `#n`: `docs/KOREKSI-DOKUMEN-2026-10-06.md`; nasib setiap rule:
> `docs/INVENTARIS-XML.md`.
>
> Status: ✅ terpenuhi · 🟡 sebagian (bagian yang kurang disebut) · ⛔ tidak, sadar · ⏳ menunggu WO · 🔜 menunggu
> migrasi atau langkah manual. *Bukti*: uji (`berkas:baris NamaUji`; vitest `berkas` "nama") atau kode (`berkas`
> `fungsi`) bila belum diuji. Jalur relatif `modul/edmtreatyin/`.
>
> ⛔ **Nol uji Oracle.** Uji seam HTTP berjalan di atas gudang tiruan (`backend/tiruan`: transaksi batal sungguhan,
> UNIQUE `OLD_POLIS_ID`, generasi tertutup, proyeksi `'PEGA'` beku); SQL diuji sebagai teks. SQL baca baru sudah dicoba
> baca-saja di DEV (hitungan saja); SQL tulis ke 360-363 belum dapat dicoba karena migrasinya belum dijalankan.

## Isi

1. Ringkasan
2. `spec-penyimpanan-relasional.md` — 57 AC berlaku
3. `spec.md` — 59 AC
4. Pemuat dokumen lama (tiket 09-10)
5. Penyimpangan sadar dan perubahan tinjauan kode
6. Butir menunggu WO dan langkah manual
7. Perintah verifikasi

## 1 · Ringkasan

| | ✅ | 🟡 | ⛔ | ⏳ | 🔜 | jumlah |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `spec-penyimpanan-relasional.md` | **43** | 10 | 0 | 0 | 4 | **57** |
| `spec.md` | **46** | 4 | 9 | 0 | 0 | **59** |

Cacah penyimpanan = nomor 1-58 tanpa AC 56 (ditarik 23-09 sore). Rule korpus (163, `INVENTARIS-XML.md`): **109 ✅** ·
**10 🟡** · **44 ⛔**.

Penahan setiap AC yang belum ✅:

| Penahan | penyimpanan | `spec.md` |
| --- | --- | --- |
| ⛔ **konversi Arasapas** — tetap tidak disambung (keputusan WO 07-10-2026 butir 4) | — | 37, 38, 39, 40, 41, 44, 45, 52, 53 |
| 🔜 **migrasi 360-363** belum dijalankan WO — logika lulus di tiruan, SQL tulis belum pernah berjalan | 24, 25, 26, 27 | — |
| 🟡 **penyimpangan sadar / keputusan** — generasi tolak dilepas, rumus XOL dipertahankan, salinan antarmuka NB; cacah tabel + survei/riwayat produksi | 2, 16, 53, 55 | — |
| 🟡 **cakupan sebagian** — keutuhan baris hanya spreading, NOURUT posisi, daftar kolom khas disusun ulang, pembulatan Oracle belum diuji | 8, 10, 11, 32, 49, 58 | 35, 42, 46, 57 |

## 2 · `spec-penyimpanan-relasional.md` — 57 AC berlaku

| AC | Ringkas | Status | Bukti | Catatan |
| ---: | --- | :---: | --- | --- |
| 1 | Baris endorsemen `PRODKE ≥ 1` | ✅ | `backend/handlers/alur_test.go:209` TestCreateMelahirkanGenerasiBerikut, `backend/repository/edm_sql_test.go:76` TestSisipGenerasiMembawaKunciEndorsemen, `backend/repository/edm_sql_test.go:38` TestDaftarEDMHanyaGenerasiEndorsemen | T 02 |
| 2 | `OLD_POLIS_ID` terisi, menunjuk generasi sebelumnya | 🟡 | `backend/handlers/alur_test.go:209` TestCreateMelahirkanGenerasiBerikut, `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga, `backend/handlers/alur_test.go:326` TestAdminMenolakMelepasGenerasi | Penyimpangan sadar: generasi kasus yang DITOLAK Admin dilepas (`OLD_POLIS_ID` NULL, `LepasGenerasi`) supaya polis dapat diendorse ulang; generasi berjalan dan selesai selalu terisi |
| 3 | Dua baris tidak berbagi `OLD_POLIS_ID` | ✅ | `UQ_GP_TREATY_OLD` (NB 320); tiruan `ErrGenerasiSudahDiendorse`; `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi (percabangan ditolak `SisipKasus`) | Oracle belum dicoba |
| 4 | Dua endorsemen serentak: yang kedua ditolak | ✅ | `backend/services/layanan.go` `BuatKasus` (UNIQUE `OLD_POLIS_ID` → 409) + `AdaEDMBerjalan`; `backend/handlers/alur_test.go:242` TestEDMKeduaSelamaBerjalanDitolak | Yang diuji berurutan (422 `PesanEDMBelumSelesai`); balapan sungguhan dijaga indeks unik Oracle - belum dicoba |
| 5 | Nomor `NOPOLIS/E` + sekurangnya dua digit, di `NOENDORS` | ✅ | `backend/models/edm_buat.go` `NomorEDM`; `backend/handlers/alur_test.go:209` TestCreateMelahirkanGenerasiBerikut, `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga, `backend/repository/edm_sql_test.go:76` TestSisipGenerasiMembawaKunciEndorsemen, `backend/handlers/alur_test.go:473` TestCreateDitolakBilaJSONPolisLebihMaju | #44, #25; ≥ 100 tiga digit tanpa uji khusus. Create ditolak bila json_polis lebih maju (mencegah EDMNo kembar) |
| 6 | Sunting generasi berpenerus ditolak | ✅ | `backend/repository/polis.go` `SimpanHalaman` → `tulisInduk` (`syaratTerbuka` → `ErrGenerasiTertutup`), `services.kerjakan`; `backend/repository/lama_edm_test.go:106` TestUbahGeneralPolisPemuatHanyaGenerasiTerbuka | Ditegakkan juga tiruan; uji HTTP penyuntingan generasi tertutup tidak ada |
| 7 | Selisih terhadap generasi tepat sebelumnya | ✅ | `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga, `backend/handlers/alur_test.go:438` TestSelisihDihitungUlangSaatSubmitTanpaTombolHitung | T 03; 1700 − 1500 = 200 (bukan 1700 − 500 / 1700 − 1000) |
| 8 | Generasi `n+1` kehilangan `NOURUT` generasi `n` ditolak di services | 🟡 | `backend/models/edm_bentuk.go` `BarisSpreadingHilang` dari `validasiKirim`; `backend/handlers/alur_test.go:309` TestBarisSpreadingGenerasiLamaWajibAdaSaatSubmit; pemuat `ErrKeutuhan` (`backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi) | Hanya SpreadingRiskList. Angsuran (jumlah termin dapat diubah `FillPaymentInstallment`, XML) dan XOL tidak diperiksa; polis NonProp baru tidak diperiksa (grid tanpa Add) |
| 9 | Tidak ada tabel salinan nilai lama | ✅ | `backend/migrasi_test.go:62` TestMigrasiHanyaEmpatTabelProyeksiSelisih | #13 |
| 10 | Baris rincian tidak dapat dihapus | 🟡 | `frontend/components/TabData.test.tsx` "grid spreading: Add, dropdown Type Treaty ("Choose"), %Share tersunting; TANPA Delete", `frontend/pages/LayarKasus.test.tsx` "tanpa tombol Delete di grid spreading mana pun; Type Treaty hanya-baca menampilkan teks acuan"; `BarisSpreadingHilang` | Spreading ✅; rincian angsuran dibangun ulang `FillPaymentInstallment` / `FillPaymentInstallmentEDMT` (XML) - jumlah baris dapat berubah |
| 11 | `NOURUT` lama terbawa apa adanya | 🟡 | `NOURUT` = posisi baris (katalog); `BarisSpreadingHilang` | `FillSpreading` (XML, EDMChooseBusiness_Act 9) hanya membawa baris 1; baris 2..n ditambah Admin (Add) - berpasangan benar hanya bila urutannya sama |
| 12 | Baris baru `NOURUT = maks + 1` | ✅ | Add menambah di akhir daftar (`frontend/components/TabData.tsx`); `NOURUT` = posisi |  |
| 13 | `PASANGAN_BERGESER = 1` = anomali | ✅ | `backend/models/penanda_migrasi.go` `HitungPenandaMigrasi`; `backend/models/penanda_migrasi_test.go:32` TestPenandaPasanganBergeserMenurutNourut, `backend/models/laporanlama_test.go:118` TestRingkasanEDMSelesaiHanyaBilaNolGalatDanNolBelumDiputuskan | Dilaporkan ringkasan pemuat (cacah bergeser) |
| 14 | Pembatalan = generasi baru bernilai nol | ✅ | `backend/models/edm_pilih.go` `SetEDMTCancel` (EDMChooseBusiness_Act 5); `backend/handlers/alur_test.go:380` TestPembatalanNolkanDataBaru, `backend/models/edm_pilih_test.go:240` TestSetEDMTCancelProporsional, `backend/models/edm_pilih_test.go:270` TestSetEDMTCancelNonProp, `backend/handlers/alur_test.go:438` TestSelisihDihitungUlangSaatSubmitTanpaTombolHitung | Generasi baru, nol penghapusan baris ✅. Cakupan nol = 16 medan XML + angsuran + spreading (GrossPremium, NetPremium, BalanceDueTo, PPN/PPH tidak dinolkan) - butir WO #47 · **WO 07-10-2026** butir 6 *"ikuti XML dulu"*: 16 medan `SetEDMTCancel` |
| 15 | Jenis endorsemen di kolomnya sendiri, dipilih di awal | ✅ | `EDM_TYPE` (`sqlSisipGenerasi`); `backend/repository/edm_sql_test.go:76` TestSisipGenerasiMembawaKunciEndorsemen; Source of Change `frontend/components/BuatEDM.tsx` | Label kode 1-4 tidak ada di korpus - kode ditampilkan |
| 16 | Satu rumus `baris_ini − baris(OLD_POLIS_ID)` untuk seluruh generasi | 🟡 | `backend/models/edm_selisih.go` `EDMTCalculateTreatyDifference` (langkah 1-3) lewat `HitungSelisihGenerasi`; `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga | Prop ✅ (keputusan WO 23-09, ID-28/30; varian HasEDMNo tidak ditiru, #40). XOL `CalculateDifferenceEDM_act` dipertahankan apa adanya (batas bawah 0, prorata, pajak) - ditetapkan, sehingga rumus kedua ada di NonProp (#41) |
| 17 | Medan uang selisih = hasil pengurangan | ✅ | `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga, `backend/handlers/alur_test.go:380` TestPembatalanNolkanDataBaru, `backend/models/edm_xol_test.go:276` TestCalculateDifferenceEDMProrata |  |
| 18 | Medan persentase disalin | ✅ | `backend/models/edm_selisih.go` `EDMTCalculateTreatyDifference` langkah 1 (RiComm*, OveriddingComm*, Installment, TotalSharePercentage*), 2.1 (ClaimPercentage, SharePercentage), 3.1 (InstallmentPercentage) | Tanpa uji khusus |
| 19 | Medan kunci disalin | ✅ | `backend/models/edm_selisih.go` `EDMTCalculateTreatyDifference` langkah 2.1 (TreatyType, TreatyName), 3.1 (DueDate, InstallmentNo) | Currency/CurrencyID tidak disalin di Prop (#39); tanpa uji khusus |
| 20 | `DUE_TO` lapisan XOL dari tanda selisih | ✅ | `backend/models/edm_xol_selisih.go` `CalculateDifferenceEDM`; `backend/models/edm_xol_test.go:325` TestCalculateDifferenceEDMDueToIndukBasiDanBarisLama | #37 |
| 21 | `DUE_TO` per lapisan (bunyi koreksi) | ✅ | `backend/models/edm_xol_test.go:325` TestCalculateDifferenceEDMDueToIndukBasiDanBarisLama | #38; induk basi, nol pembaca, tidak bertabel |
| 22 | NonProp: angka sama sebelum / sesudah penyeragaman | ✅ | `backend/models/edm_xol_test.go:254` TestCalculateDifferenceEDMBatasBawahNol, `backend/models/edm_xol_test.go:276` TestCalculateDifferenceEDMProrata, `backend/models/edm_xol_test.go:290` TestCalculateDifferenceEDMPajakAdjPremi, `backend/models/edm_xol_test.go:307` TestCalculateDifferenceEDMBatalMentah | Terpenuhi karena rumus XML NonProp dipertahankan (#41) |
| 23 | Rumus selisih tidak di basis data | ✅ | `backend/migrasi_test.go:62` TestMigrasiHanyaEmpatTabelProyeksiSelisih, `backend/migrasi_test.go:131` TestUangSelisihBertipeNumber3810TanpaFloat (larangan PROCEDURE / TRIGGER / COMMIT) |  |
| 24 | Tabel selisih hanya ditulis aplikasi | 🔜 | `backend/repository/selisih.go` `SimpanSelisih` / `BuangSelisihGenerasi` - penulis tunggal (jalur biasa `simpanGenerasi`, pemuat) | Migrasi 360-363 belum dijalankan; SQL tulis belum pernah berjalan di Oracle |
| 25 | Ditulis dalam transaksi generasinya | 🔜 | `backend/services/tindakan.go` `simpanGenerasi` (SimpanHalaman + SimpanSelisih di satu `tulis`); `backend/handlers/alur_test.go:438` TestSelisihDihitungUlangSaatSubmitTanpaTombolHitung, `backend/handlers/alur_test.go:487` TestAdminMenolakMembuangSelisih | Logika ✅ di tiruan (selisih ditulis pada setiap simpan, tanpa bergantung tombol); Oracle menunggu migrasi |
| 26 | Bangun ulang hanya `SUMBER = 'GO'` | 🔜 | `backend/repository/selisih.go` `SimpanSelisih` (`ErrSelisihBeku` untuk setiap penulis, termasuk pemuat); `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi | Menunggu migrasi 360-363 |
| 27 | Isi beda dari hitung ulang ⇒ tabelnya salah | 🔜 | `backend/models/edm_bentuk.go` `HitungSelisihGenerasi` - proyeksi `'GO'` dihitung ulang setiap kali generasi ditulis; `backend/handlers/alur_test.go:438` TestSelisihDihitungUlangSaatSubmitTanpaTombolHitung | Belum ada alat banding tabel lawan hitung ulang; menunggu migrasi |
| 28 | Kunci saring `NOPOLIS PRODKE EDM_NO IDPEGA` | ✅ | DDL 360; `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga, `backend/handlers/alur_test.go:438` TestSelisihDihitungUlangSaatSubmitTanpaTombolHitung, `backend/migrasi_test.go:83` TestKatalogSelisihSepakatDenganDDL | #16 (nama `EDM_NO` dipertahankan); NOPOLIS kosong selama berjalan, terisi saat selesai |
| 29 | Nol tabel selisih induk XOL | ✅ | `backend/migrasi_test.go:62` TestMigrasiHanyaEmpatTabelProyeksiSelisih; induk = `BangunIndukSelisihXOL` (turunan) | #48 |
| 30 | Nol tabel selisih rincian angsuran | ✅ | `backend/migrasi_test.go:62` TestMigrasiHanyaEmpatTabelProyeksiSelisih |  |
| 31 | Repository ambil dua baris, tanpa agregasi | ✅ | `backend/repository/polis.go` `BacaGenerasi` (generasi `OLD_POLIS_ID`) + hitung di `models`; nol SUM / GROUP BY di SQL selisih | Tanpa uji khusus |
| 32 | Kolom khas endorsemen kosong di polis baru | 🟡 | `NOENDORS`, `OLD_POLIS_ID` hanya ditulis generasi EDM (`sqlSisipGenerasi`); `backend/repository/edm_sql_test.go:76` TestSisipGenerasiMembawaKunciEndorsemen | #26: daftar 21 kolom tidak berlaku (NB menulis sebagian besar, PRODKE NB = 0). Disusun ulang: yang khas endorsemen = `NOENDORS`, `OLD_POLIS_ID` (+ `PRODKE ≥ 1`); kosong di baris NB dijaga modul NB |
| 33 | Tiga kolom khas NB kosong di endorsemen, tidak dihapus | ✅ | Katalog EDM memuat `IS_EDM_INPUT_ON_NB`, `HAS_FAC_OUT`, `SHARE_CURRENCY`; nol penulis di jalur EDM (pengisi salinan NB tanpa pemanggil) | #27; tanpa uji khusus |
| 34 | Keadaan layar tidak tersimpan | ✅ | `backend/models/katalog_test.go:10` TestProyeksiKatalogHanyaMedanBerkolom; Show / ViewState / pyExpanded tidak berkolom |  |
| 35 | Spreading EDM presisi 20, NB 10 | ✅ | `backend/models/angsuran_edm_test.go:57` TestCountSpreadingEDMPresisiDuaPuluh, `backend/models/hitung_test.go:293` TestSpreadingEDMTidakDibagiRataSepertiNB |  |
| 36 | Baris pertama spreading bawaan 100 | ✅ | `backend/models/angsuran_edm_test.go:14` TestCountSpreadingEDMBawaanSeratus, `backend/models/angsuran_edm_test.go:34` TestCountSpreadingEDMBawaanHanyaBarisTunggal |  |
| 37 | Rincian angsuran bertingkat tersimpan (NonProp) | ✅ | `backend/models/edm_angsuran.go` `FillPaymentInstallmentEDMT`; `backend/models/edm_pilih_test.go:417` TestPilihBisnisEDMNonProp, `backend/models/dokumenlama_test.go:109` TestPecahDokumenEDMNonProporsional; `RapikanBentukSimpan` hanya membuang di Prop |  |
| 38 | Rincian dari salinan master tersimpan sama | ✅ | `backend/models/edm_pilih.go` `FillMasterInstallment` + `SetInstallmentValue`, lalu `FillPaymentInstallmentEDMT` membangun ulang (XML); `backend/models/edm_pilih_test.go:296` TestFillMasterInstallment, `backend/models/edm_pilih_test.go:417` TestPilihBisnisEDMNonProp | Butir WO #46 (sumber rincian yang disimpan) · **WO 07-10-2026** butir 13: mengikuti XML (`FillPaymentInstallmentEDMT` membangun ulang) |
| 39 | Nilai migrasi tidak dihitung ulang | ✅ | `backend/models/dokumenlama.go` `PecahDokumenEDM`; `backend/models/dokumenlama_test.go:29` TestPecahDokumenEDMProporsional, `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi | Digit galat utuh, varian berlapis Pega (130) tidak dihitung ulang. Pemuat belum dijalankan (bab 4) |
| 40 | Pemasangan antar generasi lewat `NOURUT` | ✅ | `backend/models/penanda_migrasi_test.go:32` TestPenandaPasanganBergeserMenurutNourut, `backend/models/penanda_migrasi_test.go:79` TestPenandaLapisanXOLMenurutSubskripMataUangDanLapisan |  |
| 41 | `PASANGAN_BERGESER` tanpa mengubah angka | ✅ | `backend/models/penanda_migrasi_test.go:32` TestPenandaPasanganBergeserMenurutNourut |  |
| 42 | `RUMUS_BERLAPIS = 1` bila generasi lama ber-EDMNo | ✅ | `backend/models/penanda_migrasi_test.go:59` TestPenandaRumusBerlapisSetiapBaris, `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi | #17 - cakupan dinyatakan: baris anak 361 / 362 (induk 360 tanpa kolom; 363 tanpa varian) |
| 43 | Penanda migrasi hanya baris `SUMBER = 'PEGA'` | ✅ | `backend/repository/lama_edm_test.go:48` TestPenandaHanyaBarisSumberPega, `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi | UPDATE ke 361-363 belum dicoba (migrasi) |
| 44 | Pemuat menulis lewat antarmuka repository yang sama | ✅ | `backend/repository/lama_edm_test.go:75` TestPemuatEDMTanpaJalurTulisTerpisah; `backend/services/pemuat.go` `GudangPemuat` |  |
| 45 | Satu transaksi per generasi | ✅ | `backend/services/tindakan.go` `tulis`, `BuatKasus`; `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi (dokumen ditolak nol baris) | Tiruan bersnapshot; Oracle belum |
| 46 | Skema `POOLDATA.` eksplisit | ✅ | `g.nama` (Qualify) + `db.PeriksaSQL`; `backend/repository/lama_edm_test.go:18` TestSQLPemuatEDMBerskemaTanpaCommit; penjaga repo `TestNolNamaTabelTelanjangDiQuery` | Penjaga `inti/backend/penjaga` tidak dijalankan pada pemeriksaan ini |
| 47 | Kolom pelaku dari identitas login | ✅ | `backend/services/tindakan.go` `kirim` (`CatatRiwayat` OPERATORID = akun login); `SisipKasus` pembuat = akun login | Tanpa uji khusus |
| 48 | Nol kolom uang `float` | ✅ | `backend/migrasi_test.go:131` TestUangSelisihBertipeNumber3810TanpaFloat, `backend/models/nilaipega_test.go:11` TestTeksSkalarJSONTanpaFloat |  |
| 49 | Skala sepuluh; kelebihan dibulatkan saat dimuat | 🟡 | `backend/migrasi_test.go:131` TestUangSelisihBertipeNumber3810TanpaFloat; `backend/repository/kolom.go` `nilaiTulis` (koefisien / skala, tanpa penolakan) | #28; pembulatan > 10 desimal oleh Oracle belum diuji |
| 50 | Kode / penanda teks, nol di depan utuh | ✅ | Katalog `kKode` / `kPenanda` → VARCHAR2 (`backend/models/katalog.go`) | Tanpa uji khusus |
| 51 | Teks kosong → `NULL` | ✅ | `backend/repository/kolom.go` `nilaiTulis` | Tanpa uji khusus |
| 52 | Arah `handlers → services → repository` | ✅ | Impor: handlers → services, models; services → repository, models; repository → models; penjaga repo `TestHandlersTidakMengimporRepository` |  |
| 53 | Satu antarmuka repository | 🟡 | Satu `repository.Gudang` di modul; pemuat memakai antarmuka yang sama (`backend/repository/lama_edm_test.go:75` TestPemuatEDMTanpaJalurTulisTerpisah) | Antarmuka itu SALINAN milik NB (modul tidak saling impor), bukan antarmuka NB yang sama |
| 54 | Bentuk tabel dasar tidak berubah | ✅ | `backend/migrasi_test.go:155` TestTabelGenerasiNBTidakDibuatUlang | #23 |
| 55 | Prop 10 tabel, NonProp 14 | 🟡 | — | #24, #53: generasi EDM juga menulis `T_POLIS_SURVEY` dan `HISTORYAKSEPTASIPRODUCTION` - cacah belum pasti · **WO 07-10-2026** butir 9 *"YA"*: generasi EDM juga menulis `T_POLIS_SURVEY` dan `HISTORYAKSEPTASIPRODUCTION` (pola NB) - cacah berbeda dari bunyi AC |
| 57 | `REMARK` ≥ 128 utuh pulang-pergi | ✅ | Katalog `Remark` → `REMARK` 128; `frontend/medan.test.ts` "Remark: dapat diisi hanya di posisi admin (disabled bila PositionNote != ReasTreatyInAdmin)" | T 12; uji pulang-pergi Oracle tidak ada |
| 58 | `NUMBER(38,10)`; kelebihan dibulatkan, bukan ditolak | 🟡 | Sama dengan AC 49 | Uji `ShareValue` 24 desimal tidak ada di modul ini |

AC 56 ditarik dari lingkup (23-09-2026 sore; diperkuat #54: `BreakDownSpreading_Act` di `CountSpreading_Act` 7 ber-`//`).

## 3 · `spec.md` — 59 AC

Tidak satu pun tiket EDM menutup AC `spec.md` (`REKONSILIASI-AC.md` bab 2); bukti di bawah langsung ke kode.

| AC | Ringkas | Status | Bukti | Catatan |
| ---: | --- | :---: | --- | --- |
| 1 | Endorsemen = jenis kasus tersendiri | ✅ | `backend/handlers/alur_test.go:209` TestCreateMelahirkanGenerasiBerikut, `backend/repository/edm_sql_test.go:38` TestDaftarEDMHanyaGenerasiEndorsemen, `backend/repository/edm_sql_test.go:89` TestKeadaanHanyaGenerasiEndorsemen | ID `EDMT-<n>`, `PRODKE ≥ 1` |
| 2 | Tangga tiga jenjang | ✅ | `backend/models/tangga.go` `Langkah`; `backend/handlers/alur_test.go:352` TestAtasanMenolakKembaliKeAdmin | Backend + seam HTTP lulus (Sec Head setuju → Dept Head; penyimpangan sadar #42). Layar atasan belum dapat membuka berkas: `antreanBeranda` tidak didaftarkan - menunggu konfirmasi WO · **WO 07-10-2026** butir 1 *"YA"*: `antreanBeranda` didaftarkan; cek Chrome 07-10: Sec Head Beranda → kasus → Accept → Submit |
| 3 | Tolak atasan ⇒ kembali ke Admin | ✅ | `backend/handlers/alur_test.go:352` TestAtasanMenolakKembaliKeAdmin | Transition8 / Transition4 lulus di seam; sama dengan AC 2 - layar atasan menunggu WO · **WO 07-10-2026** butir 1 *"YA"*: `antreanBeranda` didaftarkan; cek Chrome 07-10: Sec Head Beranda → kasus → Accept → Submit |
| 4 | Tolak Admin ⇒ berkas tertutup ditolak | ✅ | `backend/handlers/alur_test.go:326` TestAdminMenolakMelepasGenerasi, `backend/handlers/alur_test.go:487` TestAdminMenolakMembuangSelisih | `Resolved-Rejected`, tanpa Utility1; generasi dilepas dan proyeksi selisihnya dibuang (penyimpangan sadar) |
| 5 | `IsApproved` dibandingkan sebagai teks | ✅ | `backend/models/tangga.go` `Disetujui`; `backend/handlers/alur_test.go:326` TestAdminMenolakMelepasGenerasi | Uji ejaan lain (" 0", "0.0") tidak ada |
| 6 | Wewenang lewat peran | ✅ | `anggota` (`backend/services/layanan.go`), `TombolUntuk`; `backend/handlers/alur_test.go:401` TestWewenangDanPosisi | `IsSPVCreate` / `IsSPVTreaty1` / `IsTreaty1` dan tombol ber-ID operator tidak dibangun |
| 7 | Data lama disalin saat berkas dibuat | ✅ | `backend/models/edm_buat.go` `PasangOldData`; `backend/handlers/alur_test.go:209` TestCreateMelahirkanGenerasiBerikut, `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga | #13, #33: snapshot = baris generasi `OLD_POLIS_ID` yang dibekukan, tanpa JSON |
| 8 | Perubahan polis induk tak mengubah nilai lama | ✅ | Generasi tertutup tidak dapat disunting (`ErrGenerasiTertutup`, AC SP 6) | #13 |
| 9 | Snapshot bagian berkas | ✅ | Sama dengan AC 7-8 | #13 |
| 10 | Selisih terhadap keadaan tepat sebelumnya | ✅ | `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga | Diuji sampai lapis ketiga |
| 11 | Pemilih = isi `OldData.EDMNo` | ✅ | `frontend/medan.test.ts` "tab Old Data = PropOldData bila OldData.EDMNo kosong, PropOldData2 bila terisi; New Data menurut posisi"; `backend/models/penanda_migrasi_test.go:59` TestPenandaRumusBerlapisSetiapBaris | Cara hitung varian tidak ditiru (ID-28); pemilih dipakai tab Old Data dan penanda `RUMUS_BERLAPIS` |
| 12 | Satu polis boleh diendorse berkali-kali | ✅ | `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga | Tanpa batas |
| 13 | Nomor = induk + `/E` + dua digit | ✅ | `backend/models/edm_buat.go` `NomorEDM`; `backend/handlers/alur_test.go:209` TestCreateMelahirkanGenerasiBerikut | #44: sekurangnya dua digit |
| 14 | Nomor urut +1 dari terakhir | ✅ | `backend/repository/generasi.go` `GenerasiTerakhir`; `backend/handlers/alur_test.go:262` TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga, `backend/handlers/alur_test.go:326` TestAdminMenolakMelepasGenerasi, `backend/handlers/alur_test.go:473` TestCreateDitolakBilaJSONPolisLebihMaju | Generasi tolak tidak menahan nomor; Create ditolak bila json_polis lebih maju; #56 |
| 15 | < 10 dibubuhi nol | ✅ | `backend/handlers/alur_test.go:209` TestCreateMelahirkanGenerasiBerikut (`/E01`) | ≥ 100 tidak dipotong (`NomorEDM`), tanpa uji khusus |
| 16 | Adendum dan premi tambahan satu jalur nomor | ✅ | `NomorEDM` untuk semua EDMType; `GeneratePolicyNoTreatyAddendum_Act` / `GenerateNoEDMTreaty` tidak dibangun | #45; `[terbuka]` S 9.2 #3 tetap |
| 17 | Nomor urut di transaksi penyimpanan | ✅ | `backend/services/layanan.go` `BuatKasus` (generasi terakhir → ID → sisip dalam satu transaksi) | Tiruan; Oracle belum |
| 18 | Pembatalan = jenis endorsemen di awal | ✅ | `backend/handlers/alur_test.go:380` TestPembatalanNolkanDataBaru, `backend/models/edm_pilih_test.go:546` TestPilihBisnisEDMCancel | EDMType 4 |
| 19 | Pembatalan lewat jalur yang sama | ✅ | `SetEDMTCancel` di dalam `PilihBisnisEDM`; `backend/models/edm_pilih_test.go:546` TestPilihBisnisEDMCancel, `backend/handlers/alur_test.go:438` TestSelisihDihitungUlangSaatSubmitTanpaTombolHitung | #47 (cakupan nol - butir WO, lihat SP AC 14) |
| 20 | 1 baris dan % kosong/0 ⇒ 100 | ✅ | `backend/models/angsuran_edm_test.go:14` TestCountSpreadingEDMBawaanSeratus |  |
| 21 | `SplitRNMSharePct / RNMShare` presisi 20 | ✅ | `backend/models/angsuran_edm_test.go:77` TestCountSpreadingEDMPembagiRNMShare, `backend/models/angsuran_edm_test.go:57` TestCountSpreadingEDMPresisiDuaPuluh | Maksud dagang P60 `[terbuka]` |
| 22 | Ketidakseragaman presisi ditiru | ✅ | `backend/models/hitung_test.go:293` TestSpreadingEDMTidakDibagiRataSepertiNB |  |
| 23 | Spreading dihitung ulang saat premi berubah | ✅ | `backend/models/hitung.go` `CountOGPONP` langkah 9 → `CountSpreading` (aksi tab New Data, `validasiKirim`) | Tanpa uji khusus |
| 24 | Uang tidak pernah `float` | ✅ | `backend/models/nilaipega_test.go:11` TestTeksSkalarJSONTanpaFloat, `backend/migrasi_test.go:131` TestUangSelisihBertipeNumber3810TanpaFloat | `apd` di seluruh jalur uang |
| 25 | `/1,022` bila `Inclusive` persis | ✅ | `backend/models/hitung_test.go:90` TestPajakBrokerageInclusiveDibagi1022, `backend/models/hitung_test.go:100` TestPajakBrokerageSelainInclusiveApaAdanya, `backend/models/hitung_test.go:135` TestTypeTaxHurufKecilMenggeserBalance |  |
| 26 | PPH 2 %, PPN 2,2 % dari potongan sesudah pembagian | ✅ | `backend/models/setppnpph_test.go:20` TestSetPPNPPHNilaiXML, `backend/models/hitung_test.go:60` TestSetPPNPPHBersyaratFlagAtauPKP |  |
| 27 | Empat medan potongan/bagian = persentase | ✅ | Katalog: `Deduction1/2` `kUang`, RiComm / OveriddingComm `kPersen` (`backend/models/katalog.go`) | #31: `DEDUCTION1/2` = uang (NB K3) |
| 28 | Presisi penuh; bulat hanya di penyajian | ✅ | `backend/repository/kolom.go` `nilaiTulis` (tanpa pembulatan); `frontend/sajian.test.ts` "2 desimal: dipadankan dan dibulatkan setengah ke atas, titik ribuan koma desimal" | #14: > 10 desimal dibulatkan Oracle saat dimuat (lihat SP AC 49) |
| 29 | Uang ke pencapaian sebagai desimal tetap | ✅ | `backend/models/produksi_test.go:384` TestProduksiCapaianPropDanNonPropEDM, `backend/repository/produksi_test.go:100` TestSQLProduksiSelisihNegatifUtuh | ACHIEVEMENT INSERT langsung, nol prosedur |
| 30 | Jenis treaty dari medan penentu sebelum urai | ✅ | `backend/models/dokumenlama.go` `PecahDokumenEDM` (`ProportionalType`); `backend/models/dokumenlama_test.go:29` TestPecahDokumenEDMProporsional, `backend/models/dokumenlama_test.go:109` TestPecahDokumenEDMNonProporsional | Pemuat (bab 4) |
| 31 | `ListInstallment` bersarang terbaca | ✅ | `backend/models/dokumenlama_test.go:109` TestPecahDokumenEDMNonProporsional |  |
| 32 | `ListInstallment` datar terbaca | ✅ | `backend/models/dokumenlama_test.go:29` TestPecahDokumenEDMProporsional |  |
| 33 | Kedua bentuk diuji terpisah | ✅ | AC 31 dan 32 - dua uji terpisah |  |
| 34 | Medan hilang = kosong | ✅ | `backend/models/dokumenlama_test.go:29` TestPecahDokumenEDMProporsional (EndDate kosong tetap kosong) |  |
| 35 | Tiga bentuk contoh produksi terbaca | 🟡 | Fixture `backend/models/testdata/dokumenlama_edm_prop.json`, `dokumenlama_edm_nonprop.json` | Dua bentuk (proporsional, non-proporsional XOL); proporsional SOA tanpa fixture |
| 36 | FacOut tidak dibangun | ✅ | `backend/models/konversi.go` `RakitMuatanKonversi` (satu muatan); langkah 7 tidak dibangun | Butir WO #49 (`IsFacRetro` dapat menyala) |
| 37 | Pemeriksaan sukses menguji kedua penanda | ⛔ | — | `IsSuccessHitService` tidak dibangun: sambungan konversi belum disambung (sama NB) · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 38 | Gagal konversi ⇒ hapus data produksi | ⛔ | — | `DeleteDataProduction` tidak dibangun; #50 (`IsPEGAPROD`) · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 39 | Tidak menghapus bila `STS_KONVERSI = 1` | ⛔ | — | `CekSTSKonversiJson` tidak dibangun · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 40 | Penghapusan lewat satu fungsi repository | ⛔ | — | Nol penghapusan data produksi di kode; menunggu AC 38 · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 41 | Pengguna diberi tahu sebelum hapus | ⛔ | — | Menunggu AC 38 · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 42 | Pesan galat konversi sampai ke pengguna | 🟡 | `backend/services/konversi.go` `konversikan` → `HasilKirim.PesanKonversi` (`FlagErrorKonversi`) → pesan portal (`frontend/pages/LayarKasus.tsx`) | Hanya di produksi; pengirim bawaan gagal terang sehingga di produksi pesan selalu muncul. #50 (`IsTreatyIn`) |
| 43 | Penutupan paksa tidak dibangun | ✅ | `ASMForceCaseClose` (langkah 18 `//`) tidak dibangun |  |
| 44 | Tanpa autentikasi ke produksi | ⛔ | — | Sambungan belum ada - nol kredensial di kode · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 45 | Alamat tujuan dari variabel lingkungan | ⛔ | — | `LinkService` / `GetLinkService` tidak dibangun; nol alamat harfiah (penjaga repo `TestNolAlamatLayananDiKode`, tidak dijalankan di sini) · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 46 | Uji dan produksi beda basis data | 🟡 | Koneksi dari `inti.Dasar` (`backend/services/gudang.go` `DariDasar`); modul tanpa alamat basis data | Konfigurasi lingkungan milik inti - tidak diperiksa di modul ini |
| 47 | Satu transaksi | ✅ | `backend/services/tindakan.go` `tulis`; Utility1 di transaksi Dept Head (`kirim`); konversi sesudah transaksi | Tiruan; Oracle belum |
| 48 | Skema eksplisit | ✅ | Sama dengan SP AC 46 |  |
| 49 | Tanggal tutup buku satu baris global | ✅ | `backend/repository/riwayat.go` `HariClosing` (penomor inti, bawaan 25); `backend/models/tutupbuku_test.go:26` TestPraprosesTanggalBatasTutupBuku, `backend/models/tutupbuku_test.go:58` TestTanggalProduksiNomorBatasTutupBuku | #52 |
| 50 | Pelaku riwayat dari login | ✅ | Sama dengan SP AC 47 | Tanpa uji khusus |
| 51 | Tanpa pencarian nomor urut antrean | ✅ | `PraprosesAdmin` (`TempEmail.CARI28` tidak dibangun), `anggota`; `backend/handlers/alur_test.go:401` TestWewenangDanPosisi |  |
| 52 | Pemantauan idempoten | ⛔ | — | `INSERTJSON_JSONPOLISMONITORING_FACIN` tidak dibangun (konversi 15-16) · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 53 | Pemantauan tidak memperbarui | ⛔ | — | Sama dengan AC 52 · **WO 07-10-2026** butir 4 *"YA"*: konversi Arasapas tetap tidak disambung |
| 54 | Anti-duplikat pencapaian tanpa nilai uang | ✅ | `backend/repository/produksi_test.go:149` TestSusunTulisanProduksiPenjagaIDPega (penjaga = IDPEGA; ACHIEVEMENT tanpa pembanding uang) | Kunci pengenal `[terbuka]` 9.2 #2 |
| 55 | Pesan galat penyimpanan tanpa markah | ✅ | Galat `repository: ...` teks polos (`backend/repository`) | Tanpa uji khusus |
| 56 | Penanda lingkungan tak mengubah pesan | ✅ | Repository tidak membaca penanda lingkungan | Tanpa uji khusus |
| 57 | Ke-380 langkah `Property-Set` dibangun | 🟡 | `docs/INVENTARIS-XML.md` bab 2 | Langkah tidak dibangun hanya yang mati (`//`, sel disabled, syarat tak pernah benar), keputusan WO, atau menunggu WO - dirinci per rule; cacah 380 per langkah tidak diukur ulang |
| 58 | Rumus dari isi langkah | ✅ | Komentar kode mengutip langkah XML (`edm_selisih.go`, `edm_xol_selisih.go`, `edm_angsuran.go`, `edm_pilih.go`, `hitung.go`); uji berharapan XML (mis. `backend/models/setppnpph_test.go:20` TestSetPPNPPHNilaiXML) |  |
| 59 | Penelusuran sampai enam tingkat sarang | ✅ | Pemindaian korpus 06-10-2026 (pohon langkah bersarang); port mengutip langkah bersarang (mis. `CalculateDifferenceEDM_act` 1.2.5, `FillPaymentInstallmentEDMT` 2.2.2.3) |  |

## 4 · Pemuat dokumen lama (tiket 09-10)

**Dinilai dari kode** — berkasnya ada dan ujinya lulus: `backend/models/dokumenlama.go` (`PecahDokumenEDM`,
`UrutKunciGenerasi`), `backend/models/penanda_migrasi.go` (`HitungPenandaMigrasi`), `backend/models/laporanlama.go`,
`backend/models/usulanlama.go`, `backend/repository/lama_edm.go` (pembaca JSON_POLIS, `SetelPenandaMigrasi`,
`SetelKolomDatarLamaEDM`, `SalinUsulanLamaEDM`), `backend/services/pemuat.go` (`GudangPemuat`), dan alat
`backend/alat/pemuatlama` (uji-kering bawaan; `-jalankan` menulis, ditolak bila `IS_PEGA_PROD=true`).

| Yang dituntut | Bukti |
| --- | --- |
| Rantai tiga generasi, `OLD_POLIS_ID` = generasi tepat sebelumnya, `PRODKE` / `NOENDORS` dari dokumen | `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi |
| Penjaga jalur biasa ikut berjalan: percabangan (UNIQUE), keutuhan baris spreading, generasi sebelumnya tidak ada | `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi (`periksaGalatRantai`) |
| Selisih `'PEGA'` apa adanya, digit galat utuh; beku terhadap penulis mana pun | `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi, `backend/models/dokumenlama_test.go:29` TestPecahDokumenEDMProporsional |
| Penanda `PASANGAN_BERGESER` / `RUMUS_BERLAPIS` hanya baris `'PEGA'` | `backend/repository/lama_edm_test.go:48` TestPenandaHanyaBarisSumberPega, `backend/models/penanda_migrasi_test.go:59` TestPenandaRumusBerlapisSetiapBaris |
| Uji-kering nol tulisan; jalankan ulang aman; ID kasus bentrok tidak menimpa | `backend/services/pemuat_test.go:430` TestPemuatUjiKeringNolTulisan, `backend/services/pemuat_test.go:448` TestPemuatIDKasusBentrokTidakMenimpa |
| Bentuk dokumen (Prop datar, NonProp bersarang), penolakan berkas cacat | `backend/models/dokumenlama_test.go:109` TestPecahDokumenEDMNonProporsional, `backend/models/dokumenlama_test.go:147` TestPecahDokumenEDMLewatDanGalat, `backend/models/dokumenlama_test.go:230` TestPenggolongEDMMenolakBerkasCacat |

Status AC: SP 39-44 ✅ (bab 2); `spec.md` 30-34 ✅, 35 🟡 (bab 3). 🔜 **Belum pernah dijalankan** terhadap DEV: urutan
manual = migrasi NB 320-328 dan EDM 360-363 (WO) → pemuat NB → uji-kering pemuat EDM → `-jalankan` (WO/DBA).

## 5 · Penyimpangan sadar dan perubahan tinjauan kode

Daftar lengkap tiga belas penyimpangan sadar beserta dasarnya: `docs/INVENTARIS-XML.md` bab 1.3 (tangga tiga tingkat;
identitas operator → posisi; satu rumus selisih Prop, XOL apa adanya; selisih dihitung ulang setiap simpan; generasi
tolak dilepas; tanpa JSON; Create ditolak bila json_polis lebih maju; SuggestList → `HISTORYAKSEPTASIPRODUCTION` `'EDMT'`;
konversi tidak disambung; status / NBStatus dari data; wewenang = keanggotaan antrean; baris spreading tanpa hapus).
Label EDMType kini dari DT `TreatyEDMListType` (screenshot WO 07-10-2026) - bukan penyimpangan lagi.

Tinjauan kode 06-10-2026 (sudah tercermin di status bab 2-3):

1. `simpanGenerasi` SELALU menghitung ulang selisih Prop (`HitungSelisihGenerasi`), juga sebelum `SusunSimpananPolis` di
   Utility1 — selisih dan produksi tidak lagi bergantung tombol *Calculate Value Difference*: `backend/handlers/alur_test.go:438` TestSelisihDihitungUlangSaatSubmitTanpaTombolHitung.
2. Create ditolak bila json_polis memuat generasi lebih banyak dari rantai relasional (`CacahGenerasiJSONPolis`):
   `backend/handlers/alur_test.go:473` TestCreateDitolakBilaJSONPolisLebihMaju.
3. Admin menolak juga membuang proyeksi selisih `'GO'` (`BuangSelisihGenerasi`); `NOPOLIS` kunci selisih kosong sampai
   selesai: `backend/handlers/alur_test.go:487` TestAdminMenolakMembuangSelisih.
4. Baris selisih `'PEGA'` ditolak untuk setiap penulis, termasuk pemuat: `backend/services/pemuat_test.go:295` TestPemuatRantaiTigaGenerasi.
5. Frontend `frontend/penjagaIsian.ts`: permintaan diurutkan, isian selama permintaan diterapkan ulang; layar Create
   mengabaikan jawaban periksa-polis yang basi — `frontend/penjagaIsian.test.ts` "isian yang diketik selama permintaan berjalan diterapkan ulang di atas jawaban".

Keputusan work owner 07-10-2026 (jawaban 13 pertanyaan laporan 06-10):

| # | Jawaban | Penerapan |
| ---: | --- | --- |
| 1 | Kotak masuk Beranda untuk EDM: **YA** | `frontend/menu.ts` `antreanBeranda` + `daftarBeranda`; `frontend/Beranda.test.ts` bersama (satu baris); cek Chrome 07-10 jalur Sec Head |
| 2 | Aturan portal NB berlaku: **YA** | switch In Progress / Resolved (`?status=selesai`, `SaringanKasus.Selesai`); In Progress = buatan akun, Resolved = semua berkas selesai hanya-baca: `backend/handlers/alur_test.go` TestPortalInProgressDanResolved |
| 3 | Label EDMType: screenshot DT `TreatyEDMListType` | 1 Internal · 2 External · 3 Adjustment Premium · 4 Cancel Input — `models.LabelJenisEDM`, `frontend/labels.ts` `LABEL_JENIS_EDM` (dijaga sama: `frontend/labels.test.ts`); TestAcuanJenisEDMBerlabel |
| 4 | Konversi Arasapas tetap tidak disambung: **YA** | AC konversi ⛔ sadar |
| 5 | RNMShare: *"bukannya strukturnya sama dengan NB?"* → *"ikuti rekomendasi"* (b) | baris spreading ber-%Share kosong (baris tambahan) = pesan *"Spreading row n: % Share is required"*, bukan pembagian `PolicyTreatyIn.RNMShare` (tak diisi rule mana pun); Submit Admin tertahan lewat `CountOGPONP` langkah 9: `backend/handlers/alur_test.go` TestSubmitDitolakBarisSpreadingTambahanTanpaShare, `backend/models/angsuran_edm_test.go` TestCountSpreadingEDMShareKosongDitolakDenganPesan |
| 6 | Pembatalan: **ikuti XML dulu** | tetap 16 medan `SetEDMTCancel` |
| 7 | XOL Retro SOBName kosong: **IYA** (apa adanya) | tidak diubah |
| 8 | Popup Choose abaikan baris terpilih: **BETUL** | tidak diubah |
| 9 | SuggestList lama ke `HISTORYAKSEPTASIPRODUCTION`: **YA** | pemuat tidak diubah |
| 10 | Penanda induk 360: *"maksudnya?"* → *"ikuti rekomendasi"* | tanpa kolom baru; cara membaca (`SUMBER = 'PEGA' AND PRODKE >= 2`) di `docs/STRUKTUR-TABEL-EDM-TREATY-IN.md` |
| 11 | Atasan ubah With Tax / Type Tax / Overiding Commision / Marketing Officer: **TIDAK** | `medanHeader` + aksi header hanya Admin (backend), sel terkunci (frontend): TestAtasanTidakBolehMengubahHeader |
| 12 | Kode mati: **HAPUS** | 63 deklarasi Go + 20 uji, 3 ekspor TS + ujinya dibuang (analisis konservatif; dipertahankan: `Router`, `ProyeksiKatalog`, `MedanMasterXOL`, kepala iota `tanpaPemicu`) |
| 13 | Butir lain: *"?"* → *"ikuti rekomendasi"* | ditutup: P50 dan `IsPEGAPROD`/`IsTreatyIn` gugur (konversi tidak disambung; bila kelak disambung ikuti XML), rincian angsuran NonProp ikut XML, nilai lama disalin apa adanya (NB ID-15) - `KEADAAN-EDM-TREATY-IN.md` bab 11.6 |
| uji DEV | Current Premium NonProp tanpa pajak (*"UNTUK EDMT NON PROP KENAPA INI TU KOSONG??"*) | With Tax kosong diwarisi generasi lama (`WarisPajakLama`); With Tax / Type Tax berubah sesudah Choose Business -> rantai Choose Business diulang dengan master sama, isian dikembalikan (`HitungUlangPajakNonPropEDM`): `backend/handlers/nonprop_pajak_test.go` TestWithTaxSesudahChooseBusinessMenghitungPajakXOL, `backend/models/edm_pilih_test.go` TestHitungUlangPajakNonPropEDM* |
| uji DEV | *"Class Of Business GANTI JADII Source Of Business, SAMAIN DENGAN NB NYA!"* (lewat sesi NB TREATY) | `frontend/medan.ts` MEDAN_KIRI `.SOBName` tanpa syarat di tempat `.BizName`; dibuang dari MEDAN_KANAN: `frontend/medan.test.ts` "Source Of Business SEKALI di header" |
| uji DEV | *"JANGAN NULL TAPI 0"* + *"BUAT 4 ANGKA BELAKANG KOMA"* (screenshot Current Premium) | `frontend/xol.ts` `UANG4` (4 desimal, `nolPolos`) di KOLOM_XOL + KOLOM_XOL_RINCI; `components/Grid.tsx` "0" rata kanan: `frontend/pages/LayarKasus.test.tsx` "grid XOL dan rinciannya: angka 4 desimal" |
| uji DEV | *"KALO DAH RESOLVE STATUS NYA PAKE STATUS RESOLVE"* (screenshot tab Resolved) | `frontend/pages/PortalEDMTreatyIn.tsx` `statusPortal`: STATUS_WORK `Resolved-*`, selain itu NBStatus (pola portal NB): `PortalEDMTreatyIn.test.tsx` "kolom Status: berkas Resolved memakai STATUS_WORK" |
| uji DEV | *"BUATKAN TOMBOL COPY OLD SAMA SEPERTI MASTER PRODUCTNAME LIFE, KHUSUS BUAT SUPERUSER"* | `GET /hak`, `GET /lama`, `POST /lama/salin` (superadmin = Kelola User + hak penuh); `services/copyold.go` memakai `Pemuat.muat`, urutan generasi: `backend/services/copyold_test.go` TestCopyOld* (4 uji, 3 mutasi merah), `frontend/components/DialogCopyOld.test.tsx`, `frontend/lama.test.ts` — lihat MODUL.md bab Pemuat |
| uji DEV | *"pencarian ... buat bisa mencari nomor nb/edm, insured name dll, intinya buat searchnya itu sangat berguna"* | 14 kolom, AND antar-kata / OR antar-kolom, ESCAPE: `backend/repository/edm_sql_test.go` TestCariPortalBanyakKolomPerKata, `backend/handlers/cari_test.go` TestCariPortalNomorKasusPolisInsured, `backend/models/kasus_cari_test.go` (2 mutasi merah); DEV COUNT 07-10 lancar |

## 6 · Butir menunggu WO dan langkah manual

| # | Butir | AC terdampak |
| ---: | --- | --- |
| 1 | Filter D `TeamGroup` RD `InboxEDM_RD2` (pengisi parameter tidak ada di korpus) | — |
| 2 | `TreatyRealizationCheckXOLList` langkah 7 (OldData XOL ditimpa master) — `INVENTARIS-XML.md` bab 3 butir 1 | — |
| 3 | 🔜 Jalankan `-migrate` 360-363 dan 970 di DEV, lalu coba SQL tulis 360-363 | SP 24-27 |
| 4 | 🔜 Jalankan pemuat dokumen lama (bab 4) | SP 39-44 (di data nyata) |

## 7 · Perintah verifikasi

Dari akar repo (07-10-2026, sesudah keputusan WO 07-10 dan pembuangan kode mati):

```powershell
go test ./modul/edmtreatyin/...      # 6 paket ok (alat/pemuatlama tanpa uji); 150 fungsi Test
npx vitest run modul/edmtreatyin     # 15 berkas, 111 uji lulus
```

Penjaga lintas modul `inti/backend/penjaga` (nama tabel telanjang, arah impor, alamat layanan) dirujuk di atas tetapi
tidak dijalankan pada pemeriksaan ini.

