# 04a: Nomor akseptasi — lahir sekali, di keputusan final

**Status:** ready-for-agent

**Blocked by:** 02 (mesin tangga) — tingkat final harus dapat dikenali lebih dulu

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin nomor akseptasi lahir **sekali saja** — pada keputusan Setuju di
tingkat tertinggi — sehingga satu klaim tidak pernah memperoleh dua nomor, dan tingkat-tingkat di
bawahnya tidak membakar sequence. *(User story 17, 19 di spec)*

## Area codebase

`internal/repository` (pemanggilan kedua SQL penomoran), `internal/services` (gerbang tingkat final;
perakitan nomor per `Type`), `internal/handlers` (nomor tampil pada respons keputusan),
`frontend/` (nomor akseptasi terlihat setelah keputusan final).

## Rule Pega sumber

Rantai **aktif** di dalam step 4 `KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, versi 2026-09-15).
Nomor sub-step dari `<pyStepPageReference>`:

| Sub-step | Rule | Deskripsi |
| --- | --- | --- |
| 4.1 | `RDB-List` | get tanggal produksi |
| **4.7** | `GetKodeProdLife_SQL` (`ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETKODEPRODLIFE_SQL` / `RULE-CONNECT-SQL`) | **AMBIL KODE PROD** — prefix |
| **4.9** | `GetSequenceNumber_SQL` (`ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL`) | **generate MM.YYYY DAN SEQUENCE** |
| 4.11 | `Property-Set` | cabang `QR,QP` |
| 4.12 | `Property-Set` | cabang `TR,TP` |

`[terverifikasi]` Gerbang tingkat final: `pyWorkPage.AcceptStatus = 1 && pyWorkPage.KomiteCount ==
pyWorkPage.KomiteLoop` (baris **5695**).

`[data DBA]` `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER(p_class, p_jenis, p_proddate, OUT p_bulan,
OUT p_seq_number)` — sequence per `(class, jenis, tahun)` di tabel `GENERATE_SEQUENCE_NUMBER`
(PK komposit), memakai `SELECT … FOR UPDATE`; `p_jenis` membedakan retro/non-retro; periode digulir
lewat `POOLDATA.TANGGAL_CLOSING`, dengan aturan cutover `TRUNC(now) <= 02/01/2026` → periode
`12.2025`; keluaran `LPAD(seq, 5, '0')`.

## ADR terkait

**ADR-0006** (penomoran lewat stored procedure — **jangan replikasi logikanya**), **ADR-0011**
(keputusan per baris), **ADR-0015** (batas transaksi dipegang Go; commit segera setelah nomor
terbentuk agar lock `FOR UPDATE` lekas lepas).

## Acceptance criteria

- [ ] Nomor akseptasi dibuat **hanya** pada keputusan **Setuju** di tingkat terakhir
      (`KomiteCount == KomiteLoop`). *(AC 13 spec)*
- [ ] Tingkat bukan-terakhir **tidak** memanggil jalur penomoran sama sekali — sequence tidak
      bergerak. *(AC 5 spec)*
- [ ] Keputusan **Tolak** tidak menghasilkan nomor akseptasi.
- [ ] Nomor diperoleh lewat rantai `GetKodeProdLife_SQL` → `GetSequenceNumber_SQL`; aplikasi
      **tidak** memuat logika pembentukan format nomor. *(AC 14 spec; **ADR-0006**)*
- [ ] Prefix diperoleh lewat **lookup** ke `POOLDATA.KODE_PRODUKSI`, **tidak** ditanam sebagai
      konstanta.
- [ ] Percabangan per `Type` (`QR`/`QP` versus `TR`/`TP`) menentukan skema nomor yang dipakai.
- [ ] Commit terjadi **segera setelah nomor terbentuk**, sehingga lock `SELECT … FOR UPDATE` tidak
      menahan pemutus lain.
- [ ] Dua keputusan final berurutan menghasilkan dua nomor berbeda.
- [ ] Tidak ada padanan `Generate_NoAccept_KMT_Life` / `Generate_NoAccept_KMT_LifeRetro` di kode —
      `[terverifikasi]` keduanya `RequestType` di step ter-remark (`//`)
      `Komite Claim Life/Activity/KomitePostAdjustment.xml`
      (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`),
      baris 1769 dan 1986 (blok remark di 1725 dan 1943). *(AC 28 spec)*

## Catatan — dua rule penomoran lama tidak dimigrasikan

`[terverifikasi]` `Generate_NoAccept_KMT_Life` (step **4.4**) dan `Generate_NoAccept_KMT_LifeRetro`
(step **4.5**) **mati**, dibuktikan **dua sinyal bebas**:

| Sinyal | Hasil |
| --- | --- |
| `<pyStepsBlockName>` | berisi `//` → **REMARK** (baris 1725 dan 1943) |
| Asimetri indeks rujukan | `<RequestType>` 1 × , terindeks `<pyRuleName>` **0 ×** |

Pola yang sama sudah ditetapkan **ADR-0006** untuk Claim Life. **Jangan dimigrasikan.**

`[terverifikasi]` Step **4.6** `Set Nilai Akseptasi` juga ter-remark (baris 2160).

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Implementasi — 28-09-2026 (giliran 10)

### ⚠️ RALAT — ADR-0006 "jangan replikasi logikanya" dilampaui keputusan **o**

Tiket ini (dan ADR-0006) menuntut prosedur `PROC_GENERATE_SEQUENCE_NUMBER` dipanggil apa adanya.
Keputusan **o** (o2/o3, brief GILIRAN-3-KOMITE §0 dan brief modul) membaliknya: prosedur **tidak
dipanggil**, penghitungnya **ditiru** lewat `repository.Penomor` — penghitung yang sama yang sudah
dipakai penomoran PremiumList. AC "tidak memuat logika pembentukan format nomor" karena itu diganti:
bentuk nomor tinggal di **satu** fungsi murni (`models.NomorAkseptasiKomite`), diuji literal.

### Pembacaan ulang XML — langkah 4 sub-langkah 4.7–4.12

| Sub | Isi (`sed -e 's/></>\n</g'`) |
| --- | --- |
| 4.7 b2392 | `GetKodeProdLife_SQL` → `ParamSeq.HASIL3` |
| 4.8 b2576 | `CARI1 = pyWorkPage.pxObjClass`, `CARI2 = HASIL3 + "A"` |
| 4.9 b2737 | `GetSequenceNumber_SQL` → `HASIL1` (`MM.YYYY`), `HASIL2` (urut) |
| 4.10 b2919 | `HASIL1 = substring(0,2) + "." + substring(5,7)` (`MM.YY`) |
| 4.11 b3120 | **QR,QP**: `ACCEPTEDNO = HASIL3+"A"+BusinessCode+"."+HASIL1+"."+HASIL2`; `STS_REJECT = 1`; `ACCEPTATION_DATE = @CurrentDateTime()`; peserta `STS_REJECT = 1`, `IsCheck = "true"` |
| 4.12 b3398 | **TR,TP**: sama, dengan `"AR"` |

Semua sub-langkah bergerbang `ACCEPTEDNO == ""`. ⛔ **Satu penghitung** untuk kedua cabang (`JENIS`
disusun di 4.8, sebelum cabang) — `(ASM-FW-GCNMFW-Work-KomiteLife, RNML-A)`, pasangan yang memang ada
di `GENERATE_SEQUENCE_NUMBER` (SUMBER-PENOMORAN-DBA).

### ⛔ OQ-K-04a `[terbuka — work owner]` — tabrakan lintas jalur MUNGKIN

`GetAcceptedNoCL` yang brief sebut **tidak ada** di jalur Komite: `KomitePostAdjustment` punya nol
`RequestType` itu (grep `RequestType>`: `GETTanggalClosing_SQL`, `Generate_NoAccept_KMT_*` (mati),
`GetKodeProdLife_SQL`, `GetSequenceNumber_SQL`, `UpdateOsAkseptasiClaimLife_sql`). Ia milik jalur Claim
Life (`SaveAdjustment_Act` 1.6.1). Padahal kedua jalur menerbitkan nomor berbentuk **sama**
(`RNML-A…`/`RNML-AR…` + kode bisnis + `.MM.YY.` + 5 digit) dari **penghitung berbeda** (Komite:
`GENERATE_SEQUENCE_NUMBER` per tahun; Claim Life: `ACCEPTATIONNOLIFE_SEQ` global). Tabrakan karena itu
mungkin, dan korpus tidak menjawab kebijakannya.

Yang dibangun: tabrakan **gagal terang** — nomor diperiksa di tabel datar warisan (cara Claim Life) **dan**
di `T_CLAIMLF_ADJUSTMENT.ACCEPTED_NO` (tempat kedua jalur kini menulis); bila dipakai, `409`,
transaksinya — termasuk kenaikan penghitung — batal. ⚠️ Akibatnya: bila tabrakan terjadi, keputusan
akhir itu tertahan sampai work owner memutuskan (lewati nomor, atau pisahkan seri).

### Yang dibangun

- `models/komite_nomor.go` — `NomorAkseptasiKomite` (memakai ulang `RakitNomorPL`/`PeriodeNomorPL`),
  `KodeCabangAkseptasiKomite` (A / AR), `JenisPenghitungKomite`, `ClassPenghitungKomiteLife`.
- `services/komite_akseptasi.go` — `PenyelesaiAkhirKomiteOracle.Akseptasi`: gerbang anggota berjalan
  diulang (penjaga penulis status), `Type`/`BusinessCode` dari klaim induk, penghitung, pemeriksa
  keunikan dua tempat, lalu stempel lewat fungsi Claim Life yang **ada** (`PerbaruiStatusBaris` —
  penjaganya `STS_REJECT = 0` sekaligus gerbang "lahir sekali" —, `CerminkanHeader`), dan jejak
  Outstanding → Aksep di fungsi yang sama. `Tolak` tetap gagal terang (tiket 05).
- Nomor tampil di jawaban keputusan (`nomorAkseptasi`) dan di layar.
- Penjaga Claim Life yang menagih: jejak per penulis transisi (`TestSetiapPenulisTransisiMerekamJejak`)
  dan pendaftaran penulis status dengan gerbangnya (`komite_akseptasi.go` →
  `periksaGiliran(kasus, pelaku.AkunID)`).

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| nomor hanya pada Setuju di tingkat akhir | ✅ `AkseptasiAkhir` (model) → satu-satunya pemanggil |
| tingkat bukan-akhir tidak menyentuh penomoran | ✅ |
| Tolak tanpa nomor | ✅ |
| rantai kode prod → penghitung; tanpa logika format tersebar | ⚠️ diralat (keputusan o) — satu fungsi murni |
| awalan lewat lookup `KODE_PRODUKSI` | ✅ `AwalanProduksi` |
| cabang per `Type` | ✅ A / AR |
| commit segera sesudah nomor | ⚠️ nomor dan stempel satu transaksi dengan keputusan; transaksinya pendek (nol panggilan luar) |
| dua keputusan akhir → dua nomor berbeda | ✅ penghitung `FOR UPDATE` + pemeriksa keunikan |
| nol `Generate_NoAccept_KMT_*` | ✅ (`TestPenghitungKomiteBukanSequenceClaimLife`) |

### Angka

Go **577 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **357** · tsc bersih.
