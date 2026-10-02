# -*- coding: utf-8 -*-
"""Bangun Excel skema Claim Life + Komite Life (struktur final D1-D8, 2026-09-18).
Sumber: tiket 14 Claim Life, keputusan-struktur-komite.md, keputusan D1-D8.
Semua nama tabel/kolom = keputusan work owner + [terverifikasi] korpus.
"""
from openpyxl import Workbook
from openpyxl.styles import Font, PatternFill, Alignment, Border, Side

wb = Workbook()

# ---------- gaya ----------
HDR_FILL = PatternFill("solid", fgColor="1F4E78")
HDR_FONT = Font(bold=True, color="FFFFFF", size=11)
PK_FILL  = PatternFill("solid", fgColor="FFF2CC")   # kuning muda = PK
FK_FILL  = PatternFill("solid", fgColor="DDEBF7")   # biru muda  = FK
TBL_FILL = PatternFill("solid", fgColor="E2EFDA")   # hijau muda = baris tabel
thin = Side(style="thin", color="BFBFBF")
BORDER = Border(left=thin, right=thin, top=thin, bottom=thin)
WRAP = Alignment(vertical="top", wrap_text=True)
TOP  = Alignment(vertical="top")

def style_header(ws, ncol):
    for c in range(1, ncol + 1):
        cell = ws.cell(row=1, column=c)
        cell.fill = HDR_FILL; cell.font = HDR_FONT
        cell.alignment = Alignment(vertical="center", horizontal="center", wrap_text=True)
        cell.border = BORDER

def apply_body(ws, ncol, first=2):
    for row in ws.iter_rows(min_row=first, max_row=ws.max_row, max_col=ncol):
        for cell in row:
            cell.border = BORDER
            if cell.alignment.wrap_text is None:
                cell.alignment = TOP

# ================= SHEET 1: TABEL =================
ws1 = wb.active
ws1.title = "Tabel"
ws1.append(["No", "Nama Tabel", "Tingkat", "PK (ID)", "Lintas-lini?", "Keterangan"])
tabel = [
    [1, "T_WORK_CLAIM", "1 (AKAR)", "ID (teks: CLM-xxxxxx / KMT-xxxxxx)", "YA",
     "Akar. 1 baris per work object. Baris klaim ID=CLM-..., baris komite ID=KMT-... COVER_KEY menunjuk induk."],
    [2, "T_GENERAL_CLAIM", "2", "ID (= T_WORK_CLAIM.ID baris klaim, shared PK)", "tidak",
     "Header klaim. 1:1 dgn baris klaim T_WORK_CLAIM lewat ID sama."],
    [3, "T_CLAIMLF_PREMIUMLIST_DETAIL", "3", "ID (sequence)", "tidak",
     "Peserta klaim. Anak T_GENERAL_CLAIM."],
    [4, "T_CLAIMLF_ADJUSTMENT", "4", "ID (sequence)", "tidak",
     "Baris adjustment. Anak peserta. Simpan KOMITE_ID (rujukan kasus komite)."],
    [5, "T_CLAIMLF_ADJUSTMENT_SPREADING", "5", "ID (sequence)", "tidak",
     "Spreading per treaty-year. Anak adjustment. Dibekukan."],
    [6, "T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "6", "ID (sequence)", "tidak",
     "Spreading retro per reinsurer. Anak spreading. Dibekukan."],
    [7, "DOCUMENT_CLAIM", "4", "ID (sequence)", "YA",
     "Dokumen per peserta. Lintas-lini. ON DELETE di Go."],
    [8, "T_GENERAL_KOMITE", "2", "ID (= T_WORK_CLAIM.ID baris komite, shared PK)", "YA",
     "Header kasus komite. 1:1 dgn baris komite T_WORK_CLAIM. Simpan ADJUSTMENT_ID (tutup lingkar)."],
    [9, "T_KOMITE_KOMITELIST", "3", "ID (sequence)", "YA",
     "List komite per jenjang + keputusan per anggota. Anak T_GENERAL_KOMITE."],
]
for r in tabel:
    ws1.append(r)
