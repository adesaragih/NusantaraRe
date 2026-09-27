# PROMPT — GILIRAN TIGA MODUL 1 *(sesi tunggal di `OUTPUT_HASIL_RNM`, cabang `main`; sesi paralel lihat §5)*: **penyatuan → Claim Life sisa (bd diagnosa, be dokumen, A4 kode) → PremiumList Life tiket 01–09 (+ av) → Komite Claim Life tiket 01–09**

> Baca `PROMPT-INDUK-TIGA-MODUL.md` *(§0.3 baru)*, brief lanjutan 12 *(kelompok yang sudah tuntas tidak
> diulang)*, brief PremiumList Life §1–§3, brief Komite §0–§3. Mekanisme giliran = lanjutan 8 §1.
> Asumsi Anda *"brief modul yang tertempel di sesi ini = kerjakan dari sesi ini"* **diterima** dan kini
> aturan *(induk §0.3)*. Berhenti untuk melapor hanya pada **batas tiket** dengan pohon bersih dan
> `LAPORAN-GILIRAN-F0.md` terbarui; *"giliran sudah panjang"* = catat, lalu lanjut.

---

## 0. KEADAAN SESUDAH `main` `639b85e` · `modul/premiumlist-life` `83658a8` · `modul/komite-claim-life` `7709c0c` — DIVERIFIKASI ULANG

| Klaim laporan *(bagian 1–3)* | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 7 + 1 + 1 commit; tiga pohon bersih | `96a85c1`…`639b85e`; `83658a8`; `7709c0c`; basis keduanya `8907fba`; `git status` kosong ×3 | ✅ |
| main Go 329 · JS 242; PL Go 303 · JS 201; Komite Go 302 · JS 201; vet/gofmt/tsc bersih | dijalankan ulang di **tiga** worktree: main **329 · 0 · 34**, vitest **242**, **50** modul; PL **303 · 0 · 34**, **201**, 47; Komite **302 · 0 · 34**, **201**, 47; vet, gofmt, tsc bersih ×3 | ✅ |
| `SaveAttachLife` b596 kunci `DL-` | `Primary.DOCUMENT = @if(Primary.DOCUMENT == "","DL-"+@pxReplaceAllViaRegex(@CurrentDateTime(),"[^0-9]",""),Primary.DOCUMENT)`; `Local.IDDoc` b615; `Call InsertDocument_Act` b1466 | ✅ ralat OQ-J diterima |
| pengenal dokumen `hh` 12 jam ditiru | `InsertDocument_Act.xml` b648 `yyyyMMddhhmmssSSS`; `InsertGoogleStorage_Act.xml` b1449 `hhmmss` juga | ✅ |
| pencarian diagnosa 500/50, `AND` | `BrowseDiseaseLife_RD.xml`: kelas `Int-DISEASE_LIFE` b40; halaman **50** b514/b747; maks **500** b659; urut `.Number` ASC b598–b601 | ✅ |
| `.Number` → `NUMBER_` "tebakan" | RD sendiri melabeli `.Number` **`ID`** b604/b677; katalog DEV: kolom `ID` | ⛔ tebakannya keliru — **§1 butir 1** |
| `.DiagnoseList` RepeatGrid banyak lawan kolom tunggal | `ClaimLifeDetailGCNM.xml` b3914 kelas `Data-DiagnoseLife`, b3923 `.DiagnoseList`; **`Add` b4690 → `addRow` b4700/b4841**; `Find Disease` b5061 → `Diagnose_Harness` b5081; kolom `.DISEASE` b5422 *(read-only)*, `.ICDCODE` b5616 *(read-only)*, `.GROUPDIAGNOSE` b5860 *(`pxDropdown` b5863)*; **`Delete` b6160 → `deleteRow` b6170** | ✅ dan terjawab: banyak diagnosa adalah **rancangan** *(Add/Delete baris)*, bukan tampilan — **§2 A** |
| `SetDisease` b260/b307 pada `Data-DiagnoseLife`; `SetSTS_Reject` atas `.DiagnoseList` | b61 kelas; b260 `.DISEASE = Param.Disease`; b307 `.ICDCODE = Param.ICD_Code`; b389 `Obj-Save pyWorkPage`. `SetSTS_Reject` kelas `Int-LIFE_PREMIUM_DETAIL` b67; b241 `.DiagnoseList`; b257 `.STS_REJECT = Primary.STS_REJECT` | ✅ |
| Komite: 013 FK tanpa `ON DELETE`; ID 40 lawan 32 | `013_tabel_komite.sql` b39–40 `VARCHAR2(40)`, b71–72; FK b81–82 tanpa `ON DELETE`; STRUKTUR komite menyebut kaskade 5 × | ✅ |
| PL: 219 kolom dibangkitkan dari STRUKTUR | 4 + 56 + 82 + 14 + 16 + 39 + 8 = **219** | ✅ |
| kedua cabang modul mengubah `migrasi_test.go` dan `strukturkolom_test.go` | ya, keduanya → **konflik** saat disatukan; diselesaikan di §5 | ⚠️ |
| nol kebocoran | seluruh berkas yang berubah, tiga cabang | ✅ |

