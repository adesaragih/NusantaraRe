# -*- coding: utf-8 -*-
"""
tulis.py - menyapu SETIAP JENIS PENULIS sebuah properti di ekspor Pega.

Lahir dari MA-05: sapuan penulis memakai nama tag yang salah (`pyTarget`, `pyName`)
mengembalikan NOL, dan nol itu sempat dibaca sebagai jawaban. Perkakas ini menutup
lubang itu dengan dua cara:

  1. Ia menyapu SEMUA jenis penulis yang ada di ekspor, bukan satu.
  2. Ia punya mode KALIBRASI: sebelum mempercayai hasil nol, buktikan dulu bahwa
     penyapu jenis itu menemukan kasus positif yang sudah diketahui.

Dan sesuai aturan proyek "perkakas yang menolak sesuatu harus melaporkan apa yang
ditolaknya": setiap sapuan mencetak berapa kemunculan yang dibuang sebagai SALINAN
terbungkus (`pyIncludedRuleXML`, `pyRuleVersionsList`), dan berapa berkas yang gagal
diurai.

JENIS PENULIS YANG DIKENALI
---------------------------
  PS   Activity Property-Set      pyParamArray/rowdata/PropertiesName + PropertiesValue
  DT   Data Transform             pyProperties/rowdata/pyPropertiesName + pyPropertiesValue
                                  (menghormati pyDisabled = true - penanda mati KETIGA)
  SEC  Kontrol Section            pyPropertyTarget (sasaran tambahan sebuah kontrol)
  BIND Pengikatan kontrol         pyValue pada sel yang pyReadOnly != true
  RD   Report Definition          pyTargetProperty (kolom hasil -> properti halaman hasil)

YANG TIDAK DAPAT DIBACA, dan karena itu dilaporkan sebagai TITIK BUTA
---------------------------------------------------------------------
  JAVA langkah Java              teksnya bebas; penulisan di dalamnya tidak terurai
  SQL  RDB-List / Connect-REST   menulis ke basis data, bukan ke clipboard

PEMAKAIAN
---------
  python tulis.py sapu  TreatyIn.EDMState           # semua jenis
  python tulis.py sapu  .CedingID --jenis PS,DT,SEC
  python tulis.py nilai TreatyIn.Position           # daftar nilai yang ditulis
  python tulis.py kalibrasi                          # uji setiap penyapu atas kasus positif
"""
import os, sys, xml.etree.ElementTree as ET, collections

AKAR = [r"D:\XML_NURE\Treaty In", r"D:\XML_NURE\Treaty In Adjustment"]
BUNGKUS = ('pyIncludedRuleXML', 'pyRuleVersionsList')

# Kasus positif yang SUDAH DIKETAHUI, satu per jenis penulis.
# Dipakai mode kalibrasi: penyapu yang tidak menemukannya tidak boleh dipercaya.
KALIBRASI = {
    'PS':   ('TreatyIn.RevisionState', 'SetTreatyIn_Act langkah 7 menulis 1'),
    'DT':   ('TreatyIn.Position',      'Akseptasi_DT menulis lima nilai peran'),
    'SEC':  ('TreatyIn.CedingID',      'autocomplete Ceding di TreatyInNONProportional'),
    'BIND': ('TreatyIn.EDMState',      'radio group di PickerTreatyInMasterRevisi'),
    'RD':   ('.ID',                    'BrowseTREATY_IN memetakan kolom ID'),
}


def _isi(el):
    """Jalan-jalan di BADAN aturan; subpohon salinan terbungkus dilewati dan dihitung."""
    ditolak = [0]
    def jalan(e):
        for c in e:
            if c.tag in BUNGKUS:
                ditolak[0] += sum(1 for _ in c.iter())
                continue
            yield c
            for x in jalan(c):
                yield x
    return jalan(el), ditolak


def _cocok(teks, sasaran):
    """Cocokkan nama properti. '.X' cocok dengan 'apa pun.X'; nama penuh cocok persis."""
    if not teks:
        return False
    t = teks.strip()
    if sasaran.startswith('.'):
        return t == sasaran or t.endswith(sasaran)
    return t == sasaran or t.endswith('.' + sasaran.split('.')[-1]) and t.split('.')[-1] == sasaran.split('.')[-1] and t == sasaran


def _berkas():
    for rt in AKAR:
        for dp, dn, fn in os.walk(rt):
            for f in fn:
                if f.endswith('.xml'):
                    yield os.path.join(dp, f)


