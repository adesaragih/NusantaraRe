# Keputusan: sasaran tulis `Save` dan `Submit`

**Pemilik proses, 6 Oktober 2026.** Dicatat apa adanya:

> "submit dan save di sistem sekrang itu akan menyimpan data kedalam masing
> masing table yang ditentukan, jadi tidak mengikuti pega lagi, dan ada satu
> rule lagi apabila user melakukan inputan itu tidak langsung masuk ke table
> nya, harus melakukan save atau submit dulu kemudian datanya masuk ke db"

Keputusan ini **menutup butir 1 daftar kebutuhan** — penghalang yang selama
empat ronde menahan `Save`, `Submit`, `Approve/Decline`, dan seluruh rute
tulis Adjustment.

---

## 1 · Dua aturan, dan apa artinya bagi kode

### Aturan A — sasaran tulis adalah tabel kita sendiri

`Save` dan `Submit` menulis ke **tabel yang ditentukan untuk tiap bagian
layar** — tabel pendaratan `T_TREATY_*` yang migrasi `430`–`439`/`444`
dirikan.

⛔ **Jalur Pega TIDAK diikuti lagi.** Ekspor menunjukkan `SaveTreatyIn_Act`
bermuara ke prosedur `POOLDATA.PEGA_TREATY_IN`, yang menulis ke
`M_TREATY_IN` dan `TREATY_IN`. Keduanya dilarang keras dipakai aplikasi, dan
keputusan ini menyatakan jalur itu memang **tidak dipakai** — bukan ditunda.

⭐ Artinya larangan `M_TREATY_IN` dan larangan `JSONDATA` kini berlaku penuh
di **kedua arah**: tidak dibaca, tidak ditulis.

### Aturan B — isian tidak masuk basis data sampai ditekan

Mengetik di layar **tidak** menyentuh basis data. Data masuk hanya ketika
pemakai menekan `Save` atau `Submit`.

⭐ Hari ini aturan B **sudah terpenuhi dengan sendirinya**: seluruh tab
menyimpan suntingan di keadaan React setempat dan nol jalur tulis ada. Yang
diperlukan bukan membangunnya, melainkan **menjaganya tetap begitu** ketika
jalur tulis lahir — lihat §3.

---

## 2 · ⛔ SATU AKIBAT YANG HARUS DIPUTUSKAN SEBELUM MENULIS SEBARIS KODE

Tabel pendaratan punya **penulis kedua**: pemuat.

`repository.MuatKontrak` menjalankan `KosongkanKontrak` — `DELETE FROM … WHERE
MASTERID = :1` — pada setiap tabel, lalu menyisipkan ulang dari dokumen Pega.
Itu disengaja dan benar untuk pemuatan: ia yang membuat pemuatan **idempoten**.

⚠️ Tetapi sesudah aplikasi ikut menulis ke tabel yang sama, pemuatan
berikutnya akan **MENGHAPUS suntingan pemakai tanpa satu pun galat**. Tidak
ada yang akan tahu sampai seseorang membuka kontraknya dan menemukan
pekerjaannya hilang.

Ini bukan alasan menolak keputusan di atas — keputusan itu sudah diambil.
Ini hal yang harus **diputuskan bersamanya**, dan ada tiga jalan:

| | Jalan | Harganya |
|---|---|---|
| **1** | Pemuat berhenti dipakai sesudah pemindahan selesai | paling sederhana; menuntut pernyataan "pemindahan selesai" yang belum ada |
| **2** | Pemuat melewati kontrak yang pernah disunting aplikasi | menuntut penanda per kontrak (mis. kolom `DISUNTING_APLIKASI`) dan satu migrasi |
| **3** | Suntingan aplikasi mendarat di baris terpisah | menuntut pembeda pada kunci unik tiap tabel — perubahan terbesar |

⛔ Sampai salah satunya dipilih, **jalur tulis tidak boleh diikat
(`Commit`)** ke tabel pendaratan. Menulis lebih dulu dan memutuskan kemudian
berarti data pemakai yang pertama hilang adalah data sungguhan.

