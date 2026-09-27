# PROMPT — lanjutan 10 *(sesi Claim Life, cabang `main`)*: **Register(2) → Detail & Tutup → Dokumen → Medis → Akseptasi → Komite-handoff → A4**, satu giliran; berjalan berdampingan dengan dua sesi modul lain

> Baca dulu `PROMPT-INDUK-TIGA-MODUL.md` *(tata letak, pembagian, kontrak)*. Brief modul dan lanjutan
> 1–9 tetap berlaku. Mekanisme giliran = lanjutan 8 §1. Berhenti "agar tidak mengulang kesalahan"
> **bukan** berhenti sah — yang sah hanya gerbang yang manusia sendiri harus buka.

---

## 0. KEADAAN SESUDAH `cfc4824` — DIVERIFIKASI ULANG

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 4 commit *(F0.6, ralat av, Outstanding aw, + docs)* | `git rev-list --count` = 4; 19 + 3 + 12 berkas | ✅ |
| Go 274 · 0 · 34; JS 150; 46 modul | dijalankan ulang: **274 PASS · 0 FAIL · 34 SKIP**; vitest **150**; **46** modul | ✅ |
| dropdown rujukan dibuang tiga lapis; `Masuk` dibuang | nol `rujukan` di handlers/repository/tsx; `pages/` tanpa `Masuk` | ✅ |
| ralat av di tiket 02, OQ-C, PARITAS 226 | tiket 02 baris 495–517; `OQ-untuk-tim.md:134`; PARITAS 226 `⛔ RALAT` | ✅ |
| `VITE_STUB_PELAKU`/`VITE_STUB_PERAN` | `.env.example` 31–32 | ✅ |
| nol kebocoran | seluruh berkas yang berubah | ✅ |

Yang belum *(laporan executor, benar)*: §3 Register(2) empat activity + `UploadCSV_ClaimLife`; §5
Detail & Tutup. Keduanya **dikerjakan sekarang**, lalu tiga kelompok berikutnya.

Catatan untuk layar: backend di mesin work owner sempat berjalan **tanpa `.env`** dan tabel belum
dimigrasi — bukan cacat kode. Perintahnya ada di induk §0.

---

## 1. REGISTER(2) — commit `claim-life: A3 — Register (2)`

Bekal langkah tiap activity ada di **lanjutan 9 §3** *(dibaca dari blok `pySteps` utuh, bukan grep
maju)*. Ringkasan yang harus jadi kode dan uji:

| Aksi | Bangun |
| --- | --- |
| `SaveInsuredClaim_Act` *(`Select Insured` 20008/20060)* | memindahkan peserta terpilih ke daftar klaim **di layar**; `AGE` = `AGE` → `ENTRY_AGE` → `CURRENT_AGE`; `SHARE_NUSANTARA_RE` fallback `SHARE_NUSANTARA_RE_GROSS` bila `"0"`/`""`; uang tidak dihitung ulang di layar *(ADR-U-0003)* — bandingkan dengan `kolomSalin` tiket 02/03: yang sudah disalin backend saat `POST` **tidak** diduplikasi |
| `DeletePesertaClaimLife` | sebelum `POST` = hapus dari pilihan; klaim yang sudah ada di tahap Input Register = `DELETE /api/klaim-life/{id}/peserta/{pesertaId}` bergerbang tahap + pemegang, kaskade baris adjustment peserta itu, jejak; uji murni + handler |
| `LoadDataPesertaSpesifik_Act` *(`Find Insured` dengan nama)* | `GET /api/peserta-life?pl=&nama=` — nama di-uppercase, `LIKE` berbatas, **wajib** bersama `pl`; SQL dari `GetPesertaClaim_sql1.xml`; +7 jam = zona Pega, tidak ditiru, dicatat |
| `SelectAllClaimLife_act` *(`Select All` 24489)* | pilih/lepas seluruh baris hasil pencarian; uji |
| `UploadCSV_ClaimLife` *(flow action 12)* | baca `UploadCSVClaimLife_Act` sebagai pohon: bila menulis peserta → rute unggah berbatas + validasi kolom + fixture `UJI-*`; bila hanya UI → dinyatakan di PARITAS |
| PARITAS | baris 74, 81, 83, 84, 129 diperbarui; `ValidasiClaimReceived_Act`/`ValidasiSTNC_Act` **dipindah** ke kelompok Detail *(sheet baris 115, 118)* |

