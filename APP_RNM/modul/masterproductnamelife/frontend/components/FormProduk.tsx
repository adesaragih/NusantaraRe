// Form produk - `InboxProductName` mode form (`OutputParam.DATASHOW = 1`), urutan wadah XML:
//   b2573 sisi umum + `LIEN CLAUSE` + `DOCUMENT CLAIM` + sisi inward   (PARITAS §3.1, §3.2, §5)
//   b30995 `PLAN LIST`, `FINANCIAL UNDERWRITING`, `UNDERWRITING LIMIT`, checkbox `On Retention`
//   b57867 tombol bawah `Close` `Save` `Edit` `Copy` `Generate`            (PARITAS §3.3)
//   b61044 grid komentar (baca-saja)                                       (PARITAS §5)
//   b64133 lampiran                                                        (PARITAS §6)
//
// Mode lihat (`ProductName.IsView == 'true'`, sesudah `View` b74798 → `SetProductName` 8 b2147): medan ber-`ro`
// baca-saja, tombol `Choose*` (wadah `IsView!='true'`) dan `Save` tersembunyi, `Edit` tampil. Autocomplete,
// dropdown `Cause Of Loss`, checkbox `On Retention`, dan tombol `Add`/`Delete` grid tidak ber-`ro` di XML.
// Grid `OUTWARD` (wadah `1==2`) tidak dirender; isinya ditulis server (`hitungOutward`).

import { useEffect, useState, type ReactNode } from 'react'

