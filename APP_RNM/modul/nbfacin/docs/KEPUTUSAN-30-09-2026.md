# Keputusan work owner — 30 September 2026

Dicatat di berkas tersendiri, **bukan** sebagai K-0xx baru di `00-KEPUTUSAN-WORK-OWNER.md`: register itu
salinan byte-identik dari `D:\migrasi\RNM\OUTPUT\` dan belum ada keputusan siapa yang memberi nomor K
berikutnya di repo ini. Bila work owner ingin ketiganya masuk register, beri nomor di sana dan rujuk
berkas ini.

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 1 | Di mana register keputusan (K-001…) dan dokumen pendampingnya disimpan di repo? | `APP_RNM/modul/nbfacin/docs/` — disalin byte-identik (26 berkas, `cmp` 26/26) |
| 2 | Fac Out masuk `nbfacin` atau modul baru? | **Modul baru** `facout` (bab 5 panduan tim). **Diparkir** di `jefri/OUTPUT FIX/_usulan-modul-facout/` — lihat bab di bawah |
| 3 | Tiket NB-01/02 (`pkg/money`, `panic`, keadaan `Unknown`) bertentangan dengan `inti/backend/uang` | **Ikuti `inti/backend/uang`** apa adanya |

## Arti keputusan 3 bagi tiket yang sudah tertulis

Isi tiket **tidak diubah**; saat mengerjakannya, baca butir di bawah sebagai pengganti.

| Tiket menyebut | Dikerjakan sebagai |
| --- | --- |
| `pkg/money`, `pkg/ratio` (NB-01) | `inti/backend/uang` — `Money`, `Ratio` — **tidak dibuat ulang**, tidak disalin ke `nbfacin` |
| Aritmetika lintas mata uang → `panic` (NB-02) | `Money.Add` / `Money.Sub` mengembalikan **`uang.ErrMataUangBerbeda`**; pemanggil wajib memeriksa galatnya |
| Keadaan `Unknown` eksplisit (NB-02, ADR-F-0006) | Tidak ada di `inti/backend/uang`. `Money.Kosong()` = nilai kosong (`ErrUangKosong` bila disentuh aritmetika) — **bukan** padanan `Unknown`. ⚠️ `Currency` adalah `string` biasa: dua `Money` bermata uang `""` **dijumlahkan tanpa galat** (`uang.go` `Add`/`Sub` hanya membandingkan `m.Currency != other.Currency`). Bagaimana nilai uang tanpa mata uang dari korpus ditangani tetap **`[pertanyaan terbuka]`** sampai spec `nbfacin` menjawabnya; jangan diisi mata uang bawaan |
| Uji gagal-kompilasi `Money + Ratio` (NB-01) | Sudah dijamin bahasa: `Money` dan `Ratio` struct berbeda, Go tanpa operator struct (komentar `uang.go`, ADR-F-0004) |
| Perubahan apa pun pada `inti/backend/uang` | Pull request ke **tim inti** — paket dipakai `claimlife` dan `premiumlistlife` |

## Kenapa `facout` diparkir, dan yang harus diajukan ke tim inti

`modul/facout/` sempat dibuat, lalu dipindah ke `jefri/OUTPUT FIX/_usulan-modul-facout/` (keputusan work
owner, 30-09-2026) karena membuat **6 uji frontend merah**: `inti/frontend/uji/sumber.ts`
(`folderKorpusBelumDimigrasi`) menganggap setiap folder `modul/*` tanpa awalan `_` adalah folder korpus
dengan baris menu di isi awal 900 — dan 900 tidak pernah disunting. Penjaga Go
(`inti/backend/penjaga`) sendiri hijau dengan folder itu.

Pull request tim inti yang dibutuhkan (bab 5 `docs/bersama/PANDUAN-TIM-PER-MODUL.md`):

1. rentang migrasi dan slot menu — usulan `760-799` / `990-991` (tertulis di `MODUL.md` usulan);
2. kelompok `M_NAV_MENU` `facout` (migrasi inti baru, bukan 900) dan barisnya di `frontend/katalogKorpus.ts`;
3. pembaca uji yang menerima modul **tanpa** folder korpus (`sumber.ts`, dan uji "dua puluh"
   `Beranda.test.ts`, `Shell.test.ts`, `daftar.*.test.ts`).

Setelah disetujui: pindahkan folder itu kembali ke `APP_RNM/modul/facout/`.

## Keputusan work owner — 1 Oktober 2026 (pengerjaan NB-01)

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 4 | Seam uji NB-01 | Satu pintu `premium.Calculate` (`modul/nbfacin/backend/services/premium`) + uji gagal-kompilasi. Parser dan resolver diuji lewat pintu itu |
| 5 | Mode pembulatan `@Math.divide` belum terverifikasi | ⚠️ **Diganti butir 14.** Hitung eksak; bila pembulatan membuang digit bukan-nol → galat `ErrModePembulatanBelumTerverifikasi`, bukan tebakan |
| 6 | Tidak ada kasus PA di lima berkas kasus | Tes memakai **contoh hitung manual**, ditandai bukan rekonsiliasi; kriteria rekonsiliasi tetap terbuka (`PERTANYAAN-NB01.md` butir 1) |
| 7 | Satuan ‰/% dan perkalian uang × rasio | Lokal di `nbfacin` (tipe `rasio`, metode `kali`); `inti/backend/uang` tidak diubah |
| 8 | Langkah 7 `CalculatePremiPA_FacIn` (L1146, `pyStepsPreCondition=false`) menimpa langkah 6 (L1003) | ⚠️ **Diganti butir 13.** **Port apa adanya** menurut P-11: premi akhir PA = L1146 (`PERTANYAAN-NB01.md` butir 3) |
| 9 | Konversi teks → angka di lapisan services (temuan review, ADR-U-0022) | Tetap di awal `Calculate` sebagai **batas masukan sementara**; pindah ke pemuat `repository` saat pemuat NB lahir |
| 10 | Medan angka kosong (temuan review, ADR-U-0022) | ⚠️ **Untuk Discount diganti butir 17.** Kosong = kosong, bukan nol: uang kosong → `uang.ErrUangKosong`, rasio kosong → `ErrRasioKosong` (`PERTANYAAN-NB01.md` butir 4) |
| 11 | Cakupan parser koma K-027 (temuan review) | ⚠️ **Diganti butir 15.** **Hanya `.Rate`** (bukti K-027). TSI, ProRatePercent, Discount lewat `utils.ParseDecimal` (titik), berlabel `[dugaan]` |
| 12 | Parser K-027 dan uji gagal-kompilasi tipe `inti` | Tetap di `nbfacin` sekarang; diusulkan ke tim inti (`USULAN-PR-TIM-INTI-NB01.md`), salinan lokal dihapus setelah diterima |

**Bukan keputusan work owner, dicatat supaya terlihat (`services/premium`):** `ErrMelampauiPresisi` — hasil antara yang
melampaui presisi 38 digit `utils.DecimalPrecision` ditolak, bukan dibulatkan diam-diam oleh `apd`.
Pengaman teknis sistem baru, tidak berasal dari perilaku Pega.

## Keputusan work owner — 1 Oktober 2026, jawaban `PERTANYAAN-NB01.md`

Dasar butir 13–17: **empat kasus NB lini PA nyata** di `D:\migrasi\RNM\DDL\CONTOH\` —
`NB-90996 (PA).xml` (metode 1) dan `NB-91017/91018/91019 (PA).xml` (metode 3). Hanya lima medan yang
dibaca dan dipakai: `TSI`, `Rate`, `CalculateMethod_FacIn`, `Premium` (coverage
`PersonList/ASMCoverage`) dan `ProRatePercent` (`pagedata`). `[terverifikasi]` Premi sistem lama
dihitung ulang dua cara — Python `decimal` dan kode Go `apd` — dan cocok **nol selisih**.

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 13 | Langkah 7 (L1146) | **Ikuti data: langkah 7 tidak dijalankan.** Metode `'1'` → langkah 4 L713, `'2'` → L858 (belum diport, tiket 04), `'3'` → langkah 6 L1003. `[terverifikasi]` metode 1 cocok eksak dengan L713 dan meleset 0,0048 dari L1146; tiga kasus metode 3 cocok eksak dengan L1003. Mengganti butir 8 dan jawaban awal butir 3 ("hanya saat `==1`"). ⚠️ Bertentangan dengan P-11 — `PERTANYAAN-NB01.md` butir 5 |
| 14 | Mode pembulatan | **Tolak hanya saat seri.** `[terverifikasi]` pemotongan tersingkir (`.61275755` → `.6128`). Setengah-ke-atas dan setengah-ke-genap keduanya dihitung; bila berbeda (digit dibuang tepat setengah) → `ErrPembulatanSeri`. Mengganti butir 5. ⚠️ **Diganti A37** (01-10-2026, menunggu konfirmasi): data kini membedakan keduanya — setengah-ke-atas |
| 15 | Format angka masukan `Calculate` | **Titik desimal untuk semua medan** (`utils.ParseDecimal`), sesuai data kerja Pega. Koma K-027 adalah format kolom Oracle `FACINOFFER.RATE` → tugas pemuat `repository` kelak (ADR-U-0022). Mengganti butir 11 |
| 16 | Angka kasus nyata di tes | **Boleh, angka saja**, diambil dengan daftar-izin lima medan di atas, tanpa nomor kasus/nama/alamat (`rekonsiliasi_test.go`). Berkas mentah tidak masuk repositori |
| 17 | `.Discount` kosong | **Dikurangi nol.** `[terverifikasi]` keempat kasus tidak memuat tag `Discount`, dan premi sistem lama cocok dengan rumus berdiskon nol. Mengganti butir 10 untuk Discount; TSI, Rate, ProRatePercent kosong tetap ditolak |

## Keputusan work owner — 1 Oktober 2026, jawaban `PERTANYAAN-NB01.md` butir 5

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 18 | Langkah ber-`pyStepsPreCondition=false` yang terbukti tidak berpengaruh (langkah 7 `CalculatePremiPA_FacIn`) bertentangan dengan P-11 ("tetap jalan") | Jawaban work owner, dikutip: *"itu hanya berlaku di PA saja"*. **P-11 tetap berlaku umum**; pengecualiannya hanya untuk PA |

**Cakupan yang sudah pasti:** langkah 7 `CalculatePremiPA_FacIn` (dibuktikan data, butir 13).

**Cakupan yang belum diputuskan:** apakah pengecualian itu juga berlaku untuk langkah bertanda sama di
rule PA lain. Sensus dua cara sepakat — **21 langkah di 8 berkas** `NB FacIn\Activity\` yang namanya
memuat PA, 2 di tingkat atas (keduanya `CalculatePremiPA_FacIn`), 19 sub-langkah bersarang:
`AddCurrencyListPA_ACT` 1 · `CalculatePremiPA_FacIn` 2 · `CountGPWMarinePAMbu_Act` 1 ·
`SaveFacinProdCurrPA_Act` 6 · `SaveFacinProdEDMPA_Act` 8 · `SetDataFacOutPATravel_Act` 1 ·
`SumFacOutPA_Act` 1 · `SumTSIPremiSpreadedRNM_PA_Act` 1.

- Cara 1: `grep -c "<pyStepsPreCondition>false</pyStepsPreCondition>"` per berkas.
- Cara 2: pengurai XML, semua kedalaman `pySteps/rowdata`.
- ⚠️ Ralat sensus: pengurai versi pertama hanya membaca tingkat atas dan memberi 1 berkas / 2 langkah
  (jebakan #5, salah wadah). Diperbaiki sebelum dicatat.
- `[dugaan]` Pemilihan "rule PA" memakai nama berkas; dua di antaranya juga menyangkut MBU dan Fac Out.

Ditanyakan saat tiket yang memakai rule itu dikerjakan, bukan sekarang — tidak satu pun dipakai NB-01.

## Keputusan work owner — 1 Oktober 2026, tiket NB-08 (registry `When`)

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 19 | Cara mengisi registry | **Dibangkitkan dari korpus** oleh `modul/nbfacin/backend/services/rules/bangkit` (jalur korpus lewat flag) menjadi `registry_gen.go`, satu entri per predikat beserta asal rule-nya |
| 20 | `compareTwoValues` tanpa tipe properti (nol rule Property) | **Tolak bila tafsir berbeda**: teks selalu, angka bila kedua sisi angka; beda → `ErrTafsirBerbeda`. Kosong vs ANGKA (literal angka atau properti bernilai angka) → `ErrTafsirBerbeda`. Literal teks berkutip (`"10053"`) bertipe teks, jadi kosong dibandingkan dengannya sebagai teks |
| 21 | Bentuk seam | `rules.Eval(nama, kasus) (bool, error)`: galat untuk keraguan data, **panic** untuk kesalahan program (nama tak dikenal, IsPKSASM, bentuk belum diport, rujukan melingkar) |
| 22 | Rule berbentuk kondisi lain — 11 (platform/CRM, `IsDeducType`, `IsTypeDeductType`, `IsShowInput`, `IsPASSG`, `StepStatusFail`) — dan rujukan `pyIsIPad` | Terdaftar dengan asal, **panic "belum diport"**; diport saat tiket yang memakainya datang |
| 23 | `!` dan `&&` di `pyLogic` | **`!` = NOT, `&&` = AND**, `[dugaan]` (dikuatkan nama `IsNotPAandNotMBU`); `\|\|` = OR ikut didukung walau tidak ada di korpus |
| 24 | ~~Diganti butir 28.~~ Operand kanan tanpa kutip (`= ReasFacInDirector`, `= A1`…`F6`, `= Offer`/`Policy`, `= True`) — 43 baris di 6 predikat (`IsUW`, `IsLimitSBondKBG`, `IsLimitCreditCL`, `IsVisible`, `IsCedingConfirmOffer`, `IsCedingConfirmPolicy`); ditemukan review spec | **Panic "belum diport"** sampai satu contoh dicek di UI Pega (`PERTANYAAN-NB08.md`). Sebelumnya dibaca sebagai properti; properti tidak ada = kosong membuat `"" = ""` benar dan `IsUW` terbuka bagi operator tanpa workbasket |
| 25 | Urutan dan hubung-singkat evaluasi baris | **Semua baris yang dirujuk dinilai, urut label** (A, B, …, Z, AA): hasil dan galat deterministik, tidak bergantung urutan evaluasi Pega |

**Angka yang dipakai — dihitung dua cara:**

- `[terverifikasi]` "601 rule `When`" di tiket = **jumlah berkas tiga folder** (NB 210 + RNW 189 +
  EDM 202). Identitas unik gabungan 273. Modul ini: **NB, 210 berkas, 209 identitas** — satu identitas
  (`ISFLAGOLDDATA`) ada di dua berkas yang **identik** (hash ternormalisasi 23-tag + `pzIndexes`, dan
  isi kondisinya). Angka "226 predikat" di spec Modul 3 **belum terjelaskan**; tidak dipakai.
- `[terverifikasi]` Registry bangkitan dicocokkan dengan pembacaan independen korpus (Python):
  nama 209/209, logika identik, jumlah baris kondisi identik, predikat panic 19/19 (13 sebelum butir 24).
- ⚠️ **Ralat sensus:** sensus pertama atas `pyLogic` mencacah kata huruf saja dan melaporkan "hanya
  AND/OR dan kurung". Pembangkit menemukan `!` dan `&&`. Dihitung ulang dua cara (cara 1: berkas yang
  `pyLogic`-nya memuat karakter selain huruf/angka/spasi/kurung; cara 2: cacah karakter): `(`/`)` 16,
  `!` 4, `&` 2 — di `IsNonPropertyandNonEngineering`, `IsNotPAandNotMBU`, `isTravelTime`.
- ⚠️ **Ralat alat uji mutasi:** satu mutasi (deteksi rujukan melingkar dibuang) terlapor "lolos"
  karena skrip hanya mencari baris `--- FAIL`; sebenarnya proses tes mati dengan `stack overflow`
  (kode keluar 1). Putaran kedua (14 mutasi, kode keluar proses dihitung): seluruhnya tertangkap, setelah dua tes diperketat agar memeriksa ALASAN panic dan kata `NOT`, bukan hanya ada-tidaknya penolakan.

## Keputusan work owner — 1 Oktober 2026, tiket NB-03 (resolver lini bisnis)

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 26 | Seam untuk resolver lini bisnis → satuan rate | **Seam kedua: `premium.SatuanRate(lini) Satuan`**, diuji table-driven untuk 8 lini K-018 (kini 7, butir 30) + lini tak dikenal (panic). `Calculate` memakainya; pembagi komposit 100000 tetap diuji lewat `Calculate` (PA L713). Kelak disediakan ke RNW/EDM/Fac Out lewat kontrak |
| 27 | `ProRatePercent` diberi satuan di luar resolver (temuan review NB-03) | **Fungsi pengisi tersendiri `SatuanProRata()`** di `resolver.go` (K-018: `ProRatePercent (persen) \| 100`, faktor yang tidak bergantung lini). `rasio` hanya dibangun lewat `rasioRate` / `rasioProRata`; dijaga `TestRasioHanyaDibangunDiResolver` |

**Dasar bukti per lini — dicatat supaya label tidak melebihi buktinya:** peta 8 lini adalah keputusan
work owner K-018 ("DIKUNCI"). Bukti rumus yang dikutip K-018 hanya untuk **PA, MBU, Layering**
(`[terverifikasi]`). FIRE, ANEKA, BONDING, GOLF, MARINE CARGO bersandar pada keputusan itu; pembacaan
sekilas `NB FacIn\Activity` (01-10-2026) memberi penguat `[dugaan]`: `FillPremiGolf` ÷10000 (rate % ×
`PctShortPeriod` %), `CountCoverageMarine` ÷100, `CopyAllObjFacOutFireAneka_ACT` memuat pembagi 100000
dan 10000 (cabangnya belum dibaca). BONDING: belum ditemukan rumus pembanding.

`LiniBisnis` kini tipe sendiri (temuan review NB-01, *Primitive Obsession*); nilainya label K-018. Kode
lini di data produksi tetap `[pertanyaan terbuka]` sampai tiket 10.

## Keputusan agent A1–A12 — DIKONFIRMASI work owner 1 Oktober 2026 (butir 34)

Semula pilihan agent saat work owner meminta seluruh tiket dikerjakan sekaligus. **Dikonfirmasi** dengan
tindak lanjut di butir 34; A2 dan A6 disusul butir 29 dan 30.

| # | Tiket | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A1 | NB-10 | `premium.LiniDariPredikat` menjadi seam ketiga paket premium dan mengembalikan **daftar** lini, bukan satu | `[terverifikasi]` gerbang pemilih rumus di `CountGrossPremi_Act` 5.1–5.4 dan `CountGPWMarinePAMbu_Act` 1.1–1.3 berdiri sendiri |
| A2 | NB-10 | Gerbang `StatusBusiness != 3` di atas langkah 5/1 **tidak** diport di jembatan | milik alur pemanggil; artinya `[pertanyaan terbuka]` |
| A3 | NB-04 | Langkah 2 (diskon persen) diport dengan masukan baru `PremiSebelumnya` (`.Premium` lama); langkah 3 tidak | langkah 2 mengubah `.Discount` yang dikurangkan premi; langkah 3 hanya menulis `.DiscountPercentage` |
| A4 | NB-05 | Hanya L1144 yang diport; `.MinPremium` dan diskon MBU tidak | `[terverifikasi]` L1626 adalah syarat `When`, bukan rumus; kedua langkah lain tidak menulis `.Premium`. ⚠️ Ini **mengoreksi kutipan K-018** (berstatus DIKUNCI, mengutip L1626 sebagai rumus MBU) dan kriteria "kedua bentuk" tiket 05 — **butuh konfirmasi work owner** |
| A5 | NB-05 | `ProRatePercentCoverage` medan terpisah dari `ProRatePercent` PA | L1144 membaca `.ProRatePercent` coverage, L713 membaca `pyWorkPage.OfferFacIn.ProRatePercent` |
| A6 | NB-06 | `param.Layer` wajib bilangan bulat ≥ 0; 0 lapis → premi 0 | loop 0…Layer−1 tidak berjalan, `.Premium` tetap 0 dari langkah 2 |
| A7 | NB-15 | Nilai kunci BUANG **dikosongkan**, kuncinya tetap | tiket: "hanya nilai kunci identitas yang hilang, bukan bentuknya" |
| A8 | NB-15 | Fixture **tidak** disimpan di repositori sampai daftar BUANG diperluas | ditemukan 9 medan bocor (`PERTANYAAN-LANJUTAN.md` butir 4); CLAUDE.md §10 |
| A9 | NB-02 | Tiket **diusulkan** ditutup oleh butir 3 (penutupan wewenang work owner); dua kewajiban (hitungan tanpa mata uang, gagal keras tulis Oracle) tetap tercatat untuk pemuat/jalur tulis | butir 3 menggantikan desain `Unknown`, tetapi penanganan uang tanpa mata uang tetap `[pertanyaan terbuka]` |
| A10 | NB-13 | Efek samping Reject (`ConfirmBinding = 0`, `ReceivedRiSlip = false`) **juga** diterapkan pada Revise | `[terverifikasi]` `SetReviseProposal.xml` L179/L209 menulis keduanya; K-014 hanya menyebut Reject — diport apa adanya |
| A11 | NB-10 | Kasus Bonding menghasilkan lini **ANEKA**; `LiniBonding` tidak pernah dihasilkan jembatan | `[terverifikasi]` tidak ada gerbang Bonding sendiri; `IsAneka` baris AA merujuk `IsBondingAndCustomBonds`. Satuan sama (%), tetapi tiket 10 menyebut predikat "mengenali BONDING" |
| A12 | NB-05 | `.Loading` kosong = 0 | `[terverifikasi]` 93 baris MBU nyata tanpa Loading cocok eksak; tidak ada keputusan work owner seperti butir 17 untuk Discount |

## Keputusan work owner — 1 Oktober 2026, jawaban `PERTANYAAN-NB08.md`, `-LANJUTAN.md`, `-AKSEPTASI.md`

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 28 | Operand kanan tanpa kutip (butir 24) | **Dibaca sebagai TEKS** ("operand isuw itu harusnya pakai kutip"). Mengganti butir 24. Di pembangkit: satu kata (`^[A-Za-z_]\w*$`) → `teksLiteral`; operand bertitik tetap ditolak ("bukan literal"). Pengecualian agent A13 dan A14 di bawah |
| 29 | Arti `StatusBusiness = 3` | **Endorsement.** Gerbang `StatusBusiness != 3` di atas `CountGrossPremi_Act` langkah 5 / `CountGPWMarinePAMbu_Act` langkah 1 berarti premi endorsement tidak lewat pemilih rumus NB. Gerbang itu tetap milik alur pemanggil (A2) |
| 30 | Layering | **Tidak dipakai lagi di sistem baru.** Peta K-018 menjadi **7 lini**; `LiniLayering`, `layering.go`, `JumlahLapisan`, `ErrJumlahLapisan` dihapus; `SatuanRate("Layering")` panic seperti lini di luar peta. NB-06 ditutup; bagian Layering NB-07 dicabut; A6 gugur; permintaan kasus Layering (`PERTANYAAN-LANJUTAN.md` butir 1) dicabut |
| 31 | Ejaan `JABATAN` bentuk A (`PERTANYAAN-AKSEPTASI.md` B2) | **Isi korpus sudah sesuai: ada yang berspasi dan ada yang tanpa spasi.** Cabang berspasi (`"DIREKTUR TEKNIK"` L4635) dan tanpa spasi (`"KADIVFACULTATIVE"` L4929) keduanya hidup — **bukan kode mati**; K-023 ("tanpa spasi") tidak berlaku umum. Pencocokan diport persis per literal |
| 32 | Sembilan medan BUANG tambahan (`PERTANYAAN-LANJUTAN.md` butir 4) | **Setuju.** `pxCreateOperator`, `AccumulationDescription`, `AccumulationCode`, `PIC`, `PICSuggest`, `CommentSuggest`, `PropertiItemNote`, `SobName`, `TopRiskLocation` masuk `DaftarBuang`. Fixture disimpan di repositori **hanya setelah nol kebocoran di kelima berkas P-5 dan tinjauan daftar kandidat**. `pxUpdateOperator` / `pxUpdateOpName` tidak ditambah (A15) |
| 33 | Guard `LOGIN` (`PERTANYAAN-AKSEPTASI.md` B3) | **Setuju: limit pengguna menjadi MASUKAN tangga** dari model peran, bukan dicari lewat `LOGIN = {OperatorID.pyUserIdentifier}`. Pemetaan login → jabatan dipegang IAM |
| 34 | Konfirmasi keputusan agent A1–A12 | **Setuju, dengan tindak lanjut:** A4 — teks K-018 di `00-KEPUTUSAN-WORK-OWNER.md` diamandemen (L1626 syarat, bukan rumus); A9 — NB-02 ditutup, kewajiban "hitungan uang tanpa mata uang" pindah ke `issues/17-hitungan-uang-tanpa-mata-uang.md`; A10 — tetap dikonfirmasi dengan Underwriting (memperluas K-014) |
| 35 | Permintaan ke DBA | **Sesuai rekomendasi:** satu query hitungan produksi (PA `CalculateMethod_FacIn = 2`; `DiscountType = Percent`; MBU `Loading` terisi; pro-rata coverage ≠ 100) — hitungan nol = jalur tidak terpakai; ekspor CSV tabel `M_LIMIT_*` bentuk A dan `M_LIMIT_FINANCIALINS` tanpa `NAMA` dan `LOGIN`. Tindakan pihak luar; belum ada hasil |

**Hasil penerapan, diukur 01-10-2026** — tiap angka dengan cara hitungnya:

- Predikat panic statis di `registry_gen.go`: **19 → 17**. Cara 1: `grep -c '^\t\tpanik:'` = 17. Cara 2:
  awk per entri, mencetak nama predikat = 17 nama. Rekonsiliasi delta terhadap angka lama (L119): 19 − 5
  (`IsUW`, `IsLimitSBondKBG`, `IsLimitCreditCL`, `IsCedingConfirmOffer`, `IsCedingConfirmPolicy` kini
  teks) + 3 (`IsGroup`, `IsGroupCreate`, `IsSPVCreate`, A14) = 17 — sepakat. `IsVisible` (A13) dan
  `pyIsIpadOrDesktop` (merujuk `pyIsIPad` yang tidak ada) tetap panic.
- De-identifikasi atas kelima `DDL\P-5 *.txt` dengan 21 kunci: kebocoran **2 / 2 / 1 / 0 / 6** (EDM-13445,
  NB-176005, NB-181231, NB-184233, RNW-10579). Sisa bocor ada di **dua medan baru**: `Comment` (5 kebocoran, 3 jalur berbeda)
  dan `OperatorID` (6 kebocoran, 3 jalur berbeda, di bawah `QuotationData` / `OldData`); 5 + 6 = 11 = jumlah per berkas. Karena itu fixture **belum** disimpan
  (butir 32); keduanya menunggu keputusan work owner (tiket 15 melarang filter pola). Daftar kandidat:
  139 nama kunci unik, untuk tinjauan; tidak disalin ke sini.

## Keputusan agent A13–A15 — DIKONFIRMASI work owner 1 Oktober 2026 (butir 37)

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A13 | 28 | `= True` berhuruf besar (`IsVisible`) **tetap ditolak**, tidak dibaca sebagai teks `"True"` | `booleanLiteral` membandingkan teks huruf kecil `"true"`; bila Pega membaca `True` sebagai boolean, teks `"True"` membuat `IsVisible` selalu salah. Satu predikat, satu baris |
| A14 | 28 | Baris yang membandingkan **login operator** (`.pyUserIdentifier`, `.pxCreateOperator`, `.pxUpdateOperator`) dengan literal ditolak, dan **nilainya tidak ditulis** ke registry maupun pesan panic | CLAUDE.md §4 butir 10 (nama orang tidak disalin). Ditemukan saat regenerasi: lima nilai login sudah tertulis di `registry_gen.go` sejak NB-08 sebagai literal berkutip — kini nol (dicek dua cara: sensus `registry_gen.go` dan pencarian nilai dari korpus ke seluruh pohon `nbfacin`). Ketiga predikat (`IsGroup`, `IsGroupCreate`, `IsSPVCreate`) diport kelak lewat model peran (butir 33) |
| A15 | 32 | `pxUpdateOperator` / `pxUpdateOpName` tidak ditambahkan | rekomendasi menyebut "bila muncul"; 0 kemunculan di kelima berkas P-5 (`grep -o '"kunci" *:'` per berkas) |

## Keputusan work owner — 1 Oktober 2026, jawaban `PERTANYAAN-LANJUTAN.md` butir 7–8

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 36 | Dua medan bocor sisa | **Setuju rekomendasi:** `Comment` dan `OperatorID` masuk `DaftarBuang`; tujuh nama kunci yang mencurigakan ditinjau sekaligus |
| 37 | Keputusan agent A13–A15 | **Dikonfirmasi** |

**Hasil penerapan butir 36, diukur 01-10-2026:**

- `DaftarBuang` 23 kunci (`TestDaftarBuangDisetujui`). Dijalankan ulang atas kelima `DDL\P-5 *.txt`:
  kebocoran **0 / 0 / 0 / 0 / 0**; alat menulis kelima berkas — ke scratchpad, di luar repositori.
- ⚠️ Nol kebocoran hanya berarti nilai yang SUDAH dibuang tidak muncul lagi; nama yang hanya ada di satu
  medan tidak terdeteksi alat. Karena itu ketujuh kunci ditinjau lewat profil bentuk (huruf besar → X,
  kecil → x, angka → 9), tanpa mencetak nilai. Uji instrumen: `InsuredName` (sudah dibuang) terisi 0,
  `BusinessType` (kode) terisi 8 — sesuai harapan.
  - `OperationalDirector`, `OperationalDivHead`, `PresidentDirector`, `TechnicalDirector` — **hanya
    `true`/`false`** (5 nilai masing-masing). Flag persetujuan, bukan nama.
  - `Name` — 135 nilai, semua tiga huruf besar, 106 di bawah `Currency` + 29 di daftar mata uang/total;
    **0 dari 135** sama dengan login/nama operator yang dikenal (34 nilai dari `pxCreateOperator`,
    `OperatorID`, `pxCreateOpName`, `PIC`, `PICSuggest` berkas asli). Kode mata uang.
  - `GroupName` — hanya `-`. `SobLeader0` — kode berbentuk `X9`.
  - **Kesimpulan:** tidak satu pun dari ketujuhnya berisi identitas; tidak ditambahkan.
- ⏸ **Fixture belum disimpan di repositori** — lihat "Yang belum diputuskan".

## Keputusan work owner — 1 Oktober 2026, bentuk fixture de-identifikasi

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 38 | Fixture masih memuat nomor polis/kasus dan teks bebas | **Kosongkan dulu**: kunci nomor polis/kasus dan teks bebas ditambahkan, nama berkas dibuat netral, nol kebocoran diperiksa ulang, lalu disimpan di `backend/services/premium/testdata/kasus/` |

**Penerapan, 01-10-2026:**

- **`DaftarBuang` 33 kunci** — identitas + nomor polis/kasus `PolicyNo`, `OldPolicyNo`, `EndorsementNo`,
  `InvoiceNumber`, `NoOfferSlip`, `PolicyMasterNumber`, `PolicyMasterIDPega`, `IDNewBisnis`,
  `IDFollowingNB`, `Following` (`Following` ditemukan memuat pola nomor polis `RNM-F…`, 4 kali).
- **`JalurBuang`** `$.ID`, `$.OldData.ID`, `$.OldData.OldData.ID` — dicocokkan persis. `ID` tingkat
  dokumen berbentuk kunci kasus Pega atau angka; `ID` di jalur lain kode (`Currency.ID` 104, `Ship.ID` 104),
  sehingga kuncinya tidak dapat dibuang utuh.
- **`DaftarKosongkan` 19 kunci** — teks bebas + bagian alamat: dikosongkan tetapi **bukan sumber deteksi
  kebocoran** (A16).
- Hasil: kebocoran **0 / 0 / 0 / 0 / 0**. Cara kedua, skrip Python terpisah dari alat Go: nomor kasus 0,
  pola polis 0, nilai identitas asli ≥ 5 huruf sebagai substring 0 kecuali satu kata umum enam huruf dari
  `CommentSuggest` di dalam nilai `IsCedingConfirm` (bukan identitas; alat Go yang memeriksa kata utuh
  benar tidak menandainya), daun di luar ketiga daftar tidak berubah 0, jumlah daun sama. Uji instrumen
  cara kedua: dijalankan atas berkas asal → kelimanya GAGAL, sesuai harapan.
- Fixture tersimpan: `edm-fire-1`, `nb-fire-1`, `nb-kredit-1`, `nb-marinecargo-1`, `rnw-fire-1` — identik
  byte dengan hasil terverifikasi. Padanan dengan nomor kasus **tidak disimpan** di repositori. Penjaga
  `TestFixtureKasusSudahDibersihkan`; diuji atas lima berkas mentah bernama netral di luar repositori →
  gagal dengan 2.775 daun "masih berisi".

## Keputusan agent — menunggu konfirmasi (1 Oktober 2026, penerapan butir 38)

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A16 | 38 | Daftar dipisah: `DaftarBuang` (identitas, diburu di medan lain) vs `DaftarKosongkan` (teks bebas, alamat — hanya dikosongkan) | Saat semua kunci baru menjadi sumber deteksi, kebocoran naik ke 36 / 2 / 5 / 0 / 42 — seluruhnya positif palsu: `Remarks` menyalin angka TSI, `OccupationNote` menyalin nama okupasi, `CoverageNote` nama bisnis (diagnosis per pasangan medan, bentuk nilai saja) |
| A17 | 38 | `ObjectName` dan `Others` **tidak** dikosongkan | `ObjectName`: label kategori satu–dua kata yang sama dengan `ObjectType`/`ItemType`; `Others`: 10 kali objek, 5 kali hanya `true`/`false` |
| A18 | 38 | `ASMCity`, `ASMDistrict`, `ASMRW` ikut dikosongkan, walau di luar "nomor polis/kasus dan teks bebas" | bagian alamat tertanggung yang sama dengan `ASMAddress`/`ASMZipCode` (sudah BUANG); `Province`, `pyCity`, `Zone` dibiarkan (tingkat wilayah) |
| A19 | 38 | Salinan angka di `Remarks` ikut hilang | `Remarks` medan teks; angka yang diuji ada di medan TSI-nya sendiri, yang utuh |

## Jawaban work owner — 1 Oktober 2026, kode transisi (`PERTANYAAN-AKSEPTASI.md` B1)

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 39 | Arti kode transisi `pyStepsPreCondParamsWhen(True\|False)` | `[terverifikasi]` lewat UI Pega — tabel di bawah. Menjawab U-1 `_DAFTAR-ISSUE-TERBUKA.md`; kode `1`, `4`, `6` (U-2) masih terbuka |

| Kode | Tampil di UI Pega | Bukti |
| :-: | --- | --- |
| `2` | **Continue Whens** | prakondisi 1 "bila salah" **dan** prakondisi 2 "bila benar" — dua baris sepakat |
| `3` | **Skip Step** | prakondisi 2 "bila salah" |
| `5` | **Skip Whens** | prakondisi 1 "bila benar" |

Tangkapan layar UI Pega dari work owner, 1 Oktober 2026: `GetLimitAkseptasi_ActFlow` langkah 12, dicocokkan
dengan `NB FacIn\Activity\GetLimitAkseptasi_ActFlow.xml` (`ASM-FW-GISFW-WORK!GETLIMITAKSEPTASI_ACTFLOW`,
`pyRuleSetVersion` 01-01-87) — kedua baris prakondisi identik: `IsLimitSBondKBG` benar=5 salah=2,
`IsLimitCustomBond` benar=2 salah=3.

- `[dugaan kuat]` **Arti tindakannya** mengikuti perilaku platform Pega, belum dibuktikan dari korpus:
  *Continue Whens* = nilai When berikutnya, dan bila tidak ada lagi, langkah dijalankan; *Skip Whens* = When
  sisanya dilewati dan langkah **dijalankan**; *Skip Step* = langkah dilewati.
- Konsekuensi `[dugaan]` (bergantung pada arti di atas): langkah 12 berjalan bila `IsLimitSBondKBG` **atau**
  `IsLimitCustomBond` benar. Langkah 19 (`IsNonPropertyandNonEngineering` benar=5 salah=2, lalu
  `IsPropertyandEngineering` benar=2 salah=5) **berjalan di setiap hasil** — prakondisinya tidak pernah
  melewati langkah. Belum dipakai kode mana pun.
- ⚠️ Ralat agent: rekomendasi B1 menduga kode 5 "tindakan loop". Dugaan itu salah, dan sudah dilarang
  `_DAFTAR-ISSUE-TERBUKA.md` (diuji dan gagal: 9,8 % langkah berkode 5 beriterasi, kode 2 30,8 %).
  Begitu juga "kode 2 = jalankan" di catatan lama hanya benar untuk When **terakhir**.

## Jawaban work owner — 1 Oktober 2026, kode transisi 1, 4, 6 (U-2)

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 40 | Arti kode `1`, `4`, `6` | `[terverifikasi]` lewat UI Pega: **`1` = Jump To Later Step**, **`4` = Exit Iteration**, **`6` = Exit Activity**. Dugaan lama "`6` = keluar" terbukti |

**Peta kode transisi kini lengkap:**

| Kode | Tindakan di UI Pega | Kemunculan* | Sumber |
| :-: | --- | ---: | --- |
| `1` | **Jump To Later Step** | 392 | butir 40 — `CopyFacRetro_ACT` langkah 1 |
| `2` | **Continue Whens** | 48.340 | butir 39 — `GetLimitAkseptasi_ActFlow` langkah 12 |
| `3` | **Skip Step** | 17.650 | butir 39 |
| `4` | **Exit Iteration** | 122 | butir 40 — `setDiscountRetro_act` langkah 2 |
| `5` | **Skip Whens** | 1.034 | butir 39 |
| `6` | **Exit Activity** | 277 | butir 40 — `AddStatusCoverageNew_Act` langkah 1 |
| *(kosong)* | `[pertanyaan terbuka]` | 5.953 | — |

\* Jendela: seluruh `*.xml` di bawah `D:\migrasi\RNM\`, elemen `pyStepsPreCondParamsWhenTrue` dan
`…WhenFalse`. Dua cara sepakat untuk `1`–`6`: `grep -rhoE '<pyStepsPreCondParamsWhen(True|False)>[^<]*</'`
dan pengurai XML (0 berkas gagal diurai). Isian kosong hanya terhitung pengurai — `grep` melewatkannya
karena ditulis sebagai tag tutup-sendiri. Tidak ada kode lain.

- *Jump To Later Step* membawa langkah tujuan di `pyStepsPreCondParamsWhenTruePrms` / `…FalsePrms`
  ("true param" / "false param" di UI) — port wajib membacanya, bukan hanya kodenya.
- `[pertanyaan terbuka]` **Arti isian kosong** (5.953): belum terlihat di UI. Contoh terdekat:
  `AddStatusCoverageNew_Act` langkah 1, kolom "if true" (kolom "if false"-nya `6`).

## Jawaban work owner — 1 Oktober 2026, isian kode transisi kosong

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 41 | Arti isian kosong `pyStepsPreCondParamsWhen(True\|False)` (5.953 kemunculan, butir 40) | `[terverifikasi]` lewat UI Pega: **kosong = Continue Whens**, sama dengan kode `2`. Peta kode transisi kini **tanpa sisa**; port wajib memperlakukan kosong dan `2` sebagai satu tindakan |

## Data dari work owner — 1 Oktober 2026, CSV tabel limit (`PERTANYAAN-AKSEPTASI.md` B4)

| # | Hal | Isi |
| ---: | --- | --- |
| 42 | Ekspor `POOLDATA.M_LIMIT_*` | Enam CSV di `DDL\` → disalin byte demi byte ke `backend/services/acceptance/testdata/limit/`. Bentuk A 37 baris × 9 kolom (kelima tabel juga memuat `WORKBASKET`/`JABATAN_ATASAN`); `M_LIMIT_FINANCIALINS` 4 baris × 4 kolom |

- **Format** `[terverifikasi]`: UTF-8, CRLF, pemisah `;`, tanpa `NAMA`/`LOGIN`, setiap baris selebar header.
  Angka bertitik ribuan: 541 nilai bertitik semuanya ≥ 2 titik berpola ribuan, 581 bulat polos, nol bentuk
  lain (pola `\d{1,3}(\.\d{3})+` dan `\d+` atas sembilan kolom angka) — tidak ada desimal.
- **`JABATAN`** `[terverifikasi]`: lima tabel bentuk A — 10 nilai, semuanya **tanpa spasi**
  (`DEPHEADUNDERWRITER`, `DIREKTURMARKETING`, `DIREKTURTEKNIK`, `JUW_B`, `KADIVFACULTATIVE`, `KADIVTEKNIK`,
  `LEADER`, `MANAGERTEKNIK`, `SENIORUW`, `UNDERWRITER`); `M_LIMIT_FINANCIALINS` — 4 nilai, semuanya
  **berspasi** (`DIREKTUR MARKETING`, `DIREKTUR TEKNIK`, `KADIV KEUANGAN`, `SENIOR UNDERWRITER`).
- **Letak literal berspasi** di `GetLimitAkseptasi_ActFlow` (`.CARI1 == "…"`, pengurai XML per langkah):
  17.1, 19.1.2, 19.1.8, 19.1.9, 21.1.7, 21.1.8, 22.2.2–22.2.4. Hanya **langkah 22** ("Limit Bond & KBG &
  Kredit & Trade") yang membaca daftar FINANCIALINS — cocok dengan data. Langkah **19 dan 21 berlabel `//`**
  (label = `pyStepsBlockName`, terbukti dari tangkapan layar langkah 12 berlabel `NON1`). Langkah **17**
  ("hapus …") bergerbang satu nomor kasus literal `pyWorkPage.pyID` — tambalan untuk satu kasus.
- **Label `//`** — 11 langkah (pengurai dan `grep -c '<pyStepsBlockName>//</'` sepakat): tingkat atas
  **15, 16, 19, 21**, dan 19.1.2/.3/.5/.6/.7, 21.1.2/.4. `[dugaan kuat]` `//` = langkah dinonaktifkan
  (konvensi Pega; U-6 masih terbuka).

⚠️ **Ralat agent atas catatan butir 31.** Catatan "cabang berspasi dan tanpa spasi keduanya hidup — bukan
kode mati" adalah perluasan agent, bukan kata work owner, dan **data membantahnya untuk bentuk A**: literal
berspasi yang cocok dengan data hanya di langkah 22 (FINANCIALINS). Yang di langkah 19 dan 21 berada di
langkah berlabel `//`; yang di 17.1 tidak pernah cocok dengan `JABATAN` bentuk A. Jawaban work owner sendiri
("ada yang pakai spasi dan tanpa spasi") **benar** — pada tingkat tabel.

⚠️ **Konsekuensi bagi NB-11, bila `//` = nonaktif:** blok 19 tidak berjalan, sehingga keputusan tangga
bentuk A hanya **langkah 20** — `LetterNo = LimitAkseptasi.pxResults(1).CARI1`; ke-15 tautologi di
19.1.x / 21.1.x dan kode `5` di langkah 19 tidak relevan. Belum dipakai kode mana pun sampai U-6 dikonfirmasi.

## Jawaban work owner — 1 Oktober 2026, label `//` (`PERTANYAAN-AKSEPTASI.md` H1)

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 43 | Arti label langkah `//` (`pyStepsBlockName`) | **"Tandanya di-remark, tidak dipakai lagi."** Langkah berlabel `//` — beserta seluruh sub-langkahnya — tidak diport. Berlaku umum, bukan hanya `GetLimitAkseptasi_ActFlow` |

**Konsekuensi yang langsung berlaku:**

- **NB-11:** di `GetLimitAkseptasi_ActFlow` langkah **15, 16, 19, 21** tidak diport (11 langkah berlabel
  `//`, butir 42). Keputusan tangga bentuk A = **langkah 20** (`LetterNo = LimitAkseptasi.pxResults(1).CARI1`)
  atas daftar langkah 6–10. Ikut gugur: 15 tautologi di 19.1.x/21.1.x, literal berspasi di 19/21, dan
  catatan butir 39 "langkah 19 berjalan di setiap hasil".
- **Premi PA (sudah diport):** `[terverifikasi]` di keempat activity premi yang diport
  (`CalculatePremiPA_FacIn`, `FillPremiMBU_FacIn`, `CountGrossPremi_Act`, `CountGPWMarinePAMbu_Act`),
  satu-satunya langkah berlabel `//` adalah **langkah 7 `CalculatePremiPA_FacIn` (L1146)** — pengurai XML dan
  `grep -c '<pyStepsBlockName>//</'` sepakat (1 / 0 / 0 / 0). Itulah sebab langkah 7 tidak berpengaruh
  (butir 13); konsisten dengan butir 18. Tidak ada perubahan perilaku, hanya komentar `premium.go`.
- **U-6** `_DAFTAR-ISSUE-TERBUKA.md` (remark `//` pada pemanggil `GetLimitAkseptasi_Act`/`_Act2`): arti `//`
  kini terjawab; pemanggil mana yang ber-`//` masih harus dibaca per pemanggil.

## Keputusan work owner — 1 Oktober 2026, tiket NB-11 (tangga akseptasi bentuk A)

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 44 | Rancangan tangga | **(a) Seam** `acceptance.Next(kasus rules.Kasus, tabel TabelLimit, pengguna Pengguna) (Transisi, error)`: membaca properti Pega yang sama dan menilai `IsFire`, `IsEngineering`, `IsNonPropertyandNonEngineering`, `IsLimit*`, `To*` lewat `rules.Eval`. **(b) Limit pengguna** dari **jabatan** pengguna (model peran), dicari di tabel dan team group yang sama; baris ganda diterima hanya bila limitnya identik. **(c) Jabatan pengguna tidak ada** di tabel/team group → **ikuti sistem lama**: tidak ada kandidat → tangga selesai. **(d) Tiga keadaan tak pasti → galat**: kandidat pertama seri antar-jabatan berbeda (`ErrUrutanTakPasti`), `IsFire` dan `IsEngineering` sama-sama benar (`ErrDaftarGanda`), tidak ada daftar dimuat (`ErrTabelTakTerpilih`) |
| 45 | Tiga temuan tinjauan kode NB-11 | **Langkah 17 diabaikan di NB-11** — "DIREKTUR TEKNIK" berspasi tidak ada di tabel bentuk A; empat nomor polis literal dibahas di NB-12. **Langkah 1** `CountTotalTSIPremiNusaRe_Act` = **tanggung jawab pemanggil** (kasus wajib berisi total TSI); port-nya tiket tersendiri. **Ketiga field state** cukup terpisah lewat tipe — `Next` tidak menerima/mengembalikan `Keputusan` |

**Fakta yang ditemukan saat mengerjakan NB-11** `[terverifikasi]`:

- `ToKadivFacultative` (Transition134 → `Assignment8`) **tidak** menulis `PositionNote`, tetapi assignment
  tujuannya merutekan ke workbasket `ReasFacInFacultativeDivHead` (parameter `Workbasket`). Ketujuh assignment
  tujuan Decision23 dibaca; pada enam lainnya workbasket = `PositionNote` yang ditulis. Ralat atas catatan
  lama "tidak menulis antrean": kasus tetap berpindah antrean, hanya `PositionNote` yang basi.
- Decision23 tidak punya konektor `SENIORUW`/`UNDERWRITER`/`JUW_*`: jabatan tujuan itu jatuh ke Else = selesai.
- `TeamGroup = 5` ada di 40 dari 115 kasus `DDL\CONTOH`; tabel limit hanya team group 1–4 → tangga selesai
  tanpa kandidat (butir 44 c).
- Seri `LIMIT_BOTTOM` KADIVTEKNIK / MANAGERTEKNIK / KADIVFACULTATIVE ada di setiap team group; langkah 18
  menghilangkannya di jalur biasa, tidak di jalur banding + reject di atas LIMIT_BOTTOM2 KADIVTEKNIK.

## Keputusan agent — menunggu konfirmasi (1 Oktober 2026, NB-11)

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A20 | 44 | `Pengguna.AnggotaGrup` menggantikan predikat `IsGroup` (langkah 3) | `IsGroup` membandingkan login literal (A14); butir 33 memindahkan identitas ke model peran |
| A21 | 44 | Langkah 11 (varian Banding) **mengganti** daftar langkah 6–10 | `[dugaan]` RDB-List ke BrowsePage yang sama mengganti isinya; deskripsi langkah "UNTUK AMBIL LIMIT KADIV FACULTATIVE KALO BANDING" |
| A23 | 44 | Langkah 18 membuang **semua** baris yang cocok | `[dugaan]` Property-Remove di dalam loop atas pxResults tidak melompati baris berikutnya |
| A24 | 44 | `team_group` dibandingkan sebagai teks persis | `[dugaan]` kolom `VARCHAR2(20)` (`DDL\M_LIMIT_PROPERTYY.txt`); `{TeamGroup}` disisipkan tanpa kutip |
| A25 | 44 | Daftar kosong → `LetterNo = ""` | `[dugaan]` `pxResults(1).CARI1` atas daftar kosong; langkah 2 sudah menyetel `LetterNo = ""` |

*(A22 tidak dipakai: digantikan butir 45.)*

## Keputusan work owner — 1 Oktober 2026, tiket NB-12 (bentuk B, financial)

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 46 | K-026/tiket 12 ("tidak bereskalasi, keluaran daftar jabatan") bertentangan dengan langkah 22 `GetLimitAkseptasi_ActFlow` / langkah 23 `GetLimitAkseptasi_Act` | **Ikuti korpus.** Port langkah 22 sebagai jalur TERPISAH (`acceptance.NextFinancial`): filter `WHERE <kolom> <= {CARID2}` dipertahankan, lalu satu langkah naik menurut antrean saat ini, lalu Decision23. K-026 dan tiket 12 diamandemen |
| 47 | Langkah 17 (buang "DIREKTUR TEKNIK" untuk empat nomor polis literal) | **"Di sistem baru nanti itu dibuang saja, harusnya itu tidak dipakai lagi."** Tidak diport, untuk bentuk A maupun B; nomor polis tidak masuk kode. Menggantikan sikap butir 45 untuk langkah 17 |

**Tangga bentuk B yang diport** `[terverifikasi]` (langkah 22.2.2–22.2.4; `_Act` langkah 23 identik isinya):

| Antrean saat ini | Baris di daftar (JABATAN berspasi) | LetterNo → antrean (Decision23) |
| --- | --- | --- |
| `ReasFacInUnderwritingFinancial` | `KADIV KEUANGAN` | `KADIVFINANCIAL` → `ReasFacInFinDivHead` |
| `ReasFacInFinDivHead` | `DIREKTUR MARKETING` | `DIREKTURMARKETING` → `ReasFacInMarketingDirector` |
| `ReasFacInMarketingDirector` | `DIREKTUR TEKNIK` | `DIREKTURTEKNIK` → `ReasFacInTechnicalDirector` |
| lainnya / tidak ada baris | — | `""` → Else = selesai |

Kolom limit per daftar: Bond (`IsLimitSBondKBG`, `IsLimitCustomBond`) → `LIMITBOND_BOTTOM`; Kredit CL dan Trade
Credit (`IsLimitCreditCL`, `IsLimitTradeCredit`) → `LIMITCREDITCL_BOTTOM`; Kredit NCL → `LIMITCREDITNCL_BOTTOM`.
`[terverifikasi]` kode `BusinessOldId` kelima predikat saling lepas (sensus atas `registry_gen.go`: irisan 0).

## Keputusan agent — menunggu konfirmasi (1 Oktober 2026, NB-12)

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A26 | 46 | Syarat langkah 22 `Local.TotalTSI >= @toDecimal(Local.Limit)` dinilai dengan DUA tafsir; berbeda hasil → `ErrTafsirLimit` | `[terverifikasi]` 22.2.1 `Local.Limit = .CARID2`, tetapi SQL bentuk B memetakan limit ke `CARI2`, bukan `CARID2` → `[dugaan]` Limit kosong = 0. Hanya berpengaruh bila `Local.TotalTSI` ≠ `CARID2` (endorsement dengan selisih top risk > 0, langkah 5) |
| ~~A27~~ | 46 | ~~Dua baris cocok dengan tujuan berbeda → `ErrUrutanTakPasti`~~ — **dicabut** (tinjauan kode): tiap antrean membuka paling banyak satu gerbang 22.2, jadi semua tujuan selalu sama; pemeriksaannya kode mati dan dihapus | — |
| A28 | 46 | Kolom limit kosong (bentuk A dan B) → `ErrLimitKosong` | Apakah kolom nullable `[pertanyaan terbuka]` (DDL tidak menyebut NOT NULL); NULL di SQL menyaring baris diam-diam — tidak ditebak |

## Keputusan work owner — 1 Oktober 2026, konfirmasi keputusan agent NB-11/NB-12

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 48 | Keputusan agent A20, A21, A23, A24, A25 (NB-11), A26, A28 (NB-12); teks `NBStatus` | **"Sesuai rekomendasi."** Ketujuh pilihan agent **dikonfirmasi**. `NBStatus` ("… IS IN <nama>'S INBOX", CARI5 = kolom `NAMA`): bila masih dibutuhkan, **nama diambil dari model peran**, bukan dari tabel limit — belum diport |

## Keputusan agent — menunggu konfirmasi (1 Oktober 2026, NB-14 dan NB-16)

| # | Tiket | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A29 | NB-16 | Pembaca fixture kasus mengisi `pyWorkPage.Quotation.*` dari `QuotationData.*` (akar fixture = halaman OfferFacIn) | Rekomendasi `PERTANYAAN-LANJUTAN.md` butir 5 ("pemuat mengisi kedua jalur … tetap mencatat bila berbeda"); pertanyaan jalur ganda masih terbuka. `[terverifikasi]` gerbang lini hanya membaca empat jalur (direkam 01-10-2026). Tanpa ini ketiga kasus FIRE tidak dikenali lininya |
| A30 | NB-14 | Tiket 14 diusulkan ditutup (`wontfix`) | `[terverifikasi]` semua pemanggil 1SA / 2SA / `GetLimitAkseptasi_Act2` di langkah berlabel `//` (butir 43); tabel bukti di tiket 14 |

## Keputusan work owner — 1 Oktober 2026, A29 dan A30

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 49 | `pyWorkPage.Quotation` vs `pyWorkPage.OfferFacIn.QuotationData` (`PERTANYAAN-LANJUTAN.md` butir 5; A29) | **"Iya, ganti dari QuotationData, itu isinya sama aja."** Kedua jalur berisi sama; `pyWorkPage.Quotation.*` diisi dari `QuotationData.*`. A29 dikonfirmasi; pertanyaan jalur ganda **terjawab** |
| 50 | Tiket 14 ditutup (A30) | **Dikonfirmasi** — `wontfix`, kedua cabang di-remark (butir 43) |

## Keputusan agent — menunggu konfirmasi (1 Oktober 2026, NB-07)

| # | Tiket | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A31 | NB-07 | Agregat `CurrencyList.Premium` dibandingkan dengan **kesamaan nilai desimal eksak** (`Cmp == 0`), bukan kesamaan teks | Nilai tersimpan tanpa nol di belakang koma (`22536540`) sedangkan premi coverage berskala 4 (`6158880.0000`). Bukan toleransi: satu digit pun tetap selisih |

**Pertanyaan terbuka baru (NB-07)** — untuk DBA / pemilik aplikasi Pega:

- ~~`[pertanyaan terbuka]` **Tipe properti `CurrencyList.Premium`**~~ → ✅ **terjawab butir 52: Decimal.**
- `[pertanyaan terbuka]` **Kasus MBU NB #5** (fixture `mbu_mata_uang.json`): keenam premi coverage cocok L1144, jumlahnya
  201.301.628, tetapi `CurrencyList.Premium` tersimpan 309.703.644,8. ⚠️ Ralat agent: catatan awal "= jumlah ÷ 0,65" keliru —
  309.703.644,8 × 0,65 = 201.307.369,12, rasio sebenarnya 0,64998… Work owner: "ikuti sistem lama" (butir 54).
  `[terverifikasi]` hipotesis yang diuji dan gugur (01-10-2026): gross-up pangsa `SetTSIPremiCedant_Act` (PercentShare kasus
  ini 5, bukan 65; kasus #1–#3 sejalur cocok tanpa gross-up); jumlah `PremiumGrossDiscountFleet` (204.676.628), `MinPremium`,
  premi ± `Discount` (diskon kosong), premi + `AdditionalCoverage` (204.676.628), TSI × (Rate + Loading) / 100 — semuanya
  201.301.628 atau 204.676.628. Nilai 309.703.644,8 hanya ada di `CurrencyList.Premium`; `OldData` tanpa coverage lain;
  `StatusBusiness = 1`. **Langkah 2.6 sistem lama atas data tersimpan pun menghasilkan 201.301.628** — nilai tersimpan ditulis
  jalur lain (mis. unggah) atau sisa sebelum coverage diubah: hanya pemilik aplikasi Pega yang dapat memastikan.

## Jawaban work owner — 1 Oktober 2026, tipe CurrencyList.Premium

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 52 | Tipe properti `CurrencyList.Premium` | **Decimal.** Tidak ada kehilangan presisi Double; teks tersimpan tanpa nol di belakang koma hanyalah bentuk tampil nilai Decimal yang sama. Jawaban ini menguatkan dasar A31 (perbandingan nilai eksak); A31 sendiri belum dikonfirmasi |
| 53 | Tambahkan agregat MBU ke kerangka rekonsiliasi tiket 16 | **"Tambahkan sesuai rekomendasi terbaik."** |

## Jawaban work owner — 1 Oktober 2026, kasus MBU NB #5

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 54 | Selisih agregat kasus MBU NB #5 | **"Untuk itu ikutin sistem lama."** Port langkah 2.6 tidak diubah — ia sudah mengikuti rule sistem lama, dan atas data tersimpan memberi 201.301.628. Nilai tersimpan 309.703.644,8 tidak terturunkan dari berkas kasus (lihat pertanyaan terbuka NB-07); tetap dilaporkan sebagai selisih sampai pemilik Pega menunjuk jalur penulisnya |

## Keputusan work owner — 1 Oktober 2026, mesin NB bersama dan R01

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 55 | Jalur mesin NB untuk modul lain; lingkup R01 | **Kontrak di inti** (`inti/backend/kontrak/facin.go`; penyedia `backend/services/kontrakfacin`). R01 **sebatas layanan** (tanpa migrasi DB, tanpa layar). Perintah lanjutan `PROMPT-LANJUT-IMPLEMENT-TANPA-TANYA-NB-RNW.md` (root repo): kontrak mencakup `rules.Eval`, `premium.Calculate`, `acceptance.Next`, resolver lini |

**Kontrak** — `TanggaAkseptasiFacIn.Langkah` (bentuk A, beralih ke bentuk B bila `ErrBentukB`), `PenilaiPredikatFacIn.Eval`,
`MesinPremiFacIn.Hitung / AsalRumus / LiniDariPredikat`; tipe `KasusFacIn`, `JabatanFacIn` ≠ `AntreanFacIn`, `PenggunaFacIn`,
`TransisiFacIn`, `MasukanPremiFacIn`, `HasilLiniFacIn`. Pemakai pertama: `rnwfacin` R01; bukti integrasi
`uji/lintasmodul/renewal_tangga_test.go` (lewat RNW = NB langsung, termasuk satu kasus yang naik ke KADIVTEKNIK).

## Keputusan agent — menunggu konfirmasi (1 Oktober 2026, kontrak dan R01)

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A32 | 55 | **Penyambungan lewat `Pendaftaran()` + perakit ditunda** sampai modul mendapat layar pertama; sampai itu pemakai merakit kontrak sendiri | `[terverifikasi]` `inti/backend/penjaga/menu_test.go` `TestMenuDimigrasiSamaDenganModulBackend`: setiap modul ber-`backend/modul.go` wajib `DIMIGRASI='1'` di menunya, yang menurut `MODUL.md` dinyalakan "saat modul mendapat layar pertamanya" (migrasi slot menu). `modul.go` sekarang = grup sidebar tanpa halaman | *(02-10-2026: tertutup untuk `PenilaiPredikatFacIn` dan `MesinPremiFacIn` — disediakan lewat `Pendaftaran()` nbfacin, tiket 20; `TanggaAkseptasiFacIn` tetap dirakit pemakai: butuh tabel limit per permintaan)*
| A33 | 55 | Gerbang masuk renewal ditetapkan: kasus baru = salinan data polis lama, `StatusBusiness = 2`, `OldPolicyNo` = nomor polis | Tiket R01: perilaku masuk tidak terekam (gerbang di konfigurasi work type/portal tak terekspor; `GetData_ACT` hanya `Obj-Open-By-Handle`). `[terverifikasi]` jalur diikat `RNW Fac In\Section\PeriodeRenewal.xml` (`.QuotationData.OldPolicyNo`, `.QuotationData.StatusBusiness`) |
| A34 | 55 | `RNWDate` tidak ditulis ke halaman kerja di R01; "Note" tanpa properti | Format clipboard `yyyyMMddTHHmmss.SSS GMT` (fixture) menuntut konversi zona waktu belum terverifikasi → R02. `[pertanyaan terbuka]` properti "Note" (tidak terikat di `PeriodeRenewal`) |
| A35 | 55 | Sumber polis lama = antarmuka `repository.PembacaPolisLama`, implementasi menunggu | `[pertanyaan terbuka]` tabel sumber (DBA) |

## Keputusan agent — menunggu konfirmasi (1 Oktober 2026, tiket 18 rumus premi lini lain)

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A36 | 18 | `PctAdjustment` FIRE **kosong = 0** (tidak adjustable) | `[terverifikasi]` gerbang langkah 21/29/37/45 `CountPremi_ACT` = `.PctAdjustment!=0`; properti kosong tidak lolos `!=0` di Pega `[dugaan]`. `[terverifikasi]` 54/54 coverage FIRE nyata cocok dengan tafsir ini |
| A37 | 14 | `@Math.divide` = **setengah-ke-atas** (premi `bagiBulat` dan `bulatNol` tangga). `ErrPembulatanSeri` dihapus. **Mengganti butir 14** | `[terverifikasi]` dua coverage FIRE kasus EDM nyata yang hasil baginya **tepat seri** di desimal ke-21, digit sebelumnya **genap**: tersimpan `…660\|5 → …661` dan `…214\|5 → …215` — setengah-genap akan memberi `…660` dan `…214`. Dihitung dua cara: Go (`TestBerkasKasusP5`) dan Python `decimal` independen (54 coverage FIRE, 2 seri, tersimpan = setengah-ke-atas 2, = setengah-genap 0). Hanya nilai positif yang terbukti |
| A38 | 18 | `.Premium*Local.Losslimit/100` (ANEKA/GOLF) dihitung eksak, **skala hasil = skala pembilang − skala pembagi**, diperpanjang hanya bila perlu | `[terverifikasi]` kasus Kredit nyata tersimpan `….0000` (skala pembilang 4), bukan 38 digit. `[dugaan]` aturan skala pilihan pembagian eksak Java `BigDecimal`; hasil yang tidak eksak ditolak, bukan dibulatkan |
| A39 | 16, 18 | Kasus EDM (`StatusBusiness = 3`) direkonsiliasi dengan **rumus coverage NB**, bertanda catatan | `CountGrossPremi_Act` langkah 5 bergerbang `StatusBusiness != 3`; jalur EDM `CountGrossPremiEDM_Act` milik `endorsmentfacin`, tidak diport di sini. `[terverifikasi]` 45/45 coverage FIRE kasus EDM tetap cocok — bukti untuk kasus ini saja, bukan untuk jalur EDM |
| A42 | 18 | `.Loading` ANEKA/GOLF **kosong = 0**, meniru A12 MBU | Belum teruji data: `[terverifikasi]` satu-satunya coverage ANEKA nyata ber-`Loading` `"0"` (urai dan grep `nb-kredit-1.json`) |
| A43 | 18 | Nol berbentuk teks lain dari `"0"` (mis. `"0.0"`) pada **pembanding teks** sistem lama → galat `ErrTeksNolAmbigu`: FIRE `.NetRate!="0"` (`CountPremi_ACT.xml` L4681) dan `Local.LossLimit=="0"` (L3572). ANEKA/GOLF `Local.Losslimit==0` (L2480 / `FillPremiGolf.xml` L1580) pembanding **angka** → `"0.0"` = bawaan 100 | `[terverifikasi]` Pega menulis desimal nol sebagai `"0.0"` (`PctAdjustment` 49/54 coverage FIRE nyata); `NetRate` nyata selalu `"0"` (54/54). `[pertanyaan terbuka]` apakah Pega membandingkan properti desimal dengan `"0"` sebagai teks: bila ya, `"0.0"` memakai NetRate 0 (premi 0) |

⚠️ Batas bukti A37: `bulatNol` tangga (`@Math.divide(x,1,0)`) memakai mode yang sama **berdasarkan analogi** fungsi
`@Math.divide` yang sama — `[dugaan]`; seri pada nilai dasar akseptasi belum pernah teramati di data.

**Penjaga tambahan tiket 18 (bukan keputusan, perilaku korpus):** FIRE `CoverageBasis` kosong/di luar 1–5 →
`ErrCoverageBasisTakDikenal` (tidak ada langkah premi yang terbuka, premi lama dibiarkan — tidak ditiru); ANEKA metode 1
pro-rata 0/kosong → `ErrBentukBelumDiport` (langkah 14.1 menghitung ulang dari tanggal, belum diport; GOLF tidak punya
langkah itu).

## Keputusan agent — menunggu konfirmasi (1 Oktober 2026, RNW R05–R06)

| # | Tiket | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A40 | R05 | `RenewalList_RD` diport sebagai **daftar kasus renewal milik pembuatnya yang belum selesai** *(ralat 01-10 malam: "di team group yang sama, kecuali Resolved-Completed dan Resolved-Rejected" — lihat bab Ralat)*, bukan "daftar kandidat jatuh tempo berkunci `OldPolicyNo` + tanggal berakhir" seperti premis tiket | `[terverifikasi]` `RNW Fac In\ReportDefinition\RenewalList_RD.xml` `pyContent/pyFilters`: `B AND A AND C AND D` = `.pxCreateOperator = Param.UserNameID`, `.pyStatusWork != "Resolved-Completed"`, `!= "Resolved-Rejected"`, `.Quotation.TeamGroup = Param.TeamGroup`; tanpa filter tanggal. Wadah `pyUI` memuat empat filter yang sama (urai: 2 wadah × 4; grep: 8 `pyFilterName`). Urut `pxCreateDateTime` DESC lalu `pyID` DESC; `pyMaxRecords` 500 (satu-satunya). Parameter kosong → galat (`pyUseNullIfEmpty` tidak diisi; perilaku Pega belum terverifikasi) |
| A41 | R06 | **R06 tidak diport**; diusulkan `wontfix` *(ralat 01-10 malam: dua pemanggil NB `PeriodeEndorsement` L2530/L2996 ditambahkan — lihat bab Ralat)* | `[terverifikasi]` keluaran `GetBusinessGroup_Act` = `ParamBis.CARI10` (langkah 4 hanya menghapus `OutBis`). Pembacanya, dua cara (grep teks dan urai elemen, jendela seluruh `D:\migrasi\RNM\**\*.xml`): di RNW **hanya** `ReportDefinition\BrowseAccountInsuredEDM` (filter `.GrupBusiness`); kedua pemanggil (`Section\PeriodeRenewal.xml` L1695, L2112) adalah tombol `click` pembuka harness `ChooseInsured` "Change Insured Name". Ketiganya fitur pilih-tertanggung yang **dibuang K-038** — memport R06 menghasilkan nilai tanpa pemakai |


## Ralat dan temuan lanjutan — 1 Oktober 2026 malam (verifikasi independen sesi `nusantarare-0f`, dicek ulang agent)

⚠️ **Status keputusan.** Sesi `nusantarare-0f` meneruskan pesan bahwa work owner **mengonfirmasi A37 (untuk premi `bagiBulat`
saja), A40, dan A41**, memilih **opsi (b)** untuk registry predikat EDM, serta **menahan** A36, A38, A39, A42, A43. Pesan itu
datang dari sesi agent lain, **bukan langsung dari work owner**, jadi di sini dicatat sebagai **dilaporkan, menunggu
konfirmasi langsung work owner** — tidak ditandai dikonfirmasi. `bulatNol` tangga tetap `[dugaan]` (analogi).

**Ralat A37 — jumlah seri.** Tertulis "dua coverage FIRE … tepat seri". Yang benar: **3 seri tepat** di desimal ke-21,
**2 di antaranya membedakan mode** (digit ke-20 genap: `edm-fire-1` Loc0.Item6.Cov4 digit 0, Loc0.Item8.Cov3 digit 4); yang
ketiga `nb-fire-1` Loc0.Item0.Cov0 (eksak `2040.50958904109589044755500`, digit ke-20 **5**, ganjil) tidak membedakan mode.
**Cara yang keliru:** skrip Python pertama menyaring `ROUND_HALF_UP != ROUND_HALF_EVEN`, sehingga hanya menghitung seri
yang membedakan mode, lalu dilaporkan sebagai jumlah seri. Dihitung ulang dua cara: Python `decimal` (sisa = ½ ulp) dan
`fractions` eksak (penyebut ×10²⁰ = 2) — keduanya 3. Kesimpulan A37 tidak berubah: tersimpan = setengah-ke-atas 2/2.

**Ralat A40 — ringkasan filter.** Ringkasan wajib memuat keempat filter: `.pxCreateOperator = Param.UserNameID`,
`.Quotation.TeamGroup = Param.TeamGroup`, `.pyStatusWork != "Resolved-Completed"`, `.pyStatusWork != "Resolved-Rejected"`.
Frasa "**belum selesai**" di tiket R05 dan `services/daftar.go` diganti "**kecuali Resolved-Completed dan Resolved-Rejected**"
— status `Resolved-*` lain tetap lolos.

**Ralat A41 — pemanggil NB terlewat.** Selain `RNW Fac In\Section\PeriodeRenewal.xml` L1695 dan L2112, `GetBusinessGroup_Act`
juga dipanggil `NB FacIn\Section\PeriodeEndorsement.xml` **L2530 dan L2996** — keduanya tombol `click` pembuka harness
`ChooseInsured` "Change Insured Name" (dicek: `pyEvent`, `pyWindowName`, `pyHarnessName`). Kesimpulan tetap.

**Ralat laporan — "fillPremiAneka dkk hanya dipanggil dari layar".** Benar untuk `fillPremiAneka` (pemanggil:
`Section\ViewCoverageAneka`), `fillPremiBond` (DataTransform) dan `CountNetPremiFOFire` (layar FacOut). **Salah** untuk
`FillPremiGolf` — dipanggil `CountGrossPremi_Act` L7498 (langkah 5.3, yang diport tiket 18) dan `CountGrossPremiEDM_Act`
L9924 — dan untuk `PremiPaymentMarine` (jalur unggah CSV, tiket 19). **Cara yang keliru:** pencarian pemanggil memakai pola
`>NamaRule<`, padahal pemanggilan activity tertulis `Call NamaRule` di `pyStepsActivityName`.

**A39 — dasar baru `[terverifikasi]`.** `CountGrossPremi_Act` langkah 6 (L9506) `call CountGrossPremiEDM_Act` memang **tanpa
prakondisi**, tetapi `CountGrossPremiEDM_Act` (`NB FacIn\Activity\CountGrossPremiEDM_Act.xml`, ASM-FW-GISFW-WORK!COUNTGROSSPREMIEDM_ACT)
hanya punya **satu langkah akar** yang membawa kelima sub-langkah, berprakondisi
`pyWorkPage.OfferFacIn.QuotationData.StatusBusiness==3`, salah → kode 3 Skip Step (L15916). Dua cara: pembongkar activity
dan urai `pySteps` akar langsung. **Kasus NB tidak diubah olehnya.** Untuk EDM, cabang 1.1 FIRE memanggil **`CountPremi_ACT`
yang sama** (1.1.2.1.4.2), dengan dua beda: `.TSI` ditimpa `Local.TSIObjectItem` (1.1.2.1.4) dan premi dinegasikan bila
`.FlagDelete=="1" && Local.deleteobjitem!="1"` (1.1.2.1.4.3). Pada `edm-fire-1` keduanya tidak terpicu (45/45 tanpa
`FlagDelete`, TSI coverage = `TSIObjectItem`) — itulah sebab 45/45 cocok. Varian NB, RNW, dan EDM rule ini **identik isinya**
(hash ternormalisasi 18 tag berbeda NB↔EDM hanya karena awalan `pyStepPageReference` `RH_6`/`RH_1` — jejak handle editor;
diusulkan sebagai tag volatil ke-19, `[pertanyaan terbuka]` pemilik register OQ-011).

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A44 | 16, A39 | Cabang EDM `CountGrossPremiEDM_Act` **tidak diport di nbfacin**; rekonsiliasi kasus EDM: beda 1.1.2.1.4 / 1.1.2.1.4.3 terpicu → **galat** *(ralat kedua: hanya 1.1.2.1.4.3; 1.1.2.1.4 `[dugaan]` → catatan)*, lini EDM selain FIRE → **belum tercakup** | Satu-satunya efeknya berlaku saat `StatusBusiness == 3` (siklus endorsement, jatah `endorsmentfacin`); varian di `Endorsment Fac In\Activity\` identik isinya, sehingga sisi EDM memportnya dari foldernya sendiri (CLAUDE.md §4.6). Bagi NB = no-op, terverifikasi di atas |
| A45 | 19 | `@sum` premi/diskon kargo **lintas mata uang ditolak** (`uang.ErrMataUangBerbeda`); nilai kosong dan daftar kosong ditolak | Pega menjumlahkan angka tanpa melihat mata uang; `[terverifikasi]` fixture MARINE nyata 104/104 coverage IDR, `Discount` `"0.0000"` 104/104 | *(dikonfirmasi work owner 02-10-2026, butir 64)*

## Keputusan work owner — 1 Oktober 2026 malam, konfirmasi langsung (AskUserQuestion)

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 56 | Keputusan yang diteruskan sesi `nusantarare-0f` | **"Ya, benar keputusan saya."** → **A37 DIKONFIRMASI untuk premi (`bagiBulat`)**; `bulatNol` tangga tetap `[dugaan]` (analogi). **A40 DIKONFIRMASI. A41 DIKONFIRMASI** (R06 → `wontfix`). Registry predikat EDM = **opsi (b)**. **DITAHAN:** A36, A38, A39, A42, A43 |
| 57 | Mesin `rules` pindah ke `inti` agar dipakai EDM tanpa impor `modul/nbfacin`? | **"EDM DAN NB MENU YANG TERPISAH."** Tafsiran agent: NB dan EDM modul terpisah — mesin **tidak** dipindah ke `inti` untuk dipakai bersama; generator `bangkit` (`-folder`, `-paket`) dijalankan sisi EDM ke paketnya sendiri, dengan mesin evaluator miliknya sendiri. ⚠️ Bila tafsiran ini keliru, butir ini dibuka lagi |

Bab "Ralat dan temuan lanjutan" di atas: status keputusan di paragraf pertamanya kini **diganti butir 56**.

## Ralat kedua — 1 Oktober 2026 malam (verifikasi independen sesi `nusantarare-0f`, dicek ulang agent)

Perintah: `PROMPT-LANJUT-IMPLEMENT-NB-SAJA.md` langkah 1. Kelima koreksi dicek ulang ke berkas sebelum dicatat.

1. **Uji mutasi tiket 18/19.** Tertulis "45/48, 3 lolos semuanya ekuivalen" — **tidak tepat**. Yang ekuivalen sungguhan
   hanya `bayar: kosong diterima` (`utils.ParseDecimal("")` juga gagal). Dua lainnya (`rekon: IsMBD selalu false`,
   `rekon: GOLF memakai jalur ANEKA`) adalah **celah tes**: data nyata tidak menjangkau cabangnya. **Cara yang keliru:**
   label "ekuivalen terhadap data" di skrip menyamakan "tidak terbedakan oleh fixture" dengan "tidak terbedakan oleh
   perilaku". Yang benar: **45/48 = 1 ekuivalen + 2 celah tes**. Ditutup `TestBerkasKasusCabangBuatan` (masukan buatan:
   GOLF di `RiskLocation.AnekaList`, `BusinessType` MBD + `IndemnityPercentage`). Angka kini **49/50**, 1 ekuivalen
   (skrip + log: `alat/mutasi_tiket18_19.*`; ditambah mutan A46 dan dua mutan tanda `[dugaan]` TSI).
2. **CountGrossPremiEDM_Act — "TSI diambil dari item" diturunkan ke `[dugaan]`.** Pasangan `.TSI = Local.TSIObjectItem`
   (L1273) tersimpan di `pyParamArray` langkah 1.1.2.1.4 yang `pyStepsActivityName`-nya **kosong** (loop EMBEDDED atas
   `.CoverageList`) — kemungkinan parameter sisa yang tidak dieksekusi. **Cara yang keliru:** pembongkar activity agent
   mencetak pasangan `pyParamArray` sebagai `set` tanpa memeriksa bahwa metode langkahnya `Property-Set`. Alat sudah
   diperbaiki (pasangan seperti itu kini ditandai "PARAM SISA"). **Dampak:** dari sembilan activity dasar port
   (`CountPremi_ACT`, `CountPremiCoverageAneka`, `FillPremiGolf`, `CountGPWMarinePAMbu_Act`, `PremiPaymentMarine`,
   `FillPremiMBU_FacIn`, `CalculatePremiPA_FacIn`, `CountGrossPremi_Act`, `InputDtlPayment_PreAct`) **nol** memuat parameter
   sisa — rumus yang sudah diport tidak terkena; `CountGrossPremiEDM_Act` memuat 4. Yang terbukti:
   1.1.2.1.4.2 (L1468) memanggil `ASM-FW-GISFW-DATA-COVERAGE!COUNTPREMI_ACT` dengan `Param.CallAct="Upload"`; 1.1.2.1.4.3
   menegasikan `.TSI` (L1623), `.TSILiability` (L1677), `.Premium` (L1697) bila `.FlagDelete=="1" && Local.deleteobjitem!="1"`.
   **A44 disesuaikan:** TSI coverage ≠ `TSIObjectItem` kini **catatan `[dugaan]`, bukan galat** (`catatanTSIItem`); galat
   tinggal untuk `FlagDelete`. Bab A39 di atas yang menyebut "`.TSI` ditimpa" dibaca dengan ralat ini.
3. **"Renewal MARINE tidak melewati cabang mana pun" — sebagian.** Benar tidak ada cabang renewal (`PremiPaymentMarine`
   hanya 2 langkah akar, `IsNB` dan `IsEDM`). **Batas bukti:** adaptor fixture `rekonsiliasi.kasusJSON.Nilai` memetakan
   `pyWorkPage.Quotation.*` ke `QuotationData.*`, sehingga `IsNB` (`Quotation.StatusBusiness`) dan `IsEDM`
   (`OfferFacIn.QuotationData.StatusBusiness`) membaca **kolom yang sama** di fixture — kondisi "keduanya benar" hanya
   teruji dengan masukan buatan (`TestMarineGerbang`), tidak dari data nyata.
4. **A45 — tes mata uang campur** kini `errors.Is(err, uang.ErrMataUangBerbeda)` (sebelumnya hanya `err != nil`); seluruh
   kasus `TestMarineDitolak` memeriksa galat yang tepat.
5. **Tag volatil ke-19 — usulan diralat.** Membuang seluruh `pyStepPageReference` menghilangkan nomor langkah. Yang tepat:
   **samakan awalan `RH_n` saja**. Diverifikasi dua cara (baris terurut + pohon kanonik): NB = RNW = EDM identik. **Cara
   yang keliru (agent, percobaan pertama):** pola `RH_[0-9]+\.` mensyaratkan titik, sehingga rujukan langkah akar
   `RH_6` (tanpa titik) lolos dan EDM tampak berbeda; pola yang benar `RH_[0-9]+`.

**Dua temuan untuk dicatat (bukan diport):**
- a) `[dugaan]` `CurrencyList[0].Policy.Payment.Premium` `1319760.1` di `nb-marinecargo-1` = 26.395.202 × `PercentShare` 5 %
  (akar `PercentShare` = "5"). Wadah lain dari yang dicari tiket 19 — petunjuk, **bukan** pembanding.
- b) **Bahan OQ RBAC:** gerbang `OperatorID.pyUserIdentifier == <literal ID operator>` (nilai tidak disalin) di
  `pyStepsPreCondParamsWhen`: `CountGrossPremiEDM_Act` **8** (antara lain L2155, L2379 — langkah 1.1.2.1.4.5.1–2 spreading)
  dan `CountGrossPremi_Act` **6** — dua cara (grep baris dan urai elemen per wadah), seluruhnya di wadah prakondisi.
  `[terverifikasi]` kode NB **tidak** memuat gerbang ID operator literal: `registry_gen.go` 0 baris merujuk medan login
  (`medanLogin` generator → panic tanpa nilai, `TestLoginTidakDisalin`); di luar itu `pyUserIdentifier` hanya muncul di
  komentar `acceptance/tangga.go`.

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A46 | 19 | Coverage kargo **tanpa mata uang ditolak** (`ErrMasukan`) | Dua coverage bermata uang kosong lolos `Money.Add` (mata uangnya "sama"); `[terverifikasi]` fixture MARINE nyata 104/104 ber-`Currency.Name` "IDR" | *(dikonfirmasi work owner 02-10-2026, butir 64)*

## Temuan data identitas di registry predikat — 1 Oktober 2026 malam

Dilaporkan sesi `nusantarare-55` (dua nomor polis), diperluas agent dengan sensus medan pembanding × literal (nama medan
saja, nilai tidak dicetak). `rules/registry_gen.go` sempat memuat **literal tidak-kosong untuk medan identitas**:
`Quotation.OldPolicyNo` 2 (`ISERRORSPREADING`), `MarketingName` 3 (`ISTBONDING`, bentuk huruf berspasi — kemungkinan nama
orang), `OperatorID.pyTelephone` 2 (`ISTREATY1`, `ISSPVTREATY1` — ternyata **kode peran**, dibahas luas di discovery, bukan
telepon pribadi). **Diperbaiki:** generator `bangkit` `medanIdentitas` → panic tanpa nilai (literal kosong tetap diport);
registry dibangkitkan ulang (penjaga korpus hijau); penjaga tanpa korpus `TestRegistryTanpaLiteralIdentitas`.
**Sebaran lain** (dicari dengan skrip, nilai tidak dicetak): dua nomor polis juga di
`docs/_ARSIP-lintas-siklus/04-aturan/01-katalog-when.md` → **disamarkan**; tiga nama di registry `endorsmentfacin` →
**dikabarkan ke sesi `nusantarare-55`, tidak disentuh**. `[terverifikasi]` kelima nilai sensitif (2 nomor polis, 3 nama)
**tidak pernah masuk riwayat git** (`git log --all -S`, dijalankan dari skrip).

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A47 | 07 | `TotalTSIPremiGrossList.TSI` dibandingkan sebagai **nilai desimal eksak**; `Name`, `Premium`, `Rate` sebagai teks persis | Tersimpan `8400000.0` untuk `TSIObjectItem` `8400000` (bentuk `Double.toString` Java); `CurrencyList` akar sekelas juga `8400000.0`. `[dugaan]` properti TSI kelas `OfferFacIn-Currency` bertipe Double — `[pertanyaan terbuka]` tipe properti (pemilik Pega). Preseden A31 | *(dikonfirmasi work owner 02-10-2026, butir 64)*

## Ralat ketiga — 2 Oktober 2026 (code review dua sumbu kelompok koreksi + tiket 07)

1. **Total lokasi FIRE (tiket 07) melewati `double` di sistem lama.** `[terverifikasi]` `SumTotalTSIPremiGross_Act`
   mendeklarasikan `Local.TSI` dan `Local.Premi` **`double`** (L293–295, L299–301; `Local.PremiCov` `Decimal`, L305–307),
   dan cabang 1.3.4 (`.Premium = Local.Premi + .Premium`) memakai keduanya. `[terverifikasi]` model "item kedua dst.
   melewati double" mengulang nilai tersimpan `edm-fire-1` (9 item USD) **sampai digit terakhir**; penjumlahan eksak
   tidak — diperiksa dua pihak (review spec; simulasi Python agent `repr(float(x))`). **Cara yang keliru:** port pertama
   menganggap seluruh penjumlahan eksak karena `Local.PremiCov` Decimal, tanpa membaca tipe `Local.TSI`/`Local.Premi`.
   **Atribusi yang keliru:** komentar tiket 07 menyebut total lokasi EDM "dari CountGrossPremiEDM_Act 1.1.4–1.1.5" — total
   itu ternyata dihasilkan **rule yang sama** (model di atas mengulangnya). Kode kini: entri cabang 1.3.4 bertanda
   `LewatDouble` dan dilaporkan **belum tercakup**; kasus EDM tidak lagi diistimewakan (9 item cocok, 1 lokasi belum
   tercakup). ⚠️ Data nyata NB/RNW hanya punya 1 item per lokasi — kriteria "daftar cukup panjang" belum terpenuhi untuk
   FIRE; cabang 1.3.4 hanya teruji masukan buatan.
2. **A47 — alasan diralat.** "Bentuk `Double.toString`, properti TSI `[dugaan]` Double" dicabut: TSI tersimpan EDM
   `286367200` tidak berbentuk `Double.toString` (yang akan mencetak `2.863672E8`). Yang tersisa: perbedaan **format**
   angka (`8400000.0` lawan `8400000`), sebabnya `belum terverifikasi`. Perbandingan nilai tetap.
3. **Fallback `.Currency.Name` di `barisTotal` dibuang** — tanpa dasar XML (`Local.Currency = .Currency`, 1.3.3);
   `[terverifikasi]` 11/11 item FIRE fixture ber-`.Currency` teks.
4. **Tiket 19 — wadah `CargoList`.** `PremiPaymentMarine` berkelas `ASM-FW-GISFW-Work` (L66), jadi `.CargoList()` =
   `pyWorkPage.CargoList`; rekonsiliasi menjumlahkan `OfferFacIn.CargoList`. `[dugaan]` keduanya berisi sama — belum
   terverifikasi; baris tetap belum tercakup, bukan pembanding.
5. **`catatanTSIItem` membandingkan teks** → kini nilai (`8400000` lawan `8400000.0` bukan perbedaan; tes ditambah,
   mutan "dibanding teks" kini tertangkap).
6. **Generator — tiga celah ditutup:** ekspresi yang **menyebut** medan identitas/login dalam bentuk apa pun (mis.
   `@String.equals`) tidak disalin ke pesan panik; nama medan dicocokkan tanpa peduli huruf besar-kecil; `pilihProfil`
   menerima `NB FacIn/When`. Penjaga registry kini juga memindai teks panik. Registry dibangkitkan ulang: **0 baris
   berubah** (tidak ada kebocoran yang lolos sebelumnya).
7. **Alat mutasi — akhir baris.** Skrip membaca dengan normalisasi dan menulis LF, sehingga "pemulihan" mengubah CRLF
   menjadi LF tanpa terdeteksi (pemeriksaan "pulih" membaca lewat normalisasi yang sama). Kini bita asli disimpan dan
   dipulihkan persis, "pulih" dibandingkan per bita. Angka sekarang **59/60** tertangkap, 1 ekuivalen sungguhan.
8. **Standar:** label `[keputusan agent]` (bukan label §4.9) di `pembayaran/marine.go` diganti "Keputusan agent Axx";
   nilai uji berawalan `UJI-` (README-BACA-DULU L111); komentar `kasusJSON` yang basi diperbarui. ⏸ Satu label
   `[keputusan agent]` tersisa di `rnwfacin/backend/services/daftar.go` — modul itu dibekukan, tidak disentuh.

| # | Butir | Pilihan agent | Dasar |
| --- | --- | --- | --- |
| A48 | 07 | Penjumlahan cabang 1.3.4 `SumTotalTSIPremiGross_Act` **tetap eksak**; entri `LewatDouble` dilaporkan belum tercakup | Meniru aritmetika `double` untuk uang dilarang tanpa keputusan ADR (CLAUDE.md §7). `[pertanyaan terbuka]` ikuti hasil `double` sistem lama (paritas) atau eksak (perubahan perilaku sadar)? | *(diputuskan work owner 02-10-2026, butir 63: tetap eksak)*

## Keputusan work owner — 2 Oktober 2026, NB sebagai modul berjalan (AskUserQuestion)

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 58 | Endpoint pertama vs model peran (OQ RBAC belum diputuskan) | **"Akseptasi, jabatan dari body"** — jabatan pengguna dikirim dalam permintaan. ⚠️ Agent mencatat: **tidak aman untuk produksi** (klien dapat mengaku jabatan mana pun); ditandai di kode dan respons, diganti saat model peran diputuskan |
| 59 | Berkas slot menu `962` (di luar rentang 180–219 berkas perintah) | **"Ya, tulis 962"** — slot 962–963 milik nbfacin di `MODUL.md`; nilai slot tidak diubah; berkas tidak dijalankan |
| 60 | Layar pertama | **"Tiru section Pega"** — port section NB apa adanya (tiket 21), bukan layar rekaan |
| 61 | Layar pertama, sesudah temuan: data kasus NB belum bersumber (tiket 17) dan blok coverage Pega read-only | **"Section + isian uji premi"** — port setia blok pertama `InputCoverageCargo_FacIn`; Name, Rate (%), TSI dapat diisi untuk memanggil hitung premi, ditandai **alat periksa** (penyimpangan sadar dari read-only Pega) |
| 62 | Menyalakan NB memerahkan 3 uji perakit bersama `frontend/` (jumlah modul aktif dikunci 4; uji gigit memakai nbfacin sebagai contoh belum dimigrasi) — boleh disunting? | **"Ya, sunting minimal"** — 4 → 5; contoh "belum dimigrasi" diganti `nbtreatyin` (tanpa `modul.go`, `DIMIGRASI='0'`). Ditinjau tim inti (CODEOWNERS) |

## Ralat keempat — 2 Oktober 2026 (code review tiket 20–21)

- **Tipe kolom limit:** `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`, `LIMITBOND_BOTTOM`, `LIMITCREDITCL_BOTTOM`,
  `LIMITCREDITNCL_BOTTOM` bertipe **`NUMBER(*,0)`** (bilangan bulat) di keenam DDL, bukan "NUMBER tanpa presisi".
  **Cara yang keliru:** pola ekstraksi agent `\([0-9, ]*\)` tidak mengizinkan `*`, sehingga `(*,0)` terbuang; diulang
  dengan `\([^)]*\)`. Dampak: hanya teks dokumen (`STRUKTUR-TABEL-NB-FACIN.md`, tiket 20, komentar repository) —
  pembacaan tetap teks → desimal.
- **Tiket 21:** "semua `pyEditOptions = Read-only`" berlaku untuk **medan**; kedua tombol `Auto`.

## Keputusan work owner — 2 Oktober 2026, diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` dengan kutipan jawaban work owner **"oke sesuai rekomendasi km aja"** atas
rekomendasinya (pola yang sama dengan butir 56, yang work owner konfirmasi langsung).

| # | Pertanyaan | Keputusan |
| ---: | --- | --- |
| 63 | A48 — paritas `double` sistem lama atau eksak? | **Tetap eksak** — perubahan perilaku yang disengaja. Lokasi FIRE ≥2 item bermata uang sama tetap dilaporkan **belum tercakup** di rekonsiliasi, tidak disesuaikan. Meniru `double` untuk uang hanya lewat ADR bila kelak paritas diwajibkan (CLAUDE.md §7) |
| 64 | A45, A46, A47 | **Dikonfirmasi** |
| 65 | Slot 962 / `-migrate` | **Agent tidak menjalankan `-migrate`.** Hanya work owner sendiri, ke Oracle DEV/lokal miliknya; **tidak** ke basis data bersama sebelum tim inti meninjau `facin.go`, `962_menu_nbfacin.sql`, dan tiga uji menu bersama; **tidak** ke produksi sebelum RBAC diputuskan dan endpoint akseptasi berhenti mempercayai jabatan dari isian (butir 58). ⚠️ `[terverifikasi]` `-migrate` menjalankan migrasi tertunda **semua modul terdaftar**, bukan hanya 962 (`cmd/api/main.go` `jalankanMigrasi`, L188: "SEMUA modul terdaftar, tidak bergantung modul aktif") |

**Rumusan A48 yang tepat** (menggantikan "menjumlahkan lewat double" di laporan): item ke-2 dst. **dibulatkan ke
`double`** — `Local.Premi`/`Local.TSI` bertipe `double` (`NB FacIn\Activity\SumTotalTSIPremiGross_Act.xml`
`pyLocalParameters` L293–301), `Local.PremiCov` `Decimal` (L305–307) — lalu **diakumulasi `Decimal`**. **Batas bukti:**
model itu cocok **2/2 lokasi** (`edm-fire-1` akar + `OldData`; eksak 0/2, `double` murni 0/2 — verifikasi independen
sesi `nusantarare-0f`), dari **satu kasus** saja.

## Yang belum diputuskan


- Nama modul `facout` dan nomor usulan `760-799` / `990-991` — perlu persetujuan tim inti.
- ~~Mesin bersama NB … lewat `inti/backend/kontrak` (saran), atau ditaruh di `inti/`.~~ → ✅ **butir 55**: kontrak
  `inti/backend/kontrak/facin.go`. Tinggal: persetujuan tim inti atas berkas kontrak itu, dan penyambungan perakit (A32).
