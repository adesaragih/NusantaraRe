# 06: Daftar plan dan batas underwriting

**Status:** selesai (01-10-2026) — paket 6 `05ca173`, layar paket 10 (`1ada8d7`); `View Rate` / `Choose R/I Rate` ⏸️ OQ-MPNL-03; uji `db` ditulis dan MELEWATI di sesi implementasi (tanpa `ORACLE_DSN`; POOLDATA/DEV bukan sasaran)

**Blocked by:** 03 (baris anak ikut dalam transaksi atomik produk)

## Hasil & nilai pengguna

Sebagai **underwriter**, saya mencatat **beberapa plan** dalam satu produk — masing-masing dengan
benefit dan RI Rate-nya — dan **beberapa batas underwriting** dengan status medis serta rentang usia
dan uang pertanggungan, sehingga satu produk dapat menawarkan beberapa paket dengan syarat
bertingkat. *(User story 18–19, 23 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Baris plan; baris batas underwriting |
| `internal/repository` | Tulis/baca kedua tabel anak **di dalam transaksi produk** |
| `internal/services` | Aturan per baris; keterkaitan RI Rate pada plan |
| `internal/handlers` | Endpoint kedua daftar |
| `frontend/` | Dua grid di layar produk |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ProteksiPlanListLife` | `ASM-FW-GISFW-…` / `PROTEKSIPLANLISTLIFE` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/ProteksiPlanListLife.xml` | daftar plan |
| `CopyUnderWritingLimit` | `ASM-FW-GISFW-…` / `COPYUNDERWRITINGLIMIT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/CopyUnderWritingLimit.xml` | salin batas underwriting |
| `BrowseUnderwritingList` | `ASM-FW-GISFW-…` / `BROWSEUNDERWRITINGLIST` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/BrowseUnderwritingList.xml` | baca daftar |
| `ChooseRIRate` | `ASM-FW-GISFW-DATA-PLAN` / `CHOOSERIRATE` / `RULE-OBJ-FLOW-ACTION` | `Master Product Name Life/FlowAction/ChooseRIRate.xml` | pemilih RI Rate (tiket 04) |
| `CountMaxReasured_Act`, `CountMaxSumReasured_Act` | `ASM-FW-GISFW-…` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/` | hitungan terkait batas |

`[data DBA]` Isi baris, dari contoh data nyata:

| Tabel anak | Kolom baris |
| --- | --- |
| `product_life_plan` | plan, id plan, nama, **benefit**, **RI Rate + id-nya** |
| `product_life_uw_limit` | deskripsi, **status medis** (`FCL` / `NM` / `MEDIS`), usia min/maks, uang pertanggungan min/maks |

⚠️ `[fakta bisnis — work owner]` **RI Rate melekat pada baris plan**, bukan pada produk induk —
berbeda dari **RI Risk** yang melekat pada induk.

## ADR terkait

**ADR-0003** (uang non-float — batas uang pertanggungan pada baris underwriting),
**ADR-0007** (jejak audit), **ADR-0015** (kegagalan ditangani eksplisit).

## Acceptance criteria

- [ ] Produk dapat memuat **beberapa plan**, masing-masing dengan benefit dan RI Rate-nya.
      *(AC 22 spec)*
- [ ] Produk dapat memuat **beberapa batas underwriting** dengan deskripsi, status medis, rentang
      usia, dan rentang uang pertanggungan. *(AC 23 spec)*
- [ ] Status medis hanya menerima **`FCL`**, **`NM`**, atau **`MEDIS`**; nilai lain **ditolak**.
      *(`[data DBA]` dari contoh nyata)*
- [ ] Baris dapat **ditambah dan dihapus** tanpa menyentuh produk induk — kecuali ketika penyimpanan
      dilakukan sebagai satu kesatuan. *(AC 26 spec)*
- [ ] Daftar yang **kosong** tersimpan sebagai daftar kosong, **bukan** kegagalan. *(AC 27 spec)*
- [ ] ⚠️ Kegagalan pada baris plan atau batas underwriting **membatalkan seluruh simpan** (baris
      `product_life` + seluruh anak) — baris anak ikut dalam transaksi produk.
      *(AC 7 spec; `[keputusan work owner]`)*
- [ ] Batas uang pertanggungan pada baris underwriting diperlakukan sebagai **desimal presisi
      arbitrer**; **tidak** melewati `float`. *(AC 28 spec; **ADR-0003**)*
- [ ] Batas bawah yang **lebih besar** dari batas atas — usia maupun uang pertanggungan — **ditolak**.
      *(AC 32 spec)*
- [ ] ⚠️ **RI Rate pada baris plan** berasal dari master RI Rate, **bukan** dari master RI Risk.
      *(AC 19 spec; `[fakta bisnis — work owner]`)*
- [ ] Urutan baris dalam daftar **dipertahankan** saat dibaca kembali.

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Dua tabel anak paling berisi.** `[data DBA]` Dari ketujuh daftar bersarang, hanya `PlanList` dan
`UnderwritingLimitList` yang **terisi** pada contoh produk nyata — karena itu bentuk kolomnya
terbaca. Lima sisanya ditangani tiket **07**.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — keikutsertaan baris anak dalam transaksi
produk hanya berperilaku benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 01-10-2026 — sesi implementasi (paket 6)

> Sumber: `../RALAT-DEV-30-09-2026.md` (P1, R16, R17, R18) dan `../PARITAS-LAYAR-DAN-AKSI.md` §5. Kalimat di atas **tidak dihapus**.

| Kalimat lama | Ralat |
| --- | --- |
| *"Tulis/baca kedua tabel anak **di dalam transaksi produk**"* | **P1**: bukan tabel anak — larik `PlanList` dan `UnderwritingLimitList` di `M_PRODUCT_LIFE.JSONDATA`, ditulis bersama produk di satu transaksi; urutan baris dipertahankan, daftar kosong = `[]` |
| *"Status medis hanya menerima **`FCL`**, **`NM`**, atau **`MEDIS`**; nilai lain **ditolak**"* | **R16**: tidak ditegakkan — `.Medical` b44817 teks bebas; tiga nilai itu tebakan dari satu contoh (OQ-MPNL-12) |
| *"**RI Rate pada baris plan** berasal dari master RI Rate"* | **R18**: sumbernya view atas JSON rate — `Choose R/I Rate` b34548 dan `View Rate` b34067 menjawab 503 berkalimat sampai OQ-MPNL-03; pasangan `RIRATE`/`RIRATEID` yang sudah tersimpan diterima, yang **baru** ditolak 422 menyebut OQ-MPNL-03. ⚠️ Akibatnya baris plan baru belum dapat disimpan (`RI/RATE tidak boleh kosong`) |
| Rule sumber `ProteksiPlanListLife` *"daftar plan"* | gerbang VERBATIM per baris: `Plan tidak boleh kosong` (b442), `Plan tidak boleh sama` (b421), `RI/RATE tidak boleh kosong` (b463), diawali `PLAN LIST row N:`. Pega menjalankannya pada onChange `.Plan` b33163 dan setiap `Delete` baris; di sini saat simpan (pesan di halaman memblokir submit Pega dengan efek yang sama) |
| (tidak disebut) autocomplete `Plan Name` b33124 | `GET /api/master-product-name-life/master-plan?cari=` — RD `BrowseProductTypeLife_RD` (kelas `PRODUCT_TYPE_LIFE`, OQ-MPNL-04): `Plan ← CoverName`, `PlanID ← ID`, `Name ← Business`, `Benefit ← Benefit`; plan yang berubah diverifikasi dan namanya dari master |
| *"Batas bawah yang **lebih besar** dari batas atas … **ditolak**"* | ditegakkan per baris `UNDERWRITING LIMIT`: `Min Insured` ≤ `Max Insured`, `Min Age` ≤ `Max Age` (R17) |
| (tidak disebut) ikon salin b45290 `CopyUnderWritingLimit`, tombol `Add` b43598 / `Delete` b45633 | aksi baris di layar (paket 10): salin = baris baru berisi keenam medan baris itu (`CopyUnderWritingLimit` 1 b224) |