## 1. JAWABAN ATAS "MENUNGGU ANDA" — dari katalog DEV *(27-09-2026, akun POOLDATA, hanya agregat)* dan XML

| # | Hal | Jawaban |
| ---: | --- | --- |
| 1 | **OQ-K.1** nama kolom `DISEASE_LIFE` | `[data DBA — katalog DEV]` `POOLDATA.DISEASE_LIFE` = **`ID` VARCHAR2(100)**, **`ICD_CODE` VARCHAR2(100)**, **`DISEASE` VARCHAR2(1000)**; 97.586 baris; panjang isi maksimum `ICD_CODE` **7**, `DISEASE` **290**. Ditambah label RD b604/b677: `.Number` = `ID`. → `kolomNomorPenyakit = "ID"` *(`repository/penyakit.go:43`)*; OQ-K.1 **ditutup** dengan blok bertanggal; uji yang "menagih OQ-K.1" diganti uji katalog bertag `db` *(SKIP tanpa skema uji)*. Tidak ada `Rule-Obj-Class` di ekspor tetap benar — jawabannya datang dari katalog, bukan dari ekspor |
| 2 | **OQ-K.2** satu atau banyak diagnosa | XML menjawab **banyak** *(Add/Delete baris, §0)*. Keputusan **bd** di §2 A; OQ-K.2 ditutup dengan blok bertanggal |
| 3 | `-migrate` cabang baru; tabrakan nama | katalog DEV: `T_WORK_POLIS`, `T_PREMIUM_LIST`, `T_PREMIUM_LIST_DETAIL`, `T_PREMIUM_LIST_SPREADING`, `T_PREMIUM_LIST_SPREADING_RETRO`, `T_PREMIUM_LIST_SUMMARY`, `T_VIEW_SUGGEST`, `SEQ_WORK_POLIS`, `T_CLAIMLF_DIAGNOSE` **tidak ada** di skema mana pun → **nol tabrakan**. `T_GENERAL_KOMITE` dan `T_KOMITE_KOMITELIST` **0 baris** → `030 MODIFY` aman. `T_MIGRASI` DEV berisi `001`–`016`; `017` belum. Executor **tidak** menjalankan `-migrate`; work owner menjalankannya **sekali dari `main`** sesudah §5 *(perintah induk §0)* |
| 4 | "penyambungan Google Storage / email nyata menuntut izin" | **Bukan izin yang ditunggu.** DEV **tidak pernah** menyambung ke layanan nyata — sudah diputuskan *(aq, brief 10 §3, km4)*. Rutenya dibangun **sekarang** dengan pelaksana **stub** yang sudah ada *(`services.PelaksanaEfek`, `antrean.go:188`; pekerja `PekerjaEfek`)*. Keputusan **be** di §2 B |
| 5 | A4 migrasi data | kodenya dibangun *(uji bertag `db` SKIP)*; **eksekusi** terhadap DEV menunggu work owner dan skema uji DBA *(G1)* |
| 6 | "lanjut di PremiumList, atau kembali ke Claim Life?" | **Keduanya, berurutan, di sesi ini**: §2 → §3 → §4. Tiap modul minimal seluruh tiketnya |

