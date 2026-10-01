# Laporan uji manual — Master Contract Retro Life (butir A4)

> 01-10-2026, gelombang 6 brief `PROMPT-IMPLEMENTASI-TIGA-MODUL-LIFE-GELOMBANG-2.md` §3 A4. Dijalankan asisten di pohon kerja `dev`
> sesudah commit `0f71e05`.
>
> ⚠️ **SEBAGIAN.** Sesi ini tidak punya peramban: layar **tidak diklik**. Yang diuji adalah backend yang menyala, `npm run dev`
> yang melayani aplikasi, dan **setiap rute baca** yang dipanggil kelima layar, terhadap data DEV. Tombol simpan/hapus **tidak**
> dijalankan terhadap DEV (nol penulisan ke DEV) — perilakunya dikunci uji otomatis (bab 4). Uji klik di peramban tetap
> pekerjaan work owner (bab 5).

## 1. Lingkungan

| Hal | Nilai |
| --- | --- |
| Backend | `go build ./cmd/api` lalu dijalankan **tanpa** `-migrate`; `MODUL_AKTIF=mastercontractretrolife` (hanya modul ini, nol pekerja latar), `AUTH_STUB=true` lokal, `IS_PEGA_PROD=false` |
| Basis data | `ORACLE_DSN` dari `.env` lokal — skema `POOLDATA` (DEV). Alamat, akun, dan sandinya **tidak** dicatat di mana pun |
| Frontend | `npx vite` (`npm run dev`) dengan `DEV_PROXY_TARGET` ke backend uji |
| Identitas | header `X-Pelaku: UJI-AKUN-1` |
| Yang dicetak | kode status, nama kunci JSON, dan cacah baris saja — nol isi baris |

## 2. Rute baca per layar (terhadap DEV)

| Layar (PARITAS) | Rute | Hasil |
| --- | --- | --- |
| §2 halaman awal `MASTER CONTRACT RETRO LIFE` (`GridRetrocessionLife`) | `GET /tahun`, `GET /tahun?halaman=1` | 200 — 3 tahun treaty; kunci baris `id, treatyYear, underwritingYear, startDate, endDate, userId, tglUpdate` |
| §3 `Reins Type` (`InboxRetroLimitReinsurers`) | `GET /tahun/{id}/kontrak` (dua tahun) | 200 — 5 dan 4 kontrak |
| §3 dropdown `REINS TYPE` | `GET /jenis-reasuransi` | 200 — 5 jenis |
| §4 `Reinsurer List` (`InboxRetroLifeReinsurersList`) | `GET /kontrak/{id}/reinsurer` | 200 — 4 dan 3 reinsurer, `totalShare` + `totalBukan100` terisi |
| §4 autocomplete `REINSURER NAME` | `GET /master-reinsurer?cari=A` | 200 — 83 saran |
| §5 `Security Reinsurer` (`InboxSecurityReinsurerLife`) | `GET /reinsurer/{id}/security` | 200 — 0 baris pada reinsurer pertama kedua kontrak (data DEV), `eksposur` terisi |
| §6 `Business List` (`InboxBusinessLifeReinsurers`) | `GET /kontrak/{id}/business` | 200 — 2 business per kontrak |
| §6 autocomplete `BUSINESS NAME` | `GET /master-business?cari=A` | 200 — 15 saran |
| §6 `Copy to all Reinstype` (pratinjau) | `GET /business/{id}/salin-semua` | 200 — 4 dan 3 sasaran |
| `Delete` (popup dampak) | `GET /kontrak/{id}/dampak-hapus`, `/reinsurer/{id}/…`, `/business/{id}/…` | 200 — bentuk `{jenis, id, dampak}` |
| §6 autocomplete `R/I RATE`, `View Rate` | `GET /ringkasan-rate?cari=A`, `GET /rate?idusedby=1` | **503** berkalimat *"the life rate tables are not read yet - reading them requires work owner approval (OQ-MCRL-13)"* — sesuai A1 (izin belum tercatat) |
| tanpa layar (OQ-MCRL-07) | `GET /laporan/total-share-bukan-100?tahun=2025` | 200 — 0 baris |
| gerbang identitas | `GET /tahun` tanpa header | 401 |

## 3. Frontend (`npm run dev`)

| Pemeriksaan | Hasil |
| --- | --- |
| `GET /` | 200 — `<title>Nusantara Re</title>`, `#root` |
| proksi `GET /api/master-contract-retro-life/tahun` lewat Vite | 200 |
| halaman `modul/mastercontractretrolife/frontend/pages/MasterContractRetroLife.tsx` ditransformasi Vite | 200 (nol galat kompilasi) |
| `GET /api/menu` | memuat `mastercontractretrolife` berlabel `Master Contract Retro Life`, `dimigrasi: true` |

## 4. Yang dikunci uji otomatis (bukan klik)

| Hal | Uji |
| --- | --- |
| 33 sel `pxButton` di 8 Section, 31 unik, label VERBATIM per baris | `frontend/labels.test.ts` (sensus PARITAS §0) |
| lima jalur simpan, hapus berdampak, salin-semua, galat sampai layar | uji `services`/`handlers` modul (tiruan) dan uji bertag `db` (skema uji, SKIP tanpa `ORACLE_DSN`) |
| gerbang lengkap gelombang ini | `go build`/`vet`/`test` (+ `-tags db -p 1`), `gofmt`, `tsc`, `vitest`, `vite build` — hijau kecuali `TestNolAlamatLayananDiKode` di ekspor bersih (menuntut `.env` lokal) |

## 5. Belum dilakukan — untuk work owner di peramban

1. Buka menu **Master Contract Retro Life**; cocokkan judul `MASTER CONTRACT RETRO LIFE` dan grid tahun.
2. Klik setiap tombol kelima layar menurut PARITAS §2–§6 (31 tombol unik) — terutama `Add`/`Edit`/`Save`, `ReinsType`,
   `Reinsurer List`, `Total Share -->>`, `Security Reinsurer`, `Business List`, `Copy to all Reinstype`, `Delete`, `View Rate`.
3. Simpan dan hapus hanya dengan data `UJI-` di skema yang boleh ditulisi — **bukan** DEV.
4. `R/I RATE` dan `View Rate` tetap 503 sampai OQ-MCRL-13 diizinkan.