---

## 3 · Yang dijaga uji

**Aturan B** ditegakkan `TestInputanTidakMasukDBTanpaSaveAtauSubmit`
(`modul/treatyin/frontend/`): nol komponen tab boleh memanggil jalur tulis
dari `onChange`. Suntingan tinggal di keadaan setempat sampai tombolnya
ditekan.

⚠️ Penjaga itu ada **sebelum** jalur tulisnya lahir, dan itu disengaja:
aturan yang baru dijaga sesudah dilanggar menjaga hal yang sudah rusak.

---

## 4 · Yang BELUM diputuskan, dan tidak boleh ditebak

- **Pembagian tabel per bagian layar.** "masing masing table yang ditentukan"
  — pemetaan tab → tabel sudah ada di `repository.PetaPendaratan`, dan jalur
  tulis dapat memakainya. Tetapi tiga hal belum punya tabel sama sekali:
  kedua tab teks panjang (`Exclusions`, `Special Conditions`) kini mendarat
  di `T_TREATY_REVISION`, sementara `Rate of Exchange` membaca
  `TREATYEXCHANGEYEARLY` — **tabel warisan yang hanya boleh dibaca**.
  ⛔ Jadi `Save` pada grid Rate of Exchange tidak punya sasaran. Itu
  pertanyaan terbuka, bukan hal yang dapat diputuskan di berkas ini.
- ~~**Beda `Save` dan `Submit`.**~~ ✅ **TERJAWAB** — pemilik proses: *"save
  hanya menyimpan di db sedangkan submit adalah alur dalam akseptasi"*.
  Tangganya dibangun di `models.LangkahBerikut`; uraiannya di
  [`ALUR-AKSEPTASI-DARI-EKSPOR.md`](ALUR-AKSEPTASI-DARI-EKSPOR.md).
- ~~**Siapa boleh menekan.**~~ ✅ **TERJAWAB** — `ReasTreatyInAdmin >
  ReasTreatyInSecHead > ReasTreatyInDeptHead > ReasTreatyInDirector`. Tangganya
  sudah jadi logika murni dan diuji; yang belum ada hanya PENGIKATNYA ke
  pemakai yang sedang masuk.
- ⛔ **`PositionUsername` pada jalur NAIK.** Ekspor mengisinya dengan nama
  orang tertanam yang hanya dicapai lewat cabang kode regu yang pemilik proses
  buang. Tanpa cabang itu nol sumber mengisinya. ⚠️ Jalur naik tidak dapat
  diikat sebelum ini dijawab — lihat §2.1 berkas alur.


---

## 5 · ⭐ DIBANGUN 7 Oktober 2026 — Save, Submit, Actions, Decline offer

Perintah pemilik proses: *"nyalakan tombol save dan submit dimana save nya itu masuk kedalam
masing masing table yg dia punya"*. Tiga keputusan tambahan hari itu:

| Pertanyaan | Jawaban pemilik proses | Yang dibangun |
|---|---|---|
| `TREATY_IN` disentuh? | *"saat save boleh menyentuh table treatyin"* | kesembilan belas kolom `PEGA_TREATY_IN` ditulis (`INSERT` kontrak baru, `UPDATE` yang lama) |
| Kurs `TREATYEXCHANGEYEARLY` | tambah & ubah saja | baris ber-ID diperbarui, baris baru ditambah (ID = situs ‖ `LPAD(TREATYEXCHANGE_SEQ, 4)`); ⛔ nol `DELETE` |
| `PositionUsername` jalur naik | *"ambil namanya dari kelola user dan sesuaikan posisinya"* | `LOGIN_ID` akun aktif pemegang workbasket posisi berikutnya (`M_LOGIN_GO_WORKBASKET`), dipisah koma |