## 2. CLAIM LIFE — sisa, di `main` sesudah §5

### A. Keputusan **bd** — diagnosa banyak per peserta `[DIPUTUSKAN — bukti §0; veto work owner]`

| Sisi | Isi |
| --- | --- |
| Skema | migrasi **`018_t_claimlf_diagnose.sql`** *(+ `_down`)*: `T_CLAIMLF_DIAGNOSE` — `ID NUMBER(19)` PK *(`SEQ_CLAIMLF_DIAGNOSE`)*, `PREMIUM_LIST_DETAIL_ID VARCHAR2(32) NOT NULL` FK → `T_CLAIMLF_PREMIUMLIST_DETAIL(ID)` **`ON DELETE CASCADE`**, `URUTAN NUMBER(5) NOT NULL` *(= `pxListSubscript`; urutan grid dan `SetSTS_Reject`)*, `ICD_CODE VARCHAR2(100)`, `DISEASE VARCHAR2(1000)` *(lebar = sumber nilainya `DISEASE_LIFE`; 255 terbukti kurang untuk 290)*, `GROUP_DIAGNOSE` *(tipe dan lebar dari **sumber dropdown** b5863 — baca `pyOptions`/data page-nya, catat barisnya)*, `STS_REJECT` *(tipe sama dengan kolom peserta, migrasi 003)*; indeks pada FK. Juga di 018: `T_CLAIMLF_PREMIUMLIST_DETAIL.DISEASE` → `VARCHAR2(1000)`. Kaskade: `.DiagnoseList` hidup **di dalam** halaman peserta *(Obj-Save `pyWorkPage` b389)* — hapus peserta = hapus daftarnya; daftar `TestKaskadeHanyaPadaEmpatRelasi` **diperbarui dengan bukti**, bukan dilawan *(pelajaran Komite tiket 00)* |
| Backend | `POST /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa` *(`Add` b4690: baris kosong, `URUTAN` berikut)*; `PUT …/diagnosa/{diagId}` *(`Choose` b2509 → `SetDisease`: `DISEASE`, `ICD_CODE` dari baris hasil pencarian; `GROUP_DIAGNOSE` dari dropdown)*; `DELETE …/diagnosa/{diagId}` *(`Delete` b6160; `URUTAN` dirapatkan)*; daftar ikut `GET /api/klaim-life/{id}`. `SetSTS_Reject` b241/b257: pada rute tolak/akseptasi yang **ada**, `STS_REJECT` setiap diagnosa = `STS_REJECT` peserta — satu transaksi. Gerbang: grid ada di `ClaimLifeDetailGCNM`, yang dimuat oleh `EditDateClaimLife_Section`, `RejectOSClaimLife_Sec`, dan flow action `ViewClaimDetailLifeGCNM` — **bukan** `MedicalCheckClaimLife`; ketiga tombol **tanpa** `pyCondition` sendiri *(b4600–b6180)* → gerbangnya gerbang pemuatnya *(tahap + pemegang)*, dibaca dan dicatat |
| Rekam akseptasi warisan | `RDBList/InsertJsonKlaimLife_sql.xml` dan `UpdateOsAkseptasiClaimLife_sql.xml` memuat `DISEASE`/`ICD_CODE` *(`OS_AKSEPTASI_KLAIM_LIFE.ICD_CODE` VARCHAR2(10), `DISEASE` 1000)* — baca **bind**-nya: dari `.DISEASE`/`.ICD_CODE` peserta *(diisi `SavePesertaClaim` b3340/b3360 dari baris sumber)* atau dari `.DiagnoseList(…)`. Ikuti XML; catat di tiket 04 |
| Frontend | `CariDiagnosa.tsx` → tombol `Choose` aktif *(PUT)*; grid diagnosa di Detail dengan `Add`/`Delete`, kolom `DISEASE`/`ICDCODE` read-only *(b5374/b5566)*, `GROUPDIAGNOSE` dropdown; penanda `KlaimLife.tsx:536` dicabut; label VERBATIM |
| Dokumen | tiket 08 *(bab bertanggal bd)*, tiket 14 *(018)*, `STRUKTUR-TABEL-CLAIM-LIFE.md` bab tabel baru, PARITAS, OQ-K ditutup |

