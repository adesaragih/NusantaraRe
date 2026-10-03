# 18: Tipe kolom, konversi masuk, dan presisi uang

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* ⭐ Presisi dinaikkan ke **`NUMBER(38,8)`** — **30 digit di depan koma**, delapan di belakang, batas tertinggi Oracle. Dasarnya sapuan korpus: ambang dagang nyata sudah **tepat di batas** 12 digit *(`181500000000.00` · `150000000000.00`)*, dan ada sentinel **17 digit** *(`99999999999999999.99`)*. Penjumlahan lintas mata uang dapat melewati keduanya. ⭐ `NUMBER` di Oracle berpanjang **berubah-ubah** — hanya digit bermakna yang tersimpan, sehingga pelebaran ini **tidak memakan ruang tambahan**. Butir ditutup oleh bukti, bukan oleh DBA.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** selesai — penahan tersisa hanya pihak luar: **K11** (skema uji Oracle; uji bertag `db` AC 12–13 sudah ditulis, belum dijalankan) *(putaran 2, paket P11 04-10-2026: AC 25 bukti uji + RALAT ID-20, AC 15 bunyi baru RALAT, AC 12 uji `db` kolom `GROUP_PANEL`; semula: sebagian — konsolidasi P10 04-10-2026; implementasi 2026-10-03; awalnya ready-for-agent)*
~~**Blocked by:** **16** · ⛔ `[data DBA]` **presisi fisik belum diuji terhadap nilai terbesar** — dua belas digit di depan koma belum dibuktikan cukup~~ ⛔ **gugur 23-09-2026 sore**
**Menutup:** NB AC **12–25** *(14 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-14..ID-20

## Hasil & nilai pengguna

Nilai yang masuk dari dokumen lama seluruhnya bertipe teks — termasuk uang, tanggal, dan penanda.
Tiket ini menetapkan bagaimana teks itu menjadi kolom bertipe, **sekali saat masuk**, bukan setiap
kali dibaca.

⛔⛔ **Yang paling mudah salah, dan akibatnya besar:** kode `"006"` yang disimpan sebagai bilangan
lalu dibaca balik menjadi `"6"` akan **lolos uji pulang-pergi** tetapi memecahkan penggolong jenis
usaha 36 baris — dan kegagalannya diam.

## Yang dibangun

| Golongan | Ketetapan |
| --- | --- |
| **uang dan persen** | angka presisi tetap; ⛔ **tidak pernah** bilangan mengambang |
| **kode** | ⭐ **tetap teks** — nol di depan membawa makna |
| **penanda** | ⭐ **tetap teks** — kosong adalah keadaan sah yang **berbeda** dari nol |
| **tanggal** | dua format masuk: delapan digit, dan cap waktu bersufiks zona |
| **teks kosong** | menjadi **kosong**, bukan nol |

⚠️ Pembandingan dua nilai uang memakai **toleransi atau bentuk terbulatkan**, bukan kesamaan
persis — data produksi terbukti membawa galat pecahan.

⛔ **RALAT P11 (04-10-2026)** — bunyi lama dikutip di atas → bunyi baru: port rumus **tidak pernah**
membandingkan dua nilai uang satu sama lain (XML terjangkau pun tidak); pembandingan uang **lawan nol**
mengikuti XML eksak. Bukti di bab P11.

## Batas — yang TIDAK termasuk

⛔ `CREATE TABLE` — presisi fisik dicocokkan DBA **di dalam** tiket ini.
⛔ Pemecahan dokumen menjadi baris — tiket **19**.

## Cara mengujinya

Lewat seam `repository`, dan ⚠️ **sebagian test memeriksa nilai kolom langsung** — pulang-pergi
saja tidak cukup, sebab tulis dan baca yang sama-sama salah simetris tetap hijau.

⭐ Uji yang wajib: kode `"006"` disimpan lalu **dibaca dari kolomnya**, dan penggolong jenis usaha
dijalankan atasnya — hasil bawaan berarti gagal.

## Acceptance criteria

- [ ] 🟡 **AC 12** — kode tiga digit berawalan nol tersimpan dan terbaca utuh *(P11: `repository/penyimpanan_db_test.go` `TestKodeBernolDepanUtuhDiKolom` membaca kolom `GROUP_PANEL` langsung — K11)*
- [ ] 🟡 **AC 13** — kode dua digit berawalan nol tersimpan utuh *(P11: kolom `BUSINESS_OLD_ID` dibaca langsung, uji yang sama — K11)*
- [x] **AC 14** — penggolong jenis usaha menemukan barisnya, bukan nilai bawaan
- [x] **AC 15** — ~~penanda kosong tersimpan kosong, **bukan** nol dan bukan tak-bernilai~~ ⛔ RALAT P11: penanda kosong tersimpan tak-bernilai (Oracle `''` ≡ NULL) dan terbaca kembali kosong; **tidak pernah** menyatu dengan nol *(bab P11)*
- [x] **AC 16** — penanda bernilai nol menghasilkan keputusan ditolak; nilai lain disetujui
- [x] **AC 17** — teks kosong pada medan uang tersimpan tak-bernilai
- [x] **AC 18** — teks kosong pada medan tanggal tersimpan tak-bernilai
- [x] **AC 19** — nilai uang berdesimal sembilan **dibulatkan pada desimal kedelapan**, bukan dipotong ke dua
- [x] **AC 20** — kolom uang berskala **delapan desimal**, ⭐ **tiga puluh digit di depan koma** *(`NUMBER(38,8)`; semula ~~dua belas~~ — dinaikkan 23-09-2026 sore)*
- [x] **AC 21** — tanggal delapan digit terurai benar *(putaran 2, pemuat tiket 22: `models.BacaTanggalLama` — `TestBacaTanggalLama`)*
- [x] **AC 22** — cap waktu bersufiks zona terurai benar *(putaran 2: `20170930T170000.000 GMT` → `2017-10-01 00:00:00` Asia/Jakarta — `TestBacaTanggalLama`; tanggal ambigu tidak ditebak, K15 — `TestTanggalAmbiguTidakDitebak`)*
- [x] **AC 23** — pengurutan menurut tanggal menghasilkan urutan kronologis, bukan leksikal
- [x] **AC 24** — nol kolom uang bertipe mengambang
- [x] **AC 25** — pembandingan uang memakai toleransi, bukan kesamaan persis *(RALAT P11: port tidak membandingkan dua nilai uang; lawan nol eksak seperti XML — bab P11)*

## ⛔ Kenapa tiket ini `blocked`

~~`[data DBA]` **Dua belas digit di depan koma belum diuji terhadap nilai terbesar.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* Presisi dinaikkan ke **`NUMBER(38,8)`** — tiga puluh digit di depan koma. Dasarnya sapuan korpus: ambang dagang nyata **tepat di batas** dua belas digit *(`181500000000.00`)*, dan ada sentinel **tujuh belas digit** *(`99999999999999999.99`)*. ⭐ `NUMBER` di Oracle berpanjang berubah-ubah, jadi pelebaran ini **tidak memakan ruang tambahan**.

⚠️ **Dan satu pertentangan di dalam spec, dicatat di sini karena spec tidak boleh disunting:**
~~`[terverifikasi]` spec EDM **AC 49** menuntut sembilan desimal, bertentangan dengan AC 58.~~ ✅ **DITUTUP 23-09-2026 sore.** AC 49 diselaraskan ke **delapan**; bunyi lamanya dikutip di spec EDM. ⭐ Pertentangan ini ditemukan **pelaksana ronde tiket**, dan laporannya terbukti tepat..

## ⭐ Penerapan KEPUTUSAN-RONDE-12 dan catatan — 2026-10-03

- **Butir 8** — nomor generasi `PRODKE NUMBER(10) DEFAULT 0` (bilangan bulat). ID-14: `INSTALLMENT_NO`
  juga bilangan bulat `NUMBER(10)`.
- **AC 15** — Oracle menyimpan `''` sebagai NULL; penanda kosong dibaca kembali sebagai `""` (setara di
  halaman, tidak setara di SQL) — sebagian.
- **AC 21-22** (format dokumen lama `YYYYMMDD`, cap waktu ` GMT`) milik pemuat dokumen lama — tiket 22,
  belum dibangun.
  ⛔ **RALAT putaran 2 (03-10-2026):** bunyi lama *"belum dibangun"* → **dibangun** (tiket 22,
  `models.BacaTanggalLama`).

## ⭐ Putaran 2 — paket penyimpanan (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 11–12, bab 2 K4/K16/K17; rincian kolom `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

- Golongan tipe tidak berubah (uang/persen `NUMBER(38,8)`, tanggal `DATE`, cacah `NUMBER(10)`, kode dan
  penanda teks); kosong → `NULL` di kolom angka/tanggal (diagram F25) tetap `repository.nilaiTulis`.
- Kolom baru mengikuti golongannya: `IS_EDM_INPUT_ON_NB` penanda, `ID_NEW_BISNIS` kode, `PREMIUM_AFTER_*`
  uang, `CEDING_CO_ID` kode. `DEDUCTION1/2` milik paket layar (K3) — tidak disentuh paket ini.

## ⭐ Putaran 2 — paket P11 (04-10-2026): AC 12, 15, 25

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0, bab 2 (K11), bab 7; `docs/HASIL-IMPLEMENTASI.md` bab 9 (semula).

### AC 25 — sisir setiap pembandingan uang, port lawan XML

Sisir: setiap ekspresi berpembanding (`==` `!=` `>` `<` `>=` `<=` `compareTwoValues`) di 176 rule terjangkau
(`docs/alat/graf.py` `terjangkau()`; prasyarat langkah, transisi, parameter `Property-Set`, syarat sel Section,
When, DecisionTable — RDBList/RD dikecualikan, SQL-nya tidak membandingkan nilai halaman), lalu setiap `.Cmp` /
`.Sign()` / `.IsZero()` / `==` teks di `backend/models`. **Nol pembandingan dua nilai uang yang hidup.**

| Golongan pembandingan | XML (berkas · langkah · ekspresi) | Port (`backend/models`) |
| --- | --- | --- |
| uang lawan nol (tanda) | `Activity/SetDueTo_act.xml` 1 `.BalanceDueTo>=0`, 2 `.BalanceDueTo<0` | `hitung.go` `SetDueTo` `b.Sign() >= 0` |
| uang lawan nol | `Activity/CountNetPremi_act.xml` 6 `.Deduction1!="" \|\| .Deduction1!=0` | `hitung.go` `CountNetPremi` langkah 6 `Sign() != 0` |
| uang lawan nol / teks `"0"` | `Activity/CountOGPONP_Act.xml` 1–2 `.PremiOgp == "" \|\|.PremiOgp ==  "0"` (dan `.PremiOnp`); `CountResult1_Act` 1 dsb. | `hitung.go` `p.teks("PremiOgp") == "0"` |
| uang lawan nol | `Activity/CountOGPONP_Act.xml` 8 `.Claim!=0&&.Claim!=""`, `.SalvageValue!=0&&…` | `hitung.go` `terisi` `Sign() != 0` |
| uang lawan nol | `Activity/InputPolicyTreatyInDetail_preACT.xml` 17 `@if(.Limit>0,"IDR","")`, `@if(.Limit2>0,"USD","")` | `nonprop_detail.go` `Sign() > 0` |
| uang lawan nol (layar) | `Section/DetailPolicyTreatyIn.xml` wajib `.Claim != '' && .Claim != 0`; `Section/DetailDeptHeadTreatyIn_UW.xml` tampil `.BalanceDueTo < 0` / `>= 0` | `layar.go` `adaKlaim` `Sign() != 0` |
| persen lawan 100 | `CountResult1_Act`, `CountResult2Ogp_act`, `CountResult1Onp_Act`, `CountResult2Onp_act` `.RiCommOgp>100`, `<100` dst.; `SetValidateInstallment_Act.xml` 4 `local.pcttotal>100`; angsuran 3.3 | `hitung.go` `Cmp(seratus)` (8 tempat), `angsuran.go` `Cmp(seratus)` (2) |
| persen lawan 0/1 | `FacultativeShare != 0` / `==0` (`InputPolicyTreatyInDetail_NonProp` 20–21), `Section/SpreadingRiskList.xml` | `nonprop*.go` `IsZero()` |
| penanda lawan 1 | `Activity/SetPPNPPH.xml` 4 `ListAgent.pxResults(1).STS_PKP == 1`; `IsNewPolicyNonProp==1` | `hitung.go` `samaDenganSatu` `Cmp(apd.New(1, 0))` |
| ⛔ rasio dua uang bertoleransi | `((.ResultOgp2/.PremiOgp)-.OveriddingCommOgp)<=0.01` dan saudaranya — **hanya** di langkah berlabel `//`: `CountOverridingCommOgp_Act` 1/2/4, `CountOverridingCommOnp_Act` 1/2, `CountRiCommOgp_act` 4, `CountRiCommOnp_act` 1 | tidak diport (dinonaktifkan di rule) |
| ⛔ uang lawan ambang | `Activity/CekLimitTreatyAcc_Act.xml` 2 `Local.NETPREMI>200000000.00`, 1 `Local.NETPREMI<0.0` | tidak dibangun — **K2** |

⇒ RALAT ID-20 dan AC 25 di `spec-penyimpanan-relasional.md` (bunyi lama dikutip): pembandingan **dua** nilai uang
tidak boleh sama-persis — dan tidak ada; pembandingan lawan nol eksak seperti XML. Uji
`backend/models/pembandingan_uang_test.go`: `TestPortTidakMembandingkanDuaNilaiUang` (penjaga AST — setiap `.Cmp`
lawan tetapan rule, nol teks medan uang katalog dibandingkan persis dengan teks medan lain; diuji merah dengan
mutasi sementara `NetPremium.Cmp(BalanceDueTo)` dan `teks(PremiOgp) == teks(PremiOnp)`) dan
`TestTandaUangLawanNolEksakSepertiXML` (`-0.00000001` → DueTo `0`, `0.00000001` dan `0` → `1`).

### AC 15 — RALAT bunyi (Oracle `''` ≡ NULL)

Bunyi lama: *"`IsApproved` bernilai `""` tersimpan sebagai `""`, bukan `NULL` dan bukan `"0"`"* → bunyi baru (ditulis
di `spec-penyimpanan-relasional.md` AC 15): `""` tersimpan NULL dan terbaca kembali `""`; `"0"` tersimpan `'0'` dan
terbaca `"0"`; keduanya tidak pernah menyatu. Dasar XML yang dijaga: `DecisionTable/isApproved.xml` (kolom `text`,
`= 0` → `No`). Uji: `repository/kolom_test.go` `TestNilaiTulisKosongJadiNULL`, `TestNilaiBaca`;
`repository/penyimpanan_db_test.go` `TestIsApprovedKosongDanNolTetapBerbeda` (bertag `db`, K11).

### AC 12–13 — uji `db` nilai kolom langsung

`repository/penyimpanan_db_test.go` `TestKodeBernolDepanUtuhDiKolom`: `Quotation.GroupPanel "006"` dan
`BusinessOldId "01"` disimpan lewat `SimpanHalaman`, lalu kolom `T_POLIS_QUOTATION.GROUP_PANEL` /
`BUSINESS_OLD_ID` dibaca **langsung** (= `006` / `01`) beserta `DATA_TYPE` `VARCHAR2` dari `ALL_TAB_COLUMNS`.
⛔ Belum dijalankan — **K11**.

Status: **selesai** — sisa penahan hanya K11 (AC 12, 13 🟡).
