# PROMPT — lanjutan 12 *(sesi Claim Life, cabang `main`)*: **ralat OQ-H + total keenam (bc) → Close Claim menutup (bb) + DocumentLife → Dokumen → Medis → Akseptasi → Komite-handoff → A4**, satu giliran; dua sesi modul lain BELUM pernah mulai

> Baca `PROMPT-INDUK-TIGA-MODUL.md` *(§0.2 baru)*, lalu brief lanjutan 10 §3–§7 *(isi Dokumen, Medis,
> Akseptasi, Komite, A4 tetap berlaku)* dan lanjutan 11 §1 *(az, ba)*. Mekanisme giliran = lanjutan 8 §1.
> ⛔ *"Konteks sesi terbatas"* **bukan** berhenti sah: catat ke `LAPORAN-GILIRAN-F0.md`, lanjutkan
> *(lanjutan 3 §1 butir 4)*. Brief modul lain yang tertempel di sesi ini → jawab **satu baris**
> *"milik sesi worktree X"*, nol pekerjaan — itu yang sudah Anda lakukan dua kali, dan benar.

---

## 0. KEADAAN SESUDAH `9352e45` — DIVERIFIKASI ULANG

| Klaim laporan *(bagian 1–3)* | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 3 + 2 + 2 commit *(+2 docs)* | 9 commit `b18de5b`…`9352e45`, pohon bersih | ✅ |
| Go 300 · 0 · 34; JS 201; tsc; vet; 47 modul; nol migrasi baru | dijalankan ulang: **300 PASS · 0 FAIL · 34 SKIP**; vitest **201**; tsc bersih; vet bersih; **47** modul; migrasi tetap `001`–`016` | ✅ |
| `go.mod` go 1.22; `t.Context()` dibuang | `go 1.22`; nol `t.Context()` di `*_test.go` | ✅ |
| envelope `galat` dua sisi | `api.ts:183–196` membaca `galat`; `RegisterKlaim.tsx` hanya menyebut bentuk axios di komentar; `klien.test.ts` tanpa `"error"` | ✅ |
| `CloseClaim_Section` b1081 → `ProtectCloseClaim_act` b1101; `closeContainer` b1129; konfirmasi b499 | b1081 `Close Claim`, b1101, b1129, b499 `Are you sure want to Close Claim?` | ✅ |
| gerbang tutup: `.STS_REJECT!=1`, pesan, `CARI1`, `FinishAssignment` | langkah 1 `ProtectLife.CARI1=""` b264; langkah 2 perulangan **`pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail`** b379 *(kelas `Int-LIFE_PREMIUM_DETAIL` b388 = PESERTA)*; 2.1 prasyarat `.STS_REJECT!=1` b608 → `CARI1=1` b512, `local.errmsg = local.name + " is not approved yet, on list " + local.idx` b558; 2.2 prasyarat `CARI1==1` b756 → `Page-Set-Messages` b648; langkah 3 prasyarat `CARI1==""` b992 → `Call FinishAssignment` b838. Go memutar `klaim.Peserta` *(`services/tutup.go:59`)* — halaman yang sama | ✅ |
| `ClaimLifeDetailGCNM` kelas `Int-LIFE_PREMIUM_DETAIL` b84; total ber-TITIK | b84; `.TotalCedingRetention` b20636, `.TotalShareRNM` b20921, `.TotalSumInsured` b21208, `.TotalSumReasured` b21495, `.TotalShareRetro` b21782, `.TotalClaimAmount` b22069 | ✅ *(tetapi lihat §0.1 butir 2: ada **enam**, bukan lima)* |
| `CheckTotalAdjustmentClaim` dirujuk 10 ×, berkasnya nol | b20970, b21091, b21262, b21377, b21546, b21664, b21836, b21951, b22120, b22238; `find` + `grep -rl` atas **19 folder modul** korpus → hanya section itu | ✅ *(tetapi lihat §0.1 butir 1: kesimpulannya keliru)* |
| `IsCheck` = `"true"` dari XML | `SetIndexAdjustmentList.xml` b328; `SavePesertaClaim.xml` b1812, b1997, b2413 *(+ b2555, b2704, b3631, b3919)*; `models.PenandaDipilih = "true"` `klaimlife.go:369` | ✅ |
| `SetSTS_Reject` halaman langkah `.DiagnoseList` | b241 `pyStepsObjectName .DiagnoseList`; b257 `.STS_REJECT = Primary.STS_REJECT` | ✅ |
| `SetIndexAdjustmentList` sudah ditiru | `adjustment.go:112`, `WarisiKolom` `:58` | ✅ ralat diri diterima |
| `Find Disease` → `BelumTersedia` bernama | b5061 → `showHarness` `Diagnose_Harness` b5081/b5225; `KlaimLife.tsx:436` | ✅ |
| pesan gerbang memakai nomor sertifikat, bukan nama | penyimpangan sadar, tercatat PARITAS 130 | ✅ diterima |
| nol kebocoran | seluruh berkas yang berubah | ✅ |

