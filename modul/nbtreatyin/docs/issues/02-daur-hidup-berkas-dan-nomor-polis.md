# 02: Daur hidup berkas realisasi dan nomor polis — satu nomor, sekali, tanpa bentrok

**Status:** selesai — tertahan hanya pihak luar: AC 31 uji `db` ditulis, belum dijalankan (K11 skema uji Oracle) *(putaran 2, konsolidasi P10 04-10-2026 — rincian `docs/HASIL-IMPLEMENTASI.md` bab 9; semula: sebagian, implementasi 2026-10-03; awalnya ready-for-agent)*
**Blocked by:** 01
**Menutup:** AC 31 · 59 · 73 · 74 *(4 AC)* — US 1 · 4 · 5

## Hasil & nilai pengguna

Hari ini admin treaty membuka penawaran yang sudah disetujui dan melengkapinya menjadi realisasi.
`[terverifikasi]` Nomor polis dibentuk dari sebuah deret, dan ⚠️ salah satu bagiannya diambil dari
slot parameter generik yang **tidak terlihat diisi** di aktivitas mana pun.

Sesudah tiket ini, sebuah realisasi treaty **lahir dari penawaran yang disetujui**, mendapat
**satu nomor polis yang tidak bentrok**, dan ⭐ pengguna **diperingatkan** ketika ia hendak membuat
penawaran yang sudah pernah ada.

## Area codebase

- Lapisan service: daur hidup berkas realisasi
- Lapisan service: pembentukan nomor polis
- Lapisan handler: peringatan penawaran ganda

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Titik masuk alur | `Flow\InputRealizationTreatyIn.xml` — `pyStartActivity` = `Start1` |
| Pembentukan nomor polis | `RDBList\GenerateNoPolicy.xml` |
| Peringatan penawaran ganda | `Activity\CheckDuplicateOffer.xml` langkah 6 dan 8, lewat `RDBList\GetCountClaim.xml` |
| Rantai pemanggilnya | `TreatyRealizationCheckXOLList` → `SetTreatyIn_Act` → `CheckDuplicateOffer` |

## ADR terkait

- **ADR-0006** — penomoran lewat deret basis data

## Acceptance criteria

- [x] **AC 73** — nomor polis memuat awalan tetap, penanda treaty, bulan-tahun, dan nomor urut
      berdigit tetap
- [x] **AC 74** — nomor dibentuk **sekali** per berkas; tidak berubah pada penyimpanan berikutnya
- [ ] 🟡 **AC 31** — dua berkas **tidak pernah** bernomor sama
- [x] **AC 59** — ~~peringatan muncul ketika jumlah berkas klaim terhubung **lebih dari nol**~~ ⇒
      RALAT putaran 2: submit admin yang menyetujui tertahan bila ada polis produksi serupa
      (`TreatyRealizationCheckDuplicate`) — lihat RALAT putaran 2 butir 1

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **P7** | asal salah satu bagian nomor polis **tidak terlihat diisi** | ⚠️ tidak menahan — polanya diketahui |

## Perintah verifikasi

1. Realisasikan dua penawaran berbarengan — ⭐ nomor polisnya **berbeda**.
2. Simpan ulang salah satunya — ⭐ nomornya **tidak berubah**.
3. Buat penawaran yang sudah pernah ada — ⭐ **peringatan muncul**.

## ⛔ RALAT implementasi 2026-10-03

1. **Rule pembentuk nomor polis.** Bunyi lama: *"Pembentukan nomor polis | `RDBList\GenerateNoPolicy.xml`"*.
   ⛔ Keliru: `GenerateNoPolicy` dipanggil `SaveJsonPolisTreatyIn_Act` langkah 4 dan **hasilnya tidak
   dipakai**. Nomor polis dibentuk `Activity\GeneratePolicyNoTreaty_Act` (langkah efektif 3-13, 25-30;
   14-24 berlabel `//`): `KODE_PRODUKSI(NONLIFE) + QR/QP/TP + ".T" + OJKBusinessID + "." + MM.YYYY +
   "." + urut5`, urut dari `PROC_GENERATE_SEQUENCE_NUMBER` (kini `inti/backend/penomor`). Langkah 28
   bersyarat `PolicyNo == ""` — nomor sekali (AC 74).
2. **P7 terjawab.** `InputData.CARI20` diisi langkah 7-9 (`DueTo` 1 → QR, 0 → QP, `ClaimType`
   "XOL Retro" → TP).
