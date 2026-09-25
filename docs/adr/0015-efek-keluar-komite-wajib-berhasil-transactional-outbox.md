---
status: accepted
tanggal: 2026-09-15
sumber: grilling Komite Claim Life Ronde 1 + Ronde 2 (`.scratch/komite-claim-life/`), keputusan work owner
menyimpang-dari: ADR-0008 (berlaku untuk Claim — Life)
---

# Efek keluar Komite Claim Life **wajib berhasil** — transactional outbox, at-least-once

Keempat efek keluar keputusan komite **wajib berhasil**. Kegagalan **tidak boleh diam**: keputusan
komite **tidak dianggap tuntas** sampai seluruh efek berhasil terkirim, atau ditandai **perlu
intervensi**.

Jaminan: **at-least-once**. Pola: **transactional outbox** — keputusan dan daftar efeknya disimpan
dalam **satu transaksi database**; worker terpisah mengirim tiap efek dengan retry sampai sukses.

> ⚠️ **Ini MENYIMPANG dari ADR-0008**, yang untuk **Claim — Life** memutuskan efek keluar bersifat
> asinkron dan **tidak memblokir** alur. Kedua ADR berlaku bersamaan pada konteks yang berbeda —
> lihat §"Mengapa Komite berbeda dari Claim".

## Keempat efek keluar `[terverifikasi]`

`Komite Claim Life/Activity/KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, versi **2026-09-15**,
515.675 byte). Nomor langkah dibaca dari `<pyStepPageReference>RH_1.pySteps(n)` — yang muncul
**setelah** `<pyStepsActivityName>` milik langkahnya:

| Step | Efek | Catatan |
| ---: | --- | --- |
| 6 | `Obj-Save` | **persist** — seluruh efek berjalan sesudahnya |
| 8 | `Call InsertJsonClaimLife_Act` | |
| 9 | *(gerbang EXIT retro)* | dibuang di sistem baru — lihat §Interaksi |
| 10 | `Call serviceInsertArasapasClaimLife_act` | endpoint di-lookup runtime dari `M_LINK_SERVICE`, kunci `Klaim` / `insertClaimLife` (**ADR-0013**) |
| 11 | `Call SendEmailKlaimLife` | |
| 12 | `Call HitServiceToKasirKMTLife_Act` | **Kasir — pembayaran.** `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `HITSERVICETOKASIRKMTLIFE_ACT` |
| 14 | `Call SetInformationData` | **di luar lingkup** — temporary, lihat §Lingkup final |

`[terverifikasi]` Step 10 dan step 12 digerbangi `AcceptStatus = 1 && KomiteCount == KomiteLoop`
(baris **8648** dan **8887**) — keduanya hanya berjalan pada **keputusan aksep di tingkat final**.

## Mengapa Komite berbeda dari Claim

**Pemicunya Kasir.** `[terverifikasi]` `HitServiceToKasirKMTLife_Act` **tidak ada di Claim Life** —
Claim Life punya empat efek keluar, Komite punya lima (termasuk `SetInformationData`).

Kasir adalah **integrasi keuangan**: ia memicu jalur pembayaran atas klaim yang baru saja disetujui
komite tingkat tertinggi. Kegagalan diam di sini berarti **klaim disetujui tetapi tidak pernah
sampai ke pembayaran** — dan tidak ada yang tahu.

Bandingkan Claim — Life (**ADR-0008**): di sana efek keluar adalah unggah berkas, email, Arasapas,
dan konversi — tidak satu pun memindahkan uang, dan `[terverifikasi]` Pega sendiri hanya
**mencatat** kegagalannya (`InsertLogServiceClaim` → `monitoring_klaim_log`) tanpa gerbang
keberhasilan. Di sana "tidak memblokir" adalah paritas yang masuk akal.

Di Komite, taruhannya berbeda. Karena itu jaminannya dinaikkan.

## Considered Options

- **Transactional outbox + retry, at-least-once** — dipilih
- Paritas dengan ADR-0008 (asinkron, tidak memblokir) — ditolak: kegagalan Kasir menjadi tak terlihat
- Panggil sinkron di dalam transaksi keputusan — ditolak: keputusan komite akan gagal hanya karena
  layanan luar sedang mati, dan transaksi database menggantung selama panggilan jaringan

## Consequences

- **Keputusan komite punya dua keadaan yang harus dibedakan**: *tersimpan* dan *tuntas*. UI dan
  laporan tidak boleh menyamakan keduanya.
- Diperlukan **keadaan "perlu intervensi"** untuk efek yang gagal berulang — beserta tempat orang
  melihatnya. Tanpa itu, "wajib berhasil" hanya menggeser kegagalan diam ke antrean.
- **At-least-once berarti penerima harus tahan duplikat.** Arasapas, Kasir, dan email dapat menerima
  kiriman yang sama lebih dari sekali. `[terbuka]` Apakah ketiganya idempoten **belum diketahui** —
  ini risiko nyata yang dibawa keputusan ini.
- Outbox disimpan dalam transaksi yang sama dengan keputusan → sejalan **OQ-013** (Go memegang batas
  transaksi untuk jalur Life).
- Retry harus membedakan **kunci kategori tidak ditemukan** di `M_LINK_SERVICE` (kesalahan
  konfigurasi data, tidak layak diulang berkali-kali) dari **jaringan gagal** (layak) — **ADR-0013**.

