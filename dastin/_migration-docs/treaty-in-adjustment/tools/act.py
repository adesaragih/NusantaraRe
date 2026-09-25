# -*- coding: utf-8 -*-
"""Peringkas Activity Pega. SELALU cetak penanda MATI (pyStepsBlockName == '//')."""
import os, sys, xml.etree.ElementTree as ET
BASE=r"D:\XML_NURE"
def find(name, root=os.path.join(BASE,"Treaty In Adjustment")):
    if os.path.exists(name): return name
    for dp,dn,fn in os.walk(root):
        for f in fn:
            if f.lower() in (name.lower(), name.lower()+".xml"): return os.path.join(dp,f)
def txt(n,t):
    e=n.find(t); return (e.text or '').strip() if e is not None and e.text else ''
def dump(p, verbose=True):
    r=ET.parse(p).getroot()
    print("### %s | %s | %s %s | avail=%s"%(txt(r,'pyRuleName'),txt(r,'pyClassName'),txt(r,'pyRuleSet'),txt(r,'pyRuleSetVersion'),txt(r,'pyRuleAvailable')))
    m=txt(r,'pyNote') or txt(r,'pyMemo')
    if m: print("memo: %s"%m[:250])
    pr=r.find('pyParameters')
    if pr is not None:
        ps=[txt(x,'pyParametersParamName') for x in pr.findall('rowdata')]
        ps=[x for x in ps if x]
        if ps: print("param: %s"%", ".join(ps))
    def walk(cont,pref):
        i=0
        for row in cont.findall('rowdata'):
            i+=1; lbl=pref+str(i)
            blk=txt(row,'pyStepsBlockName'); dead = blk=='//'
            meth=txt(row,'pyStepsActivityName'); page=txt(row,'pyStepsObjectName')
            desc=txt(row,'pyStepsDescription')
            pre=txt(row,'pyStepsPreCondParamsWhen'); prelbl=txt(row,'pyStepsPreCondLabel')
            precond=txt(row,'pyStepsPreCondition')
            loop=txt(row,'pyStepsLoopType') or txt(row,'pyStepsIterateType')
            print("%-8s %-4s %-34s %-26s %s"%(lbl,"MATI" if dead else "",meth[:34],page[:26],desc[:80]))
            if pre: print("            when: %s"%pre[:170])
            if prelbl and prelbl!=pre: print("            when-label: %s"%prelbl[:120])
            if blk and not dead: print("            block: %s"%blk)
            if loop: print("            loop: %s"%loop[:120])
            tw=txt(row,'pyStepsTransParamsWhen')
            if tw: print("            trans-when: %s -> T:%s F:%s"%(tw[:90],txt(row,'pyStepsTransParamsWhenTrue'),txt(row,'pyStepsTransParamsWhenFalse')))
            if verbose:
                cp=row.find('pyStepsCallParams')
                if cp is not None:
                    for ch in cp:
                        v=(ch.text or '').strip()
                        if ch.tag not in ('pxObjClass','pyTempPlaceHolder') and v:
                            print("            call-param %s = %s"%(ch.tag,v[:140]))
                for pa in row.findall('pyParamArray'):
                    for rr in pa.findall('rowdata'):
                        nm=txt(rr,'PropertiesName'); vl=txt(rr,'PropertiesValue')
                        if nm: print("            set   %-62s = %s"%(nm[:62],vl[:150]))
                        for extra in ('Property','Value','Param','PageName','PropertyName'):
                            pass
                pp=row.find('pySteps')
            sub=row.find('pySteps')
            if sub is not None: walk(sub,lbl+".")
    c=r.find('pySteps')
    if c is None: print("(tidak ada pySteps)")
    else: walk(c,"")
    print()
for n in sys.argv[1:]:
    p=find(n)
    if not p: print("TIDAK KETEMU",n); continue
    dump(p)
