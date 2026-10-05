// Form produk - `InboxProductName` mode form (`OutputParam.DATASHOW = 1`), urutan wadah XML:
//   b2573 sisi umum + `LIEN CLAUSE` + `DOCUMENT CLAIM` + sisi inward   (PARITAS §3.1, §3.2, §5)
//   b30995 `PLAN LIST`, `FINANCIAL UNDERWRITING`, `UNDERWRITING LIMIT`, checkbox `On Retention`
//   b57867 tombol bawah `Close` `Save` `Edit` `Copy` `Generate`            (PARITAS §3.3)
//   b61044 grid komentar (baca-saja)                                       (PARITAS §5)
//   b64133 lampiran                                                        (PARITAS §6)
//
// Mode lihat (`ProductName.IsView == 'true'`, sesudah `View` b74798 → `SetProductName` 8 b2147): medan ber-`ro`
// baca-saja, dropdown master (pengganti tombol `Choose*`, wadah `IsView!='true'`) dan `Save` tidak dapat dipakai,
// `Edit` tampil. Checkbox `On Retention` dan tombol `Add`/`Delete` grid tidak ber-`ro` di XML - tetapi keputusan work
// owner 03-10-2026 ("jika view tidak tambah/edit/delete, saat klik edit baru bisa"): di mode lihat SEMUA aksi ubah
// tersembunyi (Add/Delete/Copy row grid, Copy, Add attachment/Retry/Delete lampiran) dan On Retention mati; yang tetap:
// Close, Edit, Generate, View Rate, dan aksi baca lampiran.
// ⛔ Audit 02-10-2026: enam medan pemilih master SELALU baca-saja di XML (`pyReadOnly` true, `pyEditOptions`
// Read-only, `pyReadOnlyCondition` KOSONG - beda dengan `Product Name` b3585 yang bersyarat `IsView`): Ceding
// b4040, SOB b4428, R/I Risk Name b7362, Cause Of Loss b10626, Policy Holder b17062, Currency b28105. Nilainya
// hanya dari daftar master → `set*_DT`, tidak diketik. Begitu pula sel `Bussines` (`.Name` b33504) dan `Benefit`
// b33658 grid `PLAN LIST`: diisi pilihan `Plan Name` (dropdown, 03-10-2026), tidak diketik.
// Keputusan work owner 02-10-2026 ("perubahan pada tampilan untuk semua Choose ubah jadi dropdown saja"): ketujuh
// tombol `Choose*` + popup FlowAction-nya diganti `DropdownMaster` (termasuk `Choose R/I Rate` baris `PLAN LIST`).
// Grid `OUTWARD` (wadah `1==2`) tidak dirender; isinya ditulis server (`hitungOutward`).

import { useState, type ReactNode } from 'react'

import { Gagal, Kosong } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import {
  cariPlan,
  simpanProduk,
  unduhGenerate,
  type BarisDokumen,
  type BarisFinUW,
  type BarisLien,
  type BarisPlan,
  type BarisUWLimit,
  type JenisPlan,
  type Produk,
  type ProdukInward,
  type ProdukUmum,
} from '../api'
import {
  PILIHAN_DOKUMEN_KLAIM,
  PILIHAN_PEMBAYARAN,
  hitungMaxSumReasured,
  potongPilihan,
  namaTreaty,
  salinBaris,
  salinProduk,
  tampilPremiumFactor,
  tampilViewRate,
  waktuPega,
} from '../bentuk'
import {
  DOKUMEN_MPNL,
  FINUW_MPNL,
  INWARD_MPNL,
  KOMENTAR_MPNL,
  LAIN_MPNL,
  LIEN_MPNL,
  PEMILIH_MPNL,
  PESAN_MPNL,
  PLAN_MPNL,
  SARAN_PLAN_MPNL,
  TOMBOL_MPNL,
  UMUM_MPNL,
  UWLIMIT_MPNL,
} from '../labels'
import { NAMA_MPNL } from '../menu'
import { DialogEdit, DialogSimpan } from './Dialog'
import DropdownCari from './DropdownCari'
import DropdownMaster from './DropdownMaster'
import Medan, { PilihanMedan, TANDA_KOSONG } from './Medan'
import ModalRate from './ModalRate'
import PanelLampiran from './PanelLampiran'

