# -*- coding: utf-8 -*-
r"""
DIGANTIKAN 25 September 2026 oleh cocok-enam-sumber-struktur.py. JANGAN DIPAKAI SENDIRI.

    Perkakas ini mencetak "GERBANG 0: LULUS", dan ia lulus karena NILAI_SELISIH dan
    BESARAN_DAPAT_DISESUAIKAN DIKECUALIKAN lebih dulu ke daftar LUAR, atas dasar
    embargo SPEC-MODEL-DATA sec 11.3. Embargo itu SUDAH LEWAT -- modul Adjustment
    sudah to-spec dan to-ticket, dan tiket 06 adalah PEMBUAT PERTAMA NILAI_SELISIH.

    Pemeriksa yang mengecualikan apa yang seharusnya diperiksanya menjawab "LULUS"
    untuk kedua kemungkinan, dan jawaban semacam itu bukan jawaban.

    Ia disimpan, bukan dihapus: keluarannya adalah bukti bagaimana selisih itu
    sempat tidak terlihat. Yang dipakai sekarang cocok-enam-sumber-struktur.py.

cocok-tiga-sumber-entitas.py -- GERBANG 0: tiga sumber harus menyebut himpunan yang sama.

    2-to-spec/KAMUS-KOLOM.md            nama dan tipe kolom     (MENGIKAT, urutan wewenang butir 5)
    2-to-spec/ddl-usulan/*.sql          tabel dan kunci asing   (yang akan dibangun)
    4-erd-dan-tabel-datar/STRUKTUR-DATA.md   daftar entitas     (MENGIKAT, urutan wewenang butir 4)

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa ketiga sumber menyebut NAMA yang sama, dan mencetak selisihnya PER NAMA.

APA YANG TIDAK
    Ia TIDAK memeriksa isi -- kolom, tipe, maupun kardinalitas. Dua sumber yang
    menyebut nama yang sama dengan kolom yang berbeda akan LULUS perkakas ini.
    Ia juga TIDAK dapat menyatakan himpunannya LENGKAP: semesta yang dipakai
    menyusun sec 10 masih kurang 340 properti titik buta (L-8), dan tidak ada
    pemeriksaan di sini yang dapat melihatnya.

    STRUKTUR-DATA.md memuat entitas yang SENGAJA di luar penyerahan pertama
    (GEL-2, GEL-3) dan entitas milik modul lain. Itu BUKAN selisih, dan perkakas
    ini memisahkannya dengan DAFTAR BERNAMA -- bukan dengan tebakan.

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/cocok-tiga-sumber-entitas.py
"""
import io, os, re, json
from collections import Counter

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
KAMUS = os.path.join(AKAR, '2-to-spec', 'KAMUS-KOLOM.md')
DDL = os.path.join(AKAR, '2-to-spec', 'ddl-usulan')
STRU = os.path.join(AKAR, '4-erd-dan-tabel-datar', 'STRUKTUR-DATA.md')

tolak = Counter()

# -- 1. KAMUS: judul "## `NAMA`" ---------------------------------------------
kam_teks = io.open(KAMUS, encoding='utf-8').read()
kamus, kam_bukan = set(), []
for l in kam_teks.split('\n'):
    if not l.startswith('## '):
        continue
    m = re.match(r'^## `([A-Z][A-Z0-9_]*)`\s*$', l)
    if m:
        kamus.add(m.group(1))
    else:
        kam_bukan.append(l[3:].strip())
        tolak['judul "## " di kamus yang BUKAN nama entitas'] += 1

# -- 2. DDL: CREATE TABLE ----------------------------------------------------
ddl, ddl_berkas = set(), {}
fk = []          # (anak, kolom, induk, on_delete, berkas)
tanpa_create = []
for f in sorted(os.listdir(DDL)):
    if not f.lower().endswith('.sql'):
        tolak['berkas di ddl-usulan/ yang bukan .sql'] += 1
        continue
    s = io.open(os.path.join(DDL, f), encoding='utf-8').read()
    m = re.search(r'CREATE TABLE\s+\S+\.(\w+)\s*\(', s)
    if not m:
        tanpa_create.append(f)
        tolak['berkas .sql tanpa CREATE TABLE (berkas kerangka, bukan tabel)'] += 1
        continue
    nm = m.group(1)
    ddl.add(nm); ddl_berkas[nm] = f
    for fm in re.finditer(
            r'ADD CONSTRAINT\s+(\w+)\s+FOREIGN KEY\s*\((\w+)\)\s*\n?\s*REFERENCES\s+\S+\.(\w+)\s*\((\w+)\)\s*([^;]*);',
            s):
        od = re.search(r'ON DELETE\s+([A-Z ]+)', fm.group(5) or '')
        fk.append(dict(anak=nm, kolom=fm.group(2), induk=fm.group(3),
                       kolom_induk=fm.group(4),
                       on_delete=(od.group(1).strip() if od else '(tidak dinyatakan)'),
                       constraint=fm.group(1), berkas=f))

# -- 3. STRUKTUR-DATA: baris tabel "| `NAMA` | <angka> |" --------------------
stru_teks = io.open(STRU, encoding='utf-8').read()
stru = set(re.findall(r'^>?\s*\|\s*\*{0,2}`([A-Z][A-Z0-9_]+)`\*{0,2}\s*\|\s*\d+\s*\|',
                      stru_teks, re.M))

