# DRAF keputusan Fac Out — K-057 … K-062

> # ✅ SUDAH DISETUJUI — DIPINDAHKAN KE REGISTER
>
> **Work owner menyetujui keenamnya pada 21 September 2026.** Teks yang **berlaku** kini ada di
> **`..\00-KEPUTUSAN-WORK-OWNER.md`** sebagai **K-057 … K-062**. **Nomor berikutnya = K-063.**
>
> ⚠️ **Dokumen ini dipertahankan sebagai arsip bahan**, bukan dihapus (`PANDUAN-KERJA` §7). Bila isi
> di sini berbeda dari register, **register yang berlaku**.
>
> **Perubahan work owner terhadap draf — catat, jangan dibaca sebagai salah-salin:**
>
> | Draf | Keputusan yang diambil |
> | --- | --- |
> | K-058 rekomendasi **opsi D** | dipilih **opsi A** + pengecualian: rule yang sudah diekspor ulang ke `DDL\` **menang**. Pertanyaan metode deteksi drift **tetap TERBUKA** |
> | K-059 usul "39 bernama + 16 pembawa" | disetujui, **dengan alasan tegas**: `10015` **adalah** kode Fac Out, jadi keduanya **diperlakukan sama** |
> | K-060 dua sisi (Seam 3 vs seam sendiri) | ditetapkan: tetap `CountRateRetroCov`, **salinan `DDL\` yang berlaku**. ⛔ Konsekuensinya keanehan `RIComIN` **gugur** — lihat K-060 di register |
> | K-062 usul "bangun ulang di atas riwayat sendiri" | ditetapkan: **dihapus, tidak dipakai**; kolom pelaksana **kosong**. ⛔ **Pengecualian pertama yang disengaja terhadap K-006** |
>
> **Tanggal draf:** 21 September 2026 · **Disetujui:** 21 September 2026

---

## Daftar isi

| Draf | Pokok | Sifat |
| --- | --- | --- |
| **K-057** | Pemicu Fac Out = flag `.IsFacRetro`, bukan `When IsFacout` | koreksi temuan |
| **K-058** | Sumber kebenaran spec saat versi sama tapi isi beda + metode deteksi drift | pilihan A/B/C/D |
| **K-059** | Populasi kerja Fac Out ≈ 55 activity per folder | lingkup |
| **K-060** | Premi retro: bentuk Seam 3 atau seam tersendiri | seam |
| **K-061** | Cetak RI Slip + email: lingkup sekarang atau ditunda | lingkup |
| **K-062** | Pengganti `GetOPFacOut_Sql` | rancangan |

---

## DRAF K-057 · Pemicu Fac Out adalah flag `.IsFacRetro`, bukan `When IsFacout`

**Pertanyaan:** apa yang sesungguhnya memicu kasus masuk ke menu Fac Out.

**Usulan keputusan:** pemicu Fac Out ditetapkan sebagai **flag `.IsFacRetro`**, yang di-*reset* ke `0`
lalu dinyalakan ke `1` **per lini bisnis** oleh penanda spreading Fac Out saat UW Accept. Rumusan lama
**"UW Accept (`ProposalAcceptStatus = 4` → When `IsFacout`)"** di `09-facout\01` §0, §3 dan §4.2
**dicabut**.

**Dasar `[terverifikasi]`** — rantai lima langkah, tag pembawa disebut:

| # | Berkas | Tag | Isi |
| ---: | --- | --- | --- |
| 1 | `Activity\SetValidateDateUW_PostAct` | `pyStepPageReference` = `RH_1.pySteps(3)` | `Call SetDataFacOut_Act`, **tanpa precondition** |
| 2 | `Activity\SetDataFacOut_Act` | `<Property>` ×4 pada `RH_1.pySteps(2)` · `<PropertiesName>` pada `RH_1.pySteps(3)` | `RH_1.pySteps(2)` `Property-Remove` atas `OfferFacIn.FacRetro.`{`LocationList`,`PersonList`,`CargoList`,`VehicleList`}; `RH_1.pySteps(3)` `Property-Set` `.IsFacRetro = 0`; `RH_1.pySteps(4..7)` panggil per COB. ⚠️ `OfferFacIn.FacRetroList` muncul di `RH_1.pySteps(8).pySteps(1)` sebagai `<pyStepsObjectName>` pada langkah **bermetode kosong** — **bukan** sasaran `Property-Remove` |
| 3 | `Activity\SetDataFacOutFire_Act` | `pyStepsPreCondParamsWhen` · `PropertiesName` | `.IsFacRetro = 1`, bergerbang `.TreatyType=="10015"` / `Local.isFacOut==1` / `Local.isFacOutLoc==1` |
| 4 | `Flow\InputInwardFacultativeOffer` | `pyTaskWhen` = `IsFacRetro`, `pyTaskStatusOrWhen` = `WHEN`, `pyLikelihood` = 100 | transisi `Transition92` → shape `pyMOName` = `[Fac Out]` |
| 5 | shape `[Fac Out]` | **`pyImplementation`** = `OfferFacRetro` | **bukan** `pySubFlowName` |

**Dua lapis kekeliruan rumusan lama:**

1. `[terverifikasi]` `IsFacout` muncul **0 kali** di `Flow\InputInwardFacultativeOffer.xml`, peka huruf
   maupun tidak. `IsFacRetro` **17**, `IsInputFacRetro` **9**.
2. `[terverifikasi]` `ProposalAcceptStatus = 4` berarti **Banding**, bukan Accept — **K-029** dan
   `GLOSARIUM.md`.

**Konsekuensi bila disetujui:**

- `09-facout\01` §0.1 (sudah ditulis) menjadi rujukan pemicu; §0/§3/§4.2 lama tinggal sebagai catatan
  dibatalkan.
- `IsFacout` **keluar** dari registry predikat Fac Out; porting-nya tetap di tiket New Business
  `05-tickets\09-predikat-sikap-khusus.md` di bawah K-019.
- Tiket **F01** dan **F02** memakai pemicu ini.

**`[pertanyaan terbuka]` yang tidak ditutup keputusan ini:** indeks rujukan flow menuliskan resolusi
`IsFacRetro` ke app **SFAGIS** / workType **ClaimLife**, dan `OfferFacRetro` ke app **GISFW** /
workType **LIFE**, sementara `Embed-Reference-Rule` menyebut kelas `ASM-FW-GISFW-Work` dan definisi
ekspornya berkelas `ASM-FW-GISFW-Data-OfferFacIn`. **Butuh UI Pega. Jangan disimpulkan.**

---

## DRAF K-058 · Sumber kebenaran saat `pyRuleSetVersion` sama tetapi isi berbeda

**Pertanyaan:** bila satu identitas rule hadir di lebih dari satu folder dengan **nomor versi sama**
tetapi **isi berbeda**, salinan mana yang menjadi sumber kebenaran spec — dan apakah metode deteksi
drift diganti.

**Dasar `[terverifikasi]`** — sensus drift 21 September 2026, kontrak 23 tag:

| Pasangan | Dibanding | Versi sama + isi sama | **Versi sama + isi BEDA** | Versi beda |
| --- | ---: | ---: | ---: | ---: |
| NB ↔ RNW | 1.907 | 1.907 | **0** | 0 |
| NB ↔ EDM | 1.707 | 1.353 | **267** | 87 |
| RNW ↔ EDM | 1.674 | 1.330 | **259** | 85 |
| **Total** | 5.288 | — | **526** | 172 |

Pemeriksaan silang: `1.353 + 267 + 87 = 1.707` ✅ dan `267 + 87 = 354` ✅ — cocok angka resmi K-043.

⛔ **Dari 354 perbedaan NB↔EDM, 267 (75,4 %) ber-`pyRuleSetVersion` IDENTIK.** Metode deteksi drift
yang membandingkan nomor versi **melewatkan tiga dari empat perbedaan nyata**. Angka "95 rule berbeda
versi" di `PANDUAN-KERJA` §5 bukan sekadar batas bawah — ia melewatkan mayoritas.

⚠️ `[pertanyaan terbuka]` **Angka 95 tidak tereproduksi.** Pengukuran ini memberi 87 (NB↔EDM),
85 (RNW↔EDM), 0 (NB↔RNW). Metode yang menghasilkan 95 tidak terdokumentasi.

**Contoh paling tajam:** `Flow\OfferFacRetro` — **empat salinan, `pyRuleSetVersion` `01-01-95`
seluruhnya**, isinya 2 tingkat (NB/RNW) vs 4 tingkat (Endorsment/DDL).

### Opsi

| Opsi | Isi | Konsekuensi |
| --- | --- | --- |
| **A — per siklus** | Spec NB pakai salinan NB, RNW pakai RNW, EDM pakai EDM | Paling setia pada perilaku terekam. ⚠️ Membekukan versi ekspor yang kebetulan tertangkap; bila di produksi ketiganya sebenarnya **satu rule yang sama**, opsi ini melembagakan perbedaan yang tidak nyata |
| **B — paling baru menang** | Pakai `pxCommitDateTime` terbesar | ⛔ **Tidak dapat dipakai sendirian.** `pxCommitDateTime` adalah stempel **ekspor**, bukan stempel perubahan rule, dan `pyRuleSetVersion` sudah terbukti tidak andal. Berisiko memilih salinan yang salah secara diam-diam |
| **C — `DDL\` menang** | Ekspor ulang dari work owner jadi rujukan | Hanya berlaku bagi rule yang ada di `DDL\` — segelintir dari 526. Tidak menyelesaikan masalah; **berguna sebagai pemutus untuk kasus tertentu saja** |
| **D — per siklus + daftar tinjau** ✅ | Opsi A sebagai dasar, **plus** 526 pasangan dijadikan daftar tinjau eksplisit; yang menyentuh **uang, gerbang, atau jalur produksi** diverifikasi ke UI Pega sebelum dipakai | Biaya tinjauan terbatas dan terukur; tidak ada perbedaan yang hilang diam-diam. Lebih mahal daripada A di muka, jauh lebih murah daripada selisih tak terjelaskan saat paralel run |

**Rekomendasi: D.** Alasannya bukan kehati-hatian abstrak — ADR-0001 menuntut **nol selisih sampai
digit terakhir**, dan 267 perbedaan yang tak terlihat metode lama adalah tepat jenis hal yang muncul
sebagai selisih tak terjelaskan. **Keputusan tetap milik work owner.**

**Usulan menyertai:** metode deteksi drift di `_ARSIP-lintas-siklus\_BACA-INI.md` diganti dari
**banding `pyRuleSetVersion`** menjadi **banding hash 23 tag**. Ukurannya sederhana: metode lama
melihat 87 dari 354; metode baru melihat 354.

---

## DRAF K-059 · Populasi kerja Fac Out ≈ 55 activity per folder

**Pertanyaan:** berapa luas permukaan Fac Out yang harus diport — menentukan lingkup **F11**.

**Dasar `[terverifikasi]`:**

| Ukuran | NB | RNW | EDM |
| --- | ---: | ---: | ---: |
| Activity **bernama** `FacOut`/`FacRetro` (Ordinal, peka huruf) | 41 | 40 | 39 |
| Hadir di ketiga folder | — | **39** | — |
| Activity memuat penanda `10015` **tanpa** nama Fac Out | 16 | 15 | 16 |
| **Permukaan fungsional** | **≈57** | **≈55** | **≈55** |

⛔ **Angka lama 18/18 dan 10/18 DIBATALKAN** — populasinya tidak pernah terdefinisi.

Perbandingan atas 39 yang hadir di ketiganya: **NB vs RNW 39/39 identik**; **NB vs EDM 30/39 identik,
9 berbeda** (`CopyAllObjFacOutFireAneka_ACT` · `CopyFacRetroAnekaGolf_ACT` · `InsertFacoutProd` ·
`InsertFacoutProduction` · `InsertFacoutProductionEDM` · `InsertFacoutProductionEDMCurr` ·
`OfferFacOut_PreAct` · `SetDataFacOutAnekaGolf_Act` · `SetDataFacOutFire_Act`).

**Usulan keputusan:** populasi kerja Fac Out ditetapkan **≈55 activity per folder**, bukan 39. Yang 16
tanpa nama Fac Out — termasuk `CopyToAllSpreading_ACT`, `CopyToAllLocSpreading_ACT`,
`cekSpreadingFactIn`, `CountPremiumNet`, `SaveJsonPolicyFacIn_Act`, `SetReinsurerEndorsement_Act` —
**wajib ikut ditinjau**, karena membawa penanda `10015`.

**Alasan:** `PANDUAN-KERJA` §3 — nama bukan bukti. Populasi bernama adalah **batas kemudahan, bukan
batas fungsi**. Memakai 39 akan melewatkan mesin spreading yang justru memuat penanda umum
`10015|SPL|10007`.

---

## DRAF K-060 · Premi retro — bentuk Seam 3 atau seam Fac Out tersendiri

**Pertanyaan:** `Activity\CountRateRetroCov` menghitung uang. Apakah ia **bentuk terparameterisasi
Seam 3** (`services/premium.Calculate`) atau **seam Fac Out tersendiri**.

**Dasar `[terverifikasi]`** — rumus intinya:

```
.PremiumRetro     = @Math.divide((ShareOffered * .Rate * Local.prorate), 100000, 20)
.PremiumRetro     = .PremiumRetro - @Math.divide((.PremiumRetro * .DiscountPercentage), 100, 20)
.RICommPercentage = @Math.divide(Local.RIComIN, Local.LengCov, 20)
.RIComm           = @Math.divide((.PremiumRetro * .RICommPercentage), 100, 20)
```

Pembagi **100.000** = `1.000 × 100`, konsisten aturan **pembagi komposit K-018**.

| Sisi | Argumen |
| --- | --- |
| **Bentuk Seam 3** | Bentuk dasarnya sama dengan keluarga `nilai × rate × prorata ÷ pembagi-komposit`; pembaginya patuh K-018; preseden sudah ada — premi Life juga bentuk 1 terparameterisasi (K-052, tiket `edm\E19`). **Tidak menambah seam** (tetap 6). Satu tempat perubahan bila aturan pembulatan berubah |
| **Seam tersendiri** | Masukannya **bukan** TSI melainkan `ShareOffered` hasil penskalaan mata uang; ada **pengurangan diskon** dan **komisi RI** yang tidak ada di keluarga Fac In; ada **rumus koreksi kedua** yang melewati penskalaan; keluarannya properti retro, bukan premi polis. Menekannya ke Seam 3 berisiko membuat pintu Seam 3 bercabang terlalu lebar |

**Rekomendasi: bentuk terparameterisasi Seam 3**, dengan syarat penskalaan mata uang dan pengurangan
diskon dimodelkan sebagai **pra-pemrosesan di sisi pemanggil**, bukan cabang di dalam rumus. Ini
menjaga jumlah seam tetap **6** dan konsisten dengan cara premi Life ditangani. **Keputusan tetap
milik work owner** — ia menambah atau tidak menambah satu seam.

**`[pertanyaan terbuka]` yang harus ikut dijawab:**

1. ~~`Local.RIComIN` **dibaca tetapi tidak pernah disetel di seluruh korpus**. `.RICommPercentage`
   dibagi dari variabel yang tak pernah diisi. Diport apa adanya; **kandidat perbaikan** — apakah
   dibiarkan atau diperbaiki seperti A.5 (K-046)?~~
   ⛔ **PERTANYAAN INI GUGUR (K-060).** Keanehan itu **tidak ada** pada `DDL\CountRateRetroCov.xml`
   yang ditetapkan berlaku: di sana `Local.RIComIN` **diakumulasi** di `pySteps(8)`
   (`Local.RIComIN := Local.RIComIN + .RIComm`), sejajar `Local.RateIN`. Ia **akibat salinan korpus
   yang basi**. Penggantinya: `..\00-KEPUTUSAN-WORK-OWNER.md` **K-060** dan
   `01-pohon-langkah-countrateretrocov.md` §6.4.
2. Tiga penugasan pertama `Local.prorate` di langkah 1 saling menimpa, yang terakhir `:= 1`. Bila
   pasangan dalam satu `Property-Set` dieksekusi menurut `REPEATINGINDEX`, dua penugasan pertama
   **mati**. Apakah urutan pasangan = urutan eksekusi? **[di luar korpus]** — butuh UI Pega.

---

## DRAF K-061 · Cetak RI Slip dan email — lingkup sekarang atau ditunda

**Pertanyaan:** apakah **F13** (cetak RI Slip + kirim email) dikerjakan sekarang, atau ditandai
menunggu seperti **R07** (`05-tickets\rnw\R07-konversi-produksi-renewal.md`).

**Dasar `[terverifikasi]`:** `Flow\OfferFacRetro` memuat shape `PRINT R/I SLIP`,
`Send Email Print RI Slip`, `PrintRISlip_FlowAction`, dan `IsRISlip`. Empat activity PDF ada di korpus
(`GeneratePDFFacOutMemoPlacing` · `GeneratePDFFacOutOfferStatus` · `GeneratePDFFacOutViewOffer` ·
`DeletePDFFacOutMemoPlacing`).

⛔ **Penghambatnya:** `CLAUDE.md` §4.4 — daftar endpoint sesungguhnya ada di tabel Oracle
`M_LINK_SERVICE` (kunci `KATEGORI_1` + `KATEGORI_2`), yang **isinya tidak ada di korpus**. Tanpa itu,
tujuan kirim tidak dapat ditentukan, dan endpoint **tidak boleh** ditulis literal.

| Opsi | Konsekuensi |
| --- | --- |
| **Masuk lingkup sekarang** | Tiket ditulis penuh, tetapi eksekusinya tetap berhenti di titik yang sama. Risiko: tiket tampak siap padahal tidak |
| **Ditandai menunggu (pola R07)** ✅ | Lingkupnya tetap terlihat dan tidak hilang; statusnya jujur. Preseden sudah ada — R07 ditulis penuh, hanya eksekusinya menunggu |

**Rekomendasi: pola R07** — tiket **F13 ditulis penuh** dengan blocker eksternal dinyatakan tegas,
**tidak dihapus**. Menghapusnya menyembunyikan pekerjaan yang pasti ada. **Keputusan milik work
owner.**

---

## DRAF K-062 · Pengganti `GetOPFacOut_Sql`

**Pertanyaan:** `RDBList\GetOPFacOut_Sql` tidak dapat diport. Apa penggantinya.

**Dasar `[terverifikasi]`** — tag `<pyBrowseSQL>`, `<pxObjClass>` = `Rule-Connect-SQL`:

```sql
select PYPERFORMER as "pyOwnerUserID"
  from datapega.pc_history_asm_fw_gisfw_work
 where (PYMESSAGEKEY like '%OfferFacOut%' or PYMESSAGEKEY like '%OfferRetro%')
   and PYHISTORYTYPE='F'
   and PXHISTORYFORREFERENCE={pyWorkPage.pzInsKey}
 order by PXCOMMITDATETIME
