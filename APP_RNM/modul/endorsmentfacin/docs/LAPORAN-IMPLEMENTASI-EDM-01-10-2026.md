# Laporan implementasi EDM — putaran "lanjut tanpa tanya" (01-10-2026)

> Atas `PROMPT-LANJUT-IMPLEMENT-TANPA-TANYA-EDM.md` (perintah work owner, 01-10-2026). Lanjutan
> `LAPORAN-IMPLEMENTASI-BEFORE-IMAGE.md` (commit `9ca6f06`). **Belum di-commit** — sesuai perintah.
> Korpus: `D:\migrasi\RNM\Endorsment Fac In\` (READ-ONLY). Label mengikuti `CLAUDE.md` §4.
> Nol nomor polis, nol nama orang, nol data medis di berkas ini dan di kode modul
> (`grep -rniE 'RNM-[A-Z0-9]' APP_RNM/modul/endorsmentfacin` → nol kecocokan).
>
> ⛔ **RALAT 01-10-2026** (temuan sesi nusantarare-c3, diperiksa ulang ke berkas): klaim "nol nama
> orang" di atas **sempat keliru**. `predikat/registry_gen.go` memuat 3 literal `MarketingName` di
> `ISTBONDING` (18/14/6 karakter, dua berspasi — hampir pasti nama orang). Yang keliru caranya:
> saringan dan grep penjaga hanya mencari POLA nomor polis, sedangkan nama tidak punya pola. Kini
> diperbaiki tiga lapis: generator nbfacin terbaru menolak medan identitas; `saring_registry.py`
> menyaring MEDAN identitas (`OldPolicyNo`, `PolicyNo`, `MarketingName`, `pyTelephone`,
> `pyUserIdentifier`) apa pun nilainya; penjaga Go `TestRegistryTanpaLiteralIdentitas` (tanpa korpus).
> Berkas itu belum pernah di-commit (`git log` — sesi ini tanpa commit).

## 1. Ringkasan per tiket

| Tiket | Status | Isi |
| --- | --- | --- |
| E06 | sebagian | + **rekonsiliasi eksak 78 nilai lapis B/C** terhadap Pega (fixture NB-15 `edm-fire-1.json`); porsi periode hari bulat, **HALF_UP** (A37 dikonfirmasi), galat dipisah (`GalatPorsiPeriode`); ⚠️ kasus nyata satu-satunya DITOLAK (selisih 150,5 hari) |
| E07 | selesai | — (putaran sebelumnya) |
| E08 / E09 | sebagian | sisa: satuan ‰/% rate (resolver NB) dan pembanding rekonsiliasi E22 |
| E10 | sebagian | korpus menyimpang dari tiket (laporan before-image §2) |
| E11 | selesai | — |
| **E03** | sebagian | `OpenCase` (Seam 5): kasus `EDM-`, tautan tiga arah, fase Policy, penolakan "sudah batal" + pengecualian konfigurasi |
| **E04** | sebagian | `PeriksaGerbangPenolakan`: 6 gerbang, 4 klep, bypass RI slip (dua lompatan), langkah 7–19 jenis bisnis |
| **E05** | sebagian | `ValidasiTanggalEndorsemen`: periode, pengecualian konfigurasi, lompat ke `END` |
| **E18** | sebagian | rute `Decision19` (Life melewati tangga), rencana + susunan nomor endorsement |
| **E20** | sebagian | paket `medis`: `ScoreMedical`, `PeriksaFisik`, `NilaiRujukanLab` — tabel dibangkitkan dari korpus, dijalankan penafsir |
| **E01** | sebagian | paket `predikat`: registry varian EDM (201 predikat, generator nbfacin opsi (b) + saring medan identitas), mesin salinan nbfacin, `KasusEDM` (urutan query mengikat), `Sumber` per-rule; sisa: resolver K-018 |
| **E02** | sebagian | IsUW +Marketing, IsClaim (`PropertyHasValue`, sikap khusus), IsTravel dua sumber, IsNotEDM bukan negasi — ketiga K-046 |
| **E04** | sebagian | + langkah 10–16 dari registry (`PredikatDari`); IsLife EDM membaca `BusinessOldId` |
| **E12** | sebagian | `premifire.go`: `CalculatePremiFire` + `SetLocalNonMbuProrate` FIRE — bentuk 1/2 FIRE EDM, HALF_UP, K-046 apa adanya; sisa: satu pintu empat bentuk, K-018, bentuk 3/4 |
| **E13** | sebagian | `porsipembayaran.go`: `CountPaymentEdm_Act` langkah 1/3/4/9/11 + blok 13–15, `CountPaymentEdmTSIObj_Act`; gerbang dari registry; sisa: rumus pembayaran |
| E14, E15 | blocked | rumus pembayaran E13 + pendekatan port DataTransform (§4 butir 3) |
| E16 | blocked | pendaftaran modul (`modul.go`/menu) — pemegang modul memilih **belum** |
| E19 | blocked | OQ-037/040/046 — ambang `"3000000000"` (§4 butir 4) |
| E22 | blocked | E14, E15, E19 |
| E17, E21 | dilewati | perintah work owner |

## 2. Keputusan agent — menunggu konfirmasi

| # | Keputusan | Bukti |
| --- | --- | --- |
| A01 | Fixture NB-15 dibaca DI TEMPATNYA (`modul/nbfacin/.../testdata/kasus/edm-fire-1.json`), tidak disalin — satu sumber yang dijaga penjaga de-identifikasi NB; tanpa impor Go lintas modul | `rekonsiliasi_test.go` |
| A02 | Seam 4 dua fungsi (`PrepareBeforeImage`, `IsiNilaiLama`) — dipertahankan sesuai perintah | laporan before-image §3 |
| A03 | E18 memodelkan cabang `Decision19`, bukan flow `InputEDMLife` | `grep -rl InputEDMLife "Endorsment Fac In"` → nol berkas |
| A04 | E20: generator (`docs/alat/bangkit_medis.py`) + penafsir subset ekspresi Pega — ambang tidak pernah disalin tangan | 728 ekspresi terurai (`TestTabelBangkitanUtuh`) |
| A05 | Predikat lini dan `IsEDM` masukan (seam 4, gerbang E04) sampai registry varian EDM ada | §4 butir 1 |
| A06 | Gerbang (E04), validasi tanggal (E05), dan `OpenCase` (E03) tiga fungsi terpisah; urutan milik pemanggil | `SetErrorBatalEndorsement_Act` & `CheckEDMPolisDate` dipanggil dari `Section/WorkPrimaryDetails.xml`; `SetValueToEDMWork` dari `Section/crmNewHarnessButtons.xml` |
| A07 | Kode transisi `5` = jalankan tanpa memeriksa baris berikutnya `[dugaan]` — nol `WhenTruePrms`, hanya bacaan ini yang membuat `param.JN=="All"` bermakna; `1` = lompat ke label `[terverifikasi]` (`WhenTruePrms` = `END`/`jmp`/`jmp2`/`TO`) | `medis.go`, `CheckEDMPolisDate` langkah 2 |
| A08 | Perbandingan teks-vs-angka yang hasilnya berbeda → galat (`ErrTafsirBerbeda`, `ErrTafsirKodeBerbeda`), mengikuti sikap registry NB butir 20 | `ekspresi.go`, `gerbangtolak.samaKode` |
| A09 | Porsi periode dihitung dalam HARI BULAT `[dugaan satuan]`: lokal langkah 15 bertipe `int` `[terverifikasi]`; rule lain membagi selisih DateTime dengan 365. Pembulatan HALF_UP = A37 (dikonfirmasi, bukan lagi keputusan agent) | `porsiperiode.go` |
| A10 | Konstanta kode status di `models/` (`StatusBusinessEDM`, `StatusEDMPolisBatal`, kini juga `FlagDihapus`, `FlagPolisBerjalan`, `NilaiShortPeriod`, `NilaiFixRate`, `NilaiBenar`, `CalculateMethodProRata`) dan `EdmType` bertipe `models.JenisEndorsemen` (K-029) | penjaga `claimlife` `TestKodeStatusLiteralHanyaDiModels` |
| A11 | Halaman primer `CalculatePremiFire` = satu coverage di `pyWorkPage.PropertyList(n).CoverageList` `[dugaan]` — sehingga TSI langkah 6 dipakai langkah 7/10; dipilih lewat indeks | pemanggil: aksi klik `Section/InputCoverageFire_IsUW.xml` (kelas `Data-Property`, 4×), `InputPerCoverageFire_IsUW.xml` (2×), `ViewCoverageFire.xml` (2×) |
| A12 | `.IsAdjustableFlag="true"` (`=` tunggal, 6.2.3) ditafsir perbandingan; `TSIAdjustment` 6.2.1 dibiarkan kosong bila faktornya kosong, galat hanya bila dipakai | `=` tunggal di prakondisi 62 berkas (`grep -rlE '<pyStepsPreCondParamsWhen>[^<]*[^=!<>]="'`), termasuk deret `X="a" ∨ X="b"` yang hanya bermakna sebagai perbandingan |
| A13 | Faktor kosong yang TIDAK dibungkus `@if` di rumus premi → `ErrNilaiPremiKosong`; perkalian/penjumlahan/pembagian eksak, ditolak bila melampaui 38 digit | `premifire.go`, `pembagian.go` |
| A14 | Predikat yang membandingkan MEDAN identitas (nomor polis, nama marketing, telepon/login operator) dengan literal tidak-kosong → panic tanpa nilai (`ISERRORSPREADING`, `ISTBONDING`); diport lewat konfigurasi (pola `DaftarPolis` E03) / model peran | `saring_registry.py --uji` 9/9; atas keluaran generator lama menyaring tepat 2 entri itu; `TestRegistryTanpaLiteralIdentitas` (diuji dengan literal sintetis lewat `-overlay` → merah) |
| A15 | `IsClaim` EDM diport tangan (`sikapKhusus`); arti `PropertyHasValue` atas halaman diserahkan pemanggil (`PemeriksaNilai`) | `When/IsClaim.xml` `pyConditionValue1` |
| A16 | `EDMDay` kosong sah (15.2 mengujinya `local.day==""`), teks bukan bulat → `ErrEDMDayBukanBulat` | `CountPaymentEdm_Act` 15.2 |
| A17 | Rekonsiliasi `TestPorsiPembayaranKasusNyata`: hasil query jenis bisnis polis lama tidak ada di fixture, diisi `OldData.QuotationData.BusinessType` `[dugaan]` — hanya IsMarineCargo yang membacanya, dan blok 13 tidak | `rekonsiliasi_test.go` |
| A18 | Mesin evaluator predikat = SALINAN `nbfacin/rules` (asal + sha256 di `predikat.go`), satu penyesuaian EDM (`sikapKhusus`). Dasar: jawaban work owner **"EDM DAN NB MENU YANG TERPISAH"** (register nbfacin butir 57) — ⚠️ **ditafsir** sesi nusantarare-c3 dan nusantarare-0f sebagai "modul terpisah penuh, mesin tidak ke inti"; itu tafsiran, bukan kutipan | `predikat.go` |

## 3. Temuan korpus — menyimpang dari tiket/spec

1. **Label `//` (`pyStepsBlockName`) = di-remark** (butir 43). Berdampak: 16 langkah skoring medis
   (SGOT/SGPT/Gamma GT Prodia, AFP & Urea N Pramita, Urea N Biotest — dua jenis kelamin); langkah 12–13
   penomoran `SaveEDMToJsonPolicy_Act`. Itulah sebab `GenerateEndorsementNo` dan `GenerateEDMNoLife`
   "nihil di korpus" — keduanya sudah tidak dipakai. Spec Life §2.2 perlu diralat.