/** Daftar `Plan Name` (RD `BrowseProductTypeLife_RD`, dicari pada CoverName dan Business), dipotong BATAS_DROPDOWN. */
const cariJenisPlan = async (kata: string) => potongPilihan((await cariPlan(kata)).daftar)

/** Sel grid bersarang yang dapat disunting (`ro = ProductName.IsView=='true'`). */
function SelIsi({ nilai, onUbah, readOnly, label }: { nilai: string; onUbah: (v: string) => void; readOnly: boolean; label: string }) {
  // Mode lihat: teks seperti baris `Medan` (foto layar Pega work owner 02-10-2026), bukan isian baca-saja.
  if (readOnly) return <span className={nilai === '' ? 'mpnl-nilai mpnl-nilai--kosong' : 'mpnl-nilai'}>{nilai === '' ? TANDA_KOSONG : nilai}</span>
  return (
    <input
      className="field__input mpnl-sel-isi"
      aria-label={label}
      value={nilai}
      readOnly={readOnly}
      onChange={(e) => {
        onUbah(e.target.value)
      }}
    />
  )
}

function ganti<T>(daftar: readonly T[], i: number, baru: Partial<T>): T[] {
  return daftar.map((b, j) => (j === i ? { ...b, ...baru } : b))
}

function buang<T>(daftar: readonly T[], i: number): T[] {
  return daftar.filter((_, j) => j !== i)
}

