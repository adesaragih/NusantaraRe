# 29: Tombol Create opportunity di form Opportunity — membuat case NB

> Disusun agent (sesi `nusantarare-0f`) atas permintaan work owner 03-10-2026 — **bukan hasil `/to-tickets`**:
> *"tambah tombol create opportunity nya di kanan atas. saat di klik tombolnya akan lari ke flow nya dan
> membuat caseid, contohnya NB-12123"*.

**What to build:** tombol `Create opportunity` di kanan atas form Opportunity (tiket 26). Saat diklik: medan
wajib diperiksa; bila lengkap, isian dikirim ke backend yang **membuat case NB** dan mengembalikan nomornya
(`NB-<angka>`); lalu layar berlanjut ke flow NB.

## Bukti

- `[terverifikasi]` Label tombol `Create opportunity` = `SFAPortalOpportunitiesHeader.xml` sel 72, L2486 (sama
  dengan tombol portal, tiket 25).
- `[terverifikasi]` Layar pertama flow NB = flow action `InwardFacultative` → section **`InputInwardFacultative`**
  (`D:\migrasi\RNM\NB FacIn\FlowAction\InwardFacultative.xml`, `pySectionReference`), flow
  `InputInwardFacultativeOffer` (`ASM-FW-GISFW-WORK!INPUTINWARDFACULTATIVEOFFER`).
- `[terverifikasi]` Nomor NB di data nyata = satu deret global tanpa tahun: nama berkas ekspor kasus di
  `D:\XML NURE\Groupbusiness\` sampai **NB-184351**.
- `[terverifikasi]` ADR-0043: penomoran dijalankan aplikasi, penghitung di Oracle lewat `SELECT … FOR UPDATE`.

## Keputusan work owner 03-10-2026

- ~~**Penyimpanan: tunggu tabel flat (tiket 23).**~~ ⛔ **DIRALAT 03-10-2026** (sesi 0f): pilihan itu diambil
  work owner karena sesi 0f keliru menyebut penyimpanan case seluruhnya tertahan tiket 23. Faktanya
  `[terverifikasi]` **`T_WORK_POLIS` sudah ada** — dibuat modul premiumlistlife (`050_t_work_polis.sql`, diubah
  057/058/059/063) — dan register K-064 (butir 1b) memutuskan Fac In **menyambung** ke tabel itu (pembeda: kolom
  penanda lini). **Keputusan pengganti (work owner "setuju", 03-10-2026):** case dibuat sekarang di `T_WORK_POLIS`
  yang ada (ID `NB-<n>`, sequence NB sendiri), isian opportunity ke tabel flat pasangannya yang dibuat sekarang;
  tiket 23 diselaraskan dengan `T_WORK_POLIS` yang ada (skema loader dari DDL draf bentrok: ID NUMBER vs
  VARCHAR2(32), POSISI vs POSITION, dst.). Dikerjakan sesi c3; migrasi ditulis, dijalankan work owner di DEV.
- **Nomor NB: lanjut dari nomor terakhir Pega** (MAX nomor NB yang ada + 1), supaya tidak bentrok dengan kasus
  lama. Angka pastinya diisi saat migrasi dijalankan (work owner/DBA).

## Yang sudah dibangun (frontend)

- Tombol di kanan atas, sebaris judul; pemeriksaan medan wajib urut layar (Estimated Closing Date, Business
  Prospect Name, Class Of Business, Type Of Inward, Type Of Facultative bila tampil, Phase, Business Status) —
  yang kosong disebut, tanpa permintaan ke backend.
- Kontrak `POST /api/nbfacin/opportunity` (`api.ts` `IsianOpportunity` → `{"caseId":"NB-…"}`); tanggal dalam
  bentuk kabel inti `DD-MM-YYYY`. Sesudah berhasil: nomor case ditampilkan dan tombol nonaktif.
- **Estimated Closing Date berformat dd/mm/yyyy** (permintaan work owner 03-10-2026): kotak teks dd/mm/yyyy +
  tombol kalender (`components/TanggalDMY.tsx`) — `<input type="date">` bawaan mengikuti bahasa browser.

## Tertahan

- Endpoint `POST /api/nbfacin/opportunity` (backend sesi c3): sedang dikerjakan; migrasinya menunggu dijalankan
  work owner di DEV, dengan nilai awal penghitung NB (nomor NB Pega terakhir + 1) yang diisi work owner/DBA.
  Sampai itu, tombol menampilkan galat backend apa adanya.
- "Lari ke flow": port section `InputInwardFacultative` — tiket berikutnya; sekarang sesudah case dibuat, nomor
  case ditampilkan di form.

**Status:** in-progress — frontend selesai 03-10-2026; backend dikerjakan sesi c3 (T_WORK_POLIS yang ada)

## Comments
