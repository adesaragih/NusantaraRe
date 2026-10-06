// Layar Adjustment — `Section/InputTreatyInAdjustment.xml`.
//
// ⛔ DUA mode, keduanya dari ekspor:
//
//   daftar  `OutputParam.DATASHOW != 1` @481804 — tombol Add Revision /
//           Add Adjustment Premium, penyaring, grid dua belas kolom.
//   detail  `OutputParam.DATASHOW = 1` — kepala (@40422), lalu Old Data dan
//           New Data BERDAMPINGAN (@181268), Attachment (@276129), deret
//           tombol (@336222), History (@380295).
//
// ⛔ Datanya dari `TREATY_IN_EDM` + `M_TREATY_IN_EDM` — tabel WARISAN, baca
// saja. Bukan `VERSI_KONTRAK`: diukur nol baris.
//
// ⚠️ Panel Warning (`TreatyWarning.CARI1 != ''` @317776) TIDAK dirender:
// isinya diisi Activity saat jalan, dan halaman `TreatyWarning` tidak
// tersimpan di dokumen — syaratnya tidak pernah terpenuhi di sini.

import { useEffect, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat, Panel } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilDaftarPenyesuaian,
  ambilPenyesuaian,
  type BarisPenyesuaian,
  type BarisRiwayatWarisan,
  type Penyesuaian,
  type SisiPenyesuaian,
} from '../api'
import SisiForm, { selNilai, type ModeLayar } from '../komponen/SisiPenyesuaian'
import {
  JENIS_DAFTAR,
  KOLOM_DAFTAR,
  LEBAR_DAFTAR,
  LEBAR_TOMBOL_BARIS,
  MEDAN_KANAN_BARU,
  MEDAN_KANAN_LAMA,
  MEDAN_KIRI_BARU,
  MEDAN_KIRI_LAMA,
  PENYESUAIAN,
  TAB_BARU_NP,
  TAB_BARU_P,
  TAB_LAMA_NP,
  TAB_LAMA_P,
  TOMBOL,
} from '../labelsPenyesuaian'
import { PanelLampiranKontrak, PanelRiwayat } from './LampiranKontrak'

/** Baris per halaman grid daftar — `pyGridPaginator` @651683, ukuran bawaan Pega 10. */
const UKURAN_HALAMAN = 10

/** Nilai grid daftar, urut sesuai `KOLOM_DAFTAR`. */
export function selDaftar(b: BarisPenyesuaian): string[] {
  return [
    b.id, b.idAsal, b.jenisPenyesuaian, b.jenisMaterial, b.namaKontrak, b.sifatProporsi,
    b.asalBisnis, b.cedant, b.tanggalMulai, b.tanggalBerakhir, b.posisi, b.statusAkseptasi,
  ].map((v, i) => selNilai(JENIS_DAFTAR[i] ?? 'teks', v))
}

/** Lebar ekspor → persen kolom; satu tempat untuk seluruh grid layar ini. */
function persenLebar(lebar: readonly number[], i: number): string {
  const total = lebar.reduce((a, b) => a + b, 0)
  return `${(((lebar[i] ?? 0) / total) * 100).toFixed(2)}%`
}

/**
 * Baris `CommentList` → baris panel History.
 *
 * ⛔ Diukur: `T_VIEW_COMMENT` NOL baris berpengenal penyesuaian, sedang
 * `CommentList` berisi di 280 dari 280 dokumen. Itulah yang harness ikat
 * (`TreatyIn.CommentList` @407001).
 */
export function riwayatDari(sisi: SisiPenyesuaian): BarisRiwayatWarisan[] {
  return (sisi.larik.CommentList ?? []).map((b) => ({
    tanggal: b.Date ?? '',
    operator: b.OperatorName ?? '',
    disetujui: b.IsApproved ?? '',
    catatan: b.Suggest ?? '',
  }))
}

/**
 * Syarat TAMPIL tombol Save — @119719:
 * `TreatyIn.ViewState !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'`.
 *
 * ⚠️ Kunci `StatusAkseptasi` yang TIDAK ADA bukan `Resolve Complete` — `!=`
 * Pega bernilai benar, jadi tombolnya tampil.
 */
export function simpanTampil(mode: ModeLayar, medan: Readonly<Record<string, string>>): boolean {
  return mode !== '1' && medan.StatusAkseptasi !== 'Resolve Complete'
}

