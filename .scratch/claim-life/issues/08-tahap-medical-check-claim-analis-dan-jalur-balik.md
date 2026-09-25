# 08: Tahap Medical Check & Claim Analis + jalur balik

**Status:** ready-for-agent

**Blocked by:** 07 (penegakan peran) — tiap tahap milik peran tertentu

## Hasil & nilai pengguna

Sebagai **ReasLifeMedicalAdvisor**, saya menerima klaim yang menunggu telaah medis dan dapat
mencatat hasilnya; sebagai **ReasLifeSPV**, saya menerima klaim yang siap dianalisis. Keduanya dapat
**mengembalikan** kasus ke admin bila data kurang, dan SPV dapat mengembalikan ke medis bila telaah
perlu diulang — sehingga tidak ada keputusan yang diambil di atas data tidak lengkap.
*(User story 8, 11–14, 23, 24 di spec)*

## Area codebase

`internal/models` (posisi tahap pada klaim), `internal/services` (transisi antar tahap + jalur
balik), `internal/handlers` (endpoint submit dan kembalikan), `frontend/` (layar Medical Check dan
Claim Analis, kontrol kembalikan).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Flow/Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` empat tahap: `Assignment2` "Input Register", `Assignment1` "Outstanding Claim", `Assignment3` "Medical Check", `Assignment4` "Claim Analis" |
| `Claim Life/When/IsSendtoAdmin.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN` | `[terverifikasi]` `pyWorkPage.SendtoAdmin = 1`; menggerbangi pengembalian dari **tiga titik** di `Register_Flow` |
| `Claim Life/When/IsSendtoMedical.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN` | pengembalian SPV → medis; `[keputusan work owner]` artinya — kondisinya tidak terbaca dari tag |
| `Claim Life/Section/MedicalCheckClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `MEDICALCHECKCLAIMLIFE` / `RULE-OBJ-HTML-SECTION` | layar telaah medis; memuat gerbang `pyPosition` |
| FlowAction `MEDICALCHECK`, `AKSEPTASICLAIMLIFE` | — | `[terverifikasi]` tindakan pada kedua tahap |

## ADR terkait

**ADR-0002** (peran per tahap), **ADR-0011** (tahap 2 dan 3 **tidak mengubah** status baris —
baris tetap Outstanding sepanjang Medical Check dan sampai keputusan Komite).

## Acceptance criteria

- [ ] Klaim dapat berpindah Register → Outstanding → Medical Check → Claim Analis, masing-masing
      hanya oleh peran yang berhak.
- [ ] Status baris adjustment **tetap Outstanding** sepanjang perpindahan tahap — perpindahan tahap
      bukan keputusan akseptasi.
- [ ] `ReasLifeMedicalAdvisor` dapat mengembalikan kasus ke `ReasLifeAdmin`.
- [ ] `ReasLifeSPV` dapat mengembalikan kasus ke `ReasLifeAdmin`.
- [ ] `ReasLifeSPV` dapat mengembalikan kasus ke `ReasLifeMedicalAdvisor`.
- [ ] Kasus yang dikembalikan muncul kembali di antrean peran tujuan.
- [ ] Pengembalian **tidak** mengubah status baris adjustment mana pun.

## Catatan

`[terverifikasi]` Pengembalian bukan jalur langka: `SendtoAdmin` menggerbangi pengembalian dari
**tiga titik** di `Register_Flow`. Perekaman pelakunya ditangani tiket **09** — di sistem lama kedua
penanda hanya menyimpan nilai `1`, tanpa pelaku dan tanpa waktu.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