# -- 4. DI LUAR PENYERAHAN PERTAMA -- daftar BERNAMA, bukan tebakan ----------
# Sumbernya dinyatakan satu per satu; perkakas ini tidak menyimpulkannya dari bentuk.
LUAR = {
    'RETRO_KELUAR':              'GEL-2 -- arah keluar, SPEC-MODEL-DATA sec 2.3',
    'PENCAPAIAN':                'GEL-3 -- dimodelkan, tidak dibangun sekarang',
    'NILAI_SELISIH':             'modul Adjustment -- SPEC-MODEL-DATA sec 11.3',
    'BESARAN_DAPAT_DISESUAIKAN': 'modul Adjustment -- tabel acuan, sec 11.3',
}

# -- 5. selisih --------------------------------------------------------------
sel_kamus_tanpa_ddl = sorted(kamus - ddl)
sel_ddl_tanpa_kamus = sorted(ddl - kamus)
sel_kamus_tanpa_stru = sorted(kamus - stru)
sel_stru_tanpa_kamus = sorted(x for x in (stru - kamus) if x not in LUAR)
stru_luar = sorted(x for x in (stru - kamus) if x in LUAR)

# FK menggantung
fk_yatim = [x for x in fk if x['induk'] not in ddl]

print('=== cocok-tiga-sumber-entitas.py -- GERBANG 0 ===')
print()
print('CACAH per sumber:')
print('   %-46s %d' % ('KAMUS-KOLOM.md  (judul "## `NAMA`")', len(kamus)))
print('   %-46s %d' % ('ddl-usulan/     (CREATE TABLE)', len(ddl)))
print('   %-46s %d' % ('STRUKTUR-DATA.md (baris tabel entitas)', len(stru)))
print('   %-46s %d' % ('  -- di antaranya SENGAJA di luar penyerahan 1', len(stru_luar)))
print('   %-46s %d' % ('  -- sisanya, yang harus cocok dengan kamus', len(stru) - len(stru_luar)))
print()
print('DITOLAK / TIDAK DIHITUNG -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-62s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
if kam_bukan:
    print('   judul "## " yang ditolak, apa adanya:')
    for x in kam_bukan:
        print('      %s' % x[:70])
if tanpa_create:
    print('   berkas .sql tanpa CREATE TABLE, apa adanya:')
    for x in tanpa_create:
        print('      %s' % x)
print()
print('SELISIH PER NAMA -- tiga arah:')


def cetak(judul, daftar, arti):
    print('   %s : %d' % (judul, len(daftar)))
    if daftar:
        print('      ARTINYA: %s' % arti)
        for x in daftar:
            print('        - %s' % x)


cetak('di KAMUS, tidak ada DDL-nya ', sel_kamus_tanpa_ddl,
      'tabel yang TIDAK AKAN TERBANGUN -- LUBANG')
cetak('ada DDL, tidak ada di KAMUS ', sel_ddl_tanpa_kamus,
      'kolomnya TIDAK BERNAMA -- LUBANG YANG BERLAWANAN')
cetak('di KAMUS, tidak di STRUKTUR ', sel_kamus_tanpa_stru,
      'berkas yang MENGIKAT soal daftar entitas kehilangan entitas ini')
cetak('di STRUKTUR, tidak di KAMUS ', sel_stru_tanpa_kamus,
      'entitas yang didaftar tetapi tidak berkolom dan tidak berDDL')
print()
print('   di STRUKTUR dan SENGAJA di luar penyerahan pertama : %d' % len(stru_luar))
for x in stru_luar:
    print('        - %-28s %s' % (x, LUAR[x]))
print()
print('KUNCI ASING dibaca dari DDL : %d' % len(fk))
print('   menggantung (induknya tidak ada tabelnya) : %d' % len(fk_yatim))
for x in fk_yatim:
    print('        - %s.%s -> %s   [%s]' % (x['anak'], x['kolom'], x['induk'], x['berkas']))
print()
cocok = (not sel_kamus_tanpa_ddl and not sel_ddl_tanpa_kamus
         and not sel_kamus_tanpa_stru and not sel_stru_tanpa_kamus and not fk_yatim)
print('GERBANG 0 : %s' % ('LULUS -- ketiga sumber menyebut himpunan yang sama'
                          if cocok else 'GAGAL -- lihat selisih di atas'))
print()
print('-- BATAS PERKAKAS INI, dan ia bagian dari keluarannya --')
print('   Ia mencocokkan NAMA, bukan ISI. Dua sumber yang menyebut nama sama dengan')
print('   kolom berbeda LULUS perkakas ini. Dan "ketiganya cocok" BUKAN "himpunannya')
print('   lengkap": semesta sec 10 masih kurang 340 properti titik buta (L-8, M-4).')

json.dump({'kamus': sorted(kamus), 'ddl': sorted(ddl), 'struktur': sorted(stru),
           'ddl_berkas': ddl_berkas, 'fk': fk, 'luar': LUAR,
           'selisih': {'kamus_tanpa_ddl': sel_kamus_tanpa_ddl,
                       'ddl_tanpa_kamus': sel_ddl_tanpa_kamus,
                       'kamus_tanpa_struktur': sel_kamus_tanpa_stru,
                       'struktur_tanpa_kamus': sel_stru_tanpa_kamus},
           'fk_yatim': fk_yatim, 'lulus': cocok},
          io.open(os.path.join(AKAR, 'alat', 'himpunan-entitas.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)
