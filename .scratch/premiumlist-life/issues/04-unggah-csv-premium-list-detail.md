# 04: Unggah CSV premium list detail — staging, validasi, tinjau, simpan

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 03 (premium list detail — unggahan mengisi struktur yang dibentuk di sana)

## Hasil & nilai pengguna

Sebagai **inputor Life**, saya ingin mengunggah berkas CSV berisi ratusan baris peserta dan
**melihat dulu** apa yang lolos dan apa yang ditolak beserta **nama kolomnya**, sebelum apa pun
tersimpan permanen — supaya satu baris rusak tidak mencemari premium list dan saya tidak perlu
menebak kolom mana yang salah. *(User story 17–24 di spec)*

## Area codebase

`internal/handlers` (endpoint unggah, endpoint tinjau, endpoint simpan permanen), `internal/services`
(mesin validasi per kolom; normalisasi uang; deteksi duplikat), `internal/repository` (tabel staging
+ pembersihannya), `frontend/` (form unggah, tabel hasil validasi berlabel kolom, tombol simpan).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `UploadCSVLifePremium_Act` | `@BASECLASS` / `UPLOADCSVLIFEPREMIUM_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/UploadCSVLifePremium_Act.xml` | impor berkas — **rule base-class, dipakai bersama** |
| `ValidasiUploadPL_act` | `ASM-FW-GISFW-WORK-LIFE` / `VALIDASIUPLOADPL_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/ValidasiUploadPL_act.xml` | **43 langkah** validasi |
| `InsertDataUploadLife` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!INSERTDATAUPLOADLIFE` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertDataUploadLife.xml` | `insert into POOLDATA.M_TEMPUPLOADLIFE (…)` — **staging** |
| `DeleteTempUploadDataLife` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!DELETETEMPUPLOADDATALIFE` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/DeleteTempUploadDataLife.xml` | `… M_TEMPUPLOADLIFE where idpega = {pyWorkPage.pyID}` |
| `CekDoubleInsured` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!CEKDOUBLEINSURED` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/CekDoubleInsured.xml` | duplikat peserta + jumlah retensi ceding |


⚠️ **Di luar cakupan tiket ini:** jalur unggah CSV **endorsement**
(`Endorsement Life/Activity/UploadCSVEDMLifePremium_Act.xml`, `@BASECLASS` /
`UPLOADCSVEDMLIFEPREMIUM_ACT`, dan `Endorsement Life/Activity/SaveCSVEDMLife.xml`,
`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE`). `[keputusan work owner]` **Endorsement Life
adalah konteks terpisah** — lihat `.scratch/endorsement-life/`.

`[terverifikasi]` **Staging nyata**: `POOLDATA.M_TEMPUPLOADLIFE`, berkunci `idpega`, dibersihkan per
case. Validasi berjalan **atas staging**, bukan atas tabel permanen.

`[terverifikasi]` **`CekDoubleInsured`** mendeteksi duplikat dengan `upper(INSURED)` + `DOB`,
mengembalikan `min(NO)` baris pertama, dan menjumlahkan `to_number(CEDING_RETENTION)` per peserta.
⚠️ Variabel `v_Count2` dan `v_PLRetensiCeding` **dideklarasikan tetapi tidak pernah diisi** — selalu
`NULL`, dinetralkan `nvl(…,0)`. Sisa pemeriksaan kedua yang dicabut; **jangan** direplikasi.

`[terverifikasi]` **33 kolom uang** divalidasi (`<pyStepsDescription>` per langkah): `NET_PREMIUM`,
`GROSS_PREMIUM`, `SHARE_NUSANTARA_RE`, `SUM_INSURED`, `CEDING_RETENTION`, `SUM_REASURED`, `CLAIM`,
`TAX`, `BROKERAGE_FEE`, `OVR_COMM`, `PROF_COMM`, `EM_PERCENT`, `COMM`, `FLEET_DISCOUNT`, seluruh
kelompok `*_REFUND`, seluruh kelompok `*_RETRO`, `*_REFUND_RETRO`, dan `CLAIM_AMOUNT`.

