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

## Ralat kelima — 2 Oktober 2026 (tiket 17 tidak menunggu tabel flat)

Laporan agent 02-10 menulis "tiket 17 menunggu keputusan tabel flat dari work owner". **Keliru.** Register NB
`00-KEPUTUSAN-WORK-OWNER.md` sudah menutupnya: **K-069** ("seluruh butir rekonsiliasi tabel flat TERTUTUP"; 78 tabel,
1.329 kolom, DDL draf lolos delapan pemeriksaan), **K-073** (Jalan B disetujui), **K-074** (seam `loader.Flatten`
disetujui; seam `repository` "kini aktif"); `_DAFTAR-ISSUE-TERBUKA.md` **F-3**: empat tiket yang dulu terhalang struktur
tabel flat "kini tidak terhalang lagi". **Cara yang keliru:** agent bersandar pada baris "Blocked by" tiket 17 tanpa
memeriksa register; dan spec pemuatan (`04-spec\11-spec-pemuatan-data-lama.md`) beserta bahannya (`08-flat\`) **tidak
ikut dipindah** ke `nbfacin/docs` saat struktur satu folder per modul, sehingga tidak terbaca dari folder modul.

**Sumber yang ditetapkan work owner (02-10-2026, diteruskan sesi `nusantarare-0f`), READ-ONLY:**
`D:\migrasi\RNM\OUTPUT\04-spec\11-spec-pemuatan-data-lama.md` (md5 `0d0e3c6a…`), `08-flat\BAHAN-SPEC-PEMUATAN.md`
(md5 `0905233f…`), `08-flat\DDL-tabel-flat-draf.sql`, `08-flat\Tabel-Flat-Lintas-Siklus.xlsx` — `[terverifikasi]` md5
kedua berkas pertama cocok dengan yang disebut; DDL 78 `CREATE TABLE`, 78 PK, 64 FK, 77 indeks (dua cara: grep baris dan
urai pernyataan; cocok K-069). ⚠️ Salinan `jefri/OUTPUT FIX/` sempat hilang dan disalin ulang dari `D:\migrasi\RNM\OUTPUT\`
(md5 cocok 19/19, menurut sesi `nusantarare-0f`); versi yang hilang mungkin sempat berubah sesudah 16–25 September.

**Dua sumber skema berselisih — dicatat, tidak dipilih:** DDL draf diturunkan dari `Tabel-Flat-Lintas-Siklus.xlsx`
(**78** tabel, **1.329** kolom); `D:\migrasi\RNM\Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx` (sumber menurut K-065)
memuat **75** tabel, **1.290** kolom (lembar `Daftar Tabel` dan `Kolom`, dibaca langsung dari XML xlsx). `[terverifikasi]`
lembar `BACA-INI` workbook pertama menyatakan dirinya "Turunan dari `Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx`",
ditambah dimensi siklus dan cabang retro (3 wadah, 24-09). `[pertanyaan terbuka]` mana yang mengikat bila keduanya
dipakai bersamaan.

> **Diukur 02-10-2026 (tiket 22), dua cara — pengurai xlsx vs pencacah baris DDL:** `Tabel-Flat-Lintas-Siklus.xlsx` =
> `DDL-tabel-flat-draf.sql` **persis** (78 / 1.329, selisih nol dua arah). `Claude outputs` (75 / 1.290) adalah
> **himpunan bagian murni**: 39 kolom selisih = 3 tabel wadah V-27 (20 kolom sistem) + 18 `CURRENCY_CODE`
> (K-063/K-069) + `OLD_POLIS_ID` (K-068); **nol** kolom ke arah sebaliknya. Perselisihan "75 lawan 78" karena itu
> **bukan dua rancangan berbeda**, melainkan rancangan lama dan rancangan sesudah keputusan 24–25 September. Pertanyaan
> "mana yang mengikat" tetap terbuka secara resmi; generator `loader` membaca Lintas-Siklus karena itulah yang
> ditetapkan bersama DDL draf (bab ini, paragraf sumber).

## Keputusan work owner — 2 Oktober 2026, presisi DDL tabel flat (AskUserQuestion)

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 66 | DDL draf memakai 415 kolom `NUMBER` polos + 1 `NUMBER(3)`; penjaga inti `TestNolNumberTanpaPresisi` hanya mengizinkan `NUMBER(38,8)`/`(5)`/`(10)`/`(19)`; uang `NUMBER(38,8)` memotong premi FIRE 20 desimal | **"Flatten dulu, migrasi ditahan."** `loader.Flatten` (murni) dikerjakan; berkas migrasi 180–219 ditulis setelah tim inti memutuskan presisi tabel flat. Penjaga inti tidak diubah |

## Keputusan agent A49–A65 — DIKONFIRMASI work owner 2 Oktober 2026 (butir 68.7) (tiket 22 `loader.Flatten`)

| # | Tiket | Keputusan agent | Dasar |
| --- | ---: | --- | --- |
| A49 | 22 | `ROW_UID` fase 1 **deterministik**: UUID v5 (SHA-1) atas `IDPEGA` + tabel + posisi mentah | K-073 hanya menyebut "dibangkitkan, sementara"; deterministik menjaga Flatten murni dan uji repository "dapat diulang" (spec 11). Larangan K-073 (tidak dijadikan sandaran di luar tabel flat) tetap |
| A50 | 22 | `SEQ_NO` = **posisi di larik sumber**; unsur yang dibatalkan V-27 meninggalkan celah | V-40 menjodohkan versi menurut posisi; V-43 baris tidak pernah dihapus. Penomoran ulang akan menggeser pasangan |
| A51 | 22 | V-16 "ada" = **larik dengan ≥ 1 unsur** | `[dugaan]` V-16 diukur di XML; di JSON larik kosong lazim (`nb-marinecargo-1`: `LocationList`, `VehicleList` kosong). Hasil 115 contoh cocok lembar Grup Bisnis |
| A52 | 22 | BusinessType yang mengikat = `QuotationData.BusinessType` **akar** | Kemunculan lain di fixture seluruhnya di `OldData` (V-19); 128 = 115 + 13 berkas EDM `[dugaan]` |
| A53 | 22 | Medan milik `CurrencyList/Policy` **tidak** dilipat; hanya `Policy/Payment` (awalan `Pay`) | V-22 melipat Policy karena "nol kolom terisi" — tidak memutuskan medannya. `[terverifikasi]` `FacOfferList/CurrencyList/Policy.TSI` 15 unsur menimpa `CurrencyList.TSI` bila dilipat |
| A54 | 22 | `T_FR_CURRENCYLIST` mendapat lipatan `Policy/Payment` yang sama | Rancangannya memuat sembilan kolom `PAY_` yang sama dengan `T_CURRENCYLIST` |
| A55 | 22 | V-33: `ASMCoverage` di bawah `PersonList` → `T_COVERAGELIST` | `[dugaan]` lokasi dari Jalur Sumber "PersonList/CoverageList \| Life PA"; nol fixture PA |
| A56 | 22 | V-30: `FACTOR_GROUP` / `FACTOR_NAME` = **nama halaman apa adanya** | Peta "23 faktor" (lembar ScoringRisk) tidak ada di kedua workbook; data memuat **24** objek faktor |
| A57 | 22 | Medan `py*` dibuang sebagai metadata | Lembar Kolom memuat **nol** FIELD ASLI ber-awalan `py`; terhitung `Dibuang` |
| A58 | 22 | Angka > 38 digit **tidak dibulatkan**, hanya dihitung | ADR-0005; bahan keputusan presisi (butir 66) |
| A59 | 22 | Bentuk masukan XML (spec 11 uji 11) **tidak dibangun** | K-066 masukan produksi JSON; contoh XML tidak boleh masuk repositori |
| A60 | 22 | `CurrencyList.ID` **tidak** dipetakan ke `CURRENCY_CODE`; `CURRENCY_CODE ← Name` | `[terverifikasi]` `ID` 5/116 terisi, seluruhnya angka lima digit; `Name` kode tiga huruf 116/116 — bunyi V-22 tidak cocok dengan data |
| A61 | 22 | Berhenti keras spec 11 "kode enumerasi yang artinya belum dijawab tercapai" (`BusinessCode`, `BusinessOldId`) dibaca: berlaku bila sebuah langkah **menafsirkan** kode itu. Flatten tidak menafsirkannya — disimpan apa adanya (V-6) — jadi tidak berhenti | `[dugaan]` atas maksud spec; tafsir harfiah (berhenti setiap kali kode muncul) menghentikan hampir setiap dokumen. Temuan code review sumbu spec |
| A62 | 22 | V-16 langkah 5: BusinessType **dikenal tetapi bukan Life/PA** → `ErrLiniBisnis` | V-16 "sisanya → Life atau PA" tidak memberi kelompok untuk 16 nilai lainnya; menebak kelompok dilarang |
| A63 | 22 | Metadata dibuang menurut **awalan** `px`/`pz` (+`py`, A57), bukan daftar "23 tag" | Rujukan BAHAN "kontrak 23 tag K-042/K-043" tidak dapat ditelusuri: K-042/K-043 di register NB soal lingkup EDM. Semua yang dibuang terhitung |
| A64 | 22 | Empat berhenti keras di luar tabel spec: `ErrIDPega`, `ErrDokumen`, `ErrKolomGanda`, `ErrBentuk` | Keadaan mesin yang hanya dapat "diselesaikan" dengan menebak (IDPEGA tak berbentuk, JSON rusak, satu kolom dua nilai, bentuk simpul salah) |
| A65 | 22 | Keluaran Flatten **generik** (`map` kolom → `Nilai{Teks, Angka}`); pasangan uang–mata uang lewat `CURRENCY_CODE` baris; bentuk baris per tabel di `models` diserahkan ke tiket 24 | Spec 11 *Modul yang dibangun* menyebut `models` diperluas; 78 tabel × 1.329 kolom sebagai struct tangan = data clump besar tanpa pemakai sebelum repository ada |

**Temuan tiket 22 yang menunggu work owner** (rinci di tiket 22 bab *Butir terbuka*): ⛔ tipe `NUMBER` delapan kolom yang
isinya teks — Flatten berhenti keras di **106 dari 115** contoh dan di kelima fixture; pembawa mata uang
`T_COVERAGELIST`/`T_ANEKALIST` (`Name` sendiri tidak ada); `T_PROPERTY.CURRENCY_CODE` selalu UNKNOWN; 29 FK V-47 tanpa
aturan isi; `TSI_TOP_RISK` tanpa rumus; medan JSON yang tidak ada di rancangan (`CoverageInitial`, `AdditionalShip`);
induk ganda 12 vs 13 (Daftar Relasi tanpa tiga tabel wadah); **K-069 (7b) `IsCedingConfirm` "kolom sendiri" tidak ada
di 78 tabel maupun DDL draf** — medannya di `ViewSuggest`, yang V-31 arahkan ke `HISTORYAKSEPTASIPRODUCTION`.

**Ralat agent (tiket 22, sebelum dilaporkan):** hitungan tangan "26 kolom FK V-47" keliru — **29** (tertangkap
`TestSetiapKolomBerasal`); pembaca xlsx pertama menggeser sel kosong (tertangkap saat analisis, diperbaiki, diuji
`TestLembarMenempatkanSelMenurutKolom`); cara kedua `EDM_CHARGE_FEE` memberi 0 karena ejaan `EDMChargeFee` (benar:
`EdmChargeFee`, jebakan sensus no. 4).

## Keputusan work owner — 2 Oktober 2026, butir terbuka tiket 22, diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` dengan kutipan jawaban work owner **"setuju"** atas rekomendasinya (pola butir
56 dan 63–65). Aturan `PROMPT-LANJUT-IMPLEMENT-NB-SAJA.md` tetap berlaku.

| # | Butir | Keputusan |
| ---: | --- | --- |
| 68.1 | Delapan kolom `NUMBER` berisi teks | **`VARCHAR2`, teks apa adanya** — `%` tidak ditafsirkan, koma desimal tidak dinormalkan, spasi di ujung `PctLimit` **tidak** dipangkas. Penyimpangan sadar dari DDL draf, baris DDL-nya di tiket 22 |
| 68.2 | Presisi tabel flat | **Tetap keputusan tim inti.** Agent menulis usulan (`docs/USULAN-PRESISI-TABEL-FLAT.md`): teks mentah untuk kolom uang/rate yang dapat > 38 digit atau > 8 desimal, `NUMBER` turunan bila perlu; preseden ADR-0034. Tiket 23 tetap ditahan |
| 68.3 | Kode mata uang `T_COVERAGELIST` / `T_ANEKALIST` (dan `T_PROPERTY` bila leluhurnya punya halaman Currency) | **dari `Currency/Name`** |
| 68.4 | 29 FK V-47, `TSI_TOP_RISK` | **dibiarkan kosong** sampai aturannya tertulis |
| 68.5 | `CoverageInitial`, `AdditionalShip` | **ke penampung medan tak dikenal (ADR-0023)**, tidak dibuang |
| 68.6 | Induk ganda / V-30 | **ikuti data (24 faktor)**; selisih 23 vs 24 dicatat sebagai ralat V-30 |
| 68.7 | A49–A65 | **Dikonfirmasi** |
| 68.8 | Tiket 24, penawaran gagal | **lewati, catat galat, lanjut; ringkasan di akhir muat.** Medan tak terpetakan ke penampung, bukan alasan menolak muat. J-16 tetap pertanyaan terbuka |

**Cara agent menerapkannya** (tiket 22): 68.3 **hanya** dua tabel yang disebut — kembaran `T_FR_*` tidak diperluas
(pelajaran A53: ekstrapolasi ke cabang FR pernah keliru). `T_PROPERTY`: syaratnya `[terverifikasi]` tidak terpenuhi —
nol halaman `Currency` di bawah `LocationList` (5 fixture; 115 contoh) — jadi tetap UNKNOWN. *(Ralat butir 69: rumusan
itu terlalu luas — di bawah `PropertyItemList`/`DeductibleList`/`AnekaList` ada halaman `Currency`. Yang benar: nol
halaman `Currency` sebagai anak langsung `LocationList` maupun `LocationList/Property` — 0/4 fixture, 0/297 korpus;
`Currency` di akar ada, 3/5 dan 41/115, tetapi `Name`-nya tidak pernah terisi.)* 68.5: **semua** medan tak
terpetakan masuk penampung (ADR-0023 akibat 1 berlaku umum), termasuk penunjuk `Idx*`/`Index*` V-47 yang belum
dikonversi, agar nilainya tidak hilang. ⚠️ Usulan 68.2 **melampaui** ADR-0034: tabel ADR itu menyimpan uang di kolom basis
data sebagai desimal berskala tetap — dicatat terang di usulan.

**Keputusan agent — menunggu konfirmasi (sesudah butir 68):**

| # | Tiket | Keputusan agent | Dasar |
| --- | ---: | --- | --- |
| A66 | 22 | 68.3: kode `CURRENCY_CODE` milik baris sendiri yang **berbeda** dari `Currency/Name` anaknya, atau dua halaman `Currency` berkode berbeda di bawah satu baris → `ErrKolomGanda` (penawaran dilewati menurut 68.8) | 68.3 tidak menyebut keadaan itu; memilih salah satu = menebak. Temuan code review sumbu spec · **dikonfirmasi work owner 02-10-2026 (butir 69)** |

**Ralat V-30 (68.6):** teks V-30 menyebut "23 tabel faktor"; `[terverifikasi]` data memuat **24** objek faktor per
`DataScoringRiskList` (kelima fixture ber-skoring: 2 × 24; termasuk `Others/TotalSumInsured`, yang V-28 sebut "faktor
skoring FIRE"). Peta "23 faktor" (lembar ScoringRisk) tidak ada di kedua workbook. Yang berlaku: data.

**Ralat jendela (verifikasi independen sesi `nusantarare-0f`, dicek ulang agent):** (a) "dua kolom" vs "delapan kolom"
konflik tipe — benar keduanya, jendela berbeda: 5 fixture **2** kolom (`PPN_CHECK` 5/5, `SHARE_OF_CEDING` 4 — di
`edm-fire-1` hanya terisi di bawah `OldData`), 115 contoh korpus **8** kolom; (b) bukti A60 "`CurrencyList.ID` angka
lima digit" hanya dari **korpus** — di fixture `CurrencyList` tidak punya medan `ID`; (c) `NET_RATE` > 38 digit 185 di
korpus (wadah `LocationList/Property/PropertyItemList/CoverageList`), **0** di fixture; fixture `edm-fire-1` > 38 digit
di `T_CURRENCYLIST.PAY_EDM_PREMI_MENJADI`, `PAY_NET_PREMIUM`, `SUM_TOTAL_PAYMENT` (masing-masing 1). Ketiganya
diperiksa ulang agent 02-10-2026 dan cocok. *(Butir (c) kemudian diralat butir 69: angka itu koefisien mentah termasuk
nol ujung; dalam digit bermakna hanya `PAY_NET_PREMIUM` — 3 di korpus, 1 di fixture.)* **Instrumen agent yang sempat keliru:** ukuran "> 8 desimal" putaran pertama
ikut menghitung nol di belakang (`0.000…0`); usulan memakai putaran kedua (koefisien dinormalkan).

## Keputusan work owner — 2 Oktober 2026, verifikasi loader putaran 2, diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` dengan kutipan jawaban work owner **"mau"** atas rekomendasinya. Setiap klaim
verifikasinya diperiksa ulang agent ke berkas sebelum ditulis.

| # | Butir | Keputusan |
| ---: | --- | --- |
| 69.1 | A66 | **Dikonfirmasi** — berhenti keras; saat muat dilewati dan dicatat (68.8) |
| 69.2 | Tempat penampung di Oracle | **Tidak disimpan di Oracle** (ADR-0023: penampung bukan tempat simpan akhir). Setiap medannya harus mendapat kolom atau keputusan "dibuang" eksplisit. Untuk `CoverageInitial` dan `AdditionalShip` agent mengajukan usulan kolom — **tidak** menambah kolom sendiri (`docs/USULAN-KOLOM-PENAMPUNG.md`) |
| 69.3 | Kembaran FR `T_FR_COVERAGELIST` / `T_FR_ANEKALIST` | **Ikut 68.3** (kode dari `Currency/Name`), supaya konsisten |
| 69.4 | `IsCedingConfirm` | Jalankan K-069 (7b) "kolom sendiri"; bila K-069 tidak menyebut tabel/tipe, tulis usulan, tandai menunggu work owner, jangan menebak. Migrasi tetap tidak ditulis/dijalankan |
| 69.5 | Presisi | Tetap tim inti, sesudah tabel B usulan dibetulkan |
| 69.6 | J-16 | Tetap pertanyaan terbuka untuk work owner |

**Cara agent menerapkannya:** 69.4 — `[terverifikasi]` K-069 (7b) hanya berbunyi "kolom sendiri", tanpa tabel dan
tipe. *(Ralat butir 70: catatan semula "rujukan L3788–3794 di pesan sesi 0f tidak cocok dengan berkas yang dibaca agent"
**keliru sebabnya** — ada **dua salinan** register. Agent membaca salinan lama `D:\migrasi\RNM\OUTPUT\` (K-069 7b di
L3770–3782); sesi 0f mengutip salinan repo `docs/00-KEPUTUSAN-WORK-OWNER.md` (L3788–3798). Kutipannya benar.)* Kolom karena itu **belum** ditambahkan ke peta; yang dijalankan: `IsCedingConfirm` **tidak lagi terbuang** bersama
`ViewSuggest` (V-31) — masuk penampung (`medanDiselamatkan`), fixture 34 nilai = pengurai Python independen; usulan P4
memuat tiga pilihan tabel. 69.3 — kembaran FR membaca baris `T_FR_CURRENCY`; `T_FR_SPREADINGLIST` ikut mewarisi.

**Ralat butir 68.8** — rekomendasi sesi 0f yang lama "medan tak terpetakan masuk penampung, **bukan alasan menolak
muat**" **kurang tepat**, menurut koreksi sesi 0f sendiri: ADR-0023 "Penampung itu wajib kosong sebelum pekerjaan
dinyatakan selesai" dan akibat 2 "Penampung berisi = pekerjaan belum selesai". Yang berlaku: **muat untuk uji boleh**
dengan penampung berisi; **muat produksi / fase 1 baru boleh dinyatakan selesai bila penampung kosong** (tiket 24).

**Ralat tabel B usulan presisi** — `[terverifikasi]` diukur ulang agent: cara hitung lama memakai **koefisien mentah
termasuk nol di belakang koma** (217 nilai korpus, 3 fixture). Dengan nol ujung dibuang hanya
`T_CURRENCYLIST.PAY_NET_PREMIUM` yang sungguh > 38 digit bermakna: **3** di korpus, **1** di fixture. Kalimat
"`NUMBER` polos pun membulatkannya" keliru untuk 214 dari 217 nilai itu (`[dugaan]` atas perilaku Oracle, tidak diuji
ke Oracle). `Diagnostik.LebihDari38Digit` Flatten diralat sama.

**Ralat agent — artefak instrumen:** medan korpus "tak ada di rancangan" `LocationList.TableOfLimit` (186),
`…OfferFacIn` (22), `LocationList/Property.Property` (3), `…SurveyAgent` (2) yang dilaporkan sebelumnya **bukan data**.
`[terverifikasi]` semuanya elemen XML kosong berisi spasi saja (316 / 34 / 4 / 2), yang diubah pengubah XML→JSON agent
menjadi medan bernilai spasi.

## Keputusan work owner — 2 Oktober 2026, usulan kolom penampung, diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` dengan kutipan jawaban work owner **"setuju"** atas rekomendasinya, atas
`docs/USULAN-KOLOM-PENAMPUNG.md`.

| # | Butir | Keputusan |
| ---: | --- | --- |
| 70.P1 | `CoverageInitial` | **Setuju** `T_COVERAGELIST.COVERAGE_INITIAL VARCHAR2(500)` |
| 70.P2–P3 | `AdditionalShip` | **Setuju** tabel baru `T_ADDITIONALSHIP` (anak `T_SHIP`, berulang: `SEQ_NO`, `ROW_UID`; `DWT`/`GRT`/`NRT` `VARCHAR2(50)`, `ADDITIONAL_SHIP_REF_ID VARCHAR2(50)`) — **bukan** dilipat ke `T_SHIP`, karena di sumber ia daftar berulang (`AdditionalShip[n]`) |
| 70.P4 | `IsCedingConfirm` | **Pilihan (b)**: kolom sendiri di `POOLDATA.HISTORYAKSEPTASIPRODUCTION`, di samping `POSISI`. ⚠️ **Tafsiran** atas K-069 (7b), bukan kutipan — lihat catatan di bawah. Tabel lama: **menunggu DBA** (nama kolom + tipe); agent **tidak** menulis DDL/migrasi |
| 70.P5 | `CurrencyList.ID` | **Setuju** `T_CURRENCYLIST.CURRENCY_REF_ID VARCHAR2(50)` |
| 70.P6 | FR `CurrencyList/Policy.TSI` | **Setuju** `T_FR_CURRENCYLIST.POLICY_TSI`, tipe mengikuti keputusan presisi tim inti |
| 70.P7 | penunjuk `Idx*`/`Index*` | **Setuju** klasifikasi menurut V-47: penunjuk induk langsung dibuang; penunjuk leluhur → FK, menunggu aturan |
| 70.R | Rancangan | DDL draf / workbook `D:\migrasi\RNM\OUTPUT\08-flat\` tetap **READ-ONLY**; kolom/tabel baru dicatat sebagai **amandemen rancangan** (tiket 22/23, register). Flatten + generator dikembangkan untuk P1, P2–P3, P5, P6, P7; P4 tetap di penampung sampai DBA menjawab |

**Catatan P4 — dasar tafsiran** (salinan repo `docs/00-KEPUTUSAN-WORK-OWNER.md`): pertanyaan asal 7b (L3790–3791)
adalah usul V-32 memetakan `IsCedingConfirm` ke `HISTORYAKSEPTASIPRODUCTION.POSISI`; keputusannya (L3793) "kolom
sendiri", dengan alasan "satu kolom tidak boleh memikul dua arti" (L3795); L3796–3797 menyebut juga
`T_WORK_POLIS.POSISI` — kedua tabel itu yang tersangkut. Alasan menolak pilihan lain, dari usulan agent sendiri: (a)
`T_WORK_POLIS.IS_CEDING_CONFIRM` menyimpan satu nilai per polis sehingga **riwayat per baris hilang**; (c) tabel riwayat
flat baru bertentangan dengan V-31 — `[terverifikasi]` dibaca ulang 02-10-2026: "riwayat akseptasi memakai tabel lama
POOLDATA.HISTORYAKSEPTASIPRODUCTION, jadi T_VIEWSUGGEST tidak dibuat". ⚠️ Teks V-31 **tidak ada** di register (salinan
mana pun); sumbernya lembar BACA-INI `D:\migrasi\RNM\Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx`, yang tidak punya
salinan di repositori.

**Cara agent menerapkannya** (tiket 22 bab *Amandemen rancangan*): `loader/amandemen.go` menggabungkan kolom/tabel baru
ke skema bangkitan saat paket dimuat — `skema_gen.go` dan workbook tidak disunting. P6: tipe ditulis `NUMBER` tanpa
presisi sebagai **penanda** "menunggu tim inti"; nilainya desimal eksak. *(Penerapan P7 di bawah **dicabut butir 71**: tidak satu pun penunjuk dibuang.)* P7: generator kini juga membaca lembar
**Kandidat Hapus** (`Claude outputs`), baris penunjuk — status "SUDAH DIHAPUS" → dibuang (yang berstatus tunggal ini
hanya **R1 indeks-diri**, 13 kunci — A67); "SUDAH JADI FK (V-47)" (R3) → dibuang bila sasarannya induk langsung (A68),
selain itu tetap di penampung. `[terverifikasi]` lembar itu **berselisih** untuk **dua** kunci —
`T_FR_ANEKALIST.IdxOccupation` dan `T_FR_DEDUCTIBLELIST.IndexProperty` (R2 "SUDAH DIHAPUS" dan R3 "SUDAH JADI FK") —
tidak dipilih, tetap di penampung. R3b "DIPERTAHANKAN" (tiga kunci) sudah berkolom kunci di rancangan
(`IDX_LOCATION`/`INDEX_LOCATION`), jadi tidak masuk penampung. `[dugaan]` tabel sasaran penunjuk dibaca dari namanya, melanjutkan contoh V-47 ("dan seterusnya").

**Ralat sumber — dua salinan register** (`[terverifikasi]` 02-10-2026): salinan repo
`APP_RNM/modul/nbfacin/docs/00-KEPUTUSAN-WORK-OWNER.md` (4.383 baris, 1 Okt 17:34) **lebih baru** daripada
`D:\migrasi\RNM\OUTPUT\00-KEPUTUSAN-WORK-OWNER.md` (4.365 baris, 25 Sep): `diff` (tanpa CR) — 18 baris hanya di salinan
repo, tiga sisipan (sesudah L794, L829, L1263 salinan lama) = amandemen 1 Oktober butir 30, 34/A4, 46. **Mulai sekarang
register dibaca HANYA dari salinan repo**; `D:\migrasi\RNM\OUTPUT\` tetap untuk spec 11 dan berkas `08-flat`. Rujukan
nomor baris ke salinan lama yang ditemukan dan dibetulkan: catatan 69.4 di register ini; `docs/USULAN-KOLOM-PENAMPUNG.md`
P4; komentar `backend/services/loader/aturan.go` (`medanDiselamatkan`). Rujukan lain ke register di tiket 22–24 dan register
ini tidak memakai nomor baris (K-0xx saja), dan isi K-0xx yang dirujuk tidak termasuk 18 baris yang berbeda.

**Keputusan agent — menunggu konfirmasi (sesudah butir 70; temuan code review sumbu spec):**

| # | Tiket | Keputusan agent | Dasar |
| --- | ---: | --- | --- |
| A67 | 22 | ~~P7: penunjuk **R1 indeks-diri** berstatus "SUDAH DIHAPUS" (13 kunci) ikut **dibuang**~~ → **dicabut butir 71.2** (korpus 39 beda) | Teks P7 hanya "penunjuk induk langsung dibuang"; dasar R1 adalah V-50/V-41 (lembar Kandidat Hapus "SUDAH DIHAPUS", "sudah diwakili kolom urutan baris") — keputusan rancangan, bukan teks butir 70. 69.2 menuntut keputusan "dibuang" eksplisit; ini dibaca sebagai eksplisit |
| A68 | 22 | ~~P7: penunjuk R3 yang sasarannya **induk langsung** (16 kunci) **dibuang**~~ → **dicabut butir 71.1** (fixture 103 beda, korpus 8), walau **15** di antaranya punya kolom FK di rancangan (mis. `T_COVERAGELIST.IndexCargo` → `CARGO_ID`, `T_PROPERTYITEMLIST.IndexProperty` → `PROPERTY_ID`) | ⚠️ **Pertentangan** teks P7/V-47 ("INDUK LANGSUNG dibuang karena PARENT_ID sudah menyatakan hal yang sama") dengan rancangan (kolom FK untuk induk langsung, terutama di tabel berinduk ganda). Nilai FK itu dapat diturunkan dari `PARENT_ID` + `PARENT_TABLE` bila aturan isi FK ditulis — jadi tidak hilang; tetapi **nilai penunjuk mentahnya** hilang. Bila work owner memilih menyimpannya, 15 kunci itu kembali ke penampung (satu baris kode) |

## Keputusan work owner — 2 Oktober 2026, verifikasi butir 70, diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` dengan kutipan jawaban work owner **"setuju"** atas rekomendasinya. Setiap
angka diperiksa ulang agent dua cara (mesin Go vs pengurai Python atas pohon mentah) sebelum ditulis.

| # | Butir | Keputusan |
| ---: | --- | --- |
| 71.1 | A68 — penunjuk induk langsung | **Jangan dibuang dulu.** Ke-16 kunci (15 berkolom FK di rancangan) kembali ke penampung. Data membantah premis V-47 (BACA-INI **r32**: "Yang menunjuk INDUK LANGSUNG dibuang karena PARENT_ID sudah menyatakan hal yang sama"). Diajukan ke work owner: apakah V-47 diubah. Arti penunjuk **tidak** ditafsirkan |
| 71.2 | A67 — R1 indeks-diri (Kandidat Hapus r2–r14) | **Ukur dulu**; bila selalu sama dengan posisi baris → boleh dibuang, bila ada yang beda → kembali ke penampung. Catat: V-50 (BACA-INI **r19**) hanya mencatat status per kelompok, tidak menyebut R1 khusus |
| 71.3 | Dua kunci berselisih | Tetap di penampung |

**Hasil pengukuran** (02-10-2026; jendela fixture 5 kasus dan korpus 115 contoh, cabang yang dibuang tidak ikut; "posisi"
= `SEQ_NO` = posisi di larik sumber, 1-based):

| | Fixture | Korpus |
| --- | --- | --- |
| A68 induk langsung — sama dengan posisi induk | 105 | 5.658 |
| A68 — **beda** | **103** — seluruhnya `T_COVERAGELIST.IndexCargo` di kasus marine (104 unsur, 1 sama; 1 lagi sama bila 0-based) | **8** — `T_ANEKALIST.IdxOccupation` 2, `T_COVERAGELIST.IndexCargo` 1, `T_DEDUCTIBLELIST.IndexCoverage` 5 |
| A68 — induk berupa halaman tanpa `SEQ_NO` (`PropertyItemList.IndexProperty`, induknya `Property`) | 11 | 472 |
| A68 — jumlah | 219 | 6.138 (= jumlah yang dibuang butir 70) |
| A67 R1 — sama dengan `SEQ_NO` sendiri | 227 | 6.094 |
| A67 R1 — **beda** | **0** | **39** — `IndexCoverage` 18, `OccupationList.IdxOccupation` 16, `IndexDeductible` 5 |

⚠️ **Dua definisi, dua angka — dicatat, tidak dipilih:** sesi 0f menghitung "sama" fixture **116** dan korpus **6.130**;
agent **105** dan **5.658**. Selisihnya tepat baris yang induknya halaman tanpa `SEQ_NO` (11 dan 472): sesi 0f
menghitungnya "sama", agent "tak terbandingkan". Angka "beda" kedua pihak sama (fixture 103 vs sesi 0f "cocok 1 dari 104"
di marine; korpus 8). ⚠️ **Cara Go per kunci sempat keliru** untuk `T_FR_CARGOLIST.IdxFacRetro` (5, karena ikut jalur
`FacOfferList/CargoList` yang induknya bukan `FacRetroList`); dibatasi ke jalur yang induknya tabel sasaran → 3, sama
dengan Python.

**Akibatnya:** A68 — penunjuk induk langsung ke penampung. A67 — **ada yang beda** (korpus 39) → R1 juga ke penampung
(seluruh 13 kunci; rincian per kunci untuk keputusan per kunci bila work owner mau). **Tidak satu pun penunjuk dibuang**
lagi; P7 butir 70 dicabut penerapannya. Bukti dikunci `TestKasusPenunjukBukanPosisi` (fixture: `IndexCargo` 1/103, R1
227/0).

**Ralat (verifikasi sesi 0f, diperiksa ulang agent):** (1) "penunjuk leluhur V-47 8.752" (sisa butir 70) = **8.723**
penunjuk leluhur R3 + **25** dua kunci berselisih + **4** `T_FR_PERSONLIST.IdxPerson` yang **tidak tercantum** di lembar
Kandidat Hapus. (2) "0 di fixture" untuk `CurrencyList.ID` dan FR `Policy.TSI` tidak membuktikan apa pun — kedua medan
**tidak ada** di fixture; buktinya hanya korpus (5 dan 15 kemunculan di luar `OldData`, penampung 0). (3)
`mutasi_loader.log` butir 70 lebih tua dari suntingan terakhir `aturan.go`; mutasi dijalankan ulang sesudah perubahan
butir 71 — **60/60**. Rujukan BACA-INI sesi 0f (r19, r32, r53) `[terverifikasi]` = baris Excel; pembaca agent menomori
elemen `<row>` (17, 30, 51) karena Excel melewatkan baris kosong.

## Keputusan work owner — 2 Oktober 2026, satu aturan untuk semua penunjuk, diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` dengan kutipan jawaban work owner **"setuju"** atas rekomendasinya, sesudah
verifikasi independen butir 71 (semua klaim terbukti, dua cara).

| # | Butir | Keputusan |
| ---: | --- | --- |
| 72.1 | **V-47 diubah** — penunjuk induk langsung (16 kunci) | **Tidak dibuang**; disimpan APA ADANYA sebagai kolom teks mentah di tabel barisnya, tidak ditafsirkan. Amandemen V-47 (BACA-INI r32) dengan dasar data A68: premis "PARENT_ID menyatakan hal yang sama" dibantah — fixture 103 beda, korpus 8 beda |
| 72.2 | R1 indeks-diri (13 kunci) | Sama — kolom teks mentah; **tidak** diputuskan per kunci. Dasar: A67 korpus 39 beda |
| 72.3 | Penunjuk leluhur | Sama — kolom teks mentah |
| 72.4 | Dua kunci berselisih, `T_FR_PERSONLIST.IdxPerson` | Sama — kolom teks mentah |
| 72.5 | Kolom FK V-47 | **Tidak berubah** — tetap kosong sampai aturan isinya tertulis (68.4). Kolom penunjuk mentah **bukan** pengganti FK dan **tidak** dipakai mengisinya |

**Fakta tambahan (verifikasi sesi 0f, diperiksa ulang agent dengan Python, tidak ditafsirkan):** ke-483 penunjuk
`IndexProperty` yang induknya `T_PROPERTY` / `T_FR_PROPERTY` (fixture 11, korpus 470 + 2 FR) **selalu sama** dengan posisi
unsur `LocationList` kakeknya (483 / 483); nilainya 1..181, 181 nilai berbeda — bukan konstan.

**Penerapan** (`loader/amandemen.go`, `amandemenPenunjuk`): **48 kolom** = tepat 48 kunci penunjuk yang teramati di
penampung, `[terverifikasi]` 5 fixture + 115 contoh korpus (47 kunci lembar Kandidat Hapus + `T_FR_PERSONLIST.IdxPerson`;
tiga kunci lembar lainnya — R3b — sudah berkolom kunci di rancangan). Per kelompok: R1 indeks-diri **13**, R3 induk
langsung **16**, R3 leluhur **16**, berselisih **2**, `IdxPerson` **1**. Nama: SNAKE_CASE dari FIELD ASLI — pola kolom kunci
rancangan yang sudah ada (`IdxLocation` → `IDX_LOCATION`, `IndexPropertyItem` → `INDEX_PROPERTY_ITEM`); nama terpanjang 19
karakter (V-20). Tipe: **`VARCHAR2(50)`** — bentuk nilai terukur seluruhnya bilangan bulat, terpanjang **3** karakter;
"teks mentah" menurut keputusan, pola V-6 kolom kode. ⚠️ Pola V-49 (`*_REF_ID`) **tidak** dipakai: ia untuk medan Pega
bernama `ID`, bukan penunjuk. Sebaran per tabel:

| Tabel | Kolom baru |
| --- | ---: |
| `T_COVERAGELIST` | 8 |
| `T_DEDUCTIBLELIST` | 6 |
| `T_FR_COVERAGELIST` | 6 |
| `T_FR_DEDUCTIBLELIST` | 4 |
| `T_ANEKALIST`, `T_FR_ANEKALIST`, `T_SPREADINGLIST` | 3 masing-masing |
| `T_ADDITIONALCOVERAGE`, `T_FR_CARGOLIST`, `T_FR_PERSONLIST`, `T_FR_PROPERTYITEMLIST`, `T_OCCUPATIONLIST`, `T_PROPERTYITEMLIST` | 2 masing-masing |
| `T_CARGOLIST`, `T_FR_OCCUPATIONLIST`, `T_VEHICLELIST` | 1 masing-masing |
| **16 tabel** | **48** |

Hasil: kelima fixture — terpetakan naik tepat sebanyak penunjuknya (600); korpus 21.023 nilai penunjuk masuk kolom.
**Sisa penampung: hanya `IsCedingConfirm`** — fixture 34, korpus 805 (menunggu DBA; P4). Di korpus juga 214 entri
bernilai spasi — artefak pengubah XML→JSON agent, bukan data (butir 69). Skema sesudah amandemen: 79 tabel, **1.390**
kolom (1.329 + 61). Bukti `IndexCargo` dikunci `TestKasusPenunjukBukanPosisi`: teks "7" tersimpan apa adanya di 104
baris, sama dengan posisi induk 1 dari 104.

## Keputusan work owner — 2 Oktober 2026, penjaga claimlife (AskUserQuestion)

| # | Pertanyaan | Jawaban |
| ---: | --- | --- |
| 67 | Penjaga claimlife `TestKolomTakDibawaHanyaAdaDiKatalog` (lingkup tertulis claimlife, pindaian seluruh `APP_RNM`) menuduh loader NB atas `STS_KONVERSI`/`TGL_KONVERSI` | **"FOKUS KE NB FACIN SAJA"** — claimlife tidak disunting; `go test ./...` aplikasi tetap merah di paket itu sampai pemiliknya memutuskan |

## Keputusan work owner — 2 Oktober 2026, lookup akun tiket 27, diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` dengan kutipan jawaban work owner **"harus peka besar kecil dong, 15 baris per
halaman"** atas keputusan agent A69–A73 (`issues/27-lookup-akun-choose-account.md`).

| # | Butir | Keputusan |
| ---: | --- | --- |
| 73.1 | A70 — peka huruf besar-kecil pencarian `GET /api/nbfacin/account` | **DIUBAH: PEKA huruf.** "Mengandung" tanpa `UPPER` di `INSUREDID`, `INSUREDNAME`, `GROUPBUSINESS` apa adanya; escape `\` `%` `_` tetap |
| 73.2 | A71 — ukuran halaman | **DIUBAH: 15** (semula 20) |
| 73.3 | A69 (tiga kolom dicari), A72 (urut `INSUREDID`, `ID`), A73 (`cari` > 255 → 400) | **Belum diputus** — tetap keputusan agent, menunggu konfirmasi |

**Penerapan:** `repository/akun.go` (`PolaCari` tanpa `ToUpper`, `saringAkun` tanpa `UPPER`), `services/layanan.go`
(`UkuranHalamanAkun = 15`); uji `TestSQLAkun` (menolak `UPPER`/`LOWER`), `TestPolaCari` (`uji` ≠ `UJI`), `TestCariAkun`
services (halaman 3 → offset 30) dan handlers (`"ukuran":15`). `[dugaan]` peka-huruf `LIKE` bergantung `NLS_COMP` instance
— `belum terverifikasi`.

## Keputusan agent A74–A78 — tiket 28 Class Of Business, menunggu konfirmasi

Atas brief sesi `nusantarare-0f` dan DDL `BUSINESS.txt` dari work owner (`issues/28-pilihan-class-of-business.md`).
`GET /api/nbfacin/class-of-business?groupBusinessId=` → `{"baris":[{"id","note"}]}`, semua baris, tanpa paging.

| # | Keputusan | Dasar |
| --- | --- | --- |
| A74 | Urut `NOTE`, **bukan** `.ID` DESC seperti RD `BrowseBusiness_RD` | brief sesi 0f; tangkapan layar Pega tampak alfabetis NOTE |
| A75 | `ID` pemutus seri sesudah `NOTE` | urutan deterministik (pola A72) |
| A76 | `groupBusinessId` kosong/spasi atau > 4000 byte → 400; nilai tidak dipangkas | brief (kosong); lebar `BUSINESSGROUPID` (pola A73) |
| A77 | Hanya filter C RD (`.BusinessGroupID = Param.Group`) yang dibangun; filter A/B tidak | brief "semua baris"; di satu-satunya section pemakai di folder `NB FacIn` (`InputLossRecord_Sec`) parameter `ID`/`Note`/`Group` kosong (`RNW Fac In` punya salinan bernama sama, tidak dibandingkan) |
| A78 | `NOTE IS NOT NULL` | brief sesi 0f; **tidak ada di RD** — selisih dengan Pega, dicatat |

**Ralat atas brief** `[terverifikasi]` `BrowseBusiness_RD.xml`: filter B adalah `.Note Contains Param.Note` dengan
`pyCaseInsensitive` true — bukan `=`. `[dugaan]` `groupBusinessId` = `T_M_ACCOUNT.GROUPBUSINESSID` akun terpilih — belum
terverifikasi dari korpus (section form Opportunity pemakai RD tidak ada di folder `NB FacIn`).

## Keputusan work owner — 3 Oktober 2026, Create opportunity (tiket 29), diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` (pilihan work owner lewat AskUserQuestion di sesi itu). Konteks: tombol
`Create opportunity` di form Opportunity (tiket 29, frontend + tiket di-commit sesi 0f, `4064d0b`) memanggil kontrak
`POST /api/nbfacin/opportunity` → 201 `{"caseId":"NB-…"}`; badan = `IsianOpportunity` (`frontend/api.ts`), tanggal bentuk
kabel `DD-MM-YYYY`. ⛔ **Endpoint belum dibangun** — dicatat saja, atas permintaan sesi 0f ("jangan dibangun dulu").
⛔ **74.1 DIRALAT butir 76** (3 Oktober 2026) — endpoint kini dibangun; lihat bab butir 76.

| # | Butir | Pilihan work owner |
| ---: | --- | --- |
| 74.1 | Penyimpanan opportunity/case NB | **"Tunggu tabel flat (tiket 23)"** — **tidak ada** tabel sementara. `POST /api/nbfacin/opportunity` baru dibangun sesudah tabel flat tiket 23 ada; tiket 23 sendiri masih ⏸ ditahan (butir 66, presisi tim inti). → ⛔ **DIRALAT butir 76**: pilihan ini diambil atas premis keliru yang diteruskan sesi 0f (penyimpanan case dikira seluruhnya tertahan tiket 23), padahal `T_WORK_POLIS` sudah ada (premiumlistlife 050) |
| 74.2 | Nomor case NB | **"Lanjut dari nomor terakhir Pega"** — penghitung mulai dari MAX nomor NB yang ada + 1; angka pastinya diisi saat migrasi dijalankan oleh work owner/DBA (⛔ agent tidak menjalankan `-migrate`, butir 65) |
| 74.3 | Layar sesudah create ("lari ke flow nya") | Flow action `InwardFacultative` → section `InputInwardFacultative` — port berikutnya, **belum dimulai** |

**Bukti yang dikutip sesi 0f:**
- `[terverifikasi]` oleh agent ini: ADR-0043 (`docs/bersama/adr/0043-penomoran-di-aplikasi.md`) — penomoran di aplikasi,
  `SELECT … FOR UPDATE`; tiket 29 dan commit `4064d0b` ada.
- Nomor NB = satu deret global tanpa tahun, sampai `NB-184351` (nama berkas ekspor kasus di `D:\XML NURE\Groupbusiness\`)
  — **klaim sesi 0f, tidak diverifikasi ulang agent ini** (folder ekspor kasus tidak dibuka).

⚠️ **Terbuka, tidak ditafsirkan:** ADR-0043 mengunci penghitung atas `(CLASS, JENIS, TAHUN)` dengan aturan periode/tahun,
sedangkan deret NB menurut klaim di atas **tanpa tahun**. Bagaimana nomor NB dipetakan ke kunci itu belum diputuskan —
`belum terverifikasi`.

## Keputusan work owner — 3 Oktober 2026, lookup akun kembali tidak peka huruf (tiket 27), diteruskan sesi `nusantarare-0f`

Diteruskan sesi agent `nusantarare-0f` dengan kutipan work owner (disertai tangkapan layar kotak Search popup
ChooseAccount): **"pada saat search Group Business, itukan ada isian untuk search, itu buatin tanpa liat huruf besar atau
kecil"**.

| # | Butir | Keputusan |
| ---: | --- | --- |
| 75 | Peka huruf pencarian `GET /api/nbfacin/account` | **TIDAK peka huruf** — `UPPER(kolom) LIKE` pola huruf besar, `ESCAPE` tetap, parameter terikat, atas `INSUREDID`/`INSUREDNAME`/`GROUPBUSINESS`. **Membatalkan 73.1** — keputusan lama dikutip utuh: **"harus peka besar kecil dong, 15 baris per halaman"** (73.1 "DIUBAH: PEKA huruf"). Bagian "15 baris per halaman" (**73.2**) **tetap** |

**Penerapan:** `repository/akun.go` (`PolaCari` kembali `strings.ToUpper`, `saringAkun` kembali `UPPER(kolom)`); uji
`TestSQLAkun` (tepat tiga `UPPER(`), `TestPolaCari` (`uji`/`Uji`/`UJI` berpola sama). `[dugaan]` lama tetap: huruf besar
dibuat `strings.ToUpper` Go sedangkan kolom `UPPER` Oracle — untuk huruf non-ASCII keduanya dapat berbeda.

## Keputusan work owner — 3 Oktober 2026, case NB di `T_WORK_POLIS` yang ada (tiket 29 backend, tiket 23)

Rencana sesi `nusantarare-0f` dijawab work owner **"setuju"** (diteruskan sesi 0f): case NB dibuat sekarang di
`T_WORK_POLIS` yang ada, tiket 23 diselaraskan. Butir yang tidak boleh ditebak ditanyakan langsung ke work owner di sesi
ini (AskUserQuestion, 3 Oktober 2026); pilihannya dikutip.

| # | Butir | Pilihan work owner |
| ---: | --- | --- |
| 76.1 | `T_WORK_POLIS.ID`: tabel yang ada `VARCHAR2(32)` berisi pengenal work, rancangan `NUMBER` + `IDPEGA` | **"Follow existing table"** — ID = `NB-<n>` `VARCHAR2(32)`; `T_GENERAL_POLIS.ID` dan `PARENT_ID` 10 tabel anaknya ikut `VARCHAR2(32)`; `IDPEGA`/`JENIS_WORK`/`NO_WORK` tetap kolom tambahan; loader + uji disesuaikan |
| 76.2 | Nilai `LINI` Fac In | **"FAC"** — preseden `T_WORK_CLAIM.LINI` (K-064) |
| 76.3 | Penyimpanan isian opportunity | **"Own opportunity table"** — `T_NB_OPPORTUNITY` (migrasi 180), 1:1 dengan `T_WORK_POLIS.ID`. Dasar `[terverifikasi]`: opportunity di Pega kelas work tersendiri (`NB FacIn\ReportDefinition\GetListOpportunity.xml`: `ASM-FW-SFAGISFW-Work-Opportunity` INNER JOIN `ASM-FW-GISFW-Work-NB` pada `A.pzInsKey = .NBHandle`); rancangan flat tidak punya tabelnya |
| 76.4 | Kolom rancangan yang bertumpuk dengan kolom yang ada | **"Merge into existing"** — `POSISI`→`POSITION`, `STATUS_PROSES`→`STATUS_WORK`, `TGL_INPUT`→`TGL_CREATE`, `USERNAME`→`CREATE_OP` |
| 76.5 | Angka awal `SEQ_WORK_POLIS_NB` (tidak ada di repo) | **"Placeholder in migration"** — `181` memakai `START WITH {NB_MULAI}`; Oracle menolaknya sampai work owner/DBA menggantinya dengan (nomor NB Pega terakhir + 1); `-migrate` seluruh aplikasi berhenti di langkah itu sampai diisi — disengaja |

**Penerapan:** `loader/amandemen.go` `selaraskanWorkPolis` + `aturan.go` (`TestAmandemenWorkPolis`; skema 79 tabel /
**1.391** kolom; mutasi **67/67**); migrasi `180_t_nb_opportunity.sql`, `181_seq_work_polis_nb.sql` (+ `_down`) — **ditulis,
tidak dijalankan**; `repository/casenb.go`, `services/opportunity.go`, `handlers` `POST /api/nbfacin/opportunity`
→ 201 `{"caseId":"NB-<n>"}`. Pemetaan 18 kolom: tiket 23 bab *Penyelarasan*.

**Keputusan agent — menunggu konfirmasi:**

| # | Keputusan | Dasar |
| --- | --- | --- |
| A79 | `T_NB_OPPORTUNITY` berbagi PK dengan `T_WORK_POLIS` **tanpa** constraint FK | pola `T_PREMIUM_LIST` premiumlistlife 050/051; constraint lintas modul tidak diputuskan siapa pun |
| ~~A80~~ | ~~Baris `T_WORK_POLIS` NB: `POSITION`, `STATUS_WORK`, `FLAG_ONGOING_POLICY`, `COVER_KEY` **kosong**~~ → **DIGANTI butir 77** (`COVER_KEY` tetap kosong) | nilai awal case NB di Pega `belum terverifikasi` |
| A81 | Tipe `T_NB_OPPORTUNITY`: salinan `T_M_ACCOUNT` bertipe sumber (CHAR), `CLASS_OF_BUSINESS` 4000 BYTE (`BUSINESS.NOTE`), teks lain 255, `DESCRIPTION` 4000; isian melebihi lebar → 400 | pola A73/A76 |
| A82 | Tanpa identitas → **401**; pembuat = pengenal akun di `CREATE_OP` **dan** `CREATE_OP_NAME`; sesi login didahulukan, stub `X-Pelaku` hanya bila `AUTH_STUB` | pola `KasusPolis.Buat` premiumlistlife, ADR-U-0030 |
| A83 | Isian disimpan **apa adanya** (tidak dipangkas; medan opsional berisi spasi saja tersimpan spasi, bukan NULL); wajib-isi = bukan hanya spasi; Type Of Facultative wajib bila Type Of Inward **persis** `Facultative`; tanggal `DD-MM-YYYY` ketat (31-02 ditolak), disimpan `DATE` tanpa jam | tiket 29; nilai `Facultative` `[dugaan]` dari frontend |
| A84 | Nomor dari **sequence** (`SEQ_WORK_POLIS_NB`, `NEXTVAL`), bukan tabel penghitung ADR-0043 `(CLASS, JENIS, TAHUN)` | brief sesi 0f; preseden `SEQ_WORK_POLIS` premiumlistlife; deret NB tanpa tahun |
| A85 | Kelahiran case NB **tidak** direkam ke jejak | satu-satunya penyimpan jejak (`inti/backend/jejak`) menulis `T_CLAIMLF_JEJAK` milik claimlife |

**Risiko yang ditemukan — menunggu work owner:**
- **R1** `[terverifikasi]` kotak masuk PremiumList Life membaca **seluruh** `T_WORK_POLIS` tanpa saringan `LINI`
  (`premiumlistlife/backend/repository/polis_inbox.go`: cacah `WHERE (:1 IS NULL OR w.POSITION = :1)`, daftar `LEFT JOIN`
  `T_PREMIUM_LIST`). Case `NB-…` (POSITION kosong) akan **muncul di tab "semua" kotak masuk Life**. endorsementlife
  (`481_seq_work_edm_life.sql`) sengaja tidak menulis `T_WORK_POLIS` karena alasan yang sama. Perbaikannya di modul
  premiumlistlife (saringan `LINI`) — **di luar lingkup sesi ini, tidak disunting**.
- **R2** nbfacin MENULIS `T_WORK_POLIS` milik premiumlistlife dengan SQL sendiri (`repository/casenb.go`, tanpa kontrak
  lintas modul) dan kelak `ALTER` kolom Fac In atasnya (tiket 23): perubahan premiumlistlife atas tabel itu (mis. kolom
  NOT NULL baru) dapat mematahkan NB tanpa peringatan — koordinasi pemilik / tim inti.

## Keputusan work owner — 3 Oktober 2026, keadaan awal case NB (tiket 29), diteruskan sesi `nusantarare-0f`

Diteruskan sesi `nusantarare-0f`; konteks: work owner sudah menjalankan migrasi 180/181 di DEV (181 diisi work owner:
`START WITH 151028`, belum di-commit saat dicatat) dan case NB berhasil dibuat. Kutipan: **"saat berhasil create position =
Offer, FLAG_ONGOING_POLICY=0, STATUS_WORK=Pending-Policy."**

| # | Butir | Keputusan |
| ---: | --- | --- |
| 77 | Keadaan awal baris `T_WORK_POLIS` case NB — **menggantikan A80** ("kosong") | `POSITION = 'Offer'`, `FLAG_ONGOING_POLICY = '0'`, `STATUS_WORK = 'Pending-Policy'` — teks verbatim, konstanta `PosisiAwalCaseNB` / `FlagAwalCaseNB` / `StatusAwalCaseNB` di `repository/casenb.go`; `COVER_KEY` tetap kosong |

**R1 melebar** `[terverifikasi]`: `'Offer'` adalah nilai posisi Life yang sama (`premiumlistlife/backend/models/polis_penawaran.go`
L112 `PosisiOffer`; kasus Life baru juga `Offer`, `models/polis_kasus.go` L84). Kotak masuk Life menyaring hanya
`w.POSITION` (`repository/polis_inbox.go` L111/L125, kueri `?posisi=` dari `handlers/rute_premiumlist.go` L53), jadi case NB
kini tampil di tab "semua" **dan** di saringan `posisi=Offer`, ber-status `Pending-Policy`. Layar Life hari ini hanya memanggil
tanpa posisi (`frontend/pages/InboxPremiumList.tsx` L99). `Pending-Policy`: nol kemunculan di premiumlistlife — hanya tampil,
bukan saringan. Tidak disunting (di luar lingkup).

⚠️ **Dicatat, tidak ditafsirkan:** angka awal 151028 **lebih kecil** dari contoh `NB-184351` yang dikutip sesi 0f dari
ekspor uji (butir 74, klaim sesi 0f, tidak diverifikasi agent ini). Bila nomor NB Pega yang ada memang mencapai 184351,
nomor baru dapat bertabrakan dengan kasus lama — `belum terverifikasi`, menunggu work owner/DBA.

## Keputusan work owner — 3 Oktober 2026, layar Inward Facultative tahap 2 (tiket 31)

Tugas diteruskan sesi `nusantarare-0f` (work owner "commit dan lanjut"); butir yang perlu tafsiran ditanyakan langsung
ke work owner di sesi ini (AskUserQuestion, 3 Oktober 2026). Bukti dan kontrak: `issues/31-baca-simpan-case-inward-general.md`.

| # | Butir | Pilihan work owner |
| ---: | --- | --- |
| 78.1 | Bentuk simpan Begin/Offering/End date (rancangan: teks Pega VARCHAR2(30)) | **"Pega text, WIB"** — Offering `YYYYMMDD`; baca GMT → WIB → tanggal. Jam Begin/End: semula ~~00:00 WIB (`…T170000.000 GMT` hari sebelumnya)~~ → **diralat: "05:00 GMT / 12:00 WIB"** (pertanyaan diajukan ulang setelah agent meralat hitungannya: nilai terkini fixture Start 050000 ×4 / 170000 ×1, End ×3 / ×2 — angka awal menjumlahkan salinan `OldData`, jebakan sensus 5) |
| 78.2 | Policy Type: data Pega berkode `0`/`1`/`2`, layar berlabel `Individual Policy`/`Master Policy`, pemetaan tidak ada di korpus | **"Label text for now"** — label disimpan apa adanya (penyimpangan); dikonversi ke kode bila pemetaan diketahui |
| 78.3 | Type facultative: data `FacultativeIn`, layar `Facultative In` | **"As sent by the screen"** — tanpa konversi tebakan |
| 78.4 | `T_GENERAL_POLIS` / `T_QUOTATIONDATA` tidak dapat dibuat utuh (presisi tim inti; panjang butir 68.1) | **"Create with needed columns"** — migrasi 182/183 kolom sistem + kolom layar; sisanya tiket 23 lewat `ALTER` |

**Penerapan:** migrasi `182_t_general_polis.sql`, `183_t_quotationdata.sql` (+ `SEQ_T_QUOTATIONDATA`) — **ditulis, tidak
dijalankan**; `GET /api/nbfacin/kasus/{caseId}`, `PUT /api/nbfacin/kasus/{caseId}/general`, `GET /api/nbfacin/marketing-officer`.
Keputusan agent A86–A94 (menunggu konfirmasi): tiket 31 bab *Keputusan agent*.

**Ralat atas brief** `[terverifikasi]` `Periode.xml`: Source of business = `.QuotationData.SobName`; Old Policy Number =
`.Following` (akar, `T_GENERAL_POLIS.FOLLOWING`). `PEGA_MARKETINGOFFICER.txt` adalah prosedur, tabelnya `MARKETINGOFFICER`.

## Tiket 32 — daftar case NB di portal Opportunity, diteruskan sesi `nusantarare-0f`

Permintaan work owner (dikutip sesi 0f): **"nb yang sudah di create, muncul disini … harus ada case id nya, group business,
insured name, status."** Dibangun `GET /api/nbfacin/opportunity` (`issues/32-daftar-case-nb-portal-opportunity.md`): hanya
`LINI = 'FAC'`, status = `STATUS_WORK` (bukan `NBStatusNew` grid Pega — selisih dicatat), cari tidak peka huruf atas case id +
Business Prospect Name. Keputusan agent A95–A98 menunggu konfirmasi; A97 mencatat lima selisih dengan
`GetListOpportunityF` `[terverifikasi]` (saringan pembuat/team group/Resolved tidak diterapkan). ⚠️ Rute butuh migrasi 183.

## Yang belum diputuskan


- Nama modul `facout` dan nomor usulan `760-799` / `990-991` — perlu persetujuan tim inti.
- ~~Mesin bersama NB … lewat `inti/backend/kontrak` (saran), atau ditaruh di `inti/`.~~ → ✅ **butir 55**: kontrak
  `inti/backend/kontrak/facin.go`. Tinggal: persetujuan tim inti atas berkas kontrak itu, dan penyambungan perakit (A32).
