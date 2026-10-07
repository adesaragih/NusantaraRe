# 33: Popup Change SOB — cari SOB dari tabel AGENT, pilih, simpan bersama General

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026:
> *"selanjutnya tombol change sob … diambil dari tabel agent, dimana yang ditarik dengan kondisi StatusActive=1,
> AgentType2 != "LIFE INSURANCE", ClientID is not null dan dapat di search tanpa memperhatikan huruf besar dan kecil"*,
> disertai tangkapan layar popup (tidak disalin ke repo).

**What to build:** tombol **Change SOB** di layar Inward Facultative membuka popup pencarian SOB; **Choose** mengisi
*Source of business*; kodenya ikut tersimpan lewat *Save for later* (tiket 31).

**Blocked by:** — (~~DDL tabel AGENT belum ada~~ → `D:\migrasi\RNM\DDL\AGENT.txt` ditambahkan work owner 03-10-2026)

**Status:** frontend selesai 03-10-2026 (uji hijau); backend → sesi c3; migrasi kolom kode SOB **ditulis, tidak dijalankan**.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\`

- `NB FacIn\Section\Periode.xml` sel 52 `Change SOB` → `showHarness` HarnessName `SOB`, WindowName `Change SOB`,
  Target `popup`, activity `ASM-FW-GISFW-Data-Quotation.InputQuotation_PreAct`. Sel 48 menampilkan `.QuotationData.SobName`.
- `NB FacIn\Harness\SOB.xml`: kotak `Search` (sel 294, `pyLabelFieldValue`, properti `SearchSOB.CARI1`); grid
  `TempBusinessSource.pxResults` kelas `ASM-FW-GISFW-Data-Agent`, kolom `.ID` dan `.ClientName`, `pyPageSize` 20 (L3885);
  tombol `Choose` (`pyLabel`, L3510) → runActivity `SearchHierarkiSourceBizAgentTreatyIn_Act` → `opener.location.reload`
  → `window.close`.
- `SearchHierarkiSourceBizAgentTreatyIn_Act` **tidak ada di korpus** — syarat penarikan data diambil dari perintah work
  owner, bukan dari Pega.
- `DDL\FACINPRODUCTION.txt`: pasangan kolom `SOB` VARCHAR2(1000) + `SOBID` VARCHAR2(100) — nama dan kode SOB disimpan
  berdua. Ekspor case: `QuotationData.SourceOfBusiness` berisi kode (bentuk `G0000075`, `B0000070`).

## Keputusan agent

- **E-1** Choose mengisi kode + nama di layar dan menutup popup; tersimpan saat *Save for later* (Pega: activity
  pengisi tak ada di korpus, lalu reload layar pembuka). Belum ada simpan langsung saat Choose.
- **E-2** Judul kolom = tangkapan layar (ID · Client ID · Name). SOB.xml hanya memuat `.ID` dan `.ClientName`;
  kolom Client ID tidak terbaca di XML.
- **E-3** Pencarian lewat Enter di kotak Search — SOB.xml tidak memuat tombol Search terpisah.
- **E-4** Server menulis nama SOB dari tabel AGENT menurut kode yang dikirim (bukan nama dari klien); kode yang tidak
  lolos syarat ditolak 400.

## Kontrak

- `GET /api/nbfacin/sob?cari=&halaman=` → `{ baris: [{ id, clientId, name }], total, halaman, ukuran }`; ukuran 20.
  Syarat: `StatusActive = 1`, `AgentType2 <> 'LIFE INSURANCE'`, `ClientID IS NOT NULL`; cari = mengandung, tidak peka
  huruf, atas ID / Client ID / Name. Urutan `belum terverifikasi` (usul: ID).
  ⚠️ ~~`AgentType2 <> 'LIFE INSURANCE'` di Oracle membuang baris ber-AgentType2 NULL — perlu dikonfirmasi work owner.~~
  → **Diputus work owner 03-10-2026 (register butir 79.1): baris NULL IKUT** — syarat `(AGENTTPYE2 IS NULL OR AGENTTPYE2 <>
  'LIFE INSURANCE')` untuk cari, cacah, dan cek PUT. Kolom di DDL dieja `AGENTTPYE2`.
- `GeneralInward.sourceOfBusinessId` (baru) + `sourceOfBusiness` (nama); `PUT …/general` menerima `sourceOfBusinessId`.
- Kolom kode di `T_QUOTATIONDATA` (migrasi baru nbfacin, rentang 180–219) — nama kolom `belum terverifikasi`.

## Acceptance criteria

- [x] Tombol Change SOB hidup, membuka popup berjudul *Change SOB*.
- [x] Popup: kotak Search, paging, kolom ID · Client ID · Name, tombol Choose per baris; label diuji ke SOB.xml.
- [x] Choose mengisi Source of business; Save for later mengirim `sourceOfBusinessId`.
- [x] Endpoint `GET /api/nbfacin/sob` dengan syarat di atas, pencarian tidak peka huruf.
- [x] `PUT …/general` memeriksa kode ke AGENT dan menyimpan kode + nama; `GET kasus` mengembalikan keduanya.
- [x] Migrasi kolom kode SOB ditulis (tidak dijalankan agent) — 184, kolom `SOURCE_OF_BUSINESS`.

## Backend (sesi c3, 03-10-2026) — disusun agent

- `GET /api/nbfacin/sob?cari=&halaman=` → `{"baris":[{"id","clientId","name"}],"total","halaman","ukuran":20}`; 400 halaman /
  `cari` > 255; 503; 500. Tanpa identitas (pola lookup tiket 27/28).
- `GET /api/nbfacin/kasus/{caseId}` → `general.sourceOfBusinessId` (kode) + `general.sourceOfBusiness` (nama).
- `PUT …/general` menerima `sourceOfBusinessId`: kosong → kode + nama dikosongkan; terisi → dicek ke `AGENT` dengan
  syarat yang sama di DALAM transaksi; tidak lolos → **400**; lolos → tulis kode + `SOB_NAME` dari `AGENT` (E-4).
- Migrasi **184** `ALTER TABLE T_QUOTATIONDATA ADD (SOURCE_OF_BUSINESS VARCHAR2(50))` (+ `_down`) — **ditulis, tidak
  dijalankan**; butuh 183.

**Bukti tambahan `[terverifikasi]`:**
- `NB FacIn\ReportDefinition\BrowseAgentNonLife_RD.xml` (kelas `ASM-FW-GISFW-Int-AGENT`, dipakai
  `Section\SourceHierarki.xml`): logika `B AND E AND A AND C` = B `.AgentType2 != "LIFE INSURANCE"`, E `.ClientName Contains`
  (tidak peka huruf), A `.StatusActive = "1"` (Text), C `.ClientID IS NOT NULL`; urut `.ID` DESC; `pyMaxRecords` 200. Sama
  dengan syarat work owner.
- Rancangan flat sudah punya kolom kode: `T_QUOTATIONDATA.SOURCE_OF_BUSINESS VARCHAR2(50)` ← `.QuotationData.SourceOfBusiness`;
  5 fixture: huruf + 7 digit (`G…` ×7, `B…` ×1, lintas wadah termasuk `OldData`).

**Keputusan agent (menunggu konfirmasi):**

| # | Keputusan | Dasar |
| --- | --- | --- |
| A99 | Kolom kode = **`SOURCE_OF_BUSINESS`** rancangan (bukan `SOB_ID` usulan); lebar **50 tetap** (butir 79.2) walau `AGENT.ID` VARCHAR2(1000) dan `FACINPRODUCTION.SOBID` 100 — kode > 50 bita ditolak 400 | kolom rancangan untuk `.QuotationData.SourceOfBusiness`, isinya kode 8 karakter (fixture) |
| A100 | Urut **`ID` DESC** | RD `BrowseAgentNonLife_RD` |
| ~~A101~~ | ~~`AGENTTYPE2 <> 'LIFE INSURANCE'` apa adanya — baris NULL ikut terbuang~~ → **DIGANTI butir 79.1**: NULL ikut | semula: RD Pega memakai `!=` yang sama (`[dugaan]` SQL `<>`) |
| ~~A102~~ | ~~Nama tabel/kolom `[dugaan]`~~ → **`[terverifikasi]` DDL `AGENT.txt`**: `POOLDATA.AGENT`, `ID`, `CLIENTID`, `CLIENTNAME`, `STATUSACTIVE` (VARCHAR2(1000 BYTE)); kolom tipe-2 dieja **`AGENTTPYE2`** (dugaan semula `AGENTTYPE2` SALAH — akan ORA-00904; temuan code review + sesi 0f) | satu tempat di `repository/sob.go` |
| A103 | Cari atas ID + ClientID + nama (RD hanya `ClientName`); `cari` > 255 → 400; tanpa batas 200 RD (berhalaman 20) | perintah work owner / kontrak sesi 0f; pola A73 |
| A104 | `GET /api/nbfacin/sob` tanpa identitas | pola lookup tiket 27/28 |

⚠️ **Risiko tercatat:** `AGENT.CLIENTNAME` VARCHAR2(1000) lebih lebar dari `T_QUOTATIONDATA.SOB_NAME` VARCHAR2(500) (rancangan) —
nama > 500 bita membuat simpan gagal 500 (ORA-12899). Tidak ditangani (data belum terlihat).
