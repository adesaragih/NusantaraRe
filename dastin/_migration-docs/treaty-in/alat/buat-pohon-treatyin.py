# -*- coding: utf-8 -*-
r"""
buat-pohon-treatyin.py — membangun pohon clipboard `TreatyIn` dari ekspor XML Pega.

Sumber : D:\XML_NURE\Treaty In            (329 berkas)  — POHON DIBANGUN DARI SINI
         D:\XML_NURE\Treaty In Adjustment (379 berkas)  — HANYA untuk menandai simpul
                                                          mana yang juga muncul di sana.

EMBARGO. Modul Adjustment dibaca hanya untuk mengetahui SIMPUL APA YANG ADA dan APA
YANG MENYAMBUNGKANNYA ke Treaty In. Perilakunya tidak diadili di berkas ini.

Keluaran (semuanya TURUNAN; bila beda dari XML, ALAT INI yang salah):
  4-erd-dan-tabel-datar/datar-treatyin-lama.csv            satu baris per simpul
  4-erd-dan-tabel-datar/datar-treatyin-kelas.csv           peta kelas Pega
  4-erd-dan-tabel-datar/datar-treatyin-muatan-keluar.csv   isian ke prosedur Oracle
  4-erd-dan-tabel-datar/pohon-treatyin.txt                 pohon siap tempel

Jalankan:  PYTHONIOENCODING=utf-8 python alat/buat-pohon-treatyin.py
"""
import re, os, csv, glob, collections

AKAR     = r'D:\XML_NURE'
SRC_IN   = os.path.join(AKAR, 'Treaty In')
SRC_ADJ  = os.path.join(AKAR, 'Treaty In Adjustment')
KELUARAN = os.path.join(AKAR, '_migration-docs', 'treaty-in', '4-erd-dan-tabel-datar')
ROOT     = 'TreatyIn'
KELAS_AKAR = 'ASM-FW-GISFW-INT-TREATY_IN'


def berkas(d):
    return sorted(glob.glob(os.path.join(d, '**', '*.xml'), recursive=True))


def isi(f):
    with open(f, encoding='utf-8', errors='replace') as h:
        return h.read()


def nisbi(f, base):
    return os.path.relpath(f, base).replace('\\', '/')


# --------------------------------------------------------------- 1. panen jalur
PAT_JALUR = re.compile(
    r'\b(TreatyIn|Primary)((?:\.[A-Za-z][A-Za-z0-9_]*(?:\([^()]{0,40}\))?)+)')
PAT_RUAS = re.compile(r'\.([A-Za-z][A-Za-z0-9_]*)(\([^()]{0,40}\))?')

# ruas yang bukan properti bisnis: metadata Pega
BUANG_RUAS = {
    'pxResults', 'pxObjClass', 'pzInsKey', 'pxInsName', 'pxCreateDateTime',
    'pxUpdateDateTime', 'pxCommitDateTime', 'pxCreateOperator', 'pxUpdateOperator',
    'pxCreateOpName', 'pxUpdateOpName', 'pxCreateSystemID', 'pxUpdateSystemID',
    'pxSaveDateTime', 'pxInsId', 'pzIndexes', 'pxSubscript', 'pxPages',
    'pxIsMappedProperty', 'pxIndexCount', 'pxListSubscript', 'pxError',
}


class Simpul(object):
    def __init__(self, jalur):
        self.jalur = jalur
        self.nama = jalur.split('.')[-1]
        self.kedalaman = jalur.count('.')
        self.list_ = False
        self.punya_anak = False
        self.ref = 0
        self.berkas = set()
        self.kelas = ''
        self.sumber_kelas = ''
        self.di_adj = False
        self.hanya_adj = False
        self.contoh = []
        self.lewat_primary = False

    @property
    def jenis(self):
        if self.list_:
            return 'Page List'
        if self.punya_anak:
            return 'Page'
        return 'Skalar'


