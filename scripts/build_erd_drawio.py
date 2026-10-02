# -*- coding: utf-8 -*-
"""Bangun ERD .drawio (bisa diedit di draw.io / diagrams.net) untuk Claim Life + Komite Life.
Kotak = tabel (PK/FK di dalam), garis = relasi PK->FK. Struktur final D1-D8 (2026-09-18).
"""
import html

OUT = r"d:\XML\RNM_BRD\OUTPUT_HASIL_RNM\ERD-Claim-Komite-Life.drawio"

# tiap tabel: id, judul, baris kolom [(nama, tag)], posisi x,y, warna header
# tag: 'PK','FK','PKFK','' 
tabel = {
 "work": ("T_WORK_CLAIM  (AKAR, lintas-lini)", [
     ("ID  (CLM-xxxxxx / KMT-xxxxxx)","PK"),
     ("COVER_KEY -> T_WORK_CLAIM.ID","FK"),
     ("LINI / PY_POSITION / ACCEPT_STATUS",""),
     ("CASEID / CREATE_OP / TGL_UPDATE",""),
 ], 40, 40, "#1F4E78"),

 "gclaim": ("T_GENERAL_CLAIM  (header klaim)", [
     ("ID = T_WORK_CLAIM.ID (shared PK)","PKFK"),
     ("CASEID_POLICY  (kunci ke polis)",""),
     ("POLICY_NO / ENDORSMENT_NO",""),
     ("CLAIM_NO / RISLIPRNM / BUSINESS_NAME",""),
 ], 40, 240, "#2E75B6"),

 "det": ("T_CLAIMLF_PREMIUMLIST_DETAIL  (peserta)", [
     ("ID","PK"),
     ("CLAIM_ID -> T_GENERAL_CLAIM.ID","FK"),
     ("IS_CHECK + 8 kolom audit",""),
 ], 40, 430, "#2E75B6"),

 "adj": ("T_CLAIMLF_ADJUSTMENT  (adjustment)", [
     ("ID","PK"),
     ("PREMIUM_LIST_DETAIL_ID -> ...DETAIL.ID","FK"),
     ("KOMITE_ID -> T_WORK_CLAIM.ID (komite)","FK"),
     ("CLAIM_AMOUNT / STS_REJECT / bank",""),
 ], 40, 610, "#2E75B6"),

 "spr": ("T_CLAIMLF_ADJUSTMENT_SPREADING", [
     ("ID","PK"),
     ("ADJUSTMENT_ID -> ...ADJUSTMENT.ID","FK"),
     ("TREATY_TYPE / RATE / IDR / USD",""),
 ], 40, 800, "#2E75B6"),

 "sprr": ("T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", [
     ("ID","PK"),
     ("SPREADING_ID -> ...SPREADING.ID","FK"),
     ("PREMIUM_SPREADED_NET = GROSS-Diskon-Komisi",""),
 ], 40, 980, "#2E75B6"),

 "doc": ("DOCUMENT_CLAIM  (lintas-lini)", [
     ("ID","PK"),
     ("PREMIUM_LIST_DETAIL_ID -> ...DETAIL.ID","FK"),
 ], 470, 430, "#548235"),

 "gkom": ("T_GENERAL_KOMITE  (header komite)", [
     ("ID = T_WORK_CLAIM.ID (shared PK)","PKFK"),
     ("ADJUSTMENT_ID -> ...ADJUSTMENT.ID","FK"),
     ("KOMITE_LOOP / KOMITE_COUNT / ACCEPT_STATUS",""),
 ], 560, 240, "#548235"),

 "klist": ("T_KOMITE_KOMITELIST  (per anggota)", [
     ("ID","PK"),
     ("DATA_KOMITE_ID -> T_GENERAL_KOMITE.ID","FK"),
     ("KOMITE_URUT / KOMITE_ID / APROVAL / COMMENT",""),
 ], 560, 430, "#548235"),
}

# relasi: (dari_id, ke_id, label, gaya_exit_entry)
# arah panah dari INDUK ke ANAK (mengikuti FK)
relasi = [
 ("work","work","1 COVER_KEY (self)"),
 ("work","gclaim","2 shared PK 1:1"),
 ("gclaim","det","3 CLAIM_ID (CASCADE)"),
 ("det","adj","4 PREMIUM_LIST_DETAIL_ID (CASCADE)"),
 ("adj","spr","5 ADJUSTMENT_ID (CASCADE)"),
 ("spr","sprr","6 SPREADING_ID (CASCADE)"),
 ("det","doc","7 PREMIUM_LIST_DETAIL_ID (Go)"),
 ("work","gkom","8 shared PK 1:1"),
 ("gkom","klist","9 DATA_KOMITE_ID (CASCADE)"),
 ("adj","gkom","10 ADJUSTMENT_ID (1:1)"),
 ("adj","work","11 KOMITE_ID (penunjuk N:1)"),
]

