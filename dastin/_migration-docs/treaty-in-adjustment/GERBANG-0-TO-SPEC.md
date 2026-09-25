> Modul  : Treaty In Adjustment · **gerbang 0 fase to-spec**
> Dibuat : 2026-09-24
> Sifat  : **laporan gerbang.** Bukan spesifikasi. Tidak satu baris spesifikasi ditulis di atasnya.
> Hasil  : **DUA DARI TIGA MEMANTULKAN.** To-spec tidak dibuka.

# Gerbang 0 — hasil

| Gerbang | Isi | Hasil |
|---|---|---|
| **0.1** | tiga jawaban bisnis yang belum masuk grilling | **MEMANTULKAN** — `GRL-12` batal, `GRL-04` pembatalnya menyala |
| **0.2** | `7-2-ACTUALVALUE.md` terbalik; apakah `GRL-14` ikut? | **LOLOS** — dasar `GRL-14` **tidak** berubah, dan sekarang terverifikasi |
| **0.3** | paket `REV-1` … `REV-6` | **TERBUKA** — tidak ada tanggapan tercatat |

---

## 0.2 — LOLOS. `GRL-14` bertahan, dan dasarnya kini terverifikasi mandiri

`tools/pre.py SaveTreatyIn_EDM_Act` mencetak prasyarat langkah 2 apa adanya:

```
2  Page-Copy   "when edmtype= 2 copy data actual to treatyin.actualv"
   WHEN[TreatyIn.EDMState=="3"]   benar -> lompat -> jmp     salah -> lanjut
```

**Benar → lompat.** Maka `EDMState = 3` **melewati** pemotretan, dan yang dipotret adalah 1 dan 2.

| Sumber | Bacaannya | Status |
|---|---|---|
| `PENGETAHUAN.md` §2.2 (Adjustment) | 1 dan 2 dipotret; 3 dilompati | **benar** |
| `GRL-14` (ronde C) | idem — dibangun di atas §2.2 | **benar, dasarnya tidak berubah** |
| `treaty-in/4-erd-dan-tabel-datar/7-2-ACTUALVALUE.md` §1.2 | kebalikannya | **terbalik** — usulan koreksi untuk induk |

**`GRL-14` tidak kembali ke grilling.** Yang perlu dikoreksi adalah artefak induk, bukan putusan
ini — dan koreksinya masuk daftar usulan (GRL-01 butir 5), bukan suntingan langsung.

**Batas yang jujur:** yang saya verifikasi adalah **arah prasyarat langkah 2**. Pemasangan
`WhenTrue`/`WhenFalse` **per langkah untuk seluruh aktivitas** belum disapu, dan prompt Anda sendiri
menyebutnya belum. Klaim ini berlaku untuk langkah 2, bukan untuk seluruh berkas.

Dan satu hal yang ikut terbaca: label langkahnya berbunyi *"when edmtype= **2**"* sementara
kondisinya `!= 3`. **Kebohongan label keenam** di modul ini.

---

## 0.1 — MEMANTULKAN. Dua putusan terkunci tersentuh pembatalnya

### a. `GRL-12` — **BATAL**. Jawaban C adalah pembatalnya, kata demi kata

Pembatal yang ditulis di `GRILL-B/06-PUTUSAN.md`:

> *"Bila dibantah dengan **satu pemakaian yang harus diketahui SEBELUM perubahannya dibuat**,
> keputusan ini ditinjau ulang ke arah **(c)** — sebab hanya niat yang dinyatakan di muka yang dapat
> dipakai sebelum akibatnya ada."*

Jawaban pemilik proses, 24 September:

> *"Dipilih untuk menentukan mana dan tidak boleh diubah"* — materialitas menentukan **field mana
> yang boleh disunting**, dipilih **sebelum** perubahan dikerjakan.

