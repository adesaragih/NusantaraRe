# PROMPT — GILIRAN 3 · SESI A **Claim Life** *(folder `OUTPUT_HASIL_RNM`, cabang `main`)*: **paket 0 penyimpanan (019 vs `Update_T_Storage_SQL`) → A4 kode → sensus paritas akhir Claim Life → av saat PremiumList tiket 04 menyatu**

> Tiga sesi berjalan **bersamaan** mulai giliran ini *(induk §0.4)*: sesi A = brief ini; sesi B = `…-GILIRAN-3-PREMIUMLIST.md`
> di `.worktrees\premiumlist-life`; sesi C = `…-GILIRAN-3-KOMITE.md` di `.worktrees\komite-claim-life`. Sesi A **tidak**
> menyentuh berkas `polis_*`/`komite_*` dan **tidak** mengerjakan tiket modul lain. Bila work owner tetap memakai **satu**
> sesi: tempel brief A, lalu B, lalu C — semuanya di `main` *(induk §0.3)*. Brief GILIRAN-1 dan GILIRAN-2 tetap berlaku.
> Mekanisme giliran = lanjutan 8 §1; batas paket = titik lapor sah bila pohon bersih.

---

## 0. KEADAAN SESUDAH `4f40272` — DIVERIFIKASI ULANG

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| paket 0 `06bf711`, 1 `cf35d55`, 2 `4f40272`; pohon bersih | ya; kedua worktree modul **di-fast-forward asisten ke `4f40272`** | ✅ |
| Go 373 · 0 · 37; JS 254; gofmt, vet, vet `-tags db`, tsc, build bersih | dijalankan ulang: **373 PASS · 0 FAIL · 37 SKIP**; vitest **254**; **51** modul; semuanya bersih | ✅ |
| pengenal dokumen 17 angka > `Number.MAX_SAFE_INTEGER` → menyeberang sebagai teks | `InsertDocument_Act.xml` b648 `yyyyMMddhhmmssSSS` = 17 angka ≈ 2,0e16 > 9,0e15; `models.IDDokumenBaru` mengembalikan `string` | ✅ temuan sah, dan penting |
| gerbang `.STS_REJECT=='1' \|\| .STS_REJECT=='2'` **tujuh** kali | b2628 `.ADMIN_NOTES`, b4682, b5059, b5870, b6152, b7335 `.RECOMMENDATION`, b15234 `.NOTES` | ✅ ralat atas brief 2 diterima |
| `MedicalCheckClaimLife` **membuka** Detail lewat flow action | b17416 dan b22925 `ViewClaimDetailLifeGCNM` | ✅ ralat diterima |
| `SetSTS_Reject` = activity pra-muat, bukan dipanggil rute | `ClaimLifeDetailGCNM.xml` b3224 `pyDeferLoadRetrievalActivity SetSTS_Reject`, wilayah `S5` b3225 | ✅ ralat diterima; invariannya ditiru |
| migrasi 019 `T_CLAIMLF_STORAGE` | cermin `POOLDATA.T_STORAGE_IMAGE` `[katalog DEV]` = `IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE, TANGGAL_UPLOAD` — 019 memuat **tujuh**, **tanpa `TANGGAL_UPLOAD`**; kolom itu ditulis `RDBList/Update_T_Storage_SQL.xml` *(satu-satunya berkas yang menyebutnya)* | ⚠️ **paket 0** |
| `T_MIGRASI` DEV | masih `016` — work owner belum `-migrate` | ⚠️ di luar sesi |
| `UNGGAHAN_DIR` | ada di `.env.example`; **belum** ada di `.env` work owner | ⚠️ di luar sesi |
| nol kebocoran | seluruh berkas yang berubah | ✅ |

## 1. PAKET 0 — penyimpanan: 019 lawan `Update_T_Storage_SQL`, dan `GenerateImageID_SQL`

- Baca `Update_T_Storage_SQL.xml` utuh: kolom yang di-`SET` *(termasuk `TANGGAL_UPLOAD`)* dan `where imageid = {UploadDoc.ImageID}`.
  Bila `TANGGAL_UPLOAD` ditulis rule itu → migrasi **`020_kolom_tanggal_upload_storage.sql`** *(+ down)* menambah
  `TANGGAL_UPLOAD DATE` ke `T_CLAIMLF_STORAGE`, dan efek `storage-unggah` mengisinya saat pelaksana menandai selesai.
  Bila tidak → catat di tiket 14 kenapa kolom warisan itu tidak dibawa.