## Lingkup final dan syarat pengaman (2026-09-15)

### Lingkup: **empat** efek keluar `[terverifikasi]`

Dibaca ulang dari `KomitePostAdjustment.xml` versi **2026-09-15** (515.675 byte). Nomor langkah dari
`<pyStepPageReference>`, status dari `<pyStepsBlockName>` (kosong = aktif):

| Urutan | Step | Efek | Status |
| ---: | ---: | --- | --- |
| 1 | **8** | `Call InsertJsonClaimLife_Act` | AKTIF |
| 2 | **10** | `Call serviceInsertArasapasClaimLife_act` | AKTIF — endpoint via `M_LINK_SERVICE`, kunci `Klaim`/`insertClaimLife` (**ADR-0013**) |
| 3 | **11** | `Call SendEmailKlaimLife` | AKTIF |
| 4 | **12** | `Call HitServiceToKasirKMTLife_Act` | AKTIF — **Kasir/pembayaran** |

Seluruhnya **setelah `Obj-Save` (step 6)**.

**Di luar lingkup** `[keputusan work owner]`: step **14** `Call SetInformationData` — hanya
**temporary**, menampung hasil submit (`"Approve"`/`"Reject"` + `AcceptedNo` + `CLAIM_NO`). Bukan
efek keluar, bukan penyimpanan permanen, **tidak** wajib-berhasil.

### Syarat pengaman wajib `[keputusan work owner]`

Keempat efek **boleh di-retry berkali-kali**. Karena itu:

1. **ID idempoten unik per kiriman.** Tiap kiriman membawa identifier yang memungkinkan penerima
   **menolak duplikat**. Tanpa ini, at-least-once berarti duplikasi, bukan keandalan.
2. **Cek status sebelum kirim ulang — wajib untuk Email dan Kasir.** Sebelum retry, sistem
   memeriksa apakah kiriman sebelumnya **sudah sukses**. Email dobel mengganggu; **pembayaran dobel
   merugikan**.
3. **Status "perlu intervensi" terlihat.** Kegagalan yang tidak pulih setelah retry ditampilkan
   sebagai **status eksplisit di UI Komite** dan masuk **laporan harian**. Tanpa tempat orang
   melihatnya, "wajib berhasil" hanya memindahkan kegagalan diam ke antrean.

Butir 1 dan 2 **menjawab** risiko idempotensi yang dicatat di §Consequences: pertanyaannya bukan
lagi "apakah penerima idempoten" melainkan **"sistem baru wajib membuatnya aman"** — dengan ID unik
di sisi pengirim dan pemeriksaan status sebelum kirim ulang.

⚠️ Sisa risiko yang diterima: bila penerima **tidak** menghormati ID idempoten, butir 2 adalah
pertahanan terakhir — dan ia bergantung pada catatan status milik kita sendiri, bukan konfirmasi
penerima.

## Interaksi dengan pembuangan identitas retro ter-hardcode

⚠️ `[terverifikasi]` Di korpus versi 2026-09-15, **step 9 adalah langkah gerbang tersendiri** —
tanpa `pyStepsActivityName`, membawa `<pyStepsPreCondition>true` dan
`<pyStepsDescription>EXIT JIKA RETROID "L0000141"</pyStepsDescription>` (baris 8483), dengan
precondition `RetroID=="L0000141" || SecurityReinsurerID=="L0000134"` (baris 8525) dan
`RetroID=="1000013"` (baris 8548). Ia berdiri **di antara** efek 1 (step 8) dan efek 2 (step 10),
sehingga saat memicu, **step 10–14 tidak berjalan**: Arasapas, email, Kasir, dan
`SetInformationData` semuanya di-skip.

`[dugaan]` Bahwa kode transisi numeriknya berarti "Exit Activity" **tidak dapat dibuktikan dari
ekspor** — nilai yang terbaca (`2`) muncul juga di step 8 yang bukan EXIT. Yang terbukti adalah
label penulisnya dan bentuk langkahnya.

`[keputusan work owner]` Ketiga identitas itu **dibuang** (**OQ-064**). Akibatnya: klaim yang selama
ini dikecualikan akan **mulai menerima keempat efek keluar** — termasuk **Kasir**.

Digabung dengan ADR ini ("wajib berhasil"), perubahan itu berarti klaim-klaim tersebut kini
**menuntut keberhasilan pembayaran** yang sebelumnya tidak pernah dipicu. **Ini perubahan perilaku
yang menyentuh uang**, dan layak dikonfirmasi Product+UW sebelum rilis — lihat **OQ-064**.

## OQ yang masih terbuka dan menyentuh ADR ini
| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-064** | **TERTUTUP** 2026-09-15 — gerbang EXIT retro dibuang; semua klaim menjalankan keempat efek. Arti ketiga identitas tetap tidak diketahui dan tidak dicari lagi |
| **OQ-065** (baru) | Langkah "Tukar SecurityReinsurer dengan RetroName" (step 4.14 & 5.5) **ter-remark** — apa yang kini masuk ke kolom retro pada rekam akseptasi |
| **OQ-035** | `serviceInsertArasapasClaimLife_act` satu salinan dipakai dua konteks |
| **OQ-002** | Kontrak layanan **Kasir** tidak ada di korpus — bentuk permintaan, makna jawaban, dan apakah ia menghormati ID idempoten belum diketahui |
