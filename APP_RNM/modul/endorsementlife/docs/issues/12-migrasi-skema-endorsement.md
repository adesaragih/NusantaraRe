# 12: Migrasi skema — tabel premium Life dan kolom penanda `EDMSTATUS`

> **Ralat 01-10-2026** (gelombang 2 brief, `../RALAT-DEV-01-10-2026.md` — ralat mengalahkan isi di bawah). Teks lama yang tidak berlaku:
> - **E3** — *"**Status:** wontfix — **digantikan tiket 00 + PremiumList Life tiket 00**"* → **di luar lingkup gelombang 2, tetap `needs-info`** — keputusan memindah `JSON_POLIS` belum ada (OQ-EDM-001, bawaan tidak dipindah).

**Status:** wontfix — **digantikan tiket 00 + PremiumList Life tiket 00**

**Blocked by:** —

> ⛔ **JANGAN KERJAKAN TIKET INI.** `[keputusan work owner]` Premisnya — merancang skema target
> sebagai salinan tabel premium existing, menunggu DDL fisik — **dibatalkan 2026-09-16**.
>
> **Ke mana isinya pindah:**
>
> | Yang dulu di sini | Sekarang |
> | --- | --- |
> | DDL tabel premium target | **PremiumList Life tiket `00`** — tujuh tabel dirancang sendiri |
> | Kolom penanda `EDMSTATUS` + kolom EDM lain | **tiket `00` konteks ini** — `ALTER` kolom nullable |
> | Pemindahan data endorsement lama | **PremiumList Life tiket `00`** — endorsement adalah **versi** polis, jadi ia pindah bersama polisnya (AC 55, 56 spec) |
> | `JSON_POLIS` / CLOB | **dibuang** — tidak dibawa sama sekali (AC 54 spec) |
>
> ⚠️ **OQ-001 ditutup** — ia yang dulu menahan tiket ini. Termasuk pertanyaan **"apakah kolom uang
> menerima nilai negatif"**: terjawab oleh rancangan sendiri — kolom uang **desimal bertanda**, dan
> jurnal balik memang menulis minus.
>
> Berkas dipertahankan sebagai jejak keputusan, **bukan** sebagai pekerjaan.

---

<details>
<summary>Isi asli (sudah tidak berlaku)</summary>

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin skema Oracle sistem baru menyimpan setiap nilai endorsement
dengan presisi yang sama persis — **termasuk nilai negatif hasil jurnal balik** — dan menyimpan
penanda status peserta dengan bentuk yang sama, sehingga rekonsiliasi tidak menemukan selisih dan
penyaringan peserta hidup berperilaku identik sebelum dan sesudah pindah.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL tabel target + index |
| `internal/repository` | Pemetaan tipe kolom ↔ desimal presisi arbitrer; bentuk penyaring `EDMSTATUS` |
| — | Skrip rekonsiliasi |

## Rule Pega sumber

| Tabel | Penulis yang terbukti | Bukti |
| --- | --- | --- |
| `POOLDATA.M_LIFE_PREMIUM_DETAIL` | `SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL`), PK `M_LIFE_PREMIUM_DETAIL_SEQ.nextval` — **satu rule, dua jalur** | `Endorsement Life/RDBList/SaveMasterLPDet.xml` **dan** `PremiumList Life/RDBList/SaveMasterLPDet.xml` |
| `POOLDATA.M_LIFE_PREMIUM_SUMMARY` | `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` lewat `InsertPLSummary` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY`) | kedua modul |
| `POOLDATA.JSON_POLIS` | `InsertJsonPolisEDM` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM`) — **`INSERT` polos**; jalur NB memakai procedure upsert `INSERTJSONPOLISLIFE` | `Endorsement Life/RDBList/InsertJsonPolisEDM.xml` |
| `POOLDATA.JSON_OFFER_LIFE` | `POOLDATA.INSERTJSONOFFERLIFE` lewat `SaveOfferJsonLife_SQL` | `PremiumList Life/RDBList/SaveOfferJsonLife_SQL.xml` |
| `POOLDATA.LIFEINPRODUCTION` | `SaveLifeinProduction_SQL` (`ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL`) | `Endorsement Life/RDBList/SaveLifeinProduction_SQL.xml` |
| `ARASAPAS.DETAIL_INVOICE` | **dibaca saja**, lintas skema | `Endorsement Life/RDBList/SearcStatusBayarArasaps_SQL.xml` |