Itu bukan "mirip dengan" pembatalnya. Itu **pembatalnya**. Materialitas dipakai **sebelum** ada
perubahan untuk dipakai sebagai dasar — sehingga ia **tidak dapat** diturunkan dari akibat, karena
pada saat ia dibutuhkan akibatnya belum ada.

**Dan ia dikuatkan ekspor**, tidak hanya keterangan: **kondisi penguncian yang menyebut `EDMMaterialType` (hitungannya `GRILL-D/01-TEMUAN.md` `TD-02`)**. Penguncian itu nyata dan berfungsi.

> **ANGKA DICABUT 24 September 2026 (`MA-13`).** Baris ini semula menyebut *"220 kondisi
> `pyDisabledWhen` di badan 18 seksi, terbanyak `DetailLimits` 60 dan `Layers` 38"*. Sapuan ronde D
> **tidak dapat mereproduksinya** — tiap angkanya kira-kira separuh, konsisten dengan **satu ekspor
> saja** yang disapu. **Arah klaimnya tidak berubah**; angkanya diganti rujukan, sebab angka yang
> tidak dapat direproduksi tidak boleh berdiri sebagai bukti.

> Alasan penolakan (c) di `GRL-12` berbunyi: *"(c) memasang **penolakan baru** pada hari peralihan"*.
> Alasan itu **masih benar** — tetapi ia alasan tentang **ongkos**, bukan tentang **kebenaran**. Kini
> diketahui bahwa (a) tidak dapat menyatakan sesuatu yang dituntut bisnis, dan ongkos tidak
> mengalahkan ketidakmampuan.

### b. `GRL-18` ikut tersentuh — dan ia baru dikunci hari ini

`GRL-18` menghapus `TreatyIn.EDMMaterialType` **tanpa pengganti**, atas dasar `GRL-12`:
*"materialitas turunan, tidak ada atributnya"*. Bila `GRL-12` batal, **atribut itu harus kembali**,
dan bersamanya tiga hal yang `GRL-18` belum pernah jawab:

| Pertanyaan yang lahir | Kenapa ia baru |
|---|---|
| siapa menetapkan materialitas, dan sampai kapan ia dapat diubah | `GRL-18` menjawabnya **hanya untuk jenis**; materialitas dianggap tidak ada |
| jawaban C berbunyi *"**tidak boleh diubah**"* — apakah beku sejak **lahir**, bukan sejak diajukan? | berbeda dari `GRL-18` bagian 2, dan bila ya maka kedua sumbu **berbeda titik bekunya** |
| apa yang dikunci materialitas di sistem **baru** | penguncian layar adalah mekanisme lama; padanannya **kini diputuskan** — `GRL-20` butir a, dengan daftar field per arah di `GRILL-D/01-TEMUAN.md` `TD-02` |

**Bagian 1 dan 3 `GRL-18` (jenis adalah masukan; tiap perubahan berjejak) tidak tersentuh.**
Yang tersentuh bagian yang menyatakan sumbu kedua tidak ada.

### c. `GRL-04` — pembatalnya **MENYALA PADA KEDUA CABANGNYA**

Pembatal yang ditulis di `GRILL-A/06-PUTUSAN.md`:

> *"satu dokumen addendum mengubah lebih dari satu kontrak atau versi, **atau** satu revisi
> menggabungkan lebih dari satu dokumen"*

| Butir | Jawaban | Cabang pembatal |
|---|---|---|
| `DB-3` | **TIDAK** — satu dokumen menyentuh beberapa kontrak | cabang pertama |
| `DB-4` | **TIDAK** — dua addendum dapat digabung jadi satu revisi | cabang kedua |

Keduanya, bukan salah satu. Ditambah pendalamannya: *"1 kontrak di aksep 1 per 1"* — **persetujuan
tetap per kontrak** — dan *"punya nomor sendiri, tapi tetap ada key merujuk ke ID treaty
addendum"*.

**Saya tidak memutuskan penggantinya.** Prompt Anda melarangnya, dan larangan itu benar: keputusan
pengganti menentukan jumlah entitas, dan itu pekerjaan grilling.

