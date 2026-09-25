# 05b: ~~Simpan JSON polis dan rekam produksi~~ — ⛔ **DIBATALKAN 2026-09-16**

**Status:** wontfix — **digantikan tiket 00 + 05a**

**Blocked by:** —

> ⛔ **JANGAN KERJAKAN TIKET INI.** `[keputusan work owner]` **Seluruh JSON dibuang** (spec §12).
> Tiket ini seluruhnya tentang menulis `JSON_POLIS` dan `JSON_OFFER_LIFE` **sesudah** commit, dan
> tentang menangani titik potong yang ditimbulkan keduanya. Kedua procedure itu **tidak lagi
> dipanggil**, sehingga **seluruh alasan keberadaan tiket ini hilang**.
>
> **Ke mana isinya pindah:**
>
> | Yang dulu di sini | Sekarang |
> | --- | --- |
> | Tulis `JSON_POLIS` / `JSON_OFFER_LIFE` | **dibuang** — tidak ada padanannya (AC 32, 34 spec) |
> | Urutan "sesudah commit" & titik potong | **lenyap** — satu polis = **satu transaksi** (tiket 05a, AC 35 spec) |
> | Penulisan peserta ke `M_LIFE_PREMIUM_DETAIL` | ke **`T_PREMIUM_LIST_DETAIL`**, di dalam transaksi yang sama (tiket 03 + 05a) |
> | `SaveLifeinProduction_SQL` sebagai titik potong | `LIFEINPRODUCTION` **tidak ditulis** lagi; hanya dibaca saat migrasi (tiket 00, AC 33 spec) |
> | Deteksi keadaan separuh & pengulangan | **tidak perlu** — atomisitas menggantikannya (AC 35 spec) |
> | Efek keluar Arasapas | tetap di tiket **06** |
>
> Berkas dipertahankan sebagai jejak keputusan, **bukan** sebagai pekerjaan.

---

<details>
<summary>Isi asli (sudah tidak berlaku)</summary>

**Blocked by (asli):** 05a (transaksi summary harus commit lebih dulu — itulah inti tiket ini)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin representasi JSON polis dan rekam produksi ditulis **setelah**
premium list aman tersimpan, dan saya ingin penulisan itu **dapat diulang tanpa menggandakan data**
bila gagal di tengah — supaya kegagalan jaringan atau basis data tidak pernah meninggalkan premium
list yang benar dengan polis yang tidak pernah terbit. *(User story 33–36 di spec)*

## Area codebase

`internal/repository` (pemanggilan procedure yang commit sendiri, di luar transaksi Go; deteksi
keadaan separuh), `internal/services` (perakitan payload JSON polis; orkestrasi urutan + pengulangan),
`internal/handlers` (status "tersimpan tetapi polis tertunda" terbaca API).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` (327.332 byte, tersimpan `20260728T024422`, **17 langkah**) | orkestrator |
| `InsertJsonPolis` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLIS` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertJsonPolis.xml` | `POOLDATA.INSERTJSONPOLISLIFE(p_IDPEGA, p_NOPOLIS, P_JSONDATA, P_TGL_PROD, P_USERNAME, {OutputData.HASIL1 out})` + `COMMIT;` (102) |
| `SaveLifeinProduction_SQL` | `ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/SaveLifeinProduction_SQL.xml` | `INSERT INTO POOLDATA.LIFEINPRODUCTION` + `COMMIT;` (164) |
| `GetNopolisByIDPega` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetNopolisByIDPega.xml` | baca balik untuk verifikasi |
| `GetJsonProductLife` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetJsonProductLife.xml` | data produk untuk payload |

`[terverifikasi]` Rantai rujukan `<RequestType>` di `InsertJsonPolisLife_Act`: `GetJsonProductLife`
(1444) → `InsertPLSummary` (3112, step **8**) → `InsertJsonPolis` (4094, step **10**) →
`SaveLifeinProduction_SQL` (4271, step **11**) → `GetNopolisByIDPega` (4449, step **12**).
**`SaveOfferJsonLife_SQL` tidak ada di rantai ini** — lihat tiket 01 dan **OQ-067**.

`[data DBA]` `POOLDATA.INSERTJSONPOLISLIFE` — **upsert** ke `POOLDATA.JSON_POLIS` berkunci `IDPEGA`
(`UPDATE` bila ada, `INSERT` bila belum). Kolom: `IDPEGA`, `DATA_JSON` (**CLOB**), `TGL_INPUT`,
`NOPOLIS`, `PRODKE`, `TGL_PROD`, `USERNAME`. Gagal → jatuh ke `JSON_POLIS_ERROR`.
**COMMIT di dalam procedure.**

`[terverifikasi]` Step **12** "Cek sudah masuk atau blm datanya" membaca balik lewat
`GetNopolisByIDPega`; step **13** menyalakan penanda bila `OutDataLife.pxResults(1).PL_NUMBER==""`
(baris 5105). Inilah **deteksi keadaan separuh** milik Pega — ia dipertahankan sebagai konsep, dan
menjadi pemicu alarm di tiket **06**.

## ⚠️ Urutan yang mengikat `[keputusan desain]`

1. Transaksi Go (penomoran + summary) **sudah commit** — tiket **05a**.
2. **Baru** `INSERTJSONPOLISLIFE` dipanggil, **di luar** transaksi Go, karena ia commit sendiri.
3. `SaveLifeinProduction_SQL` juga commit sendiri (baris 164) — ia **titik potong ketiga**.
4. **Penulisan detail peserta NB ke `M_LIFE_PREMIUM_DETAIL` termasuk di sini** — dipanggil **INLINE**
   sebagai bagian alur simpan polis, lewat logika `InsertLifePremiumDetail_act`
   (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY`) →
   `SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`), yang
   **commit sendiri** (`COMMIT;` baris 252) — **titik potong keempat**. Rincian dan AC-nya di tiket
   **08**. ⚠️ `[keputusan work owner]` **Tanpa job** — di Pega, NB dipicu batch; di sistem baru
   inline, meniru EDM.