## 2. DETAIL & TUTUP — commit `claim-life: A3 — Detail & Tutup`

`ClaimLifeDetailGCNM` sebagai **pohon** *(himpunan field, `<pyCondition>`, tombol → activity)*:
`Save Adjustment` 22590 *(rute akseptasi ada)*; `Diagnose_Harness` popup 5080/5224 → `BelumTersedia`
bernama sampai kelompok Medis; `ShowEditClaimLife` + `EditDateClaimLife_Section` *(rute tanggal
kejadian ada — `ValidasiDOL_Act`)*; `ValidasiClaimReceived_Act` dan `ValidasiSTNC_Act` *(lanjutan 9
§3: `GetProductName` → `MAXEXPIREDCLAIM`/`MAXDATARECEIVE`; selisih hari; `.MAXCLAIM_RECEIVED`/`.STNC`;
pesan `"Max Claim invalid"`/`"STNC invalid"` VERBATIM; kolom `STNC_CLAIM` ada)*; `SetSTS_Reject`;
`SetIndexAdjustmentList`; `CloseClaim` *(`CloseClaim_Section` 1028, `ProtectCloseClaim_act`)*;
`DocumentLife` tampil daftar *(unggah/hapus = §3)*. Rute yang ada *(tolak, putaran, akseptasi, hapus +
dampak)* mendapat kontrolnya dengan gerbang yang **services** tegakkan.

## 3. DOKUMEN — commit `claim-life: A3 — Dokumen`

`AttachDocumentLife` + `AttachDocScreenLife` *(kategori wajib **ar1** `CATEGORY_ATTACH_CLAIMLIFE`)*,
`NewAttachLife`, `SaveAttachLife` → `InsertDocument_Act` → `InsertGoogleStorage_Act`
*(`Insert_T_Storage_SQL`, `GetTokenStorage_SQL` → token **an** dari env salt, `ServiceGoogle` →
resolver `M_LINK_SERVICE` Google/upload, `GenerateImageID_SQL`, `GetAppName_SQL`, decision table
`GetMimeType` → **tabel data** di `models`)*; `DownloadDocumentClaim` → `GetUrlGoogleStorage_Act`
*(geturl)*; `ConfirmDeleteAttachment` → `DeleteDocument_Act` → `DeleteGoogleStorage_Act` *(delete)*.
⛔ Endpoint penyimpanan sungguhan **tidak** dipanggil tanpa persetujuan manusia: jalur keluar lewat
outbox `T_LOG_SERVICE_RNM` *(aq)* dan pengirim yang di DEV memakai **stub** yang mencatat, bukan mengirim.
Berkas diunggah lewat rute berbatas ukuran dan tipe *(dari `GetMimeType`)*, disimpan sementara di
folder aplikasi yang dinyatakan di `.env`, **tidak** ke korpus.

## 4. MEDIS — commit `claim-life: A3 — Medis`

`MedicalCheck` + `MedicalCheckClaimLife`; `Diagnose_Harness`/`Diagnose_Section` *(popup; `SearchDiagnose_act`,
`SetDisease`; sumber `DISEASE_LIFE` `[data DBA]` 3 kolom `ID`, `ICD_CODE`, `DISEASE`, **97.586** baris →
pencarian **berbatas** per `ICD_CODE`/`DISEASE`, tombol `Choose` 8 ×; hasil ke `.DiagnoseList` class
`Data-DiagnoseLife` — keputusan **al** dari bukti di sini)*; `SendtoAdmin` dari Medical Check
*(posisi Medical → `SendtoAdmin_Act` **berjalan** sesuai XML; rute `tahap` yang ada)*;
`SendtoMedical` *(`SendtoMedical_Act`, dari SPV)*.

