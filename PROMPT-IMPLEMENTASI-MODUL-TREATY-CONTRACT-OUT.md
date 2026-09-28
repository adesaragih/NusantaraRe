# PROMPT — MODUL **Treaty Contract Out** *(sesi BARU, folder `OUTPUT_HASIL_RNM\.worktrees\treaty-contract-out`, cabang `modul/treaty-contract-out` @ `725636a`)*: **tiket 01 → 12 seluruhnya, dari XML, dalam satu giliran**

> Buka **sesi Claude Code baru** di folder worktree di atas, lalu `/mattpocock-skills:implement` dan tempel brief ini utuh.
> Aturan yang berlaku: `PROMPT-EKSEKUSI-HULU-HILIR.md` §4–§5 *(tata letak, templat)*, `PROMPT-INDUK-TIGA-MODUL.md` §0.3–§0.5 dan §2
> *(pembagian berkas)*, lanjutan 8 §1 *(mekanisme giliran)*, aturan berhenti GILIRAN-6, larangan keamanan, bab TELEMETRI EKSEKUSI.
> Sesi ini **tidak** menyentuh berkas modul lain dan **tidak** merge ke `main`.

---

## 0. KEADAAN AWAL — DISIAPKAN DAN DIVERIFIKASI ASISTEN 28-09-2026

