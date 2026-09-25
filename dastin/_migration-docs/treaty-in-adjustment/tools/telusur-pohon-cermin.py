# -*- coding: utf-8 -*-
"""
telusur-pohon-cermin.py — TO-SPEC langkah 1: satu baris per jalur tiga pohon cermin.

LINGKUP (yang Prompt A kecualikan):
    OLDDATA, ActualValue, ValueDifference, ValueBeforeProrate
    + skalar EDM*, RevisionState, ViewState, IsEditData, AddendumPremi

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa sebuah jalur PUNYA nasib tertulis, dan sisinya (IRISAN / 56-KHAS) dihitung
    dari SELISIH HIMPUNAN BERKAS antara kedua ekspor -- bukan dari nama berkas.

APA YANG TIDAK
    Ia TIDAK dapat menyatakan semesta LENGKAP. Peta telusur induk berlabel
    "semestanya kurang" akibat titik buta L-8, dan 340 properti BELUM DIPERIKSA
    SIAPA PUN. Jalur yang ditambahkan di sini TIDAK BOLEH dihitung sebagai
    kelengkapan, dan kalimat itu ikut dicetak ke dalam keluarannya.

    Ia juga TIDAK menentukan nasib. Nasib datang dari putusan grilling yang
    dirujuk per baris; perkakas ini hanya memasangkannya dan memastikan
    TIDAK ADA JALUR TANPA NASIB.

PENOLAKAN dilaporkan satu baris per sebab.
"""
import io, os, re, json
from collections import Counter, defaultdict

AKAR = os.path.dirname(os.path.abspath(__file__))
PETA = os.path.join(AKAR, '..', '..', 'treaty-in', 'PETA-TELUSUR-JSON.md')
EKS_A = r'D:\XML_NURE\Treaty In Adjustment'
EKS_T = r'D:\XML_NURE\Treaty In'

CERMIN = ('ActualValue', 'ValueDifference', 'OLDDATA', 'ValueBeforeProrate')
SKALAR = ['EDMState', 'EDMMaterialType', 'EDMDate', 'EDMEffective', 'EDMType',
          'RevisionState', 'RevisionDate', 'ViewState', 'IsEditData', 'AddendumPremi', 'OLDID']

# ── 1. SISI: selisih himpunan berkas, bukan nama berkas ──────────────────────
def rels(root):
    s = set()
    for dp, dn, fn in os.walk(root):
        for f in fn:
            if f.lower().endswith('.xml'):
                s.add(os.path.relpath(os.path.join(dp, f), root))
    return s

hanya_adj = rels(EKS_A) - rels(EKS_T)          # 56-KHAS
bersama = rels(EKS_A) & rels(EKS_T)            # IRISAN

# indeks: nama daun -> himpunan berkas yang menyebutnya
idx = defaultdict(set)
tolak = Counter()
for akar, tag in ((EKS_A, 'A'), (EKS_T, 'T')):
    for dp, dn, fn in os.walk(akar):
        for f in fn:
            if not f.lower().endswith('.xml'):
                tolak['bukan .xml'] += 1
                continue
            rel = os.path.relpath(os.path.join(dp, f), akar)
            try:
                s = io.open(os.path.join(dp, f), encoding='utf-8', errors='ignore').read()
            except Exception:
                tolak['berkas tidak terbaca'] += 1
                continue
            for t in ('pyIncludedRuleXML', 'pyRuleVersionsList'):
                s = re.sub(r'<%s\b[^>]*>.*?</%s>' % (t, t), '', s, flags=re.S)
            for nm in set(re.findall(r'\b([A-Za-z][A-Za-z0-9_]{2,})\b', s)):
                idx[nm].add(rel)

def sisi(daun):
    berkas = idx.get(daun, set())
    if not berkas:
        return 'TIDAK TERBACA'
    if berkas & hanya_adj:
        return '56-KHAS' if not (berkas & bersama) else 'IRISAN+56'
    return 'IRISAN'