# Setiap jalur yang ditolak penjaga `primary_ok`, beserta berkas dan kelas aturannya.
DITOLAK = []


def panen(files, base, pohon, tandai_adj=False):
    for f in files:
        s = isi(f)
        rn = nisbi(f, base)
        # `Primary` hanya berarti TreatyIn bila ATURANNYA SENDIRI applies-to kelas itu.
        # Ambil pyClassName PERTAMA: itu kelas aturannya. pyClassName berikutnya di
        # berkas yang sama datang dari entri indeks yang merujuk ATURAN LAIN, dan
        # memakainya akan memasukkan simpul palsu (mis. TreatyIn.SpreadingListXOL di
        # tingkat akar, padahal ia hidup di bawah Share[]). CONTEXT.md §2.8.
        m = re.search(r'<pyClassName>(.*?)</pyClassName>', s, re.S)
        kelas_rule = m.group(1).strip().upper() if m else ''
        primary_ok = (kelas_rule == KELAS_AKAR)
        for mm in PAT_JALUR.finditer(s):
            via_primary = (mm.group(1) == 'Primary')
            if via_primary and not primary_ok:
                # ATURAN PERKAKAS: penolakan tidak pernah diam.
                # Yang ditolak di sini adalah jalur milik KELAS ANAK. Ia benar
                # ditolak dari tingkat akar — tetapi bila tidak dihitung, yang
                # tersisa terbaca sebagai keseluruhan. Itulah yang menyembunyikan
                # 414 properti selama tujuh sesi (L-8).
                DITOLAK.append((rn, kelas_rule, mm.group(2)[:80]))
                continue
            jalur = ROOT
            for ruas, idx in PAT_RUAS.findall(mm.group(2)):
                if ruas in BUANG_RUAS:
                    break
                induk = jalur
                jalur = jalur + '.' + ruas
                n = pohon.setdefault(jalur, Simpul(jalur))
                if idx:
                    n.list_ = True
                if tandai_adj:
                    n.di_adj = True
                else:
                    n.ref += 1
                    n.berkas.add(rn)
                    if via_primary:
                        n.lewat_primary = True
                    if len(n.contoh) < 3 and rn not in n.contoh:
                        n.contoh.append(rn)
                pohon[induk].punya_anak = True


pohon = {ROOT: Simpul(ROOT)}
F_IN = berkas(SRC_IN)
F_ADJ = berkas(SRC_ADJ)
panen(F_IN, SRC_IN, pohon)
panen(F_ADJ, SRC_ADJ, pohon, tandai_adj=True)

# simpul yang HANYA ada di Adjustment dicatat terpisah, tidak dibuang diam-diam
hanya_adj = sorted(k for k, v in pohon.items() if v.ref == 0 and k != ROOT)
for k in hanya_adj:
    pohon[k].hanya_adj = True

pohon[ROOT].kelas = 'ASM-FW-GISFW-Int-TREATY_IN'
pohon[ROOT].sumber_kelas = 'dideklarasikan'


# --------------------------------------------------------------- 2. peta kelas
def petakan_kelas(files):
    peta = collections.defaultdict(collections.Counter)
    pasang = [
        ('pyPagesAndClassesPage', 'pyPagesAndClassesClass'),
        ('pxPageName', 'pxPageClass'),
        ('pyStepPageReference', 'pyStepsClassName'),
    ]
    for f in files:
        s = isi(f)
        for blk in re.findall(r'<rowdata[^>]*>(.*?)</rowdata>', s, re.S):
            for tp, tc in pasang:
                p = re.search('<%s>(.*?)</%s>' % (tp, tp), blk, re.S)
                c = re.search('<%s>(.*?)</%s>' % (tc, tc), blk, re.S)
                if p and c and p.group(1).strip() and c.group(1).strip():
                    peta[p.group(1).strip()][c.group(1).strip()] += 1
    return peta