| Hal | Keadaan |
| --- | --- |
| Worktree | `modul/treaty-contract-out` dari `main` `725636a`; `.env` `HTTP_ADDR=:8093`, frontend `DEV_PROXY_TARGET=http://localhost:8093`, Vite `5185`; `go build` OK; `npm install` OK |
| Posisi di rumpun | **HULU** rumpun Treaty dan Fac In. `spec.md` b106: modul ini **satu-satunya penulis** `TREATYREINSURER`, `TREATYBUSINESS`, `PROPORTIONALARRG`, `TREATYCONTRACT`, `TREATYYEAR`, `MTREATYSECURITY`; hilir yang membaca: Claim Prop, Komite Claim Prop, Claim Fac In *(b224–b231, OQ-042)*. Masuk = master dibaca saja *(b107)*: `REINSURANCETYPE`, `TREATYDESC`, `TREATYGROUP`, mata uang, `TREATYEXCHANGEYEARLY`, `CATEGORY_ATTACH_REAS` |
| Bahan | `.scratch/treaty-contract-out/spec.md`, `grilling-ronde-1.md`, `grilling-ronde-1-jawaban.md`, `dba-procedures.md`, `issues/01`–`12` *(12 × `ready-for-agent`)* |
| Blocker tiket 01 `CL-01` | **sudah ada**: kerangka aplikasi + seam API adalah `APP_RNM` itu sendiri *(Claim Life tiket 01)* |
| Katalog DEV *(agregat)* | keenam tabel warisan **ada** di `POOLDATA`: `TREATYYEAR` 182 baris, `TREATYCONTRACT` 480, `PROPORTIONALARRG` 3.188; nama `T_TREATY%` **belum dipakai** |
| Migrasi | rentang modul **300–319**; `main` terakhir `056` |
| Korpus | `D:\XML\RNM_BRD\Treaty Contract Out\` *(READ-ONLY)*: 3 harness portal, 2 flow action, 49 section, 168 activity, `Struktur_InboxTreatyContract.xlsx` |
| REFERENSI_UI | `REFERENSI_UI/frontend/src/pages/master/AksiTreatyContract.tsx` dan keluarga `master/*` — rujukan tampilan terdekat; **disalin, tidak disunting** |

## 1. KEPUTUSAN YANG MENGIKAT GILIRAN INI

| Butir | Isi |
| --- | --- |
| **tco1** `[DIPUTUSKAN; veto work owner]` — nama tabel skema baru | Aplikasi memakai satu skema `POOLDATA` *(`ORACLE_SCHEMA`)*, sedangkan keenam nama warisan **sudah dipakai** tabel hidup. Skema baru tiket 01 karena itu memakai awalan **`T_`** *(`T_TREATYYEAR`, `T_TREATYCONTRACT`, `T_TREATYREINSURER`, `T_TREATYBUSINESS`, `T_PROPORTIONALARRG`, `T_MTREATYSECURITY`, + dua tabel lain dari tiket 01, + sequence)*, **nama kolom yang dibaca hilir dipertahankan** *(kontrak b230–b232)*; tabel warisan **tidak disentuh**, tidak ditulis, tidak di-`ALTER`. Pola yang sama dengan `T_CLAIMLF_*`/`T_PREMIUM_LIST*` |
| **tco2** | Skrip migrasi dan rekonsiliasi tiket 01 = **kode + uji bertag `db`** *(SKIP tanpa skema uji)*; **tidak** dijalankan terhadap DEV. `-migrate` hanya oleh work owner dari `main` sesudah penyatuan |
| **tco3** | Hilir *(Claim Prop, Komite Claim Prop, Claim Fac In)* belum dimigrasi: kontraknya ditulis sebagai `repository` baca-saja berkolom VERBATIM + uji kontrak, untuk dipakai saat mereka dibangun |
| **o** | Prosedur **tidak dipanggil** *(`PEGA_PROPORTIONALARRG`, `PEGA_M_PROPORTIONALARRG_CHILD`, `PROSESCOPY`, dst. di `dba-procedures.md`)*: logikanya ditiru dari `ALL_SOURCE` yang tercatat. Nol `COMMIT` di teks SQL *(ADR-U-0029)* |
| **Lampiran** *(tiket 12)* | jalur dokumen yang **sudah ada** di `main` dipakai ulang *(outbox `T_LOG_SERVICE_RNM`, pelaksana stub, `UNGGAHAN_DIR`)*; endpoint penyimpanan nyata tidak dipanggil |
| **Kurs** *(tiket 11)* | `TREATYEXCHANGEYEARLY` dibaca saja; uang `apd.Decimal`, nol `float` |

## 2. XML — TITIK MASUK DAN MENU *(nomor baris = `sed -e 's/></>\n</g'`)*

| Harness portal *(kelas `Data-Portal`)* | Section utama | Catatan |
| --- | --- | --- |
| `InboxTreatyContract` | `GridTreatyContract` *(3.205 baris)*, `InputTreatyContract` | daftar + form kontrak |
| `InboxTreatyContractDescription` | `ViewDetailDescription` *(14.193 baris)*, `NitipKurs`, `ViewDetailDescriptionProp`, `ViewDetailDescriptionNonProp` | klausul 25 jenis; tombol `runActivity` → `BrowseDescriptionLimit` b5049/b5457/b8150, `testingKurs` b5110, `SetKirimIDDesc` b5154, `PanggilID` b5228, `RefreshErrorProportionalarrg` b5348/b5730 |
| `InboxTreatyContractReinsType` | `GridTreatyContractReinsType`, `PanggilReinsType` | saringan jenis reasuransi |

Flow action: `DetailTreatyExclustion` → `DetailTreatyExclustion_Sec`; `TreatyOutAttachContent` *(lampiran)*. Tidak ada flow
proses *(folder `Flow\` kosong)* → modul ini **master**, bukan alur tugas; tidak ada tahap/inbox per peran.

**Menu `[dari bukti]`**: kelompok **Treaty Contract Out** *(sudah ada di sidebar sebagai "belum dimigrasi")* → tiga butir
bernama VERBATIM dari `pyLabel`/judul ketiga harness portal *(baca dan catat barisnya; jangan mengarang nama)*. Keterangan
"belum dimigrasi" dicabut untuk kelompok ini saja.

Sebelum tiket 01: baca `Struktur_InboxTreatyContract.xlsx` dengan `.scratch/alat/baca-xlsx.ps1` untuk pohon penyarangan, lalu
tulis peta **harness → section → tombol → activity → RDB/SQL** di tiket 01 bab *"Pembacaan ulang XML"*. Tiap activity yang
**menulis** ditelusuri dari metode langkahnya; `WHEN` dibaca dari urutan aksi; penulis medan dicari, bukan hanya pembacanya.

## 3. URUTAN — satu commit per tiket `treaty-contract-out: tiket NN — <judul>`, tanpa pesan di antaranya

Urutan mengikuti *Blocked by* di tiap tiket:

| # | Tiket | Bergantung pada |
| ---: | --- | --- |
| 1 | **01** skema relasional + migrasi + tipe dirapikan *(PREFACTOR; tco1, tco2; migrasi `300`+)* | — |
| 2 | **02** jenis reasuransi master dibaca + saringannya | 01 |
| 3 | **03** tahun treaty, periode, anti-dobel | 01 |
| 4 | **12** lampiran di tahun treaty | 03 |
| 5 | **04** kontrak treaty di dalam tahun | 02, 03 |
| 6 | **05** reinsurer + total share | 04 |
| 7 | **07** business + nonaktif | 04 |
| 8 | **08** klausul satu tabel 25 jenis + validasi | 04 |
| 9 | **06** security reinsurer | 05 |
| 10 | **11** kurs USD ke IDR | 08 |
| 11 | **09** simpan atomik lintas enam tabel | 05, 06, 07, 08 |
| 12 | **10** kaskade hapus, popup, klausul yang tetap hidup | 06, 07, 09 |

Tiap tiket: bab **Pembacaan ulang XML** *(path + baris)*; ralat bertanggal bila XML membantah tiket; uji murni + handler + JS;
layar di `frontend/src/pages/treaty-contract-out/`, komponen di `components/treaty-contract-out/`, label di
`assets/labels.treaty-contract-out.ts`, rute di `handlers/rute_treaty_contract_out.go`, berkas Go berawalan `tco_`;
`.scratch/treaty-contract-out/PARITAS-LAYAR-DAN-AKSI.md` dan `LAPORAN-GILIRAN.md` +1 bab per tiket; penjaga kata cadangan
Oracle dan penjaga nama orang hijau; kode bersama *(Shell, `labels.ts`, penjaga lintas modul)* hanya **ditambah**, dilaporkan.

## 4. BERHENTI · LAPORAN · TELEMETRI

Laporan pertama paling cepat sesudah **tiket 04**; berhenti lebih awal hanya dengan baris pertama *"konteks menipis: kira-kira N%
tersisa"*. Laporan: tabel **tiket → commit → harness/tombol XML → rute/komponen**; ralat tiket; OQ dibuka/ditutup; kontrak hilir
yang ditulis; kode bersama yang disentuh; angka uji tiap commit *(Go, JS, vet, gofmt, tsc, build)*; bab **TELEMETRI EKSEKUSI**
per tiket. Asisten memverifikasi lalu menyatukan ke `main`; work owner menjalankan `-migrate`.

---

*Disusun 28 September 2026 dari `spec.md` (Keluar/Masuk b106–b107, kontrak hilir b224–b232), dua belas tiket (*Blocked by*),
tiga harness portal dan section-nya (tombol → activity), `REFERENSI_UI/frontend/src/pages/master/`, dan katalog DEV
(keenam tabel warisan dan cacah barisnya — agregat saja). Worktree dibuat dan diuji build oleh asisten.*
