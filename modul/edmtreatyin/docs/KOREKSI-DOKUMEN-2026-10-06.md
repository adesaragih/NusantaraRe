# Koreksi dokumen EDM Treaty In — 06-10-2026

Rekonsiliasi dokumen modul `edmtreatyin` (spec, spec-penyimpanan, tiket 00–12, KEADAAN, PERTANYAAN-RONDE-1)
yang ditulis 22–23 September 2026, **sebelum** NB Treaty In dibangun (03–06 Oktober), terhadap:

- **NB nyata** — `modul/nbtreatyin/backend/migrations/320–328`, `backend/models/*.go`, dokumen NB;
- **XML** — korpus `D:\XML\RNM_BRD\EDM Treaty In\` (hanya-baca), langkah activity dibaca per nomor;
- **konsistensi** — pertentangan antar-dokumen EDM sendiri.

Aturan: bunyi lama **dikutip/dicoret, tidak dihapus**; `grilling-ronde-1.md` **tersegel** (koreksinya di
`KEADAAN-EDM-TREATY-IN.md` Bab 11.3); ⛔ nol butir `[terbuka]` milik WO/DBA ditutup tanpa keputusan tertulis —
butir yang ditandai tertutup selalu menyebut letak keputusan tertulisnya. Peta per AC: `REKONSILIASI-AC.md`.

Singkatan: **SP** = `spec-penyimpanan-relasional.md` · **S** = `spec.md` · **K** = `KEADAAN-EDM-TREATY-IN.md` ·
**PR1** = `PERTANYAAN-RONDE-1.md` · **T nn** = `issues/nn-…md` · **320** = `modul/nbtreatyin/backend/migrations/320_t_general_polis_treaty.sql`.

## 0 · Sensus korpus (tugas F0)

| Ukuran | 22-09 | 06-10 | Cara |
| --- | ---: | ---: | --- |
| berkas `.xml` | 163 | **163** | `glob` = `os.walk` (+1 `.xlsx`) |
| byte | 17.639.259 | ⛔ **17.640.859** | `getsize` |
| md5 gabungan | `24afd5726e083c9df6a692da94728139` | ⛔ **tidak tereproduksi** — korpus berubah | 5 varian resep (K Bab 11.1) |
| langkah `Property-Set` | 380 | **380 · 380** | A ElementTree · B pola teks |
| pasangan nama=nilai | 1.309 | **1.309 · 1.309** | A satu `rowdata` · B pola teks bersebelahan |

Berkas berubah sesudah 22-09: `Activity/CountSpreading_Act.xml` (22-09 17:50), `Activity/InsetTreatyInProdAddendum_Act.xml`
(23-09 14:31). Invarian baru (resep tertulis): md5 concat byte urut jalur penuh = `57926cc1b3f4524cb0cf73e1995ca211`.

## 1 · Log koreksi

| # | Dokumen | Bunyi lama (dikutip) | Bunyi baru | Bukti | Jenis |
| ---: | --- | --- | --- | --- | --- |
| 1 | T 00 kepala + tabel Ukuran | *"11 ⛔ tetap blocked"* · *"12 ⛔ blocked"* · *"⇒ 10 siap · 2 tertahan"* | 11 dan 12 `ready-for-agent`; **12 siap · 0 tertahan** | T 11 baris 3–14; T 12 baris 3–14; SP bab *KEPUTUSAN 23-09-2026 sore* | konsistensi |
| 2 | T 00 catatan audit | `ID-27b` *"tidak tersentuh satu tiket pun"*; `ID-27c` *"tidak terbawa ke tiket mana pun"* | keduanya **T 12** (bagian 1 / 2) | T 12 baris 18, 37–53 | konsistensi |
| 3 | T 00 tabel tiket & AC→tiket; T 11; T 12 | AC 57 di T 11 (*"45–53 · 57–58 (11 AC)"*); T 12 *"dua implementation decision"* tanpa AC | **AC 57 → T 12** (isinya `REMARK`, ID-27b); T 11 = 10 AC; baris T 12 ditambah | SP AC 57 *(ID-27b)*; T 12 bagian 1 | konsistensi |
| 4 | T 00, T 08 | *"AC 35–38 · AC 56 (5 AC)"* | **4 AC**; AC 56 ditarik | SP AC 56; T 08 baris 60 | konsistensi |
| 5 | T 00 | *"Tujuh tiket blocked"*; *"Empat siap dikerjakan — 01, 03, 04, 06"* | nol `blocked`; dua belas siap | tabel tiket T 00 | konsistensi |
| 6 | T 00 | *"⚠️ Wajib diputuskan sebelum pekerjaan dimulai"* (tumpang tindih NB 24–28) | sudah diputuskan: NB 24–28 **DIGANTIKAN** | `modul/nbtreatyin/docs/issues/24…28` baris 3–6, 32 | konsistensi |
| 7 | T 11 | *"spec ini **tertinggal**"* (AC 49 sembilan lawan AC 58 delapan) | sudah **diselaraskan** 23-09 sore | SP AC 49 blok *DISELARASKAN* | konsistensi |
| 8 | T 12 | *"Butir ini tidak ditutup di tiket ini. Surat permintaannya sudah disusun."* | dicoret — bertentangan dengan baris *gugur* di atasnya | T 12 baris 3–9 | konsistensi |
| 9 | T 02, 05, 07, 09, 10 | bab *"Kenapa tiket ini `blocked`"* masih berbunyi menahan | dicoret; rujuk keputusan 23-09 | kepala masing-masing tiket; SP bab KEPUTUSAN butir 1–4 | konsistensi |
| 10 | SP Bab 9 #1, #2, #4, #6, #7, #8 | `[terbuka]` | ✅ ditandai tertutup — keputusan WO **sudah tertulis** di SP bab KEPUTUSAN (butir 1–4, sore); #3 dan #5 tetap terbuka | SP bab *KEPUTUSAN 23-09-2026* dan *sore* | konsistensi |
| 11 | SP ID-46 | *"⛔ `[terbuka]` Dua belas digit di depan koma belum diuji"* | dicoret — ditutup keputusan WO 23-09 sore | SP bab *KEPUTUSAN … sore — Presisi* | konsistensi |
| 12 | S 9.2 #6 | *"Batas berapa kali … batas teknis 99"* `[terbuka]` | ✅ ditutup keputusan WO tertulis: tanpa batas | SP bab KEPUTUSAN butir 1 | konsistensi |
| 13 | S AC 7–9 lawan SP ID-8 / AC 9 | *"snapshot … bukan rujukan yang dibaca saat ditampilkan"* lawan *"nilai lama dibaca lewat `OLD_POLIS_ID`"* | didamaikan: snapshot = baris generasi sebelumnya yang **dibekukan** (generasi tertutup), tanpa tabel salinan | SP ID-13; 320 baris 18–20 | konsistensi |
| 14 | S AC 28 | *"pembulatan hanya di titik penyajian; … pembulatan di lapisan repository gagal"* | presisi penuh sampai 10 desimal; > 10 dibulatkan **saat dimuat** | SP AC 58 `[keputusan work owner]`; 320 baris 26–28 | konsistensi |
| 15 | SP ID-20 | *"`SPLIT_RNM_SHARE_PCT` — … **hanya EDM yang menulis**"* | NB juga menulis (pengisinya `TreatyNonPropSetSpreading` ada di NB) | K Bab 6.5; NB `backend/models/nonprop.go` | konsistensi |
| 16 | SP ID-27, T 07 | *"Gaya penamaan mengikuti `TREATYINPRODUCTION`"* | catatan: di `TREATYINPRODUCTION` dan tabel dasar nomor endorsemen = `NOENDORS`; kolom 360 `EDM_NO` **tidak diubah** | `RDBList/InsertTreatyInProdEDMT_SQL.xml`; 320 baris 34 | konsistensi |
| 17 | SP AC 42, T 09 | *"`RUMUS_BERLAPIS = 1` pada **setiap** baris selisih"* | catatan: kolom hanya di 361/362, induk 360 tidak punya — cakupan dinyatakan asisten utama | migrasi EDM 360–362 | konsistensi |
| 18 | S 8.2 | *"Tiga pilihan pada P29 — spec ini tidak memilih"* | penunjuk (bukan penutupan): NB SP ID-15 `[keputusan work owner]` *"galat … data lama diikuti apa adanya"*, diagram F20, EDM SP ID-33 — **butir WO** | `modul/nbtreatyin/docs/spec-penyimpanan-relasional.md` ID-15 | konsistensi |
| 19 | SP 1.1, 1.4, 7.5, bab KEPUTUSAN sore; S 1.1, 6.3; K Bab 1, 6.1; T 00, 02, 05, 07, 08, 09, 10, 11, 12 | `..\nb-treaty-in\…`, `..\claim-life\…`, `KEPUTUSAN-RONDE-12-…md` / `DAFTAR-MEDAN-…md` / `KEADAAN-NB-…md` tanpa jalur | `modul/nbtreatyin/docs/…`, `modul/claimlife/docs/…` | berkas ada di folder itu | konsistensi |
| 20 | PR1 baris 191–192, 780–781 | `` `..⏎b-treaty-in\PERTANYAAN-untuk-DBA.md` `` · `` `..⏎b-treaty-in\VERIFIKASI-P18.md` `` (jalur terpotong: `\n` terbaca baris baru) | `modul/nbtreatyin/docs/PERTANYAAN-untuk-DBA.md` · `…/VERIFIKASI-P18.md` | isi baris | konsistensi |
| 21 | NB T 24–28 *(tidak disunting)* | menunjuk `.scratch\edm-treaty-in\issues\…` | **dicatat saja** — berkas NB; jalur kini `modul/edmtreatyin/docs/issues/` | NB T 24–28 baris 6 | konsistensi |
| 22 | SP ID-5, ID-6 (diagram), AC 57; blok kepala | `T_GENERAL_POLIS` | **`T_GENERAL_POLIS_TREATY`** | 320 baris 1–6; `modul/nbtreatyin/MODUL.md` baris 75 | NB nyata |
| 23 | SP 1.1, ID-5, AC 54, Bab 8 #1; T 01; T 00 | *"Sembilan tabel dasar"* · *"Sepuluh tabel dasar dibuat tiket NB"* | `T_WORK_POLIS` + `T_GENERAL_POLIS_TREATY` + delapan anak 321–328 (**`T_POLIS_SURVEY`** ikut, di luar diagram) | migrasi NB 320–328; 328 baris 1–3 | NB nyata |
| 24 | SP AC 55; T 01 | *"proporsional **10** tabel; non-proporsional **14**"* | cacah = daftar diagram (memuat `HISTORYAKSEPTASIPRODUCTION`, tanpa `T_POLIS_SURVEY`) ⇒ **belum pasti** — butir WO | diagram *EDM Treaty In Prop* K146–K156, *NonProp* K171–K185 | NB nyata |
| 25 | SP ID-12, Bab 5.5, ID-20; T 01, 02 | `EDM_NO` · `PROD_KE` | **`NOENDORS`** · **`PRODKE NUMBER(10)`** | 320 baris 33–34 | NB nyata |
| 26 | SP Bab 5.5 + AC 32; T 01 | *"Dua puluh satu kolom khas endorsemen **kosong** pada baris polis baru"* | NB menulis sebagian besar (PRODKE = 0, OGP/ONP, `QUARTAL`, `STATEMENT_TYPE`, `ID_NEW_BISNIS`); kosong hanya `NOENDORS`, `OLD_POLIS_ID` ⇒ disusun ulang — asisten utama | NB `models/katalog.go` 137–184; `hitung.go` 186–188; `layar.go` 128, 334; `produksi.go` 63, 199 | NB nyata |
| 27 | SP ID-21, AC 33; T 01 | *"**empat** kolom hanya dipakai NB — … `IS_APPROVEDTO_DEPT_HEAD` …"* | **tiga**; `IS_APPROVEDTO_DEPT_HEAD` tidak ada (P36) | NB `spec.md` AC 64; `PERTANYAAN-untuk-Product-dan-Underwriting.md` baris 362; 320 | NB nyata |
| 28 | SP blok RALAT kepala #1–3, ID-46, AC 49, AC 58, bab KEPUTUSAN sore; T 11 | *"`NUMBER(38,8)` … 30 digit di depan koma … desimalnya tetap delapan"*; *"dibulatkan > 8"* | **`NUMBER(38,10)`**, 28 digit, desimal **10**; dibulatkan hanya > 10 | 320 baris 26–28; NB SP baris 5–17; NB `HASIL-IMPLEMENTASI.md` baris 257 | NB nyata |
| 29 | SP bab KEPUTUSAN sore; T 12 bagian 2 | *"Penampung medan tak dikenal tetap wajib, dan wajib **kosong**"* | tanpa penampung; medan diputuskan per medan; wajib nol = **belum diputuskan** (RALAT F3, `[keputusan work owner]` 04-10) | NB T 19 baris 45–54; NB T 22 baris 50 | NB nyata |
| 30 | SP ID-27b; T 12 bagian 1 | medan catatan *"tidak pernah masuk rancangan tabel"* | kolom **ada**: `REMARK VARCHAR2(128)` | 320 baris 80; NB `katalog.go` 165 | NB nyata |
| 31 | S AC 27 | *"Empat medan potongan dan bagian adalah **persentase**"* | `DEDUCTION1/2` = **uang** (K3, `[keputusan work owner]` 03-10-2026) | NB `spec.md` AC 26 blok RALAT; NB `katalog.go` (`kUang` Deduction1/2) | NB nyata |
| 32 | S 1.4, Bab 3, 9.2 #4 | *"bentuk tabelnya menunggu sensus properti"*; *"Sensus properti … belum dikerjakan"* | tabel ada (NB 320–328, EDM 360–363); sensus ada — #4 gugur oleh bukti | NB `SENSUS-PROPERTI-POLICYTREATYIN.md`, `DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`; `STRUKTUR-TABEL-EDM-TREATY-IN.md` | NB nyata |
| 33 | S 5.2 | data lama = salinan dokumen (`JSON_POLIS.DATA_JSON`) | sistem baru **tanpa JSON**: data lama = baris generasi yang ditunjuk `OLD_POLIS_ID` | NB `models/produksi.go` baris 63; SP ID-8 | NB nyata |
| 34 | T 00 *Gerbang NB* | *"Sisanya bergerbang NB 18 atau NB 19, yang keduanya menunggu DBA"* | NB 16–20 dan 22 **selesai** (sisa K11) | NB T 16 baris 3, 17 baris 3, 18 baris 20, 19 baris 14, 20 baris 3, 22 baris 14 | NB nyata |
| 35 | SP 7.5 | *"kode Go belum ada"* | NB dibangun 03–06 Oktober | NB `HASIL-IMPLEMENTASI.md` | NB nyata |
| 36 | T 10 *(tambahan)* | — | pemuat NB menolak `PRODKE > 0` (`ErrGenerasiEndorsemen` *"milik pemuat EDM tiket 10"*) | NB `models/dokumenlama.go` 113–115 | NB nyata |
| 37 | SP ID-29, AC 20; T 06 | *"`DUE_TO` diturunkan dari **tanda** selisih"* (seluruh selisih) | **hanya lapisan XOL**; jalur proporsional tidak menulis `DueTo` | `EDMTCalculateTreatyDifference` 1–6 (26 medan tanpa `DueTo`); `CalculateDifferenceEDM_act` 1.2.1/1.2.3/1.2.5; pembaca `InsetTreatyInProdAddendum_Act` 9.4/9.5 | XML |
| 38 | SP ID-32, AC 21; T 06 | *"`local.duetovalue` dijumlahkan sepanjang `ValueList`, lalu `@If(> 0…)`"*; *"dari **jumlah** … bukan dari satu lapisan"* | lapisan: tanda selisih **lapisan itu**; induk: diset 1.1 dari `local.duetovalue` **iterasi sebelumnya**; jumlah hanya ke `DueToValue` induk (2.3); **nol pembaca** | `CalculateDifferenceEDM_act` 1.1, 1.2.1, 2.1–2.3; `InsetTreatyInProdAddendum_Act` 10.1 (`"DUE TO US"` tetap) | XML |
| 39 | diagram `Diagram-Skema-Tabel-NusantaraRe.xlsx` *(tidak disunting — di luar wilayah)* | R104 / R119 *"kunci disalin TREATY_NAME · TREATY_TYPE · CURRENCY_ID"*; J94 / J109 *"DUE_TO = dari TANDA selisih"*; J93 / J108 *"nonprop … sudah seragam"* | XML 2.1 **tidak** menyalin `CurrencyID` (361 tanpa kolom itu benar); J94 → #37; J93 → #41 | `EDMTCalculateTreatyDifference` langkah 2.1 (6 medan: `ClaimPercentage SharePercentage ClaimSpreaded PremiumSpreaded TreatyType TreatyName`) | XML |
| 40 | SP ID-30, Bab 9.3 #9, 10.5 #4; T 06 | *"kedua blok … parameter percabangan kosong … pemilih hanya tertulis di keterangan langkah"* | **berpemilih**: langkah 1 `.OldData.EDMNo==""` salah ⇒ lompat `HasEDMNo`; langkah 4 benar ⇒ keluar. #9 gugur oleh bukti | `EDMTCalculateTreatyDifference` langkah 1 (`WhenFalse=1`, `HasEDMNo`), 4 (`WhenTrue=6`) | XML |
| 41 | SP ID-31, ID-28, AC 22; T 06 | *"Sisi non-proporsional … memang **sudah seragam**"* | bukan pengurangan polos: batas bawah 0, prorata, pajak ulang, mentah untuk jenis 4/2 ⇒ cakupan ID-28 atas NonProp — **butir WO** | `CalculateDifferenceEDM_act` 1.2.1–1.2.5 (`ProRatePercent,100,8` ×23) | XML |
| 42 | K Bab 3; S 5.1, AC 2; `grilling-ronde-1.md` baris 342 (via K 11.3) | *"Sejalan … tangga tiga jenjang … tidak ada penyimpangan"*; AC 2 `[terverifikasi]` | XML: Sec Head dapat menyelesaikan sendiri (`LetterNo` kosong, ≤ 200 juta); tiga jenjang = **keputusan** (brief F0-F1 06-10, prompt WO; NB K2) ⇒ penyimpangan sadar | `Flow/InputAddendumTreatyIn.xml` `Decision9`/`Transition24`; `When/ToTREATYDEPTHEAD.xml`; `CekLimitTreatyAcc_Act` 2.4 (2.2 `//`) | XML |
| 43 | K Bab 3; S 5.1 | *"22 sambungan"* | **24** baris `pyConnectors` (23 `TransitionN` + 1 sisa) | `Flow/InputAddendumTreatyIn.xml` | XML |
| 44 | SP ID-12, AC 5; S 5.4, AC 13, AC 15; T 02; PR1 P57 | *"dua digit"*; *"batas teknis 99"* | **sekurangnya** dua digit; ≥ 100 tiga digit; tanpa batas | `SetEDMTNoPolis` langkah 3 | XML |
| 45 | S 10.1; SP 10.2 | *"tidak ditemukan syarat penjaga apa pun … jalur kedua yang hidup"* | berpenjaga `EDMNo==""` (langkah 10, 12, 13); `EDMNo` selalu terisi ⇒ praktis tak menerbitkan `[dugaan]`; `[terbuka]` S 9.2 #3 **tidak ditutup** | `GeneratePolicyNoTreatyAddendum_Act` 10, 12, 13; `CreateEDMT` 15; `SetEDMTNoPolis` 3 | XML |
| 46 | SP ID-40, AC 38; T 08 | sumber *"`pyWorkPage.TreatyIn.Installment(n).InstallmentList`"*; *"rincian dari salinan master … tersimpan sama"* | sumber `TreatyIn.ValueDifference.Installment(n).InstallmentList`; pada `EDMT-` **dihapus & dibangun ulang** dari selisih XOL — **butir WO / asisten utama** | `DataTransform/SetInstallmentValue.xml`; `EDMChooseBusiness_Act` 7, 11, 12; `FillPaymentInstallmentEDMT` 1–3 | XML |
| 47 | SP ID-19, AC 14; S 5.5; T 05 | *"menolkan **seluruh** kolom uang"*; *"`CreateEDMT` langkah 16 memanggil `SetEDMTCancel`"* | 16 medan + angsuran + spreading; tidak: `GrossPremium NetPremium BalanceDueTo PPN PPH BalanceBefore* ClaimSpreaded ClaimPercentage`; panggilan efektif `EDMChooseBusiness_Act` 5 — cakupan **butir WO** | `SetEDMTCancel` 1.1–1.4, 2.x; `EDMChooseBusiness_Act` 5; `CreateEDMT` 16, 18 | XML |
| 48 | SP ID-7; T 07 | *"dirujuk hanya lewat delapan medan skalar, masing-masing satu kali … tidak membawa informasi baru"* | cacah literal berjangkar (8 / 85 terulang); rujukan **relatif** ada — induk = salinan kunci + **jumlah** lapisan (turunan); kesimpulan *tanpa tabel* tetap | `CalculateDifferenceEDM_act` 2.1–2.3; `FillPaymentInstallmentEDMT` 2.1; `FillSpreading` 3; `SetEDMAchivementValue`; `Section/DetailPolicyTreatyInAddPremi.xml` | XML |
| 49 | PR1 P50; K 5.3; S 5.7, AC 36, 8.1 #1 | *"pengiriman kedua **tidak pernah berjalan**"*; *"`FacOut` hanya muncul di dua berkas"* | `IsFacRetro` diset di jalur EDM (master `FacultativeShare > 0`) ⇒ langkah 7 dapat jalan; keputusan P50 tidak diubah — **butir WO** | `InputPolicyTreatyEDMDetail_NP` / `_AdjPremi` langkah 3; `SetValueEDM_Act` 9/10; `serviceInsertArasapas_act` 7; `When/IsFacRetro.xml` | XML |
| 50 | S 5.7, AC 38, AC 42 | penghapusan bersyarat *"hanya bila `STS_KONVERSI` belum 1"*; *"langkah 10 dihidupkan"* | juga bersyarat **`IsPEGAPROD`**; langkah 8 dan 10 dilewati bila **`IsTreatyIn`** — **butir WO** | `serviceInsertArasapas_act` 8, 10, 12 | XML |
| 51 | S 5.2 | *"`Activity\CreateEDMT` — 20 langkah"* | **21** (18 aktif, 3 `//`: 4, 5, 6) | `Activity/CreateEDMT.xml` | XML |
| 52 | S 5.8, AC 49 | *"dibaca `WHERE ROWNUM = 1`"* | naskah `SELECT * FROM POOLDATA.TANGGAL_CLOSING`; baris pertama `pxResults(1)`, kosong → 25; P54 tetap | `RDBList/GETTanggalClosing_SQL.xml`; `SaveJsonPolisTreatyInEDM_Act` 2–3 | XML |
| 53 | SP ID-5; AC 55 | *"ditambah `POOLDATA.HISTORYAKSEPTASIPRODUCTION`"* | korpus EDM **nol rujukan** (juga nol `SaveViewSuggest`/`InsertViewSuggest`) — EDM Pega tidak menulisnya; ikut NB K4 atau tidak — **butir WO** | sapuan 163 `.xml` + `.xlsx` | XML |
| 54 | T 08 / SP AC 56 *(penguat)* | AC 56 ditarik atas keputusan WO | diperkuat: `CountSpreading_Act` EDM langkah 7 `Call BreakDownSpreading_Act` ber-`//` | `Activity/CountSpreading_Act.xml` langkah 7 | XML |
| 55 | K 2.1, Bab 8; S 11.1 | 17.639.259 B; md5 `24afd572…` | 17.640.859 B; md5 berubah; 380 / 1.309 tetap | bab 0 di atas; K Bab 11.1 | XML |
| 56 | S 5.4 *(catatan)* | — | `PRODKE` ke `JSON_POLIS` = `COUNT(*)` per `NOPOLIS`, berbeda dari maks+1 `substr(…,1,24)` nomor EDM | `RDBList/TreatyInSearchProdKe.xml`; `SaveJsonPolisTreatyInEDM_Act` 11–12; `RDBList/SelectProdKe.xml` | XML |