**§2 (pemuat menghapus suntingan) TERJAWAB oleh jalan 1:** pemuat membaca `M_TREATY_IN.JSONDATA`,
dan pemilik proses menutup `M_TREATY_IN` dan JSON TOTAL pada 7 Oktober 2026. Pemuat tidak lagi
boleh dijalankan, jadi satu-satunya penulis tabel pendaratan adalah jalur tombol ini.

**Rute** (`handlers/rute_simpan.go`): `POST /api/treaty-in/kontrak/simpan` (Save) dan
`POST /api/treaty-in/kontrak/kirim` (`aksi`: `submit` · `akseptasi` + `pilihan` · `decline`).

**Satu transaksi** (`repository/simpan_kontrak.go`): `MuatKontrak` (kosongkan lalu sisip — peta
yang sama dengan pembaca), kepala `TREATY_IN`, kurs. Dibuktikan di Oracle lewat
`SimpanLaluBatalkanUntukUji` (dibatalkan — nol Commit di uji).

**"Clipboard"** (`services/simpan_kontrak.go`): kepala `TREATY_IN` ⊕ dokumen pendaratan
(`BacaDokumenPendaratan`, kebalikan `MuatKontrak`) ⊕ kiriman layar. Tab yang tidak dibuka tidak
dikirim layar dan TIDAK terhapus. `StatusAkseptasi`, `Position`, `PositionUsername`,
`ChooseStatusAkseptasi`, `CommentList` milik server — kiriman layar untuknya diabaikan.

| Tombol | Aturan ekspor yang dijalankan |
|---|---|
| Save | DT `TreatyInAddNew` (`Position = ReasTreatyInAdmin`); ditolak bila `Resolve Complete` |
| Submit | `TreatyInCheckError` ("Please input Ceding" / "Please input Source of Business (SoB)") → `Akseptasi_DT` (Accept) → `AddCommentList_Act` → simpan |
| Actions | modal `TreatyInAction` (radio Accept/Reject/Decline + Comment) → `Akseptasi_DT` → `AddCommentList_Act` → simpan |
| Decline offer | konfirmasi "Are you sure you want to DECLINE this offer?" → komentar `IsApproved = Decline` → `StatusAkseptasi = Decline` → simpan |

Penekan WAJIB memegang workbasket `Position` berkas (kosong = Admin) — selain itu 403 berpesan.
Tombol Actions tampil bagi pemegang workbasket posisi penyetuju (SecHead/DeptHead/Director).

**Kontrak baru:** pengenal = situs ‖ 6 digit, angka tertinggi `TREATY_IN`/`T_TREATY_REVISION` + 1,
baris situs dikunci `FOR UPDATE` (penjaga pengenal kembar — `TREATY_IN` tanpa kunci unik).
⛔ `M_TREATY_IN_SEQ` milik prosedur Pega TIDAK dipakai: namanya berawalan tabel yang ditutup total.

### 5.1 ⭐ Properti tanpa kolom — DIBUATKAN migrasi `448` (7 Oktober 2026)

Pemilik proses: *"ini dibuatkan saja"*. Migrasi `448_kolom_simpan_lengkap.sql` (modul Treaty In
Adjustment, pemilik migrasi tabel pendaratan):

| Cabang | Properti | Rumah |
|---|---|---|
| Prop · Share | `RNMShareP` `BrokeragePercentP` `OptionLimit` | kolom `T_TREATY_REVISION` |
| Prop · Share | `TotalShareRnmProp` `TotalSpreadedRnmProp` `TotalSpreadedRnmRIProp` | `T_TREATY_TOTAL` (`LarikGabung`, kolom `JENIS`) — tanpa DDL |
| Non-Prop · Installment | `InstallmentNo` · `TotalInstallmentNP` | kolom `T_TREATY_REVISION` · `T_TREATY_TOTAL` |
| Non-Prop · Share | `LimitShareSummaryList` `LimitFacShareSummaryList` | tabel baru `T_TREATY_SHARE_SUMMARY` (+ `SEQ_TT_SHARE_SUMMARY`) |
| Jalur revisi | `RevisionState` `ViewState` | kolom `T_TREATY_REVISION` |

