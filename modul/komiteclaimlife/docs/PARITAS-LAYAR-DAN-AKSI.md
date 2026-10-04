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
| 5.1 timpa seluruh tangga | `penyelesaiAkhirOracle.Tolak` → `TanggaSebelumDitimpa`, `jejakTimpaTangga`, `TimpaTanggaTolakAkhir` | ✅ **ditiru — OQ-K-05 ditutup 29-09-2026 (GILIRAN-17)**; keputusan lama dijejaki dulu, tingkat dilewati eskalasi tidak ditimpa |
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

## Sensus remark 28-09-2026 (GILIRAN-12 paket 2)

Setiap activity Komite Claim Life dicetak beserta `pyStepsBlockName`; `//` = ter-remark (beserta
anaknya). 18 activity, **7** ber-remark, **23** langkah `//` (label `EXIT` bukan remark). Nol langkah
mati yang ditiru di kode.

**Cara hitung** `[terverifikasi]` — dua cara, jendela = seluruh `Komite Claim Life/Activity/*.xml`:

| Cara | Perintah | activity | ber-remark | langkah `//` |
| --- | --- | ---: | ---: | ---: |
| A — pengurai XML (expat), per `rowdata` ber-`pyStepPageReference` yang `pyStepsBlockName` tepat `//` | skrip sensus sesi (cetak activity, nomor langkah, bNNN) | 18 | 7 | 23 |
| B — grep mentah per berkas | `grep -o "<pyStepsBlockName>//</pyStepsBlockName>" "<berkas>" \| wc -l`, dijumlah | 18 | 7 | 23 |

Alat diuji lebih dulu pada item yang sudah diketahui: `Claim Life/Activity/SaveOutStandingLife_Act.xml`
— 8 langkah `//` (dibaca ulang pemeriksa independen GILIRAN-11) — kedua cara menjawab 8. Label
lompatan lain (`EXIT`, `send`, `AA`, …) bukan remark dan tidak terhitung.

⚠️ **Cakupan celah langkah hidup DINYATAKAN:** sensus ini mendaftar langkah hidup yang tidak ditiru
untuk activity **ber-remark** saja. Activity tanpa remark tidak diaudit ulang langkah demi langkah di
giliran ini; statusnya tetap seperti baris paritasnya (bila ada) — **belum** audit penuh.


| Activity | Langkah `//` (b) | Ditiru? | Keputusan |
| --- | --- | --- | --- |
| `HitServiceToKasirKMTLife_Act` | 11 b3033 (`Connect-REST`, satu-satunya panggilan keluar) | tidak — pengirim Kasir stub `ErrKasirBelumDisetujui` | ⛔ **Sistem lama tidak pernah memanggil Kasir** dari jalur ini. ADR-0015, spec, CONTEXT, tiket 07/08 diralat buktinya; **OQ-K-06 ditutup 29-09-2026 (GILIRAN-17): Kasir tetap stub, sambungan nyata = fitur baru** |
| `KomitePostAdjustment` | 4.4 b1725 · 4.5 b1943 (nomor lewat RDB) · 4.6 b2160 · 4.14 b4947 · 5.5 b7673 (tukar retro) | tidak — nomor 4.7–4.12 (`models/komite_nomor.go`), stempel 4.11/4.12 | tercatat sejak tiket 04a; 4.6 isinya sama dengan 4.11/4.12 yang hidup; tukar → keputusan ronde 1 #4 |
| `LoadDocumentKomiteLife_ACT` | 1 b245 (+1.1–1.8) | tidak | — |
| `LoadDocumentLife_ACT` | 1.4 b913 (+1.4.1, 1.4.2) | tidak | dicatat di PARITAS Claim Life |
| `SaveAttachLife` | 1.1, 1.3–1.6, 2, 3 | tidak (lampiran Komite belum tersambung, baris 23) | — |
| `SendEmailKlaimLife` | 1 b365 · 2 b499 (CC) | tidak | — |
| `SetRemarkKomiteLife` | 3 b539 `Obj-Save` | tidak | langkah 1–2 hidup = celah tercatat (baris 25) |

**Langkah HIDUP yang tidak ditiru dan belum tercatat** (kini dicatat sebagai celah):

- `HitServiceToKasirKMTLife_Act` langkah 13 b3395 — `INSERT POOLDATA.DIRECTTOKASIR_LOG` (JSON b2774,
  `KET` kosong), bergerbang `KomitePostAdjustment` 12 (final Setuju, bukan TP/TR, KPR, produksi).
- `SendEmailKlaimLife` — cabang penerima: anggota berikutnya (b1882), pembuat saat setuju (b2496) dan
  tolak (b2706), akun syariah (b3228); CC hidup langkah 3 (b704–705), BCC tetap (b3540). `EfekEmail`
  masih stub.
- Daftar dokumen di layar keputusan `ShowTransfer` (b22515 → `LoadDocumentKomiteLife` langkah 2
  b1572 → `LoadDocumentLife` 1.1–1.3; grid b24240) — baris 23/24 hanya menyebut Add attachment dan
  unduhan.
- `KomitePostAdjustment` 4.11/4.12 juga menyetel `IsCheck = "true"` (b3282/b3560); aplikasi tidak
  menulisnya — nilainya sudah `"true"` sejak pendaftaran (`models.PenandaDipilih`), jadi tanpa akibat.