### 0.1 Tiga hal yang laporan LEWATKAN — paket 0 giliran ini

**1. OQ-H keliru pada kesimpulannya.** Laporan menulis *"kelima total tidak dihitung; baris mana
yang ikut belum terjawab"*. XML menjawab: totalnya **dihitung oleh dua activity yang ADA di korpus**.

| Activity | Pemanggil *(tombol → activity)* | Langkah yang menghitung |
| --- | --- | --- |
| `Activity/SavePesertaClaim.xml` | `InputRegisterClaimLife.xml` `Submit` b27369 → b27393 *(dan b27574)* | langkah 8 b4002: perulangan **`pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail`** b4008; reset enam `local.Total*` = `0` b4027–b4159; `Local.IndexPremium = .pxListSubscript` b4180. Sub-langkah 8.1 b4221: perulangan **`.AdjustmentList`** b4226 **tanpa prasyarat** — `local.TotalCedingRetention = .CEDING_RETENTION + local.TotalCedingRetention` b4246, `.SHARE_NUSANTARA_RE` b4292, `.SUM_INSURED` b4312, `.SUM_REASURED` b4332, `.SHARE_RETRO` b4352, `.CLAIM_AMOUNT` b4372; lalu total berjalan ditulis pada **baris adjustment** b4391–b4491 dan `.IndexPremiumList` b4511. Sub-langkah 8.2 b4592: total akhir ditulis pada **halaman PESERTA** — `.TotalCedingRetention` b4616, `.TotalShareRNM` b4662, `.TotalSumInsured` b4682, `.TotalSumReasured` b4702, `.TotalShareRetro` b4722, `.TotalClaimAmount` b4742 |
| `Activity/SaveOutStandingLife_Act.xml` | `InputOSClaimLife.xml` `Save to RNM` b21102 → b21126 *(dan b21234)* | langkah 23 b10638: perulangan dan rumus yang **sama** — reset b10663–b10794, sub-langkah 23.1 `.AdjustmentList` b10841 *(b10860 dst.)*, sub-langkah 23.2 b11067 menulis `.Total*` peserta b11091–b11197 |

Jawaban *"baris mana yang ikut"*: **seluruh baris `.AdjustmentList` peserta itu**, tanpa memandang
`STS_REJECT`. Yang hilang hanyalah `CheckTotalAdjustmentClaim` — pemanggil saat medan di-refresh di
layar Detail; **nilainya** datang dari kedua activity di atas. Executor membaca b4002–b4800 dan
b10638–b11250 **utuh** dan mencatat prasyarat langkah 8 / langkah 23 sendiri *(asisten tidak
membacanya sampai selesai)*.

