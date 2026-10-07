# 34: Popup Change Ceding Co — daftar ceding (Add / Select, Delete, Submit) dari tabel AGENT

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026:
> *"tombol change ceding co, tampilannya seperti ini dan saat di klik add/select ceding, akan muncul tampilan sama
> seperti tampilan change sob, dan jika dipilih akan muncul ke tampilan yg saya kirim sekarang ini"*, disertai
> tangkapan layar popup *Ceding Co List* (tidak disalin ke repo).

**What to build:** tombol **Change Ceding Co** di layar Inward Facultative membuka popup **Ceding Co List**. Di
dalamnya, **Add / Select Ceding** membuka pencari AGENT yang sama dengan Change SOB; setiap Choose menambah satu
baris; **Delete** membuang baris; **Submit** menyerahkan daftar ke layar. *Ceding co name* = gabungan nama
berpemisah `;`; kode-kodenya tersimpan lewat *Save for later*.

**Blocked by:** tiket 33 (endpoint `GET /api/nbfacin/sob`, kolom AGENT terverifikasi DDL).

**Status:** frontend selesai 03-10-2026 (uji hijau); backend → sesi c3; migrasi **ditulis, tidak dijalankan**.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\Periode.xml` sel 53 `Change Ceding Co` → showHarness `ShowCedingCoList`, Target `popup`.
- `Harness\ShowCedingCoList.xml`: pyLabel `Ceding Co List`; tombol `Submit` → `SetCedingCo_Act` → `opener.location.reload`
  → `window.close`.
- `Section\ShowCedingCoList.xml`: grid atas `pyWorkPage.Quotation.CedingCoList`, 20 per halaman; kepala sel 14
  `Ceding Co` (wajib) + sel 15 `Add / Select Ceding` (runActivity `AddCedingList_act` → refreshList → showHarness
  `CedingCompany`, WindowName `Ceding Company`); baris sel 17 `.CedingCoName` (baca-saja, wajib) + sel 18 `Delete`
  (`DeleteCeding_Act`, parameter index / cedingconame / cedingco).
- `Section\CedingCoHierarki.xml` (isi harness `CedingCompany`): grid ber-Choose atas **`BrowseAgentNonLife_RD`** — RD
  yang sama dengan popup SOB; judul kolom `ID` · `Client ID` · `Name`; Choose → `SetDataSobCeding_Act`.
- `Activity\AddCedingList_act.xml`: `CedingCoList(<APPEND>)` — menambah baris kosong.
- `Activity\SetDataSobCeding_Act.xml` (btnQuotation `CedingCo`): `CedingCoList(<LAST>).CedingCo` = `pxResults(1).ID`,
  `.CedingCoName` = `pxResults(1).ClientName` (`BrowseSearchSobCeding_RD`, Param.ID); lalu `Quotation.CedingCo` /
  `CedingCoName` dikosongkan dan disusun ulang per baris: `@If(x=="", .CedingCo, x+";"+.CedingCo)`; diakhiri
  `Obj-Save pyWorkPage`.
- `Activity\SetCedingCo_Act.xml`: penggabungan `;` yang sama (bila `CedingCoName` kosong) lalu
  `OfferFacIn.QuotationData = Quotation`.
- `Activity\DeleteCeding_Act.xml`: `@replaceAll(CedingCo, cedingco+";", "")`, idem nama; `Page-Remove` baris;
  `Call CheckCeding_Act` (tidak ditelusur).
- DDL `FACINOFFER.txt` / `FACINPRODUCTION.txt`: `CEDINGCO` VARCHAR2(1000); FACINPRODUCTION juga `CEDINGID` VARCHAR2(100).
  Rancangan flat `STRUKTUR-TABEL-NB-FACIN.md` hanya punya `T_QUOTATIONDATA.CEDING_CO_NAME`.

## Keputusan agent

- **E-5** Baris ditambah hanya saat Choose — tidak ada baris kosong (Pega: `AddCedingList_act` menambah baris kosong
  dulu; bila popup ditutup tanpa memilih, baris kosong tertahan oleh tanda wajib sel 17).
- **E-6** Popup daftar dan popup pencari tampil bergantian, bukan bertumpuk: Modal inti memasang Escape di
  `document` dan tidak memakai portal, sehingga dua modal bertumpuk akan tertutup bersamaan oleh satu Escape.
- **E-7** Submit menyerahkan daftar ke layar; tersimpan lewat Save for later (sejalan E-1 tiket 33; Pega
  `Obj-Save` saat Choose). Cancel/X membuang perubahan. Ceding ganda tidak ditolak — Pega tidak memeriksanya.
- **E-8** Gabungan nama memakai `;` tanpa spasi, urut sesuai urutan pilih (`SetCedingCo_Act`).

## Kontrak

- Pencari: `GET /api/nbfacin/sob` (tiket 33) — RD yang sama.
- `GeneralInward.cedingList: [{ id, name }]` (urut) + `cedingCoName` (gabungan `;`).
- `PUT …/general` menerima `cedingIds: string[]` (urut, boleh kosong). Server memeriksa tiap kode ke AGENT dengan
  syarat tiket 33; ada yang tidak lolos → 400; lolos → simpan daftar dengan nama dari AGENT, plus gabungan `;` kode
  dan nama di `T_QUOTATIONDATA`.
- Penyimpanan daftar: tabel anak baru (migrasi nbfacin, rentang 180–219) — bentuk dan nama `belum terverifikasi`,
  usul `T_QUOTATION_CEDING (CASE_ID, URUTAN, CEDING_CO, CEDING_CO_NAME)`; kolom gabungan kode di `T_QUOTATIONDATA`
  `belum terverifikasi` (usul `CEDING_CO`, mengikuti DDL `CEDINGCO`).

## Acceptance criteria

- [x] Tombol Change Ceding Co hidup, membuka popup *Ceding Co List*; label diuji ke ShowCedingCoList / CedingCoHierarki.
- [x] Daftar kosong "No items"; Add / Select Ceding → pencari *Ceding Company* (ID · Client ID · Name, Choose).
- [x] Choose menambah baris; Delete membuang baris; Submit mengisi Ceding co name (gabungan `;`).
- [x] Save for later mengirim `cedingIds`.
- [x] `GET kasus` mengembalikan `cedingList` + `cedingCoName`; `PUT …/general` memeriksa dan menyimpan `cedingIds`.
- [x] Migrasi penyimpanan daftar ditulis (tidak dijalankan agent) — 185, tabel rancangan `T_CEDINGCOLIST`.

## Backend (sesi c3, 03-10-2026) — disusun agent

- **Penyimpanan = bentuk RANCANGAN** (permintaan: "kalau rancangan sudah punya bentuk lain, pakai itu"): `[terverifikasi]`
  `loader/skema_gen.go` sudah memuat **`T_CEDINGCOLIST`** (jalur `QuotationData/CedingCoList`, induk `T_QUOTATIONDATA`: ID,
  IDPEGA, COB_GROUP, PARENT_ID, SEQ_NO, ROW_UID, `CEDING_CO` VARCHAR2(50), `CEDING_CO_NAME` VARCHAR2(500)) dan
  `T_QUOTATIONDATA.CEDING_CO` / `CEDING_CO_NAME` — usulan `T_QUOTATION_CEDING` **tidak** dipakai. Fixture: setiap case
  tepat 1 ceding (kode 8 karakter), tidak pernah bergabung `;`.
- Migrasi **185** (ditulis, tidak dijalankan; butuh 183): `CREATE TABLE T_CEDINGCOLIST` utuh + `SEQ_T_CEDINGCOLIST` +
  indeks `(PARENT_ID, SEQ_NO)`; `T_QUOTATIONDATA ADD (CEDING_CO VARCHAR2(1000))`, `MODIFY (CEDING_CO_NAME VARCHAR2(4000))`
  — keputusan work owner butir 80 (lebar gabungan; DDL Pega `CEDINGCO` VARCHAR2(1000)); loader ikut (`amandemenLebar`).
- `GET kasus` → `general.cedingList [{id, name}]` urut `SEQ_NO` (selalu larik) + `general.cedingCoName` (`T_QUOTATIONDATA.CEDING_CO_NAME`).
- `PUT …/general` menerima `cedingIds: string[]` — daftar **diganti utuh** di dalam transaksi: tiap kode dicek ke `AGENT`
  dengan syarat tiket 33 (termasuk butir 79.1, AgentType2 NULL ikut); ada yang tidak lolos → **400** (pesan menyebut
  `cedingIds[n]`); lolos → baris `T_CEDINGCOLIST` urut pilih + gabungan `;` kode/nama di `T_QUOTATIONDATA`. Kosong / tidak
  dikirim → daftar dan gabungan dikosongkan. Ganda tidak ditolak (E-7).

**Keputusan agent (menunggu konfirmasi):**

| # | Keputusan | Dasar |
| --- | --- | --- |
| A105 | Tabel rancangan `T_CEDINGCOLIST` + `T_QUOTATIONDATA.CEDING_CO`, bukan `T_QUOTATION_CEDING` usulan | rancangan sudah punya bentuknya |
| A106 | `ROW_UID` baris aplikasi = UUID v4 acak | loader memakai v5 deterministik; K-073: sementara, tidak dijadikan sandaran |
| A107 | Kode kosong / > 50 bita → 400 sebelum basis data; nama AGENT > 500 bita atau gabungan > 1000/4000 bita → 400 | lebar kolom `T_CEDINGCOLIST` (rancangan) dan butir 80 |
| A108 | `cedingIds` tidak dikirim = dikosongkan (sama dengan medan General lain, A90) | PUT menimpa seluruh General |

Catatan: gabungan mengikuti `@If` Pega persis (nama kosong di depan tanpa `;` awal, di tengah `;;`). Jalur mundur 185
gagal (ORA-01441) bila ada `CEDING_CO_NAME` > 500 bita. `[dugaan]` langkah `ProtectCeding.CARI22 = 0`
(`SetDataSobCeding_Act.xml` L1594, cabang `btnQuotation=="CedingCo"`) belum ditelusur — tidak diport. Uji menyimpan
daftar (urutan DELETE/INSERT, SEQ_NO) hanya lewat teks SQL — tanpa Oracle.