2. **`pyStepsRepeatDefHasRepeat=REPEAT`** di 50 langkah blok bernilai Start=1 · Limit=1 · Iteration=1
   (56 definisi) — satu kali jalan.
3. **E04**: prakondisi yang dieksekusi tidak memuat satu pun nomor polis. 557 dari 560 kemunculan ada di
   cache editor `pyExpressionGadget/pyExpressionMapNew` (`langkah.py --lokasi-pola`). Spec §4.1 mencacah 502
   di berkas itu — cacah cache.
4. **E05**: perbandingan `@substring(…,0,8)` hanya di langkah "temporary utk monitoring"; keputusan
   memakai CARI20 dan `@CompareDates`. Kejanggalan yang sebenarnya diport: tanggal mulai dibaca dari
   tanggal kalender **GMT** — sehari lebih awal dari WIB untuk polis 17:00 GMT.
5. **E12**: rumus FIRE endorsement (`CalculatePremiFire` L2564/L2811: `…×FirstLossScale×IndemnityPercentage
   / 1e9`, skala 4, dua bagian dibulatkan terpisah) berbeda dari rumus FIRE NB di balik kontrak
   (`CountPremi_ACT`: `…×Indemnity×LossLimit / 1e9`, skala 20).
6. **E20**: `SetParamLab_Act` memuat 56 nilai `Negatif/0` TANPA kutip — Pega membacanya sebagai
   ekspresi; di sini galat per properti. `IsLife` `pyConditionViewer` memuat cache `BusinessCode =
   "10164"` — yang dieksekusi `BusinessOldId = "L1"…"L16"`.
