# 31: Backend layar Inward Facultative tahap 2 — baca case, simpan General, pilihan Marketing Name

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah: permintaan sesi `nusantarare-0f`
> 03-10-2026 (work owner "commit dan lanjut") untuk tiket 30 tahap 2; butir yang perlu tafsiran ditanyakan langsung ke
> work owner di sesi ini (AskUserQuestion 03-10-2026, register butir 78).

**What to build:** endpoint baca case NB, simpan blok **General** (tombol *Save for later*), dan pilihan **Marketing
Name** untuk layar `frontend/pages/InwardFacultative.tsx`. **Submit tidak** di tiket ini (pasca-proses
`AddToListSuggestOfferFacIn_DT` + `SetValidateDate_PostAct` + Decision3 — tahap terpisah).

**Blocked by:** — (tiket 29 selesai; migrasi 180/181 sudah dijalankan work owner di DEV)

**Status:** ready-for-human — dibangun 03-10-2026; uji tanpa Oracle hijau; migrasi 182/183 **ditulis, tidak dijalankan**

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\`

- `NB FacIn\Section\Periode.xml` — sel → properti: 9 `.QuotationData.NoOfferSlip`, 10 `.QuotationData.StatusBusiness`,
  16 `.QuotationData.InsuredName`, 20 `.QuotationData.QQName`, 21 `.PolicyData.StartDateTime`, 22 `.PolicyData.OfferingDate`,
  26 `.QuotationData.PolicyType` (radiogroup), 42 `.QuotationData.BusinessName`, 43 `.QuotationData.TypeFacultative`
  (dropdown, `associated`), 48 `.QuotationData.SobName`, 49 `.QuotationData.CedingCoName`, 56 `.QuotationData.GroupName`,
  60 `.PolicyData.EndDateTime`, 72 `.Following` (Old Policy Number), 75 `.QuotationData.MOID`, 78 `.QuotationData.EDMDay`.
  ⚠️ Ralat atas brief: Source of business = `.QuotationData.SobName` (bukan SourceOfBusiness); Old Policy Number =
  `.Following` di akar (bukan QuotationData).
- Rancangan (`loader/skema_gen.go`): `.PolicyData.*` dilipat ke `T_GENERAL_POLIS` (V-39) — `START_DATE_TIME`,
  `OFFERING_DATE`, `END_DATE_TIME` VARCHAR2(30), `FOLLOWING` VARCHAR2(50); `.QuotationData.*` → `T_QUOTATIONDATA`
  (induk `T_GENERAL_POLIS`) — `NO_OFFER_SLIP`, `QQ_NAME`, `POLICY_TYPE`, `MOID`, `EDM_DAY`, `TYPE_FACULTATIVE`, `SOB_NAME`,
  `CEDING_CO_NAME`, `GROUP_NAME`.
- Bentuk nilai di 5 fixture (`services/premium/testdata/kasus`), **per wadah** (dihitung dua cara: penjelajah JSON
  per jalur vs grep mentah 8 = 5 terkini + 3 salinan `OldData`):
  - wadah terkini `PolicyData`/`QuotationData`: `OfferingDate` 8 digit ×5; `StartDateTime` `…T050000.000 GMT` ×4,
    `…T170000.000 GMT` ×1; `EndDateTime` 050000 ×3, 170000 ×2; `PolicyType` `"0"` ×4, `"2"` ×1; `TypeFacultative`
    `"FacultativeIn"` ×3; `EDMDay` `"365"` ×4;
  - salinan `OldData` / `OldData/OldData`: tanggal seluruhnya 170000 (×3 masing-masing Start/End), `PolicyType` `"0"` ×3.
  - ⚠️ **Ralat agent:** pertanyaan 78.1 semula menyebut "17:00 GMT, sebagian 05:00" — angka itu menjumlahkan lintas wadah
    (jebakan sensus 5). Hitungan per wadah di atas diajukan ulang ke work owner (78.1 diralat).
  - Kondisi rule korpus memakai `PolicyType` 0/1/2.
- `DataTransform\InwardFacultative_PreDT.xml` langkah **5–5.1** (aktif): WHEN
  `pyWorkPage.OfferFacIn.PolicyData.OfferingDate==""` SET `@getCurrentDateStamp()`. Langkah 2–2.1
  (`@DateTime.CurrentDate("dd/MM/yyyy","")`) ber-`pyDisabled` true — tidak berlaku (ralat; temuan code review).
- Nilai bawaan sel (`pyDefaultValue`, `Periode.xml`): sel 26 Policy Type `0`, sel 43 Type facultative `FacultativeIn`
  (juga sel 10/11 `New Bisnis`/`New Business`) — kode Pega; **tidak** diterapkan backend (78.2/78.3 menyimpan label).
- Dropdown sel 75: `pyListSource` reportdefinition `BrowseMarketingOfficer_RD`, nilai `.ID`, pyPrompt `.ClientName`,
  "Choose". RD (kelas `ASM-FW-GISFW-Int-marketingofficer`): filter `.MOStatus = "1"` (Text), urut `.ID` DESC,
  `pyMaxRecords` 500. Tabel `POOLDATA.MARKETINGOFFICER` (DDL 15 kolom, dihitung dua cara);
  `PEGA_MARKETINGOFFICER.txt` = **prosedur** penulis tabel itu, bukan tabel.

## Kontrak

| Rute | Jawaban |
| --- | --- |
| `GET /api/nbfacin/kasus/{caseId}` | 200 `{"caseId","position","statusWork","opportunity":{14 medan IsianOpportunity},"insuredName","general":{"reffNumber","qqName","beginDate","offeringDate","endDate","policyType","marketingId","day","typeFacultative","sourceOfBusiness","cedingCoName","groupName","oldPolicyNumber"}}`; 404 case tidak ada / LINI bukan `FAC`; 503; 500 |
| `PUT /api/nbfacin/kasus/{caseId}/general` | badan = sembilan medan `general` yang dapat diubah (tanpa empat medan tampil-saja); 200 objek case seperti GET; 400 JSON / tanggal / panjang; 401 tanpa identitas; 404; 503; 500 |
| `GET /api/nbfacin/marketing-officer` | 200 `{"baris":[{"id","nama"}]}` (urut `ID` DESC, ≤ 500); 503; 500 |

Tanggal kabel `DD-MM-YYYY`, kosong = `""`. `offeringDate` kosong di basis data dikembalikan **hari ini (WIB)** (PreDT).
Teks tersimpan yang bukan bentuk Pega dikembalikan apa adanya (A89). ⚠️ PUT **menimpa kesembilan medan**: medan yang
tidak dikirim atau `""` mengosongkan kolomnya — kirim seluruh objek `general` (dari GET) beserta perubahannya.

## Yang dibangun

- [x] Migrasi `182_t_general_polis.sql`, `183_t_quotationdata.sql` (+ `SEQ_T_QUOTATIONDATA`) — sebagian (butir 78.4),
      nama/tipe = rancangan (`TestMigrasiFlatSebagianCocokRancangan`)
- [x] `repository/kasus.go` (baca: `T_WORK_POLIS` ⋈ `T_NB_OPPORTUNITY` ⋈ `T_GENERAL_POLIS` ⋈ `T_QUOTATIONDATA`, nama
      tertanggung `T_M_ACCOUNT`; simpan: sentuh `TGL_UPDATE` case → UPDATE-atau-INSERT dua tabel), `repository/marketing.go`
- [x] `services/kasus.go` (konversi tanggal WIB, PreDT, pemeriksaan), `handlers/kasus.go`
- [x] `MODUL.md` + `docs/STRUKTUR-TABEL-NB-FACIN.md` (`T_GENERAL_POLIS`, `T_QUOTATIONDATA`, warisan `MARKETINGOFFICER`)
- [ ] Uji terhadap Oracle sungguhan — ⛔ tidak dijalankan

## Keputusan work owner (butir 78, AskUserQuestion 03-10-2026)

- **78.1** tanggal disimpan **teks bentuk Pega** (kolom rancangan VARCHAR2(30)): Offering `YYYYMMDD`; Begin/End =
  ~~00:00 WIB (`…T170000.000 GMT` hari sebelumnya)~~ → **diralat: 12:00 WIB = `…T050000.000 GMT`** (mayoritas nilai terkini;
  delapan digit teks = tanggal pilihan); baca GMT → WIB → tanggal.
- **78.2** Policy Type disimpan **label layar apa adanya** (mis. `Individual Policy`) — menyimpang dari kode Pega 0/1/2
  sampai pemetaannya diketahui.
- **78.3** Type facultative disimpan **seperti dikirim layar** (`Facultative In`), bukan `FacultativeIn`.
- **78.4** `T_GENERAL_POLIS` / `T_QUOTATIONDATA` **dibuat sekarang dengan kolom yang diperlukan saja**; sisanya tiket 23.

## Keputusan agent — menunggu konfirmasi

| # | Keputusan | Dasar |
| --- | --- | --- |
| A86 | GET tidak meminta identitas; PUT meminta (401) | pola: baca nbfacin tanpa identitas (tiket 27/28), tulis dengan pembuat (tiket 29) |
| A87 | Satu baris `T_QUOTATIONDATA` per case — `UQ_T_QUOTATIONDATA_PARENT` | halaman tunggal `QuotationData` di rancangan; loader 1 baris/case (kasus_test) |
| A88 | Nama tertanggung = `MAX(INSUREDNAME)` `T_M_ACCOUNT` per `ACCOUNT_ID` | DDL `T_M_ACCOUNT` tanpa PK — `[dugaan]` ID unik |
| A89 | Tanggal tersimpan yang bukan bentuk Pega dikembalikan apa adanya; simpan berikutnya menolaknya (400), tidak menimpa diam-diam | data lama belum dimuat |
| A90 | Save for later menyimpan **sebagian** — tidak ada medan wajib; ditolak hanya tanggal tak sah dan isian melebihi lebar kolom; kesembilan medan selalu ditulis (yang kosong = NULL) | brief sesi 0f: wajib-isi Pega ditegakkan saat Submit `[dugaan]` |
| A91 | PreDT diterapkan **saat baca** (Offering kosong → hari ini WIB); tidak ditulis sampai disimpan | PreDT berjalan tiap kali layar dibuka |
| A92 | `T_QUOTATIONDATA.ID` NUMBER(19) dari `SEQ_T_QUOTATIONDATA` | pola identitas-dari-sequence yang penjaga izinkan; presisi kolom ID belum diputus tim inti |
| A93 | `MOID` disimpan tanpa memeriksa keberadaannya di `MARKETINGOFFICER`; lebar `MOID` 50 < `MARKETINGOFFICER.ID` 100 → ID > 50 bita ditolak 400 | rancangan `MOID VARCHAR2(50)` |
| A94 | Simpan juga memperbarui `T_WORK_POLIS.TGL_UPDATE` (sekaligus mengunci baris dan memastikan case ber-LINI `FAC`) | pola premiumlistlife (setiap ubah baris kasus) |

⚠️ **Tidak dibangun (selisih dengan Pega, dicatat):** Pega juga menyimpan `.QuotationData.StatusBusiness`,
`.InsuredName`, `.BusinessName` (sel 10/16/42) — layar mengambilnya dari opportunity/akun; kolomnya belum dibuat.
Langkah PreDT lain (`ViewOfferFacInUW_PreDT`, `SetObjNoDT`, cek `RICommision`) tidak menyentuh blok General — tidak
diport.

## Comments

*Tanpa nama orang, tanpa data pelanggan — uji memakai data sintetis berawalan `UJI-`.*
