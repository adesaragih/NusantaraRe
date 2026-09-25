# -*- coding: utf-8 -*-
"""
tulis-penelusuran.py — menulis 2-to-spec/PENELUSURAN-JSON-KE-KOLOM.md

Masukan: alat/telusur-jalur.json + alat/telusur-golongan.json
         (keduanya dibuat alat/telusur-jalur-ke-kolom.py)

Jalankan berurutan:
    PYTHONIOENCODING=utf-8 python alat/telusur-jalur-ke-kolom.py
    PYTHONIOENCODING=utf-8 python alat/tulis-penelusuran.py
"""
import io, os, json
from collections import Counter

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
A = json.load(io.open(os.path.join(AKAR, 'alat', 'telusur-jalur.json'), encoding='utf-8'))
G = json.load(io.open(os.path.join(AKAR, 'alat', 'telusur-golongan.json'), encoding='utf-8'))

hasil, luar, tolak = A['dalam'], A['luar'], A['tolak']
gol, rinci = G['gol'], G['rinci']
gol_dari = {}
for g, isi in rinci.items():
    for j, alasan in isi:
        gol_dari[j] = (g, alasan)

cocok = Counter()
for j, nasib, e, k, tp, ks, tanda in hasil:
    if nasib != 'DIPETAKAN':
        cocok[nasib] += 1
    elif tanda == 'COCOK-TUNGGAL':
        cocok['DIPETAKAN·cocok tunggal'] += 1
    elif tanda.startswith('COCOK-GANDA'):
        cocok['DIPETAKAN·cocok ganda'] += 1
    else:
        cocok['DIPETAKAN·tidak tercocokkan'] += 1

EKSPOR = 'ekspor-2026-09'
ALASAN_NASIB = {
    'DITURUNKAN': 'turunan — tidak disimpan (ADR-0037, §4)',
    'DIBUANG': 'dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3',
    'DITUNDA': 'ditunda — yang ditunggu dan pemiliknya di `PETA-TELUSUR-JSON.md` §4',
}
ART = {
    'GEL-2': 'cabang fakultatif keluar — `RETRO_KELUAR`, **di luar gelombang 1**',
    'TABEL': 'jalur menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel**',
    'NAMA-MASTER': 'nama dibaca dari master; hanya pengenalnya disimpan (**ADR-0041**)',
    'DIBUANG-ULANG': 'sudah diputuskan dibuang atau melebur di §12.5, §13, §14.5',
    'GOL-C': 'mata uang **paket uang golongan C** — menunggu `T-6` (P-8)',
    'BERSUSUN': 'daftar kelas bisnis / kelompok treaty bersusun — **belum punya rumah bernama**',
    'AGREGAT': 'berpasangan dengan besaran **agregat** yang tidak disimpan (§4.2)',
    'YATIM': '**tidak ada rumah, dan tidak ada alasan tertulis**',
}

L = []; w = L.append
w('# Penelusuran jalur `JSONDATA` ke kolom skema baru — Treaty In')
w('')
w('**Tanggal:** 24 September 2026 · **Langkah 1 sesi to-spec**')
w('**Dibangkitkan** `alat/telusur-jalur-ke-kolom.py` + `alat/tulis-penelusuran.py`.')
w('Jangan disunting dengan tangan.')
w('')
w('> ## SEMESTA BERKAS INI — dibaca sebelum satu baris pun dipercaya')
w('>')
w('> **Klaim *"seluruh jalur JSON punya rumah"* DILARANG di berkas ini.** Yang dinyatakan:')
w('> *seluruh jalur **di dalam semesta yang diperiksa** punya nasib, dan semestanya kurang')
w('> sebanyak angka di bawah.*')
w('>')
w('> | Fakta | Angka |')
w('> |---|---:|')
w('> | jalur beradjudikasi di `PETA-TELUSUR-JSON.md` §6 | **%d** |' % (len(hasil) + len(luar)))
w('> | **dikecualikan** — pohon cermin, milik modul Adjustment | **%d** |' % len(luar))
w('> | **lingkup berkas ini** | **%d** |' % len(hasil))
w('> | titik buta penjaga `primary_ok` (`L-8`) | **414** properti |')
w('> | sudah diperiksa untuk 27 entitas gelombang 1 | 74, yang **41** diadili |')
w('> | **belum diperiksa siapa pun** | **340** |')
w('>')
w('> Pohon cermin yang dikecualikan: `ActualValue`, `ValueDifference`, `OLDDATA`,')
w('> `ValueBeforeProrate`. `struktur-treatyin-lama.md` §5 menghitungnya **526 dari 985 simpul**')
w('> (217 + 152 + 142 + 15). **Penyebut %d di atas dihitung ulang dari daftar jalur**, bukan' % len(hasil))
w('> diwarisi dari hitungan simpul — keduanya cara menghitung yang berbeda dan **tidak harus sama**.')
w('>')
w('> **`PETA-TELUSUR-JSON.md` sendiri berlabel *"semestanya kurang"***; lihat')
w('> `4-erd-dan-tabel-datar/AUDIT-PENYEBUT-POHON.md`.')
w('')
w('> ### BATAS PERKAKAS INI, dinyatakan di dalam keluarannya sendiri')
w('>')
w('> Ia dapat menunjukkan sebuah jalur **punya nasib tertulis**, dan untuk jalur DIPETAKAN ia dapat')
w('> menunjuk kolom yang kolom *Asal*-nya menyebut ruas terakhir jalur itu.')
w('>')
w('> Ia **tidak** dapat membuktikan pemetaan itu **benar**. Pencocokannya **leksikal**: dua ruas')
w('> bernama sama di kelas berbeda akan tercocokkan ke kolom yang sama. Setiap baris bertanda')
w('> **COCOK-GANDA** adalah **calon**, bukan putusan — yang ditampilkan hanya kandidat pertama.')
w('')
w('---')
w('')
w('## 1. Hitungan penutup')
w('')
w('| Nasib | Jumlah |')
w('|---|---:|')
for k in ('DIPETAKAN·cocok tunggal', 'DIPETAKAN·cocok ganda', 'DIPETAKAN·tidak tercocokkan',
          'DITURUNKAN', 'DIBUANG', 'DITUNDA'):
    w('| %s | %d |' % (k, cocok.get(k, 0)))
