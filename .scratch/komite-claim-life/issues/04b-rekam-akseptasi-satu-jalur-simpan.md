# 04b: Rekam akseptasi — satu jalur simpan berparameter status

**Status:** ready-for-agent

**Blocked by:** **00 (skema penyimpanan komite — PREFACTOR)**, 04a (nomor akseptasi)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin hasil keputusan komite tersimpan sebagai **rekam akseptasi**
lengkap dengan nomornya — lewat **satu** jalur simpan, apa pun keputusannya — sehingga jalur aksep
dan jalur tolak tidak pernah menyimpang diam-diam satu sama lain.
*(User story 17, 20, 21 di spec)*

⚠️ **Penyimpangan sadar dari paritas struktural.**

## Area codebase

`internal/models` (bentuk rekam akseptasi), `internal/repository` (penulisan ke tabel akseptasi),
`internal/services` (satu jalur simpan berparameter status), `frontend/` (dokumen akseptasi dapat
diunduh setelah keputusan final).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | blok PL/SQL `INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE (…)` — **55 kolom terbaca langsung** |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` step **4.15** | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | **Insert ke OS** — jalur aksep, gerbang `AcceptStatus = 1 && KomiteCount == KomiteLoop` (baris 5695) |
| idem, step **5.6** | idem | **Insert ke OS** — jalur reject, gerbang `AcceptStatus==2 && KomiteCount == KomiteLoop` (baris 8119) |
| idem, step **4.17 / 4.18** | `Call PrintAkseptasiPDF` / `Call LoadDocumentLife_ACT` | dokumen akseptasi |

⚠️ `[terverifikasi]` **Rule ini dipakai bersama Claim — Life** — hash ternormalisasi `c50bfd9a12`
identik di kedua modul, dan **tidak terdaftar** di register konflik OQ-011. Mengubah bentuknya adalah
**perubahan kontrak lintas konteks**.

`[data DBA]` `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`: `STS_REJECT NUMBER(38)`; delapan kolom uang
bertipe `NUMBER` **tanpa presisi**; `CURRENCY VARCHAR2(100)`; `TYPE VARCHAR2(10)`.

## ADR terkait

**ADR-0003** (uang non-float — `NUMBER` tanpa presisi menuntut desimal presisi arbitrer),
**ADR-0011** (unit keputusan = baris), **ADR-0001** (tabel akseptasi adalah kontrak bersama),
**ADR-0015** (batas transaksi dipegang Go).

## Acceptance criteria

- [ ] Rekam akseptasi ditulis **hanya** pada tingkat terakhir — baik untuk keputusan Setuju maupun
      Tolak. *(AC 13 spec)*
- [ ] Penyimpanan memakai **satu jalur berparameter status**; **tidak ada dua jalur kembar** di kode.
      *(AC 16 spec)*
- [ ] Tidak ada nilai uang sebagai *binary floating point* di lapisan mana pun maupun di JSON.
      *(AC 17 spec)*
- [ ] Nilai retro yang ditulis adalah nilai **apa adanya dari data policy** — tidak ada logika
      penukaran. *(AC 27 spec; **OQ-065**)*
- [ ] Dokumen akseptasi dapat dihasilkan dan diunduh setelah keputusan Setuju final.
- [ ] Perubahan bentuk rekam akseptasi diperlakukan sebagai **perubahan kontrak lintas konteks**,

### Penyimpanan keputusan final ⚠️ BARU 2026-09-16 — spec §9

- [ ] ⚠️ Keputusan final **juga menetapkan `T_GENERAL_KOMITE.ACCEPT_STATUS`** (`1` aksep / `2` tolak).
      *(AC 30 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Keputusan final **mengisi `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`** sehingga baris adjustment
      menunjuk kasus komite yang memutuskannya. *(AC 32 spec; User story 42)*
- [ ] ⚠️ Ketiganya — rekam akseptasi, `ACCEPT_STATUS`, dan `KOMITE_ID` — ditulis dalam **satu
      transaksi**; kegagalan pada salah satunya **membatalkan seluruhnya**. *(AC 32 spec)*
      dan ditandai demikian di kode.

## Catatan — mengapa dua blok disatukan

`[terverifikasi]` Kedua blok tulis didahului rangkaian precondition **yang sama persis**, dan
mengisi himpunan properti **identik** — diff strukturalnya **nol beda**. Yang membedakan hanya nilai
status yang ditulis dan satu gerbang tambahan.

Itu adalah **salin-tempel jalur aksep menjadi jalur tolak**, sepanjang ribuan baris.
`[keputusan work owner]` **Disatukan** — dua blok kembar adalah tempat divergensi diam tumbuh.

## Catatan — langkah tukar-RetroName tidak direplikasi

`[terverifikasi]` Step **4.14** (jalur aksep, baris 4937–5185) dan step **5.5** (jalur reject,
baris 7663–7911), keduanya berjudul **"Tukar SecurityReinsurer dengan RetroName"**, ber-
`<pyStepsBlockName>//` → **REMARK**.