style_header(ws1, 6)
for r in range(2, ws1.max_row + 1):
    ws1.cell(row=r, column=2).fill = TBL_FILL
    ws1.cell(row=r, column=2).font = Font(bold=True)
    ws1.cell(row=r, column=4).fill = PK_FILL
    ws1.cell(row=r, column=6).alignment = WRAP
apply_body(ws1, 6)
for col, w in zip("ABCDEF", [4, 34, 10, 40, 12, 55]):
    ws1.column_dimensions[col].width = w
ws1.freeze_panes = "A2"

# ================= SHEET 2: KOLOM & FK =================
ws2 = wb.create_sheet("Kolom & FK")
ws2.append(["Tabel", "Kolom", "Kunci", "Tipe", "Menunjuk / Sumber", "Catatan"])
# Kunci: PK / FK / PK+FK / (kosong)
kolom = [
 # T_WORK_CLAIM
 ["T_WORK_CLAIM", "ID", "PK", "teks CLM-/KMT-", "-", "Identitas work object. Baris klaim CLM-..., baris komite KMT-..."],
 ["T_WORK_CLAIM", "COVER_KEY", "FK", "teks", "T_WORK_CLAIM.ID (self)", "ID baris induk; NULL bila tak punya induk. Baris komite -> ID baris klaim."],
 ["T_WORK_CLAIM", "LINI", "", "teks", "-", "[terbuka-Non-Life] enum lintas-lini; Life = konstanta lini Life."],
 ["T_WORK_CLAIM", "PY_POSITION", "", "teks", "-", "posisi/status tangga"],
 ["T_WORK_CLAIM", "ACCEPT_STATUS", "", "teks", "-", ""],
 ["T_WORK_CLAIM", "SENDTO_ADMIN", "", "teks", "-", ""],
 ["T_WORK_CLAIM", "SENDTO_MEDICAL", "", "teks", "-", ""],
 ["T_WORK_CLAIM", "TYPE", "", "teks", "-", ""],
 ["T_WORK_CLAIM", "CASEID", "", "teks", "-", "PINDAHAN dari header klaim (D2)."],
 ["T_WORK_CLAIM", "CREATE_OP", "", "teks", "-", "pindahan dari header klaim"],
 ["T_WORK_CLAIM", "CREATE_OP_NAME", "", "teks", "-", "pindahan"],
 ["T_WORK_CLAIM", "TGL_UPDATE", "", "DATE", "-", "pindahan"],
 # T_GENERAL_CLAIM
 ["T_GENERAL_CLAIM", "ID", "PK+FK", "teks CLM-", "= T_WORK_CLAIM.ID (baris klaim)", "Shared PK 1:1 (D1). Sekaligus PK dan FK ke T_WORK_CLAIM."],
 ["T_GENERAL_CLAIM", "CASEID_POLICY", "", "teks", "CaseID polis (baca)", "KUNCI ke polis (D7/Q2). Polis dibaca hidup lewat CASEID ini."],
 ["T_GENERAL_CLAIM", "POLICY_NO", "", "teks", "PolicyDataLife.PremiumListSummary.PL_NUMBER", "Atribut tampilan (bekas PL_NUMBER). [terverifikasi] SetClaimXOL_Act."],
 ["T_GENERAL_CLAIM", "ENDORSMENT_NO", "", "teks", "versi polis via CASEID_POLICY", "Versi polis yang ditunjuk (D7)."],
 ["T_GENERAL_CLAIM", "CLAIM_NO", "", "teks", "-", "nomor klaim"],
 ["T_GENERAL_CLAIM", "RISLIPRNM", "", "teks", "-", "field PremiumListSummary"],
 ["T_GENERAL_CLAIM", "BUSINESS_NAME", "", "teks", "-", "field PremiumListSummary"],
 ["T_GENERAL_CLAIM", "CLAIM_RETRO", "", "desimal", "-", ""],
 # T_CLAIMLF_PREMIUMLIST_DETAIL
 ["T_CLAIMLF_PREMIUMLIST_DETAIL", "ID", "PK", "sequence", "-", "peserta"],
 ["T_CLAIMLF_PREMIUMLIST_DETAIL", "CLAIM_ID", "FK", "teks CLM-", "T_GENERAL_CLAIM.ID", "ON DELETE CASCADE (relasi 3)."],
 ["T_CLAIMLF_PREMIUMLIST_DETAIL", "IS_CHECK", "", "teks", "-", "salah satu dari 9 kolom audit"],
 ["T_CLAIMLF_PREMIUMLIST_DETAIL", "(+8 kolom audit)", "", "-", "-", "IS_CHECK + 3 tanggal per peserta + lainnya (spec AC 39-41)."],
 # T_CLAIMLF_ADJUSTMENT
 ["T_CLAIMLF_ADJUSTMENT", "ID", "PK", "sequence", "-", "baris adjustment"],
 ["T_CLAIMLF_ADJUSTMENT", "PREMIUM_LIST_DETAIL_ID", "FK", "sequence", "T_CLAIMLF_PREMIUMLIST_DETAIL.ID", "ON DELETE CASCADE (relasi 4)."],
 ["T_CLAIMLF_ADJUSTMENT", "KOMITE_ID", "FK", "teks KMT-", "T_WORK_CLAIM.ID (baris komite)", "Rujukan kasus komite (N:1, penunjuk). NULL bila belum kirim komite."],
 ["T_CLAIMLF_ADJUSTMENT", "CLAIM_AMOUNT", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT", "STS_REJECT", "", "angka", "-", "1 aksep / 2 tolak"],
 ["T_CLAIMLF_ADJUSTMENT", "ACCEPTEDNO", "", "teks", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT", "ACCEPTATION_DATE", "", "DATE", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT", "NAME_OF_BANK", "", "teks", "-", "kolom bank"],
 ["T_CLAIMLF_ADJUSTMENT", "ID_BANK", "", "teks", "-", "kolom bank"],
 ["T_CLAIMLF_ADJUSTMENT", "ACCOUNT_NO", "", "teks", "-", "kolom bank"],
 # T_CLAIMLF_ADJUSTMENT_SPREADING
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "ID", "PK", "sequence", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "ADJUSTMENT_ID", "FK", "sequence", "T_CLAIMLF_ADJUSTMENT.ID", "ON DELETE CASCADE (relasi 5)."],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "TREATY_TYPE_ID", "", "teks", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "TREATY_TYPE_NAME", "", "teks", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "TREATY_YEAR_LIFE", "", "teks", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "RETROCADED_SHARE", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "RATE", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "IDR", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "USD", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING", "CURRENCY", "", "teks", "-", ""],
 # SPREADING_RETRO
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "ID", "PK", "sequence", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "SPREADING_ID", "FK", "sequence", "T_CLAIMLF_ADJUSTMENT_SPREADING.ID", "ON DELETE CASCADE (relasi 6)."],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "REINSURER_NAME", "", "teks", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "PERCENT_SHARE", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "AMOUNT", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "RATE", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "PREMIUM_SPREADED_GROSS", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "PREMIUM_SPREADED_NET", "", "desimal", "= GROSS - Diskon - Komisi", "Rumus terbaru (D6/Q1)."],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "COMMISION", "", "desimal", "-", "sic (ejaan korpus)"],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "OVR_COMM", "", "desimal", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "TREATY_TYPE_ID", "", "teks", "-", ""],
 ["T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "TREATY_TYPE_NAME", "", "teks", "-", ""],
 # DOCUMENT_CLAIM
 ["DOCUMENT_CLAIM", "ID", "PK", "sequence", "-", "lintas-lini"],
 ["DOCUMENT_CLAIM", "PREMIUM_LIST_DETAIL_ID", "FK", "sequence", "T_CLAIMLF_PREMIUMLIST_DETAIL.ID", "ON DELETE di Go (D3, lintas-lini)."],
 # T_GENERAL_KOMITE
 ["T_GENERAL_KOMITE", "ID", "PK+FK", "teks KMT-", "= T_WORK_CLAIM.ID (baris komite)", "Shared PK 1:1 (D1). Sekaligus PK dan FK ke T_WORK_CLAIM."],
 ["T_GENERAL_KOMITE", "ADJUSTMENT_ID", "FK", "sequence", "T_CLAIMLF_ADJUSTMENT.ID", "Tutup lingkar (1:1, di Go). Komite ini milik adjustment mana."],
 ["T_GENERAL_KOMITE", "KOMITE_LOOP", "", "angka", "-", "jumlah tingkat = COUNT roster aktif"],
 ["T_GENERAL_KOMITE", "KOMITE_COUNT", "", "angka", "-", "tingkat berjalan"],
 ["T_GENERAL_KOMITE", "ACCEPT_STATUS", "", "angka", "-", "1 aksep / 2 tolak"],
 # T_KOMITE_KOMITELIST
 ["T_KOMITE_KOMITELIST", "ID", "PK", "sequence", "-", ""],
 ["T_KOMITE_KOMITELIST", "DATA_KOMITE_ID", "FK", "teks KMT-", "T_GENERAL_KOMITE.ID", "ON DELETE CASCADE (relasi 9)."],
 ["T_KOMITE_KOMITELIST", "KOMITE_URUT", "", "angka", "-", "jenjang/urutan tangga"],
 ["T_KOMITE_KOMITELIST", "KOMITE_ID", "", "teks", "Policy..OPERATOR_ID", "[terverifikasi] CreateKMTLife_Act. Anggota pemutus."],
 ["T_KOMITE_KOMITELIST", "ID_KOMITE", "", "teks", "JABATAN", ""],
 ["T_KOMITE_KOMITELIST", "KOMITE_EMAIL", "", "teks", "EMAIL", ""],
 ["T_KOMITE_KOMITELIST", "KOMITE_APROVAL", "", "angka", "-", "0 belum / 1 setuju / 2 tolak"],
 ["T_KOMITE_KOMITELIST", "KOMITE_COMMENT", "", "teks", "-", ""],
 ["T_KOMITE_KOMITELIST", "DATE_APPROVE", "", "DATE", "-", ""],
]
for r in kolom:
    ws2.append(r)