ROW_H = 26
HDR_H = 30
W = 380

def esc(s): return html.escape(s, quote=True)

cells = []
cid = 2  # 0,1 dipakai root

# buat kotak tabel (pakai swimlane style: header + baris)
tbl_cell_id = {}
for key,(judul, kolom, x, y, warna) in tabel.items():
    h = HDR_H + ROW_H*len(kolom)
    tid = f"t_{key}"
    tbl_cell_id[key] = tid
    style = (f"swimlane;html=1;startSize={HDR_H};horizontal=1;fontStyle=1;fontColor=#FFFFFF;"
             f"fillColor={warna};strokeColor=#333333;fontSize=11;align=center;")
    cells.append(f'<mxCell id="{tid}" value="{esc(judul)}" style="{style}" vertex="1" parent="1">'
                 f'<mxGeometry x="{x}" y="{y}" width="{W}" height="{h}" as="geometry"/></mxCell>')
    for i,(nama,tag) in enumerate(kolom):
        rid = f"{tid}_r{i}"
        if tag == "PK":
            fill="#FFF2CC"; pfx="[PK] "
        elif tag == "FK":
            fill="#DDEBF7"; pfx="[FK] "
        elif tag == "PKFK":
            fill="#FCE4D6"; pfx="[PK+FK] "
        else:
            fill="#FFFFFF"; pfx=""
        rstyle=(f"text;html=1;align=left;verticalAlign=middle;spacingLeft=6;fontSize=10;"
                f"fillColor={fill};strokeColor=#BFBFBF;")
        cells.append(f'<mxCell id="{rid}" value="{esc(pfx+nama)}" style="{rstyle}" vertex="1" parent="{tid}">'
                     f'<mxGeometry x="0" y="{HDR_H+ROW_H*i}" width="{W}" height="{ROW_H}" as="geometry"/></mxCell>')

# buat garis relasi
for i,(a,b,label) in enumerate(relasi):
    eid=f"e{i}"
    if a==b:
        estyle=("edgeStyle=orthogonalEdgeStyle;rounded=1;html=1;endArrow=block;"
                "strokeColor=#C55A11;strokeWidth=2;fontSize=9;exitX=1;exitY=0.2;entryX=1;entryY=0.5;")
    else:
        estyle=("edgeStyle=orthogonalEdgeStyle;rounded=1;html=1;endArrow=block;"
                "strokeColor=#1F4E78;strokeWidth=1.5;fontSize=9;")
    cells.append(f'<mxCell id="{eid}" value="{esc(label)}" style="{estyle}" edge="1" parent="1" '
                 f'source="{tbl_cell_id[a]}" target="{tbl_cell_id[b]}">'
                 f'<mxGeometry relative="1" as="geometry"/></mxCell>')

# legenda
leg=("[PK] kuning = Primary Key   |   [FK] biru = Foreign Key   |   [PK+FK] oranye = Shared PK 1:1\n"
     "Garis biru = FK anak->induk   |   Garis oranye = self-reference\n"
     "CASCADE = anak ikut terhapus (DB)   |   Go = penghapusan diatur aplikasi   |   penunjuk = rujukan saja\n"
     "Struktur final D1-D8 (2026-09-18). Biru=Claim Life, Hijau=Komite (lintas-lini).")
cells.append(f'<mxCell id="legenda" value="{esc(leg)}" '
             f'style="text;html=1;align=left;verticalAlign=top;spacingLeft=8;spacingTop=6;'
             f'fillColor=#F2F2F2;strokeColor=#999999;fontSize=10;" vertex="1" parent="1">'
             f'<mxGeometry x="470" y="640" width="470" height="110" as="geometry"/></mxCell>')

body="\n".join(cells)
xml=(f'<mxfile host="app.diagrams.net">\n'
     f'<diagram id="erd1" name="ERD Claim + Komite Life">\n'
     f'<mxGraphModel dx="1200" dy="800" grid="1" gridSize="10" guides="1" '
     f'tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" '
     f'pageWidth="1100" pageHeight="1500" math="0" shadow="0">\n'
     f'<root>\n<mxCell id="0"/>\n<mxCell id="1" parent="0"/>\n{body}\n</root>\n'
     f'</mxGraphModel>\n</diagram>\n</mxfile>\n')

with open(OUT,"w",encoding="utf-8") as f:
    f.write(xml)
print("SELESAI:", OUT)
