# 09: Jejak audit — siapa + kapan untuk setiap transisi dan setiap jalur balik

**Status:** ready-for-agent

**Blocked by:** 08 (tahap + jalur balik) — seluruh transisi harus ada dulu untuk dapat direkam

## Hasil & nilai pengguna

Sebagai **auditor**, saya dapat mengetahui **siapa** dan **kapan** untuk setiap transisi status dan
setiap pengembalian kasus — sehingga setiap keputusan dapat dipertanggungjawabkan, dan pengembalian
kasus dapat ditelusuri. Sebagai **ReasLifeAdmin**, saya melihat siapa yang mengembalikan kasus
kepada saya dan kapan, sehingga saya tahu apa yang diminta.
*(User story 10, 29, 30 di spec)*

**Ini penyimpangan sadar dari sistem lama — sebuah perbaikan, bukan paritas.**

## Area codebase

`internal/models` (entri jejak audit), `internal/repository` (penyimpanan jejak),
`internal/services` (perekaman pada setiap transisi — satu tempat, bukan tersebar),
`internal/handlers` (identitas pelaku), `frontend/` (tampilan riwayat pada klaim).

Jejak direkam **per baris `AdjustmentList`**, bukan per klaim — karena unit statusnya baris
(**ADR-0011**).

## Rule Pega sumber

| Rule | Identitas | Keadaan sekarang |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` hanya `CREATEOPNAME` + empat kolom tanggal: `ACCEPTATION_DATE`, `CONFIRMATION_DATE`, `CLAIM_RECEIVED_DATE`, `COMPLETE_DATE` |
| `Claim Life/When/IsSendtoAdmin.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN` | `[terverifikasi]` hanya menyimpan nilai `1` — **tanpa pelaku, tanpa waktu** |
| `Claim Life/When/IsSendtoMedical.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN` | idem |
| `Claim Life/RDBList/InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` / `RNM!INSERTLOGSERVICECLAIM` / `RULE-CONNECT-SQL` | `[terverifikasi]` `INSERT INTO pooldata.monitoring_klaim_log` — **log layanan**, bukan jejak keputusan |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml`, `SendEmailKlaimLF.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `…` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` memakai `OperatorID.pyUserIdentifier` / `pyUserName` sebagai data |

## ADR terkait

**ADR-0007** (jejak audit setiap transisi dan setiap jalur balik — penyimpangan sadar),
**ADR-0002** (pelaku diidentifikasi lewat akun berperan, bukan nama ter-hardcode),
**ADR-0011** (jejak per baris).

## Acceptance criteria

- [ ] Setiap transisi status baris menghasilkan catatan berisi **pelaku dan waktu**. *(AC 16 spec)*
- [ ] Setiap pengembalian (`SendtoAdmin`, `SendtoMedical`) menghasilkan catatan berisi **pelaku dan
      waktu**. *(AC 17 spec)*
- [ ] Perubahan nilai `Type` menghasilkan catatan berisi pelaku dan waktu. *(AC 18 spec)* — `Type`
      menyentuh keamanan, bukan sekadar data (**ADR-0012**).
- [ ] Jejak melekat pada **baris** yang bersangkutan, dan riwayat satu klaim dapat dibaca utuh
      lintas seluruh barisnya.
- [ ] Rekam akseptasi lama tetap ditulis sebagaimana adanya — kontrak dengan Komite tidak berubah
      karena tiket ini.
- [ ] Pelaku dicatat sebagai identitas akun; **tidak ada nama orang ter-hardcode**.

## Catatan

`[keputusan work owner 2026-09-14]` Jejak lama **tidak dapat direkonstruksi ke belakang** — data
sebelum cutover hanya punya `CREATEOPNAME` + empat tanggal. Riwayat transisi lengkap hanya ada untuk
kejadian **setelah** cutover. Konsekuensi ini disadari dan diterima.

`[terbuka]` **OQ-013** (pemilik **DBA**) — `COMMIT` berada di dalam blok PL/SQL, sehingga
atomisitas "tulis akseptasi + tulis jejak audit" belum dapat dipastikan. **Tidak memblokir** tiket
ini; memblokir jaminan atomisitasnya.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
