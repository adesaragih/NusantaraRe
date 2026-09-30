# 06: Daftar plan dan batas underwriting

**Status:** ready-for-agent

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
