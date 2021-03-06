// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package localtype

import (
	"fmt"
	"io"
	"sort"

	types "github.com/filecoin-project/lotus/chain/types"
	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = sort.Sort

func (t *BeaconEntryInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{162}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Round (uint64) (uint64)
	if len("Round") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Round\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Round"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Round")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Round)); err != nil {
		return err
	}

	// t.Entry (types.BeaconEntry) (struct)
	if len("Entry") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Entry\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Entry"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Entry")); err != nil {
		return err
	}

	if err := t.Entry.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *BeaconEntryInfo) UnmarshalCBOR(r io.Reader) error {
	*t = BeaconEntryInfo{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("BeaconEntryInfo: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Round (uint64) (uint64)
		case "Round":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.Round = uint64(extra)

			}
			// t.Entry (types.BeaconEntry) (struct)
		case "Entry":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}
					t.Entry = new(types.BeaconEntry)
					if err := t.Entry.UnmarshalCBOR(br); err != nil {
						return xerrors.Errorf("unmarshaling t.Entry pointer: %w", err)
					}
				}

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