import { Area, Field, Gagal, Kosong, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import {
  BATAS_SARAN,
  cariMaster,
  cariPlan,
  simpanProduk,
  unduhGenerate,
  type BarisDokumen,
  type BarisFinUW,
  type BarisLien,
  type BarisPlan,
  type BarisUWLimit,
  type JenisMaster,
  type JenisPlan,
  type NilaiMaster,
  type Produk,
  type ProdukInward,
  type ProdukUmum,
} from '../api'
import {
  PILIHAN_PEMBAYARAN,
  hitungMaxSumReasured,
  namaTreaty,
  salinBaris,
  salinProduk,
  tampilPremiumFactor,
  tampilViewRate,
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
  TOMBOL_MPNL,
  UMUM_MPNL,
  UWLIMIT_MPNL,
} from '../labels'
import { DialogEdit, DialogSimpan } from './Dialog'
import ModalRate from './ModalRate'
import PanelLampiran from './PanelLampiran'
import PemilihMaster from './PemilihMaster'
import Saran from './Saran'

/** Pemilih yang sedang terbuka: tombol pembukanya, jenis RD, dan penerima `set*_DT`. */
interface PemilihTerbuka {
  judul: string
  jenis: JenisMaster
  kolomNama?: string
  pilih: (v: NilaiMaster) => void
}

const cariDari = (jenis: JenisMaster) => async (kata: string) => (await cariMaster(jenis, kata, BATAS_SARAN)).daftar
const cariJenisPlan = async (kata: string) => (await cariPlan(kata)).daftar

/** `pxDateTime` - tanggal `YYYY-MM-DD`; teks lama yang bukan tanggal tampil apa adanya (tidak dibuang). */
function MedanTanggal({ label, value, onChange, readOnly }: { label: string; value: string; onChange: (v: string) => void; readOnly: boolean }) {
  const iso = value === '' || /^\d{4}-\d{2}-\d{2}$/.test(value)
  return (
    <div className="field">
      <label className="field__label">{label}</label>
      <input
        className="field__input"
        type={iso ? 'date' : 'text'}
        value={value}
        readOnly={readOnly}
        onChange={(e) => {
          onChange(e.target.value)
        }}
      />
    </div>
  )
}

/** Sel grid bersarang yang dapat disunting (`ro = ProductName.IsView=='true'`). */
function SelIsi({ nilai, onUbah, readOnly, label }: { nilai: string; onUbah: (v: string) => void; readOnly: boolean; label: string }) {
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
  const [lihat, setLihat] = useState(lihatAwal)
  const [pesan, setPesan] = useState<string | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [dialog, setDialog] = useState<'simpan' | 'edit' | null>(null)
  const [pemilih, setPemilih] = useState<PemilihTerbuka | null>(null)
  const [rate, setRate] = useState<string | null>(null)
  const [penyebab, setPenyebab] = useState<NilaiMaster[]>([])
  const [galatPenyebab, setGalatPenyebab] = useState<unknown>(null)

  // `Cause Of Loss` b10661: dropdown RD `BrowseCauseofLossLife_RD` b10781 (nilai/tampil `.CauseofLoss`).
  useEffect(() => {
    let batal = false
    cariMaster('penyebab', '')
      .then((d) => {
        if (!batal) setPenyebab(d.daftar)
      })
      .catch((e: unknown) => {
        if (!batal) setGalatPenyebab(e)
      })
    return () => {
      batal = true
    }
  }, [])

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

  function bukaPemilih(judul: string, jenis: JenisMaster, pilih: (v: NilaiMaster) => void, kolomNama?: string): void {
    setPemilih({ judul, jenis, pilih, kolomNama })
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

  const opsiPenyebab: Opsi[] = penyebab.map((v) => ({ value: v.nama, label: v.nama }))

  /** Medan autocomplete + tombol `Choose*` (wadah `IsView!='true'`) - tombolnya ditulis pemanggil. */
  function medanMaster(
    label: string,
    jenis: JenisMaster,
    nama: string,
    ketik: (v: string) => void,
    pilih: (v: NilaiMaster) => void,
    tombol: (buka: (judul: string) => void) => ReactNode,
    pilihTombol: (v: NilaiMaster) => void = pilih,
  ) {
    return (
      <div className="mpnl-medan-pilih">
        <Saran<NilaiMaster>
          label={label}
          nilai={nama}
          onKetik={ketik}
          cari={cariDari(jenis)}
          teks={(v) => v.nama}
          kunci={(v) => v.id}
          onPilih={pilih}
        />
        {!lihat && tombol((judul) => bukaPemilih(judul, jenis, pilihTombol))}
      </div>
    )
  }

  return (
    <section className="panel">
      {pesan !== null && <p className="mpnl-pesan">{pesan}</p>}
      {galat !== null && <Gagal galat={galat} />}

      {/* ---- wadah b2573: sisi umum ---- */}
      <div className="form-grid">
        <Field
          label={UMUM_MPNL.productName}
          value={u.productName}
          readOnly={lihat}
          onChange={(v) => {
            // onChange `SetTreatyName_Act` b3686.
            ubahUmum({ productName: v, inwardName: namaTreaty(v, w.policyHolderName) })
          }}
        />
        <Field label={UMUM_MPNL.productCode} value={p.id} onChange={() => undefined} readOnly />
        {medanMaster(
          UMUM_MPNL.ceding,
          'ceding',
          u.ceding,
          (v) => {
            ubahUmum({ ceding: v, cedingId: '' })
          },
          (v) => {
            ubahUmum({ ceding: v.nama, cedingId: v.id })
          },
          (buka) => (
            <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(UMUM_MPNL.chooseCeding)}>
              {UMUM_MPNL.chooseCeding}
            </button>
          ),
        )}
        {medanMaster(
          UMUM_MPNL.sob,
          'sob',
          u.sobName,
          (v) => {
            ubahUmum({ sobName: v, sobId: '' })
          },
          (v) => {
            ubahUmum({ sobName: v.nama, sobId: v.id })
          },
          (buka) => (
            <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(UMUM_MPNL.chooseSob)}>
              {UMUM_MPNL.chooseSob}
            </button>
          ),
        )}
        <Field label={UMUM_MPNL.deduction} value={u.riComm} readOnly={lihat} onChange={medanUmum('riComm')} />
        {medanMaster(
          UMUM_MPNL.riRisk,
          'ri-risk',
          u.riRisk,
          (v) => {
            ubahUmum({ riRisk: v, riRiskId: '' })
          },
          (v) => {
            ubahUmum({ riRisk: v.nama, riRiskId: v.id })
          },
          (buka) => (
            <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(UMUM_MPNL.chooseRiRisk)}>
              {UMUM_MPNL.chooseRiRisk}
            </button>
          ),
        )}
        <Field label={UMUM_MPNL.treatyName} value={u.inwardName} readOnly={lihat} onChange={medanUmum('inwardName')} />
        <Field label={UMUM_MPNL.treatyNumber} value={u.treatyNumber} readOnly={lihat} onChange={medanUmum('treatyNumber')} />
        <div className="mpnl-medan-pilih">
          <Pilih
            label={UMUM_MPNL.causeOfLoss}
            value={u.cause}
            opsi={opsiPenyebab}
            onChange={(v) => {
              ubahUmum({ cause: v, causeId: penyebab.find((x) => x.nama === v)?.id ?? '' })
            }}
          />
          {!lihat && (
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              onClick={() =>
                bukaPemilih(UMUM_MPNL.chooseCause, 'penyebab', (v) => {
                  ubahUmum({ cause: v.nama, causeId: v.id })
                })
              }
            >
              {UMUM_MPNL.chooseCause}
            </button>
          )}
        </div>
      </div>
      {galatPenyebab !== null && <Gagal galat={galatPenyebab} />}

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

      {/* `DOCUMENT CLAIM` b14601 - `Document List` dropdown `associated` (daftar tak ikut ekspor, OQ-MPNL-05). */}
      <h4 className="mpnl-judul-grid">{DOKUMEN_MPNL.judul}</h4>
      <GridSederhana<BarisDokumen>
        baris={p.documentClaim}
        lihat={lihat}
        kolom={[[DOKUMEN_MPNL.documentList, 'document']]}
        kosong={{ document: '' }}
        onUbah={(d) => {
          setP((x) => ({ ...x, documentClaim: d }))
        }}
      />

      {/* ---- sisi inward (halaman `ProductNameInward`) ---- */}
      <div className="form-grid mpnl-bagian">
        {medanMaster(
          INWARD_MPNL.policyHolder,
          'pemegang-polis',
          w.policyHolderName,
          (v) => {
            // onChange `SetTreatyName_Act` b17168.
            setP((x) => ({
              ...x,
              inward: { ...x.inward, policyHolderName: v, policyHolder: '' },
              umum: { ...x.umum, inwardName: namaTreaty(x.umum.productName, v) },
            }))
          },
          (v) => {
            setP((x) => ({
              ...x,
              inward: { ...x.inward, policyHolderName: v.nama, policyHolder: v.id },
              umum: { ...x.umum, inwardName: namaTreaty(x.umum.productName, v.nama) },
            }))
          },
          (buka) => (
            <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(INWARD_MPNL.choosePolicyHolder)}>
              {INWARD_MPNL.choosePolicyHolder}
            </button>
          ),
          (v) => {
            // `setPolicyHolder_DT` b2416 - tanpa `SetTreatyName_Act`.
            ubahInward({ policyHolderName: v.nama, policyHolder: v.id })
          },
        )}
        <Field label={INWARD_MPNL.insured} value={w.insured} readOnly={lihat} onChange={medanInward('insured')} />
        <Field label={INWARD_MPNL.addendumNo} value={w.addendumNo} readOnly={lihat} onChange={medanInward('addendumNo')} />
        <Field label={INWARD_MPNL.addendum} value={w.addendumWord} readOnly={lihat} onChange={medanInward('addendumWord')} />
        <Field label={INWARD_MPNL.amandementNo} value={w.amandementNo} readOnly={lihat} onChange={medanInward('amandementNo')} />
        <Field label={INWARD_MPNL.amandement} value={w.amandementSchd} readOnly={lihat} onChange={medanInward('amandementSchd')} />
        <Field label={INWARD_MPNL.maxExpiredClaim} value={w.maxExpiredClaim} readOnly={lihat} onChange={medanInward('maxExpiredClaim')} />
        <MedanTanggal label={INWARD_MPNL.begin} value={w.begin} readOnly={lihat} onChange={medanInward('begin')} />
        <MedanTanggal label={INWARD_MPNL.stnc} value={w.stnc} readOnly={lihat} onChange={medanInward('stnc')} />
        <Field label={INWARD_MPNL.cedingRetention} value={w.cedingRetentionNum} readOnly={lihat} onChange={medanInward('cedingRetentionNum')} />
        <Field label={INWARD_MPNL.cedingLimit} value={w.cedingLimit} readOnly={lihat} onChange={ubahHitung('cedingLimit')} />
        <Field label={INWARD_MPNL.brokerage} value={w.brokerage} readOnly={lihat} onChange={medanInward('brokerage')} />
        <Field label={INWARD_MPNL.minAge} value={w.minAge} readOnly={lihat} onChange={medanInward('minAge')} />
        <Field label={INWARD_MPNL.maxAge} value={w.maxAge} readOnly={lihat} onChange={medanInward('maxAge')} />
        <Field label={INWARD_MPNL.expiryAge} value={w.expiryAge} readOnly={lihat} onChange={medanInward('expiryAge')} />
        <Field label={INWARD_MPNL.extraPremi} value={w.extraPremi} readOnly={lihat} onChange={medanInward('extraPremi')} />
        <Field label={INWARD_MPNL.minSumInsured} value={w.minSumInsured} readOnly={lihat} onChange={medanInward('minSumInsured')} />
        <Field label={INWARD_MPNL.maxSumInsured} value={w.maxSumInsured} readOnly={lihat} onChange={ubahHitung('maxSumInsured')} />
        <Field label={INWARD_MPNL.maxSumReasured} value={w.maxSumReasured} readOnly={lihat} onChange={medanInward('maxSumReasured')} />
        <Field label={INWARD_MPNL.rnmShare} value={w.rnmShare} readOnly={lihat} onChange={medanInward('rnmShare')} />
        <span className="mpnl-catatan-medan">{INWARD_MPNL.ofSumReasured}</span>
        <Field label={INWARD_MPNL.rnmLimit} value={w.rnmLimitNum} readOnly={lihat} onChange={medanInward('rnmLimitNum')} />
        {tampilPremiumFactor(w.payment) && (
          <Field label={INWARD_MPNL.premiumFactor} value={w.premiumFactor} readOnly={lihat} onChange={medanInward('premiumFactor')} />
        )}
        {lihat ? (
          // `ro = ProductName.IsView=='true'` b25611 - `Pilih` bersama tidak punya mode baca-saja.
          <Field
            label={INWARD_MPNL.payment}
            value={PILIHAN_PEMBAYARAN.find((o) => o.value === w.payment)?.label ?? w.payment}
            onChange={() => undefined}
            readOnly
          />
        ) : (
          <Pilih label={INWARD_MPNL.payment} value={w.payment} opsi={[...PILIHAN_PEMBAYARAN]} onChange={medanInward('payment')} />
        )}
        <Area label={INWARD_MPNL.subjectTo} value={w.subjectTo} onChange={lihat ? () => undefined : medanInward('subjectTo')} />
        <Field label={INWARD_MPNL.annuityInterest} value={w.annuityInterest} readOnly={lihat} onChange={medanInward('annuityInterest')} />
        <Field
          label={INWARD_MPNL.premiumRefundFactor}
          value={w.premiumRefundFactor}
          readOnly={lihat}
          onChange={medanInward('premiumRefundFactor')}
        />
        <Field label={INWARD_MPNL.maxDataReceive} value={w.maxDataReceive} readOnly={lihat} onChange={medanInward('maxDataReceive')} />
        <MedanTanggal label={INWARD_MPNL.mature} value={w.mature} readOnly={lihat} onChange={medanInward('mature')} />
        {/* `Birthday` b27960 - radio `associated` (pilihan tak ikut ekspor, OQ-MPNL-05): isian teks. */}
        <Field label={INWARD_MPNL.birthday} value={w.birthday} readOnly={lihat} onChange={medanInward('birthday')} />
        {medanMaster(
          INWARD_MPNL.currency,
          'mata-uang',
          w.currency,
          (v) => {
            ubahInward({ currency: v, currencyId: '' })
          },
          (v) => {
            ubahInward({ currency: v.nama, currencyId: v.id })
          },
          (buka) => (
            <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(INWARD_MPNL.chooseCurrency)}>
              {INWARD_MPNL.chooseCurrency}
            </button>
          ),
        )}
        <Field label={INWARD_MPNL.extraMortality} value={w.extraMortality} readOnly={lihat} onChange={medanInward('extraMortality')} />
        <Field label={INWARD_MPNL.maxContract} value={w.maxContract} readOnly={lihat} onChange={medanInward('maxContract')} />
        <Field label={INWARD_MPNL.proportionalTable} value={w.proportionalTable} readOnly={lihat} onChange={medanInward('proportionalTable')} />
      </div>

      {/* ---- wadah b30995: PLAN LIST ---- */}
      <h4 className="mpnl-judul-grid">{PLAN_MPNL.judul}</h4>
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
                  <Saran<JenisPlan>
                    labelAria={PLAN_MPNL.planName}
                    nilai={b.plan}
                    readOnly={lihat}
                    onKetik={(v) => {
                      setP((x) => ({ ...x, planList: ganti<BarisPlan>(x.planList, i, { plan: v, planId: '' }) }))
                    }}
                    cari={cariJenisPlan}
                    teks={(t) => t.coverName}
                    kunci={(t) => t.id}
                    onPilih={(t) => {
                      // Autocomplete b33198: `.CoverName` → `.Plan`, `.ID` → `.PlanID`, `.Business` → `.Name`, `.Benefit` → `.Benefit`.
                      setP((x) => ({
                        ...x,
                        planList: ganti<BarisPlan>(x.planList, i, { plan: t.coverName, planId: t.id, name: t.business, benefit: t.benefit }),
                      }))
                    }}
                  />
                </td>
                <td>
                  <SelIsi
                    label={PLAN_MPNL.bussines}
                    nilai={b.name}
                    readOnly={lihat}
                    onUbah={(v) => {
                      setP((x) => ({ ...x, planList: ganti<BarisPlan>(x.planList, i, { name: v }) }))
                    }}
                  />
                </td>
                <td>
                  <SelIsi
                    label={PLAN_MPNL.benefit}
                    nilai={b.benefit}
                    readOnly={lihat}
                    onUbah={(v) => {
                      setP((x) => ({ ...x, planList: ganti<BarisPlan>(x.planList, i, { benefit: v }) }))
                    }}
                  />
                </td>
                <td>{b.riRate}</td>
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
                      onClick={() =>
                        bukaPemilih(
                          PLAN_MPNL.chooseRiRate,
                          'ri-rate',
                          (v) => {
                            // `SetRIRate` 1 b249: `.RIRATEID ← id`, `.RIRATE ← usedby`.
                            setP((x) => ({ ...x, planList: ganti<BarisPlan>(x.planList, i, { riRate: v.nama, riRateId: v.id }) }))
                          },
                          PEMILIH_MPNL.kolomRiRateName,
                        )
                      }
                    >
                      {PLAN_MPNL.chooseRiRate}
                    </button>
                  )}{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      setP((x) => ({ ...x, planList: buang(x.planList, i) }))
                    }}
                  >
                    {PLAN_MPNL.delete}
                  </button>
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
          onChange={(e) => {
            const v = e.target.checked
            setP((x) => ({ ...x, umum: { ...x.umum, isOrs: v }, hitungOutward: true }))
          }}
        />
        {UMUM_MPNL.onRetention}
      </label>

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
        {lihat && (
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
        </button>{' '}
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
                <td>{k.date}</td>
                <td>{k.operatorName}</td>
                <td>{k.suggest}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* ---- wadah b64133: lampiran (produk tersimpan) ---- */}
      {p.id !== '' && <PanelLampiran produkId={p.id} />}

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
      {pemilih !== null && (
        <PemilihMaster
          judul={pemilih.judul}
          jenis={pemilih.jenis}
          kolomNama={pemilih.kolomNama}
          onPilih={pemilih.pilih}
          onTutup={() => {
            setPemilih(null)
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
  kolom: ReadonlyArray<readonly [label: string, medan: keyof T & string]>
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
      <div className="aksi-baris">{tambah}</div>
      {baris.length === 0 ? (
        <Kosong pesan={LAIN_MPNL.kosong} />
      ) : (
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
                <td className="table__actions">
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => onUbah(salinBaris(baris, i))}>
                    {LAIN_MPNL.salinBaris}
                  </button>{' '}
                  {hapus(i)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </>
  )
}
