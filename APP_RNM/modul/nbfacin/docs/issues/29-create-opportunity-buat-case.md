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

- **Penyimpanan: tunggu tabel flat (tiket 23).** Data opportunity dan case NB baru disimpan sesudah tabel flat
  (`T_WORK_POLIS` dkk.) dibuat — tiket 23 masih menunggu keputusan presisi tim inti. Tidak ada tabel sementara.
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

- Endpoint `POST /api/nbfacin/opportunity` (backend sesi c3): menunggu tabel flat (tiket 23) + nilai awal
  penghitung NB. Sampai itu, tombol menampilkan galat backend apa adanya.
- "Lari ke flow": port section `InputInwardFacultative` — tiket berikutnya; sekarang sesudah case dibuat, nomor
  case ditampilkan di form.

**Status:** blocked — frontend selesai 03-10-2026; backend menunggu tiket 23 dan nilai awal penghitung NB

## Comments
