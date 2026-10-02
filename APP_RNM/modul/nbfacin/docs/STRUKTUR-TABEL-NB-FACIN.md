# Struktur Tabel — NB FacIn: PETA TABEL WARISAN yang dibaca

Modul ini **tidak membuat satu tabel pun** (tiket 20). Ia hanya **membaca** enam tabel limit akseptasi yang sudah ada
di `POOLDATA`, dengan nama tabel dan kolom **verbatim**. Berkas ini **peta**, bukan DDL: hanya kolom yang dibaca
repository (`backend/repository/limit.go`). Ke-enamnya dinyatakan di `MODUL.md` bab "Tabel warisan: dibaca, tidak
dibuat".

Sumber tipe `[terverifikasi]`: DDL `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` (diberikan work owner, 01-10-2026; butir 42).
Kolom lain di DDL (`ID`, `MAX_LIMIT_*`, `BATAS_WAKTU`, `TGL_UPDATE`, `EFFECTIVE_DATE`, `WORKBASKET`, `JABATAN_ATASAN`,
`LIMITTRADE_BOTTOM`) tidak dibaca. ⛔ `NAMA` dan `LOGIN` **tidak pernah** dibaca (CLAUDE.md §4 butir 10, K-025).

## M_LIMIT_PROPERTYY

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga (`LetterNo`) |
| `TEAM_GROUP` | VARCHAR2(20) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_PROPERTY_NON_PREFERREDD

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(20) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_ENGINEERINGG

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_NONPROPANDENGG

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga |
| `TEAM_GROUP` | VARCHAR2(10) | team group kasus |
| `LIMIT_BOTTOM` | NUMBER(*,0) | batas bawah limit |
| `LIMIT_BOTTOM2` | NUMBER(*,0) | batas bawah limit jalur banding |

## M_LIMIT_FINANCIALINS

| Kolom | Tipe DDL | Dibaca untuk |
| --- | --- | --- |
| `JABATAN` | VARCHAR2(100) | jabatan tangga bentuk B (berspasi, mis. `DIREKTUR TEKNIK`) |
| `LIMITBOND_BOTTOM` | NUMBER(*,0) | limit Bond |
| `LIMITCREDITCL_BOTTOM` | NUMBER(*,0) | limit Kredit CL (juga Trade Credit, A26) |
| `LIMITCREDITNCL_BOTTOM` | NUMBER(*,0) | limit Kredit NCL |
