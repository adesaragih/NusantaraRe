# Indeks Tiket — Lima Modul Terverifikasi (16 tiket)

> ## 🗺️ Peta indeks — empat modul, **60 tiket**
>
> Berkas ini adalah indeks **New Business saja**. Tiga modul lain punya indeksnya sendiri:
>
> | Modul | Indeks | Tiket | Status |
> | --- | --- | ---: | --- |
> | **New Business** | **berkas ini** | **16** | tertulis, belum dikerjakan |
> | Renewal | [`rnw\00-INDEKS-RNW.md`](rnw/00-INDEKS-RNW.md) | **8** | tertulis, belum dikerjakan |
> | Endorsement | [`edm\00-INDEKS-EDM.md`](edm/00-INDEKS-EDM.md) | **22** | tertulis, seluruhnya blocked |
> | Fac Out (retrosesi keluar) | [`facout\00-INDEKS-FACOUT.md`](facout/00-INDEKS-FACOUT.md) | **14** | tertulis, seluruhnya blocked |
> | | | **60** | `[terverifikasi]` **nol berkas `.go`** — nol yang dapat dimulai hari ini |
>
> Urutan kerja: **NB → RNW → EDM → Fac Out**. Yang membuka seluruhnya adalah tiket
> `01-tracer-money-ratio-premi-pa.md` di bawah ini.
>
> ⚠️ Catatan `_HANDOFF-EDM.md` yang menyebut "54 tiket" dan "Fac Out F01..F08" sudah **dibatalkan**
> di berkas itu — angka yang berlaku adalah **60** dan **F01…F14**.
>
> *Blok ini ditambahkan 21 September 2026. Isi indeks New Business di bawahnya tidak diubah.*

> **Indeks pengaman.** Publikasi per-tiket lengkap menyusul lewat **`/to-tickets` manual**
> (`disable-model-invocation: true`). File ini **arsip breakdown, bukan tiket final**.
> **Nomor baris presisi menyusul di tiket lengkap.**

**Status breakdown:** disetujui work owner 17 September 2026 · **16 tiket, urut dependency**
(blocker lebih dulu).

---

## Daftar tiket

| # | Slug | Blocked by | Yang dihasilkan |
| ---: | --- | --- | --- |
| 01 | `tracer-money-ratio-premi-pa` | — | `pkg/money` (decimal, mata uang wajib, keadaan `Unknown`, parser koma **K-027**) + `pkg/ratio` (skala melekat pada nilai, `Money × Ratio`) + resolver berisi PA + **Seam 1 hidup** + uji **gagal-kompilasi** `Money + Ratio`. Kriteria: **rekonsiliasi eksak** satu fixture PA |
| 02 | `mata-uang-unknown` | 01 | Tabel perilaku `Unknown` penuh; **aritmetika lintas mata uang → `panic`**; "`panic` saat tulis Oracle" dicatat sebagai **kewajiban tertunda** (seam repository belum ada) — ADR-0006 / K-012 |
| 03 | `resolver-cob-skala` | 01 | **8 lini bisnis (COB) terkunci**; aturan **pembagi komposit** (`100000 = 1.000 × 100`); **COB tak dikenal → `panic`** — K-018 / ADR-0004 |
| 04 | `rumus-premi-pa` | 03 | **Keempat bentuk rumus PA**; presisi pembulatan **literal per langkah** + komentar asal rule Pega — K-011 |
| 05 | `rumus-premi-mbu` | 03 | **Bentuk bersarang MBU** — keterikatan faktor→pembagi terbaca dari strukturnya sendiri |
| 06 | `rumus-premi-layering` | 03 | Rumus premi **Layering** (rule `GenerateLayerList_ACT`) |
| 07 | `pembulatan-dalam-loop-akumulasi` | 04, 05, 06 | Kasus uji **lintas-COB** untuk galat yang **menumpuk per iterasi** — wajib ADR-0005 |
| 08 | `registry-rules-eval` | — | **Satu seam untuk 601 rule `When`**; baca **kedua** tag kondisi; nama tak dikenal → `panic`; table-driven |
| 09 | `predikat-sikap-khusus` | 08 | `IsPKSASM` **tetap `panic`** · `IsOfferFacIn` pakai ekspresi tersimpan (**K-002**) · `IsSpreadingDepan` catat asal salinan (**K-003**) · `IsFacout` diport apa adanya (**K-019**) |
| 10 | `jembatan-predikat-cob` | 08, 03 | Titik integrasi **Seam 3 → resolver skala** (`IsFire` / `IsPA` / `IsMBU` mengisi COB) |
| 11 | `tangga-akseptasi-bentuk-a` | 01 | **Satu keputusan = satu transisi**; tiga field state terpisah; **tangga selesai = penyelesaian normal**; fixture dari data nyata **tanpa kolom `NAMA`/`LOGIN`** (**K-025**) |
| 12 | `akseptasi-bentuk-b-financial` | 11 | Query **per jenis pertanggungan**; ejaan jabatan **berspasi**; **tanpa eskalasi & tanpa antrean**; **`WHERE` dipertahankan**; dua mekanisme **tidak disatukan** — K-026 |
| 13 | `domain-hasil-keputusan-flag-fase` | 11 | Domain **{1, 2, 3, 4, 7, 9}** (**K-021**) · `Reject` ≠ `Decline` (**K-014**) · **satu flag fase** `panic`/`decline` (**K-008** / ADR-0002) |
| 14 | `special-acceptance-putaran-kedua` | 11, 13 | Cabang **1SA/2SA** lewat `RequestType` langkah RDB + `GetLimitAkseptasi_Act2` — **cabang PULIH** via amandemen **K-006**. ⚠️ Catatan: 1SA/2SA adalah **rule integrasi/SQL, bukan Activity** |
| 15 | `de-identifikasi-berkas-kasus` | — | Identitas dibuang, **seluruh angka dipertahankan**. Kriteria penerimaan tersendiri — lihat blok di bawah |
| 16 | `rekonsiliasi-eksak-tahap-1` | **15**, 04, 05, 06, 11 | Pembanding **nol-selisih** atas berkas kasus. ⛔ **Tidak pernah menyentuh berkas mentah ber-PII** — ADR-0001 |