```

⛔ `datapega.pc_history_*` adalah **tabel internal Pega**. `CLAUDE.md` §4.3: `DATAPEGA.PC_*`
**hilang bersama Pega — jangan dimigrasikan apa adanya.**

**Yang dibacanya, dinyatakan netral:** identitas pelaksana (`PYPERFORMER`) dari **riwayat penugasan**
kasus, disaring pada langkah yang namanya memuat `OfferFacOut`/`OfferRetro`, diurutkan menurut waktu
commit. ⚠️ **Nilainya adalah identitas orang** — K-025: tidak pernah disalin ke artefak; yang dicatat
hanya mekanismenya.

**Usulan arah:** sistem baru **menyimpan riwayat penugasannya sendiri**, dan klep ini dibangun ulang
di atas riwayat itu — bukan di atas tabel Pega. `HISTORYAKSEPTASIPEGA` sudah dipakai untuk mengambil
keputusan (`CLAUDE.md` §4.3) dan **bentuknya dipertahankan**, sehingga ia kandidat pertama.

**`[pertanyaan terbuka]` untuk work owner + DBA:**

1. Apakah `HISTORYAKSEPTASIPEGA` memuat langkah Fac Out (`OfferFacOut` / `OfferRetro`), atau riwayat
   itu **hanya** ada di `DATAPEGA.PC_*`? Bila hanya di sana, riwayat Fac Out historis **hilang saat
   migrasi** dan itu keputusan bisnis, bukan teknis.
2. Untuk apa nilai `PYPERFORMER` dipakai di layar Fac Out — menampilkan pelaksana, atau
   menggerbangi sesuatu? Menentukan apakah penggantinya wajib atau boleh kosong.
3. Apakah riwayat Fac Out lama perlu **dimigrasikan** ke tabel baru, atau cukup dimulai kosong?

⚠️ Sampai ketiganya terjawab, **F12** memuat klep ini sebagai **blocker eksternal**, dan jalur yang
memanggilnya **ditangguhkan, bukan dihapus** (pola K-006).

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
*Seluruh isi dokumen ini DRAF — menunggu persetujuan work owner.*
