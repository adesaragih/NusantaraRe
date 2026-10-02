# F11: Jalur Fac Out pada siklus endorsement

**What to build:** Endorsement yang mengubah porsi retrosesi menghasilkan **selisih** Fac Out yang
tertulis berdampingan dengan nilai sesudahnya, memakai jalur produksi Fac Out yang sudah hidup (F10).

## Populasi kerja — **55, bukan 39**

⛔ **Populasi berbasis nama TIDAK cukup.** `[terverifikasi]`:

| Ukuran | NB | RNW | EDM |
| --- | ---: | ---: | ---: |
| Activity **bernama** `FacOut`/`FacRetro` (Ordinal, peka huruf) | 41 | 40 | 39 |
| Hadir di ketiga folder | — | **39** | — |
| Activity membawa penanda `10015` **tanpa** nama Fac Out | 16 | 15 | 16 |
| **Permukaan fungsional** | ≈57 | **≈55** | **≈55** |

Yang tanpa nama Fac Out termasuk `CopyToAllSpreading_ACT`, `CopyToAllLocSpreading_ACT`,
`cekSpreadingFactIn`, `CountPremiumNet`, `SaveJsonPolicyFacIn_Act`, `SetReinsurerEndorsement_Act`.

✅ **K-059:** keduanya **diperlakukan sama**, karena **`10015` adalah kode Fac Out**. Pembawa kode itu
**bukan** kelompok kelas dua yang "perlu ditinjau" — ia **bagian penuh** dari populasi kerja.

⛔ Angka lama **18/18 dan 10/18 DIBATALKAN** — populasinya tidak pernah terdefinisi.

📌 **Sumber angka:** `..\..\10-audit\03-struktur-facretro-dan-populasi.md` §3 (populasi + banding 23
tag, lengkap dengan uji silang wajib) dan `..\..\10-audit\07-hitung-berkas-dan-drift.md` §1 (jumlah
berkas per subfolder) serta §3 (drift). Keduanya memuat perintah audit yang dapat dijalankan ulang.

## Sembilan activity yang berbeda di Endorsement

`[terverifikasi]` Dari **39** yang hadir di ketiga folder, kontrak 23 tag (K-042/K-043, termasuk
amandemen "isi blok `pzIndexes` diabaikan"): **NB vs RNW 39/39 identik**; **NB vs EDM 30/39 identik,
9 berbeda**:

`CopyAllObjFacOutFireAneka_ACT` · `CopyFacRetroAnekaGolf_ACT` · `InsertFacoutProd` ·
`InsertFacoutProduction` · `InsertFacoutProductionEDM` · `InsertFacoutProductionEDMCurr` ·
`OfferFacOut_PreAct` · `SetDataFacOutAnekaGolf_Act` · `SetDataFacOutFire_Act`

📌 **NB dan RNW berbagi satu implementasi** — nol perbedaan. Hanya Endorsement yang butuh cabang.

## Isi khusus endorsement

`[terverifikasi]` `Activity\InsertFacoutProd` adalah pengirim: `Call InsertFacoutProduction` dan
`Call InsertFacoutProductionEDM` (bergerbang `IsEDM`), lalu `Obj-Refresh-And-Lock`, lalu
`Call UpdatePolisAddendum_Act`. Gerbang prefiks `@contains(pyWorkPage.pyWorkIDPrefix,"EDM-")` ada di
`InsertFacoutProduction`.

`[terverifikasi]` Kolom selisih di `FACOUTPRODUCTION`: `SHAREOFFERED_SELISIH`, `OBJECTPREMI_SELISIH`,
`COMMISION_SELISIH`, `PREMI_COVERAGE_MENJADI`/`_SELISIH`, `COMMISION_COVERAGE_MENJADI`/`_SELISIH`,
`RATE_COVERAGE`.

⚠️ `[terverifikasi]` **Ketidaksetangkupan pasangan `_MENJADI` vs `_SELISIH`** — tidak setiap kolom
selisih punya pasangan "menjadi". Sudah tercatat sebagai butir K-046 (nomor 7 pada daftar 12 kode
usang). **Diport apa adanya.**

⚠️ **K-046** `K046_GuardFacOut_HanyaType7` — gerbang idempotensi fac out hanya dipakai keluar pada
`Type=="7"`. Untuk `Type` lain, penulisan ganda **tidak dicegah**. `[terverifikasi]`
`InsertFacoutProductionEDM`. Diport apa adanya; risikonya tercatat sebagai E-Q32.

**Asal (Pega).** `Activity\InsertFacoutProductionEDM` · `InsertFacoutProductionEDMCurr` ·
`CopyFacRetroEDMLoc_ACT` · `ProtectionEDMFacout_Act` · `DataTransform\CountPremiEDMFacOut_DT`
(`Rule-Obj-Model`) · sembilan activity berbeda di atas

**Keputusan.** **K-059** (populasi ≈55 — nama Fac Out **dan** pembawa `10015` diperlakukan sama,
karena `10015` **adalah** kode Fac Out) · K-056 · **K-046** ·
K-042/K-043 (kontrak 23 tag) · `CLAUDE.md` §4.3, §4.6

**Blocked by:** F10 · `..\edm\E21-jalur-produksi-edm.md`
⛔ `E21` sendiri **blocked** menunggu tabel flat (P-10) — lihat indeks EDM.

**Status:** blocked

- [ ] Populasi kerja **≈55 per folder**, bukan 39 — pembawa `10015` **diperlakukan sama** dengan yang bernama Fac Out (K-059)
- [ ] **Sembilan** activity yang berbeda di Endorsement punya cabang tersendiri; **tiga puluh** sisanya dipakai ulang
- [ ] ⛔ **NB dan RNW berbagi satu implementasi** — tidak ada cabang siklus di antara keduanya
- [ ] Kolom selisih terisi berpasangan dengan nilai sesudahnya, **sesuai ketaksetangkupan aslinya**
- [ ] **K-046** `K046_GuardFacOut_HanyaType7` — gerbang idempotensi hanya aktif pada satu jenis penyesuaian
- [ ] **K-046** `K046_MenjadiSelisih_TidakSetangkup`
- [ ] Nilai berkoma desimal (K-027); uang bertipe `Money`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