3. **AC 59 tidak dapat dipenuhi seperti tertulis.** `CheckDuplicateOffer` langkah 1-4 berlabel `//`
   dan dipanggil `SetTreatyIn_Act` (pembongkar JSON, tidak dimigrasi — AC 62) TANPA parameter, sehingga
   `GetCountClaim` selalu menghitung `masterid = NULL` — peringatan **tidak pernah menyala** di sistem
   lama. Yang dibangun: `TreatyRealizationCheckDuplicate` (pasca-submit admin, IsApproved 1).

## ⛔ RALAT putaran 2 (P4, 03-10-2026)

1. **AC 59 — bunyi baru.** Bunyi lama: *"peringatan muncul ketika jumlah berkas klaim terhubung
   **lebih dari nol**"* (= `Activity\CheckDuplicateOffer.xml` langkah 6-8, `RDBList\GetCountClaim`).
   Bukti XML bahwa bunyi itu tidak pernah terjadi: langkah 1-4 ber-`pyStepsBlockName` `//`; langkah 5
   `TreatyWarning.CAIREINSFACIN = Param.ID`; satu-satunya pemanggil, `SetTreatyIn_Act` langkah 14,
   memanggil **tanpa parameter** dengan `pyPassCurrentParameterPage = false` ⇒ `Param.ID` kosong ⇒
   `GetCountClaim` `where masterid = NULL` ⇒ `count = 0` ⇒ syarat langkah 7-8
   `HasilClaim.pxResults(1).CARI1>0` tidak pernah benar (bukti (a): cabang pesannya tidak terjangkau).
   Bunyi baru: **"Submit admin yang menyetujui (IsApproved 1) tertahan dengan pesan *Protect Duplicate
   Policy; data is similar to &lt;nopolis&gt; …* bila `TREATYINPRODUCTION` memuat polis serupa dan
   `ClaimType` bukan `XOL`"** — `Activity\TreatyRealizationCheckDuplicate.xml` langkah 1-6, dipanggil
   `InputPolicyTreatyInPost_Act` langkah 3 (syarat `.PolicyTreatyIn.IsApproved==1`), SQL
   `RDBList\TreatyRealizationCheckDuplicate.xml`. Dibangun: `repository.PolisSerupa` +
   `services.cekDuplikat`. Uji: `handlers/alur_test.go` TestPolisSerupaMenahan, `handlers/logika_test.go`
   TestCekPolisSerupaHanyaSaatAdminMenyetujui (admin menolak dan putusan atasan tidak diperiksa;
   pesan merangkai setiap nopolis). ⇒ AC 59 ✅.
2. **Hari tutup buku.** `Activity\InputPolicyTreatyInPre_Act.xml` langkah 9 ("Set Production Date kalau
   diatas tanggal 25", syarat `@substring(StatementDate,6,2)>25`) menanam 25;
   `Activity\GeneratePolicyNoTreaty_Act.xml` langkah 3-5.3 membaca `POOLDATA.TANGGAL_CLOSING`
   (`GETTanggalClosing_SQL` → `Local.TglProd`). Keduanya dibangun dengan hari dari
   `TANGGAL_CLOSING` (`inti/backend/penomor.HariClosing`, preseden WO PremiumList Life — penjaga
   `TestNolAmbangTutupBukuTertanam`), pembanding tetap `>`: hari = batas **tidak** digeser. Uji batas
   24/25/26 dengan batas tiruan 25: `models/tutupbuku_test.go`, `handlers/logika_test.go`
   TestTanggalProduksiDariTanggalClosing (batas 27 → hari 26 tidak digeser: bukan 25 tertanam).
3. **RALAT kode — jam ProductionDate.** Bunyi lama (`models.GeserTanggalProduksi`): *"tanggal 1 bulan
   berikut"* pukul 00:00:00. Langkah 9 menulis `@substring(ProductionDate,0,6) + "01" +
   @substring(ProductionDate,8)` — bagian waktu StatementDate **dipertahankan**. Diperbaiki.
4. ⚠️ Catatan, tidak ditiru: `@substring(StatementDate,6,2)` membaca hari dari teks internal DateTime
   Pega (GMT); sistem baru membacanya di zona jam aplikasi. Langkah 5.3 nomor polis membaca hari di
   `Asia/Jakarta` (ditiru).