Yang saya catat hanya **bentuk pertanyaan** yang harus dijawab, karena ia tidak sama dengan
pertanyaan yang dulu ditolak: yang kembali **bukan** `PENYESUAIAN`. Yang muncul adalah **benda di
atas versi** — dokumen berpenomoran sendiri, satu dokumen ke banyak kontrak, sementara persetujuan
tetap melekat pada versi per kontrak. Apakah itu membatalkan `GRL-04` atau mempersempit lingkupnya
adalah **pertanyaan grilling**, dan saya menyerahkannya apa adanya.

### d. Sapuan yang diminta sebelum `GRL-04` diputuskan — **hasilnya NOL**

**Pertanyaan:** adakah properti yang menyimpan nomor dokumen addendum eksternal?

**Yang disapu:** `datar-treatyin-lama.csv` (985 simpul pohon) dan `datar-properti-per-kelas.csv`
(properti per kelas Pega), atas dua belas akar kata: `Addendum`, `Endors`, `Dokumen`, `Document`,
`Doc`, `Slip`, `Letter`, `Surat`, `NoRef`, `RefNo`, `Reference`, `Nomor`.

**Yang ditemukan, dan tidak satu pun nomor dokumen:**

| Nama | Apa ia sebenarnya |
|---|---|
| `AddendumPremi` | **bendera**, disetel `TreatyInSetEditPre` untuk `EDMState = 3` |
| `AddendumStatus` | status, bukan nomor |
| `ContractRefNo` | **nol penugasan** dalam bentuk apa pun — sudah disapu ulang dan dicatat di `treaty-in/KEPUTUSAN-TANPA-VERIFIKASI.md` §4c; `pyReadOnly`, tidak ada di DDL |
| `DOCUMENT`, `EndorseCost`, `StampEndorse` | lampiran dan biaya, bukan pengenal dokumen |

**Sumber kedua yang mandiri**, dan ia lebih kuat daripada sapuan kata: kelas integrasi addendum
`ASM-FW-GISFW-Int-treaty_in_edm` memuat **tepat 19 properti**, dan daftarnya lengkap —

```
Ceding · CedingID · ChooseStatusAkseptasi · ClassOfBusiness · Commencement · EDMDate ·
EDMMaterialType · EDMState · ID · Information · LeadingReinsSource · LeadingReinsSourceID ·
OLDID · Position · PositionUsername · ProportionType · StatusAkseptasi · Termination ·
TreatyContractName
```

**Tidak ada nomor dokumen di antaranya.**

**Kalibrasi penyapu** (wajib sebelum klaim negatif): dijalankan atas tiga kasus positif yang sudah
diketahui — `EDMState`, `EDMMaterialType`, `OLDID` — dan **ketiganya ditemukan** di kedua berkas.
Penyapunya bekerja; nolnya bermakna.

**Batasnya, dinyatakan:** kedua CSV itu turunan dari sapuan pohon dan indeks rujukan aturan.
**32 langkah Java dan 191 langkah SQL/REST tidak terurai** oleh sapuan mana pun, dan nol di sini
berarti nol **di antara bentuk yang terurai**.

**Maka, menurut aturan putusan yang Anda tetapkan di muka:**

> Nol → **sistem lama tidak pernah merekam nomor dokumen addendum eksternal** → entitasnya
> **BARU**, dan **migrasi tidak punya sumber untuk mengisinya**.

Akibat yang harus ikut diputuskan grilling, dan ia bukan hal kecil: **seluruh baris warisan akan
punya kolom dokumen yang kosong**, dan tidak ada cara mengisinya kecuali orang mengetiknya ulang
dari arsip kertas.

---

## 0.3 — TERBUKA. Tidak ada tanggapan atas paket `REV`

