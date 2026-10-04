# Modul `masterdata` — Master Data

Modul di LUAR dua puluh folder korpus (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 5), dibuat atas keputusan work
owner 04-10-2026: "biar data nya ke update km buatin menu master dengan nama masing2 tabel tadi itu. misal manu master
province, isi nya untu insert, update, aktif, non aktif kan data dari master itu" — penempatan "Modul baru
'masterdata'", cakupan "Yang dipakai NB FacIn dulu".

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah nilainya hanya
lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `masterdata` |
| Folder korpus | `Master Data` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERDATA` |
| Status | dimigrasi |
| Rentang migrasi | `760-799` |
| Slot menu | `990-991` |
| Prefix rute API | `/api/masterdata` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | rencana dan tiket (`issues/NN-<slug>.md`), `STRUKTUR-TABEL-MASTER-DATA.md` |
| `backend/` | `models/` (daftar master + kolom), `repository/`, `services/`, `handlers/` (`GET /api/masterdata`, `GET` / `POST /api/masterdata/{tabel}`, `PUT /api/masterdata/{tabel}/{id}`, `PUT /api/masterdata/{tabel}/{id}/status`), `migrations/` (760 view → tabel flat + `STS_AKTIF`, 761 `T_MASTER_STATUS` (status NATION), slot menu 990), `modul.go` |
| `frontend/` | `menu.ts`, `rute.tsx`, `api.ts`, `labels.ts`, `masterdata.css`, `pages/HalamanMasterData.tsx` (tab per master), `components/` (`DaftarMaster`, `FormMaster`, `PilihRujukan`) — sesi 0f |

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/masterdata/...
npx vitest run modul/masterdata
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per butir.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang `docs/STRUKTUR-TABEL-MASTER-DATA.md` gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris = kepemilikan tabel
berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `NATION` | tabel negara warisan POOLDATA; Master Data menulis isinya (menu Nation) tanpa mengubah strukturnya - statusnya di T_MASTER_STATUS (761); nbfacin membacanya untuk saran Country (tiket 46) |
| `OBJECTITEMTYPE` | tabel jenis item objek warisan POOLDATA (sumber view V_JN_OBJ_ITEM); Master Data menulis isinya dan memakai ISACTIVE-nya |
| `BRANCH` | tabel cabang warisan POOLDATA; Master Data hanya MEMBACA ID untuk rujukan City |