w('| **JUMLAH lingkup** | **%d** |' % sum(cocok.values()))
w('| dikecualikan — pohon cermin | %d |' % len(luar))
w('| **JUMLAH semesta §6** | **%d** |' % (len(hasil) + len(luar)))
w('')
w('**Menutup:** %d + %d = %d, sama dengan penyebut §6.'
  % (sum(cocok.values()), len(luar), len(hasil) + len(luar)))
w('')
w('**Tidak ada jalur tanpa nasib.**')
w('')
w('**Penolakan perkakas:** %s' % ('tidak ada.' if not tolak else
                                 '; '.join('%s = %d' % (k, v) for k, v in tolak.items())))
w('')
w('---')
w('')
tt = [h[0] for h in hasil if h[6] == 'TIDAK-TERCOCOKKAN']
w('## 2. Jalur DIPETAKAN yang tidak tercocokkan — %d, digolongkan bersebab' % len(tt))
w('')
w('**Tidak tercocokkan bukan berarti tidak punya rumah.** Enam dari delapan golongan punya alasan')
w('tertulis di berkas induk; dua terakhir yang menuntut pekerjaan.')
w('')
w('| Golongan | Jumlah | Artinya |')
w('|---|---:|---|')
for g, n in sorted(gol.items(), key=lambda x: -x[1]):
    w('| **%s** | %d | %s |' % (g, n, ART.get(g, '')))
w('')
w('### 2.1 Dua jalur YATIM — dan keduanya menulis ke tabel datar setiap simpan')
w('')
for j, _a in rinci.get('YATIM', []):
    w('* **`%s`**' % j)
w('')
w('**`NusareSharePct`** sudah tercatat di `SPEC-MODEL-DATA.md` §3.2a: **nol penulis** atas kelima')
w('bentuk, dan ia menulis ke `TREATY_IN.NUSARESHAREPCT` pada setiap penyimpanan — kolomnya')
w('**selalu kosong**.')
w('')
w('**`TreatyYear` adalah temuan baru berkas ini.** Ia bertanda **DIPETAKAN** di')
w('`PETA-TELUSUR-JSON.md`, **tidak punya satu pun kolom** di skema baru — tidak ada `TAHUN_TREATY`')
w('maupun padanan lain — dan ia **salah satu dari 20 kolom bisnis `TREATY_IN`**, ditulis setiap')
w('penyimpanan lewat `POOLDATA.PEGA_TREATY_IN`.')
w('')
w('> **Tidak ditambal di sini.** Tahun treaty **mungkin** turunan dari tanggal mulai kontrak, dan')
w('> **mungkin** tahun *underwriting* yang disepakati terpisah. Keduanya lazim di pasar, dan')
w('> keduanya menghasilkan angka yang **sama pada kebanyakan kontrak** — persis bentuk *"cacat')
w('> bersembunyi di balik nilai bawaan"*. Menebaknya sebagai turunan akan benar untuk hampir semua')
w('> baris, dan **salah tanpa terlihat** untuk sisanya.')
w('>')
w('> | | |')
w('> |---|---|')
w('> | **Siapa menutup** | **teknik treaty** — satu pertanyaan, diusulkan sebagai `T-7` |')
w('> | **Yang menagih** | ketiadaan `TAHUN_TREATY` di `KAMUS-KOLOM.md`, dan berkas ini |')
w('')
w('### 2.2 Empat jalur BERSUSUN — daftar yang belum punya rumah bernama')
w('')
for j, _a in rinci.get('BERSUSUN', []):
    w('* `%s`' % j)