⚠️ **Migrasi `448` dipasang lewat migrator pemakai, bukan dari sesi ini** (`.env` menunjuk
`POOLDATA`). Sampai terpasang, pembaca dan penulis TOLERAN (`repository/kolom_terpasang.go`):
kolom/tabel yang belum ada dilewati, dan Save melaporkannya di `kunciTakTersimpan`. Tombol
Revision MENOLAK berjalan sebelum `REVISIONSTATE`/`VIEWSTATE` ada — keadaan revisi yang tidak
tersimpan akan menjatuhkan tangganya ke cabang biasa.

⭐ **Migrasi `449` (8 Oktober 2026)** — laporan Save kontrak 1002305: *"Belum punya kolom di
tabel, jadi TIDAK tersimpan: Currency, NusaReLimit, RNMShareList, RNMSpreadedList,
RNMSpreadedListRI, SpreadingList"*; pemilik proses: *"berikan apa yang anda butuhkan"*.

| Properti `Limits[].Detail[]` (tab Share Prop) | Rumah | DDL |
|---|---|---|
| `RNMShareList` `RNMSpreadedList` `RNMSpreadedListRI` | `T_TREATY_LIMIT_AMOUNT` (`LarikGabung`, `JENIS`) | TIDAK — cukup peta |
| `Currency` `NusaReLimit` | kolom `T_TREATY_LIMIT_DETAIL` | `449` |
| `SpreadingList` (anak susunan `PROPORTIONALARRG`) | tabel baru `T_TREATY_LIMIT_SPREADING` (+ `SEQ_TT_LIMIT_SPREADING`) | `449` |

⭐ **Migrasi `450` (8 Oktober 2026)** — Save berikutnya kontrak 1002305 melaporkan sembilan larik
yang pembaca pohon Limits sisipkan KOSONG karena tanpa tabel (`larikTanpaTabelDetail/Limit`);
grid Deduction/Reserve/PLA/Cash Loss/Claim Coop dapat diisi, jadi isiannya memang hilang.

| Larik | Rumah | DDL |
|---|---|---|
| `ReserveList` `PLAList` `CashLossList` `ClaimCoopList` `DeductionTotalList` (Detail) | `T_TREATY_LIMIT_AMOUNT` (`JENIS`) | TIDAK |
| `MDPMinList` (layer Non-Prop) | `T_TREATY_LIMIT_MEASURE` (`JENIS`) | TIDAK |
| `DeductionList` (Detail) | `T_TREATY_LIMIT_DEDUCTION` | `450` |
| `CurrencyList` (Detail — grid Parameter Achievement) | `T_TREATY_LIMIT_ACH_PARAM` | `450` |
| `Reinstatement_List` (layer Non-Prop) | `T_TREATY_LIMIT_REINSTATEMENT` | `450` |

Laporan `kunciTakTersimpan` kini dihitung atas properti BERISI saja (`repository.TanpaKosong`):
larik kosong dan teks kosong tidak membawa apa pun yang dapat hilang, jadi tidak dilaporkan.

Pembaca pohon Limits (`bacaEntri`) dan pelapor `KunciBelumTerpasang` kini toleran juga untuk
tabel ANAK: sebelum `449` terpasang kontrak tetap terbuka, dan Save melaporkan `Currency`,
`NusaReLimit`, `SpreadingList` alih-alih menelannya.

Pembacaan kembali: keenam kolom `T_TREATY_REVISION` dibaca `BacaRevisiPendaratan` (disaring kolom
terpasang) dan dikirim di `KontrakWarisan.penampung`; form menyemaikannya ke penampung halaman
sesudah `kosongkan()`, sehingga isian yang di-Save tampil kembali.

### 5.2 ⭐ Jalur REVISI — DIBANGUN 7 Oktober 2026

`Section/InputTreatyInOffer.xml` cell 994 → `SetTreatyIn_Act(viewstate=1, revisionstate=1)`.

