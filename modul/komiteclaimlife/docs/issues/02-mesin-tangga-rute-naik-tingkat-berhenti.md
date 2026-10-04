# 02: Mesin tangga — rute, naik tingkat, berhenti saat Tolak

**Status:** selesai — 28-09-2026, `ae77877` (tingkat akhir tersambung di `fe50a13`/`0ab7622`)

**Blocked by:** **00 (skema penyimpanan komite — PREFACTOR)**, 01 (terima kasus + inbox per posisi)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin tangga persetujuan **berjalan sendiri**: kasus berpindah ke
anggota berikutnya saat disetujui, dan berhenti saat ditolak — tanpa siapa pun perlu mengatur
urutannya manual. *(User story 3, 13, 14, 15 di spec)*

Ini **inti konteks Komite**: tangga yang di Pega tidak terlihat di graf, dinyatakan eksplisit.

## Area codebase

`internal/models` (`KomiteCount`, `KomiteLoop`, entri `KomiteList` per tingkat),
`internal/services` (mesin tangga: pemilihan tingkat, transisi, penghentian),
`internal/handlers` (endpoint simpan keputusan), `frontend/` (layar keputusan).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`, 26.387 byte | `[terverifikasi]` `param.AssignTo = .KomiteID` (baris ~294), bergerbang `.KomiteAproval == 0` (baris ~382) |
| `Komite Claim Life/When/IsKomiteLoop.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `ISKOMITELOOP` / `RULE-OBJ-WHEN` | `[terverifikasi]` `.AcceptStatus = "1"` **DAN** `.KomiteCount <= .KomiteLoop` |
| `Komite Claim Life/Flow/KomiteLife_Flow.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITELIFE_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` 1 Assignment + 1 Decision + 4 connector — **satu assignment yang di-loop** |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, 515.675 byte | `[terverifikasi]` jejak per tingkat (baris ~792, ~908); `KomiteCount + 1` (baris ~9020) |
| `Komite Claim Life/Section/ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION` | `[terverifikasi]` dropdown `AcceptStatus` wajib, baris 32607 — `pyFormat = pxDropdown`, `pyRequired = true` |

## ADR terkait

**ADR-0011** (unit keputusan = baris `AdjustmentList`), **ADR-0007** (jejak tiap transisi),
**ADR-0001** (tangga adalah isi konteks ini; batasnya ke Claim Life di tiket 05).

## Acceptance criteria

- [x] Kasus baru dirutekan ke baris roster **pertama** yang ber-`KomiteAproval == 0`. *(AC 1 spec)* — bukti: `repository/komite_inbox.go:sqlSaringInboxKomite` (`MIN(KOMITE_URUT)` ber-approval `0`), uji `TestInboxKomiteHanyaAnggotaBerjalan`; arti kode transisi `6` `[dugaan kuat]`
- [x] Keputusan **Setuju** pada tingkat bukan-terakhir menaikkan `KomiteCount` satu dan **tidak**
      menyentuh tabel akseptasi. *(AC 5 spec)* — bukti: `models/komite_tangga.go:TerapkanKeputusanKomite`, uji `TestTanggaTigaTingkatSetujuSeluruhnya`
- [x] Keputusan **Tolak** menghentikan tangga pada tingkat mana pun ia terjadi. *(AC 6 spec)* — bukti: uji `TestTolakMenghentikanDiTingkatManaPun`, uji `TestTanggaBerhentiTidakMenerimaKeputusan`, uji `TestInboxKomiteMengikutiIsKomiteLoop`
- [x] Tangga berlanjut **hanya** bila keputusan Setuju **dan** `KomiteCount <= KomiteLoop`. — bukti: `models/komite_tangga.go:TanggaBerlanjut`, uji `TestIsKomiteLoopVERBATIM`
- [x] Tiap tingkat menghasilkan satu entri berisi **keputusan, komentar, dan waktu**. *(AC 7 spec)* — bukti: `repository/komite_keputusan.go:InboxKomite.CatatKeputusan` (approval, komentar, `DATE_APPROVE`) + jejak per tingkat di `KeputusanKomite.Putuskan`
- [x] Nilai keputusan adalah **enum tertutup `{1 = Setuju, 2 = Tolak}`**; nilai lain **ditolak
      terang-terangan**, bukan menghentikan tangga diam-diam. *(AC 35 spec; `[keputusan work owner]`)* — bukti: uji `TestKeputusanDiLuarEnumDitolakTerang`, uji `TestKeputusanAsingDitolakSebelumBasisData`
- [x] Tidak ada padanan `TransferType` di kode — lihat catatan. *(AC 26 spec)* — bukti: nol kemunculan `TransferType` di `internal/` dan `frontend/src/` (grep 28-09-2026)

### Penyimpanan tangga ⚠️ BARU 2026-09-16 — spec §9

