# PROMPT — MODUL **Komite Claim Life** *(sesi 3, cabang `modul/komite-claim-life`, worktree `.worktrees\komite-claim-life`)*: Bagian B seluruhnya dalam satu giliran, dari XML

> Baca dulu `PROMPT-INDUK-TIGA-MODUL.md`. Aturan modul Claim Life berlaku: brief modul §1–§2,
> mekanisme giliran lanjutan 8 §1, larangan keamanan, bab TELEMETRI EKSEKUSI. Bekal Bagian B dari
> lanjutan 4 §4 dan lanjutan 5 §5 **tetap berlaku** dan dirangkum di bab 0.

---

## 0. KONTEKS YANG SUDAH ADA

| Sumber | Isi |
| --- | --- |
| `.scratch/komite-claim-life/spec.md` | tangga persetujuan dari **roster** *(jumlah tingkat = roster, bukan struktur proses)*; wewenang di services *(ADR-0014)* dengan satu pengecualian: eskalasi manual naik satu tingkat; efek keluar wajib berhasil lewat outbox *(ADR-0015; Kasir memindahkan uang)*; satu jalur simpan berparameter status; kode mati dibuang *(`TransferType` + `ShowTransfer`, gerbang EXIT retro, lima langkah ter-remark)*; batas dengan Claim Life: **masuk** = `serahkanKomite` *(child work `ASM-FW-GCNMFW-Work-KomiteLife`, `KomiteLoop`, `KomiteList`, penunjuk baris, nilai, `CURRENCY`)*, **keluar** = `STS_REJECT` `1`/`2` pada tingkat akhir *(tiket Claim Life 11)*; `OS_AKSEPTASI_KLAIM_LIFE` ditulis rule yang sama dari dua sisi |
| `STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md` + `keputusan-struktur-komite.md` | `T_WORK_CLAIM` baris komite `KMTLF-` *(ralat A2)*, `COVER_KEY` → klaim induk; `T_GENERAL_KOMITE` *(`ID` = shared PK, `ADJUSTMENT_ID`)*; `T_KOMITE_KOMITELIST` *(FK `DATA_KOMITE_ID`)*; **tidak ada** `WORK_CLAIM_ID`/`KMT_NO` |
| Sudah ada di `main` *(Claim Life A2, migrasi `013`)* | `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, `SEQ_KOMITE_KOMITELIST`; `KasusKomite` + `Roster` *(af; roster `EMAILKOMITE` 18 kolom — **data orang, nol baris disalin**)*; outbox `T_EFEK_KELUAR` *(aq)*; resolver `M_LINK_SERVICE` *(`Kasir/*`, `SendEmail/*`)*; `PenomorCounter` *(o1–o3)*; sensus `KomitePostAdjustment` *(`1` ×3, `2` ×1 pada baris **dan** peserta; `ACCEPTATION_DATE` ×3; gerbang `KomiteCount == KomiteLoop`)*; decision table 1 berkas → tabel data di `models` |
| `issues/` | `00` skema *(verifikasi + penjaga + yang kurang — **bukan** membuat ulang `013`)* → `01` terima kasus + inbox per posisi → `02` mesin tangga → `03` wewenang + eskalasi → `04a` nomor akseptasi sekali di tingkat final → `04b` rekam akseptasi satu jalur → `05` jalur balik `STS_REJECT` dua tingkat → `06` outbox → `07` worker pengirim retry anti-dobel → `08` status perlu intervensi + laporan harian → `09` riwayat tangga di UI |
| ⚠️ Nomor akseptasi Komite | `SUMBER-PENOMORAN-DBA.md`: pasangan **counter procedure** `(ASM-FW-GCNMFW-Work-KomiteLife, RNML-A)`, sedangkan Claim Life memakai sequence `ACCEPTATIONNOLIFE_SEQ` *(ap)*. Keputusan **o** berlaku: prosedur **tidak dipanggil**, counter ditiru lewat `PenomorCounter` *(o2/o3)*; keduanya harus **tidak** bertabrakan nomornya — baca `GetAcceptedNoCL` dan uji keunikan lintas jalur; bila XML tidak menjawab, `[terbuka — work owner]` dengan bukti, bukan tebakan |

## 1. XML — TITIK MASUK DAN MENU *(korpus `D:\XML\RNM_BRD\Komite Claim Life\`, dibaca 27 September 2026)*

| # | Bukti | Isi |
| ---: | --- | --- |
| 1 | `Struktur_KomiteClaimLife.xlsx` r2–r4, r65 | root = Flow `KomiteLife_Flow` "(root / entry point)"; `KomiteRouter` *(Activity, "Router assignment", Assignment1)*; `ViewTransferDtl` *(Flow Action assignment)*; `IsKomiteLoop` *(When, Decision1 KomiteLoop)*. **Tidak ada** folder `Harness\` → tidak ada portal modul |
| 2 | `Flow\KomiteLife_Flow.xml` | `pyWorkTypeName KomiteLife` 220; assignment `pyImplementation WorkList` 594, `pyRouteTo Custom` 619, router `KomiteRouter` 662; task `KomiteRouter` 233, `ViewTransferDtl` 248, `KomiteLoop` 218, `IsKomiteLoop` 263; ⚠️ `pyWorkTypeName KomiteTreaty` 250 di sekitar `ViewTransferDtl` — baca pohonnya: rule salinan dari modul Komite lain atau bukan |
| 3 | `FlowAction\` · `Section\` | `ViewTransferDtl` *(layar keputusan)*, `AttachDocumentLife`, `RetroLife`; section `ShowTransfer` *(mati — `TransferType` tidak pernah diisi)*, `RetroDetailLife`, `AttachDocScreenLife` |
| 4 | `Activity\` *(18)* | antara lain `KomitePostAdjustment`, `KomiteRouter`, `SetRemarkKomiteLife`, `DownloadDocumentClaim`, `HitServiceToKasirKMTLife_Act` — tiap yang **menulis** ditelusuri dari metode langkahnya *(PARITAS Koreksi 3)* |
| 5 | `ReportDefinition\` | `GetEmailUser_RD`, `GetAllAttachmentsLife`, `BrowseDiseaseLife_RD` |

**Menu `[dari bukti 1–2]`**: kelompok **Komite Claim Life** → satu butir **`Inbox Komite`**
`[tidak ada di korpus — worklist `KomiteRouter`]` berisi kasus pada **tingkat yang pelaku miliki**
*(roster)*; kolom Pega-standar seperti `InboxPremiumList` *(dipinjam, ditandai)* + kolom komite
*(nomor klaim induk, tingkat/`KomiteLoop`, nilai + `CURRENCY`)*. Layar keputusan
`ViewTransferDtl` dan lampiran dibuka **dari baris**, bukan dari menu. `ShowTransfer` **tidak** dibawa.

## 2. KEPUTUSAN YANG MENGIKAT GILIRAN INI

| Butir | Isi |
| --- | --- |
| **km1** | migrasi hanya `030`+ untuk yang **kurang** dari `013` *(mis. kolom status tangga, indeks, outbox komite bila `T_EFEK_KELUAR` tidak cukup — dibuktikan dulu)*; tiket 00 = verifikasi bentuk `013` terhadap STRUKTUR-TABEL + penjaga; tidak ada CREATE ulang |
| **km2** | peran anggota komite = **roster** *(`EMAILKOMITE` → `Roster` A2)*, bukan `X-Peran` Claim Life; identitas stub dari env tetap *(`VITE_STUB_PELAKU` = akun roster `UJI-*` di fixture)*; wewenang ditegakkan di services *(tiket 03)*; eskalasi manual satu tingkat — satu-satunya pengecualian, diuji |
| **km3** | `STS_REJECT` ke baris Claim Life hanya pada **tingkat akhir** *(tiket 05 ↔ Claim Life 11)*, lewat fungsi yang **sudah ada** di `services` Claim Life *(dipanggil, tidak disalin)*; bila fungsinya belum ada, tambah **aditif** dan laporkan |
| **km4** | Kasir dan email = outbox `T_EFEK_KELUAR` + worker *(tiket 06–07)*; endpoint sungguhan **tidak** dipanggil di DEV tanpa persetujuan manusia — pengirim stub mencatat; retry/anti-dobel diuji dengan kunci idempoten |
| **km5** | riwayat tangga *(tiket 09)* dibaca dari `T_KOMITE_KOMITELIST` + jejak; status "perlu intervensi" *(tiket 08)* adalah kata, kode tetap di `models` |

## 3. URUTAN KERJA — satu giliran, satu commit per tiket `komite-claim-life: tiket NN — <judul>`

`00 → 01 → 02 → 03 → 04a → 04b → 05 → 06 → 07 → 08 → 09`. Tiap tiket: bab **Pembacaan ulang XML**
*(activity sebagai pohon; `KomitePostAdjustment` dan `KomiteRouter` utuh)*, ralat bila tiket
menyimpang, uji murni + handler + `db` *(SKIP)*; frontend di Shell yang sama: **Inbox Komite**,
**layar keputusan** *(`ViewTransferDtl` sebagai pohon: `KomiteAproval`, `KomiteComment`, tombol
VERBATIM + baris; dropdown keputusan **wajib** — spec §"Problem": tidak ada rule yang menulisnya)*,
**riwayat tangga**; `labels.komite.ts` berbukti; `PARITAS-LAYAR-DAN-AKSI.md` milik modul ini.

Batas dengan sesi Claim Life yang berjalan bersamaan: berkas `komite_*` milik sesi ini; berkas Claim
Life **tidak** disunting; kontrak masuk/keluar persis A2 — bila kontrak perlu berubah, tulis usulan di
`.scratch/komite-claim-life/KONTRAK-USULAN.md` dan laporkan, jangan mengubah sepihak.

## 4. LANGKAH 0 · LAPORAN · TELEMETRI

Langkah 0: worktree + `.env` sesuai induk §1 *(`:8083`, Vite `5175`)*; `git log -1` ≥ `cfc4824`;
uji hijau. Laporan akhir *(satu pesan)*: tiket ter-commit, ralat, tabel **menu → bukti → halaman →
rute**, kontrak yang dipakai/diusulkan, kode bersama aditif, bab TELEMETRI per tiket.

---

*Disusun 27 September 2026 dari spec, struktur tabel, keputusan struktur, sebelas tiket, sheet
struktur, `KomiteLife_Flow.xml` (routing), daftar flow action/section/activity, dan bekal A2 di
`main` (roster, kasus, outbox, resolver, penomoran).*