7. **E18**: `Decision39` titik temu jalur Life dan non-Life (non-Life tiba sesudah tangga).
8. **E13 — dua langkah yang mengubah angka** (temuan code review, diverifikasi ke korpus):
   `CountPaymentEdm_Act` langkah 4 menormalkan `OfferFacIn.PolicyData.Start/EndDateTime` ke
   `FormatDateTime(…,"yyyyMMdd","Asia/Jakarta")+"T050000.000 GMT"` (tanggal Jakarta pukul 05:00 GMT);
   langkah 9 menyetel `ProtectSpreading.CARI2 = "FIX RATE"` untuk IsMarineCargo ∨ IsEdmAdjShareCedant —
   Adj Share Cedant karenanya SELALU porsi 1/0 (14.3). `CountPremiEDM_DT` hanya membaca CARI2.
   `EdmDate` tidak dinormalisasi, tetapi fixture menyimpannya `T050000.000 GMT` — jam yang sama dengan
   hasil langkah 4 — sehingga selisih blok 14/TSIObj atas tanggal nyata **bulat**
   (`TestPorsiPembayaranTanggalNyataBulat`, predikat sintetis; bukan rekonsiliasi). ⛔ Ralat: draf
   laporan ini sempat menulis `EdmDate` fixture `T170000.000 GMT` → n,5 hari; caranya keliru — jam
   `PolicyData` (`T170000`) dikira jam `EdmDate` tanpa dibaca. Langkah 15 `SetValueToEDMWork` tetap
   ditolak: ia membaca `OldData.PolicyData` TANPA normalisasi.