- [x] ⚠️ **`KOMITE_COUNT` dan `KOMITE_LOOP` di-persist di `T_GENERAL_KOMITE`** — bukan hanya hidup di
      halaman kerja. Tingkat berjalan terbaca kembali setelah proses dimulai ulang. *(AC 30 spec;
      penyimpangan sadar 1)* — bukti: `repository/komite_keputusan.go:sqlMajukanTangga` + `repository/komite_inbox.go:InboxKomite.Kasus` (dibaca ulang dari `T_GENERAL_KOMITE`)
- [x] ⚠️ Keputusan di setiap tingkat **menulis baris `T_KOMITE_KOMITELIST`** yang bersesuaian —
      `KOMITE_APROVAL`, `KOMITE_COMMENT`, `DATE_APPROVE`. *(AC 31 spec)* — bukti: `repository/komite_keputusan.go:sqlCatatAnakTangga`, uji `TestKeputusanKomiteBersyaratDuaBaris`
- [x] ⚠️ Baris yang ditulis dipilih lewat **`KOMITE_URUT` + `DATA_KOMITE_ID`**, bukan lewat indeks
      posisi. *(AC 34 spec; penyimpangan sadar 2)* — bukti: uji `TestKeputusanKomiteBersyaratDuaBaris` (`DATA_KOMITE_ID` + `KOMITE_URUT`); ⚠️ mengandaikan `KOMITE_URUT` (= `DEGREE` roster) berturutan 1..n

## Catatan — `TransferType` tidak direplikasi

`[terverifikasi]` `TransferType` **tidak pernah diisi** di korpus: **nol `Property-Set`** terhadapnya
di `Claim Life` maupun `Komite Claim Life`. Ia hanya **dibaca** di dua tempat — precondition kedua
`KomiteRouter` (baris ~442) dan visible-when `.TransferType==2` di `ShowTransfer.xml`
(baris ~29177). Karena tak pernah di-set, `TransferType == '2'` **selalu FALSE** → cabang mati.

`[keputusan work owner]` **Dibuang**, beserta bagian UI `ShowTransfer` yang bergantung padanya.
Routing tingkat digerakkan **hanya** oleh `.KomiteAproval == 0`. (**OQ-033** tertutup.)

### ✅ Catatan ini DIVERIFIKASI ULANG dan TETAP BENAR — 28 September 2026

Brief giliran 3 §1 menyuruh meralat tiket 02 *"bila mengulanginya"*. Ia **tidak**
mengulanginya: kalimat di atas berbunyi *"beserta **bagian UI** `ShowTransfer` yang
bergantung padanya"* — **bagian**, bukan seluruh section — dan ia menyebut **kedua**
pembacanya dengan benar.

Yang keliru ada di `spec.md` butir 34, yang berbunyi *"layar `ShowTransfer` ... ikut
dibuang"*. Butir itu **ditarik** 28-09-2026 dengan rantai keterjangkauannya:

```
Flow/KomiteLife_Flow.xml       b763  pyMOName ViewTransferDtl
  -> FlowAction/ViewTransferDtl.xml  b91   pySectionReference ShowTransfer
       -> Section/ShowTransfer.xml         LAYAR KEPUTUSAN KOMITE
```

`TransferType` muncul di section 1.125.234-byte itu **tepat sekali** (b29177).

⚠️ Dan satu penajaman untuk tiket ini sendiri: pembaca **pertama** —
`KomiteRouter` b442 — menggerbangi **PEROUTEAN**, yaitu siapa yang menerima
pekerjaan. Itu lebih berakibat daripada wadah tersembunyi, dan `[terbuka — pemilik
ekspor]` apakah produksi mengisi `TransferType` dari tempat yang tidak diekspor.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Implementasi — 28-09-2026 (giliran 10)

### Pembacaan ulang XML — `KomitePostAdjustment` sebagai pohon (14 langkah tingkat atas)

Perintah: `sed -e 's/></>\n</g' Activity/KomitePostAdjustment.xml` (11.730 baris), lalu
`<pyStepPageReference>RH_1.pySteps(n)` + metode + deskripsi.

| # | Metode · deskripsi | Yang ditiru di tiket ini |
| ---: | --- | --- |
| 1 | `Obj-Open-By-Handle` "Open Claim" → `TempOpenPage` | klaim induk dibaca lewat `COVER_KEY` |
| 2 | "Set Index": `Local.IndexAdjustment`, `Local.IndexPremium`, **`Local.Komite = pyWorkPage.KomiteCount`** | anak tangga yang ditulis = `KOMITE_URUT = KOMITE_COUNT` |
| 3 | "set acc / reject adjustment in PNC and komite": `KomiteList(Local.Komite).KomiteAproval = .AcceptStatus`, `.KomiteComment = .Comment`, `.DateApprove = @CurrentDateTime()` | ✅ `CatatKeputusan` |
| 4 | "Approve Last Komite" — `AcceptStatus = 1 && KomiteCount == KomiteLoop` (b5695) | ⛔ **digerbang** → tiket 04a/04b |
| 5 | "Reject" — `AcceptStatus==2 && KomiteCount == KomiteLoop` (b8119) | ⛔ **digerbang** → tiket 05 |
| 6–12 | `UpdateWorkObject`, `Obj-Save`, `InsertJsonClaimLife_Act`, "EXIT JIKA RETROID", Arasapas, email, kasir | tiket 04b/06/07 |
| 13 | `KomiteCount = KomiteCount + 1` (b9028–9029), **tanpa syarat** | ✅ `CatatKeputusan` (satu tulisan dengan langkah 3) |
| 14 | `SetInformationData` | tiket 04b |

