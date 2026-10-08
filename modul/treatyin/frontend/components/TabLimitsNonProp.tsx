// Tab **Limits** cabang NON-PROPORSIONAL — dibaca dari ekspor
// (`labelsLimitsNP.ts`): grid layer → rincian `Layers` (Treaty Group →
// `CoBList`, limit dua mata uang, Reinstatement, Premium Earned, MDP, ROL)
// → Summary of Limit → Total All Layers.
//
// ⭐ Rumusnya hidup di services (`hitung_limit_np.go`), disalin dari
// Activity dan diukur terhadap data Pega. Layar ini hanya mengirim isian
// dan menggabungkan medan yang Activity TULIS — nol rumus di sini.
//
// ⭐ Tampilan (permintaan pemakai 6 Oktober 2026, "UI/UX hebat"): tiap
// layer satu KARTU bernomor dengan ringkasan limit/deductible/ROL di
// kepalanya; rinciannya dikelompokkan per bagian, dan limit dua mata uang
// disusun sebagai matriks (baris = 100 % Limit · Agregate · Deductible,
// kolom = mata uang 1 · mata uang 2) — tiga baris Pega yang mengikat SATU
// properti mata uang kini memilihnya SEKALI di kepala kolom.
//
// ⛔ Hasil suntingan dan hitungan hidup di salinan pohon layar ini; jalur
// Save menunggu keputusan pemilik proses.

import { useEffect, useRef, useState, type ReactNode } from 'react'

import { Field, FieldAngka, Kosong, Panel, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilKelasBisnis,
  ambilOpsiLimits,
  hitungLimitNP,
  type AksiLimitNP,
  type EgnpiLimitNP,
  type KursLimitNP,
  type LimitsAkar,
  type OpsiLimits,
  type OpsiPilihan,
  type PilihanWarisan,
  type SimpulLimit,
} from '../api'
import { GRID_TOTAL, KOLOM_LAYER, KOLOM_REINSTATEMENT, KOLOM_RINGKASAN, LIMITS_NP, kontrakRevisi, type KolomNP } from '../labelsLimitsNP'
import type { ModeForm } from '../mode'
import { DropdownDaftar } from './IsianAuto'
import { Bagian, Chip, KartuLipat, KepalaBagian, TombolHapus, TombolTambah } from './limitsUI'
import { formatLimit, teksDari } from './TabLimitsProp'

function larikDari(s: SimpulLimit, kunci: string): SimpulLimit[] {
  const v = s[kunci]
  return Array.isArray(v) ? v : []
}

function ganti<T>(larik: readonly T[], i: number, baru: T): T[] {
  return larik.map((x, j) => (j === i ? baru : x))
}

/** Angka tampil (mode lihat) — pemformat modul, bukan pemformat kedua. */
const tampil = (nilai: string, desimal: number | null) => formatLimit('uang', desimal, nilai)

/**
 * Medan yang tiap aksi TULIS — hanya ini yang digabung kembali ke layer.
 * Menggabung seluruh jawaban akan menimpa suntingan yang tidak dikirim.
 */
export const DITULIS_LIMIT_NP: Readonly<Record<AksiLimitNP, readonly string[]>> = {
  egnpi: ['EgnpiTotalList'],
  reinstatement: ['Reinstatement_List'],
  'reinst-jumlah': ['Reinstatement_List'],
  'reinst-persen': ['Reinstatement_List'],
  'reinst-tambahan': ['Reinstatement_List'],
  adj: ['PremiumEarnedList', 'ROLPct'],
  mdp: ['MDPList', 'MDPMinList', 'ROLPct'],
  total: [],
  'nilai-list': ['EgnpiTotalList', 'PremiumEarnedList', 'MDPList', 'MDPMinList', 'ROLPct'],
}

/** Aksi yang menulis SEMUA layer; selainnya hanya layer `indeks`. */
const SEMUA_LAYER: readonly AksiLimitNP[] = ['egnpi', 'nilai-list']

/** Dropdown `associated` — nilai tersimpan apa adanya bila di luar domain. */
function PilihNP({
  label,
  nilai,
  opsi,
  kosong,
  bisaUbah,
  onUbah,
}: {
  label: string
  nilai: string
  opsi: readonly OpsiPilihan[]
  kosong?: string
  bisaUbah: boolean
  onUbah: (v: string) => void
}) {
  if (!bisaUbah) {
    const t = opsi.find((o) => o.value === nilai)?.label ?? nilai
    return <Field label={label} value={t} readOnly onChange={() => undefined} />
  }
  return <Pilih label={label} value={nilai} kosong={kosong ?? LIMITS_NP.pilihKosong} opsi={[...opsi]} onChange={onUbah} />
}