export default function FormProduk({
  awal,
  lihatAwal,
  onTutup,
  onTersimpan,
}: {
  awal: Produk
  /** `IsView` awal: `View` = true, `Add` = false. */
  lihatAwal: boolean
  onTutup: () => void
  onTersimpan: () => void
}) {
  const [p, setP] = useState<Produk>(awal)
  // Menu View only (M_LOGIN_GO_MENU.HAK, 04-10-2026): form selalu mode lihat dan tanpa tombol Edit.
  const bolehUbah = useBolehUbah(NAMA_MPNL)
  const [lihat, setLihat] = useState(lihatAwal || !bolehUbah)
  const [pesan, setPesan] = useState<string | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [dialog, setDialog] = useState<'simpan' | 'edit' | null>(null)
  const [rate, setRate] = useState<string | null>(null)

  const u = p.umum
  const w = p.inward

  function ubahUmum(baru: Partial<ProdukUmum>): void {
    setP((x) => ({ ...x, umum: { ...x.umum, ...baru } }))
  }
  function ubahInward(baru: Partial<ProdukInward>): void {
    setP((x) => ({ ...x, inward: { ...x.inward, ...baru } }))
  }
  const medanUmum = (k: keyof ProdukUmum) => (v: string) => {
    ubahUmum({ [k]: v })
  }
  const medanInward = (k: keyof ProdukInward) => (v: string) => {
    ubahInward({ [k]: v })
  }

  /** `CountMaxSumReasured_Act` - onChange `Ceding's Limit` b22725 dan `Max Sum Insured` b24024. */
  function ubahHitung(k: 'cedingLimit' | 'maxSumInsured') {
    return (v: string) => {
      setP((x) => {
        const inward = { ...x.inward, [k]: v }
        const hasil = hitungMaxSumReasured(inward.maxSumInsured, inward.cedingLimit)
        return { ...x, inward: hasil === null ? inward : { ...inward, maxSumReasured: hasil } }
      })
    }
  }

  async function simpan(): Promise<void> {
    if (menyimpan) return
    setMenyimpan(true)
    setGalat(null)
    try {
      await simpanProduk(p)
      setDialog(null)
      // `SaveProductName_Act` 11 b2209: `DATASHOW = 0` - form tertutup, grid dimuat ulang.
      onTersimpan()
    } catch (e) {
      setDialog(null)
      setGalat(e)
    } finally {
      setMenyimpan(false)
    }
  }

  async function generate(): Promise<void> {
    setGalat(null)
    try {
      await unduhGenerate(p)
    } catch (e) {
      setGalat(e)
    }
  }

  return (
    <section className="panel mpnl-form">
      {pesan !== null && <p className="mpnl-pesan">{pesan}</p>}
      {galat !== null && <Gagal galat={galat} />}

      {/* ---- wadah b2573 (foto layar Pega work owner 02-10-2026): kartu TREATY NAME kiri, INWARD kanan ---- */}
      <div className="mpnl-tata">
        <section className="mpnl-kartu">
          <h3 className="mpnl-kartu__judul">{UMUM_MPNL.judul}</h3>
          {/* `pyRequired` true b3585 (Product Name), b25362, b26054, b26268: tanda wajib `*` (keputusan work owner 02-10-2026). */}
          <Medan
            label={UMUM_MPNL.productName}
            required
            nilai={u.productName}
            lihat={lihat}
            onUbah={(v) => {
              // onChange `SetTreatyName_Act` b3686.
              ubahUmum({ productName: v, inwardName: namaTreaty(v, w.policyHolderName) })
            }}
          />
          <Medan label={UMUM_MPNL.productCode} kode nilai={p.id} lihat={lihat} />
          <Medan label={UMUM_MPNL.ceding} nilai={u.ceding} lihat={lihat}>
            <DropdownMaster
              labelAria={UMUM_MPNL.ceding}
              jenis="ceding"
              nilai={u.ceding}
              lihat={lihat}
              onPilih={(v) => {
                ubahUmum({ ceding: v.nama, cedingId: v.id })
              }}
            />
          </Medan>
          <Medan label={UMUM_MPNL.sob} nilai={u.sobName} lihat={lihat}>
            <DropdownMaster
              labelAria={UMUM_MPNL.sob}
              jenis="sob"
              nilai={u.sobName}
              lihat={lihat}
              onPilih={(v) => {
                ubahUmum({ sobName: v.nama, sobId: v.id })
              }}
            />
          </Medan>
          <Medan label={UMUM_MPNL.deduction} jenis="angka" nilai={u.riComm} lihat={lihat} onUbah={medanUmum('riComm')} />
          <Medan label={UMUM_MPNL.riRisk} nilai={u.riRisk} lihat={lihat}>
            <DropdownMaster
              labelAria={UMUM_MPNL.riRisk}
              jenis="ri-risk"
              nilai={u.riRisk}
              lihat={lihat}
              onPilih={(v) => {
                ubahUmum({ riRisk: v.nama, riRiskId: v.id })
              }}
            />
          </Medan>
          <Medan label={UMUM_MPNL.treatyName} nilai={u.inwardName} lihat={lihat} onUbah={medanUmum('inwardName')} />
          <Medan label={UMUM_MPNL.treatyNumber} nilai={u.treatyNumber} lihat={lihat} onUbah={medanUmum('treatyNumber')} />
          {/* `Cause Of Loss` b10626: SELALU baca-saja - diisi dari daftar master → `setCauseofLoss_DT`. */}
          <Medan label={UMUM_MPNL.causeOfLoss} nilai={u.cause} lihat={lihat}>
            <DropdownMaster
              labelAria={UMUM_MPNL.causeOfLoss}
              jenis="penyebab"
              nilai={u.cause}
              lihat={lihat}
              onPilih={(v) => {
                ubahUmum({ cause: v.nama, causeId: v.id })
              }}
            />
          </Medan>

          {/* `LIEN CLAUSE` b12201 - ikon grid bawaan b12373 (vis `IsView!='true'`). */}
          <h4 className="mpnl-judul-grid">{LIEN_MPNL.judul}</h4>
          <GridSederhana<BarisLien>
            baris={p.lienClause}
            lihat={lihat}
            kolom={[
              [LIEN_MPNL.usia, 'usia'],
              [LIEN_MPNL.manfaat, 'manfaat'],
            ]}
            kosong={{ usia: '', manfaat: '' }}
            onUbah={(d) => {
              setP((x) => ({ ...x, lienClause: d }))
            }}
          />

          {/* `DOCUMENT CLAIM` b14601 - `Document List` dropdown: PromptList properti `.Document` (XML dikirim work owner
          03-10-2026, `PILIHAN_DOKUMEN_KLAIM`); nilai lama di luar daftar tetap tampil, tidak dibuang. */}
          <h4 className="mpnl-judul-grid">{DOKUMEN_MPNL.judul}</h4>
          <GridSederhana<BarisDokumen>
            baris={p.documentClaim}
            lihat={lihat}
            kolom={[[DOKUMEN_MPNL.documentList, 'document', PILIHAN_DOKUMEN_KLAIM]]}
            kosong={{ document: '' }}
            onUbah={(d) => {
              setP((x) => ({ ...x, documentClaim: d }))
            }}
          />
        </section>

        {/* ---- sisi inward (halaman `ProductNameInward`, wadah b16621) ---- */}
        <section className="mpnl-kartu mpnl-kartu--lebar">
          <h3 className="mpnl-kartu__judul">{INWARD_MPNL.judul}</h3>
          {/* Pilihan Policy Holder → `setPolicyHolder_DT` b2416 - TANPA `SetTreatyName_Act` (onChange b17168 milik
              autocomplete yang selalu baca-saja, jadi tidak pernah terpicu). */}
          <Medan label={INWARD_MPNL.policyHolder} nilai={w.policyHolderName} lihat={lihat}>
            <DropdownMaster
              labelAria={INWARD_MPNL.policyHolder}
              jenis="pemegang-polis"
              nilai={w.policyHolderName}
              lihat={lihat}
              onPilih={(v) => {
                ubahInward({ policyHolderName: v.nama, policyHolder: v.id })
              }}
            />
          </Medan>
          <Medan label={INWARD_MPNL.insured} nilai={w.insured} lihat={lihat} onUbah={medanInward('insured')} />
          {/* Pasangan sebaris b18777 / b20095: `Addendum No.` + `Addendum`, `Amandement No.` + `Amandement`. */}
          <div className="mpnl-dua-kolom">
            <Medan label={INWARD_MPNL.addendumNo} jenis="angka" nilai={w.addendumNo} lihat={lihat} onUbah={medanInward('addendumNo')} />
            <Medan label={INWARD_MPNL.addendum} nilai={w.addendumWord} lihat={lihat} onUbah={medanInward('addendumWord')} />
            <Medan label={INWARD_MPNL.amandementNo} jenis="angka" nilai={w.amandementNo} lihat={lihat} onUbah={medanInward('amandementNo')} />
            <Medan label={INWARD_MPNL.amandement} nilai={w.amandementSchd} lihat={lihat} onUbah={medanInward('amandementSchd')} />
          </div>
          <div className="mpnl-dua-kolom mpnl-dua-kolom--pisah">
            <div className="mpnl-kolom">
              {/* b21723 kolom kiri */}
              <Medan label={INWARD_MPNL.maxExpiredClaim} jenis="angka" nilai={w.maxExpiredClaim} lihat={lihat} onUbah={medanInward('maxExpiredClaim')} />
              <Medan label={INWARD_MPNL.begin} jenis="tanggal" nilai={w.begin} lihat={lihat} onUbah={medanInward('begin')} />
              <Medan label={INWARD_MPNL.stnc} jenis="tanggal" nilai={w.stnc} lihat={lihat} onUbah={medanInward('stnc')} />
              <Medan label={INWARD_MPNL.cedingRetention} jenis="angka" nilai={w.cedingRetentionNum} lihat={lihat} onUbah={medanInward('cedingRetentionNum')} />
              <Medan label={INWARD_MPNL.cedingLimit} jenis="angka" nilai={w.cedingLimit} lihat={lihat} onUbah={ubahHitung('cedingLimit')} />
              <Medan label={INWARD_MPNL.brokerage} jenis="angka" nilai={w.brokerage} lihat={lihat} onUbah={medanInward('brokerage')} />
              <Medan label={INWARD_MPNL.minAge} jenis="angka" nilai={w.minAge} lihat={lihat} onUbah={medanInward('minAge')} />
              <Medan label={INWARD_MPNL.maxAge} jenis="angka" nilai={w.maxAge} lihat={lihat} onUbah={medanInward('maxAge')} />
              <Medan label={INWARD_MPNL.expiryAge} jenis="angka" nilai={w.expiryAge} lihat={lihat} onUbah={medanInward('expiryAge')} />
              <Medan label={INWARD_MPNL.extraPremi} jenis="angka" nilai={w.extraPremi} lihat={lihat} onUbah={medanInward('extraPremi')} />
              <Medan label={INWARD_MPNL.minSumInsured} jenis="angka" nilai={w.minSumInsured} lihat={lihat} onUbah={medanInward('minSumInsured')} />
              <Medan label={INWARD_MPNL.maxSumInsured} jenis="angka" nilai={w.maxSumInsured} lihat={lihat} onUbah={ubahHitung('maxSumInsured')} />
              <Medan label={INWARD_MPNL.maxSumReasured} jenis="angka" nilai={w.maxSumReasured} lihat={lihat} onUbah={medanInward('maxSumReasured')} />
              <Medan
                label={INWARD_MPNL.rnmShare}
                jenis="angka"
                nilai={w.rnmShare}
                lihat={lihat}
                onUbah={medanInward('rnmShare')}
                catatan={<span className="mpnl-catatan-medan">{INWARD_MPNL.ofSumReasured}</span>}
              />
              <Medan label={INWARD_MPNL.rnmLimit} jenis="angka" nilai={w.rnmLimitNum} lihat={lihat} onUbah={medanInward('rnmLimitNum')} />
              {tampilPremiumFactor(w.payment) && (
                <Medan label={INWARD_MPNL.premiumFactor} jenis="angka" required nilai={w.premiumFactor} lihat={lihat} onUbah={medanInward('premiumFactor')} />
              )}
              {/* `ro = ProductName.IsView=='true'` b25611. */}
              <Medan
                label={INWARD_MPNL.payment}
                nilai={w.payment}
                tampil={PILIHAN_PEMBAYARAN.find((o) => o.value === w.payment)?.label ?? w.payment}
                lihat={lihat}
              >
                <PilihanMedan labelAria={INWARD_MPNL.payment} value={w.payment} opsi={PILIHAN_PEMBAYARAN} onChange={medanInward('payment')} />
              </Medan>
              <Medan label={INWARD_MPNL.subjectTo} nilai={w.subjectTo} lihat={lihat}>
                <textarea
                  className="field__input mpnl-area"
                  aria-label={INWARD_MPNL.subjectTo}
                  rows={3}
                  value={w.subjectTo}
                  onChange={(e) => {
                    ubahInward({ subjectTo: e.target.value })
                  }}
                />
              </Medan>
              <Medan label={INWARD_MPNL.annuityInterest} jenis="angka" required nilai={w.annuityInterest} lihat={lihat} onUbah={medanInward('annuityInterest')} />
              <Medan
                label={INWARD_MPNL.premiumRefundFactor}
                jenis="angka"
                required
                nilai={w.premiumRefundFactor}
                lihat={lihat}
                onUbah={medanInward('premiumRefundFactor')}
              />
            </div>
            <div className="mpnl-kolom">
              {/* b27004 kolom kanan - `%` b27585 dan `X + n` b27773 ber-vis `OTHER 1=2` (mati, PARITAS §3.2). */}
              <Medan label={INWARD_MPNL.maxDataReceive} jenis="angka" nilai={w.maxDataReceive} lihat={lihat} onUbah={medanInward('maxDataReceive')} />
              <Medan label={INWARD_MPNL.mature} jenis="tanggal" nilai={w.mature} lihat={lihat} onUbah={medanInward('mature')} />
              {/* `Birthday` b27960 - radio `associated` (pilihan tak ikut ekspor, OQ-MPNL-05): isian teks. */}
              <Medan label={INWARD_MPNL.birthday} nilai={w.birthday} lihat={lihat} onUbah={medanInward('birthday')} />
              <Medan label={INWARD_MPNL.currency} nilai={w.currency} lihat={lihat}>
                <DropdownMaster
                  labelAria={INWARD_MPNL.currency}
                  jenis="mata-uang"
                  nilai={w.currency}
                  lihat={lihat}
                  onPilih={(v) => {
                    ubahInward({ currency: v.nama, currencyId: v.id })
                  }}
                />
              </Medan>
              <Medan label={INWARD_MPNL.extraMortality} jenis="angka" nilai={w.extraMortality} lihat={lihat} onUbah={medanInward('extraMortality')} />
              <Medan label={INWARD_MPNL.maxContract} jenis="angka" nilai={w.maxContract} lihat={lihat} onUbah={medanInward('maxContract')} />
              <Medan label={INWARD_MPNL.proportionalTable} jenis="angka" nilai={w.proportionalTable} lihat={lihat} onUbah={medanInward('proportionalTable')} />
            </div>
          </div>
        </section>
      </div>

      {/* ---- wadah b30995: PLAN LIST, FINANCIAL UNDERWRITING, UNDERWRITING LIMIT, On Retention - satu kartu ---- */}
      <section className="mpnl-kartu mpnl-bagian">
      <h4 className="mpnl-judul-grid">{PLAN_MPNL.judul}</h4>
      {!lihat && (
        <div className="aksi-baris">
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setP((x) => ({ ...x, planList: [...x.planList, { plan: '', planId: '', name: '', benefit: '', riRate: '', riRateId: '' }] }))
            }}
          >
            {PLAN_MPNL.add}
          </button>
        </div>
      )}
      {p.planList.length === 0 ? (
        <Kosong pesan={LAIN_MPNL.kosong} />
      ) : (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{PLAN_MPNL.planName}</th>
              <th>{PLAN_MPNL.bussines}</th>
              <th>{PLAN_MPNL.benefit}</th>
              <th>{PLAN_MPNL.riRate}</th>
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {p.planList.map((b, i) => (
              <tr key={i} className="inbox__baris">
                <td>
                  {/* `Plan Name` b33121 sebagai dropdown (keputusan work owner 03-10-2026 "tolong ubah jadi model dropdown";
                      XML: pxAutoComplete berisian bebas b33137). Kolom = `pyAdditionalFields` ber-`pyShow` true. */}
                  <DropdownCari<JenisPlan>
                    labelAria={PLAN_MPNL.planName}
                    nilai={b.plan}
                    lihat={lihat}
                    lebar
                    cari={cariJenisPlan}
                    kunci={(t) => t.id}
                    kolom={{
                      judul: [SARAN_PLAN_MPNL.kolomId, SARAN_PLAN_MPNL.kolomCoverName, SARAN_PLAN_MPNL.kolomBusiness, SARAN_PLAN_MPNL.kolomBenefit],
                      isi: (t) => [t.id, t.coverName, t.business, t.benefit],
                      lebar: '6.5rem minmax(9rem, 2fr) minmax(9rem, 2fr) minmax(5rem, 1fr)',
                    }}
                    terpilih={(t) => t.id === b.planId}
                    onPilih={(t) => {
                      // b33198: `.CoverName` → `.Plan`, `.ID` → `.PlanID`, `.Business` → `.Name`, `.Benefit` → `.Benefit`.
                      setP((x) => ({
                        ...x,
                        planList: ganti<BarisPlan>(x.planList, i, { plan: t.coverName, planId: t.id, name: t.business, benefit: t.benefit }),
                      }))
                    }}
                  />
                </td>
                {/* `Bussines` (`.Name` b33504) dan `Benefit` b33658 SELALU baca-saja: diisi pilihan `Plan Name`. */}
                <td>{b.name}</td>
                <td>{b.benefit}</td>
                <td>
                  {/* `Choose R/I Rate` → `SetRIRate` 1 b249: `.RIRATEID ← id`, `.RIRATE ← usedby`. */}
                  <DropdownMaster
                    labelAria={PLAN_MPNL.riRate}
                    jenis="ri-rate"
                    kolomNama={PEMILIH_MPNL.kolomRiRateName}
                    nilai={b.riRate}
                    lihat={lihat}
                    onPilih={(v) => {
                      setP((x) => ({ ...x, planList: ganti<BarisPlan>(x.planList, i, { riRate: v.nama, riRateId: v.id }) }))
                    }}
                  />
                </td>
                <td className="table__actions">
                  {tampilViewRate(b.riRate) && (
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setRate(b.riRateId)
                      }}
                    >
                      {PLAN_MPNL.viewRate}
                    </button>
                  )}{' '}
                  {!lihat && (
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setP((x) => ({ ...x, planList: buang(x.planList, i) }))
                      }}
                    >
                      {PLAN_MPNL.delete}
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* `FINANCIAL UNDERWRITING` b37148. */}
      <h4 className="mpnl-judul-grid">{FINUW_MPNL.judul}</h4>
      <GridBerangka<BarisFinUW>
        baris={p.financialUnderwriting}
        lihat={lihat}
        kolom={[
          [FINUW_MPNL.minInsured, 'minInsured'],
          [FINUW_MPNL.maxInsured, 'maxInsured'],
          [FINUW_MPNL.employee, 'employee'],
          [FINUW_MPNL.nonEmployee, 'nonEmployee'],
        ]}
        tambah={
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setP((x) => ({ ...x, financialUnderwriting: [...x.financialUnderwriting, { minInsured: '', maxInsured: '', employee: '', nonEmployee: '' }] }))
            }}
          >
            {FINUW_MPNL.add}
          </button>
        }
        hapus={(i) => (
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setP((x) => ({ ...x, financialUnderwriting: buang(x.financialUnderwriting, i) }))
            }}
          >
            {FINUW_MPNL.delete}
          </button>
        )}
        onUbah={(d) => {
          setP((x) => ({ ...x, financialUnderwriting: d }))
        }}
      />

      {/* `UNDERWRITING LIMIT` b42075 - `Medical` teks bebas (R16). */}
      <h4 className="mpnl-judul-grid">{UWLIMIT_MPNL.judul}</h4>
      <GridBerangka<BarisUWLimit>
        baris={p.underwritingLimit}
        lihat={lihat}
        kolom={[
          [UWLIMIT_MPNL.minInsured, 'minInsured'],
          [UWLIMIT_MPNL.maxInsured, 'maxInsured'],
          [UWLIMIT_MPNL.minAge, 'minAge'],
          [UWLIMIT_MPNL.maxAge, 'maxAge'],
          [UWLIMIT_MPNL.medical, 'medical'],
          [UWLIMIT_MPNL.description, 'description'],
        ]}
        tambah={
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setP((x) => ({
                ...x,
                underwritingLimit: [...x.underwritingLimit, { minInsured: '', maxInsured: '', minAge: '', maxAge: '', medical: '', description: '' }],
              }))
            }}
          >
            {UWLIMIT_MPNL.add}
          </button>
        }
        hapus={(i) => (
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setP((x) => ({ ...x, underwritingLimit: buang(x.underwritingLimit, i) }))
            }}
          >
            {UWLIMIT_MPNL.delete}
          </button>
        )}
        onUbah={(d) => {
          setP((x) => ({ ...x, underwritingLimit: d }))
        }}
      />

      {/* Checkbox `On Retention` b47312 - onChange `GetReinsTypeOR_Life` b47488 (dihitung server saat simpan). */}
      <label className="mpnl-centang">
        <input
          type="checkbox"
          checked={u.isOrs}
          disabled={lihat}
          onChange={(e) => {
            const v = e.target.checked
            setP((x) => ({ ...x, umum: { ...x.umum, isOrs: v }, hitungOutward: true }))
          }}
        />
        {UMUM_MPNL.onRetention}
      </label>
      </section>

      {/* ---- wadah b57867: tombol bawah ---- */}
      <div className="aksi-baris mpnl-bagian">
        <button type="button" className="btn btn--ghost" onClick={onTutup}>
          {TOMBOL_MPNL.close}
        </button>{' '}
        {!lihat && (
          <button
            type="button"
            className="btn btn--primary"
            onClick={() => {
              setDialog('simpan')
            }}
          >
            {TOMBOL_MPNL.save}
          </button>
        )}{' '}
        {lihat && bolehUbah && (
          <button
            type="button"
            className="btn btn--primary"
            onClick={() => {
              setDialog('edit')
            }}
          >
            {TOMBOL_MPNL.edit}
          </button>
        )}{' '}
        {!lihat && (
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              // DataTransform `CopyProduct` b59965; pesan 4 b236 di wadah b1269 (`STSSAVE==99`).
              setP((x) => salinProduk(x))
              setPesan(PESAN_MPNL.copy)
              setGalat(null)
            }}
          >
            {TOMBOL_MPNL.copy}
          </button>
        )}{' '}
        <button type="button" className="btn btn--ghost" onClick={() => void generate()}>
          {TOMBOL_MPNL.generate}
        </button>
      </div>

      {/* ---- wadah b61044: komentar (baca-saja) ---- */}
      {p.commentList.length > 0 && (
        <table className="inbox__tabel mpnl-bagian">
          <thead>
            <tr>
              <th>{KOMENTAR_MPNL.date}</th>
              <th>{KOMENTAR_MPNL.pic}</th>
              <th>{KOMENTAR_MPNL.comment}</th>
            </tr>
          </thead>
          <tbody>
            {p.commentList.map((k, i) => (
              <tr key={i} className="inbox__baris">
                <td>{waktuPega(k.date)}</td>
                <td>{k.operatorName}</td>
                <td>{k.suggest}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* ---- wadah b64133: lampiran (produk tersimpan) ---- */}
      {p.id !== '' && <PanelLampiran produkId={p.id} lihat={lihat} />}

      {dialog === 'simpan' && (
        <DialogSimpan
          komentar={u.comment}
          onKomentar={medanUmum('comment')}
          menyimpan={menyimpan}
          onSimpan={() => void simpan()}
          onTutup={() => {
            setDialog(null)
          }}
        />
      )}
      {dialog === 'edit' && (
        <DialogEdit
          onEdit={() => {
            // `SetViewEdit` 1 b151: `IsView := "false"`.
            setLihat(false)
            setDialog(null)
          }}
          onTutup={() => {
            setDialog(null)
          }}
        />
      )}
      {rate !== null && (
        <ModalRate
          riRateId={rate}
          onTutup={() => {
            setRate(null)
          }}
        />
      )}
    </section>
  )
}

/** Grid bersarang dengan ikon grid bawaan `pzPegaDefaultGridIcons` (tambah/hapus baris, vis `IsView!='true'`). */
function GridSederhana<T extends { asli?: string }>({
  baris,
  lihat,
  kolom,
  kosong,
  onUbah,
}: {
  baris: readonly T[]
  lihat: boolean
  /** Kolom: label, medan, dan (opsional) daftar pilihan - sel berpilihan menjadi pilihan di mode sunting. */
  kolom: ReadonlyArray<readonly [label: string, medan: keyof T & string, opsi?: readonly string[]]>
  kosong: T
  onUbah: (d: T[]) => void
}) {
  return (
    <>
      {!lihat && (
        <div className="aksi-baris">
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => onUbah([...baris, { ...kosong }])}>
            {LAIN_MPNL.tambahBaris}
          </button>
        </div>
      )}
      {baris.length === 0 ? (
        <Kosong pesan={LAIN_MPNL.kosong} />
      ) : (
        <div className="mpnl-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {kolom.map(([l]) => (
                  <th key={l}>{l}</th>
                ))}
                <th className="table__actions" />
              </tr>
            </thead>
            <tbody>
              {baris.map((b, i) => (
                <tr key={i} className="inbox__baris">
                  {kolom.map(([l, m, opsi]) => (
                    <td key={m}>
                      {opsi !== undefined && !lihat ? (
                        <>
                          <PilihanMedan
                            labelAria={l}
                            value={String(b[m] ?? '')}
                            opsi={opsi.map((v) => ({ value: v, label: v }))}
                            kelas="field__input mpnl-sel-isi"
                            onChange={(v) => onUbah(ganti<T>(baris, i, { [m]: v } as Partial<T>))}
                          />
                          {/* Kotak pilihan bawaan tidak membungkus teks: nama utuh nilai terpilih di bawahnya. */}
                          {String(b[m] ?? '') !== '' && <div className="mpnl-sel-teks">{String(b[m] ?? '')}</div>}
                        </>
                      ) : (
                        <SelIsi
                          label={l}
                          nilai={String(b[m] ?? '')}
                          readOnly={lihat}
                          onUbah={(v) => onUbah(ganti<T>(baris, i, { [m]: v } as Partial<T>))}
                        />
                      )}
                    </td>
                  ))}
                  <td className="table__actions">
                    {!lihat && (
                      <button type="button" className="btn btn--ghost btn--sm" onClick={() => onUbah(buang(baris, i))}>
                        {LAIN_MPNL.hapusBaris}
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </>
  )
}

/** Grid ber-tombol `Add` / `Delete` korpus + ikon salin baris (`CopyFinancialWriting` / `CopyUnderWritingLimit`). */
function GridBerangka<T extends { asli?: string }>({
  baris,
  lihat,
  kolom,
  tambah,
  hapus,
  onUbah,
}: {
  baris: readonly T[]
  lihat: boolean
  kolom: ReadonlyArray<readonly [label: string, medan: keyof T & string]>
  tambah: ReactNode
  hapus: (i: number) => ReactNode
  onUbah: (d: T[]) => void
}) {
  return (
    <>
      {!lihat && <div className="aksi-baris">{tambah}</div>}
      {baris.length === 0 ? (
        <Kosong pesan={LAIN_MPNL.kosong} />
      ) : (
        <div className="mpnl-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {kolom.map(([l]) => (
                  <th key={l}>{l}</th>
                ))}
                {!lihat && <th className="table__actions" />}
              </tr>
            </thead>
            <tbody>
              {baris.map((b, i) => (
                <tr key={i} className="inbox__baris">
                  {kolom.map(([l, m]) => (
                    <td key={m}>
                      <SelIsi
                        label={l}
                        nilai={String(b[m] ?? '')}
                        readOnly={lihat}
                        onUbah={(v) => onUbah(ganti<T>(baris, i, { [m]: v } as Partial<T>))}
                      />
                    </td>
                  ))}
                  {!lihat && (
                    <td className="table__actions">
                      <button type="button" className="btn btn--ghost btn--sm" onClick={() => onUbah(salinBaris(baris, i))}>
                        {LAIN_MPNL.salinBaris}
                      </button>{' '}
                      {hapus(i)}
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </>
  )
}