`When/IsKomiteLoop.xml`: `.AcceptStatus = "1"` **dan** `.KomiteCount <= .KomiteLoop` →
`models.TanggaBerlanjut`. `Flow/KomiteLife_Flow.xml`: `End1` **tanpa `pyWorkStatus`** — "tangga
berhenti" karena itu diturunkan dari `IsKomiteLoop` atas `T_GENERAL_KOMITE.ACCEPT_STATUS` +
`KOMITE_COUNT`, bukan dari status kerja karangan (`models.KasusDiTangga`).

⚠️ **Urutan menentukan**: langkah 4/5 membaca `KomiteCount` **sebelum** langkah 13 menaikkannya;
`IsKomiteLoop` membacanya **sesudah**. Tingkat akhir = `count == loop`; berlanjut bila `count+1 ≤ loop`
(`TestTanggaTigaTingkatSetujuSeluruhnya`).

⚠️ `AcceptStatus` (dropdown `ShowTransfer` b32607, `pyRequired true`) ber-`pyListSource associated`:
pilihannya milik aturan properti yang **tidak diekspor**. Enum `{1 Setuju, 2 Tolak}` = `[keputusan work
owner]` (AC 35); nilai lain **ditolak terang** sebelum basis data. `.KomiteComment` b31001
`pyRequired false`.

### Yang dibangun

- `models/komite_tangga.go` — enum, `TanggaBerlanjut` (= `IsKomiteLoop`), `KasusDiTangga`,
  `TerapkanKeputusanKomite` (tingkat diputus, count baru, berlanjut, akseptasi-akhir, tolak-akhir).
- `repository/komite_keputusan.go` — dua tulisan **bersyarat** dalam satu transaksi: anak tangga
  (`DATA_KOMITE_ID` + `KOMITE_URUT` + `KOMITE_OPERATORID` + approval masih `0`) dan kepala
  (`KOMITE_COUNT` = count yang dibaca). Dua keputusan serentak: satu kalah, `409`.
- Inbox tiket 01 kini juga menegakkan `IsKomiteLoop`: **Tolak di tingkat tengah** meninggalkan anggota
  berikutnya ber-approval `0`; tanpa syarat ini kasusnya jatuh ke inbox mereka.
- `services/komite_keputusan.go` — hanya **anggota berjalan** (403 selain itu), gerbang klaim induk
  tertutup (`PastikanKasusTerbuka`, butir bb), jejak per tingkat dalam transaksi yang sama.
  ⛔ **Tingkat akhir digerbang**: `PenyelesaiAkhirBelumAda` menolak langkah 4/5 (`501`) dan
  transaksinya batal utuh, sampai tiket 04a/04b/05 menggantinya — tidak ada kasus "disetujui" tanpa
  nomor akseptasi, atau "ditolak" tanpa baris klaimnya tahu.
- `POST /api/komite/{id}/keputusan`; terdaftar di penjaga butir bb (`rutePengubah`,
  `layananPengubah`). Penjaga batas konteks Claim Life: pengecualian jalur penuh kedua untuk
  `repository/komite_keputusan.go` — ia persis "milik konteks Komite" yang penjaga itu maksud.
- Layar: formulir keputusan di `KasusKomite.tsx`, **hanya pada giliran pelaku** — dropdown wajib,
  `Comment`, `Submit`/`Cancel` VERBATIM, kalimat konfirmasi VERBATIM.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| rute ke anggota pertama ber-approval 0 | ✅ (tiket 01, `[dugaan kuat]` kode transisi 6) |
| Setuju bukan-akhir naik satu, tak menyentuh akseptasi | ✅ |
| Tolak menghentikan di tingkat mana pun | ✅ — di tengah tanpa langkah 5 (tiket 05) |
| berlanjut hanya Setuju ∧ count ≤ loop | ✅ `IsKomiteLoop` VERBATIM |
| satu entri per tingkat: keputusan, komentar, waktu | ✅ anak tangga + jejak |
| enum tertutup; nilai lain ditolak terang | ✅ |
| nol `TransferType` | ✅ nol di kode |
| `KOMITE_COUNT`/`LOOP` di-persist; baris dipilih lewat `KOMITE_URUT` + `DATA_KOMITE_ID` | ✅ |
| syarat tampil blok rincian `ShowTransfer` (`IsTreatyIn`, `Type TP/TR`, `SwiftCode`, `RetrocadedShare`) | ⚠️ belum — layar kini hanya kepala + tangga + keputusan |

### Angka

Go **567 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **356** · tsc bersih.
