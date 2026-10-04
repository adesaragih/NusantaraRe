# Modul `companydetail` — Company Detail

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
03/04-10-2026: *"aku mau buat modul COMPANY DETAIL bentuk menu seperti gambar ini"*, *"buat dalam modul baru ya,
jangan di inti, nama modul companydetail"*. Padanan Pega: layar SFAGIS Company Detail (organisasi CRM) — **tidak ada
di korpus XML**; acuannya screenshot Pega dari work owner dan katalog DEV. Baris menunya dibuat migrasi inti `907`
(slot menu modul tidak boleh membuat kelompok), dan `Folder korpus` di bawah = label menu `M_NAV_MENU.LABEL`
(keputusan work owner: "Company Detail", kelompok MASTER).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — diambil dari
cadangan `760-899` dan `990-999`.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `companydetail` |
| Folder korpus | `Company Detail` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-COMPANYDETAIL` |
| Status | dimigrasi |
| Rentang migrasi | `800-839` |
| Slot menu | `992-993` |
| Prefix rute API | `/api/company-detail` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-COMPANYDETAIL.md` — tabel yang dibuat dan diubah, peta tabel warisan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `alat/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `companydetail.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas)

- **Tanpa JSON** (*"aku tidak mau ada json lagi"*): sumber data = `CLIENT` (`FLAG` `Org`), `CLIENT_PICLIST`,
  `CLIENT_ADDRESS`. Dokumen `M_CLIENT` hanya dibaca SEKALI oleh alat pindah (`backend/alat`), disetujui work owner.
  Prosedur `RDBINSERTCLIENT` tidak dipanggil; kolom barisnya ditiru.
- Layar: Company Detail (NPWP, Parent organization, COUNTRY*, Title, Organization Name*, Business Field*, Note), grid
  PIC (Name*, Position*, Gender, Email, Date of birth, Phone number), grid Address (Type, Address, Phone and Fax),
  tombol Create. Client Status, Established Date, Number of Employees, Owner **dihapus** (work owner 04-10-2026).
- Nomor ORG baru dari sequence `SEQ_CLIENT_ORG` (810), yang dimulai dari nomor ORG tertinggi Pega (`CLIENT.IDVIEW`
  dan `M_CLIENT.ID`) + 1; tanpa indeks unik (work owner).
- COUNTRY dari tabel `NATION` (dulu view, kini tabel datar — migrasi 802-804); `COUNTRY` = `OLDID`, `COUNTRYNAME` =
  `NOTE`. Dropdown lain dari `M_ENUMERASI` (tabel baru modul ini; *"jangan ada pake datapega"*).
- Phone and Fax: satu baris `CLIENT_ADDRESS` per nomor (*"CLIENT_ADDRESS itu kan list bisa banyak row"*).
- Baris PIC dan alamat boleh dihapus dari grid; organisasi tidak pernah dihapus. `BU_NOTE` tidak disentuh.
- Rinciannya: `docs/STRUKTUR-TABEL-COMPANYDETAIL.md` dan dokumentasi paket `backend/services`.

## Migrasi

Rentang `800-839`:

| Berkas | Isi |
| --- | --- |
| `800_m_enumerasi.sql` | tabel `M_ENUMERASI` |
| `801_m_enumerasi_isi.sql` | 213 baris pilihan, disalin sekali dari enumerasi Pega DEV 04-10-2026 |
| `802_nation_lepas_view.sql` | buang view `NATION` |
| `803_nation.sql` | tabel datar `NATION`, nama dan kolom sama dengan view-nya |
| `804_nation_isi.sql` | salin isi `NATION` dari `M_NATION` sekali |
| `805_client_kolom.sql` | `CLIENT`: `PARENT_ID`, `NOTE`, jejak dibuat/diubah |
| `806_client_piclist_kolom.sql` | `CLIENT_PICLIST`: `GENDER` |
| `807_client_address_kolom.sql` | `CLIENT_ADDRESS`: `TELFAX_TYPE`, `TELFAX_CODE`, `TELFAX_NO` |
| `808_m_client_trigger.sql` | `TRG_M_CLIENT` dan `TRG_M_CLIENT_PIC` dinonaktifkan |
| `809_m_enumerasi_telfax_tanpa_email.sql` | Phone and Fax tanpa pilihan EMAIL (barisnya tetap, AKTIF 0) |
| `810_seq_client_org.sql` | sequence `SEQ_CLIENT_ORG` nomor ORG; nilai awal = nomor ORG tertinggi + 1, dihitung saat migrasi |

Slot menu `992`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`. Barisnya dibuat migrasi inti
`907_m_nav_menu_companydetail.sql`. Karena 802-808 mengubah tabel warisan, skema uji (`uji/skemauji`) membuat
tiruannya SEBELUM migrasi dan membongkarnya SESUDAH migrasi mundur.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/companydetail/...
npx vitest run modul/companydetail
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `M_CLIENT` | dokumen organisasi Pega SFAGIS; kolom `ID`-nya dibaca migrasi 810 (nilai awal `SEQ_CLIENT_ORG`) dan dua triggernya dinonaktifkan (808); isinya dibaca sekali oleh alat pindah |
| `M_NATION` | sumber salinan sekali `NATION` (804); tidak diubah |