/** Medan angka/teks satu layer; `onLepas` = peristiwa `change` Pega. */
function MedanNP({
  label,
  nilai,
  desimal,
  bisaUbah,
  onUbah,
  onLepas,
}: {
  label: string
  nilai: string
  desimal: number | null
  bisaUbah: boolean
  onUbah: (v: string) => void
  onLepas?: () => void
}) {
  return (
    <div className="trin__limit-medan" onBlur={bisaUbah ? onLepas : undefined}>
      {/* ⭐ PEMISAH RIBUAN SAAT MENGETIK — permintaan pemilik proses
          7 Oktober 2026. `FieldAngka` menerima dan mengembalikan bentuk
          KABEL (titik desimal), jadi nol pemanggil berubah selain nama
          komponennya.

          ⛔ RALAT 7 Oktober 2026, pukul berikutnya. Bentuk pertama menulis
          `desimal ?? 2` — memaksa SETIAP medan ber-`desimal === null`
          menjadi uang berdua desimal. Akibatnya `Layer` dan `Part` di kartu
          Layers tampil `1,00`, dan pemilik proses melaporkannya: *"pada tab
          limits di layer nya tidak perlu format nya krn itu memang angka
          biasa"*.

          ⚠️ `desimal === null` memang BUKAN "pakai bawaan". Ia berarti
          presisi kolomnya TIDAK TERBACA di gambar mana pun — dan medan
          seperti `Layer`/`Part` bukan uang sama sekali, melainkan nomor
          urut. Yang presisinya tidak diketahui LEWAT APA ADANYA, persis
          seperti sebelum pemisah ribuan ada. */}
      {desimal === null ? (
        <Field label={label} value={bisaUbah ? nilai : tampil(nilai, desimal)} readOnly={!bisaUbah} onChange={onUbah} />
      ) : (
        <FieldAngka label={label} value={nilai} desimal={desimal} readOnly={!bisaUbah} onChange={onUbah} />
      )}
    </div>
  )
}

/** Kotak centang bernilai teks `"true"`/`"false"`, seperti dokumen. */
function Centang({ label, nilai, bisaUbah, onUbah }: { label: string; nilai: string; bisaUbah: boolean; onUbah: (v: string) => void }) {
  return (
    <label className="trin__limit-medan">
      <input
        type="checkbox"
        checked={nilai === 'true'}
        disabled={!bisaUbah}
        onChange={(e) => {
          onUbah(e.target.checked ? 'true' : 'false')
        }}
      />{' '}
      {label}
    </label>
  )
}

