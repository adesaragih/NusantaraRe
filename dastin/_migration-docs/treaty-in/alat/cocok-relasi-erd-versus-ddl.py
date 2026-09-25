# -*- coding: utf-8 -*-
r"""
cocok-relasi-erd-versus-ddl.py -- mencocokkan RELASI dua arah:

    4-erd-dan-tabel-datar/ERD.md sec 2   (naratif, MENGIKAT soal relasi)
    2-to-spec/ddl-usulan/*.sql           (FOREIGN KEY, yang akan dibangun)

KENAPA ADA
    GERBANG 0 mencocokkan NAMA ENTITAS, dan batas itu dinyatakan di dalam
    keluarannya sendiri: "ia mencocokkan NAMA, bukan ISI". Sebuah relasi yang
    hilang TIDAK terlihat olehnya -- kedua sumber tetap menyebut kedua entitasnya.
    Perkakas ini menutup celah itu untuk relasi.

APA YANG DAPAT DIBUKTIKAN
    Bahwa tiap relasi di ERD.md sec 2 punya FOREIGN KEY-nya, dan sebaliknya; dan
    bahwa perilaku hapus yang ERD.md putuskan benar-benar sampai ke DDL.

APA YANG TIDAK
    Ia TIDAK memeriksa kolom bukan-kunci, dan TIDAK memeriksa kardinalitas yang
    ERD.md tulis (lambang 1/o/<) terhadap UNIQUE di DDL -- itu pemeriksaan ketiga
    yang belum ada, dan ketiadaannya dicetak di bawah.

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/cocok-relasi-erd-versus-ddl.py
"""
import io, os, re, json
from collections import Counter, OrderedDict

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
ERD = os.path.join(AKAR, '4-erd-dan-tabel-datar', 'ERD.md')
DDL = os.path.join(AKAR, '2-to-spec', 'ddl-usulan')
KAMUS = os.path.join(AKAR, '2-to-spec', 'KAMUS-KOLOM.md')

tolak = Counter()

# -- 1. relasi dari ERD.md sec 2 --------------------------------------------
teks = io.open(ERD, encoding='utf-8').read()
a = teks.index('## 2. Relasi')
b = teks.index('## 3. Atribut bukan-kunci')
seg = teks[a:b]

# bentuk:  INDUK  <lambang>  ANAK  [hapus: X]   <catatan opsional>
#          ANAK   >o--1      INDUK [hapus: X]   (arah terbalik, sec 2.7 dan 2.8)
POLA = re.compile(
    r'^([A-Z][A-Z0-9_]+)\s+([1o<>─—o-]+)\s+([A-Z][A-Z0-9_]+)\s*'
    r'(\[luar\])?\s*(\[G[23]\])?\s*\[hapus:\s*([^\]]+)\]\s*(.*)$', re.M)

erd_rel = []
n_baris_blok = 0
for blok in re.findall(r'```\n(.*?)```', seg, re.S):
    for l in blok.split('\n'):
        if not l.strip():
            continue
        n_baris_blok += 1
        m = POLA.match(l.strip())
        if not m:
            tolak['baris di dalam blok sec 2 yang bukan baris relasi'] += 1
            continue
        kiri, lam, kanan, luar, gel, hapus, ekor = m.groups()
        # arah: ">" di awal lambang berarti KIRI adalah ANAK (sec 2.7, 2.8)
        if lam.startswith('>'):
            induk, anak = kanan, kiri
        else:
            induk, anak = kiri, kanan
        erd_rel.append(dict(induk=induk, anak=anak, lambang=lam,
                            hapus=hapus.strip(), luar=bool(luar), gel=bool(gel),
                            catatan=ekor.strip()))

# -- 2. FK dari ddl-usulan/ --------------------------------------------------
ddl_fk = []
tabel = set()
for f in sorted(os.listdir(DDL)):
    if not f.lower().endswith('.sql'):
        continue
    s = io.open(os.path.join(DDL, f), encoding='utf-8').read()
    m = re.search(r'CREATE TABLE\s+\S+\.(\w+)\s*\(', s)
    if not m:
        tolak['berkas .sql tanpa CREATE TABLE (kerangka)'] += 1
        continue
    nm = m.group(1)
    tabel.add(nm)
    for fm in re.finditer(r'ADD CONSTRAINT\s+\w+\s+FOREIGN KEY\s*\((\w+)\)\s*\n?\s*'
                          r'REFERENCES\s+\S+\.(\w+)\s*\((\w+)\)\s*([^;]*);', s):
        od = re.search(r'ON DELETE\s+([A-Z ]+)', fm.group(4) or '')
        ddl_fk.append(dict(anak=nm, induk=fm.group(2), kolom=fm.group(1),
                           on_delete=(od.group(1).strip() if od else None)))