`[terverifikasi]` Validasi non-uang: `CERTIFICATE_NO` tidak boleh duplikat; `NAME_OF_INSURED` dan
`POLICY_HOLDER` **harus ada di `M AGENT`**; `DOB`, `BEGIN_DATE`, `EXPIRED_DATE` format `dd/mm/yyyy`;
`PLAN`, `CURRENCY` wajib; `MEDICAL_STATUS` salah satu dari `FCL` / `M` / `NM`; lampiran wajib
("BELUM ADA LAMPIRAN").

⚠️ `[terverifikasi]` **Normalisasi desimal Pega berbahaya.** Langkah "Rubah decimal dari koma jadi
titik" (baris 12480) memakai `@replaceAll(.KOLOM, ",", ".")` atas tiap kolom uang — **mengganti
setiap koma**, tanpa membedakan pemisah desimal dari pemisah ribuan. Nilai `1,234,567.89` menjadi
`1.234.567.89`, yang bukan angka. **Jangan** direplikasi apa adanya.

## ADR terkait

**ADR-0003** (uang non-float, desimal presisi arbitrer — inti tiket ini), **ADR-0010** (penyimpanan
berkas tetap Google Storage untuk lampiran), **ADR-0007** (jejak audit unggahan).

## Acceptance criteria

- [ ] Berkas yang diunggah masuk **staging** lebih dulu; kegagalan validasi **tidak** menyentuh tabel
      permanen mana pun. *(AC 18 spec)*
- [ ] Pengguna dapat **meninjau** hasil unggahan — baris lolos dan baris ditolak — sebelum menyimpan
      permanen. *(AC 19 spec)*
- [ ] Setiap penolakan menyebut **nama kolom** yang salah dan **nomor baris** CSV-nya. *(AC 17 spec)*
- [ ] Nilai uang di-parse dengan **format yang dinyatakan eksplisit** (pemisah desimal dan pemisah
      ribuan ditentukan, bukan ditebak). `1,234,567.89` dan `1.234.567,89` **tidak** boleh keduanya
      diterima diam-diam sebagai angka yang sama. Test wajib memuat kasus pemisah ribuan.
- [ ] Nilai uang **tidak** melewati `float` pada tahap mana pun — parse langsung ke desimal presisi
      arbitrer. *(AC 14–15 spec; **ADR-0003**)*
- [ ] Nilai uang yang lolos validasi, dibaca kembali dari staging, **identik** dengan yang diunggah —
      tidak ada pembulatan diam. *(AC 16 spec)*
- [ ] Duplikat peserta terdeteksi dengan pembandingan nama **tanpa peduli huruf besar/kecil** +
      tanggal lahir, dan pesannya menunjuk baris pertama yang bentrok.
- [ ] `CERTIFICATE_NO` ganda di dalam satu berkas ditolak.
- [ ] `NAME_OF_INSURED` dan `POLICY_HOLDER` diverifikasi terhadap master agen; baris yang tidak
      ditemukan ditolak dengan pesan yang menyebut nilainya.
- [ ] Tanggal diverifikasi berformat `dd/mm/yyyy`; `MEDICAL_STATUS` hanya menerima `FCL`, `M`, `NM`;
      `CURRENCY` wajib ada.
- [ ] Unggahan tanpa lampiran ditolak.
- [ ] Staging dibersihkan per case setelah simpan permanen **maupun** setelah pembatalan — tidak ada
      sisa baris menggantung milik case lain.
- [ ] Mengunggah berkas kedua atas case yang sama **mengganti** isi staging, tidak menumpuk.

## Blocker

⚠️ **OQ-069 terbuka** — pesan Pega `"NET PREMIUM HARUS ADA DAN LEBIH BESAR DARI GROSS PREMIUM"`
menjanjikan aturan yang **tidak pernah diperiksa**: satu-satunya precondition atas `NET_PREMIUM` di
seluruh berkas adalah `@PropertyHasValue(.NET_PREMIUM)`, dan **tidak ada** perbandingan terhadap
`GROSS_PREMIUM`. Arah perbandingannya pun janggal — lazimnya net lebih **kecil** dari gross.

**Selama OQ-069 terbuka:** tiket ini memeriksa **keberadaan** `NET_PREMIUM` saja, persis seperti
korpus. **JANGAN** mengarang aturan perbandingan. Pesan Pega dicatat apa adanya di kode sebagai
rujukan, dengan penunjuk ke OQ-069.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```
