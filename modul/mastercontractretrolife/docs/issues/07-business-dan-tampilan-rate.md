# 07: Business — cakupan treaty, rate reasuransi, dan tampilan tabel rate

**Status:** ready-for-agent

**Blocked by:** 02 (business lahir di bawah kontrak), 04 (jenis reasuransi diwarisi dari kontrak)

## Hasil & nilai pengguna

Sebagai **admin master**, saya menambahkan **jenis business** yang tercakup sebuah kontrak beserta
rate reasuransinya, dan saya dapat **melihat tabel rate** lebih dulu supaya memilih yang benar —
dengan rate tersimpan **persis seperti saya menuliskannya**. *(User story 21–23 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas business; **`RIRATE` sebagai teks** |
| `internal/repository` | Pemanggilan procedure penulis business; pembacaan tabel rate |
| `internal/services` | Wajib-isi; pewarisan jenis reasuransi dari kontrak |
| `internal/handlers` | Endpoint CRUD business; endpoint lihat tabel rate |
| `frontend/` | Grid business; dua layar tampilan rate |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyBusiness_Life_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatyBusiness_Life_SQL.xml` | `POOLDATA.INSERTBUSINESS_LIFE` |
| `SaveBusinessLife_Act` | `ASM-FW-GISFW-…` / `SAVEBUSINESSLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveBusinessLife_Act.xml` | orkestrator simpan |
| `NewInputBusinessLife_Act` | `ASM-FW-GISFW-…` / `NEWINPUTBUSINESSLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/NewInputBusinessLife_Act.xml` | baris baru |
| `SetBusinessListLife_Act` | `ASM-FW-GISFW-…` / `SETBUSINESSLISTLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetBusinessListLife_Act.xml` | isi daftar |
| `SetParamRate` / `SetParamRateTable` | `ASM-FW-GISFW-…` / `SETPARAMRATE`, `SETPARAMRATETABLE` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/` | menyiapkan tampilan rate |
| `ViewRate` / `ViewRateTable` | `@BASECLASS` / `VIEWRATE`, `VIEWRATETABLE` / `RULE-OBJ-FLOW-ACTION` | `Master Contract Retro Life/FlowAction/` | **tampilan saja** |
| `BrowseRateLife_RD` | `ASM-FW-GISFW-INT-M_RATE_LIFE` / `BROWSERATELIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseRateLife_RD.xml` | tabel rate |
| `BrowseRateLifeSummary` | `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` / `BROWSERATELIFESUMMARY` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseRateLifeSummary.xml` | ringkasan rate |
| `BrowseBusinessLife_RD` | `ASM-FW-GISFW-INT-BUSINESS` / `BROWSEBUSINESSLIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseBusinessLife_RD.xml` | daftar business |
| `InboxBusinessLifeReinsurers` | `DATA-PORTAL` / `INBOXBUSINESSLIFEREINSURERS` / `RULE-HTML-HARNESS` | `Master Contract Retro Life/Harness/InboxBusinessLifeReinsurers.xml` (486.911 byte) | layar grid |

`[terverifikasi]` `SetParamRate` menyalin `InputBusinessLife.RIRATEID` dan `RIRATE` ke `ParamID.*`;
`SetParamRateTable` menyalin dari `Param.*`. Keduanya menyiapkan **popup tampilan**, tidak mengubah
apa pun.

⚠️ `[data DBA]` **`RIRATE` bertipe `VARCHAR2(1000)` — TEKS**, bukan `NUMBER`. Rate reasuransi
disimpan sebagai string; kemungkinan memuat format bertingkat. **Jangan paksa menjadi angka.**

`[data DBA]` `TREATYBUSINESS_LIFE` adalah **satu-satunya tabel yang sudah punya PK sejak semula**
(`TREATYBUSINESS_LIFE_PK`, unique index, `ENABLE VALIDATE`).

## ADR terkait

**ADR-0003** (uang non-float — **tidak berlaku pada `RIRATE`** yang memang teks),
**ADR-0006**, **ADR-0007**.

## Acceptance criteria

- [ ] Business lahir **di bawah** satu kontrak, dengan kode business, nama, dan rate reasuransi.
      *(AC 25 spec)*
- [ ] ⚠️ **`RIRATE` tersimpan sebagai teks apa adanya** — **tidak** dikonversi, **tidak** diformat
      ulang, **tidak** dipaksa menjadi angka. Nilai yang ditulis dan dibaca kembali **identik
      karakter demi karakter**. *(AC 26 spec; `[data DBA]` kolom `VARCHAR2(1000)`)*
- [ ] Tabel rate dapat **dilihat** sebelum memilih, **tanpa mengubah apa pun** — popup murni baca.
      *(AC 27 spec)*
- [ ] Ringkasan rate juga dapat dilihat, dari sumber yang terpisah dari tabel rate.
- [ ] `REINSTYPEID` dan `TREATYYEAR` **tidak ditulis** pada baris business — diwarisi dari kontrak
      dan tahun. *(AC 41–42 spec; tiket 04)*
- [ ] Menyimpan business yang sudah ada = **upsert**, bukan baris kedua; identitas baru dibuat basis
      data. *(AC 45, 4 spec)*
- [ ] `o_message` diperiksa; kegagalan ditampilkan. *(AC 46 spec — HTML dan konformansi di
      **tiket 10**)*
- [ ] Daftar business yang dapat dipilih dibaca dari master business, bukan ditanam sebagai
      konstanta.

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Jangan ubah tipe `RIRATE`.** `[data DBA]` Ia `VARCHAR2(1000)`. Mengubahnya menjadi angka
menuntut **konfirmasi format terlebih dahulu** dari Product+UW — di luar cakupan tiket ini, dan
tercatat di §Out of Scope spec.

⚠️ **OQ-066 — jangan buang perilaku karena memo.** `[data DBA]` `SetOutputParam_DT` (`@BASECLASS` /
`SETOUTPUTPARAM_DT`) bermemo **`not used`** **tetapi masih dirujuk tiga Harness dan dua Section**.
Bila tiket ini menyentuh Harness business, **telusuri rujukannya** sebelum menyimpulkan ia mati.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 30-09-2026 — sesi implementasi (paket 0)

> Sumber: `RALAT-DEV-30-09-2026.md` (K1–K8 katalog DEV, R1–R12 pembacaan ulang XML) dan `PARITAS-LAYAR-DAN-AKSI.md`. Kalimat di atas **tidak dihapus**; yang berlaku adalah ralat ini.

| Kalimat lama | Ralat |
| --- | --- |
| TAMBAHAN-TIKET: *"⛔ **Tiket ini MACET** sampai **Pertanyaan A** dijawab"* | **terjawab dari korpus** (R7): `RIRATE` = nama tabel rate (`USEDBY` dari `BrowseRateLifeSummary`), `RIRATEID` = ID tabel rate; tetap teks. Tiket tidak macet |
| *"`REINSTYPEID` dan `TREATYYEAR` **tidak ditulis** pada baris business"* | tetap **ditulis** sebagai salinan induk (K4) |
| — | `View Rate` form (tampil bila `RIRATEID` terisi) dan `View Rate` baris sama-sama membuka section `ViewRate` (`Rate List`, view `RATE_LIFE` disaring `IDUSEDBY`); objek ringkasan rate `[dugaan]` — OQ-MCRL-05 |

## Status 01-10-2026 (paket 11)

**Status:** ⏸️ **sebagian** — business simpan/ubah/hapus paket 6/7 (`b3e097d`, `81b78bd`) dan panel `Business List` paket 9+10 (`a3c07bc`) dibangun; autocomplete `R/I RATE` dan `Rate List` dibangun tetapi datanya **menunggu OQ-MCRL-13** (503 berkalimat) — business BARU belum dapat disimpan (`RIRATEID` wajib).

## Status 01-10-2026 — K1 keputusan work owner 01-10-2026 (OQ-MCRL-13 + OQ-MCRL-05)

**Status:** ✅ **selesai** — kalimat lama *"autocomplete `R/I RATE` dan `Rate List` dibangun tetapi datanya **menunggu OQ-MCRL-13**
(503 berkalimat) — business BARU belum dapat disimpan (`RIRATEID` wajib)"* tidak berlaku lagi.

- `R/I RATE` membaca view `RATE_LIFE_SUMMARY` (`BrowseRateLifeSummary`: `ID`, `USEDBY`, urut `ID ASC`); `Rate List` membaca *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)*
  view `RATE_LIFE` (`BrowseRateLife_RD`: enam kolom grid, `IDUSEDBY = :1`, urut `ID DESC, RATE ASC`, 500 baris + `terpotong`).
- Business **baru** dapat disimpan; RIRATEID pilihan baru wajib ada di `RATE_LIFE_SUMMARY` (penyimpangan sadar — Pega tidak *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)*
  memeriksa); `RIRATE` tetap teks apa adanya (AC 26).
- AC 27 *"popup murni baca"* kini dijaga dua lapis: `periksaBacaSaja` (runtime) dan `TestMCRLMasterHanyaDibacaSelect` (statik).
- DEV baca-saja 01-10-2026: `GET /ringkasan-rate` 200 (100 saran), `GET /rate` 200 (1 dan 59 baris), nol tulisan.