`[terverifikasi]` **Kolom `M_LIFE_PREMIUM_DETAIL` terbaca lengkap** dari daftar `INSERT`
`SaveMasterLPDet` — termasuk `PL_NUMBER`, `PL_NUMBER_EDM`, `CERTIFICATE_NO`, `NAME_OF_INSURED`,
`CURRENCY`, `RATE`, `PRORATETYPE`, `SUM_AT_RISK_GROSS`, `SUM_AT_RISK_RETRO`, `RETROCEDED_SHARE`,
seluruh kelompok uang gross / `*_REFUND` / `*_RETRO` / `*_REFUND_RETRO`, **`IDPEGA`**,
**`EDMSTATUS`**, **`STATUSOLD`**, **`STATUS`**, `EM_PERCENT`, `RISK`.

`[data DBA]` Kolom `M_LIFE_PREMIUM_SUMMARY` (**37 kolom**), `JSON_POLIS`, dan `JSON_OFFER_LIFE`
sudah diketahui dari body procedure yang diserahkan 2026-09-15.

## ADR terkait

**ADR-0003** (uang non-float; DDL `NUMBER` tanpa presisi → desimal presisi arbitrer di aplikasi),
**ADR-0009** (migrasi penuh; koeksistensi ditolak), **ADR-0011**, **ADR-0010** (lampiran tetap di
Google Storage — **bukan** bagian migrasi tabel ini).

## Acceptance criteria

*(belum dapat difinalkan — menunggu OQ-001; disusun agar siap dijalankan begitu DDL turun)*

- [ ] DDL tabel target menyalin **tipe, presisi, PK, index, dan nullability** dari tabel sumber apa
      adanya; tidak ada kolom uang yang menjadi `FLOAT`/`BINARY_DOUBLE`. (**ADR-0003**)
- [ ] ⚠️ **Kolom uang menerima nilai NEGATIF tanpa kehilangan presisi** — jurnal balik endorsement
      menulis 32 kolom bertanda minus. Rekonsiliasi membandingkan nilai lama dan baru **secara
      tepat**, bukan dengan toleransi.
- [ ] ⚠️ **Tipe dan nullability `EDMSTATUS` ditetapkan** — dan bentuk penyaring "peserta hidup"
      disesuaikan: `IS NULL` atau `= ''` untuk baris new business. Test penyaring dijalankan ulang
      terhadap skema hasil migrasi.
- [ ] Kolom `CLOB` (`DATA_JSON`, `JSONDATA`) pindah utuh, termasuk isi yang panjang.
- [ ] Sequence (`M_LIFE_PREMIUM_DETAIL_SEQ`, `M_LIFE_PREMIUM_SUMMARY_SEQ`, `JSON_OFFER_SEQ`)
      dipindahkan dengan **nilai berjalan yang benar**, sehingga nomor pasca-migrasi tidak pernah
      bertabrakan dengan nomor lama.
- [ ] Index yang menopang kueri hilir ada sejak hari pertama — khususnya
      `M_LIFE_PREMIUM_DETAIL(PL_NUMBER)` yang dipakai Claim Life, dan kolom `EDMSTATUS` bila
      penyaringan menuntutnya.
- [ ] `PRODKE` pindah utuh, dan penomoran endorsement pasca-migrasi melanjutkan urutan yang benar.
- [ ] Migrasi dapat dijalankan ulang dengan aman dan punya jalur mundur yang diuji.
- [ ] Skema uji yang dipakai seluruh tiket lain dibangun **dari DDL yang sama** — bukan dari tiruan
      yang ditulis terpisah.

## Blocker

🚧 **`needs-info` — OQ-001 (sisa) terbuka.** Pemilik: **DBA**.

**Yang sudah ada:** nama kolom keempat tabel — `M_LIFE_PREMIUM_DETAIL` terbaca lengkap dari korpus;
`M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS`, `JSON_OFFER_LIFE` dari body procedure `[data DBA]`.

**Yang belum ada:**

1. **DDL fisik** — tipe, presisi, PK, index, nullability. ⚠️ Seluruh parameter
   `PEGA_M_LIFE_PREMIUM_SUMMARY` bertipe `VARCHAR2` **termasuk kolom uang**, sehingga tipe kolom
   sebenarnya belum diketahui.
2. **Tipe dan nullability `EDMSTATUS`** — apakah baris new business menyimpan `NULL` atau string
   kosong `''`. Salah pilih akan **membuang seluruh peserta new business** dari layar klaim.

**JANGAN paksa `ready`.** Pola sama dengan **CL-13** (Claim Life) dan **PL-09** (PremiumList Life).

**Tidak memblokir tiket lain** — kolomnya sudah cukup untuk seluruh tiket 01–11.

## Perintah verifikasi

```
go test ./internal/...
```

</details>
