# Laporan uji manual — Endorsement Life

> 01-10-2026, K6 brief `PROMPT-LANJUTAN-TIGA-MODUL-LIFE-KEPUTUSAN-OQ.md`. **Baca-saja**: `Save` dan setiap `Submit` **tidak**
> ditekan, `Add CSV Data` / `DELETE ALL` / `Generate Data Detail` tidak ditekan, nol berkas dipilih. Peramban Chrome headless lewat
> CDP dengan daftar klik terlarang; setiap permintaan `/api/` dicatat. Nol nama orang, nol nomor polis nyata: fixture `UJI-` saja.

## 1. Lingkungan

| Bagian | Isi |
| --- | --- |
| Backend DEV | `cmd/api` dari HEAD `bf45753`, `MODUL_AKTIF=endorsementlife`, `127.0.0.1:18084`, `.env` hanya di proses itu, `AUTH_STUB`; modul ini tanpa pekerja latar (`JalankanPekerja` = tanpa pekerja) |
| Server uji layar kasus | sementara, tidak di-commit: router modul ini (`handlers.RouterDengan`) di atas gudang tiruan berisi empat kasus `UJI-` (`EDMLF-901` QR, `-902` QP, `-903` TR, `-904` TP, masing-masing versi NB `UJI-PL-<Type>` + dua peserta), `127.0.0.1:18085`; permintaan non-modul (shell: menu, modul aktif) diteruskan ke backend DEV dan **hanya GET** yang diteruskan |
| Frontend | `vite --port 5176 --strictPort`, `DEV_PROXY_TARGET` = server uji. ⚠️ Penyimpangan dari brief (port **5174**): 5174 dipakai proses sesi lain (PID 168032) dan konfigurasi `strictPort` menolak berbagi port; proses itu tidak dihentikan |
| Mengapa server uji | DEV memuat **0** kasus endorsement (`GET /inbox` → `total 0`); membuka layar kasus dan popup `View Old Policy` di DEV menuntut membuat kasus (`POST /kasus` = tulisan) — dilarang |

## 2. DEV baca-saja

| Aksi | Rute | Hasil |
| --- | --- | --- |
| Kotak masuk | `GET /api/endorsement-life/inbox?halaman=1` | **200** — `{"baris":[],"total":0,"halaman":1,"ukuran":20}`, 0,05 detik |
| Gerbang kelayakan, polis `UJI-TIDAK-ADA`, EDM Type `1` | `POST /api/endorsement-life/kelayakan` (lima gerbang **tanpa tulis**) | **200** — `{"pesan":["Policy no Not Found !"],"boleh":false}`, 0,15 detik |
| EDM Type `3` (Batal) | — | **tidak dicoba**: gerbang 5 membaca invoice Arasapas, izinnya menunggu DBA (OQ-EDM-012) |

Log backend DEV: 2 baris (skema, alamat dengar), nol `INSERT`/`UPDATE`/`DELETE`/`MERGE`.

## 3. Layar dan tombol — dicocokkan dengan sensus PARITAS bab 10 (22 tombol)