Commit: `claim-life: A3 — Diagnosa (bd), migrasi 018`.

### B. Keputusan **be** — unggah / unduh / hapus dokumen lewat outbox `[DIPUTUSKAN — turunan aq/an/ar1; veto work owner]`

| Sisi | Isi |
| --- | --- |
| Konfigurasi | kunci baru `UNGGAHAN_DIR` di `.env.example` + `config` *(folder lokal aplikasi; **bukan** korpus; wajib ada saat unggah)* |
| Unggah | `POST /api/klaim-life/{id}/peserta/{pesertaId}/dokumen` *(multipart; kategori wajib **ar1**; batas ukuran; MIME dari tabel 48 baris yang ada, pemanggil menang b586)* → berkas ke `UNGGAHAN_DIR` → baris `T_CLAIMLF_DOCUMENT` *(aturan `DL-`, `hh`, "unggah dulu baris kemudian" b1283 yang sudah ada; `T_STORAGE_ID` kosong sampai efeknya selesai)* → satu baris outbox `T_LOG_SERVICE_RNM` jenis `storage-unggah` *(`InsertGoogleStorage_Act`)* yang pelaksana **stub** tandai selesai dan isi `URLPUBLIC` = `GET /api/dokumen/{id}/isi` *(berbatas identitas)* |
| Unduh | `DownloadDocumentClaim` b3006/b3519 → `GetUrlGoogleStorage_Act` → di DEV = URL lokal di atas; `View Office Online` b3502 = tautan yang sama |
| Hapus | `ConfirmDeleteAttachment` b4317 → `DeleteDocument_Act` *(`Obj-Delete` b513 tanpa prasyarat; penyimpanan dilewati bila kosong b472)* → baris dihapus + efek `storage-hapus` ke outbox + berkas lokal dihapus |
| Email | `SendEmailWithAttachments`, `SendEmailKlaimLF` = efek outbox `email`, pelaksana stub mencatat; `EMAILKOMITE` dibaca saat jalan, nol baris disalin |
| Pelaksana nyata | implementasi `PelaksanaEfek` lain, dipilih env *(`PELAKSANA_STORAGE=nyata`)*, **tidak ada** di giliran ini; resolver `M_LINK_SERVICE` tetap dipakai untuk **nama** layanan saja |
| Layar | `DocumentLife`: `Add attachment` b1245 → form `AttachDocScreenLife` *(grid, `deleteRow` b3177/b3233)*; tautan baris + `View Office Online` → unduh; `Delete` b4288 → konfirmasi → hapus; penanda "URL menunggu penyambungan" dicabut |

Commit: `claim-life: A3 — Dokumen (2), unggah/unduh/hapus lewat outbox (be)`.

### C. A4 — kode saja *(brief 10 §7; lanjutan 5 §4)*, tanpa eksekusi. Commit `claim-life: A4 — migrasi data (kode)`.

## 3. PREMIUMLIST LIFE — tiket **01 → 09**, di `main`

Per brief modul §1–§3 *(pl1–pl5)*. Tambahan dari verifikasi: nol tabrakan nama; `T_MIGRASI` belum memuat `050`–`056`
*(work owner)*; migrasi berikut `057`+; menu = kelompok **PremiumList Life** → satu butir **`PremiumList`**
*(harness portal `PremiumLife_harness`, grid `InboxPremiumList`, tombol `Input Offer`/`Input Premium`)*
di Shell yang sama; berkas `polis_*`, rute `handlers/rute_premiumlist.go`, `src/pages/premiumlist/`,
`labels.premiumlist.ts`. Commit `premiumlist-life: tiket NN — <judul>`.