/** Grid `Currency · Value` — Premium Earned, Min Premium Amt, MDP, total. */
function GridNilai({
  judul,
  baris,
  bisaUbahMataUang,
  bisaUbahNilai,
  tambah,
  mataUang,
  onUbah,
}: {
  judul: string
  baris: readonly SimpulLimit[]
  bisaUbahMataUang: boolean
  bisaUbahNilai: boolean
  tambah: boolean
  mataUang: readonly PilihanWarisan[]
  onUbah: (b: SimpulLimit[]) => void
}) {
  return (
    <div className="trin__limit-grid">
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              <th scope="col">{judul}</th>
              <th scope="col">{LIMITS_NP.nilai}</th>
              {tambah && (
                <th scope="col">
                  <TombolTambah
                    label={LIMITS_NP.tambahMDP}
                    onClick={() => {
                      onUbah([...baris, { Currency: '', CurrencyID: '', Value: '' }])
                    }}
                  />
                </th>
              )}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={tambah ? 3 : 2}>
                  <Kosong pesan={LIMITS_NP.tanpaBaris} />
                </td>
              </tr>
            )}
            {baris.map((b, r) => (
              <tr key={r}>
                <td>
                  {bisaUbahMataUang ? (
                    <DropdownDaftar
                      label=""
                      nilai={teksDari(b, 'Currency')}
                      pilihan={mataUang}
                      bisaUbah
                      onPilih={(nama, id) => {
                        onUbah(ganti(baris, r, { ...b, Currency: nama, CurrencyID: id }))
                      }}
                    />
                  ) : (
                    teksDari(b, 'Currency')
                  )}
                </td>
                <td className="trin__angka">
                  {bisaUbahNilai ? (
                    /* ⭐ Sel angka berpemisah ribuan — sama dengan grid
                       `100% Limit` / `Retention` / `Cession to R/I` cabang
                       Prop. Keduanya grid nilai bermata uang; membiarkan
                       yang satu mentah sementara yang lain terformat adalah
                       jenis selisih yang tidak akan dilaporkan siapa pun
                       sampai angkanya salah dibaca. */
                    <FieldAngka
                      label=""
                      value={teksDari(b, 'Value')}
                      desimal={2}
                      onChange={(v) => {
                        onUbah(ganti(baris, r, { ...b, Value: v }))
                      }}
                    />
                  ) : (
                    tampil(teksDari(b, 'Value'), 2)
                  )}
                </td>
                {tambah && (
                  <td>
                    <TombolHapus
                      label={LIMITS_NP.hapusMDP}
                      labelAkses={`${LIMITS_NP.hapusMDP} ${judul} ${r + 1}`}
                      onClick={() => {
                        onUbah(baris.filter((_, j) => j !== r))
                      }}
                    />
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

/** Grid baca-saja berkolom ekspor. */
function GridBaca({ kolom, baris }: { kolom: readonly KolomNP[]; baris: readonly SimpulLimit[] }) {
  return (
    <div className="table-wrap">
      <table className="trin__tabel">
        <thead>
          <tr>
            {kolom.map((k) => (
              <th key={k.kunci} scope="col">
                {k.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td colSpan={kolom.length}>
                <Kosong pesan={LIMITS_NP.tanpaBaris} />
              </td>
            </tr>
          )}
          {baris.map((b, r) => (
            <tr key={r}>
              {kolom.map((k) => (
                <td key={k.kunci} className={k.desimal === null ? undefined : 'trin__angka'}>
                  {k.desimal === null ? teksDari(b, k.kunci) : tampil(teksDari(b, k.kunci), k.desimal)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

/** Satu Treaty Group layer — kartu; rinciannya `Section/CoBList.xml`. */
function KartuGrup({
  g,
  nomor,
  opsi,
  bisaUbah,
  onUbah,
  onHapus,
  onGantiGrup,
}: {
  g: SimpulLimit
  nomor: number
  opsi: OpsiLimits
  bisaUbah: boolean
  onUbah: (g: SimpulLimit) => void
  onHapus: () => void
  onGantiGrup: (g: SimpulLimit) => void
}) {
  const [kelas, setKelas] = useState<PilihanWarisan[]>([])
  const idGrup = teksDari(g, 'TreatyGroupID')
  useEffect(() => {
    if (!bisaUbah || idGrup === '') return
    let dibuang = false
    ambilKelasBisnis(idGrup)
      .then((d) => {
        if (!dibuang) setKelas(d)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [bisaUbah, idGrup])
  const pilihGrup = (nama: string, id: string) => {
    onGantiGrup({ ...g, TreatyGroup: nama, TreatyGroupID: id })
  }
  const daftarCoB = larikDari(g, 'ClassOfBusinessList')
  return (
    <KartuLipat
      nomor={nomor}
      judul={teksDari(g, 'TreatyGroup')}
      judulKosong={LIMITS_NP.belumDipilih}
      bukaAwal={teksDari(g, 'TreatyGroup') === ''}
      meta={
        <>
          <Chip label={LIMITS_NP.kelasBisnis} nilai={String(daftarCoB.length)} />
          {teksDari(g, 'IsROLProfile') === 'true' && <span className="tl-chip">{LIMITS_NP.rolProfile}</span>}
        </>
      }
      aksi={bisaUbah && <TombolHapus label={LIMITS_NP.hapus} labelAkses={`${LIMITS_NP.hapus} ${LIMITS_NP.treatyGroup} ${nomor}`} onClick={onHapus} />}
    >
      <div className="form-grid">
        <DropdownDaftar label={LIMITS_NP.treatyGroup} nilai={teksDari(g, 'TreatyGroup')} pilihan={opsi.kelompokTreaty} bisaUbah={bisaUbah} onPilih={pilihGrup} />
        <Centang
          label={LIMITS_NP.rolProfile}
          nilai={teksDari(g, 'IsROLProfile')}
          bisaUbah={bisaUbah}
          onUbah={(v) => {
            onUbah({ ...g, IsROLProfile: v })
          }}
        />
      </div>
      {/* Add CoB / Delete MATI di ekspor (`ViewState !='1' && 1=2`). */}
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              <th scope="col">{LIMITS_NP.kelasBisnis}</th>
            </tr>
          </thead>
          <tbody>
            {daftarCoB.length === 0 && (
              <tr>
                <td>
                  <Kosong pesan={LIMITS_NP.tanpaBaris} />
                </td>
              </tr>
            )}
            {daftarCoB.map((c, r) => (
              <tr key={r}>
                <td>
                  <DropdownDaftar
                    label=""
                    nilai={teksDari(c, 'ClassOfBusiness')}
                    pilihan={kelas}
                    bisaUbah={bisaUbah}
                    onPilih={(nama, id) => {
                      // DT `SetCoB` — nama dan kode Class of Business.
                      onUbah({ ...g, ClassOfBusinessList: ganti(daftarCoB, r, { ...c, ClassOfBusiness: nama, ClassOfBusinessID: id }) })
                    }}
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </KartuLipat>
  )
}

/** Rincian satu layer — `Section/Layers.xml`. Diekspor untuk uji. */
export function RincianLayer({
  l,
  opsi,
  bisaUbah,
  modeUbah,
  onUbah,
  hitung,
  onGantiGrup,
}: {
  l: SimpulLimit
  opsi: OpsiLimits
  bisaUbah: boolean
  /** Mode Edit saja — No RIP dan ROL % tidak ikut kunci materialitas. */
  modeUbah: boolean
  onUbah: (l: SimpulLimit) => void
  hitung: (aksi: AksiLimitNP, baris?: number) => void
  onGantiGrup: (l: SimpulLimit) => void
}) {
  const set = (kunci: string) => (v: string) => {
    onUbah({ ...l, [kunci]: v })
  }
  const grup = larikDari(l, 'TreatyGroupList')
  const reinst = larikDari(l, 'Reinstatement_List')
  const limit2 = teksDari(l, 'Limit2')
  const kolomReinst = KOLOM_REINSTATEMENT.filter((k) => k.kunci !== 'ReinstatementAmount2' || (limit2 !== '' && Number(limit2) !== 0))

  /** Sel mata uang di kepala kolom matriks — satu properti per kolom. */
  const mataUang = (kunciMU: 'Currency' | 'Currency2', label: string) => (
    <DropdownDaftar
      label={label}
      nilai={teksDari(l, kunciMU)}
      pilihan={opsi.mataUang}
      bisaUbah={bisaUbah}
      onPilih={(nama, id) => {
        // Pega: ketiga baris satu mata uang mengikat SATU properti, dan
        // keduanya menimpa `.CurrencyID` (`col:.ID->.CurrencyID`).
        onUbah({ ...l, [kunciMU]: nama, CurrencyID: id })
      }}
    />
  )
  const selMatriks = (kunci: string, desimal: number | null, label: string) =>
    bisaUbah ? (
      <input
        className="field__input"
        value={teksDari(l, kunci)}
        aria-label={label}
        onChange={(e) => {
          onUbah({ ...l, [kunci]: e.target.value })
        }}
      />
    ) : (
      tampil(teksDari(l, kunci), desimal)
    )
  const barisMatriks: readonly { label: string; k1: string; k2: string; desimal: number | null }[] = [
    { label: LIMITS_NP.limit100, k1: 'Limit', k2: 'Limit2', desimal: null },
    { label: LIMITS_NP.agregat, k1: 'AgregateLimit', k2: 'AgregateLimit2', desimal: null },
    { label: LIMITS_NP.deductible, k1: 'Deductible', k2: 'Deductible2', desimal: 2 },
  ]

  return (
    <div className="tl-rincian">
      <Bagian judul={LIMITS_NP.layers}>
        <div className="form-grid">
          <PilihNP label={LIMITS_NP.layers} nilai={teksDari(l, 'LayerType')} opsi={opsi.jenisLayer ?? []} bisaUbah={bisaUbah} onUbah={set('LayerType')} />
          <MedanNP label={LIMITS_NP.layer} nilai={teksDari(l, 'Layer')} desimal={null} bisaUbah={bisaUbah} onUbah={set('Layer')} />
          <PilihNP label={LIMITS_NP.partOf} nilai={teksDari(l, 'LayerPartType')} opsi={opsi.jenisLayer ?? []} bisaUbah={bisaUbah} onUbah={set('LayerPartType')} />
          <MedanNP label={LIMITS_NP.part} nilai={teksDari(l, 'LayerPart')} desimal={null} bisaUbah={bisaUbah} onUbah={set('LayerPart')} />
        </div>
      </Bagian>

      {/* Grid Treaty Group — expandPane `CoBList`. */}
      <Bagian
        judul={LIMITS_NP.treatyGroup}
        aksi={
          bisaUbah && (
            <TombolTambah
              label={LIMITS_NP.tambahGrup}
              onClick={() => {
                onUbah({ ...l, TreatyGroupList: [...grup, { TreatyGroup: '', TreatyGroupID: '', IsROLProfile: 'false', ClassOfBusinessList: [] }] })
              }}
            />
          )
        }
      >
        {grup.length === 0 ? (
          <Kosong pesan={LIMITS_NP.tanpaBaris} />
        ) : (
          <div className="tl-daftar">
            {grup.map((g, r) => (
              <KartuGrup
                key={r}
                g={g}
                nomor={r + 1}
                opsi={opsi}
                bisaUbah={bisaUbah}
                onUbah={(baru) => {
                  onUbah({ ...l, TreatyGroupList: ganti(grup, r, baru) })
                }}
                onHapus={() => {
                  // deleteRow → `TotalEgnpi`
                  onGantiGrup({ ...l, TreatyGroupList: grup.filter((_, j) => j !== r) })
                }}
                onGantiGrup={(baru) => {
                  // DT `SetIndexLayer_DT` + `TotalEgnpi`
                  onGantiGrup({ ...l, TreatyGroupList: ganti(grup, r, baru) })
                }}
              />
            ))}
          </div>
        )}
        <GridNilai judul={LIMITS_NP.egnpiLayer} baris={larikDari(l, 'EgnpiTotalList')} bisaUbahMataUang={false} bisaUbahNilai={false} tambah={false} mataUang={opsi.mataUang} onUbah={() => undefined} />
      </Bagian>

      <Bagian judul={LIMITS_NP.bagianLimit}>
        <div className="form-grid">
          <PilihNP label={LIMITS_NP.cover} nilai={teksDari(l, 'Cover')} opsi={opsi.cover ?? []} bisaUbah={bisaUbah} onUbah={set('Cover')} />
          <PilihNP
            label={LIMITS_NP.relasi}
            nilai={teksDari(l, 'CurrencyRelation')}
            opsi={opsi.relasiMataUang ?? []}
            kosong={LIMITS_NP.pilihRelasi}
            bisaUbah={bisaUbah}
            onUbah={set('CurrencyRelation')}
          />
        </div>
        <div className="table-wrap">
          <table className="trin__tabel tl-matriks">
            <thead>
              <tr>
                <th scope="col" />
                <th scope="col">{mataUang('Currency', LIMITS_NP.mataUang1)}</th>
                <th scope="col">{mataUang('Currency2', LIMITS_NP.mataUang2)}</th>
              </tr>
            </thead>
            <tbody>
              {barisMatriks.map((b) => (
                <tr key={b.k1}>
                  <th scope="row">{b.label}</th>
                  <td className="trin__angka">{selMatriks(b.k1, b.desimal, `${b.label} 1`)}</td>
                  <td className="trin__angka">{selMatriks(b.k2, b.desimal, `${b.label} 2`)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Bagian>

      <Bagian judul={LIMITS_NP.reinstatement}>
        <Centang label={LIMITS_NP.noRIP} nilai={teksDari(l, 'NoRIPCalculation')} bisaUbah={modeUbah} onUbah={set('NoRIPCalculation')} />
        {/* `ReinstatementValue` → `SetReinstatementPct` saat berubah. */}
        <MedanNP
          label={LIMITS_NP.reinstatement}
          nilai={teksDari(l, 'ReinstatementValue')}
          desimal={0}
          bisaUbah={bisaUbah}
          onUbah={set('ReinstatementValue')}
          onLepas={() => {
            hitung('reinstatement')
          }}
        />
        <div className="table-wrap">
          <table className="trin__tabel">
            <thead>
              <tr>
                {kolomReinst.map((k) => (
                  <th key={k.kunci} scope="col">
                    {k.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {reinst.length === 0 && (
                <tr>
                  <td colSpan={kolomReinst.length}>
                    <Kosong pesan={LIMITS_NP.tanpaBaris} />
                  </td>
                </tr>
              )}
              {reinst.map((b, r) => (
                <tr key={r}>
                  {kolomReinst.map((k) => {
                    const v = teksDari(b, k.kunci)
                    const ubahSel = (x: string) => {
                      onUbah({ ...l, Reinstatement_List: ganti(reinst, r, { ...b, [k.kunci]: x }) })
                    }
                    // DT per sel: % Additional Premium → ReCalculateReinstatement;
                    // Reinstatement % → CalculateReinstatement; Amount IDR →
                    // CalculateReinstatementPct.
                    const pemicu: Partial<Record<string, AksiLimitNP>> = {
                      ReinstatementPct: 'reinst-tambahan',
                      AdditionalPct: 'reinst-jumlah',
                      ReinstatementAmount1: 'reinst-persen',
                    }
                    const aksi = pemicu[k.kunci]
                    if (k.kunci === 'ReinstatementValue' || !bisaUbah) {
                      return (
                        <td key={k.kunci} className={k.desimal === null ? undefined : 'trin__angka'}>
                          {k.desimal === null ? v : tampil(v, k.desimal)}
                        </td>
                      )
                    }
                    if (k.kunci === 'ReinstatementNote') {
                      return (
                        <td key={k.kunci}>
                          <PilihNP label="" nilai={v} opsi={opsi.catatanReinstatement ?? []} bisaUbah onUbah={ubahSel} />
                        </td>
                      )
                    }
                    return (
                      <td key={k.kunci} className="trin__angka">
                        <input
                          className="field__input"
                          value={v}
                          aria-label={k.label}
                          onChange={(e) => {
                            ubahSel(e.target.value)
                          }}
                          onBlur={() => {
                            if (aksi !== undefined) hitung(aksi, r)
                          }}
                        />
                      </td>
                    )
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Bagian>

      <Bagian judul={LIMITS_NP.bagianPremi}>
        <div className="tl-tiga">
          <div>
            <MedanNP
              label={LIMITS_NP.adjRate}
              nilai={teksDari(l, 'AdjRate')}
              desimal={null}
              bisaUbah={bisaUbah}
              onUbah={set('AdjRate')}
              onLepas={() => {
                hitung('adj')
              }}
            />
            <GridNilai
              judul={LIMITS_NP.premiEarned}
              baris={larikDari(l, 'PremiumEarnedList')}
              bisaUbahMataUang={false}
              bisaUbahNilai={bisaUbah}
              tambah={false}
              mataUang={opsi.mataUang}
              onUbah={(b) => {
                onUbah({ ...l, PremiumEarnedList: b })
              }}
            />
          </div>
          <div>
            <MedanNP
              label={LIMITS_NP.mdpMinPct}
              nilai={teksDari(l, 'MDPMinPct')}
              desimal={null}
              bisaUbah={bisaUbah}
              onUbah={set('MDPMinPct')}
              onLepas={() => {
                hitung('mdp')
              }}
            />
            <GridNilai
              judul={LIMITS_NP.mdpMinAmt}
              baris={larikDari(l, 'MDPMinList')}
              bisaUbahMataUang={bisaUbah}
              bisaUbahNilai={bisaUbah}
              tambah={bisaUbah}
              mataUang={opsi.mataUang}
              onUbah={(b) => {
                onUbah({ ...l, MDPMinList: b })
              }}
            />
          </div>
          <div>
            <MedanNP
              label={LIMITS_NP.mdpPct}
              nilai={teksDari(l, 'MDPPct')}
              desimal={null}
              bisaUbah={bisaUbah}
              onUbah={set('MDPPct')}
              onLepas={() => {
                hitung('mdp')
              }}
            />
            <GridNilai
              judul={LIMITS_NP.mdp}
              baris={larikDari(l, 'MDPList')}
              bisaUbahMataUang={bisaUbah}
              bisaUbahNilai={bisaUbah}
              tambah={bisaUbah}
              mataUang={opsi.mataUang}
              onUbah={(b) => {
                onUbah({ ...l, MDPList: b })
              }}
            />
          </div>
        </div>
        <div className="form-grid">
          <Centang label={LIMITS_NP.combineMDP} nilai={teksDari(l, 'IsCombineMDP')} bisaUbah={bisaUbah} onUbah={set('IsCombineMDP')} />
          <MedanNP label={LIMITS_NP.rol} nilai={teksDari(l, 'ROLPct')} desimal={null} bisaUbah={modeUbah} onUbah={set('ROLPct')} />
        </div>
      </Bagian>
    </div>
  )
}

/** Judul kartu layer — `LayerType Layer of LayerPartType LayerPart` (kolom Layers). */
function judulLayer(l: SimpulLimit): string {
  const kiri = `${teksDari(l, 'LayerType')} ${teksDari(l, 'Layer')}`.trim()
  const kanan = `${teksDari(l, 'LayerPartType')} ${teksDari(l, 'LayerPart')}`.trim()
  if (kiri === '' && kanan === '') return ''
  return `${kiri} ${LIMITS_NP.of} ${kanan}`.trim()
}

/** Ringkasan kepala kartu layer — kolom grid layer ekspor + ROL %. */
function ringkasLayer(l: SimpulLimit): ReactNode {
  return (
    <>
      {KOLOM_LAYER.map((k) => (
        <Chip key={k.kunci} label={k.label} nilai={teksDari(l, k.kunci) === '' ? '' : tampil(teksDari(l, k.kunci), k.desimal)} />
      ))}
      <Chip label={LIMITS_NP.rol} nilai={teksDari(l, 'ROLPct') === '' ? '' : formatLimit('persen', 2, teksDari(l, 'ROLPct'))} />
    </>
  )
}

const OPSI_KOSONG: OpsiLimits = { jenisTreaty: [], kelompokTreaty: [], mataUang: [] }

/** Larik akar kosong — kontrak tanpa Summary/Total. */
const AKAR_KOSONG: LimitsAkar = { LimitSummaryList: [], Total: {}, TotalLimitsROL: '' }

export default function TabLimitsNonProp({
  petunjukKosong,
  pohon,
  akar: akarAwal,
  egnpi,
  kurs,
  edmState = '',
  edmJenisMaterial = '',
  mode,
  opsi: opsiAwal,
  onUbah,
}: {
  /** Petunjuk bila kontrak tak punya layer — layar membedakan sebabnya. */
  petunjukKosong?: string
  pohon: readonly SimpulLimit[]
  akar?: LimitsAkar
  egnpi: readonly EgnpiLimitNP[]
  kurs: readonly KursLimitNP[]
  edmState?: string
  /**
   * `EDMMaterialType` — bila `2`, hampir seluruh kontrol tab ini dimatikan
   * (`pyDisabledWhen TreatyIn.EDMMaterialType = 2` di setiap sel kecuali
   * No RIP dan ROL %).
   */
  edmJenisMaterial?: string
  mode: ModeForm
  /** Isi dropdown; bila tidak diberikan, diminta sendiri. */
  opsi?: OpsiLimits
  /**
   * Keadaan terkini layer dan larik akar, dilaporkan tiap kali berubah —
   * form menyimpannya supaya ia hidup lintas pindah tab dan tab Share
   * membaca layer yang SAMA (`TreatyIn.Limits` di clipboard Pega).
   */
  onUbah?: (layers: SimpulLimit[], akar: LimitsAkar) => void
}) {
  const [layers, setLayers] = useState<SimpulLimit[]>(() => [...pohon])
  const [akar, setAkar] = useState<LimitsAkar>(akarAwal ?? AKAR_KOSONG)
  const onUbahTerkini = useRef(onUbah)
  onUbahTerkini.current = onUbah
  useEffect(() => {
    onUbahTerkini.current?.(layers, akar)
  }, [layers, akar])
  const [opsi, setOpsi] = useState<OpsiLimits>(opsiAwal ?? OPSI_KOSONG)
  const [pesan, setPesan] = useState<string[]>([])
  const [gagal, setGagal] = useState('')
  // `TreatyIn.ViewState != 1` — tombol dan isian hidup di mode Edit saja.
  const modeUbah = mode === 'ubah'
  // … dan dimatikan pula oleh kunci materialitas.
  const bisaUbah = modeUbah && edmJenisMaterial.trim() !== '2'

  useEffect(() => {
    if (opsiAwal !== undefined) return
    let dibuang = false
    ambilOpsiLimits()
      .then((o) => {
        if (!dibuang) setOpsi(o)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [opsiAwal])

  /**
   * Menjalankan satu Activity di services lalu menggabung HANYA medan yang
   * Activity itu tulis (`DITULIS_LIMIT_NP`).
   */
  const hitung = (aksi: AksiLimitNP, indeks = 0, baris = 0, dasar: readonly SimpulLimit[] = layers) => {
    if (!modeUbah) return
    setGagal('')
    hitungLimitNP({ aksi, layers: dasar, egnpi, kurs, indeks, baris })
      .then((h) => {
        setPesan(h.pesan)
        const ditulis = DITULIS_LIMIT_NP[aksi]
        if (ditulis.length > 0) {
          setLayers((kini) =>
            kini.map((l, i) => {
              if (!SEMUA_LAYER.includes(aksi) && i !== indeks) return l
              const hasil = h.layers[i]
              if (hasil === undefined) return l
              const baru: SimpulLimit = { ...l }
              for (const k of ditulis) {
                const v = hasil[k]
                if (v !== undefined) baru[k] = v
              }
              return baru
            }),
          )
        }
        if (h.Total !== undefined) {
          setAkar({
            LimitSummaryList: h.LimitSummaryList ?? [],
            Total: h.Total,
            TotalLimitsROL: h.TotalLimitsROL ?? '',
          })
        }
      })
      .catch((e: unknown) => {
        setGagal(e instanceof Error ? e.message : String(e))
      })
  }

  return (
    <Panel judul={LIMITS_NP.judul}>
      <div className="tl-rincian">
        <KepalaBagian
          judul={LIMITS_NP.layers}
          jumlah={layers.length}
          aksi={
            bisaUbah && (
              <TombolTambah
                label={LIMITS_NP.tambahLayer}
                onClick={() => {
                  // `TreatyInNonAddItem(limits)`: ID = jumlah layer.
                  // `CopyLastLimitNP` tidak tercapai (Activity keluar
                  // sesudah langkah 5).
                  setLayers([
                    ...layers,
                    {
                      ID: String(layers.length + 1),
                      TreatyGroupList: [],
                      EgnpiTotalList: [],
                      PremiumEarnedList: [],
                      MDPList: [],
                      MDPMinList: [],
                      Reinstatement_List: [],
                    },
                  ])
                }}
              />
            )
          }
        />
        {layers.length === 0 ? (
          <Kosong pesan={petunjukKosong ?? LIMITS_NP.tanpaBaris} />
        ) : (
          <div className="tl-daftar">
            {layers.map((l, i) => (
              <KartuLipat
                key={i}
                nomor={i + 1}
                judul={judulLayer(l)}
                judulKosong={LIMITS_NP.layerBaru}
                bukaAwal={judulLayer(l) === ''}
                meta={ringkasLayer(l)}
                aksi={
                  bisaUbah && (
                    <TombolHapus
                      label={LIMITS_NP.hapus}
                      labelAkses={`${LIMITS_NP.hapus} ${LIMITS_NP.layers} ${i + 1}`}
                      onClick={() => {
                        setLayers(layers.filter((_, j) => j !== i))
                      }}
                    />
                  )
                }
              >
                <RincianLayer
                  l={l}
                  opsi={opsi}
                  bisaUbah={bisaUbah}
                  modeUbah={modeUbah}
                  onUbah={(baru) => {
                    setLayers(ganti(layers, i, baru))
                  }}
                  hitung={(aksi, baris = 0) => {
                    hitung(aksi, i, baris)
                  }}
                  onGantiGrup={(baru) => {
                    // Treaty Group berubah/dihapus → `TotalEgnpi` SEMUA layer;
                    // ID layer = indeksnya (`TreatyTypeSetIndex`).
                    const semua = ganti(layers, i, { ...baru, ID: String(i + 1) })
                    setLayers(semua)
                    hitung('egnpi', i, 0, semua)
                  }}
                />
              </KartuLipat>
            ))}
          </div>
        )}
        {pesan.length > 0 && (
          <ul className="tl-pesan" role="alert">
            {pesan.map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        )}
        {gagal !== '' && (
          <p className="tl-pesan" role="alert">
            {gagal}
          </p>
        )}

        <Bagian judul={LIMITS_NP.ringkasan}>
          <GridBaca kolom={KOLOM_RINGKASAN} baris={akar.LimitSummaryList} />
        </Bagian>

        <Bagian
          judul={LIMITS_NP.totalSemua}
          aksi={
            modeUbah && (
              <div className="trin__aksi">
                <button
                  type="button"
                  className="btn btn--primary btn--sm"
                  disabled={!bisaUbah}
                  onClick={() => {
                    hitung('total')
                  }}
                >
                  {LIMITS_NP.perbaruiTotal}
                </button>
                {kontrakRevisi(edmState) && (
                  <button
                    type="button"
                    className="btn btn--sm"
                    disabled={!bisaUbah}
                    onClick={() => {
                      hitung('nilai-list')
                    }}
                  >
                    {LIMITS_NP.perbaruiList}
                  </button>
                )}
              </div>
            )
          }
        >
          <div className="tl-tiga">
            {GRID_TOTAL.map((g) => (
              <GridNilai
                key={g.kunci}
                judul={g.judul}
                baris={akar.Total[g.kunci] ?? []}
                bisaUbahMataUang={false}
                bisaUbahNilai={false}
                tambah={false}
                mataUang={[]}
                onUbah={() => undefined}
              />
            ))}
          </div>
          <div className="trin__limit-medan">
            <Field label={LIMITS_NP.totalROL} value={tampil(akar.TotalLimitsROL, 2)} readOnly onChange={() => undefined} />
          </div>
        </Bagian>
      </div>
    </Panel>
  )
}
