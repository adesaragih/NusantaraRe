# Modul `bordereaux` — Bordereaux

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
04-10-2026: folder korpus `D:\XML\RNM_BRD\Bordereaux` sebagai modul/menu baru, kelompok MASTER TREATY. Padanan Pega:
kelas `ASM-FW-GISFW-Int-BORDEREAUX` dan 29 kelas detail (`-PREMI_*`, `-CLAIM_*`, `-SUBROGATION_BONDING`); harness
`PortalBordereaux`, `InputBordereaux`; activity `UploadCSVBordereaux_Act`, 29 `Mapping*Bdx*`, `BdxSave_Act`,
`SaveDetailBordereaux_Act`, `ActionSubmit`, `CountSummaryBdx`, `GetZipCode_AI_Bdx_Act`. Baris menunya dibuat migrasi
inti `913` (slot menu modul tidak boleh membuat kelompok), dan `Folder korpus` di bawah = label menu `M_NAV_MENU.LABEL`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — diambil dari
cadangan `760-899` dan `990-999` (Aggregate dipersempit ke `880-889`, keputusan work owner 04-10-2026).

| Kunci | Nilai |
| --- | --- |
| Nama modul | `bordereaux` |
| Folder korpus | `Bordereaux` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-BORDEREAUX` |
| Status | dimigrasi |
| Rentang migrasi | `890-899` |
| Slot menu | `998-999` |
| Prefix rute API | `/api/bordereaux` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-BORDEREAUX.md` — tabel yang ditulis dan dibaca, sumber setiap kolom |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `templat/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `bordereaux.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas) — keputusan work owner 04-10-2026

- XML Pega patokan dasarnya; **bug Pega diperbaiki di Go** (format ID tetap mengikuti XML).
- **Tanpa JSON**: Pega menyimpan seluruh kasus di `M_BORDEREAUX.DATA_JSON` dan baru menulis tabel detail saat
  Resolve-Complete. Di sini setiap Save menulis header ke `BORDEREAUX` dan detail ke tabel detailnya; riwayat
  persetujuan ke tabel baru `BORDEREAUX_HISTORY`. JSON lama hanya DIBACA oleh fitur Copy Old Data (superadmin).
- Akses: menu Bordereaux ber-hak **PENUH** = Input Data, Edit/Delete/Submit berkas sendiri; **View only** = lihat dan
  unduh (`M_LOGIN_GO_MENU.HAK`, keputusan work owner 04-10-2026, pengganti `ReasBordereauxAdmin` yang dibuang 893).
  Peran alur kerja tetap workbasket: `ReasBordereauxChecker` (bukan pembuat berkasnya), `ReasBordereauxSupervisor` -
  cukup View only. Superadmin (pemegang Kelola User) = pengganti `IT Developer` Pega.
- Upload CSV: pemisah `;` (atau `,`), baris pertama kepala, kolom menurut posisi (29 kombinasi Type x Business,
  `models/kombinasi_gen.go`); nilai dibaca menurut TIPE KOLOM tabel; validasi per baris dan kolom.
- Templat unduhan dikelola **Template Manager** (inti): 29 slot `bordereaux.*`; berkas bawaan `backend/templat/`
  berisi BARIS HEADER SAJA (data contoh tidak masuk repositori - unggah berkas lengkap lewat Template Manager).
- Kode pos AI (Gemini) untuk FIRE dan ENGINEERING berjalan di belakang lewat outbox.
- Lampiran (keputusan work owner 08-10-2026, "kaya XML nya"): panel `AttachmentsBdx` di form berkas tersimpan - grid
  kategori `M_KATEGORIBORDEREAUX` + jumlah, Upload File (banyak berkas), View File (unduh, View Office Online untuk
  xls..pptx, Delete). Berkas ke Google Storage lewat komponen bersama `inti/backend/penyimpanan` (padanan kelas Pega
  `T_STORAGE_IMAGE`: Insert/GetUrl/Delete/GeminiAI); rekam di `M_ATTACHMENTBORDEREAUX`. Upload File dan Delete hanya
  saat berkas dibuka Edit oleh yang berhak (`HakAtas.Ubah`, menu PENUH); mode View dan menu View only tidak bisa,
  termasuk superadmin (keputusan work owner 08-10-2026 - pengecualian `IT Developer` Pega dibuang).

## Migrasi

Rentang `890-899`:

| Berkas | Isi |
| --- | --- |
| `890_bordereaux_history.sql` | tabel `BORDEREAUX_HISTORY` + `SEQ_BORDEREAUX_HISTORY` |
| `891_workbasket_bordereaux.sql` | tiga baris `M_WORKBASKET`: `ReasBordereauxAdmin`, `ReasBordereauxChecker`, `ReasBordereauxSupervisor` |
| `892_spread_aviation.sql` | kolom spread `BORDEREAUX_CLAIM_AVIATION` dan `BORDEREAUX_PREMI_AVIATION` dilebarkan `NUMBER(10,4)` ke `NUMBER(38,8)` (isinya nominal) |
| `893_hapus_workbasket_admin.sql` | baris `M_WORKBASKET` `ReasBordereauxAdmin` dan pemegangnya dibuang - Input Data kini hak menu PENUH |

Slot menu `998`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`. Barisnya dibuat migrasi inti
`913_m_nav_menu_bordereaux.sql`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/bordereaux/...
npx vitest run modul/bordereaux
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `BORDEREAUX` | tabel warisan Pega (header berkas); modul ini menyisipkan, mengubah, membaca, dan menghapus barisnya tanpa mengubah strukturnya |
| `M_BORDEREAUX` | tabel warisan Pega (JSON kasus); hanya DIBACA Copy Old Data dan barisnya dihapus bersama berkasnya |
| `M_WORKBASKET` | master peran warisan; migrasi 891 menambah tiga baris data |
| `M_ATTACHMENTBORDEREAUX` | tabel warisan Pega (lampiran) |
| `M_KATEGORIBORDEREAUX` | tabel warisan Pega (kategori dokumen), hanya dibaca |
| `M_PROMPT_AI` | prompt kode pos AI, hanya dibaca |
| `TREATY_IN` | Master Treaty (popup Choose Master Treaty), hanya dibaca |
| `AGENT` | Cedant (popup Choose Master Treaty), hanya dibaca |
| `BORDEREAUX_PREMI_FIRE` | detail warisan Pega |
| `BORDEREAUX_PREMI_ENGINEERING` | detail warisan Pega |
| `BORDEREAUX_PREMI_MARINE_CARGO` | detail warisan Pega |
| `BORDEREAUX_PREMI_MARINE_HULL` | detail warisan Pega |
| `BORDEREAUX_PREMI_MOTOR_VEHICLE` | detail warisan Pega |
| `BORDEREAUX_PREMI_BONDING` | detail warisan Pega |
| `BORDEREAUX_PREMI_LIABILITY` | detail warisan Pega |
| `BORDEREAUX_PREMI_PERSONAL_ACCIDENT` | detail warisan Pega |
| `BORDEREAUX_PREMI_OFFSHORE` | detail warisan Pega |
| `BORDEREAUX_PREMI_CREDIT` | detail warisan Pega |
| `BORDEREAUX_PREMI_ONSHORE` | detail warisan Pega |
| `BORDEREAUX_PREMI_HIO` | detail warisan Pega |
| `BORDEREAUX_PREMI_AVIATION` | detail warisan Pega |
| `BORDEREAUX_PREMI_MONEY_INSURANCE` | detail warisan Pega |
| `BORDEREAUX_CLAIM_FIRE` | detail warisan Pega |
| `BORDEREAUX_CLAIM_ENGINEERING` | detail warisan Pega |
| `BORDEREAUX_CLAIM_MARINE_CARGO` | detail warisan Pega |
| `BORDEREAUX_CLAIM_MARINE_HULL` | detail warisan Pega |
| `BORDEREAUX_CLAIM_MOTOR_VEHICLE` | detail warisan Pega |
| `BORDEREAUX_CLAIM_BONDING` | detail warisan Pega |
| `BORDEREAUX_CLAIM_LIABILITY` | detail warisan Pega |
| `BORDEREAUX_CLAIM_PERSONAL_ACCIDENT` | detail warisan Pega |
| `BORDEREAUX_CLAIM_OFFSHORE` | detail warisan Pega |
| `BORDEREAUX_CLAIM_CREDIT` | detail warisan Pega |
| `BORDEREAUX_CLAIM_ONSHORE` | detail warisan Pega |
| `BORDEREAUX_CLAIM_HIO` | detail warisan Pega |
| `BORDEREAUX_CLAIM_AVIATION` | detail warisan Pega |
| `BORDEREAUX_CLAIM_MONEY_INSURANCE` | detail warisan Pega |
| `BORDEREAUX_SUBROGATION_BONDING` | detail warisan Pega |
