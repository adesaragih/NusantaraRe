---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 2 Q14 (`.scratch/claim-life/grilling-ronde-2.md`), keputusan work owner — disetujui 2026-09-14
---

# Penyimpanan berkas klaim tetap di Google Storage

Berkas klaim Life **tetap disimpan di Google Storage**. Migrasi tidak memindahkan penyimpanan
berkas; yang berganti hanya aplikasi yang memanggilnya.

## Jalur yang ada sekarang `[terverifikasi]`

Delapan berkas `Claim Life` menyentuh penyimpanan:

| Rule | Peran |
| --- | --- |
| `Claim Life/Activity/InsertGoogleStorage_Act.xml` (`ASM-FW-GISFW-INT-T_STORAGE_IMAGE!INSERTGOOGLESTORAGE_ACT`, 160.027 byte) | unggah |
| `Claim Life/Activity/GetUrlGoogleStorage_Act.xml` | ambil URL |
| `Claim Life/Activity/DeleteGoogleStorage_Act.xml` | hapus |
| `Claim Life/Activity/InsertDocument_Act.xml`, `DeleteDocument_Act.xml`, `LoadDocumentLife_ACT.xml`, `DownloadDocumentClaim.xml` | pembungkus tingkat dokumen |
| `Claim Life/RDBList/GetTokenStorage_SQL.xml` (`ASM-FW-GISFW-INT-T_STORAGE_IMAGE!RNM!GETTOKENSTORAGE_SQL`) | **penerbitan token** |

Token **tidak** diterbitkan aplikasi. Ia datang dari Oracle:

```sql
BEGIN
  pooldata.GET_TOKEN_STORAGE ( {UploadDoc.App}, {OperatorID.pyUserIdentifier},
                               {UploadDoc.Kodestring OUT}, {DocAPI.ResponseMsg OUT});
  COMMIT;
END;
```

`[terverifikasi]` Dua hal terbaca langsung dari tanda tangan itu:

1. Procedure menerima **`OperatorID.pyUserIdentifier`** — identitas pengguna ikut masuk, sehingga
   wewenang atas berkas **mungkin** ditentukan per pengguna. Tanpa body procedure hal ini tidak
   dapat dipastikan (**OQ-002**).
2. Procedure menerima **`UploadDoc.App`** — token tampaknya dicakup per aplikasi, bukan global.

`[terverifikasi]` Ini **bukan milik Claim — Life**: 97 berkas di **15 modul** memakai jalur
penyimpanan yang sama, dan kelasnya `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` berada di luar kelas modul.
Penyimpanan berkas adalah **infrastruktur bersama korpus**.

## Considered Options

- **Tetap di Google Storage** — dipilih
- Pindah ke penyimpanan lain (S3-compatible, on-prem) — ditolak: berkas yang sudah ada harus ikut
  pindah, dan penyimpanan dipakai bersama 15 modul yang **tidak** ikut dimigrasi pada gelombang ini

## Consequences

- Sistem baru harus dapat memanggil `POOLDATA.GET_TOKEN_STORAGE` — **ketergantungan Oracle kedua**
  di luar penomoran klaim (**ADR-0006**). Keduanya bukan sekadar penyimpanan data; keduanya
  **logika** yang berada di database.
- Kontrak procedure harus diperoleh dari **DBA** sebelum unggah/unduh dapat dispesifikasikan:
  masa berlaku token, cakupan (`App`), dan apakah `pyUserIdentifier` menentukan wewenang.
- Unggah berkas termasuk **efek keluar asinkron** (**ADR-0008**) — kegagalannya dicatat, tidak
  memblokir alur klaim, dan diantre ulang.
- Karena penyimpanan dipakai 15 modul yang masih di Pega, selama masa transisi **dua aplikasi**
  membaca dan menulis berkas yang sama. Ini berbeda dari data klaim, yang dipindah penuh tanpa
  koeksistensi (**ADR-0009**).
- Alamat layanan penyimpanan **di-lookup runtime** dari `M_LINK_SERVICE` (**ADR-0013**, menggantikan
  ADR-0004); kredensial dan koneksi database tetap env var. Gerbang lingkungan **ADR-0005**.
## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-002** | Body `POOLDATA.GET_TOKEN_STORAGE` tidak ada di korpus — kontrak token tidak terbaca. **Pemilik: DBA** |
| **OQ-018** | Apakah bucket/proyek berbeda antara production dan dev belum dinyatakan |
| **OQ-047** | Isi `M_LINK_SERVICE` — daftar endpoint yang sebenarnya tidak diketahui |
