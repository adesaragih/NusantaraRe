# -*- coding: utf-8 -*-
"""Precondition + blok mati per langkah. Precondition tersarang di
   pyStepsPreCondParams/rowdata/pyStepsPreCondParamsWhen -> harus dicari sebagai keturunan,
   TANPA masuk ke pySteps bersarang.  3 = lewati langkah, 2 = lanjut."""
import os, sys, xml.etree.ElementTree as ET
BASE=r"D:\XML_NURE\Treaty In Adjustment"
def find(n):
    if os.path.exists(n): return n
    for dp,dn,fn in os.walk(BASE):
        for f in fn:
            if f.lower() in (n.lower(), n.lower()+'.xml'): return os.path.join(dp,f)
def own(row, tag):
    """cari tag di keturunan row tanpa melintasi pySteps"""
    out=[]
    def rec(n):
        for c in n:
            if c.tag=='pySteps': continue
            if c.tag==tag and (c.text or '').strip(): out.append(c.text.strip())
            rec(c)
    rec(row); return out
ACT={'2':'lanjut','3':'LEWATI','4':'lompat','5':'keluar-aktivitas','6':'keluar-iterasi'}
for name in sys.argv[1:]:
    p=find(name)
    if not p: print("TIDAK KETEMU",name); continue
    r=ET.parse(p).getroot()
    print("### %s  (%s %s)"%((r.findtext('pyRuleName') or '').strip(),(r.findtext('pyRuleSet') or '').strip(),(r.findtext('pyRuleSetVersion') or '').strip()))
    def walk(c,pref):
        i=0
        for row in c.findall('rowdata'):
            i+=1; lbl=pref+str(i)
            blk=(row.findtext('pyStepsBlockName') or '').strip()
            meth=(row.findtext('pyStepsActivityName') or '').strip()
            desc=(row.findtext('pyStepsDescription') or '').strip()
            w=own(row,'pyStepsPreCondParamsWhen'); wt=own(row,'pyStepsPreCondParamsWhenTrue'); wf=own(row,'pyStepsPreCondParamsWhenFalse')
            wtp=own(row,'pyStepsPreCondParamsWhenTruePrms'); wfp=own(row,'pyStepsPreCondParamsWhenFalsePrms')
            en=(row.findtext('pyStepsPreCondition') or '').strip()
            tag="MATI" if blk=='//' else ""
            cond=""
            if w:
                t=ACT.get(wt[0] if wt else '2','?'); f=ACT.get(wf[0] if wf else '2','?')
                if wt and wt[0]=='1': t='lompat->'+(wtp[0] if wtp else '?')
                if wf and wf[0]=='1': f='lompat->'+(wfp[0] if wfp else '?')
                if en!='true': cond="  [PRECOND NONAKTIF] sisa teks: %s"%w[0][:70]
                else: cond="  WHEN[%s] benar->%s salah->%s"%(w[0][:75],t,f)
            print("%-8s %-4s %-32s %-52s%s"%(lbl,tag,meth[:32],desc[:52],cond))
            s=row.find('pySteps')
            if s is not None: walk(s,lbl+'.')
    walk(r.find('pySteps'),'')
    print()