def bersih(h):
    return re.sub(r'\([^()]*\)', '', h).strip()


peta_hal = petakan_kelas(F_IN)
for hal, ctr in peta_hal.items():
    j = bersih(hal)
    if j in pohon and not pohon[j].kelas:
        pohon[j].kelas = ctr.most_common(1)[0][0]
        pohon[j].sumber_kelas = 'dideklarasikan'


# --------------------------------------------------------------- 3. kelas -> properti
def inventaris_kelas(files):
    inv = collections.defaultdict(collections.Counter)
    for f in files:
        s = isi(f)
        for blk in re.findall(r'<rowdata[^>]*>(.*?)</rowdata>', s, re.S):
            # entri indeks properti dikenali dari pxRuleObjClass, BUKAN pxRuleFamilyName
            # (pxRuleFamilyName di ekspor ini berisi NAMA aturannya, bukan familinya)
            oc = re.search(r'<pxRuleObjClass>(.*?)</pxRuleObjClass>', blk, re.S)
            if not oc or 'Property' not in oc.group(1):
                continue
            c = re.search(r'<pxRuleClassName>(.*?)</pxRuleClassName>', blk, re.S)
            n = re.search(r'<pyRuleName>(.*?)</pyRuleName>', blk, re.S)
            if c and n:
                inv[c.group(1).strip()][n.group(1).strip()] += 1
    return inv


inv_kelas = inventaris_kelas(F_IN)

# simpul list tanpa kelas: coba samakan lewat inventaris kelas
# (nama anaknya cocok dengan daftar properti satu kelas)
anak_dari = collections.defaultdict(list)
for j, n in pohon.items():
    if '.' in j:
        anak_dari[j.rsplit('.', 1)[0]].append(n.nama)

# Cabang rekursif: halaman berbentuk TreatyIn di dalam TreatyIn.
# TIGA dinyatakan pyPagesAndClasses; OLDDATA TIDAK dinyatakan di mana pun —
# kelasnya disimpulkan dari bentuk anaknya, dan itu ditulis apa adanya.
REKURSIF = {
    'TreatyIn.ActualValue':       'dideklarasikan (27 kali)',
    'TreatyIn.ValueDifference':   'dideklarasikan (12 kali)',
    'TreatyIn.ValueBeforeProrate': 'dideklarasikan (1 kali)',
    'TreatyIn.OLDDATA':           'TIDAK dideklarasikan - bentuk anak sama dengan akar',
}
for j, asal in REKURSIF.items():
    if j in pohon:
        pohon[j].kelas = 'ASM-FW-GISFW-Int-TREATY_IN'
        pohon[j].sumber_kelas = asal

# simpul wadah tanpa kelas: samakan lewat inventaris properti per kelas.
# Syaratnya ketat — mayoritas anaknya harus ada di daftar properti kelas itu,
# DAN kelas itu harus yang paling banyak menutupi. Di bawah ambang: dibiarkan kosong.
for j, n in sorted(pohon.items()):
    if n.kelas or n.jenis == 'Skalar':
        continue
    anak = set(anak_dari.get(j, []))
    if len(anak) < 2:
        continue
    terbaik, skor, tumpang = '', 0.0, 0
    for kelas, props in inv_kelas.items():
        p = set(props)
        if not p:
            continue
        t = len(anak & p)
        cocok = t / float(len(anak))
        if (cocok, t) > (skor, tumpang):
            terbaik, skor, tumpang = kelas, cocok, t
    if terbaik and skor >= 0.75 and tumpang >= 2:
        n.kelas = terbaik
        n.sumber_kelas = 'identitas-kelas (%d/%d anak cocok)' % (tumpang, len(anak))


# --------------------------------------------------------------- 4. muatan keluar
PROSEDUR = re.compile(r'\b(PEGA_[A-Z0-9_]+)\s*\(', re.I)
# kendali UI Pega, bukan jalur data — tidak ikut muatan keluar
BUKAN_DATA = {'PEGA_CONTROL_CHECKBOX', 'PEGA_CKEDITOR_CLOSE'}