style_header(ws2, 6)
cur = None
for r in range(2, ws2.max_row + 1):
    tbl = ws2.cell(row=r, column=1).value
    key = ws2.cell(row=r, column=3).value
    # tebalkan nama tabel hanya di baris pertama tiap tabel
    if tbl != cur:
        ws2.cell(row=r, column=1).font = Font(bold=True)
        cur = tbl
    if key in ("PK", "PK+FK"):
        ws2.cell(row=r, column=2).fill = PK_FILL
        ws2.cell(row=r, column=3).fill = PK_FILL
    elif key == "FK":
        ws2.cell(row=r, column=2).fill = FK_FILL
        ws2.cell(row=r, column=3).fill = FK_FILL
    ws2.cell(row=r, column=5).alignment = WRAP
    ws2.cell(row=r, column=6).alignment = WRAP
apply_body(ws2, 6)
for col, w in zip("ABCDEF", [38, 26, 8, 14, 40, 46]):
    ws2.column_dimensions[col].width = w
ws2.freeze_panes = "A2"

# ================= SHEET 3: RELASI =================
ws3 = wb.create_sheet("Relasi")
ws3.append(["#", "Induk", "Anak", "Kunci Tamu / Cara", "Kardinalitas", "ON DELETE", "Dibuat oleh"])
relasi = [
 [1,  "T_WORK_CLAIM", "T_WORK_CLAIM (self)", "COVER_KEY", "1:N", "di Go", "tiket 00 Komite"],
 [2,  "T_WORK_CLAIM", "T_GENERAL_CLAIM", "SHARED PK (ID sama)", "1:1", "di Go", "tiket 14 (D1)"],
 [3,  "T_GENERAL_CLAIM", "T_CLAIMLF_PREMIUMLIST_DETAIL", "CLAIM_ID", "1:N", "CASCADE", "tiket 14"],
 [4,  "T_CLAIMLF_PREMIUMLIST_DETAIL", "T_CLAIMLF_ADJUSTMENT", "PREMIUM_LIST_DETAIL_ID", "1:N", "CASCADE", "tiket 14"],
 [5,  "T_CLAIMLF_ADJUSTMENT", "T_CLAIMLF_ADJUSTMENT_SPREADING", "ADJUSTMENT_ID", "1:N", "CASCADE", "tiket 14"],
 [6,  "T_CLAIMLF_ADJUSTMENT_SPREADING", "T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO", "SPREADING_ID", "1:N", "CASCADE", "tiket 14"],
 [7,  "T_CLAIMLF_PREMIUMLIST_DETAIL", "DOCUMENT_CLAIM", "PREMIUM_LIST_DETAIL_ID", "1:N", "di Go (lintas-lini, D3)", "tiket 14"],
 [8,  "T_WORK_CLAIM", "T_GENERAL_KOMITE", "SHARED PK (ID sama)", "1:1", "di Go", "tiket 00 Komite (D1)"],
 [9,  "T_GENERAL_KOMITE", "T_KOMITE_KOMITELIST", "DATA_KOMITE_ID", "1:N", "CASCADE", "tiket 00 Komite"],
 [10, "T_CLAIMLF_ADJUSTMENT", "T_GENERAL_KOMITE", "ADJUSTMENT_ID", "1:1", "di Go", "tiket 00 Komite"],
 [11, "T_CLAIMLF_ADJUSTMENT", "T_WORK_CLAIM", "KOMITE_ID (penunjuk)", "N:1", "penunjuk (tak cascade)", "tiket 14"],
]
for r in relasi:
    ws3.append(r)