# ── 2. semesta: jalur cermin dari PETA §6 + skalar ───────────────────────────
teks = io.open(PETA, encoding='utf-8').read()
GOL = [('DIPETAKAN', '### DIPETAKAN — '), ('DITURUNKAN', '### DITURUNKAN — '),
       ('DIBUANG', '### DIBUANG — '), ('DITUNDA', '### DITUNDA — ')]
semula = []
for nama, kepala in GOL:
    m = re.search(r'^' + re.escape(kepala) + r'.*?$', teks, re.M)
    if not m:
        tolak['bagian §6 tak ditemukan: ' + nama] += 1
        continue
    mulai = m.end()
    n = re.search(r'^### ', teks[mulai:], re.M)
    seg = teks[mulai:mulai + n.start()] if n else teks[mulai:]
    for b in seg.split('\n'):
        b = b.strip().strip('`')
        if not b or b.startswith(('```', '#', '|', '>', '*', '-')):
            continue
        if not re.match(r'^[A-Za-z][A-Za-z0-9_\[\]\.]*$', b):
            tolak['baris §6 bukan jalur'] += 1
            continue
        semula.append((nama, b))

def cermin(j):
    return j.split('[')[0].split('.')[0] in CERMIN or any(('.' + c) in j for c in CERMIN)

lingkup = [(g, j) for g, j in semula if cermin(j)]
luar = [(g, j) for g, j in semula if not cermin(j)]

# skalar yang ADA di §6 dan sudah terhitung di `luar` -- ditarik masuk lingkup
tarik = []
for g, j in list(luar):
    daun = re.split(r'[\.\[]', j.rstrip('[]'))[-1]
    if daun in SKALAR or j in SKALAR:
        tarik.append((g, j)); luar.remove((g, j))
lingkup += tarik

# ── 3. NASIB per jalur -- datang dari putusan, bukan dari perkakas ───────────
def nasib(j):
    akar = j.split('[')[0].split('.')[0]
    daun = re.split(r'[\.\[]', j.rstrip('[]'))[-1]
    if akar == 'ActualValue':
        return ('DIBUANG', 'GRL-14', '',
                'pohon dibuang seluruhnya; premi aktual menjadi nilai versinya sendiri')
    if akar == 'OLDDATA':
        return ('DIBUANG', 'ADR-0048 butir 2 · GRL-03 · GRL-10', '',
                'tidak ada tabel: sisi lama di-SELECT ke versi dasar lewat ID_VERSI_KONTRAK_DASAR')
    if akar == 'ValueBeforeProrate':
        return ('DIBUANG', 'GRL-15', '',
                'nol penulis di sistem lama; atribut "berlaku sejak" dibawa, mesin pro rata sengaja tidak dibangun')
    if akar == 'ValueDifference':
        if re.match(r'^(Total|Sum)', daun) or 'SummaryList' in j:
            return ('TURUNAN', 'ADR-0037 · §4.2', '',
                    'agregat -- tidak disimpan; dihitung saat dibaca')
        if j.count('.') == 1 and not j.endswith('[]'):
            return ('DISIMPAN', 'ADR-0048 butir 3 · METODE §8.3', 'NILAI_SELISIH',
                    'besaran pokok berselisih; berkunci KUNCI_PADANAN yang MEMUAT mata uang')
        return ('DISIMPAN', 'ADR-0048 butir 3 · METODE §8.3', 'NILAI_SELISIH',
                'ruas pembentuk baris selisih')
    # skalar
    P = {
        'EDMState': ('DISIMPAN', 'GRL-13 · GRL-18', 'VERSI_KONTRAK.JENIS_ADDENDUM',
                     'enum dua nilai untuk versi baru; nilai warisan dibawa apa adanya'),
        'EDMMaterialType': ('DISIMPAN', 'GRL-20 · KTV-1', 'VERSI_KONTRAK.SIFAT_MATERIAL_ADDENDUM',
                            'atribut MASUKAN -- GRL-12 batal; beku sejak AJUKAN, berjejak selama DRAFT'),
        'EDMEffective': ('DISIMPAN', 'GRL-15', 'VERSI_KONTRAK.TANGGAL_BERLAKU_ADDENDUM',
                         'atribut "berlaku sejak" dibawa; mesin pro rata tidak dibangun'),
        'EDMDate': ('DIBUANG', 'ADR-0006', '',
                    'EDMDATE = SYSDATE pada sisip DAN perbarui -- ia tanggal sentuh terakhir, bukan tanggal addendum'),
        'EDMType': ('DIBUANG', 'GRL-13', '', 'tidak pernah ditulis; jenis dibawa EDMState'),
        'RevisionState': ('DIBUANG', 'ADR-0046 · CONTEXT §3.4', '',
                          'dilebur ke satu keadaan siklus hidup'),
        'RevisionDate': ('DIBUANG', 'ADR-0006', '', 'jejak mesin, bukan fakta kontrak'),
        'ViewState': ('DIBUANG', 'ADR-0046 · CONTEXT §3.4', '', 'dilebur ke satu keadaan siklus hidup'),
        'IsEditData': ('DIBUANG', 'ADR-0046 · CONTEXT §3.4', '', 'dilebur ke satu keadaan siklus hidup'),
        'AddendumPremi': ('TURUNAN', 'GRL-13', '',
                          'bendera yang dapat dibaca dari JENIS_ADDENDUM; satu fakta satu penulis (ADR-0041)'),
        'OLDID': ('DISIMPAN', 'GRL-10 · GRL-17', 'VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR',
                  'rujukan versi dasar, disimpan TEPAT DI SATU TEMPAT; arti ketiganya dipisah (7-1-OLDID)'),
    }
    if daun in P:
        return P[daun]
    if j in P:
        return P[j]
    return ('DITUNDA', '', '', 'BELUM DIADILI -- perkakas tidak menentukan nasib')