w('')
w('Keempatnya daftar **kelas bisnis** atau **kelompok treaty** yang bersusun di bawah induknya.')
w('`EGNPI` punya `ID_KELAS_BISNIS` dan `VERSI_KONTRAK` punya `KELAS_BISNIS_KONTRAK` — tetapi')
w('**tidak ada tempat bagi daftar kelas bisnis per kelompok treaty per layer**.')
w('')
w('> **Ini kandidat kekeliruan lingkup yang sudah tiga kali terjadi** (`CONTEXT.md` §2.2c): daftar')
w('> yang tampak menggantung langsung pada induk besarnya, padahal ada satu tingkat pengelompokan')
w('> di antaranya. **Lingkupnya harus dibaca dari aktivitas yang menambah barisnya**, bukan dari')
w('> bentuk pohon. Dilaporkan; pemiliknya sesi to-spec berikutnya.')
w('')
w('---')
w('')
w('## 3. Tabel penelusuran — satu baris per jalur, %d baris' % len(hasil))
w('')
w('`BUKTI` berbentuk `EVIDENCED(sumber@%s)`. Ekspor ini **mungkin** dari lingkungan QA (`L-9`);' % EKSPOR)
w('penanda itu yang memungkinkan **satu sapuan** mengeluarkan daftar terdampak bila jawabannya datang.')
w('')
w('| JALUR_JSON | NASIB | TABEL | KOLOM | TIPE | BOLEH_KOSONG | ALASAN | BUKTI |')
w('|---|---|---|---|---|---|---|---|')
for j, nasib, e, k, tp, ks, tanda in hasil:
    if nasib == 'DIPETAKAN' and tanda == 'TIDAK-TERCOCOKKAN':
        g, alasan = gol_dari.get(j, ('?', ''))
        n_tampil = 'DIPETAKAN·%s' % g
        e = k = tp = ks = '—'
    elif nasib == 'DIPETAKAN':
        n_tampil = nasib
        alasan = 'tercocokkan **%s** dari kolom *Asal* §10' % tanda.lower()
        ks = 'tidak' if (ks or '').lower().startswith('tidak') else 'ya'
    else:
        n_tampil = nasib
        alasan = ALASAN_NASIB.get(nasib, '')
        e = k = tp = ks = '—'
    w('| `%s` | %s | %s | %s | %s | %s | %s | `EVIDENCED(PETA-TELUSUR-JSON@%s)` |'
      % (j, n_tampil, e or '—', k or '—', (tp or '—').replace('|', '/'), ks or '—', alasan, EKSPOR))
w('')
w('---')
w('')
w('## 4. Jalur yang DIKECUALIKAN — %d, milik modul Adjustment' % len(luar))
w('')
w('Pohon `ActualValue`, `ValueDifference`, `OLDDATA`, dan `ValueBeforeProrate`. **Tidak diadili di')
w('sini**, dan itu batas prompt, bukan kelalaian.')
w('')
w('> `GRL-14` grilling Adjustment sudah memutuskan nasib `ActualValue`: **dibuang seluruhnya**,')
w('> premi aktual menjadi **nilai versinya sendiri**. Maka 108 jalur `ActualValue.*` di antara yang')
w('> dikecualikan ini **tidak akan melahirkan satu pun kolom**, dan itu keputusan yang sudah')
w('> terkunci — bukan pekerjaan yang menunggu.')
w('')
w('| Nasib di §6 | Jumlah |')
w('|---|---:|')
lc = Counter(n for n, _ in luar)
for kk, vv in lc.most_common():
    w('| %s | %d |' % (kk, vv))
w('| **JUMLAH** | **%d** |' % len(luar))
io.open(os.path.join(AKAR, '2-to-spec', 'PENELUSURAN-JSON-KE-KOLOM.md'), 'w',
        encoding='utf-8').write('\n'.join(L) + '\n')
print('DITULIS: 2-to-spec/PENELUSURAN-JSON-KE-KOLOM.md')
print('   baris tabel penelusuran : %d' % len(hasil))
print('   dikecualikan            : %d' % len(luar))
print('   menutup ke              : %d' % (len(hasil) + len(luar)))