/** Kepala mode detail — kelima sel ber-`pyReadOnly=true`. */
function Kepala({ p }: { p: Penyesuaian }) {
  const m = p.baru.medan
  const kode = (label: string, kunci: string) =>
    Object.prototype.hasOwnProperty.call(m, kunci) ? (
      <div className="field">
        <label className="field__label">{label}</label>
        <input className="field__input field__input--readonly" type="text" value={m[kunci] ?? ''} readOnly />
        <span className="tria__redup">{PENYESUAIAN.kodeBelumBerteks}</span>
      </div>
    ) : (
      <div className="field">
        <label className="field__label">{label}</label>
        <span className="tria__tak-ada">{PENYESUAIAN.takAdaDiWarisan}</span>
      </div>
    )
  return (
    <Panel judul={PENYESUAIAN.judul}>
      <div className="tria__kepala">
        <div className="field">
          <label className="field__label">{PENYESUAIAN.idAsal}</label>
          <input className="field__input field__input--readonly" type="text" value={m.OLDID ?? p.idAsal} readOnly />
        </div>
        <div className="field">
          <label className="field__label">{PENYESUAIAN.id}</label>
          <input className="field__input field__input--readonly" type="text" value={m.ID ?? p.id} readOnly />
        </div>
        <fieldset className="tria__radio">
          <legend>{PENYESUAIAN.jenisReasuransi}</legend>
          {[PENYESUAIAN.proporsional, PENYESUAIAN.nonProporsional].map((v) => (
            <label key={v}>
              <input type="radio" name="tria-jenis-reasuransi" value={v} checked={m.ProportionType === v} disabled readOnly />
              {v}
            </label>
          ))}
        </fieldset>
        {/* `EDMState = 3` → LABEL "Adjustment Premium" @82379 menggantikan
            dropdown @76962 (`EDMState != 3`). Keduanya tidak pernah tampil
            bersamaan. */}
        {m.EDMState === '3' ? (
          <h4 className="tria__label-jenis">{PENYESUAIAN.premiPenyesuaian}</h4>
        ) : (
          kode(PENYESUAIAN.jenisPenyesuaian, 'EDMState')
        )}
        {kode(PENYESUAIAN.jenisMaterial, 'EDMMaterialType')}
      </div>
    </Panel>
  )
}

/**
 * Deret tombol — blok `TreatyMasterInEDM` @103200 (EDMState 1/2/3; terukur
 * 280 dari 280 penyesuaian memenuhinya).
 *
 * ⛔ Tombol yang belum punya jalur tulis DIMATIKAN, bukan dihilangkan.
 * `Close` HIDUP — ia menutup layar, tidak menulis apa pun. Blok dev
 * (`OperatorID.pyOrgDivision = 'IT'` @167976) tidak dibangun.
 */
function DeretTombol({ mode, medan, onTutup }: { mode: ModeLayar; medan: Readonly<Record<string, string>>; onTutup: () => void }) {
  return (
    <div className="tria__aksi" role="group" aria-label={PENYESUAIAN.judul}>
      {simpanTampil(mode, medan) && (
        <button type="button" className="btn btn--primary" disabled title={PENYESUAIAN.tombolTulisMati}>
          {TOMBOL.simpan}
        </button>
      )}
      <button type="button" className="btn" onClick={onTutup}>
        {TOMBOL.tutup}
      </button>
      <button type="button" className="btn" disabled title={TOMBOL.syaratPeranBelum}>
        {TOMBOL.aksi}
      </button>
      <span className="tria__redup">{PENYESUAIAN.tombolTulisMati}</span>
    </div>
  )
}