## 2 · Cacah per jenis

| Jenis | Baris | Nomor |
| --- | ---: | --- |
| konsistensi | **21** | 1–21 |
| NB nyata | **15** | 22–36 |
| XML | **20** | 37–56 |
| **Jumlah** | **56** | |

## 3 · Butir yang butuh keputusan work owner (tidak ditutup di sini)

1. Cakupan *"satu rumus untuk semua"* (SP ID-28) atas NonProp — XML memakai batas bawah 0 / prorata / pajak ulang (#41); AC 22 menuntut angka NonProp tidak berubah.
2. P50 bila `IsFacRetro` = 1 (#49).
3. Pembatalan: 16 medan XML atau seluruh uang (#47).
4. Sumber rincian angsuran NonProp — salinan master atau dibangun ulang dari selisih XOL (#46).
5. `T_POLIS_SURVEY` dan `HISTORYAKSEPTASIPRODUCTION` pada generasi EDM; cacah AC 55 (#24, #53).
6. Gerbang `IsPEGAPROD` pada penghapusan dan `IsTreatyIn` pada pesan galat (#50).
7. P29 (a) untuk EDM — penunjuk ke NB ID-15 (#18).

Untuk asisten utama (bukan WO): daftar kolom khas AC 32 (#26), cakupan `RUMUS_BERLAPIS` induk (#17).

## 4 · Batas ronde

| Larangan | Dipatuhi | Bukti |
| --- | :---: | --- |
| sunting hanya `modul/edmtreatyin/docs/` | ✅ | 15 berkas disunting (SP, S, K, PR1, T 00–02, 05–12) + 2 berkas baru, semuanya di folder itu |
| `grilling-ronde-1.md` tidak disentuh | ✅ | koreksinya di K Bab 11.3 |
| kode, `modul/nbtreatyin`, `inti`, `frontend`, korpus | ✅ | **0** tulis (dibaca saja) |
| DB, `.env`, git | ✅ | tidak dipakai |
| nama orang · nomor polis · data pribadi | ✅ | **0** — identitas operator, alamat surel dan host di korpus tidak disalin |
| butir `[terbuka]` WO/DBA ditutup tanpa keputusan tertulis | ✅ | **0** |

⚠️ Token, durasi, biaya ronde ini **tidak diukur** — angka token sejati tidak terlihat dari dalam sesi.
