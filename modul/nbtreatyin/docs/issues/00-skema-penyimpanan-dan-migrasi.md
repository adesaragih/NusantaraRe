# 00: Skema penyimpanan dan migrasi — menunggu SENSUS PROPERTI, bukan lagi bahan dari DBA

**Status:** selesai — sensus dicocokkan ulang ke XML; butir DBA tetap `[terbuka]` *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: needs-info — ⚠️ **yang ditunggu berubah**)*
**Blocked by:** ⭐ **sensus properti `PolicyTreatyIn`** — ronde tersendiri, pekerjaan tim migrasi

> ⭐⭐ **P1 dan P29 SUDAH TERJAWAB.** `[terverifikasi]` 2026-09-22 — dua baris ini semula berbunyi:
> > *"**Status:** needs-info · **Blocked by:** P29 *(contoh isi JSON polis)* · P1 *(naskah tiga
> > stored procedure)*"*
>
> Naskah **empat** stored procedure dan **dua contoh** `DATA_JSON` diterima. ⛔ **Tiket ini tetap
> belum dapat dikerjakan**, tetapi sebabnya berubah: bukan lagi bahan dari `[DBA]`, melainkan
> **sensus properti yang belum dijalankan**.
**Menutup:** — *(0 AC)*

## Hasil & nilai pengguna

Hari ini data realisasi treaty tersimpan sebagai **satu dokumen teks** di dalam satu kolom.
`[terverifikasi]` Fungsi pembentuknya — `@ASM.GetPageJSONString()`, naskahnya diterima pada P15 —
**tidak memilih apa pun**: ia memotret seluruh halaman kerja pada saat penyimpanan. ⛔ Karena itu
bentuk datanya bukan skema melainkan potret.

> ⚠️ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat di atas semula berakhir dengan
> ⛔ *"… dan **tidak dapat dirancang ulang tanpa melihat contoh isinya.**"*
>
> ⛔⛔ **Contoh saja ternyata TIDAK CUKUP, dan itu terbukti justru setelah contohnya datang.**
> Dua contoh diterima; keduanya berselisih **30 medan skalar** dan **tidak berbagi satu pun daftar
> bersarang**. ⭐ Yang dibutuhkan adalah **sensus properti korpus**, dengan contoh sebagai penguji.

Sesudah tiket ini, penyimpanan realisasi treaty punya bentuk yang dapat dinyatakan — tabel datar,
bukan dokumen. ⛔ **Tiket ini belum dapat dikerjakan**, dan nomornya dipesan supaya tidak diambil
tiket lain.

## Blocker

| Butir | Yang ditunggu | Pemilik |
| --- | --- | --- |
| ⛔ **sensus properti `PolicyTreatyIn`** | daftar lengkap properti dan simpul bersarang, dari korpus — **bukan** dari contoh | tim migrasi |

> ✅ ~~`P29` — contoh isi dokumen polis~~ **DITERIMA** 2026-09-22, dua contoh.
> ✅ ~~`P1` — naskah stored procedure~~ **DITERIMA** 2026-09-22, **empat** naskah.

⛔⛔ **Kenapa contoh saja tidak cukup.** `[terverifikasi]` **Tiga** contoh diterima — dua
proporsional, satu **non-proporsional XOL** — dan ketiganya hanya berbagi **23 dari 74** medan
gabungan *(37 · 64 · 47 skalar)*. Daftar bersarangnya berbeda-beda; contoh kedua bahkan tidak punya
`LocationList`. Sementara korpus menunjukkan sedikitnya **94 properti** dan **9 simpul bersarang**,
empat di antaranya — `OldData`, `TreatyXOLList`, `TreatyXOLDifferenceList`, `TreatyDifference` —
**tidak muncul di ketiga contoh** padahal dirujuk ratusan kali.

⛔⛔ **Dan bentuknya mengikuti JENIS TREATY.** `QuotationData.ProportionalType` dan
`IsNewPolicyNonProp` menentukan susunannya. Yang paling berbahaya: **`ListInstallment` bersarang
pada satu kasus dan datar pada kasus lain** — nama sama, kedalaman berbeda. `[terverifikasi]`
Korpus memuat **kedua** bentuk: `ListInstallment(n).InstallmentList` **58 rujukan di 6 berkas**,
`ListInstallment(n).<skalar>` **101 rujukan di 11 berkas**.

⛔ **Pengurai yang menganggap bentuknya seragam akan patah**, dan rancangan tabel tidak dapat
dibuat sebagai satu bentuk tunggal tanpa lebih dulu memutuskan bagaimana kedua jenis treaty
ditampung. ⛔ **Itu keputusan rancangan — ronde berikutnya, sesudah sensus.**

⭐ **Urutan yang benar:** rancangan dibangun dari sensus korpus, **contoh sebagai penguji**.

## Yang sudah diketahui — ⭐ bertambah banyak 2026-09-22

`[keputusan work owner]` **P29** — dokumen JSON **dibuang**; pembacaan pindah ke **view relasional
data treaty** yang sudah ada, dan penyimpanan menjadi **tabel datar**.

⭐ **Yang kini sudah pasti, dan memperkecil pekerjaan:** `[terverifikasi]`

| Hal | Keadaan |
| --- | --- |
| ⭐ **data kontrak treaty SUDAH relasional** | `POOLDATA.TREATY_IN` — `ID` + **20 kolom datar**. Tidak perlu dirancang ulang |
| ⭐ **tujuh dari delapan kolom `json_polis` sudah datar** | `IDPEGA` `TGL_INPUT` `NOPOLIS` `NOENDORS` `PRODKE` `TGL_PROD` `USERNAME` — dapat dipertahankan apa adanya |
| ⛔ **hanya `DATA_JSON` yang perlu dipecah** | satu kolom CLOB; **inilah** seluruh sisa pekerjaan perancangan |
| ⭐ sisi baca | view `POOLDATA.TREATYINDETAILJOINEDM`, **39 kolom**, kecukupan sudah diuji |

