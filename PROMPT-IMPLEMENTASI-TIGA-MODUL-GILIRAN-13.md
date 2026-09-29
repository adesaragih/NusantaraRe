# PROMPT — GILIRAN 13 *(sesi baru **atau** sesi yang sama, folder `OUTPUT_HASIL_RNM`, cabang `main` @ `0b05391` atau lebih baru)*: **tiga titik buta uji layar — PremiumList membuat kasus (bn) → Claim Life baris adjustment pertama → data uji sintetis `UJI-*` untuk skema uji**

> Hanya konteks **Claim Life, PremiumList Life, Komite Claim Life**. GILIRAN-3 *(A/B/C)*, 4–12 tetap rujukan; aturan berhenti
> GILIRAN-6 berlaku; **setiap pembacaan activity mencetak `pyStepsBlockName`** *(GILIRAN-12)*. Berkas `tco_*` tidak disentuh.
> Sebab giliran ini: work owner mencoba layar dan **tidak ada tombol pembuat kasus di PremiumList** — panduan uji bab 0 sudah
> mencatat tiga titik buta; dua di antaranya ternyata **celah kode**, bukan data awal.

## 0. BUKTI — DIPERIKSA ASISTEN 29-09-2026

| Titik buta | XML | Kode sekarang | Sebab |
| --- | --- | --- | --- |
| **PremiumList tidak dapat membuat kasus** | `Section/PremiumList.xml` tombol `Input Offer` b16964 dan `Input Premium` b16472 → `runActivity` `CreateInputLife` b3291/b3578 *(`FlagPolicy = "0"`)* dan b3939/b4214 *(`FlagPolicy = "1"`)*. `Activity/CreateInputLife.xml`: `param.classname = "ASM-FW-GISFW-Work-LIFE"` b274, `FlowType = "pyStartCase"` b363, `Call svcAddWorkObject` b444, `curWorkPage.FlagOnGoingPolicy = Param.FlagPolicy` b618, `Obj-Save` b726, `Commit` b858, catatan `"ASSIGN-WORKLIST " + pzInsKey + "!InputPolicyHolder"` b982 | `pages/premiumlist/InboxPremiumList.tsx:99–100` merender keduanya sebagai `BelumTersedia`; **nol** rute `POST /api/polis-life` | komentar kode: *"`SEQ_WORK_POLIS` tidak dibuat migrasi mana pun; kolom `FlagOnGoingPolicy` belum ada; keputusan skema"*. ⛔ **`SEQ_WORK_POLIS` SUDAH diputuskan** di brief modul PremiumList **pl3** *(baris 43: pengenal dirakit dari `SEQ_WORK_POLIS`, pola `PengenalWorkBerikut`)* — terlewat, bukan terbuka |
| **Claim Life tanpa baris adjustment pertama** | `Section/ClaimLifeDetailGCNM.xml` grid `.AdjustmentList` b17126: `Add` b17937 → `addRow` b17947 + `SetIndexAdjustmentList` b17991/b18123 *(langkah 1 `.IsCheck = true` b328; langkah 3 mewarisi delapan kolom dari `.AdjustmentList(1)`)*; `Delete` b19120 → `deleteRow` b19130 | rute yang ada: `…/putaran` *(mewarisi dari baris pertama — **menuntut** baris pertama sudah ada)*, `…/akseptasi`, `…/tolak`, `…/komite`; **nol** rute membuat baris saat grid kosong | `Add` pada grid kosong = baris pertama *(`.AdjustmentList(1)` adalah dirinya sendiri)*; penulis lain `AdjustmentList` yang harus dibaca: `SaveInsuredClaim_Act`, `SavePesertaClaim`, `SaveOutStandingLife_Act`, `SpreadingClaimLife_Act` |
| **Komite tanpa roster dan baris siap serah** | roster `EMAILKOMITE` berisi **data orang** *(nol baris boleh disalin)* | — | bukan celah kode: butuh **data sintetis** di skema uji |

## 1. KEPUTUSAN