| Bagian | Aturan |
|---|---|
| Tombol `Revision` (daftar) | tampil bila workbasket `ReasTreatyInAdmin` && `Position = ''` && `Resolve Complete`; event **`doubleclick`** persis ekspor |
| `POST /api/treaty-in/kontrak/revisi` | [6]–[7] `ViewState=1`, `RevisionState=1`, `StatusAkseptasi=""`, `PositionUsername=operator` (Position TIDAK diubah); [9] komentar "Create Revision" tanpa `IsApproved`; [10–11] disimpan SEKETIKA |
| Form mode lihat ber-`RevisionState=1` | Comment hidup (cell 9), Submit revisi tampil (cell 21); Submit biasa dan Decline offer tersembunyi (`ViewState=1`) |
| Tangga | cabang revisi `Akseptasi_DT` (`LangkahBerikut` revisi): Admin → SecHead → Resolve Complete, `RevisionState`/`ViewState` dikosongkan |

`RevisionState`/`ViewState` kini milik server — kiriman layar untuknya diabaikan.

### 5.3 ⭐ Save/Submit ADJUSTMENT — DIBANGUN 7 Oktober 2026

Penulisnya di modul ini (`repository/simpan_penyesuaian.go`, `services/simpan_penyesuaian.go`) —
tabel pendaratannya sama, dan modul tidak boleh saling impor. Layar Adjustment memanggil:

| Rute | Tombol | Aturan ekspor |
|---|---|---|
| `POST /api/treaty-in/penyesuaian/simpan` | Save | pra-DT `TreatyInAddNew(Status=1)` → `SaveTreatyIn_EDM_Act` |
| `POST /api/treaty-in/penyesuaian/kirim` | Submit · Actions | `TreatyInSubmitEDM` (nol validasi — [2]–[4] mati; [7] Value Difference EDM 1/2) · `TreatyInAkseptasiEDM_Act` |
| `POST /api/treaty-in/penyesuaian/hapus` | Decline offer | `TreatyInDeclineConfirmation_postactEDM` [6]–[7]: baris `TREATY_IN_EDM` + pendaratan kedua sisi DIHAPUS fisik |

Sasaran tulis, satu transaksi: `TREATY_IN_EDM` (INSERT/UPDATE, `EDMDATE = SYSDATE` berbentuk
`DD-MON-RR`) dan `T_TREATY_*` sisi New (`MASTERID = ID`) — draf juga sisi Old (`ID#LAMA`).

⚠️ Beda yang DISENGAJA dari ekspor:
- draf berpengenal yang sudah ada DITOLAK (Pega menimpanya — rumus `/R0n` dapat bertabrakan);
- Decline offer menolak penyesuaian `Resolve Complete` (Pega hanya bersyarat `ViewState != 1`;
  penghapusan fisik tidak dapat dibatalkan).

Yang tidak dibangun: salinan `ActualValue` EDM 1/2 (`SaveTreatyIn_EDM_Act` [2]–[4]),
`TreatyRevisionCopyAttachment` (unggah ulang ke Google Storage lewat Java), dan penulisan kurs
(prosedur EDM tidak menulis `TREATYEXCHANGEYEARLY`; grid yang diubah dilaporkan tak tersimpan).

### 5.5 ⭐ Upload lampiran — DIBANGUN 8 Oktober 2026

