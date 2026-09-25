# Indeks Tiket — Siklus Endorsement (EDM)

> **22 tiket, lingkup 708 berkas.** Sumber: enam spec di `..\..\04-spec\05`–`10-spec-edm-*.md` ·
> register **K-043…K-052** · banding premi `..\..\07-edm\10`/`11-banding-*.md` · ADR-0001…0006.
>
> ⛔ **16 tiket NB (`..\01`…`..\16`) dan 8 tiket RNW (`..\rnw\R01`…`R08`) tidak disentuh.**

**Lingkup:** **708 berkas** (K-043) = 354 berbeda + 354 EDM-only. Dari EDM-only, **350 diport**
(K-051 — 4 berkas pilih-tertanggung **dibuang**, K-044).

---

## ⛔ 0. URUTAN KERJA SEBENARNYA: NB → RNW → EDM

> **Tidak satu pun tiket EDM dapat dimulai sekarang.**

`[terverifikasi]` **Nol berkas `.go` ada di repositori.** Tidak ada `go.mod`, `cmd/`, `internal/`,
`pkg/`, maupun `services/`. Ketujuh belas tiket NB dan sembilan tiket RNW **sudah tertulis tetapi
belum dikerjakan**.

| Tahap | Status | Isi |
| --- | --- | --- |
| **1. New Business** | ⛔ **belum dikerjakan** | 16 tiket — fondasi: `pkg/money`, `pkg/ratio`, Seam 1, Seam 2, Seam 3, rekonsiliasi |
| **2. Renewal** | ⛔ belum dikerjakan | 8 tiket — memakai ulang seluruh modul NB |
| **3. Endorsement** | ⛔ **menunggu tahap 1** | 22 tiket di berkas ini |

⛔ **Berkas ini adalah PETA RENCANA, bukan antrean kerja yang siap diambil.** Seluruh tiket EDM
ter-block oleh tiket NB yang belum ada implementasinya. Menampilkannya seolah dapat dimulai akan
menyesatkan siapa pun yang membacanya.

**Fondasi NB yang harus mendarat lebih dulu** — disebut dengan **nama berkas tiketnya**, bukan kode:

| Berkas tiket NB | Memberi apa kepada EDM |
| --- | --- |
| `..\01-tracer-money-ratio-premi-pa.md` | `pkg/money` + `pkg/ratio`, uji gagal-kompilasi, resolver PA, pintu perhitungan premi |
| `..\03-resolver-cob-skala.md` | Resolver skala per lini bisnis (K-018) — pembagi komposit |
| `..\04-rumus-premi-pa.md` · `..\05-rumus-premi-mbu.md` · `..\06-rumus-premi-layering.md` | Ketiga rumus premi yang Seam 3 perluas |
| `..\08-registry-rules-eval.md` | Seam 1 penuh — registry 601 rule `When` |
| `..\11-tangga-akseptasi-bentuk-a.md` | Seam 2 — **tidak dipakai EDM Life**, tetapi dipakai lini umum |
| `..\15-de-identifikasi-berkas-kasus.md` | Fixture ter-de-identifikasi |
| `..\16-rekonsiliasi-eksak-tahap-1.md` | Kerangka pembanding nol-selisih |

---

## 1. Enam seam — tidak bertambah

| # | Seam | K | Dihidupkan tiket |
| ---: | --- | --- | --- |
| 1 | `internal/rules.Eval` — **diperluas**, sumber data per-rule | K-050 | **E01** |
| 2 | `services/acceptance.Next` | (NB) | — ⛔ **Life melewatinya** (E18) |
| 3 | `services/premium.Calculate` — **satu pintu, 4 bentuk** | K-051 | **E12** |
| 4 | `services/endorsement.PrepareBeforeImage` | K-047 | **E06** |
| 5 | `services/endorsement.OpenCase` | K-049 | **E03** |
| 6 | `services/underwriting.ScoreMedical` | K-052 | **E20** |

⛔ **Seam 3 adalah SATU pintu dengan empat bentuk** — dasar · dua-bagian endorsement · tabel tarif
Travel · dekomposisi delta coverage. **Bukan** seam per bentuk.

---

## 2. Urutan dependency