- `GenerateImageID_SQL.xml` = `SELECT STANDARD_HASH(…)`: pastikan `ImageID` di Go lahir dari rumus yang sama *(fungsi
  hash + masukan yang sama, huruf besar/kecil sama)*, bukan pengenal lain; uji dengan contoh terhitung dari Oracle
  `[db]` dan literal tetap `[murni]`.
- Ralat bertanggal di tiket 14 dan PARITAS; commit `claim-life: A3 — Dokumen (3), TANGGAL_UPLOAD dan IMAGEID sesuai RDB`.

## 2. A4 — MIGRASI DATA, KODE SAJA *(brief 10 §7; lanjutan 5 §4)*

`BongkarBarisLama` + pengisian `TAHAP`/`TGL_CREATE`/`STATUS_WORK` baris lama *(at/au/bb)*; dokumen warisan →
`T_CLAIMLF_DOCUMENT` lewat kunci kelompok `KATEGORI_1 = DL-…` *(OQ-J: yang tidak cocok dilaporkan, bukan dipaksa)*;
diagnosa warisan **tidak ada** sumber tabelnya *(hidup di BLOB Pega)* → dicatat, bukan dikarang; laporan Temuan; uji
bertag `db` SKIP; **tidak** dijalankan terhadap DEV. Commit `claim-life: A4 — migrasi data (kode)`.

## 3. SENSUS PARITAS AKHIR CLAIM LIFE — supaya "yang di-skip" terlihat, lalu dibangun

Tabel **satu baris per rule** untuk **seluruh** isi korpus `D:\XML\RNM_BRD\Claim Life\` — folder `Activity`, `FlowAction`,
`Section`, `Harness`, `RDBList`, `ReportDefinition`, `DataTransform`, `DecisionTable`, `DataPage`, `When`, `ConnectREST`,
`SystemSettings` — kolom: rule · dipanggil dari *(berkas + baris)* · padanan Go/React *(berkas:baris)* · keadaan
*(**ada** / **tidak perlu** + sebab / **celah**)*. Ditulis ke `PARITAS-LAYAR-DAN-AKSI.md` bab baru *"Sensus akhir 28-09-2026"*.
Setiap **celah** yang punya pemanggil di XML **dibangun di giliran ini** *(satu commit per kelompok celah)*; yang tidak
punya pemanggil dicatat sebagai residu dengan bukti *(seperti `setDetailClaim_act`)*. Tidak ada baris "nanti".

## 4. av — SESUDAH PREMIUMLIST TIKET 04 MENYATU KE `main`

Asisten menyatukan cabang PremiumList ke `main` sesudah memverifikasi lognya; saat `repository.PolisRingkas` *(pl4)* ada di
`main`: `PanelDataPolis` sepuluh medan dari `T_PREMIUM_LIST`, penanda "menunggu modul PremiumList Life" dicabut, ambang
**ba** berkunci `ProductNameID` disambungkan. Bila belum ada saat §3 selesai → laporkan; jangan membuat pembaca sendiri.
Commit `claim-life: av — PolicyDataLife dari PremiumList Life`.

## 5. ATURAN · LAPORAN · TELEMETRI

XML sebagai pohon, path + baris di tiket; ralat bertanggal bila tiket dibantah; migrasi Claim Life `020`+ hanya dari
keputusan tercatat; `-migrate` **tidak** dijalankan executor; nol kebocoran; satu pesan pada batas paket dengan tabel
**paket → tombol XML → rute/kontrol**, angka uji tiap commit, bab **TELEMETRI EKSEKUSI** per paket *(lanjutan 8 §6–§7)*.

---

*Disusun 28 September 2026 sesudah verifikasi `4f40272` (uji, vet, vet -tags db, gofmt, tsc, build dijalankan ulang),
pembacaan ulang `ClaimLifeDetailGCNM.xml` (tujuh gerbang, b3224), `MedicalCheckClaimLife.xml` b17416, `InsertDocument_Act.xml`
b648, migrasi 019, `Insert_T_Storage_SQL.xml`, `Update_T_Storage_SQL.xml`, `GenerateImageID_SQL.xml`, dan katalog DEV
(`T_STORAGE_IMAGE`, `T_MIGRASI` — agregat saja).*