**Sesudah tiket 04** *(pl4 `repository.PolisRingkas`)* selesai — **paket kecil Claim Life av**:
`PanelDataPolis` sepuluh medan dari `T_PREMIUM_LIST` lewat pl4; penanda "menunggu modul PremiumList Life"
dicabut; ambang **ba** *(`PRODUCTINWARD_LIFE` berkunci `ProductNameID`)* disambungkan. Commit
`claim-life: av — PolicyDataLife dari PremiumList Life`.

## 4. KOMITE CLAIM LIFE — tiket **01 → 09** *(01, 02, 03, 04a, 04b, 05, 06, 07, 08, 09)*, di `main`

Per brief modul §0–§3 *(km1–km5)*; `030` sudah menyatu; migrasi berikut `031`+; menu = kelompok
**Komite Claim Life** → **`Inbox Komite`** `[tidak ada di korpus — worklist KomiteRouter]`; berkas
`komite_*`, `handlers/rute_komite.go`, `src/pages/komite/`, `labels.komite.ts`. Commit
`komite-claim-life: tiket NN — <judul>`.

## 5. LANGKAH 0 — PENYATUAN; sesudahnya satu cabang untuk satu sesi

```
Set-Location D:\XML\RNM_BRD\OUTPUT_HASIL_RNM
git merge --no-ff modul/komite-claim-life  -m "merge: komite-claim-life tiket 00 (013 diverifikasi, migrasi 030)"
#   konflik migrasi_test.go / strukturkolom_test.go: ambil KEDUA sisi, uji hijau, commit
git merge --no-ff modul/premiumlist-life   -m "merge: premiumlist-life tiket 00 (migrasi 050-056)"
#   sama; uji hijau di main (Go >= 329, JS 242, tsc, vet, gofmt) — angka baru dicatat
git -C .worktrees\komite-claim-life merge --ff-only main
git -C .worktrees\premiumlist-life  merge --ff-only main
```

Sesudah itu **seluruh** pekerjaan giliran ini di `main` *(folder `OUTPUT_HASIL_RNM`)*; kedua worktree
hanya di-fast-forward pada akhir giliran. Langkah 0 brief modul *(worktree, `.env`, `git log -1`)*
**digantikan** oleh Langkah 0 ini. Sesi paralel, bila work owner membukanya: tiap sesi mengambil
bagiannya *(§2 / §3 / §4)* di worktree-nya dan menyatukannya ke `main` pada akhir tiap tiket dengan
cara yang sama; pembagian berkas induk §2 tetap berlaku di kedua mode.

Work owner, sesudah Langkah 0 di-commit: `-migrate` **sekali dari `main`** *(mengisi `017`, `030`,
`050`–`056`; `018` menyusul sesudah §2 A)* — perintah di induk §0.

## 6. LAPORAN · TELEMETRI

Satu pesan di akhir, atau pada **batas tiket** bila konteks menipis *(pohon bersih; LAPORAN dan tiket
terbarui; SHA disebut)*. Isi: per modul tabel **tiket → commit → menu/tombol XML → rute/kontrol**;
ralat tiket bertanggal; OQ yang ditutup/dibuka; angka uji tiap commit; bab **TELEMETRI EKSEKUSI**
per tiket *(lanjutan 8 §6–§7)*.

---

*Disusun 27 September 2026 sesudah verifikasi tiga cabang (uji, vet, gofmt, tsc, build dijalankan
ulang di tiga worktree), pembacaan ulang `SaveAttachLife.xml`, `InsertDocument_Act.xml`,
`BrowseDiseaseLife_RD.xml`, `SearchDiagnose_act.xml`, `Diagnose_Section.xml`, `SetDisease.xml`,
`SetSTS_Reject.xml`, `ClaimLifeDetailGCNM.xml` b3900–b6250, `SavePesertaClaim.xml`,
`SaveOutStandingLife_Act.xml`, migrasi `013`, `030`, `050`–`056`, dan katalog DEV (`DISEASE_LIFE`,
tabrakan nama, `T_MIGRASI`, cacah tabel komite — agregat saja, nol baris data disalin).*