## 5. AKSEPTASI — commit `claim-life: A3 — Akseptasi`

`AkseptasiClaimLife` + `InputAkseptasiClaimLife` *(`Send Back to Medical` 20467, `Send Back to
Admin` 20221, `Close Claim`, gerbang `Type` TP/TR; `pyEditAction` `ViewClaimDetailLifeGCNM`
bergerbang `pyWorkPage.Save = 1`)*; `SaveAdjustment_Act` *(ada: nomor akseptasi sequence, `IsCheck`,
`<LAST>`; dua cacat Pega yang sudah dilaporkan tetap dilaporkan, tidak ditiru)*.

## 6. KOMITE-HANDOFF — commit `claim-life: A3 — Komite`

`ClaimComite` + `Committe_Life` popup *(`Send Claim to Committee`, `Cancel`, medan `Date`, `PIC`,
`Remarks`, `TanggalComitee`, `KomiteList` grid `Committee Name`/`Date Approve`/`Comment`/`Status`)*
→ rute `serahkanKomite` yang ada + roster `FilterEmailKomiteWithLimit` *(A2)*; `SendEmailKlaimLF`
→ outbox. **Jangan** menyentuh berkas `komite_*` milik sesi Komite; kontraknya tetap A2.

## 7. A4 — migrasi data *(lanjutan 5 §4)* bila giliran masih hidup

`BongkarBarisLama` + pengisian `TAHAP`/`TGL_CREATE` baris lama *(at/au)*; laporan Temuan; tidak
menyentuh tabel warisan selain membaca.

## 8. KOORDINASI DENGAN DUA SESI LAIN

Tidak menyunting `polis_*`/`komite_*`; perubahan kode bersama hanya aditif dan dilaporkan; nomor
migrasi `017`–`029`; sesudah PremiumList Life di-merge, paket kecil berikutnya menyambungkan
`PolicyDataLife` ke `T_PREMIUM_LIST` *(av)* — **bukan** di giliran ini.

## 9. LANGKAH 0 · LAPORAN · TELEMETRI

Langkah 0: `git add PROMPT-INDUK-TIGA-MODUL.md PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-10.md
PROMPT-IMPLEMENTASI-MODUL-PREMIUMLIST-LIFE.md PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-LIFE.md
.gitignore` → commit `docs: induk tiga modul; brief Claim Life 10, PremiumList Life, Komite Claim
Life` → uji hijau *(274 · 34 SKIP · 150 JS · 46 modul)*. **Langkah 0.5 — sebelum paket mana pun dan
sebelum worktree dibuat**: keputusan **ax** induk §0.1 **sudah diterapkan asisten di disk** *(27-09-2026:
`015_t_log_service_rnm.sql` + `_down` disunting di tempat dan diganti nama; `repository/efekkeluar.go`,
`services/antrean.go`, `migrasi_test.go`; blok ralat di tiket 12; gofmt/vet/uji hijau 274 · 0 · 34)* —
**sudah di-commit asisten sebagai `e3d537a`** `claim-life: ax — T_EFEK_KELUAR menjadi T_LOG_SERVICE_RNM`; executor hanya **memverifikasi** *(`git log -1`, `go test -tags=db ./...` hijau)* dan memakai nama baru sejak paket pertama. Laporan akhir dan telemetri persis lanjutan
8 §6–§7, ditambah tabel **kelompok → tombol XML → rute/kontrol** per kelompok.

---

*Disusun 27 September 2026 sesudah verifikasi `cfc4824` dan pembacaan ulang XML yang tercantum di
lanjutan 9; bahan Dokumen/Medis/Akseptasi/Komite dari lanjutan 4 §7, lanjutan 5 §2–§3, sensus paritas,
dan katalog DEV (`DISEASE_LIFE`, `CATEGORY_ATTACH_CLAIMLIFE`, `M_LINK_SERVICE` — definisi dan cacah).*
