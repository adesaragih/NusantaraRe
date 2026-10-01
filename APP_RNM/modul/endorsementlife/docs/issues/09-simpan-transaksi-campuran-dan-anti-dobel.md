# 09: Simpan endorsement — batas transaksi campuran dan penjaga anti-dobel

> **Ralat 01-10-2026** (gelombang 2 brief, `../RALAT-DEV-01-10-2026.md` — ralat mengalahkan isi di bawah). Teks lama yang tidak berlaku:
> - **E2** — *"Nama tabel lama (`M_LIFE_PREMIUM_DETAIL` / `_SUMMARY`) adalah sistem lama — target tulis sistem baru adalah tujuh tabel PremiumList Life"* → keduanya **juga** ditulis di transaksi yang sama (kontrak Claim Life): rekap `M_LIFE_PREMIUM_SUMMARY` ditulis; peserta `M_LIFE_PREMIUM_DETAIL` menunggu OQ-EDM-016 (**R29**).
> - **R17/R18** — `LIFEINPRODUCTION` dan `JSON_POLIS` tidak ditulis (OQ-EDM-010; AC 54).
> - Anti-dobel `(NOPOLIS, PRODKE)` → `T_PREMIUM_LIST (NO_POLIS, PROD_KE)` index `IDX_PL_NOPOLIS_PRODKE`, ditambah pemeriksaan versi berjalan: kasus yang salinannya basi (versi lebih baru resmi sesudah kasus dibuat) ditolak 409.