→ **Keputusan bc** `[DIPUTUSKAN — bukti di atas; veto work owner]`: keenam total dihitung di
`services` dengan rumus persis *(jumlah SELURUH baris adjustment peserta)*, disajikan pada pembacaan
Detail `GET /api/klaim-life/{id}` per peserta, **tidak disimpan** *(nol kolom baru; `T_CLAIMLF_ADJUSTMENT`
dan `T_CLAIMLF_PREMIUMLIST_DETAIL` tidak punya `TOTAL_*`)*. Penyimpangan sadar yang **dicatat**:
Pega menghitung saat `Submit`/`Save to RNM`, sehingga layar Pega dapat basi sesudah putaran atau
akseptasi sampai `Save to RNM` ditekan lagi; di Go tidak pernah basi. Uji: contoh terhitung dengan
literal bebas *(bukan rumus ulang)*; baris `STS_REJECT = 2` **ikut** dijumlah *(itu yang XML lakukan)*;
peserta tanpa baris adjustment → nol. Ralat bertanggal di: `OQ-untuk-tim.md` OQ-H *(blok
"Ralat kedua" — kesimpulan dicabut, pertanyaan ke pemilik ekspor tinggal "kenapa rule refresh-nya
tidak ikut diekspor")*, PARITAS baris 18, komentar `labels.ts:291–296`, `LAPORAN-GILIRAN-F0.md`, dan
tiket 03 *(bab Pembacaan ulang XML — `SaveOutStandingLife_Act` adalah miliknya)*.

**2. Total keenam terlewat.** `Total Ceding Retention` b20629 / `.TotalCedingRetention` b20636 ada di
section *(tanpa refresh-activity, sebab itu tidak ikut hitungan "10 pemanggilan")*, dihitung b4246
dan b4616. Tidak ada di `labels.ts` maupun `PanelTotalPeserta`. Tambahkan; uji cacah label = **6**.

**3. Halaman gerbang tutup benar** *(peserta, bukan adjustment)* — dicatat supaya tidak "diperbaiki"
ke arah yang salah.

## 1. KEPUTUSAN **bb** — `Close Claim` MENUTUP kasus `[DIPUTUSKAN — niat nyata, seperti aw; veto work owner; OQ-I]`

Peta konektor `Flow/Register_Flow.xml` *(dibaca utuh 27-09-2026; `pyTo` mendahului `pyFrom` di DOM)*:

| Konektor | Dari → Ke | Nama / syarat |
| --- | --- | --- |
| Transition2 b1933 | Start2 → Assignment2 `Input Register` b1501 | `[Always]` |
| Transition3 b1721 | Assignment2 → Assignment1 `Outstanding Claim` b1297 | `InputRegisterClaimLife` |
| Transition4 b1583 | Assignment1 → Decision3 | `Os ClaimLife` |
| Transition11 b1785 | Decision3 → Assignment2 | `IsSendtoAdmin` |
| Transition10 b1653 | Decision3 → Assignment3 `Medical Check` b990 | `Else` |
| Transition1 b1995 | Assignment3 → Decision1 | `MedicalCheck` |
| Transition7 b2140 | Decision1 → Assignment1 | `IsSendtoAdmin` |
| Transition6 b2214 | Decision1 → Assignment4 `Claim Analis` b1116 | `Else` |
| Transition5 b2282 | Assignment4 → Decision2 | `AkseptasiClaimLife` |
| Transition12 b1859 | Decision2 → Assignment3 | `IsSendtoMedical` |
| Transition8 b2066 | Decision2 → Assignment1 | `IsSendtoAdmin` |
| Transition9 b2356 | Decision2 → **End1** b886, **`pyWorkStatus Resolved-Completed`** b899 | `Else` |

Fakta: **tidak ada** konektor bernama `CloseClaim`. `CloseClaim` adalah **local action**
*(`FlowAction/CloseClaim.xml`; `pyLocalAction CloseClaim` di `InputOSClaimLife.xml` b22838 dan
`InputAkseptasiClaimLife.xml` b21457/b21602)* yang membuka `CloseClaim_Section`; activity-nya
berakhir dengan `Call FinishAssignment` b838 yang **seluruh parameternya kosong** *(`TaskStatus`,
`PerformFormName`, `ReviewFormName`, `pyHarness`, `DisplayHarness false` b884–b890)* dan membawa
parameter page yang sedang berjalan b846. Hanya End1 yang menetapkan status kerja; shape lain kosong
*(b859, b942, b1004, b1130, b1248, b1311, b1444, b1499, b1608)*.

Yang ekspor **tidak** jawab: dari Assignment1 *(Outstanding)* tidak ada jalur ke End1 tanpa melewati
Medical Check dan Claim Analis; perilaku mesin untuk `FinishAssignment` dari local action tanpa
konektor senama tidak dapat diturunkan dari ekspor → **OQ-I** untuk pemilik ekspor / pengembang Pega:
*"Di produksi, `Close Claim` ditekan pada Outstanding: kasus selesai, atau pindah ke Medical Check?"*
Catat di `OQ-untuk-tim.md`.

**Keputusan bb** — niat nyata *(label, konfirmasi b499, gerbang "not approved yet")* = menutup:

| Sisi | Isi |
| --- | --- |
| Skema | migrasi **`017_kolom_status_work.sql`** *(+ `_down`)*: `ALTER TABLE {skema}.T_WORK_CLAIM ADD (STATUS_WORK VARCHAR2(32))`. Nilai hanya **`Resolved-Completed`** VERBATIM b899 saat tutup; **NULL** untuk kasus terbuka — status bawaan Pega tidak ada di ekspor dan **tidak dikarang**; kolom `Work Status` inbox tetap seperti sekarang. Blok bertanggal di tiket 14 |
| Backend | `POST /api/klaim-life/{id}/tutup`: gerbang `PenghalangTutupKlaim` yang ada → **409** berisi seluruh penghalang *(pesan rule yang sama)*; lolos → satu transaksi: `STATUS_WORK`, `TAHAP` **dikosongkan** *(assignment selesai; inbox = worklist, kasus hilang dari keempat tab)*, `TGL_UPDATE`, jejak audit *(aturan tiket 09)*. Diizinkan hanya dari tahap **`Outstanding Claim`** dan **`Claim Analis`** *(dua section yang menawarkan local action-nya)* dan pemegang/posisi yang sudah menggerbang rute `tahap`; selain itu 403/409 berkata-kata. Sesudah tutup, **setiap** rute pengubah *(peserta, adjustment, tahap, tolak, tanggal kejadian, akseptasi, putaran, hapus)* ditolak `services` dengan satu kata pesan — uji tabel |
| Frontend | tombol `Close Claim` di Detail **hanya** pada kedua tahap itu; konfirmasi VERBATIM b499; berhasil → kembali ke Inbox *(`closeContainer` b1129)*; gagal → daftar penghalang yang sudah ada. Label "belum menutup" dicabut |
| Dokumen | PARITAS baris 13 dan 27 → **ada**; tiket 08 mendapat bab bertanggal *"Penutupan kasus (Close Claim) — bukti XML"* sebab **nol** tiket menyebut `CloseClaim` |

## 2. URUTAN GILIRAN INI — paket 0 lalu enam kelompok, tanpa pesan di antaranya

| # | Kelompok | Isi | Commit |
| ---: | --- | --- | --- |
| 0 | **Ralat OQ-H + total keenam** *(bc)* | §0.1 | `fix: enam total peserta dihitung seperti SavePesertaClaim b4221–b4742; OQ-H diralat` |
| 1 | **Close Claim menutup** *(bb)* + **DocumentLife tampil** | §1. DocumentLife: `Activity/LoadDocumentLife_ACT.xml` — `Param.inskey = @If(pyWorkCover.pzInsKey="",pyWorkPage.pzInsKey,pyWorkCover.pzInsKey)` b755–b756; `Obj-Browse` b861 *(kelas warisan `DOCUMENT_CLAIM`; padanan Go = `T_CLAIMLF_DOCUMENT`, STRUKTUR-TABEL bab itu; `pohonklaim.go:513` sudah mengenalnya)* → `.DocumentClaimList` b1164; per baris `Call GetUrlGoogleStorage_Act` b1377 *(URL dari penyimpanan luar → di DEV **penanda** "URL menunggu pengirim", bukan panggilan)* → `IMAGEID` b1516, `URLPUBLIC` b1562, `PNOTE` b1583. Tombol `Section/DocumentLife.xml`: `Refresh` b611; `Add attachment` b1245 → localAction `AttachDocumentLife` b1273/b1462; tautan baris dan `View Office Online` b3502 → `DownloadDocumentClaim` b3006/b3519; `Delete` b4288 → localAction `ConfirmDeleteAttachment` b4317/b4555 + `closeContainer` b4441/b4672. Kelompok ini: **daftar** dari `T_CLAIMLF_DOCUMENT` per klaim; ketiga tombol hadir sebagai `BelumTersedia` bernama sampai kelompok 2 | `claim-life: A3 — Detail & Tutup (5), Close Claim menutup; DocumentLife tampil` |
| 2 | **Dokumen** | brief 10 §3. Flow action `AttachDocumentLife` → section `AttachDocScreenLife` *(grid `deleteRow` b3177/b3233)* → activity `SaveAttachLife`; `ConfirmDeleteAttachment` → `DeleteDocument_Act`. Activity dibaca sebagai pohon: `NewAttachLife`, `SaveAttachLife`, `InsertDocument_Act`, `InsertGoogleStorage_Act`, `GetUrlGoogleStorage_Act`, `DownloadDocumentClaim`, `DeleteDocument_Act`, `DeleteGoogleStorage_Act`, `SendEmailWithAttachments` | `claim-life: A3 — Dokumen` |
| 3 | **Medis** | brief 10 §4. Flow action `MedicalCheck` → section `MedicalCheckClaimLife`: `Save` b18572; `Send Back to Admin` b20256 → localAction `SendtoAdmin` b20285/b20409; `Send to Claim Analyst` b21151 → `SendtoAdmin_Act1` b21174/b21359 + `finishAssignment` b21202/b21387 → Transition1 → Decision1 *(`IsSendtoAdmin` → Outstanding; `Else` → Claim Analis)*. `Find Disease` b5061 → `Diagnose_Harness` *(section `Diagnose_Section`)*: `SearchDiagnose_act` b568/b696/b850/b973; `Choose` b2509 → `SetDisease` b2528/b2646; `SetSTS_Reject` atas `.DiagnoseList` b241 | `claim-life: A3 — Medis` |
| 4 | **Akseptasi** | brief 10 §5. Flow action `AkseptasiClaimLife` → section `InputAkseptasiClaimLife`: `Save` b18543; `Send Back to Admin` b20221 → `SendtoAdmin` b20250; `Send Back to Medical` b20467 → `SendtoMedical` b20496/b20652; `Close Claim` b21428 → `CloseClaim` b21457/b21602 *(bb)*; `Save Adjustment` b22641 *(`ClaimLifeDetailGCNM`)* → rute akseptasi yang ada. Transition5 → Decision2: `IsSendtoMedical` → Medical Check; `IsSendtoAdmin` → Outstanding; `Else` → End1 — menyelesaikan Akseptasi tanpa kedua penanda = kasus selesai, **sama** dengan bb | `claim-life: A3 — Akseptasi` |
| 5 | **Komite-handoff** | brief 10 §6. Harness `Committe_Life` → section `ClaimComite` *(`ClaimComiteeLife` b139)*: `Send Claim to Committee` b7033 → runActivity `CreateKMTLife_Act` b7051/b7163; `Cancel` b7715; `GetListKomiteLife`; `SendEmailKlaimLF` → outbox. **Tanpa** menyentuh `komite_*` | `claim-life: A3 — Komite` |
| 6 | **A4** *(bila giliran masih hidup)* | brief 10 §7 | `claim-life: A4 — migrasi data` |

Tiap kelompok: bab **Pembacaan ulang XML** *(path + baris)* di tiket terkait; ralat bertanggal bila
XML membantah tiket; `PARITAS-LAYAR-DAN-AKSI.md`; uji murni + handler + JS; `LAPORAN-GILIRAN-F0.md`
+1 bab. Migrasi hanya dari keputusan tercatat *(`017` = bb)*; nomor `017`–`029`. Nomor baris di atas
adalah hasil `sed -e 's/></>\n</g'`; bila bacaan Anda berbeda satu-dua baris, yang menang adalah
bacaan Anda, dicatat.

## 3. DUA SESI LAIN — BELUM PERNAH MULAI

- Fakta: cabang `modul/premiumlist-life` dan `modul/komite-claim-life` **nol commit**; kedua brief
  tertempel di sesi ini, dan jawaban Anda *("milik sesi worktree …")* benar. Asisten sudah
  fast-forward kedua worktree ke **`9352e45`** *(27-09-2026)* sehingga keduanya memuat perbaikan
  envelope `galat`.
- Sesi ini **tidak** mengerjakan kedua modul itu, kecuali work owner menulis di chat sesi ini
  *"kerjakan modul X dari sesi ini"*. Bila itu terjadi: `Set-Location .worktrees\<modul>\APP_RNM`,
  kerjakan pada cabangnya sesuai brief modulnya, **sesudah** kelompok 6 di atas, laporan terpisah.

## 4. LANGKAH 0 · LAPORAN · TELEMETRI

Langkah 0: brief ini dan induk **sudah di-commit asisten** *(lihat `git log -1`)*; executor hanya
memverifikasi `git status --porcelain` kosong → uji hijau *(300 · 34 SKIP · 201 JS · 47 modul)* →
langsung paket 0. Laporan akhir dan telemetri persis lanjutan 8 §6–§7, ditambah tabel
**kelompok → tombol XML → rute/kontrol**; satu pesan sesudah kelompok 5 *(atau 6)*; bab
**TELEMETRI EKSEKUSI** per paket.

---

*Disusun 27 September 2026 sesudah verifikasi `9352e45` (uji, vet, tsc, build dijalankan ulang; 9 commit
dibaca statistiknya), pembacaan ulang `ProtectCloseClaim_act.xml` utuh, `CloseClaim_Section.xml`,
`FlowAction/CloseClaim.xml`, `Register_Flow.xml` (12 konektor, End1), `ClaimLifeDetailGCNM.xml`
(enam total, sepuluh rujukan), `SavePesertaClaim.xml` b3900–b4760, `SaveOutStandingLife_Act.xml`
b10560–b11200, `SetIndexAdjustmentList.xml`, `SetSTS_Reject.xml`, `LoadDocumentLife_ACT.xml`,
`DocumentLife.xml`, `MedicalCheckClaimLife.xml`, `Diagnose_Section.xml`, `InputAkseptasiClaimLife.xml`,
`ClaimComite.xml`, pencarian `CheckTotalAdjustmentClaim` atas seluruh korpus, dan kode
`tutupklaim.go`, `tutup.go`, `adjustment.go`, `klaimlife.go`, `api.ts`, `KlaimLife.tsx`, `labels.ts`.*