| Butir | Isi |
| --- | --- |
| **pl3** *(sudah ada)* | `SEQ_WORK_POLIS` + pengenal berformat dirakit di repository. Awalan pengenal **dibaca dari data**: pola nilai `IDPEGA`/`pzInsKey` kelas `ASM-FW-GISFW-Work-LIFE` di tabel warisan *(agregat — bentuk awalan saja, nol baris disalin)*; bila tidak dapat ditentukan → `[terbuka]` berbukti, bukan dikarang |
| **bn** `[DIPUTUSKAN; veto work owner]` | migrasi **`057_seq_work_polis_dan_flag_ongoing.sql`** *(+ down)*: `SEQ_WORK_POLIS` dan kolom `T_WORK_POLIS.FLAG_ONGOING_POLICY VARCHAR2(1)` bernilai VERBATIM `"0"` *(Input Offer)* / `"1"` *(Input Premium)* dari b3291/b3939; dipakai decision `IsFlagOnGoingPolicy` *(flow `InputPolicyHolder.xml` b909)* — baca decision itu utuh untuk memastikan tahap awal yang dipilih tiap nilai. Penjaga kata cadangan Oracle hijau |
| **bo** `[DIPUTUSKAN; veto work owner]` | baris adjustment pertama dibuat oleh tombol `Add` **seperti XML**: baris kosong + `IsCheck = true` *(`models.PenandaDipilih`)*; bila sudah ada baris pertama, `Add` = jalur `putaran` yang ada *(satu rute atau dua, pilih yang lebih sederhana dan catat)*; gerbang tahap + pemegang + tujuh gerbang `.STS_REJECT` yang sudah ada |

## 2. URUTAN — satu commit per paket

| # | Paket | Isi |
| ---: | --- | --- |
| 1 | **PremiumList membuat kasus** | migrasi `057` *(bn)*; `POST /api/polis-life` berparameter `flag` *(`0`/`1`)* → baris `T_WORK_POLIS` *(pengenal pl3, `POSITION` dan `STATUS` sesuai konektor pertama `InputPolicyHolder.xml` dan decision `IsFlagOnGoingPolicy`)* + baris `T_PREMIUM_LIST` kosong yang dibutuhkan layar berikut; satu transaksi; tombol `Input Offer` dan `Input Premium` di kotak masuk menjadi aktif dan membuka layar tahapnya; `BelumTersedia` dicabut; uji murni + handler + JS; tiket 00/01 diralat bertanggal *(pl3 terlewat)*. Commit `premiumlist-life: Input Offer dan Input Premium membuat kasus (pl3, bn)` |
| 2 | **Claim Life baris adjustment pertama** | baca kelima penulis `AdjustmentList` sebagai pohon; bangun **bo**; tombol `Add`/`Delete` pada grid adjustment di layar Detail; tiket 03 diralat bertanggal. Commit `claim-life: bo — Add membuat baris adjustment pertama` |
| 3 | **Data uji sintetis** | `APP_RNM/internal/repository/skemauji/data_uji_tiga_modul.sql` *(atau berkas di `testdata/`)*: satu polis per tahap PremiumList, satu klaim per tahap Claim Life **dengan** baris adjustment, roster `EMAILKOMITE` **sintetis** *(nama `UJI-KOMITE-1…n`, surel `uji-…@contoh.invalid`)*, satu baris siap serah Komite; seluruh pengenal `UJI-*`; berkas menolak dijalankan bila skema = `POOLDATA` *(pagar yang sama dengan `-migrate-down`)*; **tidak** dijalankan executor. Panduan uji bab 0 diperbarui: tiga titik buta ditutup, cara memuat data uji ke skema uji |
| 4 | **Status tiket + panduan** | tiket yang tersentuh; `PANDUAN-UJI-LAYAR-TIGA-MODUL.md` bab PremiumList §1.2 dan Claim Life §1.3 diperbarui |

## 3. LAPORAN · TELEMETRI

Baris pertama alasan berhenti; tabel **paket → commit → tombol/rule XML → rute/komponen**; ralat tiket; OQ; angka uji tiap commit
dengan dan tanpa tag `db`; bab **TELEMETRI EKSEKUSI**. `-migrate` tidak dijalankan executor; work owner menjalankan `057` dari `main`.

---

*Disusun 29 September 2026 dari `Section/PremiumList.xml` (b3291/b3578/b3939/b4214, b16472/b16964), `Activity/CreateInputLife.xml`
(b274/b363/b444/b618/b726/b858/b982), `Section/ClaimLifeDetailGCNM.xml` (b17126/b17937/b17991/b19120), brief PremiumList pl3, dan
kode `InboxPremiumList.tsx`, `handlers` (rute `polis-life` dan `adjustment`).*
