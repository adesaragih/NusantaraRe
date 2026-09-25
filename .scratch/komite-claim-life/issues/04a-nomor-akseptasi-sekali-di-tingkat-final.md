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