baris = []
for g6, j in lingkup:
    n, dasar, tabel, alasan = nasib(j)
    daun = re.split(r'[\.\[]', j.rstrip('[]'))[-1]
    baris.append(dict(jalur=j, gol6=g6, nasib=n, dasar=dasar, tabel=tabel,
                      alasan=alasan, sisi=sisi(daun)))

hit = Counter(b['nasib'] for b in baris)
hit_sisi = Counter(b['sisi'] for b in baris)

print('=== telusur-pohon-cermin.py ===')
print('jalur §6 seluruhnya            : %d' % len(semula))
print('  LINGKUP modul ini            : %d  (cermin %d + skalar ditarik %d)'
      % (len(lingkup), len(lingkup) - len(tarik), len(tarik)))
print('  di luar lingkup (pohon utama): %d  -> milik Prompt A' % len(luar))
print()
print('DITOLAK -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-38s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
print()
print('NASIB -- empat, tidak ada yang kelima:')
for k in ('DISIMPAN', 'TURUNAN', 'DIBUANG', 'DITUNDA'):
    print('   %-12s %d' % (k, hit.get(k, 0)))
print('   %-12s %d' % ('JUMLAH', sum(hit.values())))
print('   menutup ke lingkup: %s' % ('YA' if sum(hit.values()) == len(lingkup) else 'TIDAK'))
print()
print('SISI -- dihitung dari SELISIH HIMPUNAN BERKAS:')
for k, v in hit_sisi.most_common():
    print('   %-14s %d' % (k, v))
print()
print('berkas hanya di ekspor Adjustment (56-KHAS): %d' % len(hanya_adj))
print('berkas di kedua ekspor (IRISAN)            : %d' % len(bersama))

json.dump({'baris': baris, 'luar': len(luar), 'tolak': dict(tolak)},
          io.open(os.path.join(AKAR, 'telusur-cermin.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)
