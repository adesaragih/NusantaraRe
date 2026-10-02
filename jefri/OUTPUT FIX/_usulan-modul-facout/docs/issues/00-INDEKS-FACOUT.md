# Indeks Tiket — Fac Out (Retrosesi Keluar)

> **Empat belas tiket.** Sumber: `..\..\09-facout\01-temuan-dan-rancangan-facout.md` · register
> **K-053…K-062** (seluruhnya **sudah diputuskan**, ada di `..\..\00-KEPUTUSAN-WORK-OWNER.md`;
> bahan draf diarsipkan di `..\..\10-audit\05-draf-keputusan-facout.md`) ·
> audit `..\..\10-audit\04-perujuk-isfacout-peka-huruf.md`.
>
> ⛔ **46 tiket siklus Fac In tidak disentuh** — NB (`..\01`…`..\16`), RNW (`..\rnw\R01`…`R08`),
> EDM (`..\edm\E01`…`E22`).

---

## §0 ⛔ Seluruh 14 tiket berstatus **blocked**

`[terverifikasi]` **Nol berkas `.go` ada** di repo — tidak ada `go.mod`, `cmd/`, `internal/`, `pkg/`.
Tiket ini adalah **peta rencana**, bukan antrean kerja yang dapat diambil hari ini.

Urutan kerja tetap **NB → RNW → EDM → Fac Out**. Yang membuka seluruhnya adalah tiket New Business
pertama, `..\01-tracer-money-ratio-premi-pa.md`.

---

## Daftar tiket

| # | Berkas | Blocked by | Yang dihasilkan |
| ---: | --- | --- | --- |
| **F01** | `F01-registry-predikat-facout.md` | `..\08-registry-rules-eval.md` · `..\09-predikat-sikap-khusus.md` | `IsFacRetro` (ekspresi tersimpan) + `IsInputFacRetro` + **tiga bentuk penanda spreading** yang tidak digabung. ⛔ `IsFacout` **tidak** di sini |
| **F02** | `F02-tracer-derivasi-facretro-fire.md` | F01 · `..\01-tracer-money-ratio-premi-pa.md` · `..\11-tangga-akseptasi-bentuk-a.md` | Tracer lini Fire: UW Accept → derivasi → flag `.IsFacRetro` menyala → alur membelok ke menu Fac Out |
| **F03** | `F03-pengecualian-okupasi.md` | F02 | Daftar **26 alternatif** okupasi (20 awalan + 6 kesetaraan), diport apa adanya |
| **F04** | `F04-derivasi-tiga-cob-sisanya.md` | F02 | Derivasi untuk Aneka/Golf, Marine Cargo/MBU, PA/Travel + penjumlah PA |
| **F05** | `F05-empat-struktur-staging.md` | F02 | **Empat** struktur penampung + 4 properti `FacRetroDetails` + 6 koleksi anak `FacRetroList` |
| **F06** | `F06-premi-retro.md` | F05 · `..\01-tracer-money-ratio-premi-pa.md` · `..\03-resolver-cob-skala.md` | Premi retro, komisi RI, diskon. **Urutan terkunci: diskon sebelum komisi** |
| **F07** | `F07-validasi-facout.md` | F05 | Validasi daftar objek + reasuradur + ambang 50 lokasi |
| **F08** | `F08-tangga-retro-empat-tingkat.md` | F05 | Tangga retro 4 tingkat + percabangan grup. **Bukan Seam 2** |
| **F09** | `F09-dua-generator-nomor-retro.md` | F08 | **Dua** generator nomor, sequence berbeda, tidak diseragamkan |
| **F10** | `F10-produksi-facoutproduction.md` | F06 · F08 · F09 | Insert 65 kolom ke `FACOUTPRODUCTION` + gerbang idempotensi + parameter terikat |
| **F11** | `F11-jalur-facout-endorsement.md` | F10 · `..\edm\E21-jalur-produksi-edm.md` | Selisih Fac Out endorsement. Populasi **55**, 9 activity berbeda |
| **F12** | `F12-pendukung-layar-facout.md` | F01 · F05 | Menu + daftar Fac Out. ⛔ `GetOPFacOut_Sql` **tidak dapat diport** |
| **F13** | `F13-cetak-rislip-dan-email.md` | F08 · ⛔ **eksternal** | Cetak R/I Slip + kirim. Menunggu isi `M_LINK_SERVICE` |
| **F14** | `F14-rekonsiliasi-facout.md` | F10 · `..\15-de-identifikasi-berkas-kasus.md` · `..\16-rekonsiliasi-eksak-tahap-1.md` | Nol selisih; 9 penanda K-046 wajib nol |

