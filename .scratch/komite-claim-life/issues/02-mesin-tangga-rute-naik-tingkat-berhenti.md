# 02: Mesin tangga — rute, naik tingkat, berhenti saat Tolak

**Status:** ready-for-agent

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

- [ ] Kasus baru dirutekan ke baris roster **pertama** yang ber-`KomiteAproval == 0`. *(AC 1 spec)*
- [ ] Keputusan **Setuju** pada tingkat bukan-terakhir menaikkan `KomiteCount` satu dan **tidak**
      menyentuh tabel akseptasi. *(AC 5 spec)*
- [ ] Keputusan **Tolak** menghentikan tangga pada tingkat mana pun ia terjadi. *(AC 6 spec)*
- [ ] Tangga berlanjut **hanya** bila keputusan Setuju **dan** `KomiteCount <= KomiteLoop`.
- [ ] Tiap tingkat menghasilkan satu entri berisi **keputusan, komentar, dan waktu**. *(AC 7 spec)*
- [ ] Nilai keputusan adalah **enum tertutup `{1 = Setuju, 2 = Tolak}`**; nilai lain **ditolak
      terang-terangan**, bukan menghentikan tangga diam-diam. *(AC 35 spec; `[keputusan work owner]`)*
- [ ] Tidak ada padanan `TransferType` di kode — lihat catatan. *(AC 26 spec)*

### Penyimpanan tangga ⚠️ BARU 2026-09-16 — spec §9

- [ ] ⚠️ **`KOMITE_COUNT` dan `KOMITE_LOOP` di-persist di `T_GENERAL_KOMITE`** — bukan hanya hidup di
      halaman kerja. Tingkat berjalan terbaca kembali setelah proses dimulai ulang. *(AC 30 spec;
      penyimpangan sadar 1)*
- [ ] ⚠️ Keputusan di setiap tingkat **menulis baris `T_KOMITE_KOMITELIST`** yang bersesuaian —
      `KOMITE_APROVAL`, `KOMITE_COMMENT`, `DATE_APPROVE`. *(AC 31 spec)*
- [ ] ⚠️ Baris yang ditulis dipilih lewat **`KOMITE_URUT` + `DATA_KOMITE_ID`**, bukan lewat indeks
      posisi. *(AC 34 spec; penyimpangan sadar 2)*

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