def muatan_keluar(files, base):
    baris = []
    for f in files:
        s = isi(f)
        rn = nisbi(f, base)
        for m in PROSEDUR.finditer(s):
            nama = m.group(1).upper()
            if nama in BUKAN_DATA:
                continue
            # panggilan menjangkau banyak baris; ambil sampai kurung tutup terakhir
            seg = s[m.end(): m.end() + 20000]
            tutup = seg.find(');')
            if tutup > 0:
                seg = seg[:tutup]
            args = re.findall(r'\{([^{}]{1,200})\}', seg)
            for urut_arg, a in enumerate(args, 1):
                baris.append((nama, urut_arg, a.strip(), rn))
    seen, out = set(), []
    for r in baris:
        if r in seen:
            continue
        seen.add(r)
        out.append(r)
    return sorted(out, key=lambda r: (r[0], r[1]))


keluar = muatan_keluar(F_IN, SRC_IN)


# --------------------------------------------------------------- 5. tulis CSV
if not os.path.isdir(KELUARAN):
    os.makedirs(KELUARAN)


def tulis_csv(nama, kepala, baris):
    p = os.path.join(KELUARAN, nama)
    with open(p, 'w', newline='', encoding='utf-8-sig') as h:
        w = csv.writer(h)
        w.writerow(kepala)
        w.writerows(baris)
    return p, len(baris)


urut = sorted(pohon.values(), key=lambda n: n.jalur)
b1 = [[n.jalur, n.kedalaman, n.nama, n.jenis, n.kelas or '(tidak dideklarasikan)',
       len(anak_dari.get(n.jalur, [])), n.ref, len(n.berkas),
       'hanya-adjustment' if n.hanya_adj else
       ('alias-primary' if n.lewat_primary and n.ref else 'jalur-penuh'),
       n.sumber_kelas, 'ya' if n.di_adj else '', ' | '.join(n.contoh)]
      for n in urut]
p1, c1 = tulis_csv('datar-treatyin-lama.csv',
                   ['JALUR', 'KEDALAMAN', 'NAMA', 'JENIS', 'KELAS_PEGA', 'CACAH_ANAK',
                    'REF', 'CACAH_BERKAS', 'BUKTI', 'SUMBER_KELAS', 'ADA_DI_ADJUSTMENT',
                    'BERKAS_CONTOH'], b1)

pakai = collections.defaultdict(list)
for n in urut:
    if n.kelas:
        pakai[n.kelas].append(n.jalur)
b2 = [[k, len(v), ' | '.join(v)] for k, v in sorted(pakai.items(),
                                                    key=lambda kv: -len(kv[1]))]
p2, c2 = tulis_csv('datar-treatyin-kelas.csv',
                   ['KELAS_PEGA', 'CACAH_PEMAKAIAN', 'DIPAKAI_SEBAGAI'], b2)

p3, c3 = tulis_csv('datar-treatyin-muatan-keluar.csv',
                   ['PROSEDUR', 'URUT_ARGUMEN', 'DIISI_DARI', 'BERKAS'], keluar)


