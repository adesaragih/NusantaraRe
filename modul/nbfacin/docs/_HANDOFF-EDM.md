# HANDOFF — Status Discovery EDM + Tabel Flat + Fac Out (untuk dilanjutkan)

> Diperbarui 21 September 2026. Lanjutkan dari bagian "LANGKAH BERIKUTNYA".
> Peran Kiro: verifikator hasil Claude + pembuat prompt (tidak jalankan skill sendiri). Balas Bahasa
> Indonesia. Tiap klaim faktual diverifikasi ke korpus sebelum dipakai (anti-halusinasi CLAUDE.md §3).

---

## ⏩ TAMBAHAN 21 SEPT 2026 — TABEL FLAT & FAC OUT

### Tabel flat (pengganti JSON_POLIS.DATA_JSON)
- Rancangan awal dari 5 contoh kasus nyata (`DDL\P-5 *.txt`): `08-flat\01-rancangan-tabel-flat-dari-contoh-nb.md`.
- [terverifikasi] Struktur NB/RNW/EDM SAMA (pxObjClass identik; EDM +Int-INWARDSCALE) → **satu skema
  flat lintas siklus**, dibedakan PRODKE. CoverageList+SpreadingList = simpul konvergen semua COB.
  Field before-image EDM (`*Old`) ada di level Coverage.
- Menunggu: kunci versi (IDPEGA vs IDPEGA+PRODKE), skema Layer/AdditionalCoverage/Person (kosong di
  contoh), bentuk final tabel scoring. Semua ditandai [pertanyaan terbuka].