9. **E12**: rumus FIRE EDM memakai porsi sebagai **persen** (`@Math.divide(hari,365,6)*100`, pembagi
   1e9) — tiket menulis "pecahan, dipakai langsung". Polis berjalan: `+12 jam` hanya pada argumen kedua
   `@DateTimeDifference` → dengan jam sama di kedua tanggal, selisih SELALU n,5 hari → port menolak.
10. **E01**: `IsLife` EDM membaca `pyWorkPage.Quotation.BusinessOldId` (L1…L16), bukan jenis bisnis
    "Life" yang ditulis `SetErrorBatalEndorsement_Act` langkah 8. `IsEDM` membaca
    `OfferFacIn.QuotationData`, `IsNotEDM` membaca `Quotation` — keduanya dapat benar/salah bersamaan
    (diuji dua arah).
11. **Instrumen** (atas ralat sesi nbfacin): pasangan `pyParamArray` di langkah TANPA metode mungkin sisa
    tak dieksekusi. Diaudit atas 17 activity yang diport modul ini: **nol** penugasan di langkah tanpa
    metode; instrumen diuji pada kasus positif yang diketahui (`CountGrossPremiEDM_Act` 1.1.2.1.4 `.TSI`
    → tertangkap, + 3 lainnya). `langkah.py` kini menandainya `SET? (langkah tanpa metode)`.

## 4. Menunggu work owner

1. ✅ **Registry predikat varian EDM — opsi (b) DIJALANKAN.** Generator nbfacin (`rules/bangkit -folder
   "Endorsment Fac In\When" -paket predikat`) → 202 berkas, 201 predikat → `saring_registry.py` →
   `predikat/registry_gen.go`; bangkit ulang + saring = identik (`diff -q`). Larangan "jangan tiru mesin NB"
   dicabut sesi nusantarare-0f untuk mesin predikat → A18. **Konfirmasi yang diminta:** tafsiran A18 atas
   jawaban "EDM DAN NB MENU YANG TERPISAH". Generator kini versi terbaru c3 (menolak medan identitas
   sendiri); saringan EDM tinggal lapis kedua (0 entri tersaring). Sesi c3 melaporkan kelima nilai
   sensitif (2 nomor polis, 3 nama) tidak pernah masuk riwayat git.