`USULAN-REVISI-ADR.md` disapu untuk kata `tanggapan`, `ditanggapi`, `disetujui pemilik ADR`,
`jawaban pemilik`. **Satu-satunya kemunculan adalah pernyataan prasyaratnya sendiri** — tidak ada
tanggapan tercatat.

`REV-1` … `REV-6` menyentuh ADR-0037, 0048, 0049 (dua kali), 0052, 0055. **Bagian to-spec mana pun
yang bersandar pada kelimanya ditahan**, dan daftarnya akan ditulis saat to-spec dibuka.

---

## Satu hal lagi yang diperiksa dan TIDAK lolos — di luar gerbang 0

Daftar entitas yang **mengikat** adalah `treaty-in/4-erd-dan-tabel-datar/STRUKTUR-DATA.md` (L-6).

| Yang diperiksa | Hasil |
|---|---|
| `PEMULIHAN_LIMIT` ada di `STRUKTUR-DATA.md`? | **TIDAK — nol kemunculan** |
| `STRUKTUR-DATA.md` menyebut cacah entitas? | **tidak menyebut angka sama sekali** |

`SPEC-MODEL-DATA.md` §10.23c menaikkan cacahnya ke **28** dengan menerima `PEMULIHAN_LIMIT`.
**Sumber yang mengikat belum memuatnya.**

Sesuai gerbang selesai butir 3, ini **harus diselaraskan sebelum entitas baru ditambahkan** — dan
ia artefak **induk**, jadi ia menuntut diff yang disetujui lebih dulu.

---

## Yang TIDAK saya kerjakan, dan kenapa

| Langkah | Kenapa dilewati |
|---|---|
| 1 · perbarui `CABANG-K` | isinya bergantung pada `GRL-12` dan `GRL-18` yang baru saja memantul. Memperbaruinya sekarang **membekukan peta di atas putusan yang sedang batal** |
| 2 · `SPEC-MODEL-DATA.md` | `GRL-12` menentukan ada-tidaknya atribut materialitas; `GRL-04` menentukan jumlah entitas. Keduanya terbuka |
| 3 · `SPEC-INVARIAN` + uji | invarian materialitas dan invarian dokumen keduanya belum punya keputusan untuk ditegakkan |
| 4 · `PETA-TELUSUR-JSON` | nasib `EDMMaterialType` berubah arah bila `GRL-12` kembali ke (c) |
| 5 · `4-erd-dan-tabel-datar/` | gerbang 0.1 **memang** menghasilkan calon entitas dokumen; relasinya tidak dapat digambar sebelum entitasnya diputuskan |
| 6 · prosedur addendum | **satu-satunya yang tidak bergantung pada gerbang 0.1** — lihat catatan di bawah |

### Langkah 6 tidak terhalang, dan ia membawa satu temuan tanpa rumah

Pemetaan `PEGA_M_TREATY_IN_EDM` dan `PEGA_M_TREATY_IN_DETAIL_EDM` ke tiga golongan — **PINDAH KE
GOLANG · TURUN JADI CONSTRAINT · TIDAK DIBAWA** — tidak menyentuh materialitas maupun entitas
dokumen. Ia dapat dikerjakan begitu Anda memerintahkannya.

Dan **temuan kedua Anda belum punya nomor TDA**: `TREATYINDETAILEDM` tertinggal tujuh kolom dari
`TREATYINDETAIL` — `DEDUCTIBLE`, `DEDUCTIBLE2`, `ADJ_RATE`, `PREMIUM_EARNED`, `MDP_PCT`, `ROL_PCT`,
`MDP` — dan ketujuhnya justru yang paling mungkin berubah lewat addendum.

> **Ia lahir di sesi induk SESUDAH ronde TDA ditutup, dan karena itu tidak pernah melewati satu pun
> sidang.** Enam belas TDA yang saya nyatakan "enam belas nasib" kemarin **tidak memuatnya**.
> Dicatat di sini sebagai calon **`TDA-17`**, belum diadili, dan **tidak diperbaiki diam-diam**.