def sapu(sasaran, jenis=None):
    jenis = jenis or ['PS', 'DT', 'SEC', 'BIND', 'RD']
    hasil = collections.defaultdict(list)
    salinan = 0
    gagal = []
    buta = collections.Counter()
    for p in _berkas():
        try:
            r = ET.parse(p).getroot()
        except Exception as ex:
            gagal.append((p, str(ex)[:60]))
            continue
        it, tolak = _isi(r)
        rel = os.path.relpath(p, r"D:\XML_NURE")
        induk = {}
        akar_anak = list(r)
        pohon = [r] + [e for e in _isi(r)[0]]
        par = {}
        for e in pohon:
            for c in e:
                if c.tag not in BUNGKUS:
                    par[c] = e
        for e in pohon:
            if e.tag == 'rowdata':
                # PS - Activity Property-Set
                if 'PS' in jenis:
                    n = e.findtext('PropertiesName')
                    if _cocok(n, sasaran):
                        hasil['PS'].append((rel, n.strip(), (e.findtext('PropertiesValue') or '').strip(), ''))
                # DT - Data Transform
                if 'DT' in jenis:
                    n = e.findtext('pyPropertiesName')
                    if _cocok(n, sasaran):
                        mati = (e.findtext('pyDisabled') or '').strip() == 'true'
                        hasil['DT'].append((rel, (e.findtext('pyPropertyStepId') or '').strip(),
                                            (e.findtext('pyPropertiesValue') or '').strip(),
                                            'MATI (pyDisabled)' if mati else ''))
                # RD
                if 'RD' in jenis:
                    n = e.findtext('pyTargetProperty')
                    if _cocok(n, sasaran):
                        hasil['RD'].append((rel, n.strip(), (e.findtext('pyColumnName') or '').strip(), ''))
                # BIND - pengikatan kontrol yang dapat disunting
                if 'BIND' in jenis and (e.findtext('pxObjClass') or '') == 'Embed-Display-Table-Cell':
                    n = e.findtext('pyValue')
                    if _cocok(n, sasaran):
                        ro = (e.findtext('pyReadOnly') or '').strip()
                        hasil['BIND'].append((rel, n.strip(), (e.findtext('pyFormat') or '').strip(),
                                              'hanya-baca' if ro == 'true' else 'dapat disunting'))
            if 'SEC' in jenis and e.tag == 'pyPropertyTarget' and _cocok(e.text, sasaran):
                hasil['SEC'].append((rel, (e.text or '').strip(), '', ''))
        # titik buta
        for e in pohon:
            if e.tag == 'pyStepsActivityName' and e.text:
                t = e.text.strip()
                if t == 'Java':
                    buta['JAVA'] += 1
                elif t in ('RDB-List', 'Connect-REST', 'Obj-Browse'):
                    buta['SQL/REST'] += 1
        salinan += tolak[0]
    return hasil, salinan, gagal, buta


def _cetak(sasaran, hasil, salinan, gagal, buta):
    print("=" * 78)
    print("SASARAN: %s" % sasaran)
    tot = 0
    for j in ('PS', 'DT', 'SEC', 'BIND', 'RD'):
        baris = hasil.get(j, [])
        tot += len(baris)
        print("  %-5s %d penulisan" % (j, len(baris)))
        for b in baris:
            print("        %-52s %-26s %-28s %s" % (b[0][:52], b[1][:26], b[2][:28], b[3]))
    print("  TOTAL %d penulisan di badan aturan" % tot)
    print("  DITOLAK: %d simpul di dalam salinan terbungkus (tidak dihitung)" % salinan)
    print("  GAGAL URAI: %d berkas" % len(gagal))
    for g in gagal:
        print("        %s - %s" % g)
    print("  TITIK BUTA (tidak terurai, bukan berarti nol): %s" %
          (", ".join("%s %d langkah" % (k, v) for k, v in sorted(buta.items())) or "tidak ada"))


def kalibrasi():
    print("KALIBRASI - setiap penyapu diuji atas kasus positif yang sudah diketahui")
    print("Penyapu yang menemukan NOL di sini tidak boleh dipercaya (MA-05, TA-04).\n")
    ok = True
    for j, (sasaran, catatan) in sorted(KALIBRASI.items()):
        hasil, _, _, _ = sapu(sasaran, [j])
        n = len(hasil.get(j, []))
        tanda = "LULUS" if n > 0 else "GAGAL"
        if n == 0:
            ok = False
        print("  %-5s %-6s %-26s %3d penulisan   (%s)" % (j, tanda, sasaran, n, catatan))
    print("\n%s" % ("Semua penyapu terkalibrasi." if ok else "ADA PENYAPU YANG GAGAL - jangan pakai hasilnya."))
    return ok


def nilai(sasaran):
    hasil, salinan, gagal, buta = sapu(sasaran)
    c = collections.Counter()
    for j, baris in hasil.items():
        for b in baris:
            if j in ('PS', 'DT'):
                c[(b[2], b[3])] += 1
    print("NILAI YANG DITULIS ke %s:" % sasaran)
    for (v, m), n in c.most_common():
        print("   %-40s %3d %s" % (v[:40], n, m))
    _cetak(sasaran, hasil, salinan, gagal, buta)


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print(__doc__)
    elif sys.argv[1] == 'kalibrasi':
        kalibrasi()
    elif sys.argv[1] == 'nilai':
        nilai(sys.argv[2])
    else:
        j = None
        if '--jenis' in sys.argv:
            j = sys.argv[sys.argv.index('--jenis') + 1].split(',')
        h, s, g, b = sapu(sys.argv[2], j)
        _cetak(sys.argv[2], h, s, g, b)
