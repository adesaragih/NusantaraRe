# Paritas layar dan aksi — Komite Claim Life

> Disusun 28-09-2026 (giliran 10, tiket 01–09). Satu baris per tombol/aksi korpus → rute/kontrol sistem
> baru. Nomor baris = `sed -e 's/></>\n</g'` atas berkas korpus `D:\XML\RNM_BRD\Komite Claim Life\`.

## Menu

| Korpus | Sistem baru | Catatan |
| --- | --- | --- |
| worklist `KomiteRouter` (`Flow/KomiteLife_Flow.xml` b594 `WorkList`, b619 `pyRouteTo Custom`) | sidebar **Komite Claim Life → Inbox Komite** → `GET /api/komite` | `[tidak ada di korpus]` sebagai menu; nama butir kosakata kami (bg). Anggota **berjalan** = anggota pertama ber-`KomiteAproval` 0 (`KomiteRouter` b382; kode transisi `6` `[dugaan kuat]`) |

## Layar keputusan `ShowTransfer` (lewat `FlowAction/ViewTransferDtl.xml` b91)

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| dropdown `.AcceptStatus` b32607 (`pyRequired`, `pyListSource associated`) | `<select required>` Setuju/Tolak di `KasusKomite.tsx` | ✅ enum `{1,2}` `[keputusan work owner]` — daftar properti tidak diekspor |
| label "Are you sure to accept this document?" | VERBATIM | ✅ |
| `.KomiteComment` b31001 (`pyRequired false`) | `Comment` | ✅ |
| `Submit` b34722 → `finishAssignment` → `KomitePostAdjustment` | `POST /api/komite/{id}/keputusan` | ✅ tiket 02 (+04a/04b/05/06 di transaksi yang sama) |
| `Cancel` b33880 | kembali ke Inbox Komite | ✅ |
| grid `.KomiteList` b29727 — `Committee` / `Status` / `Date Approve` / `Comment` | tabel tangga, judul VERBATIM, dari `GET /api/komite/{id}/riwayat` | ✅ tiket 09 |
| `Refresh` b22882 | muat ulang sesudah keputusan | ✅ (otomatis) |
| `Add attachment` b23511 → `AttachDocumentLife` → `SaveAttachLife` | — | ⚠️ belum — jalur dokumen Claim Life ada di `main`, belum disambung ke layar ini |
| `DownloadDocumentClaim` b25128/b25366/b25646/b25804, `View Office Online` b25629 | — | ⚠️ belum |
| `SetRemarkKomiteLife` b32664/b32800, `Remarks` (`AdjustmentDetail.DataCommitteeTreaty.Remarks`) | — | ⚠️ belum |
| syarat tampil `IsTreatyIn==0` ×12, `Type TP‖TR` ×3, `SwiftCode`, `RetrocadedShare` | — | ⚠️ blok rincian belum dibawa |
| bagian bergerbang `.TransferType==2` b29177 | — | ➖ mati (tiket 02) |

## Aksi di balik `Submit` — `KomitePostAdjustment`

| Langkah | Sistem baru | Tiket |
| --- | --- | --- |
| 2–3, 13 (tingkat, anak tangga, `KomiteCount + 1`) | `CatatKeputusan` (dua tulisan bersyarat) | 02 |
| 4 Approve Last Komite: 4.7–4.12 nomor, 4.15 Insert ke OS | `penyelesaiAkhirOracle.Akseptasi` + `RekamAkhirWarisan` | 04a, 04b |
| 4.17/4.18 `PrintAkseptasiPDF` | — | ⚠️ 04b belum |
| 5 Reject: 5.3, 5.6 | `penyelesaiAkhirOracle.Tolak` | 05 |
| 5.1 timpa seluruh tangga | — | ⛔ tidak ditiru (OQ-K-05) |
| 8 `InsertJsonClaimLife_Act` | — | ⛔ JSON dibuang |
| 9 EXIT retro | — | ⛔ dibuang (OQ-064) |
| 10 Arasapas · 11 email · 12 Kasir | outbox `T_LOG_SERVICE_RNM` (`KOMITELIFE`) + `PelaksanaKomite` | 06, 07 |
| 14 `SetInformationData` | — | ➖ bukan efek keluar (keputusan work owner) |

## Aksi baru tanpa padanan korpus

| Aksi | Rute | Tiket |
| --- | --- | --- |
| eskalasi naik satu tingkat (admin) | `POST /api/komite/{id}/eskalasi` | 03 |
| efek keluar per kasus + laporan harian "perlu intervensi" | `GET /api/komite/{id}` (`efek`), `GET /api/komite/laporan-harian` | 08 |
| riwayat tangga untuk siapa pun | `GET /api/komite/{id}/riwayat` | 09 |