| Layar | Teramati | PARITAS | Cocok |
| --- | --- | --- | --- |
| Kotak masuk (`InboxEndorsementLife`) | judul `Endorsement Life`; 13 kolom `Case ID` … `Status`; tombol `Create Addendum`, pengenal kasus (4), `Previous`/`Next` (mati) | `Create Addendum` b6620 dibangun; `Create` b3780 / `Cancel` b4228 (R11 panel tak terjangkau) dan `Create Case Endorsement` b5879 (`VIS=never`) tidak dibangun | ✅ |
| Buat endorsement (`EndorsmentLife_Section`) | medan `Policy No`, `EDM Type`, `Description`, `EDM Date`; tombol tutup (`aria-label` Close); **`Submit` tidak tampil** sebelum gerbang lolos | `Submit` b4226 (VIS `Protect.CARI1 = '0' && EdmType 1/3`) dibangun; `Process Policy No` b2414 `VIS=never` tidak dibangun | ✅ — `Submit` tidak ditekan |
| Kasus (`InputEDMLife`), keempat Type | judul `Endorsement Life Detail`; 14 medan kepala; grid 13 kolom (`DELETE ALL`, `POLICY NO` … `EXPIRED DATE`, rinci), 2 baris; tombol `View Old Policy`, `Save`, `Back`, `Upload CSV`, `View Upload`, `Add CSV Data` (mati — belum ada berkas), `DELETE ALL`, `Details` ×2, pager | `Upload CSV` b8973, `View Upload` b9340, `Add CSV Data` b10405, `DELETE ALL` b13607, `Save` b37202 dibangun; `Submit` b37494 / b38109 tampil hanya bila `.IsJsonPolis=1` (kasus sudah `Save`) — **tidak tampil**, benar untuk kasus belum disimpan; `View Premium` b36494 (`VIS 1=2`) dan `Process Policy No` b5233 (`VIS=never`) tidak dibangun | ✅ — `Save` tidak ditekan |
| Popup `View Old Policy` — QR (`EDMLF-901`) | judul `View Old Policy`; **12** kolom `NAME_OF_INSURED`, `CURRENCY`, `SUM_INSURED`, `CEDING_RETENTION`, `SUM_REASURED`, `SHARE_NUSANTARA_RE_GROSS`, `SUM_AT_RISK_GROSS`, `GROSS_PREMIUM`, `DEDUCTION`, `NET_PREMIUM`, `FACTOR`, `CLAIM_AMOUNT`; 2 baris; `Close` | `ShowLifePremiumSummary_EDM` b64965 | ✅ |
| Popup — QP (`EDMLF-902`) | **11** kolom (… `GROSS_PREMIUM_REFUND`, `DEDUCTION_REFUND`, `NET_PREMIUM_REFUND`, `CLAIM_AMOUNT`); 2 baris | b65522 | ✅ |
| Popup — TR (`EDMLF-903`) | **19** kolom (… `SHARE_NUSANTARA_RE`, `SUM_AT_RISK_RETRO`, `RETROCEDED_SHARE`, `SHARE_RETRO`, `RATE`, `FACTOR`, `GROSS_PREMIUM_REFUND_RETRO`, `DISCOUNT_PREMIUM_RETRO`, `DISCOUNT_PREMIUM_REFUND_RETRO`, `RI_ADMIN_FEE_REFUND_RETRO`, `BROKERAGE_FEE_RETRO`, `NET_PREMIUM_REFUND_RETRO`, `CLAIM_AMOUNT`); 2 baris | b66083 | ✅ |
| Popup — TP (`EDMLF-904`) | **15** kolom (… `SUM_AT_RISK_RETRO`, `RETROCEDED_SHARE`, `SHARE_RETRO`, `RATE`, `GROSS_PREMIUM_RETRO`, `DISCOUNT_PREMIUM_RETRO`, `RI_ADMIN_FEE_RETRO`, `BROKERAGE_FEE_RETRO`, `NET_PREMIUM_RETRO`); 2 baris | b66640 | ✅ |
| Dialog `Upload CSV` | dibuka lalu ditutup, nol berkas dipilih (nol `POST /unggah`) | b8973 | ✅ |
| Dialog `View Upload` | `Close`, **`Generate Data Detail`** tampil — tidak ditekan | `ViewCSVResult_LifeEDM` b14322 | ✅ |
| Rincian peserta `Details` | `GET /kasus/{id}/peserta/{pid}` → 200 `{peserta, spreading}` di keempat kasus | `PL_Detail_Sec` + `RetroDetailLife` | ✅ |
| Halaman terima kasih (`ConfirmSubmitEDM`) | **tidak dibuka** — hanya terjangkau sesudah `Submit` keputusan | `Close` b1366 | ⏸ menunggu uji di skema yang boleh ditulisi |
| `ShowLifePremiumSummary_EDM` `Submit` b67733 | — | tidak dibangun (R02 layar yatim; efeknya = jalur Confirm) | ✅ |

Sensus: **15** tombol dibangun — **11** teramati di layar (`Create Addendum`, `Upload CSV`, `View Upload`, `Add CSV Data`,
`DELETE ALL`, `Save`, `Generate Data Detail`, dan keempat `View Old Policy`, satu per kasus sesuai Type), **3** tersembunyi sesuai
`VIS` (`Submit` buat — gerbang belum lolos; `Submit` keputusan ×2 — kasus belum `Save`), **1** tidak terjangkau tanpa `Submit`
(`Close` terima kasih); **7** tombol tidak dibangun — nol yang muncul.

## 4. Sensus jaringan

56 permintaan `/api/` dari peramban, **seluruhnya GET**: `inbox`, `kasus/{id}`, `kasus/{id}/peserta`, `kasus/{id}/peserta/{pid}`,
`kasus/{id}/polis-lama` (keempat kasus), dan shell `menu`, `modul-aktif`, `klaim-life`. Server uji mencatat **0** permintaan non-GET ke
modul ini. Bagian DEV §2: satu `POST /kelayakan` tanpa tulis, dikirim langsung (bukan dari peramban).

## 5. Belum dilakukan — untuk work owner, di skema yang boleh ditulisi (bukan DEV)

1. `Submit` buat endorsement atas polis nyata (lahirnya kasus + salinan versi), lalu `Save`, `Upload CSV` + `Add CSV Data`, `DELETE ALL`.
2. `Submit` keputusan Confirm/Decline → halaman `ConfirmSubmitEDM` (`Close`), nomor `<polis>/01` (K3), baris `LIFEINPRODUCTION` (K4) dan
   `M_LIFE_PREMIUM_DETAIL` (K5).
3. EDM Type `3` (Batal) sesudah izin baca invoice Arasapas (OQ-EDM-012).