**Jangan** menempatkan data yang harus atomik pada dua sisi procedure yang commit sendiri.

## ADR terkait

**ADR-0015** (batas transaksi dipegang Go; efek yang tidak boleh hilang ditangani eksplisit),
**ADR-0003** (uang non-float di dalam payload JSON), **ADR-0013** (endpoint di-resolve runtime —
berlaku pada efek keluar di tiket 06), **ADR-0011**.

## Acceptance criteria

- [ ] `INSERTJSONPOLISLIFE` dipanggil **setelah** transaksi summary commit, tidak pernah di dalamnya.
      *(AC 21 spec)*
- [ ] Kegagalan pada `INSERTJSONPOLISLIFE` meninggalkan nomor + rekam summary **utuh**; API tetap
      melaporkan premium list tersimpan, dengan penanda bahwa polis **tertunda**. *(AC 23 spec)*
- [ ] Pemanggilan ulang `INSERTJSONPOLISLIFE` untuk `IDPEGA` yang sama **tidak** menggandakan baris —
      dibuktikan dengan memanggilnya dua kali dan menghitung baris `JSON_POLIS`.
- [ ] Keadaan separuh **terdeteksi**: setelah penulisan, sistem membaca balik dan menandai kasus yang
      belum lengkap. Penanda itu terbaca lewat API dan menjadi masukan tiket 06.
- [ ] `SaveLifeinProduction_SQL` diperlakukan sebagai **titik potong** tersendiri; kegagalan
      sesudahnya tidak merusak apa yang sudah ter-commit dan dapat diulang dengan aman.
- [ ] Uang di dalam payload JSON ditulis sebagai desimal presisi arbitrer — **tidak** lewat `float`
      dan tidak lewat pembulatan diam. *(AC 14 spec; **ADR-0003**)*
- [ ] Ada test yang **gagal** bila urutan 05a → 05b dibalik. *(AC 24 spec)*
- [ ] Kegagalan yang jatuh ke `JSON_POLIS_ERROR` **terlihat**: sistem tidak menganggapnya sukses.
- [ ] Penulisan detail peserta NB ke `M_LIFE_PREMIUM_DETAIL` terjadi **di dalam alur simpan ini**,
      bukan dijadwalkan. Setelah respons simpan berhasil, baris untuk `PL_NUMBER` itu **sudah ada**.
      Tidak ada job/cron/worker terjadwal di jalur ini. *(AC lengkap di tiket 08;
      `[keputusan work owner]`)*

## Blocker

**Tidak ada pemblokir.** Catatan terbuka yang tidak memblokir: **OQ-067** (penempatan
`INSERTJSONOFFERLIFE` — ditulis di tahap penawaran, tiket 01).

## Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` Di `InsertJsonPolisLife_Act`, hanya **dua** `<pyStepsBlockName>` berisi `//`:

| Bagian | Baris | Nasib |
| --- | ---: | --- |
| step 16 `Commit` eksplisit | 5294 | **mati** — Go memegang transaksi (**ADR-0015**) |
| step 17 `Connect-REST` `ConvertJsonNusareToProduction` | 5383 | **mati** — tidak dipakai |

`[keputusan work owner]` **Juga tidak dimigrasikan** meski di korpus **AKTIF**: empat gerbang treaty
ID pada step 6–7 — `@contains(.ID,"1000032")` **QS** (1714), `"1000033"` **2nd QS** (1856),
`"1000034"` **SURPLUS** (1998), `"1000035"` **2nd SURPLUS** (2140). Itu logika polis lama. Jenis
treaty diambil dari data **`TypeCeding`** (`1`=QS, `2`=SURPLUS, `3`=QS+SURPLUS, `4`=XOL — ekspresi
baris ~3787). *(AC 30–31 spec)*

`[keputusan work owner]` Step **4** dengan precondition `@toDecimal(Local.currentdate)>25` (baris
1170) **tidak** direplikasi — ambang dibaca dari tabel, lihat tiket **02**.

⚠️ **OQ-066**: penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul ini —
berkas ini justru buktinya. Konfirmasikan ke work owner sebelum menyimpulkan langkah lain.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

</details>