---

## Graf ketergantungan

```
F01 ─┬─ F02 ─┬─ F03
     │       ├─ F04
     │       └─ F05 ─┬─ F06 ──┐
     │               ├─ F07   │
     │               └─ F08 ─┬┴─ F10 ─┬─ F11  (+ edm\E21)
     │                       │ F09 ───┘       └─ F14
     │                       └─ F13
     └─ F12  (+ F05)
```

**Akar: F01.** Asiklik — setiap blocker bernomor lebih kecil, nol rujukan maju.
**Daun penutup: F11, F13, F14.**

**Frontier setelah blocker NB selesai:** F01. Setelah F02 hijau, tiga tiket terbuka sekaligus
(F03, F04, F05).

---

## Ketergantungan ke tiket siklus lain

| Tiket Fac Out | Tiket luar | Mengapa |
| --- | --- | --- |
| F01 | `..\08-registry-rules-eval.md` | Predikat Fac Out masuk registry yang sama (K-050), bukan registry kedua |
| F01 | `..\09-predikat-sikap-khusus.md` | `IsFacout` ditangani di sana (K-019); `IsFacRetro` memakai pola `IsOfferFacIn` (K-002) |
| F02 | `..\01-tracer-money-ratio-premi-pa.md` · `..\11-tangga-akseptasi-bentuk-a.md` | Derivasi dipicu dari jalur UW Accept dan menyalin nilai uang |
| F06 | `..\01-tracer-money-ratio-premi-pa.md` · `..\03-resolver-cob-skala.md` | Pembagi 100.000 diturunkan resolver K-018 |
| F11 | `..\edm\E21-jalur-produksi-edm.md` | Jalur produksi endorsement; **E21 sendiri blocked** menunggu tabel flat (P-10) |
| F14 | `..\15-de-identifikasi-berkas-kasus.md` · `..\16-rekonsiliasi-eksak-tahap-1.md` | Fac Out tidak punya pembanding sendiri |

---

## ⛔ Blocker di luar kedua set tiket

| Butir | Menunggu | Kena |
| --- | --- | --- |
| Isi `M_LINK_SERVICE` (endpoint) | DBA | **F13** |
| ~~Pengganti `GetOPFacOut_Sql`~~ | ⛔ **BUKAN blocker lagi** — **K-062**: klep **dihapus**, kolom pelaksana kosong | — |
| Tabel flat (P-10) lewat `edm\E21` | work owner | **F11** |

---

## Keputusan yang tercermin di tiket

| Keputusan | Tercermin di |
| --- | --- |
| **K-053** menu Fac Out terpisah namun terhubung | F02 · F12 |
| **K-054** tiga bentuk penanda spreading tidak digabung | F01 · F04 |
| **K-055** tangga retro tersendiri 4 tingkat, bukan Seam 2 | F08 |
| **K-056** produksi ke `FACOUTPRODUCTION` existing, insert = Connect-SQL | F10 |
| **K-019** fitur usang, kode tetap diport | F01 · F14 |
| **K-046** kode usang diport apa adanya — **9 penanda** | F03 · F05 · F06 · F09 · F11 · F14 |
| **K-048** §8.1 satu properti ditimpa berurutan | F06 |
| **K-018** pembagi komposit | F06 |
| **K-025** identitas tidak pernah disalin | F08 · F10 · F12 · F13 · F14 |

---

## ✅ K-057 … K-062 — SUDAH DIPUTUSKAN 21 September 2026

Teks yang berlaku ada di `..\..\00-KEPUTUSAN-WORK-OWNER.md`. **Nomor berikutnya = K-063.**