/** Mode detail — satu penyesuaian. */
function Detail({ id, mode, onTutup }: { id: string; mode: ModeLayar; onTutup: () => void }) {
  const [p, setP] = useState<Penyesuaian | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let dibuang = false
    setP(null)
    setGalat(null)
    ambilPenyesuaian(id)
      .then((x) => {
        if (!dibuang) setP(x)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [id])

  if (galat !== null) return <Gagal galat={galat} />
  if (p === null) return <Memuat />

  // ⛔ Cabang dibaca SEKALI dari halaman akar, dan KEDUA panel memakainya:
  // `TreatyInNONProportionalOldData.xml` memilih tab lewat
  // `TreatyIn.ProportionType` @278100 @292774 — akar, bukan `OLDDATA`.
  const cabang = p.baru.medan.ProportionType ?? ''
  const np = cabang === 'NonProportional'

  return (
    <>
      <Kepala p={p} />
      {/* ⭐ BERDAMPINGAN: Old di kiri, New di kanan — @181268. Wadahnya
          bersyarat `(NonProportional && DATASHOW=1) || (Proportional &&
          DATASHOW=1)`: kontrak tanpa cabang yang dikenal tidak membukanya. */}
      {cabang === 'NonProportional' || cabang === 'Proportional' ? (
        <div className="tria__bandingan">
          <SisiForm
            key={`lama-${p.id}`}
            judul={PENYESUAIAN.panelLama}
            sisi={p.lama}
            akar={p.baru}
            bacaSaja
            mode={mode}
            cabang={cabang}
            medanKiri={MEDAN_KIRI_LAMA}
            medanKanan={MEDAN_KANAN_LAMA}
            bagian={np ? 'TreatyInTabsNonProportionalOldData' : 'TreatyInTabsProportionalOldData'}
            tab={np ? TAB_LAMA_NP : TAB_LAMA_P}
          />
          <SisiForm
            key={`baru-${p.id}-${mode}`}
            judul={PENYESUAIAN.panelBaru}
            sisi={p.baru}
            akar={p.baru}
            bacaSaja={false}
            mode={mode}
            cabang={cabang}
            medanKiri={MEDAN_KIRI_BARU}
            medanKanan={MEDAN_KANAN_BARU}
            bagian={np ? 'TreatyInTabsNonProportional' : 'TreatyInTabsProportional'}
            tab={np ? TAB_BARU_NP : TAB_BARU_P}
          />
        </div>
      ) : null}

      <PanelLampiranKontrak masterID={p.id} />
      <DeretTombol mode={mode} medan={p.baru.medan} onTutup={onTutup} />
      <PanelRiwayat riwayat={riwayatDari(p.baru)} petunjuk={PENYESUAIAN.petunjukHistory} />
    </>
  )
}

/** Mode daftar — grid `TREATY_IN_EDM`. */
function Daftar({ onBuka }: { onBuka: (id: string, mode: ModeLayar) => void }) {
  const [baris, setBaris] = useState<BarisPenyesuaian[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [halaman, setHalaman] = useState(1)

  useEffect(() => {
    let dibuang = false
    ambilDaftarPenyesuaian()
      .then((d) => {
        if (!dibuang) setBaris(d)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [])

  const semua = baris ?? []
  const tampil = semua.slice((halaman - 1) * UKURAN_HALAMAN, halaman * UKURAN_HALAMAN)
  const lebar = [...LEBAR_DAFTAR, LEBAR_TOMBOL_BARIS, LEBAR_TOMBOL_BARIS]

  return (
    <Panel judul={PENYESUAIAN.judul}>
      {/* ⛔ Ketiga tombol kepala MATI: dua membuka pemilih yang MEMBUAT
          penyesuaian baru (jalur tulis), satu penyaring yang belum dibangun. */}
      <div className="tria__aksi" role="group" aria-label={PENYESUAIAN.judul}>
        <button type="button" className="btn btn--primary" disabled title={PENYESUAIAN.tombolTulisMati}>
          {PENYESUAIAN.tambahRevisi}
        </button>
        <button type="button" className="btn btn--primary" disabled title={PENYESUAIAN.tombolTulisMati}>
          {PENYESUAIAN.tambahPremi}
        </button>
        <button type="button" className="btn btn--ghost" disabled title={PENYESUAIAN.belumDibangun}>
          {PENYESUAIAN.saring}
        </button>
      </div>

      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && (
        <>
          <div className="table-wrap">
            <table className="tria__tabel">
              <colgroup>
                {lebar.map((_, i) => (
                  <col key={i} style={{ width: persenLebar(lebar, i) }} />
                ))}
              </colgroup>
              <thead>
                <tr>
                  {KOLOM_DAFTAR.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                  <th scope="col" aria-label={PENYESUAIAN.ubah} />
                  <th scope="col" aria-label={PENYESUAIAN.lihat} />
                </tr>
              </thead>
              <tbody>
                {semua.length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_DAFTAR.length + 2}>
                      <Kosong pesan={PENYESUAIAN.tanpaBaris} petunjuk={PENYESUAIAN.petunjukDaftarKosong} />
                    </td>
                  </tr>
                )}
                {tampil.map((b) => (
                  <tr key={b.id}>
                    {selDaftar(b).map((v, i) => (
                      <td key={i}>{v}</td>
                    ))}
                    {/* `Edit` @782051 dan `View` @798870 sama-sama MEMBUKA —
                        bedanya `ViewState` 0 lawan 1. Membuka bukan menulis;
                        yang menulis tombol Save, dan ia mati. */}
                    <td>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          onBuka(b.id, '0')
                        }}
                      >
                        {PENYESUAIAN.ubah}
                      </button>
                    </td>
                    <td>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          onBuka(b.id, '1')
                        }}
                      >
                        {PENYESUAIAN.lihat}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN} total={semua.length} onPindah={setHalaman} />
        </>
      )}
    </Panel>
  )
}

export default function PenyesuaianKontrak() {
  const [buka, setBuka] = useState<{ id: string; mode: ModeLayar } | null>(null)
  return (
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{PENYESUAIAN.judulMenu}</h2>
        {buka !== null && (
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              setBuka(null)
            }}
          >
            {PENYESUAIAN.kembali}
          </button>
        )}
      </header>
      {buka === null ? (
        <Daftar
          onBuka={(id, mode) => {
            setBuka({ id, mode })
          }}
        />
      ) : (
        <Detail
          id={buka.id}
          mode={buka.mode}
          onTutup={() => {
            setBuka(null)
          }}
        />
      )}
    </div>
  )
}