---

## Kriteria penerimaan tiket 15 — **terkunci**

⛔ **Daftar kunci EKSPLISIT, ditinjau satu per satu. Penyaringan berbasis pola/substring DILARANG** —
`TotalSumInsured` dan `CedingRetention` akan ikut terbuang dan **merusak angka rekonsiliasi**.

**BUANG (12 kunci):**

`InsuredName` · `InsuredID` · `ASMAddress` · `SelectedLocationAddress` · `RoadName` · `ASMZipCode` ·
`RiskZipCode` · `Email` · `pxCreateOpName` · `MarketingName` · `CedingCoName` · `CoinsName`

**PERTAHANKAN:**

Seluruh field bernilai **angka** — termasuk `TotalSumInsured`, `CedingRetention`, `ShareCeding`,
`ShareOfCeding`.

**Bila ada field identitas lain saat implementasi:** tambahkan lewat **tinjauan manual**.
**JANGAN beralih ke filter pola.**

**Lulus bila:** nol nama, nol alamat, dan **setiap angka identik dengan sumbernya**.

---

## Ditandai menunggu — **bukan tiket**

| Butir | Menunggu |
| --- | --- |
| `services/spreading` | Spec + seam tersendiri — yang sendirinya menunggu penjelasan gerbang parameter dari **IT + Underwriting** |
| Jalur simpan produksi · `repository` tulis · `services/production` · `models` tulis | Struktur **tabel flat** (butir E) + **`ALL_SOURCE` 32 prosedur** |
| Rujukan menggantung `SetErrorMessage_Act` vs `CountASMGrossPremi_ACT` / `CountASMNetPremi_ACT` | **Pertanyaan terbuka** — perlu keputusan, bukan tiket |

---

## Sumber breakdown

`04-spec\03-spec-modul-terverifikasi.md` (5 modul + 3 seam) · register keputusan
**K-001…K-029** · **ADR-0001…0006**.

Kosakata mengikuti `steering\GLOSARIUM.md`. Tiga seam yang disetujui —
`premium.Calculate` · `acceptance.Next` · `rules.Eval` — **tidak boleh ditambah** tanpa keputusan baru.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