# --------------------------------------------------------------- 6. pohon teks
def gambar():
    anak = collections.defaultdict(list)
    for n in urut:
        if '.' in n.jalur:
            anak[n.jalur.rsplit('.', 1)[0]].append(n)
    keluar_baris = []

    def tanda(n):
        return {'Page List': '[]', 'Page': '{}', 'Skalar': ''}[n.jenis]

    def rek(j, pre, batas=4):
        kids = sorted(anak.get(j, []), key=lambda n: (-n.ref, n.nama))
        wadah = [k for k in kids if k.jenis != 'Skalar']
        skalar = [k for k in kids if k.jenis == 'Skalar']
        if skalar:
            nama = ' \u00b7 '.join(k.nama for k in skalar[:8])
            sisa = ' \u2026 (+%d)' % (len(skalar) - 8) if len(skalar) > 8 else ''
            keluar_baris.append('%s\u251c\u2500 %d skalar: %s%s'
                                % (pre, len(skalar), nama, sisa))
        for i, k in enumerate(wadah):
            akhir = (i == len(wadah) - 1)
            cab = '\u2514\u2500 ' if akhir else '\u251c\u2500 '
            label = '%s%s' % (k.nama, tanda(k))
            kel = k.kelas.replace('ASM-FW-', '') if k.kelas else '(tidak dideklarasikan)'
            # cabang rekursif digambar SEKALI, lalu dirujuk \u2014 kalau tidak, pohon
            # ini mencetak seluruh isi TreatyIn tiga kali dan tak terbaca
            if k.jalur in REKURSIF:
                keluar_baris.append(
                    '%s%s%-34s \u2192 %-40s n=%-4d \u25c4\u2500 SALINAN UTUH pohon TreatyIn '
                    '[%s]' % (pre, cab, label, kel, k.ref, REKURSIF[k.jalur]))
                continue
            keluar_baris.append('%s%s%-34s \u2192 %-40s n=%d' % (pre, cab, label, kel, k.ref))
            rek(k.jalur, pre + ('     ' if akhir else '\u2502    '), batas)

    r = pohon[ROOT]
    keluar_baris.append('%-37s \u2192 %-40s n=%d'
                        % (ROOT + '{}', r.kelas.replace('ASM-FW-', ''),
                           sum(1 for n in urut if n.kedalaman == 1)))
    rek(ROOT, '')
    return '\n'.join(keluar_baris)


teks = gambar()
p4 = os.path.join(KELUARAN, 'pohon-treatyin.txt')
with open(p4, 'w', encoding='utf-8') as h:
    h.write(teks + '\n')

# --------------------------------------------------------------- 7. ringkasan
n_simpul = len([n for n in urut if not n.hanya_adj])
n_skalar1 = len([n for n in urut if n.kedalaman == 1 and n.jenis == 'Skalar'])
n_wadah1 = len([n for n in urut if n.kedalaman == 1 and n.jenis != 'Skalar'])
print('berkas Treaty In        : %d' % len(F_IN))
print('berkas Adjustment       : %d' % len(F_ADJ))
print('simpul pohon TreatyIn   : %d' % n_simpul)
print('  hanya di Adjustment   : %d' % len(hanya_adj))
print('  kedalaman maksimum    : %d' % max(n.kedalaman for n in urut))
print('  skalar langsung       : %d' % n_skalar1)
print('  wadah langsung        : %d' % n_wadah1)
print('kelas Pega berbeda      : %d' % len(pakai))
print('tanpa deklarasi kelas   : %d' % len([n for n in urut
                                            if n.jenis != 'Skalar' and not n.kelas]))
print('muatan keluar (baris)   : %d' % c3)
# --- Yang DITOLAK, dilaporkan. Lihat 4-erd-dan-tabel-datar/TITIK-BUTA-POHON.md ---
print('---')
print('DITOLAK penjaga primary_ok: %d rujukan jalur' % len(DITOLAK))
print('  dari %d berkas, %d kelas aturan bukan-akar'
      % (len({d[0] for d in DITOLAK}), len({d[1] for d in DITOLAK})))
print('  -> jalur ini milik KELAS ANAK dan TIDAK ada di pohon mana pun.')
print('  -> daftar propertinya: jalankan alat/sapu-properti-per-kelas.py')
for kls in sorted({d[1] for d in DITOLAK}):
    print('     %-44s %d' % (kls, len([d for d in DITOLAK if d[1] == kls])))
print('---')
print('%s (%d)' % (p1, c1))
print('%s (%d)' % (p2, c2))
print('%s (%d)' % (p3, c3))
print('%s' % p4)
