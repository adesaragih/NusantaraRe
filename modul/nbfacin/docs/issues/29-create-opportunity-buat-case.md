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

## Backend (sesi c3, selesai 03-10-2026 — diverifikasi sesi 0f)

- `POST /api/nbfacin/opportunity` → 201 `{"caseId":"NB-<n>"}`; 400 medan wajib/tanggal DD-MM-YYYY/panjang kolom;
  **401 tanpa identitas login** (tambahan atas kontrak awal); 503 tanpa DB; 500 tanpa rincian. Ke-14 nama medan
  `IsianOpportunity` cocok dengan `api.ts`.
- Satu transaksi: nomor dari `SEQ_WORK_POLIS_NB`; baris `T_WORK_POLIS` yang ADA (ID `NB-<n>`, `LINI = 'FAC'`,
  **`POSITION = 'Offer'`, `FLAG_ONGOING_POLICY = '0'`, `STATUS_WORK = 'Pending-Policy'`** — butir 77, menggantikan A80
  "kosong"; pembuat = pengenal akun, `TGL_CREATE/TGL_UPDATE`); isian ke tabel sendiri **`T_NB_OPPORTUNITY`** (opportunity = work class
  tersendiri di Pega, `GetListOpportunity.xml`, ditaut lewat NBHandle).
- Keputusan work owner yang ditanyakan langsung di sesi c3 (03-10-2026, register butir 76): ID ikut tabel yang ada
  (VARCHAR2(32)), LINI `FAC`, tabel opportunity sendiri, empat kolom rancangan digabung ke kolom yang ada,
  placeholder `{NB_MULAI}` di migrasi.

## Yang harus dijalankan work owner (Oracle DEV)

✅ **Sudah dijalankan work owner di DEV 03-10-2026** (diteruskan sesi 0f): 180 dan 181 berjalan, case NB berhasil dibuat;
181 diisi work owner `START WITH 151028` (⚠️ lebih kecil dari contoh `NB-184351` — register butir 77).

1. `180_t_nb_opportunity.sql` — membuat `T_NB_OPPORTUNITY`.
2. `181_seq_work_polis_nb.sql` — ⛔ **ganti `{NB_MULAI}` dengan nomor NB Pega terakhir + 1 SEBELUM `-migrate`.**
   Runner hanya mengganti `{skema}`; selama `{NB_MULAI}` masih ada, Oracle menolak dan `-migrate` SELURUH aplikasi
   berhenti di 181 (pilihan work owner). Angka produksi tidak ada di repo (NB-184351 hanya contoh data uji).

## Risiko terbuka

- **R1 `[terverifikasi]`, MELEBAR sejak butir 77:** kotak masuk PremiumList Life membaca SELURUH `T_WORK_POLIS` tanpa
  saringan `LINI` (`modul/premiumlistlife/backend/repository/polis_inbox.go` L104 `SELECT w.ID, w.STATUS_WORK,
  w.POSITION`, L111 `WHERE (:1 IS NULL OR w.POSITION = :2)`, L125 cacah `w.POSITION = :1`; posisi dari kueri `?posisi=`,
  `handlers/rute_premiumlist.go` L53). Case NB kini ber-`POSITION = 'Offer'` — **nilai posisi Life yang sama**
  (`models/polis_penawaran.go` L112 `PosisiOffer = "Offer"`; kasus Life baru juga `Offer`, `models/polis_kasus.go` L84).
  Jadi case NB ikut tampil di tab "semua" **dan** di saringan `posisi=Offer`, tampak seperti kasus Life ber-status
  `Pending-Policy`. Layar Life kini hanya memanggil tanpa posisi (`frontend/pages/InboxPremiumList.tsx` L99 →
  `frontend/api.ts` L44-50, `posisi = ''`), jadi yang terlihat hari ini tetap tab "semua"; panggilan/tab `Offer` mana pun
  ikut kena. `Pending-Policy`: **nol** kemunculan di premiumlistlife (status Life: `Input Offer Life` / `Input Premium
  Detail` / `Input Premium Summary`, `models/polis_penawaran.go` L76/L78/L91) — tidak menjadi saringan, hanya tampil.
  Perbaikannya di modul premiumlistlife (pemiliknya); tidak disunting.
- **R2:** nbfacin menulis tabel milik premiumlistlife dengan SQL sendiri, tanpa kontrak lintas modul.
- A79, A81–A85 (tanpa FK, panjang kolom, 401, tanpa pangkas, sequence vs tabel penghitung ADR-0043,
  tanpa jejak audit) menunggu konfirmasi work owner.
- "Lari ke flow": port section `InputInwardFacultative` — tiket berikutnya; sesudah case dibuat, nomor case
  ditampilkan di form.

**Status:** ready-for-human — frontend + backend selesai 03-10-2026; migrasi 180/181 sudah dijalankan work owner di DEV; keadaan awal case butir 77

## Comments
