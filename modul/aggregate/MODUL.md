# Modul `aggregate` — Aggregate

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
04-10-2026: *"hanya modul aggregate"*, *"Aggregate - Master Treaty"*. Padanan Pega: kelas
`ASM-FW-GISFW-Int-AGGREGATE`, folder korpus `D:\XML\RNM_BRD\Aggregate` (7 rule: `ShowAggregateList`,
`GridDasbordAgg`, `ChooseMasterID`, `GetMasterIDAgg_Act`, `SetMasterID`, `UploadCSVAggregate_Act`,
`SaveAggregate_Act`). Baris menunya dibuat migrasi inti `911` (slot menu modul tidak boleh membuat kelompok), dan
`Folder korpus` di bawah = label menu `M_NAV_MENU.LABEL`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — diambil dari
cadangan `760-899` dan `990-999`. Dipersempit 04-10-2026 dari `880-899` ke `880-889` (modul ini memakai 880 saja) supaya
Bordereaux mendapat `890-899` - keputusan work owner "aku ikuti rekomendasi kamu".

| Kunci | Nilai |
| --- | --- |
| Nama modul | `aggregate` |
| Folder korpus | `Aggregate` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-AGGREGATE` |
| Status | dimigrasi |
| Rentang migrasi | `880-889` |
| Slot menu | `996-996` |
| Prefix rute API | `/api/aggregate` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-AGGREGATE.md` — tabel yang ditulis dan dibaca, sumber setiap kolom |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `aggregate.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas) — keputusan work owner 04-10-2026

- Langkah XML Pega diikuti apa adanya. Langkah bertanda `//` (tidak pernah jalan di Pega) tidak dipindahkan.
- **Layar daftar** (`GridDasbordAgg`): satu baris per Tanggal Input (hari), Ceding Code, Ceding Name, Treaty Type,
  As At, UW Year; klik ganda = rincian; Delete = hapus seluruh baris kunci itu (setelah konfirmasi). Pie chart
  RNM Value (USD) per Coverage per As At dirancang sendiri (work owner: *"pie chart kamu buat sendiri aja"*).
- **Layar unggah** (`ShowAggregateList`): Master ID (popup `ChooseMasterID`), Template, Upload CSV (versi Go
  `UploadCSVAggregate_Act`), grid pratinjau, Save (`SaveAggregate_Act`).
- Upload CSV: kolom menurut nomor urut, baris pertama kepala; Assessment Zone lewat `ASSESSMENT_ZONE`; Ceding dari
  Master ID pertama; Treaty Year dari periode `TREATYYEAR` yang memuat As At; kurs dari `TREATYEXCHANGEYEARLY` (IDR =
  1 / TOIDR baris USD); beberapa Master ID = share dijumlah, ID digabung `;`; satu baris "Total :" per mata uang.
- Save: pesan penolakan Pega apa adanya; seluruh baris dalam satu transaksi; `ID` = `AGG-<SEQ_AGGREGATE>`;
  `COMMENCEMENT` tidak diisi; data kembar tidak dicegah (seperti Pega).

## Migrasi

Rentang `880-889`:

| Berkas | Isi |
| --- | --- |
| `880_seq_aggregate.sql` | `SEQ_AGGREGATE` mulai dari nomor AGG terbesar + 1, dihitung saat migrasi berjalan (blok sequence-dari-kueri) |

Slot menu `996`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`. Barisnya dibuat migrasi inti
`911_m_nav_menu_aggregate.sql`. Karena 880 membaca `AGGREGATE`, skema uji (`uji/skemauji`) membuat tiruannya sebelum
migrasi.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/aggregate/...
npx vitest run modul/aggregate
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `AGGREGATE` | tabel warisan Pega (`ASM-FW-GISFW-Int-AGGREGATE`); modul ini menyisipkan, membaca, dan menghapus barisnya tanpa mengubah strukturnya |
