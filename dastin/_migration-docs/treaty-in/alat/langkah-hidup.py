# -*- coding: utf-8 -*-
r"""
langkah-hidup.py — mencetak daftar langkah sebuah Activity Pega BESERTA status bloknya.

Ada karena CONTEXT §2.0a: `pySteps` bersarang, dan langkah ber-pyStepsBlockName "//"
MATI. Peringkas apa pun yang menyembunyikan penanda itu menyesatkan.

Pakai:  python alat/langkah-hidup.py "Treaty In/Activity/NamaActivity.xml"
"""
import re, sys, os, xml.etree.ElementTree as ET

AKAR = r'D:\XML_NURE'


def teks(el, tag):
    c = el.find(tag)
    return (c.text or '').strip() if c is not None else ''


def langkah(el, dalam=0, mati_induk=False, keluar=None):
    """pySteps bersarang: satu langkah mati mematikan seluruh anaknya."""
    if keluar is None:
        keluar = []
    for ps in el.findall('pySteps'):
        for row in ps.findall('rowdata'):
            blok = teks(row, 'pyStepsBlockName')
            mati = mati_induk or blok.strip() == '//'
            keluar.append({
                'dalam': dalam,
                'mati': mati,
                'blok': blok,
                'metode': teks(row, 'pyMethodName') or teks(row, 'pyStepsMethod'),
                'hal': teks(row, 'pyStepsPage'),
                'ket': teks(row, 'pyStepsDescription'),
                'when': teks(row, 'pyStepsPreCondParamsWhen'),
                'rule': teks(row, 'pyStepsRuleName') or teks(row, 'pyActivityName'),
            })
            langkah(row, dalam + 1, mati, keluar)
    return keluar


def main(path):
    if not os.path.isabs(path):
        path = os.path.join(AKAR, path)
    s = open(path, encoding='utf-8', errors='replace').read()
    try:
        root = ET.fromstring(s)
    except ET.ParseError as e:
        print('XML tidak terurai:', e)
        return
    ls = langkah(root)
    hidup = sum(1 for x in ls if not x['mati'])
    print('%s\n  langkah: %d  (hidup %d / mati %d)\n' %
          (os.path.basename(path), len(ls), hidup, len(ls) - hidup))
    for i, x in enumerate(ls, 1):
        tanda = 'MATI' if x['mati'] else 'hidup'
        pre = '  ' * x['dalam']
        print('%3d %-5s %s%s %s' % (i, tanda, pre, x['metode'] or '-', x['rule']))
        if x['ket']:
            print('      %s  %s' % (pre, x['ket'][:100]))
        if x['when']:
            print('      %s  when: %s' % (pre, x['when'][:100]))


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print(__doc__)
    else:
        main(sys.argv[1])
