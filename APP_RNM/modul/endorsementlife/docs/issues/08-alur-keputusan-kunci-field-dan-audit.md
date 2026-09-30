# 08: Alur keputusan `Confirm`/`Decline`, kunci field permanen, dan jejak audit

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 01 (gerbang — `Decline` harus membebaskan polis untuk di-endorse ulang),
02 (case + kunci field ter-set saat pemetaan)

## Hasil & nilai pengguna

Sebagai **atasan**, saya menyetujui atau menolak endorsement sebelum tersimpan permanen; dan sebagai
**inputor**, bila saya salah memilih polis saya menolaknya lalu memulai ulang — sementara field
penentu **terkunci** agar endorsement tidak diam-diam berubah sasaran.
*(User story 27–30 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Tahap endorsement; status akhir; rekam jejak audit |
| `internal/repository` | Penulisan jejak audit; penutupan case |
| `internal/services` | Mesin alur `Confirm`/`Decline`; **penegakan kunci field di sisi server** |
| `internal/handlers` | Endpoint keputusan; penolakan perubahan field terkunci |
| `frontend/` | Tombol keputusan; field terkunci **beserta pesan penjelas**; riwayat transisi |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InputEDMLife` (Flow) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-FLOW` | `Endorsement Life/Flow/InputEDMLife.xml` (68.374 byte) | alur |
| `IsLifeAccepted` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `ISLIFEACCEPTED` / `RULE-OBJ-DECISIONTABLE` | `Endorsement Life/DecisionTable/IsLifeAccepted.xml` (15.096 byte) | **keluaran `Confirm` / `Decline`** |
| `InputEDMLife` (Section) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/InputEDMLife.xml` (1.255.106 byte) | **kunci field** |
| `CancelCreateCaseEDML` | `DATA-PORTAL` / `CANCELCREATECASEEDML` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/CancelCreateCaseEDML.xml` (20.134 byte) | ⚠️ **bukan** mekanisme batal case |
| `InboxEDMLife` | `ASSIGN-WORKLIST` / `INBOXEDMLIFE` / `RULE-OBJ-REPORT-DEFINITION` | `Endorsement Life/ReportDefinition/InboxEDMLife.xml` | daftar tugas |

`[terverifikasi]` Flow memuat assignment **Input EDM Detail** dan **Input EDM Summary**, keputusan
`IsLifeAccepted`, utility simpan, routing **`WorkList`**, status akhir **`Resolved-Completed`** dan
**`Resolved-Rejected`**.

⚠️ `[terverifikasi]` **`IsLifeAccepted` versi Endorsement hanya mengeluarkan `Confirm` dan
`Decline`** — **tanpa `Reject`**, berbeda dari versi PremiumList Life (`ASM-FW-GISFW-WORK-LIFE` /
`ISLIFEACCEPTED`) yang punya tiga keluaran. **Endorsement tidak punya jalur "kembali ke input".**

### Kunci field — mekanismenya bukan rule `When`

`[terverifikasi]` Modul ini hanya punya **dua** rule `When` (`@BASECLASS` / `ISPEGAPROD`,
`@BASECLASS` / `RECORDEVENT`), keduanya **bukan** soal editabilitas. Mekanismenya **kondisi sebaris
`<pyDisabledWhen>`** di Section:

| Baris | Kondisi | Kendali atas |
| ---: | --- | --- |
| 3166 | `.EditInput==1` | **`.PolicyNo`** |
| 7537 | `.EditInput==1` | **`.EdmTypeBatal`** |
| 15763 | `.EditInput==1` | **`.EdmBatal`** |
| 8965, 10403 | `.EditInput1=1` | kelompok field setelah CSV disimpan |

`[terverifikasi]` **`.EditInput` adalah kunci satu arah** — tiga penulis, **semuanya ke `1`**, tidak
ada yang mengembalikan ke `0`: `MappingEDMLife` (2956), `SetPremi_EDM` (1338),
`SaveCSVEDMLife` (2899, `.EditInput1`).

## ADR terkait

**ADR-0007** (jejak audit setiap transisi — **inti tiket ini**), **ADR-0002** (RBAC — peran pemutus;
penegakannya mengikuti pola Claim Life), **ADR-0001**.

## Acceptance criteria

- [ ] `Confirm` melanjutkan endorsement ke penyimpanan. *(AC 28 spec)*
- [ ] `Decline` **menutup** case; case tertutup **tidak dapat** dilanjutkan maupun diputuskan ulang.
      *(AC 29 spec)*
- [ ] **Tidak ada keluaran `Reject`** di konteks ini — test yang menemukan jalur "kembali ke input"
      **gagal**. *(AC 30 spec)*
- [ ] Setelah `Decline`, polis yang sama **dapat di-endorse ulang** — gerbang "satu EDM terbuka per
      polis" (tiket 01) **tidak lagi menolak**. *(AC 31 spec)*
- [ ] Setiap transisi tahap menulis **jejak audit**: siapa, kapan, dari tahap apa ke tahap apa.
      *(AC 32 spec; **ADR-0007**)*
- [ ] ⚠️ **Penyimpangan sadar — kunci field permanen + pesan.** Nomor polis dan jenis
      endorsement/batal **terkunci permanen** setelah case dibuat; percobaan mengubahnya **ditolak di
      sisi server**, bukan hanya di layar. *(AC 33 spec; `[keputusan work owner + desain]`)*
- [ ] Layar **menampilkan pesan penjelas** mengapa field terkunci — **bukan** sekadar mematikannya
      seperti Pega. *(AC 34 spec)*
- [ ] Riwayat transisi **terlihat pengguna**, bukan hanya tersimpan.

## Catatan — `CancelCreateCaseEDML` bukan mekanisme batal

⚠️ `[keputusan work owner]` **Tidak ada tombol batal terpisah.** Membatalkan endorsement =
**`Decline`**. `[terverifikasi]` `CancelCreateCaseEDML` (`DATA-PORTAL` / `CANCELCREATECASEEDML`)
hanya **satu langkah** `Property-Set` — ia membersihkan halaman dialog, **bukan** menutup case.
**Jangan** memigrasikannya sebagai pembatalan.

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