| Keputusan | Isi | Kena |
| --- | --- | --- |
| **K-057** | Pemicu Fac Out = flag `.IsFacRetro`; rumusan `IsFacout` **dicabut** | F01 · F02 |
| **K-058** | Sumber kebenaran = salinan per siklus, **kecuali** rule yang sudah diekspor ulang ke `DDL\` → `DDL\` **menang**. Tangga 4 tingkat berlaku **menyeluruh**; **layar keempat peran sama** | F08 · F11 |
| **K-059** | Populasi ≈55 — nama Fac Out **dan** pembawa `10015` **diperlakukan sama** | F11 |
| **K-060** | Premi retro tetap `CountRateRetroCov`; salinan berlaku = **`DDL\CountRateRetroCov.xml`** | F06 |
| **K-061** | F13 ditulis penuh, ditandai menunggu, pola R07 | F13 |
| **K-062** | `GetOPFacOut_Sql` **dihapus**; kolom pelaksana kosong. ⛔ **pengecualian pertama terhadap K-006** | F12 |

### ⚠️ Satu pertanyaan yang TIDAK ditutup

**Apakah metode deteksi drift diganti** dari banding `pyRuleSetVersion` menjadi banding **hash 23
tag** — **TETAP TERBUKA** (K-058). `[terverifikasi]` metode lama melewatkan **267 dari 354**
perbedaan NB↔EDM.

### ⛔ Dua konsekuensi yang mengubah tiket

1. **K-060 mencabut keanehan `Local.RIComIN`.** Salinan korpus **basi**; pada `DDL\` variabel itu
   **diakumulasi** di `pySteps(8)`. Penanda `K046_RIComIN_TidakPernahDisetel` **dicabut** dari F06
   dan F14 — sisa **delapan** penanda K-046, bukan sembilan.
2. **K-062 menghapus satu blocker eksternal.** F12 tidak lagi menunggu apa pun untuk klep itu.

📌 Daftar rule yang layak diekspor ulang lebih dulu: `..\..\10-audit\08-daftar-ekspor-ulang-prioritas.md`.

---

## Di luar lingkup tiket Fac Out

| Butir | Alasan |
| --- | --- |
| `IsFacout` | Ditangani tiket NB `..\09-predikat-sikap-khusus.md` di bawah K-019 — `[terverifikasi]` **tidak terpanggil** di jalur Fac Out |
| Mesin premi Fac In, registry, tangga akseptasi, spreading | Diwarisi NB; sudah punya tiket 01–16 |
| Rancangan tabel flat | Butir tertunda lintas siklus (P-10) |
| Isi riwayat Pega `DATAPEGA.PC_*` | `CLAUDE.md` §4.3 — hilang bersama Pega |

---

## ⚠️ Catatan privasi — nama kolom **bukan** nilai

`DOB`, `LICENSEPLATE`, `ENGINENUMBER` dan `ADDRESS` muncul di **F10**, **F12** dan **F14**. Seluruhnya
adalah **nama kolom tabel `FACOUTPRODUCTION`**, bukan nilai data pelanggan.

⛔ **Menghapusnya justru merusak F14**, yang tugasnya **mende-identifikasi kolom-kolom itu**: sebuah
tiket tidak dapat memerintahkan perlindungan sebuah kolom tanpa menyebut namanya.

📌 **Pemindaian privasi berikutnya WAJIB membedakan nama kolom dari nilai.** Yang dilarang K-025
adalah **nilainya** — tanggal lahir sungguhan, plat nomor sungguhan, nomor mesin sungguhan, alamat
sungguhan. Pemindai yang hanya mencari string `DOB` akan menandai F10/F12/F14 sebagai pelanggaran
padahal keduanya justru mekanisme perlindungannya.

Pola pemindaian yang benar: cari **bentuk nilai** (mis. tanggal `dd/mm/yyyy`, pola plat nomor), bukan
**nama kolom**.

---

## Status ketiga siklus + Fac Out

| Siklus | Tiket | Status |
| --- | ---: | --- |
| New Business | 16 | tertulis, **belum dikerjakan** |
| Renewal | 8 | tertulis, **belum dikerjakan** |
| Endorsement | 22 | tertulis, seluruhnya blocked |
| **Fac Out** | **14** | **tertulis, seluruhnya blocked** |
| **TOTAL** | **60** | **nol yang dapat dimulai hari ini** |

📌 Catatan `_HANDOFF-EDM.md` baris 20–21 dan 46 yang menyebut "tiket F01..F08 sudah ada" dan
"54 tiket" **tidak benar** — folder ini kosong sebelum penerbitan ini. Total sebelumnya **46**,
sekarang **60**.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
