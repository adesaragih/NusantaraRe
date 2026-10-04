# Checklist Kesiapan — Migrasi Facultative Inward

> ## 📍 Status siklus — 17 September 2026 (K-030)
>
> | Siklus | Keadaan |
> | --- | --- |
> | **New Business** (`NB FacIn\`) | ✅ **SELESAI sampai tahap tiket** — discovery · grilling · ADR · spec · **16 tiket** di `05-tickets\`. Final; tidak diulang |
> | **Renewal** (`RNW Fac In\`) | ✅ **SELESAI sampai tahap tiket** — discovery · keputusan · spec · **8 tiket** di `05-tickets\rnw\`. Final |
> | **Endorsement** (`Endorsment Fac In\`) | 🔵 **FOKUS AKTIF sejak 18 September 2026 (K-042)** — mulai dari discovery. Tiga pertanyaan matang di `04-kuesioner\_DITUNDA-fase-endorsement.md` kini **dibuka** |
>
> **Batas baca kini termasuk `RNW Fac In\`.** Batas tulis tidak berubah: hanya `OUTPUT\`.
>
> ⚠️ **Isi checklist di bawah ini disusun untuk siklus NB.** Sebagian besar butirnya (DDL, prosedur,
> tabel limit, enumerasi) **berlaku lintas siklus** dan tidak perlu diulang untuk RNW. Butir yang
> khas RNW ditambahkan setelah discovery RNW berjalan.

> Daftar semua yang ditunggu sebelum modul di luar 5-modul-terverifikasi dapat dibangun.
> Sumber: `04-spec\03-spec-modul-terverifikasi.md` (Out of Scope 13 butir) + `_EKSTRAKSI-PEGA-SELAGI-HIDUP.md`.
> Semua daftar sudah terverifikasi dari korpus. Tandai `[x]` saat sebuah item kembali.
>
> Cara pakai: kirim bagian per pihak ke pemiliknya, centang saat jawaban masuk, lalu catat sebagai
> keputusan (K-nnn) bila mengubah rancangan.

Legenda status: ⬜ belum dikirim · 📤 sudah dikirim, menunggu jawaban · ✅ jawaban masuk

---

## A. DBA — buka database, jalankan query/ekstraksi

### A.1 · Isi kode 35 prosedur/fungsi POOLDATA  (butir 3 · membuka `services/production`, `repository`)
Perintah: `ALL_SOURCE` atau `DBMS_METADATA.GET_DDL`. **Definisi kode saja, tanpa data.**
Tersimpan di `D:\migrasi\RNM\DDL\*.txt`.

Status: ✅ MASUK (32 dari 35; 2 dikonfirmasi tidak ada di DB; 1 perlu konfirmasi)

Paling kritis (jalur simpan utama):
- [x] `PEGA_M_JSON_OFFER`
- [x] `INSERTJSONPOLIS`
- [x] `INSERTUPDATECEDINGPRODUCTION`

Produksi Fac In & Life:
- [x] `FACINOFFER`
- [ ] `FACINOFFERLIFE` — ❎ N/A: dikonfirmasi TIDAK ADA di DB (2026-09-16)
- [x] `FACINLIFE`
- [x] `FACINSPREADLIFE`
- [x] `FACINPERFORMANCE`
- [x] `FACINFORBACKUP` — MASUK (ditambah 2026-09-16)
- [x] `ERRORFACINPROD`
- [x] `TREATYPRODUCTION_BACKUP`
- [x] `PEGA_TREATY_IN`
- [x] `PEGA_JSON_POLIS_TREATYIN`

Riwayat & monitoring:
- [x] `HISTORYAKSEPTASIPRODUCTION`
- [x] `INSERTJSONPOLISMONITORING`
- [x] `MONITORING_PROD_LOG`
- [x] `PEGA_DELETE_ERROR_KONVERSI`

Master data:
- [x] `RDBINSERTCLIENT`
- [x] `RDBMASTERACCUMULATEDTYPE`
- [x] `RDBMASTERACCUMULATION`
- [x] `RDBMASTERBRANCH`
- [x] `RDBMASTERCITY`
- [x] `RDBMASTERDISTRICT`
- [x] `RDBMASTERNATION`
- [x] `RDBMASTERPROVINCE`
- [x] `RDBMASTERRW`
- [x] `RDBMASTERSHIP`
- [x] `PEGA_M_ACCUMULATION_LIFE`
- [x] `PEGA_MARKETINGOFFICER`
- [x] `INSERTUPDATERISKADDRESS`

Utilitas/pendukung:
- [ ] `PROSESCOPY` — ❎ N/A: dikonfirmasi TIDAK ADA di DB (2026-09-16)
- [x] `GENERATE_FACRETRO_NO`
- [x] `GET_TOKEN_STORAGE`
- [x] `GETCURRENCYSTANDARD`
- [x] `PROC_GENERATE_SEQUENCE_NUMBER`

> Catatan: 3 objek dikonfirmasi TIDAK ADA di DB (bukan kekurangan): `FACINOFFERLIFE`, `PROSESCOPY`,
> `M_LIMIT_LIFE`. Jadi target efektif 32 prosedur, bukan 35. `FACINFORBACKUP` masih perlu konfirmasi.

### A.2 · Isi + DDL tabel limit  (butir 2 · membuka `services/acceptance` versi final)
Tersimpan di `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` (DDL) + `.xls` (isi).

**Status: ✅ LENGKAP — 6 dari 6 target efektif** (bukan 7; `M_LIMIT_LIFE` dikonfirmasi tidak ada di DB).

⛔ **`[terverifikasi]` Ada DUA BENTUK yang tidak kompatibel, bukan satu.** Dikunci **K-023** +
amandemen `adr/0003-…` (17 September 2026).

**Bentuk A — standar, 14 kolom** · ejaan `JABATAN` **TANPA SPASI** (cocok token routing) ·
punya `WORKBASKET` (terisi), `JABATAN_ATASAN`, `TEAM_GROUP`, `EFFECTIVE_DATE`, `LOGIN`
Kolom limit: `MAX_LIMIT_IDR`, `LIMIT_BOTTOM`, `MAX_LIMIT_USD`, `LIMIT_BOTTOM2`
- [x] `M_LIMIT_PROPERTYY` — 14 kolom
- [x] `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` — 14 kolom
- [x] `M_LIMIT_PROPERTY_NON_PREFERREDD` — 14 kolom
- [x] `M_LIMIT_NONPROPANDENGG` — 14 kolom
- [x] `M_LIMIT_ENGINEERINGG` — 14 kolom

**Bentuk B — financial, 7 kolom** · ejaan `JABATAN` **PAKAI SPASI** (TIDAK cocok token routing) ·
**tanpa** `WORKBASKET` / `JABATAN_ATASAN` / `TEAM_GROUP` / `EFFECTIVE_DATE` / `LOGIN`
Kolom limit: `LIMITBOND_BOTTOM`, `LIMITCREDITCL_BOTTOM`, `LIMITCREDITNCL_BOTTOM`, `LIMITTRADE_BOTTOM`
- [x] `M_LIMIT_FINANCIALINS` — 7 kolom

- [ ] `M_LIMIT_LIFE` — ❎ **N/A: dikonfirmasi TIDAK ADA di DB** (2026-09-16). Bukan tertunggak; jangan
      dihitung lagi sebagai kekurangan.

#### Yang masih menahan `services/acceptance` versi final

- [x] Ejaan `JABATAN` bentuk A — `[terverifikasi]` **cocok** dengan token routing. Asumsi ADR-0003 benar.
- [x] ✅ **Eskalasi bentuk B — SELESAI 17 September 2026 (K-026).** **Tidak ada eskalasi berjenjang.**
      `[terverifikasi]` Tiga rule di `NB FacIn\RDBList\` (`GetLimitAkseptasiBond_SQL`,
      `…KreditCL_SQL`, `…KreditNCL_SQL`) hanya menyaring ambang:
      `WHERE <kolom_limit> <= {DataSearch.CARID2}` — tanpa `JOIN`, tanpa `JABATAN_ATASAN`,
      tanpa `ORDER BY`. ⚠️ **`WHERE` tetap diport**; yang tidak ada hanya langkah naiknya.
- [x] ✅ **Antrean bentuk B — SELESAI 17 September 2026 (K-026).** Memang **tidak ada antrean per
      jabatan**: kolom `WORKBASKET` tidak ada di tabelnya dan tidak satu pun query membacanya.
- [x] ✅ **Baris sisa `JUW_A` — SELESAI 17 September 2026.** Work owner **menghapus `JUW_A` dari
      `M_LIMIT_PROPERTYY`**. `[terverifikasi]` (LastWrite 17 Sept 2026 13:36, 37 baris): `JUW_A`
      **0 kemunculan**; `JUW_B` masih ada; token jabatan lain utuh. Pertanyaan "disaring vs
      dibiarkan" **gugur** — tidak ada lagi baris tujuan `JUW_A`. **K-024 ditutup.**
      ⚠️ Rule `ToJUW_A` **TETAP DIPORT apa adanya** — ini konsekuensi **data**, bukan penonaktifan
      rule. `[terverifikasi]` alur masih menulis `LetterNo = "JUW_A"` secara literal di
      `Activity\GetLimitAkseptasi_JUW_UW.xml` L3105, **tanpa membaca tabel limit**. Rincian di K-024.
- [x] **Kolom `NAMA` & `LOGIN` wajib dibuang** dari fixture — memuat nama orang (**K-025**).

**Konsekuensi status:** ✅ **keduanya siap.** Bentuk A memakai **tangga** `JABATAN` → `JABATAN_ATASAN`;
**bentuk B: eskalasi & antrean sudah diputus — filter ambang limit saja** (K-026), tanpa langkah naik
dan tanpa antrean. Keduanya boleh memakai data nyata, **masing-masing dengan mekanismenya sendiri** —
tidak disatukan di balik satu abstraksi.

### A.3 · DDL semua tabel (tipe kolom)  (butir 13 · membuka `models`)
Status: ✅ MASUK (17 Sept 2026) — DDL produksi + lookup lengkap; `models` terbuka
- [x] DDL `FACINOFFER` — MASUK. `[terverifikasi]` RATE = VARCHAR2(100) (TEKS!), TSI/PREMI = NUMBER(20,4)
- [x] DDL 6 tabel `M_LIMIT_*` — MASUK (dua bentuk: A 14 kolom, B 7 kolom — lihat A.2)
- [x] DDL 11 tabel produksi/penyimpanan — MASUK: `JSON_POLIS`, `JSON_OFFER`, `FACINPRODUCTION`,
      `FACOUTPRODUCTION`, `FACINPRODUCTION_BACKUP`, `TREATYINPRODUCTION(_BACKUP)`, `CEDING_FACINPRODUCTION`,
      `JSON_POLIS_MONITORING`, `JSON_KLAIM`, `C_COUNTER_PRODKE`
- [x] DDL 15 tabel lookup/master — MASUK: `M_CURRENCYSTANDARD`, `M_CLIENT`, `M_CITY`, `M_DISTRICT`,
      `M_NATION`, `M_PROVINCE`, `RISKADDRESS`, `SHIP`, `MARKETINGOFFICER`, `M_ACCUMULATION(_LIFE)`,
      `M_ACCUMULATEDTYPE`, `M_TREATY_IN`, `M_SITE_DATABASE`, `M_LINK_SERVICE`
- [x] Konfirmasi tipe: **RATE bertipe TEKS (VARCHAR2)** — perbandingan angka di atasnya = perbandingan
      string. Ini MEMBUKTIKAN §4.1/R2 yang tadinya [dugaan] → kini [terverifikasi]. Uang = NUMBER(20,4).
- [ ] **Tindak lanjut Kiro:** baca DDL `JSON_POLIS` + `FACINPRODUCTION` saat merancang tabel flat (butir E)

### A.4 · Riwayat akseptasi  (ekstraksi butir 5)
Status: ✅ TERSELESAIKAN — keputusan work owner 2026-09-16 (→ catat sebagai K-022)

`[terverifikasi]` DUA tabel riwayat yang BERBEDA — dibuktikan dari isi rule:
- `POOLDATA.HISTORYAKSEPTASIPEGA` — **TETAP DIPAKAI, tidak berubah.** Sumber:
  · `GetAksepBanding_SQL` → `FROM HISTORYAKSEPTASIPEGA` + join `POOLDATA.M_LIMIT_PROPERTYY` (jalur banding)
  · `GetFlagReject_SQL` → `FROM HISTORYAKSEPTASIPEGA` (flag reject)
- `DATAPEGA.pc_History_ASM_FW_GISFW_Work` — ❎ **TIDAK DIPAKAI LAGI** (tabel internal Pega, musnah
  bersama Pega). Rule pembacanya `GetHistoryAccPega_SQL` juga **TIDAK DIPAKAI LAGI**.

⚠️ Koreksi: TIDAK ada penggantian ke `HISTORYAKSEPTASIPRODUCTION`. HISTORYAKSEPTASIPEGA tetap seperti semula.

- [x] `DATAPEGA.pc_History_*` + `GetHistoryAccPega_SQL` → tidak dipakai lagi
- [x] `HISTORYAKSEPTASIPEGA` → tetap dipakai (jalur banding + flag reject)
- [x] **DBA**: DDL `POOLDATA.HISTORYAKSEPTASIPEGA` — MASUK (`DDL\HISTORYAKSEPTASIPEGA.txt`, 17 Sept 2026)
- [x] ✅ Arsip: **tidak ada di arsip** (keputusan work owner 17 Sept 2026)
- [x] ✅ Kolom `ID_KOMITE`: **diabaikan** — dulu dipakai untuk klaim; **ke depan tidak dipakai lagi**
      (keputusan work owner 17 Sept 2026)

**Status A.4: ✅ tuntas — diformalkan sebagai K-028** (17 Sept 2026). Konsekuensi yang dicatat di sana:
karena riwayatnya utuh (tidak pernah dipangkas), **hitungan riwayat berstatus reject** yang menentukan
ketersediaan jalur banding dapat dipercaya apa adanya.

### A.5 · Pengukuran (kuesioner DBA — menutup pertanyaan tanpa pendapat)
Status: 🟡 TABEL & KOLOM SUDAH DIIDENTIFIKASI dari DDL → tinggal DBA jalankan query
- [x] Tabel & kolom target teridentifikasi dari DDL yang sudah masuk (lihat rincian di bawah)
- [x] ✅ **D1 MASUK (17 Sept 2026)** — `DDL\D1.xml`: total **691.925** baris, **4 tanpa mata uang
      (0,0006 %)**. Memperkuat **K-012**: `Unknown` nyata tetapi sangat langka → default IDR tetap
      ditolak, `panic` tanpa kecuali tetap ditolak. ⚠️ Justru karena langka, cacatnya **tidak akan
      tertangkap pengujian sampel**.
- [x] ✅ **D2 MASUK (17 Sept 2026)** — `DDL\D2.xml`: 30 nilai `RATE` terbanyak. Sebaran **campuran
      dalam satu kolom** (dari `0,00000476581` sampai `1,43`) — konsisten dengan dua skala hidup
      berdampingan. ⚠️ **Tidak dipecah per COB**, jadi **tidak** mengubah peta skala K-018; penguat
      saja. Memunculkan temuan terpisah → **K-027** (koma desimal).
- [ ] ~~**D1**: hitungan nilai uang tanpa mata uang.~~ *(arsip — rincian query asli)* Tabel + kolom mata uang dari DDL:
      · `FACINOFFER` → kolom uang TSI/PREMI/TSITOP `_MENJADI`/`_SELISIH` (NUMBER(20,4)); mata uang = `CURR_ID`
      · `FACINLIFE` → 7 kolom uang; ada 1 kolom currency
      · `FACINSPREADLIFE` → 6 kolom uang; ⚠️ **0 kolom currency** (otomatis "tanpa mata uang")
      · `TREATYPRODUCTION_BACKUP` → 13 kolom uang; 2 currency
      Query: `SELECT COUNT(*), SUM(CASE WHEN CURR_ID IS NULL OR TRIM(CURR_ID)='' THEN 1 ELSE 0 END) FROM POOLDATA.FACINOFFER;`
- [ ] **D2**: sampel nilai `RATE` per COB (membuktikan satuan ‰ vs % — konfirmasi silang K-018).
      ⚠️ `[terverifikasi]` `FACINOFFER.RATE` bertipe **VARCHAR2(100)** (teks) — waspada koma vs titik.
      Query: `SELECT RATE, COUNT(*) FROM POOLDATA.FACINOFFER WHERE RATE IS NOT NULL GROUP BY RATE ORDER BY COUNT(*) DESC FETCH FIRST 30 ROWS ONLY;`

---

## B. IT — ekspor ulang dari Pega (selagi masih hidup)

### B.1 · Ekspor produksi tunggal  (butir 1 & 11 · membuka baris keputusan final + Special Acceptance)
Status: ✅ **11 ACTIVITY K-006 TERSELESAIKAN (17 Sept 2026)** — 9 masuk, 2 dikeluarkan; **6 cabang PULIH**.
🟡 Ekspor rule produksi tunggal masih ditunggu untuk ketidaksepadanan versi antar folder — bukan untuk K-006.

Sampel kasus (17 Sept 2026) — 5 berkas kasus di `DDL\`:
`EDM-13445 (FIRE)`, `NB-176005 (AS. KREDIT)`, `NB-181231 (FIRE)`, `NB-184233 (MARINE CARGO)`, `RNW-10579 (FIRE)`
- [x] Sampel kasus lintas siklus (NB/RNW/EDM) + beberapa COB — berguna untuk rekonsiliasi
- [x] ✅ **11 activity "hilang" TERSELESAIKAN (17 Sept 2026)** — **9 rule MASUK** ke `DDL\`,
      **2 dikeluarkan** atas keputusan work owner. Diformalkan sebagai **amandemen K-006**.
- [x] ✅ **Keenam cabang K-006 PULIH** dari status ⏸ DITANGGUHKAN — masing-masing diverifikasi ke
      korpus dengan bukti baris rujukan:

      | # | Cabang | Bukti |
      | ---: | --- | --- |
      | 1 | Special Acceptance | `GetLimitAkseptasi_Act` L4008/L4265 · `_ActFlow` L4003/L4260 (`<RequestType>`) |
      | 2 | Tangga akseptasi putaran kedua | `SetValidateDate_PostAct` L2031 (`Call`) |
      | 3 | Limit tambahan treaty type | `SpreadingAdditionalProtection` L1963 (`Call`) |
      | 4 | Kapasitas treaty | `SumTreatyCapacity_Act` L2562 (`call`) |
      | 5 | Simpan produksi endorsement bonding | `SaveFacinProdAllEDM_Act` L884 (`Call`) |
      | 6 | Penanganan galat konversi produksi | rantai 2 langkah lewat `UpdateErrorNoteJsonPolisMonitoring` |

      ⚠️ **Dua koreksi tercatat di K-006:** (a) `GetLimitAkseptasi1SA/2SA_Act` dan
      `UpdateErrorNoteJsonPolis` **bukan Activity** melainkan rule integrasi/SQL kelas `Int-policyjson`;
      (b) cabang 6 **rantai dua langkah**, bukan rujukan langsung — catatan lama tertipu kecocokan awalan.
- [ ] ⚠️ **`CountASMGrossPremi_ACT` & `CountASMNetPremi_ACT` dikeluarkan — meninggalkan RUJUKAN
      MENGGANTUNG.** `[terverifikasi]` keduanya masih dipanggil `SetErrorMessage_Act.xml` L518 & L618
      (`<pyStepsActivityName>`). `[pertanyaan terbuka]`: langkah pemanggilnya ikut usang (diport apa
      adanya, §1) **atau** `SetErrorMessage_Act` perlu ditinjau? ⛔ Jangan hapus langkah pemanggil
      diam-diam; jangan simpulkan activity itu mati tanpa bukti. Pola sama K-019/K-024.
- [ ] Ekspor **rule** produksi tunggal (satu titik waktu, 3 siklus) — masih berguna untuk
      ketidaksepadanan versi antar folder (95 rule), **tetapi tidak lagi memblokir 6 cabang K-006**

### B.2 · DecisionTable berbaris  (butir 1 · arah keputusan underwriting)
Status: 🟡 SEBAGIAN MASUK (17 Sept 2026)
- [x] `IsUWAccepted.xml` — MASUK (berisi)
- [x] `isApproved.xml` — MASUK (berisi)
- [ ] DecisionTable lain (bila masih ada yang relevan) — menyusul bila diperlukan

### B.3 · ~~Ekspor rule **RNW** yang lengkap~~  ⛔ **DIBATALKAN 18 September 2026**
Status: ✅ **TIDAK DIPERLUKAN — dijawab work owner** (K-032 **tertutup**; P-9 dibatalkan)

> Dugaan "ekspor RNW tidak lengkap" **tidak terbukti**:
> `GenerateNoPolicy` **usang** (simpan produksi NB & renewal memakai `SaveJsonPolicyFacIn_Act` yang
> `[terverifikasi]` identik byte-per-byte — panggilan tersisa = sisa kelewat, **diport apa adanya**) ·
> `SumTreatyCapacity_Act` ✅ **sudah ada di `DDL\`** (139,9 KB, 17 langkah) ·
> `SearchJobID` **usang, sudah lama dihapus**.
>
> 📌 **Jalur simpan produksi renewal = jalur simpan NB.** Renewal tidak perlu penghasil nomor polis
> tersendiri. Badan permintaan di bawah dipertahankan sebagai jejak.

*Rumusan asli (arsip):*

`[terverifikasi]` Tiga rule **dirujuk dari berkas yang ada di RNW**, tetapi **berkasnya tidak ada di
RNW** — ketiganya **ada di NB**:

- [ ] **`GenerateNoPolicy`** (RDBList, lewat `<RequestType>`) — ⛔ **penghasil nomor polis**, dirujuk
      `SaveJsonPolicyFacIn_Act` L2204 dan `GenerateNopolis_Act` L3179. **Paling berdampak**: tanpa ini
      jalur simpan produksi renewal kehilangan penghasil nomor polis
- [ ] **`SumTreatyCapacity_Act`** (Activity, lewat `Call`) — dirujuk `SetValidateDate_Act` L1129 dan
      `CountASMShareTotal_ACT` L444. **Menyentuh cabang 4 K-006** (kapasitas treaty)
- [ ] `SearchJobID` (RDBList, lewat `<RequestType>`) — dirujuk `UploadCSVPerson_PostAct` L2404
      *(dilaporkan Claude; belum diperiksa ulang work owner)*

⚠️ **Yang diminta:** ekspor rule RNW **lengkap**, analog B.1. Cukup untuk memastikan ketiganya memang
tidak ada di ruleset renewal, atau sekadar **tidak ikut terekspor**.

⛔ **Jangan divonis usang.** `CLAUDE.md` §4.5 — dan di sini bobot "ekspor tidak lengkap" **lebih
besar** dari biasanya: `[terverifikasi]` 1.907 dari 1.907 berkas bersama identik **byte-per-byte**,
metadata ekspor termasuk, sehingga kedua folder berasal dari keadaan ruleset yang sama.

---

## C. Product — arti kode & daftar nilai

### C.1 · Arti enumerasi  (butir 6 · membuka klasifikasi produk)
Status: ✅ **TUTUP — diformalkan K-029 (17 Sept 2026). Tidak ada lagi enumerasi menunggu Product.**

> `[terverifikasi]` Keenam berkas memuat `<pyPromptTableList>` berisi pasangan `<pyStandardValue>`
> (kode) + `<pyLocalizedValue>` (arti). **Nilai dan artinya terbaca langsung dari ekspor.**
> Pemetaan lengkap: **K-029** (otoritas) dan `steering\GLOSARIUM.md` (kosakata).

- [x] `Type.xml` — ✅ **10 kode + artinya terbaca** (jenis penyesuaian: Extend Period, Adjustment TSI, …)
- [x] `EdmType.xml` — ✅ **3 kode + artinya terbaca** (`1` Batal Sejak Semula · `2` Batal Prorata ·
      `4` Penambahan/Pengurangan/Perubahan)
- [x] `TypeDeductible.xml` + `TypeDeductible2.xml` — ✅ **terbaca; ternyata DUA property berbeda**,
      dan kode `4` berarti hal berlainan di masing-masing. Jangan disatukan.
- [x] `TeamGroup.xml` — ✅ terbaca. ⚠️ `5` = **Bonding**, memutus pola "Group N"
- [x] `ProRateType.xml` — ✅ **4 kode + artinya terbaca**
- [x] `B2B.xml` — ✅ `[terverifikasi]` **4 nilai**: `ASM` · `KBRU` · `BDX` · `SRB`.
      ⛔ **TIDAK menutup `IsPKSASM`** — lihat catatan di bawah.
- [x] ✅ **`IsB2B` DITUTUP (K-029)** — **flag penanda**; nilainya sendiri adalah isinya, tidak ada arti
      bisnis lebih dalam. Diperlakukan apa adanya sebagai kode string.
- [x] ✅ **Selisih daftar TERJAWAB (K-029)** — `Type = 3`/`5` dan `EdmType = 3` adalah **kode usang
      yang sudah dihapus**; jejak di rule = sisa yang kelewat. Pola sama K-019/K-024.
      ⚠️ **Sikap porting BERBEDA:** `Type 3/5` efektif setara `IsEdmAdjTSI` (aman);
      **`EdmType = 3` cabang `× −1` WAJIB diport apa adanya** — `[terverifikasi]` masih ada di
      **7 berkas folder NB**, menghapusnya mengubah perilaku record lama.

#### ⛔ `IsPKSASM` TETAP `panic` — daftar nilai `IsB2B` tidak mencabutnya

`[terverifikasi]` Penyebab `panic`-nya **bukan** nilai yang tidak dikenal, melainkan tiga hal yang
tidak satu pun disentuh daftar nilai: `<pyConditionString>` masih placeholder
`[Double click to add condition]`, **label belum ter-resolve**, `<pyTempText>` = `true` — di ketiga
folder. Mengetahui nilai apa saja yang sah **tidak** membuktikan kondisinya dieksekusi.
Dicabut hanya lewat ekspor produksi yang labelnya sudah ter-resolve. Rincian di
`04-spec\03-spec-modul-terverifikasi.md`.

### C.2 · Isi dropdown  (butir 7 · membuka layar input)
Status: 🟡 HAMPIR TUNTAS (17 Sept 2026) — hanya `TABLEOFLIMIT` (perlu re-save) tersisa
- [x] MASUK: `CITY`, `DISTRICT`, `RW`, `LST_BANK_GROUP`, `OCCUPATION`, `BRANDDETAIL`,
      `M_KLAUSUL_PLAN`, `M_TYPE_PROPERTY_PLAN`
- [x] ⛔ Dikonfirmasi TIDAK ADA lagi (kelewat saat hapus): `V_ZIPCODE`, `V_COINS`,
      `VIEW_BENEFIT_PROPERTY`, `VJ_M_TYPE_PROPERTY_PLAN` (tidak dipakai lagi)
- [x] ⛔ Tabel enumerasi `Data-Enumeration` — **tidak dipakai lagi**; panggilan sisa di Section = kelewat.
      Tidak perlu diekstrak. (Membatalkan rencana "satu tabel enum banyak dropdown" di `C2-dropdown-sumber.md`.)
- [x] ✅ `TABLEOFLIMIT` (view) — `DDL\TABLEOFLIMIT.xml` MASUK & terbaca. `[terverifikasi]` kolom:
      `ID`, `BIZCODE`, `NOTE`, `TAHUN`, `CATEGORY`, `DESCRIPTION`, `PCTLIMIT` (1 record contoh).
      `PCTLIMIT` pakai koma desimal. Cukup untuk merancang bentuk flat.
- [x] Nama tabel "(konfirmasi DBA)" — sebagian dikonfirmasi lewat kehadiran DDL (CITY/DISTRICT/RW dll. cocok)

**Status C.2: ✅ tuntas** — semua sumber yang masih dipakai sudah masuk; sisanya dikonfirmasi tidak dipakai lagi.

---

## D. IT + Underwriting

### D.1 · Mesin spreading terbesar  (butir 12 · membuka `services/spreading`)
Status: ✅ MASUK + dikonfirmasi (17 Sept 2026)
- [x] `GetKapasitasTreaty.xml` (±1,02 MB) — diekspor ulang, ada di `DDL\`
- [x] Dikonfirmasi work owner: **masih dipakai** → diport apa adanya (reproduksi perilaku, CLAUDE.md §1)

---

## E. Keputusan work owner (Anda)

- [ ] **Butir 4** — struktur tabel penyimpanan: flat relasional (pengganti JSON_POLIS) — desain + ADR
      *(bentuk sudah disepakati; menunggu spec + ADR dibuat)*
- [ ] Konfirmasi hubungan tabel flat dengan JSON_POLIS (pengganti — sudah diputuskan; formalkan sebagai ADR)

---

## G. Renewal (RNW) — hasil discovery, menunggu keputusan work owner

Sumber: `06-rnw\01`…`04`. Diformalkan sebagai **K-031** (temuan), **K-032** (rujukan menggantung),
**K-033** (pertanyaan terbuka).

### G.1 · Sudah selesai — tidak menunggu apa pun

- [x] ✅ `[terverifikasi]` **1.907 dari 1.927 berkas RNW identik byte-per-byte dengan NB**; delta sejati
      **20 berkas** → **renewal memakai ulang seluruh modul NB** (K-031)
- [x] ✅ Kedua puluh berkas delta **terpetakan seluruhnya** (alur masuk · layar input · periode ·
      portal · pemilih tertanggung · kelompok bisnis · sunting marketing · daftar renewal · konversi)
- [x] ✅ `[dilaporkan Claude]` **Renewal tidak punya mesin hitung ulang** — menjalankan perhitungan NB
      yang sama; `OldPolicyNo` hanya rujukan. **Menjelaskan K-015 dari sisi mekanisme**
- [x] ✅ `[terverifikasi]` Empat nama yang tampak menggantung **ternyata bukan** — `isApproved`,
      `IsUWAccepted`, `GetLimitAkseptasi_ActFlow`, `GetInsuredID`. Tidak perlu diperiksa ulang

### G.2 · Menunggu pihak luar

- [x] ✅ **TIDAK ADA LAGI.** Permintaan ekspor rule RNW ke IT (**B.3 / P-9**) **dibatalkan
      18 September 2026** — ketiga rujukan dijawab work owner (K-032 tertutup). **Tidak ada butir RNW
      yang menunggu pihak luar.**

### G.3 · Menunggu keputusan work owner  (K-033)

- [x] ✅ **RNW-1 SELESAI** — ketiga rujukan terjawab: `GenerateNoPolicy` usang · `SumTreatyCapacity_Act`
      sudah masuk `DDL\` · `SearchJobID` usang. **Bukan ekspor tidak lengkap.** Dua yang usang
      **diport apa adanya**, status usangnya dicatat (K-032 tertutup)
- [x] ✅ **RNW-2 SELESAI — K-034.** Gerbang masuk renewal **ditetapkan**: No. Polis + Renewal Date +
      Note → OK → activity pembuat renewal (salin data polis). `[terverifikasi]` didukung
      `Renewal_FlowAct` (`pySectionReference`=`InputRenewal`, pre=`InputOfferFacInEngineer_preACT`,
      post=`SetValidateDate_PostAct`). **Keputusan bisnis**, mengisi celah yang korpus tak merekam
- [x] ✅ **RNW-3 SELESAI — K-035.** Pola konversi produksi renewal **sama seperti NB**.
      ⚠️ Celah K-004 **tetap terbuka** (rule asli kelas `Work` masih hilang)
- [x] ✅ **RNW-4 SELESAI — K-036.** `Section\EditMarketing` + `FlowAction\EditMarketing`
      **dikeluarkan dari lingkup port**; pemilihan marketing mengikuti NB.
      📌 **Temuan keamanan `.Password` GUGUR** — komponennya tidak diport, tinjauan §6 tidak relevan.
      Permukaan delta renewal turun dari **20 → 18 berkas** yang perlu dirancang
- [x] ✅ **RNW-5 SELESAI — K-037.** Blok `pySaveSQL` ber-`PROSESCOPY` **DIBUANG**.
      ⚠️ Dicatat sebagai **perubahan perilaku yang disengaja**, bukan porting. `pyBrowseSQL` pada rule
      yang sama **tetap diport**
- [x] ✅ **RNW-6 SELESAI.** xlsx dibaca (kolom `XML` **tidak disentuh**) →
      `06-rnw\05-konfirmasi-silang-peta-struktur.md`

**Status G.3: ✅ keenam butir tertutup. Tidak ada lagi pertanyaan terbuka RNW.**

### G.5 · Keputusan lanjutan tahap spec (18 September 2026)

- [x] ✅ **K-038 — fitur pilih-tertanggung DIBUANG** (4 berkas: `Harness\ChooseInsured` ·
      `Section\ChooseInsuredDtl` · `ReportDefinition\BrowseAccountInsuredEDM` ·
      `Activity\SetDataInsuredEDM_Act`). Renewal mengambil tertanggung dari **polis lama**, bukan
      pemilihan CRM/SFA. `[terverifikasi]` dua di antaranya berkelas **SFA**
- [x] ✅ **K-039 — `.OldTSI` terjawab.** Data polis lama **disalin jadi nilai awal saat pembuatan
      kasus**; perhitungan lalu berjalan seperti NB atas hasil salinan. K-031 tetap berdiri.
      ⚠️ **Atribusi teknis dikoreksi**: `.OldTSI` (prefiks) ditulis activity **perhitungan**
      (`CountPremiumNet`, `SumTSIPremiSpreadedRNM_Act`), sedangkan snapshot nilai awal memakai
      properti ber-**sufiks** `Old` lewat `Activity\SetOldData` (63 langkah, identik NB↔RNW).
      **Dua properti berbeda** — `TSIOld` ≠ `OldTSI`
- [x] ⏸ `SetOldData` (via `GetPolicyData_ACT`) — **DITUNDA ke fase Endorsement (EDM)** atas keputusan
      work owner. `[dugaan]` kemungkinan milik fase EDM; tidak masuk lingkup delta renewal, tidak
      dianalisis lebih jauh. ⚠️ `TSIOld` (sufiks) ≠ `.OldTSI` (properti perhitungan)
- [x] ✅ **Lima section menggantung — TERJAWAB 18 September 2026** (temuan saat menulis spec RNW):
      **K-040** — `InputInwardFacultativeSuggest` · `InwardFacIn` · `OfferFacIn_NusaRe` ·
      `OfferFacIn_NusaRe_IsUW` **diwarisi dari NB** (konsekuensi K-031, bukan menggantung; delta tetap 14).
      **K-041** — `InputDtlObject` **usang**, digantikan `InputDtlObject_FacIn`; 4 sisipan tersisa
      **diport apa adanya**. ⚠️ Keduanya **section berbeda**, dibedakan hanya sufiks

**Status G.5: ✅ tuntas. Tidak ada pertanyaan terbuka RNW yang tersisa.**

📌 **Delta renewal final: 20 − 2 (K-036) − 4 (K-038) = 14 berkas.**

### G.4 · Belum dikerjakan (bila diperlukan)

- [ ] Isi `Section\InputRenewalDtl` + `_IsUW` (4,65 MB) — per-field; belum diperlukan sampai tahap spec
- [x] ✅ Isi `Struktur_…xlsx` — **selesai** (RNW-6)

---

## H. Endorsement (EDM) — **fokus aktif**, baru dimulai

Ditetapkan **K-042** (18 September 2026). Orientasi: `07-edm\01-orientasi-discovery-edm.md`.

### H.1 · Peta skala `[terverifikasi]`

| Ukuran | RNW (pembanding) | **EDM** |
| --- | ---: | ---: |
| Total `.xml` | 1.927 | **2.061** |
| Bernama sama dengan NB | 1.907 | **1.707** |
| — **identik byte-per-byte** | **1.907** | **0** |
| Eksklusif (tidak ada di NB) | 20 | **354** |

⛔ **EDM bukan "NB dengan pintu masuk berbeda".** Nol berkas identik → **strategi pakai-ulang NB
(K-031) TIDAK berlaku**; setiap berkas bernama sama wajib dibandingkan isinya.

### H.2 · Rencana discovery bertahap

- [ ] **E-1** Pola perbedaan pada 1.707 berkas bersama — uji hipotesis pola P1 (agregat tersimpan vs
      halaman aktif). **Daya ungkit tertinggi**: bila satu pola menjelaskan mayoritas, 1.707 berkas
      jadi satu keputusan rancangan
- [ ] **E-2** Alur endorsement — 2 Flow (paling sedikit dari ketiga siklus) + FlowAction, gerbang `EdmType`
- [ ] **E-3** Before/after image + `SetOldData` (K-039) — mekanisme **selisih**, inti endorsement
- [ ] **E-4** 354 berkas EDM-only — Activity (94) & Section (81) lebih dulu
- [ ] **E-5** Jalur produksi endorsement — `_MENJADI`/`_SELISIH`, `PRODKE`
- [ ] **E-6** Celah & penutup — rujukan menggantung, kode usang, pertanyaan terbuka

### H.3 · Pertanyaan terbuka dari orientasi

- [ ] **E-Q1** Apakah resolusi rule Pega **peka huruf**? Berdampak pada **12 berkas** EDM
      (`IsCar`↔`IsCAR`, dll.) dan menentukan angka resmi 1.707/354 vs 1.719/342 · **IT**
- [ ] **E-Q2** Mengapa EDM hanya punya **2 Flow** (NB 6)? Alur terpusat, atau sebagian tidak terekspor?
- [ ] **E-Q3** Apakah **K-038** (pilih-tertanggung dibuang) berlaku juga untuk endorsement?
      `[terverifikasi]` keempat berkasnya **ada di EDM dengan isi berbeda** · **work owner**

### H.4 · Yang dibuka kembali dari penundaan

- [ ] Tiga pertanyaan di `04-kuesioner\_DITUNDA-fase-endorsement.md` — **D-1** revise pada endorsement ·
      **D-2** banding vs ask diperlakukan sama · **D-3** nilai mutlak pada selisih. Kini relevan
- [ ] **K-039** — `Activity\SetOldData` + properti sufiks `TSIOld`, ditunda dari RNW ke fase ini

---

## F. Sudah selesai — tidak menunggu apa pun (arsip status)

- [x] **Butir 9** — predikat fac out + 10 perujuk: diport apa adanya (K-019)
- [x] **Butir 10** — alamat email di kode: pengecualian tercatat (K-020)
- [x] Arti `ProposalAcceptStatus` {1,2,3,4,7,9}: tertutup dari korpus (K-021)
- [x] Satuan rate per COB: terkunci (K-018) — PA/Layering/FIRE=‰, MBU/ANEKA/BONDING/GOLF/MARINE CARGO=%
- [x] Dasar akseptasi renewal = TSI penuh (K-015)
- [x] 5 modul terverifikasi: spec selesai (money, ratio, rules, acceptance, premium)

---

## Peta: apa membuka apa

| Jika ini masuk | Modul yang terbuka |
| --- | --- |
| A.1 (35 prosedur) + A.3 (DDL) | `services/production`, `repository`, `models` |
| A.2 (tabel limit) | ✅ **TERBUKA PENUH** — `services/acceptance` boleh ganti fixture → tabel nyata untuk **kedua bentuk**: bentuk A lewat tangga `JABATAN`→`JABATAN_ATASAN` (`JUW_A` sudah bersih, K-024), bentuk B lewat **filter ambang limit** tanpa eskalasi (K-026). Dua mekanisme terpisah, bukan satu (K-023) |
| B.1 (9 rule K-006 masuk) | ✅ **TERBUKA** — **Special Acceptance** + **tangga akseptasi putaran kedua** pulih untuk `services/acceptance` (amandemen K-006), berikut limit tambahan treaty type, kapasitas treaty, produksi EDM bonding, penanganan galat konversi |
| B.2 (DecisionTable) | baris keputusan final — `IsUWAccepted.xml` & `isApproved.xml` sudah masuk, **belum dibaca** |
| C.1 + C.2 (Product) | layar input lengkap, klasifikasi produk |
| D.1 (spreading) | `services/spreading` |
| E (keputusan Anda) | struktur penyimpanan → lalu `to-tickets` mencakup semuanya |

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan. Semua daftar terverifikasi dari korpus.*