### Fac Out / Retrosesi keluar — DISCOVERY + 14 TIKET SELESAI
- Dokumen: `09-facout\01-temuan-dan-rancangan-facout.md`.
- ~~Tiket: `05-tickets\facout\F01..F08` + `00-INDEKS-FACOUT.md` [terverifikasi Kiro: label lengkap,
  nol rujukan maju, nol PII bernilai].~~
  ⛔ **DIBATALKAN 21 September 2026.** Pernyataan itu **tidak benar terhadap filesystem**: saat
  ditulis, `05-tickets\facout\` **kosong (0 berkas)**, sehingga label "[terverifikasi Kiro]" tidak
  dapat dipertahankan. **Penggantinya:** tiket Fac Out kini **F01…F14** (14 tiket) + `00-INDEKS-FACOUT.md`
  = **15 berkas**, diterbitkan 21 September 2026.
- **Berkas audit pendukung** (angka yang dikutip tiket, dapat diaudit ulang):
  `10-audit\01-pohon-langkah-countrateretrocov.md` · `02-peta-kode-transisi.md` ·
  `03-struktur-facretro-dan-populasi.md` · `04-perujuk-isfacout-peka-huruf.md` ·
  `05-draf-keputusan-facout.md` · `06-tinjau-ulang-suntingan-facout.md` ·
  `07-hitung-berkas-dan-drift.md`.
- **DRAF K-057…K-062 menunggu work owner** — belum satu pun diputuskan.
- Keputusan **K-053..K-056** dicatat. K berikutnya = **K-057**.
- **[terverifikasi]** rantai: UW Accept (`IsFacout`=ProposalAcceptStatus=4) → `SetValidateDateUW_PostAct`
  step 3 `Call SetDataFacOut_Act` → per-COB `SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act` →
  salin spreading ber-FacOut ke `OfferFacIn.FacRetro` (staging) & `FacRetroList` (list berindeks) +
  `FacOutObjectList{Rate,FacOutTSI}`.
- **[terverifikasi]** penanda: `TreatyType=="10015"` (ketat) / `contains 10015|SPL|10007` (umum).
  `CoverageBasis==5` = Layering Basis (`DDL\CoverageBasis.xml`).
- **[keputusan work owner K-054]** `10015`=FACOUT, `10007`=ORS (belum ada file pemeta angka→label;
  `M_TREATY_IN` = ID+JSONDATA CLOB). `"SPL"` = [pertanyaan terbuka].
- **[terverifikasi]** tangga retro 4 tingkat di ekspor ulang `DDL\OfferFacRetro.xml`:
  Admin→Head→GroupLeader→TechnicalDirector (0 shape Approval → bukan tangga limit Fac In/Seam 2).
- **[terverifikasi]** produksi = **Connect-SQL langsung** `DDL\InsertTreatyProd_Sql.xml`
  (`INSERT INTO FACOUTPRODUCTION` 65 kolom) — BUKAN stored procedure. `FACOUTPRODUCTION` sudah flat &
  ada. Daftar: `GetFacoutList_SQL` (by policyno). Nomor slip: `GENERATE_FACRETRO_NO` → `RNM-Y...`.
- **Tidak butuh ALL_SOURCE untuk FacOut** (koreksi: insert SQL lengkap di korpus).

### Soal `/to-tickets` di Claude (FacOut) — JANGAN salah paham
- `/to-tickets` = proses interaktif (quiz user dulu), format default `Status: ready-for-agent`,
  TANPA `Asal (Pega)`/`Keputusan`, publish ke `.scratch/<slug>/issues/`. Kalau dijalankan polos,
  Claude bisa MENGHASILKAN VERSI BEDA & memindah lokasi.
- **Tiket FacOut SUDAH JADI** di `05-tickets\facout\` (pola sama EDM). Tidak wajib `/to-tickets`.
  Bila tetap mau lewat skill: instruksikan Claude pakai bahan di folder itu, pertahankan
  `Asal (Pega)`/`Keputusan`, status `blocked`, jangan pindah ke `.scratch/`.

### ~~Status proyek: **54 tiket** = NB 16 + RNW 8 + EDM 22 + FacOut 8.~~ Nol berkas .go. K s/d K-056.

⛔ **ANGKA 54 DIBATALKAN 21 September 2026** — ia menghitung 8 tiket Fac Out yang **tidak ada**.
**Penggantinya:** `[terverifikasi]` **60 tiket** = NB **16** + RNW **8** + EDM **22** + Fac Out **14**.
`[terverifikasi]` **Nol berkas `.go`** tetap berlaku. Register: **K-001…K-056** ditetapkan;
**K-057…K-062 masih DRAF** menunggu work owner (`10-audit\05-draf-keputusan-facout.md`).

---

## POSISI SAAT INI

Discovery EDM **E-1…E-6 SELESAI & terverifikasi**. Sekarang di tahap **menutup blocker EDM** lewat
dokumen `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md`.

- Lingkup EDM final = **708 berkas** (354 berbeda + 354 EDM-only), setelah **amandemen K-043**
  (kontrak volatile = 23 tag + isi blok `pzIndexes` diabaikan; angka 158 tereproduksi bit-per-bit
  oleh Kiro: 1195/512 → 1353/354).
- 7 dokumen discovery di `OUTPUT\07-edm\` (01-orientasi + 02-e1 … 07-e6? — cek: sampai 06-e5; E-6
  laporannya sudah masuk dari Claude, verifikasi Kiro selesai).
- E-Q1…E-Q32 lengkap tanpa bolong.

---

## STATUS BLOCKER EDM (di `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md`)

### P-10 (DBA) — ✅ TERTUTUP semua
- `POOLDATA.INSERTJSONPOLIS` → tak dipakai di sistem baru (insert ke tabel flat bareng NB/RNW)
- `PEGA_DELETE_ERROR_KONVERSI` & `PROC_GENERATE_SEQUENCE_NUMBER` → file sudah diperbarui, **[terverifikasi Kiro] ADA di `DDL\`** (3.019 B & 3.015 B)
- `MachingDataFacin_Sql` → tak dipakai di sistem baru
- skema `facinproduction` = `pooldata.facinproduction` → **objek SAMA** (E-Q29 tutup)
- PRODKE → `GetProdKeOldData_SQL` sudah sesuai (E-Q30 tutup)
- `OPENPROTEKSI_EDM` → di menu lain, **di luar lingkup EDM**

### P-11 (IT) — ✅ TERTUTUP semua
- `pyStepsPreCondition=false` (E-Q16) → **TIDAK di-skip, langkah TETAP JALAN**. Diverifikasi work
  owner via UI Pega (langkah "Get ProdKe" `SaveEDMToJsonPolicy_Act` tampil aktif: Loop/When/RDB-List).
  → seluruh pemetaan alur EDM valid apa adanya.
- rule produksi EDM → nopolis+insert di `SaveEDMToJsonPolicy_Act`; NB/RNW di `SaveJsonPolicyFacIn_Act`.
  Menutup `GenerateEndorsementNo`/`GenerateEDMNoLife` (bukan panic).

### P-12 (WORK OWNER) — ✅ SEMUA TERTUTUP (dicatat K-044/K-045/K-046)
- **W-1** ✅ EDM Life **IKUT dikerjakan** (K-044).
- **W-2** ✅ Pilih-tertanggung **dibuang** (K-038 berlaku EDM; K-044). 4 berkas pilih-tertanggung tidak diport.
- **W-3** ✅ (K-045) `GetLimitAkseptasi_Act2` "Putaran 2" **MATI** — di-remark di KEDUA pemanggil
  (`InputOfferFacInEngineerUW_preACT` step 6 + `SetValidateDate_PostAct`). `GetLimitAkseptasi_Act`
  "Putaran 1" **TETAP HIDUP** (3 pemanggil lain). Bonding (`SaveFacinProdEDMBonding_Act` di DDL) status
  jalur produksi **tetap terbuka**, tidak dihapus (K-006).
- **W-4** ✅ (K-046) 9 diport apa adanya · **A.5 `PremiNusantaraReOld` "0"→angka 0 DIPERBAIKI** di sistem
  baru (satu-satunya) · B.6 (PCT_BROKERGARE_FEE kolom persen) & C.12 (IsMBU ganda) **bukan keanehan**.
  ⚠️ A.5 = perbaikan sadar → saat rekonsiliasi bisa tampak beda tipe, bukan bug.

---

## 12 KODE USANG (W-4) — sudah dijelaskan ke work owner

Sumber: `07-edm\06-e5-jalur-produksi.md` §8 + E-3. Default = diport apa adanya (CLAUDE.md §1).
1. RateOld self-referential `@If(.RateOld!="",.Rate,0)` (6/6) → rate lama jadi 0 saat pertama buka.
2. EDMPremiMenjadi 11 rumus berbeda di 7 file (tidak seragam).
3. Blok EdmType==1 hanya jalan untuk Life (non-Life syaratnya ditulis ==2).
4. Label "batal sejak semula" (harusnya EdmType 1) bergerbang ==2.
5. PremiNusantaraReOld diisi string `"0"` (50 field lain numerik 0).
6. Kolom salah eja `PCT_BROKERGARE_FEE` (skema tak berubah → apa adanya).
7. Ketidaksetangkupan kolom `*_MENJADI` vs `*_SELISIH` (5v8 facin, 2v5 facout).
8. Dua semantik PRODKE (COUNT-1 rentan vs TGL_INPUT desc aman) dalam satu rantai.
9. Label langkah "MBU" bergerbang IsMarineCargo (MBU tak jalan, MC dobel).
10. Guard idempotensi fac out hanya aktif Type==7.
11. JSON_POLIS_MONITORING hanya ditulis untuk prefiks EDMT-.
12. Cabang IsMBU ganda.

---

## LANGKAH BERIKUTNYA (saat lanjut)

✅ **SELESAI:** discovery EDM lengkap; blocker tuntas (K-044/045/046); **spec modul 1 (before-image)
selesai & terverifikasi** di `04-spec\05-spec-edm-before-image.md` (28 KB); **Seam 4 disetujui (K-047)**.

- Spec before-image sudah memuat: kontrak 3 lapis, 45 field lapis A berpola X=OldData.X + 9 menyimpang
  (koreksi arsip terverifikasi Kiro), 52/52 penugasan *Old pakai @If, A.5 (PremiNusantaraReOld→angka 0,
  perbaikan sadar), Life sebagai alur terpisah, pilih-tertanggung Out of Scope, §4.3 ditahan [belum diuji].
- Seam: 4 total (3 NB + Seam 4 `services/endorsement.PrepareBeforeImage`); selisih lewat Seam 3 diperluas.

0. ✅ **Modul SELISIH — bahan SELESAI & terverifikasi** (`04-spec\06-spec-edm-selisih.md`).
   Verifikasi Kiro: 2 tingkat delta (pembayaran + baris produksi), 5 properti lama baru dari
   `CountPremiEDM_DT` (Rule-Obj-Model, belum dibedah E-3/E-5) terkonfirmasi, seam tetap 4 (K-047),
   §4.3 ditahan, 7 keanehan K-046 jadi test. Dua keputusan ditutup di **K-048**:
   - §8.1 `ProrateEDMEnd` **port apa adanya** (1 properti ditimpa berurutan) → spec modul 1 WAJIB
     diberi catatan "nilai ProrateEDMEnd sementara, dapat ditimpa modul selisih".
   - §8.2 `pyDisabled=true` (DataTransform) = **DILEWATI** (beda dari pyStepsPreCondition=false yang
     tetap jalan!). 22 aksi disabled di CountPremiEDM_DT L478–L1075 diabaikan spec.
   ⚠️ TINDAK LANJUT: pastikan spec before-image (05) diberi catatan §8.1 sebelum /to-tickets.
1. **`/to-tickets`** (work owner jalankan MANUAL — disable-model-invocation) untuk memecah spec
   (before-image, lalu selisih) jadi tiket. Kiro TIDAK menjalankannya; Kiro verifikasi hasil tiket.
2. **Spec modul EDM berikutnya** setelah selisih (pola sama, siapkan bahan → /to-spec manual):
   - alur masuk EDM (2 Flow, gerbang IsEDM/StatusBusiness=3, 6 gerbang penolakan)
   - predikat When EDM (registry per-rule; COB baca OutData.pxResults)
   - 354 EDM-only (klasifikasi pxObjClass di E-4)
   - jalur Life terpisah (K-044)
3. **Jalur produksi EDM DITUNDA** → menunggu tabel flat (INSERTJSONPOLIS diganti insert-ke-flat, K-046/P-10).
4. Pola kerja: Kiro buat prompt (verifikasi dulu, "tidak halu", SATU blok kode), work owner jalankan di
   Claude, Kiro verifikasi hasil. Skill to-spec/to-tickets: Claude siapkan bahan → BERHENTI → work owner manual.
5. Register: K-001..K-052 + amandemen. K berikutnya = **K-053**.

## SPEC EDM — status modul (per 19 Sept 2026)
- Modul 1 before-image: ✅ spec `04-spec\05-spec-edm-before-image.md` + Seam 4 (K-047) + catatan §8.1 (K-048)
- Modul 2 selisih: ✅ bahan `04-spec\06-spec-edm-selisih.md` + K-048 (§8.1 port apa adanya, §8.2 pyDisabled=dilewati)
- Modul 3 alur masuk: ✅ bahan `04-spec\07-spec-edm-alur-masuk.md` + **Seam 5 OpenCase (K-049)** + koreksi §9 OPENPROTEKSI
- SEAM TOTAL = 5 (3 NB + Seam 4 before-image + Seam 5 open-case). Modul lain tanpa seam baru.
- Modul 4 predikat When EDM: ✅ bahan `04-spec\08-spec-edm-predikat-when.md` + **K-050** (E-Q12=A satu registry sumber-data-per-rule; E-Q5 tuntas GetBusinessType_Sql; E-Q30 PRODKE aman/tak ada penghapusan; 37 cabang dieksekusi bukan 203). COB dipilih di halaman depan NB (field Class Of Business), RNW/EDM copy dari NB.
- Modul 5 (354 EDM-only) + banding premi: ✅ K-051. 350 diport (−4 pilih-tertanggung). Rumus premi dasar SAMA lintas siklus (EDM pakai ulang mesin NB); CalculatePremi_Act NB = kelas Treaty (bukan baseline, jebakan ke-8). Seam 3 = SATU PINTU 4 BENTUK (dasar/dua-bagian/tabel-Travel/dekomposisi-delta-coverage), tanpa seam baru. ProRatePercent = PECAHAN dipakai langsung (bukan bug 100×). ProRatePercent=StartEDM+EndEDM (L6209).
- Modul 6 (jalur Life): ✅ K-052. Lingkup L1 inti 16 + L2 medis ±45; flow InputEDMLife (DITAMBAH work owner, nihil di korpus); TANPA tangga akseptasi (GetLimitAkseptasiLife_Act nihil); CalculateScorLife_Act = underwriting medis (Decline/Postpone/Standard), BUKAN premi; skoring risiko bersama NB vs medis EDM-only; premi Life = parameterisasi bentuk 1 Seam 3. **Seam 6 ScoreMedical disetujui**.
- ✅ **6 MODUL SPEC EDM LENGKAP** (before-image/S4, selisih, alur-masuk/S5, predikat-when, edm-only, life/S6). Seam total = 6.
- Belum tuntas sebelum IMPLEMENTASI (bukan spec): rumus PA/Travel penuh; CalculatePremiPA_FacIn identik NB↔EDM (23-tag); 5/6 salinan mesin premi beda isi; 4-bentuk Seam 3 = satu fungsi bercabang atau 4 impl (keputusan desain terbuka); ambang klinis Life (diport apa adanya, domain medis).
- ✅ **22 TIKET EDM SUDAH ADA & TERVERIFIKASI** di `05-tickets\edm\E01..E22` + `00-INDEKS-EDM.md`.
  [terverifikasi Kiro] file lengkap, 9 rujukan tiket NB cocok (nama file, bukan kode NB-NN), 6 seam,
  E17 & E21 BLOCKED, privasi bersih (E20 skoring medis wajib nilai sintetis). Spot-check E09/E12/E14
  konsisten dgn K-046/K-048/K-051.
  ⚠️ Catatan: Claude menulis "/to-tickets belum dijalankan" tapi FILE SUDAH ADA (23 berkas). Praktis
  tiket sudah jadi; jalankan /to-tickets manual hanya bila mau resmikan via skill (hasil sama).
  ⚠️ Rekomendasi kecil (opsional, saat implementasi): catatan E12 tulis PA "terbukti identik" NB↔EDM —
  di K-051 statusnya [dugaan] (belum 23-tag). Selaraskan jadi "verifikasi 23-tag sebelum dipakai ulang".
- **SIKLUS EDM TUNTAS SAMPAI TIKET** — sejajar NB (16 tiket) & RNW (8 tiket). Ketiga siklus lengkap.
- Jalur produksi EDM ditunda (tabel flat, K-046/P-10) — E21 blocked.
- Register: K-001..K-056 (K-053..K-056 = Fac Out). K berikutnya = **K-057**. Seam Fac In total = 6;
  Fac Out punya seam tersendiri (FO-1..FO-4, lihat `05-tickets\facout\00-INDEKS-FACOUT.md`).

## STATUS PROYEK KESELURUHAN (per 19 Sept 2026)
- NB: discovery+spec+16 tiket ✅ | RNW: discovery+spec+8 tiket ✅ | EDM: discovery+6 spec+22 tiket ✅
- **Belum ada implementasi .go** [terverifikasi] — nol berkas .go, tak ada go.mod/cmd/internal/pkg/services.
  Tiket = PETA RENCANA. Urutan kerja: NB → RNW → EDM.
- Menunggu pihak luar: tabel flat + ALL_SOURCE (DBA) untuk jalur produksi 3 siklus.
- Tabel flat (pengganti JSON_POLIS) BELUM dirancang — butir tertunda, dirancang setelah gambaran 3 siklus lengkap (sekarang sudah lengkap → bisa mulai).

## CATATAN PENTING
- Prompt ke Claude: kirim sebagai SATU blok kode (ada tombol salin). Gaya: mengalir, penanda
  `Konteks:` / `⛔ JANGAN` / `BATAS KERJA:` / `KERJAKAN:` bernomor / `BERHENTI dan laporkan`.
- UI Pega work owner MENGALAHKAN pembacaan XML (status remark `//` & pyStepsPreCondition tidak
  terserialisasi jelas ke ekspor). Sudah terbukti 2x.
- File temp PowerShell: selalu di `OUTPUT\_tmp_*.ps1`, tulis hasil ke file lalu baca, HAPUS setelah selesai.
  Shell = PowerShell (bukan cmd). Hindari em-dash & `\"` di skrip (merusak parser).
- File .md di OUTPUT ber-line-ending LF → grep_search tool sering gagal; pakai skrip PowerShell untuk cari.
- Register keputusan: K-001..K-056. K berikutnya = **K-057**.
- Skill `to-tickets`/`to-spec`/`to-questionnaire` = `disable-model-invocation:true`: Claude siapkan
  bahan → BERHENTI → work owner jalankan manual. `/to-tickets` interaktif (quiz user) + format default
  beda dari tiket kita (lihat catatan Fac Out di atas).