Pemakai: *"untuk upload file seharusnya kesini SELECT * FROM M_ATTACHMENTTREATY_2 … pelajari xml ini
WorkAttachments.xml"*. Rantai ekspor: tombol `Upload file` per kategori (`ViewState !='1' ||
RevisionState='1'`) → `SetkategoriDoc` → modal **ASM Attach Content** (`TreatyAttachContent`) →
**Attach** (`TreatySaveAttachment`) → Refresh (`GetMasterTreatyCategory_Act`).

| Langkah ekspor | Di sini |
|---|---|
| `InsertGoogleStorage_Act` → `ServiceGoogle` (Google/upload, `M_LINK_SERVICE`; token `GCP_IMAGE`; App `T_FOLDER_IMAGE`) | `services/lampiran_unggah.go` — Folder `Contract/Doc/YYYY/MM/`, Namafile `yyyyMMdd-hhmmss-S - <nama>`, Durasi 1800 |
| `Insert_T_Storage_SQL` | `T_STORAGE_IMAGE` (STORAGE `standard`) |
| `InsertAttachment2_Sql` | `M_ATTACHMENTTREATY_2` — kolom JSON-nya TIDAK disebut (Pega mengisinya kosong) |

Rute: `GET`/`POST /api/treaty-in/kontrak/{id}/lampiran` (multipart `kategori` + `berkas`).
Kedua INSERT satu transaksi; berkas berjenis tak dikenal (`GetMimeType`) dilaporkan, berkas lain
tetap diunggah. `M_ATTACHMENTTREATY_2` dikeluarkan dari daftar baca-saja `TestWarisanHanyaDibaca`.
⚠️ Token baru menuntut `STORAGE_TOKEN_SALT` di `.env`. ⭐ 8 Oktober 2026 (izin pemakai): procedure
`GET_TOKEN_STORAGE` ternyata TIDAK memakai garam (`MD5('ASMAPP' || waktu)`), jadi tidak ada yang
disalin — `.env` diisi nilai ACAK 64 heksa (tidak pernah dicetak); tokennya tetap disimpan di
`GCP_IMAGE` seperti buatan procedure.

**Modal View File** (`ShowAttachmentTreaty`, 8 Oktober 2026):

| Tombol | Rantai ekspor | Rute |
|---|---|---|
| nama berkas | `DownloadAttachmentTreaty` → `GetUrlGoogleStorage_Act` (URL tersimpan selama EXPDATE berlaku; selain itu Google/geturl + `Update_T_Storage_SQL`); isi berkas DIALIRKAN backend (https saja, tanpa pengalihan) dan diunduh lewat `fetch` beridentitas | `GET …/lampiran/{lid}/isi` |
| View Office Online (xls/xlsx/doc/docx/ppt/pptx) | idem, URL penampil kantor dimuat di bingkai POPUP (pola Master Product Name Life) | `GET …/lampiran/{lid}/tautan?office=1` |
| Delete (`ViewState !='1' \|\| RevisionState='1'`) | `Delete_act`: `DeleteGoogleStorage_Act` (Google/delete, `DeleteStorage_SQL`) → `DeleteAttachment2_Sql` | `POST …/lampiran/{lid}/hapus` |
| Change Category → Save (status bukan Resolve Complete/Decline) | `ChangeDokument_Act` → `ChangeKateAttachment2_Sql` per baris | `POST …/lampiran/kategori` |

⛔ Ralat 8 Oktober 2026: bentuk pertama membuka URL bertanda tangan di tab baru (`window.open`) —
dilanggar penjaga lintas modul `modul/claimlife/frontend/unduhdokumen.test.ts` (nol navigasi lewat
skrip). Kini unduhan lewat rute `/isi` dan penampil kantor di bingkai popup.

Penyimpangan yang dinyatakan: hapus di storage gagal → baris TETAP; jawaban geturl tanpa
`appfolder` tidak mengosongkan APPFOLDER; `CATEGORY` dari katalog yang dirapikan (nama Non-Prop
untuk `00007`).

**Bentuk modal Attach = Master Product Name Life** (pemakai, 8 Oktober 2026: *"contoh yang ada di
menu product name life itu aja ditiru samakan yah krn tim saya menggunakan hal tersebut"*): kotak
seret-lepas, pilihan digabung (nama kembar dilewati), **Remove** per berkas, unggah SATU berkas per
permintaan secara berurutan dengan progres `Uploading 2/5 nama`, berkas gagal (termasuk yang ditolak
`GetMimeType`) tinggal di pilihan beserta alasannya, modal tertutup bila semua berhasil. Helper
`gabungBerkas` / `unggahBerurutan` DISALIN ke `frontend/unggahBerkas.ts` — impor lintas modul
dilarang `lapisan.guard.test.ts`. Rute dan backend tidak berubah.

### 5.4 Yang sengaja TIDAK dibangun

- `SaveTreatyInDetail_Act` / `SaveTreatyInDetailEdm_Act` (`M_TREATY_IN_DETAIL`,
  `TREATYINDETAILEDM` saat Resolve Complete) — pemilik proses 7 Oktober 2026: *"biarkan data nya
  ditarik dari table nya masing masing saja … untuk melihat status nya bisa dari treaty_in"*.

### 5.5 Tombol `Copy` daftar kontrak (8 Oktober 2026)

`Section/InputTreatyInOffer.xml` cell 993 (label `Copy`, event `click`): `refresh thisSection` +
`SetTreatyIn_Act(ID=.ID)` (tanpa `viewstate`/`revisionstate` — langkah simpan [10–11] tidak jalan),
lalu `refresh thisSection` + `TreatyInCopy`. Syarat tampilnya IDENTIK dengan `Revision` (cell 994):
WB `ReasTreatyInAdmin` && `Position = ''` && `StatusAkseptasi = 'Resolve Complete'`.

`Activity/TreatyInCopy.xml` — empat langkah, **nol `SaveTreatyIn`**:

| Langkah | Isi |
|---|---|
| [1] | `CommentList = ""`, `RevisionState = ""` |
| [2] | `RDB-List GetCurrentDate` → halaman `time` |
| [3] | `CommentList(<APPEND>)`: `Date`, `OperatorName = OperatorID.pxInsName`, `Suggest = "Copied from ID " + TreatyIn.ID` (tanpa `IsApproved`); `OLDID = TreatyIn.ID` |
| [4] | `ID = "UnknownId"`, `Position = "ReasTreatyInAdmin"`, `StatusAkseptasi = ""`, `ViewState = 0`, `IsEditData = 0` |

Jadi Copy **membuka draf**, tidak menyimpan; `UnknownId` adalah penanda kontrak baru yang sama
dengan `Add` (`TreatyInInputVis` Add=1), dan pengenal baru lahir ketika draf di-Save.

| Tombol | Rute | Tulisan |
|---|---|---|
| Copy (daftar) | `GET /api/treaty-in/kontrak-warisan/{id}/salin` | NOL — draf tampil (status/posisi/riwayat sudah `TreatyInCopy`) |
| Save draf | `POST /api/treaty-in/kontrak/salin` | clipboard SUMBER (kepala + `T_TREATY_*`) → `TreatyInCopy` → isian layar → `tulis` dengan ID kosong (`idKontrakBaru`) |
| Submit / Decline draf | `POST /api/treaty-in/kontrak/salin` + `aksi` | salinan disimpan lalu `KirimKontrak` atas ID barunya; `TreatyInCheckError` diperiksa SEBELUM menulis |

Mengapa bukan `/kontrak/simpan` biasa: kontrak baru di sana hanya memegang kiriman layar (tab yang
pernah dibuka), dan `CommentList` milik server. Sasaran tulisnya tetap yang sama dengan Save.

Penyimpangan/yang tidak disalin (dinyatakan): lampiran (`M_ATTACHMENTTREATY_2` berkunci
`TreatyIn.ID` — `UnknownId` memberi nol), log achievement, polis produksi (panel dikosongkan bila
`EDMState` kosong, persis `FetchTreatyExistingProduction`); kurs `TREATYEXCHANGEYEARLY` berkunci
tahun. Submit/Decline draf = DUA transaksi (Pega: satu `SaveTreatyIn`) — bila langkah kedua gagal,
salinan tetap tersimpan sebagai draf Admin dan pesannya menyebut ID-nya. Kode:
`backend/services/salin_kontrak.go`, `backend/handlers/rute_salin_kontrak.go`.
