# Triage Labels

Skill berbicara dalam lima peran triage kanonik. Berkas ini memetakan peran itu ke string label
yang benar-benar dipakai di issue tracker proyek ini.

Proyek ini memakai **label default** — string label sama dengan nama perannya.

| Label in mattpocock/skills | Label in our tracker | Meaning                                  |
| -------------------------- | -------------------- | ---------------------------------------- |
| `needs-triage`             | `needs-triage`       | Maintainer needs to evaluate this issue  |
| `needs-info`               | `needs-info`         | Waiting on reporter for more information |
| `ready-for-agent`          | `ready-for-agent`    | Fully specified, ready for an AFK agent  |
| `ready-for-human`          | `ready-for-human`    | Requires human implementation            |
| `wontfix`                  | `wontfix`            | Will not be actioned                     |

Bila sebuah skill menyebut sebuah peran (mis. "apply the AFK-ready triage label"), pakai string
label dari kolom kanan.

Edit kolom kanan bila kelak Anda memakai kosakata lain.

## Cara label dicatat di tracker ini

`[terverifikasi]` Tracker proyek ini adalah **local markdown** (`docs/agents/issue-tracker.md`),
bukan GitHub/GitLab — jadi **tidak ada label API**. Status triage dicatat sebagai baris teks di
dekat bagian atas berkas issue:

```
Status: needs-triage
```

Satu baris `Status:` per berkas. Mengubah triage berarti mengubah baris itu, bukan memanggil CLI.

## Catatan khusus proyek ini

`[pertanyaan terbuka]` Sebagian besar tiket FASE B kemungkinan akan berstatus **`needs-info`** pada
awalnya: dari **57 pertanyaan terbuka** di `discovery/open-questions.md`, **38 memblokir FASE B**
(`discovery/D3-D4-CLOSING-REPORT.md` §3). Tiket yang menunggu jawaban pemilik peran
(DBA / Product+Underwriting / Actuarial / Finance / IAM) **bukan** `ready-for-agent` sampai OQ-nya
terjawab.

Bila sebuah tiket bergantung pada OQ tertentu, sebut nomornya di badan tiket — diambil dari
register, bukan dari ingatan.