# -- 3. kolom dari KAMUS -----------------------------------------------------
kam = {}
ent = None
for l in io.open(KAMUS, encoding='utf-8').read().split('\n'):
    m = re.match(r'^## `([A-Z][A-Z0-9_]*)`\s*$', l)
    if m:
        ent = m.group(1); kam[ent] = set(); continue
    if ent and l.startswith('| `'):
        k = re.match(r'^\| `([A-Z][A-Z0-9_]*)`', l)
        if k:
            kam[ent].add(k.group(1))

# -- 4. cocokkan -------------------------------------------------------------
DALAM = lambda r: not r['luar'] and not r['gel'] and r['induk'] in tabel and r['anak'] in tabel
erd_dalam = [r for r in erd_rel if DALAM(r)]
erd_luar = [r for r in erd_rel if not DALAM(r)]

pas_erd = Counter((r['induk'], r['anak']) for r in erd_dalam)
pas_ddl = Counter((r['induk'], r['anak']) for r in ddl_fk)

hilang_di_ddl = []
for r in erd_dalam:
    k = (r['induk'], r['anak'])
    if pas_ddl[k] < pas_erd[k]:
        if not any(x is r for x in hilang_di_ddl):
            hilang_di_ddl.append(r)
        pas_ddl[k] += 1        # supaya pasangan berganda terhitung sekali-sekali

hilang_di_erd = []
pas_erd2 = Counter((r['induk'], r['anak']) for r in erd_dalam)
for r in ddl_fk:
    k = (r['induk'], r['anak'])
    if pas_erd2[k] > 0:
        pas_erd2[k] -= 1
    else:
        hilang_di_erd.append(r)

hapus_tak_sampai = [r for r in erd_dalam if r['hapus']]
od_ada = [r for r in ddl_fk if r['on_delete']]

# kolom yang ERD sebut tetapi tidak ada di kamus
kolom_hilang = []
for r in erd_rel:
    m = re.search(r'\b(ID_[A-Z0-9_]+)\b', r['catatan'] or '')
    if m and r['anak'] in kam and m.group(1) not in kam[r['anak']]:
        kolom_hilang.append((r['anak'], m.group(1), r['induk']))

print('=== cocok-relasi-erd-versus-ddl.py ===')
print()
print('CACAH:')
print('   baris di dalam blok sec 2        : %d' % n_baris_blok)
print('   terbaca sebagai relasi           : %d' % len(erd_rel))
print('     -- di dalam skema ini          : %d' % len(erd_dalam))
print('     -- ke luar skema / di luar gelombang / lintas sekat : %d' % len(erd_luar))
print('   FOREIGN KEY di ddl-usulan/       : %d' % len(ddl_fk))
print()
print('DITOLAK -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-58s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
print()
print('SELISIH DUA ARAH:')
print('   relasi di ERD.md sec 2, TIDAK ADA FOREIGN KEY-nya : %d' % len(hilang_di_ddl))
for r in hilang_di_ddl:
    print('      - %-22s -> %-22s [hapus: %-10s] %s'
          % (r['induk'], r['anak'], r['hapus'], r['catatan'][:40]))
print('   FOREIGN KEY di DDL, TIDAK ADA di ERD.md sec 2      : %d' % len(hilang_di_erd))
for r in hilang_di_erd:
    print('      - %-22s -> %-22s (%s)' % (r['induk'], r['anak'], r['kolom']))
print()
print('PERILAKU HAPUS:')
print('   relasi dalam-skema yang ERD.md PUTUSKAN hapusnya : %d dari %d'
      % (len(hapus_tak_sampai), len(erd_dalam)))
print('   FOREIGN KEY yang MEMBAWA ON DELETE di DDL        : %d dari %d'
      % (len(od_ada), len(ddl_fk)))
if len(od_ada) == 0 and hapus_tak_sampai:
    print('   -> KEPUTUSANNYA ADA, TETAPI TIDAK SAMPAI KE DDL.')
    print('      INV-18 menuntut perilaku hapus DITETAPKAN SADAR. Ia ditetapkan')
    print('      di ERD.md sec 2 dan HILANG di jalan menuju berkas yang akan dibangun.')
print()
print('KOLOM YANG ERD SEBUT, TIDAK ADA DI KAMUS-KOLOM.md : %d' % len(kolom_hilang))
for anak, kol, induk in kolom_hilang:
    print('      - %s.%s  (relasi ke %s)' % (anak, kol, induk))
print()
print('-- BATAS PERKAKAS INI --')
print('   Ia TIDAK memeriksa kardinalitas yang ERD.md tulis (lambang 1 / o / <)')
print('   terhadap UNIQUE di DDL. Itu pemeriksaan KETIGA, dan ia BELUM ADA.')
print('   Ia juga tidak memeriksa kolom bukan-kunci sama sekali.')

json.dump({'erd': erd_rel, 'ddl': ddl_fk,
           'hilang_di_ddl': hilang_di_ddl, 'hilang_di_erd': hilang_di_erd,
           'kolom_hilang': kolom_hilang,
           'on_delete_di_ddl': len(od_ada), 'tolak': dict(tolak)},
          io.open(os.path.join(AKAR, 'alat', 'cocok-relasi.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)