2. ✅ **Rumus premi FIRE endorsement — DIPUTUSKAN: port di modul ini** (perintah 1); dikerjakan
   (`premifire.go`). Yang tersisa untuk work owner: **satuan `@DateTimeDifference(…,D)` atas pecahan
   hari** — menghalangi langkah 15 `SetValueToEDMWork` (kasus nyata 150,5 hari) dan
   `SetLocalNonMbuProrate` FIRE (`+12 jam`, temuan 3.9). Blok 14/TSIObj `CountPaymentEdm_Act` lolos
   berkat normalisasi langkah 4 (temuan 3.8).
3. **Pendekatan port E13–E15** (171 langkah / 406 penugasan + 137 aksi DataTransform): rekomendasi agent —
   pola E20 (pembangkit + penafsir atas clipboard berjalur, sejalan `kontrak.KasusFacIn`), dengan konversi
   ke `uang.Money` di batas keluar.
4. **OQ-037/040/046** — mata uang ambang `"3000000000"` (`ReCountPremiLifeEDM` 4.1.4/4.1.5) membuka E19.
5. ✅ **A37 DIKONFIRMASI** (HALF_UP untuk `@Math.divide`, disampaikan sesi nusantarare-0f) — porsi periode
   E06 kini dihitung HALF_UP menjauhi nol, eksak lewat bilangan bulat (`TestPorsiPeriodeSetengahKeAtas`,
   `TestPorsiPeriodeNegatifSetengahMenjauhiNol`; mutasi tanpa-pembulatan dan tanpa-tanda tertangkap).
   `[terverifikasi]` A37 berlaku langsung: `SetValueToEDMWork.xml` L5341/L5368 `@Math.divide(…, 20)`.

   > ⛔ **RALAT 01-10-2026** (atas verifikasi sesi nusantarare-0f, diperiksa ulang ke berkas). Klaim
   > sebelumnya — "fixture NB-15 menyimpan 151/150 = bukti sejalan" dan "porsi periode kini terisi untuk
   > data nyata" — **keliru**. Yang keliru caranya: test 151/150 memakai **tanggal buatan**, yang cocok
   > hanya bentuk angka; dan 151/150 pada 20 desimal sama untuk HALF_UP dan HALF_EVEN (hanya pemotongan
   > tersingkir), jadi bukan bukti HALF_UP. Fakta kasus nyata `edm-fire-1.json` (dihitung ulang dari
   > `/OldData/PolicyData/StartDateTime`, `/EndDateTime`, `/QuotationData/EdmDate`): EdmToStart = 150,5
   > hari, EdmToEnd = −0,5 hari, TotalPeriod = 150 → port ini **MENOLAK** kasus nyata satu-satunya
   > (`TestPorsiPeriodeKasusNyataEDMFire`). Nilai tersimpan `ProrateEDMEnd = 1` juga bukan keluaran
   > langkah 15 (−0,5/150) — ditulis/ditimpa rule lain (§7).

   Yang TETAP terbuka: satuan selisih DateTime (A09). `[dugaan]` 150,5 hari → 151 mengisyaratkan Pega
   membulatkan selisih DateTime atau menghitung per tanggal kalender GMT — bahan pertanyaan terbuka, bukan
   dasar kode.
   `@divide` (fungsi lain, skoring medis) TIDAK ikut diubah: A37 menyebut `@Math.divide`. ⛔ A39 (rekonsiliasi
   kasus EDM dengan rumus NB) DITAHAN work owner — tidak dipakai di modul ini.
6. **Ralat tiket/spec**: E10 (`K046_Life_DuaPenandaOldData`, `LayerList`), E18 (`InputEDMLife`), E05
   (K-046 substring), E12 (premis bentuk 1 tunggal), spec Life §2.2, spec alur masuk §4.1.
7. **Semantik Pega yang belum terverifikasi** (di kode: galat eksplisit atau `[dugaan]` berlabel):
   `@toDate` atas `""` dan `dd/MM/yyyy`; `@CompareDates(a,b)` = a sesudah b; konteks titik di dalam
   `UPDATE_PAGE` (`DataToEDM`); `Page-Remove` di dalam loop; `Property-Set` pada nomor baris yang belum ada
   (Travel, lapis B); penjumlahan properti desimal kosong; skala bawaan `@divide` dua argumen; `@if`
   malas; nilai tak berkutip `Negatif/0`.