⛔⛔ **Satu butir `[terbuka]` yang menyentuh tiket ini secara tidak langsung:** data produksi memuat
**galat angka uang yang sudah tersimpan permanen** — ekor `2,76 x 10^-7` pada `NetPremium`. Tiga
pilihan penanganannya tercatat di **P29**, pemutusnya `[work owner]` dan `[Finance]`.
⚠️ Ia **tidak menahan tiket ini**, tetapi menentukan **apa yang dianggap benar saat data lama
dipindahkan** — dan **ambang uji paritas** di tiket **15**.

⚠️ `[terverifikasi]` **Satu procedure melayani lebih dari satu modul.**
`RDBList\SavePolisTreatyIn_SQL.xml` memanggil `POOLDATA.PEGA_JSON_POLIS_TREATYIN` dengan
**delapan parameter**, dan procedure yang sama dipanggil **NB Treaty In**, **NB FacIn**, dan
**EDM Treaty In**. ⭐ **EDM mengisi dua parameter yang NB kirim sebagai `NULL` dan `'0'`.**
⛔ **Rancangan tabel kelak harus menampung keduanya** — bila dirancang hanya dari NB, jalur EDM
akan kehilangan dua nilainya.

## ⭐ Tiket lain TIDAK bergantung pada tiket ini

⛔ **Penting, dan disebut tegas:** ke-15 tiket lain **tidak menunggu** tiket ini. Seluruhnya ditulis
dalam bahasa domain — *"data realisasi treaty"*, *"riwayat akseptasi"*, *"putusan per tingkat"* —
⭐ sehingga pekerjaan alur, tahap, wewenang, kronologi, layar, dan penggolongan **dapat berjalan
sekarang** dan menyesuaikan bentuk penyimpanan ketika ia tiba.

## ADR terkait

- **ADR-0003** — uang tidak `float`
- **ADR-0009** — migrasi penuh, tanpa koeksistensi dua penulis

## Catatan

⚠️ Tiket ini **satu-satunya** yang boleh menyebut nama procedure dan parameter, dan itu pun hanya
sebagai keterangan **apa yang ditunggu**. ⛔ Nol DDL, nol `CREATE TABLE`.

## ⭐ Hasil implementasi 2026-10-03 — sensus dicocokkan ulang ke XML

- Sensus properti dari rule TERJANGKAU (DataTransform, When, DecisionTable ikut disapu):
  `docs/SENSUS-PROPERTI-POLICYTREATYIN.md` (bangkitan `docs/alat/sensus.py`).
- Katalog kolom tunggal `backend/models/katalog.go` → migrasi 320-328 dan
  `docs/STRUKTUR-TABEL-NB-TREATY-IN.md` dibangkitkan darinya (`docs/alat/skema.py`); uji
  `repository/kolom_test.go` TestKatalogSepakatDenganDDL menagih kesepakatannya.
- `T_GENERAL_POLIS` memuat **72** kolom katalog (69 medan `PolicyTreatyIn` + 3 halaman kerja:
  `PositionNote`, `NBStatus`, `TreatyIn.ID`) ditambah kolom kunci/generasi. Tidak dibuat, dan sebabnya:
  `TOTAL_*` (turunan baris; penjaga repo melarang nama ber-awalan TOTAL_ di migrasi), `LAYER*` (ID-22),
  `isApprovedtoDeptHead` (P36, AC 64), `IsEDMInputOnNB` (hanya dipakai pembongkar JSON, AC 62),
  `IDNewBisnis` (nol rule).

## ⭐ Putaran 2 — paket penyimpanan (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 11–12, bab 2 K4/K16/K17; rincian kolom `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

⛔ **RALAT** atas *Hasil implementasi 2026-10-03* di atas. Bunyi lama, dikutip: *"`T_GENERAL_POLIS` memuat
**72** kolom katalog (69 medan `PolicyTreatyIn` + 3 halaman kerja …) … Tidak dibuat, dan sebabnya: … `IsEDMInputOnNB`
(hanya dipakai pembongkar JSON, AC 62), `IDNewBisnis` (nol rule)."* dan *"migrasi 320-328"*.

Bunyi baru:
- Migrasi **320–327 = tepat delapan tabel diagram grilling**; 328 `T_POLIS_SUGGEST`, 329 `T_POLIS_MEDAN_LAIN`,
  330 `M_NBTRIN_PERAN_TEMPAT` **dihapus** (K4, K17, K16). Ditagih `TestTabelDanKolomMengikutiDiagramGrilling`.
- `T_GENERAL_POLIS` tetap **72** kolom katalog, isinya berubah: **+** `IS_EDM_INPUT_ON_NB` (rancangan §4.1;
  jalur XOL kini dibangun, K8) dan `ID_NEW_BISNIS` (rancangan §4.1, data guide); **−** `IS_OJK_NOPOLIS`,
  `BROKERAGE_FEE` (hanya ditulis). `TGL_TUTUP` dibuang (di luar diagram).
- `T_POLIS_QUOTATION` 10 medan diagram + 6 RALAT (dari 26); `T_POLIS_CEDING` di bawah quotation
  (`QUOTATION_ID`, `CEDING_CO_ID`); angsuran mengikuti rancangan §4.3 + diagram R61.
- Status tetap **selesai**.