| # | Slug | Blocked by | Seam | Yang dihasilkan |
| ---: | --- | --- | :---: | --- |
| **E01** | `registry-predikat-edm-sumber-data` | `..\08-registry-rules-eval` | **1** | Registry satu, **sumber data per-rule**. **37 cabang** COB, bukan 203 |
| **E02** | `predikat-perilaku-beda-edm` | E01 | 1 | `IsUW` · `IsClaim` · `IsTravel` · **`IsNotEDM` bukan negasi** |
| **E03** | `tracer-open-case-edm` | `..\08-registry-rules-eval` · E01 | **5** | **Seam 5 hidup.** Kasus `EDM-`, tautan tiga arah, langsung fase Policy |
| **E04** | `gerbang-penolakan-enam` | E03 | 5 | 6 gerbang + **4 klep** + bypass RI Slip |
| **E05** | `validasi-tanggal-edm` | E03 | 5 | Tanggal dalam periode; daftar pengecualian jadi konfigurasi |
| **E06** | `tracer-before-image-fire` | `..\01-tracer-money-ratio-premi-pa` · E01 | **4** | **Seam 4 hidup.** Lapis A+B+C FIRE + prorata → fixture rekonsiliasi eksak |
| **E07** | `lapis-a-sembilan-penyimpangan` | E06 | 4 | 6 cermin + 3 non-`OldData` |
| **E08** | `lapis-b-enam-lini` | E06 | 4 | Golf · Aneka · PA · MC · MBU · Travel + cedant |
| **E09** | `guard-k046-dan-a5` | E08 | 4 | `RateOld` 6/6 · `PremiumOld` FIRE · **A.5** angka `0` |
| **E10** | `lapis-c-tujuh-varian` | E06 | 4 | 7 varian · `IsProRate` FIRE · `FlagOldData` LIFE |
| **E11** | `gerbang-keluar-before-image` | E08 | 4 | Keluar: Life · non-EDM · >100 lokasi |
| **E12** | `prefactor-seam3-empat-bentuk` | `..\01-tracer-money-ratio-premi-pa` · `..\03-resolver-cob-skala` · `..\04-rumus-premi-pa` · `..\05-rumus-premi-mbu` · `..\06-rumus-premi-layering` | **3** | **Prefactor.** Satu pintu 4 bentuk. `ProRatePercent` **pecahan, dipakai langsung** |
| **E13** | `delta-pembayaran-mata-uang` | E12 · E06 | 3 | `SumTotalPayment = EDMPremiMenjadi − EDMOldPayment` |
| **E14** | `tabel-edmpremimenjadi` | E13 | 3 | **22 penugasan · 11 rumus** per (lini × `EdmType` × jalur) |
| **E15** | `delta-baris-spreading` | E12 · E06 | 3 | `(baru × ProrateEDMEnd) − lama_TreatyType_sama` + 9 pasangan kolom |
| **E16** | `pendukung-layar-endorsement` | E01 · E03 | — | **350 diport**; 4 pilih-tertanggung **tidak** |
| **E17** | `repository-query-edm` | ⛔ **seam repository** | — | 22 Connect-SQL + 14 DataPage + 28 RD |
| **E18** | `jalur-life-terpisah` | E06 · E08 · E10 | 4 | Flow **`InputEDMLife`** · **tanpa tangga akseptasi** |
| **E19** | `premi-life` | E12 · E18 | 3 | `TSILiability × RateLifeAverage × (1+Rate/100) ÷ 1.000` |
| **E20** | `seam6-skoring-medis` | E18 | **6** | **Seam 6 hidup.** `Standard`/`Postpone`/`Decline` |
| **E21** | `jalur-produksi-edm` | ⛔ **tabel flat (P-10)** | — | `facinproduction` · `FACOUTPRODUCTION` · versi polis |
| **E22** | `rekonsiliasi-edm` | `..\15-de-identifikasi-berkas-kasus` · `..\16-rekonsiliasi-eksak-tahap-1` · E09 · E14 · E15 · E19 | — | Nol selisih memakai kerangka NB |

### Frontier — setelah fondasi NB mendarat

**E01** dan **E12** terbuka lebih dulu. Setelah **E01** hijau → **E02**, **E03**, **E06**.
Setelah **E06** hijau → **E07**, **E08**, **E10** sekaligus.

⛔ **Sebelum tahap NB selesai, frontier-nya KOSONG.**

---

## 3. Daftar BLOCKED oleh pihak luar

| Tiket | Menunggu | Catatan |
| --- | --- | --- |
| **E21** | ⛔ **tabel flat (P-10)** | `POOLDATA.INSERTJSONPOLIS` **diganti insert-ke-flat**, dirancang bersama NB/RNW |
| **E17** | ⛔ **seam repository** | Belum ada; kewajiban tertunda sejak spec NB |

⚠️ **Keduanya tetap di daftar.** Menghapusnya menyembunyikan lingkup yang sudah diketahui.

---

## 4. Belum tuntas sebelum IMPLEMENTASI — catatan tiket, bukan blocker

| Butir | Muncul di |
| --- | --- |
| Rumus PA & Travel penuh — `ViewPremiTravel_SQL` nihil; tarif jadi masukan repository | E12 · E17 |
| `CalculatePremiPA_FacIn` identik NB↔EDM (23-tag lulus) — tinggal dipakai ulang | E12 |
| 4 bentuk Seam 3: satu fungsi bercabang atau 4 implementasi — keputusan desain | E12 |
| 5 dari 6 salinan mesin premi berbeda isi — dibedah sebelum dipakai ulang | E12 · E13 |

---

## 5. Peta seam → tiket

| Seam | Dihidupkan | Dipakai lagi |
| --- | --- | --- |
| **1** `rules.Eval` | **E01** | E02 · E03 · E06 · E16 |
| **2** `acceptance.Next` | (`..\11-tangga-akseptasi-bentuk-a`) | ⛔ **Life melewatinya** (E18) |
| **3** `premium.Calculate` | **E12** | E13 · E14 · E15 · E19 |
| **4** `PrepareBeforeImage` | **E06** | E07 · E08 · E09 · E10 · E11 · E18 |
| **5** `OpenCase` | **E03** | E04 · E05 |
| **6** `ScoreMedical` | **E20** | — |

---

## 6. Yang tidak dijalankan

⛔ **`/to-tickets` belum dijalankan** (`disable-model-invocation: true`). Berkas tiket E01–E22 di
folder ini adalah **bahan siap-publikasi**; penerbitan resminya lewat `/to-tickets` manual oleh work
owner.

⚠️ **§4.3 arsip** (*"nilai dasar akseptasi = selisih TSI"*) tetap **`[belum diuji]`** — **E18** dan
**E22** tidak bersandar padanya.

---

*Tanpa nama orang, tanpa alamat email, tanpa nomor polis, tanpa data pelanggan, tanpa data medis.*