style_header(ws3, 7)
for r in range(2, ws3.max_row + 1):
    ws3.cell(row=r, column=4).alignment = WRAP
    dele = ws3.cell(row=r, column=6).value
    if dele and "CASCADE" in dele:
        ws3.cell(row=r, column=6).fill = PatternFill("solid", fgColor="FCE4D6")  # oranye muda
apply_body(ws3, 7)
for col, w in zip("ABCDEFG", [4, 32, 36, 26, 12, 22, 20]):
    ws3.column_dimensions[col].width = w
ws3.freeze_panes = "A2"

# ================= SHEET 4: LEGENDA =================
ws4 = wb.create_sheet("Legenda")
ws4.append(["Keterangan warna & istilah"])
ws4["A1"].font = Font(bold=True, size=12)
legenda = [
 ("", ""),
 ("PK (kuning)", "Primary Key"),
 ("FK (biru)", "Foreign Key"),
 ("PK+FK (kuning)", "Shared Primary Key: ID sekaligus PK dan FK 1:1 ke induk"),
 ("CASCADE (oranye)", "ON DELETE CASCADE - anak ikut terhapus di database"),
 ("di Go", "Penghapusan diatur aplikasi Go, bukan cascade DB"),
 ("penunjuk", "FK rujukan saja (N:1), tidak memicu penghapusan"),
 ("", ""),
 ("Akar", "T_WORK_CLAIM - 1 baris per work object. Baris klaim (CLM-) & baris komite (KMT-)."),
 ("Shared PK", "T_GENERAL_CLAIM.ID = T_WORK_CLAIM.ID (klaim); T_GENERAL_KOMITE.ID = T_WORK_CLAIM.ID (komite)."),
 ("COVER_KEY", "Baris komite menunjuk baris klaim induknya lewat COVER_KEY."),
 ("", ""),
 ("Terbuka (DBA)", "Generator nomor CLM-/KMT-; apakah COVER_KEY/KOMITE_ID pakai REFERENCES."),
 ("Dibaca hidup", "BUSINESS_ID, TEAM_GROUP, 3 TANGGAL, POLICY_NO dibaca dari PolicyDataLife (polis), tidak disimpan di klaim."),
 ("Sumber", "Struktur final D1-D8 (2026-09-18): tiket 14 Claim Life + keputusan-struktur-komite.md."),
]
for r in legenda:
    ws4.append(r)
ws4.column_dimensions["A"].width = 20
ws4.column_dimensions["B"].width = 90
for r in range(2, ws4.max_row + 1):
    ws4.cell(row=r, column=1).font = Font(bold=True)
    ws4.cell(row=r, column=2).alignment = WRAP

wb.save(r"d:\XML\RNM_BRD\OUTPUT_HASIL_RNM\Skema-Claim-Komite-Life.xlsx")
print("SELESAI: Skema-Claim-Komite-Life.xlsx (4 sheet)")