**Status:** sebagian 01-10-2026 — c83bf68; peserta warisan `M_LIFE_PREMIUM_DETAIL` menunggu OQ-EDM-016

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 04 (nomor endorsement), 06 (nilai baris, termasuk yang negatif)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin endorsement yang disetujui tersimpan lengkap — rekam polis,
produksi, detail peserta, dan ringkasan premium — dan bila gagal di tengah, keadaannya
**terdeteksi dan dapat diulang tanpa menggandakan apa pun**. *(User story 43 di spec; §11 spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Endorsement + baris detail versi baru (`T_PREMIUM_LIST_DETAIL`) beserta `PARENT_ID`; rekap mata uang |
| `internal/repository` | Pemanggilan keempat penulis; **penjaga anti-dobel `(NOPOLIS, PRODKE)`** |
| `internal/services` | Orkestrasi urutan; deteksi keadaan separuh; pengulangan aman |
| `internal/handlers` | Status "tersimpan" versus "tertunda" terbaca API |
| `frontend/` | Penanda endorsement tersimpan / tertunda |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Commit sendiri? |
| --- | --- | --- | --- |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` (444.634 byte, **16 langkah**) | orkestrator |
| `InsertJsonPolisEDM` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/InsertJsonPolisEDM.xml` | **ya** — `COMMIT;` baris 107 |
| `SaveLifeinProduction_SQL` | `ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SaveLifeinProduction_SQL.xml` | **ya** — `COMMIT;` baris 164 |
| `SaveMasterLPDet` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SaveMasterLPDet.xml` | **ya** — `COMMIT;` baris 252 |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/InsertPLSummary.xml` | **ya** — `COMMIT;` baris 125 |

`[terverifikasi]` **Rantai rujukan `<RequestType>` di `InsertJsonPolisLife_Act`:**
`GetProdKeOldData_SQL` (911) → `Generate_NoEndorsmentLife` (1304) → **`InsertJsonPolisEDM` (2543)**
→ **`SaveLifeinProduction_SQL` (2720)** → **`SaveMasterLPDet` (5154, step 11.6)** →
**`InsertPLSummary` (6226, step 12.2)** → `GetNopolisByIDPega` (6444).

`[terverifikasi]` Dua langkah **REMARK**: step **13** "Cek sudah masuk atau blm datanya" (6399) dan
step **15** `SendEmailNotification` (6752) — ⚠️ **keduanya justru DIHIDUPKAN** di sistem baru, lihat
tiket **10**.

`[terverifikasi]` Ditambah **dua `Commit` aktif** yang berjalan **sebelum** pengguna menekan simpan:
`CreateCaseEMDL` step 4 dan `MappingEDMLife` step 14 (tiket 02).

⚠️ `[terverifikasi]` **`InsertJsonPolisEDM` adalah `INSERT` polos, bukan upsert:**
`INSERT INTO POOLDATA.JSON_POLIS (IDPEGA, DATA_JSON, TGL_INPUT, NOPOLIS, NOENDORS, PRODKE, TGL_PROD,
USERNAME) VALUES (…)`. Jalur new business tidak punya masalah ini karena memakai procedure **upsert**
`INSERTJSONPOLISLIFE` berkunci `IDPEGA` (**PL-05b**).

## ADR terkait

**ADR-0015** (Go memegang batas transaksi; efek yang tidak boleh hilang ditangani eksplisit),
**ADR-0003** (uang non-float di kolom tabel relasional maupun JSON API), **ADR-0011**.

## Acceptance criteria

- [ ] Setiap penulis yang **commit sendiri** diperlakukan sebagai **titik potong**; data yang harus
      atomik **tidak dipisahkan** olehnya. *(AC 39 spec; **ADR-0015**)*
- [ ] Kegagalan di tengah meninggalkan keadaan yang **terdeteksi dan dapat dipulihkan** — bukan
      senyap. *(AC 40 spec)*
- [ ] ⚠️ **Penyimpangan sadar — penjaga anti-dobel.** Menjalankan penyalinan dua kali untuk
      **`(NOPOLIS, PRODKE)`** yang sama **tidak** menggandakan versi — dibuktikan dengan
      memanggilnya dua kali lalu **menghitung baris versi di `T_PREMIUM_LIST`** untuk kombinasi itu.
      ⚠️ **Bukan** dengan menghitung baris `JSON_POLIS`: JSON **tidak lagi ditulis** (§16).
      *(AC 41, 71 spec; `[keputusan work owner]`)*
- [ ] Urutan pemanggilan terdokumentasi di kode **sebagai bagian kebenaran**, bukan kebetulan; ada
      test yang **gagal bila urutannya diubah**. *(AC 42 spec)*
- [ ] Commit dini di pembuatan case dan pemetaan (tiket 02) **dipertahankan** — gerbang "satu EDM
      terbuka per polis" membutuhkannya. `[keputusan work owner]`
- [ ] Baris detail **termasuk yang bernilai negatif** tertulis utuh ke tabel peserta
      (`T_PREMIUM_LIST_DETAIL`); rekap mata uang **dihitung ulang** dari peserta versi baru.
      ⚠️ Nama tabel lama (`M_LIFE_PREMIUM_DETAIL` / `_SUMMARY`) adalah **sistem lama** — target
      tulis sistem baru adalah **tujuh tabel PremiumList Life** (§16). *(AC 18, 63 spec)*
- [ ] Nilai uang di **kolom tabel** maupun di **kontrak API** ditulis sebagai **desimal presisi
      arbitrer** — tidak lewat `float`, tidak ada pembulatan diam. (**ADR-0003**)
- [ ] `PL_NUMBER_EDM` tertulis pada rekam summary **di samping** `PL_NUMBER`, sebagai **dua nilai
      terpisah**.
- [ ] Kegagalan setelah salah satu commit meninggalkan bagian yang sudah ter-commit **utuh**, dan

### Satu transaksi ⚠️ BARU 2026-09-16 — spec §11, §16

- [ ] ⚠️ Satu endorsement ditulis dalam **satu transaksi** — penyalinan versi, kolom EDM pada header,
      seluruh peserta beserta `PARENT_ID`, spreading & retro, dan rekap mata uang yang dihitung
      ulang; kegagalan di mana pun **membatalkan seluruhnya**. *(AC 70 spec)*
- [ ] ⚠️ **Tidak ada procedure JSON yang dipanggil.** Test yang menemukan pemanggilan
      `InsertJsonPolisEDM` atau padanan perakit payload **gagal**. *(AC 54 spec; penyimpangan
      sadar 8)*
- [ ] ⚠️ **Penjaga anti-dobel tetap berlaku, dengan alasan baru**: menjalankan penyalinan dua kali
      untuk maksud endorsement yang sama menghasilkan **dua versi**. Pemeriksaan `(NOPOLIS, PRODKE)`
      sebelum menulis menahannya. *(AC 71 spec; spec §11)*
- [ ] Dua `Commit` dini (`CreateCaseEMDL`, `MappingEDMLife`) tetap **di luar** transaksi simpan —
      gerbang "satu endorsement terbuka per polis" membutuhkannya. *(spec §11)*
      pengulangan **aman**.

## Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` `Commit` eksplisit dan `Connect-REST` yang ter-remark di rantai simpan **tidak
direplikasi** — **Go memegang transaksi** (**ADR-0015**).

⚠️ **OQ-066 berlaku.** Berkas ini justru buktinya: dua langkah ber-`//` (13 dan 15) **bukan** kode
mati yang dibuang, melainkan jalur yang **sengaja dihidupkan kembali** (tiket 10). **Jangan**
menyimpulkan hidup/mati dari `<pyStepsBlockName>` saja.

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