8. **Underwriting**: arti teks skor medis ("Invalid Age", "150%", "175", tiga ejaan "Postpone…");
   penggabungan menjadi satu keputusan akhir tidak ada di activity mana pun; jiwa tanpa tangga akseptasi.
9. **RBAC** (dari sesi nbfacin): `CountGrossPremiEDM_Act` memuat 8 prakondisi
   `OperatorID.pyUserIdentifier == <ID operator literal>` — tidak diport sebagai literal; bahan OQ RBAC
   (OQ-021 dkk.). Belum dipakai modul ini.
10. **Kandidat inti** (di luar jatah): pembagian HALF_UP `@Math.divide` kini ada di beberapa modul
    (nbfacin, treatycontractout, endorsmentfacin `pembagian.go`) — layak satu fungsi di `inti/backend/utils`.
11. **Pemilik `claimlife`**: medan `QuotationData.Type` EDM bertipe `models.JenisPenyesuaian` sehingga
   tidak tertagih `satutype_test.go` — boleh didaftarkan eksplisit.

## 5. Uji

| | Hasil |
| --- | --- |
| `go vet ./...` | bersih |
| `go test ./...` (APP_RNM) | **46 paket ok, nol gagal** pada putaran pertama (termasuk `inti/backend/penjaga` — impor lintas modul). Putaran penutup: **45 ok, 1 gagal build** — `modul/nbfacin/backend/services/premium`: berkas `total_test.go` yang muncul 20:15 merujuk `TotalFireLokasi` yang belum ada (pekerjaan sesi nbfacin yang sedang berjalan; di luar jatah, tidak disentuh). Seluruh paket endorsmentfacin, penjaga, dan claimlife ok |
| `npx vitest run` | 8 gagal — seluruhnya `modul/treatycontractout/frontend` (sudah merah sebelum putaran ini; modul ini tanpa frontend) |
| Alat audit `docs/alat/langkah.py --uji` | 6/6 butir yang jawabannya diketahui |
| Alat `docs/alat/saring_registry.py --uji` | 9/9 |
| Uji mutasi (instrumen test) | premi FIRE 3/3 tertangkap (HALF_DOWN, PremiRp tak tertukar, FLS dibungkus); porsi pembayaran 3/4 tertangkap (langkah 9, 4, 1), 1 mutan ekuivalen (15.2 `kosong ∨ d==0` vs `d==0`: `hariEDM` memberi 0 untuk kosong) |
| Code review dua sumbu | tiga kelompok (E06/E18/E20; E03/E04/E05; **E12/E13**) — temuannya diperbaiki tanpa bertanya. E12/E13: langkah 1/4/9 diport, `EDMDay` tak lagi ditelan, `EdmType` bertipe, literal ke `models`, pembagian ke `pembagian.go` (+ jaga 38 digit, label `[dugaan]` nilai negatif), rekonsiliasi lewat registry, header usang diralat. Tidak diperbaiki: satu pintu empat bentuk & dua keluarga rasio bertipe (E12, desain), pembagian ke `inti` (di luar jatah), test membaca fixture nbfacin (A01) |

## 6. Berkas (belum di-commit)

Baru: `backend/services/{bukakasus,gerbangtolak,validasitanggal,jalurlife,porsipembayaran,premifire,pembagian,predikatkasus}.go` (+ test),
paket `backend/services/predikat/` (`predikat.go` + `logika/` salinan, `kasus.go`, `registry_gen.go`
bangkitan-tersaring, test), `docs/alat/saring_registry.py`,
`backend/services/rekonsiliasi_test.go`, paket `backend/services/medis/` (`ekspresi.go`, `medis.go`,
`aturan_gen.go` bangkitan, test), `docs/alat/bangkit_medis.py`, laporan ini. Diubah terhadap `9ca6f06`
(`git status --short APP_RNM/modul/endorsmentfacin`): `backend/models/penawaran.go`,
`docs/alat/langkah.py`, `docs/LAPORAN-IMPLEMENTASI-BEFORE-IMAGE.md` (penunjuk), baris `Status:` tiket
E01–E22 dan `00-INDEKS-EDM.md` (keduanya bagian rename dokumen yang ter-stage milik pemegang modul).
⚠️ `MODUL.md` dan `docs/PETA-ASAL.md` berubah/baru oleh pihak lain sebelum sesi ini — tidak disentuh.