Ketiga precondition retro berada **di dalam** langkah yang mati itu — termasuk cutover
`ProdDateTime < "20250207T000000.000 GMT"` (baris 5144 dan 7870).

`[keputusan work owner]` Nilai `RetroID`/`RetroName` **sudah di-set di langkah sebelumnya**, mentah
dari data policy. Sistem baru memakainya **apa adanya**; **tidak ada logika penukaran yang
direplikasi**. (**OQ-065** tertutup, **OQ-034** diperkuat.)

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Implementasi — 28-09-2026 (giliran 10)

### Pembacaan ulang XML

- `RDBList/UpdateOsAkseptasiClaimLife_sql.xml` — `BEGIN INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE
  (CASEID, NO_CLAIM, POLICY_NO, …) VALUES (…); COMMIT; END;` (`COMMIT` b117 pecahan baris).
- `KomitePostAdjustment` 4.15 (aksep, gerbang b5695) dan 5.6 (tolak, gerbang b8119) — "Insert ke OS",
  precondition dan properti identik; beda hanya status.
- 4.14 / 5.5 "Tukar SecurityReinsurer dengan RetroName" ber-`//` — tidak ditiru (catatan tiket tetap benar).
- Langkah 8 `InsertJsonClaimLife_Act` — JSON; dibuang (keputusan work owner 2026-09-16, sama dengan
  Claim Life). Langkah 4.17/4.18 `PrintAkseptasiPDF`/`LoadDocumentLife_ACT` dan 14 `SetInformationData`
  — **belum** (lihat di bawah).

### ⚠️ RALAT — `INSERT` warisan menjadi `UPDATE` di sistem ini

Rule warisan meng-`INSERT` baris baru di tingkat akhir. Di sistem ini baris datar adjustment itu
**sudah ada** sejak pendaftaran (Claim Life `PohonKlaim.Simpan`, `ID = adjustment ID`, 18 dari 55 kolom).
`INSERT` kedua menggandakannya (dan bertabrakan PK). Yang berubah di tingkat akhir karena itu
**diperbarui**: `STS_REJECT`, `NO_ACCEPTATION`, `ACCEPTATION_DATE`. Nol baris tersentuh = gagal terang
(hilir Arasapas membaca tabel ini). Kolom lain yang rule warisan isi di titik ini (atribut polis) tetap
milik catatan `[terbuka — tiket 02/03 Claim Life]` di `pohonklaim.go`.

### Yang dibangun

- `repository.RekamAkhirWarisan(adjID, status, nomor, saat)` — **satu** jalur berparameter status
  (AC 16): aksep `1` + nomor, tolak `2` (tiket 05). Status dipagari `{1, 2}` sebelum bind
  (`STS_REJECT NUMBER(38)`, `[data DBA]`) — pola `PeriksaNilaiWarisan`.
- Dipanggil dari `Akseptasi` sesudah stempel relasional, sebelum jejak — satu transaksi dengan
  keputusan tingkat akhir, `ACCEPT_STATUS` (`CatatKeputusan`), dan nomor.
- `KOMITE_ID` di baris adjustment: sudah terisi saat penyerahan (A2, satu transaksi dengan kelahiran
  kasus) — tidak ditulis ulang.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| rekam hanya di tingkat akhir, Setuju maupun Tolak | ✅ aksep · Tolak menyusul tiket 05 lewat fungsi yang sama |
| satu jalur berparameter status | ✅ `TestRekamAkhirSatuJalurDalamAkseptasi` (satu pemanggil) |
| nol float | ✅ |
| retro apa adanya, tanpa tukar | ✅ tidak ada logika penukaran |
| dokumen akseptasi dapat dihasilkan dan diunduh | ⚠️ **belum** — `PrintAkseptasiPDF` merakit HTML → PDF; mesin PDF dan templatenya belum ada di repo. Jalur dokumen Claim Life (`T_CLAIMLF_DOCUMENT` + storage stub) siap menampung berkasnya |
| `ACCEPT_STATUS` final | ✅ (tiket 02) |
| `KOMITE_ID` terisi | ✅ (A2) |
| ketiganya satu transaksi | ✅ |
| bentuk rekam = kontrak lintas konteks | ✅ kolom yang ditulis sama dengan penulis Claim Life; `RekamAkhirWarisan` tidak menambah kolom |

### Angka

Go **579 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **357** · tsc bersih.
