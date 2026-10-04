# 01: Modul Master Data — rencana (insert, update, aktif / nonaktif per tabel master)

> ⚠️ **Disusun agent sesi `nusantarare-0f` atas permintaan work owner 04-10-2026 — bukan hasil `/to-tickets`.**
> Kutipan: "nanti di sistem baru itu, biar data nya ke update km buatin menu master dengan nama masing2 tabel tadi itu.
> misal manu master province, isi nya untu insert, update, aktif, non aktif kan data dari master itu."
> Jawaban AskUserQuestion: penempatan **Modul baru 'masterdata'**; cakupan **Yang dipakai NB FacIn dulu**.

**Status:** rencana. Backend + pendaftaran modul (inti) = sesi c3; frontend = sesi 0f.

## Latar

View warisan PROVINCE / ACCUMULATEDTYPE / CZONE / ACCUMULATION sudah diganti tabel flat (nbfacin migrasi 196, tiket
46); CITY / DISTRICT / V_JN_OBJ_ITEM menunggu DDL tabel sumbernya (CITYINPUT, DISTRICTINPUT, RWINPUT, BRANCH,
OBJECTITEMTYPE). Tabel flat = salinan saat migrasi; agar datanya tetap terbarui, aplikasi baru menyediakan layar master.

## Cakupan tahap 1 (tabel yang dipakai NB FacIn)

| Master | Tabel | Status aktif sekarang | Catatan |
| --- | --- | --- | --- |
| Nation | `NATION` (tabel) | tidak ada | ID, OLDID, NOTE, NATIONINITIAL |
| Province | `PROVINCE` (flat 196) | tidak ada | NATIONNAME = turunan NATION |
| City | `CITYINPUT` (view JSON `m_city` → flat) | tidak ada | `CITY` tetap view gabungan (M-8) |
| District | `DISTRICTINPUT` (view JSON `M_DISTRICT` → flat) | tidak ada | `DISTRICT` tetap view gabungan (M-8) |
| CZone | `CZONE` (flat 196) | tidak ada | GROUPOFNAME turunan; TGLUPDATE / USERID = jejak |
| Accumulated Type | `ACCUMULATEDTYPE` (flat 196) | tidak ada | |
| Accumulation | `ACCUMULATION` (flat 196) | baris nonaktif tidak tersalin (view `IsActive IS NULL`) | ACCUMULATIONNAME turunan |
| Object Item Type | `OBJECTITEMTYPE` (tabel) | `ISACTIVE` VARCHAR2(10) | `V_JN_OBJ_ITEM` tetap view atasnya (usul, dinilai c3) |
| (rujukan) | `BRANCH` (tabel) | `STATUS` / `BRANCHSTATUS` | dibaca untuk BRANCHID City |

Tahap berikut: BRANDDETAIL, VJ_M_TYPE_PROPERTY_PLAN, COVERAGE_FACIN (A179 → pemilik = modul ini).

## Keputusan agent (menunggu konfirmasi)

- **M-1** Modul `masterdata`, GROUPMENU `MASTER`, migrasi `760-799`, slot menu `990-991`, prefix `/api/masterdata`.
  Satu menu "Master Data" (menu datar, 1 modul 1 menu); di dalam halaman, satu tab / daftar per tabel master.
- **M-2** Kepemilikan tabel master pindah ke modul ini; nbfacin membacanya sebagai tabel warisan baca-saja (pola
  `TabelCoverage`). Migrasi flat yang belum dijalankan dipertimbangkan pindah ke rentang 760 (keputusan c3 + penjaga).
- **M-3** Aktif / nonaktif: tabel tanpa kolom status mendapat `STS_AKTIF VARCHAR2(1) DEFAULT '1' NOT NULL` ('1' aktif,
  '0' nonaktif); tabel yang sudah punya kolom (V_JN_OBJ_ITEM `ISACTIVE`) memakai kolomnya sendiri sesudah nilainya
  diverifikasi. Tidak ada hapus — hanya nonaktif. Popup / saran di modul lain hanya menampilkan baris aktif.
- **M-4** ID baris baru: tabel kode (NATION, PROVINCE, CITY, DISTRICT, CZONE, ACCUMULATEDTYPE) = diisi pengguna, unik.
  ACCUMULATION = aturan `RDBMASTERACCUMULATION` (`COUNTRY-ZIPCODE-lpad(seq,6)`) diport ke Go (ADR-0043: tanpa
  prosedur).
- **M-5** Kolom turunan (NATIONNAME, GROUPOFNAME, ACCUMULATIONNAME, PROVINCENAME, CITYNAME, …) diisi backend dari
  tabel rujukannya saat simpan — pengguna memilih rujukan, bukan mengetik namanya.
- **M-6** Akses lewat Kelola User (gerbang menu `M_LOGIN_GO_MENU`), tanpa aturan peran baru.
- **M-7** Jejak ubah (`inti/backend/jejak`) untuk setiap insert / update / ubah status.
- **M-8** (keputusan work owner 04-10-2026, AskUserQuestion: "Flat-kan CITYINPUT & DISTRICTINPUT saja") — view
  `CITY` / `DISTRICT` adalah gabungan CITYINPUT + DISTRICTINPUT + RWINPUT + PROVINCE + BRANCH (satu baris per kode pos);
  yang di-flat-kan dan diedit = `CITYINPUT` / `DISTRICTINPUT`, sedangkan `CITY` / `DISTRICT` tetap view sehingga
  perubahan master langsung terlihat. `RWINPUT` tidak disentuh (work owner: "rwinput tidak usah").

## Kontrak API (usulan, ditetapkan c3)

- `GET /api/masterdata/{tabel}?q=&status=&halaman=` → `{ baris: [...], total }` (q = Contains tanpa peka huruf atas
  kolom ID / nama).
- `POST /api/masterdata/{tabel}` → `{ id }`; `PUT /api/masterdata/{tabel}/{id}`; `PUT /api/masterdata/{tabel}/{id}/status`
  badan `{ aktif: boolean }`.
- `{tabel}` = `nation` · `province` · `city` · `district` · `czone` · `accumulatedtype` · `accumulation` · `objectitemtype`.

## Langkah

1. Pendaftaran modul ke-21 (inti): baris `M_NAV_MENU` (migrasi inti 900-949), `frontend/katalogKorpus.ts`, uji perakit
   (`daftar.datar`, `Beranda`, `Shell`) — tanpa ini 6 uji perakit merah begitu folder `modul/masterdata` ada.
2. Backend kerangka (`modul.go`, `go generate ./inti/backend/daftar`), migrasi status aktif, endpoint per tabel.
3. Frontend: halaman Master Data (template Kelola User: kepala, toolbar, tabel) + form tambah / ubah + tombol aktif /
   nonaktif.
4. Penyaring aktif di pembaca nbfacin (saran / popup).

## Menunggu work owner

- ~~DDL CITYINPUT, DISTRICTINPUT, BRANCH, OBJECTITEMTYPE~~ — dikirim 04-10-2026 (RWINPUT tidak diperlukan).
- Konfirmasi M-1 … M-7.
